# 2006 audio preparation, waveform synthesis, and output walk

Target: original Kate16 `vt_eng.dll`, SHA-256
`00fc9375d08bd8cec303845c992481d79c8f616d1bfb6390a5322cd06f69273d`.
All newly collected observations are static REA/Ghidra 12.1.4.
No engine code executed or bytes changed. Evidence: `audio-evidence.json.gz`.
Decompiled names, types and function boundaries are provisional.

## End-to-end static call spine

`0x10020510` (caller-facing PCM fill) ->
`0x10020460` (fill/advance loop) ->
`0x10022f60` (selection/synthesis mode dispatch) ->
`0x10022fd0` (timeline loop) ->
`0x100232d0` (selected rows to parallel timeline arrays),
`0x10023670` (synchronization records),
`0x10023a10` (DAT read/decode/crop/gain),
`0x10023d80` (direct or modified grain synthesis).

`0x10022f60` invokes selector `0x1001e470` unless context byte +0x20
is 1 or 2, or context +0x48 indicates resumed synthesis. Mode 3 bypasses the
synthesis call. Recovered return bits/signatures here are unreliable;
the context fields carry actual loop state.

## Selected ID to original index and descriptor (follow-up)

Additional live REA evidence is in `descriptor-evidence.json.gz`; session
closed after export. This closes the handoff between selection and timeline.

`0x10013cf0` dispatches on model byte +0x84: zero ->
`0x10013a40` (unversioned interleaved record), nonzero ->
`0x10013bc0` (compact versioned record).
Evidence: `ev_d382be4307d8264789e9f24529019d642af8637c7ac457ecae28f6d956f47500`.

Both readers find the index bank for global selected ID G using records at
model +0x88, stride 0x68, with bank count at model +0xe98.
Bank 0 contains [0,bank0.+0x50); later banks compare consecutive global
boundaries at bank +0x60 and next bank +0x60. Per-bank +0x60 is cumulative
global start, +0x50 local row count, +0x54 on-disk record-array byte start,
+0x5c index file handle. Local ID = G - bank_start; read offset is
record_start + model.short(+0x76)*local_ID. Both cache selected index handle
at model +0x7c and read one complete stride through 0x1001f5a0.

The resulting descriptor bank byte (+0x13) comes from the global per-unit
map at model +0x6c indexed by G. This is the DAT/UPM asset bank; it need
not equal the index-bank ordinal. The loader rebases per-bank asset IDs by
cumulative asset-file counts when assembling that global map.

For original Kate unit-gen2, index bank start 179995, local count 98133,
record start 46 and stride 19; selected G=272822 becomes local 92827 and
file offset 1763759. Its raw compact row is
`18212604980262026c02644819000a08272024`:

| Compact source bytes | Selected descriptor bytes | Width and role |
| --- | --- | --- |
| +0x00..03 | +0x00..03 | u32 DAT payload byte offset (sample 69607704) |
| +0x04..05 | +0x08..09 | i16 first-part decoded sample count (664) |
| +0x06..07 | +0x0a..0b | i16 second-part decoded sample count (610) |
| +0x08..09 | +0x0c..0d | i16 compressed DAT byte length (620) |
| +0x0a..0d | +0x04..07 | u32 UPM byte offset (1656932) |
| +0x0e..0f | +0x0e..0f | two u8 UPM interval counts (10,8) |
| +0x10..12 | +0x10..12 | three u8 period/spans (39,32,36), later doubled |
| global bank map[G] | +0x13 | u8 DAT/UPM asset bank |

For this row's whole-unit mode, timeline sample count is
664+610-2*32=1210, matching the existing old Hello runtime capture.
Its trailing span is 2*36=72. This is a direct link from file bytes through
selected descriptor and timeline arithmetic to the earlier observed output.

