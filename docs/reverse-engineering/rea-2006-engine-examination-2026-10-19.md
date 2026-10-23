# REA examination of the 2006 VoiceText engine

## Target and scope

This static examination used REA's Ghidra 12.1.4 provider on the original
`tools/revkit/work/corpus-parity/kate-msi/Program Files/NeoSpeech/Kate16/lib/vt_eng.dll`.
SHA-256: `00fc9375d08bd8cec303845c992481d79c8f616d1bfb6390a5322cd06f69273d`.
REA identifies it as an x86 PE shared library. No target code was executed or
patched. The session was closed after exporting the
[complete evidence bundle](../../tools/revkit/work/rea-2006-2026-10-19/evidence.json).

Questions examined: API inventory, buffer dispatch and loading, legacy index
tail reads, and candidate continuity/scoring. These are bounded observations,
not a complete reconstruction or a new runtime parity result. All code below
describes decompiler pseudocode; recovered types and names may be imperfect.

## Observations

### Inventory and API

REA reports 1,028 procedures and 2,285 defined strings. Its `VT_` procedure
search returns 44 names, matching the prior PE export count, including file
and buffer synthesis, playback, model and user dictionary loading,
configuration, synchronization, information, and CSV helpers.

The [earlier static report](lead5-kate-common-resource-comparison-2026-09-26.md)
records 947 discovered functions for the same hash. Procedure counts depend
on analysis configuration and counting conventions; this examination did not
reconcile the two inventories. The difference does not establish newly
discovered functionality. Defined strings are not an exhaustive search of raw
artifact text.

`VT_TextToBuffer_ENG` at `0x10016740` dispatches its first argument:

| Value | Callee |
| --- | --- |
| 0 | `0x10018ba0` |
| 1 | `0x10018f60` |
| 2 | `0x100193e0` |
| 3 | `0x10019860` |
| other | returns `0xffffffff` |

These branches establish dispatch, not the semantic names of the modes.

`VT_LOADTTS_ENG` at `0x10021570` manages shared initialization and two indexed
contexts. It computes a context index clamped to 0 for arguments outside 0–1,
initializes a critical section and shared state, then calls `0x1001f8c0` and
`0x1001f950`. One earlier branch still indexes a flag with the original argument;
the clamp alone therefore does not establish safe handling of arbitrary inputs.
Helper behavior and runtime error contracts were not followed in this pass.

### Legacy index tail

`0x100123a0` skips `record_stride * N` bytes, where the stride is read from
model `+0x56` and N from `+0x50`, then invokes `0x10012030` in this order:

| Destination in model | Element width | Count |
| --- | --- | --- |
| `+0x4c` | 1 | N |
| `+0x48` | 1 | N |
| `+0x44` | 1 | 5N |
| `+0x40` | 1 | N |
| `+0x1c`, `+0x20`, `+0x24` | 2 | N each |
| `+0x28`, `+0x2c`, `+0x30` | 1 | N each |
| `+0x34`, `+0x38`, `+0x3c` | 1 | N each |

The last nine reads are interleaved by group: `+0x1c`, `+0x28`, `+0x34`,
then `+0x20`, `+0x2c`, `+0x38`, then `+0x24`, `+0x30`, `+0x3c`.
This supports the documented `attr_4c[N], attr_48[N], key[5N], attr_40[N]`
prefix and 12N bytes of metric arrays. These labels describe existing research
roles, not recovered vendor field names.

`0x10012030` handles buffered reads/refills. Its copy helper `0x10001480`
conditionally reverses individual multi-byte elements under a global flag.
Consequently, the width arguments alone do not establish file endianness.
This pass did not inspect that flag's producer or sample index bytes; use the
existing index reconstruction evidence for file-level claims.

### Candidate continuity and local scoring

`0x1001cd60` walks neighboring positions and class memberships around each
candidate unit. A separate byte array at model `+0x68` controls whether it
can step across adjacent units. It writes signed 16-bit candidate fields:

- `+0x0c`: left match count L.
- `+0x0e`: right match count R.
- `+0x10`: span S = L + 1 + R.

It writes a floating priority at `+4` equal to the negative of
`2*S + (L>0 and R>0 ? 1 : 0) + boundary_bonus`. The bonus is 2 when the
current position marker is 1 with a right match, or 2 with a left match.
It then calls `0x100140c0`; that helper's sorting implementation was not
examined here.

Its caller `0x1001cff0` gathers candidates from class lists, stops expansion
at 10,000, and calls this continuity helper. If the leading candidate spans
the full position count, it retains the leading full-span group. Otherwise,
when more than 30 candidates remain, it derives a split within the first 30.
It evaluates candidates with `0x1001c860`, orders the remaining suffix via
`0x100140c0`, and caps the output at 30. Thus continuity affects candidate
retention as well as the local score's weight.

`0x1001c860` compares five-byte keys, applies category-table penalties,
splits key bytes 3 and 4 into decimal quotient/remainder components, and
combines a weighted floating distance with integer penalties. Ghidra assigns
some arguments float types even where the body uses them as IDs or pointers;
the recovered signature is unsuitable as a reimplementation contract without
assembly or runtime corroboration. The pseudocode also shows 1,000-point
penalties for several invalid or disallowed feature combinations. No score
values were measured during execution in this pass.

## Evidence ledger

All IDs below are retained in the linked bundle with full dossiers and
provider limitations.

| Question | Status | Evidence ID |
| --- | --- | --- |
| Binary overview | Answered for this analysis | `ev_599b58c3217d406175512700571f88fe23a3ac19fd21331d73eab36f10d7f92c` |
| Named API inventory | Answered for `VT_` procedure search | `ev_6378647cc87352c32819f06a6ebbdf2f7e8fda014743316415c59ce8b1ed27b0` |
| Load orchestration | Partial; helpers not traced | `ev_ae3ea25ff132e4fec2a57de6c2d50d2c4e8da113bbb79c2293b87f7c15b0cdf3` |
| Buffer dispatch | Answered; mode semantics unresolved | `ev_98b18e9980a8936776b5ef38a4053733a5162c40edaffe13a18889240376266a` |
| Index tail widths/order | Answered statically | `ev_3789a8a655c306f8523a9df06fb9e3fde2e28410476c71b9c3c0e00b7a66b4c6` |
| Buffered read | Answered statically | `ev_b97eee7e2a2d37c6c2e7e1d10d3c32d602fc426c9d50fe4b30757a849a4fa2da` |
| Conditional byte reversal | Partial; flag provenance unresolved | `ev_d7ea88f8f41625616a67fcb7131c8cefbeacd28efdfb27b02b8a6dc0dfb71a22` |
| Continuity fields and priority | Answered statically | `ev_46c0aff315fd6f9967ce229718e71f9347b1f70d8a3714972863500ff4d6db8f` |
| Candidate retention and scoring calls | Answered statically | `ev_d3be410f18f6ce76ae32fcbda0b0af2fbc99d0bbb3dcf1d62b00bc3595bb64c3` |
| Local scorer | Partial; types and live values unresolved | `ev_0283012cdfbe4adacaea11a105e9e835fd60f8bf190b3716683cfed078e78b20` |

## Validation boundary

REA imported and analyzed the original target, returned the overview and eight
function dossiers, and exported ten evidence records. The target hash was
checked again after analysis. No implementation code changed, so no unit tests,
lint, or type checks were run. Documentation whitespace and final Git status
were checked. No runtime compatibility, synthesis, or whole-engine parity is
claimed.
