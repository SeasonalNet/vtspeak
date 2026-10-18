# Generation-only source extracted from tools/revkit/work/stage21/run-makeinfo-ascii-whitespace-class-grid.sh
#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
trace="$probe/trace-makeinfo-ascii-whitespace-class-grid.gdb"
log="$probe/makeinfo-ascii-whitespace-class-grid-api.log"
xvfb_log="$probe/xvfb-ascii-whitespace-class-grid.log"
listing="$probe/makeinfo-ascii-whitespace-class-grid-files.txt"
input_backup="$probe/ascii-whitespace-class-original-input1.txt"
output_backup="$probe/ascii-whitespace-class-original-output.wav"
xpid=
left_words=(a a hello paul good hi a weather)
right_words=(a hello world hi morning kate good today)
masks=(UU UL LU LL)
marks=(comma period ellipsis)
seps=($'\v' $'\f' $'\r')
sep_names=(vt ff cr)
sides=(before after)

case_index=0
end_case=$(( ${start_case:-0} + ${case_limit:-62208} ))
if [ "$end_case" -gt 62208 ]; then end_case=62208; fi
{
  printf '%s\n' 'set pagination off' 'set confirm off' \
    'set debuginfod enabled off' 'handle SIGSEGV nostop noprint pass' '' \
    'break *0x1001da50' 'commands' '  silent' '  disable 1'
  for ((pair = 0; pair < ${#left_words[@]}; pair++)); do
    left_base=${left_words[pair]}
    right_base=${right_words[pair]}
    for mark in "${marks[@]}"; do
      case "$mark" in
        comma) punctuation=, ;;
        period) punctuation=. ;;
        ellipsis) punctuation=... ;;
      esac
      for mask in "${masks[@]}"; do
        left=$left_base
        right=$right_base
        if [ "${mask:0:1}" = U ]; then left=${left^^}; fi
        if [ "${mask:1:1}" = U ]; then right=${right^^}; fi
        for side in "${sides[@]}"; do
          for ((sep_index = 0; sep_index < ${#seps[@]}; sep_index++)); do
            sep_name=${sep_names[sep_index]}
            if [ "$side" = before ]; then
              text="$left${seps[sep_index]}$punctuation $right"
            else
              text="$left$punctuation${seps[sep_index]}$right"
            fi
            varname="case${case_index}_${side}_${sep_name}_${mark}_${mask}"
            path="Z:/work/stage21/mi-word-wsclass-$side-$sep_name-$mark-$mask-$left_base-$right_base"
            printf '  set $text_%s = ((char *(*)(unsigned int))0x1001d9c0)(256)\n' "$varname"
            for ((offset = 0; offset < ${#text}; offset++)); do
              character=${text:offset:1}
              printf -v code '%d' "'$character"
              printf '  set {char}($text_%s + %d) = %d\n' "$varname" "$offset" "$code"
            done
            printf '  set {char}($text_%s + %d) = 0\n' "$varname" "${#text}"
            printf '  set $path_%s = ((char *(*)(unsigned int))0x1001d9c0)(256)\n' "$varname"
            for ((offset = 0; offset < ${#path}; offset++)); do
              character=${path:offset:1}
              printf -v code '%d' "'$character"
              printf '  set {char}($path_%s + %d) = %d\n' "$varname" "$offset" "$code"
            done
            printf '  set {char}($path_%s + %d) = 0\n' "$varname" "${#path}"
            printf '  set $raw_%s = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de90)($text_%s, $path_%s, 1, -1, -1, -1, -1, -1, -1)\n' \
              "$varname" "$varname" "$varname"
            printf '  printf "MAKEINFO_ASCII_WSCLASS pair=%d side=%s sep=%s mark=%s mask=%s raw_eax=%%#x\\n", $raw_%s\n' \
              "$pair" "$side" "$sep_name" "$mark" "$mask" "$varname"
            case_index=$((case_index + 1))
          done
        done
      done
    done
  done
  printf '%s\n' '  kill' '  quit' 'end' '' 'continue'
}
