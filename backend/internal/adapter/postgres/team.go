package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/AknurOE/sen-support-platform/internal/application"
	"github.com/AknurOE/sen-support-platform/internal/domain"
)

// ListTeam returns guardians first, then specialist assignments.
func (r *Repo) ListTeam(ctx context.Context, child uuid.UUID) ([]application.TeamMember, error) {
	rows, err := r.q.Query(ctx, `
		SELECT u.id, u.first_name, u.last_name, u.role, g.is_primary, g.relationship, NULL::timestamptz, NULL::timestamptz
		  FROM child_guardians g JOIN users u ON u.id = g.guardian_id WHERE g.child_id = $1
		UNION ALL
		SELECT u.id, u.first_name, u.last_name, u.role, false, '', s.assigned_at, s.ended_at
		  FROM child_specialists s JOIN users u ON u.id = s.specialist_id WHERE s.child_id = $1
		ORDER BY 4, 2`, child)
	return collect(rows, err, func(row pgx.Row) (application.TeamMember, error) {
		var m application.TeamMember
		err := row.Scan(&m.UserID, &m.FirstName, &m.LastName, &m.Role, &m.IsPrimary, &m.Relationship, &m.AssignedAt, &m.EndedAt)
		return m, err
	})
}

// ListApprovedSpecialists lists active, approved specialists.
func (r *Repo) ListApprovedSpecialists(ctx context.Context) ([]domain.User, error) {
	rows, err := r.q.Query(ctx, `SELECT `+userCols+` FROM users WHERE role='specialist' AND status='active' AND specialist_approved ORDER BY last_name, first_name, id`)
	return collect(rows, err, scanUser)
}
