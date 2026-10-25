// Fails if initial JS (entry + modulepreload) exceeds the gzip budget (NFR: <200 KB).
import { readFileSync } from 'node:fs'
import { gzipSync } from 'node:zlib'
const BUDGET = 200 * 1024
const html = readFileSync('dist/index.html', 'utf8')
const files = [...html.matchAll(/(?:src|href)="\/(assets\/[^"]+\.js)"/g)].map((m) => m[1])
let total = 0
for (const f of files) {
  const n = gzipSync(readFileSync(`dist/${f}`)).length
  total += n
  console.log(`${(n / 1024).toFixed(1).padStart(7)} KB  ${f}`)
}
console.log(`initial JS gzip: ${(total / 1024).toFixed(1)} KB (budget ${BUDGET / 1024} KB)`)
if (total > BUDGET) process.exit(1)
