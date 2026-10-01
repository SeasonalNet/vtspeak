#!/bin/bash
set -euo pipefail

probe=/work/stage21
work=/work/stage5
trace="$probe/trace-james-licensed-userdict-effect.gdb"
log="$probe/james-licensed-userdict-effect-api.log"
xvfb_log="$probe/xvfb-james-licensed-userdict-effect.log"
xpid=
paths=(
  "Z:/work/"
  "Z:/work/data-common/verify/verification.txt"
  "Z:/work/stage21/james-userdict-p.csv"
  "Z:/work/stage21/james-userdict-a.csv"
  "Z:/work/stage21/james-userdict-control.wav"
  "Z:/work/stage21/james-userdict-p.wav"
  "Z:/work/stage21/james-userdict-a.wav"
  "Z:/work/stage21/james-userdict-restored.wav"
  "hello"
)
names=(dbroot license p_csv a_csv control p_output a_output restored text)
outputs=(
  "$probe/james-userdict-control.wav"
  "$probe/james-userdict-p.wav"
  "$probe/james-userdict-a.wav"
  "$probe/james-userdict-restored.wav"
)

if [ "$PWD" != "$work" ]; then
  echo "run with Compose working directory $work" >&2
  exit 2
fi
for output in "$trace" "$log" "$xvfb_log" "${outputs[@]}"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
for fixture in "$probe/james-userdict-p.csv" "$probe/james-userdict-a.csv"; do
  if [ ! -f "$fixture" ]; then
    echo "required fixture is missing: $fixture" >&2
    exit 4
  fi
done

emit_string() {
  local variable=$1 value=$2 offset character code
  printf '  set $%s = ((char *(*)(unsigned int))0x1001d9c0)(%d)\n' \
    "$variable" "$(( ${#value} + 1 ))"
  for ((offset = 0; offset < ${#value}; offset++)); do
    character=${value:offset:1}
    printf -v code '%d' "'$character"
    printf '  set {char}($%s + %d) = %d\n' "$variable" "$offset" "$code"
  done
  printf '  set {char}($%s + %d) = 0\n' "$variable" "${#value}"
}

{
  printf '%s\n' 'set pagination off' 'set confirm off' \
    'set debuginfod enabled off' 'handle SIGSEGV nostop noprint pass' '' \
    'break *0x1001da50' 'commands' '  silent' \
    '  set $sample_speaker = *(int *)($esp + 16)' '  disable 1'
  for i in "${!paths[@]}"; do
    emit_string "${names[$i]}" "${paths[$i]}"
  done
  printf '%s\n' \
    '  set $load_result = ((short (*)(void *, int, char *, unsigned int, unsigned int, char *, char *, unsigned int))0x10027af0)(0, 4, $dbroot, 0, 0xffffffff, $license, 0, 0xffffffff)' \
    '  set $slot_state = *(unsigned int *)(0x100a0464 + 16)' \
    '  set $gate = *(unsigned char *)(0x100a7488 + 4)' \
    '  set $capacity = *(int *)($slot_state + 0x4d14)' \
    '  printf "JAMES_DICT model_load_ax=%d sample_speaker=%d gate=%d capacity=%d\n", $load_result, $sample_speaker, $gate, $capacity' \
    '  set $baseline = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $control, 4, -1, -1, -1, -1, 0, 0)' \
    '  printf "JAMES_DICT case=control synth=%d\n", $baseline' \
    '  set $load_p = ((short (*)(int, char *))0x10027960)(0, $p_csv)' \
    '  set $synth_p = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $p_output, 4, -1, -1, -1, -1, 0, 0)' \
    '  set $unload_p = ((short (*)(int))0x10027a80)(0)' \
    '  printf "JAMES_DICT case=p load=%d synth=%d unload=%d\n", $load_p, $synth_p, $unload_p' \
    '  set $load_a = ((short (*)(int, char *))0x10027960)(0, $a_csv)' \
    '  set $synth_a = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $a_output, 4, -1, -1, -1, -1, 0, 0)' \
    '  set $unload_a = ((short (*)(int))0x10027a80)(0)' \
    '  printf "JAMES_DICT case=a load=%d synth=%d unload=%d\n", $load_a, $synth_a, $unload_a' \
    '  set $restored_result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $restored, 4, -1, -1, -1, -1, 0, 0)' \
    '  printf "JAMES_DICT case=restored synth=%d\n", $restored_result' \
    '  call ((void (*)(int))0x10027ea0)(4)' \
    '  kill' '  quit' 'end' '' 'continue'
} > "$trace"

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
trap 'if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi' EXIT
trap 'exit 129' HUP
trap 'exit 130' INT
trap 'exit 143' TERM
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 600s \
  winedbg --gdb /samples/voicetext_paul.exe < "$trace" > "$log" 2>&1
