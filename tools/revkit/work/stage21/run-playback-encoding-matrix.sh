#!/bin/bash
set -euo pipefail

probe=/work/stage21
log="$probe/playback-encoding-matrix-api.log"
xlog="$probe/xvfb-playback-encoding-matrix.log"
for output in "$log" "$xlog"; do
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

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xlog" 2>&1 &
xpid=$!
export DISPLAY=:99
export ALSA_CONFIG_PATH="$probe/../stage16/alsa-null.conf"
sleep 1
for case in \
  CP1252_CAFE:hex:636166e9206e6f6972 \
  UTF8_CAFE:hex:636166c3a9206e6f6972 \
  UTF8_E_ACUTE_WORD:hex:c3a9206e6f6972; do
  label=${case%%:*}
  input=${case#*:}
  printf '\n=== CASE %s ===\n' "$label" >> "$log"
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 45s \
    wine /work/stage21/probe-playback-natural-completion.exe "$input" \
    >> "$log" 2>&1
done
grep -E '^(=== CASE|PROBE_BEGIN|LOAD |PLAY |STATE phase=(loaded_idle|play_return|after_dispatch|before_stop|after_stop)|MESSAGE |PLAYBACK_END|PROBE_END)' \
  "$log"
