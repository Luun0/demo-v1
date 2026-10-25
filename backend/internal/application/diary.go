package application

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/AknurOE/sen-support-platform/internal/domain"
	"github.com/AknurOE/sen-support-platform/internal/domain/policy"
)

// EntryInput creates or replaces a diary entry.
type EntryInput struct {
	Type       domain.EntryType
	Visibility domain.Visibility
	Content    string
	ObservedAt time.Time
}

// ListQuery filters a diary list.
type ListQuery struct {
	Query      string
	Visibility domain.Visibility
	Limit      int
	After      *Cursor
}

// EntryPage is a visibility-filtered page; Total counts visible rows only.
type EntryPage struct {
	Items []domain.DiaryEntry
	Total int
	Next  *Cursor
}

func validateEntry(in EntryInput) error {
	var v domain.Validator
	v.Check(in.Type.Valid(), "type", "invalid")
	v.Check(in.Visibility.Valid(), "visibility", "invalid")
	v.Check(domain.ValidText(in.Content, domain.MaxContentLen), "content", "invalid")
	v.Check(!in.ObservedAt.IsZero(), "observed_at", "required")
	return v.Err()
}

func versionOf(e domain.DiaryEntry, n int, by uuid.UUID, at time.Time) *domain.DiaryEntryVersion {
	return &domain.DiaryEntryVersion{ID: newID(), EntryID: e.ID, VersionNumber: n, ChangedBy: by, Type: e.Type,
		Visibility: e.Visibility, Content: e.Content, ObservedAt: e.ObservedAt, ChangedAt: at}
}

// ListEntries lists entries; parents get shared entries only, filtered in SQL.
func (s *Service) ListEntries(ctx context.Context, a *Actor, child uuid.UUID, q ListQuery) (EntryPage, error) {
	if _, err := s.guard(ctx, a, child, policy.DiaryList, "diary_entry", nil, nil); err != nil {
		return EntryPage{}, err
	}
	if q.Limit <= 0 || q.Limit > 100 {
		q.Limit = 50
	}
	f := DiaryFilter{ChildID: child, SharedOnly: a.Role != domain.RoleSpecialist, Visibility: q.Visibility,
		Query: q.Query, Limit: q.Limit + 1, After: q.After}
	r := s.uow.Repo()
	items, err := r.ListEntries(ctx, f)
	if err != nil {
		return EntryPage{}, err
	}
	f.After, f.Limit = nil, 0
	total, err := r.CountEntries(ctx, f)
	if err != nil {
		return EntryPage{}, err
	}
	page := EntryPage{Items: items, Total: total}
	if len(items) > q.Limit {
		page.Items = items[:q.Limit]
		last := page.Items[q.Limit-1]
		page.Next = &Cursor{At: last.ObservedAt, ID: last.ID}
	}
	return page, s.readAudit(ctx, a, &child, policy.DiaryList, "diary_entry", nil)
}

// CreateEntry writes an entry plus version 1 and an audit event atomically.
func (s *Service) CreateEntry(ctx context.Context, a *Actor, child uuid.UUID, in EntryInput) (domain.DiaryEntry, error) {
	if err := validateEntry(in); err != nil {
		if _, gerr := s.guard(ctx, a, child, policy.DiaryList, "diary_entry", nil, nil); gerr != nil {
			return domain.DiaryEntry{}, gerr // never reveal validation details without access
		}
		return domain.DiaryEntry{}, err
	}
	if _, err := s.guard(ctx, a, child, policy.DiaryCreate, "diary_entry", nil, func(r *policy.Resource) error {
		r.Visibility = in.Visibility
		return nil
	}); err != nil {
		return domain.DiaryEntry{}, err
	}
	now := s.now()
	e := domain.DiaryEntry{ID: newID(), ChildID: child, AuthorID: a.ID, Type: in.Type, Visibility: in.Visibility,
		Content: in.Content, ObservedAt: in.ObservedAt.UTC(), CreatedAt: now, UpdatedAt: now}
	err := s.uow.Within(ctx, func(r Repo) error {
		if err := r.CreateEntry(ctx, &e); err != nil {
			return err
		}
		if err := r.InsertVersion(ctx, versionOf(e, 1, a.ID, now)); err != nil {
			return err
		}
		return r.InsertAudit(ctx, s.event(ctx, a, &child, string(policy.DiaryCreate), "diary_entry", &e.ID, domain.OutcomeSuccess))
	})
	return e, err
}

// EditEntry replaces an own entry; a new version is stored (AUD-03).
func (s *Service) EditEntry(ctx context.Context, a *Actor, child, id uuid.UUID, in EntryInput) (domain.DiaryEntry, error) {
	var e domain.DiaryEntry
	_, err := s.guard(ctx, a, child, policy.DiaryEdit, "diary_entry", &id, func(res *policy.Resource) error {
		var err error
		e, err = s.uow.Repo().EntryByID(ctx, child, id, false)
		res.Author, res.Visibility = err == nil && e.AuthorID == a.ID, e.Visibility
		return err
	})
	if err != nil {
		return domain.DiaryEntry{}, err
	}
	if err := validateEntry(in); err != nil {
		return domain.DiaryEntry{}, err
	}
	if a.Role == domain.RoleParent && !domain.ParentVisibilityAllowed(in.Visibility) {
		s.auditOutside(ctx, a, &child, string(policy.DiaryEdit), "diary_entry", &id, domain.OutcomeDenied)
		return domain.DiaryEntry{}, domain.ErrForbidden
	}
	now := s.now()
	e.Type, e.Visibility, e.Content, e.ObservedAt, e.UpdatedAt = in.Type, in.Visibility, in.Content, in.ObservedAt.UTC(), now
	err = s.uow.Within(ctx, func(r Repo) error {
		n, err := r.NextVersionNumber(ctx, id)
		if err != nil {
			return err
		}
		if err := r.UpdateEntry(ctx, &e); err != nil {
			return err
		}
		if err := r.InsertVersion(ctx, versionOf(e, n, a.ID, now)); err != nil {
			return err
		}
		return r.InsertAudit(ctx, s.event(ctx, a, &child, string(policy.DiaryEdit), "diary_entry", &id, domain.OutcomeSuccess))
	})
	return e, err
}

// ListVersions returns the edit history of a visible entry.
func (s *Service) ListVersions(ctx context.Context, a *Actor, child, id uuid.UUID) ([]domain.DiaryEntryVersion, error) {
	_, err := s.guard(ctx, a, child, policy.DiaryVersions, "diary_entry", &id, func(res *policy.Resource) error {
		e, err := s.uow.Repo().EntryByID(ctx, child, id, false)
		res.Visibility = e.Visibility
		return err
	})
	if err != nil {
		return nil, err
	}
	vs, err := s.uow.Repo().ListVersions(ctx, id)
	if err != nil {
		return nil, err
	}
	if a.Role == domain.RoleParent { // defence in depth: drop non-shared history
		out := vs[:0]
		for _, v := range vs {
			if v.Visibility == domain.VisibilityShared {
				out = append(out, v)
			}
		}
		vs = out
	}
	return vs, s.readAudit(ctx, a, &child, policy.DiaryVersions, "diary_entry", &id)
}
