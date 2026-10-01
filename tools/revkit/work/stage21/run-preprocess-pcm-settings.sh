#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
log="$probe/preprocess-pcm-settings-api.log"
xlog="$probe/xvfb-preprocess-pcm-settings.log"
xpid=

for output in "$log" "$xlog" "$work/test.pcm"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 2
  fi
done
for variant in {0..10}; do
  name=$(printf '%02d' "$variant")
  if [ -e "$probe/preprocess-pcm-settings-v$name.pcm" ]; then
    echo "refusing to overwrite preprocess-pcm-settings-v$name.pcm" >&2
    exit 2
  fi
done

cp "$work/input1.txt" /tmp/preprocess-pcm-settings-input1.txt
cp "$work/output.wav" /tmp/preprocess-pcm-settings-output.wav
restore() {
  cp /tmp/preprocess-pcm-settings-input1.txt "$work/input1.txt"
  cp /tmp/preprocess-pcm-settings-output.wav "$work/output.wav"
  cmp -s /tmp/preprocess-pcm-settings-input1.txt "$work/input1.txt"
  cmp -s /tmp/preprocess-pcm-settings-output.wav "$work/output.wav"
  if [ -e "$work/test.pcm" ] && [ ! -e "$probe/preprocess-pcm-settings-partial.pcm" ]; then
    mv "$work/test.pcm" "$probe/preprocess-pcm-settings-partial.pcm"
  fi
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
  < "$probe/trace-preprocess-pcm-settings.gdb" > "$log" 2>&1

if [ "$(grep -c '^PREPROCESS_PCM_SETTINGS ' "$log")" -ne 11 ]; then
  echo "expected 11 API calls in $log" >&2
  exit 1
fi
grep -q '\[Inferior .*exited normally\]' "$log"
for variant in {0..10}; do
  name=$(printf '%02d' "$variant")
  if [ ! -s "$probe/preprocess-pcm-settings-v$name.pcm" ]; then
    echo "missing or empty PCM output: preprocess-pcm-settings-v$name.pcm" >&2
    exit 1
  fi
done
grep '^PREPROCESS_PCM_SETTINGS ' "$log"
