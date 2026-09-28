#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
trace="$probe/trace-userdict-malformed-quote-v2-effect.gdb"
log="$probe/userdict-malformed-quote-v2-effect-api.log"
xvfb_log="$probe/xvfb-userdict-malformed-quote-v2-effect.log"
backup_input="$probe/userdict-malformed-quote-v2-effect-original-input1.txt"
backup_output="$probe/userdict-malformed-quote-v2-effect-original-output.wav"
xpid=

restore() {
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
  if [ -e "$backup_input" ]; then cp "$backup_input" "$work/input1.txt"; fi
  if [ -e "$backup_output" ]; then cp "$backup_output" "$work/output.wav"; fi
  if [ -e "$backup_input" ]; then cmp -s "$backup_input" "$work/input1.txt"; fi
  if [ -e "$backup_output" ]; then cmp -s "$backup_output" "$work/output.wav"; fi
}
trap restore EXIT

outputs=(
  "$trace" "$log" "$xvfb_log" "$backup_input" "$backup_output"
  "$probe/userdict-malformed-quote-v2-control-hello.wav"
  "$probe/userdict-malformed-quote-v2-control-quoted-text.wav"
  "$probe/userdict-malformed-quote-v2-plain-row.wav"
  "$probe/userdict-malformed-quote-v2-unmatched-type.wav"
  "$probe/userdict-malformed-quote-v2-embedded-hello.wav"
  "$probe/userdict-malformed-quote-v2-embedded-quoted-text.wav"
  "$probe/userdict-malformed-quote-v2-escaped-quoted-text.wav"
  "$probe/userdict-malformed-quote-v2-hashes.txt"
)
for output in "${outputs[@]}"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done
for fixture in userdict-validation-plain-p.csv userdict-validation-unmatched-type-quote.csv \
  userdict-validation-embedded-source-quote.csv userdict-validation-escaped-source-quote.csv; do
  if [ ! -s "$probe/$fixture" ]; then
    echo "missing nonempty fixture $probe/$fixture" >&2
    exit 4
  fi
done

