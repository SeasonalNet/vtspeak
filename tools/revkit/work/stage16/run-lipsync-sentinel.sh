#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage16
xpid=
backup_input=/tmp/vtspeak-lipsync-sentinel-input1.txt
backup_output=/tmp/vtspeak-lipsync-sentinel-output.wav
shopt -s nullglob
existing=("$work"/length-sync-*)
if [ "${#existing[@]}" -ne 0 ]; then
  echo 'refusing to run with an existing length-sync report' >&2
  exit 1
fi
if [ -e "$probe/lipsync-sentinel.log" ] || [ -e "$probe/lipsync-sentinel-output.txt" ]; then
  echo 'refusing to overwrite sentinel captures' >&2
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
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-lipsync-sentinel.log" 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99
if WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-lipsync-sentinel.gdb" > "$probe/lipsync-sentinel.log" 2>&1; then
  process_status=0
else
  process_status=$?
fi
if ! grep -Fq 'LIPSYNC_SENTINEL_CALL' "$probe/lipsync-sentinel.log"; then
  echo "GDB did not reach the sentinel API call (process status $process_status)" >&2
  exit 1
fi
printf 'LIPSYNC_SENTINEL_PROCESS_STATUS %s\n' "$process_status" >> "$probe/lipsync-sentinel.log"
outputs=("$work"/length-sync-*)
if [ "${#outputs[@]}" -eq 0 ]; then
  printf 'LIPSYNC_SENTINEL_OUTPUT none\n' >> "$probe/lipsync-sentinel.log"
elif [ "${#outputs[@]}" -eq 1 ]; then
  printf 'LIPSYNC_SENTINEL_OUTPUT %s\n' "${outputs[0]##*/}" >> "$probe/lipsync-sentinel.log"
  mv "${outputs[0]}" "$probe/lipsync-sentinel-output.txt"
else
  echo "expected one default lip-sync report; observed ${#outputs[@]}" >&2
  exit 1
fi
