#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage10
scratch=/work/corpus-parity
trace=${TRACE_SCRIPT:-$probe/trace-context.gdb}
suffix=${TRACE_SUFFIX:-}
xpid=
mkdir -p "$scratch"
cp "$work/input1.txt" "$scratch/stage10-original-input1.txt"
cp "$work/output.wav" "$scratch/stage10-original-output.wav"
restore() {
  cp "$scratch/stage10-original-input1.txt" "$work/input1.txt"
  cp "$scratch/stage10-original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then
    kill "$xpid" 2>/dev/null || true
  fi
}
trap restore EXIT

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$scratch/stage10-xvfb.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1

if [ "$#" -eq 0 ]; then
  inputs=("$probe"/inputs/*.txt)
else
  shopt -s nullglob
  inputs=()
  for pattern in "$@"; do
    prior_count=${#inputs[@]}
    for input in "$probe"/inputs/$pattern.txt; do
      if [ -f "$input" ]; then
        inputs+=("$input")
      fi
    done
    if [ "${#inputs[@]}" -eq "$prior_count" ]; then
      printf 'no Stage 10 input matches: %s\n' "$pattern" >&2
      exit 2
    fi
  done
fi

for input in "${inputs[@]}"; do
  name=${input##*/}
  name=${name%.txt}
  cp "$input" "$work/input1.txt"
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
    winedbg --gdb /samples/voicetext_paul.exe \
    < "$trace" > "$probe/$name$suffix.log" 2>&1
done
