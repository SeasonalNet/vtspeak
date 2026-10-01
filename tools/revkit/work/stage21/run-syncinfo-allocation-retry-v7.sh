#!/bin/bash
set -euo pipefail

repo=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../.." && pwd)
bash "$repo/tools/revkit/work/stage21/build-syncinfo-allocation-retry.sh"
docker compose -f "$repo/tools/revkit/work/stage8/compose.yaml" run --rm runtime \
  /bin/bash /work/stage21/run-syncinfo-allocation-retry-v7-container.sh
