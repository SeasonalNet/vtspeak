#!/bin/bash
set -euo pipefail

backup_dir=$(mktemp -d)
names=(prose numbers address)
for name in "${names[@]}"; do
  cp "/probe/runtime-original-msi-2006-${name}.wav" "$backup_dir/${name}.wav"
done

xpid=
restore() {
  if [[ -n "$xpid" ]]; then kill "$xpid" 2>/dev/null || true; fi
  for name in "${names[@]}"; do
    cp "$backup_dir/${name}.wav" "/probe/runtime-original-msi-2006-${name}.wav"
  done
}
trap restore EXIT

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > /probe/xvfb-license-success-trace-2006.log 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /probe/probe-original-kate-2006.exe \
  < /work/stage19/license-control/trace-parser.gdb \
  > /probe/runtime-license-success-2006-parser-gdb.log 2>&1
