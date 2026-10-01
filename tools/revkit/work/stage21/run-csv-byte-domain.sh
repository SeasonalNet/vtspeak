#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
xpid=
restore() {
  cp /tmp/csv-byte-domain-input1.txt "$work/input1.txt"
  cp /tmp/csv-byte-domain-output.wav "$work/output.wav"
  cmp -s /tmp/csv-byte-domain-input1.txt "$work/input1.txt"
  cmp -s /tmp/csv-byte-domain-output.wav "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
cp "$work/input1.txt" /tmp/csv-byte-domain-input1.txt
cp "$work/output.wav" /tmp/csv-byte-domain-output.wav
trap restore EXIT
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-csv-byte-domain.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-csv-byte-domain.gdb" > "$probe/csv-byte-domain-api.log" 2>&1
grep -E '^CSV_BYTE_DOMAIN|^CSV_BYTE_FAIL' "$probe/csv-byte-domain-api.log"
grep -q 'CSV_BYTE_DOMAIN passed=256 failed=0' "$probe/csv-byte-domain-api.log"
grep -q '\[Inferior .*exited normally\]' "$probe/csv-byte-domain-api.log"
