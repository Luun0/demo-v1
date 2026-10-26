# Improvement queue progress

| Task | Change / gap closed | Verification |
|---|---|---|
| 1. Code health sweep | Reviewed `docs/openapi.yaml`, Go and frontend source for TODO/FIXME markers and oversized files. No OpenAPI markers or hand-written source files over 252 lines were found; generated `schema.gen.ts` is excluded from the source-file limit. Error wrapping uses `%w` at underlying-error boundaries; production handlers/use cases propagate request contexts. No behavior changes were warranted. | `make check` — passed; coverage 84.9%, initial JS gzip 186.9 KB |

## Deferred tasks

Tasks 2–7 remain queued and will be recorded here when completed. No coverage, benchmark, query-plan, load-test, or evaluation result is claimed in advance.
