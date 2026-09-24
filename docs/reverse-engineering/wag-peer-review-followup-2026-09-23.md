# Follow-up peer review for the vtpaul agents

Prepared by the vtspeak analysis on 2026-09-23. This note is a technical
handoff for Wag to share with the separate vtpaul agent team. The two projects
remain independent; the checks below were run against the local VoiceText
package and our own prototype.

## Results from our side

| Claim or check | Result | Evidence and limits |
| --- | --- | --- |
| Shared UPM boundary byte | **Independently consistent** | Our runtime trace showed the two 7-byte reads at `0x0041a461` and `0x0041a467`; the ranges overlap by one byte. The index reader uses `first_count + second_count - 1` for the combined span. |
| Decoded sample counts | **Consistent on a small cross-bank sample** | We decoded the first, second, middle, and last record in each of `gen`, `num`, `etc`, and `alp` (16 records total). Every output count matched `2 * sum(combined UPM periods)`. This checks lengths and termination, not sample values for those 16 records. |
| DAT decoder against engine PCM | **Exact for one unit** | In an isolated Wine/GDB run, a breakpoint at `vt_pau.dll!0x10001b30` captured an 800-byte input and 3,036 output bytes. The input is `unit-gen.idx` record 0, DAT offset 0. Our `decode_dat.py` produced the same 1,518 samples, byte-for-byte across all 3,036 PCM bytes. |
| Shorten profile across all records | **Not independently reproduced here** | Wag's review reports headerless mono Shorten with block size 256, `nmean=4`, no QLPC, and exact engine PCM for all 580,474 payloads. We attribute that corpus-wide result to the vtpaul agents. |

The first two checks align with Wag's shared-byte correction and period-length
interpretation. We have not independently verified the meanings of index
record bytes 4–7 or 16–18, or the reported cross-fade use of the cached UPM
periods.

## Resolution: `nmean=4`

The Stage 2 pass subsequently reconciled this against the decoder pseudocode
and the broader captured set. The pseudocode at `0x10001c19` sums four
32-bit history slots, applies signed rounding and output-shift compensation,
and the shared frame tail shifts in the current rounded, shifted mean. The
decoder invocation initializes the four slots to zero. The captured checkpoint
contains 32 DLL calls for 27 unique payloads; every captured PCM output matches
the prototype byte-for-byte, and the set includes repeated mode-0 frames with
four-block history. This resolves the profile and the captured behavior for
the tested cases.

The Shorten reference decoder maintains a history of block means for its
mode-0 offset ([FFmpeg `shorten.c`](https://ffmpeg.org/doxygen/8.0/shorten_8c_source.html)).
This is a comparison reference, not proof that VoiceText initializes and
updates its state identically. The VoiceText implementation should remain the
authority for this profile.

The following questions motivated the original review. The Stage 2 evidence
now answers them as follows:

1. The four-slot profile is established by the decoder's reads and updates,
   not inferred from the corpus.
2. Mode 0 uses the rounded average of the four stored means, with current
   output-shift compensation; slots initialize to zero and shift at frame
   boundaries. Decoder state is per invocation/payload.
3. The 27-payload capture exercises repeated mode-0 frames and matches exact
   outputs. The checked-in evidence does not separately identify a frame whose
   prior four means are all distinct and nonzero, so that narrow adversarial
   state-coverage question remains unproven.
4. The external 580,474-payload report's exact comparison boundary is not
   established by this repository's evidence. It remains external corroboration
   and is not counted as local full-corpus parity.

## Current interpretation

The four-block mean behavior is resolved to the level supported by the local
decompiler and 27-payload runtime capture. Full-corpus local PCM parity, mode 8
runtime coverage, and proof of a distinct-nonzero four-mean state remain
separate validation questions; they do not reopen the recovered `nmean=4`
profile.

Our capture files and the standalone decoder remain in the vtspeak workspace.
The small captured DAT/PCM fixtures and `capture-many.gdb` are tracked under
`tools/revkit/work/stage2-copy/`; the full voice/model assets remain external.
