#!/bin/bash
set -euo pipefail

probe=/work/stage21
date_file=/work/db_build.date
trace="$probe/trace-info23-capacity-edges.gdb"
log="$probe/info23-capacity-edges-api.log"
xvfb_log="$probe/xvfb-info23-capacity-edges.log"
xpid=
created_date_file=0
cleanup() {
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
  if [ "$created_date_file" -eq 1 ]; then rm -f "$date_file"; fi
}
trap cleanup EXIT

if [ -e "$date_file" ]; then
  echo "refusing to overwrite $date_file" >&2
  exit 3
fi
for output in "$log" "$xvfb_log"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
printf 'D23-OK\n' > "$date_file"
created_date_file=1

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" > "$log" 2>&1
