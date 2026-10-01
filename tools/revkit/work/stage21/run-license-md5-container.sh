#!/bin/bash
set -euo pipefail

exe=/work/stage21/probe-license-md5.exe
log=/work/stage21/license-md5-vectors-api.log
if [ -e "$log" ]; then
  echo "refusing to overwrite the license MD5 vector capture" >&2
  exit 2
fi
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 60s \
  wine "$exe" > "$log" 2>&1
case_count=$(grep -c '^MD5_VECTOR name=' "$log")
if [ "$case_count" -ne 10 ]; then
  echo "expected 10 MD5 vectors, captured $case_count" >&2
  exit 1
fi
grep -q '^MD5_VECTOR_SUMMARY cases=10 failures=0$' "$log"
