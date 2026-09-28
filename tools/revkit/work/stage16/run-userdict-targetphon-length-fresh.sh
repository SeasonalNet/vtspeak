#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
trace="$probe/trace-userdict-targetphon-length-fresh.gdb"
log="$probe/userdict-targetphon-length-fresh-api.log"
xvfb_log="$probe/xvfb-userdict-targetphon-length-fresh.log"
backup_input="$probe/userdict-targetphon-length-fresh-original-input1.txt"
backup_output="$probe/userdict-targetphon-length-fresh-original-output.wav"
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
{
  cat <<'GDB'
set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
break *0x1001da50
commands
  silent
  disable 1
GDB
  for name in n130 n130_space n131; do
    printf '  set $%s = ((char *(*)(unsigned int))0x1001d9c0)(300)\n' "$name"
  done
  cat <<'GDB'
  set $token = 0
  set $position = 0
  while $token < 130
    set {char}($n130 + $position) = 66
    set $position = $position + 1
    set $token = $token + 1
    if $token < 130
      set {char}($n130 + $position) = 32
      set $position = $position + 1
    end
  end
  set {char}($n130 + $position) = 0
  printf "TARGETPHON fresh-case=n130 bytes=%d result=%d\n", $position, ((short (*)(char *))0x1002a590)($n130)

  set $position = 0
  while $token > 0
    set {char}($n130_space + $position) = 66
    set $position = $position + 1
    set $token = $token - 1
    if $token > 0
      set {char}($n130_space + $position) = 32
      set $position = $position + 1
    end
  end
  set {char}($n130_space + $position) = 32
  set $position = $position + 1
  set {char}($n130_space + $position) = 0
  printf "TARGETPHON fresh-case=n130-trailing-space bytes=%d result=%d\n", $position, ((short (*)(char *))0x1002a590)($n130_space)

  set $token = 0
  set $position = 0
  while $token < 131
    set {char}($n131 + $position) = 66
    set $position = $position + 1
    set $token = $token + 1
    if $token < 131
      set {char}($n131 + $position) = 32
      set $position = $position + 1
    end
  end
  set {char}($n131 + $position) = 0
  printf "TARGETPHON fresh-case=n131 bytes=%d result=%d\n", $position, ((short (*)(char *))0x1002a590)($n131)
  call ((void (*)(void *))0x1001da30)($n130)
  call ((void (*)(void *))0x1001da30)($n130_space)
  call ((void (*)(void *))0x1001da30)($n131)
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
