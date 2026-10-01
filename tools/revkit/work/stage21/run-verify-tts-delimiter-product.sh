#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
xpid=
restore() {
  cp /tmp/verify-tts-delimiter-product-input1.txt "$work/input1.txt"
  cp /tmp/verify-tts-delimiter-product-output.wav "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

cp "$work/input1.txt" /tmp/verify-tts-delimiter-product-input1.txt
cp "$work/output.wav" /tmp/verify-tts-delimiter-product-output.wav
trap restore EXIT
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-verify-tts-delimiter-product.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-verify-tts-delimiter-product.gdb" > "$probe/verify-tts-delimiter-product-api.log" 2>&1
grep '^VERIFY_TTS_DELIMITER_PRODUCT ' "$probe/verify-tts-delimiter-product-api.log" | tail -n 4
