package integration

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func explain(t *testing.T, conn *pgx.Conn, query string, args ...any) string {
	t.Helper()
	rows, err := conn.Query(context.Background(), query, args...)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var plan []string
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			t.Fatal(err)
		}
		plan = append(plan, line)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return strings.Join(plan, "\n")
}

func TestDiaryAndSessionListPlansUseSupportingIndexes(t *testing.T) {
	e := newEnv(t, nil)
	w := e.world()
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, e.dbURL)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(ctx)

	children := []uuid.UUID{w.child1, w.child2}
	childRows := make([][]any, 998)
	base := time.Now().UTC()
	for i := range childRows {
		children = append(children, uuid.New())
		childRows[i] = []any{children[len(children)-1], "Synthetic", fmt.Sprintf("Child-%03d", i), base.AddDate(-6, 0, 0), base, base}
	}
	if _, err := conn.CopyFrom(ctx, pgx.Identifier{"children"}, []string{"id", "first_name", "last_name", "birth_date", "created_at", "updated_at"}, pgx.CopyFromRows(childRows)); err != nil {
		t.Fatal(err)
	}

	const diaryPerChild = 100
	const sessionPerChild = 10
	diaryRows := make([][]any, 0, len(children)*diaryPerChild)
	sessionRows := make([][]any, 0, len(children)*sessionPerChild)
	for ci, child := range children {
		for i := 0; i < diaryPerChild; i++ {
			at := base.Add(time.Duration(ci*diaryPerChild+i) * time.Minute)
			visibility := "shared"
			if i%4 == 0 {
				visibility = "specialists_only"
			}
			diaryRows = append(diaryRows, []any{uuid.New(), child, w.spec1.id, "observation", visibility,
				"Synthetic indexed diary note", at, at, at})
		}
		for i := 0; i < sessionPerChild; i++ {
			at := base.Add(time.Duration(ci*sessionPerChild+i) * time.Hour)
			sessionRows = append(sessionRows, []any{uuid.New(), child, w.spec1.id, w.parent1.id,
				"Synthetic session", at, at.Add(time.Hour), "online", "scheduled", at, at})
		}
	}
	if _, err := conn.CopyFrom(ctx, pgx.Identifier{"diary_entries"}, []string{"id", "child_id", "author_id", "type", "visibility", "content", "observed_at", "created_at", "updated_at"}, pgx.CopyFromRows(diaryRows)); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.CopyFrom(ctx, pgx.Identifier{"support_sessions"}, []string{"id", "child_id", "specialist_id", "created_by", "title", "starts_at", "ends_at", "mode", "status", "created_at", "updated_at"}, pgx.CopyFromRows(sessionRows)); err != nil {
		t.Fatal(err)
	}
	for _, table := range []string{"diary_entries", "support_sessions"} {
		if _, err := conn.Exec(ctx, "ANALYZE "+table); err != nil {
			t.Fatal(err)
		}
	}

	diaryQuery := `EXPLAIN (ANALYZE, BUFFERS, COSTS OFF)
		SELECT id, child_id, author_id, type, visibility, content, observed_at, created_at, updated_at
		FROM diary_entries WHERE child_id=$1 AND (NOT $2 OR visibility='shared')
		AND ($3 = '' OR visibility=$3) AND ($4 = '' OR content ILIKE '%' || $4 || '%' ESCAPE '\')
		ORDER BY observed_at DESC, id DESC LIMIT $5`
	parentPlan := explain(t, conn, diaryQuery, children[0], true, "", "", 51)
	specialistPlan := explain(t, conn, diaryQuery, children[0], false, "", "", 51)
	sessionPlan := explain(t, conn, `EXPLAIN (ANALYZE, BUFFERS, COSTS OFF)
		SELECT id, child_id, specialist_id, created_by, title, starts_at, ends_at, mode, status, created_at, updated_at
		FROM support_sessions WHERE child_id=$1 ORDER BY starts_at, id`, children[0])
	t.Logf("parent diary list plan:\n%s", parentPlan)
	t.Logf("specialist diary list plan:\n%s", specialistPlan)
	t.Logf("session list plan:\n%s", sessionPlan)
	if !strings.Contains(parentPlan, "diary_entries_access_idx") {
		t.Fatalf("parent diary plan did not use its visibility index:\n%s", parentPlan)
	}
	if !strings.Contains(specialistPlan, "diary_entries_child_order_idx") {
		t.Fatalf("specialist diary plan did not use ordered index:\n%s", specialistPlan)
	}
	if !strings.Contains(sessionPlan, "support_sessions_child_idx") {
		t.Fatalf("session plan did not use its child/start-time index:\n%s", sessionPlan)
	}
}
