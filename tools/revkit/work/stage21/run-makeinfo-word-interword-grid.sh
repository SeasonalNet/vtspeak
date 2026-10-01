#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
trace="$probe/trace-makeinfo-word-interword-grid.gdb"
log="$probe/makeinfo-word-interword-grid-api.log"
xvfb_log="$probe/xvfb-word-interword-grid.log"
listing="$probe/makeinfo-word-interword-grid-files.txt"
input_backup="$probe/word-interword-grid-original-input1.txt"
output_backup="$probe/word-interword-grid-original-output.wav"
xpid=
words=(a i hello world hi kate paul good morning weather today voice)
masks=(UU UL LU LL)
marks=(comma period ellipsis)

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
for mark in "${marks[@]}"; do
  for mask in "${masks[@]}"; do
    for left in "${words[@]}"; do
      for right in "${words[@]}"; do
        stem="mi-word-interword-$mark-$mask-$left-$right"
        if [ -e "$probe/$stem.asc.dtt" ] || [ -e "$probe/$stem.bin.dtt" ]; then
          echo "refusing to overwrite capture for $stem" >&2
          exit 3
        fi
      done
    done
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
  printf '%s\n' 'set pagination off' 'set confirm off' \
    'set debuginfod enabled off' 'handle SIGSEGV nostop noprint pass' '' \
    'break *0x1001da50' 'commands' '  silent' '  disable 1'
  for mark in "${marks[@]}"; do
    case "$mark" in
      comma) punctuation=, ;;
      period) punctuation=. ;;
      ellipsis) punctuation=... ;;
    esac
    for mask in "${masks[@]}"; do
      for left_base in "${words[@]}"; do
        for right_base in "${words[@]}"; do
          left=$left_base
          right=$right_base
          if [ "${mask:0:1}" = U ]; then left=${left^^}; fi
          if [ "${mask:1:1}" = U ]; then right=${right^^}; fi
          text="$left$punctuation $right"
          varname="${mark}_${mask}_${left_base}_${right_base}"
          path="Z:/work/stage21/mi-word-interword-$mark-$mask-$left_base-$right_base"
          printf '  set $text_%s = ((char *(*)(unsigned int))0x1001d9c0)(256)\n' "$varname"
          for ((offset = 0; offset < ${#text}; offset++)); do
            character=${text:offset:1}
            code=$(printf '%d' "'$character")
            printf '  set {char}($text_%s + %d) = %d\n' "$varname" "$offset" "$code"
          done
          printf '  set {char}($text_%s + %d) = 0\n' "$varname" "${#text}"
          printf '  set $path_%s = ((char *(*)(unsigned int))0x1001d9c0)(256)\n' "$varname"
          for ((offset = 0; offset < ${#path}; offset++)); do
            character=${path:offset:1}
            code=$(printf '%d' "'$character")
            printf '  set {char}($path_%s + %d) = %d\n' "$varname" "$offset" "$code"
          done
          printf '  set {char}($path_%s + %d) = 0\n' "$varname" "${#path}"
          printf '  set $raw_%s = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de90)($text_%s, $path_%s, 1, -1, -1, -1, -1, -1, -1)\n' \
            "$varname" "$varname" "$varname"
          printf '  printf "MAKEINFO_WORD_INTERWORD mark=%s mask=%s left=%s right=%s raw_eax=%%#x\\n", $raw_%s\n' \
            "$mark" "$mask" "$left_base" "$right_base" "$varname"
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
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 1800s \
  winedbg --gdb /samples/voicetext_paul.exe < "$trace" > "$log" 2>&1

find "$probe" -maxdepth 1 -type f -name 'mi-word-interword-*.*.dtt' \
  -printf '%f %s bytes\n' | sort > "$listing"
expected=$((${#marks[@]} * ${#masks[@]} * ${#words[@]} * ${#words[@]} * 2))
actual=$(wc -l < "$listing")
if [ "$actual" -ne "$expected" ]; then
  echo "expected $expected paired files, found $actual" >&2
  exit 4
fi
