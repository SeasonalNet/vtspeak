#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
trace="$probe/trace-pause-precedence-grid.gdb"
log="$probe/pause-precedence-grid-api.log"
xvfb_log="$probe/xvfb-pause-precedence-grid.log"
input_backup="$probe/pause-precedence-grid-original-input1.txt"
output_backup="$probe/pause-precedence-grid-original-output.wav"
xpid=
stored_values=(0 925)
call_values=(-1 0 250)
contexts=("Hello. World." "Hello, world." "Hello world.")
context_names=(period comma control)
axes=(sentence comma)

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
for axis in "${axes[@]}"; do
  for stored in "${stored_values[@]}"; do
    for call_pause in "${call_values[@]}"; do
      for context in "${context_names[@]}"; do
        output="$probe/pause-precedence-$axis-$stored-$call_pause-$context.wav"
        if [ -e "$output" ]; then
          echo "refusing to overwrite $output" >&2
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
    'break *0x1001da50' 'commands' '  silent' \
    '  set $fmt = *(int *)($esp + 4)' \
    '  set $path = *(char **)($esp + 12)' \
    '  set $speaker = *(int *)($esp + 16)' \
    '  disable 1' '  set $text = (char *)malloc(64)'
  for axis in "${axes[@]}"; do
    for stored in "${stored_values[@]}"; do
      if [ "$axis" = sentence ]; then
        printf '  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, %d, $speaker)\n' "$stored"
        printf '  call ((void (*)(int, int))0x100281b0)(200, $speaker)\n'
      else
        printf '  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 200, $speaker)\n'
        printf '  call ((void (*)(int, int))0x100281b0)(%d, $speaker)\n' "$stored"
      fi
      for call_pause in "${call_values[@]}"; do
        for context_index in "${!contexts[@]}"; do
          context=${contexts[$context_index]}
          name=${context_names[$context_index]}
          for ((offset = 0; offset < ${#context}; offset++)); do
            character=${context:offset:1}
            printf -v code '%d' "'$character"
            printf '  set {char}($text + %d) = %d\n' "$offset" "$code"
          done
          printf '  set {char}($text + %d) = 0\n' "${#context}"
          printf '  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, %d, -1, -1)\n' "$call_pause"
          printf '  printf "PAUSE_PRECEDENCE axis=%s stored=%d call=%d context=%s ret=%%d\\n", $result\n' "$axis" "$stored" "$call_pause" "$name"
          printf '  shell cp %s/output.wav %s/pause-precedence-%s-%d-%d-%s.wav\n' "$work" "$probe" "$axis" "$stored" "$call_pause" "$name"
        done
      done
    done
  done
  printf '%s\n' '  kill' '  quit' 'end' '' 'continue'
} > "$trace"

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 600s \
  winedbg --gdb /samples/voicetext_paul.exe < "$trace" > "$log" 2>&1
