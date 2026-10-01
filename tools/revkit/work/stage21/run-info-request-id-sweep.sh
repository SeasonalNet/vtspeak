#!/bin/bash
set -euo pipefail

probe=/work/stage21
work=/work/stage21/sandbox/stage5
trace="$probe/trace-info-request-id-sweep.gdb"
log="$probe/info-request-id-sweep-api.log"
xvfb_log="$probe/xvfb-info-request-id-sweep.log"
backup_input="$probe/info-request-id-sweep-original-input1.txt"
backup_output="$probe/info-request-id-sweep-original-output.wav"
xpid=

restore() {
  if [ -f "$backup_input" ]; then
    cp "$backup_input" "$work/input1.txt"
    cmp -s "$backup_input" "$work/input1.txt"
  fi
  if [ -f "$backup_output" ]; then
    cp "$backup_output" "$work/output.wav"
    cmp -s "$backup_output" "$work/output.wav"
  fi
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT

for output in "$log" "$xvfb_log" "$backup_input" "$backup_output"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
test -f "$work/input1.txt"
test -f "$work/output.wav"
cp "$work/input1.txt" "$backup_input"
cp "$work/output.wav" "$backup_output"
cp "$probe/input.txt" "$work/input1.txt"
cd "$work"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 240s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" > "$log" 2>&1
