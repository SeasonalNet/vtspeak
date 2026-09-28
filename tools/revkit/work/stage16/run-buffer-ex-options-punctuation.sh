#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage16
out="$probe/options-punctuation"
xpid=
restore() {
  cp "$probe/buffer-ex-options-punctuation-original-input1.txt" "$work/input1.txt"
  cp "$probe/buffer-ex-options-punctuation-original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

mkdir -p "$out"
cp "$work/input1.txt" "$probe/buffer-ex-options-punctuation-original-input1.txt"
cp "$work/output.wav" "$probe/buffer-ex-options-punctuation-original-output.wav"
trap restore EXIT
cp "$probe/input-ex-options-punctuation.txt" "$work/input1.txt"
sed 's|/work/stage16/buffer-ex-options-|/work/stage16/options-punctuation/buffer-ex-options-|g' \
  "$probe/trace-buffer-ex-options.gdb.in" > /tmp/trace-buffer-ex-options-punctuation.gdb
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$out/xvfb.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < /tmp/trace-buffer-ex-options-punctuation.gdb > "$out/runtime.log" 2>&1
