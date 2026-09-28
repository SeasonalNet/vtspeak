#!/bin/bash
set -euo pipefail

probe=/work/stage16
log="$probe/load-ext-slot-matrix-alias-api.log"
xvfb_log="$probe/xvfb-load-ext-slot-matrix-alias.log"
xpid=
cleanup() {
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap cleanup EXIT
for output in "$log" "$xvfb_log"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99
cd /work/stage5
for slot in 0 1 2 3 4 5; do
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
    wine "$probe/probe-load-ext-state.exe" "$slot" >> "$log" 2>&1
done
sed -i 's/[[:blank:]]*$//' "$log"
