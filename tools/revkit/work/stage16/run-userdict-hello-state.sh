#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
trace="$probe/trace-userdict-hello-state.gdb"
log="$probe/userdict-hello-state-api.log"
xvfb_log="$probe/xvfb-userdict-hello-state.log"
backup_input="$probe/userdict-hello-state-original-input1.txt"
backup_output="$probe/userdict-hello-state-original-output.wav"
xpid=

restore() {
  cp "$backup_input" "$work/input1.txt"
  cp "$backup_output" "$work/output.wav"
  cmp -s "$backup_input" "$work/input1.txt"
  cmp -s "$backup_output" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT

for output in "$trace" "$log" "$xvfb_log" "$backup_input" "$backup_output" \
  "$probe/userdict-hello-state-control.wav" "$probe/userdict-hello-state-p.wav" \
  "$probe/userdict-hello-state-a.wav" "$probe/userdict-hello-state-hashes.txt"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
for candidate in userdict-hello-plain.csv userdict-effect-alpha.csv; do
  if [ ! -s "$probe/$candidate" ]; then
    echo "missing nonempty candidate $probe/$candidate" >&2
    exit 4
  fi
done

write_gdb_string() {
  local var_name=$1 value=$2 index character
  printf '  set $%s = ((char *(*)(unsigned int))0x1001d9c0)(%d)\n' \
    "$var_name" "$(( ${#value} + 1 ))"
  for ((index = 0; index < ${#value}; index++)); do
    character=${value:index:1}
    printf '  set {char}($%s + %d) = %d\n' \
      "$var_name" "$index" "'${character}"
  done
  printf '  set {char}($%s + %d) = 0\n' "$var_name" "${#value}"
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
  write_gdb_string path_p 'Z:/work/stage16/userdict-hello-plain.csv'
  write_gdb_string path_a 'Z:/work/stage16/userdict-effect-alpha.csv'
  write_gdb_string text 'hello'
  write_gdb_string path_control 'Z:/work/stage16/userdict-hello-state-control.wav'
  write_gdb_string path_p_output 'Z:/work/stage16/userdict-hello-state-p.wav'
  write_gdb_string path_a_output 'Z:/work/stage16/userdict-hello-state-a.wav'
  cat <<'GDB'
  set $control = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $path_control, 1, -1, -1, -1, -1, 0, 0)
  printf "USERDICT_SYNTH control dictidx=-1 status=%d\n", $control
  set $load_p = ((short (*)(int, char *))0x10027960)(27, $path_p)
  printf "USERDICT_LOAD type=P index=27 status=%d path=%s\n", $load_p, $path_p
  printf "USERDICT_STATE speaker1_gate=%u index27_pointer=%#x\n", *(unsigned char *)0x100a7489, *(unsigned int *)(0x100a647c + 27 * 4)
  if $load_p == 1
    set $synth_p = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $path_p_output, 1, -1, -1, -1, -1, 27, 0)
    printf "USERDICT_SYNTH type=P index=27 status=%d\n", $synth_p
    set $unload_p = ((short (*)(int))0x10027a80)(27)
    printf "USERDICT_UNLOAD type=P index=27 status=%d\n", $unload_p
  end
  set $load_a = ((short (*)(int, char *))0x10027960)(28, $path_a)
  printf "USERDICT_LOAD type=A index=28 status=%d path=%s\n", $load_a, $path_a
  if $load_a == 1
    set $synth_a = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $path_a_output, 1, -1, -1, -1, -1, 28, 0)
    printf "USERDICT_SYNTH type=A index=28 status=%d\n", $synth_a
    set $unload_a = ((short (*)(int))0x10027a80)(28)
    printf "USERDICT_UNLOAD type=A index=28 status=%d\n", $unload_a
  end
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
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 240s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" \
  > "$log" 2>&1

for output in "$probe/userdict-hello-state-control.wav" \
  "$probe/userdict-hello-state-p.wav" "$probe/userdict-hello-state-a.wav"; do
  if [ -e "$output" ]; then
    sha256sum "$output"
  fi
done > "$probe/userdict-hello-state-hashes.txt"
