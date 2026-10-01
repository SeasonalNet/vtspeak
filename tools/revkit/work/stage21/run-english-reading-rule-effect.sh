#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
xpid=
restore() {
  cp /tmp/english-reading-rule-input1.txt "$work/input1.txt"
  cp /tmp/english-reading-rule-output.wav "$work/output.wav"
  cmp -s /tmp/english-reading-rule-input1.txt "$work/input1.txt"
  cmp -s /tmp/english-reading-rule-output.wav "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
cp "$work/input1.txt" /tmp/english-reading-rule-input1.txt
cp "$work/output.wav" /tmp/english-reading-rule-output.wav
trap restore EXIT
cp "$probe/input-reading-rule.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-english-reading-rule.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1

for mode in 0 1; do
  sed "s/@MODE@/$mode/g" "$probe/trace-english-reading-rule-effect.gdb.in" > "/tmp/english-reading-rule-$mode.gdb"
  : > "$work/output.wav"
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
    winedbg --gdb /samples/voicetext_paul.exe \
    < "/tmp/english-reading-rule-$mode.gdb" > "$probe/english-reading-rule-$mode-api.log" 2>&1
  grep -E '^(READING_RULE_|#)' "$probe/english-reading-rule-$mode-api.log"
  grep -q '\[Inferior .*exited normally\]' "$probe/english-reading-rule-$mode-api.log"
  cp "$work/output.wav" "$probe/english-reading-rule-$mode.wav"
done

sha256sum "$probe/english-reading-rule-0.wav" "$probe/english-reading-rule-1.wav"
if cmp -s "$probe/english-reading-rule-0.wav" "$probe/english-reading-rule-1.wav"; then
  echo READING_RULE_OUTPUT_IDENTICAL=1
else
  echo READING_RULE_OUTPUT_IDENTICAL=0
fi
