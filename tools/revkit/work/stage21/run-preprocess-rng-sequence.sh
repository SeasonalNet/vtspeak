#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
log="$probe/preprocess-rng-sequence-api.log"
xlog="$probe/xvfb-preprocess-rng-sequence.log"
xpid=

for output in "$log" "$xlog" \
  "$probe/preprocess-rng-r1" "$probe/preprocess-rng-r2" "$probe/preprocess-rng-r3"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 2
  fi
done

cp "$work/input1.txt" /tmp/preprocess-rng-sequence-input1.txt
cp "$work/output.wav" /tmp/preprocess-rng-sequence-output.wav
restore() {
  cp /tmp/preprocess-rng-sequence-input1.txt "$work/input1.txt"
  cp /tmp/preprocess-rng-sequence-output.wav "$work/output.wav"
  cmp -s /tmp/preprocess-rng-sequence-input1.txt "$work/input1.txt"
  cmp -s /tmp/preprocess-rng-sequence-output.wav "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT

cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xlog" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-preprocess-rng-sequence.gdb" > "$log" 2>&1
for output in "$probe/preprocess-rng-r1" "$probe/preprocess-rng-r2" "$probe/preprocess-rng-r3"; do
  if [ ! -s "$output" ]; then
    echo "missing or empty flag-7 output: $output" >&2
    exit 1
  fi
  prefix_len=$(awk 'NR == 1 { match($0, /^[.]+/); print RLENGTH; exit }' "$output")
  printf 'PREPROCESS_RNG_PREFIX file=%s periods=%s\n' "${output##*/}" "$prefix_len"
done
grep '^PREPROCESS_RNG ' "$log"
grep -q 'PREPROCESS_RNG_CALL n=3 raw=0x1' "$log"
grep -q '\[Inferior .*exited normally\]' "$log"
