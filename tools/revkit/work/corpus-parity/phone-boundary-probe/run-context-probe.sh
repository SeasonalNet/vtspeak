#!/bin/bash
set -euo pipefail
work=/work/stage5
probe=/work/corpus-parity/phone-boundary-probe
out=/probe
prefix=${TRACE_PREFIX:-phone-boundary-context}
probe_input=${TRACE_INPUT:-$probe/input.txt}
input_backup=$(mktemp "$out/.$prefix-input.XXXXXX")
output_backup=$(mktemp "$out/.$prefix-output.XXXXXX")
cp "$work/input1.txt" "$input_backup"
cp "$work/output.wav" "$output_backup"
input_hash=$(sha256sum "$input_backup" | cut -d ' ' -f 1)
output_hash=$(sha256sum "$output_backup" | cut -d ' ' -f 1)
xpid=
restore() {
  status=$?
  cp "$input_backup" "$work/input1.txt"
  cp "$output_backup" "$work/output.wav"
  cmp -s "$input_backup" "$work/input1.txt"
  input_restored=$?
  cmp -s "$output_backup" "$work/output.wav"
  output_restored=$?
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
  printf 'fixture_input_before=%s\nfixture_input_after=%s\nfixture_output_before=%s\nfixture_output_after=%s\n' \
    "$input_hash" "$(sha256sum "$work/input1.txt" | cut -d ' ' -f 1)" \
    "$output_hash" "$(sha256sum "$work/output.wav" | cut -d ' ' -f 1)" > "$out/$prefix-fixture-restore.txt"
  rm -f "$input_backup" "$output_backup"
  if [ "$input_restored" -ne 0 ] || [ "$output_restored" -ne 0 ]; then status=1; fi
  exit "$status"
}
trap restore EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
cp "$probe_input" "$work/input1.txt"
: > "$work/output.wav"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp >"$out/xvfb-$prefix.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
set +e
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 300s \
  winedbg --gdb /samples/voicetext_kate.exe \
  < /work/stage19/trace-phone-boundary-context.gdb >"$out/$prefix-winedbg.log" 2>&1
rc=$?
set -e
bytes=$(stat -c %s "$work/output.wav")
hash=$(sha256sum "$work/output.wav" | cut -d ' ' -f 1)
if [ "$bytes" -gt 44 ]; then cp "$work/output.wav" "$out/$prefix.wav"; fi
printf 'exit=%s\noutput_bytes=%s\noutput_sha256=%s\n' "$rc" "$bytes" "$hash" > "$out/$prefix-result.txt"
