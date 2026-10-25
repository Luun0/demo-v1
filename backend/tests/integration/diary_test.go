package integration

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

const secret = "SPECIALIST-ONLY-SECRET-7f3a"

func TestSpecialistsOnlyNeverReachesParent(t *testing.T) {
	e := newEnv(t, nil)
	w := e.world()
	c := "/api/v1/children/" + w.child1.String()
	e.entry(w.parent1, w.child1, "shared", "shared note alpha")
	so := e.entry(w.spec1, w.child1, "specialists_only", secret)
	e.entry(w.spec1, w.child1, "shared", "shared note beta")

	for _, p := range []*client{w.parent1, w.coParent} {
		checks := []resp{
			must(t, p.do("GET", c+"/diary-entries", nil), 200),
			must(t, p.do("GET", c+"/diary-entries?q=SECRET", nil), 200),
			must(t, p.do("GET", c+"/diary-entries?visibility=specialists_only", nil), 200),
			must(t, p.do("GET", c+"/diary-entries/"+so.String()+"/versions", nil), 404),
			must(t, p.do("PUT", c+"/diary-entries/"+so.String(), map[string]any{"type": "observation", "visibility": "shared", "content": "x", "observed_at": "2026-10-20T09:00:00Z"}), 404),
			p.do("POST", c+"/exports", nil),
			must(t, p.do("GET", c+"/audit-events", nil), 200),
		}
		for i, r := range checks {
			if strings.Contains(r.body, secret) || strings.Contains(r.body, "specialists_only") {
				t.Errorf("leak in check %d: %s", i, r.body)
			}
		}
		if got := checks[0].hdr.Get("X-Total-Count"); got != "2" {
			t.Errorf("count leaks hidden entries: X-Total-Count=%s", got)
		}
		if strings.Count(checks[0].body, `"id"`) != 2 || checks[1].body != "[]\n" || checks[2].body != "[]\n" {
			t.Errorf("unexpected parent lists: %s | %s | %s", checks[0].body, checks[1].body, checks[2].body)
		}
	}
	// the specialist sees all three, including search
	r := must(t, w.spec1.do("GET", c+"/diary-entries?q=secret", nil), 200)
	if !strings.Contains(r.body, secret) || must(t, w.spec1.do("GET", c+"/diary-entries", nil), 200).hdr.Get("X-Total-Count") != "3" {
		t.Fatalf("specialist view wrong: %s", r.body)
	}
	// parent export contains only shared entries
	ex := must(t, w.parent1.do("POST", c+"/exports", nil), 200)
	if strings.Count(ex.body, `"visibility":"shared"`) != 2 {
		t.Fatalf("export: %s", ex.body)
	}
}

func TestParentCanOnlyWriteShared(t *testing.T) {
	e := newEnv(t, nil)
	w := e.world()
	c := "/api/v1/children/" + w.child1.String() + "/diary-entries"
	body := map[string]any{"type": "observation", "visibility": "specialists_only", "content": "x", "observed_at": "2026-10-20T09:00:00Z"}
	must(t, w.parent1.do("POST", c, body), 403)
	id := e.entry(w.parent1, w.child1, "shared", "mine")
	must(t, w.parent1.do("PUT", c+"/"+id.String(), body), 403)
	if auditCount(t, e, "diary:edit", "denied") != 1 || auditCount(t, e, "diary:create", "denied") != 1 {
		t.Fatal("denied writes must be audited exactly once")
	}
	// co-parent cannot edit someone else's entry
	body["visibility"] = "shared"
	must(t, w.coParent.do("PUT", c+"/"+id.String(), body), 403)
	// invalid enum / missing fields -> 422
	must(t, w.parent1.do("POST", c, map[string]any{"type": "rumor", "visibility": "public", "content": "", "observed_at": "yesterday"}), 422)
}

