#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
xpid=
backup_input=/tmp/vtspeak-csv-iscsv-input1.txt
backup_output=/tmp/vtspeak-csv-iscsv-output.wav
restore() {
  cp "$backup_input" "$work/input1.txt"
  cp "$backup_output" "$work/output.wav"
  cmp -s "$backup_input" "$work/input1.txt"
  cmp -s "$backup_output" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

cp "$work/input1.txt" "$backup_input"
cp "$work/output.wav" "$backup_output"
trap restore EXIT
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-csv-iscsv-matrix.log" 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-csv-iscsv-matrix.gdb" > "$probe/csv-iscsv-matrix-api.log" 2>&1

if [ "$(grep -c '^CSV_ISCSV case=' "$probe/csv-iscsv-matrix-api.log")" -ne 33 ]; then
  echo "did not capture all 33 CSV IsCsv cases" >&2
  exit 1
fi
if [ "$(grep -c '^CSV_PARSE_CAPACITY ' "$probe/csv-iscsv-matrix-api.log")" -ne 4 ]; then
  echo "did not capture all 4 direct parser capacity checks" >&2
  exit 1
fi
grep -Fq 'CSV_PARSE_CAPACITY input_fields=101 low_ax=-4 stored_fields=100' "$probe/csv-iscsv-matrix-api.log"
grep -Fq 'CSV_ISCSV case=field_capacity fields=101 expected=100 limit=512 low_ax=1' "$probe/csv-iscsv-matrix-api.log"
grep -Fq 'CSV_ISCSV case=field_capacity fields=101 expected=101 limit=512 low_ax=0' "$probe/csv-iscsv-matrix-api.log"
