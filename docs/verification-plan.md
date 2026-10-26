# Verification and Evaluation Plan

**Project:** SEN Support Platform
**Owner:** Ondasyn Aknur (SE-2420), Security, Privacy & Quality Assurance
**Status:** Verification plan for the university MVP prototype (26 October 2026), revision 5. Password reset and parent-owned diary deletion are implemented; Trivy passed in CI; the first ZAP baseline ran with warnings. Field-level encryption, warning triage, and evaluations remain open as identified below.

## 1. Aim of verification

To demonstrate that the platform meets the security, privacy and quality requirements in [security-privacy-requirements.md](security-privacy-requirements.md), and in particular that a user can access only the child data they are entitled to ([rbac-and-threat-model.md](rbac-and-threat-model.md)).

## 2. Test levels

| Level | Tools | What is verified | Requirements |
| --- | --- | --- | --- |
| Unit | Go `testing` | Password hashing, token creation and validation, RBAC decision function, visibility filter, input validation | SEC-01, 02, 05, 10, 12, 20 |
| Integration | Go `testing` with `net/http/httptest` and a test PostgreSQL instance (`/backend/tests/integration`) | Full request path: authentication, authorization, database, audit events | SEC-03, 04, 08, 11, 13, AUD-01, 02, 03, 05 |
| Access-control matrix | Table-driven Go test generated from the access matrix | Every role × resource × action cell, including object-level cases | SEC-10, 11, 12, 15 |
| Security (OWASP) | Manual test cases, scripted requests, OWASP ZAP baseline scan, `govulncheck`, `npm audit` | OWASP Top 10 categories (section 4) | SEC-20 to 25 |
| Load | `k6` or `vegeta` | Behaviour under load and rate limiting | SEC-03, 23 |
| Frontend | Component tests (Vitest) and manual checks | Escaping of user content, role-dependent UI, no tokens in logs | SEC-21, PRV-04 |
| Acceptance | Manual walkthrough of key scenarios with the supervisor | Parent, specialist and admin scenarios end to end | all "M" requirements |

## 3. Metrics, targets and measurement conditions

Targets are proposals for the prototype and will be reviewed with the supervisor.

| Metric | Target | How it is measured | Conditions |
| --- | --- | --- | --- |
| Access-matrix test pass rate | 100% of cells | Table-driven test result in CI | Fresh database with seed data for two families, two specialists, one admin |
| Cross-family access attempts blocked | 100% (0 successful reads or writes), all answered with 404 | Integration tests against every endpoint that takes a child or entry ID | Authenticated as a parent of another family |
| `specialists_only` leakage to guardians | 0 occurrences | Visibility tests over list, search and export | Seed data with mixed visibility |
| Admin access to diary text | 0 occurrences | Tests on every admin endpoint, including export and deletion | Seed data with diary entries of both visibility types |
| Audit completeness (data access) | 100% of create/read/update/delete/export actions on child-related resources produce exactly one event (AUD-01) | Count of events versus actions in an integration scenario | Includes denied attempts |
| Audit completeness (authentication) | 100% of login success, login failure, logout, refresh and lockout actions produce exactly one event (AUD-05) | Count of events versus scripted actions | Includes unknown-email attempts |
| Login response uniformity | Identical status and body; difference of median response times below 50 ms (SEC-04) | 200 attempts with an unknown email and 200 with a wrong password | Same account state, single node |
| Open high-severity findings at release | 0 | OWASP checklist, ZAP baseline, dependency scans | Run on the release candidate |
| Unit and integration line coverage of `internal/service` and `internal/transport/http` | 70% or more (security-critical packages 85% or more) | `go test -cover` | CI on main |
| Authentication endpoint rate limit | 429 after the configured limit (default 5 failures per 15 minutes) | Scripted login attempts | Single IP, distributed accounts, and unknown emails (lockout must behave identically for existing and unknown emails, SEC-04) |
| API latency under normal load | p95 below 300 ms for read endpoints | `k6` run | 50 concurrent virtual users for 5 minutes on a single-node test deployment, audit of reads enabled |
| Error rate under normal load | below 1% | `k6` run | Same as above |
| Behaviour under overload | No crash; excess requests answered with 429 or 503 | `k6` ramp to 5× normal load | Same hardware, reported with the hardware description |
| Cost of auditing under load | Reported, not a pass/fail target: p95 and error rate with audit of reads enabled versus disabled | Two `k6` runs with identical scenario | Every read writes an audit row, so this is the most likely bottleneck |

**Basis for the load figures (assumption, to be confirmed with the supervisor).** The pilot is assumed to serve about 200 registered families and 20 specialists, with at most 20 users active at the same moment. 50 virtual users is therefore 2.5× the assumed peak, and the 5× overload ramp (250 virtual users) is 12.5× the assumed peak. Targets are for a single-node test deployment and are compared relative to the naive baseline (section 5), not presented as production capacity.

