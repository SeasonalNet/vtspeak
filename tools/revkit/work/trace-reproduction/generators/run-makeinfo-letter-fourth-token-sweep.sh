# Generation-only source extracted from tools/revkit/work/stage21/run-makeinfo-letter-fourth-token-sweep.sh
#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
trace="$probe/trace-makeinfo-letter-fourth-token-sweep.gdb"
log="$probe/makeinfo-letter-fourth-token-sweep-api.log"
xvfb_log="$probe/xvfb-letter-fourth-token-sweep.log"
listing="$probe/makeinfo-letter-fourth-token-sweep-files.txt"
masks=(UUUU UUUL ULUU ULUL)

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
  for context in aa ag; do
    for letter in {a..z}; do
      if [ "$letter" = d ]; then continue; fi
      for mask in "${masks[@]}"; do
        chars=(A a "${context:1:1}" "$letter")
        text=
        for index in 0 1 2 3; do
          character=${chars[$index]}
          if [ "${mask:$index:1}" = U ]; then
            character=${character^^}
          else
            character=${character,,}
          fi
          if [ -n "$text" ]; then text+=' '; fi
          text+=$character
        done
        for form in plain question; do
          case_name="$context-$letter-$mask-$form"
          varname=${case_name//-/_}
          call_text=$text
          if [ "$form" = question ]; then call_text+='?'; fi
          path="Z:/work/stage21/mi-letter-fourth-sweep-$case_name"
          printf '  set $text_%s = ((char *(*)(unsigned int))0x1001d9c0)(128)\n' "$varname"
          for ((offset = 0; offset < ${#call_text}; offset++)); do
            character=${call_text:offset:1}
            printf -v code '%d' "'$character"
            printf '  set {char}($text_%s + %d) = %d\n' "$varname" "$offset" "$code"
          done
          printf '  set {char}($text_%s + %d) = 0\n' "$varname" "${#call_text}"
          printf '  set $path_%s = ((char *(*)(unsigned int))0x1001d9c0)(128)\n' "$varname"
          for ((offset = 0; offset < ${#path}; offset++)); do
            character=${path:offset:1}
            printf -v code '%d' "'$character"
            printf '  set {char}($path_%s + %d) = %d\n' "$varname" "$offset" "$code"
          done
          printf '  set {char}($path_%s + %d) = 0\n' "$varname" "${#path}"
          printf '  set $raw_%s = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de90)($text_%s, $path_%s, 1, -1, -1, -1, -1, -1, -1)\n' \
            "$varname" "$varname" "$varname"
          printf '  printf "MAKEINFO_FOURTH_SWEEP case=%s raw_eax=%%#x text=%%s path=%%s\\n", $raw_%s, $text_%s, $path_%s\n' \
            "$case_name" "$varname" "$varname" "$varname"
        done
      done
    done
  done
  printf '%s\n' '  kill' '  quit' 'end' '' 'continue'
}
