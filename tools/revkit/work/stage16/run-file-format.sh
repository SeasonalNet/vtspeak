#!/bin/bash
set -euo pipefail

if [ "$#" -ne 1 ] || ! [[ "$1" =~ ^[0-9]+$ ]]; then
  echo "usage: $0 FORMAT" >&2
  exit 2
fi

format=$1
case "$format" in
  0|1|2|3|4|5|6|7|8|9|10) ;;
  *) echo "format outside probe set: $format" >&2; exit 2 ;;
esac

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
sed "s/@FORMAT@/$format/g" "$probe/trace-file-format.gdb.in" > "/tmp/trace-file-format-$format.gdb"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-$format.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "/tmp/trace-file-format-$format.gdb" > "$probe/file-format-$format.log" 2>&1
cp "$work/output.wav" "$probe/file-format-$format.out"
