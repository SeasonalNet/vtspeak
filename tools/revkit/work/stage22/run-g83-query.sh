#!/bin/bash
set -euo pipefail

probe=/work/stage22
work=/work/stage5
xpid=
had_output=0
mode=${1:?usage: run-g83-query.sh control|g83-zero fixture}
fixture=${2:?usage: run-g83-query.sh control|g83-zero fixture}
case "$mode" in
  control|g83-zero) trace="$probe/g83-query-$mode.gdb" ;;
  *) printf 'unknown mode: %s\n' "$mode" >&2; exit 2 ;;
esac

mkdir -p "$probe" "$work"
cd "$work"
cp "$work/input1.txt" /tmp/stage22-original-input1.txt
if [ -e "$work/output.wav" ]; then
  had_output=1
  cp "$work/output.wav" /tmp/stage22-original-output.wav
fi
restore() {
  cp /tmp/stage22-original-input1.txt "$work/input1.txt"
  if [ "$had_output" = 1 ]; then
    cp /tmp/stage22-original-output.wav "$work/output.wav"
  elif [ -e "$work/output.wav" ]; then
    mv "$work/output.wav" "$probe/$fixture-g83-query-$mode-partial.wav"
  fi
  if [ -n "$xpid" ]; then
    kill "$xpid" 2>/dev/null || true
  fi
}
trap restore EXIT

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/g83-query-xvfb.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
cp "$probe/../stage12/inputs/$fixture.txt" "$work/input1.txt"
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" > "$probe/$fixture-g83-query-$mode.log" 2>&1
if [ "$had_output" = 1 ]; then
  cp "$work/output.wav" "$probe/$fixture-g83-query-$mode.wav"
else
  mv "$work/output.wav" "$probe/$fixture-g83-query-$mode.wav"
fi

log="$probe/$fixture-g83-query-$mode.log"
if ! grep -Fq '[Inferior 1 (Remote target) exited normally]' "$log"; then
  printf 'GDB inferior did not exit normally; see %s\n' "$log" >&2
  exit 1
fi
if ! grep -q '^TPP_CODE83 ' "$log"; then
  printf 'No code-83 token write was recorded; see %s\n' "$log" >&2
  exit 1
fi
if [ "$mode" = g83-zero ] && ! grep -q '^TPP_CODE83_ABLATED ' "$log"; then
  printf 'The code-83 intervention was not recorded; see %s\n' "$log" >&2
  exit 1
fi
