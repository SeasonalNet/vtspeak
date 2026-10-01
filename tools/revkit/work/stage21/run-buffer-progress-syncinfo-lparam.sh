#!/bin/bash
set -euo pipefail

probe=/work/stage21
work=/work/stage21/sandbox/stage5
log="$probe/buffer-progress-syncinfo-lparam-api.log"
xlog="$probe/xvfb-buffer-progress-syncinfo-lparam.log"
xpid=

for output in "$log" "$xlog"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 2
  fi
done
test -f "$work/input1.txt"
test -f "$work/output.wav"

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xlog" 2>&1 &
xpid=$!
restore() {
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-buffer-progress-syncinfo-lparam.gdb" > "$log" 2>&1

if [ "$(grep -c 'source_matches=1' "$log")" -ne 4 ]; then
  echo "expected four PostMessage wParam/progress-context matches in $log" >&2
  exit 1
fi
if [ "$(grep -c 'lparam_matches_row_end=1' "$log")" -ne 4 ]; then
  echo "expected four PostMessage lParam/SyncInfo row-end matches in $log" >&2
  exit 1
fi
grep -q 'PROGRESS_COMPLETE raw=1 polls=3 notifications=4 total_bytes=183224' "$log"
grep 'PROGRESS_POST\|PROGRESS_RETURN\|PROGRESS_COMPLETE' "$log"
