#!/bin/bash
set -euo pipefail

probe=/work/stage21
for output in "$probe/playback-natural-completion-host-api.log" \
              "$probe/xvfb-playback-natural-completion.log"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
xpid=
restore() {
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-playback-natural-completion.log" 2>&1 &
xpid=$!
export DISPLAY=:99
export ALSA_CONFIG_PATH="$probe/../stage16/alsa-null.conf"
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 45s \
  wine /work/stage21/probe-playback-natural-completion.exe \
  > "$probe/playback-natural-completion-host-api.log" 2>&1
grep -E '^(PROBE_|LOAD |PLAY |STATE phase=(loaded_idle|play_return|after_dispatch|before_stop|after_stop)|MESSAGE |PLAYBACK_END)' \
  "$probe/playback-natural-completion-host-api.log"
