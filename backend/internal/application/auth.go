package application

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/AknurOE/sen-support-platform/internal/domain"
)

// RegisterInput is the public registration payload. Role is limited to
// parent/specialist; admin is impossible here (SEC-13, SEC-14).
type RegisterInput struct {
	Email, Password, FirstName, LastName string
	Role                                 domain.Role
}

// Session is the result of login/refresh.
type Session struct {
	AccessToken  string
	ExpiresIn    time.Duration
	RefreshToken string // raw value, only for the cookie
}

func validateNewUser(email, pw, first, last string, role domain.Role, allowed ...domain.Role) error {
	var v domain.Validator
	v.Check(domain.ValidEmail(email), "email", "invalid")
	ok, code := domain.ValidPassword(pw)
	v.Check(ok, "password", code)
	v.Check(domain.ValidName(first), "first_name", "invalid")
	v.Check(domain.ValidName(last), "last_name", "invalid")
	roleOK := false
	for _, r := range allowed {
		roleOK = roleOK || r == role
	}
	v.Check(roleOK, "role", "invalid")
	return v.Err()
}

// Register creates a parent or (unapproved) specialist account.
func (s *Service) Register(ctx context.Context, in RegisterInput) (domain.User, error) {
	in.Email = domain.NormalizeEmail(in.Email)
	if err := validateNewUser(in.Email, in.Password, in.FirstName, in.LastName, in.Role, domain.RoleParent, domain.RoleSpecialist); err != nil {
		return domain.User{}, err
	}
	return s.createUser(ctx, nil, in.Email, in.Password, in.FirstName, in.LastName, in.Role, "register")
}

func (s *Service) createUser(ctx context.Context, by *Actor, email, pw, first, last string, role domain.Role, action string) (domain.User, error) {
	hash, err := s.hasher.Hash(pw)
	if err != nil {
		return domain.User{}, err
	}
	now := s.now()
	u := domain.User{ID: newID(), Email: email, PasswordHash: hash, FirstName: strings.TrimSpace(first),
		LastName: strings.TrimSpace(last), Role: role, Status: domain.StatusActive, CreatedAt: now, UpdatedAt: now}
	err = s.uow.Within(ctx, func(r Repo) error {
		if err := r.CreateUser(ctx, &u); err != nil {
			return err
		}
		return r.InsertAudit(ctx, s.event(ctx, by, nil, action, "user", &u.ID, domain.OutcomeSuccess))
	})
	return u, err
}

func identifierKey(email string) string {
	sum := sha256.Sum256([]byte("login:" + email))
	return hex.EncodeToString(sum[:])
}

// Login authenticates with uniform errors and timing (SEC-03, SEC-04).
func (s *Service) Login(ctx context.Context, email, password string) (Session, error) {
	email = domain.NormalizeEmail(email)
	now := s.now()
	if !s.loginLimiter.Allow(identifierKey(email), now) {
		s.auditOutside(ctx, nil, nil, "auth:login", "session", nil, domain.OutcomeDenied)
		return Session{}, domain.ErrRateLimited
	}
	u, err := s.uow.Repo().UserByEmail(ctx, email)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return Session{}, err
	}
	found := err == nil
	hash := s.dummyHash
	if found {
		hash = u.PasswordHash
	}
	ok, rehash, verr := s.hasher.Verify(hash, password)
	if verr != nil || !found || !ok || !u.Active(now) {
		if found && !ok {
			s.recordFailure(ctx, u, now)
		}
		var actor *Actor
		if found {
			actor = &Actor{ID: u.ID, Role: u.Role}
		}
		s.auditOutside(ctx, actor, nil, "auth:login", "session", nil, domain.OutcomeFailure)
		return Session{}, domain.ErrUnauthorized
	}
	u.FailedLoginAttempts, u.LockedUntil = 0, nil
	if rehash {
		if h, err := s.hasher.Hash(password); err == nil {
			u.PasswordHash = h
		}
	}
	u.UpdatedAt = now
	var sess Session
	err = s.uow.Within(ctx, func(r Repo) error {
		if err := r.UpdateUser(ctx, &u); err != nil {
			return err
		}
		var err error
		sess, err = s.issue(ctx, r, u.ID, newID(), now)
		if err != nil {
			return err
		}
		return r.InsertAudit(ctx, s.event(ctx, &Actor{ID: u.ID, Role: u.Role}, nil, "auth:login", "session", nil, domain.OutcomeSuccess))
	})
	return sess, err
}