write_gdb_string() {
  local name=$1 value=$2 index character
  printf '  set $%s = ((char *(*)(unsigned int))0x1001d9c0)(%d)\n' \
    "$name" "$(( ${#value} + 1 ))"
  for ((index = 0; index < ${#value}; index++)); do
    character=${value:index:1}
    printf '  set {char}($%s + %d) = %d\n' "$name" "$index" "'${character}"
  done
  printf '  set {char}($%s + %d) = 0\n' "$name" "${#value}"
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
  write_gdb_string path_plain 'Z:/work/stage16/userdict-validation-plain-p.csv'
  write_gdb_string path_unmatched 'Z:/work/stage16/userdict-validation-unmatched-type-quote.csv'
  write_gdb_string path_embedded 'Z:/work/stage16/userdict-validation-embedded-source-quote.csv'
  write_gdb_string path_escaped 'Z:/work/stage16/userdict-validation-escaped-source-quote.csv'
  write_gdb_string text_hello 'hello'
  write_gdb_string text_quoted 'he"llo'
  write_gdb_string out_control_hello 'Z:/work/stage16/userdict-malformed-quote-v2-control-hello.wav'
  write_gdb_string out_control_quoted 'Z:/work/stage16/userdict-malformed-quote-v2-control-quoted-text.wav'
  write_gdb_string out_plain 'Z:/work/stage16/userdict-malformed-quote-v2-plain-row.wav'
  write_gdb_string out_unmatched 'Z:/work/stage16/userdict-malformed-quote-v2-unmatched-type.wav'
  write_gdb_string out_embedded_hello 'Z:/work/stage16/userdict-malformed-quote-v2-embedded-hello.wav'
  write_gdb_string out_embedded_quoted 'Z:/work/stage16/userdict-malformed-quote-v2-embedded-quoted-text.wav'
  write_gdb_string out_escaped_quoted 'Z:/work/stage16/userdict-malformed-quote-v2-escaped-quoted-text.wav'
  cat <<'GDB'
  set $gate_before = *(unsigned char *)0x100a7489
  set {unsigned char}0x100a7489 = 1
  printf "USERDICT_QUOTE_EFFECT gate_before=%u gate_forced=1\n", $gate_before
  set $control_hello = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text_hello, $out_control_hello, 1, -1, -1, -1, -1, -1, 0)
  printf "USERDICT_QUOTE_EFFECT case=control-hello synth_ax=%d\n", $control_hello
  set $control_quoted = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text_quoted, $out_control_quoted, 1, -1, -1, -1, -1, -1, 0)
  printf "USERDICT_QUOTE_EFFECT case=control-quoted-text synth_ax=%d\n", $control_quoted
  set $load_plain = ((short (*)(int, char *))0x10027960)(199, $path_plain)
  printf "USERDICT_QUOTE_EFFECT case=plain-row load_ax=%d\n", $load_plain
  if $load_plain == 1
    set $synth_plain = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text_hello, $out_plain, 1, -1, -1, -1, -1, 199, 0)
    printf "USERDICT_QUOTE_EFFECT case=plain-row synth_ax=%d\n", $synth_plain
    set $unload_plain = ((short (*)(int))0x10027a80)(199)
    printf "USERDICT_QUOTE_EFFECT case=plain-row unload_ax=%d\n", $unload_plain
  end
  set $load_unmatched = ((short (*)(int, char *))0x10027960)(200, $path_unmatched)
  printf "USERDICT_QUOTE_EFFECT case=unmatched-type load_ax=%d\n", $load_unmatched
  if $load_unmatched == 1
    set $synth_unmatched = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text_hello, $out_unmatched, 1, -1, -1, -1, -1, 200, 0)
    printf "USERDICT_QUOTE_EFFECT case=unmatched-type synth_ax=%d\n", $synth_unmatched
    set $unload_unmatched = ((short (*)(int))0x10027a80)(200)
    printf "USERDICT_QUOTE_EFFECT case=unmatched-type unload_ax=%d\n", $unload_unmatched
  end
  set $load_embedded = ((short (*)(int, char *))0x10027960)(201, $path_embedded)
  printf "USERDICT_QUOTE_EFFECT case=embedded-source load_ax=%d\n", $load_embedded
  if $load_embedded == 1
    set $synth_embedded_hello = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text_hello, $out_embedded_hello, 1, -1, -1, -1, -1, 201, 0)
    printf "USERDICT_QUOTE_EFFECT case=embedded-source-hello synth_ax=%d\n", $synth_embedded_hello
    set $synth_embedded_quoted = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text_quoted, $out_embedded_quoted, 1, -1, -1, -1, -1, 201, 0)
    printf "USERDICT_QUOTE_EFFECT case=embedded-source-quoted-text synth_ax=%d\n", $synth_embedded_quoted
    set $unload_embedded = ((short (*)(int))0x10027a80)(201)
    printf "USERDICT_QUOTE_EFFECT case=embedded-source unload_ax=%d\n", $unload_embedded
  end
  set $load_escaped = ((short (*)(int, char *))0x10027960)(202, $path_escaped)
  printf "USERDICT_QUOTE_EFFECT case=escaped-source load_ax=%d\n", $load_escaped
  if $load_escaped == 1
    set $synth_escaped_quoted = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text_quoted, $out_escaped_quoted, 1, -1, -1, -1, -1, 202, 0)
    printf "USERDICT_QUOTE_EFFECT case=escaped-source-quoted-text synth_ax=%d\n", $synth_escaped_quoted
    set $unload_escaped = ((short (*)(int))0x10027a80)(202)
    printf "USERDICT_QUOTE_EFFECT case=escaped-source unload_ax=%d\n", $unload_escaped
  end
  set {unsigned char}0x100a7489 = $gate_before
  printf "USERDICT_QUOTE_EFFECT gate_restored=%u\n", *(unsigned char *)0x100a7489
  continue
end

continue
GDB
} > "$trace"

cp "$work/input1.txt" "$backup_input"
cp "$work/output.wav" "$backup_output"
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 300s \
  winedbg --gdb /samples/voicetext_paul.exe < "$trace" > "$log" 2>&1

sha256sum "$probe/userdict-malformed-quote-v2-control-hello.wav" \
  "$probe/userdict-malformed-quote-v2-control-quoted-text.wav" \
  "$probe/userdict-malformed-quote-v2-plain-row.wav" \
  "$probe/userdict-malformed-quote-v2-unmatched-type.wav" \
  "$probe/userdict-malformed-quote-v2-embedded-hello.wav" \
  "$probe/userdict-malformed-quote-v2-embedded-quoted-text.wav" \
  "$probe/userdict-malformed-quote-v2-escaped-quoted-text.wav" \
  > "$probe/userdict-malformed-quote-v2-hashes.txt"
