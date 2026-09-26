#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/probe
input_backup=/tmp/lead5-james-package-input-original
output_backup=/tmp/lead5-james-package-output-original

cp "$work/input1.txt" "$input_backup"
cp "$work/output.wav" "$output_backup"
restore() {
  cp "$input_backup" "$work/input1.txt"
  cp "$output_backup" "$work/output.wav"
}
trap restore EXIT

for fixture in prose numbers address; do
  case "$fixture" in
    prose) input="$probe/input.txt" ;;
    numbers) input="$probe/numbers.txt" ;;
    address) input="$probe/address.txt" ;;
  esac
  cp "$input" "$work/input1.txt"
  : > "$work/output.wav"
  set +e
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=+file timeout 180s \
    wine /samples/voicetext_james.exe >"$probe/james-package-$fixture.log" 2>&1
  rc=$?
  set -e
  bytes=$(stat -c %s "$work/output.wav")
  printf '%s exit=%s output_bytes=%s\n' "$fixture" "$rc" "$bytes"
  if [ "$rc" -ne 0 ] || [ "$bytes" -le 44 ]; then
    exit 1
  fi
  cp "$work/output.wav" "$probe/james-package-$fixture.wav"
done

(cd "$probe" && sha256sum ./james-package-*.wav > james-package-wav-sha256.txt)
