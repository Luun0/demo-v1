# ADR-0003: SQL layer — pgx with hand-written parameterised queries (sqlc deferred)

- Status: Accepted (SSP-51). Deviation from the master prompt (which suggested sqlc).

## Context
sqlc adds a code-generation step and a tool dependency for ~25 queries. The security requirement is parameterised queries only (SEC-20) and visibility filtering inside SQL (SEC-12).

## Decision
`pgx/v5` (`pgxpool`) with hand-written queries that only use `$n` placeholders; no string concatenation of user input (gosec G201/G202 enforced). Migrations are plain SQL files applied by a tiny embedded migrator (`cmd/api migrate`), up-only with checksums recorded in `schema_migrations`.
Visibility: parent queries always add `AND visibility = 'shared'`, including COUNT, search and export.

## Consequences
Switching to sqlc later is mechanical because queries already live in one adapter package.
