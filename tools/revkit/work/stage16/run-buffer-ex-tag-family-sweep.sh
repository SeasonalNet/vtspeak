#!/bin/bash
set -euo pipefail

if [ "$#" -lt 1 ] || [ "$#" -gt 2 ] || ! [[ "$1" =~ ^[0-2]$ ]]; then
  echo "usage: $0 SELECTOR [vtml|legacy|mark-forms|mark-name-edges]" >&2
  exit 2
fi

selector=$1
tag_family=${2:-vtml}
case "$tag_family" in
  vtml) input=input-ex-records-vtml-tag-family-sweep.txt ;;
  legacy) input=input-ex-records-legacy-tag-family-sweep.txt ;;
  mark-forms) input=input-ex-records-vtml-mark-unnamed.txt ;;
  mark-name-edges) input=input-ex-records-vtml-mark-name-edges.txt ;;
  *) echo "unknown tag family: $tag_family" >&2; exit 2 ;;
esac
work=/work/stage5
probe=/work/stage16
log="$probe/buffer-ex-tag-family-sweep-$tag_family-$selector.log"
xvfb_log="$probe/xvfb-buffer-ex-tag-family-sweep-$tag_family-$selector.log"
backup_input="$probe/buffer-ex-tag-family-sweep-$tag_family-original-input1.txt"
backup_output="$probe/buffer-ex-tag-family-sweep-$tag_family-original-output.wav"
xpid=
restore() {
  cp "$backup_input" "$work/input1.txt"
  cp "$backup_output" "$work/output.wav"
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
sed "s/@SELECTOR@/$selector/g" "$probe/trace-buffer-ex-records.gdb.in" \
  > "/tmp/trace-buffer-ex-tag-family-sweep-$tag_family-$selector.gdb"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "/tmp/trace-buffer-ex-tag-family-sweep-$tag_family-$selector.gdb" > "$log" 2>&1
