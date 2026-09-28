#!/bin/bash
set -euo pipefail

probe=/work/stage16
mode=${1:?usage: run-load-ext-error2-static-container.sh control|force-allocation|force-speaker}
case "$mode" in
  control|force-allocation|force-speaker) ;;
  *) echo "unsupported probe mode: $mode" >&2; exit 2 ;;
esac
log="$probe/load-ext-error2-static-v5-$mode-api.log"
xvfb_log="$probe/xvfb-load-ext-error2-static-v5-$mode.log"
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
cd /work
if [ "$mode" = control ]; then
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
    wine "$probe/probe-load-ext-error2-static.exe" > "$log" 2>&1
else
  trace="$probe/trace-load-ext-error2-allocation-static.gdb"
  if [ "$mode" = force-speaker ]; then
    trace="$probe/trace-load-ext-error2-speaker-allocation-static.gdb"
  fi
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
    winedbg --gdb "$probe/probe-load-ext-error2-static.exe" \
    < "$trace" > "$log" 2>&1
fi
sed -i 's/[[:blank:]]*$//' "$log"
