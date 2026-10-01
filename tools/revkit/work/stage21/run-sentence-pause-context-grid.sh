#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
trace="$probe/trace-sentence-pause-context-grid.gdb"
log="$probe/sentence-pause-context-grid-api.log"
xvfb_log="$probe/xvfb-sentence-pause-context-grid.log"
input_backup="$probe/sentence-pause-context-grid-original-input1.txt"
output_backup="$probe/sentence-pause-context-grid-original-output.wav"
xpid=
values=(0 1 199 200 201 250 500 924 925 926 65534 65535)
contexts=("Hello. World." "Hello. world." "hello. world." "Hello world.")
context_names=(upper upper_lower lower no_period)

if [ "$PWD" != "$work" ]; then
  echo "run this probe with Compose working directory $work" >&2
  exit 2
fi
for output in "$trace" "$log" "$xvfb_log" "$input_backup" "$output_backup"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
for value in "${values[@]}"; do
  for context in "${context_names[@]}"; do
    output="$probe/sentence-pause-$value-$context.wav"
    if [ -e "$output" ]; then
      echo "refusing to overwrite $output" >&2
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
  printf '%s\n' 'set pagination off' 'set confirm off' \
    'set debuginfod enabled off' 'handle SIGSEGV nostop noprint pass' '' \
    'break *0x1001da50' 'commands' '  silent' \
    '  set $fmt = *(int *)($esp + 4)' \
    '  set $path = *(char **)($esp + 12)' \
    '  set $speaker = *(int *)($esp + 16)' \
    '  disable 1' '  set $out = (int *)malloc(16)'
  for value in "${values[@]}"; do
    printf '  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, %d, $speaker)\n' "$value"
    printf '  call ((void (*)(int, int))0x100281b0)(200, $speaker)\n'
    printf '  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)\n'
    printf '  printf "SENT_PAUSE value=%d getter_ret=%%d pitch=%%d speed=%%d volume=%%d getter_value=%%d\\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3)\n' "$value"
    for context_index in "${!contexts[@]}"; do
      context=${contexts[$context_index]}
      name=${context_names[$context_index]}
      varname="sentence_${value}_${name}"
      printf '  set $text_%s = ((char *(*)(unsigned int))0x1001d9c0)(64)\n' "$varname"
      for ((offset = 0; offset < ${#context}; offset++)); do
        character=${context:offset:1}
        printf -v code '%d' "'$character"
        printf '  set {char}($text_%s + %d) = %d\n' "$varname" "$offset" "$code"
      done
      printf '  set {char}($text_%s + %d) = 0\n' "$varname" "${#context}"
      printf '  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_%s, $path, $speaker, -1, -1, -1, -1, -1, -1)\n' "$varname"
      printf '  printf "SENT_PAUSE value=%d context=%s synth_ret=%%d\\n", $result\n' "$value" "$name"
      printf '  shell cp %s/output.wav %s/sentence-pause-%d-%s.wav\n' "$work" "$probe" "$value" "$name"
    done
  done
  printf '%s\n' '  kill' '  quit' 'end' '' 'continue'
} > "$trace"

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 1800s \
  winedbg --gdb /samples/voicetext_paul.exe < "$trace" > "$log" 2>&1
