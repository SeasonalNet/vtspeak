#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
trace="$probe/trace-userdict-extended-argument-matrix-v2.gdb"
log="$probe/userdict-extended-argument-matrix-v2-api.log"
xvfb_log="$probe/xvfb-userdict-extended-argument-matrix-v2.log"
backup_input="$probe/userdict-extended-argument-matrix-v2-original-input1.txt"
backup_output="$probe/userdict-extended-argument-matrix-v2-original-output.wav"
xpid=

restore() {
  cp "$backup_input" "$work/input1.txt"
  cp "$backup_output" "$work/output.wav"
  cmp -s "$backup_input" "$work/input1.txt"
  cmp -s "$backup_output" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT

for output in "$trace" "$log" "$xvfb_log" "$backup_input" "$backup_output"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
cmp -s "$probe/preprocess-flags-original-input1.txt" "$work/input1.txt"
cmp -s "$probe/preprocess-flags-original-output.wav" "$work/output.wav"

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
    '  disable 1' \
    '  set $path = ((char *(*)(unsigned int))0x1001d9c0)(35)' \
    '  set {char}($path + 0) = 90' \
    '  set {char}($path + 1) = 58' \
    '  set {char}($path + 2) = 47' \
    '  set {char}($path + 3) = 119' \
    '  set {char}($path + 4) = 111' \
    '  set {char}($path + 5) = 114' \
    '  set {char}($path + 6) = 107' \
    '  set {char}($path + 7) = 47' \
    '  set {char}($path + 8) = 115' \
    '  set {char}($path + 9) = 116' \
    '  set {char}($path + 10) = 97' \
    '  set {char}($path + 11) = 103' \
    '  set {char}($path + 12) = 101' \
    '  set {char}($path + 13) = 49' \
    '  set {char}($path + 14) = 54' \
    '  set {char}($path + 15) = 47' \
    '  set {char}($path + 16) = 117' \
    '  set {char}($path + 17) = 115' \
    '  set {char}($path + 18) = 101' \
    '  set {char}($path + 19) = 114' \
    '  set {char}($path + 20) = 100' \
    '  set {char}($path + 21) = 105' \
    '  set {char}($path + 22) = 99' \
    '  set {char}($path + 23) = 116' \
    '  set {char}($path + 24) = 45' \
    '  set {char}($path + 25) = 112' \
    '  set {char}($path + 26) = 114' \
    '  set {char}($path + 27) = 111' \
    '  set {char}($path + 28) = 98' \
    '  set {char}($path + 29) = 101' \
    '  set {char}($path + 30) = 46' \
    '  set {char}($path + 31) = 99' \
    '  set {char}($path + 32) = 115' \
    '  set {char}($path + 33) = 118' \
    '  set {char}($path + 34) = 0'
  # Allocate the CSV bytes at the exact text length plus its terminator.
  printf '%s\n' '  set $csv = ((char *(*)(unsigned int))0x1001d9c0)(11)' \
    '  set {char}($csv + 0) = 104' '  set {char}($csv + 1) = 101' \
    '  set {char}($csv + 2) = 108' '  set {char}($csv + 3) = 108' \
    '  set {char}($csv + 4) = 111' '  set {char}($csv + 5) = 44' \
    '  set {char}($csv + 6) = 72' '  set {char}($csv + 7) = 72' \
    '  set {char}($csv + 8) = 44' '  set {char}($csv + 9) = 80' \
    '  set {char}($csv + 10) = 0'
  case_number=0
  for path_mode in null file; do
    if [ "$path_mode" = null ]; then path_arg=0; else path_arg='$path'; fi
    for buffer_mode in null data; do
      if [ "$buffer_mode" = null ]; then buffer_arg=0; else buffer_arg='$csv'; fi
      for length in 0 1 9 10 11 -1; do
        index=$((300 + case_number))
        case_number=$((case_number + 1))
        suffix="p_${path_mode}_b_${buffer_mode}_n_${length#-}"
        [ "$length" -ge 0 ] || suffix="${suffix}_neg"
        printf '  set $result_%s = ((short (*)(int, char *, char *, int))0x10027880)(%d, %s, %s, %d)\n' \
          "$case_number" "$index" "$path_arg" "$buffer_arg" "$length"
        printf '  printf "USERDICT_EXT_MATRIX case=%s index=%d ax=%%d\\n", $result_%s\n' \
          "$suffix" "$index" "$case_number"
        printf '  if $result_%s == 1\n' "$case_number"
        printf '    set $unload_%s = ((short (*)(int))0x10027980)(%d)\n' \
          "$case_number" "$index"
        printf '    printf "USERDICT_EXT_MATRIX case=%s index=%d unload_ax=%%d\\n", $unload_%s\n' \
          "$suffix" "$index" "$case_number"
        printf '%s\n' '  else' \
          "    set \$failed_unload_$case_number = ((short (*)(int))0x10027980)($index)" \
          "    printf \"USERDICT_EXT_MATRIX case=$suffix index=$index empty_slot_unload_ax=%d\\n\", \$failed_unload_$case_number" \
          '  end'
      done
    done
  done
  printf '%s\n' '  continue' 'end' '' 'continue'
} > "$trace"

cp "$work/input1.txt" "$backup_input"
cp "$work/output.wav" "$backup_output"
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 300s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" > "$log" 2>&1
