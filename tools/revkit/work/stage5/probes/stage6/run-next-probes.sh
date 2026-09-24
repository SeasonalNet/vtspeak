#!/bin/bash
set -euo pipefail

work=/work/stage5
probe_root="$work/probes/stage6"
xpid=

restore_inputs() {
  cp "$probe_root/original-input1.txt" "$work/input1.txt"
  cp "$probe_root/original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then
    kill "$xpid" 2>/dev/null || true
  fi
}

trap restore_inputs EXIT
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe_root/xvfb.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1

for name in numbers abbreviations; do
  cp "$probe_root/$name/input1.txt" "$work/input1.txt"
  timeout 120s winedbg --gdb /samples/voicetext_paul.exe \
    < "$probe_root/$name/trace.gdb" > "$probe_root/$name/gdb-run.log" 2>&1
  cp "$work/output.wav" "$probe_root/$name/output.wav"
done
