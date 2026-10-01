#!/bin/bash
set -euo pipefail

repo=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../.." && pwd)
docker run --rm --network none --user "$(id -u):$(id -g)" \
  -v "$repo:/src:ro" \
  -v "$repo/tools/revkit/work/stage21:/out" \
  vtspeak-pe32-builder:local \
  i686-w64-mingw32-gcc -O2 -Wall -Wextra -Werror -std=c11 \
  /src/tools/revkit/work/stage21/license-token-byte-matrix.c \
  -o /out/license-token-byte-matrix.exe
docker compose -f "$repo/tools/revkit/work/stage8/compose.yaml" run --rm \
  runtime /bin/bash /work/stage21/run-license-token-byte-matrix-container.sh
python3 "$repo/tools/revkit/work/stage21/analyze_license_token_byte_matrix.py" \
  "$repo/tools/revkit/work/stage21/license-token-byte-matrix-v3-api.log"
