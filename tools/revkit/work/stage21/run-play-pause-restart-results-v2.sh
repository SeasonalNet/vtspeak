#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
log="$probe/play-pause-restart-results-v2-api.log"
xlog="$probe/xvfb-play-pause-restart-results-v2.log"
xpid=
for output in "$log" "$xlog"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 2
  fi
done
restore() {
  cp /tmp/play-stop-lifecycle-input1.txt "$work/input1.txt"
  cp /tmp/play-stop-lifecycle-output.wav "$work/output.wav"
  cmp -s /tmp/play-stop-lifecycle-input1.txt "$work/input1.txt"
  cmp -s /tmp/play-stop-lifecycle-output.wav "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
cp "$work/input1.txt" /tmp/play-stop-lifecycle-input1.txt
cp "$work/output.wav" /tmp/play-stop-lifecycle-output.wav
trap restore EXIT
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xlog" 2>&1 &
xpid=$!
export DISPLAY=:99
export ALSA_CONFIG_PATH="$probe/../stage16/alsa-null.conf"
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-play-pause-restart-results-v2.gdb" > "$log" 2>&1
test "$(grep -c 'PAUSE_WAVEOUT_RESULT mmresult=0 ' "$log")" -eq 2
test "$(grep -c 'RESTART_WAVEOUT_RESULT mmresult=0 ' "$log")" -eq 2
grep -E '^(PLAY_|STOP_|IDLE_|PAUSE_WAVEOUT_RESULT|RESTART_WAVEOUT_RESULT)' "$log"
