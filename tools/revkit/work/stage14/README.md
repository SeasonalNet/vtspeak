# Stage 14: Lead 4 abbreviation context probes

This stage expands Lead 4 of the follow-up investigation. It contains 48
fresh-process probes for 16 selected keys across the H/T/C abbreviation sets.
`lead4-probe-manifest.json` maps each fixture to its key, context, and exact
input. `inputs/` contains the generated one-phrase fixtures; the GDB logs
capture the upstream token records and ordered returned phone/context rows.
`lead4-results.json` contains each returned row's fields, surface, and 24 code
bytes.

Reproduce from the repository root:

```sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/scripts/generate_lead4_probes.py
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage14/run-lead4.sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/scripts/summarize_lead4.py
```

The runner backs up and restores `stage5/input1.txt` and `stage5/output.wav`
using `/work/corpus-parity/lead4-stage14/`. Its Wine prefix and vendor inputs
are reused from the existing isolated setup; the container keeps the same
read-only mounts and dropped capabilities.

The curated observations and limits are in
[`lead4-abbreviation-context-results-2026-09-25.md`](../../../../docs/reverse-engineering/lead4-abbreviation-context-results-2026-09-25.md).
