package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/AknurOE/sen-support-platform/internal/domain"
)

const childCols = `c.id, c.first_name, c.last_name, c.birth_date, c.created_at, c.updated_at`

func scanChild(row pgx.Row) (domain.Child, error) {
	var c domain.Child
	err := row.Scan(&c.ID, &c.FirstName, &c.LastName, &c.BirthDate, &c.CreatedAt, &c.UpdatedAt)
	return c, mapErr(err)
}

// CreateChild inserts a child.
func (r *Repo) CreateChild(ctx context.Context, c *domain.Child) error {
	return r.exec(ctx, `INSERT INTO children (id, first_name, last_name, birth_date, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6)`,
		c.ID, c.FirstName, c.LastName, c.BirthDate, c.CreatedAt, c.UpdatedAt)
}

// ChildByID loads a child.
func (r *Repo) ChildByID(ctx context.Context, id uuid.UUID) (domain.Child, error) {
	return scanChild(r.q.QueryRow(ctx, `SELECT `+childCols+` FROM children c WHERE c.id=$1`, id))
}

// UpdateChild persists profile fields.
func (r *Repo) UpdateChild(ctx context.Context, c *domain.Child) error {
	return r.execOne(ctx, `UPDATE children SET first_name=$2, last_name=$3, birth_date=$4, updated_at=$5 WHERE id=$1`,
		c.ID, c.FirstName, c.LastName, c.BirthDate, c.UpdatedAt)
}

// DeleteChild deletes a child and cascades all child data.
func (r *Repo) DeleteChild(ctx context.Context, id uuid.UUID) error {
	return r.execOne(ctx, `DELETE FROM children WHERE id=$1`, id)
}

// ListChildrenForGuardian lists linked children.
func (r *Repo) ListChildrenForGuardian(ctx context.Context, user uuid.UUID) ([]domain.Child, error) {
	rows, err := r.q.Query(ctx, `SELECT `+childCols+` FROM children c JOIN child_guardians g ON g.child_id=c.id
		WHERE g.guardian_id=$1 ORDER BY c.created_at, c.id`, user)
	return collect(rows, err, scanChild)
}

// ListChildrenForSpecialist lists children with an active assignment.
func (r *Repo) ListChildrenForSpecialist(ctx context.Context, user uuid.UUID, now time.Time) ([]domain.Child, error) {
	rows, err := r.q.Query(ctx, `SELECT `+childCols+` FROM children c JOIN child_specialists s ON s.child_id=c.id
		WHERE s.specialist_id=$1 AND s.assigned_at <= $2 AND (s.ended_at IS NULL OR s.ended_at > $2) ORDER BY c.created_at, c.id`, user, now)
	return collect(rows, err, scanChild)
}

// ListAllChildren lists every child (admin metadata view).
func (r *Repo) ListAllChildren(ctx context.Context) ([]domain.Child, error) {
	rows, err := r.q.Query(ctx, `SELECT `+childCols+` FROM children c ORDER BY c.created_at, c.id`)
	return collect(rows, err, scanChild)
}

// AddGuardian links a guardian.
func (r *Repo) AddGuardian(ctx context.Context, g domain.ChildGuardian) error {
	return r.exec(ctx, `INSERT INTO child_guardians (child_id, guardian_id, relationship, is_primary, created_at) VALUES ($1,$2,$3,$4,$5)`,
		g.ChildID, g.GuardianID, g.Relationship, g.IsPrimary, g.CreatedAt)
}

// Guardian loads a guardian link.
func (r *Repo) Guardian(ctx context.Context, child, user uuid.UUID) (domain.ChildGuardian, error) {
	var g domain.ChildGuardian
	err := r.q.QueryRow(ctx, `SELECT child_id, guardian_id, relationship, is_primary, created_at FROM child_guardians
		WHERE child_id=$1 AND guardian_id=$2`, child, user).Scan(&g.ChildID, &g.GuardianID, &g.Relationship, &g.IsPrimary, &g.CreatedAt)
	return g, mapErr(err)
}

