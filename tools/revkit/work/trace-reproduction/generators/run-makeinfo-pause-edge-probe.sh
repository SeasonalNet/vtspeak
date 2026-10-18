# Generation-only source extracted from tools/revkit/work/stage21/run-makeinfo-pause-edge-probe.sh
#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
trace="$probe/trace-makeinfo-pause-edges.gdb"
log="$probe/makeinfo-pause-edges-api.log"
xvfb_log="$probe/xvfb-makeinfo-pause-edges.log"
listing="$probe/makeinfo-pause-edges-files.txt"
input_backup="$probe/makeinfo-pause-edges-original-input1.txt"
output_backup="$probe/makeinfo-pause-edges-original-output.wav"
xpid=
pauses=(-2147483648 -100 -2 -1 65535 65536 65537 100000 2147483647)
texts=("period-lower|a. b" "ellipsis-lower|a... b" "ellipsis-upper-positive|A... A")

case_index=0
end_case=$(( ${start_case:-0} + ${case_limit:-62208} ))
if [ "$end_case" -gt 62208 ]; then end_case=62208; fi
{
  printf '%s\n' 'set pagination off' 'set confirm off' \
    'set debuginfod enabled off' 'handle SIGSEGV nostop noprint pass' '' \
    'break *0x1001da50' 'commands' '  silent' '  disable 1'
  for item in "${texts[@]}"; do
    name=${item%%|*}
    text=${item#*|}
    for pause in "${pauses[@]}"; do
      safe_pause=${pause//-/neg}
      varname="${name//-/_}_p$safe_pause"
      path="Z:/work/stage21/mi-pause-edge-$name-$pause"
      printf '  set $text_%s = ((char *(*)(unsigned int))0x1001d9c0)(128)\n' "$varname"
      for ((offset = 0; offset < ${#text}; offset++)); do
        character=${text:offset:1}
        printf -v code '%d' "'$character"
        printf '  set {char}($text_%s + %d) = %d\n' "$varname" "$offset" "$code"
      done
      printf '  set {char}($text_%s + %d) = 0\n' "$varname" "${#text}"
      printf '  set $path_%s = ((char *(*)(unsigned int))0x1001d9c0)(128)\n' "$varname"
      for ((offset = 0; offset < ${#path}; offset++)); do
        character=${path:offset:1}
        printf -v code '%d' "'$character"
        printf '  set {char}($path_%s + %d) = %d\n' "$varname" "$offset" "$code"
      done
      printf '  set {char}($path_%s + %d) = 0\n' "$varname" "${#path}"
      printf '  set $raw_%s = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de90)($text_%s, $path_%s, 1, -1, -1, -1, %s, -1, -1)\n' \
        "$varname" "$varname" "$varname" "$pause"
      printf '  printf "MAKEINFO_PAUSE_EDGE case=%s pause=%s raw_eax=%%#x\\n", $raw_%s\n' \
        "$name" "$pause" "$varname"
    done
  done
  printf '%s\n' '  kill' '  quit' 'end' '' 'continue'
}
