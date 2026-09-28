#!/bin/bash
set -euo pipefail

work=/work/stage5
evidence=${EVIDENCE_DIR:-/work/corpus-parity/stage20/continuity-bit7-unit-major-overlay}
run_timeout=${RUN_TIMEOUT:-300s}
wine_debug=${WINE_DEBUG_LEVEL:--all}
gdb_script=${GDB_SCRIPT:-/work/stage20/trace-2013-continuity-gate-bit7-unit-major.gdb}
backup=$(mktemp -d /tmp/stage20-marker-provenance-restore.XXXXXX)
xpid=

mkdir -p "$evidence"

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
    > "$evidence/fixture-restore.txt"
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

cp /work/corpus-parity/phone-boundary-probe/pah0-input.txt "$work/input1.txt"
: > "$work/output.wav"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$evidence/xvfb.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
set +e
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG="$wine_debug" timeout "$run_timeout" \
  winedbg --gdb /samples/voicetext_kate.exe \
  < "$gdb_script" \
  > "$evidence/winedbg.log" 2>&1
rc=$?
set -e

bytes=$(stat -c %s "$work/output.wav")
hash=$(sha256sum "$work/output.wav" | cut -d ' ' -f 1)
if [ "$bytes" -gt 44 ]; then cp "$work/output.wav" "$evidence/output.wav"; fi
printf 'exit=%s\noutput_bytes=%s\noutput_sha256=%s\n' "$rc" "$bytes" "$hash" \
  > "$evidence/result.txt"
if [ "$rc" -ne 0 ] || [ "$bytes" -le 44 ]; then exit 1; fi
