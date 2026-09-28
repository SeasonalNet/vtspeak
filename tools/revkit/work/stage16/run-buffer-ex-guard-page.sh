#!/bin/bash
set -euo pipefail

if [ "$#" -ne 1 ] || ! [[ "$1" =~ ^[0-2]$ ]]; then
  echo "usage: $0 SELECTOR" >&2
  exit 2
fi
selector=$1
work=/work/stage5
probe=/work/stage16
if [ "$selector" = 0 ]; then
  suffix=
else
  suffix="-$selector"
fi
trace="/tmp/trace-buffer-ex-guard-page-$selector.gdb"
log="$probe/buffer-ex-guard-page$suffix.log"
xvfb_log="$probe/xvfb-buffer-ex-guard-page$suffix.log"
xpid=
restore() {
  cp "$probe/buffer-ex-guard-page-original-input1.txt" "$work/input1.txt"
  cp "$probe/buffer-ex-guard-page-original-output.wav" "$work/output.wav"
  cmp -s "$probe/buffer-ex-guard-page-original-input1.txt" "$work/input1.txt"
  cmp -s "$probe/buffer-ex-guard-page-original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

for output in "$log" "$xvfb_log"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
trap restore EXIT
cp "$work/input1.txt" "$probe/buffer-ex-guard-page-original-input1.txt"
cp "$work/output.wav" "$probe/buffer-ex-guard-page-original-output.wav"
cp "$probe/input-ex-records-vtml-mark-position-2.txt" "$work/input1.txt"
sed "s/@SELECTOR@/$selector/g" "$probe/trace-buffer-ex-guard-page.gdb" > "$trace"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp \
  > "$xvfb_log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" > "$log" 2>&1
