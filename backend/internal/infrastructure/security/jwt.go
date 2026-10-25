package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWT issues HS256 access tokens with pinned algorithm, iss and aud.
type JWT struct {
	Secret   []byte
	Issuer   string
	Audience string
	TTL      time.Duration
}

// Issue signs a token for user.
func (j JWT) Issue(user uuid.UUID, now time.Time) (string, time.Duration, error) {
	claims := jwt.RegisteredClaims{
		Subject: user.String(), Issuer: j.Issuer, Audience: jwt.ClaimStrings{j.Audience},
		IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(j.TTL)), ID: uuid.NewString(),
	}
	s, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.Secret)
	return s, j.TTL, err
}

// Verify validates signature, alg, iss, aud, exp and returns the subject.
func (j JWT) Verify(token string, now time.Time) (uuid.UUID, error) {
	var claims jwt.RegisteredClaims
	_, err := jwt.ParseWithClaims(token, &claims, func(*jwt.Token) (any, error) { return j.Secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithIssuer(j.Issuer),
		jwt.WithAudience(j.Audience), jwt.WithExpirationRequired(), jwt.WithTimeFunc(func() time.Time { return now }))
	if err != nil {
		return uuid.Nil, err
	}
	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, errors.New("invalid subject")
	}
	return id, nil
}
