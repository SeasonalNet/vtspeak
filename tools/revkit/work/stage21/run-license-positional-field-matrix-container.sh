#!/bin/bash
set -euo pipefail

probe=/work/stage21
log="$probe/license-positional-field-matrix-api.log"
if [ -e "$log" ]; then
  echo "refusing to overwrite $log" >&2
  exit 2
fi

WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all \
  timeout 120s wine "$probe/license-positional-field-matrix.exe" > "$log" 2>&1
grep -q '^LICENSE_FIELD baseline=0 ' "$log"
for case_name in baseline_repeat expiry_position_only expdate_attribute_only \
  expiry_position_and_expdate vw_vtapi_token channel_position_only \
  channel_attribute_only channel_position_and_attribute; do
  grep -q "^LICENSE_FIELD case=$case_name result=" "$log"
done
