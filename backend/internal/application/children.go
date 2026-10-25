package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/AknurOE/sen-support-platform/internal/domain"
	"github.com/AknurOE/sen-support-platform/internal/domain/policy"
)

// ChildInput creates a child with guardian consent (PRV, ADR-0007).
type ChildInput struct {
	FirstName, LastName, PolicyVersion string
	BirthDate                          time.Time
}

// ChildPatch updates optional child fields.
type ChildPatch struct {
	FirstName, LastName *string
	BirthDate           *time.Time
}

// CreateChild creates a child; the caller becomes primary guardian.
func (s *Service) CreateChild(ctx context.Context, a *Actor, in ChildInput) (domain.Child, error) {
	if err := s.guardGlobal(ctx, a, policy.ChildCreate, "child"); err != nil {
		return domain.Child{}, err
	}
	now := s.now()
	var v domain.Validator
	v.Check(domain.ValidName(in.FirstName), "first_name", "invalid")
	v.Check(domain.ValidName(in.LastName), "last_name", "invalid")
	v.Check(domain.ValidBirthDate(in.BirthDate, now), "birth_date", "invalid")
	v.Check(domain.ValidText(in.PolicyVersion, 50), "policy_version", "required")
	if err := v.Err(); err != nil {
		return domain.Child{}, err
	}
	c := domain.Child{ID: newID(), FirstName: strings.TrimSpace(in.FirstName), LastName: strings.TrimSpace(in.LastName),
		BirthDate: in.BirthDate, CreatedAt: now, UpdatedAt: now}
	err := s.uow.Within(ctx, func(r Repo) error {
		if err := r.CreateChild(ctx, &c); err != nil {
			return err
		}
		if err := r.AddGuardian(ctx, domain.ChildGuardian{ChildID: c.ID, GuardianID: a.ID, Relationship: domain.RelParent, IsPrimary: true, CreatedAt: now}); err != nil {
			return err
		}
		if err := r.AddConsent(ctx, domain.ConsentRecord{ID: newID(), ChildID: c.ID, GuardianID: a.ID, PolicyVersion: in.PolicyVersion, GrantedAt: now}); err != nil {
			return err
		}
		return r.InsertAudit(ctx, s.event(ctx, a, &c.ID, string(policy.ChildCreate), "child", &c.ID, domain.OutcomeSuccess))
	})
	return c, err
}

// ListChildren returns children the caller is related to (admin: all).
func (s *Service) ListChildren(ctx context.Context, a *Actor) ([]domain.Child, error) {
	if a == nil {
		return nil, domain.ErrUnauthorized
	}
	r := s.uow.Repo()
	switch a.Role {
	case domain.RoleParent:
		return r.ListChildrenForGuardian(ctx, a.ID)
	case domain.RoleSpecialist:
		return r.ListChildrenForSpecialist(ctx, a.ID, s.now())
	case domain.RoleAdmin:
		return r.ListAllChildren(ctx)
	}
	return nil, domain.ErrForbidden
}

// GetChild returns the child profile.
func (s *Service) GetChild(ctx context.Context, a *Actor, id uuid.UUID) (domain.Child, error) {
	if _, err := s.guard(ctx, a, id, policy.ChildView, "child", &id, nil); err != nil {
		return domain.Child{}, err
	}
	c, err := s.uow.Repo().ChildByID(ctx, id)
	if err != nil {
		return c, err
	}
	return c, s.readAudit(ctx, a, &id, policy.ChildView, "child", &id)
}

// UpdateChild edits the profile (primary guardian only).
func (s *Service) UpdateChild(ctx context.Context, a *Actor, id uuid.UUID, p ChildPatch) (domain.Child, error) {
	if _, err := s.guard(ctx, a, id, policy.ChildEdit, "child", &id, nil); err != nil {
		return domain.Child{}, err
	}
	var c domain.Child
	err := s.uow.Within(ctx, func(r Repo) error {
		var err error
		if c, err = r.ChildByID(ctx, id); err != nil {
			return err
		}
		var v domain.Validator
		if p.FirstName != nil {
			v.Check(domain.ValidName(*p.FirstName), "first_name", "invalid")
			c.FirstName = strings.TrimSpace(*p.FirstName)
		}
		if p.LastName != nil {
			v.Check(domain.ValidName(*p.LastName), "last_name", "invalid")
			c.LastName = strings.TrimSpace(*p.LastName)
		}
		if p.BirthDate != nil {
			v.Check(domain.ValidBirthDate(*p.BirthDate, s.now()), "birth_date", "invalid")
			c.BirthDate = *p.BirthDate
		}
		if err := v.Err(); err != nil {
			return err
		}
		c.UpdatedAt = s.now()
		if err := r.UpdateChild(ctx, &c); err != nil {
			return err
		}
		return r.InsertAudit(ctx, s.event(ctx, a, &id, string(policy.ChildEdit), "child", &id, domain.OutcomeSuccess))
	})
	return c, err
}

