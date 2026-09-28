#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage16
if [ "$#" -ne 1 ] || ! [[ "$1" =~ ^[0-2]$ ]]; then
  echo "usage: $0 SELECTOR" >&2
  exit 2
fi
selector=$1
xpid=
restore() {
  cp "$probe/buffer-ex-mark-options-original-input1.txt" "$work/input1.txt"
  cp "$probe/buffer-ex-mark-options-original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

cp "$work/input1.txt" "$probe/buffer-ex-mark-options-original-input1.txt"
cp "$work/output.wav" "$probe/buffer-ex-mark-options-original-output.wav"
trap restore EXIT
cp "$probe/input-ex-records-vtml-mark-boundary.txt" "$work/input1.txt"
sed "s/@SELECTOR@/$selector/g" "$probe/trace-buffer-ex-mark-options.gdb.in" \
  > /tmp/trace-buffer-ex-mark-options.gdb
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-buffer-ex-mark-options-$selector.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < /tmp/trace-buffer-ex-mark-options.gdb > "$probe/buffer-ex-mark-options-$selector.log" 2>&1
