#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
xpid=
backup_input=/tmp/vtspeak-unit-history-load-input.txt
backup_output=/tmp/vtspeak-unit-history-load-output.wav
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

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-unit-history-load.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-unit-history-load.gdb" > "$probe/unit-history-load-api.log" 2>&1
cp "$work/output.wav" "$probe/unit-history-fallback-output.wav"

grep -Fq 'UNIT_HISTORY_LOAD_SET value=1 loaded=0' "$probe/unit-history-load-api.log"
test "$(grep -c '^UNIT_HISTORY_ENTRY ' "$probe/unit-history-load-api.log")" -eq "$(grep -c '^UNIT_HISTORY_RETURN ' "$probe/unit-history-load-api.log")"
test "$(grep -c '^UNIT_HISTORY_SOURCE handle=(nil) ' "$probe/unit-history-load-api.log")" -eq 4
test "$(grep -c '^UNIT_HISTORY_FALLBACK count=.* a_first=0 a_last=0 b_first=0 b_last=0$' "$probe/unit-history-load-api.log")" -eq 4
grep -Fq 'UNIT_HISTORY_TEXT_RETURN ax=1' "$probe/unit-history-load-api.log"
