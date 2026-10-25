# ADR-0004: Tokens, cookies and CSRF

- Status: Accepted (SSP-51). Refines `authentication-and-session-decisions.md`; resolves a conflict with openapi v0.1 (refresh token in JSON).

## Decision
- Access token: JWT HS256, 15 min, claims `sub`, `iss`, `aud`, `exp`, `iat`, `jti`. `alg` is pinned; `none`/RS*/other algs rejected. Role and account status are re-read from the database on every request (immediate effect of deactivation; role is never trusted from the client).
- Refresh token: 32 random bytes, base64url, cookie `ssp_refresh` with `HttpOnly; Secure; SameSite=Strict; Path=/api/v1/auth`. Only SHA-256 is stored. Every refresh rotates; presenting an already-rotated token revokes the whole family (reuse detection). Logout revokes the family. SameSite=Strict instead of Lax: the SPA and API share one origin, so nothing breaks.
- CSRF: double-submit cookie `ssp_csrf` (readable) must equal header `X-CSRF-Token` on `/auth/refresh` and `/auth/logout`, plus Origin allowlist check. Bearer endpoints are not CSRF-able.
- `JWT_SECRET` must be ≥ 32 bytes, otherwise the API refuses to start.
- Passwords: Argon2id (m=64 MiB, t=3, p=2 by default, configurable), PHC string; rehash on login when parameters change. ≥10 chars and not in a common-password list.

## Consequences
openapi v0.2: login/refresh return `AccessToken`; `TokenPair` and `RefreshRequest` removed.
