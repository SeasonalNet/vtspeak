# Generation-only source extracted from tools/revkit/work/stage21/run-makeinfo-ascii-whitespace-class-grid-all-pairs.sh
#!/bin/bash
set -euo pipefail

probe=/work/stage21
work="$probe/sandbox/stage5"
start_case=${START_CASE:-0}
suffix=
if [[ "$start_case" =~ ^[0-9]+$ ]] && [ "$start_case" -gt 0 ]; then suffix="-resume-$start_case"; fi
trace="$probe/trace-makeinfo-ascii-whitespace-class-grid-all-pairs$suffix.gdb"
log="$probe/makeinfo-ascii-whitespace-class-grid-all-pairs-api$suffix.log"
xvfb_log="$probe/xvfb-ascii-whitespace-class-grid-all-pairs$suffix.log"
listing="$probe/makeinfo-ascii-whitespace-class-grid-all-pairs-files.txt"
input_backup="$probe/ascii-whitespace-class-all-pairs-original-input1.txt"
output_backup="$probe/ascii-whitespace-class-all-pairs-original-output.wav"
xpid=
words=(a i hello world hi kate paul good morning weather today voice)
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
        for side in "${sides[@]}"; do
          for ((sep_index = 0; sep_index < ${#seps[@]}; sep_index++)); do
            sep_name=${sep_names[sep_index]}
            if [ "$case_index" -lt "$start_case" ]; then
              case_index=$((case_index + 1))
              continue
            fi
            if [ "$side" = before ]; then
              text="$left${seps[sep_index]}$punctuation $right"
            else
              text="$left$punctuation${seps[sep_index]}$right"
            fi
            varname=reused
            path="Z:/work/stage21/mi-ws12-$case_index"
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
  done
  printf '%s\n' '  kill' '  quit' 'end' '' 'continue'
}
