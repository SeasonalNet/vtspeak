#!/bin/bash
set -euo pipefail
xpid=
cleanup() { if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi; }
trap cleanup EXIT
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > /probe/xvfb-original-kate-parser.log 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /probe/probe-original-kate-2006.exe \
  < /probe/trace-original-kate-2006-parser.gdb \
  > /probe/runtime-original-msi-2006-parser-gdb.log 2>&1
