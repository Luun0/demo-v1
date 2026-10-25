// Package config loads typed configuration from the environment and fails
// fast on missing or weak values (no default secrets).
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config is the API configuration.
type Config struct {
	Addr                  string
	DatabaseURL           string
	JWTSecret             []byte
	JWTIssuer             string
	JWTAudience           string
	AccessTTL             time.Duration
	RefreshTTL            time.Duration
	AllowedOrigins        []string
	CookieSecure          bool
	MaxBodyBytes          int64
	LoginPerMinute        int
	APIPerMinute          int
	MaxFailedLogins       int
	LockoutDuration       time.Duration
	Argon2MemoryKiB       uint32
	Argon2Time            uint32
	Argon2Threads         uint8
	AuditRetentionDays    int
	AutoMigrate           bool
	PasswordResetBaseURL  string
	LogPasswordResetLinks bool
}

// MinSecretLen is the minimum JWT secret length in bytes (ADR-0004).
const MinSecretLen = 32

// Load reads the environment via getenv (os.Getenv in production).
func Load(getenv func(string) string) (Config, error) {
	if getenv == nil {
		getenv = os.Getenv
	}
	var errs []error
	str := func(k, def string) string {
		if v := strings.TrimSpace(getenv(k)); v != "" {
			return v
		}
		return def
	}
	num := func(k string, def int) int {
		v := getenv(k)
		if v == "" {
			return def
		}
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			errs = append(errs, fmt.Errorf("%s must be a positive integer", k))
		}
		return n
	}
	dur := func(k string, def time.Duration) time.Duration {
		v := getenv(k)
		if v == "" {
			return def
		}
		d, err := time.ParseDuration(v)
		if err != nil || d <= 0 {
			errs = append(errs, fmt.Errorf("%s must be a positive duration", k))
		}
		return d
	}
	c := Config{
		Addr:                  str("API_ADDR", ":8080"),
		DatabaseURL:           str("DATABASE_URL", ""),
		JWTSecret:             []byte(getenv("JWT_SECRET")),
		JWTIssuer:             str("JWT_ISSUER", "sen-support-platform"),
		JWTAudience:           str("JWT_AUDIENCE", "sen-support-platform-web"),
		AccessTTL:             dur("ACCESS_TOKEN_TTL", 15*time.Minute),
		RefreshTTL:            dur("REFRESH_TOKEN_TTL", 7*24*time.Hour),
		CookieSecure:          str("COOKIE_SECURE", "true") != "false",
		MaxBodyBytes:          int64(num("MAX_BODY_BYTES", 64*1024)),
		LoginPerMinute:        num("LOGIN_RATE_PER_MINUTE", 10),
		APIPerMinute:          num("API_RATE_PER_MINUTE", 600),
		MaxFailedLogins:       num("MAX_FAILED_LOGINS", 5),
		LockoutDuration:       dur("LOCKOUT_DURATION", 15*time.Minute),
		Argon2MemoryKiB:       uint32(num("ARGON2_MEMORY_KIB", 64*1024)), //nolint:gosec // bounded positive int
		Argon2Time:            uint32(num("ARGON2_TIME", 3)),             //nolint:gosec // bounded positive int
		Argon2Threads:         uint8(num("ARGON2_THREADS", 2)),           //nolint:gosec // bounded positive int
		AuditRetentionDays:    num("AUDIT_RETENTION_DAYS", 365),
		AutoMigrate:           str("AUTO_MIGRATE", "true") != "false",
		PasswordResetBaseURL:  str("PASSWORD_RESET_BASE_URL", ""),
		LogPasswordResetLinks: str("PASSWORD_RESET_LOG_LINK", "false") == "true",
	}
	for _, o := range strings.Split(str("CORS_ALLOWED_ORIGINS", ""), ",") {
		if o = strings.TrimSpace(o); o != "" {
			if o == "*" {
				errs = append(errs, errors.New("CORS_ALLOWED_ORIGINS must not contain *"))
			}
			c.AllowedOrigins = append(c.AllowedOrigins, o)
		}
	}
	if c.PasswordResetBaseURL == "" && len(c.AllowedOrigins) > 0 {
		c.PasswordResetBaseURL = c.AllowedOrigins[0]
	}
	if c.LogPasswordResetLinks && c.CookieSecure {
		errs = append(errs, errors.New("PASSWORD_RESET_LOG_LINK requires COOKIE_SECURE=false; reset links must not be logged in secure deployments"))
	}
	if c.DatabaseURL == "" {
		errs = append(errs, errors.New("DATABASE_URL is required"))
	}
	if len(c.JWTSecret) < MinSecretLen {
		errs = append(errs, fmt.Errorf("JWT_SECRET must be at least %d bytes", MinSecretLen))
	}
	if c.AccessTTL > 15*time.Minute {
		errs = append(errs, errors.New("ACCESS_TOKEN_TTL must be <= 15m"))
	}
	return c, errors.Join(errs...)
}
