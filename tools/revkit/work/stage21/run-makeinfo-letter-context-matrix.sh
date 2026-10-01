#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
trace="$probe/trace-makeinfo-letter-context-matrix.gdb"
log="$probe/makeinfo-letter-context-matrix-api.log"
xvfb_log="$probe/xvfb-letter-context-matrix.log"
listing="$probe/makeinfo-letter-context-matrix-files.txt"
xpid=
names=(a a-period a-comma a-question a-bang b b-period b-comma b-question b-bang c c-period c-comma c-question c-bang ab bc abc bcd dotted-spaced dotted-tight comma-pair)
texts=(A 'A.' 'A,' 'A?' 'A!' B 'B.' 'B,' 'B?' 'B!' C 'C.' 'C,' 'C?' 'C!' 'A B' 'B C' 'A B C' 'B C D' 'A. B. C.' 'A.B.C.' 'A, B')

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
  for index in "${!names[@]}"; do
    name=${names[$index]}
    varname=${name//-/_}
    text=${texts[$index]}
    path="Z:/work/stage21/mi-alphabet-context-$name"
    printf '  set $text_%s = ((char *(*)(unsigned int))0x1001d9c0)(128)\n' "$varname"
    for ((offset = 0; offset < ${#text}; offset++)); do
      character=${text:offset:1}
      printf '  set {char}($text_%s + %d) = %d\n' "$varname" "$offset" "'$character"
    done
    printf '  set {char}($text_%s + %d) = 0\n' "$varname" "${#text}"
    printf '  set $path_%s = ((char *(*)(unsigned int))0x1001d9c0)(128)\n' "$varname"
    for ((offset = 0; offset < ${#path}; offset++)); do
      character=${path:offset:1}
      printf '  set {char}($path_%s + %d) = %d\n' "$varname" "$offset" "'$character"
    done
    printf '  set {char}($path_%s + %d) = 0\n' "$varname" "${#path}"
    printf '  set $raw_%s = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de90)($text_%s, $path_%s, 1, -1, -1, -1, -1, -1, -1)\n' \
      "$varname" "$varname" "$varname"
    printf '  printf "MAKEINFO_CONTEXT case=%s raw_eax=%%#x text=%%s path=%%s\\n", $raw_%s, $text_%s, $path_%s\n' \
      "$name" "$varname" "$varname" "$varname"
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

find "$probe" -maxdepth 1 -type f -name 'mi-alphabet-context-*.*.dtt' \
  -printf '%f %s bytes\n' | sort > "$listing"
