# Load test (k6)

Scenario `k6-diary.js`: 50 constant VUs for 60 s; parents and specialists list the diary (limit 20), list sessions, and every 5th iteration create a shared entry. Thresholds: p95 < 300 ms, errors < 1 %.

## Result — 2026-10-24 23:47 (UTC+5), box: local API binary + local PostgreSQL 17, single instance

| Metric | Value |
|---|---|
| Requests | 13 103 (216 req/s) |
| http_req_failed | 0.00 % |
| p95 latency | 6.19 ms |
| p90 / median / max | 4.98 ms / 2.41 ms / 115.96 ms |
| Checks (incl. "parent never sees specialists_only") | 16 075 / 16 075 passed |

Raw summary: `k6-summary.json`. Note: `API_RATE_PER_MINUTE` was raised for the run because all VUs share one IP; numbers come from a dev box, not a production-like environment.

## Run

```bash
# against a throwaway, seeded database (the run creates ~1 200 entries)
API_RATE_PER_MINUTE=1000000 /path/to/api serve &
k6 run loadtest/k6-diary.js   # handleSummary writes k6-summary.json without setup tokens
```
