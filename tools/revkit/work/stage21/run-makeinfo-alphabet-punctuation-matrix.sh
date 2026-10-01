#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
trace="$probe/trace-makeinfo-alphabet-punctuation-matrix.gdb"
log="$probe/makeinfo-alphabet-punctuation-matrix-api.log"
xvfb_log="$probe/xvfb-alphabet-punctuation-matrix.log"
listing="$probe/makeinfo-alphabet-punctuation-matrix-files.txt"
xpid=
letter_case=${LETTER_CASE:-upper}
if [ "$letter_case" = lower ]; then
  capture_prefix=mi-alphabet-punct-lower
  trace="$probe/trace-makeinfo-lowercase-punctuation-matrix.gdb"
  log="$probe/makeinfo-lowercase-punctuation-matrix-api.log"
  xvfb_log="$probe/xvfb-lowercase-punctuation-matrix.log"
  listing="$probe/makeinfo-lowercase-punctuation-matrix-files.txt"
elif [ "$letter_case" = upper ]; then
  capture_prefix=mi-alphabet-punct
else
  echo "LETTER_CASE must be upper or lower" >&2
  exit 2
fi
punctuation_names=(period comma bang question)
punctuation_chars=(. , ! ?)

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
for letter in {a..z}; do
  for punctuation in "${punctuation_names[@]}"; do
    stem="$capture_prefix-$letter-$punctuation"
    if [ -e "$probe/$stem.asc.dtt" ] || [ -e "$probe/$stem.bin.dtt" ]; then
      echo "refusing to overwrite capture for $stem" >&2
      exit 3
    fi
  done
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
    if [ "$letter_case" = upper ]; then upper=$(printf '%s' "$letter" | tr '[:lower:]' '[:upper:]'); else upper=$letter; fi
    for index in "${!punctuation_names[@]}"; do
      name="$letter-${punctuation_names[$index]}"
      varname=${name//-/_}
      text="$upper${punctuation_chars[$index]}"
      path="Z:/work/stage21/$capture_prefix-$name"
      printf '  set $text_%s = ((char *(*)(unsigned int))0x1001d9c0)(128)\n' "$varname"
      for ((offset = 0; offset < ${#text}; offset++)); do
        character=${text:offset:1}
        printf '  set {char}($text_%s + %d) = %d\n' "$varname" "$offset" "'${character}"
      done
      printf '  set {char}($text_%s + %d) = 0\n' "$varname" "${#text}"
      printf '  set $path_%s = ((char *(*)(unsigned int))0x1001d9c0)(128)\n' "$varname"
      for ((offset = 0; offset < ${#path}; offset++)); do
        character=${path:offset:1}
        printf '  set {char}($path_%s + %d) = %d\n' "$varname" "$offset" "'${character}"
      done
      printf '  set {char}($path_%s + %d) = 0\n' "$varname" "${#path}"
      printf '  set $raw_%s = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de90)($text_%s, $path_%s, 1, -1, -1, -1, -1, -1, -1)\n' \
        "$varname" "$varname" "$varname"
      printf '  printf "MAKEINFO_PUNCT case=%s raw_eax=%%#x text=%%s path=%%s\\n", $raw_%s, $text_%s, $path_%s\n' \
        "$name" "$varname" "$varname" "$varname"
    done
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

find "$probe" -maxdepth 1 -type f -name "$capture_prefix-*.*.dtt" \
  -printf '%f %s bytes\n' | sort > "$listing"
