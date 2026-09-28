#!/bin/bash
set -euo pipefail

probe=/work/stage16
log="$probe/load-ext-slot4-filetrace-v3-api.log"
xvfb_log="$probe/xvfb-load-ext-slot4-filetrace-v3.log"
raw="/tmp/vtspeak-load-ext-slot4-filetrace-$$.log"
xpid=
cleanup() {
  rm -f "$raw"
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
set +e
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=+file timeout 180s \
  wine "$probe/probe-load-ext-state.exe" 4 > "$raw" 2>&1
wine_status=$?
set -e
printf 'WINE_PROCESS_STATUS=%d\n' "$wine_status" > "$log"
grep -Ei 'data-jame|data-james|tree3|error|status=|^SLOT ' "$raw" >> "$log" || true
sed -i 's/[[:blank:]]*$//' "$log"
