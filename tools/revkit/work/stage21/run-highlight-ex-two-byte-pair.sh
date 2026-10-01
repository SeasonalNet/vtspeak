#!/bin/bash
set -euo pipefail

for flag in 0 1; do
  /bin/bash /work/stage21/run-highlight-ex-two-byte-mark.sh 0 "$flag"
done
