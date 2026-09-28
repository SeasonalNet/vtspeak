#!/bin/bash
set -euo pipefail

GDB_SCRIPT=/work/stage20/trace-natural-tail-only.gdb \
  EVIDENCE_DIR=/work/corpus-parity/stage20/natural-tail-only \
  /bin/bash /work/stage20/run-adapted-forced-pah0.sh
GDB_SCRIPT=/work/stage20/trace-natural-last-row-only.gdb \
  EVIDENCE_DIR=/work/corpus-parity/stage20/natural-last-row-only \
  /bin/bash /work/stage20/run-adapted-forced-pah0.sh
