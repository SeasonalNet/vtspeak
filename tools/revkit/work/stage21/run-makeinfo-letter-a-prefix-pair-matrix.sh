#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
trace="$probe/trace-makeinfo-letter-a-prefix-pair-matrix.gdb"
log="$probe/makeinfo-letter-a-prefix-pair-matrix-api.log"
xvfb_log="$probe/xvfb-letter-a-prefix-pair-matrix.log"
listing="$probe/makeinfo-letter-a-prefix-pair-matrix-files.txt"
masks=(UUUU UUUL ULUU ULUL)

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

for third in {a..z}; do
  for fourth in {a..z}; do
    for mask in "${masks[@]}"; do
      for form in plain question; do
        stem="mi-letter-a-prefix-pair-$third$fourth-$mask-$form"
        if [ -e "$probe/$stem.asc.dtt" ] || [ -e "$probe/$stem.bin.dtt" ]; then
          echo "refusing to overwrite capture for $stem" >&2
          exit 3
        fi
      done
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
  for third in {a..z}; do
    for fourth in {a..z}; do
      for mask in "${masks[@]}"; do
        chars=(A a "$third" "$fourth")
        text=
        for index in 0 1 2 3; do
          character=${chars[$index]}
          if [ "${mask:$index:1}" = U ]; then
            character=${character^^}
          else
            character=${character,,}
          fi
          if [ -n "$text" ]; then text+=' '; fi
          text+=$character
        done
        for form in plain question; do
          case_name="$third$fourth-$mask-$form"
          varname=${case_name//-/_}
          call_text=$text
          if [ "$form" = question ]; then call_text+='?'; fi
          path="Z:/work/stage21/mi-letter-a-prefix-pair-$case_name"
          printf '  set $text_%s = ((char *(*)(unsigned int))0x1001d9c0)(128)\n' "$varname"
          for ((offset = 0; offset < ${#call_text}; offset++)); do
            character=${call_text:offset:1}
            code=$(printf '%d' "'$character")
            printf '  set {char}($text_%s + %d) = %d\n' "$varname" "$offset" "$code"
          done
          printf '  set {char}($text_%s + %d) = 0\n' "$varname" "${#call_text}"
          printf '  set $path_%s = ((char *(*)(unsigned int))0x1001d9c0)(128)\n' "$varname"
          for ((offset = 0; offset < ${#path}; offset++)); do
            character=${path:offset:1}
            code=$(printf '%d' "'$character")
            printf '  set {char}($path_%s + %d) = %d\n' "$varname" "$offset" "$code"
          done
          printf '  set {char}($path_%s + %d) = 0\n' "$varname" "${#path}"
          printf '  set $raw_%s = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de90)($text_%s, $path_%s, 1, -1, -1, -1, -1, -1, -1)\n' \
            "$varname" "$varname" "$varname"
          printf '  printf "MAKEINFO_A_PREFIX_PAIR case=%s raw_eax=%%#x text=%%s path=%%s\\n", $raw_%s, $text_%s, $path_%s\n' \
            "$case_name" "$varname" "$varname" "$varname"
        done
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
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 1800s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" \
  > "$log" 2>&1

find "$probe" -maxdepth 1 -type f -name 'mi-letter-a-prefix-pair-*.*.dtt' \
  -printf '%f %s bytes\n' | sort > "$listing"
