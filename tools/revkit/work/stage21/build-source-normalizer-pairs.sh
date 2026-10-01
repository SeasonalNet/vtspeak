#!/bin/bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../.." && pwd)
docker run --rm --network none --user "$(id -u):$(id -g)" \
  -v "$repo_root:/src:ro" \
  -v "$repo_root/tools/revkit/work/stage21:/out" \
  vtspeak-pe32-builder:local \
  i686-w64-mingw32-gcc -O2 -Wall -Wextra -Werror -std=c11 \
  -Wl,--no-insert-timestamp \
  /src/tools/revkit/work/stage21/probe-source-normalizer-pairs.c \
  -o /out/probe-source-normalizer-pairs.exe
