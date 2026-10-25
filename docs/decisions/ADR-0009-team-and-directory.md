# ADR-0009: Team view and specialist directory endpoints

- Status: Accepted (SSP-54)

## Context
The UI needs to show who is linked to a child and to let a primary guardian pick a specialist, but the v0.1 contract had no read endpoint for assignments and parents could not discover specialist IDs.

## Decision
- `GET /children/{childId}/team` — guardians and assignments with names and dates only (no emails). Authorized as `child:view` (404 without relationship).
- `GET /specialists` — approved, active specialists (id + names). New policy action `specialist:directory`: parent and admin allowed, specialist denied. Added to the access-matrix test.

## Consequences
Specialist names become visible to any parent account. Accepted: they are professionals offering services on the platform. Open question for the privacy review.
