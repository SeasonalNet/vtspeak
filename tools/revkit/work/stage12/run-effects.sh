#!/bin/bash
set -euo pipefail

work=${VTSPEAK_WORK_DIR:-/work/stage5}
probe=/work/stage12
capture_dir=${VTSPEAK_CAPTURE_DIR:-$probe}
scratch=${VTSPEAK_SCRATCH_DIR:-/work/corpus-parity}
xpid=
mode=${1:?usage: run-effects.sh control|effects-control|numeric-zero|numeric-83-zero|numeric-74-zero|numeric-swap|numeric-83-to95|numeric-consumers|numeric-consumers83|candidate-path-control|candidate-path-83-zero|acoustic-control|acoustic-83-zero|category-override|place-calls|place-calls-ax-zero|place-calls-ax-miss|place-gate-first-d|place-gate-second-d|place-gate-first-d-bit1|marker-producer|ax-one|ax-all-zero|ax-all-one|ax-all-miss fixture...}
shift

case "$mode" in
  control) trace="$probe/trace.gdb" ;;
  effects-control) trace="$probe/trace-effects-control.gdb" ;;
  numeric-zero) trace="$probe/trace-numeric-zero.gdb" ;;
  numeric-83-zero) trace="$probe/trace-numeric-83-zero.gdb" ;;
  numeric-74-zero) trace="$probe/trace-numeric-74-zero.gdb" ;;
  numeric-swap) trace="$probe/trace-numeric-swap.gdb" ;;
  numeric-83-to95) trace="$probe/trace-numeric-83-to95.gdb" ;;
  numeric-consumers) trace="$probe/trace-numeric-consumers.gdb" ;;
  numeric-consumers83) trace="$probe/trace-numeric-consumers83.gdb" ;;
  candidate-path-control) trace="$probe/trace-candidate-path-control.gdb" ;;
  candidate-path-83-zero) trace="$probe/trace-candidate-path-83-zero.gdb" ;;
  acoustic-control) trace="$probe/trace-acoustic-control.gdb" ;;
  acoustic-83-zero) trace="$probe/trace-acoustic-83-zero.gdb" ;;
  category-override) trace="$probe/trace-category-override.gdb" ;;
  place-calls) trace="$probe/trace-place-call.gdb" ;;
  place-calls-ax-zero) trace="$probe/trace-place-call-ax-zero.gdb" ;;
  place-calls-ax-miss) trace="$probe/trace-place-call-ax-miss.gdb" ;;
  place-gate-first-d) trace="$probe/trace-place-gate-first-d.gdb" ;;
  place-gate-second-d) trace="$probe/trace-place-gate-second-d.gdb" ;;
  place-gate-first-d-bit1) trace="$probe/trace-place-gate-first-d-bit1.gdb" ;;
  marker-producer) trace="$probe/trace-replacement-marker.gdb" ;;
  ax-one) trace="$probe/trace-ax-one.gdb" ;;
  ax-all-zero) trace="$probe/trace-ax-all-zero.gdb" ;;
  ax-all-one) trace="$probe/trace-ax-all-one.gdb" ;;
  ax-all-miss) trace="$probe/trace-ax-all-miss.gdb" ;;
  *) printf 'unknown mode: %s\n' "$mode" >&2; exit 2 ;;
esac

mkdir -p "$scratch" "$capture_dir"
cd "$work"
cp "$work/input1.txt" "$scratch/stage12-effects-original-input1.txt"
cp "$work/output.wav" "$scratch/stage12-effects-original-output.wav"
restore() {
  cp "$scratch/stage12-effects-original-input1.txt" "$work/input1.txt"
  cp "$scratch/stage12-effects-original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then
    kill "$xpid" 2>/dev/null || true
  fi
}
trap restore EXIT

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$capture_dir/effects-xvfb.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1

for input in "$@"; do
  cp "$probe/inputs/$input.txt" "$work/input1.txt"
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
    winedbg --gdb /samples/voicetext_paul.exe \
    < "$trace" > "$capture_dir/$input-effect-$mode.log" 2>&1
  cp "$work/output.wav" "$capture_dir/$input-effect-$mode.wav"
done
