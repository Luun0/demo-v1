package policy

import (
	"fmt"
	"strings"
	"testing"

	"github.com/AknurOE/sen-support-platform/internal/domain"
)

// matrix encodes docs/rbac-and-threat-model.md ("Access matrix") plus
// ADR-0008. Columns: role rel action author assigned visibility => decision.
// rel: none | guardian | primary | active (specialist). '*' matches anything.
// The FIRST matching row wins. Every enumerated cell must match a row and
// every row must be used, so the test covers 100% of cells.
const matrix = `
guest     *        *                  * * *                => unauthenticated
admin     *        user:manage        * * *                => allow
admin     *        audit:system       * * *                => allow
admin     *        child:view         * * *                => allow
admin     *        child:delete       * * *                => allow
admin     *        guardian:link      * * *                => allow
admin     *        assignment:manage  * * *                => allow
admin     *        audit:child        * * *                => allow
admin     *        specialist:directory * * *              => allow
admin     *        *                  * * *                => deny
parent    *        child:create       * * *                => allow
parent    *        specialist:directory * * *              => allow
parent    *        user:manage        * * *                => deny
parent    *        audit:system       * * *                => deny
parent    none     *                  * * *                => hide
parent    *        diary:edit         * * specialists_only => hide
parent    *        diary:versions     * * specialists_only => hide
parent    *        child:view         * * *                => allow
parent    *        session:list       * * *                => allow
parent    *        session:create     * * *                => allow
parent    *        session:update     * * *                => allow
parent    *        diary:list         * * *                => allow
parent    *        diary:versions     * * shared           => allow
parent    *        audit:child        * * *                => allow
parent    primary  child:edit         * * *                => allow
parent    primary  child:delete       * * *                => allow
parent    primary  child:export       * * *                => allow
parent    primary  guardian:link      * * *                => allow
parent    primary  assignment:manage  * * *                => allow
parent    guardian child:edit         * * *                => deny
parent    guardian child:delete       * * *                => deny
parent    guardian child:export       * * *                => deny
parent    guardian guardian:link      * * *                => deny
parent    guardian assignment:manage  * * *                => deny
parent    *        diary:create       * * shared           => allow
parent    *        diary:create       * * specialists_only => deny
parent    *        diary:edit         yes * shared         => allow
parent    *        diary:edit         no  * shared         => deny
parent    *        diary:delete       yes * *              => allow
parent    *        diary:delete       no  * *              => deny
specialist *       child:create       * * *                => deny
specialist *       specialist:directory * * *              => deny
specialist *       user:manage        * * *                => deny
specialist *       audit:system       * * *                => deny
specialist none    *                  * * *                => hide
specialist active  child:view         * * *                => allow
specialist active  session:list       * * *                => allow
specialist active  session:create     * * *                => allow
specialist active  session:update     * yes *              => allow
specialist active  session:update     * no  *              => deny
specialist active  diary:list         * * *                => allow
specialist active  diary:create       * * *                => allow
specialist active  diary:versions     * * *                => allow
specialist active  diary:edit         yes * *              => allow
specialist active  diary:edit         no  * *              => deny
specialist active  diary:delete       * * *                => deny
specialist active  child:edit         * * *                => deny
specialist active  child:delete       * * *                => deny
specialist active  child:export       * * *                => deny
specialist active  guardian:link      * * *                => deny
specialist active  assignment:manage  * * *                => deny
specialist active  audit:child        * * *                => deny
`

type row struct {
	f    [6]string
	want string
	used bool
}

func parse(t *testing.T) []*row {
	t.Helper()
	var rows []*row
	for _, line := range strings.Split(strings.TrimSpace(matrix), "\n") {
		parts := strings.Split(line, "=>")
		f := strings.Fields(parts[0])
		if len(f) != 6 || len(parts) != 2 {
			t.Fatalf("bad row %q", line)
		}
		r := &row{want: strings.TrimSpace(parts[1])}
		copy(r.f[:], f)
		rows = append(rows, r)
	}
	return rows
}

func yn(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func TestAccessMatrix100PercentCells(t *testing.T) {
	rows := parse(t)
	roles := map[string][]string{
		"guest": {"none"}, "admin": {"none"},
		"parent": {"none", "guardian", "primary"}, "specialist": {"none", "active"},
	}
	cells := 0
	for role, rels := range roles {
		for _, rel := range rels {
			for _, a := range AllActions {
				for _, author := range []bool{false, true} {
					for _, assigned := range []bool{false, true} {
						for _, vis := range []domain.Visibility{domain.VisibilityShared, domain.VisibilitySpecialistsOnly} {
							cells++
							key := [6]string{role, rel, string(a), yn(author), yn(assigned), string(vis)}
							var match *row
							for _, r := range rows {
								ok := true
								for i := range key {
									if r.f[i] != "*" && r.f[i] != key[i] {
										ok = false
										break
									}
								}
								if ok {
									match = r
									break
								}
							}
							if match == nil {
								t.Fatalf("cell not covered by matrix: %v", key)
							}
							match.used = true
							var actor *Actor
							if role != "guest" {
								actor = &Actor{Role: domain.Role(role)}
							}
							res := Resource{
								Guardian: rel == "guardian" || rel == "primary", PrimaryGuardian: rel == "primary",
								ActiveSpecialist: rel == "active", Author: author, AssignedSession: assigned, Visibility: vis,
							}
							if got := Decide(actor, a, res).String(); got != match.want {
								t.Errorf("%v: got %s want %s", key, got, match.want)
							}
						}
					}
				}
			}
		}
	}
	for _, r := range rows {
		if !r.used {
			t.Errorf("matrix row never exercised: %v", r.f)
		}
	}
	t.Logf("access matrix: %d cells, %d rows, all covered", cells, len(rows))
}

func TestInvalidRoleIsUnauthenticated(t *testing.T) {
	if d := Decide(&Actor{Role: "root"}, ChildView, Resource{Guardian: true}); d != Unauthenticated {
		t.Fatalf("got %v", d)
	}
}

func TestUnknownActionDefaultsToDeny(t *testing.T) {
	cases := []struct {
		name  string
		actor *Actor
		want  Decision
	}{
		{name: "guest", want: Unauthenticated},
		{name: "admin", actor: &Actor{Role: domain.RoleAdmin}, want: Deny},
		{name: "guardian", actor: &Actor{Role: domain.RoleParent}, want: Deny},
		{name: "assigned specialist", actor: &Actor{Role: domain.RoleSpecialist}, want: Deny},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Decide(tc.actor, Action("unknown:action"), Resource{Guardian: true, ActiveSpecialist: true}); got != tc.want {
				t.Fatalf("unknown action decision = %s, want %s", got, tc.want)
			}
		})
	}
}

func ExampleDecide() {
	parent := &Actor{Role: domain.RoleParent}
	fmt.Println(Decide(parent, ChildView, Resource{}))
	fmt.Println(Decide(parent, DiaryCreate, Resource{Guardian: true, Visibility: domain.VisibilitySpecialistsOnly}))
	// Output:
	// hide
	// deny
}
