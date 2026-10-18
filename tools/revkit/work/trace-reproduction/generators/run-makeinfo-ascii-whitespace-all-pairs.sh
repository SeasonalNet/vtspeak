# Generation-only source extracted from tools/revkit/work/stage21/run-makeinfo-ascii-whitespace-all-pairs.sh
#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
start_case=${START_CASE:-0}
case_limit=${CASE_LIMIT:-62208}
suffix=
trace=
log=
xvfb_log=
listing=
input_backup="$probe/ascii-whitespace-all-pairs-original-input1.txt"
output_backup="$probe/ascii-whitespace-all-pairs-original-output.wav"
xpid=
words=(a i hello world hi kate paul good morning weather today voice)
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
  for ((left_index = 0; left_index < ${#words[@]}; left_index++)); do
    left_base=${words[left_index]}
    for ((right_index = 0; right_index < ${#words[@]}; right_index++)); do
      right_base=${words[right_index]}
      pair=$((left_index * ${#words[@]} + right_index))
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
              if [ "$case_index" -lt "$start_case" ] || [ "$case_index" -ge "$end_case" ]; then
                case_index=$((case_index + 1))
                continue
              fi
              pre_name=${sep_names[pre_index]}
              post_name=${sep_names[post_index]}
              text="$left${seps[pre_index]}$punctuation${seps[post_index]}$right"
              varname=reused
              path="Z:/work/stage21/mi-wsfull12-$case_index"
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
              printf '  printf "MAKEINFO_ASCII_WS_FULL pair=%d mark=%s mask=%s pre=%s post=%s raw_eax=%%#x\\n", $raw_%s\n' \
                "$pair" "$mark" "$mask" "$pre_name" "$post_name" "$varname"
              case_index=$((case_index + 1))
            done
          done
        done
      done
    done
  done
  printf '%s\n' '  kill' '  quit' 'end' '' 'continue'
}
