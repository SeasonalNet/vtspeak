#!/bin/bash
set -euo pipefail

case_name=${1:-}
case "$case_name" in
  baseline) pitch=-1; volume=-1; pause=-1; dictionary=-1; type=-1 ;;
  pitch-low) pitch=50; volume=-1; pause=-1; dictionary=-1; type=-1 ;;
  pitch-high) pitch=200; volume=-1; pause=-1; dictionary=-1; type=-1 ;;
  volume-zero) pitch=-1; volume=0; pause=-1; dictionary=-1; type=-1 ;;
  pause-250) pitch=-1; volume=-1; pause=250; dictionary=-1; type=-1 ;;
  dictionary-zero) pitch=-1; volume=-1; pause=-1; dictionary=0; type=-1 ;;
  text-type-4) pitch=-1; volume=-1; pause=-1; dictionary=-1; type=4 ;;
  text-type-0) pitch=-1; volume=-1; pause=-1; dictionary=-1; type=0 ;;
  text-type-1) pitch=-1; volume=-1; pause=-1; dictionary=-1; type=1 ;;
  text-type-2) pitch=-1; volume=-1; pause=-1; dictionary=-1; type=2 ;;
  text-type-3) pitch=-1; volume=-1; pause=-1; dictionary=-1; type=3 ;;
  text-type-5) pitch=-1; volume=-1; pause=-1; dictionary=-1; type=5 ;;
  text-type-6) pitch=-1; volume=-1; pause=-1; dictionary=-1; type=6 ;;
  text-type-7) pitch=-1; volume=-1; pause=-1; dictionary=-1; type=7 ;;
  extremes) pitch=50; volume=0; pause=250; dictionary=0; type=-1 ;;
  type4-extremes) pitch=50; volume=0; pause=250; dictionary=0; type=4 ;;
  type6-extremes) pitch=50; volume=0; pause=250; dictionary=0; type=6 ;;
  *) echo 'usage: run-lipsync-option.sh CASE (see case table in runner)' >&2; exit 2 ;;
esac

work=/work/stage5
probe=/work/stage16
xpid=
backup_input=/tmp/vtspeak-lipsync-option-input1.txt
backup_output=/tmp/vtspeak-lipsync-option-output.wav
name="vtspeak-lipoption-$case_name.txt"
output="$work/$name"
log="$probe/lipsync-option-$case_name.log"
capture="$probe/lipsync-option-$case_name-report.txt"

shopt -s nullglob
existing_default=("$work"/length-sync-*)
if [ "${#existing_default[@]}" -ne 0 ] || [ -e "$output" ] || [ -e "$log" ] || [ -e "$capture" ]; then
  echo "refusing to overwrite an existing lip-sync output or capture for $case_name" >&2
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
trace=/tmp/trace-lipsync-option-$case_name.gdb
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
    line=${line//@PITCH@/$pitch}
    line=${line//@VOLUME@/$volume}
    line=${line//@PAUSE@/$pause}
    line=${line//@DICTIONARY@/$dictionary}
    line=${line//@TYPE@/$type}
    printf '%s\n' "$line"
  fi
done < "$probe/trace-lipsync-option.gdb.in" > "$trace"

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-lipsync-option-$case_name.log" 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99
if WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe < "$trace" > "$log" 2>&1; then
  process_status=0
else
  process_status=$?
fi
if ! grep -Fq "LIPSYNC_OPTION_CALL case=$case_name" "$log"; then
  echo "GDB did not reach the option call for $case_name (process status $process_status)" >&2
  exit 1
fi
printf 'LIPSYNC_OPTION_PROCESS_STATUS %s\n' "$process_status" >> "$log"
if [ ! -f "$output" ]; then
  echo "expected report $output" >&2
  exit 1
fi
mv "$output" "$capture"
