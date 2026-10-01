#!/bin/bash
set -euo pipefail

probe=/work/stage21
work=/work/stage5
trace="$probe/trace-natural-userdict-inuse-unload-v3.gdb"
log="$probe/natural-userdict-inuse-unload-v3-api.log"
xvfb_log="$probe/xvfb-natural-userdict-inuse-unload-v3.log"
outputs=(
  "$probe/natural-userdict-inuse-control-v3.wav"
  "$probe/natural-userdict-inuse-active-v3.wav"
  "$probe/natural-userdict-inuse-restored-v3.wav"
)
xpid=
paths=(
  "Z:/work/stage21/james-userdict-p.csv"
  "Z:/work/stage21/natural-userdict-inuse-control-v3.wav"
  "Z:/work/stage21/natural-userdict-inuse-active-v3.wav"
  "Z:/work/stage21/natural-userdict-inuse-restored-v3.wav"
  "hello"
)
names=(dict_csv control output restored text)

if [ "$PWD" != "$work" ]; then
  echo "run with Compose working directory $work" >&2
  exit 2
fi
for output_path in "$trace" "$log" "$xvfb_log" "${outputs[@]}"; do
  if [ -e "$output_path" ]; then
    echo "refusing to overwrite $output_path" >&2
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
  set $api_return = *(unsigned int *)$esp
  set $incoming_selector = *(int *)($esp + 4)
  set $gate_before = *(unsigned char *)0x100a7489
  set $control_result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $control, 1, -1, -1, -1, -1, -1, 0)
  printf "NATURAL_INUSE control_synth=%d\n", $control_result
  set $dict_load = ((short (*)(int, char *))0x10027960)(0, $dict_csv)
  set {unsigned char}0x100a7489 = 1
  printf "NATURAL_INUSE setup selector=%d dict_load=%d gate_before=%d gate_forced=%d return=%p\n", $incoming_selector, $dict_load, $gate_before, *(unsigned char *)0x100a7489, $api_return
  set *(int *)($esp + 4) = 4
  set *(char **)($esp + 8) = $text
  set *(char **)($esp + 12) = $output
  set *(int *)($esp + 16) = 1
  set *(int *)($esp + 20) = -1
  set *(int *)($esp + 24) = -1
  set *(int *)($esp + 28) = -1
  set *(int *)($esp + 32) = -1
  set *(int *)($esp + 36) = 0
  set *(int *)($esp + 40) = 0
  break *0x100260c2
  commands 2
    silent
    set $active_context = *(unsigned int *)0x100a147c
    set $active_dictionary = *(unsigned int *)($active_context + 0x1312c0)
    set $dictionary_pointer = *(unsigned int *)0x100a647c
    printf "NATURAL_INUSE context=%p context_dictionary=%p dictionary=%p\n", $active_context, $active_dictionary, $dictionary_pointer
    set $active_unload = ((short (*)(int))0x10027a80)(0)
    printf "NATURAL_INUSE active_unload=%d\n", $active_unload
    disable 2
    continue
  end
  break *$api_return
  commands 3
    silent
    set $synth_status = $eax
    printf "NATURAL_INUSE synthesis_return=%d\n", $synth_status
    set $idle_unload = ((short (*)(int))0x10027a80)(0)
    printf "NATURAL_INUSE idle_unload=%d\n", $idle_unload
    set {unsigned char}0x100a7489 = $gate_before
    printf "NATURAL_INUSE gate_restored=%d\n", *(unsigned char *)0x100a7489
    set $restored_result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $restored, 1, -1, -1, -1, -1, -1, 0)
    printf "NATURAL_INUSE restored_synth=%d\n", $restored_result
    disable 3
    continue
  end
  continue
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