All compact multibyte values above are little-endian. Reader
`0x10001480` copies width*count bytes when global 0x10070974=0 and
otherwise reverses each width-2/4 element. `VT_LOADTTS_ENG` initializes that
global through `0x1001f8b0`, which simply returns zero in this x86 build.
The byte-order helper differs from DAT bit-reader's big-endian words; those
are separate formats. Evidence:
`ev_aa2919f019a96fd4d4c80fd89289269b2499c424c3e7b9aa1a13d79badfdd231`,
`ev_ae3ea25ff132e4fec2a57de6c2d50d2c4e8da113bbb79c2293b87f7c15b0cdf3`,
`ev_36e4859f3e5af23564f4f76252445b900ab1b0276e3dee9ded9f1d8b18528108`.

`0x10012960` marks an index versioned when its first byte gives a string
length and the following bytes begin `ver.` (four-byte comparison).
It sets per-bank +0x64=1 and header skip +0x66=length+1; otherwise both zero.
Raw unit-gen2 begins `17 76 65 72 2e 32 30 30 35 00 ...`,
i.e. length 0x17 and `ver.2005`. The loader `0x100129f0` requires all
bank flags equal, copies the flag to model +0x84, reads the versioned
record stride, and copies it to model +0x76. The code tests the presence of
the prefix, not a numeric year/version policy.

On zero-flag/unversioned branch the loader sets model.short(+0x78)=6 and
short(+0x7a)=4. Thus `0x10013a40` reads the equivalent fields from
raw offsets +7 (DAT u32), +11/+13 (sample i16 pair), +15 (compressed i16),
+17 (UPM u32), +21/+22 (UPM counts), and +24/+29/+34 (span bytes).
This branch supports the older interleaved metadata layout; the inspected
Kate originals use the compact nonzero branch.
Evidence:
`ev_78be00cf87258d25aaccb0b5efcb0449655e4b018fdd3727284e6a696bc4cb25`,
`ev_a0c158e67299c673436d7114b636ce247da04a6f2b3141cdf6ec2e846238bcaf`.
Loader IDs:
`ev_0833dd4e55bde5ca36e59c6eec9b7b79222d71fd35dba3a45bd41e6059e1b594`,
`ev_9f8b1b6486a32c3427027d38638d12f928718cb65ed67d5824da1216e0debb94`.
Prefix bytes Evidence:
`ev_286a8d79bd95d4d5be216f91aec40667e818d0ad4eaa8f623752072e72db7584`.

### Grain record type confirmation

The five dwords produced by `0x10024520` are integers, including the last
two scale fields. Assembly sign-extends source shorts with MOVSX, writes
dwords with MOV, advances destination by 0x14, and performs no floating
conversion. Consumers use IMUL/IDIV for scale arithmetic. Floating values
enter later through the coefficient table in `0x10023d80`.
Evidence:
`ev_b65727602c8d93350092f9c8cff75dd85c4684f4e632bbb64a9cb1ed15cd7875`.


## Selected descriptor to timeline

`0x100232d0` walks 20-byte selected records (stride 0x14), plus six-byte
companion selection rows. It constructs arrays relative to timeline T
(which caller passes as context +0xf8ecc):

| Offset | Width | Observed role |
| --- | --- | --- |
| +0x1130 + 2*i | 16 | asset bank from descriptor byte +0x13 |
| +0x3c28 + 4*i | pointer | original selected 20-byte record |
| +0x4d58 + i | 8 | portion mode (1/2/other) |
| +0x2af8 + 4*i | 32 | cropped decoded sample count |
| +0x51a8 + 4*i | 32 | UPM read offset |
| +0x62d8 + i | 8 | number of UPM interval bytes |
| +0x6724 + 2*i | 16 | leading period/span, descriptor byte doubled |
| +0x6fbc + 2*i | 16 | trailing period/span, descriptor byte doubled |
| +0xabe4 + 4*i | 32 | pitch scale (role inferred from grain spacing) |
| +0xbd14 + 4*i | 32 | duration scale (role inferred from source progression) |
| +0xce44 + 4*i | 32 | gain percentage, directly used to scale PCM |
| +0xdf74 | 16 | populated timeline count |

