package httpapi

import (
	"encoding/base64"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/AknurOE/sen-support-platform/internal/application"
	"github.com/AknurOE/sen-support-platform/internal/domain"
)

// userDTO never exposes password hash, failed attempts or lock state (S1).
type userDTO struct {
	ID                 uuid.UUID         `json:"id"`
	Email              string            `json:"email"`
	FirstName          string            `json:"first_name"`
	LastName           string            `json:"last_name"`
	Role               domain.Role       `json:"role"`
	Status             domain.UserStatus `json:"status"`
	SpecialistApproved *bool             `json:"specialist_approved,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

func toUser(u domain.User) userDTO {
	d := userDTO{ID: u.ID, Email: u.Email, FirstName: u.FirstName, LastName: u.LastName, Role: u.Role, Status: u.Status, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt}
	if u.Role == domain.RoleSpecialist {
		ok := u.SpecialistApproved
		d.SpecialistApproved = &ok
	}
	return d
}

type childDTO struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	BirthDate string    `json:"birth_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toChild(c domain.Child) childDTO {
	return childDTO{c.ID, c.FirstName, c.LastName, c.BirthDate.Format(time.DateOnly), c.CreatedAt, c.UpdatedAt}
}

type sessionDTO struct {
	ID           uuid.UUID            `json:"id"`
	ChildID      uuid.UUID            `json:"child_id"`
	SpecialistID uuid.UUID            `json:"specialist_id"`
	CreatedBy    uuid.UUID            `json:"created_by"`
	Title        string               `json:"title"`
	StartsAt     time.Time            `json:"starts_at"`
	EndsAt       time.Time            `json:"ends_at"`
	Mode         domain.SessionMode   `json:"mode"`
	Status       domain.SessionStatus `json:"status"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
}

func toSession(s domain.SupportSession) sessionDTO {
	return sessionDTO{s.ID, s.ChildID, s.SpecialistID, s.CreatedBy, s.Title, s.StartsAt, s.EndsAt, s.Mode, s.Status, s.CreatedAt, s.UpdatedAt}
}

type entryDTO struct {
	ID         uuid.UUID         `json:"id"`
	ChildID    uuid.UUID         `json:"child_id"`
	AuthorID   uuid.UUID         `json:"author_id"`
	Type       domain.EntryType  `json:"type"`
	Visibility domain.Visibility `json:"visibility"`
	Content    string            `json:"content"`
	ObservedAt time.Time         `json:"observed_at"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

func toEntry(e domain.DiaryEntry) entryDTO {
	return entryDTO{e.ID, e.ChildID, e.AuthorID, e.Type, e.Visibility, e.Content, e.ObservedAt, e.CreatedAt, e.UpdatedAt}
}

type versionDTO struct {
	ID            uuid.UUID         `json:"id"`
	EntryID       uuid.UUID         `json:"entry_id"`
	VersionNumber int               `json:"version_number"`
	ChangedBy     uuid.UUID         `json:"changed_by"`
	Type          domain.EntryType  `json:"type"`
	Visibility    domain.Visibility `json:"visibility"`
	Content       string            `json:"content"`
	ObservedAt    time.Time         `json:"observed_at"`
	ChangedAt     time.Time         `json:"changed_at"`
}

type auditDTO struct {
	ID         uuid.UUID      `json:"id"`
	ActorID    *uuid.UUID     `json:"actor_id"`
	ChildID    *uuid.UUID     `json:"child_id"`
	Action     string         `json:"action"`
	Resource   string         `json:"resource"`
	ResourceID *uuid.UUID     `json:"resource_id,omitempty"`
	Outcome    domain.Outcome `json:"outcome"`
	RequestID  string         `json:"request_id"`
	OccurredAt time.Time      `json:"occurred_at"`
}

func mapSlice[T, U any](in []T, f func(T) U) []U {
	out := make([]U, 0, len(in))
	for _, v := range in {
		out = append(out, f(v))
	}
	return out
}

func encodeCursor(c *application.Cursor) string {
	if c == nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString([]byte(strconv.FormatInt(c.At.UnixNano(), 10) + "|" + c.ID.String()))
}

func decodeCursor(s string) (*application.Cursor, error) {
	if s == "" {
		return nil, nil
	}
	bad := &domain.ValidationError{Fields: []domain.FieldError{{Field: "cursor", Code: "invalid"}}}
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, bad
	}
	ts, id, ok := strings.Cut(string(b), "|")
	n, err1 := strconv.ParseInt(ts, 10, 64)
	u, err2 := uuid.Parse(id)
	if !ok || err1 != nil || err2 != nil {
		return nil, bad
	}
	return &application.Cursor{At: time.Unix(0, n).UTC(), ID: u}, nil
}

func parseLimit(s string) (int, error) {
	if s == "" {
		return 50, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 || n > 100 {
		return 0, &domain.ValidationError{Fields: []domain.FieldError{{Field: "limit", Code: "invalid"}}}
	}
	return n, nil
}
