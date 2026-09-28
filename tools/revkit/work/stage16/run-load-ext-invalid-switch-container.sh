#!/bin/bash
set -euo pipefail

probe=/work/stage16
slot=${1:?usage: run-load-ext-invalid-switch-container.sh slot}
case "$slot" in
  -2|-1|6|1000) ;;
  *) echo "unsupported isolated slot: $slot" >&2; exit 2 ;;
esac
log="$probe/load-ext-invalid-switch-recovery-$slot-api.log"
xvfb_log="$probe/xvfb-load-ext-invalid-switch-recovery-$slot.log"
xpid=
cleanup() {
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap cleanup EXIT
for output in "$log" "$xvfb_log"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99
cd /work
set +e
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  wine "$probe/probe-load-ext-state.exe" invalid-switch "$slot" > "$log" 2>&1
status=$?
set -e
printf 'PROCESS_STATUS slot=%s status=%s\n' "$slot" "$status" >> "$log"
sed -i 's/[[:blank:]]*$//' "$log"
