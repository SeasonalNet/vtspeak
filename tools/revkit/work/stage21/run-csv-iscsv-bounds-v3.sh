#!/bin/bash
set -euo pipefail

probe=/work/stage21
work=/work/stage21/sandbox/stage5
trace="$probe/trace-csv-iscsv-bounds-v3.gdb"
log="$probe/csv-iscsv-bounds-v3-api.log"
xvfb_log="$probe/xvfb-csv-iscsv-bounds-v3.log"
backup_input=/tmp/vtspeak-csv-iscsv-bounds-v3-input1.txt
backup_output=/tmp/vtspeak-csv-iscsv-bounds-v3-output.wav
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

if [ "$(grep -c '^CSV_BOUND pointer=' "$log")" -ne 14 ]; then
  echo "did not capture all 14 IsCsv bound cases" >&2
  exit 1
fi
grep -Eq 'CSV_BOUND pointer=[^ ]+ limit=-2147483648 end=[^ ]+ low_ax=1 ' "$log"
grep -Eq 'CSV_BOUND pointer=[^ ]+ limit=-1024 end=[^ ]+ low_ax=0 ' "$log"
grep -Eq 'CSV_BOUND pointer=[^ ]+ limit=5 end=[^ ]+ low_ax=1 ' "$log"
grep -Eq 'CSV_BOUND pointer=[^ ]+ limit=4 end=[^ ]+ low_ax=0 ' "$log"
grep -Eq 'CSV_BOUND pointer=[^ ]+ limit=-2147483648 end=0x[0-9a-f]+ low_ax=1 ' "$log"
grep -Eq 'CSV_BOUND pointer=[^ ]+ limit=-1 end=0x[0-9a-f]+ low_ax=0 ' "$log"
grep -Eq 'CSV_BOUND pointer=[^ ]+ limit=0 end=0x[0-9a-f]+ low_ax=0 ' "$log"
grep -Eq 'CSV_BOUND pointer=[^ ]+ limit=1 end=0x[0-9a-f]+ low_ax=0 ' "$log"
grep -Eq 'CSV_BOUND pointer=[^ ]+ limit=6 end=0x[0-9a-f]+ low_ax=1 ' "$log"
grep -Eq 'CSV_BOUND pointer=[^ ]+ limit=2147483647 end=0x[0-9a-f]+ low_ax=1 ' "$log"
for endpoint in 0 0x1 0x4 0x5; do
  grep -Eq "CSV_BOUND pointer=[^ ]+ limit=-[0-9]+ end=$endpoint low_ax=0 " "$log"
done
grep -Eq 'CSV_BOUND pointer=[^ ]+ limit=-[0-9]+ end=0xffffffff low_ax=1 ' "$log"
