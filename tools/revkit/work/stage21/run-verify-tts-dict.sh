#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
xpid=
restore() {
  cp /tmp/verify-tts-dict-input1.txt "$work/input1.txt"
  cp /tmp/verify-tts-dict-output.wav "$work/output.wav"
  cmp -s /tmp/verify-tts-dict-input1.txt "$work/input1.txt"
  cmp -s /tmp/verify-tts-dict-output.wav "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
cp "$work/input1.txt" /tmp/verify-tts-dict-input1.txt
cp "$work/output.wav" /tmp/verify-tts-dict-output.wav
trap restore EXIT
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-verify-tts-dict.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-verify-tts-dict.gdb" > "$probe/verify-tts-dict-api.log" 2>&1
grep '^VERIFY_TTS_DICT' "$probe/verify-tts-dict-api.log"
