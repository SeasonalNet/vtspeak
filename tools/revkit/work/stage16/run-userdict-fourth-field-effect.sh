#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
trace="$probe/trace-userdict-fourth-field-effect-v2.gdb"
log="$probe/userdict-fourth-field-effect-v2-api.log"
xvfb_log="$probe/xvfb-userdict-fourth-field-effect-v2.log"
backup_input="$probe/userdict-fourth-field-effect-v2-original-input1.txt"
backup_output="$probe/userdict-fourth-field-effect-v2-original-output.wav"
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
  "$probe/userdict-fourth-field-plain-v2.wav" "$probe/userdict-fourth-field-extra-v2.wav" \
  "$probe/userdict-fourth-field-empty-v2.wav" "$probe/userdict-fourth-field-quoted-v2.wav" \
  "$probe/userdict-fourth-field-hashes-v2.txt"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
for fixture in userdict-validation-plain-p.csv userdict-validation-four-p.csv \
  userdict-validation-trailing-empty.csv userdict-validation-quoted-p.csv; do
  if [ ! -s "$probe/$fixture" ]; then
    echo "missing nonempty fixture $probe/$fixture" >&2
    exit 4
  fi
done

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
  write_gdb_string path_plain 'Z:/work/stage16/userdict-validation-plain-p.csv'
  write_gdb_string path_extra 'Z:/work/stage16/userdict-validation-four-p.csv'
  write_gdb_string path_empty 'Z:/work/stage16/userdict-validation-trailing-empty.csv'
  write_gdb_string path_quoted 'Z:/work/stage16/userdict-validation-quoted-p.csv'
  write_gdb_string text 'hello'
  write_gdb_string out_plain 'Z:/work/stage16/userdict-fourth-field-plain-v2.wav'
  write_gdb_string out_extra 'Z:/work/stage16/userdict-fourth-field-extra-v2.wav'
  write_gdb_string out_empty 'Z:/work/stage16/userdict-fourth-field-empty-v2.wav'
  write_gdb_string out_quoted 'Z:/work/stage16/userdict-fourth-field-quoted-v2.wav'
  cat <<'GDB'
  set $gate_before = *(unsigned char *)0x100a7489
  set {unsigned char}0x100a7489 = 1
  printf "USERDICT_FOURTH gate_before=%u gate_forced=%u\n", $gate_before, *(unsigned char *)0x100a7489
  set $load_plain = ((short (*)(int, char *))0x10027960)(200, $path_plain)
  printf "USERDICT_FOURTH case=plain load_ax=%d\n", $load_plain
  if $load_plain == 1
    set $synth_plain = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $out_plain, 1, -1, -1, -1, -1, 200, 0)
    printf "USERDICT_FOURTH case=plain synth_ax=%d\n", $synth_plain
    set $unload_plain = ((short (*)(int))0x10027a80)(200)
    printf "USERDICT_FOURTH case=plain unload_ax=%d\n", $unload_plain
  end
  set $load_extra = ((short (*)(int, char *))0x10027960)(200, $path_extra)
  printf "USERDICT_FOURTH case=extra load_ax=%d\n", $load_extra
  if $load_extra == 1
    set $synth_extra = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $out_extra, 1, -1, -1, -1, -1, 200, 0)
    printf "USERDICT_FOURTH case=extra synth_ax=%d\n", $synth_extra
    set $unload_extra = ((short (*)(int))0x10027a80)(200)
    printf "USERDICT_FOURTH case=extra unload_ax=%d\n", $unload_extra
  end
  set $load_empty = ((short (*)(int, char *))0x10027960)(200, $path_empty)
  printf "USERDICT_FOURTH case=empty load_ax=%d\n", $load_empty
  if $load_empty == 1
    set $synth_empty = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $out_empty, 1, -1, -1, -1, -1, 200, 0)
    printf "USERDICT_FOURTH case=empty synth_ax=%d\n", $synth_empty
    set $unload_empty = ((short (*)(int))0x10027a80)(200)
    printf "USERDICT_FOURTH case=empty unload_ax=%d\n", $unload_empty
  end
  set $load_quoted = ((short (*)(int, char *))0x10027960)(200, $path_quoted)
  printf "USERDICT_FOURTH case=quoted load_ax=%d\n", $load_quoted
  if $load_quoted == 1
    set $synth_quoted = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $out_quoted, 1, -1, -1, -1, -1, 200, 0)
    printf "USERDICT_FOURTH case=quoted synth_ax=%d\n", $synth_quoted
    set $unload_quoted = ((short (*)(int))0x10027a80)(200)
    printf "USERDICT_FOURTH case=quoted unload_ax=%d\n", $unload_quoted
  end
  set {unsigned char}0x100a7489 = $gate_before
  printf "USERDICT_FOURTH gate_restored=%u\n", *(unsigned char *)0x100a7489
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

for output in "$probe/userdict-fourth-field-plain-v2.wav" \
  "$probe/userdict-fourth-field-extra-v2.wav" "$probe/userdict-fourth-field-empty-v2.wav" \
  "$probe/userdict-fourth-field-quoted-v2.wav"; do
  sha256sum "$output"
done > "$probe/userdict-fourth-field-hashes-v2.txt"
