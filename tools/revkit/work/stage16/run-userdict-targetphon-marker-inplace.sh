#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
trace="$probe/trace-userdict-targetphon-marker-inplace.gdb"
log="$probe/userdict-targetphon-marker-inplace-api.log"
xvfb_log="$probe/xvfb-userdict-targetphon-marker-inplace.log"
backup_input="$probe/userdict-targetphon-marker-inplace-original-input1.txt"
backup_output="$probe/userdict-targetphon-marker-inplace-original-output.wav"
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
  if [ -e "$output" ]; then echo "refusing to overwrite $output" >&2; exit 3; fi
done
write_gdb_string() {
  local name=$1 value=$2 index character
  printf '  set $%s = ((char *(*)(unsigned int))0x1001d9c0)(%d)\n' "$name" "$(( ${#value} + 1 ))"
  for ((index = 0; index < ${#value}; index++)); do
    character=${value:index:1}
    printf '  set {char}($%s + %d) = %d\n' "$name" "$index" "'${character}"
  done
  printf '  set {char}($%s + %d) = 0\n' "$name" "${#value}"
}
{
  printf '%s\n' 'set pagination off' 'set confirm off' 'set debuginfod enabled off' \
    'handle SIGSEGV nostop noprint pass' 'break *0x1001da50' 'commands' \
    '  silent' '  disable 1'
  write_gdb_string padded $'\t HH \r'
  write_gdb_string marker_space 'HH [CI] '
  write_gdb_string marker_adjacent 'HH[CI]'
  write_gdb_string marker_extra 'HH [CI] HH'
  write_gdb_string invalid_marker 'HH [SKIP]'
  cat <<'GDB'
  set $r = ((short (*)(char *))0x1002a590)($padded)
  printf "TARGETPHON inplace=padded result=%d text=<%s>\n", $r, $padded
  set $r = ((short (*)(char *))0x1002a590)($marker_space)
  printf "TARGETPHON inplace=marker-space result=%d text=<%s>\n", $r, $marker_space
  set $r = ((short (*)(char *))0x1002a590)($marker_adjacent)
  printf "TARGETPHON inplace=marker-adjacent result=%d text=<%s>\n", $r, $marker_adjacent
  set $r = ((short (*)(char *))0x1002a590)($marker_extra)
  printf "TARGETPHON inplace=marker-extra result=%d text=<%s>\n", $r, $marker_extra
  set $r = ((short (*)(char *))0x1002a590)($invalid_marker)
  printf "TARGETPHON inplace=invalid-marker result=%d text=<%s>\n", $r, $invalid_marker
  call ((void (*)(void *))0x1001da30)($padded)
  call ((void (*)(void *))0x1001da30)($marker_space)
  call ((void (*)(void *))0x1001da30)($marker_adjacent)
  call ((void (*)(void *))0x1001da30)($marker_extra)
  call ((void (*)(void *))0x1001da30)($invalid_marker)
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
