// Package testpg provides real PostgreSQL databases for tests. It uses
// TEST_DATABASE_URL when set (CI service container); otherwise it starts a
// throw-away cluster with initdb/pg_ctl. Tests never silently skip: if no
// PostgreSQL is available, the test binary fails.
package testpg

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/AknurOE/sen-support-platform/internal/adapter/postgres"
)

var (
	baseURL string
	counter atomic.Int64
	cleanup func()
)

func pgBin(name string) (string, error) {
	if dir := os.Getenv("PG_BIN"); dir != "" {
		return filepath.Join(dir, name), nil
	}
	if p, err := exec.LookPath(name); err == nil {
		return p, nil
	}
	m, _ := filepath.Glob("/usr/lib/postgresql/*/bin/" + name)
	sort.Strings(m)
	if len(m) == 0 {
		return "", fmt.Errorf("%s not found: install PostgreSQL or set TEST_DATABASE_URL", name)
	}
	return m[len(m)-1], nil
}

func freePort() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer func() { _ = l.Close() }()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func startLocal() (string, func(), error) {
	initdb, err := pgBin("initdb")
	if err != nil {
		return "", nil, err
	}
	pgctl, err := pgBin("pg_ctl")
	if err != nil {
		return "", nil, err
	}
	dir, err := os.MkdirTemp("", "ssp-pg-")
	if err != nil {
		return "", nil, err
	}
	data := filepath.Join(dir, "data")
	if out, err := exec.Command(initdb, "-D", data, "-U", "postgres", "--auth=trust", "-E", "UTF8", "--no-sync").CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("initdb: %w: %s", err, out)
	}
	port, err := freePort()
	if err != nil {
		return "", nil, err
	}
	opts := fmt.Sprintf("-p %d -k %s -c listen_addresses=127.0.0.1 -c fsync=off -c synchronous_commit=off -c full_page_writes=off", port, dir)
	if out, err := exec.Command(pgctl, "-D", data, "-o", opts, "-w", "-l", filepath.Join(dir, "log"), "start").CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("pg_ctl start: %w: %s", err, out)
	}
	stop := func() {
		_ = exec.Command(pgctl, "-D", data, "-m", "immediate", "stop").Run()
		_ = os.RemoveAll(dir)
	}
	return fmt.Sprintf("postgres://postgres@127.0.0.1:%d/postgres?sslmode=disable", port), stop, nil
}

// Main runs the tests of a package against PostgreSQL; call from TestMain.
func Main(m *testing.M) int {
	baseURL = os.Getenv("TEST_DATABASE_URL")
	if baseURL == "" {
		u, stop, err := startLocal()
		if err != nil {
			fmt.Fprintln(os.Stderr, "testpg:", err)
			return 1
		}
		baseURL, cleanup = u, stop
	}
	if err := prepareTemplate(); err != nil {
		fmt.Fprintln(os.Stderr, "testpg template:", err)
		if cleanup != nil {
			cleanup()
		}
		return 1
	}
	code := m.Run()
	if cleanup != nil {
		cleanup()
	}
	return code
}

func withDB(u, db string) string {
	p, _ := url.Parse(u)
	p.Path = "/" + db
	return p.String()
}

const template = "ssp_template"

func adminExec(ctx context.Context, sql string) error {
	c, err := pgx.Connect(ctx, baseURL)
	if err != nil {
		return err
	}
	defer func() { _ = c.Close(ctx) }()
	_, err = c.Exec(ctx, sql)
	return err
}

func prepareTemplate() error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	tmpl := fmt.Sprintf("%s_%d", template, os.Getpid())
	_ = adminExec(ctx, "DROP DATABASE IF EXISTS "+tmpl)
	if err := adminExec(ctx, "CREATE DATABASE "+tmpl); err != nil {
		return err
	}
	s, err := postgres.Open(ctx, withDB(baseURL, tmpl))
	if err != nil {
		return err
	}
	defer s.Close()
	return s.Migrate(ctx)
}

// NewDB returns the DSN of a fresh, migrated database dropped after t.
func NewDB(t testing.TB) string {
	t.Helper()
	ctx := context.Background()
	name := fmt.Sprintf("ssp_t_%d_%d", os.Getpid(), counter.Add(1))
	tmpl := fmt.Sprintf("%s_%d", template, os.Getpid())
	if err := adminExec(ctx, fmt.Sprintf("CREATE DATABASE %s TEMPLATE %s", name, tmpl)); err != nil {
		t.Fatalf("create db: %v", err)
	}
	t.Cleanup(func() { _ = adminExec(context.Background(), "DROP DATABASE IF EXISTS "+name+" WITH (FORCE)") })
	return withDB(baseURL, name)
}
