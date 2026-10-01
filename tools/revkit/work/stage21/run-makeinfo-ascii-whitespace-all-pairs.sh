#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
start_case=${START_CASE:-0}
case_limit=${CASE_LIMIT:-62208}
suffix=
trace=
log=
xvfb_log=
listing=
input_backup="$probe/ascii-whitespace-all-pairs-original-input1.txt"
output_backup="$probe/ascii-whitespace-all-pairs-original-output.wav"
xpid=
words=(a i hello world hi kate paul good morning weather today voice)
masks=(UU UL LU LL)
marks=(comma period ellipsis)
seps=("" " " "  " $'\t' $'\n' $'\r\n')
sep_names=(none space double tab lf crlf)

if [ "$PWD" != "$work" ]; then
  echo "run this probe with Compose working directory $work" >&2
  exit 2
fi
if ! [[ "$start_case" =~ ^[0-9]+$ ]] || [ "$start_case" -ge 62208 ]; then
  echo "START_CASE must be an integer from 0 through 62207" >&2
  exit 2
fi
if ! [[ "$case_limit" =~ ^[0-9]+$ ]] || [ "$case_limit" -lt 1 ]; then
  echo "CASE_LIMIT must be a positive integer" >&2
  exit 2
fi
end_case=$((start_case + case_limit))
if [ "$end_case" -gt 62208 ]; then end_case=62208; fi
suffix="-$start_case-$end_case"
trace="$probe/trace-makeinfo-ascii-whitespace-all-pairs$suffix.gdb"
log="$probe/makeinfo-ascii-whitespace-all-pairs-api$suffix.log"
xvfb_log="$probe/xvfb-ascii-whitespace-all-pairs$suffix.log"
listing="$probe/makeinfo-ascii-whitespace-all-pairs-files$suffix.txt"
outputs=("$trace" "$log" "$xvfb_log" "$listing")
if [ "$start_case" -eq 0 ]; then outputs+=("$input_backup" "$output_backup"); fi
for output in "${outputs[@]}"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
if [ "$start_case" -eq 0 ] && compgen -G "$probe/mi-wsfull12-*.dtt" > /dev/null; then
  echo "refusing to overwrite existing all-pairs whitespace captures" >&2
  exit 3
fi

if [ "$start_case" -eq 0 ]; then
  cp -p "$work/input1.txt" "$input_backup"
  cp -p "$work/output.wav" "$output_backup"
else
  expected_existing=$((start_case * 2))
  existing=$(find "$probe" -maxdepth 1 -type f -name 'mi-wsfull12-*.dtt' | wc -l)
  if [ "$existing" -ne "$expected_existing" ]; then
    echo "resume expected $expected_existing existing captures, found $existing" >&2
    exit 3
  fi
  cmp -s "$input_backup" "$work/input1.txt"
  cmp -s "$output_backup" "$work/output.wav"
fi
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

case_index=0
{
  printf '%s\n' 'set pagination off' 'set confirm off' \
    'set debuginfod enabled off' 'handle SIGSEGV nostop noprint pass' '' \
    'break *0x1001da50' 'commands' '  silent' '  disable 1'
  for ((left_index = 0; left_index < ${#words[@]}; left_index++)); do
    left_base=${words[left_index]}
    for ((right_index = 0; right_index < ${#words[@]}; right_index++)); do
      right_base=${words[right_index]}
      pair=$((left_index * ${#words[@]} + right_index))
      for mark in "${marks[@]}"; do
        case "$mark" in
          comma) punctuation=, ;;
          period) punctuation=. ;;
          ellipsis) punctuation=... ;;
        esac
        for mask in "${masks[@]}"; do
          left=$left_base
          right=$right_base
          if [ "${mask:0:1}" = U ]; then left=${left^^}; fi
          if [ "${mask:1:1}" = U ]; then right=${right^^}; fi
          for ((pre_index = 0; pre_index < ${#seps[@]}; pre_index++)); do
            for ((post_index = 0; post_index < ${#seps[@]}; post_index++)); do
              if [ "$case_index" -lt "$start_case" ] || [ "$case_index" -ge "$end_case" ]; then
                case_index=$((case_index + 1))
                continue
              fi
              pre_name=${sep_names[pre_index]}
              post_name=${sep_names[post_index]}
              text="$left${seps[pre_index]}$punctuation${seps[post_index]}$right"
              varname=reused
              path="Z:/work/stage21/mi-wsfull12-$case_index"
              printf '  set $text_%s = ((char *(*)(unsigned int))0x1001d9c0)(256)\n' "$varname"
              for ((offset = 0; offset < ${#text}; offset++)); do
                character=${text:offset:1}
                printf -v code '%d' "'$character"
                printf '  set {char}($text_%s + %d) = %d\n' "$varname" "$offset" "$code"
              done
              printf '  set {char}($text_%s + %d) = 0\n' "$varname" "${#text}"
              printf '  set $path_%s = ((char *(*)(unsigned int))0x1001d9c0)(256)\n' "$varname"
              for ((offset = 0; offset < ${#path}; offset++)); do
                character=${path:offset:1}
                printf -v code '%d' "'$character"
                printf '  set {char}($path_%s + %d) = %d\n' "$varname" "$offset" "$code"
              done
              printf '  set {char}($path_%s + %d) = 0\n' "$varname" "${#path}"
              printf '  set $raw_%s = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de90)($text_%s, $path_%s, 1, -1, -1, -1, -1, -1, -1)\n' \
                "$varname" "$varname" "$varname"
              printf '  printf "MAKEINFO_ASCII_WS_FULL pair=%d mark=%s mask=%s pre=%s post=%s raw_eax=%%#x\\n", $raw_%s\n' \
                "$pair" "$mark" "$mask" "$pre_name" "$post_name" "$varname"
              case_index=$((case_index + 1))
            done
          done
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

find "$probe" -maxdepth 1 -type f -name 'mi-wsfull12-*.dtt' \
  -printf '%f %s bytes\n' | sort > "$listing"
expected=$((end_case * 2))
actual=$(wc -l < "$listing")
if [ "$actual" -ne "$expected" ]; then
  echo "expected $expected paired files, found $actual" >&2
  exit 4
fi
