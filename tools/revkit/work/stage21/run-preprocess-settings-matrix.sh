#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
log="$probe/preprocess-settings-matrix-api.log"
xlog="$probe/xvfb-preprocess-settings-matrix.log"
xpid=

for output in "$log" "$xlog"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 2
  fi
done
for flag in 3 5 7 a; do
  for variant in {0..10}; do
    name=$(printf '%02d' "$variant")
    if [ -e "$probe/ptf${flag}v$name" ]; then
      echo "refusing to overwrite $probe/ptf${flag}v$name" >&2
      exit 2
    fi
  done
done

cp "$work/input1.txt" /tmp/preprocess-settings-matrix-input1.txt
cp "$work/output.wav" /tmp/preprocess-settings-matrix-output.wav
restore() {
  cp /tmp/preprocess-settings-matrix-input1.txt "$work/input1.txt"
  cp /tmp/preprocess-settings-matrix-output.wav "$work/output.wav"
  cmp -s /tmp/preprocess-settings-matrix-input1.txt "$work/input1.txt"
  cmp -s /tmp/preprocess-settings-matrix-output.wav "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT

cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xlog" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 300s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-preprocess-settings-matrix.gdb" > "$log" 2>&1

if [ "$(grep -c '^PREPROCESS_SETTINGS ' "$log")" -ne 44 ]; then
  echo "expected 44 API calls in $log" >&2
  exit 1
fi
grep -q '\[Inferior .*exited normally\]' "$log"
for flag in 3 5 7 a; do
  for variant in {0..10}; do
    name=$(printf '%02d' "$variant")
    if [ ! -f "$probe/ptf${flag}v$name" ]; then
      echo "missing output: $probe/ptf${flag}v$name" >&2
      exit 1
    fi
  done
done
grep '^PREPROCESS_SETTINGS ' "$log"
