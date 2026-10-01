#!/bin/bash
set -euo pipefail

probe=/work/stage21
run_id=${1:-v1}
if [[ ! "$run_id" =~ ^[A-Za-z0-9_-]+$ ]]; then
  echo "run id may contain only letters, digits, underscores, and hyphens" >&2
  exit 2
fi
trace="$probe/trace-speaker-metadata-wrap-alias-v1.gdb"
log="$probe/speaker-metadata-wrap-alias-$run_id-api.log"
xvfb_log="$probe/xvfb-speaker-metadata-wrap-alias-$run_id.log"
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
  < "$trace" > "$log" 2>&1

test "$(grep -c '^PATHKEY_WRAP ' "$log")" -eq 42
test "$(grep -c '^SPEAKERSINFO_WRAP ' "$log")" -eq 42
keys=(
  'SOFTWARE\VW\VT\Kate\M16'
  'SOFTWARE\VW\VT\Paul\M16'
  'SOFTWARE\VW\VT\em001\M16'
  'SOFTWARE\VW\VT\Julie\M16'
  'SOFTWARE\VW\VT\James\M16'
  'SOFTWARE\VW\VT\Ashley\M16'
)
ids=(kate paul em001 julie james ashley)
paths=(
  d:/eng/db/susan/pcm/
  d:/eng/db/isaac/pcm/
  d:/eng/db/em001/pcm/
  d:/eng/db/jennifer/pcm/
  d:/eng/db/lee/pcm/
  d:/eng/db/casey/pcm/
)
for k in {1..7}; do
  for slot in {0..5}; do
    selector=$((slot + k * 536870912))
    if [ "$selector" -ge 2147483648 ]; then
      selector=$((selector - 4294967296))
    fi
    grep -Fq "PATHKEY_WRAP k=$k slot=$slot selector=$selector value=${keys[$slot]}" "$log"
    grep -Fq "SPEAKERSINFO_WRAP k=$k slot=$slot selector=$selector result=6 name=${ids[$slot]} path=${paths[$slot]}" "$log"
  done
done
