#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
input_backup="$probe/highlight-two-byte-classifier-original-input.txt"
output_backup="$probe/highlight-two-byte-classifier-original-output.wav"
log="$probe/highlight-two-byte-classifier-api.log"
xvfb_log="$probe/xvfb-highlight-two-byte-classifier.log"
xpid=
restore() {
  cp "$input_backup" "$work/input1.txt"
  cp "$output_backup" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

for path in "$input_backup" "$output_backup" "$log" "$xvfb_log"; do
  if [ -e "$path" ]; then
    echo "refusing to overwrite $path" >&2
    exit 3
  fi
done
cp "$work/input1.txt" "$input_backup"
cp "$work/output.wav" "$output_backup"
trap restore EXIT
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
cd "$work"
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 600s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-highlight-two-byte-classifier.gdb" > "$log" 2>&1