func TestAdminNeverSeesDiaryContent(t *testing.T) {
	e := newEnv(t, nil)
	w := e.world()
	c := "/api/v1/children/" + w.child1.String()
	id := e.entry(w.spec1, w.child1, "shared", "admin-must-not-see-this")
	for _, r := range []resp{
		w.admin.do("GET", c+"/diary-entries", nil),
		w.admin.do("GET", c+"/diary-entries?q=admin", nil),
		w.admin.do("GET", c+"/diary-entries/"+id.String()+"/versions", nil),
		w.admin.do("POST", c+"/exports", nil),
		w.admin.do("GET", c+"/audit-events", nil),
		w.admin.do("GET", "/api/v1/audit-events", nil),
		w.admin.do("GET", c, nil),
		w.admin.do("GET", "/api/v1/children", nil),
		w.admin.do("GET", "/api/v1/users", nil),
	} {
		if strings.Contains(r.body, "admin-must-not-see-this") || strings.Contains(r.body, `"content"`) {
			t.Errorf("admin saw diary content: %d %s", r.code, r.body)
		}
	}
}

func TestEveryEditCreatesVersionAndExactlyOneAuditEvent(t *testing.T) {
	e := newEnv(t, nil)
	w := e.world()
	c := "/api/v1/children/" + w.child1.String() + "/diary-entries"
	id := e.entry(w.spec1, w.child1, "shared", "v1")
	before := auditCount(t, e, "diary:edit", "success")
	for i, txt := range []string{"v2", "v3"} {
		must(t, w.spec1.do("PUT", c+"/"+id.String(), map[string]any{"type": "progress", "visibility": "specialists_only", "content": txt, "observed_at": "2026-10-21T09:00:00Z"}), 200)
		if got := auditCount(t, e, "diary:edit", "success"); got != before+i+1 {
			t.Fatalf("want exactly one audit event per edit, got %d", got-before)
		}
	}
	var vs []struct {
		VersionNumber int    `json:"version_number"`
		Content       string `json:"content"`
	}
	must(t, w.spec1.do("GET", c+"/"+id.String()+"/versions", nil), 200).json(t, &vs)
	if len(vs) != 3 || vs[0].Content != "v1" || vs[2].VersionNumber != 3 {
		t.Fatalf("versions: %+v", vs)
	}
	// entry is now specialists_only: parent gets 404 on history
	must(t, w.parent1.do("GET", c+"/"+id.String()+"/versions", nil), 404)
	// another specialist (spec2 after assignment) cannot edit spec1's entry
	must(t, w.parent1.do("POST", "/api/v1/children/"+w.child1.String()+"/specialist-assignments", map[string]any{"specialist_id": w.spec2.id}), 201)
	must(t, w.spec2.do("PUT", c+"/"+id.String(), map[string]any{"type": "progress", "visibility": "shared", "content": "x", "observed_at": "2026-10-21T09:00:00Z"}), 403)
	// reads are audited too
	if auditCount(t, e, "diary:versions", "success") < 1 {
		t.Fatal("reads must be audited")
	}
}

func TestDiaryPaginationAndXSSStoredAsText(t *testing.T) {
	e := newEnv(t, nil)
	w := e.world()
	c := "/api/v1/children/" + w.child1.String() + "/diary-entries"
	xss := `<script>alert(1)</script><img src=x onerror=alert(2)>`
	for i := 0; i < 5; i++ {
		e.entry(w.parent1, w.child1, "shared", xss)
	}
	r := must(t, w.parent1.do("GET", c+"?limit=2", nil), 200)
	if r.hdr.Get("Content-Type") != "application/json" || !strings.Contains(r.body, `\u003cscript\u003e`) {
		t.Fatalf("payload must be JSON-escaped text: %s", r.body)
	}
	seen := strings.Count(r.body, `"id"`)
	cur := r.hdr.Get("X-Next-Cursor")
	for cur != "" {
		r = must(t, w.parent1.do("GET", c+"?limit=2&cursor="+cur, nil), 200)
		seen += strings.Count(r.body, `"id"`)
		cur = r.hdr.Get("X-Next-Cursor")
	}
	if seen != 5 {
		t.Fatalf("paginated %d of 5", seen)
	}
	must(t, w.parent1.do("GET", c+"?limit=0", nil), 422)
	must(t, w.parent1.do("GET", c+"?cursor=@@", nil), 422)
	must(t, w.parent1.do("GET", c+"?q=%25_%5C", nil), 200) // LIKE wildcards are escaped
}

