# ADR-0002: Frontend follows Feature-Sliced Design

- Status: Accepted (SSP-51)

## Decision
Layers `app → pages → widgets → features → entities → shared`, imports only downward, public API via `index.ts`, enforced with `eslint-plugin-boundaries` in CI. Server state via TanStack Query, routing via React Router, types generated from `docs/openapi.yaml` (`openapi-typescript`).
`AppContext.tsx` / `Workspace.tsx` are removed once all screens use the API.

## Consequences
Migration happens screen by screen (Phase 4). Remaining gaps are listed honestly in the phase report.
