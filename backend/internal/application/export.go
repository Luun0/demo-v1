package application

import (
	"context"

	"github.com/google/uuid"

	"github.com/AknurOE/sen-support-platform/internal/domain"
	"github.com/AknurOE/sen-support-platform/internal/domain/policy"
)

// Export is the guardian export: profile, sessions and SHARED entries only.
type Export struct {
	Child    domain.Child
	Sessions []domain.SupportSession
	Entries  []domain.DiaryEntry
}

// ExportChild exports permitted records (primary guardian, PRV-05).
func (s *Service) ExportChild(ctx context.Context, a *Actor, child uuid.UUID) (Export, error) {
	if _, err := s.guard(ctx, a, child, policy.ChildExport, "child_export", &child, nil); err != nil {
		return Export{}, err
	}
	r := s.uow.Repo()
	var out Export
	var err error
	if out.Child, err = r.ChildByID(ctx, child); err != nil {
		return out, err
	}
	if out.Sessions, err = r.ListSessions(ctx, child); err != nil {
		return out, err
	}
	// SharedOnly is enforced in SQL; the exporter is always a parent.
	if out.Entries, err = r.ListEntries(ctx, DiaryFilter{ChildID: child, SharedOnly: true, Limit: 100000}); err != nil {
		return out, err
	}
	return out, s.readAudit(ctx, a, &child, policy.ChildExport, "child_export", &child)
}
