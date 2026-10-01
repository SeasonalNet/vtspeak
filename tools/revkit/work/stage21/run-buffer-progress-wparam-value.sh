#!/bin/bash
set -euo pipefail

probe=/work/stage21
work=/work/stage21/sandbox/stage5
log="$probe/buffer-progress-wparam-value-api.log"
xlog="$probe/xvfb-buffer-progress-wparam-value.log"
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
  < "$probe/trace-buffer-progress-wparam-value.gdb" > "$log" 2>&1

if [ "$(grep -c 'source_matches=1' "$log")" -ne 4 ]; then
  echo "expected four PostMessage wParam/progress-context matches in $log" >&2
  exit 1
fi
grep -q 'PROGRESS_COMPLETE raw=1 polls=3 notifications=4 total_bytes=183224' "$log"
grep -q 'PROGRESS_CONTEXT_INIT ' "$log"
grep -q 'PROGRESS_HELPER_STRING ' "$log"
helper_bytes=$(sed -n 's/.*PROGRESS_HELPER_STRING .* bytes=\([0-9][0-9]*\) text=.*/\1/p' "$log")
test -n "$helper_bytes"
helper_hex=$(printf '%x' "$helper_bytes")
grep -q "PROGRESS_CONTEXT_WRITE eip=0x10026ad6 old=0 new=0x$helper_hex$" "$log"
grep 'PROGRESS_POST\|PROGRESS_RETURN\|PROGRESS_COMPLETE' "$log"
