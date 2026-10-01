#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
letter_case=${LETTER_CASE:-upper}
case "$letter_case" in
  upper)
    case_label=upper
    first_text=A
    middle_upper=1
    right_upper=1
    ;;
  lower)
    case_label=lower
    first_text=a
    middle_upper=0
    right_upper=0
    ;;
  upper-lower)
    case_label=mixed-upper-lower
    first_text=A
    middle_upper=0
    right_upper=0
    ;;
  lower-upper)
    case_label=mixed-lower-upper
    first_text=a
    middle_upper=1
    right_upper=1
    ;;
  *)
    echo "LETTER_CASE must be upper, lower, upper-lower, or lower-upper" >&2
    exit 2
    ;;
esac

if [ "$PWD" != "$work" ]; then
  echo "run this probe with Compose working directory $work" >&2
  exit 2
fi

trace="$probe/trace-makeinfo-letter-triple-a-context-$case_label.gdb"
log="$probe/makeinfo-letter-triple-a-context-$case_label-api.log"
xvfb_log="$probe/xvfb-letter-triple-a-context-$case_label.log"
listing="$probe/makeinfo-letter-triple-a-context-$case_label-files.txt"
for output in "$trace" "$log" "$xvfb_log" "$listing"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done

for middle in {a..z}; do
  for right in {a..z}; do
    for form in plain question; do
      stem="mi-letter-triple-a-$case_label-$middle$right-$form"
      if [ -e "$probe/$stem.asc.dtt" ] || [ -e "$probe/$stem.bin.dtt" ]; then
        echo "refusing to overwrite capture for $stem" >&2
        exit 3
      fi
    done
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
  for middle in {a..z}; do
    middle_text=$middle
    if [ "$middle_upper" -eq 1 ]; then
      middle_text=$(printf '%s' "$middle" | tr '[:lower:]' '[:upper:]')
    fi
    for right in {a..z}; do
      right_text=$right
      if [ "$right_upper" -eq 1 ]; then
        right_text=$(printf '%s' "$right" | tr '[:lower:]' '[:upper:]')
      fi
      for form in plain question; do
        name="$middle$right-$form"
        varname=${name//-/_}
        text="$first_text $middle_text $right_text"
        if [ "$form" = question ]; then text+='?'; fi
        path="Z:/work/stage21/mi-letter-triple-a-$case_label-$middle$right-$form"
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
        printf '  set $raw_%s = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de90)($text_%s, $path_%s, 1, -1, -1, -1, -1, -1, -1)\n' \
          "$varname" "$varname" "$varname"
        printf '  printf "MAKEINFO_TRIPLE case=%s raw_eax=%%#x text=%%s path=%%s\\n", $raw_%s, $text_%s, $path_%s\n' \
          "$name" "$varname" "$varname" "$varname"
      done
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

find "$probe" -maxdepth 1 -type f -name "mi-letter-triple-a-$case_label-*.*.dtt" \
  -printf '%f %s bytes\n' | sort > "$listing"
