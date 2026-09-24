#!/bin/bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../.." && pwd)
scratch="$repo_root/tools/revkit/work/corpus-parity"
mkdir -p "$scratch"

docker build -f "$repo_root/tools/revkit/Dockerfile.pe32-builder" \
  -t vtspeak-pe32-builder:local "$repo_root/tools/revkit"
docker run --rm --network none --user "$(id -u):$(id -g)" \
  -v "$repo_root:/src:ro" -v "$scratch:/out" \
  vtspeak-pe32-builder:local \
  i686-w64-mingw32-gcc -O2 -Wall -Wextra -Werror -std=c11 \
  -Wl,--no-insert-timestamp \
  /src/tools/revkit/scripts/probe_dll_decoder.c \
  -o /out/probe_dll_decoder.exe -ladvapi32
