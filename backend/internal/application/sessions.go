package application

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/AknurOE/sen-support-platform/internal/domain"
	"github.com/AknurOE/sen-support-platform/internal/domain/policy"
)

// SessionInput creates a support session.
type SessionInput struct {
	SpecialistID     uuid.UUID
	Title            string
	StartsAt, EndsAt time.Time
	Mode             domain.SessionMode
	Status           domain.SessionStatus
}

// SessionPatch updates a support session.
type SessionPatch struct {
	Title            *string
	StartsAt, EndsAt *time.Time
	Mode             *domain.SessionMode
	Status           *domain.SessionStatus
}

func validateSession(s domain.SupportSession) error {
	var v domain.Validator
	v.Check(domain.ValidText(s.Title, domain.MaxTitleLen), "title", "invalid")
	v.Check(!s.StartsAt.IsZero() && s.EndsAt.After(s.StartsAt), "ends_at", "before_start")
	v.Check(s.Mode.Valid(), "mode", "invalid")
	v.Check(s.Status.Valid(), "status", "invalid")
	return v.Err()
}

// ListSessions lists the child's sessions.
func (s *Service) ListSessions(ctx context.Context, a *Actor, child uuid.UUID) ([]domain.SupportSession, error) {
	if _, err := s.guard(ctx, a, child, policy.SessionList, "support_session", nil, nil); err != nil {
		return nil, err
	}
	out, err := s.uow.Repo().ListSessions(ctx, child)
	if err != nil {
		return nil, err
	}
	return out, s.readAudit(ctx, a, &child, policy.SessionList, "support_session", nil)
}

// CreateSession schedules a session with an actively assigned specialist.
func (s *Service) CreateSession(ctx context.Context, a *Actor, child uuid.UUID, in SessionInput) (domain.SupportSession, error) {
	if _, err := s.guard(ctx, a, child, policy.SessionCreate, "support_session", nil, nil); err != nil {
		return domain.SupportSession{}, err
	}
	now := s.now()
	ss := domain.SupportSession{ID: newID(), ChildID: child, SpecialistID: in.SpecialistID, CreatedBy: a.ID,
		Title: strings.TrimSpace(in.Title), StartsAt: in.StartsAt.UTC(), EndsAt: in.EndsAt.UTC(), Mode: in.Mode,
		Status: in.Status, CreatedAt: now, UpdatedAt: now}
	if err := validateSession(ss); err != nil {
		return ss, err
	}
	err := s.uow.Within(ctx, func(r Repo) error {
		as, err := r.Assignment(ctx, child, in.SpecialistID)
		if err != nil || !as.ActiveAt(now) {
			return &domain.ValidationError{Fields: []domain.FieldError{{Field: "specialist_id", Code: "not_assigned"}}}
		}
		if err := r.CreateSession(ctx, &ss); err != nil {
			return err
		}
		return r.InsertAudit(ctx, s.event(ctx, a, &child, string(policy.SessionCreate), "support_session", &ss.ID, domain.OutcomeSuccess))
	})
	return ss, err
}

// UpdateSession edits a session (specialist: only own sessions).
func (s *Service) UpdateSession(ctx context.Context, a *Actor, child, id uuid.UUID, p SessionPatch) (domain.SupportSession, error) {
	var ss domain.SupportSession
	_, err := s.guard(ctx, a, child, policy.SessionUpdate, "support_session", &id, func(res *policy.Resource) error {
		var err error
		ss, err = s.uow.Repo().SessionByID(ctx, child, id)
		res.AssignedSession = err == nil && a != nil && ss.SpecialistID == a.ID
		return err
	})
	if err != nil {
		return ss, err
	}
	if p.Title != nil {
		ss.Title = strings.TrimSpace(*p.Title)
	}
	if p.StartsAt != nil {
		ss.StartsAt = p.StartsAt.UTC()
	}
	if p.EndsAt != nil {
		ss.EndsAt = p.EndsAt.UTC()
	}
	if p.Mode != nil {
		ss.Mode = *p.Mode
	}
	if p.Status != nil {
		ss.Status = *p.Status
	}
	if err := validateSession(ss); err != nil {
		return ss, err
	}
	ss.UpdatedAt = s.now()
	err = s.uow.Within(ctx, func(r Repo) error {
		if err := r.UpdateSession(ctx, &ss); err != nil {
			return err
		}
		return r.InsertAudit(ctx, s.event(ctx, a, &child, string(policy.SessionUpdate), "support_session", &id, domain.OutcomeSuccess))
	})
	return ss, err
}
