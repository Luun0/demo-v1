// Command api runs the SEN Support Platform API.
//
//	api [serve]                      run the HTTP server (migrates first)
//	api migrate                      apply migrations
//	api seed                         insert synthetic demo data (PRV-08)
//	api bootstrap-admin              create the first admin from ADMIN_EMAIL / ADMIN_PASSWORD
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AknurOE/sen-support-platform/internal/app"
	"github.com/AknurOE/sen-support-platform/internal/infrastructure/config"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	cmd := "serve"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}
	if cmd == "healthcheck" { // for distroless HEALTHCHECK
		resp, err := http.Get("http://127.0.0.1:8080/healthz") //nolint:noctx // tiny probe
		if err != nil {
			return err
		}
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("status %d", resp.StatusCode)
		}
		return nil
	}
	cfg, err := config.Load(os.Getenv)
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	a, err := app.Build(ctx, cfg, app.Options{Logs: os.Stdout})
	if err != nil {
		return err
	}
	defer a.Store.Close()
	if cmd == "migrate" || cfg.AutoMigrate {
		if err := a.Store.Migrate(ctx); err != nil {
			return err
		}
	}
	switch cmd {
	case "migrate":
		return nil
	case "seed":
		return app.Seed(ctx, a.Service)
	case "bootstrap-admin":
		_, err := a.Service.BootstrapAdmin(ctx, os.Getenv("ADMIN_EMAIL"), os.Getenv("ADMIN_PASSWORD"), "Platform", "Admin")
		return err
	case "retention":
		_, err := a.Service.RunRetention(ctx)
		return err
	case "serve":
	default:
		return fmt.Errorf("unknown command %q", cmd)
	}
	srv := &http.Server{Addr: cfg.Addr, Handler: a.Handler, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second,
		WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 16 << 10}
	errc := make(chan error, 1)
	go func() { errc <- srv.ListenAndServe() }()
	a.Log.Info("listening", "addr", cfg.Addr)
	select {
	case err := <-errc:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-ctx.Done():
	}
	shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdown)
}
