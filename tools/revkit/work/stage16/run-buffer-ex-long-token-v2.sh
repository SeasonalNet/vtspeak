#!/bin/bash
set -euo pipefail

if [ "$#" -ne 1 ] || ! [[ "$1" =~ ^[0-2]$ ]]; then
  echo "usage: $0 SELECTOR" >&2
  exit 2
fi
selector=$1
probe=/work/stage16
trace="$probe/trace-buffer-ex-long-token-v2-$selector.gdb"
log="$probe/buffer-ex-long-token-v2-$selector.log"
xvfb_log="$probe/xvfb-buffer-ex-long-token-v2-$selector.log"
xpid=
if [ -e "$trace" ] || [ -e "$log" ] || [ -e "$xvfb_log" ]; then
  echo "refusing to overwrite long-token v2 probe output" >&2
  exit 3
fi
sed "s/@SELECTOR@/$selector/g" \
  "$probe/trace-buffer-ex-long-token-v2.gdb.in" > "$trace"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
trap 'kill "$xpid" 2>/dev/null || true' EXIT
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" > "$log" 2>&1
