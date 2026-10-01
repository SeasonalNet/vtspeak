#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
trace="$probe/trace-negative-pause-setter-grid.gdb"
log="$probe/negative-pause-setter-grid-api.log"
xvfb_log="$probe/xvfb-negative-pause-setter-grid.log"
input_backup="$probe/negative-pause-setter-grid-original-input1.txt"
output_backup="$probe/negative-pause-setter-grid-original-output.wav"
xpid=
negative_values=(-2147483648 -2 -1)
fields=(pitch speed volume sentence comma)

if [ "$PWD" != "$work" ]; then
  echo "run this probe with Compose working directory $work" >&2
  exit 2
fi
for output in "$trace" "$log" "$xvfb_log" "$input_backup" "$output_backup"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done

cp -p "$work/input1.txt" "$input_backup"
cp -p "$work/output.wav" "$output_backup"
restore() {
  cp -p "$input_backup" "$work/input1.txt"
  cp -p "$output_backup" "$work/output.wav"
  cmp -s "$input_backup" "$work/input1.txt"
  cmp -s "$output_backup" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT
trap 'exit 129' HUP
trap 'exit 130' INT
trap 'exit 143' TERM

{
  printf '%s\n' 'set pagination off' 'set confirm off' \
    'set debuginfod enabled off' 'handle SIGSEGV nostop noprint pass' '' \
    'break *0x1001da50' 'commands' '  silent' \
    '  set $speaker = *(int *)($esp + 16)' '  disable 1' \
    '  set $out = (int *)malloc(16)' '  set $comma_out = (int *)malloc(4)' \
    '  call ((void (*)(int, int, int, int, int))0x10027fe0)(123, 234, 345, 456, $speaker)' \
    '  call ((void (*)(int, int))0x100281b0)(567, $speaker)' \
    '  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)' \
    '  set $comma_result = ((int (*)(int *, int))0x100281f0)($comma_out, $speaker)' \
    '  printf "NEG_PAUSE stage=baseline getter_ret=%d pitch=%d speed=%d volume=%d sentence=%d comma_ret=%d comma=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $comma_result, *$comma_out'
  for field in "${fields[@]}"; do
    for value in "${negative_values[@]}"; do
      case "$field" in
        pitch) args="$value, 234, 345, 456" ;;
        speed) args="123, $value, 345, 456" ;;
        volume) args="123, 234, $value, 456" ;;
        sentence) args="123, 234, 345, $value" ;;
      esac
      if [ "$field" = comma ]; then
        printf '  call ((void (*)(int, int))0x100281b0)(%d, $speaker)\n' "$value"
      else
        printf '  call ((void (*)(int, int, int, int, int))0x10027fe0)(%s, $speaker)\n' "$args"
      fi
      printf '  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)\n'
      printf '  set $comma_result = ((int (*)(int *, int))0x100281f0)($comma_out, $speaker)\n'
      printf '  printf "NEG_PAUSE field=%s input=%d getter_ret=%%d pitch=%%d speed=%%d volume=%%d sentence=%%d comma_ret=%%d comma=%%d\\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $comma_result, *$comma_out\n' "$field" "$value"
    done
  done
  printf '%s\n' '  call ((void (*)(int, int, int, int, int))0x10027fe0)(123, 234, 345, 456, $speaker)' \
    '  call ((void (*)(int, int))0x100281b0)(567, $speaker)' \
    '  kill' '  quit' 'end' '' 'continue'
} > "$trace"

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 300s \
  winedbg --gdb /samples/voicetext_paul.exe < "$trace" > "$log" 2>&1
