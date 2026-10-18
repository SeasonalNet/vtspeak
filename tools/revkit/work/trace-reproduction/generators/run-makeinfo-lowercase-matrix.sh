# Generation-only source extracted from tools/revkit/work/stage21/run-makeinfo-lowercase-matrix.sh
#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
trace="$probe/trace-makeinfo-lowercase-matrix.gdb"
log="$probe/makeinfo-lowercase-matrix-api.log"
xvfb_log="$probe/xvfb-makeinfo-lowercase-matrix.log"
listing="$probe/makeinfo-lowercase-matrix-files.txt"
xpid=

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
  for letter in {a..z}; do
    varname="letter_$letter"
    path="Z:/work/stage21/mi-alphabet-lower-$letter"
    printf '  set $text_%s = ((char *(*)(unsigned int))0x1001d9c0)(128)\n' "$varname"
    printf '  set {char}($text_%s + 0) = %d\n' "$varname" "'$letter"
    printf '  set {char}($text_%s + 1) = 0\n' "$varname"
    printf '  set $path_%s = ((char *(*)(unsigned int))0x1001d9c0)(128)\n' "$varname"
    for ((offset = 0; offset < ${#path}; offset++)); do
      character=${path:offset:1}
      printf '  set {char}($path_%s + %d) = %d\n' "$varname" "$offset" "'$character"
    done
    printf '  set {char}($path_%s + %d) = 0\n' "$varname" "${#path}"
    printf '  set $raw_%s = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de90)($text_%s, $path_%s, 1, -1, -1, -1, -1, -1, -1)\n' \
      "$varname" "$varname" "$varname"
    printf '  printf "MAKEINFO_LOWER letter=%s raw_eax=%%#x text=%%s path=%%s\\n", $raw_%s, $text_%s, $path_%s\n' \
      "$letter" "$varname" "$varname" "$varname"
  done
  printf '%s\n' '  kill' '  quit' 'end' '' 'continue'
}
