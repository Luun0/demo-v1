// Package app is the composition root shared by cmd/api and tests.
package app

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/AknurOE/sen-support-platform/internal/adapter/httpapi"
	"github.com/AknurOE/sen-support-platform/internal/adapter/postgres"
	"github.com/AknurOE/sen-support-platform/internal/application"
	"github.com/AknurOE/sen-support-platform/internal/infrastructure/config"
	"github.com/AknurOE/sen-support-platform/internal/infrastructure/logging"
	"github.com/AknurOE/sen-support-platform/internal/infrastructure/ratelimit"
	"github.com/AknurOE/sen-support-platform/internal/infrastructure/security"
)

type systemClock struct{}

func (systemClock) Now() time.Time { return time.Now() }

// App bundles the wired service and handler.
type App struct {
	Store   *postgres.Store
	Service *application.Service
	Handler http.Handler
	Log     *slog.Logger
}

// Options allow tests to override the clock.
type Options struct {
	Clock application.Clock
	Logs  io.Writer
}

// Build wires all dependencies from configuration.
func Build(ctx context.Context, cfg config.Config, opt Options) (*App, error) {
	if opt.Clock == nil {
		opt.Clock = systemClock{}
	}
	if opt.Logs == nil {
		opt.Logs = io.Discard
	}
	log := logging.New(opt.Logs)
	store, err := postgres.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	hasher := security.Argon2{Memory: cfg.Argon2MemoryKiB, Time: cfg.Argon2Time, Threads: cfg.Argon2Threads}
	tokens := security.JWT{Secret: cfg.JWTSecret, Issuer: cfg.JWTIssuer, Audience: cfg.JWTAudience, TTL: cfg.AccessTTL}
	svc, err := application.NewService(store, opt.Clock, hasher, tokens, ratelimit.New(cfg.LoginPerMinute, time.Minute),
		application.Config{RefreshTTL: cfg.RefreshTTL, MaxFailedLogins: cfg.MaxFailedLogins, LockoutDuration: cfg.LockoutDuration,
			AuditRetentionDays: cfg.AuditRetentionDays})
	if err != nil {
		store.Close()
		return nil, err
	}
	h := httpapi.New(svc, httpapi.Config{AllowedOrigins: cfg.AllowedOrigins, MaxBodyBytes: cfg.MaxBodyBytes,
		CookieSecure: cfg.CookieSecure, RefreshTTL: cfg.RefreshTTL}, log,
		ratelimit.New(cfg.APIPerMinute, time.Minute), ratelimit.New(cfg.LoginPerMinute*3, time.Minute))
	return &App{Store: store, Service: svc, Handler: h, Log: log}, nil
}
