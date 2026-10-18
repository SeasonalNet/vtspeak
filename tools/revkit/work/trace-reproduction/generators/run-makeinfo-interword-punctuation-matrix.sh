# Generation-only source extracted from tools/revkit/work/stage21/run-makeinfo-interword-punctuation-matrix.sh
#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
mark=${1:?usage: run-makeinfo-interword-punctuation-matrix.sh comma|period|ellipsis}
case "$mark" in
  comma) punctuation=, ;;
  period) punctuation=. ;;
  ellipsis) punctuation=... ;;
  *) echo "unknown punctuation: $mark" >&2; exit 2 ;;
esac
case_masks=(UU UL LU LL)
trace="$probe/trace-makeinfo-interword-punctuation-$mark.gdb"
log="$probe/makeinfo-interword-punctuation-$mark-api.log"
xvfb_log="$probe/xvfb-interword-punctuation-$mark.log"
listing="$probe/makeinfo-interword-punctuation-$mark-files.txt"
input_backup="$probe/interword-punctuation-$mark-original-input1.txt"
output_backup="$probe/interword-punctuation-$mark-original-output.wav"
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
  for first in {a..z}; do
    for second in {a..z}; do
      for mask in "${case_masks[@]}"; do
        left=$first
        right=$second
        if [ "${mask:0:1}" = U ]; then left=${left^^}; else left=${left,,}; fi
        if [ "${mask:1:1}" = U ]; then right=${right^^}; else right=${right,,}; fi
        case_name="$first$second-$mask"
        varname=${case_name//-/_}
        text="$left$punctuation $right"
        path="Z:/work/stage21/mi-interword-punct-$mark-$mask-$first$second"
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
        printf '  set $raw_%s = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de90)($text_%s, $path_%s, 1, -1, -1, -1, -1, -1, -1)\n' \
          "$varname" "$varname" "$varname"
        printf '  printf "MAKEINFO_INTERWORD_PUNCT mark=%s case=%s raw_eax=%%#x text=%%s path=%%s\\n", $raw_%s, $text_%s, $path_%s\n' \
          "$mark" "$case_name" "$varname" "$varname" "$varname"
      done
    done
  done
  printf '%s\n' '  kill' '  quit' 'end' '' 'continue'
}
