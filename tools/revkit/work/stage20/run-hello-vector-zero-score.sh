#!/bin/bash
set -euo pipefail

work=/work/stage5
evidence=/work/corpus-parity/stage20
input_fixture=${INPUT_FIXTURE:-/work/corpus-parity/phone-boundary-probe/pah0-input.txt}
backup=$(mktemp -d /tmp/stage20-adapted-pah0-restore.XXXXXX)
xpid=

cp "$work/input1.txt" "$backup/input1.txt"
cp "$work/output.wav" "$backup/output.wav"
input_before=$(sha256sum "$backup/input1.txt" | cut -d ' ' -f 1)
output_before=$(sha256sum "$backup/output.wav" | cut -d ' ' -f 1)

restore() {
  status=$?
  cp "$backup/input1.txt" "$work/input1.txt"
  cp "$backup/output.wav" "$work/output.wav"
  cmp -s "$backup/input1.txt" "$work/input1.txt" || status=1
  cmp -s "$backup/output.wav" "$work/output.wav" || status=1
  printf 'fixture_input_before=%s\nfixture_input_after=%s\nfixture_output_before=%s\nfixture_output_after=%s\n' \
    "$input_before" "$(sha256sum "$work/input1.txt" | cut -d ' ' -f 1)" \
    "$output_before" "$(sha256sum "$work/output.wav" | cut -d ' ' -f 1)" \
    > "$evidence/hello-vector-zero-score-restore.txt"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
  if [ "$status" -eq 0 ]; then
    rm -f "$backup/input1.txt" "$backup/output.wav"
    rmdir "$backup"
  fi
  exit "$status"
}
trap restore EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

cp "$input_fixture" "$work/input1.txt"
: > "$work/output.wav"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$evidence/hello-vector-zero-score-xvfb.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
set +e
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 300s \
  winedbg --gdb /samples/voicetext_kate.exe \
  < /work/stage20/trace-hello-vector-zero-score.gdb \
  > "$evidence/hello-vector-zero-score-winedbg.log" 2>&1
rc=$?
set -e

bytes=$(stat -c %s "$work/output.wav")
hash=$(sha256sum "$work/output.wav" | cut -d ' ' -f 1)
if [ "$bytes" -gt 44 ]; then cp "$work/output.wav" "$evidence/hello-vector-zero-score.wav"; fi
printf 'exit=%s\noutput_bytes=%s\noutput_sha256=%s\n' "$rc" "$bytes" "$hash" \
  > "$evidence/hello-vector-zero-score-result.txt"
if [ "$rc" -ne 0 ] || [ "$bytes" -le 44 ]; then exit 1; fi
