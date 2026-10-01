#!/bin/bash
set -euo pipefail

repo=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../.." && pwd)
docker run --rm --network none --user "$(id -u):$(id -g)" \
  -v "$repo:/src:ro" \
  -v "$repo/tools/revkit/work/stage21:/out" \
  vtspeak-pe32-builder:local \
  i686-w64-mingw32-gcc -O2 -Wall -Wextra -Werror -std=c11 \
  /src/tools/revkit/work/stage21/probe-license-md5.c \
  -o /out/probe-license-md5.exe
docker compose -f "$repo/tools/revkit/work/stage8/compose.yaml" run --rm \
  runtime /bin/bash /work/stage21/run-license-md5-container.sh
