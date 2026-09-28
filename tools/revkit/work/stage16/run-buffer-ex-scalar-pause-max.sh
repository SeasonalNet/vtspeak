#!/bin/bash
set -euo pipefail

if [ "$#" -ne 1 ] || ! [[ "$1" =~ ^[0-2]$ ]]; then
  echo "usage: $0 SELECTOR" >&2
  exit 2
fi
selector=$1
work=/work/stage5
probe=/work/stage16
trace="$probe/trace-buffer-ex-scalar-pause-max-$selector.gdb"
log="$probe/buffer-ex-scalar-pause-max-$selector.log"
xvfb_log="$probe/xvfb-buffer-ex-scalar-pause-max-$selector.log"
backup_input="$probe/buffer-ex-scalar-pause-max-$selector-original-input1.txt"
backup_output="$probe/buffer-ex-scalar-pause-max-$selector-original-output.wav"
xpid=
restore() {
  cp "$backup_input" "$work/input1.txt"
  cp "$backup_output" "$work/output.wav"
  cmp -s "$backup_input" "$work/input1.txt"
  cmp -s "$backup_output" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

for output in "$trace" "$log" "$xvfb_log" "$backup_input" "$backup_output"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
cp "$work/input1.txt" "$backup_input"
cp "$work/output.wav" "$backup_output"
trap restore EXIT
cp "$probe/input-ex-scalar-pause-sentence.txt" "$work/input1.txt"
sed "s/@SELECTOR@/$selector/g" \
  "$probe/trace-buffer-ex-scalar-pause.gdb.in" > "$trace"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" > "$log" 2>&1
