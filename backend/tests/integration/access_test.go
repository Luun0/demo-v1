package integration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// childEndpoints lists every endpoint that takes a child or entry ID.
func childEndpoints(child, entry, session, spec uuid.UUID) [][3]any {
	c := "/api/v1/children/" + child.String()
	diaryBody := map[string]any{"type": "observation", "visibility": "shared", "content": "x", "observed_at": "2026-10-20T09:00:00Z"}
	return [][3]any{
		{"GET", c, nil},
		{"PATCH", c, map[string]any{"first_name": "Hacked"}},
		{"DELETE", c, nil},
		{"POST", c + "/guardians", map[string]any{"guardian_id": uuid.New(), "relationship": "guardian", "is_primary": false}},
		{"POST", c + "/specialist-assignments", map[string]any{"specialist_id": spec}},
		{"PATCH", c + "/specialist-assignments/" + spec.String(), map[string]any{"ended_at": "2026-10-20T09:00:00Z"}},
		{"GET", c + "/sessions", nil},
		{"POST", c + "/sessions", map[string]any{"specialist_id": spec, "title": "t", "starts_at": "2030-01-01T10:00:00Z", "ends_at": "2030-01-01T11:00:00Z", "mode": "online", "status": "scheduled"}},
		{"PATCH", c + "/sessions/" + session.String(), map[string]any{"title": "x"}},
		{"GET", c + "/diary-entries", nil},
		{"POST", c + "/diary-entries", diaryBody},
		{"PUT", c + "/diary-entries/" + entry.String(), diaryBody},
		{"DELETE", c + "/diary-entries/" + entry.String(), nil},
		{"GET", c + "/diary-entries/" + entry.String() + "/versions", nil},
		{"GET", c + "/audit-events", nil},
		{"POST", c + "/exports", nil},
	}
}

func TestCrossFamilyIDORReturns404OnEveryEndpoint(t *testing.T) {
	e := newEnv(t, nil)
	w := e.world()
	entry := e.entry(w.parent1, w.child1, "shared", "family one note")
	var s struct {
		ID uuid.UUID `json:"id"`
	}
	must(t, w.parent1.do("POST", "/api/v1/children/"+w.child1.String()+"/sessions", map[string]any{"specialist_id": w.spec1.id, "title": "t",
		"starts_at": "2030-01-01T10:00:00Z", "ends_at": "2030-01-01T11:00:00Z", "mode": "online", "status": "scheduled"}), 201).json(t, &s)

	for _, who := range []struct {
		name string
		c    *client
	}{{"other family parent", w.parent2}, {"unassigned specialist", w.spec2}} {
		for _, ep := range childEndpoints(w.child1, entry, s.ID, w.spec1.id) {
			r := who.c.do(ep[0].(string), ep[1].(string), ep[2])
			if r.code != 404 {
				t.Errorf("%s %s %s: want 404 got %d %s", who.name, ep[0], ep[1], r.code, r.body)
			}
			if strings.Contains(r.body, "family one note") || strings.Contains(r.body, "Mira") {
				t.Errorf("%s %s leaked data: %s", ep[0], ep[1], r.body)
			}
		}
	}
	// non-existent and malformed IDs are indistinguishable from foreign ones
	for _, id := range []string{uuid.NewString(), "not-a-uuid", "1%27%20OR%201=1"} {
		must(t, w.parent2.do("GET", "/api/v1/children/"+id, nil), 404)
	}
	// child 1 is untouched
	r := must(t, w.parent1.do("GET", "/api/v1/children/"+w.child1.String(), nil), 200)
	if !strings.Contains(r.body, "Mira") {
		t.Fatal("child 1 modified by IDOR attempt")
	}
}

func TestDiaryDeletionOwnerOnlyAndAudited(t *testing.T) {
	e := newEnv(t, nil)
	w := e.world()
	entry := e.entry(w.parent1, w.child1, "shared", "delete me")
	path := "/api/v1/children/" + w.child1.String() + "/diary-entries/" + entry.String()
	must(t, w.coParent.do("DELETE", path, nil), 403)
	must(t, w.spec1.do("DELETE", path, nil), 403)
	must(t, w.admin.do("DELETE", path, nil), 403)
	must(t, w.parent1.do("DELETE", path, nil), 204)
	audit := must(t, w.parent1.do("GET", "/api/v1/children/"+w.child1.String()+"/audit-events", nil), 200)
	if !strings.Contains(audit.body, `"action":"diary:delete"`) || strings.Contains(audit.body, "delete me") {
		t.Fatalf("delete audit missing or contains content: %s", audit.body)
	}
	must(t, w.parent1.do("GET", "/api/v1/children/"+w.child1.String()+"/diary-entries/"+entry.String()+"/versions", nil), 404)
}

func TestGuestGets401(t *testing.T) {
	e := newEnv(t, nil)
	w := e.world()
	g := e.client()
	for _, ep := range childEndpoints(w.child1, uuid.New(), uuid.New(), w.spec1.id) {
		if r := g.do(ep[0].(string), ep[1].(string), ep[2]); r.code != 401 {
			t.Errorf("guest %s %s: got %d", ep[0], ep[1], r.code)
		}
	}
	for _, p := range []string{"/api/v1/children", "/api/v1/users", "/api/v1/audit-events", "/api/v1/auth/me"} {
		must(t, g.do("GET", p, nil), 401)
	}
}

