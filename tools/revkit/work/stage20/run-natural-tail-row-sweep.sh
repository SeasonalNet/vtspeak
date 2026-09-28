#!/bin/bash
set -euo pipefail

for row_count in 5 6 7; do
  if [ "$row_count" -eq 7 ]; then
    suffix=full7
  else
    suffix="cut${row_count}"
  fi
  GDB_SCRIPT="/work/stage20/trace-natural-tail-${suffix}.gdb" \
    EVIDENCE_DIR="/work/corpus-parity/stage20/natural-tail-${suffix}" \
    /bin/bash /work/stage20/run-adapted-forced-pah0.sh
done
