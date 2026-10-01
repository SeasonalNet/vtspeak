#!/bin/bash
set -euo pipefail

probe=/work/stage21
run_id=${1:-v1}
log="$probe/speakersinfo-guard-setup-$run_id-api.log"
xvfb_log="$probe/xvfb-speakersinfo-guard-setup-$run_id.log"
xpid=
restore() {
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT

for output in "$log" "$xvfb_log"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-speakersinfo-guard-setup-v1.gdb" > "$log" 2>&1
grep -Eq 'SPEAKERSINFO_GUARD_SETUP allocation=0x[0-9a-f]+ guard=0x[0-9a-f]+ protect_result=1 ' "$log"
grep -Fq 'SPEAKERSINFO_GUARD_RESTORE result=1' "$log"