// DeleteChild deletes all child data (primary guardian or admin).
func (s *Service) DeleteChild(ctx context.Context, a *Actor, id uuid.UUID) error {
	if _, err := s.guard(ctx, a, id, policy.ChildDelete, "child", &id, nil); err != nil {
		return err
	}
	return s.uow.Within(ctx, func(r Repo) error {
		if err := r.DeleteChild(ctx, id); err != nil {
			return err
		}
		return r.InsertAudit(ctx, s.event(ctx, a, &id, string(policy.ChildDelete), "child", &id, domain.OutcomeSuccess))
	})
}

// LinkGuardian links another parent account to the child.
func (s *Service) LinkGuardian(ctx context.Context, a *Actor, child uuid.UUID, g domain.ChildGuardian) (domain.ChildGuardian, error) {
	if _, err := s.guard(ctx, a, child, policy.GuardianLink, "child_guardian", &g.GuardianID, nil); err != nil {
		return g, err
	}
	var v domain.Validator
	v.Check(g.Relationship.Valid(), "relationship", "invalid")
	if err := v.Err(); err != nil {
		return g, err
	}
	g.ChildID, g.CreatedAt = child, s.now()
	err := s.uow.Within(ctx, func(r Repo) error {
		u, err := r.UserByID(ctx, g.GuardianID)
		if err != nil || u.Role != domain.RoleParent || u.Status != domain.StatusActive {
			return &domain.ValidationError{Fields: []domain.FieldError{{Field: "guardian_id", Code: "not_a_parent"}}}
		}
		if err := r.AddGuardian(ctx, g); err != nil {
			return err
		}
		return r.InsertAudit(ctx, s.event(ctx, a, &child, string(policy.GuardianLink), "child_guardian", &g.GuardianID, domain.OutcomeSuccess))
	})
	return g, err
}

// AssignSpecialist assigns an approved, active specialist.
func (s *Service) AssignSpecialist(ctx context.Context, a *Actor, child, spec uuid.UUID) (domain.ChildSpecialist, error) {
	if _, err := s.guard(ctx, a, child, policy.AssignmentManage, "child_specialist", &spec, nil); err != nil {
		return domain.ChildSpecialist{}, err
	}
	as := domain.ChildSpecialist{ChildID: child, SpecialistID: spec, AssignedBy: a.ID, AssignedAt: s.now()}
	err := s.uow.Within(ctx, func(r Repo) error {
		u, err := r.UserByID(ctx, spec)
		if err != nil || u.Role != domain.RoleSpecialist || u.Status != domain.StatusActive || !u.SpecialistApproved {
			return &domain.ValidationError{Fields: []domain.FieldError{{Field: "specialist_id", Code: "not_an_approved_specialist"}}}
		}
		if err := r.UpsertAssignment(ctx, as); err != nil {
			return err
		}
		return r.InsertAudit(ctx, s.event(ctx, a, &child, string(policy.AssignmentManage), "child_specialist", &spec, domain.OutcomeSuccess))
	})
	return as, err
}

// EndAssignment ends access immediately or at endedAt (not in the past).
func (s *Service) EndAssignment(ctx context.Context, a *Actor, child, spec uuid.UUID, endedAt time.Time) (domain.ChildSpecialist, error) {
	if _, err := s.guard(ctx, a, child, policy.AssignmentManage, "child_specialist", &spec, nil); err != nil {
		return domain.ChildSpecialist{}, err
	}
	now := s.now()
	if endedAt.IsZero() || endedAt.Before(now) {
		endedAt = now
	}
	var as domain.ChildSpecialist
	err := s.uow.Within(ctx, func(r Repo) error {
		var err error
		if as, err = r.Assignment(ctx, child, spec); err != nil {
			return err
		}
		if as.EndedAt != nil && !as.EndedAt.After(now) {
			return domain.ErrConflict
		}
		if err := r.EndAssignment(ctx, child, spec, endedAt); err != nil {
			return err
		}
		as.EndedAt = &endedAt
		return r.InsertAudit(ctx, s.event(ctx, a, &child, "assignment:end", "child_specialist", &spec, domain.OutcomeSuccess))
	})
	if errors.Is(err, domain.ErrNotFound) {
		return as, domain.ErrNotFound
	}
	return as, err
}
