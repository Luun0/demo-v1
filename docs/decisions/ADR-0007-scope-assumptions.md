# ADR-0007: Scope assumptions for the MVP implementation

- Status: Accepted (SSP-51)

1. Password reset endpoints are in the contract with `x-status: planned`; no mail delivery is implemented (SF-01 scope control).
2. Specialist self-registration creates an account with `approved=false`; such a specialist cannot be assigned until an admin approves (SRS: "approved Specialist").
3. Consent is recorded at child creation (`policy_version`, timestamp, guardian); a child without consent cannot be created.
4. Admin is created only via `cmd/api bootstrap-admin` or by an existing admin (SEC-14).
5. Only synthetic seed data (`cmd/api seed`), PRV-08.
6. Legal compliance is NOT claimed; all legal points stay open questions.
