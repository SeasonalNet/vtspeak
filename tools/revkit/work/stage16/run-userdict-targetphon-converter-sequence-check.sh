#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
trace="$probe/trace-userdict-targetphon-converter-sequence-check.gdb"
log="$probe/userdict-targetphon-converter-sequence-check-api.log"
xvfb_log="$probe/xvfb-userdict-targetphon-converter-sequence-check.log"
backup_input="$probe/userdict-targetphon-converter-sequence-check-original-input1.txt"
backup_output="$probe/userdict-targetphon-converter-sequence-check-original-output.wav"
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
  write_gdb_string first12 'B D F G K L M N P R S T'
  write_gdb_string rest24 'V W Y Z CH DH HH JH NG SH TH ZH AA0 AE0 AH0 AO0 AW0 AY0 EH0 ER0 EY0 IH0 IY0 OW0'
  write_gdb_string all36 'B D F G K L M N P R S T V W Y Z CH DH HH JH NG SH TH ZH AA0 AE0 AH0 AO0 AW0 AY0 EH0 ER0 EY0 IH0 IY0 OW0'
  printf '%s\n' '  set $output = ((char *(*)(unsigned int))0x1001d9c0)(66)' \
    '  set $i = 0' '  while $i < 66' '    set {char}($output + $i) = 0' \
    '    set $i = $i + 1' '  end'
  for spec in 'first12 12' 'rest24 24' 'all36 36'; do
    read -r name count <<< "$spec"
    printf '  set $result = ((short (*)(char *, char *))0x1005f710)($output, $%s)\n' "$name"
    printf '  printf "TARGETPHON sequence=%s result=%%d bytes=", $result\n' "$name"
    printf '  set $i = 0\n  while $i < %d\n' "$count"
    printf '    printf "%%02x", *(unsigned char *)($output + $i)\n'
    printf '    set $i = $i + 1\n  end\n  printf "\\n"\n'
  done
  cat <<'GDB'
  call ((void (*)(void *))0x1001da30)($first12)
  call ((void (*)(void *))0x1001da30)($rest24)
  call ((void (*)(void *))0x1001da30)($all36)
  call ((void (*)(void *))0x1001da30)($output)
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
