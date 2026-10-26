package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/AknurOE/sen-support-platform/internal/application"
	"github.com/AknurOE/sen-support-platform/internal/domain"
)

// BenchmarkDiaryListVisibilityFilter measures the Postgres list query used by
// parents (shared-only) and specialists (all assigned entries). Setup is
// synthetic and excluded from the timed section.
func BenchmarkDiaryListVisibilityFilter(b *testing.B) {
	e := newEnv(b, nil)
	w := e.world()
	repo := e.app.Store.Repo()
	ctx := context.Background()
	base := time.Now().UTC()
	for i := 0; i < 1000; i++ {
		visibility := domain.VisibilityShared
		if i%4 == 0 {
			visibility = domain.VisibilitySpecialistsOnly
		}
		entry := domain.DiaryEntry{
			ID: uuid.New(), ChildID: w.child1, AuthorID: w.spec1.id,
			Type: domain.EntryObservation, Visibility: visibility,
			Content: fmt.Sprintf("synthetic diary note %04d", i), ObservedAt: base.Add(time.Duration(i) * time.Second),
			CreatedAt: base, UpdatedAt: base,
		}
		if err := repo.CreateEntry(ctx, &entry); err != nil {
			b.Fatal(err)
		}
	}

	cases := []struct {
		name   string
		filter application.DiaryFilter
	}{
		{name: "guardian_shared_page", filter: application.DiaryFilter{ChildID: w.child1, SharedOnly: true, Limit: 51}},
		{name: "specialist_all_page", filter: application.DiaryFilter{ChildID: w.child1, Limit: 51}},
		{name: "guardian_search", filter: application.DiaryFilter{ChildID: w.child1, SharedOnly: true, Query: "synthetic note 09", Limit: 51}},
		{name: "specialist_private_filter", filter: application.DiaryFilter{ChildID: w.child1, Visibility: domain.VisibilitySpecialistsOnly, Limit: 51}},
	}
	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, err := repo.ListEntries(ctx, tc.filter); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
