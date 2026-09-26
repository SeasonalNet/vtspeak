#!/bin/bash
set -euo pipefail

if [ "$#" -ne 1 ] || ! [[ "$1" =~ ^[01]$ ]]; then
  echo "usage: $0 0|1" >&2
  exit 2
fi

mode=$1
work=/work/stage5
probe=/work/stage16
xpid=
restore() {
  cp "/tmp/unit-history-input1-$mode.txt" "$work/input1.txt"
  cp "/tmp/unit-history-output-$mode.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

cp "$work/input1.txt" "/tmp/unit-history-input1-$mode.txt"
cp "$work/output.wav" "/tmp/unit-history-output-$mode.wav"
trap restore EXIT
cp "$probe/input.txt" "$work/input1.txt"
sed "s/@MODE@/$mode/g" "$probe/trace-unit-history.gdb.in" > "/tmp/trace-unit-history-$mode.gdb"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-unit-history-$mode.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "/tmp/trace-unit-history-$mode.gdb" > "$probe/unit-history-$mode-api.log" 2>&1
cp "$work/output.wav" "$probe/unit-history-$mode.wav"
