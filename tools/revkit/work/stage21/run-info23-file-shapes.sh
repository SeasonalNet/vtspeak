#!/bin/bash
set -euo pipefail

probe=/work/stage21
date_file=/work/db_build.date
trace="$probe/trace-info23-file-shape.gdb"
embedded_log="$probe/info23-embedded-nul-api.log"
multiline_log="$probe/info23-multiline-api.log"
xvfb_log="$probe/xvfb-info23-file-shapes.log"
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
for output in "$embedded_log" "$multiline_log" "$xvfb_log"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done

printf 'A\0B\nC\n' > "$date_file"
created_date_file=1
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" > "$embedded_log" 2>&1

printf 'LINE1\nLINE2\n' > "$date_file"
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" > "$multiline_log" 2>&1
