#!/bin/bash
set -euo pipefail

probe=/work/stage16
xpid=
cleanup() {
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap cleanup EXIT
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-license-probes.log" 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 120s \
  wine /work/stage16/probe-license-api.exe > "$probe/license-paths-api.log" 2>&1
sed -i 's/[[:blank:]]*$//' "$probe/license-paths-api.log"
