#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage8
map_input=${MAP_INPUT:-map-input.txt}
trace_script=${TRACE_SCRIPT:-trace-map.gdb}
map_output=${MAP_OUTPUT:-map.wav}
xpid=
restore() {
  cp "$probe/original-map-input1.txt" "$work/input1.txt"
  cp "$probe/original-map-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

cp "$work/input1.txt" "$probe/original-map-input1.txt"
cp "$work/output.wav" "$probe/original-map-output.wav"
trap restore EXIT
cp "$probe/$map_input" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/map-xvfb.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/$trace_script" > "$probe/$map_output-run.log" 2>&1
cp "$work/output.wav" "$probe/$map_output"
