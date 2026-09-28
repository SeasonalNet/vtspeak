#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
trace="$probe/trace-userdict-targetphon-marker-boundaries.gdb"
log="$probe/userdict-targetphon-marker-boundaries-api.log"
xvfb_log="$probe/xvfb-userdict-targetphon-marker-boundaries.log"
backup_input="$probe/userdict-targetphon-marker-boundaries-original-input1.txt"
backup_output="$probe/userdict-targetphon-marker-boundaries-original-output.wav"
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
    printf '  set {char}($%s + %d) = %d\n' \
      "$name" "$index" "'${character}"
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
  write_gdb_string exact 'HH [CI]'
  write_gdb_string lowercase 'HH [ci]'
  write_gdb_string mixedcase 'HH [Ci]'
  write_gdb_string two_spaces 'HH  [CI]'
  write_gdb_string tab_separator $'HH\t[CI]'
  write_gdb_string trailing_space 'HH [CI] '
  write_gdb_string trailing_tab $'HH [CI]\t'
  write_gdb_string following_phone 'HH [CI] HH'
  write_gdb_string repeated_marker 'HH [CI][CI]'
  write_gdb_string adjacent_marker 'HH[CI]'
  cat <<'GDB'
  set $r_exact = ((short (*)(char *))0x1002a590)($exact)
  printf "TARGETPHON boundary=exact result=%d\n", $r_exact
  set $r_lowercase = ((short (*)(char *))0x1002a590)($lowercase)
  printf "TARGETPHON boundary=lowercase result=%d\n", $r_lowercase
  set $r_mixedcase = ((short (*)(char *))0x1002a590)($mixedcase)
  printf "TARGETPHON boundary=mixedcase result=%d\n", $r_mixedcase
  set $r_two_spaces = ((short (*)(char *))0x1002a590)($two_spaces)
  printf "TARGETPHON boundary=two-spaces result=%d\n", $r_two_spaces
  set $r_tab_separator = ((short (*)(char *))0x1002a590)($tab_separator)
  printf "TARGETPHON boundary=tab-separator result=%d\n", $r_tab_separator
  set $r_trailing_space = ((short (*)(char *))0x1002a590)($trailing_space)
  printf "TARGETPHON boundary=trailing-space result=%d\n", $r_trailing_space
  set $r_trailing_tab = ((short (*)(char *))0x1002a590)($trailing_tab)
  printf "TARGETPHON boundary=trailing-tab result=%d\n", $r_trailing_tab
  set $r_following_phone = ((short (*)(char *))0x1002a590)($following_phone)
  printf "TARGETPHON boundary=following-phone result=%d\n", $r_following_phone
  set $r_repeated_marker = ((short (*)(char *))0x1002a590)($repeated_marker)
  printf "TARGETPHON boundary=repeated-marker result=%d\n", $r_repeated_marker
  set $r_adjacent_marker = ((short (*)(char *))0x1002a590)($adjacent_marker)
  printf "TARGETPHON boundary=adjacent-marker result=%d\n", $r_adjacent_marker
  call ((void (*)(void *))0x1001da30)($exact)
  call ((void (*)(void *))0x1001da30)($lowercase)
  call ((void (*)(void *))0x1001da30)($mixedcase)
  call ((void (*)(void *))0x1001da30)($two_spaces)
  call ((void (*)(void *))0x1001da30)($tab_separator)
  call ((void (*)(void *))0x1001da30)($trailing_space)
  call ((void (*)(void *))0x1001da30)($trailing_tab)
  call ((void (*)(void *))0x1001da30)($following_phone)
  call ((void (*)(void *))0x1001da30)($repeated_marker)
  call ((void (*)(void *))0x1001da30)($adjacent_marker)
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
