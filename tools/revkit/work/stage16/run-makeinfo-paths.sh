#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
outdir=/tmp/vtspeak-makeinfo-paths-run
xpid=
if [ -e "$outdir" ]; then
  echo "refusing to overwrite existing MakeInfo path-probe directory" >&2
  exit 1
fi
mkdir "$outdir"
restore() {
  cp "$probe/makeinfo-original-input1.txt" "$work/input1.txt"
  cp "$probe/makeinfo-original-output.wav" "$work/output.wav"
  cmp -s "$probe/makeinfo-original-input1.txt" "$work/input1.txt"
  cmp -s "$probe/makeinfo-original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

# These checked-in files are the fixed restore baseline. Do not snapshot the
# mounted Stage 5 paths here: the runtime container may have initialized them
# with a probe fixture before this script starts.
trap restore EXIT
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-makeinfo-paths.log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-makeinfo-paths.gdb" > "$probe/makeinfo-paths-api.log" 2>&1
find "$outdir" /work/stage16 /samples -type f -name '*.dtt' -printf '%p %s bytes\n' \
  > "$probe/makeinfo-paths-files.txt"
