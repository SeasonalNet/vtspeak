#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
trace="$probe/trace-makeinfo-silence-cases.gdb"
log="$probe/makeinfo-silence-cases-api.log"
listing="$probe/makeinfo-silence-cases-files.txt"
xvfb_log="$probe/xvfb-makeinfo-silence-cases.log"
xpid=
cases=(plain comma sentence ellipsis pause120)
restore() {
  cp "$probe/makeinfo-original-input1.txt" "$work/input1.txt"
  cp "$probe/makeinfo-original-output.wav" "$work/output.wav"
  cmp -s "$probe/makeinfo-original-input1.txt" "$work/input1.txt"
  cmp -s "$probe/makeinfo-original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT

for output in "$trace" "$log" "$listing" "$xvfb_log"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
for case_id in "${cases[@]}"; do
  for suffix in bin asc; do
    if [ -e "$probe/mi-silence-$case_id-20260926.$suffix.dtt" ]; then
      echo "refusing to overwrite $probe/mi-silence-$case_id-20260926.$suffix.dtt" >&2
      exit 3
    fi
  done
done

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
    '  set $text = *(unsigned int *)($esp + 8)' \
    '  disable 1'
  for case_id in "${cases[@]}"; do
    case "$case_id" in
      plain) text='Hello world.'; pause=-1 ;;
      comma) text='Hello, world.'; pause=-1 ;;
      sentence) text='Hello. World.'; pause=-1 ;;
      ellipsis) text='Hello... World.'; pause=-1 ;;
      pause120) text='Hello, world.'; pause=120 ;;
    esac
    path="Z:/work/stage16/mi-silence-$case_id-20260926"
    printf '  set $text_%s = ((char *(*)(unsigned int))0x1001d9c0)(128)\n' "$case_id"
    for ((index = 0; index < ${#text}; index++)); do
      character=${text:index:1}
      printf '  set {char}($text_%s + %d) = %d\n' "$case_id" "$index" "'${character}"
    done
    printf '  set {char}($text_%s + %d) = 0\n' "$case_id" "${#text}"
    printf '  set $path_%s = ((char *(*)(unsigned int))0x1001d9c0)(128)\n' "$case_id"
    for ((index = 0; index < ${#path}; index++)); do
      character=${path:index:1}
      printf '  set {char}($path_%s + %d) = %d\n' "$case_id" "$index" "'${character}"
    done
    printf '  set {char}($path_%s + %d) = 0\n' "$case_id" "${#path}"
    printf '  set $raw_%s = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de90)($text_%s, $path_%s, 1, -1, -1, -1, %s, -1, -1)\n' \
      "$case_id" "$case_id" "$case_id" "$pause"
    printf '  printf "MAKEINFO_SILENCE case=%s raw_eax=%%#x text=%%s path=%%s\\n", $raw_%s, $text_%s, $path_%s\n' \
      "$case_id" "$case_id" "$case_id" "$case_id"
  done
  printf '%s\n' '  continue' 'end' '' 'continue'
} > "$trace"

cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 240s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" \
  > "$log" 2>&1
find "$probe" -maxdepth 1 -type f -name 'mi-silence-*-20260926.*.dtt' \
  -printf '%f %s bytes\n' | sort > "$listing"
for case_id in "${cases[@]}"; do
  for suffix in bin asc; do
    if [ ! -s "$probe/mi-silence-$case_id-20260926.$suffix.dtt" ]; then
      echo "missing MakeInfo output for $case_id ($suffix)" >&2
      exit 4
    fi
  done
done
