#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage16
xpid=
backup_input=/tmp/vtspeak-lipsync-sentinel-open-input1.txt
backup_output=/tmp/vtspeak-lipsync-sentinel-open-output.wav
shopt -s nullglob
existing=("$work"/length-sync-*)
if [ "${#existing[@]}" -ne 0 ]; then
  echo 'refusing to run with an existing length-sync report' >&2
  exit 1
fi
if [ -e "$probe/lipsync-sentinel-open.log" ]; then
  echo 'refusing to overwrite sentinel-open capture' >&2
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
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-lipsync-sentinel-open.log" 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-lipsync-sentinel-open.gdb" > "$probe/lipsync-sentinel-open.log" 2>&1
if grep -Fq 'will be abandoned' "$probe/lipsync-sentinel-open.log"; then
  echo 'GDB aborted the inferior API call at the nested breakpoint; this capture is not API evidence' >&2
  exit 1
fi
outputs=("$work"/length-sync-*)
if [ "${#outputs[@]}" -eq 0 ]; then
  printf 'LIPSYNC_SENTINEL_OUTPUT none\n' >> "$probe/lipsync-sentinel-open.log"
elif [ "${#outputs[@]}" -eq 1 ]; then
  printf 'LIPSYNC_SENTINEL_OUTPUT %s\n' "${outputs[0]##*/}" >> "$probe/lipsync-sentinel-open.log"
  mv "${outputs[0]}" "$probe/lipsync-sentinel-open-output.txt"
else
  echo "unexpected default lip-sync report count: ${#outputs[@]}" >&2
  exit 1
fi
