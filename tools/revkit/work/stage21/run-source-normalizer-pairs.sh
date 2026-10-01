#!/bin/bash
set -euo pipefail

probe=/work/stage21
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  wine "$probe/probe-source-normalizer-pairs.exe" \
  > "$probe/source-normalizer-pairs-api.log" 2>&1
tail -n 1 "$probe/source-normalizer-pairs-api.log"
