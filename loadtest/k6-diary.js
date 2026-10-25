// k6 load test: p95 < 300 ms at 50 VUs, errors < 1% (master prompt Phase 3).
// Run: BASE_URL=http://127.0.0.1:8080 k6 run loadtest/k6-diary.js
import http from 'k6/http'
import { check, sleep } from 'k6'

const BASE = __ENV.BASE_URL || 'http://127.0.0.1:8080'
const PASSWORD = __ENV.SEED_PASSWORD || 'Synthetic-Demo-2026'

export const options = {
  scenarios: { steady: { executor: 'constant-vus', vus: 50, duration: __ENV.DURATION || '60s' } },
  thresholds: { http_req_duration: ['p(95)<300'], http_req_failed: ['rate<0.01'] },
}

function login(email) {
  const r = http.post(`${BASE}/api/v1/auth/login`, JSON.stringify({ email, password: PASSWORD }), { headers: { 'Content-Type': 'application/json' } })
  return r.json('access_token')
}

export function setup() {
  const parent = login('parent1@example.test')
  const specialist = login('specialist1@example.test')
  const children = http.get(`${BASE}/api/v1/children`, { headers: { Authorization: `Bearer ${parent}` } }).json()
  return { parent, specialist, child: children[0].id }
}

export default function (d) {
  const asParent = __VU % 2 === 0
  const h = { headers: { Authorization: `Bearer ${asParent ? d.parent : d.specialist}`, 'Content-Type': 'application/json' } }
  const list = http.get(`${BASE}/api/v1/children/${d.child}/diary-entries?limit=20`, h)
  check(list, { 'list 200': (r) => r.status === 200 })
  if (asParent) {
    check(list, { 'parent never sees specialists_only': (r) => !r.body.includes('specialists_only') })
  }
  check(http.get(`${BASE}/api/v1/children/${d.child}/sessions`, h), { 'sessions 200': (r) => r.status === 200 })
  if (__ITER % 5 === 0) {
    const body = JSON.stringify({ type: 'observation', visibility: 'shared', content: 'Synthetic load-test note', observed_at: new Date().toISOString() })
    check(http.post(`${BASE}/api/v1/children/${d.child}/diary-entries`, body, h), { 'create 201': (r) => r.status === 201 })
  }
  sleep(0.5)
}

// Never persist tokens from setup() into the exported summary.
export function handleSummary(data) {
  delete data.setup_data
  return { 'loadtest/k6-summary.json': JSON.stringify(data, null, 4), stdout: '\n' }
}
