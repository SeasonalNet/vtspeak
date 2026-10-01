#!/bin/bash
set -euo pipefail

probe=/work/stage21
log="$probe/playback-span-boundaries-api.log"
for output in "$log" "$probe/xvfb-playback-span-boundaries.log"; do
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

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-playback-span-boundaries.log" 2>&1 &
xpid=$!
export DISPLAY=:99
export ALSA_CONFIG_PATH="$probe/../stage16/alsa-null.conf"
sleep 1
for text in "A A A" "Hello, world!" "hex:636166e9206e6f6972"; do
  printf '\n=== INPUT %s ===\n' "$text" >> "$log"
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 45s \
    wine /work/stage21/probe-playback-natural-completion.exe "$text" \
    >> "$log" 2>&1
done
grep -E '^(=== INPUT|PROBE_BEGIN|LOAD |PLAY |STATE phase=(loaded_idle|play_return|after_dispatch|before_stop|after_stop)|MESSAGE |PLAYBACK_END|PROBE_END)' \
  "$log"
