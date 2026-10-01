#!/bin/bash
set -euo pipefail

if [ "$#" -ne 2 ] || ! [[ "$1" =~ ^[01]$ ]] || ! [[ "$2" =~ ^[01]$ ]]; then
  echo "usage: $0 SELECTOR FLAG" >&2
  exit 2
fi
selector=$1
flag=$2
work=/work/stage5
probe=/work/stage21
input_backup="$probe/highlight-ex-raw-original-$selector-$flag-input.txt"
output_backup="$probe/highlight-ex-raw-original-$selector-$flag-output.wav"
trace="$probe/trace-highlight-ex-raw-$selector-$flag.gdb"
log="$probe/highlight-ex-raw-$selector-$flag-api.log"
xvfb_log="$probe/xvfb-highlight-ex-raw-$selector-$flag.log"
dump="$probe/highlight-ex-raw-$selector-$flag-descriptor.bin"
xpid=
restore() {
  cp "$input_backup" "$work/input1.txt"
  cp "$output_backup" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

for path in "$input_backup" "$output_backup" "$trace" "$log" "$xvfb_log" "$dump"; do
  if [ -e "$path" ]; then
    echo "refusing to overwrite $path" >&2
    exit 3
  fi
done
cp "$work/input1.txt" "$input_backup"
cp "$work/output.wav" "$output_backup"
trap restore EXIT
cp /work/stage16/input-ex-records-vtml-mark-boundary.txt "$work/input1.txt"
sed -e "s/@SELECTOR@/$selector/g" -e "s/@FLAG@/$flag/g" \
  "$probe/trace-highlight-ex-raw.gdb.in" > "$trace"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
cd "$work"
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe < "$trace" > "$log" 2>&1