func TestSessionsByRole(t *testing.T) {
	e := newEnv(t, nil)
	w := e.world()
	c := "/api/v1/children/" + w.child1.String() + "/sessions"
	body := map[string]any{"specialist_id": w.spec1.id, "title": "Speech", "starts_at": "2030-01-01T10:00:00Z", "ends_at": "2030-01-01T11:00:00Z", "mode": "in_person", "status": "scheduled"}
	var s struct {
		ID uuid.UUID `json:"id"`
	}
	must(t, w.spec1.do("POST", c, body), 201).json(t, &s)
	must(t, w.parent1.do("PATCH", c+"/"+s.ID.String(), map[string]any{"status": "cancelled"}), 200)
	must(t, w.spec1.do("PATCH", c+"/"+s.ID.String(), map[string]any{"status": "completed"}), 200)
	must(t, w.spec1.do("PATCH", c+"/"+s.ID.String(), map[string]any{"status": "canceled"}), 422)
	must(t, w.spec1.do("PATCH", c+"/"+s.ID.String(), map[string]any{"specialist_id": w.spec2.id}), 422)
	must(t, w.spec1.do("PATCH", c+"/"+uuid.NewString(), map[string]any{"title": "x"}), 404)
	// spec2 assigned but not the session's specialist -> 403
	must(t, w.parent1.do("POST", "/api/v1/children/"+w.child1.String()+"/specialist-assignments", map[string]any{"specialist_id": w.spec2.id}), 201)
	must(t, w.spec2.do("PATCH", c+"/"+s.ID.String(), map[string]any{"title": "x"}), 403)
	// session with a non-assigned specialist or bad times -> 422
	body["specialist_id"] = uuid.New()
	must(t, w.parent1.do("POST", c, body), 422)
	body["specialist_id"], body["ends_at"] = w.spec1.id, "2029-01-01T10:00:00Z"
	must(t, w.parent1.do("POST", c, body), 422)
	var list []any
	must(t, w.coParent.do("GET", c, nil), 200).json(t, &list)
	if len(list) != 1 {
		t.Fatalf("sessions: %v", list)
	}
}

func TestDeleteChildAndAdminFlows(t *testing.T) {
	e := newEnv(t, nil)
	w := e.world()
	var page []any
	must(t, w.admin.do("GET", "/api/v1/users?limit=2", nil), 200).json(t, &page)
	if len(page) != 2 {
		t.Fatal("pagination of users")
	}
	must(t, w.admin.do("GET", "/api/v1/audit-events?limit=5", nil), 200)
	must(t, w.admin.do("PATCH", "/api/v1/users/"+w.parent1.id.String(), map[string]any{"specialist_approved": true}), 422)
	must(t, w.admin.do("PATCH", "/api/v1/users/"+w.parent1.id.String(), map[string]any{}), 422)
	must(t, w.admin.do("PATCH", "/api/v1/users/"+uuid.NewString(), map[string]any{"status": "inactive"}), 404)
	// admin creates another admin (SEC-14)
	must(t, w.admin.do("POST", "/api/v1/users", map[string]any{"email": "a2@example.test", "password": pw, "first_name": "A", "last_name": "B", "role": "admin"}), 201)
	// primary guardian deletes child 2 data; then everything is 404
	must(t, w.parent2.do("DELETE", "/api/v1/children/"+w.child2.String(), nil), 204)
	must(t, w.parent2.do("GET", "/api/v1/children/"+w.child2.String(), nil), 404)
	// admin can delete on documented request
	must(t, w.admin.do("DELETE", "/api/v1/children/"+w.child1.String(), nil), 204)
	if auditCount(t, e, "child:delete", "success") != 2 {
		t.Fatal("deletions must be audited")
	}
}
