#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/probe
input_backup=$(mktemp "$probe/.stage5-adapter-input.XXXXXX")
output_backup=$(mktemp "$probe/.stage5-adapter-output.XXXXXX")
xpid=

cp "$work/input1.txt" "$input_backup"
cp "$work/output.wav" "$output_backup"
initial_output_sha256=$(sha256sum "$output_backup" | cut -d ' ' -f 1)
restore() {
  cp "$input_backup" "$work/input1.txt"
  cp "$output_backup" "$work/output.wav"
  if [ -n "$xpid" ]; then
    kill "$xpid" 2>/dev/null || true
  fi
  rm -f "$input_backup" "$output_backup"
}
trap restore EXIT

cp "$probe/input.txt" "$work/input1.txt"
: > "$work/output.wav"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp >"$probe/xvfb-versioned.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1

set +e
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=+file,+loaddll timeout 600s \
  wine /samples/voicetext_kate.exe >"$probe/versioned-runtime.log" 2>&1
rc=$?
set -e
bytes=$(stat -c %s "$work/output.wav")
observed_output_sha256=$(sha256sum "$work/output.wav" | cut -d ' ' -f 1)
new_output=no
if [ "$bytes" -gt 44 ] && [ "$observed_output_sha256" != "$initial_output_sha256" ]; then
  cp "$work/output.wav" "$probe/versioned-new-output.wav"
  new_output=yes
fi
printf 'exit=%s\noutput_bytes=%s\noutput_sha256=%s\ninitial_output_sha256=%s\nnew_output=%s\n' \
  "$rc" "$bytes" "$observed_output_sha256" "$initial_output_sha256" "$new_output" \
  > "$probe/versioned-runtime-result.txt"