**Audit overhead.** The audit log writes one row for every read (AUD-01), which adds database load. The load test is therefore run with auditing enabled, and the extra run with auditing disabled shows its cost. The result decides whether read events may be written asynchronously. Denied attempts and all writes stay synchronous so that no event is lost.

## 4. OWASP Top 10 (2021) test plan

| OWASP category | Planned checks | Linked threats |
| --- | --- | --- |
| A01 Broken Access Control | Access matrix test, IDOR tests on every ID parameter, revoked-assignment test, force-browsing of admin endpoints, 404 versus 403 rule check | T3, T4, T5, T6 |
| A02 Cryptographic Failures | Password hash inspection, TLS configuration scan, encryption-at-rest configuration review, check that secrets are not in the repository | T13, T14, T17 |
| A03 Injection | SQL injection payloads on all inputs, XSS payloads in diary text and titles | T7, T8 |
| A04 Insecure Design | Review of the threat model and misuse cases against the implemented flows | all |
| A05 Security Misconfiguration | Security headers, CORS, error messages without internals, default accounts absent | T9, T12 |
| A06 Vulnerable and Outdated Components | `govulncheck`, `npm audit` in CI | T16 |
| A07 Identification and Authentication Failures | Brute force, token tampering and expiry, `alg=none`, logout and refresh revocation, refresh token reuse, login response uniformity | T1, T2, T18 |
| A08 Software and Data Integrity Failures | Dependency lock files committed, CI checks, audit log immutability test | T11 |
| A09 Security Logging and Monitoring Failures | Audit completeness tests (data access and authentication), log content inspection | T10, T12 |
| A10 Server-Side Request Forgery | Confirm the API makes no server-side requests to user-supplied URLs (not applicable by design; documented) | n/a |

## 5. Baseline for comparison

The project aim is to show that the platform gives measurably better protection than a straightforward implementation without the security design. Two baselines are used:

1. **Naive baseline:** a minimal version of the API with authentication only (a valid token grants access to any child ID, no visibility filter, no rate limit). The same test suite is run against it to show which tests fail. This gives a concrete before/after comparison in the final report.
2. **Reference standard:** the OWASP Application Security Verification Standard (ASVS) Level 1 checklist. The share of applicable Level 1 controls satisfied is reported for the release candidate.

## 6. Test data and environment

- Synthetic data only (PRV-08): two families, three children, two specialists with active and ended assignments, one admin, diary entries of both visibility types.
- Isolated test database created and destroyed per test run, with migrations from `/backend/migrations`.
- CI pipeline runs unit, integration and dependency scans on each pull request; security and load tests run on the release candidate before each pre-defense tag.

## 7. Defect handling and reporting

- Each finding is recorded as a task with severity (high, medium, low), affected requirement, reproduction steps and a link to the fixing pull request.
- High severity findings block the release; medium findings need an accepted mitigation or a documented residual risk.
- A verification report (`/docs/security-report-<tag>.md`) is produced for each pre-defense tag and lists test results, metrics versus targets, open findings and known limitations.

## 8. Schedule

| Milestone | Verification activity |
| --- | --- |
| Pre-defense 1 (24 October) | This plan, requirements, access matrix and threat model; test skeleton and the first RBAC unit tests where code exists |
| After Access and Profiles backend milestones | Authentication, RBAC and IDOR integration tests |
| After Support and Diary milestones | Visibility, audit and history tests |
| Quality and launch milestone | Full OWASP pass, ZAP baseline, load tests, verification report, final-pre-defense tag |

## 9. Risks to verification

| Risk | Mitigation |
| --- | --- |
| Backend features arrive late, leaving little time to test | Write tests from the access matrix first (test-first for RBAC); agree on API contracts early with the backend developer |
| Load test hardware differs from production | Report the hardware used and compare relative to the baseline rather than claiming absolute capacity |
| Legal interpretation (health data, consent) is uncertain | Record assumptions and confirm with the supervisor |
| Audit of every read slows the API | Measure the cost (section 3); if needed write read events asynchronously, keeping denied and write events synchronous |

## 10. Revision notes (revision 2)

- Rate-limit test now includes unknown emails (SEC-03 and SEC-04 consistency); frontend row also references PRV-04.

- Audit completeness split into data access (AUD-01) and authentication (AUD-05) metrics.
- SEC-04 tolerance turned into a measurable metric (50 ms, 200 attempts per case).
- Added metrics for admin access to diary text (SEC-15) and for 404 on cross-family access (SEC-11).
- Load figures justified by an explicit usage assumption; audit overhead added as a separate measurement and as a risk.

## 11. Results — 25 October 2026 (SSP-58)

