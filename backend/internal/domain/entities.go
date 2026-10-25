package domain

import (
	"time"

	"github.com/google/uuid"
)

// User is an account. PasswordHash never leaves the backend.
type User struct {
	ID                  uuid.UUID
	Email               string
	PasswordHash        string
	FirstName           string
	LastName            string
	Role                Role
	Status              UserStatus
	SpecialistApproved  bool
	FailedLoginAttempts int
	LockedUntil         *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// Active reports whether the account may authenticate at t.
func (u User) Active(t time.Time) bool {
	return u.Status == StatusActive && (u.LockedUntil == nil || !u.LockedUntil.After(t))
}

// Child is a supported child. ConsentGranted is required at creation.
type Child struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	BirthDate time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ChildGuardian links a parent/guardian to a child.
type ChildGuardian struct {
	ChildID      uuid.UUID
	GuardianID   uuid.UUID
	Relationship Relationship
	IsPrimary    bool
	CreatedAt    time.Time
}

// ChildSpecialist is an assignment; it grants access only while active.
type ChildSpecialist struct {
	ChildID      uuid.UUID
	SpecialistID uuid.UUID
	AssignedBy   uuid.UUID
	AssignedAt   time.Time
	EndedAt      *time.Time
}

// ActiveAt reports whether the assignment grants access at t.
func (a ChildSpecialist) ActiveAt(t time.Time) bool {
	return !a.AssignedAt.After(t) && (a.EndedAt == nil || a.EndedAt.After(t))
}

// ConsentRecord stores the guardian's consent to a policy version.
type ConsentRecord struct {
	ID            uuid.UUID
	ChildID       uuid.UUID
	GuardianID    uuid.UUID
	PolicyVersion string
	GrantedAt     time.Time
}

// SupportSession is a scheduled support lesson.
type SupportSession struct {
	ID           uuid.UUID
	ChildID      uuid.UUID
	SpecialistID uuid.UUID
	CreatedBy    uuid.UUID
	Title        string
	StartsAt     time.Time
	EndsAt       time.Time
	Mode         SessionMode
	Status       SessionStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// DiaryEntry is a versioned observation about a child.
type DiaryEntry struct {
	ID         uuid.UUID
	ChildID    uuid.UUID
	AuthorID   uuid.UUID
	Type       EntryType
	Visibility Visibility
	Content    string
	ObservedAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// DiaryEntryVersion is an immutable snapshot written on every create/edit.
type DiaryEntryVersion struct {
	ID            uuid.UUID
	EntryID       uuid.UUID
	VersionNumber int
	ChangedBy     uuid.UUID
	Type          EntryType
	Visibility    Visibility
	Content       string
	ObservedAt    time.Time
	ChangedAt     time.Time
}

// RefreshToken stores only the SHA-256 of the opaque token.
type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	FamilyID  uuid.UUID
	TokenHash []byte
	ExpiresAt time.Time
	RotatedAt *time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

// PasswordResetToken stores only the hash of a one-time reset credential.
type PasswordResetToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash []byte
	ExpiresAt time.Time
	UsedAt    *time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

// AuditEvent never contains diary text, passwords or tokens (ADR-0005).
type AuditEvent struct {
	ID         uuid.UUID
	ActorID    *uuid.UUID
	ChildID    *uuid.UUID
	Action     string
	Resource   string
	ResourceID *uuid.UUID
	Outcome    Outcome
	RequestID  string
	OccurredAt time.Time
}
