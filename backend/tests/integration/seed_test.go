package integration

import (
	"context"
	"strings"
	"testing"

	"github.com/AknurOE/sen-support-platform/internal/app"
)

func TestSeedIsSyntheticAndIdempotent(t *testing.T) {
	e := newEnv(t, nil)
	ctx := context.Background()
	if err := app.Seed(ctx, e.app.Service); err != nil {
		t.Fatal(err)
	}
	if err := app.Seed(ctx, e.app.Service); err != nil {
		t.Fatal("second seed must be a no-op:", err)
	}
	p := e.client()
	must(t, p.login("parent1@example.test", app.SeedPassword), 200)
	r := must(t, p.do("GET", "/api/v1/children", nil), 200)
	if !strings.Contains(r.body, "Demo") {
		t.Fatal(r.body)
	}
	if err := e.app.Store.Migrate(ctx); err != nil {
		t.Fatal("migrations must be idempotent:", err)
	}
}
