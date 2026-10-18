# Generation-only source extracted from tools/revkit/work/stage21/run-makeinfo-letter-pair-punctuation-matrix.sh
#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
letter_case=${LETTER_CASE:-upper}
case "$letter_case" in
  upper)
    case_label=upper
    left_upper=1
    right_upper=1
    ;;
  lower)
    case_label=lower
    left_upper=0
    right_upper=0
    ;;
  upper-lower)
    case_label=mixed-upper-lower
    left_upper=1
    right_upper=0
    ;;
  lower-upper)
    case_label=mixed-lower-upper
    left_upper=0
    right_upper=1
    ;;
  *)
    echo "LETTER_CASE must be upper, lower, upper-lower, or lower-upper" >&2
    exit 2
    ;;
esac

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
  for left in {a..z}; do
    left_text=$left
    if [ "$left_upper" -eq 1 ]; then
      left_text=${left^^}
    fi
    for right in {a..z}; do
      right_text=$right
      if [ "$right_upper" -eq 1 ]; then
        right_text=${right^^}
      fi
      for mark in period comma bang; do
        case "$mark" in
          period) punctuation='.' ;;
          comma) punctuation=',' ;;
          bang) punctuation='!' ;;
        esac
        name="$left$right-$mark"
        varname=${name//-/_}
        text="$left_text $right_text$punctuation"
        path="Z:/work/stage21/mi-letter-pair-punct-$case_label-$left$right-$mark"
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
        printf '  printf "MAKEINFO_PAIR_PUNCT case=%s raw_eax=%%#x text=%%s path=%%s\\n", $raw_%s, $text_%s, $path_%s\n' \
          "$name" "$varname" "$varname" "$varname"
      done
    done
  done
  printf '%s\n' '  kill' '  quit' 'end' '' 'continue'
}
