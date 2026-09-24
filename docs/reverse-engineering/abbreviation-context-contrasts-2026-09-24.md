# Abbreviation context contrasts in the local Paul engine (2026-09-24)

These are controlled observations from the original 2013 M16 Paul executable
and DLL under Wine. Each [Stage 10 input](../../tools/revkit/work/stage10/README.md)
ran in a fresh process. A GDB breakpoint at `FUN_10007520` (`0x10007520`)
captured its upstream `0x94`-byte token surfaces (`+0x34`) and the returned
`0x70`-byte phone/context rows (surface `+0x07`, code `+0x25`). Those offsets
come from the [Stage 6 layout analysis](voice-engine-and-model-formats.md#stage-6-text-and-pronunciation-resources);
the function names and signatures are decompiler pseudocode. The
[capture verifier](../../tools/revkit/work/scripts/verify_stage10.py) confirms
all 40 complete ordered surface sequences and normal process exits.

| Controlled input | Returned surfaces, in order |
| --- | --- |
| `Dr. Smith spoke.` | `doctor / Smith / spoke / .` |
| `The doctor spoke.` | `The / doctor / spoke / .` |
| `42 Main Dr.` | `forty / two / Main / Drive / .` |
| `St. Paul spoke.` | `Saint / Paul / spoke / .` |
| `42 Main St.` | `forty / two / Main / Street / .` |
| `Apt 4.` and `Apt. 4.` in separate runs | Both `Apartment / four / .` |
| `Apt 4 is ready.` and `Apt. 4 is ready.` | Both `Apartment / four / is / ready / .` |
| `Apt. 4 and Apt 4.` | `Apartment / four / and / Apt / four / .` |
| `Apt 4 and Apt 4.` | `Apartment / four / and / Apt / four / .` |
| `Apt 4 and Apt 5.` | `Apartment / four / and / Apt / five / .` |
| `Apt 4. Apt 4.` | `Apartment / four / Apartment / four / .` |
| `Unit 4 and Apt 4.` | `Unit / four / and / Apartment / four / .` |
| `Hello and Apt 4.` | `Hello / and / Apartment / four / .` |
| `Apt 4 and Hello and Apt 4.` | Both `Apt` occurrences expand to `Apartment` |
| `Apt 4, Apt 4.` | Both `Apt` occurrences expand to `Apartment` |
| `Unit 4, Apt 4.` | `Unit / four / , / Apartment / four / .` |
| `Apt 4 and Unit 4.` | `Apartment / four / and / Unit / four / .` |
| `It is 5 p.m. today.` and `It is 5 pm. today.` | Both `It / is / five / PM / today / .` |
| `Acme Inc. hired Lee.` and `Acme, Inc. hired Lee.` | Both expand `Inc` to `incorporated` |
| `Acme LLC hired Lee.` | `LLC` remains `LLC` in this context |
| `Acme Co. hired Lee.` | `Co.` expands to `Company` |
| `No. 4 arrived.` | `No.` expands to `Number` |
| `No, thanks.` | `No` remains `No` |
| `Mr. Smith spoke.` / `Mrs. Smith spoke.` | `mister` / `misses`, respectively |
| `Prof.`, `Capt.`, and `Rev.` before `Smith` | `professor`, `captain`, and `reverend` |
| `Acme Corp. hired Lee.` / `Acme Ltd. hired Lee.` | `Corp.` → `corporation`; `Ltd.` → `limited` |
| `Acme GmbH hired Lee.` | `GmbH` remains `GmbH` in this context |
| `We met on Jan. 5.` / `It happened in Jan.` | Both produce `January` (the first also expands `5` to `fifth`) |
| `It is 10 ft. long.` / `The road is 2 mi. long.` | `ft.` → `feet`; `mi.` → `miles` |
| `The apt. is ready.` | `apt.` remains `apt` in this non-address context |

The `Dr.` and `St.` pairs demonstrate different expansions in name and address
contexts. `Mr.`, `Mrs.`, `Prof.`, `Capt.`, and `Rev.` expand before the tested
proper name. Both `p.m.` and `pm.` produce `PM` in the tested time context.
The T-set company suffix contrasts show that `Inc.`, `Co.`, `Corp.`, and
`Ltd.` expand in the tested company phrase, while `LLC` and `GmbH` remain
literal; the comma before `Inc.` does not change its expansion. The C-set
contrasts show `Jan.` before both a date and an ordinary following phrase,
`ft.`/`mi.` as measurements, and `apt.` retained outside the tested address
phrase. `No.` before a number expands to `Number`, while comma punctuation
does not.
Period presence alone does **not** explain the tested `Apt` difference:
isolated `Apt 4` expands with or without one. Repeating `Apt 4` in the same
sentence leaves the second `Apt` unexpanded, even when its number changes to
`5`. A sentence boundary, comma, or intervening word allows both occurrences
to expand. This narrows the behavior to the tested uninterrupted `Apt number
and Apt number` phrase; it is not a general rule for all abbreviations.

The repeated case already differs at the upstream token surface on entry to
`FUN_10007520`: the first surface is `Apartment`, the second is `Apt`. The
returned row for the second `Apt` also has class byte `0x41`, while the
expanded `Apartment` row has zero there in these captures. The contrast is
therefore established **before** this phone/context rule pass; assigning its
cause to `FUN_10007520` would be incorrect. A further breakpoint on the
address-alias lookup helper `FUN_10056150` shows one `Apt` match (index `0`
at caller `0x10043d75`) in each of `Apt. 4 and Apt 4.` and
`Unit 4 and Apt 4.`. A branch probe at `FUN_10043810` (`0x10043d52` and
`0x10044003`) identifies one alias branch for `Apt 4 and Apt 4.` and
`Apt 4 and Apt 5.`, but two branches for each of `Apt 4. Apt 4.`,
`Apt 4, Apt 4.`, and `Apt 4 and Hello and Apt 4.`. In the repeated
conjunction case, only the first `Apt` reaches the branch; the unexpanded
second one does not. The different-number control shows that the behavior is
not caused by matching equal numbers. The separated controls show that
repetition by itself is not sufficient. The address-phrase span mechanism is
resolved below; the meaning of row class `0x41` is not labeled by the local
corpus.
An additional breakpoint at `0x10043b61`, just before the address-rule guard
in `FUN_10043810`, sees one `Apt` key in each of the repeated and `Unit`
controls. In these captures the observed guard state for that key is
token index `0`, pending flag `0`, and kind `1`. A call/return trace at
`FUN_10043810` resolves the exclusion: for `Apt 4 and Apt 4.`, the call
receives ` Apt 4 and Apt 4.\n` and returns `14`, consuming exactly the prefix
` Apt 4 and Apt`; the next call receives ` 4.\n` and returns `-1`. The second
`Apt` is consumed as part of the matched phrase span, so it never receives an
independent address-alias lookup. For `Apt 4 and Hello and Apt 4.`, the matcher
consumes ` Apt 4 and Hello` (16 bytes), then later receives ` Apt 4.\n` and
returns `6` for the second address phrase. In the comma control, the first
`Apt 4` returns `6`, the comma is rejected, and the second `Apt 4` returns `6`
independently. The raw call/return captures are named `*-address-call.log`,
and the reproducible GDB script is
[`trace-address-call.gdb`](../../tools/revkit/work/stage10/trace-address-call.gdb).
This identifies the operational selector as the address phrase matcher's
consumed span. The internal table remains represented by raw rows, while the
accepted span and its branch behavior are established directly.
The [verifier](../../tools/revkit/work/scripts/verify_stage10.py) checks both
the guard values and the alias-branch visits.
These rows also do not label every rule flag or phone code byte.

The table is a targeted matrix across H/T/C abbreviation classes, not a
complete map of their outputs. It does not enumerate all 427 H/T/C keys or
every placement and punctuation combination. The `Apt` contrasts isolate
one phrase-boundary effect. The H/T/C resources only encode membership and case matching;
they do not provide replacement words. The [open-topics review](open-topics-review-2026-09-24.md)
retains broader linguistic-label and abbreviation coverage as open work.

## Complete H/T key sweeps in one context each

The [inventory-driven generator](../../tools/revkit/scripts/generate_abbreviation_probes.py)
creates 11 short batches containing every key from the 55-row H set and
42-row T set. Each H key is placed before `Smith spoke.`; each T key is placed
in `Acme KEY. hired Lee.`. The original engine returned normally for all
batches, and the [Stage 10 verifier](../../tools/revkit/work/scripts/verify_stage10.py)
checks the declared, ordered rows for all 87 context-return calls. These runs
show the key sets are exercised at runtime under one controlled template per
key; they are not a per-key expansion table, because many keys remain literal
or require different contexts. Additional contexts for all three sets remain
untested as a complete matrix.

## C-set baseline: every key in one neutral context

The 330 C-set keys were each placed in `We noted KEY today.`. Thirty-three
fresh-process batches produced 330 `FUN_10007520` context returns, and every
process exited normally. The [generator](../../tools/revkit/scripts/generate_abbrc_probes.py),
[ordered key manifest](../../tools/revkit/work/stage10/abbrc-probe-manifest.json),
[raw Wine/GDB logs](../../tools/revkit/work/stage10/), and
[summarizer/verifier](../../tools/revkit/work/scripts/summarize_abbrc_probes.py)
make the run reproducible. All full returned surface sequences are recorded
in [`abbrc-probe-results.json`](../../tools/revkit/work/stage10/abbrc-probe-results.json).

In this context, 237 returned the literal five-surface sequence
`We / noted / KEY / today / .`. Ninety-three differed. The observed
differences include lexical expansions such as `acctg` → `accounting`,
`approx` → `approximately`, `appt` → `appointment`, and `Jan` → `January`,
as well as capitalization changes such as `cc` → `CC` and `cgs` → `CGS`.
Compound pronunciations can return multiple surface rows: `kWh` produced
`kilowatt / hour`, `lbf` produced `pound / force`, and `Nm` produced
`newton / meter`. The result artifact retains each key and full row sequence,
including the other observed outputs.

These are outcomes for one deliberately uniform surrounding phrase. A literal
key result means only that this template did not change its returned surface.
It does not show that the entry is inactive or lacks a specialized context.
Conversely, a returned expansion in the neutral phrase does not establish how
that key behaves in a date, time, measurement, address, or case-sensitive
context. The sweep closes the missing one-context inventory for the C set;
discriminating per-key context rules remain open.
