#!/bin/bash
set -euo pipefail

probe=/work/stage21
work=/work/stage21/sandbox/stage5
trace="$probe/trace-userdict-memory-multiline-effect-v2.gdb"
log="$probe/userdict-memory-multiline-effect-v2-api.log"
xvfb_log="$probe/xvfb-userdict-memory-multiline-effect-v2.log"
backup_input="$probe/userdict-memory-multiline-effect-v2-original-input1.txt"
backup_output="$probe/userdict-memory-multiline-effect-v2-original-output.wav"
xpid=

restore() {
  if [ -f "$backup_input" ]; then
    cp "$backup_input" "$work/input1.txt"
    cmp -s "$backup_input" "$work/input1.txt"
  fi
  if [ -f "$backup_output" ]; then
    cp "$backup_output" "$work/output.wav"
    cmp -s "$backup_output" "$work/output.wav"
  fi
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT

outputs=(
  "$probe/userdict-memory-multiline-hello-control-v2.wav"
  "$probe/userdict-memory-multiline-world-control-v2.wav"
  "$probe/userdict-memory-multiline-hello-loaded-v2.wav"
  "$probe/userdict-memory-multiline-world-loaded-v2.wav"
  "$probe/userdict-memory-multiline-hashes-v2.txt"
)
for output in "$trace" "$log" "$xvfb_log" "$backup_input" "$backup_output" "${outputs[@]}"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
test -f "$work/input1.txt"
test -f "$work/output.wav"

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
    'commands' \
    '  silent' \
    '  disable 1'
  write_gdb_string text_hello hello
  write_gdb_string text_world world
  write_gdb_string out_hello_control Z:/work/stage21/userdict-memory-multiline-hello-control-v2.wav
  write_gdb_string out_world_control Z:/work/stage21/userdict-memory-multiline-world-control-v2.wav
  write_gdb_string out_hello_loaded Z:/work/stage21/userdict-memory-multiline-hello-loaded-v2.wav
  write_gdb_string out_world_loaded Z:/work/stage21/userdict-memory-multiline-world-loaded-v2.wav
  cat <<'GDB'
  set $hello_control = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text_hello, $out_hello_control, 1, -1, -1, -1, -1, 200, 0)
  printf "USERDICT_MULTILINE case=hello_control synth_ax=%d\n", $hello_control
  set $world_control = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text_world, $out_world_control, 1, -1, -1, -1, -1, 200, 0)
  printf "USERDICT_MULTILINE case=world_control synth_ax=%d\n", $world_control
  set $gate_before = *(unsigned char *)0x100a7489
  set {unsigned char}0x100a7489 = 1
  set $rows = ((char *(*)(unsigned int))0x1001d9c0)(22)
  set {unsigned char}($rows + 0) = 104
  set {unsigned char}($rows + 1) = 101
  set {unsigned char}($rows + 2) = 108
  set {unsigned char}($rows + 3) = 108
  set {unsigned char}($rows + 4) = 111
  set {unsigned char}($rows + 5) = 44
  set {unsigned char}($rows + 6) = 72
  set {unsigned char}($rows + 7) = 72
  set {unsigned char}($rows + 8) = 44
  set {unsigned char}($rows + 9) = 80
  set {unsigned char}($rows + 10) = 10
  set {unsigned char}($rows + 11) = 119
  set {unsigned char}($rows + 12) = 111
  set {unsigned char}($rows + 13) = 114
  set {unsigned char}($rows + 14) = 108
  set {unsigned char}($rows + 15) = 100
  set {unsigned char}($rows + 16) = 44
  set {unsigned char}($rows + 17) = 72
  set {unsigned char}($rows + 18) = 72
  set {unsigned char}($rows + 19) = 44
  set {unsigned char}($rows + 20) = 80
  set {unsigned char}($rows + 21) = 0
  set $load = ((short (*)(int, char *, char *, int))0x10027880)(200, 0, $rows, 21)
  printf "USERDICT_MULTILINE load_ax=%d gate_before=%u gate_forced=%u span=21\n", $load, $gate_before, *(unsigned char *)0x100a7489
  if $load == 1
    set $hello_loaded = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text_hello, $out_hello_loaded, 1, -1, -1, -1, -1, 200, 0)
    printf "USERDICT_MULTILINE case=hello_loaded synth_ax=%d\n", $hello_loaded
    set $world_loaded = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text_world, $out_world_loaded, 1, -1, -1, -1, -1, 200, 0)
    printf "USERDICT_MULTILINE case=world_loaded synth_ax=%d\n", $world_loaded
    set $unload = ((short (*)(int))0x10027a80)(200)
    printf "USERDICT_MULTILINE unload_ax=%d\n", $unload
  end
  set {unsigned char}0x100a7489 = $gate_before
  printf "USERDICT_MULTILINE gate_restored=%u\n", *(unsigned char *)0x100a7489
  continue
end

continue
GDB
} > "$trace"

cp "$work/input1.txt" "$backup_input"
cp "$work/output.wav" "$backup_output"
cp "$probe/input.txt" "$work/input1.txt"
cd "$work"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 300s \
  winedbg --gdb /samples/voicetext_paul.exe < "$trace" > "$log" 2>&1

sha256sum "${outputs[@]:0:4}" > "${outputs[4]}"
if cmp -s "${outputs[0]}" "${outputs[2]}"; then echo 'hello_effect=identical' >> "${outputs[4]}"; else echo 'hello_effect=different' >> "${outputs[4]}"; fi
if cmp -s "${outputs[1]}" "${outputs[3]}"; then echo 'world_effect=identical' >> "${outputs[4]}"; else echo 'world_effect=different' >> "${outputs[4]}"; fi
