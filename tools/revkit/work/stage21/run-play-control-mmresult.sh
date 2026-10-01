#!/bin/bash
set -euo pipefail

probe=/work/stage21
work=/work/stage21/sandbox/stage5
xpid=
restore() {
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT

test -f "$work/input1.txt"
test -f "$work/output.wav"
for output in "$probe/play-control-mmresult-v4-api.log" "$probe/xvfb-play-control-mmresult-v4.log"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-play-control-mmresult-v4.log" 2>&1 &
xpid=$!
export DISPLAY=:99
export ALSA_CONFIG_PATH="$probe/../stage16/alsa-null.conf"
sleep 1
cd "$work"
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-play-control-mmresult-v4.gdb" > "$probe/play-control-mmresult-v4-api.log" 2>&1