Measured with synthetic data on the development box and isolated containers. The full `make check` wrapper is still blocked by the host's missing Go toolchain; its component gates were run separately.

| Check | Command | Result |
|---|---|---|
| Full gate | `make check` | Not completed on this host: `go`, `gofmt`, and `golangci-lint` are absent; component gates were run separately |
| Go suite | `docker run --rm -v '/Users/luun/Downloads/demo v1:/work' -w /work/backend -e 'TEST_DATABASE_URL=postgres://postgres@host.docker.internal:55437/postgres?sslmode=disable' golang:1.26 go test ./...` | All Go unit, integration, and contract packages passed, including password-reset expiry/reuse and route-contract tests |
| Backend coverage | `backend/scripts/covergate.sh` | Not rerun after these changes; previous baseline was 85.7 % total |
| Access matrix | `TestAccessMatrix100PercentCells` | every role × relation × action cell asserted |
| IDOR | `TestCrossFamilyIDORReturns404OnEveryEndpoint` | 404 on all child-scoped endpoints |
| Frontend unit | `npm test` (Vitest) | 6 files, 15 tests passed |
| E2E + a11y | `VITE_API_PROXY=http://127.0.0.1:18084 make e2e` (Playwright + axe, WCAG 2.1 AA tags) | 6/6 passed; includes reset routes, parent deletion of own entry, specialist delete-action absence, and mobile drawer; 0 axe violations in checked screens |
| Lighthouse (mobile preset, `vite preview`) | `npx lighthouse@12` | `/`: perf 91, a11y 100, best-practices 96; `/login`: 90 / 100 / 96; LCP 2.7 s (simulated throttling) |
| Bundle | `npm run size:bundle` | initial JS 186.9 KB gzip (budget 200 KB) |
| Load | `k6 run loadtest/k6-diary.js` | 50 VU × 60 s: p95 6.19 ms, 0 % errors, 13 103 requests |
| Dependencies | Go `govulncheck ./...` in Go 1.26 container; `npm --prefix frontend audit --audit-level=high` | 0 called Go vulnerabilities (one uncalled module finding); 0 npm vulnerabilities |
| Go lint | `golangci-lint v2.15.0 run ./...` in Go 1.26 container | 0 issues |
| Go vet | `go vet ./...` in Go 1.26 container | Passed |
| `make vuln` wrapper | `make vuln` | Not completed on this host because Go is absent; Go and npm scanners were run separately |
| Secrets | `gitleaks git` / `gitleaks dir` | no leaks (2 expired local test JWTs in history allow-listed in `.gitleaksignore` with reason) |
| Workflows | `actionlint` | no findings |
| Compose images and services | `docker compose -p sspview --env-file /tmp/demo-v1-verify-hardening2.env up --build -d` | API and web images rebuilt; migrations completed; API and web healthy; synthetic seed available |
| Audit role permissions | `TestApplicationRoleCannotTamperWithAuditRows` | `ssp_app` cannot UPDATE or DELETE audit rows; retention is through the restricted function |
| Password reset | `TestPasswordResetUniformResponseSingleUseAndExpiry`, `TestPasswordResetTokenIsOpaqueAndHashed` | Generic response for known/unknown accounts; hashed one-hour token; expiry/reuse rejected; refresh sessions revoked |
| Diary deletion UI | Playwright parent/specialist diary scenario | Parent can delete own entry with confirmation; specialist has no delete action |
| Password-reset frontend | Playwright guest axe checks on request and token form | Both routes accessible; HTTP calls use the generated OpenAPI client |
| Trivy, ZAP baseline | Image scan and nightly scan | **Not verified** |

The naive baseline comparison (section 5) and the usability test with 5 users (SF-08) are not part of this run and remain open.

At the time of this 25 October run, the API, PostgreSQL integration, role and relationship authorization, diary visibility, audit trail, diary-entry deletion, retention command, password reset, audit-role restrictions, and frontend-on-API were implemented. Field-level encryption, Trivy/ZAP execution, legal review, and usability testing were still open then. In local Compose mode only, reset links are written to API logs; keep that flag off for deployments with real accounts until a mail provider is configured.

## 12. Follow-up verification (26 October 2026, SSP-61)

The following CI checks use synthetic data. The ZAP workflow was manually dispatched against `main` at commit `062f02d13026fbe75c6af60e622caccec28561df`.

