#!/bin/bash
set -euo pipefail

probe=/work/stage21
work=/work/stage21/sandbox/stage5
trace="$probe/trace-userdict-memory-boundaries-v3.gdb"
log="$probe/userdict-memory-boundaries-v3-api.log"
xvfb_log="$probe/xvfb-userdict-memory-boundaries-v3.log"
backup_input="$probe/userdict-memory-boundaries-v3-original-input1.txt"
backup_output="$probe/userdict-memory-boundaries-v3-original-output.wav"
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

for output in "$trace" "$log" "$xvfb_log" "$backup_input" "$backup_output"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
test -f "$work/input1.txt"
test -f "$work/output.wav"

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
  case_no=0
  for spec in \
    'ten_span_trailing_a5|10|104 101 108 108 111 44 72 72 44 80 165 0' \
    'eleven_span_trailing_a5|11|104 101 108 108 111 44 72 72 44 80 165 0' \
    'embedded_nul_tail|14|104 101 108 108 111 44 72 72 44 80 0 44 88 0' \
    'malformed_second_row|17|104 101 108 108 111 44 72 72 44 80 10 120 44 72 72 44 88 0'; do
    tag=${spec%%|*}
    rest=${spec#*|}
    length=${rest%%|*}
    bytes=${rest#*|}
    count=0
    for _ in $bytes; do count=$((count + 1)); done
    index=$((400 + case_no))
    case_no=$((case_no + 1))
    printf '  set $buf_%d = ((char *(*)(unsigned int))0x1001d9c0)(%d)\n' "$case_no" "$count"
    offset=0
    for byte in $bytes; do
      printf '  set {unsigned char}($buf_%d + %d) = %d\n' "$case_no" "$offset" "$byte"
      offset=$((offset + 1))
    done
    printf '  set $load_%d = ((short (*)(int, char *, char *, int))0x10027880)(%d, 0, $buf_%d, %d)\n' "$case_no" "$index" "$case_no" "$length"
    printf '  printf "USERDICT_MEMORY_BOUNDARY case=%s index=%d span=%d physical=%d load_ax=%%d\\n", $load_%d\n' "$tag" "$index" "$length" "$count" "$case_no"
    printf '  if $load_%d == 1\n' "$case_no"
    printf '    set $unload_%d = ((short (*)(int))0x10027980)(%d)\n' "$case_no" "$index"
    printf '    printf "USERDICT_MEMORY_BOUNDARY case=%s index=%d unload_ax=%%d\\n", $unload_%d\n' "$tag" "$index" "$case_no"
    printf '  else\n  set $empty_%d = ((short (*)(int))0x10027980)(%d)\n' "$case_no" "$index"
    printf '    printf "USERDICT_MEMORY_BOUNDARY case=%s index=%d empty_slot_unload_ax=%%d\\n", $empty_%d\n' "$tag" "$index" "$case_no"
    printf '  end\n'
  done
  printf '%s\n' '  continue' 'end' '' 'continue'
} > "$trace"

cp "$work/input1.txt" "$backup_input"
cp "$work/output.wav" "$backup_output"
cp "$probe/input.txt" "$work/input1.txt"
cd "$work"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 300s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" > "$log" 2>&1
