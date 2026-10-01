#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
xpid=
restore() {
  cp /tmp/csv-delimiter-input1.txt "$work/input1.txt"
  cp /tmp/csv-delimiter-output.wav "$work/output.wav"
  cmp -s /tmp/csv-delimiter-input1.txt "$work/input1.txt"
  cmp -s /tmp/csv-delimiter-output.wav "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
cp "$work/input1.txt" /tmp/csv-delimiter-input1.txt
cp "$work/output.wav" /tmp/csv-delimiter-output.wav
trap restore EXIT
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-csv-delimiter-field.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-csv-delimiter-field.gdb" > "$probe/csv-delimiter-field-api.log" 2>&1
grep '^CSV_DELIMITER' "$probe/csv-delimiter-field-api.log"
grep -q '\[Inferior .*exited normally\]' "$probe/csv-delimiter-field-api.log"
