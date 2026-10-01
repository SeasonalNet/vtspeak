#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
xpid=
restore() {
  cp /tmp/parenthesis-record-input1.txt "$work/input1.txt"
  cp /tmp/parenthesis-record-output.wav "$work/output.wav"
  cmp -s /tmp/parenthesis-record-input1.txt "$work/input1.txt"
  cmp -s /tmp/parenthesis-record-output.wav "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
cp "$work/input1.txt" /tmp/parenthesis-record-input1.txt
cp "$work/output.wav" /tmp/parenthesis-record-output.wav
trap restore EXIT
cp "$probe/input-parenthesis-single-length3.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-parenthesis-record-consumer.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-parenthesis-record-consumer.gdb" > "$probe/parenthesis-record-consumer-api.log" 2>&1
grep -E '^(PAREN_|#)' "$probe/parenthesis-record-consumer-api.log" || true
grep -q '\[Inferior .*exited normally\]' "$probe/parenthesis-record-consumer-api.log"
cp "$work/output.wav" "$probe/parenthesis-record-consumer.wav"
