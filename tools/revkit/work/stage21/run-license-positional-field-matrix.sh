#!/bin/bash
set -euo pipefail

repo=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../.." && pwd)
exe="$repo/tools/revkit/work/stage21/license-positional-field-matrix.exe"
if [ -e "$exe" ]; then
  echo "refusing to overwrite $exe" >&2
  exit 2
fi

docker run --rm --network none --user "$(id -u):$(id -g)" \
  -v "$repo:/src:ro" \
  -v "$repo/tools/revkit/work/stage21:/out" \
  vtspeak-pe32-builder:local \
  i686-w64-mingw32-gcc -O2 -Wall -Wextra -Werror -std=c11 \
  /src/tools/revkit/work/stage21/license-positional-field-matrix.c \
  -o /out/license-positional-field-matrix.exe

docker compose -f "$repo/tools/revkit/work/stage8/compose.yaml" run --rm \
  runtime /bin/bash /work/stage21/run-license-positional-field-matrix-container.sh
