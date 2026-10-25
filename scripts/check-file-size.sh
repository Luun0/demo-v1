#!/usr/bin/env bash
# DoD: no hand-written source file over 300 lines in backend/internal or frontend/src.
set -euo pipefail
cd "$(dirname "$0")/.."
max=300
bad=$(find backend/internal backend/cmd frontend/src -type f \( -name '*.go' -o -name '*.ts' -o -name '*.tsx' -o -name '*.css' \) \
  ! -path '*/generated/*' ! -name '*.gen.ts' ! -path '*/locales/*' -print0 | xargs -0 wc -l | awk -v m=$max '$2 != "total" && $1 > m {print}')
if [ -n "$bad" ]; then echo "files over $max lines:"; echo "$bad"; exit 1; fi
echo "file-size check: OK (all <= $max lines)"
