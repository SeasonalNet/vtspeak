#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
xpid=
backup_input=/tmp/vtspeak-csv-lifecycle-input1.txt
backup_output=/tmp/vtspeak-csv-lifecycle-output.wav
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
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-csv-lifecycle.log" 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-csv-lifecycle.gdb" > "$probe/csv-lifecycle-api.log" 2>&1

if [ "$(grep -c '^CSV_LIFECYCLE case=' "$probe/csv-lifecycle-api.log")" -ne 5 ]; then
  echo "did not capture all 5 CSV lifecycle cases" >&2
  exit 1
fi
grep -Fq 'size=24 p1_fields=0,0,0,100,0' "$probe/csv-lifecycle-api.log"
grep -Fq 'delim=<,>' "$probe/csv-lifecycle-api.log"
grep -Fq 'CSV_LIFECYCLE case=exit_null completed=1' "$probe/csv-lifecycle-api.log"
grep -Fq 'CSV_LIFECYCLE case=exit_unparsed completed=1' "$probe/csv-lifecycle-api.log"
grep -Fq 'CSV_LIFECYCLE case=parse_copy low_ax=1 count=2' "$probe/csv-lifecycle-api.log"
grep -Fq 'CSV_LIFECYCLE case=exit_parsed completed=1' "$probe/csv-lifecycle-api.log"
