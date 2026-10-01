#!/bin/bash
set -euo pipefail
work=/work/stage5
probe=/work/stage21
run_id=${1:-matrix1}
group=${2:-all}
[[ "$run_id" =~ ^[A-Za-z0-9_-]+$ ]] || exit 2
names=(null_array_cap1 null_first_cap1 null_first_cap2 bad_first_cap1
  bad_first_cap2 late_null_cap4 late_null_cap5 late_null_cap6
  empty_then_null empty_huge_count negative_null_array)
guards=(guard_array_cap4 guard_array_cap5 guard_array_cap6 guard_empty_huge)
case "$group" in
  all) names+=("${guards[@]}") ;;
  guards) names=("${guards[@]}") ;;
  *) echo "group must be all or guards" >&2; exit 2 ;;
esac
input_backup="/tmp/csv-pointer-order-$run_id-$group-input.txt"
output_backup="/tmp/csv-pointer-order-$run_id-$group-output.wav"
xvfb_log="$probe/xvfb-csv-pointer-order-$run_id-$group.log"
xpid=
restore() {
  cp "$input_backup" "$work/input1.txt"
  cp "$output_backup" "$work/output.wav"
  cmp -s "$input_backup" "$work/input1.txt"
  cmp -s "$output_backup" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
for path in "$input_backup" "$output_backup" "$xvfb_log"; do
  if [ -e "$path" ]; then echo "refusing to overwrite $path" >&2; exit 3; fi
done
for name in "${names[@]}"; do
  log="$probe/csv-pointer-order-$run_id-$name-api.log"
  if [ -e "$log" ]; then echo "refusing to overwrite $log" >&2; exit 3; fi
done
cp "$work/input1.txt" "$input_backup"
cp "$work/output.wav" "$output_backup"
trap restore EXIT
cd "$work"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
for name in "${names[@]}"; do
  log="$probe/csv-pointer-order-$run_id-$name-api.log"
  if WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 90s \
    winedbg --gdb /samples/voicetext_paul.exe \
    < "$probe/trace-csv-pointer-order-$name.gdb" > "$log" 2>&1; then
    rc=0
  else
    rc=$?
  fi
  printf 'DEBUGGER_EXIT code=%s\n' "$rc" >> "$log"
  grep -E '^CSV_PTR_(RETURN|FAULT)' "$log"
done
