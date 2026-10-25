package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"

	"github.com/AknurOE/sen-support-platform/internal/app"
	"github.com/AknurOE/sen-support-platform/internal/application"
	"github.com/AknurOE/sen-support-platform/internal/infrastructure/config"
	"github.com/AknurOE/sen-support-platform/internal/testpg"
)

func TestMain(m *testing.M) { os.Exit(testpg.Main(m)) }

const testOrigin = "https://app.example.test"
const pw = "Integration-Pass-42"

type syncBuf struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (s *syncBuf) Write(p []byte) (int, error) { s.mu.Lock(); defer s.mu.Unlock(); return s.b.Write(p) }
func (s *syncBuf) String() string              { s.mu.Lock(); defer s.mu.Unlock(); return s.b.String() }

type env struct {
	t     *testing.T
	srv   *httptest.Server
	app   *app.App
	logs  *syncBuf
	dbURL string
}

func newEnv(t *testing.T, overrides map[string]string) *env {
	t.Helper()
	vars := map[string]string{
		"DATABASE_URL": testpg.NewDB(t), "JWT_SECRET": strings.Repeat("k", 48), "CORS_ALLOWED_ORIGINS": testOrigin,
		"COOKIE_SECURE": "false", "ARGON2_MEMORY_KIB": "8192", "ARGON2_TIME": "1", "ARGON2_THREADS": "1",
		"LOGIN_RATE_PER_MINUTE": "1000", "API_RATE_PER_MINUTE": "100000",
	}
	for k, v := range overrides {
		vars[k] = v
	}
	cfg, err := config.Load(func(k string) string { return vars[k] })
	if err != nil {
		t.Fatal(err)
	}
	logs := &syncBuf{}
	a, err := app.Build(context.Background(), cfg, app.Options{Logs: logs})
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(a.Handler)
	t.Cleanup(func() { srv.Close(); a.Store.Close() })
	return &env{t: t, srv: srv, app: a, logs: logs, dbURL: vars["DATABASE_URL"]}
}

// client is one browser-like user agent with its own cookie jar.
type client struct {
	e     *env
	http  *http.Client
	token string
	id    uuid.UUID
	email string
}

func (e *env) client() *client {
	jar, _ := cookiejar.New(nil)
	return &client{e: e, http: &http.Client{Jar: jar}}
}

type resp struct {
	code int
	body string
	hdr  http.Header
}

func (r resp) json(t *testing.T, v any) {
	t.Helper()
	if err := json.Unmarshal([]byte(r.body), v); err != nil {
		t.Fatalf("decode %q: %v", r.body, err)
	}
}

func (c *client) do(method, path string, body any, hdr ...string) resp {
	c.e.t.Helper()
	var rd io.Reader
	switch b := body.(type) {
	case nil:
	case string:
		rd = strings.NewReader(b)
	default:
		buf, _ := json.Marshal(b)
		rd = bytes.NewReader(buf)
	}
	req, _ := http.NewRequest(method, c.e.srv.URL+path, rd)
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	for i := 0; i+1 < len(hdr); i += 2 {
		req.Header.Set(hdr[i], hdr[i+1])
	}
	res, err := c.http.Do(req)
	if err != nil {
		c.e.t.Fatal(err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	return resp{res.StatusCode, string(b), res.Header}
}

func (c *client) csrf() string {
	for _, ck := range c.http.Jar.Cookies(mustURL(c.e.srv.URL + "/")) {
		if ck.Name == "ssp_csrf" {
			return ck.Value
		}
	}
	return ""
}

func (c *client) login(email, password string) resp {
	r := c.do("POST", "/api/v1/auth/login", map[string]string{"email": email, "password": password})
	if r.code == 200 {
		var tok struct {
			AccessToken string `json:"access_token"`
		}
		r.json(c.e.t, &tok)
		c.token, c.email = tok.AccessToken, email
		var me struct {
			ID uuid.UUID `json:"id"`
		}
		c.do("GET", "/api/v1/auth/me", nil).json(c.e.t, &me)
		c.id = me.ID
	}
	return r
}

func must(t *testing.T, r resp, code int) resp {
	t.Helper()
	if r.code != code {
		t.Fatalf("want %d got %d: %s", code, r.code, r.body)
	}
	return r
}

// world is a standard fixture: admin, two families, two specialists.
type world struct {
	admin, parent1, parent2, coParent, spec1, spec2 *client
	child1, child2                                  uuid.UUID
}

func (e *env) user(email, role string, approve bool, admin *client) *client {
	r := must(e.t, admin.do("POST", "/api/v1/users", map[string]any{"email": email, "password": pw, "first_name": "Syn", "last_name": "Thetic",
		"role": role, "specialist_approved": approve}), 201)
	_ = r
	c := e.client()
	must(e.t, c.login(email, pw), 200)
	return c
}

func (e *env) world() *world {
	t := e.t
	if _, err := e.app.Service.BootstrapAdmin(context.Background(), "admin@example.test", pw, "Ad", "Min"); err != nil {
		t.Fatal(err)
	}
	w := &world{admin: e.client()}
	must(t, w.admin.login("admin@example.test", pw), 200)
	w.parent1 = e.user("p1@example.test", "parent", false, w.admin)
	w.parent2 = e.user("p2@example.test", "parent", false, w.admin)
	w.coParent = e.user("p1b@example.test", "parent", false, w.admin)
	w.spec1 = e.user("s1@example.test", "specialist", true, w.admin)
	w.spec2 = e.user("s2@example.test", "specialist", true, w.admin)
	w.child1 = e.child(w.parent1, "Mira")
	w.child2 = e.child(w.parent2, "Timur")
	must(t, w.parent1.do("POST", "/api/v1/children/"+w.child1.String()+"/guardians", map[string]any{"guardian_id": w.coParent.id, "relationship": "guardian", "is_primary": false}), 201)
	must(t, w.parent1.do("POST", "/api/v1/children/"+w.child1.String()+"/specialist-assignments", map[string]any{"specialist_id": w.spec1.id}), 201)
	return w
}

func (e *env) child(p *client, name string) uuid.UUID {
	var c struct {
		ID uuid.UUID `json:"id"`
	}
	must(e.t, p.do("POST", "/api/v1/children", map[string]any{"first_name": name, "last_name": "Demo", "birth_date": "2019-03-01", "policy_version": "2026-10"}), 201).json(e.t, &c)
	return c.ID
}

func (e *env) entry(c *client, child uuid.UUID, vis, content string) uuid.UUID {
	var out struct {
		ID uuid.UUID `json:"id"`
	}
	must(e.t, c.do("POST", "/api/v1/children/"+child.String()+"/diary-entries", map[string]any{"type": "observation", "visibility": vis,
		"content": content, "observed_at": "2026-10-20T09:00:00Z"}), 201).json(e.t, &out)
	return out.ID
}

func auditCount(t *testing.T, e *env, action string, outcome string) int {
	t.Helper()
	n := 0
	evs, err := e.app.Store.Repo().ListAudit(context.Background(), nil, 100000, (*application.Cursor)(nil))
	if err != nil {
		t.Fatal(err)
	}
	for _, ev := range evs {
		if ev.Action == action && string(ev.Outcome) == outcome {
			n++
		}
	}
	return n
}
