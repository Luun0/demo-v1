package config

import (
	"strings"
	"testing"
)

func env(m map[string]string) func(string) string { return func(k string) string { return m[k] } }

func TestLoadFailsFastOnWeakOrMissingSecrets(t *testing.T) {
	_, err := Load(env(map[string]string{}))
	if err == nil || !strings.Contains(err.Error(), "DATABASE_URL") || !strings.Contains(err.Error(), "JWT_SECRET") {
		t.Fatalf("got %v", err)
	}
	_, err = Load(env(map[string]string{"DATABASE_URL": "postgres://x", "JWT_SECRET": "short"}))
	if err == nil {
		t.Fatal("short secret must fail")
	}
	_, err = Load(env(map[string]string{"DATABASE_URL": "postgres://x", "JWT_SECRET": strings.Repeat("a", 32), "CORS_ALLOWED_ORIGINS": "*", "ACCESS_TOKEN_TTL": "1h", "MAX_BODY_BYTES": "-1", "LOCKOUT_DURATION": "bad"}))
	for _, want := range []string{"CORS", "ACCESS_TOKEN_TTL", "MAX_BODY_BYTES", "LOCKOUT_DURATION"} {
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("expected error about %s, got %v", want, err)
		}
	}
}

func TestLoadDefaults(t *testing.T) {
	c, err := Load(env(map[string]string{"DATABASE_URL": "postgres://x", "JWT_SECRET": strings.Repeat("a", 32), "CORS_ALLOWED_ORIGINS": "https://a.test, https://b.test"}))
	if err != nil {
		t.Fatal(err)
	}
	if c.AccessTTL.Minutes() != 15 || !c.CookieSecure || len(c.AllowedOrigins) != 2 || c.MaxBodyBytes != 65536 {
		t.Fatalf("%+v", c)
	}
}
