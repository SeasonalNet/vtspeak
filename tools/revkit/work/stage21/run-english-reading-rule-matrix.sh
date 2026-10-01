#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
xpid=
backup_input=/tmp/english-reading-rule-matrix-input.txt
backup_output=/tmp/english-reading-rule-matrix-output.wav
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

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-english-reading-rule-matrix.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1

run_mode() {
  local text_case=$1 mode=$2
  local log="$probe/english-reading-rule-matrix-$text_case-$mode-api.log"
  local output="$probe/english-reading-rule-matrix-$text_case-$mode.wav"
  local input="$probe/input-reading-rule-$text_case.txt"
  if [ "$text_case" = baseline ]; then input="$probe/input-reading-rule.txt"; fi
  cp "$input" "$work/input1.txt"
  : > "$work/output.wav"
  sed "s/@MODE@/$mode/g" "$probe/trace-english-reading-rule-effect.gdb.in" > /tmp/english-reading-rule-matrix.gdb
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
    winedbg --gdb /samples/voicetext_paul.exe \
    < /tmp/english-reading-rule-matrix.gdb > "$log" 2>&1
  grep -Fq "READING_RULE_SET mode=$mode " "$log"
  grep -Fq "READING_RULE_SYNTHESIS_RETURN eax=0x1 value=$mode accesses=0" "$log"
  cp "$work/output.wav" "$output"
}

results=/tmp/english-reading-rule-matrix-results.txt
: > "$results"
for text_case in baseline heteronyms lead dates abbreviations acronyms; do
  run_mode "$text_case" 0
  run_mode "$text_case" 1
  if cmp -s \
    "$probe/english-reading-rule-matrix-$text_case-0.wav" \
    "$probe/english-reading-rule-matrix-$text_case-1.wav"; then
    printf 'READING_RULE_MATRIX_COMPARE case=%s modes=0,1 identical=1\n' "$text_case" | tee -a "$results"
  else
    printf 'READING_RULE_MATRIX_COMPARE case=%s modes=0,1 identical=0\n' "$text_case" | tee -a "$results"
  fi
done

run_mode dates 2147483647
if cmp -s \
  "$probe/english-reading-rule-matrix-dates-0.wav" \
  "$probe/english-reading-rule-matrix-dates-2147483647.wav"; then
  printf 'READING_RULE_MATRIX_COMPARE case=dates modes=0,2147483647 identical=1\n' | tee -a "$results"
else
  printf 'READING_RULE_MATRIX_COMPARE case=dates modes=0,2147483647 identical=0\n' | tee -a "$results"
fi
cp "$results" "$probe/english-reading-rule-matrix-results.txt"
