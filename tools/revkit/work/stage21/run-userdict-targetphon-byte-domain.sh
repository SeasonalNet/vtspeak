#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
xpid=
restore() {
  cp /tmp/targetphon-byte-domain-input1.txt "$work/input1.txt"
  cp /tmp/targetphon-byte-domain-output.wav "$work/output.wav"
  cmp -s /tmp/targetphon-byte-domain-input1.txt "$work/input1.txt"
  cmp -s /tmp/targetphon-byte-domain-output.wav "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
cp "$work/input1.txt" /tmp/targetphon-byte-domain-input1.txt
cp "$work/output.wav" /tmp/targetphon-byte-domain-output.wav
trap restore EXIT
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-targetphon-byte-domain.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-userdict-targetphon-byte-domain.gdb" > "$probe/targetphon-byte-domain-api.log" 2>&1
grep '^TARGET_PHON_BYTE ' "$probe/targetphon-byte-domain-api.log" | tail -n 8
