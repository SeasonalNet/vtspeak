#!/bin/bash
set -euo pipefail

scratch=/work/corpus-parity
mode=${1:-smoke}
xpid=
restore() {
  if [ -n "$xpid" ]; then
    kill "$xpid" 2>/dev/null || true
  fi
}
trap restore EXIT

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$scratch/probe-xvfb.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1

if [ "$mode" = corpus ]; then
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 7200s \
    wine "$scratch/probe_dll_decoder.exe" --corpus > "$scratch/probe.log" 2>&1
elif [ "$mode" = smoke ]; then
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
    wine "$scratch/probe_dll_decoder.exe" > "$scratch/probe.log" 2>&1
else
  printf 'unsupported probe mode: %s\n' "$mode" >&2
  exit 2
fi
