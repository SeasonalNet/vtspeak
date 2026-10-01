#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
log="$probe/executable-heap-smoke-api.log"
xvfb_log="$probe/xvfb-executable-heap-smoke.log"
input_backup=/tmp/executable-heap-smoke-input1.txt
output_backup=/tmp/executable-heap-smoke-output.wav
xpid=
restore() {
  cp "$input_backup" "$work/input1.txt"
  cp "$output_backup" "$work/output.wav"
  cmp -s "$input_backup" "$work/input1.txt"
  cmp -s "$output_backup" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
for path in "$log" "$xvfb_log" "$input_backup" "$output_backup"; do
  if [ -e "$path" ]; then
    echo "refusing to overwrite $path" >&2
    exit 3
  fi
done
cp "$work/input1.txt" "$input_backup"
cp "$work/output.wav" "$output_backup"
trap restore EXIT
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 60s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-executable-heap-smoke.gdb" > "$log" 2>&1
grep -q '^HEAP_CODE_RETURNED ' "$log"
grep -q '\[Inferior .*exited normally\]' "$log"
grep '^HEAP_CODE_RETURNED ' "$log"
