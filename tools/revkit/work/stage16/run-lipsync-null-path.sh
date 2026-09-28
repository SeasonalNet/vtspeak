#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage16
xpid=
backup_input=/tmp/vtspeak-lipsync-null-path-input1.txt
backup_output=/tmp/vtspeak-lipsync-null-path-output.wav
shopt -s nullglob
existing=("$work"/length-sync-*)
if [ "${#existing[@]}" -ne 0 ]; then
  echo "refusing to run with an existing default lip-sync report" >&2
  exit 1
fi
if [ -e "$probe/lipsync-null-path-report.txt" ] || [ -e "$probe/lipsync-null-path.log" ]; then
  echo "refusing to overwrite null-path captures" >&2
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
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-lipsync-null-path.log" 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-lipsync-null-path.gdb" > "$probe/lipsync-null-path.log" 2>&1
outputs=("$work"/length-sync-*)
if [ "${#outputs[@]}" -eq 0 ]; then
  printf 'LIPSYNC_NULL_PATH_OUTPUT none\n' >> "$probe/lipsync-null-path.log"
elif [ "${#outputs[@]}" -eq 1 ]; then
  printf 'LIPSYNC_NULL_PATH_OUTPUT unexpected-default-name=%s\n' "${outputs[0]##*/}" >> "$probe/lipsync-null-path.log"
  mv "${outputs[0]}" "$probe/lipsync-null-path-report.txt"
else
  echo "unexpected default lip-sync report count: ${#outputs[@]}" >&2
  exit 1
fi