// AddConsent records consent.
func (r *Repo) AddConsent(ctx context.Context, c domain.ConsentRecord) error {
	return r.exec(ctx, `INSERT INTO consent_records (id, child_id, guardian_id, policy_version, granted_at) VALUES ($1,$2,$3,$4,$5)`,
		c.ID, c.ChildID, c.GuardianID, c.PolicyVersion, c.GrantedAt)
}

// UpsertAssignment creates or restarts an assignment.
func (r *Repo) UpsertAssignment(ctx context.Context, a domain.ChildSpecialist) error {
	return r.exec(ctx, `INSERT INTO child_specialists (child_id, specialist_id, assigned_by, assigned_at, ended_at) VALUES ($1,$2,$3,$4,NULL)
		ON CONFLICT (child_id, specialist_id) DO UPDATE SET assigned_by=EXCLUDED.assigned_by, assigned_at=EXCLUDED.assigned_at, ended_at=NULL`,
		a.ChildID, a.SpecialistID, a.AssignedBy, a.AssignedAt)
}

// Assignment loads an assignment (active or ended).
func (r *Repo) Assignment(ctx context.Context, child, spec uuid.UUID) (domain.ChildSpecialist, error) {
	var a domain.ChildSpecialist
	err := r.q.QueryRow(ctx, `SELECT child_id, specialist_id, assigned_by, assigned_at, ended_at FROM child_specialists
		WHERE child_id=$1 AND specialist_id=$2`, child, spec).Scan(&a.ChildID, &a.SpecialistID, &a.AssignedBy, &a.AssignedAt, &a.EndedAt)
	return a, mapErr(err)
}

// EndAssignment sets ended_at.
func (r *Repo) EndAssignment(ctx context.Context, child, spec uuid.UUID, at time.Time) error {
	return r.execOne(ctx, `UPDATE child_specialists SET ended_at=$3 WHERE child_id=$1 AND specialist_id=$2`, child, spec, at)
}

const sessCols = `id, child_id, specialist_id, created_by, title, starts_at, ends_at, mode, status, created_at, updated_at`

func scanSession(row pgx.Row) (domain.SupportSession, error) {
	var s domain.SupportSession
	err := row.Scan(&s.ID, &s.ChildID, &s.SpecialistID, &s.CreatedBy, &s.Title, &s.StartsAt, &s.EndsAt, &s.Mode, &s.Status, &s.CreatedAt, &s.UpdatedAt)
	return s, mapErr(err)
}

// CreateSession inserts a session.
func (r *Repo) CreateSession(ctx context.Context, s *domain.SupportSession) error {
	return r.exec(ctx, `INSERT INTO support_sessions (`+sessCols+`) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		s.ID, s.ChildID, s.SpecialistID, s.CreatedBy, s.Title, s.StartsAt, s.EndsAt, s.Mode, s.Status, s.CreatedAt, s.UpdatedAt)
}

// SessionByID loads a session scoped to the child (no cross-child IDOR).
func (r *Repo) SessionByID(ctx context.Context, child, id uuid.UUID) (domain.SupportSession, error) {
	return scanSession(r.q.QueryRow(ctx, `SELECT `+sessCols+` FROM support_sessions WHERE child_id=$1 AND id=$2`, child, id))
}

// UpdateSession persists a session.
func (r *Repo) UpdateSession(ctx context.Context, s *domain.SupportSession) error {
	return r.execOne(ctx, `UPDATE support_sessions SET title=$3, starts_at=$4, ends_at=$5, mode=$6, status=$7, updated_at=$8
		WHERE child_id=$1 AND id=$2`, s.ChildID, s.ID, s.Title, s.StartsAt, s.EndsAt, s.Mode, s.Status, s.UpdatedAt)
}

// ListSessions lists a child's sessions by start time.
func (r *Repo) ListSessions(ctx context.Context, child uuid.UUID) ([]domain.SupportSession, error) {
	rows, err := r.q.Query(ctx, `SELECT `+sessCols+` FROM support_sessions WHERE child_id=$1 ORDER BY starts_at, id`, child)
	return collect(rows, err, scanSession)
}
