package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/AknurOE/sen-support-platform/internal/domain"
)

//nolint:gosec // G101: SQL column names only; this constant contains no credentials.
const passwordResetCols = `id, user_id, token_hash, expires_at, used_at, revoked_at, created_at`

func scanPasswordReset(row interface{ Scan(...any) error }) (domain.PasswordResetToken, error) {
	var t domain.PasswordResetToken
	err := row.Scan(&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.UsedAt, &t.RevokedAt, &t.CreatedAt)
	return t, mapErr(err)
}

func (r *Repo) InsertPasswordResetToken(ctx context.Context, t *domain.PasswordResetToken) error {
	return r.exec(ctx, `INSERT INTO password_reset_tokens (`+passwordResetCols+`) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		t.ID, t.UserID, t.TokenHash, t.ExpiresAt, t.UsedAt, t.RevokedAt, t.CreatedAt)
}

func (r *Repo) PasswordResetTokenByHash(ctx context.Context, hash []byte) (domain.PasswordResetToken, error) {
	return scanPasswordReset(r.q.QueryRow(ctx, `SELECT `+passwordResetCols+` FROM password_reset_tokens WHERE token_hash=$1 FOR UPDATE`, hash))
}

func (r *Repo) MarkPasswordResetUsed(ctx context.Context, id uuid.UUID, at time.Time) error {
	return r.execOne(ctx, `UPDATE password_reset_tokens SET used_at=$2 WHERE id=$1 AND used_at IS NULL AND revoked_at IS NULL`, id, at)
}

func (r *Repo) RevokePasswordResetTokens(ctx context.Context, user uuid.UUID, at time.Time) error {
	return r.exec(ctx, `UPDATE password_reset_tokens SET revoked_at=$2 WHERE user_id=$1 AND used_at IS NULL AND revoked_at IS NULL`, user, at)
}
