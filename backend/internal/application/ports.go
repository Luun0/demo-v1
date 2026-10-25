// Package application contains use cases and the ports they depend on.
// It must not import net/http, pgx or jwt (ADR-0001, depguard).
package application

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/AknurOE/sen-support-platform/internal/domain"
)

// DiaryFilter selects diary entries. SharedOnly MUST be applied in SQL.
type DiaryFilter struct {
	ChildID    uuid.UUID
	SharedOnly bool
	Visibility domain.Visibility // optional extra filter
	Query      string
	Limit      int
	After      *Cursor
}

// Cursor is a keyset position (time, id) for descending lists.
type Cursor struct {
	At time.Time
	ID uuid.UUID
}

// Repo is the persistence port. Implementations must use parameterised SQL.
type Repo interface {
	CreateUser(ctx context.Context, u *domain.User) error
	UserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	UserByEmail(ctx context.Context, email string) (domain.User, error)
	UpdateUser(ctx context.Context, u *domain.User) error
	ListUsers(ctx context.Context, limit int, after *uuid.UUID) ([]domain.User, error)

	InsertRefreshToken(ctx context.Context, t *domain.RefreshToken) error
	RefreshTokenByHash(ctx context.Context, hash []byte) (domain.RefreshToken, error)
	MarkRefreshRotated(ctx context.Context, id uuid.UUID, at time.Time) error
	RevokeRefreshFamily(ctx context.Context, family uuid.UUID, at time.Time) error
	RevokeUserRefreshTokens(ctx context.Context, user uuid.UUID, at time.Time) error

	CreateChild(ctx context.Context, c *domain.Child) error
	ChildByID(ctx context.Context, id uuid.UUID) (domain.Child, error)
	UpdateChild(ctx context.Context, c *domain.Child) error
	DeleteChild(ctx context.Context, id uuid.UUID) error
	ListChildrenForGuardian(ctx context.Context, user uuid.UUID) ([]domain.Child, error)
	ListChildrenForSpecialist(ctx context.Context, user uuid.UUID, now time.Time) ([]domain.Child, error)
	ListAllChildren(ctx context.Context) ([]domain.Child, error)
	AddGuardian(ctx context.Context, g domain.ChildGuardian) error
	Guardian(ctx context.Context, child, user uuid.UUID) (domain.ChildGuardian, error)
	AddConsent(ctx context.Context, c domain.ConsentRecord) error

	UpsertAssignment(ctx context.Context, a domain.ChildSpecialist) error
	Assignment(ctx context.Context, child, specialist uuid.UUID) (domain.ChildSpecialist, error)
	EndAssignment(ctx context.Context, child, specialist uuid.UUID, at time.Time) error

	CreateSession(ctx context.Context, s *domain.SupportSession) error
	SessionByID(ctx context.Context, child, id uuid.UUID) (domain.SupportSession, error)
	UpdateSession(ctx context.Context, s *domain.SupportSession) error
	ListSessions(ctx context.Context, child uuid.UUID) ([]domain.SupportSession, error)

	CreateEntry(ctx context.Context, e *domain.DiaryEntry) error
	EntryByID(ctx context.Context, child, id uuid.UUID, sharedOnly bool) (domain.DiaryEntry, error)
	UpdateEntry(ctx context.Context, e *domain.DiaryEntry) error
	InsertVersion(ctx context.Context, v *domain.DiaryEntryVersion) error
	NextVersionNumber(ctx context.Context, entry uuid.UUID) (int, error)
	ListVersions(ctx context.Context, entry uuid.UUID) ([]domain.DiaryEntryVersion, error)
	ListEntries(ctx context.Context, f DiaryFilter) ([]domain.DiaryEntry, error)
	CountEntries(ctx context.Context, f DiaryFilter) (int, error)

	TeamRepo

	InsertAudit(ctx context.Context, e *domain.AuditEvent) error
	ListAudit(ctx context.Context, child *uuid.UUID, limit int, after *Cursor) ([]domain.AuditEvent, error)
}

// UnitOfWork runs fn in one transaction; Repo() is for single statements.
type UnitOfWork interface {
	Within(ctx context.Context, fn func(Repo) error) error
	Repo() Repo
	Ping(ctx context.Context) error
}

// Clock abstracts time for tests.
type Clock interface{ Now() time.Time }

// PasswordHasher hashes and verifies passwords (Argon2id).
type PasswordHasher interface {
	Hash(password string) (string, error)
	// Verify returns ok and whether the hash uses outdated parameters.
	Verify(encoded, password string) (ok, needsRehash bool, err error)
}

// TokenIssuer issues and verifies short-lived access tokens.
type TokenIssuer interface {
	Issue(user uuid.UUID, now time.Time) (token string, ttl time.Duration, err error)
	Verify(token string, now time.Time) (uuid.UUID, error)
}

// Limiter is a keyed rate limiter.
type Limiter interface {
	Allow(key string, now time.Time) bool
}

// Actor is the authenticated caller, loaded from the database per request.
type Actor struct {
	ID   uuid.UUID
	Role domain.Role
}

type ctxKey int

const requestIDKey ctxKey = 1

// WithRequestID stores the request id for audit events.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// RequestID returns the request id or "-".
func RequestID(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok && v != "" {
		return v
	}
	return "-"
}

// TeamMember is a guardian or specialist linked to a child (names only).
type TeamMember struct {
	UserID       uuid.UUID
	FirstName    string
	LastName     string
	Role         domain.Role
	IsPrimary    bool
	Relationship domain.Relationship
	AssignedAt   *time.Time
	EndedAt      *time.Time
}

// TeamRepo lists relationships for the team screen.
type TeamRepo interface {
	ListTeam(ctx context.Context, child uuid.UUID) ([]TeamMember, error)
	ListApprovedSpecialists(ctx context.Context) ([]domain.User, error)
}
