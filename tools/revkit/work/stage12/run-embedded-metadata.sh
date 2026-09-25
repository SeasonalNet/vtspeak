#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage12
scratch=$(mktemp -d /tmp/lead3-embedded.XXXXXX)
trace=${TRACE_SCRIPT:-$probe/trace-embedded-metadata.gdb}
suffix=${TRACE_SUFFIX:--embedded-metadata}
xpid=
cp "$work/input1.txt" "$scratch/input1.txt"
cp "$work/output.wav" "$scratch/output.wav"
input_hash=$(sha256sum "$scratch/input1.txt" | cut -d' ' -f1)
output_hash=$(sha256sum "$scratch/output.wav" | cut -d' ' -f1)
restore() {
  cp "$scratch/input1.txt" "$work/input1.txt"
  cp "$scratch/output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then
    kill "$xpid" 2>/dev/null || true
  fi
  test "$(sha256sum "$work/input1.txt" | cut -d' ' -f1)" = "$input_hash"
  test "$(sha256sum "$work/output.wav" | cut -d' ' -f1)" = "$output_hash"
}
trap restore EXIT

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > /dev/null 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1

for name in "$@"; do
  cp "$probe/inputs/$name.txt" "$work/input1.txt"
  debugger_rc=0
  set +e
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
    winedbg --gdb /samples/voicetext_paul.exe \
    < "$trace" > "$probe/$name$suffix.log" 2>&1
  debugger_rc=$?
  set -e
  result_file="$probe/$name$suffix.result"
  printf "Debugger wrapper status: %s\n" "$debugger_rc" > "$result_file"
  if [ -f "$work/output.wav" ]; then
    printf "WAVE bytes: " >> "$result_file"
    wc -c < "$work/output.wav" >> "$result_file"
    sha256sum "$work/output.wav" >> "$result_file"
  else
    printf "WAVE file missing\n" >> "$result_file"
  fi
done
