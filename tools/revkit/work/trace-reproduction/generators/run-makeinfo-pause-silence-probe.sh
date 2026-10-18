# Generation-only source extracted from tools/revkit/work/stage21/run-makeinfo-pause-silence-probe.sh
#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
trace="$probe/trace-makeinfo-pause-silence.gdb"
log="$probe/makeinfo-pause-silence-api.log"
xvfb_log="$probe/xvfb-makeinfo-pause-silence.log"
listing="$probe/makeinfo-pause-silence-files.txt"
input_backup="$probe/makeinfo-pause-silence-original-input1.txt"
output_backup="$probe/makeinfo-pause-silence-original-output.wav"
xpid=
pauses=(-1 0 1 119 120 121 249 250 251 1000 65534 65535)
texts=("comma|A, B" "period-lower|a. b" "period-exception|n. e" \
  "period-upper-control|A. B" "ellipsis-lower|a... b" \
  "ellipsis-upper-positive|A... A" "ellipsis-upper-negative|A... B" \
  "ellipsis-mixed-control|A... a" "plain-control|A B")

case_index=0
end_case=$(( ${start_case:-0} + ${case_limit:-62208} ))
if [ "$end_case" -gt 62208 ]; then end_case=62208; fi
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
    '  disable 1'
  for item in "${texts[@]}"; do
    name=${item%%|*}
    text=${item#*|}
    for pause in "${pauses[@]}"; do
      varname="${name//-/_}_p${pause#-}"
      if [[ "$pause" == -* ]]; then varname="${name//-/_}_neg${pause#-}"; fi
      path="Z:/work/stage21/mi-pause-silence-$name-$pause"
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
      printf '  printf "MAKEINFO_PAUSE_SILENCE case=%s pause=%s raw_eax=%%#x text=%%s path=%%s\\n", $raw_%s, $text_%s, $path_%s\n' \
        "$name" "$pause" "$varname" "$varname" "$varname"
    done
  done
  printf '%s\n' '  kill' '  quit' 'end' '' 'continue'
}
