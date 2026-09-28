#!/bin/bash
set -euo pipefail

GDB_SCRIPT=/work/stage20/trace-natural-row6-only.gdb \
  EVIDENCE_DIR=/work/corpus-parity/stage20/natural-row6-only-corrected \
  /bin/bash /work/stage20/run-adapted-forced-pah0.sh
GDB_SCRIPT=/work/stage20/trace-natural-rows5-6-only.gdb \
  EVIDENCE_DIR=/work/corpus-parity/stage20/natural-rows5-6-only \
  /bin/bash /work/stage20/run-adapted-forced-pah0.sh
