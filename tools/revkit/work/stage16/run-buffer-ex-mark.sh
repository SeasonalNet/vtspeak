#!/bin/bash
set -euo pipefail

if [ "$#" -lt 1 ] || [ "$#" -gt 2 ] || ! [[ "$1" =~ ^[0-2]$ ]]; then
  echo "usage: $0 SELECTOR [mark|boundary|boundary-sweep|position-N]" >&2
  exit 2
fi

selector=$1
fixture=${2:-mark}
case "$fixture" in
  mark) input=input-ex-records-vtml-mark.txt ;;
  boundary) input=input-ex-records-vtml-mark-boundary.txt ;;
  boundary-sweep) input=input-ex-records-vtml-mark-boundary-sweep.txt ;;
  position-[0-9]*) input="input-ex-records-vtml-mark-$fixture.txt" ;;
  *) echo "unknown fixture: $fixture" >&2; exit 2 ;;
esac
work=/work/stage5
probe=/work/stage16
xpid=
restore() {
  cp "$probe/buffer-ex-mark-original-input1.txt" "$work/input1.txt"
  cp "$probe/buffer-ex-mark-original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

cp "$work/input1.txt" "$probe/buffer-ex-mark-original-input1.txt"
cp "$work/output.wav" "$probe/buffer-ex-mark-original-output.wav"
trap restore EXIT
cp "$probe/$input" "$work/input1.txt"
sed "s/@SELECTOR@/$selector/g" "$probe/trace-buffer-ex-records.gdb.in" \
  > "/tmp/trace-buffer-ex-mark-$selector.gdb"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-buffer-ex-mark-$fixture-$selector.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "/tmp/trace-buffer-ex-mark-$selector.gdb" > "$probe/buffer-ex-mark-$fixture-$selector.log" 2>&1
