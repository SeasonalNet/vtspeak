#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage16
xpid=
backup_input=/tmp/vtspeak-lipsync-prefix-input1.txt
backup_output=/tmp/vtspeak-lipsync-prefix-output.wav
shopt -s nullglob
existing=("$work"/length-sync-* "$work"/*vtspeak-lead6*)
if [ "${#existing[@]}" -ne 0 ]; then
  echo "refusing to run with an existing matching lip-sync output" >&2
  exit 1
fi
if [ -e "$probe/lipsync-prefix-output.txt" ]; then
  echo "refusing to overwrite captured lip-sync output" >&2
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
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-lipsync-prefix.log" 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-lipsync-prefix.gdb" > "$probe/lipsync-prefix.log" 2>&1

outputs=("$work"/length-sync-* "$work"/*vtspeak-lead6*)
if [ "${#outputs[@]}" -eq 0 ]; then
  echo "LIPSYNC_PREFIX_OUTPUT none"
elif [ "${#outputs[@]}" -eq 1 ]; then
  mv "${outputs[0]}" "$probe/lipsync-prefix-output.txt"
  echo "LIPSYNC_PREFIX_OUTPUT captured"
else
  echo "LIPSYNC_PREFIX_OUTPUT unexpected-count=${#outputs[@]}" >&2
  exit 1
fi
