#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage16
log="$probe/load-ext-memory-license-state-api.log"
xvfb_log="$probe/xvfb-load-ext-memory-license-state.log"
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
cd "$work"
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  wine "$probe/probe-load-ext-memory-license.exe" > "$log" 2>&1
sed -i 's/[[:blank:]]*$//' "$log"