| Check | Result |
|---|---|
| `make check` | Passed in the PR #4 quality job; [CI run #40246060377](https://github.com/Luun0/demo-v1/actions/runs/40246060377) completed 3/3 green. |
| Trivy image scan | Passed for both API and web images, 0 HIGH/CRITICAL findings; [CI run #40242594381](https://github.com/Luun0/demo-v1/actions/runs/40242594381). |
| ZAP baseline | [Workflow run #40245962871](https://github.com/Luun0/demo-v1/actions/runs/40245962871) completed successfully: 0 FAIL, 5 WARN, 63 PASS. The workflow uses `fail_action: false`, so a successful job does not mean the warnings are resolved. |

ZAP warning categories reported: Information Disclosure - Suspicious Comments (1), Storable and Cacheable Content (8), Permissions Policy Header Not Set (5), Modern Web Application (3), and Cross-Origin-Embedder-Policy Header Missing or Invalid (9). Counts are alert instances/URLs as reported by ZAP. These warnings have not been triaged; this result is not a claim that the application has no security findings. Run and review another baseline before release.

The naive baseline comparison (section 5), usability test with 5 users (SF-08), legal review, and PRV-03a field-level encryption remain open. Shared-store rate limiting is optional and was not implemented; the configured default remains the in-memory limiter.

## 13. Test and coverage hardening (26 October 2026, SSP-64)

Added table-driven boundary checks for diary validation (including the 10,000-rune Unicode limit), default-deny cases for unknown policy actions, and integration assertions that `%`, `_`, and `\\` in diary searches are treated literally. Added a PostgreSQL benchmark with 1,000 synthetic diary entries; setup is outside the timed section.

| Check | Command | Result |
|---|---|---|
| Full quality gate | `make check` | Passed: total Go coverage 85.0%; policy 97.6%, application 86.1%, security 94.6%, HTTP adapter 82.7%; initial JS gzip 186.9 KB / 200 KB. |
| Focused tests | `go test ./internal/domain/policy ./internal/application` | Passed. |
| Diary visibility tests | `go test ./tests/integration -run 'Test(DiarySearchTreatsLikeWildcardsLiterally|SpecialistsOnlyNeverReachesParent)$' -count=1` | Passed against a fresh migrated PostgreSQL database; synthetic notes only. |
| Diary list benchmarks | `go test ./tests/integration -run '^$' -bench '^BenchmarkDiaryListVisibilityFilter$' -benchmem` | Passed on Apple M2. guardian shared page: 118,446 ns/op; specialist all page: 236,811 ns/op; guardian search: 432,104 ns/op; specialist private filter: 228,836 ns/op. These are local diagnostic measurements, not a cross-machine performance claim. |

The fresh base measurement was 84.9%, increasing to 85.0% after these focused tests. The older ~85.7% figure in section 11 was not reproduced by the current gate; it has not been presented as achieved. The additional tests protect security and validation edge cases, while the suite remains under its configured gates. No runtime, API, OpenAPI, route, or access-matrix change was made.

## 14. Read-path index review (26 October 2026, SSP-65)

Added `diary_entries_child_order_idx(child_id, observed_at DESC, id DESC)` for the specialist diary list, which returns shared and specialists-only entries and needs an ordered page across both visibility values. The existing `diary_entries_access_idx(child_id, visibility, observed_at DESC, id DESC)` continues to serve the parent-only shared list. The existing `support_sessions_child_idx(child_id, starts_at)` serves the session list; its plan still sorts the small child-scoped result by `starts_at, id`, so no duplicate session index was added. The migration intentionally makes no route or access-policy change.

`TestDiaryAndSessionListPlansUseSupportingIndexes` builds 100,000 synthetic diary rows across 1,000 synthetic children and 10,000 synthetic sessions, runs `ANALYZE`, and records `EXPLAIN (ANALYZE, BUFFERS)` plans. The parent plan used `diary_entries_access_idx`; the specialist plan used `diary_entries_child_order_idx` without a sort; the session plan used `support_sessions_child_idx` for child filtering and retained a small in-memory sort.

| Check | Command | Result |
|---|---|---|
| Full quality gate | `make check` | Passed: 85.0% total Go coverage, contract/layer/file-size gates green, frontend tests and build passed; initial JS gzip 186.9 KB / 200 KB. |
| Explain plans | `go test ./tests/integration -run '^TestDiaryAndSessionListPlansUseSupportingIndexes$' -count=1 -v` | Passed on PostgreSQL 17 with the plan paths described above. |
| Paired load test, no index | `BASE_URL=http://127.0.0.1:18095 k6 run --quiet loadtest/k6-diary.js` | 50 VU / 60 s, 207.0 req/s, p95 23.94 ms, 0% failed requests, 15,400/15,400 checks. |
| Paired load test, index present | `BASE_URL=http://127.0.0.1:18095 k6 run --quiet loadtest/k6-diary.js` | 50 VU / 60 s, 207.5 req/s, p95 22.84 ms, 0% failed requests, 15,400/15,400 checks. |

Both paired tests used the same local API binary, PostgreSQL 17 database, synthetic data and machine; only the diary ordering index differed. The previously recorded 6.19 ms p95 was not reproduced under this run's conditions, so it remains historical and is not presented as a directly comparable baseline. The index's measured paired result did not regress p95.
