#!/bin/bash
set -euo pipefail

probe=/work/stage21
case_id=${1:?case required: name-exact, name-short, path-exact, or path-short}
run_id=${2:-v1}
case "$case_id" in
  name-exact) expected='SPEAKERSINFO_GUARD_RESULT case=name-exact selector=5 result=6' ;;
  name-short) expected='SPEAKERSINFO_GUARD_RESULT case=name-short selector=5 result=6' ;;
  path-exact) expected='SPEAKERSINFO_GUARD_RESULT case=path-exact selector=3 result=6' ;;
  path-short) expected='SPEAKERSINFO_GUARD_RESULT case=path-short selector=3 result=6' ;;
  *) echo "unknown guard case: $case_id" >&2; exit 2 ;;
esac

trace="$probe/trace-speakersinfo-guard-$case_id-v1.gdb"
log="$probe/speakersinfo-guard-$case_id-$run_id-api.log"
xvfb_log="$probe/xvfb-speakersinfo-guard-$case_id-$run_id.log"
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
  echo "SPEAKERSINFO_GUARD_RESULT case=$case_id outcome=returned"
elif grep -Fq 'SPEAKERSINFO_GUARD_SETUP_FAILED' "$log"; then
  echo "SPEAKERSINFO_GUARD_RESULT case=$case_id outcome=setup-failed" >&2
  exit 1
elif grep -Eq 'exited with code 030000000005|exit process \(3221225477\)' "$log"; then
  echo "SPEAKERSINFO_GUARD_RESULT case=$case_id outcome=fault"
else
  echo "SPEAKERSINFO_GUARD_RESULT case=$case_id outcome=unclassified status=${command_status:-0}" >&2
  exit 1
fi
