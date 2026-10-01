#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
trace="$probe/trace-makeinfo-lowercase-matrix.gdb"
log="$probe/makeinfo-lowercase-matrix-api.log"
xvfb_log="$probe/xvfb-makeinfo-lowercase-matrix.log"
listing="$probe/makeinfo-lowercase-matrix-files.txt"
xpid=

if [ "$PWD" != "$work" ]; then
  echo "run this probe with Compose working directory $work" >&2
  exit 2
fi
for output in "$trace" "$log" "$xvfb_log" "$listing"; do
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
    '  disable 1'
  for letter in {a..z}; do
    varname="letter_$letter"
    path="Z:/work/stage21/mi-alphabet-lower-$letter"
    printf '  set $text_%s = ((char *(*)(unsigned int))0x1001d9c0)(128)\n' "$varname"
    printf '  set {char}($text_%s + 0) = %d\n' "$varname" "'$letter"
    printf '  set {char}($text_%s + 1) = 0\n' "$varname"
    printf '  set $path_%s = ((char *(*)(unsigned int))0x1001d9c0)(128)\n' "$varname"
    for ((offset = 0; offset < ${#path}; offset++)); do
      character=${path:offset:1}
      printf '  set {char}($path_%s + %d) = %d\n' "$varname" "$offset" "'$character"
    done
    printf '  set {char}($path_%s + %d) = 0\n' "$varname" "${#path}"
    printf '  set $raw_%s = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de90)($text_%s, $path_%s, 1, -1, -1, -1, -1, -1, -1)\n' \
      "$varname" "$varname" "$varname"
    printf '  printf "MAKEINFO_LOWER letter=%s raw_eax=%%#x text=%%s path=%%s\\n", $raw_%s, $text_%s, $path_%s\n' \
      "$letter" "$varname" "$varname" "$varname"
  done
  printf '%s\n' '  kill' '  quit' 'end' '' 'continue'
} > "$trace"

cp "$probe/input.txt" "$work/input1.txt"
: > "$work/output.wav"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
trap 'if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi' EXIT
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 600s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" \
  > "$log" 2>&1

find "$probe" -maxdepth 1 -type f -name 'mi-alphabet-lower-*.*.dtt' \
  -printf '%f %s bytes\n' | sort > "$listing"
