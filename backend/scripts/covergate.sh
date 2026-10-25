#!/usr/bin/env bash
# Coverage gate: total >= 70%, security-critical packages >= 85%.
set -euo pipefail
cd "$(dirname "$0")/.."
PROFILE=${1:-coverage.out}
go test -race -count=1 -coverpkg=./internal/... -coverprofile="$PROFILE" ./... >/dev/null
total=$(go tool cover -func="$PROFILE" | awk '/^total:/ {sub("%","",$3); print $3}')
echo "total coverage: ${total}%"
fail=0
awk "BEGIN{exit !(${total} < 70)}" && { echo "FAIL: total < 70%"; fail=1; }
# per-package statement coverage from the profile
pkgcov() {
  awk -v pkg="$1" 'NR>1 { split($1, f, ":"); file=f[1]; dir=file; sub(/\/[^\/]*$/, "", dir);
      if (dir != pkg) next; n[$1]=$2; if ($3 > c[$1]) c[$1]=$3; if (!($1 in c)) c[$1]=$3 }
    END { for (k in n) { tot+=n[k]; if (c[k]>0) cov+=n[k] } if (tot==0) print 0; else printf "%.1f", 100*cov/tot }' "$PROFILE"
}
for p in domain/policy application infrastructure/security adapter/httpapi; do
  v=$(pkgcov "github.com/AknurOE/sen-support-platform/internal/$p")
  min=85; [ "$p" = adapter/httpapi ] && min=70
  echo "  internal/$p: ${v}% (min ${min}%)"
  awk "BEGIN{exit !(${v} < ${min})}" && { echo "FAIL: internal/$p < ${min}%"; fail=1; }
done
exit $fail
