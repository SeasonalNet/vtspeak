# 2006 selection core: live REA walk

Target: original Kate `vt_eng.dll`, SHA-256
`00fc9375d08bd8cec303845c992481d79c8f616d1bfb6390a5322cd06f69273d`.
Ghidra 12.1.4 through REA; static observations, no engine execution. The
complete 29-record session is in `selection-evidence.json.gz` (deterministic
gzip of REA's original JSON export). Addresses below
refer to this DLL, not the 2013 DLL. All recovered code is pseudocode.

Question coverage: class-key and fallback production answered statically;
unit expansion/continuity and candidate filtering answered statically; target
cost answered with semantic feature-label gaps; join cost/DP/pruning/backtrack
answered statically; selected-descriptor handoff answered through dispatcher.
Runtime scores, phonetic table naming, and measured output parity remain open.

## Class production and fallback

`0x10022f60 -> 0x1001e470` is the caller edge into selection. The selector
walks phrase/syllable structures through utterance `+0x54` and emits six-byte
position descriptors at utterance `+0xb99e4`: owner index, substructure index,
phone index, position index, whole/half marker, category flag.

For a positive predicted duration it first calls `0x1001dfb0` (whole marker 0).
This gets class IDs through `0x1001da80`, sorts and removes duplicates, retains
at most 30 using `0x1001d5f0`, and accepts a whole position only when the total
unit population across those classes exceeds two (`0x1001d5c0`). Otherwise
`0x1001e110` makes two descriptors with markers 1 and 2 and performs the split
lookup. The split path is thus a class-population fallback inside the same
global selector, not a separately invoked synthesis engine.

`0x1001b2d0` maps a five-byte target key into a seven-byte class key. Its
first three bytes use categorical tables selected partly by decimal digits
of input byte 3; bytes 3 and 4 are copied; bytes 5 and 6 map right and left
context using `0x10070050`. `0x1001b3b0` derives eight-byte marker-1/marker-2
keys, selectively collapsing small decimal feature values and mapping context.

`0x1001da80` expands decimal feature alternatives (including 1/2 choices),
category substitutions through `0x1007082c`, and selected 3-to-0 or 3-to-4
fallbacks. `0x1001da00` dispatches whole versus split lookup and appends class
IDs. Whole lookup `0x1001d6f0 -> 0x10011e80` searches the sorted seven-byte
class records by prefixes of length 7, then 6, then 5, stopping once more than
two units are covered; it caps collected classes at 300. Split lookup
`0x1001d8c0 -> 0x10011fa0 -> 0x100248f0` searches the sorted split-class ID
lists by transformed eight-byte key prefixes, relaxing from length 8 and
ranking up to ten matches by class distance. Prefix length is carried in an
adjacent byte of the output descriptor; pseudocode omits that implicit coupling
from the apparent helper signatures. The split lookup's initial retry loop has
no visible lower bound before a match; no runtime safety claim is made.

Class ranking `0x1001d530` adds mismatch weights over the seven-byte keys:
`[40,100,40,20,10,2,1]`. Split keys add eight mismatch weights
`[10000,10000,5000,2000,1000,500,200,100]`. These are class filtering costs,
separate from the full per-unit target scorer.

## Unit expansion, continuity, and target cost

`0x1001cff0` expands model class counts (`+0x8c`, unsigned 16-bit) and class
unit lists (`+0x90`) into temporary candidate records, stopping at 10,000.
Candidate record fields are target cost `+0`, sort priority `+4`, unit ID `+8`,
signed 16-bit left/right/span continuity `+0x0c/+0x0e/+0x10`.

`0x1001cd60` tests neighboring position class membership against unit-to-class
IDs (model `+0x30`) while respecting unit adjacency flags (model `+0x68`)
and the whole/half markers. It computes L, R, S=L+1+R and ascending priority
`-(2*S + (L>0 && R>0) + boundary_bonus)`. Boundary bonus is two for marker 1
with a right match or marker 2 with a left match. `0x100140c0` does indeed sort
candidate pointers ascending by the float at `+4`.

When the leading span covers the entire position sequence, only the leading
full-span group remains. Otherwise a continuity-priority break within the first
30 becomes a protected prefix. Remaining candidates are ranked by target cost;
source category `d` multiplies that ranking priority by five. Output is capped
at 30. The source category also adds ten inside the target scorer itself.

`0x1001cff0 -> 0x1001c860` compares the five-byte target and unit keys. The
central byte mismatch costs 100 when either value is `0x4b`; otherwise five for
equal broad category or 1000 for unequal category. Left/right mismatches use
category penalties 2/5/10, with special cases costing 100 and half markers
disabling the corresponding outer-side comparison. Bytes 3 and 4 are split
into decimal quotient/remainder and scored with fixed 4x4 and 6x6 tables;
out-of-range or disallowed context combinations add 1000. Further three-way
category tables and source-category penalties are added.

The target's first byte feature is predicted duration, quantized by
`0x1001cfb0` as truncation of `duration/50 + 0.5`, saturated to 0..255 (confirmed
from caller assembly and the helper's FILD/FDIV/FADD sequence; its pseudocode
incorrectly shows no argument). Two additional target bytes are sampled from
the prosody structure using duration-proportional bins over a 12-bin contour.
Byte features use the model `+0x18` triangular distance table. With category
flag 0, distance is the sum of three distances; otherwise it is 0.2 times the
duration distance, and some matching-context integer penalties are halved.
Return is `integer_penalty + continuity_weight * feature_distance`.

Let E be the continuity extent chosen by the caller: full-span S when the
candidate spans all positions; R when its L reaches the left sequence edge;
otherwise L when its R reaches the right edge; otherwise min(L,R). The weight
is `1/(E + max(S/2,1))`. Several pseudocode arguments are mislabeled float
even though assembly and uses carry unit IDs and pointers; do not copy the
recovered C signature into an implementation contract.

## Join cost, dynamic programming, and backtracking

`0x1001e470` initializes the first position with its target costs. For later
positions it evaluates each current unit c against every retained predecessor p:

`D[t,c] = T[t,c] + min_p(D[t-1,p] + A[p,c] + P[p,c])`.

Acoustic join term A is zero for the same unit across marker-2 split halves,
or natural adjacency `p=c-1` when the predecessor's model `+0x68` byte is 1.
Otherwise A is `w0*d16 + w1*d8a + w2*d8b + 2`. The 16-bit endpoint IDs come
from model arrays `+0x3c/+0x40/+0x44` and index model `+0x0c`'s triangular
float table; byte endpoints use `+0x48/+0x4c/+0x50` and
`+0x54/+0x58/+0x5c`, indexing the shared model `+0x18` distance table.
Endpoint sides depend on the half markers. The five contextual weight rows at
`0x100707f0` are `(10,20,20), (10,10,2), (10,5,1), (10,0,2), (10,10,10)`.
Their row selection uses current/previous category flags and an owner/context
boundary test. Numeric row meaning is observed; phonetic naming remains open.

P is additional symbolic join cost. In a split-half transition it includes a
ten-point nonidentical-unit penalty under a category flag and 1000 for central
key mismatch. Whole transitions use the 6x6 table at `0x10070648`, indexed by
previous key byte-3 remainder and current byte-3 quotient; context mismatches
can add 100 or 1000, and a flagged nonnatural transition can add ten. The zero
acoustic shortcut does not bypass these symbolic tests.

`0x1001d4e0` selects the first minimum (strict `<`, inclusive index range).
Per-position DP costs are float arrays at utterance `+0x8d27c` with stride
`0x78`; predecessor indices are signed 16-bit entries at `+0xaa7b4` with stride
`0x3c`. After each later position `0x1001e2f0` prunes lists longer than ten:
keep costs `<= mean + (mean-min)/4` or candidates whose right continuity is
nonzero. It compacts costs, predecessor indices, and unit IDs together. Thus
this is dynamic programming over heuristically retained candidates, not an
unpruned global optimum over all vendor units.

The last-position argmin seeds backtracking through the predecessor matrix
into utterance `+0xb9214`. Selected IDs are then materialized with
`0x10013cf0`, writing 20-byte descriptors at utterance `+0x117d98`. That helper
dispatches to `0x10013a40` versus `0x10013bc0` based on model `+0x84`; root's
audio walk covers those readers. Optional counters track considered and selected
units when global `0x100942d1` is set.

## Evidence ledger and remaining limits

| Scope | Evidence |
| --- | --- |
| Selector/DP/backtrack | `ev_c64d3cba309655d57b4feb64f1eceae3eaca27cb3c5ce00a82df41d1cca67678` |
| Candidate retention | `ev_6af593f4c1598ba72031fe9ecf46e36c3c6c3e88162d6d9dea64b9c9e0c236cd` |
| Target scorer | `ev_6c0eaf9ad50eb1ede28c62ccc61b38e7064d3cc83cdb6ecdfc7824df69c35636` |
| Continuity | `ev_92097ceac7b9d6ad3b81648d98776acb832e363b145ed98cb5196e7cc4d5f6f2` |
| Class alternatives | `ev_0866eccc85d047139e986ad073c507e1f798a86565f6eda060fe4c1e22c7a6e0` |
| Whole and split class production | `ev_3698a6a3af04f32d383e99821463eb108ace46b70b1cc542d5f7c95d0f9f58d6`, `ev_fad9be24425f750ddc6a494c9b5af30179ebc67327587d44b1f83fbfc52a1d01` |
| Whole/split key transforms | `ev_2fbba78a1c8cd9b5526870ced2a074ba9f1d7c871273771ef1fa7bde71116279`, `ev_62314d0228b0239ad70fc7b05479fdc132a3d7edd8c0cc2bdaa502f73deec252` |
| Whole/split lookup | `ev_810be501031b219e88d3e91009bf52a0c3b137845d8539ed62844f6d442900e8`, `ev_93685b6241f8dc58c30017c30263b7d26ec2528184b07d433996e71f6911c163`, `ev_3adb0df29b1c8ba933a4b9d1ea596c212d8506dddda5833d3a549b8c2920aa5e` |
| Class distance | `ev_85865fe136b86966769a5525ac6e687fc3fcd977edf4178ea93f78acc78e12c6` |
| DP prune and argmin | `ev_366fa1b52ff67f6d23d59ded05b8c18584a178e496d2a92366add9d9607e4160`, `ev_87f582b7bc075436fa666852bc7acdda93ae2db2aac3c19606ba191efe04f41e` |
| Ascending sort | `ev_16b58e62cfb0f34f4255e385cd521803c7bf1d24343e967d364e04be5c722294` |
| Duration quantizer assembly | `ev_39b40260eebbd51a837dfa8f8f992c401ee08c749f368023b83646a1b65b74bc` |
| Numeric tables/constants | `ev_d951e225fe0b4cd34448ef6b03e053323652583c285dbfdbd6dd608a6b5af83c`, `ev_05e2f1f66af36cf7dbc20575d48e756a414cc6349b2984e634d84ea7e9314ba4` |

All formulas are static deductions from the above pseudocode, assembly, and
read bytes. Byte-table phonetic meanings and runtime branch frequencies remain
unresolved. No new runtime score, selected path, PCM output, reconstruction
parity, or malformed-input behavior was measured. No implementation code changed.
