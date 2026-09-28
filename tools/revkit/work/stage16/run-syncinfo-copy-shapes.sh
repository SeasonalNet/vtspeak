#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage16
log="$probe/syncinfo-copy-shapes-api.log"
xvfb_log="$probe/xvfb-syncinfo-copy-shapes.log"
xpid=
restore() {
  cp "$probe/syncinfo-copy-shapes-original-input1.txt" "$work/input1.txt"
  cp "$probe/syncinfo-copy-shapes-original-output.wav" "$work/output.wav"
  cmp -s "$probe/syncinfo-copy-shapes-original-input1.txt" "$work/input1.txt"
  cmp -s "$probe/syncinfo-copy-shapes-original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
for output in "$log" "$xvfb_log"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
cp "$work/input1.txt" "$probe/syncinfo-copy-shapes-original-input1.txt"
cp "$work/output.wav" "$probe/syncinfo-copy-shapes-original-output.wav"
trap restore EXIT
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-syncinfo-copy-shapes.gdb" > "$log" 2>&1
