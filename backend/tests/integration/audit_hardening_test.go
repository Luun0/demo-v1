package integration

import (
	"context"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
)

func TestApplicationRoleCannotTamperWithAuditRows(t *testing.T) {
	e := newEnv(t, nil)
	conn, err := pgx.Connect(context.Background(), e.dbURL)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(context.Background())
	if _, err := conn.Exec(context.Background(), `SET ROLE ssp_app`); err != nil {
		t.Fatal(err)
	}
	for _, query := range []string{
		`UPDATE audit_events SET action='tampered'`,
		`DELETE FROM audit_events`,
	} {
		_, err := conn.Exec(context.Background(), query)
		if err == nil || !strings.Contains(err.Error(), "permission denied") {
			t.Errorf("%q: expected permission denied, got %v", query, err)
		}
	}
}