Descriptor mode 1 uses signed short +8 as sample count and dword +4 as UPM
offset; mode 2 uses short +10 and advances UPM offset by byte +0xe minus 1.
Mode other uses count = short(+8)+short(+10)-2*byte(+0x11), and UPM
count = byte(+0xe)+byte(+0xf)-1. Leading/trailing spans are doubled byte
fields +0x10/+0x11/+0x12 depending on mode. The routine stops at 1,100 rows.

Pause rows use timeline +0x7854 (whereas later boundary tests use +0x7855,
the next row). Pause counts are scaled by (base duration * config +0x3830
+50)/100. Pause marker +0xa798 is 0x5a. Caller emits zero PCM directly.

## DAT preparation

`0x10023a10` selects DAT bank handle from the array at *manager, indexed by
timeline bank. It reads descriptor dword +0 as file offset, byte-size width 1,
signed short descriptor +0xc as read length into context +0x1127a2 using
`0x1001f5a0`. This helper serializes with a critical section, seeks unless
cached position already matches, reads, optionally transforms read bytes
when global 0x10070974 != 0, and updates cached position.

It calls decoder `0x10001d50`, then copies count*2 bytes from decoded scratch
context +0x10d982 into context +0x108b62. Portion mode 2 offsets the source by
2*(descriptor short(+8)-timeline leading span). It stores that PCM pointer
at timeline +4*i. Gain !=100 applies signed integer sample*gain/100, followed
by clipping to [-32767,32766]. Return-code propagation is weak: this wrapper
returns 1 without checking the read or decoder return value.

## DAT codec: bounded recovery, decompiler blind spot

REA pseudocode for `0x10001d50` is incomplete: computed jump at 0x10001e84
is treated as an indirect call; the predictor handlers are outside its recovered
body. Its recovered void parameter list is wrong: assembly reads input,
output and length arguments from the caller's stack. Do not call its
pseudocode a complete decoder.

Observed entropy primitive `0x100020b0(state,k)` counts zeros before a
one bit, appends k remainder bits and returns (zero_count << k) | remainder.
It reads 32-bit words MSB first. `0x100021c0` assembles four file bytes in
big-endian order. Signed residual `0x100022a0` maps unsigned u to
u>>1 for even u, ~(u>>1) for odd u.

REA retained jump-table bytes at 0x10002084 (Evidence
`ev_543e2f7e17ee4d91ed5deca912a8d27a36abc2898e498a8cc1a35f4461e3d209`)
select:
0 -> 0x10001ea2; 1 -> 0x10001ed2; 2 -> 0x10001f00;
3 -> 0x10001f31; 8 -> 0x10001e8b. Cases 4–7 point at 0x10001f68
in this second table but top-level command dispatch handles 4/5/6 first.
Command 4 ends decoding; 5 changes block count through `0x10002360`;
6 changes bit shift. Block count starts at 256.

The following predictor formulas are analyst readings of the retained x86
bytes, corroborated with existing objdump report
`tools/revkit/work/reports/vt_eng-2006-objdump-disassembly.txt:1414`:
command 0 residual + block mean; command 1 residual + x[n-1];
command 2 residual + 2*x[n-1]-x[n-2];
command 3 residual + 3*(x[n-1]-x[n-2])+x[n-3];
command 8 zero block. Predictor byte Evidence:
`ev_ac95b321eef00e42fe62f2cd7b34065f21cae418c952f47fe7649eb80e0883fb`.
Post block mean/history updates exist at 0x10001f68–0x10001fc0;
`0x100022c0` restores shift, `0x100022f0` writes shorts and advances byte
cursor. The latter only explicitly upper-clamps values to 32767 in pseudocode.
This is a Rice-coded differential PCM family; identifying an exact published
codec or claiming all malformed streams handled is outside this evidence.

## Direct join: pitch=duration=100

`0x10023d80` immediately delegates to `0x10023b40` when both scale arrays
are 100. A nonfinal row whose next-row pause flag !=1 emits
count - trailing_span samples; final/pre-pause row emits the whole count.
Tail samples are saved at context +0x1079fc. Pending tail count and consumed
offset are context +0x1081cc/+0x1081ce; caller's carry short is updated.
Output destination is context +0x3c plus byte cursor +0x38.

