#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/probe
fixture=${FIXTURE:?set FIXTURE to prose, numbers, or address}
prefix=${OUTPUT_PREFIX:?set OUTPUT_PREFIX to a result label}
case "$fixture" in
  prose) source_fixture=input ;;
  numbers|address) source_fixture=$fixture ;;
  *) printf 'unsupported fixture: %s\n' "$fixture" >&2; exit 2 ;;
esac
input_backup=$(mktemp "$probe/.stage5-single-input.XXXXXX")
output_backup=$(mktemp "$probe/.stage5-single-output.XXXXXX")
xpid=

cp "$work/input1.txt" "$input_backup"
cp "$work/output.wav" "$output_backup"
restore() {
  status=$?
  cp "$input_backup" "$work/input1.txt"
  cp "$output_backup" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
  rm -f "$input_backup" "$output_backup"
  exit "$status"
}
trap restore EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

cp "/work/stage17/$source_fixture.txt" "$work/input1.txt"
: > "$work/output.wav"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp >"$probe/xvfb-$prefix.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
set +e
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  wine /samples/voicetext_kate.exe >"$probe/$prefix.log" 2>&1
rc=$?
set -e
bytes=$(stat -c %s "$work/output.wav")
hash=$(sha256sum "$work/output.wav" | cut -d ' ' -f 1)
result=failed
if [ "$rc" -eq 0 ] && [ "$bytes" -gt 44 ]; then
  cp "$work/output.wav" "$probe/$prefix.wav"
  result=wav
fi
printf 'fixture\texit\toutput_bytes\tsha256\tresult\n%s\t%s\t%s\t%s\t%s\n' \
  "$fixture" "$rc" "$bytes" "$hash" "$result" > "$probe/$prefix-matrix.tsv"
if [ "$result" != wav ]; then exit 1; fi
