#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
trace="$probe/trace-text-file-path-error-grid.gdb"
log="$probe/text-file-path-error-grid-api.log"
xvfb_log="$probe/xvfb-text-file-path-error-grid.log"
xpid=
paths=(
  "Z:/work/stage21/text-file-path-probe-valid.wav"
  "Z:/work/stage21/sandbox/stage5"
  "Z:/work/stage21/text-file-path-missing-parent-v1/output.wav"
)
names=(valid directory missing_parent)

if [ "$PWD" != "$work" ]; then
  echo "run this probe with Compose working directory $work" >&2
  exit 2
fi
for output in "$trace" "$log" "$xvfb_log" \
  /work/stage21/text-file-path-probe-valid.wav \
  /work/stage21/text-file-path-missing-parent-v1; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite existing path $output" >&2
    exit 3
  fi
done

{
  printf '%s\n' 'set pagination off' 'set confirm off' \
    'set debuginfod enabled off' 'handle SIGSEGV nostop noprint pass' '' \
    'break *0x1001da50' 'commands' '  silent' \
    '  set $text = *(char **)($esp + 8)' \
    '  set $speaker = *(int *)($esp + 16)' \
    '  disable 1' '  set $probe_path = (char *)malloc(512)'
  for i in "${!paths[@]}"; do
    path=${paths[$i]}
    name=${names[$i]}
    printf '  set $probe_path = (char *)malloc(%d)\n' "$(( ${#path} + 1 ))"
    for ((offset = 0; offset < ${#path}; offset++)); do
      character=${path:offset:1}
      printf -v code '%d' "'$character"
      printf '  set {char}($probe_path + %d) = %d\n' "$offset" "$code"
    done
    printf '  set {char}($probe_path + %d) = 0\n' "${#path}"
    printf '  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $probe_path, $speaker, -1, -1, -1, -1, -1, -1)\n'
    printf '  printf "TEXT_FILE_PATH name=%s result=%%d\\n", $result\n' "$name"
  done
  printf '%s\n' '  kill' '  quit' 'end' '' 'continue'
} > "$trace"

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
trap 'if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi' EXIT
trap 'exit 129' HUP
trap 'exit 130' INT
trap 'exit 143' TERM
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 300s \
  winedbg --gdb /samples/voicetext_paul.exe < "$trace" > "$log" 2>&1
