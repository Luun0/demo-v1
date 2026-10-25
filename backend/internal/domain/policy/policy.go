// Package policy is the single, pure authorization function (ADR-0001,
// ADR-0008). Default is deny. Handlers never make access decisions.
package policy

import "github.com/AknurOE/sen-support-platform/internal/domain"

// Action is an operation subject to authorization.
type Action string

// Actions covered by the access matrix (docs/rbac-and-threat-model.md).
const (
	UserManage       Action = "user:manage"
	SystemAudit      Action = "audit:system"
	ChildCreate      Action = "child:create"
	ChildView        Action = "child:view"
	ChildEdit        Action = "child:edit"
	ChildDelete      Action = "child:delete"
	ChildExport      Action = "child:export"
	GuardianLink     Action = "guardian:link"
	AssignmentManage Action = "assignment:manage"
	SessionList      Action = "session:list"
	SessionCreate    Action = "session:create"
	SessionUpdate    Action = "session:update"
	DiaryList        Action = "diary:list"
	DiaryCreate      Action = "diary:create"
	DiaryEdit        Action = "diary:edit"
	DiaryDelete      Action = "diary:delete"
	DiaryVersions    Action = "diary:versions"
	ChildAudit       Action = "audit:child"
	SpecialistDir    Action = "specialist:directory"
)

// AllActions lists every action, used by the access-matrix test.
var AllActions = []Action{
	UserManage, SystemAudit, ChildCreate, ChildView, ChildEdit, ChildDelete, ChildExport,
	GuardianLink, AssignmentManage, SessionList, SessionCreate, SessionUpdate,
	DiaryList, DiaryCreate, DiaryEdit, DiaryDelete, DiaryVersions, ChildAudit, SpecialistDir,
}

// Decision is the outcome of Decide.
type Decision int

// Decisions map to HTTP 200 / 401 / 403 / 404.
const (
	Deny Decision = iota // 403
	Allow
	Unauthenticated // 401
	Hide            // 404: do not reveal existence (SEC-11)
)

func (d Decision) String() string {
	return [...]string{"deny", "allow", "unauthenticated", "hide"}[d]
}

// Actor is the authenticated caller; nil means guest.
type Actor struct {
	ID   [16]byte
	Role domain.Role
}

// Resource describes the caller's relationship to the target, computed by
// the application layer from the database on every request.
type Resource struct {
	Guardian         bool // caller is a linked guardian of the child
	PrimaryGuardian  bool
	ActiveSpecialist bool // caller has an active assignment
	Author           bool // caller authored the diary entry
	AssignedSession  bool // caller is the specialist of the session
	// Visibility of the target entry (create: requested; edit/versions: stored).
	Visibility domain.Visibility
}

func childScoped(a Action) bool {
	switch a {
	case UserManage, SystemAudit, ChildCreate, SpecialistDir:
		return false
	}
	return true
}

// Decide returns the access decision. It is total and side-effect free.
func Decide(actor *Actor, action Action, r Resource) Decision {
	if actor == nil || !actor.Role.Valid() {
		return Unauthenticated
	}
	switch actor.Role {
	case domain.RoleAdmin:
		return decideAdmin(action)
	case domain.RoleParent:
		return decideParent(action, r)
	case domain.RoleSpecialist:
		return decideSpecialist(action, r)
	}
	return Deny
}

func decideAdmin(a Action) Decision {
	switch a {
	case UserManage, SystemAudit, ChildView, ChildDelete, GuardianLink, AssignmentManage, ChildAudit, SpecialistDir:
		return Allow
	}
	return Deny
}

func decideParent(a Action, r Resource) Decision {
	if a == ChildCreate || a == SpecialistDir {
		return Allow
	}
	if !childScoped(a) {
		return Deny
	}
	if !r.Guardian {
		return Hide
	}
	// specialists_only entries do not exist for a parent (SEC-12).
	if (a == DiaryEdit || a == DiaryVersions) && r.Visibility != domain.VisibilityShared {
		return Hide
	}
	switch a {
	case ChildView, SessionList, SessionCreate, SessionUpdate, DiaryList, DiaryVersions, ChildAudit:
		return Allow
	case ChildEdit, ChildDelete, ChildExport, GuardianLink, AssignmentManage:
		return boolDecision(r.PrimaryGuardian)
	case DiaryCreate:
		return boolDecision(domain.ParentVisibilityAllowed(r.Visibility))
	case DiaryEdit:
		return boolDecision(r.Author && domain.ParentVisibilityAllowed(r.Visibility))
	case DiaryDelete:
		return boolDecision(r.Author)
	}
	return Deny
}

func decideSpecialist(a Action, r Resource) Decision {
	if !childScoped(a) {
		return Deny
	}
	if !r.ActiveSpecialist {
		return Hide
	}
	switch a {
	case ChildView, SessionList, SessionCreate, DiaryList, DiaryCreate, DiaryVersions:
		return Allow
	case SessionUpdate:
		return boolDecision(r.AssignedSession)
	case DiaryEdit:
		return boolDecision(r.Author)
	case DiaryDelete:
		return Deny // specialist observations are retained
	}
	return Deny
}

func boolDecision(ok bool) Decision {
	if ok {
		return Allow
	}
	return Deny
}
