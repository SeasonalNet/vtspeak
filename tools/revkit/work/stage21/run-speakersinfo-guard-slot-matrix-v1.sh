#!/bin/bash
set -euo pipefail

probe=/work/stage21
run_id=${1:-v1}
if [[ ! "$run_id" =~ ^[A-Za-z0-9_-]+$ ]]; then
  echo "run id may contain only letters, digits, underscores, and hyphens" >&2
  exit 2
fi
template="$probe/trace-speakersinfo-guard-matrix-template-v1.gdb"
xvfb_log="$probe/xvfb-speakersinfo-guard-slot-matrix-$run_id.log"
xpid=
temp_traces=()
restore() {
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
  for trace in "${temp_traces[@]}"; do rm -f "$trace"; done
}
trap restore EXIT

cases=(
  slot0-exact:0:5:21:returned
  slot0-name-short:0:4:21:fault
  slot0-path-short:0:5:20:fault
  slot1-exact:1:5:21:returned
  slot1-name-short:1:4:21:fault
  slot1-path-short:1:5:20:fault
  slot2-exact:2:6:21:returned
  slot2-name-short:2:5:21:fault
  slot2-path-short:2:6:20:fault
  slot3-exact:3:6:24:returned
  slot3-name-short:3:5:24:fault
  slot3-path-short:3:6:23:fault
  slot4-exact:4:6:19:returned
  slot4-name-short:4:5:19:fault
  slot4-path-short:4:6:18:fault
  slot5-exact:5:7:21:returned
  slot5-name-short:5:6:21:fault
  slot5-path-short:5:7:20:fault
)

for record in "${cases[@]}"; do
  IFS=: read -r case_id slot name_capacity path_capacity outcome <<< "$record"
  log="$probe/speakersinfo-guard-$case_id-$run_id-api.log"
  if [ -e "$log" ]; then
    echo "refusing to overwrite $log" >&2
    exit 3
  fi
done
if [ -e "$xvfb_log" ]; then
  echo "refusing to overwrite $xvfb_log" >&2
  exit 3
fi

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1

for record in "${cases[@]}"; do
  IFS=: read -r case_id slot name_capacity path_capacity outcome <<< "$record"
  trace="/tmp/vtspeak-speakersinfo-guard-$case_id-$run_id.gdb"
  log="$probe/speakersinfo-guard-$case_id-$run_id-api.log"
  temp_traces+=("$trace")
  sed \
    -e "s/@CASE@/$case_id/g" \
    -e "s/@SLOT@/$slot/g" \
    -e "s/@NAMECAP@/$name_capacity/g" \
    -e "s/@PATHCAP@/$path_capacity/g" \
    "$template" > "$trace"

  DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
    winedbg --gdb /samples/voicetext_paul.exe \
    < "$trace" > "$log" 2>&1 || command_status=$?

  if grep -Fq 'SPEAKERSINFO_GUARD_MATRIX_SETUP_FAILED' "$log"; then
    echo "guard setup failed for $case_id" >&2
    exit 1
  fi
  grep -Eq 'SPEAKERSINFO_GUARD_MATRIX_SETUP case=.* name=0x[0-9a-f]+/1 path=0x[0-9a-f]+/1 old_protect=0x4' "$log"
  if [ "$outcome" = returned ]; then
    grep -Fq "SPEAKERSINFO_GUARD_MATRIX case=$case_id slot=$slot name_capacity=$name_capacity path_capacity=$path_capacity result=6" "$log"
  else
    grep -Eq 'exited with code 030000000005|exit process \(3221225477\)' "$log"
  fi
  echo "SPEAKERSINFO_GUARD_MATRIX_RESULT case=$case_id outcome=$outcome"
done
