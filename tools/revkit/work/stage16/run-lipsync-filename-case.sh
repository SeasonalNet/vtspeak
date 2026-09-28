#!/bin/bash
set -euo pipefail

case_name=${1:-}
case "$case_name" in
  extension) name='vtspeak-lipname.txt' ;;
  dot-relative) name='./vtspeak-lipdot' ;;
  subdir-relative) name='lipsync-probe/vtspeak-liprelative.txt' ;;
  subdir-relative-verified) name='lipsync-probe/vtspeak-liprelative.txt' ;;
  subdir-simple) name='lipsync-probe/rel.txt' ;;
  subdir-control) name='lipsync-probe/report.txt' ;;
  backslash-relative) name='lipsync-probe\backslash.txt' ;;
  missing-parent) name='no-such-lipsync-dir/report.txt' ;;
  absolute-z) name='Z:/work/stage16/lipsync-filename-bytewise-absolute-z-output.txt' ;;
  absolute-c) name='C:\windows\temp\vtspeak-lipsync-c-drive-output.txt' ;;
  absolute-c-verified) name='C:\windows\temp\vtspeak-lipsync-c-drive-output-verified.txt' ;;
  *) echo 'usage: run-lipsync-filename-case.sh extension|dot-relative|subdir-relative|subdir-relative-verified|subdir-simple|subdir-control|backslash-relative|absolute-z|absolute-c|absolute-c-verified|missing-parent' >&2; exit 2 ;;
esac

work=/work/stage5
probe=/work/stage16
xpid=
backup_input=/tmp/vtspeak-lipsync-filename-input1.txt
backup_output=/tmp/vtspeak-lipsync-filename-output.wav
log="$probe/lipsync-filename-bytewise-$case_name.log"
capture="$probe/lipsync-filename-bytewise-$case_name-output.txt"
case "$case_name" in
  extension) expected_output="$work/vtspeak-lipname.txt" ;;
  dot-relative) expected_output="$work/vtspeak-lipdot" ;;
  subdir-*) expected_output="$work/lipsync-probe/${name##*/}" ;;
  backslash-relative) expected_output="$work/lipsync-probe/backslash.txt" ;;
  missing-parent) expected_output="$work/$name" ;;
  absolute-z) expected_output="$probe/lipsync-filename-bytewise-absolute-z-output.txt" ;;
  absolute-c) expected_output=/work/stage2-copy/wineprefix/drive_c/windows/temp/vtspeak-lipsync-c-drive-output.txt ;;
  absolute-c-verified) expected_output=/work/stage2-copy/wineprefix/drive_c/windows/temp/vtspeak-lipsync-c-drive-output-verified.txt ;;
esac
shopt -s nullglob
existing_default=("$work"/length-sync-*)
if [ "${#existing_default[@]}" -ne 0 ]; then
  echo 'refusing to run with an existing length-sync report' >&2
  exit 1
fi
if [ -e "$log" ] || [ -e "$capture" ]; then
  echo "refusing to overwrite filename-case captures for $case_name" >&2
  exit 1
fi
if [ -e "$expected_output" ]; then
  echo "refusing to overwrite $expected_output" >&2
  exit 1
fi
if [ -e "$work/lipsync-probe" ] && [[ "$case_name" == subdir-* || "$case_name" == backslash-relative ]]; then
  echo "refusing to use existing scratch directory $work/lipsync-probe" >&2
  exit 1
fi
if [ -e "$work/no-such-lipsync-dir" ] && [[ "$case_name" == missing-parent ]]; then
  echo "refusing to use existing scratch directory $work/no-such-lipsync-dir" >&2
  exit 1
fi

restore() {
  cp "$backup_input" "$work/input1.txt"
  cp "$backup_output" "$work/output.wav"
  cmp -s "$backup_input" "$work/input1.txt"
  cmp -s "$backup_output" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
  if [[ "$case_name" == subdir-* || "$case_name" == backslash-relative ]]; then rmdir "$work/lipsync-probe" 2>/dev/null || true; fi
}

cp "$work/input1.txt" "$backup_input"
cp "$work/output.wav" "$backup_output"
trap restore EXIT
cp "$probe/input.txt" "$work/input1.txt"
if [[ "$case_name" == subdir-* || "$case_name" == backslash-relative ]]; then mkdir "$work/lipsync-probe"; fi
trace="$probe/trace-lipsync-filename-$case_name.gdb"
if [ -e "$trace" ]; then
  echo "refusing to overwrite $trace" >&2
  exit 1
fi
while IFS= read -r line || [ -n "$line" ]; do
  if [[ "$line" == *'@INITIALIZE_NAME@'* ]]; then
    for ((index = 0; index < ${#name}; index++)); do
      character=${name:index:1}
      printf -v byte_value '%d' "'$character"
      printf '  set {char}($name + %d) = %d\n' "$index" "$byte_value"
    done
    printf '  set {char}($name + %d) = 0\n' "${#name}"
  else
    line=${line//@CASE@/$case_name}
    line=${line//@NAME@/$name}
    printf '%s\n' "$line"
  fi
done < "$probe/trace-lipsync-filename.gdb.in" > "$trace"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-lipsync-filename-$case_name.log" 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99
if WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe < "$trace" > "$log" 2>&1; then
  process_status=0
else
  process_status=$?
fi
if ! grep -Fq "LIPSYNC_FILENAME_CALL case=$case_name" "$log"; then
  echo "GDB did not reach the filename call for $case_name (process status $process_status)" >&2
  exit 1
fi
printf 'LIPSYNC_FILENAME_PROCESS_STATUS %s\n' "$process_status" >> "$log"
found=()
for output in "${existing_default[@]}"; do
  if [ -f "$output" ]; then found+=("$output"); fi
done
if [ -f "$expected_output" ] && [[ "$case_name" != subdir-* && "$case_name" != backslash-relative ]]; then found+=("$expected_output"); fi
literal_backslash_output="$work/$name"
if [[ "$case_name" == backslash-relative && -f "$literal_backslash_output" ]]; then found+=("$literal_backslash_output"); fi
if [ -d "$work/lipsync-probe" ]; then
  subdir_outputs=("$work"/lipsync-probe/*)
  for output in "${subdir_outputs[@]}"; do
    if [ -f "$output" ]; then found+=("$output"); fi
  done
fi
if [ "${#found[@]}" -eq 1 ]; then
  printf 'LIPSYNC_FILENAME_OUTPUT %s\n' "${found[0]#"$work"/}" >> "$log"
  if [ "${found[0]}" != "$capture" ]; then mv "${found[0]}" "$capture"; fi
elif [ "${#found[@]}" -eq 0 ]; then
  printf 'LIPSYNC_FILENAME_OUTPUT none\n' >> "$log"
else
  printf 'LIPSYNC_FILENAME_OUTPUT unexpected-count=%s\n' "${#found[@]}" >> "$log"
  exit 1
fi
