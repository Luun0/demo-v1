package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/AknurOE/sen-support-platform/internal/domain"
)

const userCols = `id, email, password_hash, first_name, last_name, role, status, specialist_approved, failed_login_attempts, locked_until, created_at, updated_at`

func scanUser(row pgx.Row) (domain.User, error) {
	var u domain.User
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role, &u.Status,
		&u.SpecialistApproved, &u.FailedLoginAttempts, &u.LockedUntil, &u.CreatedAt, &u.UpdatedAt)
	return u, mapErr(err)
}

// CreateUser inserts a user; duplicate email yields ErrConflict.
func (r *Repo) CreateUser(ctx context.Context, u *domain.User) error {
	return r.exec(ctx, `INSERT INTO users (`+userCols+`) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		u.ID, u.Email, u.PasswordHash, u.FirstName, u.LastName, u.Role, u.Status, u.SpecialistApproved,
		u.FailedLoginAttempts, u.LockedUntil, u.CreatedAt, u.UpdatedAt)
}

// UserByID loads a user.
func (r *Repo) UserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return scanUser(r.q.QueryRow(ctx, `SELECT `+userCols+` FROM users WHERE id=$1`, id))
}

// UserByEmail loads a user by normalised email.
func (r *Repo) UserByEmail(ctx context.Context, email string) (domain.User, error) {
	return scanUser(r.q.QueryRow(ctx, `SELECT `+userCols+` FROM users WHERE email=$1`, email))
}

// UpdateUser persists mutable fields.
func (r *Repo) UpdateUser(ctx context.Context, u *domain.User) error {
	return r.execOne(ctx, `UPDATE users SET password_hash=$2, first_name=$3, last_name=$4, status=$5, specialist_approved=$6,
		failed_login_attempts=$7, locked_until=$8, updated_at=$9 WHERE id=$1`,
		u.ID, u.PasswordHash, u.FirstName, u.LastName, u.Status, u.SpecialistApproved, u.FailedLoginAttempts, u.LockedUntil, u.UpdatedAt)
}

// ListUsers pages users by id.
func (r *Repo) ListUsers(ctx context.Context, limit int, after *uuid.UUID) ([]domain.User, error) {
	rows, err := r.q.Query(ctx, `SELECT `+userCols+` FROM users WHERE ($1::uuid IS NULL OR id > $1) ORDER BY id LIMIT $2`, after, limit)
	return collect(rows, err, scanUser)
}

const rtCols = `id, user_id, family_id, token_hash, expires_at, rotated_at, revoked_at, created_at`

// InsertRefreshToken stores a refresh token hash.
func (r *Repo) InsertRefreshToken(ctx context.Context, t *domain.RefreshToken) error {
	return r.exec(ctx, `INSERT INTO refresh_tokens (`+rtCols+`) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		t.ID, t.UserID, t.FamilyID, t.TokenHash, t.ExpiresAt, t.RotatedAt, t.RevokedAt, t.CreatedAt)
}

// RefreshTokenByHash locks and loads a token row.
func (r *Repo) RefreshTokenByHash(ctx context.Context, hash []byte) (domain.RefreshToken, error) {
	var t domain.RefreshToken
	err := r.q.QueryRow(ctx, `SELECT `+rtCols+` FROM refresh_tokens WHERE token_hash=$1 FOR UPDATE`, hash).
		Scan(&t.ID, &t.UserID, &t.FamilyID, &t.TokenHash, &t.ExpiresAt, &t.RotatedAt, &t.RevokedAt, &t.CreatedAt)
	return t, mapErr(err)
}

// MarkRefreshRotated marks a token as used.
func (r *Repo) MarkRefreshRotated(ctx context.Context, id uuid.UUID, at time.Time) error {
	return r.execOne(ctx, `UPDATE refresh_tokens SET rotated_at=$2 WHERE id=$1 AND rotated_at IS NULL`, id, at)
}

// RevokeRefreshFamily revokes every token in a family.
func (r *Repo) RevokeRefreshFamily(ctx context.Context, family uuid.UUID, at time.Time) error {
	return r.exec(ctx, `UPDATE refresh_tokens SET revoked_at=$2 WHERE family_id=$1 AND revoked_at IS NULL`, family, at)
}

// RevokeUserRefreshTokens revokes all of a user's tokens.
func (r *Repo) RevokeUserRefreshTokens(ctx context.Context, user uuid.UUID, at time.Time) error {
	return r.exec(ctx, `UPDATE refresh_tokens SET revoked_at=$2 WHERE user_id=$1 AND revoked_at IS NULL`, user, at)
}
