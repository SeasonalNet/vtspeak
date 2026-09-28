#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage16
row_log="$probe/syncinfo-copy-row-guard-api.log"
row_xvfb_log="$probe/xvfb-syncinfo-copy-row-guard.log"
nested_log="$probe/syncinfo-copy-nested-guard-api.log"
backup_input=/tmp/vtspeak-syncinfo-copy-undersized-input1.txt
backup_output=/tmp/vtspeak-syncinfo-copy-undersized-output.wav
xpid=
restore() {
  cp "$backup_input" "$work/input1.txt"
  cp "$backup_output" "$work/output.wav"
  cmp -s "$backup_input" "$work/input1.txt"
  cmp -s "$backup_output" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
for output in "$row_log" "$row_xvfb_log" "$nested_log"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
cp "$work/input1.txt" "$backup_input"
cp "$work/output.wav" "$backup_output"
trap restore EXIT
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$row_xvfb_log" 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99
set +e
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-syncinfo-copy-row-guard.gdb" > "$row_log" 2>&1
row_rc=$?
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-syncinfo-copy-nested-guard.gdb" > "$nested_log" 2>&1
nested_rc=$?
set -e
printf 'row_case_exit=%d nested_case_exit=%d\n' "$row_rc" "$nested_rc"
