#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
xpid=
backup_input=/tmp/vtspeak-userdict-limit-input1.txt
backup_output=/tmp/vtspeak-userdict-limit-output.wav
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
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-userdict-limit-domain.log" 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-userdict-limit-domain.gdb" > "$probe/userdict-limit-domain-api.log" 2>&1

if [ "$(grep -c '^USERDICT_LIMIT selector=' "$probe/userdict-limit-domain-api.log")" -ne 10 ]; then
  echo "did not capture all 10 user-dictionary limit selectors" >&2
  exit 1
fi
grep -Fq 'USERDICT_LIMIT selector=INT_MIN result=-1' "$probe/userdict-limit-domain-api.log"
grep -Fq 'USERDICT_LIMIT selector=-1 result=-1' "$probe/userdict-limit-domain-api.log"
grep -Fq 'USERDICT_LIMIT selector=0 result=30' "$probe/userdict-limit-domain-api.log"
grep -Fq 'USERDICT_LIMIT selector=1 result=10' "$probe/userdict-limit-domain-api.log"
grep -Fq 'USERDICT_LIMIT selector=2 result=50' "$probe/userdict-limit-domain-api.log"
grep -Fq 'USERDICT_LIMIT selector=3 result=65' "$probe/userdict-limit-domain-api.log"
grep -Fq 'USERDICT_LIMIT selector=4 result=65' "$probe/userdict-limit-domain-api.log"
grep -Fq 'USERDICT_LIMIT selector=5 result=-1' "$probe/userdict-limit-domain-api.log"
grep -Fq 'USERDICT_LIMIT selector=6 result=-1' "$probe/userdict-limit-domain-api.log"
grep -Fq 'USERDICT_LIMIT selector=INT_MAX result=-1' "$probe/userdict-limit-domain-api.log"
