// Package domain holds entities, enumerations and invariants. It has no
// knowledge of HTTP, SQL or token formats (ADR-0001).
package domain

// Role is a platform role. Guest is the absence of an authenticated actor.
type Role string

const (
	RoleParent     Role = "parent"
	RoleSpecialist Role = "specialist"
	RoleAdmin      Role = "admin"
)

func (r Role) Valid() bool { return r == RoleParent || r == RoleSpecialist || r == RoleAdmin }

// UserStatus is the account status.
type UserStatus string

const (
	StatusActive   UserStatus = "active"
	StatusInactive UserStatus = "inactive"
)

// Visibility of a diary entry.
type Visibility string

const (
	VisibilityShared          Visibility = "shared"
	VisibilitySpecialistsOnly Visibility = "specialists_only"
)

func (v Visibility) Valid() bool { return v == VisibilityShared || v == VisibilitySpecialistsOnly }

// EntryType of a diary entry.
type EntryType string

const (
	EntryObservation    EntryType = "observation"
	EntryProgress       EntryType = "progress"
	EntryRecommendation EntryType = "recommendation"
)

func (t EntryType) Valid() bool {
	return t == EntryObservation || t == EntryProgress || t == EntryRecommendation
}

// SessionMode of a support session.
type SessionMode string

const (
	ModeOnline   SessionMode = "online"
	ModeInPerson SessionMode = "in_person"
)

func (m SessionMode) Valid() bool { return m == ModeOnline || m == ModeInPerson }

// SessionStatus of a support session.
type SessionStatus string

const (
	SessionScheduled SessionStatus = "scheduled"
	SessionCompleted SessionStatus = "completed"
	SessionCancelled SessionStatus = "cancelled"
)

func (s SessionStatus) Valid() bool {
	return s == SessionScheduled || s == SessionCompleted || s == SessionCancelled
}

// Relationship of a guardian to a child.
type Relationship string

const (
	RelParent   Relationship = "parent"
	RelGuardian Relationship = "guardian"
)

func (r Relationship) Valid() bool { return r == RelParent || r == RelGuardian }

// Outcome of an audited action.
type Outcome string

const (
	OutcomeSuccess Outcome = "success"
	OutcomeDenied  Outcome = "denied"
	OutcomeFailure Outcome = "failure"
)
