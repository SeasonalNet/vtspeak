#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
run_id=${1:-v1}
if [[ ! "$run_id" =~ ^[A-Za-z0-9_-]+$ ]]; then
  echo "run id may contain only letters, digits, underscores, and hyphens" >&2
  exit 2
fi
log="$probe/csv-makecsv-byte-pairs-full-$run_id-api.log"
result="$probe/csv-makecsv-byte-pairs-full-$run_id.bin"
fault_result="$probe/csv-makecsv-byte-pairs-fault-$run_id.bin"
xvfb_log="$probe/xvfb-csv-makecsv-byte-pairs-full-$run_id.log"
trace="/tmp/trace-csv-makecsv-byte-pairs-full-$run_id.gdb"
loader="$probe/load-csv-makecsv-byte-pairs-native-$run_id.gdb"
input_backup="/tmp/csv-makecsv-byte-pairs-full-$run_id-input1.txt"
output_backup="/tmp/csv-makecsv-byte-pairs-full-$run_id-output.wav"
xpid=
restore() {
  cp "$input_backup" "$work/input1.txt"
  cp "$output_backup" "$work/output.wav"
  cmp -s "$input_backup" "$work/input1.txt"
  cmp -s "$output_backup" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
for path in "$log" "$result" "$fault_result" "$xvfb_log" "$trace" "$loader" "$input_backup" "$output_backup"; do
  if [ -e "$path" ]; then
    echo "refusing to overwrite $path" >&2
    exit 3
  fi
done
sed "s/__RUN_ID__/$run_id/g" \
  "$probe/trace-csv-makecsv-byte-pairs-full-v1.gdb.in" > "$trace"
od -An -v -t x1 "$probe/csv-makecsv-byte-pairs-native.bin" |
  awk '{ for (i = 1; i <= NF; i++) printf "set {unsigned char}($code + %d) = 0x%s\n", n++, $i }' > "$loader"
cp "$work/input1.txt" "$input_backup"
cp "$work/output.wav" "$output_backup"
trap restore EXIT
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" > "$log" 2>&1
grep -q 'CSV_MAKE_PAIR_FULL cases=65536 bytes=917510 guards_checked=65536 probe_return=0' "$log"
grep -q '\[Inferior .*exited normally\]' "$log"
test "$(stat -c '%s' "$result")" -eq 917510
grep '^CSV_MAKE_PAIR_FULL' "$log"
