#!/bin/bash
set -euo pipefail

work=/work/stage17
evidence=/probe/stage20
input_fixture=${INPUT_FIXTURE:-/probe/phone-boundary-probe/pah0-input.txt}
backup=$(mktemp -d /tmp/stage20-pah0-restore.XXXXXX)
xpid=

for name in input.txt numbers.txt address.txt; do
  cp "$work/$name" "$backup/$name"
done
for name in prose numbers address; do
  cp "/probe/runtime-original-msi-2006-$name.wav" "$backup/$name.wav"
done
input_before=$(sha256sum "$backup/input.txt" | cut -d ' ' -f 1)
output_before=$(sha256sum "$backup/prose.wav" | cut -d ' ' -f 1)

restore() {
  status=$?
  for name in input.txt numbers.txt address.txt; do
    cp "$backup/$name" "$work/$name"
  done
  for name in prose numbers address; do
    cp "$backup/$name.wav" "/probe/runtime-original-msi-2006-$name.wav"
  done
  input_ok=0
  output_ok=0
  cmp -s "$backup/input.txt" "$work/input.txt" || input_ok=1
  cmp -s "$backup/numbers.txt" "$work/numbers.txt" || input_ok=1
  cmp -s "$backup/address.txt" "$work/address.txt" || input_ok=1
  cmp -s "$backup/prose.wav" "/probe/runtime-original-msi-2006-prose.wav" || output_ok=1
  cmp -s "$backup/numbers.wav" "/probe/runtime-original-msi-2006-numbers.wav" || output_ok=1
  cmp -s "$backup/address.wav" "/probe/runtime-original-msi-2006-address.wav" || output_ok=1
  printf 'fixture_input_before=%s\nfixture_input_after=%s\nfixture_output_before=%s\nfixture_output_after=%s\n' \
    "$input_before" "$(sha256sum "$work/input.txt" | cut -d ' ' -f 1)" \
    "$output_before" "$(sha256sum "/probe/runtime-original-msi-2006-prose.wav" | cut -d ' ' -f 1)" \
    > "$evidence/hello-legacy-records-restore.txt"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
  if [ "$input_ok" -ne 0 ] || [ "$output_ok" -ne 0 ]; then status=1; fi
  exit "$status"
}
trap restore EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

cp "$input_fixture" "$work/input.txt"
cp "$input_fixture" "$work/numbers.txt"
cp "$input_fixture" "$work/address.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$evidence/hello-legacy-records-xvfb.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
set +e
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /probe/probe-original-kate-2006.exe \
  < /work/stage20/trace-hello-legacy-records.gdb \
  > "$evidence/hello-legacy-records-winedbg.log" 2>&1
rc=$?
set -e

for name in prose numbers address; do
  cp "/probe/runtime-original-msi-2006-$name.wav" "$evidence/old-hello-records-$name.wav"
done
printf 'exit=%s\n' "$rc" > "$evidence/hello-legacy-records-result.txt"
if [ "$rc" -ne 0 ]; then exit "$rc"; fi
