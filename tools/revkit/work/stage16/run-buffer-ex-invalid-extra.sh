#!/bin/bash
set -euo pipefail

if [ "$#" -ne 1 ] || ! [[ "$1" =~ ^[ab]$ ]]; then
  echo "usage: $0 a|b" >&2
  exit 2
fi
slot=$1
work=/work/stage5
probe=/work/stage16
trace="$probe/trace-buffer-ex-invalid-extra-$slot.gdb"
log="$probe/buffer-ex-invalid-extra-$slot.log"
xvfb_log="$probe/xvfb-buffer-ex-invalid-extra-$slot.log"
xpid=
restore() {
  cp "$probe/buffer-ex-invalid-extra-original-input1.txt" "$work/input1.txt"
  cp "$probe/buffer-ex-invalid-extra-original-output.wav" "$work/output.wav"
  cmp -s "$probe/buffer-ex-invalid-extra-original-input1.txt" "$work/input1.txt"
  cmp -s "$probe/buffer-ex-invalid-extra-original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

for output in "$trace" "$log" "$xvfb_log"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
trap restore EXIT
cp "$work/input1.txt" "$probe/buffer-ex-invalid-extra-original-input1.txt"
cp "$work/output.wav" "$probe/buffer-ex-invalid-extra-original-output.wav"
cmp -s "$work/input1.txt" "$probe/makeinfo-original-input1.txt"
cmp -s "$work/output.wav" "$probe/makeinfo-original-output.wav"
if [ "$slot" = a ]; then
  bad_a=1
  bad_b=0
else
  bad_a=0
  bad_b=1
fi
sed -e "s/@BAD_A@/$bad_a/g" -e "s/@BAD_B@/$bad_b/g" \
  -e "s/@SLOT@/$slot/g" \
  "$probe/trace-buffer-ex-invalid-extra.gdb.in" > "$trace"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" > "$log" 2>&1
