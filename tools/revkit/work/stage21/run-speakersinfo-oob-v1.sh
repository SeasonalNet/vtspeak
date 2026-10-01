#!/bin/bash
set -euo pipefail

probe=/work/stage21
case_id=${1:?selector case required: minus2, minus1, 6, 7, or intmax}
run_id=${2:-v1}
case "$case_id" in
  minus2) expected='SPEAKERSINFO_OOB selector=-2' ;;
  minus1) expected='SPEAKERSINFO_OOB selector=-1' ;;
  6) expected='SPEAKERSINFO_OOB selector=6' ;;
  7) expected='SPEAKERSINFO_OOB selector=7' ;;
  intmax) expected='SPEAKERSINFO_OOB selector=2147483647' ;;
  *) echo "unknown selector case: $case_id" >&2; exit 2 ;;
esac

trace="$probe/trace-speakersinfo-oob-$case_id-v1.gdb"
log="$probe/speakersinfo-oob-$case_id-$run_id-api.log"
xvfb_log="$probe/xvfb-speakersinfo-oob-$case_id-$run_id.log"
xpid=
restore() {
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT

for output in "$log" "$xvfb_log"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" > "$log" 2>&1 || command_status=$?

if grep -Fq "$expected" "$log"; then
  echo "SPEAKERSINFO_OOB_RESULT selector_case=$case_id outcome=returned"
elif grep -Eq 'exited with code 030000000005|exit process \(3221225477\)' "$log"; then
  echo "SPEAKERSINFO_OOB_RESULT selector_case=$case_id outcome=fault"
else
  echo "SPEAKERSINFO_OOB_RESULT selector_case=$case_id outcome=unclassified status=${command_status:-0}" >&2
  exit 1
fi
