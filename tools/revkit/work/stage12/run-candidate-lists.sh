#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage12
scratch=/work/corpus-parity
xpid=

mode=${1:?usage: run-candidate-lists.sh control|g83-zero fixture...}
shift
case "$mode" in
  control) trace="$probe/trace-candidate-lists.gdb" ;;
  g83-zero) trace="$probe/trace-candidate-lists-83-zero.gdb" ;;
  *) printf 'unknown mode: %s\n' "$mode" >&2; exit 2 ;;
esac

mkdir -p "$scratch"
cp "$work/input1.txt" "$scratch/stage12-candidate-lists-original-input1.txt"
cp "$work/output.wav" "$scratch/stage12-candidate-lists-original-output.wav"
restore() {
  cp "$scratch/stage12-candidate-lists-original-input1.txt" "$work/input1.txt"
  cp "$scratch/stage12-candidate-lists-original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then
    kill "$xpid" 2>/dev/null || true
  fi
}
trap restore EXIT

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/candidate-lists-xvfb.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1

for input in "$@"; do
  cp "$probe/inputs/$input.txt" "$work/input1.txt"
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
    winedbg --gdb /samples/voicetext_paul.exe \
    < "$trace" > "$probe/$input-candidate-lists-$mode.log" 2>&1
  cp "$work/output.wav" "$probe/$input-candidate-lists-$mode.wav"
done
