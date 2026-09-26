#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/probe
input_backup=/tmp/lead5-input-original
output_backup=/tmp/lead5-output-original
xpid=

cp "$work/input1.txt" "$input_backup"
cp "$work/output.wav" "$output_backup"
restore() {
  cp "$input_backup" "$work/input1.txt"
  cp "$output_backup" "$work/output.wav"
  if [ -n "$xpid" ]; then
    kill "$xpid" 2>/dev/null || true
  fi
}
trap restore EXIT

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp >"$probe/xvfb.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1

printf 'voice\tfixture\texit\toutput_bytes\tresult\n' > "$probe/runtime-matrix.tsv"
for voice in paul james julie kate; do
  for fixture in prose numbers address; do
    if { [ "$voice" = julie ] || [ "$voice" = kate ]; } && [ "$fixture" != prose ]; then
      continue
    fi
    case "$fixture" in
      prose) input="$probe/input.txt" ;;
      numbers) input="$probe/numbers.txt" ;;
      address) input="$probe/address.txt" ;;
    esac
    output_name="$voice"
    if [ "$fixture" != prose ]; then
      output_name="$voice-$fixture"
    fi
    cp "$input" "$work/input1.txt"
    : > "$work/output.wav"
    if [ "$voice" = julie ] || [ "$voice" = kate ]; then
      debug=+file
    else
      debug=-all
    fi
    set +e
    WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG="$debug" timeout 180s \
      wine "/samples/voicetext_${voice}.exe" >"$probe/${voice}-${fixture}.log" 2>&1
    rc=$?
    set -e
    bytes=$(stat -c %s "$work/output.wav")
    result=failed
    if [ "$rc" -eq 0 ] && [ "$bytes" -gt 44 ]; then
      cp "$work/output.wav" "$probe/${output_name}.wav"
      result=wav
    fi
    printf '%s\t%s\t%s\t%s\t%s\n' "$voice" "$fixture" "$rc" "$bytes" "$result" \
      >> "$probe/runtime-matrix.tsv"
  done
done

(cd "$probe" && sha256sum ./*.wav > wav-sha256.txt)
