// Package postgres implements application.Repo with pgx and parameterised
// SQL only (ADR-0003, SEC-20).
package postgres

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"sort"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AknurOE/sen-support-platform/internal/application"
	"github.com/AknurOE/sen-support-platform/internal/domain"
	"github.com/AknurOE/sen-support-platform/migrations"
)

type querier interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Repo implements application.Repo over a pool or a transaction.
type Repo struct{ q querier }

// Store is the UnitOfWork.
type Store struct{ pool *pgxpool.Pool }

var _ application.UnitOfWork = (*Store)(nil)
var _ application.Repo = (*Repo)(nil)

// Open connects and verifies the database.
func Open(ctx context.Context, dsn string) (*Store, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return &Store{pool: pool}, nil
}

// Close releases the pool.
func (s *Store) Close() { s.pool.Close() }

// Ping checks connectivity.
func (s *Store) Ping(ctx context.Context) error { return s.pool.Ping(ctx) }

// Repo returns a non-transactional repository.
func (s *Store) Repo() application.Repo { return &Repo{q: s.pool} }

// Within runs fn in a transaction, committing only if fn returns nil.
func (s *Store) Within(ctx context.Context, fn func(application.Repo) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(context.WithoutCancel(ctx)) }()
	if err := fn(&Repo{q: tx}); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// Migrate applies embedded migrations in lexical order, once each.
func (s *Store) Migrate(ctx context.Context) error {
	if _, err := s.pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (name text PRIMARY KEY, applied_at timestamptz NOT NULL DEFAULT now())`); err != nil {
		return err
	}
	names, err := fs.Glob(migrations.FS, "*.sql")
	if err != nil {
		return err
	}
	sort.Strings(names)
	for _, n := range names {
		var exists bool
		if err := s.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE name=$1)`, n).Scan(&exists); err != nil {
			return err
		}
		if exists {
			continue
		}
		body, err := migrations.FS.ReadFile(n)
		if err != nil {
			return err
		}
		if err := s.Within(ctx, func(r application.Repo) error {
			q := r.(*Repo).q
			if _, err := q.Exec(ctx, string(body)); err != nil {
				return fmt.Errorf("%s: %w", n, err)
			}
			_, err := q.Exec(ctx, `INSERT INTO schema_migrations(name) VALUES ($1)`, n)
			return err
		}); err != nil {
			return err
		}
	}
	return nil
}

func mapErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	var pg *pgconn.PgError
	if errors.As(err, &pg) {
		switch pg.Code {
		case "23505":
			return domain.ErrConflict
		case "23503", "23514":
			return &domain.ValidationError{Fields: []domain.FieldError{{Field: pg.ConstraintName, Code: "constraint"}}}
		}
	}
	return err
}

func (r *Repo) exec(ctx context.Context, sql string, args ...any) error {
	_, err := r.q.Exec(ctx, sql, args...)
	return mapErr(err)
}

func (r *Repo) execOne(ctx context.Context, sql string, args ...any) error {
	tag, err := r.q.Exec(ctx, sql, args...)
	if err != nil {
		return mapErr(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func collect[T any](rows pgx.Rows, err error, scan func(pgx.Row) (T, error)) ([]T, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []T{}
	for rows.Next() {
		v, err := scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
