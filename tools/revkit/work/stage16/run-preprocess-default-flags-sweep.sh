#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
base=preprocess-default-flags-20260928-
prefix="Z:/work/stage16/$base"
trace="$probe/trace-preprocess-default-flags-sweep.gdb"
log="$probe/preprocess-default-flags-sweep-api.log"
listing="$probe/preprocess-default-flags-sweep-files.txt"
xvfb_log="$probe/xvfb-preprocess-default-flags-sweep.log"
before=$(mktemp)
after=$(mktemp)
xpid=
restore() {
  cp "$probe/preprocess-flags-original-input1.txt" "$work/input1.txt"
  cp "$probe/preprocess-flags-original-output.wav" "$work/output.wav"
  cmp -s "$probe/preprocess-flags-original-input1.txt" "$work/input1.txt"
  cmp -s "$probe/preprocess-flags-original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
  rm -f "$before" "$after"
}
trap restore EXIT

for output in "$trace" "$log" "$listing" "$xvfb_log"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
for ((flag = 12; flag <= 253; flag++)); do
  output=$(printf '%s%03d' "$probe/$base" "$flag")
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
  for ((index = 0; index < ${#prefix}; index++)); do
    character=${prefix:index:1}
    printf '  set {char}($path + %d) = %d\n' "$index" "'${character}"
  done
  printf '  set {char}($path + %d) = 0\n' "${#prefix}"
  printf '%s\n' \
    "  set {char}(\$path + ${#prefix} + 3) = 0" \
    '  set $flag = 12' \
    '  while $flag <= 253' \
    "    set {char}(\$path + ${#prefix}) = 48 + (\$flag / 100)" \
    "    set {char}(\$path + $(( ${#prefix} + 1 ))) = 48 + ((\$flag / 10) % 10)" \
    "    set {char}(\$path + $(( ${#prefix} + 2 ))) = 48 + (\$flag % 10)" \
    '    set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, $flag, 1, -1, -1, -1, -1, 0, 0)' \
    '    printf "PREPROCESS_DEFAULT_SWEEP flag=%d raw_eax=%#x low_ax=%d path=%s\\n", $flag, $raw, ((short)$raw), $path' \
    '    set $flag = $flag + 1' \
    '  end' \
    '  continue' \
    'end' \
    '' \
    'continue'
} > "$trace"

cp "$probe/preprocess-flags-original-input1.txt" "$work/input1.txt"
find "$probe" -maxdepth 1 -type f -printf '%f %s %T@\n' | sort > "$before"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 300s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" > "$log" 2>&1
find "$probe" -maxdepth 1 -type f -printf '%f %s %T@\n' | sort > "$after"
comm -13 "$before" "$after" > "$listing"

for ((flag = 12; flag <= 253; flag++)); do
  output=$(printf '%s%03d' "$probe/$base" "$flag")
  if [ -f "$output" ]; then
    printf 'generated flag=%d bytes=%s sha256=%s\n' \
      "$flag" "$(wc -c < "$output")" "$(sha256sum "$output" | cut -d ' ' -f 1)"
  else
    printf 'generated flag=%d absent\n' "$flag"
  fi
done >> "$log"