func (s *Service) recordFailure(ctx context.Context, u domain.User, now time.Time) {
	u.FailedLoginAttempts++
	if s.cfg.MaxFailedLogins > 0 && u.FailedLoginAttempts >= s.cfg.MaxFailedLogins {
		until := now.Add(s.cfg.LockoutDuration)
		u.LockedUntil, u.FailedLoginAttempts = &until, 0
	}
	u.UpdatedAt = now
	_ = s.uow.Repo().UpdateUser(context.WithoutCancel(ctx), &u)
}

// HashToken returns the SHA-256 of an opaque refresh token.
func HashToken(raw string) []byte {
	sum := sha256.Sum256([]byte(raw))
	return sum[:]
}

func (s *Service) issue(ctx context.Context, r Repo, user, family uuid.UUID, now time.Time) (Session, error) {
	access, ttl, err := s.tokens.Issue(user, now)
	if err != nil {
		return Session{}, err
	}
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return Session{}, err
	}
	raw := base64.RawURLEncoding.EncodeToString(buf)
	rt := domain.RefreshToken{ID: newID(), UserID: user, FamilyID: family, TokenHash: HashToken(raw),
		ExpiresAt: now.Add(s.cfg.RefreshTTL), CreatedAt: now}
	if err := r.InsertRefreshToken(ctx, &rt); err != nil {
		return Session{}, err
	}
	return Session{AccessToken: access, ExpiresIn: ttl, RefreshToken: raw}, nil
}

// Refresh rotates the refresh token; reuse revokes the family (ADR-0004).
func (s *Service) Refresh(ctx context.Context, raw string) (Session, error) {
	if raw == "" {
		return Session{}, domain.ErrUnauthorized
	}
	now := s.now()
	var sess Session
	var reusedUser, family uuid.UUID
	action := "auth:refresh_reuse"
	err := s.uow.Within(ctx, func(r Repo) error {
		t, err := r.RefreshTokenByHash(ctx, HashToken(raw))
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrUnauthorized
		} else if err != nil {
			return err
		}
		if t.RotatedAt != nil || t.RevokedAt != nil {
			reusedUser, family = t.UserID, t.FamilyID
			if t.RevokedAt != nil {
				action = "auth:refresh_revoked"
			}
			return errReuse
		}
		if !t.ExpiresAt.After(now) {
			return domain.ErrUnauthorized
		}
		u, err := r.UserByID(ctx, t.UserID)
		if err != nil || !u.Active(now) {
			return domain.ErrUnauthorized
		}
		if err := r.MarkRefreshRotated(ctx, t.ID, now); err != nil {
			return err
		}
		sess, err = s.issue(ctx, r, u.ID, t.FamilyID, now)
		return err
	})
	if errors.Is(err, errReuse) {
		// Presenting a rotated/revoked token: revoke the whole family.
		if rerr := s.uow.Repo().RevokeRefreshFamily(context.WithoutCancel(ctx), family, now); rerr != nil {
			return Session{}, rerr
		}
		s.auditOutside(ctx, &Actor{ID: reusedUser}, nil, action, "session", nil, domain.OutcomeDenied)
		return Session{}, domain.ErrUnauthorized
	}
	return sess, err
}

var errReuse = errors.New("refresh token reuse")

// Logout revokes the token family. Unknown tokens are ignored (idempotent).
func (s *Service) Logout(ctx context.Context, raw string) error {
	if raw == "" {
		return nil
	}
	r := s.uow.Repo()
	t, err := r.RefreshTokenByHash(ctx, HashToken(raw))
	if errors.Is(err, domain.ErrNotFound) {
		return nil
	} else if err != nil {
		return err
	}
	now := s.now()
	if err := r.RevokeRefreshFamily(ctx, t.FamilyID, now); err != nil {
		return err
	}
	return r.InsertAudit(ctx, s.event(ctx, &Actor{ID: t.UserID}, nil, "auth:logout", "session", nil, domain.OutcomeSuccess))
}

// Authenticate verifies an access token and reloads role/status from the
// database on every request (SEC-05, SEC-13).
func (s *Service) Authenticate(ctx context.Context, token string) (*Actor, error) {
	now := s.now()
	id, err := s.tokens.Verify(token, now)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	u, err := s.uow.Repo().UserByID(ctx, id)
	if err != nil || u.Status != domain.StatusActive {
		return nil, domain.ErrUnauthorized
	}
	return &Actor{ID: u.ID, Role: u.Role}, nil
}

// Me returns the caller's own account.
func (s *Service) Me(ctx context.Context, a *Actor) (domain.User, error) {
	if a == nil {
		return domain.User{}, domain.ErrUnauthorized
	}
	return s.uow.Repo().UserByID(ctx, a.ID)
}
