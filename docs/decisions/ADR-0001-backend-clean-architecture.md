# ADR-0001: Backend follows Clean Architecture

- Status: Accepted (SSP-51, 2026-10-24)

## Context
The backend was only domain structs (audit C1). The main risk is server-enforced authorization (SEC-11/12), which must be testable without HTTP or a database.

## Decision
Layers with inward dependencies only:
`internal/domain` (entities, enums, invariants, `policy.Decide`) ← `internal/application` (use cases + ports) ← `internal/adapter/{http,postgres,memory}` ← `internal/infrastructure` (config, security, logging) ← `cmd/api` (composition root).
`domain` and `application` must not import `net/http`, `pgx` or `jwt`; enforced by `depguard` in `backend/.golangci.yml`.
Router: stdlib `net/http` ServeMux (Go 1.22+ patterns), no framework.

## Consequences
Use cases are unit tested against an in-memory adapter; the Postgres adapter is covered by integration tests against a real Postgres. The old `transport/http/dto`, `service`, `repository` folders are replaced.
