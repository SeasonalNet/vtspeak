#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
base=heap-makeinfo-path-20260926
path="Z:/work/stage16/$base"
trace="$probe/trace-makeinfo-heap-path.gdb"
log="$probe/makeinfo-heap-path-api.log"
listing="$probe/makeinfo-heap-path-files.txt"
xvfb_log="$probe/xvfb-makeinfo-heap-path.log"
xpid=
restore() {
  cp "$probe/makeinfo-original-input1.txt" "$work/input1.txt"
  cp "$probe/makeinfo-original-output.wav" "$work/output.wav"
  cmp -s "$probe/makeinfo-original-input1.txt" "$work/input1.txt"
  cmp -s "$probe/makeinfo-original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT

for output in "$trace" "$log" "$listing" "$xvfb_log" \
    "$probe/$base.bin.dtt" "$probe/$base.asc.dtt"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
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
    '  disable 1' \
    '  set $path = ((char *(*)(unsigned int))0x1001d9c0)(128)'
  for ((index = 0; index < ${#path}; index++)); do
    character=${path:index:1}
    printf '  set {char}($path + %d) = %d\n' "$index" "'$character"
  done
  printf '  set {char}($path + %d) = 0\n' "${#path}"
  printf '%s\n' \
    '  printf "MAKEINFO_HEAP_PATH_BEFORE ptr=%#x value=", $path' \
    '  x/s $path'
  printf '  set $raw = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de90)((char *)$text, $path, 1, -1, -1, -1, -1, -1, -1)\n'
  printf '  printf "MAKEINFO_HEAP_PATH_AFTER raw_eax=%%#x signed=%%d value=", $raw, $raw\n'
  printf '%s\n' \
    '  x/s $path' \
    '  continue' \
    'end' \
    '' \
    'continue'
} > "$trace"

cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" \
  > "$log" 2>&1
if [ "$(grep -F -c "\"$path\"" "$log")" -ne 2 ]; then
  echo "API call did not preserve the complete path before and after" >&2
  exit 4
fi
find "$probe" -maxdepth 1 -type f -name "$base*.dtt" -printf '%f %s bytes\n' \
  > "$listing"
if ! cmp -s "$probe/$base.bin.dtt" "$probe/lead6-makeinfo-plead6.bin.dtt" || \
    ! cmp -s "$probe/$base.asc.dtt" "$probe/lead6-makeinfo-plead6.asc.dtt"; then
  echo "heap-backed output contents differ from the earlier pair" >&2
  exit 5
fi
