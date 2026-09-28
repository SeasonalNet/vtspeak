#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
trace="$probe/trace-userdict-targetphon-pair-sweep.gdb"
log="$probe/userdict-targetphon-pair-sweep-api.log"
xvfb_log="$probe/xvfb-userdict-targetphon-pair-sweep.log"
backup_input="$probe/userdict-targetphon-pair-sweep-original-input1.txt"
backup_output="$probe/userdict-targetphon-pair-sweep-original-output.wav"
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

cat > "$trace" <<'GDB'
set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
break *0x1001da50
commands
  silent
  disable 1
  set $phon = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set $first = 65
  while $first <= 90
    set $second = 65
    while $second <= 90
      set {unsigned char}($phon) = $first
      set {unsigned char}($phon + 1) = $second
      set {char}($phon + 2) = 0
      set $result = ((short (*)(char *))0x1002a590)($phon)
      printf "TARGETPHON pair=%c%c result=%d\n", $first, $second, $result
      set $second = $second + 1
    end
    set $first = $first + 1
  end
  call ((void (*)(void *))0x1001da30)($phon)
  continue
end
continue
GDB

cp "$work/input1.txt" "$backup_input"
cp "$work/output.wav" "$backup_output"
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 300s \
  winedbg --gdb /samples/voicetext_paul.exe < "$trace" > "$log" 2>&1
