#!/bin/bash
set -euo pipefail

probe=/work/stage21
exe="$probe/license-token-byte-matrix.exe"
log=/work/stage21/license-token-byte-matrix-v3-api.log

if [ -e "$log" ]; then
  echo "refusing to overwrite an existing license-token byte matrix capture" >&2
  exit 2
fi

WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all \
  timeout 120s wine "$exe" > "$log" 2>&1
grep -q '^TOKEN_MATRIX baseline=0 ' "$log"
case_count=$(grep -c '^TOKEN_BYTE ' "$log")
if [ "$case_count" -lt 192 ]; then
  echo "expected at least 192 token-byte mutation calls, captured $case_count" >&2
  exit 1
fi
grep -q '^TOKEN_MATRIX_SUMMARY ' "$log"
grep -q '^TOKEN_CASE_COMBO mask=0 result=0$' "$log"
