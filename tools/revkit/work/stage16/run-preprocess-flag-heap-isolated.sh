#!/bin/bash
set -euo pipefail

flag=${1:?usage: run-preprocess-flag-heap-isolated.sh BYTE_FLAG}
if ! [[ "$flag" =~ ^([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$ ]]; then
  echo "flag must be an integer from 1 through 255" >&2
  exit 2
fi

probe=/work/stage16
work=/work/stage5
base="heap-flag${flag}-20260926"
path="Z:/work/stage16/${base}"
trace="$probe/trace-preprocess-heap-flag${flag}.gdb"
log="$probe/preprocess-heap-flag${flag}-api.log"
listing="$probe/preprocess-heap-flag${flag}-files.txt"
xvfb_log="$probe/xvfb-preprocess-heap-flag${flag}.log"
before=$(mktemp)
after=$(mktemp)
xpid=
restore() {
  cp "$probe/preprocess-flags-original-input1.txt" "$work/input1.txt"
  cp "$probe/preprocess-flags-original-output.wav" "$work/output.wav"
  cmp -s "$probe/preprocess-flags-original-input1.txt" "$work/input1.txt"
  cmp -s "$probe/preprocess-flags-original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
  rm -f "$before" "$after"
}
trap restore EXIT

for output in "$trace" "$log" "$listing" "$xvfb_log" "$probe/$base" "$probe/$base.0" "$probe/$base.1" "$probe/$base.2" "$probe/$base.3"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done

{
  printf '%s\n' \
    'set pagination off' \
    'set confirm off' \
    'set debuginfod enabled off' \
    'handle SIGSEGV nostop noprint pass' \
    '' \
    'break *0x1001da50' \
    'commands' \
    '  silent' \
    '  set $text = *(unsigned int *)($esp + 8)' \
    '  disable 1' \
    '  set $path = ((char *(*)(unsigned int))0x1001d9c0)(128)'
  for ((index = 0; index < ${#path}; index++)); do
    character=${path:index:1}
    printf '  set {char}($path + %d) = %d\n' "$index" "'$character"
  done
  printf '  set {char}($path + %d) = 0\n' "${#path}"
  printf '%s\n' \
    '  printf "PREPROCESS_HEAP_PATH_BEFORE value="' \
    '  x/s $path'
  printf '  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, %s, 1, -1, -1, -1, -1, 0, 0)\n' "$flag"
  printf '  printf "PREPROCESS_HEAP_FLAG flag=%s raw_eax=%%#x low_ax=%%d path=%%s\\n", $raw, ((short)$raw), $path\n' "$flag"
  printf '%s\n' \
    '  continue' \
    'end' \
    '' \
    'continue'
} > "$trace"

cp "$probe/input.txt" "$work/input1.txt"
find "$probe" -maxdepth 1 -type f -printf '%f %s %T@\n' | sort > "$before"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" \
  > "$log" 2>&1
find "$probe" -maxdepth 1 -type f -printf '%f %s %T@\n' | sort > "$after"
comm -13 "$before" "$after" > "$listing"