func TestEndedAssignmentLosesAccessImmediatelyWithOldToken(t *testing.T) {
	e := newEnv(t, nil)
	w := e.world()
	path := "/api/v1/children/" + w.child1.String() + "/diary-entries"
	must(t, w.spec1.do("GET", path, nil), 200)
	must(t, w.parent1.do("PATCH", "/api/v1/children/"+w.child1.String()+"/specialist-assignments/"+w.spec1.id.String(),
		map[string]any{"ended_at": "2000-01-01T00:00:00Z"}), 200)
	// same access token, issued before the assignment ended
	must(t, w.spec1.do("GET", path, nil), 404)
	must(t, w.spec1.do("GET", "/api/v1/children/"+w.child1.String(), nil), 404)
	var list []any
	must(t, w.spec1.do("GET", "/api/v1/children", nil), 200).json(t, &list)
	if len(list) != 0 {
		t.Fatalf("ended specialist still lists child: %v", list)
	}
	// ending twice is a conflict
	must(t, w.parent1.do("PATCH", "/api/v1/children/"+w.child1.String()+"/specialist-assignments/"+w.spec1.id.String(),
		map[string]any{"ended_at": "2000-01-01T00:00:00Z"}), 409)
}

func TestMassAssignmentCannotCreateAdmin(t *testing.T) {
	e := newEnv(t, nil)
	g := e.client()
	must(t, g.do("POST", "/api/v1/auth/register", map[string]any{"email": "evil@example.test", "password": pw, "first_name": "E", "last_name": "V", "role": "admin"}), 422)
	must(t, g.do("POST", "/api/v1/auth/register", map[string]any{"email": "evil@example.test", "password": pw, "first_name": "E", "last_name": "V", "role": "parent", "status": "active", "specialist_approved": true}), 422)
	r := must(t, g.do("POST", "/api/v1/auth/register", map[string]any{"email": "ok@example.test", "password": pw, "first_name": "O", "last_name": "K", "role": "specialist"}), 201)
	if !strings.Contains(r.body, `"specialist_approved":false`) || strings.Contains(r.body, "password") || strings.Contains(r.body, "failed_login") || strings.Contains(r.body, "locked_until") {
		t.Fatalf("unexpected user DTO: %s", r.body)
	}
	must(t, g.do("POST", "/api/v1/auth/register", map[string]any{"email": "ok@example.test", "password": pw, "first_name": "O", "last_name": "K", "role": "parent"}), 409)
	// a non-admin cannot use the admin endpoint either
	p := e.client()
	must(t, g.do("POST", "/api/v1/auth/register", map[string]any{"email": "p@example.test", "password": pw, "first_name": "P", "last_name": "P", "role": "parent"}), 201)
	must(t, p.login("p@example.test", pw), 200)
	must(t, p.do("POST", "/api/v1/users", map[string]any{"email": "x@example.test", "password": pw, "first_name": "X", "last_name": "X", "role": "admin"}), 403)
	must(t, p.do("PATCH", "/api/v1/users/"+p.id.String(), map[string]any{"status": "active"}), 403)
	// unapproved specialist cannot be assigned
	s := e.client()
	must(t, s.login("ok@example.test", pw), 200)
	child := e.child(p, "Kid")
	must(t, p.do("POST", "/api/v1/children/"+child.String()+"/specialist-assignments", map[string]any{"specialist_id": s.id}), 422)
}

func TestRolePermissionsWithinRelationship(t *testing.T) {
	e := newEnv(t, nil)
	w := e.world()
	c := "/api/v1/children/" + w.child1.String()
	// non-primary guardian: can view, cannot edit/export/delete/assign (403)
	must(t, w.coParent.do("GET", c, nil), 200)
	must(t, w.coParent.do("PATCH", c, map[string]any{"first_name": "X"}), 403)
	must(t, w.coParent.do("POST", c+"/exports", nil), 403)
	must(t, w.coParent.do("DELETE", c, nil), 403)
	must(t, w.coParent.do("POST", c+"/specialist-assignments", map[string]any{"specialist_id": w.spec2.id}), 403)
	// specialist: cannot edit child, export, or read audit
	must(t, w.spec1.do("PATCH", c, map[string]any{"first_name": "X"}), 403)
	must(t, w.spec1.do("POST", c+"/exports", nil), 403)
	must(t, w.spec1.do("GET", c+"/audit-events", nil), 403)
	must(t, w.spec1.do("GET", "/api/v1/audit-events", nil), 403)
	// admin: metadata yes, diary/sessions/export no
	must(t, w.admin.do("GET", c, nil), 200)
	must(t, w.admin.do("GET", c+"/diary-entries", nil), 403)
	must(t, w.admin.do("GET", c+"/sessions", nil), 403)
	must(t, w.admin.do("POST", c+"/exports", nil), 403)
	must(t, w.admin.do("POST", "/api/v1/children", map[string]any{"first_name": "A", "last_name": "B", "birth_date": "2019-01-01", "policy_version": "1"}), 403)
	// primary guardian edits
	must(t, w.parent1.do("PATCH", c, map[string]any{"first_name": "Mira2"}), 200)
	must(t, w.parent1.do("PATCH", c, map[string]any{"birth_date": "2999-01-01"}), 422)
	// a second primary guardian is rejected
	must(t, w.parent1.do("POST", c+"/guardians", map[string]any{"guardian_id": w.parent2.id, "relationship": "parent", "is_primary": true}), 422)
	fmt.Fprint(new(strings.Builder), "ok")
}
