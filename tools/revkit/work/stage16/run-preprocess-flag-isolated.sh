#!/bin/bash
set -euo pipefail

flag=${1:?usage: run-preprocess-flag-isolated.sh BYTE_FLAG}
if ! [[ "$flag" =~ ^([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$ ]]; then
  echo "flag must be an integer from 1 through 255" >&2
  exit 2
fi

probe=/work/stage16
work=/work/stage5
trace="$probe/trace-preprocess-flag${flag}-isolated.gdb"
log="$probe/preprocess-flag${flag}-isolated-api.log"
listing="$probe/preprocess-flag${flag}-isolated-files.txt"
xvfb_log="$probe/xvfb-preprocess-flag${flag}-isolated.log"
wine_debug=${VT_PREPROCESS_WINEDEBUG:--all}
before=$(mktemp)
after=$(mktemp)
xpid=
restore() {
  cp "$probe/preprocess-flags-original-input1.txt" "$work/input1.txt"
  cp "$probe/preprocess-flags-original-output.wav" "$work/output.wav"
  cmp -s "$probe/preprocess-flags-original-input1.txt" "$work/input1.txt"
  cmp -s "$probe/preprocess-flags-original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
  rm -f "$before" "$after"
}
trap restore EXIT

for output in "$trace" "$log" "$listing" "$xvfb_log"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done

cat > "$trace" <<GDB
set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set \$text = *(unsigned int *)(\$esp + 8)
  disable 1
  set \$raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)\$text, "Z:/work/stage16/flag${flag}-preprocess-isolated-20260926", ${flag}, 1, -1, -1, -1, -1, 0, 0)
  printf "PREPROCESS_FLAG_ISOLATED flag=${flag} raw_eax=%#x low_ax=%d\\n", \$raw, ((short)\$raw)
  continue
end

continue
GDB

cp "$probe/input.txt" "$work/input1.txt"
find "$probe" -maxdepth 1 -type f -printf '%f %s %T@\n' | sort > "$before"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG="$wine_debug" timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" \
  > "$log" 2>&1
find "$probe" -maxdepth 1 -type f -printf '%f %s %T@\n' | sort > "$after"
comm -13 "$before" "$after" > "$listing"
