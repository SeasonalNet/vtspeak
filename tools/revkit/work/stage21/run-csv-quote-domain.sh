#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
xpid=
restore() {
  cp /tmp/csv-quote-domain-input1.txt "$work/input1.txt"
  cp /tmp/csv-quote-domain-output.wav "$work/output.wav"
  cmp -s /tmp/csv-quote-domain-input1.txt "$work/input1.txt"
  cmp -s /tmp/csv-quote-domain-output.wav "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
cp "$work/input1.txt" /tmp/csv-quote-domain-input1.txt
cp "$work/output.wav" /tmp/csv-quote-domain-output.wav
trap restore EXIT
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-csv-quote-domain.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 600s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-csv-quote-domain.gdb" > "$probe/csv-quote-domain-api.log" 2>&1
case_count=$(grep -c '^CSV_QDOMAIN case=' "$probe/csv-quote-domain-api.log")
if [ "$case_count" -ne 1093 ]; then
  echo "expected 1093 quote/comma cases, captured $case_count" >&2
  exit 1
fi
grep '^CSV_QUOTE_DOMAIN' "$probe/csv-quote-domain-api.log"
grep -q 'CSV_QUOTE_DOMAIN cases=1093 low_ax_failures=0' "$probe/csv-quote-domain-api.log"
grep -q '\[Inferior .*exited normally\]' "$probe/csv-quote-domain-api.log"
