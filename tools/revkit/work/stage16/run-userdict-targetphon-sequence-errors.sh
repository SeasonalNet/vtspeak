#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
trace="$probe/trace-userdict-targetphon-sequence-errors-v2.gdb"
log="$probe/userdict-targetphon-sequence-errors-v2-api.log"
xvfb_log="$probe/xvfb-userdict-targetphon-sequence-errors-v2.log"
backup_input="$probe/userdict-targetphon-sequence-errors-v2-original-input1.txt"
backup_output="$probe/userdict-targetphon-sequence-errors-v2-original-output.wav"
xpid=

restore() {
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
  if [ -e "$backup_input" ]; then cp "$backup_input" "$work/input1.txt"; fi
  if [ -e "$backup_output" ]; then cp "$backup_output" "$work/output.wav"; fi
  if [ -e "$backup_input" ]; then cmp -s "$backup_input" "$work/input1.txt"; fi
  if [ -e "$backup_output" ]; then cmp -s "$backup_output" "$work/output.wav"; fi
}
trap restore EXIT

for output in "$trace" "$log" "$xvfb_log" "$backup_input" "$backup_output"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
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
    'break *0x1001da50' \
    'commands' \
    '  silent' \
    '  disable 1'
  write_gdb_string empty ''
  write_gdb_string spaces '   '
  write_gdb_string single 'HH'
  write_gdb_string sequence 'HH B AA1 CH'
  write_gdb_string repeated_spaces '  HH   B  '
  write_gdb_string tabs $'\tHH\tB\t'
  write_gdb_string bad_token 'HH EXAMPLE'
  write_gdb_string hash '#'
  write_gdb_string bracket '['
  cat <<'GDB'
  set $result_null = ((short (*)(char *))0x1002a590)(0)
  printf "TARGETPHON case=null result=%d\n", $result_null
  set $result_empty = ((short (*)(char *))0x1002a590)($empty)
  printf "TARGETPHON case=empty result=%d\n", $result_empty
  set $result_spaces = ((short (*)(char *))0x1002a590)($spaces)
  printf "TARGETPHON case=spaces result=%d\n", $result_spaces
  set $result_single = ((short (*)(char *))0x1002a590)($single)
  printf "TARGETPHON case=single result=%d\n", $result_single
  set $result_sequence = ((short (*)(char *))0x1002a590)($sequence)
  printf "TARGETPHON case=sequence result=%d\n", $result_sequence
  set $result_repeated = ((short (*)(char *))0x1002a590)($repeated_spaces)
  printf "TARGETPHON case=repeated-spaces result=%d\n", $result_repeated
  set $result_tabs = ((short (*)(char *))0x1002a590)($tabs)
  printf "TARGETPHON case=tabs result=%d\n", $result_tabs
  set $result_bad = ((short (*)(char *))0x1002a590)($bad_token)
  printf "TARGETPHON case=bad-token result=%d\n", $result_bad
  set $result_hash = ((short (*)(char *))0x1002a590)($hash)
  printf "TARGETPHON case=hash-marker result=%d\n", $result_hash
  set $result_bracket = ((short (*)(char *))0x1002a590)($bracket)
  printf "TARGETPHON case=open-bracket result=%d\n", $result_bracket

  set $long = ((char *(*)(unsigned int))0x1001d9c0)(263)
  set $i = 0
  while $i < 129
    set {char}($long + $i * 2) = 66
    set {char}($long + $i * 2 + 1) = 32
    set $i = $i + 1
  end
  set {char}($long + 258) = 66
  set {char}($long + 259) = 0
  set $actual_length = 0
  while *(char *)($long + $actual_length) != 0
    set $actual_length = $actual_length + 1
  end
  printf "TARGETPHON measured_length=%d\\n", $actual_length
  set $result_259 = ((short (*)(char *))0x1002a590)($long)
  printf "TARGETPHON case=length-259 result=%d\n", $result_259
  set {char}($long + 259) = 32
  set {char}($long + 260) = 0
  set $result_260 = ((short (*)(char *))0x1002a590)($long)
  printf "TARGETPHON case=length-260 result=%d\n", $result_260
  set {char}($long + 260) = 66
  set {char}($long + 261) = 0
  set $result_261 = ((short (*)(char *))0x1002a590)($long)
  printf "TARGETPHON case=length-261 result=%d\n", $result_261
  call ((void (*)(void *))0x1001da30)($long)
  call ((void (*)(void *))0x1001da30)($empty)
  call ((void (*)(void *))0x1001da30)($spaces)
  call ((void (*)(void *))0x1001da30)($single)
  call ((void (*)(void *))0x1001da30)($sequence)
  call ((void (*)(void *))0x1001da30)($repeated_spaces)
  call ((void (*)(void *))0x1001da30)($tabs)
  call ((void (*)(void *))0x1001da30)($bad_token)
  call ((void (*)(void *))0x1001da30)($hash)
  call ((void (*)(void *))0x1001da30)($bracket)
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
