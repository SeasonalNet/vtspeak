#!/bin/bash
set -euo pipefail

exe=/work/stage21/probe-syncinfo-allocation-retry.exe
log=/work/stage21/syncinfo-allocation-retry-v7-api.log
if [ -e "$log" ]; then
  echo "refusing to overwrite $log" >&2
  exit 2
fi
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 60s \
  winedbg --gdb "$exe" < /work/stage21/trace-syncinfo-allocation-retry-v7.gdb \
  > "$log" 2>&1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 60s \
  wine "$exe" >> "$log" 2>&1
grep -q '^SYNC_ALLOC_EXPORT_ENTER address=0x10026360$' "$log"
grep -q '^SYNC_ALLOC_NESTED_HELPER_ENTER size=520 return=0x100263be$' "$log"
grep -q '^SYNC_ALLOC_FORCED_TRANSIENT_NULL callsite=0x100263be size=520$' "$log"
grep -q '^SYNC_ALLOC_RETRY_RETURN callsite=0x100263be size=520 malloc_result_nonnull=1$' "$log"
grep -q '^SYNC_ALLOC_INJECTED_OBJECT nested_alloc_calls=600 valid=1 rows=600 width=65 row_array=1 first_nested=1 last_nested=1 distinct_nested=1$' "$log"
grep -q '^SYNC_ALLOC_RESULT valid=1 rows=600 width=65 row_array=1 first_nested=1 last_nested=1 distinct_nested=1$' "$log"
