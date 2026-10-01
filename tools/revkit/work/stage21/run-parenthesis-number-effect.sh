#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
fixture=${1:-number}
mode_list=${2:-"0 1 2 4"}
xpid=
restore() {
  cp /tmp/parenthesis-number-input1.txt "$work/input1.txt"
  cp /tmp/parenthesis-number-output.wav "$work/output.wav"
  cmp -s /tmp/parenthesis-number-input1.txt "$work/input1.txt"
  cmp -s /tmp/parenthesis-number-output.wav "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
cp "$work/input1.txt" /tmp/parenthesis-number-input1.txt
cp "$work/output.wav" /tmp/parenthesis-number-output.wav
trap restore EXIT
cp "$probe/input-parenthesis-$fixture.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-parenthesis-number.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1

for mode in $mode_list; do
  sed "s/@MODE@/$mode/g" "$probe/trace-parenthesis-number-effect.gdb.in" > "/tmp/parenthesis-number-$mode.gdb"
  : > "$work/output.wav"
  run_rc=0
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
    winedbg --gdb /samples/voicetext_paul.exe \
    < "/tmp/parenthesis-number-$mode.gdb" > "$probe/parenthesis-$fixture-$mode-api.log" 2>&1 || run_rc=$?
  echo "PAREN_RUN mode=$mode shell_status=$run_rc"
  grep -E '^(PAREN_|#)' "$probe/parenthesis-$fixture-$mode-api.log" || true
  cp "$work/output.wav" "$probe/parenthesis-$fixture-$mode.wav"
done

sha256sum "$probe"/parenthesis-$fixture-*.wav
