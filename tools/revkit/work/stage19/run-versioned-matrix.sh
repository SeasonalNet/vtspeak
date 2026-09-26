#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/probe
input_backup=$(mktemp "$probe/.stage5-matrix-input.XXXXXX")
output_backup=$(mktemp "$probe/.stage5-matrix-output.XXXXXX")
xpid=
failed=no
output_prefix=${OUTPUT_PREFIX:-runtime}

cp "$work/input1.txt" "$input_backup"
cp "$work/output.wav" "$output_backup"
restore() {
  cp "$input_backup" "$work/input1.txt"
  cp "$output_backup" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
  rm -f "$input_backup" "$output_backup"
}
trap restore EXIT

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp >"$probe/xvfb-matrix.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1

printf 'fixture\texit\toutput_bytes\tsha256\tresult\n' > "$probe/$output_prefix-matrix.tsv"
for fixture in prose numbers address; do
  source_fixture=$fixture
  if [ "$fixture" = prose ]; then source_fixture=input; fi
  cp "/work/stage17/$source_fixture.txt" "$work/input1.txt"
  : > "$work/output.wav"
  set +e
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
    wine /samples/voicetext_kate.exe >"$probe/$output_prefix-$fixture.log" 2>&1
  rc=$?
  set -e
  bytes=$(stat -c %s "$work/output.wav")
  hash=$(sha256sum "$work/output.wav" | cut -d ' ' -f 1)
  result=failed
  if [ "$rc" -eq 0 ] && [ "$bytes" -gt 44 ]; then
    cp "$work/output.wav" "$probe/$output_prefix-$fixture.wav"
    result=wav
  else
    failed=yes
  fi
  printf '%s\t%s\t%s\t%s\t%s\n' "$fixture" "$rc" "$bytes" "$hash" "$result" \
    >> "$probe/$output_prefix-matrix.tsv"
done

if [ "$failed" = yes ]; then exit 1; fi
