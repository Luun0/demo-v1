package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/AknurOE/sen-support-platform/internal/domain"
	"github.com/AknurOE/sen-support-platform/internal/domain/policy"
)

// Config holds use-case tunables.
type Config struct {
	RefreshTTL      time.Duration
	MaxFailedLogins int
	LockoutDuration time.Duration
}

// Service implements all use cases.
type Service struct {
	uow          UnitOfWork
	clock        Clock
	hasher       PasswordHasher
	tokens       TokenIssuer
	loginLimiter Limiter
	cfg          Config
	dummyHash    string
}

// NewService wires the use cases. dummyHash is verified when an email is
// unknown so that response time does not reveal account existence.
func NewService(uow UnitOfWork, clock Clock, h PasswordHasher, t TokenIssuer, l Limiter, cfg Config) (*Service, error) {
	dummy, err := h.Hash("dummy-password-for-timing")
	if err != nil {
		return nil, err
	}
	return &Service{uow: uow, clock: clock, hasher: h, tokens: t, loginLimiter: l, cfg: cfg, dummyHash: dummy}, nil
}

// Ready reports whether the database is reachable.
func (s *Service) Ready(ctx context.Context) error { return s.uow.Ping(ctx) }

func newID() uuid.UUID { return uuid.Must(uuid.NewV7()) }

func (s *Service) now() time.Time { return s.clock.Now().UTC() }

// event builds an audit event; diary content is never part of it.
func (s *Service) event(ctx context.Context, a *Actor, child *uuid.UUID, action, resource string, rid *uuid.UUID, o domain.Outcome) *domain.AuditEvent {
	var actor *uuid.UUID
	if a != nil {
		id := a.ID
		actor = &id
	}
	return &domain.AuditEvent{ID: newID(), ActorID: actor, ChildID: child, Action: action, Resource: resource,
		ResourceID: rid, Outcome: o, RequestID: RequestID(ctx), OccurredAt: s.now()}
}

// auditOutside records a denied/failed action outside the main transaction.
func (s *Service) auditOutside(ctx context.Context, a *Actor, child *uuid.UUID, action, resource string, rid *uuid.UUID, o domain.Outcome) {
	_ = s.uow.Repo().InsertAudit(context.WithoutCancel(ctx), s.event(ctx, a, child, action, resource, rid, o))
}

// relation computes the caller's relationship to a child from the database.
func (s *Service) relation(ctx context.Context, r Repo, a *Actor, child uuid.UUID) (policy.Resource, error) {
	var res policy.Resource
	if a == nil {
		return res, nil
	}
	if _, err := r.ChildByID(ctx, child); err != nil {
		return res, err
	}
	switch a.Role {
	case domain.RoleParent:
		g, err := r.Guardian(ctx, child, a.ID)
		if err == nil {
			res.Guardian, res.PrimaryGuardian = true, g.IsPrimary
		} else if !errors.Is(err, domain.ErrNotFound) {
			return res, err
		}
	case domain.RoleSpecialist:
		as, err := r.Assignment(ctx, child, a.ID)
		if err == nil {
			res.ActiveSpecialist = as.ActiveAt(s.now())
		} else if !errors.Is(err, domain.ErrNotFound) {
			return res, err
		}
	}
	return res, nil
}

func decisionErr(d policy.Decision) error {
	switch d {
	case policy.Allow:
		return nil
	case policy.Hide:
		return domain.ErrNotFound
	case policy.Unauthenticated:
		return domain.ErrUnauthorized
	}
	return domain.ErrForbidden
}

func toPolicy(a *Actor) *policy.Actor {
	if a == nil {
		return nil
	}
	return &policy.Actor{ID: a.ID, Role: a.Role}
}

// guard is the single authorization entry point for child-scoped use cases.
// On deny it writes exactly one denied audit event and returns the error.
// A missing child yields ErrNotFound for everyone except admin semantics are
// identical (child does not exist).
func (s *Service) guard(ctx context.Context, a *Actor, child uuid.UUID, action policy.Action, resource string, rid *uuid.UUID, mutate func(*policy.Resource) error) (policy.Resource, error) {
	r := s.uow.Repo()
	res, err := s.relation(ctx, r, a, child)
	if err == nil && mutate != nil {
		err = mutate(&res)
	}
	if err == nil {
		err = decisionErr(policy.Decide(toPolicy(a), action, res))
	}
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) || errors.Is(err, domain.ErrForbidden) {
			s.auditOutside(ctx, a, &child, string(action), resource, rid, domain.OutcomeDenied)
		}
		return res, err
	}
	return res, nil
}

// guardGlobal authorizes non-child actions (user management, system audit).
func (s *Service) guardGlobal(ctx context.Context, a *Actor, action policy.Action, resource string) error {
	if err := decisionErr(policy.Decide(toPolicy(a), action, policy.Resource{})); err != nil {
		if !errors.Is(err, domain.ErrUnauthorized) {
			s.auditOutside(ctx, a, nil, string(action), resource, nil, domain.OutcomeDenied)
		}
		return err
	}
	return nil
}

// readAudit records a successful read (AUD-01).
func (s *Service) readAudit(ctx context.Context, a *Actor, child *uuid.UUID, action policy.Action, resource string, rid *uuid.UUID) error {
	return s.uow.Repo().InsertAudit(ctx, s.event(ctx, a, child, string(action), resource, rid, domain.OutcomeSuccess))
}
