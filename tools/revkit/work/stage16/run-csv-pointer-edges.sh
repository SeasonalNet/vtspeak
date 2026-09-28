#!/bin/bash
set -euo pipefail

if [ "$#" -ne 0 ]; then
  echo "usage: $0" >&2
  exit 2
fi

work=/work/stage5
probe=/work/stage16
xpid=
backup_input=/tmp/vtspeak-csv-pointer-edges-input1.txt
backup_output=/tmp/vtspeak-csv-pointer-edges-output.wav
restore() {
  cp "$backup_input" "$work/input1.txt"
  cp "$backup_output" "$work/output.wav"
  cmp -s "$backup_input" "$work/input1.txt"
  cmp -s "$backup_output" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

cp "$work/input1.txt" "$backup_input"
cp "$work/output.wav" "$backup_output"
trap restore EXIT
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-csv-pointer-edges.log" 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99

for mode in 0 1 2; do
  trace="$probe/trace-csv-pointer-edges-$mode.gdb"
  log="$probe/csv-pointer-edges-$mode-api.log"
  sed "s/@MODE@/$mode/" "$probe/trace-csv-pointer-edges.gdb.in" > "$trace"
  if WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
    winedbg --gdb /samples/voicetext_paul.exe < "$trace" > "$log" 2>&1; then
    rc=0
  else
    rc=$?
  fi
  printf 'DEBUGGER_EXIT mode=%s code=%s\n' "$mode" "$rc" >> "$log"
done
