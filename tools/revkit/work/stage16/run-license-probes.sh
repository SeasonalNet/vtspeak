#!/bin/bash
set -euo pipefail

repo=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../.." && pwd)
docker compose -f "$repo/tools/revkit/work/stage8/compose.yaml" run --rm \
  runtime /bin/bash /work/stage16/run-license-probes-container.sh
