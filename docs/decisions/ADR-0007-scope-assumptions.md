# ADR-0007: Scope assumptions for the MVP implementation

- Status: Accepted (SSP-51)

1. Password reset is implemented with a one-hour, single-use SHA-256-hashed token and a uniform request response. No SMTP provider is configured; local Compose mode may write the active reset link to API logs when `PASSWORD_RESET_LOG_LINK=true` and `COOKIE_SECURE=false`. Treat those logs as credentials and disable this mode outside local development.
2. Specialist self-registration creates an account with `approved=false`; such a specialist cannot be assigned until an admin approves (SRS: "approved Specialist").
3. Consent is recorded at child creation (`policy_version`, timestamp, guardian); a child without consent cannot be created.
4. Admin is created only via `cmd/api bootstrap-admin` or by an existing admin (SEC-14).
5. Only synthetic seed data (`cmd/api seed`), PRV-08.
6. Legal compliance is NOT claimed; all legal points stay open questions.
