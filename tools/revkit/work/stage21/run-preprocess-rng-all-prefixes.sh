#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
log="$probe/preprocess-rng-all-prefixes-api.log"
xlog="$probe/xvfb-preprocess-rng-all-prefixes.log"
backup_input=/tmp/preprocess-rng-all-prefixes-input1.txt
backup_output=/tmp/preprocess-rng-all-prefixes-output.wav
xpid=

for output in "$log" "$xlog"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 2
  fi
done
for index in {0..7}; do
  if [ -e "$probe/p$index" ]; then
    echo "refusing to overwrite $probe/p$index" >&2
    exit 2
  fi
done

cp "$work/input1.txt" "$backup_input"
cp "$work/output.wav" "$backup_output"
restore() {
  cp "$backup_input" "$work/input1.txt"
  cp "$backup_output" "$work/output.wav"
  cmp -s "$backup_input" "$work/input1.txt"
  cmp -s "$backup_output" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT

cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xlog" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 240s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-preprocess-rng-all-prefixes.gdb" > "$log" 2>&1

periods=(108 111 140 52 33 151 106 58)
states=(134456 117649 100842 84035 67228 50421 33614 16807)
for index in {0..7}; do
  actual=$(awk 'NR == 1 { match($0, /^[.]+/); print RLENGTH; exit }' "$probe/p$index")
  if [ "$actual" != "${periods[index]}" ]; then
    echo "index $index expected ${periods[index]} periods, got $actual" >&2
    exit 1
  fi
  grep -q "PREPROCESS_RNG_PREFIX index=$index .*state=${states[index]} raw=0x1 path=Z:/work/stage21/p$index" "$log"
  printf 'PREPROCESS_RNG_PREFIX_CHECK index=%d state=%d periods=%d\n' \
    "$index" "${states[index]}" "$actual"
done
grep -q '\[Inferior .*exited normally\]' "$log"
