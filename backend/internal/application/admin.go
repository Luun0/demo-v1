package application

import (
	"context"

	"github.com/google/uuid"

	"github.com/AknurOE/sen-support-platform/internal/domain"
	"github.com/AknurOE/sen-support-platform/internal/domain/policy"
)

// UserPatch is the admin user update (deactivate, approve specialist).
type UserPatch struct {
	Status             *domain.UserStatus
	SpecialistApproved *bool
}

// CreateUserByAdmin creates any account, including admin (SEC-14).
func (s *Service) CreateUserByAdmin(ctx context.Context, a *Actor, in RegisterInput, approved bool) (domain.User, error) {
	if err := s.guardGlobal(ctx, a, policy.UserManage, "user"); err != nil {
		return domain.User{}, err
	}
	in.Email = domain.NormalizeEmail(in.Email)
	if err := validateNewUser(in.Email, in.Password, in.FirstName, in.LastName, in.Role, domain.RoleParent, domain.RoleSpecialist, domain.RoleAdmin); err != nil {
		return domain.User{}, err
	}
	u, err := s.createUser(ctx, a, in.Email, in.Password, in.FirstName, in.LastName, in.Role, string(policy.UserManage))
	if err != nil || !approved || in.Role != domain.RoleSpecialist {
		return u, err
	}
	t := true
	return s.UpdateUser(ctx, a, u.ID, UserPatch{SpecialistApproved: &t})
}

// BootstrapAdmin creates the first admin from the CLI (no actor).
func (s *Service) BootstrapAdmin(ctx context.Context, email, pw, first, last string) (domain.User, error) {
	email = domain.NormalizeEmail(email)
	if err := validateNewUser(email, pw, first, last, domain.RoleAdmin, domain.RoleAdmin); err != nil {
		return domain.User{}, err
	}
	return s.createUser(ctx, nil, email, pw, first, last, domain.RoleAdmin, "user:bootstrap_admin")
}

// UpdateUser deactivates or approves a user. Deactivation revokes sessions.
func (s *Service) UpdateUser(ctx context.Context, a *Actor, id uuid.UUID, p UserPatch) (domain.User, error) {
	if err := s.guardGlobal(ctx, a, policy.UserManage, "user"); err != nil {
		return domain.User{}, err
	}
	var v domain.Validator
	v.Check(p.Status != nil || p.SpecialistApproved != nil, "body", "empty")
	if p.Status != nil {
		v.Check(*p.Status == domain.StatusInactive || *p.Status == domain.StatusActive, "status", "invalid")
	}
	if err := v.Err(); err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err := s.uow.Within(ctx, func(r Repo) error {
		var err error
		if u, err = r.UserByID(ctx, id); err != nil {
			return err
		}
		now := s.now()
		if p.Status != nil {
			u.Status = *p.Status
			if u.Status == domain.StatusInactive {
				if err := r.RevokeUserRefreshTokens(ctx, id, now); err != nil {
					return err
				}
			}
		}
		if p.SpecialistApproved != nil {
			if u.Role != domain.RoleSpecialist {
				return &domain.ValidationError{Fields: []domain.FieldError{{Field: "specialist_approved", Code: "not_a_specialist"}}}
			}
			u.SpecialistApproved = *p.SpecialistApproved
		}
		u.UpdatedAt = now
		if err := r.UpdateUser(ctx, &u); err != nil {
			return err
		}
		return r.InsertAudit(ctx, s.event(ctx, a, nil, string(policy.UserManage), "user", &id, domain.OutcomeSuccess))
	})
	return u, err
}

// ListUsers lists accounts for the admin.
func (s *Service) ListUsers(ctx context.Context, a *Actor, limit int, after *uuid.UUID) ([]domain.User, error) {
	if err := s.guardGlobal(ctx, a, policy.UserManage, "user"); err != nil {
		return nil, err
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.uow.Repo().ListUsers(ctx, limit, after)
}
