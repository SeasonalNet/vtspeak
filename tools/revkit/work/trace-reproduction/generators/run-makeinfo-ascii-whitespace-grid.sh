# Generation-only source extracted from tools/revkit/work/stage21/run-makeinfo-ascii-whitespace-grid.sh
#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
start_case=${START_CASE:-0}
suffix=
if [ "$start_case" -gt 0 ]; then suffix="-resume-$start_case"; fi
trace="$probe/trace-makeinfo-ascii-whitespace-grid$suffix.gdb"
log="$probe/makeinfo-ascii-whitespace-grid-api$suffix.log"
xvfb_log="$probe/xvfb-ascii-whitespace-grid$suffix.log"
listing="$probe/makeinfo-ascii-whitespace-grid-files.txt"
input_backup="$probe/ascii-whitespace-grid-original-input1.txt"
output_backup="$probe/ascii-whitespace-grid-original-output.wav"
xpid=
left_words=(a a hello paul good hi a weather)
right_words=(a hello world hi morning kate good today)
masks=(UU UL LU LL)
marks=(comma period ellipsis)
seps=("" " " "  " $'\t' $'\n' $'\r\n')
sep_names=(none space double tab lf crlf)

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
        for ((pre_index = 0; pre_index < ${#seps[@]}; pre_index++)); do
          for ((post_index = 0; post_index < ${#seps[@]}; post_index++)); do
            if [ "$case_index" -lt "$start_case" ]; then
              case_index=$((case_index + 1))
              continue
            fi
            case_index=$((case_index + 1))
            pre_name=${sep_names[pre_index]}
            post_name=${sep_names[post_index]}
            text="$left${seps[pre_index]}$punctuation${seps[post_index]}$right"
            varname="pair${pair}_${mark}_${mask}_${pre_name}_${post_name}"
            path="Z:/work/stage21/mi-word-wsgrid-$mark-$mask-$left_base-$right_base-pre-$pre_name-post-$post_name"
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
            printf '  printf "MAKEINFO_ASCII_WS pair=%d mark=%s mask=%s pre=%s post=%s raw_eax=%%#x\\n", $raw_%s\n' \
              "$pair" "$mark" "$mask" "$pre_name" "$post_name" "$varname"
          done
        done
      done
    done
  done
  printf '%s\n' '  kill' '  quit' 'end' '' 'continue'
}
