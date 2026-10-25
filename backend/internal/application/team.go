package application

import (
	"context"

	"github.com/google/uuid"

	"github.com/AknurOE/sen-support-platform/internal/domain"
	"github.com/AknurOE/sen-support-platform/internal/domain/policy"
)

// Team lists guardians and specialist assignments of a child.
func (s *Service) Team(ctx context.Context, a *Actor, child uuid.UUID) ([]TeamMember, error) {
	if _, err := s.guard(ctx, a, child, policy.ChildView, "child_team", &child, nil); err != nil {
		return nil, err
	}
	out, err := s.uow.Repo().ListTeam(ctx, child)
	if err != nil {
		return nil, err
	}
	return out, s.readAudit(ctx, a, &child, policy.ChildView, "child_team", &child)
}

// Specialists is the directory of approved, active specialists (names only).
func (s *Service) Specialists(ctx context.Context, a *Actor) ([]domain.User, error) {
	if err := s.guardGlobal(ctx, a, policy.SpecialistDir, "specialist"); err != nil {
		return nil, err
	}
	return s.uow.Repo().ListApprovedSpecialists(ctx)
}
