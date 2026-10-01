#!/bin/bash
set -euo pipefail

probe=/work/stage21
xpid=
restore() {
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-version-export-mutation.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-version-export-mutation.gdb" > "$probe/version-export-mutation-api.log" 2>&1

grep -Fq 'VERSION_EXPORT_BASE first=3 second=11 third=7 fourth=1' "$probe/version-export-mutation-api.log"
grep -Fq 'VERSION_EXPORT_MUTATED first=203 second=204 third=205 fourth=206' "$probe/version-export-mutation-api.log"
grep -Fq 'DEFAULT_VERSION_WITH_MUTATED_DATA value=Paul-M16-FileIO' "$probe/version-export-mutation-api.log"
grep -Fq 'VERSION_EXPORT_RESTORED first=3 second=11 third=7 fourth=1' "$probe/version-export-mutation-api.log"
