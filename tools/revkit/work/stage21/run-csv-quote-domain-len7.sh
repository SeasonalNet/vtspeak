#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
log="$probe/csv-quote-domain-len7-api.log"
xlog="$probe/xvfb-csv-quote-domain-len7.log"
xpid=

if [ -e "$log" ] || [ -e "$xlog" ]; then
  echo "refusing to overwrite an existing length-7 quote-domain capture" >&2
  exit 2
fi

cp "$work/input1.txt" /tmp/csv-quote-domain-len7-input1.txt
cp "$work/output.wav" /tmp/csv-quote-domain-len7-output.wav
restore() {
  cp /tmp/csv-quote-domain-len7-input1.txt "$work/input1.txt"
  cp /tmp/csv-quote-domain-len7-output.wav "$work/output.wav"
  cmp -s /tmp/csv-quote-domain-len7-input1.txt "$work/input1.txt"
  cmp -s /tmp/csv-quote-domain-len7-output.wav "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT

cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xlog" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 600s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-csv-quote-domain-len7.gdb" > "$log" 2>&1
case_count=$(grep -c '^CSV_QDOMAIN7 case=' "$log")
if [ "$case_count" -ne 3280 ]; then
  echo "expected 3280 quote/comma cases through length 7, captured $case_count" >&2
  exit 1
fi
grep '^CSV_QUOTE_DOMAIN7' "$log"
grep -q 'CSV_QUOTE_DOMAIN7 cases=3280 low_ax_failures=0' "$log"
grep -q '\[Inferior .*exited normally\]' "$log"
