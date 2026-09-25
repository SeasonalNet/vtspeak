#!/bin/bash
set -euo pipefail

if [ "$#" -ne 3 ] || ! [[ "$1" =~ ^[0-3]$ ]] || ! [[ "$2" =~ ^[0-9]+$ ]] || ! [[ "$3" =~ ^[0-9]+$ ]]; then
  echo "usage: $0 FORMAT FLAG THREAD_ID" >&2
  exit 2
fi

format=$1
flag=$2
thread_id=$3
work=/work/stage5
probe=/work/stage16
xpid=
restore() {
  cp "$probe/original-input1.txt" "$work/input1.txt"
  cp "$probe/original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

cp "$work/input1.txt" "$probe/original-input1.txt"
cp "$work/output.wav" "$probe/original-output.wav"
trap restore EXIT
cp "$probe/input.txt" "$work/input1.txt"
sed -e "s/@FORMAT@/$format/g" -e "s/@FLAG@/$flag/g" -e "s/@THREAD@/$thread_id/g" \
  "$probe/trace-buffer.gdb.in" > "/tmp/trace-buffer-$format-$flag-$thread_id.gdb"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-buffer-$format-$flag-$thread_id.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "/tmp/trace-buffer-$format-$flag-$thread_id.gdb" > "$probe/buffer-$format-$flag-$thread_id.log" 2>&1
