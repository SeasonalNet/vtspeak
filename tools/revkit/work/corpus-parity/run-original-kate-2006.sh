#!/bin/bash
set -euo pipefail
xpid=
cleanup() { if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi; }
trap cleanup EXIT
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > /probe/xvfb-original-kate-2006.log 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=${WINEDEBUG:--all} timeout 120s wine /probe/probe-original-kate-2006.exe > /probe/runtime-original-msi-2006-api.log 2>&1
