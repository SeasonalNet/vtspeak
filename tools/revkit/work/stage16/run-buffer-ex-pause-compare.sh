#!/bin/bash
set -euo pipefail

if [ "$#" -ne 2 ] || ! [[ "$1" =~ ^[0-2]$ ]]; then
  echo "usage: $0 SELECTOR [control|pause]" >&2
  exit 2
fi

selector=$1
case "$2" in
  control) input=input-ex-pause-control.txt ;;
  pause0) input=input-ex-pause-tag-zero.txt ;;
  pause) input=input-ex-pause-tag.txt ;;
  pause1000) input=input-ex-pause-tag-1000.txt ;;
  *) echo "unknown case: $2" >&2; exit 2 ;;
esac
case_name=$2
work=/work/stage5
probe=/work/stage16
log="$probe/buffer-ex-pause-$case_name-$selector.log"
xvfb_log="$probe/xvfb-buffer-ex-pause-$case_name-$selector.log"
backup_input="$probe/buffer-ex-pause-$case_name-$selector-original-input1.txt"
backup_output="$probe/buffer-ex-pause-$case_name-$selector-original-output.wav"
trace="/tmp/trace-buffer-ex-pause-$case_name-$selector.gdb"
xpid=
restore() {
  cp "$backup_input" "$work/input1.txt"
  cp "$backup_output" "$work/output.wav"
  cmp -s "$backup_input" "$work/input1.txt"
  cmp -s "$backup_output" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

for output in "$log" "$xvfb_log" "$backup_input" "$backup_output"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
cp "$work/input1.txt" "$backup_input"
cp "$work/output.wav" "$backup_output"
trap restore EXIT
cp "$probe/$input" "$work/input1.txt"
sed -e "s/@SELECTOR@/$selector/g" -e "s/@CASE@/$case_name/g" \
  "$probe/trace-buffer-ex-pause-compare.gdb.in" > "$trace"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe < "$trace" > "$log" 2>&1
