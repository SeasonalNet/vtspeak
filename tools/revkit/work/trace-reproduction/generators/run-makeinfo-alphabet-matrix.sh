# Generation-only source extracted from tools/revkit/work/stage21/run-makeinfo-alphabet-matrix.sh
#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
trace="$probe/trace-makeinfo-alphabet-matrix-v2.gdb"
log="$probe/makeinfo-alphabet-matrix-v2-api.log"
xvfb_log="$probe/xvfb-makeinfo-alphabet-matrix-v2.log"
listing="$probe/makeinfo-alphabet-matrix-v2-files.txt"
xpid=

case_index=0
end_case=$(( ${start_case:-0} + ${case_limit:-62208} ))
if [ "$end_case" -gt 62208 ]; then end_case=62208; fi
names=()
texts=()
for letter in {A..Z}; do
  names+=("letter-${letter,,}")
  texts+=("$letter")
done
names+=(upper-run spaced-letters dotted-letters acronym)
texts+=("ABC" "A B C" "A. B. C." "U.S.A.")
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
  for index in "${!names[@]}"; do
    name=${names[$index]}
    varname=${name//-/_}
    text=${texts[$index]}
    path="Z:/work/stage21/mi-alphabet-v2-$name"
    printf '  set $text_%s = ((char *(*)(unsigned int))0x1001d9c0)(128)\n' "$varname"
    for ((offset = 0; offset < ${#text}; offset++)); do
      character=${text:offset:1}
      printf '  set {char}($text_%s + %d) = %d\n' "$varname" "$offset" "'$character"
    done
    printf '  set {char}($text_%s + %d) = 0\n' "$varname" "${#text}"
    printf '  set $path_%s = ((char *(*)(unsigned int))0x1001d9c0)(128)\n' "$varname"
    for ((offset = 0; offset < ${#path}; offset++)); do
      character=${path:offset:1}
      printf '  set {char}($path_%s + %d) = %d\n' "$varname" "$offset" "'$character"
    done
    printf '  set {char}($path_%s + %d) = 0\n' "$varname" "${#path}"
    printf '  set $raw_%s = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de90)($text_%s, $path_%s, 1, -1, -1, -1, -1, -1, -1)\n' \
      "$varname" "$varname" "$varname"
    printf '  printf "MAKEINFO_ALPHABET case=%s raw_eax=%%#x text=%%s path=%%s\\n", $raw_%s, $text_%s, $path_%s\n' \
      "$name" "$varname" "$varname" "$varname"
  done
  printf '%s\n' '  kill' '  quit' 'end' '' 'continue'
}
