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

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > /probe/xvfb-license-status-only-2006.log 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99
run_status=0
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 120s \
  wine /probe/probe-original-kate-2006.exe \
  > /probe/runtime-license-status-only-2006-api.log 2>&1 || run_status=$?
printf 'probe_exit=%s\n' "$run_status" >> /probe/runtime-license-status-only-2006-api.log
if [[ $run_status == 0 ]]; then
  for name in "${names[@]}"; do
    cp "/probe/runtime-original-msi-2006-${name}.wav" \
      "/probe/runtime-license-status-only-2006-${name}.wav"
  done
fi
