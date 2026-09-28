#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
tag=matrix-v3
trace="$probe/trace-userdict-validation-$tag.gdb"
log="$probe/userdict-validation-$tag-api.log"
xvfb_log="$probe/xvfb-userdict-validation-$tag.log"
backup_input="$probe/userdict-validation-$tag-original-input1.txt"
backup_output="$probe/userdict-validation-$tag-original-output.wav"
xpid=
cases=(plain-p four-p four-a quoted-p two-fields five-fields empty-file empty-source empty-target-p empty-target-a invalid-type long-type trailing-empty crlf-p)
if [ "$#" -gt 0 ]; then
  tag="case-$1"
  trace="$probe/trace-userdict-validation-$tag.gdb"
  log="$probe/userdict-validation-$tag-api.log"
  xvfb_log="$probe/xvfb-userdict-validation-$tag.log"
  backup_input="$probe/userdict-validation-$tag-original-input1.txt"
  backup_output="$probe/userdict-validation-$tag-original-output.wav"
  cases=("$@")
fi

restore() {
  cp "$backup_input" "$work/input1.txt"
  cp "$backup_output" "$work/output.wav"
  cmp -s "$backup_input" "$work/input1.txt"
  cmp -s "$backup_output" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT

for output in "$trace" "$log" "$xvfb_log" "$backup_input" "$backup_output"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
cmp -s "$probe/preprocess-flags-original-input1.txt" "$work/input1.txt"
cmp -s "$probe/preprocess-flags-original-output.wav" "$work/output.wav"
write_fixture() {
  local id=$1 path="$probe/userdict-validation-$1.csv"
  if [ -e "$path" ]; then
    case "$id" in
      plain-p) cmp -s <(printf 'hello,HH,P\n') "$path" ;;
      four-p) cmp -s <(printf 'hello,HH,P,extra\n') "$path" ;;
      four-a) cmp -s <(printf 'hello,world,A,extra\n') "$path" ;;
      quoted-p) cmp -s <(printf '"hello","HH","P"\n') "$path" ;;
      two-fields) cmp -s <(printf 'hello,HH\n') "$path" ;;
      five-fields) cmp -s <(printf 'hello,HH,P,extra,tail\n') "$path" ;;
      empty-file) cmp -s <(printf '') "$path" ;;
      empty-source) cmp -s <(printf ',HH,P\n') "$path" ;;
      empty-target-p) cmp -s <(printf 'hello,,P\n') "$path" ;;
      empty-target-a) cmp -s <(printf 'hello,,A\n') "$path" ;;
      invalid-type) cmp -s <(printf 'hello,HH,X\n') "$path" ;;
      long-type) cmp -s <(printf 'hello,HH,PP\n') "$path" ;;
      trailing-empty) cmp -s <(printf 'hello,HH,P,\n') "$path" ;;
      crlf-p) cmp -s <(printf 'hello,HH,P\r\n') "$path" ;;
      unmatched-type-quote) cmp -s <(printf 'hello,HH,"P\n') "$path" ;;
      embedded-source-quote) cmp -s <(printf 'he"llo,HH,P\n') "$path" ;;
      escaped-source-quote) cmp -s <(printf '"he""llo",HH,P\n') "$path" ;;
      invalid-target-p) cmp -s <(printf 'hello,NOTAPHONE,P\n') "$path" ;;
      invalid-target-example) cmp -s <(printf 'hello,example,P\n') "$path" ;;
    esac || { echo "fixture exists with unexpected contents: $path" >&2; return 3; }
    return 0
  fi
  case "$id" in
    plain-p) printf 'hello,HH,P\n' > "$path" ;;
    four-p) printf 'hello,HH,P,extra\n' > "$path" ;;
    four-a) printf 'hello,world,A,extra\n' > "$path" ;;
    quoted-p) printf '"hello","HH","P"\n' > "$path" ;;
    two-fields) printf 'hello,HH\n' > "$path" ;;
    five-fields) printf 'hello,HH,P,extra,tail\n' > "$path" ;;
    empty-file) : > "$path" ;;
    empty-source) printf ',HH,P\n' > "$path" ;;
    empty-target-p) printf 'hello,,P\n' > "$path" ;;
    empty-target-a) printf 'hello,,A\n' > "$path" ;;
    invalid-type) printf 'hello,HH,X\n' > "$path" ;;
    long-type) printf 'hello,HH,PP\n' > "$path" ;;
    trailing-empty) printf 'hello,HH,P,\n' > "$path" ;;
    crlf-p) printf 'hello,HH,P\r\n' > "$path" ;;
    unmatched-type-quote) printf 'hello,HH,"P\n' > "$path" ;;
    embedded-source-quote) printf 'he"llo,HH,P\n' > "$path" ;;
    escaped-source-quote) printf '"he""llo",HH,P\n' > "$path" ;;
    invalid-target-p) printf 'hello,NOTAPHONE,P\n' > "$path" ;;
    invalid-target-example) printf 'hello,example,P\n' > "$path" ;;
    *) echo "unknown fixture: $id" >&2; return 2 ;;
  esac
}
for id in "${cases[@]}"; do write_fixture "$id"; done

write_gdb_string() {
  local var_name=$1 value=$2 index character
  printf '  set $%s = ((char *(*)(unsigned int))0x1001d9c0)(%d)\n' \
    "$var_name" "$(( ${#value} + 1 ))"
  for ((index = 0; index < ${#value}; index++)); do
    character=${value:index:1}
    printf '  set {char}($%s + %d) = %d\n' \
      "$var_name" "$index" "'${character}"
  done
  printf '  set {char}($%s + %d) = 0\n' "$var_name" "${#value}"
}

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
  write_gdb_string valid_path 'Z:/work/stage16/userdict-validation-plain-p.csv'
  case_number=0
  for id in "${cases[@]}"; do
    var=${id//-/_}
    index=$((100 + case_number))
    case_number=$((case_number + 1))
    write_gdb_string "path_$var" "Z:/work/stage16/userdict-validation-$id.csv"
    printf '  set $load_%s = ((short (*)(int, char *))0x10027960)(%d, $path_%s)\n' \
      "$var" "$index" "$var"
    printf '  printf "USERDICT_MATRIX case=%s index=%d load_ax=%%d\\n", $load_%s\n' \
      "$id" "$index" "$var"
    printf '  set $unload_%s = ((short (*)(int))0x10027a80)(%d)\n' "$var" "$index"
    printf '  printf "USERDICT_MATRIX case=%s index=%d unload_ax=%%d\\n", $unload_%s\n' \
      "$id" "$index" "$var"
    printf '  if $load_%s != 1\n' "$var"
    printf '    set $recover_%s = ((short (*)(int, char *))0x10027960)(%d, $valid_path)\n' \
      "$var" "$index"
    printf '    printf "USERDICT_MATRIX case=%s index=%d valid_recovery_ax=%%d\\n", $recover_%s\n' \
      "$id" "$index" "$var"
    printf '    if $recover_%s == 1\n' "$var"
    printf '      set $recover_unload_%s = ((short (*)(int))0x10027a80)(%d)\n' \
      "$var" "$index"
    printf '      printf "USERDICT_MATRIX case=%s index=%d recovery_unload_ax=%%d\\n", $recover_unload_%s\n' \
      "$id" "$index" "$var"
    printf '%s\n' '    end' '  end'
  done
  printf '%s\n' '  continue' 'end' '' 'continue'
} > "$trace"

cp "$probe/preprocess-flags-original-input1.txt" "$backup_input"
cp "$probe/preprocess-flags-original-output.wav" "$backup_output"
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 300s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" > "$log" 2>&1
