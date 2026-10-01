#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
trace="$probe/trace-makeinfo-pause-silence.gdb"
log="$probe/makeinfo-pause-silence-api.log"
xvfb_log="$probe/xvfb-makeinfo-pause-silence.log"
listing="$probe/makeinfo-pause-silence-files.txt"
input_backup="$probe/makeinfo-pause-silence-original-input1.txt"
output_backup="$probe/makeinfo-pause-silence-original-output.wav"
xpid=
pauses=(-1 0 1 119 120 121 249 250 251 1000 65534 65535)
texts=("comma|A, B" "period-lower|a. b" "period-exception|n. e" \
  "period-upper-control|A. B" "ellipsis-lower|a... b" \
  "ellipsis-upper-positive|A... A" "ellipsis-upper-negative|A... B" \
  "ellipsis-mixed-control|A... a" "plain-control|A B")

if [ "$PWD" != "$work" ]; then
  echo "run this probe with Compose working directory $work" >&2
  exit 2
fi
for output in "$trace" "$log" "$xvfb_log" "$listing" "$input_backup" "$output_backup"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
for item in "${texts[@]}"; do
  name=${item%%|*}
  for pause in "${pauses[@]}"; do
    stem="mi-pause-silence-$name-$pause"
    if [ -e "$probe/$stem.asc.dtt" ] || [ -e "$probe/$stem.bin.dtt" ]; then
      echo "refusing to overwrite capture for $stem" >&2
      exit 3
    fi
  done
done

cp -p "$work/input1.txt" "$input_backup"
cp -p "$work/output.wav" "$output_backup"
restore() {
  cp -p "$input_backup" "$work/input1.txt"
  cp -p "$output_backup" "$work/output.wav"
  cmp -s "$input_backup" "$work/input1.txt"
  cmp -s "$output_backup" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT
trap 'exit 129' HUP
trap 'exit 130' INT
trap 'exit 143' TERM

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
  for item in "${texts[@]}"; do
    name=${item%%|*}
    text=${item#*|}
    for pause in "${pauses[@]}"; do
      varname="${name//-/_}_p${pause#-}"
      if [[ "$pause" == -* ]]; then varname="${name//-/_}_neg${pause#-}"; fi
      path="Z:/work/stage21/mi-pause-silence-$name-$pause"
      printf '  set $text_%s = ((char *(*)(unsigned int))0x1001d9c0)(128)\n' "$varname"
      for ((offset = 0; offset < ${#text}; offset++)); do
        character=${text:offset:1}
        code=$(printf '%d' "'$character")
        printf '  set {char}($text_%s + %d) = %d\n' "$varname" "$offset" "$code"
      done
      printf '  set {char}($text_%s + %d) = 0\n' "$varname" "${#text}"
      printf '  set $path_%s = ((char *(*)(unsigned int))0x1001d9c0)(128)\n' "$varname"
      for ((offset = 0; offset < ${#path}; offset++)); do
        character=${path:offset:1}
        code=$(printf '%d' "'$character")
        printf '  set {char}($path_%s + %d) = %d\n' "$varname" "$offset" "$code"
      done
      printf '  set {char}($path_%s + %d) = 0\n' "$varname" "${#path}"
      printf '  set $raw_%s = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de90)($text_%s, $path_%s, 1, -1, -1, -1, %s, -1, -1)\n' \
        "$varname" "$varname" "$varname" "$pause"
      printf '  printf "MAKEINFO_PAUSE_SILENCE case=%s pause=%s raw_eax=%%#x text=%%s path=%%s\\n", $raw_%s, $text_%s, $path_%s\n' \
        "$name" "$pause" "$varname" "$varname" "$varname"
    done
  done
  printf '%s\n' '  kill' '  quit' 'end' '' 'continue'
} > "$trace"

cp "$probe/input.txt" "$work/input1.txt"
: > "$work/output.wav"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 300s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" \
  > "$log" 2>&1

find "$probe" -maxdepth 1 -type f -name 'mi-pause-silence-*.*.dtt' \
  -printf '%f %s bytes\n' | sort > "$listing"
expected=$((${#texts[@]} * ${#pauses[@]} * 2))
actual=$(wc -l < "$listing")
if [ "$actual" -ne "$expected" ]; then
  echo "expected $expected paired files, found $actual" >&2
  exit 4
fi
