#!/bin/bash
set -euo pipefail

if [ "$#" -ne 1 ] || ! [[ "$1" =~ ^(100|200)$ ]]; then
  echo "usage: $0 100|200" >&2
  exit 2
fi
speed=$1

work=/work/stage5
probe=/work/stage16
xpid=
backup_input=/tmp/vtspeak-lipsync-speed-input1.txt
backup_output=/tmp/vtspeak-lipsync-speed-output.wav
shopt -s nullglob
existing=("$work"/vtspeak-lipunits-*)
if [ "${#existing[@]}" -ne 0 ]; then
  echo "refusing to run with an existing lipsync output" >&2
  exit 1
fi
if [ -e "$probe/lipsync-speed-$speed.log" ] || [ -e "$probe/lipsync-speed-$speed-report.txt" ]; then
  echo "refusing to overwrite speed-$speed captures" >&2
  exit 1
fi

restore() {
  cp "$backup_input" "$work/input1.txt"
  cp "$backup_output" "$work/output.wav"
  cmp -s "$backup_input" "$work/input1.txt"
  cmp -s "$backup_output" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

cp "$work/input1.txt" "$backup_input"
cp "$work/output.wav" "$backup_output"
trap restore EXIT
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-lipsync-speed.log" 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99

trace=/tmp/trace-lipsync-speed-$speed.gdb
sed "s/@SPEED@/$speed/g" "$probe/trace-lipsync-speed.gdb.in" > "$trace"
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" > "$probe/lipsync-speed-$speed.log" 2>&1
output="$work/vtspeak-lipunits-$speed"
if [ ! -f "$output" ]; then
  echo "expected explicit output file $output" >&2
  exit 1
fi
printf 'LIPSYNC_SPEED_FILE speed=%s name=%s\n' "$speed" "${output##*/}" >> "$probe/lipsync-speed-$speed.log"
mv "$output" "$probe/lipsync-speed-$speed-report.txt"
