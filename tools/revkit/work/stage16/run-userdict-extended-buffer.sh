#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
trace="$probe/trace-userdict-extended-buffer-v2.gdb"
log="$probe/userdict-extended-buffer-v2-api.log"
xvfb_log="$probe/xvfb-userdict-extended-buffer-v2.log"
backup_input="$probe/userdict-extended-buffer-v2-original-input1.txt"
backup_output="$probe/userdict-extended-buffer-v2-original-output.wav"
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
    '  set $csv = ((char *(*)(unsigned int))0x1001d9c0)(11)' \
    '  set {char}($csv + 0) = 104' \
    '  set {char}($csv + 1) = 101' \
    '  set {char}($csv + 2) = 108' \
    '  set {char}($csv + 3) = 108' \
    '  set {char}($csv + 4) = 111' \
    '  set {char}($csv + 5) = 44' \
    '  set {char}($csv + 6) = 72' \
    '  set {char}($csv + 7) = 72' \
    '  set {char}($csv + 8) = 44' \
    '  set {char}($csv + 9) = 80' \
    '  set {char}($csv + 10) = 0' \
    '  set $init_before = *(int *)0x100a0458' \
    '  printf "USERDICT_EXT init_before=0x%x data_len=10\\n", $init_before' \
    '  set $load = ((short (*)(int, char *, char *, int))0x10027880)(200, 0, $csv, 10)' \
    '  printf "USERDICT_EXT memory_load index=200 path=NULL length=10 ax=%d\\n", $load' \
    '  set $duplicate = ((short (*)(int, char *, char *, int))0x10027880)(200, 0, $csv, 10)' \
    '  printf "USERDICT_EXT duplicate_load index=200 ax=%d\\n", $duplicate' \
    '  set $unload = ((short (*)(int))0x10027980)(200)' \
    '  printf "USERDICT_EXT unload index=200 ax=%d\\n", $unload' \
    '  set {int}0x100a0458 = 0' \
    '  set $uninitialized = ((short (*)(int, char *, char *, int))0x10027880)(201, 0, $csv, 10)' \
    '  printf "USERDICT_EXT forced_uninitialized index=201 ax=%d\\n", $uninitialized' \
    '  set $uninitialized_unload = ((short (*)(int))0x10027980)(201)' \
    '  printf "USERDICT_EXT forced_uninitialized_unload index=201 ax=%d\\n", $uninitialized_unload' \
    '  set {int}0x100a0458 = $init_before' \
    '  printf "USERDICT_EXT init_restored=0x%x\\n", *(int *)0x100a0458' \
    '  continue' \
    'end' \
    '' \
    'continue'
} > "$trace"

cp "$work/input1.txt" "$backup_input"
cp "$work/output.wav" "$backup_output"
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 240s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" > "$log" 2>&1
