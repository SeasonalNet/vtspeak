#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
trace="$probe/trace-userdict-targetphon-marker-sweep.gdb"
log="$probe/userdict-targetphon-marker-sweep-api.log"
xvfb_log="$probe/xvfb-userdict-targetphon-marker-sweep.log"
backup_input="$probe/userdict-targetphon-marker-sweep-original-input1.txt"
backup_output="$probe/userdict-targetphon-marker-sweep-original-output.wav"
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
  write_gdb_string ci '[CI]'
  write_gdb_string ci_after 'HH [CI]'
  write_gdb_string ci_embedded 'HH[CI]B'
  write_gdb_string skip '[SKIP]'
  write_gdb_string skip_after 'HH [SKIP]'
  write_gdb_string bracket '['
  write_gdb_string bracket_word '[OTHER]'
  write_gdb_string hash_sequence '# HH'
  cat <<'GDB'
  set $result_ci = ((short (*)(char *))0x1002a590)($ci)
  printf "TARGETPHON marker=ci-only result=%d\n", $result_ci
  set $result_ci_after = ((short (*)(char *))0x1002a590)($ci_after)
  printf "TARGETPHON marker=ci-after-phone result=%d\n", $result_ci_after
  set $result_ci_embedded = ((short (*)(char *))0x1002a590)($ci_embedded)
  printf "TARGETPHON marker=ci-embedded result=%d\n", $result_ci_embedded
  set $result_skip = ((short (*)(char *))0x1002a590)($skip)
  printf "TARGETPHON marker=skip-only result=%d\n", $result_skip
  set $result_skip_after = ((short (*)(char *))0x1002a590)($skip_after)
  printf "TARGETPHON marker=skip-after-phone result=%d\n", $result_skip_after
  set $result_bracket = ((short (*)(char *))0x1002a590)($bracket)
  printf "TARGETPHON marker=open-bracket result=%d\n", $result_bracket
  set $result_other = ((short (*)(char *))0x1002a590)($bracket_word)
  printf "TARGETPHON marker=other-bracket result=%d\n", $result_other
  set $result_hash = ((short (*)(char *))0x1002a590)($hash_sequence)
  printf "TARGETPHON marker=hash-sequence result=%d\n", $result_hash
  call ((void (*)(void *))0x1001da30)($ci)
  call ((void (*)(void *))0x1001da30)($ci_after)
  call ((void (*)(void *))0x1001da30)($ci_embedded)
  call ((void (*)(void *))0x1001da30)($skip)
  call ((void (*)(void *))0x1001da30)($skip_after)
  call ((void (*)(void *))0x1001da30)($bracket)
  call ((void (*)(void *))0x1001da30)($bracket_word)
  call ((void (*)(void *))0x1001da30)($hash_sequence)
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
