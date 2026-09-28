#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
trace="$probe/trace-userdict-empty-source-effect.gdb"
log="$probe/userdict-empty-source-effect-api.log"
xvfb_log="$probe/xvfb-userdict-empty-source-effect.log"
backup_input="$probe/userdict-empty-source-effect-original-input1.txt"
backup_output="$probe/userdict-empty-source-effect-original-output.wav"
xpid=

restore() {
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
  if [ -e "$backup_input" ]; then cp "$backup_input" "$work/input1.txt"; fi
  if [ -e "$backup_output" ]; then cp "$backup_output" "$work/output.wav"; fi
  if [ -e "$backup_input" ]; then cmp -s "$backup_input" "$work/input1.txt"; fi
  if [ -e "$backup_output" ]; then cmp -s "$backup_output" "$work/output.wav"; fi
}
trap restore EXIT

for output in "$trace" "$log" "$xvfb_log" "$backup_input" "$backup_output" \
  "$probe/userdict-empty-source-effect-hello-control.wav" \
  "$probe/userdict-empty-source-effect-hello-dict.wav" \
  "$probe/userdict-empty-source-effect-period-control.wav" \
  "$probe/userdict-empty-source-effect-period-dict.wav" \
  "$probe/userdict-empty-source-effect-hashes.txt"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
if [ ! -s "$probe/userdict-validation-empty-source.csv" ]; then
  echo "missing nonempty fixture $probe/userdict-validation-empty-source.csv" >&2
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
    'commands' \
    '  silent' \
    '  disable 1'
  write_gdb_string path 'Z:/work/stage16/userdict-validation-empty-source.csv'
  write_gdb_string text_hello 'hello'
  write_gdb_string text_period '.'
  write_gdb_string out_hello_control 'Z:/work/stage16/userdict-empty-source-effect-hello-control.wav'
  write_gdb_string out_hello_dict 'Z:/work/stage16/userdict-empty-source-effect-hello-dict.wav'
  write_gdb_string out_period_control 'Z:/work/stage16/userdict-empty-source-effect-period-control.wav'
  write_gdb_string out_period_dict 'Z:/work/stage16/userdict-empty-source-effect-period-dict.wav'
  cat <<'GDB'
  set $gate_before = *(unsigned char *)0x100a7489
  set {unsigned char}0x100a7489 = 1
  printf "USERDICT_EMPTY_SOURCE gate_before=%u gate_forced=%u\n", $gate_before, *(unsigned char *)0x100a7489
  set $control_hello = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text_hello, $out_hello_control, 1, -1, -1, -1, -1, 200, 0)
  printf "USERDICT_EMPTY_SOURCE case=hello-control synth_ax=%d\n", $control_hello
  set $load_hello = ((short (*)(int, char *))0x10027960)(200, $path)
  printf "USERDICT_EMPTY_SOURCE case=hello load_ax=%d\n", $load_hello
  if $load_hello == 1
    set $dict_hello = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text_hello, $out_hello_dict, 1, -1, -1, -1, -1, 200, 0)
    printf "USERDICT_EMPTY_SOURCE case=hello-dict synth_ax=%d\n", $dict_hello
    set $unload_hello = ((short (*)(int))0x10027a80)(200)
    printf "USERDICT_EMPTY_SOURCE case=hello unload_ax=%d\n", $unload_hello
  end
  set $control_period = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text_period, $out_period_control, 1, -1, -1, -1, -1, 0, 0)
  printf "USERDICT_EMPTY_SOURCE case=period-control synth_ax=%d\n", $control_period
  set $load_period = ((short (*)(int, char *))0x10027960)(200, $path)
  printf "USERDICT_EMPTY_SOURCE case=period load_ax=%d\n", $load_period
  if $load_period == 1
    set $dict_period = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text_period, $out_period_dict, 1, -1, -1, -1, -1, 200, 0)
    printf "USERDICT_EMPTY_SOURCE case=period-dict synth_ax=%d\n", $dict_period
    set $unload_period = ((short (*)(int))0x10027a80)(200)
    printf "USERDICT_EMPTY_SOURCE case=period unload_ax=%d\n", $unload_period
  end
  set {unsigned char}0x100a7489 = $gate_before
  printf "USERDICT_EMPTY_SOURCE gate_restored=%u\n", *(unsigned char *)0x100a7489
  continue
end

continue
GDB
} > "$trace"

cp "$work/input1.txt" "$backup_input"
cp "$work/output.wav" "$backup_output"
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 300s \
  winedbg --gdb /samples/voicetext_paul.exe < "$trace" > "$log" 2>&1

sha256sum "$probe/userdict-empty-source-effect-hello-control.wav" \
  "$probe/userdict-empty-source-effect-hello-dict.wav" \
  "$probe/userdict-empty-source-effect-period-control.wav" \
  "$probe/userdict-empty-source-effect-period-dict.wav" \
  > "$probe/userdict-empty-source-effect-hashes.txt"