Within pending carry m, direct join weighted sample k uses:
a=floor((m+k+1)*8191/(2*m+1));
b=floor((k+1)*8191/(2*m+1));
output=trunc(clamp(tail[k]*T[a]+current[k]*T[b])).
The coefficient table starts at 0x1006119c, its count at 0x10061198
is 8192 (REA byte Evidence
`ev_1512a10b9322174d8ab5dc47474c08d59b1f44638ab30936a2a8ac8499acaaea`).
The floating expression is missing in Ghidra pseudocode for this routine;
its complete dossier assembly explicitly retains multiplies, add, clamp,
`__ftol`. Additional raw instruction Evidence:
`ev_f48e4872d46a6997047c5175163cdd4f08678d3a197359284e77f7cca64e1c71`.
After carry, samples copy verbatim. No extra onset or final envelope here.

Existing prior runtime corroboration (not newly executed):
`docs/reverse-engineering/lead5-kate-common-resource-comparison-2026-09-26.md:3015`
records Hi. source counts 1420/5910, 72-sample carry and 7258 output frames;
the formula reproduces all join samples exactly, interiors copy exactly.
This strengthens the direct path only.

## Pitch/duration modification

If either scale !=100, `0x10022fd0` calls:
- `0x10024460`: read interval bytes from UPM bank array at manager +4,
  using timeline UPM offset/count, widen each byte to ushort then double it;
  pointer stored at timeline +0x51a4 by caller.
- `0x10024520`: build five-dword records
  [source_start, left_interval, right_interval, duration_scale, pitch_scale].
  For N>1 interval bytes, N-1 records begin at cumulative interval offsets.
  For N=1 it uses [0, leading_span-1, 1, duration_scale, pitch_scale].
- `0x10023d80`: select and window grains; `0x100245e0` chooses the record
  whose center (start+left) is closest to the requested position, stopping
  once requested position >= final record start+left+right.

For a record with left interval L, right R and pitch P:
nominal output hop = trunc(100*L/P), clamped to [trunc(L/16),
trunc(31*L/16)]. Source selection request uses current generated extent *100
/duration_scale, so source progression and pulse spacing are independently
controlled. This is the static basis for pitch/duration names.

Each source grain has L+R samples. For k<L it applies rising weight
T[floor((k+1)*8191/(2*L+1))]; for k>=L falling weight
T[floor((R-L+1+k)*8191/(2*R+1))].
It converts/clamps weighted lanes to shorts in context +0x1081d0, then
adds overlapping grains with integer clipping. Grain placement can overlap,
copy leading pending carry, or insert zero gaps. At a nonfinal/non-pause
boundary, it saves the remaining right interval as carry and advances cursor
by generated hops; final/pre-pause row advances by total extent. This supports
a pitch-synchronous overlap-add interpretation; no new altered-pitch runtime
capture was made, so numerical parity for this branch remains unproven.

## PCM streaming and postfilter

`0x10022fd0` fills up to 60,000 bytes, preserving row index at shared
state +0x5285e and resumed flag context +0x48. Silence can be split exactly
at buffer capacity; non-silence reserves 2,000 bytes plus scaled count.
It adds each emitted byte count into word/phone synchronization records
constructed by `0x10023670`, including ring wrap at 200 records.

`0x10020460` resets byte count, fills until resumed/full/done, advances
text position after completed timelines and invokes `0x100202f0` for more.
`0x10020510` sets caller PCM pointer, invokes fill and returns byte count.
It optionally applies one-sample filtering controlled by manager +0x383c:
positive coefficient uses previous filtered sample,
y[n]=clip(x[n]+trunc((previous_y*c+50)/100));
negative coefficient traverses the input backwards and uses previous input
sample instead, y[n]=clip(x[n]+trunc((x[n-1]*c+50)/100)).
Previous sample is retained across calls at context +0x11cc80.
This identifies an additional synthesis output stage after joins.
Public buffer/file/playback callers all reach 0x10020510, but caller-specific
encoding/resampling/file headers were not followed in this audio pass.

