#!/bin/bash
set -euo pipefail

probe=/work/stage21
work=/work/stage21/sandbox/stage5
trace="$probe/trace-csv-iscsv-capacity-v2.gdb"
log="$probe/csv-iscsv-capacity-v2-api.log"
xvfb_log="$probe/xvfb-csv-iscsv-capacity-v2.log"
backup_input=/tmp/vtspeak-csv-iscsv-capacity-v2-input1.txt
backup_output=/tmp/vtspeak-csv-iscsv-capacity-v2-output.wav
xpid=

restore() {
  if [ -f "$backup_input" ]; then
    cp "$backup_input" "$work/input1.txt"
    cmp -s "$backup_input" "$work/input1.txt"
    rm -f "$backup_input"
  fi
  if [ -f "$backup_output" ]; then
    cp "$backup_output" "$work/output.wav"
    cmp -s "$backup_output" "$work/output.wav"
    rm -f "$backup_output"
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
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" > "$log" 2>&1

if [ "$(grep -c '^CSV_CAPACITY fields=' "$log")" -ne 35 ] || \
   [ "$(grep -c '^CSV_PARSE_CAPACITY fields=' "$log")" -ne 35 ]; then
  echo "did not capture all 35 IsCsv/parser capacity cases" >&2
  exit 1
fi
grep -Fq 'CSV_CAPACITY fields=99 bytes=197 bound=198 expected_100=0 expected_fields=1 expected_next=0' "$log"
grep -Fq 'CSV_CAPACITY fields=100 bytes=199 bound=200 expected_100=1 expected_fields=1 expected_next=0' "$log"
grep -Fq 'CSV_CAPACITY fields=4096 bytes=8191 bound=8192 expected_100=1 expected_fields=0 expected_next=0' "$log"
grep -Fq 'CSV_PARSE_CAPACITY fields=4096 low_ax=-4 stored_fields=100' "$log"
