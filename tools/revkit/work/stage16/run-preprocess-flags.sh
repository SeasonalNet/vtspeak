#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
xpid=
restore() {
  cp "$probe/preprocess-flags-original-input1.txt" "$work/input1.txt"
  cp "$probe/preprocess-flags-original-output.wav" "$work/output.wav"
  cmp -s "$probe/preprocess-flags-original-input1.txt" "$work/input1.txt"
  cmp -s "$probe/preprocess-flags-original-output.wav" "$work/output.wav"
  rm -f "$work/test.pcm"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

if [ -e "$work/test.pcm" ]; then
  echo "refusing to overwrite preexisting stage5/test.pcm" >&2
  exit 1
fi
for flag in 1 2 3 4 5 6 7 8 9 10 255; do
  if [ -e "$probe/preprocess-followup-flag-$flag.out" ]; then
    echo "refusing to overwrite existing preprocess output for flag $flag" >&2
    exit 1
  fi
done
if [ -e "$probe/preprocess-test.pcm" ]; then
  echo "refusing to overwrite captured preprocess test.pcm" >&2
  exit 1
fi
for flag in 1 2 3 4 5 6 7 8 9 10 255; do
  if [ -e "$probe/preprocess-followup-test-flag-$flag.pcm" ]; then
    echo "refusing to overwrite followup test.pcm for flag $flag" >&2
    exit 1
  fi
done

cp "$work/input1.txt" "$probe/preprocess-flags-original-input1.txt"
cp "$work/output.wav" "$probe/preprocess-flags-original-output.wav"
trap restore EXIT
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-preprocess-flags.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
: > "$probe/preprocess-flags.log"
for flag in 1 2 3 4 5 6 7 8 9 10 255; do
  sed "s/@FLAG@/$flag/g" "$probe/trace-preprocess-flags.gdb" > "$probe/trace-preprocess-one.gdb"
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
    winedbg --gdb /samples/voicetext_paul.exe \
    < "$probe/trace-preprocess-one.gdb" > "$probe/preprocess-flag-$flag.log" 2>&1
  grep 'PREPROCESS flag=' "$probe/preprocess-flag-$flag.log" >> "$probe/preprocess-flags.log"
  candidate="$probe/preprocess-followup-flag-$flag.out"
  if [ -e "$candidate" ]; then
    echo "generated_argument_file flag=$flag bytes=$(wc -c < "$candidate") sha256=$(sha256sum "$candidate" | cut -d ' ' -f 1)" >> "$probe/preprocess-flags.log"
  else
    echo "generated_argument_file flag=$flag absent" >> "$probe/preprocess-flags.log"
  fi
  if [ -e "$work/test.pcm" ]; then
    cp "$work/test.pcm" "$probe/preprocess-followup-test-flag-$flag.pcm"
    echo "generated_test_pcm flag=$flag bytes=$(wc -c < "$work/test.pcm") sha256=$(sha256sum "$work/test.pcm" | cut -d ' ' -f 1) first64=$(od -An -v -tx1 -N64 "$work/test.pcm" | tr -d '\n')" >> "$probe/preprocess-flags.log"
    rm "$work/test.pcm"
  else
    echo "generated_test_pcm flag=$flag absent" >> "$probe/preprocess-flags.log"
  fi
done
