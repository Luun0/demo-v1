package application

import (
	"context"

	"github.com/google/uuid"

	"github.com/AknurOE/sen-support-platform/internal/domain"
	"github.com/AknurOE/sen-support-platform/internal/domain/policy"
)

// AuditPage is a page of audit events.
type AuditPage struct {
	Items []domain.AuditEvent
	Next  *Cursor
}

func page(items []domain.AuditEvent, limit int) AuditPage {
	p := AuditPage{Items: items}
	if len(items) > limit {
		p.Items = items[:limit]
		last := p.Items[limit-1]
		p.Next = &Cursor{At: last.OccurredAt, ID: last.ID}
	}
	return p
}

func clamp(limit int) int {
	if limit <= 0 || limit > 100 {
		return 50
	}
	return limit
}

// ChildAudit lists access events for a child (related parent or admin).
func (s *Service) ChildAudit(ctx context.Context, a *Actor, child uuid.UUID, limit int, after *Cursor) (AuditPage, error) {
	if _, err := s.guard(ctx, a, child, policy.ChildAudit, "audit_event", nil, nil); err != nil {
		return AuditPage{}, err
	}
	limit = clamp(limit)
	items, err := s.uow.Repo().ListAudit(ctx, &child, limit+1, after)
	if err != nil {
		return AuditPage{}, err
	}
	return page(items, limit), s.readAudit(ctx, a, &child, policy.ChildAudit, "audit_event", nil)
}

// SystemAudit lists all audit events (admin).
func (s *Service) SystemAudit(ctx context.Context, a *Actor, limit int, after *Cursor) (AuditPage, error) {
	if err := s.guardGlobal(ctx, a, policy.SystemAudit, "audit_event"); err != nil {
		return AuditPage{}, err
	}
	limit = clamp(limit)
	items, err := s.uow.Repo().ListAudit(ctx, nil, limit+1, after)
	if err != nil {
		return AuditPage{}, err
	}
	return page(items, limit), s.uow.Repo().InsertAudit(ctx, s.event(ctx, a, nil, string(policy.SystemAudit), "audit_event", nil, domain.OutcomeSuccess))
}
