# ADR-0006: Pagination and error format

- Status: Accepted (SSP-51)

## Decision
- Errors: RFC 9457 `application/problem+json` with `type`, `title`, `status`, `request_id`, optional `errors[]` (field + code, never the submitted value). No stack traces.
- 404 when no active relationship exists (existence not revealed); 403 when related but not permitted; 422 for validation; 409 for conflicts; 413 body > 64 KiB; 429 with `Retry-After`.
- Lists keep the v0.1 array shape (backward compatible); pagination via `limit` + opaque `cursor`, next cursor returned in `X-Next-Cursor`. `X-Total-Count` counts only rows visible to the caller.
