# Lead 4: abbreviation context probes (2026-09-25)

## Full inventory matrix

The completed matrix covers every inventory key: 55 H keys, 42 T keys, and
330 C keys. Each key has a baseline observation plus context, punctuation, and
case contrasts. The H keys have four observations each, T keys have five, and
C keys have four, for 1,750 delimited observations in total. The Stage 10 GDB
trace captured complete ordered returned rows and 24 phone/context bytes per
row. The capture parser checked normal process exits, all delimited samples,
and exact per-set inventory coverage. It mapped 3,482 original-engine context
returns from 55 Wine processes.

The matrix and raw evidence are under
[`tools/revkit/work/stage15`](../../tools/revkit/work/stage15/). The [JSON
matrix](../../tools/revkit/work/stage15/lead4-key-matrix.json) groups all
observations by key and retains every row. The [TSV
matrix](../../tools/revkit/work/stage15/lead4-key-summary.tsv) exposes one row
per key/context with the input, focused returned surface span, and complete
ordered surface sequence. The generator also records when an alternate-case
probe equals another key in the same inventory set (21 C keys and 5 T keys),
so those observations are not mistaken for isolated case-only changes.

The context selector uses the inventory family and spelling to choose likely
contexts. H probes compare title, ordinary-sentence, punctuation, and case
forms. T probes compare company suffix, comma, missing-period, ordinary, and
case forms. C probes compare the neutral baseline, one context selected by a
spelling heuristic, comma punctuation, and alternate case. The C selector
used state-code contexts for 50 keys, calendar-date for 24, weekday for 20,
clock/zone for 14, address/building for 13, and a numeric phrase for the other
209. Those selector categories are probe choices, not recovered VoiceText
labels. A separate `Next.` sentence delimits each sample in the batched text;
the parser retains it in the captured rows and uses it to associate outputs
with manifest entries.

This closes the all-key pass at the documented contexts. It does not test
every possible punctuation combination or neighboring phrase for every key,
and it does not assign semantic names to raw phone/context codes. The
observations and per-key surfaces in the Stage 15 matrix are the full result;
interpretations remain bounded by those inputs.

## Initial selected-key probes

The original 2013 M16 Paul executable and DLL were run under Wine for 48
independent inputs covering 16 selected H/T/C abbreviation keys. Each input
ran in a fresh process. The Stage 10 GDB trace captured the upstream token
surfaces and each returned `0x70`-byte phone/context row, including its ordered
surface, class/flag fields, and 24 bytes beginning at the row's code field.
The complete captures and parsed rows are under
[`tools/revkit/work/stage14`](../../tools/revkit/work/stage14/); the exact
inputs and selection are in the [pilot manifest](../../tools/revkit/work/stage14/lead4-probe-manifest.json).

The matrix was selected for plausible title, company, date, time, measurement,
address, and state-code uses. It extends the 40 earlier targeted contrasts;
it is not a complete context matrix for all 427 H/T/C keys.

| Inventory key | Context outcomes observed |
| --- | --- |
| H `jr` | After a name, `Jr.` and `Jr` returned `Junior`; before `Smith`, `Jr.` returned `junior`. |
| H `sr` | After a name, `Sr.` and `Sr` returned `Senior`; before `Smith`, `Sr.` returned `senior`. |
| H `gen` | `Gen. Lee` returned `general`; lowercase `gen. Lee` retained `gen`. |
| H `rep` | `Rep. Jones` returned `representative`; lowercase `rep. Jones` retained `rep`. |
| T `co` | In `Acme Co.`, `Co.` returned `Company`; lowercase `co` and standalone `co.` stayed literal. |
| T `inc` | `Inc.` returned `incorporated` both with and without a comma after the company name, and in the tested standalone phrase. |
| T `LLC` | `LLC`, `llc`, and `Llc` stayed literal in the three company-name controls. |
| T `AG` | `AG` stayed literal before a company name and as a standalone token; the dotted `A.G.` form stayed `A.G`. |
| C `am` | After a number, `am` and `AM` returned `AM`; the verb in `I am ready` stayed `am`. |
| C `in` | After a number, `in.` returned `inches`; as a preposition or sentence-initial word it stayed `in`/`In`. |
| C `mi` | After a number, `mi.` returned `miles`; `MI` and `mi` in the state-name control stayed literal. |
| C `ct` | After a street name and number, `Ct.` returned `Court`; `CT` in the state-code phrase and standalone `ct` stayed literal. |
| C `CA` | `CA` and mixed-case `Ca` stayed literal in the tested state-code phrases; lowercase `ca` returned `circa`. |
| C `IN` | `IN` in the state-code phrase stayed literal. The preposition and sentence-initial controls stayed `in` and `In`. |
| C `CO` | `CO` in the state-code phrase stayed literal; the company suffix `Co.` and standalone `Co.` returned `Company`. |
| C `Jan` | The date, month-name, and person-name controls all returned `January`; the date also returned `fifth` for `5`. |

These observations establish phrase-conditioned returned surfaces for this
selected set. They do not name the raw row flags or phone-code bytes, infer
general rules from a single contrast, or establish behavior for every key,
punctuation form, or neighboring token. In particular, the `ca` result is a
single lowercase control and does not establish a general case-folding rule.

## Reproduce the selected-key pilot

Generate the 48 fixtures from the checked-in probe definitions, then run the
fresh-process Wine/GDB captures and parse every returned row:

```sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/scripts/generate_lead4_probes.py
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage14/run-lead4.sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/scripts/summarize_lead4.py
```

The runner uses the existing Stage 10 trace and restores its Stage 5 input and
output from a dedicated Stage 14 scratch path when it exits. It uses the
existing isolated Wine setup and read-only vendor mounts.

## Reproduce the full inventory matrix

```sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/scripts/generate_lead4_matrix.py
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage15/run-lead4-matrix.sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/scripts/summarize_lead4_matrix.py
```

To rerun just one inventory family, pass `H`, `T`, or `C` to the Stage 15
runner. Its dedicated scratch path restores Stage 5 input/output on exit.
