package integration

import (
	"context"
	"testing"
	"time"

	"github.com/AknurOE/sen-support-platform/internal/domain"
	"github.com/google/uuid"
)

func TestRetentionDeletesOnlyExpiredAuditEvents(t *testing.T) {
	e := newEnv(t, nil)
	now := time.Now().UTC()
	old := &domain.AuditEvent{ID: uuid.New(), Action: "test:old", Resource: "test", Outcome: domain.OutcomeSuccess, RequestID: "test", OccurredAt: now.AddDate(0, 0, -366)}
	recent := &domain.AuditEvent{ID: uuid.New(), Action: "test:recent", Resource: "test", Outcome: domain.OutcomeSuccess, RequestID: "test", OccurredAt: now}
	if err := e.app.Store.Repo().InsertAudit(context.Background(), old); err != nil {
		t.Fatal(err)
	}
	if err := e.app.Store.Repo().InsertAudit(context.Background(), recent); err != nil {
		t.Fatal(err)
	}
	n, err := e.app.Service.RunRetention(context.Background())
	if err != nil || n != 1 {
		t.Fatalf("deleted=%d err=%v", n, err)
	}
	page, err := e.app.Store.Repo().ListAudit(context.Background(), nil, 100, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, ev := range page {
		if ev.Action == "test:old" {
			t.Fatal("expired event remains")
		}
	}
}
