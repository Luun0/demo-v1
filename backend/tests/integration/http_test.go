package integration

import (
	"net/http"
	"regexp"
	"strings"
	"testing"
)

func TestSecurityHeadersCORSAndBodyLimit(t *testing.T) {
	e := newEnv(t, map[string]string{"MAX_BODY_BYTES": "1024"})
	w := e.world()
	r := w.parent1.do("GET", "/api/v1/children", nil, "Origin", testOrigin)
	for k, want := range map[string]string{
		"Content-Security-Policy": "frame-ancestors 'none'", "X-Content-Type-Options": "nosniff", "Referrer-Policy": "no-referrer",
		"Strict-Transport-Security": "max-age=", "X-Frame-Options": "DENY", "Cache-Control": "no-store",
		"Access-Control-Allow-Origin": testOrigin, "X-Request-Id": "",
	} {
		if !strings.Contains(r.hdr.Get(k), want) || r.hdr.Get(k) == "" {
			t.Errorf("header %s = %q", k, r.hdr.Get(k))
		}
	}
	evil := w.parent1.do("GET", "/api/v1/children", nil, "Origin", "https://evil.example")
	if evil.hdr.Get("Access-Control-Allow-Origin") != "" {
		t.Fatal("CORS must not allow foreign origin")
	}
	must(t, w.parent1.do("OPTIONS", "/api/v1/children", nil, "Origin", "https://evil.example", "Access-Control-Request-Method", "GET"), 403)
	must(t, w.parent1.do("OPTIONS", "/api/v1/children", nil, "Origin", testOrigin, "Access-Control-Request-Method", "POST"), 204)

	big := `{"type":"observation","visibility":"shared","content":"` + strings.Repeat("a", 4096) + `","observed_at":"2026-10-20T09:00:00Z"}`
	must(t, w.parent1.do("POST", "/api/v1/children/"+w.child1.String()+"/diary-entries", big), http.StatusRequestEntityTooLarge)
	must(t, w.parent1.do("POST", "/api/v1/children", "{not json"), 422)
	must(t, w.parent1.do("GET", "/api/v1/nope", nil), 404)
	must(t, e.client().do("GET", "/healthz", nil), 200)
	must(t, e.client().do("GET", "/api/v1/readyz", nil), 200)
}

func TestErrorsAreProblemJSONWithoutInternals(t *testing.T) {
	e := newEnv(t, nil)
	w := e.world()
	r := must(t, w.parent2.do("GET", "/api/v1/children/"+w.child1.String(), nil), 404)
	if r.hdr.Get("Content-Type") != "application/problem+json" || !strings.Contains(r.body, `"status":404`) {
		t.Fatalf("bad problem: %v %s", r.hdr, r.body)
	}
	r = must(t, w.parent1.do("POST", "/api/v1/auth/register", map[string]any{"email": "bad", "password": "tiny1", "first_name": "", "last_name": "x", "role": "parent"}), 422)
	if strings.Contains(r.body, "tiny1") || !strings.Contains(r.body, `"field":"password"`) {
		t.Fatalf("422 must name fields, never echo values: %s", r.body)
	}
}

// TestLogsContainNoPII checks the JSON log stream for emails, tokens,
// passwords and diary text (NFR / PRV logging requirement).
func TestLogsContainNoPII(t *testing.T) {
	e := newEnv(t, nil)
	w := e.world()
	e.entry(w.parent1, w.child1, "shared", "secret-diary-text-123")
	w.parent1.login("p1@example.test", "Wrong-password-1")
	logs := e.logs.String()
	if logs == "" {
		t.Fatal("expected request logs")
	}
	for _, re := range []string{`[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[a-z]{2,}`, `eyJ[A-Za-z0-9_-]{10,}`, `secret-diary-text`, pw, `Wrong-password`} {
		if regexp.MustCompile(re).MatchString(logs) {
			t.Errorf("log stream matches forbidden pattern %q", re)
		}
	}
}
