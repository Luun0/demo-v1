package postgres

import (
	"context"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/AknurOE/sen-support-platform/internal/application"
	"github.com/AknurOE/sen-support-platform/internal/domain"
)

const entryCols = `id, child_id, author_id, type, visibility, content, observed_at, created_at, updated_at`

func scanEntry(row pgx.Row) (domain.DiaryEntry, error) {
	var e domain.DiaryEntry
	err := row.Scan(&e.ID, &e.ChildID, &e.AuthorID, &e.Type, &e.Visibility, &e.Content, &e.ObservedAt, &e.CreatedAt, &e.UpdatedAt)
	return e, mapErr(err)
}

// CreateEntry inserts an entry.
func (r *Repo) CreateEntry(ctx context.Context, e *domain.DiaryEntry) error {
	return r.exec(ctx, `INSERT INTO diary_entries (`+entryCols+`) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		e.ID, e.ChildID, e.AuthorID, e.Type, e.Visibility, e.Content, e.ObservedAt, e.CreatedAt, e.UpdatedAt)
}

// EntryByID loads an entry scoped to the child; sharedOnly filters in SQL.
func (r *Repo) EntryByID(ctx context.Context, child, id uuid.UUID, sharedOnly bool) (domain.DiaryEntry, error) {
	return scanEntry(r.q.QueryRow(ctx, `SELECT `+entryCols+` FROM diary_entries
		WHERE child_id=$1 AND id=$2 AND (NOT $3 OR visibility='shared')`, child, id, sharedOnly))
}

// UpdateEntry persists an entry.
func (r *Repo) UpdateEntry(ctx context.Context, e *domain.DiaryEntry) error {
	return r.execOne(ctx, `UPDATE diary_entries SET type=$3, visibility=$4, content=$5, observed_at=$6, updated_at=$7
		WHERE child_id=$1 AND id=$2`, e.ChildID, e.ID, e.Type, e.Visibility, e.Content, e.ObservedAt, e.UpdatedAt)
}

// InsertVersion stores an immutable version.
func (r *Repo) InsertVersion(ctx context.Context, v *domain.DiaryEntryVersion) error {
	return r.exec(ctx, `INSERT INTO diary_entry_versions (id, entry_id, version_number, changed_by, type, visibility, content, observed_at, changed_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`, v.ID, v.EntryID, v.VersionNumber, v.ChangedBy, v.Type, v.Visibility, v.Content, v.ObservedAt, v.ChangedAt)
}

// NextVersionNumber returns max+1 for the entry.
func (r *Repo) NextVersionNumber(ctx context.Context, entry uuid.UUID) (int, error) {
	var n int
	err := r.q.QueryRow(ctx, `SELECT COALESCE(MAX(version_number),0)+1 FROM diary_entry_versions WHERE entry_id=$1`, entry).Scan(&n)
	return n, mapErr(err)
}

// ListVersions lists versions ascending.
func (r *Repo) ListVersions(ctx context.Context, entry uuid.UUID) ([]domain.DiaryEntryVersion, error) {
	rows, err := r.q.Query(ctx, `SELECT id, entry_id, version_number, changed_by, type, visibility, content, observed_at, changed_at
		FROM diary_entry_versions WHERE entry_id=$1 ORDER BY version_number`, entry)
	return collect(rows, err, func(row pgx.Row) (domain.DiaryEntryVersion, error) {
		var v domain.DiaryEntryVersion
		err := row.Scan(&v.ID, &v.EntryID, &v.VersionNumber, &v.ChangedBy, &v.Type, &v.Visibility, &v.Content, &v.ObservedAt, &v.ChangedAt)
		return v, err
	})
}

// diaryWhere is the single visibility-aware predicate used by list, count,
// search and export ($1 child, $2 sharedOnly, $3 visibility, $4 query).
const diaryWhere = ` WHERE child_id=$1 AND (NOT $2 OR visibility='shared')
	AND ($3 = '' OR visibility=$3) AND ($4 = '' OR content ILIKE '%' || $4 || '%' ESCAPE '\')`

func likeEscape(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
}

func diaryArgs(f application.DiaryFilter) []any {
	return []any{f.ChildID, f.SharedOnly, string(f.Visibility), likeEscape(f.Query)}
}

// ListEntries returns visible entries, newest first, keyset paginated.
func (r *Repo) ListEntries(ctx context.Context, f application.DiaryFilter) ([]domain.DiaryEntry, error) {
	args := diaryArgs(f)
	q := `SELECT ` + entryCols + ` FROM diary_entries` + diaryWhere
	if f.After != nil {
		q += ` AND (observed_at, id) < ($5, $6)`
		args = append(args, f.After.At, f.After.ID)
	}
	q += ` ORDER BY observed_at DESC, id DESC`
	if f.Limit > 0 {
		args = append(args, f.Limit)
		q += ` LIMIT $` + strconv.Itoa(len(args))
	}
	rows, err := r.q.Query(ctx, q, args...)
	return collect(rows, err, scanEntry)
}

// CountEntries counts visible entries with the same predicate.
func (r *Repo) CountEntries(ctx context.Context, f application.DiaryFilter) (int, error) {
	var n int
	err := r.q.QueryRow(ctx, `SELECT count(*) FROM diary_entries`+diaryWhere, diaryArgs(f)...).Scan(&n)
	return n, mapErr(err)
}

const auditCols = `id, actor_id, child_id, action, resource, resource_id, outcome, request_id, occurred_at`

// InsertAudit stores an audit event.
func (r *Repo) InsertAudit(ctx context.Context, e *domain.AuditEvent) error {
	return r.exec(ctx, `INSERT INTO audit_events (`+auditCols+`) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		e.ID, e.ActorID, e.ChildID, e.Action, e.Resource, e.ResourceID, e.Outcome, e.RequestID, e.OccurredAt)
}

// ListAudit lists events newest first, optionally for one child.
func (r *Repo) ListAudit(ctx context.Context, child *uuid.UUID, limit int, after *application.Cursor) ([]domain.AuditEvent, error) {
	var at, id any
	if after != nil {
		at, id = after.At, after.ID
	}
	rows, err := r.q.Query(ctx, `SELECT `+auditCols+` FROM audit_events
		WHERE ($1::uuid IS NULL OR child_id=$1) AND ($2::timestamptz IS NULL OR (occurred_at, id) < ($2, $3::uuid))
		ORDER BY occurred_at DESC, id DESC LIMIT $4`, child, at, id, limit)
	return collect(rows, err, func(row pgx.Row) (domain.AuditEvent, error) {
		var e domain.AuditEvent
		err := row.Scan(&e.ID, &e.ActorID, &e.ChildID, &e.Action, &e.Resource, &e.ResourceID, &e.Outcome, &e.RequestID, &e.OccurredAt)
		return e, err
	})
}
