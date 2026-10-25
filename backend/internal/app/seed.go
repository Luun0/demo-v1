package app

import (
	"context"
	"errors"
	"time"

	"github.com/AknurOE/sen-support-platform/internal/application"
	"github.com/AknurOE/sen-support-platform/internal/domain"
)

// SeedPassword is the shared password of all synthetic demo accounts.
const SeedPassword = "Synthetic-Demo-2026"

// Seed creates SYNTHETIC data only (PRV-08): one admin, two specialists,
// two families with one child each. It is idempotent per email.
func Seed(ctx context.Context, svc *application.Service) error {
	admin, err := svc.BootstrapAdmin(ctx, "admin@example.test", SeedPassword, "Demo", "Admin")
	if errors.Is(err, domain.ErrConflict) {
		return nil // already seeded
	} else if err != nil {
		return err
	}
	adm := &application.Actor{ID: admin.ID, Role: domain.RoleAdmin}
	mk := func(email, first string, role domain.Role) (*application.Actor, error) {
		u, err := svc.CreateUserByAdmin(ctx, adm, application.RegisterInput{Email: email, Password: SeedPassword, FirstName: first, LastName: "Demo", Role: role}, role == domain.RoleSpecialist)
		return &application.Actor{ID: u.ID, Role: u.Role}, err
	}
	spec1, err := mk("specialist1@example.test", "Dana", domain.RoleSpecialist)
	if err != nil {
		return err
	}
	if _, err := mk("specialist2@example.test", "Erlan", domain.RoleSpecialist); err != nil {
		return err
	}
	now := time.Now().UTC()
	for i, email := range []string{"parent1@example.test", "parent2@example.test"} {
		p, err := mk(email, []string{"Aida", "Bota"}[i], domain.RoleParent)
		if err != nil {
			return err
		}
		c, err := svc.CreateChild(ctx, p, application.ChildInput{FirstName: []string{"Mira", "Timur"}[i], LastName: "Demo",
			BirthDate: now.AddDate(-7-i, 0, 0).Truncate(24 * time.Hour), PolicyVersion: "2026-10"})
		if err != nil {
			return err
		}
		if i == 0 {
			if _, err := svc.AssignSpecialist(ctx, p, c.ID, spec1.ID); err != nil {
				return err
			}
			if _, err := svc.CreateEntry(ctx, p, c.ID, application.EntryInput{Type: domain.EntryObservation, Visibility: domain.VisibilityShared, Content: "Synthetic note: enjoyed the reading game.", ObservedAt: now.Add(-48 * time.Hour)}); err != nil {
				return err
			}
			if _, err := svc.CreateEntry(ctx, spec1, c.ID, application.EntryInput{Type: domain.EntryProgress, Visibility: domain.VisibilitySpecialistsOnly, Content: "Synthetic specialist-only note.", ObservedAt: now.Add(-24 * time.Hour)}); err != nil {
				return err
			}
			if _, err := svc.CreateSession(ctx, p, c.ID, application.SessionInput{SpecialistID: spec1.ID, Title: "Speech session (synthetic)", StartsAt: now.Add(72 * time.Hour), EndsAt: now.Add(73 * time.Hour), Mode: domain.ModeOnline, Status: domain.SessionScheduled}); err != nil {
				return err
			}
		}
	}
	return nil
}
