#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
history_mode=${1:-single}
case "$history_mode" in
  single) trace="$probe/trace-unit-history-present.gdb" ;;
  repeat) trace="$probe/trace-unit-history-repeat.gdb" ;;
  *) printf 'usage: %s [single|repeat]\n' "$0" >&2; exit 2 ;;
esac
overlay=/tmp/vtspeak-history-pattern-matrix
xpid=
backup_input=/tmp/vtspeak-history-pattern-input.txt
backup_output=/tmp/vtspeak-history-pattern-output.wav
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
mkdir -p "$overlay/stage5" "$model/mc_idx_tbl"
ln -s /work/data-common "$overlay/data-common"
for source in /work/data-paul/M16/*; do
  name=${source##*/}
  if [ "$name" != mc_idx_tbl ]; then ln -s "$source" "$model/$name"; fi
done
for bank in gen num etc alp; do
  ln -s "/work/data-paul/M16/mc_idx_tbl/unit-$bank.idx" "$model/mc_idx_tbl/unit-$bank.idx"
done
cp "$backup_output" "$overlay/stage5/output.wav"

write_pattern() {
  local pattern=$1
  local bank count header zero one file pair index
  for bank in gen num etc alp; do
    case "$bank" in
      gen) count=440124; header='\x3c\xb7\x06\x00' ;;
      num) count=24508; header='\xbc\x5f\x00\x00' ;;
      etc) count=115723; header='\x0b\xc4\x01\x00' ;;
      alp) count=119; header='\x77\x00\x00\x00' ;;
    esac
    zero='\x00\x00\x00\x00\x00\x00\x00\x00'
    case "$pattern" in
      a-alternating) one='\xff\xff\xff\x7f\x00\x00\x00\x00' ;;
      b-alternating) one='\x00\x00\x00\x00\xff\xff\xff\x7f' ;;
      both-alternating) one='\xff\xff\xff\x7f\xff\xff\xff\x7f' ;;
      *) printf 'unknown history pattern: %s\n' "$pattern" >&2; return 2 ;;
    esac
    file="$model/mc_idx_tbl/unit-$bank.his"
    {
      printf '%b' "$header"
      pair=$((count / 2))
      for ((index = 0; index < pair; index++)); do
        printf '%b%b' "$zero" "$one"
      done
      if ((count % 2)); then printf '%b' "$zero"; fi
    } > "$file"
    test "$(stat -c %s "$file")" -eq "$((4 + count * 8))"
  done
}

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-unit-history-pattern-matrix.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
cd "$overlay/stage5"

run_case() {
  local label=$1
  local text=$2
  printf '%s\n' "$text" > "$overlay/stage5/input1.txt"
  cp "$backup_output" "$overlay/stage5/output.wav"
  if [ "$history_mode" = repeat ]; then
    rm -f "$overlay/stage5"/repeat{1,2,3}.wav
  fi
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 240s \
    winedbg --gdb /samples/voicetext_paul.exe \
    < "$trace" > "$probe/unit-history-pattern-$label-api.log" 2>&1
  local log="$probe/unit-history-pattern-$label-api.log"
  if [ "$history_mode" = single ]; then
    test "$(grep -c '^UNIT_HISTORY_PRESENT_ENTRY ' "$log")" -eq 4
    test "$(grep -c '^UNIT_HISTORY_PRESENT_DATA ' "$log")" -eq 4
    test "$(grep -c '^UNIT_HISTORY_PRESENT_RETURN ax=1$' "$log")" -eq 4
    grep -Fq 'UNIT_HISTORY_PRESENT_LOAD ax=0' "$log"
    grep -Fq 'UNIT_HISTORY_PRESENT_TEXT ax=1' "$log"
  else
    test "$(grep -c '^UNIT_HISTORY_REPEAT_ENTRY ' "$log")" -eq 4
    test "$(grep -c '^UNIT_HISTORY_REPEAT_DATA ' "$log")" -eq 4
    test "$(grep -c '^UNIT_HISTORY_REPEAT_RETURN ax=1$' "$log")" -eq 4
    grep -Fq 'UNIT_HISTORY_REPEAT_LOAD ax=0' "$log"
    grep -Fq 'UNIT_HISTORY_REPEAT_CALL index=0 ax=1 ' "$log"
  fi
  cp "$overlay/stage5/output.wav" "$probe/unit-history-pattern-$label-output.wav"
  if [ "$history_mode" = repeat ]; then
    test "$(grep -c '^UNIT_HISTORY_REPEAT_CALL index=' "$log")" -eq 4
    test "$(grep -c '^UNIT_HISTORY_REPEAT_CALL index=[1-3] ax=1$' "$log")" -eq 3
    for repeat_index in 1 2 3; do
      test -f "$overlay/stage5/repeat$repeat_index.wav"
      cp "$overlay/stage5/repeat$repeat_index.wav" \
        "$probe/unit-history-pattern-$label-repeat-$repeat_index-output.wav"
    done
  fi
}

compare_case() {
  local text_name=$1
  local text=$2
  local pattern=$3
  rm -f "$model/mc_idx_tbl"/*.his
  run_case "$text_name-fallback" "$text"
  cp "$probe/unit-history-pattern-$text_name-fallback-output.wav" "$overlay/fallback.wav"
  write_pattern "$pattern"
  run_case "$text_name-$pattern" "$text"
  if cmp -s "$overlay/fallback.wav" "$probe/unit-history-pattern-$text_name-$pattern-output.wav"; then
    printf 'UNIT_HISTORY_PATTERN_COMPARE text=%s pattern=%s identical=1\n' "$text_name" "$pattern"
  else
    printf 'UNIT_HISTORY_PATTERN_COMPARE text=%s pattern=%s identical=0\n' "$text_name" "$pattern"
  fi
  if [ "$history_mode" = repeat ]; then
    for repeat_index in 1 2 3; do
      if cmp -s \
        "$probe/unit-history-pattern-$text_name-fallback-repeat-$repeat_index-output.wav" \
        "$probe/unit-history-pattern-$text_name-$pattern-repeat-$repeat_index-output.wav"; then
        printf 'UNIT_HISTORY_PATTERN_REPEAT_COMPARE text=%s pattern=%s repeat=%s identical=1\n' \
          "$text_name" "$pattern" "$repeat_index"
      else
        printf 'UNIT_HISTORY_PATTERN_REPEAT_COMPARE text=%s pattern=%s repeat=%s identical=0\n' \
          "$text_name" "$pattern" "$repeat_index"
      fi
    done
  fi
}

results=/tmp/vtspeak-history-pattern-matrix-results.txt
: > "$results"
for text_case in hello pangram numbers; do
  case "$text_case" in
    hello) text='Hello world.' ;;
    pangram) text='The quick brown fox jumps over the lazy dog.' ;;
    numbers) text='I saw 123 birds at 10:30.' ;;
  esac
  for pattern in a-alternating b-alternating both-alternating; do
    compare_case "$text_case" "$text" "$pattern" | tee -a "$results"
  done
done
cp "$results" "$probe/unit-history-pattern-matrix-results.txt"
