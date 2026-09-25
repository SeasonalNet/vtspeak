#!/bin/bash
set -euo pipefail

if [ "$#" -ne 1 ]; then
  echo "usage: $0 null-text|empty-text|null-path" >&2
  exit 2
fi

case "$1" in
  null-text) patch1='set *(unsigned int *)($esp + 8) = 0'; patch2='# no second patch' ;;
  empty-text) patch1='set $text = *(unsigned int *)($esp + 8)'; patch2='set *(unsigned char *)$text = 0' ;;
  null-path) patch1='set *(unsigned int *)($esp + 12) = 0'; patch2='# no second patch' ;;
  *) echo "unknown error probe: $1" >&2; exit 2 ;;
esac

case_name=$1
work=/work/stage5
probe=/work/stage16
xpid=
restore() {
  cp "$probe/original-input1.txt" "$work/input1.txt"
  cp "$probe/original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

cp "$work/input1.txt" "$probe/original-input1.txt"
cp "$work/output.wav" "$probe/original-output.wav"
trap restore EXIT
cp "$probe/input.txt" "$work/input1.txt"
sed -e "s/@PATCH1@/$patch1/" -e "s/@PATCH2@/$patch2/" "$probe/trace-file-error.gdb.in" > "/tmp/trace-file-error-$case_name.gdb"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-error-$case_name.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "/tmp/trace-file-error-$case_name.gdb" > "$probe/file-error-$case_name.log" 2>&1
