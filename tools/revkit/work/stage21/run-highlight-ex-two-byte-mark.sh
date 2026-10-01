#!/bin/bash
set -euo pipefail

if [ "$#" -lt 2 ] || [ "$#" -gt 3 ] || ! [[ "$1" =~ ^[012]$ ]] || ! [[ "$2" =~ ^[01]$ ]]; then
  echo "usage: $0 SELECTOR FLAG [single|classes]" >&2
  exit 2
fi
selector=$1
flag=$2
profile=${3:-single}
case "$profile" in
  single)
    fixture=input-ex-raw-highlight-two-byte.txt
    stem="highlight-ex-two-byte-$selector-$flag"
    ;;
  classes)
    fixture=input-ex-raw-highlight-two-byte-classes.txt
    stem="highlight-ex-two-byte-classes-$selector-$flag"
    ;;
  classes-replay)
    fixture=input-ex-raw-highlight-two-byte-classes.txt
    stem="highlight-ex-two-byte-classes-replay-$selector-$flag"
    ;;
  *)
    echo "unsupported fixture profile: $profile" >&2
    exit 2
    ;;
esac
work=/work/stage5
probe=/work/stage21
input_backup="$probe/$stem-original-input.txt"
output_backup="$probe/$stem-original-output.wav"
trace="$probe/trace-$stem.gdb"
log="$probe/$stem-api.log"
xvfb_log="$probe/xvfb-$stem.log"
dump="$probe/$stem-descriptor.bin"
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
cp "$probe/$fixture" "$work/input1.txt"
sed -e "s/@SELECTOR@/$selector/g" -e "s/@FLAG@/$flag/g" \
  -e "s|/work/stage21/highlight-ex-raw-mark-flag[01]-descriptor.bin|$dump|g" \
  "$probe/trace-highlight-ex-raw-mark.gdb.in" > "$trace"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
cd "$work"
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe < "$trace" > "$log" 2>&1
