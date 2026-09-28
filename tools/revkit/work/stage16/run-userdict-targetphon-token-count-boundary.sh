#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
trace="$probe/trace-userdict-targetphon-token-count-boundary-v2.gdb"
log="$probe/userdict-targetphon-token-count-boundary-v2-api.log"
xvfb_log="$probe/xvfb-userdict-targetphon-token-count-boundary-v2.log"
backup_input="$probe/userdict-targetphon-token-count-boundary-v2-original-input1.txt"
backup_output="$probe/userdict-targetphon-token-count-boundary-v2-original-output.wav"
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
  set $sequence = ((char *(*)(unsigned int))0x1001d9c0)(600)
  set $count = 120
  while $count <= 140
    set $position = 0
    set $token = 0
    while $token < $count
      set {char}($sequence + $position) = 66
      set $position = $position + 1
      set $token = $token + 1
      if $token < $count
        set {char}($sequence + $position) = 32
        set $position = $position + 1
      end
    end
    set {char}($sequence + $position) = 0
    set $result = ((short (*)(char *))0x1002a590)($sequence)
    printf "TARGETPHON token_count=%d bytes=%d result=%d\n", $count, $position, $result
    set $count = $count + 1
  end
  call ((void (*)(void *))0x1001da30)($sequence)
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
