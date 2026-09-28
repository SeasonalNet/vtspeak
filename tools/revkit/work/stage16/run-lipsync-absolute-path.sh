#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage16
output="$probe/lipsync-absolute-path-report.txt"
xpid=
backup_input=/tmp/vtspeak-lipsync-absolute-path-input1.txt
backup_output=/tmp/vtspeak-lipsync-absolute-path-output.wav
if [ -e "$output" ]; then
  echo "refusing to overwrite $output" >&2
  exit 1
fi
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
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-lipsync-absolute-path.log" 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-lipsync-absolute-path.gdb" > "$probe/lipsync-absolute-path.log" 2>&1
if [ -e "$output" ]; then
  printf 'LIPSYNC_ABSOLUTE_PATH_OUTPUT exists\n' >> "$probe/lipsync-absolute-path.log"
else
  printf 'LIPSYNC_ABSOLUTE_PATH_OUTPUT none\n' >> "$probe/lipsync-absolute-path.log"
fi
