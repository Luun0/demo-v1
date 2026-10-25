# ADR-0005: Audit events and retention

- Status: Accepted (SSP-51); legal retention period is an open question (see kazakhstan-privacy-legal-checklist.md).

## Decision
Exactly one audit event per authorised or denied action, written in the same transaction as the change (UnitOfWork). Fields: actor, child, action, resource, resource_id, outcome (success/denied/failure), request_id, occurred_at. Never diary content, passwords, tokens, or the attempted email in clear text (failed logins store a SHA-256 of the normalised identifier).
Retention assumption: 365 days by default, configurable with `AUDIT_RETENTION_DAYS`. The operator-run `api retention` command deletes rows older than that cutoff and can be scheduled with the deployment's existing scheduler. Partitioning is deferred until volume requires it.
