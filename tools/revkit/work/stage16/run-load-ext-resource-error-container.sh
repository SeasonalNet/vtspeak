#!/bin/bash
set -euo pipefail

probe=/work/stage16
case_name=${1:?usage: run-load-ext-resource-error-container.sh case}
case "$case_name" in
  omit-gen-dat|omit-gen-upm|omit-unit-gen-idx|omit-dblist-idx-preserve-tree|omit-cepdist-tbl|omit-pitch-nbt|omit-atmt-tree|omit-engbi-tree|omit-sbd-tree|omit-tppdict|omit-hashidx-tpp|empty-dblist|empty-gen-dat) ;;
  *) echo "unsupported resource overlay: $case_name" >&2; exit 2 ;;
esac
log="$probe/load-ext-resource-error-$case_name-api.log"
xvfb_log="$probe/xvfb-load-ext-resource-error-$case_name.log"
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
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  wine "$probe/probe-load-ext-state.exe" 1 > "$log" 2>&1
sed -i 's/[[:blank:]]*$//' "$log"
