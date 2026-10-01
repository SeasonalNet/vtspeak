#!/bin/bash
set -euo pipefail

for selector in 1 2; do
  for flag in 0 1; do
    /bin/bash /work/stage21/run-highlight-ex-two-byte-mark.sh \
      "$selector" "$flag" classes-replay
  done
done
