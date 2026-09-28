#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
trace="$probe/trace-userdict-inuse-unload-v5.gdb"
log="$probe/userdict-inuse-unload-v5-api.log"
xvfb_log="$probe/xvfb-userdict-inuse-unload-v5.log"
backup_input="$probe/userdict-inuse-unload-v5-original-input1.txt"
backup_output="$probe/userdict-inuse-unload-v5-original-output.wav"
xpid=
backed_up=0

restore() {
  if [ "$backed_up" -eq 1 ]; then
    cp "$backup_input" "$work/input1.txt"
    cp "$backup_output" "$work/output.wav"
    cmp -s "$backup_input" "$work/input1.txt"
    cmp -s "$backup_output" "$work/output.wav"
  fi
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT

for output in "$trace" "$log" "$xvfb_log" "$backup_input" "$backup_output"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
if [ ! -s "$probe/userdict-validation-plain-p.csv" ]; then
  echo "missing nonempty fixture $probe/userdict-validation-plain-p.csv" >&2
  exit 4
fi

write_gdb_string() {
  local name=$1 value=$2 index character
  printf '  set $%s = ((char *(*)(unsigned int))0x1001d9c0)(%d)\n' \
    "$name" "$(( ${#value} + 1 ))"
  for ((index = 0; index < ${#value}; index++)); do
    character=${value:index:1}
    printf '  set {char}($%s + %d) = %d\n' "$name" "$index" "'${character}"
  done
  printf '  set {char}($%s + %d) = 0\n' "$name" "${#value}"
}

{
  printf '%s\n' \
    'set pagination off' \
    'set confirm off' \
    'set debuginfod enabled off' \
    'handle SIGSEGV nostop noprint pass' \
    '' \
    'break *0x1001da50' \
    'commands 1' \
    '  silent' \
    '  disable 1'
  write_gdb_string path 'Z:/work/stage16/userdict-validation-plain-p.csv'
  cat <<'GDB'
  set $load = ((short (*)(int, char *))0x10027960)(200, $path)
  printf "USERDICT_INUSE load index=200 ax=%d\n", $load
  set $dict_pointer = *(unsigned int *)(0x100a647c + 200 * 4)
  set $speaker_state = *(unsigned int *)0x100a0468
  set $capacity = *(int *)($speaker_state + 0x4d14)
  set $reference_slot = (unsigned int *)0x100a147c
  set $reference_before = *$reference_slot
  set $active_context = ((char *(*)(unsigned int))0x1001d9c0)(0x1312e0)
  set {unsigned int}($active_context + 0x1312c0) = $dict_pointer
  printf "USERDICT_INUSE state=0x%x capacity=%d dictionary=0x%x reference_slot=0x%x old_reference=0x%x active_context=0x%x context_dict=0x%x\n", $speaker_state, $capacity, $dict_pointer, $reference_slot, $reference_before, $active_context, *(unsigned int *)($active_context + 0x1312c0)
  set *$reference_slot = $active_context
  set $busy = ((short (*)(int))0x10027980)(200)
  printf "USERDICT_INUSE forced_reference_unload index=200 ax=%d\n", $busy
  set *$reference_slot = $reference_before
  printf "USERDICT_INUSE reference_restored=0x%x\n", *$reference_slot
  call ((void (*)(void *))0x1001da30)($active_context)
  printf "USERDICT_INUSE scratch_context_freed=1\n"
  set $after = ((short (*)(int))0x10027980)(200)
  printf "USERDICT_INUSE idle_unload_after_restore index=200 ax=%d\n", $after
  continue
end

continue
GDB
} > "$trace"

cp "$work/input1.txt" "$backup_input"
cp "$work/output.wav" "$backup_output"
backed_up=1
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 300s \
  winedbg --gdb /samples/voicetext_paul.exe < "$trace" > "$log" 2>&1
