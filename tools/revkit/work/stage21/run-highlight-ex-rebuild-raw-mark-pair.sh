#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
input_backup="$probe/highlight-ex-rebuild-original-input.txt"
output_backup="$probe/highlight-ex-rebuild-original-output.wav"
xpid=
restore() {
  cp "$input_backup" "$work/input1.txt"
  cp "$output_backup" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

for path in "$input_backup" "$output_backup"; do
  if [ -e "$path" ]; then
    echo "refusing to overwrite $path" >&2
    exit 3
  fi
done
cp "$work/input1.txt" "$input_backup"
cp "$work/output.wav" "$output_backup"
trap restore EXIT

for flag in 0 1; do
  trace="$probe/trace-highlight-ex-raw-mark-rebuild-$flag.gdb"
  log="$probe/highlight-ex-raw-mark-rebuild-$flag-api.log"
  xvfb_log="$probe/xvfb-highlight-ex-raw-mark-rebuild-$flag.log"
  target="$probe/highlight-ex-raw-mark-flag$flag-descriptor.bin"
  for path in "$trace" "$log" "$xvfb_log"; do
    if [ -e "$path" ]; then
      echo "refusing to overwrite $path" >&2
      exit 3
    fi
  done
  cp /work/stage16/input-ex-records-vtml-mark.txt "$work/input1.txt"
  sed -e "s/@SELECTOR@/0/g" -e "s/@FLAG@/$flag/g" \
    -e "s|/work/stage21/highlight-ex-raw-mark-flag[01]-descriptor.bin|$target|g" \
    "$probe/trace-highlight-ex-raw-mark.gdb.in" > "$trace"
  Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
  xpid=$!
  export DISPLAY=:99
  sleep 1
  cd "$work"
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
    winedbg --gdb /samples/voicetext_paul.exe < "$trace" > "$log" 2>&1
  kill "$xpid" 2>/dev/null || true
  xpid=
done