## Evidence index

| Function | REA Evidence ID |
| --- | --- |
| `0x10023a10` | `ev_349e4841dfce3f3cb16e322994ef5f106cbc722071762104ac8ab348fa9918bf` |
| `0x10023b40` | `ev_d9ed8700557284de5249dd0715144ab24887d3b587cef3b8f3bb0b5c48c66467` |
| `0x10023d80` | `ev_211c0b407afe439e5a59de54bae121652d6b1eb9b0f9fc15cc69f1fe08562936` |
| `0x10022fd0` | `ev_f11617ca37c86029a67c213beb8eb201182f61bef301ade9d5281858d0483611` |
| `0x100245e0` | `ev_539f80e7ea718a217578e64caec2050154b84dd8c550653e0f7e5ab40f01722f` |
| `0x1001f5a0` | `ev_140ac18f29c65b5696ea4a2773ee607c5bdbce0b4a7fdee605bfd97f6cdc5c5b` |
| `0x10001d50` | `ev_53fca8fb26a16235a34c66d4b9b0696e0a9bfe3dc0b28ae42a31463563c656fa` |
| `0x10024460` | `ev_438abbfa6052be3721d87a1188b8104d8e62f75fb84df71df22125a266f4cd57` |
| `0x10024520` | `ev_b65727602c8d93350092f9c8cff75dd85c4684f4e632bbb64a9cb1ed15cd7875` |
| `0x100232d0` | `ev_a3e5e26c59e3328a999aebe3bd4b24425753a970cd6a0cb3611d2a22241587a5` |
| `0x1001f570` | `ev_7f4f6f0752af3e346f900a99244ec6a01685e7a0bf9f171f40da2b4604c5e5ea` |
| `0x1001f610` | `ev_a584d79a3d39edfea9bea5fe531a878a1c59e89bb91dd4ab2755e214c485cd9e` |
| `0x10023670` | `ev_831fbcb688b4661cba476181c6f516b6b18bb48c89128bbb12983618258c9804` |
| `0x10022f60` | `ev_69871cbf20c26dde4f5e7466f9f1bd39218705fe0b146b0b1508a377ad662204` |
| `0x100020b0` | `ev_b542a0588c30160e354ff5c94fdc4eeae335d5cc7b1c9b92ae78ddab9910c35d` |
| `0x10002360` | `ev_1dfc1b59c07bd236fa4c2a49a6ace2e6b7d55e1a81f137cc6bb5d23886f4ff5a` |
| `0x100021c0` | `ev_7931abff142fafec8b920040d4656945a87236e34e8af673eba67977a78ab1d9` |
| `0x100022c0` | `ev_cabc44ae54c8c0ac9f99b1ab6fc94c4fc603218fc465bad7aa254e7f458481ff` |
| `0x100022f0` | `ev_d3809419d8e3747a5a0a691fd49cc71c6c3004da0b0121325f0f5ab2c2e2306c` |
| `0x10020460` | `ev_25f51defaeb7a33d3e24b30e4ce3807f3521a881f878c68faa53913eca1ac710` |
| `0x10020510` | `ev_0e1236821c54803ae9072d0077edee5b19f11a73c7d4faba54e909ccc095bc07` |

Additional signed-residual function 0x100022a0:
`ev_53d9ddf531c5029d3a83e4356aace37d1563fff2b4708a5f61755c485a3fec77`.

## Remaining unknowns

- Full decoder CFG and original signature are not repaired in the database.
  Predictor formulas are disassembly-derived, not recovered pseudocode.
- Grain path semantics need runtime altered-pitch/duration PCM corroboration;
  low-level formula description does not claim parity or robustness.
- Return/error contracts, malformed asset behavior and actual sample rate must
  be established with loader/public API evidence.
- Caller output conversion, resampling, WAV header, playback lifecycle remain
  outside this branch's collected function boundary.

