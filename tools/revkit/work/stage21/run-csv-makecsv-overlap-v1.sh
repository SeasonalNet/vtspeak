#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
run_id=${1:-v1}
if [[ ! "$run_id" =~ ^[A-Za-z0-9_-]+$ ]]; then
  echo "run id may contain only letters, digits, underscores, and hyphens" >&2
  exit 2
fi
log="$probe/csv-makecsv-overlap-$run_id-api.log"
xvfb_log="$probe/xvfb-csv-makecsv-overlap-$run_id.log"
input_backup="/tmp/csv-makecsv-overlap-$run_id-input1.txt"
output_backup="/tmp/csv-makecsv-overlap-$run_id-output.wav"
xpid=
restore() {
  cp "$input_backup" "$work/input1.txt"
  cp "$output_backup" "$work/output.wav"
  cmp -s "$input_backup" "$work/input1.txt"
  cmp -s "$output_backup" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
for path in "$log" "$xvfb_log" "$input_backup" "$output_backup"; do
  if [ -e "$path" ]; then
    echo "refusing to overwrite $path" >&2
    exit 3
  fi
done
cp "$work/input1.txt" "$input_backup"
cp "$work/output.wav" "$output_backup"
trap restore EXIT
cd "$work"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 600s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-csv-makecsv-overlap-v1.gdb" > "$log" 2>&1
grep -q 'CSV_ALIAS_DONE calls=95' "$log"
grep -q '\[Inferior .*exited normally\]' "$log"
grep '^CSV_ALIAS_DONE' "$log"
