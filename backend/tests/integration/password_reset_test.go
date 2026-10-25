package integration

import (
	"context"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/AknurOE/sen-support-platform/internal/application"
)

var resetLinkPattern = regexp.MustCompile(`development password reset link: https://app\.example\.test/reset-password\?token=([A-Za-z0-9_-]+)`)

func lastResetToken(t *testing.T, logs string) string {
	t.Helper()
	matches := resetLinkPattern.FindAllStringSubmatch(logs, -1)
	if len(matches) == 0 {
		t.Fatalf("development reset link missing from configured test log")
	}
	return matches[len(matches)-1][1]
}

func TestPasswordResetUniformResponseSingleUseAndExpiry(t *testing.T) {
	e := newEnv(t, map[string]string{"PASSWORD_RESET_LOG_LINK": "true"})
	w := e.world()
	unknown := e.client().do("POST", "/api/v1/auth/password-reset/request", map[string]string{"email": "missing@example.test"})
	if unknown.code != 202 || strings.Contains(e.logs.String(), "development password reset link:") {
		t.Fatalf("unknown account should get generic 202 without delivery, got %d %s", unknown.code, e.logs.String())
	}
	known := e.client().do("POST", "/api/v1/auth/password-reset/request", map[string]string{"email": "p1@example.test"})
	if known.code != unknown.code || known.body != unknown.body {
		t.Fatalf("known and unknown account responses differ: known=%d %s unknown=%d %s", known.code, known.body, unknown.code, unknown.body)
	}
	token := lastResetToken(t, e.logs.String())
	newPassword := "Reset-Password-42"
	confirm := func(value, password string) resp {
		return e.client().do("POST", "/api/v1/auth/password-reset/confirm", map[string]string{"token": value, "password": password})
	}
	if r := confirm(token, newPassword); r.code != 204 {
		t.Fatalf("confirm reset: %d %s", r.code, r.body)
	}
	if r := confirm(token, "Another-Password-43"); r.code != 401 {
		t.Fatalf("reused token must be rejected, got %d %s", r.code, r.body)
	}
	if r := e.client().login("p1@example.test", pw); r.code != 401 {
		t.Fatalf("old password should fail after reset, got %d", r.code)
	}
	if r := e.client().login("p1@example.test", newPassword); r.code != 200 {
		t.Fatalf("new password should work after reset, got %d %s", r.code, r.body)
	}
	if r := w.parent1.do("POST", "/api/v1/auth/refresh", nil, "X-CSRF-Token", w.parent1.csrf()); r.code != 401 {
		t.Fatalf("reset should revoke refresh sessions, got %d %s", r.code, r.body)
	}

	if r := e.client().do("POST", "/api/v1/auth/password-reset/request", map[string]string{"email": "p1@example.test"}); r.code != 202 {
		t.Fatalf("request expired-token fixture: %d %s", r.code, r.body)
	}
	expired := lastResetToken(t, e.logs.String())
	conn, err := pgx.Connect(context.Background(), e.dbURL)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(context.Background())
	if _, err := conn.Exec(context.Background(), `UPDATE password_reset_tokens SET expires_at=$1 WHERE token_hash=$2`, time.Now().UTC().Add(-time.Minute), application.HashToken(expired)); err != nil {
		t.Fatal(err)
	}
	if r := confirm(expired, "Reset-Password-44"); r.code != 401 {
		t.Fatalf("expired token must be rejected, got %d %s", r.code, r.body)
	}
}
