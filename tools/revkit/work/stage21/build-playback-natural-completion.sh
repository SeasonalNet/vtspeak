#!/bin/bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../.." && pwd)
stage="$repo_root/tools/revkit/work/stage21"

docker run --rm --network none --user "$(id -u):$(id -g)" \
  -v "$repo_root:/src:ro" -v "$stage:/out" \
  vtspeak-pe32-builder:local \
  i686-w64-mingw32-gcc -O2 -Wall -Wextra -Werror -std=c11 \
  -Wl,--no-insert-timestamp \
  /src/tools/revkit/work/stage21/probe-playback-natural-completion.c \
  -o /out/probe-playback-natural-completion.exe -luser32
