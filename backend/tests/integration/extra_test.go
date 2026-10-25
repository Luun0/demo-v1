package integration

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/AknurOE/sen-support-platform/internal/application"
	"github.com/AknurOE/sen-support-platform/internal/domain"
)

func TestMiscUseCaseBranches(t *testing.T) {
	e := newEnv(t, map[string]string{"REFRESH_TOKEN_TTL": "1s"})
	w := e.world()
	svc := e.app.Service
	ctx := context.Background()
	// logout with unknown / empty token is idempotent
	if err := svc.Logout(ctx, "unknown"); err != nil {
		t.Fatal(err)
	}
	if err := svc.Logout(ctx, ""); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Refresh(ctx, ""); err == nil {
		t.Fatal("empty refresh must fail")
	}
	if _, err := svc.Me(ctx, nil); err == nil {
		t.Fatal("guest has no profile")
	}
	if _, err := svc.ListChildren(ctx, nil); err == nil {
		t.Fatal("guest lists nothing")
	}
	if _, err := svc.BootstrapAdmin(ctx, "bad", "x", "", ""); err == nil {
		t.Fatal("invalid bootstrap must fail")
	}
	// expired refresh token
	raw := ""
	for _, ck := range w.parent1.http.Jar.Cookies(mustURL(e.srv.URL + "/api/v1/auth/")) {
		if ck.Name == "ssp_refresh" {
			raw = ck.Value
		}
	}
	time.Sleep(1100 * time.Millisecond)
	if _, err := svc.Refresh(ctx, raw); err == nil {
		t.Fatal("expired refresh token accepted")
	}
	// user listing with cursor and child audit paging
	r := must(t, w.admin.do("GET", "/api/v1/users?limit=3", nil), 200)
	cur := r.hdr.Get("X-Next-Cursor")
	must(t, w.admin.do("GET", "/api/v1/users?limit=3&cursor="+cur, nil), 200)
	must(t, w.admin.do("GET", "/api/v1/users?cursor=zzz", nil), 422)
	for i := 0; i < 3; i++ {
		must(t, w.parent1.do("GET", "/api/v1/children/"+w.child1.String(), nil), 200)
	}
	r = must(t, w.parent1.do("GET", "/api/v1/children/"+w.child1.String()+"/audit-events?limit=2", nil), 200)
	if r.hdr.Get("X-Next-Cursor") == "" || strings.Contains(r.body, "@") {
		t.Fatalf("audit page: %v %s", r.hdr, r.body)
	}
	must(t, w.parent1.do("GET", "/api/v1/children/"+w.child1.String()+"/audit-events?cursor="+r.hdr.Get("X-Next-Cursor"), nil), 200)
	must(t, w.admin.do("GET", "/api/v1/audit-events?limit=1", nil), 200)
	// version history: a shared entry turned specialists_only and back hides the private version from parents
	id := e.entry(w.spec1, w.child1, "shared", "public v1")
	c := "/api/v1/children/" + w.child1.String() + "/diary-entries/" + id.String()
	must(t, w.spec1.do("PUT", c, map[string]any{"type": "observation", "visibility": "specialists_only", "content": "private v2", "observed_at": "2026-10-20T09:00:00Z"}), 200)
	must(t, w.spec1.do("PUT", c, map[string]any{"type": "observation", "visibility": "shared", "content": "public v3", "observed_at": "2026-10-20T09:00:00Z"}), 200)
	h := must(t, w.parent1.do("GET", c+"/versions", nil), 200)
	if strings.Contains(h.body, "private v2") || !strings.Contains(h.body, "public v3") {
		t.Fatalf("history leak: %s", h.body)
	}
	// user patch: re-activate and approve specialist
	st := domain.StatusActive
	if _, err := svc.UpdateUser(ctx, &application.Actor{ID: w.admin.id, Role: domain.RoleAdmin}, w.parent2.id, application.UserPatch{Status: &st}); err != nil {
		t.Fatal(err)
	}
	must(t, w.admin.do("PATCH", "/api/v1/users/"+w.spec2.id.String(), map[string]any{"specialist_approved": false}), 200)
	must(t, w.admin.do("PATCH", "/api/v1/users/bad", map[string]any{"status": "inactive"}), 404)
	must(t, w.admin.do("PATCH", "/api/v1/users/"+w.spec2.id.String(), map[string]any{"status": "banned"}), 422)
	// guardian link validation
	must(t, w.parent1.do("POST", "/api/v1/children/"+w.child1.String()+"/guardians", map[string]any{"guardian_id": w.spec1.id, "relationship": "guardian", "is_primary": false}), 422)
	must(t, w.parent1.do("POST", "/api/v1/children/"+w.child1.String()+"/guardians", map[string]any{"guardian_id": w.parent2.id, "relationship": "uncle", "is_primary": false}), 422)
	must(t, w.parent1.do("POST", "/api/v1/children", map[string]any{"first_name": "", "last_name": "x", "birth_date": "bad", "policy_version": ""}), 422)
	must(t, w.parent1.do("POST", "/api/v1/children", map[string]any{"first_name": "", "last_name": "x", "birth_date": "2019-01-01", "policy_version": ""}), 422)
	must(t, w.parent1.do("PATCH", "/api/v1/children/"+w.child1.String(), map[string]any{"last_name": "New", "birth_date": "bad"}), 422)
	must(t, w.parent1.do("PATCH", "/api/v1/children/"+w.child1.String()+"/specialist-assignments/"+w.spec2.id.String(), map[string]any{"ended_at": "2026-01-01T00:00:00Z"}), 404)
	must(t, w.parent1.do("PATCH", "/api/v1/children/"+w.child1.String()+"/specialist-assignments/"+w.spec1.id.String(), map[string]any{"ended_at": "bad"}), 422)
}

func TestTeamAndSpecialistDirectory(t *testing.T) {
	e := newEnv(t, nil)
	w := e.world()
	c := "/api/v1/children/" + w.child1.String() + "/team"
	r := must(t, w.parent1.do("GET", c, nil), 200)
	if !strings.Contains(r.body, `"is_primary":true`) || !strings.Contains(r.body, `"role":"specialist"`) || strings.Contains(r.body, "@") {
		t.Fatalf("team: %s", r.body)
	}
	must(t, w.spec1.do("GET", c, nil), 200)
	must(t, w.parent2.do("GET", c, nil), 404)
	d := must(t, w.parent1.do("GET", "/api/v1/specialists", nil), 200)
	if strings.Count(d.body, `"id"`) != 2 || strings.Contains(d.body, "@") {
		t.Fatalf("directory: %s", d.body)
	}
	must(t, w.spec1.do("GET", "/api/v1/specialists", nil), 403)
	must(t, w.admin.do("GET", "/api/v1/specialists", nil), 200)
}
