#!/bin/bash
set -euo pipefail

work=/work/stage5
records="$work/probes/stage6/records"
probe_input="${1:-$records/input1.txt}"
cp "$work/input1.txt" "$records/original-input1.txt"
cp "$work/output.wav" "$records/original-output.wav"
xpid=

restore_inputs() {
  cp "$records/original-input1.txt" "$work/input1.txt"
  cp "$records/original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then
    kill "$xpid" 2>/dev/null || true
  fi
}
trap restore_inputs EXIT

cp "$probe_input" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$records/xvfb.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
timeout 180s winedbg --gdb /samples/voicetext_paul.exe \
  < "$records/trace-phone-map.gdb" > "$records/gdb-run.log" 2>&1
cp "$work/output.wav" "$records/output.wav"
