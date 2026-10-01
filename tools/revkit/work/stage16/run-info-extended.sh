#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage16
xpid=
input_backup="$probe/info-extended-original-input1.txt"
output_backup="$probe/info-extended-original-output.wav"
restore() {
  cp "$input_backup" "$work/input1.txt"
  cp "$output_backup" "$work/output.wav"
  cmp -s "$input_backup" "$work/input1.txt"
  cmp -s "$output_backup" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

cp "$work/input1.txt" "$input_backup"
cp "$work/output.wav" "$output_backup"
trap restore EXIT
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-info-extended.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-info-extended.gdb" > "$probe/info-extended-api.log" 2>&1
