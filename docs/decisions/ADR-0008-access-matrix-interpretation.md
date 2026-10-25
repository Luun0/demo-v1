# ADR-0008: Access-matrix interpretation for implementation

- Status: Accepted (SSP-51). Source: docs/rbac-and-threat-model.md.

| Ambiguity | Decision |
|---|---|
| Admin "metadata only" on child profile | Admin may GET child (id, names, birth date, timestamps); never diary/sessions content. |
| Admin vs 404 | Admin has no relationship, but is a platform role: child metadata endpoints return 200; diary endpoints return 403 (not 404) because existence is already known to admin. |
| Specialist with ended assignment | Treated as no relationship → 404, checked on every request. |
| Parent edits entry visibility | Parent entries stay `shared`; attempting `specialists_only` → 403. |
| Specialist session update | Only sessions where they are the assigned specialist. |
| Delete diary entry | Not implemented in MVP API (contract has no endpoint); recorded as gap. |
