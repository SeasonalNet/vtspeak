#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
trace="$probe/trace-userdict-targetphon-converter-map.gdb"
log="$probe/userdict-targetphon-converter-map-api.log"
xvfb_log="$probe/xvfb-userdict-targetphon-converter-map.log"
backup_input="$probe/userdict-targetphon-converter-map-original-input1.txt"
backup_output="$probe/userdict-targetphon-converter-map-original-output.wav"
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
phones=(B D F G K L M N P R S T V W Y Z CH DH HH JH NG SH TH ZH)
for stem in AA AE AH AO AW AY EH ER EY IH IY OW OY UH UW; do
  for stress in 0 1 2; do phones+=("${stem}${stress}"); done
done
phones+=("#")
{
  printf '%s\n' 'set pagination off' 'set confirm off' 'set debuginfod enabled off' \
    'handle SIGSEGV nostop noprint pass' 'break *0x1001da50' 'commands' \
    '  silent' '  disable 1'
  for i in "${!phones[@]}"; do write_gdb_string "phone$i" "${phones[$i]}"; done
  printf '%s\n' '  set $output = ((char *(*)(unsigned int))0x1001d9c0)(66)' \
    '  set $i = 0' '  while $i < 66' \
    '    set {char}($output + $i) = 0' \
    '    set $i = $i + 1' '  end'
  for i in "${!phones[@]}"; do
    printf '  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone%d)\n' "$i"
    printf '  printf "TARGETPHON converter token=%s result=%%d code=%%02x next=%%02x\\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)\n' "${phones[$i]}"
  done
  for i in "${!phones[@]}"; do printf '  call ((void (*)(void *))0x1001da30)($phone%d)\n' "$i"; done
  cat <<'GDB'
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
