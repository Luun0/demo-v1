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

## Paired index check — 2026-10-26

Compared the same local API binary and PostgreSQL 17 database with the new diary ordering index absent/present. Each run used the existing 50-VU, 60-second scenario with synthetic seeded data and `API_RATE_PER_MINUTE=1000000`.

| Setup | Requests/s | p95 | Failed requests | Checks |
|---|---:|---:|---:|---:|
| Without `diary_entries_child_order_idx` | 207.0 | 23.94 ms | 0 % | 15,400 / 15,400 |
| With `diary_entries_child_order_idx` | 207.5 | 22.84 ms | 0 % | 15,400 / 15,400 |

These paired results show no regression from the added index in this local run. The earlier 6.19 ms p95 was not reproduced under the current run conditions and remains a separate historical measurement; do not compare the absolute values as if the environments were identical. The raw summary file contains the final indexed run.

## Run

```bash
# against a throwaway, seeded database (the run creates ~1 200 entries)
API_RATE_PER_MINUTE=1000000 /path/to/api serve &
k6 run loadtest/k6-diary.js   # handleSummary writes k6-summary.json without setup tokens
```
