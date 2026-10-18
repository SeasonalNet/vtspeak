# Generation-only source extracted from tools/revkit/work/stage21/run-makeinfo-letter-triple-a-position-matrix.sh
#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
letter_case=${LETTER_CASE:-upper}
positions=${POSITIONS:-middle,last}
positions=${positions//,/ }
case "$letter_case" in
  upper)
    case_label=upper
    first_upper=1
    middle_upper=1
    last_upper=1
    ;;
  lower)
    case_label=lower
    first_upper=0
    middle_upper=0
    last_upper=0
    ;;
  upper-lower)
    case_label=mixed-upper-lower
    first_upper=1
    middle_upper=0
    last_upper=0
    ;;
  lower-upper)
    case_label=mixed-lower-upper
    first_upper=0
    middle_upper=1
    last_upper=1
    ;;
  lower-lower-upper)
    case_label=lower-lower-upper
    first_upper=0
    middle_upper=0
    last_upper=1
    ;;
  lower-upper-lower)
    case_label=lower-upper-lower
    first_upper=0
    middle_upper=1
    last_upper=0
    ;;
  upper-lower-upper)
    case_label=upper-lower-upper
    first_upper=1
    middle_upper=0
    last_upper=1
    ;;
  upper-upper-lower)
    case_label=upper-upper-lower
    first_upper=1
    middle_upper=1
    last_upper=0
    ;;
  *)
    echo "LETTER_CASE must be a three-token case pattern such as upper-lower-upper" >&2
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
for position in $positions; do
  for left in {a..z}; do
    left_text=$left
      left_upper=$first_upper
      if [ "$position" = leading ]; then
        left_upper=$middle_upper
        # Reproduce the two quarantined pre-correction command files exactly.
        if [ "${LEGACY_LEADING_CASE:-0}" = 1 ]; then left_upper=$first_upper; fi
      fi
      if [ "$left_upper" -eq 1 ]; then
        left_text=${left^^}
      fi
      for right in {a..z}; do
        right_text=$right
        right_upper=$last_upper
        if [ "$position" = last ]; then right_upper=$middle_upper; fi
        if [ "$right_upper" -eq 1 ]; then
          right_text=${right^^}
        fi
        fixed_text=a
        if [ "$position" = leading ]; then
          if [ "$first_upper" -eq 1 ]; then fixed_text=A; fi
          first_token=$fixed_text
          second_token=$left_text
          third_token=$right_text
        elif [ "$position" = middle ]; then
          if [ "$middle_upper" -eq 1 ]; then fixed_text=A; fi
          first_token=$left_text
          second_token=$fixed_text
          third_token=$right_text
        else
          if [ "$last_upper" -eq 1 ]; then fixed_text=A; fi
          first_token=$left_text
          second_token=$right_text
          third_token=$fixed_text
        fi
        for form in plain question; do
          name="$position-$left$right-$form"
          varname=${name//-/_}
          text="$first_token $second_token $third_token"
          if [ "$form" = question ]; then text+='?'; fi
          if [ "$position" = leading ]; then
            capture_prefix="mi-letter-triple-a-$case_label"
          else
            capture_prefix="mi-letter-triple-a-$position-$case_label"
          fi
          path="Z:/work/stage21/$capture_prefix-$left$right-$form"
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
          printf '  printf "MAKEINFO_TRIPLE_POSITION case=%s raw_eax=%%#x text=%%s path=%%s\\n", $raw_%s, $text_%s, $path_%s\n' \
            "$name" "$varname" "$varname" "$varname"
        done
      done
    done
  done
  printf '%s\n' '  kill' '  quit' 'end' '' 'continue'
}
