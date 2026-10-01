#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
overlay=/tmp/vtspeak-unit-history-overlay
xpid=
backup_input=/tmp/vtspeak-unit-history-present-input.txt
backup_output=/tmp/vtspeak-unit-history-present-output.wav
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

model="$overlay/data-paul/M16"
mkdir -p "$overlay/stage5"
ln -s /work/data-common "$overlay/data-common"
mkdir -p "$model/mc_idx_tbl"
for source in /work/data-paul/M16/*; do
  name=${source##*/}
  if [ "$name" != mc_idx_tbl ]; then ln -s "$source" "$model/$name"; fi
done
for bank in gen num etc alp; do
  ln -s "/work/data-paul/M16/mc_idx_tbl/unit-$bank.idx" "$model/mc_idx_tbl/unit-$bank.idx"
done

write_history() {
  local bank=$1 count=$2 header=$3 first_a=$4 first_b=$5 last_a=$6 last_b=$7
  local file="$model/mc_idx_tbl/unit-$bank.his"
  truncate -s "$((4 + count * 8))" "$file"
  printf '%b' "$header" > "$file"
  printf '%b' "$first_a" "$first_b" | dd of="$file" bs=1 seek=4 conv=notrunc status=none
  printf '%b' "$last_a" "$last_b" | dd of="$file" bs=1 seek="$((4 + (count - 1) * 8))" conv=notrunc status=none
  printf 'UNIT_HISTORY_FIXTURE bank=%s count=%d size=%d\n' "$bank" "$count" "$(stat -c %s "$file")"
}
write_history gen 440124 '\x3c\xb7\x06\x00' '\x11\x11\x00\x00' '\x12\x12\x00\x00' '\x31\x31\x00\x00' '\x32\x32\x00\x00'
write_history num 24508 '\xbc\x5f\x00\x00' '\x21\x21\x00\x00' '\x22\x22\x00\x00' '\x41\x41\x00\x00' '\x42\x42\x00\x00'
write_history etc 115723 '\x0b\xc4\x01\x00' '\x51\x51\x00\x00' '\x52\x52\x00\x00' '\x71\x71\x00\x00' '\x72\x72\x00\x00'
write_history alp 119 '\x77\x00\x00\x00' '\x81\x81\x00\x00' '\x82\x82\x00\x00' '\xa1\xa1\x00\x00' '\xa2\xa2\x00\x00'

cp "$backup_input" "$overlay/stage5/input1.txt"
cp "$backup_output" "$overlay/stage5/output.wav"

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-unit-history-present.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
cd "$overlay/stage5"
run_trace() {
  local label=$1
  local wine_debug=-all
  if [ "$label" = header-mismatch ] || [ "$label" = short-read ]; then wine_debug=+seh; fi
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG="$wine_debug" timeout 240s \
    winedbg --gdb /samples/voicetext_paul.exe \
    < "$probe/trace-unit-history-present.gdb" > "$probe/unit-history-$label-api.log" 2>&1
}
run_trace present

present_log="$probe/unit-history-present-api.log"
cp "$overlay/stage5/output.wav" "$probe/unit-history-present-output.wav"
test "$(grep -c '^UNIT_HISTORY_PRESENT_ENTRY ' "$present_log")" -eq 4
test "$(grep -c '^UNIT_HISTORY_PRESENT_OPEN source=0x' "$present_log")" -eq 4
test "$(grep -c '^UNIT_HISTORY_PRESENT_DATA ' "$present_log")" -eq 4
test "$(grep -c '^UNIT_HISTORY_PRESENT_RETURN ax=1$' "$present_log")" -eq 4
grep -Fq 'UNIT_HISTORY_PRESENT_LOAD ax=0' "$present_log"
grep -Fq 'UNIT_HISTORY_PRESENT_TEXT ax=1' "$present_log"

gen_file="$model/mc_idx_tbl/unit-gen.his"
printf '%b' '\x3b\xb7\x06\x00' > "$gen_file"
run_trace header-mismatch
mismatch_log="$probe/unit-history-header-mismatch-api.log"
test "$(grep -c '^UNIT_HISTORY_PRESENT_ENTRY ' "$mismatch_log")" -eq 1
test "$(grep -c '^UNIT_HISTORY_PRESENT_OPEN source=0x' "$mismatch_log")" -eq 1
grep -Fq 'UNIT_HISTORY_PRESENT_PARSE_FAILURE' "$mismatch_log"
grep -Fq 'UNIT_HISTORY_PRESENT_RETURN ax=0' "$mismatch_log"
grep -Fq 'dispatch_exception code=c0000005 flags=0 addr=10027D9C ip=10027d9c' "$mismatch_log"
grep -Fq 'info[1]=00004d08' "$mismatch_log"

write_history gen 440124 '\x3c\xb7\x06\x00' '\x11\x11\x00\x00' '\x12\x12\x00\x00' '\x31\x31\x00\x00' '\x32\x32\x00\x00'
truncate -s "$((4 + 440124 * 8 - 4))" "$gen_file"
run_trace short-read
short_log="$probe/unit-history-short-read-api.log"
test "$(grep -c '^UNIT_HISTORY_PRESENT_ENTRY ' "$short_log")" -eq 1
test "$(grep -c '^UNIT_HISTORY_PRESENT_OPEN source=0x' "$short_log")" -eq 1
grep -Fq 'UNIT_HISTORY_PRESENT_PARSE_FAILURE' "$short_log"
grep -Fq 'UNIT_HISTORY_PRESENT_RETURN ax=0' "$short_log"
grep -Fq 'dispatch_exception code=c0000005 flags=0 addr=10027D9C ip=10027d9c' "$short_log"
grep -Fq 'info[1]=00004d08' "$short_log"

write_history gen 440124 '\x3c\xb7\x06\x00' '\x11\x11\x00\x00' '\x12\x12\x00\x00' '\x31\x31\x00\x00' '\x32\x32\x00\x00'
truncate -s "$((4 + 440124 * 8 + 4))" "$gen_file"
run_trace trailing-data
trailing_log="$probe/unit-history-trailing-data-api.log"
test "$(grep -c '^UNIT_HISTORY_PRESENT_RETURN ax=1$' "$trailing_log")" -eq 4
grep -Fq 'UNIT_HISTORY_PRESENT_LOAD ax=0' "$trailing_log"
grep -Fq 'UNIT_HISTORY_PRESENT_TEXT ax=1' "$trailing_log"

cd "$work"
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-unit-history-load.gdb" > "$probe/unit-history-fallback-comparison-api.log" 2>&1
cp "$work/output.wav" "$probe/unit-history-fallback-output.wav"
grep -Fq 'UNIT_HISTORY_TEXT_RETURN ax=1' "$probe/unit-history-fallback-comparison-api.log"
if cmp -s "$probe/unit-history-present-output.wav" "$probe/unit-history-fallback-output.wav"; then
  printf 'UNIT_HISTORY_OUTPUT_COMPARE identical=1\n'
else
  printf 'UNIT_HISTORY_OUTPUT_COMPARE identical=0\n'
  exit 1
fi
