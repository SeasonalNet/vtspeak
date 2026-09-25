#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage14
scratch=/work/corpus-parity/lead4-stage14
trace=/work/stage10/trace-context.gdb
xpid=
mkdir -p "$scratch"
cp "$work/input1.txt" "$scratch/original-input1.txt"
cp "$work/output.wav" "$scratch/original-output.wav"
restore() {
  cp "$scratch/original-input1.txt" "$work/input1.txt"
  cp "$scratch/original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then
    kill "$xpid" 2>/dev/null || true
  fi
}
trap restore EXIT

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$scratch/xvfb.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1

if [ "$#" -eq 0 ]; then
  mapfile -t fixtures < "$probe/lead4-fixtures.txt"
  inputs=()
  for fixture in "${fixtures[@]}"; do
    inputs+=("$probe/inputs/$fixture.txt")
  done
else
  inputs=()
  for name in "$@"; do
    input="$probe/inputs/$name.txt"
    if [ ! -f "$input" ]; then
      printf 'no Lead 4 input found: %s\n' "$name" >&2
      exit 2
    fi
    inputs+=("$input")
  done
fi

for input in "${inputs[@]}"; do
  name=${input##*/}
  name=${name%.txt}
  cp "$input" "$work/input1.txt"
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
    winedbg --gdb /samples/voicetext_paul.exe \
    < "$trace" > "$probe/$name.log" 2>&1
done
