#!/bin/bash
set -euo pipefail

probe=/work/stage21
work=/work/stage21/sandbox/stage5
trace="$probe/trace-lipsync-missing-parent-null-writer.gdb"
log="$probe/lipsync-missing-parent-null-writer-api.log"
xvfb_log="$probe/xvfb-lipsync-missing-parent-null-writer.log"
backup_input=/tmp/vtspeak-lipsync-null-writer-input1.txt
backup_output=/tmp/vtspeak-lipsync-null-writer-output.wav
xpid=

restore() {
  if [ -f "$backup_input" ]; then
    cp "$backup_input" "$work/input1.txt"
    cmp -s "$backup_input" "$work/input1.txt"
    rm -f "$backup_input"
  fi
  if [ -f "$backup_output" ]; then
    cp "$backup_output" "$work/output.wav"
    cmp -s "$backup_output" "$work/output.wav"
    rm -f "$backup_output"
  fi
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT

for output in "$log" "$xvfb_log" "$backup_input" "$backup_output"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
test -f "$work/input1.txt"
test -f "$work/output.wav"
test ! -e "$work/no-such-lipsync-dir"
cp "$work/input1.txt" "$backup_input"
cp "$work/output.wav" "$backup_output"
cp "$probe/input.txt" "$work/input1.txt"
cd "$work"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" > "$log" 2>&1

if ! grep -Fq 'LIPSYNC_NULL_WRITER eip=0x10025e4d edx=0 ' "$log"; then
  echo "did not capture the null report-writer pointer" >&2
  exit 1
fi
