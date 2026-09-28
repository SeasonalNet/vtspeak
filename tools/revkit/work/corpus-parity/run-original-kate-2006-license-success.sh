#!/bin/bash
set -euo pipefail

backup_dir=$(mktemp -d)
saved=(0 0 0)
names=(prose numbers address)
for i in 0 1 2; do
  name=${names[$i]}
  source="/probe/runtime-original-msi-2006-${name}.wav"
  if [[ -f "$source" ]]; then
    cp "$source" "$backup_dir/${name}.wav"
    saved[$i]=1
  fi
done

xpid=
restore() {
  if [[ -n "$xpid" ]]; then kill "$xpid" 2>/dev/null || true; fi
  for i in 0 1 2; do
    name=${names[$i]}
    target="/probe/runtime-original-msi-2006-${name}.wav"
    if [[ ${saved[$i]} == 1 ]]; then
      cp "$backup_dir/${name}.wav" "$target"
    fi
  done
}
trap restore EXIT

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > /probe/xvfb-license-success-2006.log 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99
run_status=0
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=${WINEDEBUG:--all} timeout 120s \
  wine /probe/probe-original-kate-2006.exe \
  > /probe/runtime-license-success-2006-api.log 2>&1 || run_status=$?

if [[ $run_status != 0 ]]; then
  printf 'patched probe exit status=%s\n' "$run_status" \
    >> /probe/runtime-license-success-2006-api.log
  exit "$run_status"
fi

for name in "${names[@]}"; do
  cp "/probe/runtime-original-msi-2006-${name}.wav" \
    "/probe/runtime-license-success-2006-${name}.wav"
done
