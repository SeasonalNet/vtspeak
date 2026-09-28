#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage16
run_id=${1:-}
xpid=
backup_input=/tmp/vtspeak-csv-parser-deep-input1.txt
backup_output=/tmp/vtspeak-csv-parser-deep-output.wav
log="$probe/csv-parser-deep${run_id:+-$run_id}-api.log"
trace="/tmp/trace-csv-parser-deep${run_id:+-$run_id}.gdb"

if [ -e "$log" ]; then
  echo "refusing to overwrite $log" >&2
  exit 1
fi

restore() {
  cp "$backup_input" "$work/input1.txt"
  cp "$backup_output" "$work/output.wav"
  cmp -s "$backup_input" "$work/input1.txt"
  cmp -s "$backup_output" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

cp "$work/input1.txt" "$backup_input"
cp "$work/output.wav" "$backup_output"
trap restore EXIT
cp "$probe/input.txt" "$work/input1.txt"

cat > "$trace" <<'GDB_HEADER'
set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
GDB_HEADER

cases=(
  'quoted_empty=2222'
  'quoted_empty_middle=612c22222c63'
  'text_after_close=226122782c62'
  'quote_before_open=6122622c63'
  'unclosed_quoted_after_comma=612c22622c63'
  'quote_only=22'
  'leading_lf=0a612c62'
  'trailing_lf=612c620a'
  'newline_in_quoted=22610a62222c63'
  'cr_only_between=612c620d632c64'
  'crlf_between=612c620d0a632c64'
  'crlf_leading=0d0a612c62'
  'lf_after_delimiter=612c0a62'
  'double_lf_between=610a0a62'
  'quoted_eof=226122'
  'quoted_comma_trailing=2261222c'
  'quoted_postspace=226122202c62'
  'space_before_quote=2022612c6222'
  'quote_after_text=616222632c64'
  'quote_at_unquoted_eof=61226222'
  'bare_cr_eof=612c620d'
  'bare_cr_start=0d612c62'
  'bare_cr_after_comma=612c0d62'
  'bare_cr_in_quoted=22610d62222c63'
  'crlf_after_comma=612c0d0a622c63'
  'crlf_after_text=610d0a622c63'
  'cr_before_comma=610d2c62'
  'crlf_trailing=612c620d0a'
  'empty_input='
  'delimiter_only=2c'
  'leading_delimiter=2c61'
  'two_trailing_delimiters=612c2c'
  'only_delimiters=2c2c'
  'whitespace_only=2020090a0d'
  'whitespace_before_delimiter=202c61'
  'tabs_around_quoted=092022612c62222009'
  'space_then_quoted=612c2022612c6222'
  'bare_cr_before_eof_after_prefix=612c620d63'
  'bare_cr_before_delimiter=612c620d2c63'
  'bare_cr_multiple_following_fields=612c620d632c642c65'
  'bare_cr_after_delimiter_two_fields=612c0d622c63'
)

for item in "${cases[@]}"; do
  case_name=${item%%=*}
  hex_bytes=${item#*=}
  byte_count=$((${#hex_bytes} / 2))
  printf '  set $parser = ((void *(*)(void))0x10016bc0)()\n' >> "$trace"
  printf '  set $line = (char *)malloc(%d)\n' "$((byte_count + 1))" >> "$trace"
  for ((index = 0; index < byte_count; index++)); do
    byte_hex=${hex_bytes:index * 2:2}
    byte_value=$((16#$byte_hex))
    printf '  set {unsigned char}($line + %d) = %d\n' "$index" "$byte_value" >> "$trace"
  done
  printf '  set {unsigned char}($line + %d) = 0\n' "$byte_count" >> "$trace"
  printf '  set $parsed = ((int (*)(void *, char *, int))0x10016bd0)($parser, $line, 0)\n' >> "$trace"
  printf '  set $count = ((int (*)(void *))0x10016c10)($parser)\n' >> "$trace"
  printf '  printf "CSV_DEEP case=%s low_ax=%%d count=%%d", $parsed, $count\n' "$case_name" >> "$trace"
  cat >> "$trace" <<'GDB_FIELDS'
  set $i = 0
  while $i < $count
    set $field = ((char *(*)(void *, int))0x10016c30)($parser, $i)
    printf " field%d=", $i
    set $j = 0
    while $j < 64 && *(unsigned char *)($field + $j) != 0
      printf "%02x", *(unsigned char *)($field + $j)
      set $j = $j + 1
    end
    set $i = $i + 1
  end
  printf "\n"
  call ((void (*)(void *))0x10016bf0)($parser)
GDB_FIELDS
done

cat >> "$trace" <<'GDB_FOOTER'
  continue
end

continue
GDB_FOOTER

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-csv-parser-deep.log" 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" > "$log" 2>&1

if [ "$(grep -c '^CSV_DEEP case=' "$log")" -ne "${#cases[@]}" ]; then
  echo "did not capture all ${#cases[@]} CSV parser cases" >&2
  exit 1
fi
