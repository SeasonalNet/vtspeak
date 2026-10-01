#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
xpid=
backup_input=/tmp/vtspeak-speaker-query-input1.txt
backup_output=/tmp/vtspeak-speaker-query-output.wav
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
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-speaker-query-edges.log" 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-speaker-query-edges.gdb" > "$probe/speaker-query-edges-api.log" 2>&1

test "$(grep -c '^SPEAKER_NAME_EDGE ' "$probe/speaker-query-edges-api.log")" -eq 10
test "$(grep -c '^DB_SIZE_EDGE ' "$probe/speaker-query-edges-api.log")" -eq 10
grep -Fq 'SPEAKER_NAME_EDGE input=-2147483648 value=Paul' "$probe/speaker-query-edges-api.log"
grep -Fq 'SPEAKER_NAME_EDGE input=6 value=Paul' "$probe/speaker-query-edges-api.log"
grep -Fq 'SPEAKER_NAME_EDGE input=2147483647 value=Paul' "$probe/speaker-query-edges-api.log"
grep -Fq 'SPEAKER_NAME_EDGE input=0 value=Kate' "$probe/speaker-query-edges-api.log"
grep -Fq 'SPEAKER_NAME_EDGE input=1 value=Paul' "$probe/speaker-query-edges-api.log"
grep -Fq 'SPEAKER_NAME_EDGE input=2 value=em001' "$probe/speaker-query-edges-api.log"
grep -Fq 'SPEAKER_NAME_EDGE input=3 value=Julie' "$probe/speaker-query-edges-api.log"
grep -Fq 'SPEAKER_NAME_EDGE input=4 value=James' "$probe/speaker-query-edges-api.log"
grep -Fq 'SPEAKER_NAME_EDGE input=5 value=Ashley' "$probe/speaker-query-edges-api.log"
grep -Fq 'DB_SIZE_EDGE input=-2147483648 result=1 value=508121688' "$probe/speaker-query-edges-api.log"
grep -Fq 'DB_SIZE_EDGE input=-1 result=1 value=508121688' "$probe/speaker-query-edges-api.log"
grep -Fq 'DB_SIZE_EDGE input=0 result=-1 value=324508639' "$probe/speaker-query-edges-api.log"
grep -Fq 'DB_SIZE_EDGE input=1 result=1 value=508121688' "$probe/speaker-query-edges-api.log"
grep -Fq 'DB_SIZE_EDGE input=2 result=-1 value=324508639' "$probe/speaker-query-edges-api.log"
grep -Fq 'DB_SIZE_EDGE input=3 result=-1 value=324508639' "$probe/speaker-query-edges-api.log"
grep -Fq 'DB_SIZE_EDGE input=4 result=-1 value=324508639' "$probe/speaker-query-edges-api.log"
grep -Fq 'DB_SIZE_EDGE input=5 result=-1 value=324508639' "$probe/speaker-query-edges-api.log"
grep -Fq 'DB_SIZE_EDGE input=6 result=1 value=508121688' "$probe/speaker-query-edges-api.log"
grep -Fq 'DB_SIZE_EDGE input=2147483647 result=1 value=508121688' "$probe/speaker-query-edges-api.log"
