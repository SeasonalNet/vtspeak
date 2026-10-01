#!/bin/bash
set -euo pipefail

probe=/work/stage21
work=/work/stage5
trace="$probe/trace-james-natural-userdict-inuse-unload-v3.gdb"
log="$probe/james-natural-userdict-inuse-unload-v3-api.log"
xvfb_log="$probe/xvfb-james-natural-userdict-inuse-unload-v3.log"
xpid=
outputs=(
  "$probe/james-userdict-inuse-control-v3.wav"
  "$probe/james-userdict-inuse-active-v3.wav"
  "$probe/james-userdict-inuse-restored-v3.wav"
)
paths=(
  "Z:/work/"
  "Z:/work/data-common/verify/verification.txt"
  "Z:/work/stage21/james-userdict-p.csv"
  "Z:/work/stage21/james-userdict-inuse-control-v3.wav"
  "Z:/work/stage21/james-userdict-inuse-active-v3.wav"
  "Z:/work/stage21/james-userdict-inuse-restored-v3.wav"
  "hello"
)
names=(dbroot license dict_csv control active restored text)

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
if [ ! -s "$probe/james-userdict-p.csv" ]; then
  echo "required dictionary fixture is missing: $probe/james-userdict-p.csv" >&2
  exit 4
fi

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
    'break *0x1001da50' 'commands 1' '  silent' '  disable 1'
  for i in "${!paths[@]}"; do
    emit_string "${names[$i]}" "${paths[$i]}"
  done
  cat <<'GDB'
  set $load_james = ((short (*)(void *, int, char *, unsigned int, unsigned int, char *, char *, unsigned int))0x10027af0)(0, 4, $dbroot, 0, 0xffffffff, $license, 0, 0xffffffff)
  set $slot_state = *(unsigned int *)(0x100a0464 + 16)
  set $gate = *(unsigned char *)(0x100a7488 + 4)
  set $capacity = *(int *)($slot_state + 0x4d14)
  printf "NATURAL_INUSE load_james_ax=%d gate=%d capacity=%d\n", $load_james, $gate, $capacity
  set $control_result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $control, 4, -1, -1, -1, -1, -1, 0)
  printf "NATURAL_INUSE case=control synth=%d\n", $control_result
  set $load_dict = ((short (*)(int, char *))0x10027960)(0, $dict_csv)
  set $dict_pointer = *(unsigned int *)0x100a647c
  printf "NATURAL_INUSE dictionary_load=%d ptr=%p\n", $load_dict, $dict_pointer
  break *0x100260c2
  commands 2
    silent
    set $active_context = *(unsigned int *)0x100a447c
    set $context_dict = *(unsigned int *)($active_context + 0x1312c0)
    printf "NATURAL_INUSE watch context=%p context_dictionary=%p dictionary=%p\n", $active_context, $context_dict, $dict_pointer
    set $active_unload = ((short (*)(int))0x10027a80)(0)
    printf "NATURAL_INUSE active_unload=%d\n", $active_unload
    disable 2
    continue
  end
  set $active_result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $active, 4, -1, -1, -1, -1, 0, 0)
  printf "NATURAL_INUSE case=active synth=%d\n", $active_result
  set $idle_unload = ((short (*)(int))0x10027a80)(0)
  printf "NATURAL_INUSE idle_unload=%d\n", $idle_unload
  set $restored_result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $restored, 4, -1, -1, -1, -1, -1, 0)
  printf "NATURAL_INUSE case=restored synth=%d\n", $restored_result
  call ((void (*)(int))0x10027ea0)(4)
  kill
  quit
end

continue
GDB
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
