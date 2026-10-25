package application

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/AknurOE/sen-support-platform/internal/domain"
)

// RequestPasswordReset always returns nil for unknown or inactive accounts.
// The caller sends a link only when raw is non-empty, but returns one generic
// response to the client in every case (SEC-04).
func (s *Service) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	email = domain.NormalizeEmail(email)
	if !domain.ValidEmail(email) {
		return "", &domain.ValidationError{Fields: []domain.FieldError{{Field: "email", Code: "invalid"}}}
	}
	raw, hash, err := newPasswordResetToken()
	if err != nil {
		return "", err
	}
	u, err := s.uow.Repo().UserByEmail(ctx, email)
	if errors.Is(err, domain.ErrNotFound) || err == nil && u.Status != domain.StatusActive {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	ttl := s.cfg.PasswordResetTTL
	if ttl <= 0 {
		ttl = time.Hour
	}
	t := &domain.PasswordResetToken{ID: newID(), UserID: u.ID, TokenHash: hash,
		ExpiresAt: s.now().Add(ttl), CreatedAt: s.now()}
	if err := s.uow.Within(ctx, func(r Repo) error {
		if err := r.RevokePasswordResetTokens(ctx, u.ID, s.now()); err != nil {
			return err
		}
		return r.InsertPasswordResetToken(ctx, t)
	}); err != nil {
		return "", err
	}
	return raw, nil
}

// ConfirmPasswordReset consumes a single unexpired token and revokes refresh
// sessions in the same transaction as the password change.
func (s *Service) ConfirmPasswordReset(ctx context.Context, raw, password string) error {
	if raw == "" {
		return domain.ErrUnauthorized
	}
	if ok, code := domain.ValidPassword(password); !ok {
		return &domain.ValidationError{Fields: []domain.FieldError{{Field: "password", Code: code}}}
	}
	hash, err := s.hasher.Hash(password)
	if err != nil {
		return err
	}
	now := s.now()
	return s.uow.Within(ctx, func(r Repo) error {
		t, err := r.PasswordResetTokenByHash(ctx, HashToken(raw))
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrUnauthorized
		}
		if err != nil {
			return err
		}
		if t.UsedAt != nil || t.RevokedAt != nil || !t.ExpiresAt.After(now) {
			return domain.ErrUnauthorized
		}
		u, err := r.UserByID(ctx, t.UserID)
		if err != nil {
			return err
		}
		u.PasswordHash, u.FailedLoginAttempts, u.LockedUntil, u.UpdatedAt = hash, 0, nil, now
		if err := r.UpdateUser(ctx, &u); err != nil {
			return err
		}
		if err := r.MarkPasswordResetUsed(ctx, t.ID, now); err != nil {
			return err
		}
		if err := r.RevokePasswordResetTokens(ctx, u.ID, now); err != nil {
			return err
		}
		if err := r.RevokeUserRefreshTokens(ctx, u.ID, now); err != nil {
			return err
		}
		return r.InsertAudit(ctx, s.event(ctx, &Actor{ID: u.ID, Role: u.Role}, nil,
			"auth:password_reset", "user", &u.ID, domain.OutcomeSuccess))
	})
}

func newPasswordResetToken() (string, []byte, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", nil, err
	}
	raw := base64.RawURLEncoding.EncodeToString(b)
	return raw, HashToken(raw), nil
}
