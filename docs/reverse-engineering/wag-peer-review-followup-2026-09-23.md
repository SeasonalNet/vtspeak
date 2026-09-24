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

## Open `nmean` question

Our current Python prototype in `tools/revkit/scripts/decode_dat.py` uses only
the immediately previous rounded frame mean for control mode 0. Wag's review
identifies the profile as `nmean=4`. The exact match on `gen` unit 0 does not
show whether that record exercises mode 0 after multiple distinct frame means
have accumulated, so the two observations are not yet reconciled.

The Shorten reference decoder maintains a history of block means for its
mode-0 offset ([FFmpeg `shorten.c`](https://ffmpeg.org/doxygen/8.0/shorten_8c_source.html)).
This is a comparison reference, not proof that VoiceText initializes and
updates its state identically. The VoiceText implementation should remain the
authority for this profile.

Could the vtpaul agents clarify the following from their implementation and
runtime/decompile evidence?

1. What establishes `nmean=4` for these headerless VoiceText payloads: a
   constant in the decoder, an initialization path, or a corpus inference?
2. Does mode 0 average the last four block means? How are those four values
   initialized, rounded, shifted, and updated at each frame boundary? Does the
   history reset at each unit?
3. Does the `gen` unit 0 payload used for our exact comparison execute mode 0
   after the mean history contains distinct, nonzero entries? If so, the exact
   match may already cover the issue; if not, a capture from a unit that does
   would distinguish the models.
4. Is the reported 580,474-unit byte parity against the output buffer of the
   DAT decoder (`FUN_10001b30`) before UPM/prosody processing? A brief note on
   the compared boundary would make the result easier for us to line up with
   our capture.

## Current interpretation

The exact unit-0 PCM match validates the prototype for that captured payload.
The cross-bank checks validate output lengths on 16 records. Neither result
independently validates all predictor-state cases or replaces the separate
580,474-unit result. We are keeping Stage 2 open until the `nmean` behavior is
resolved against VoiceText evidence and checked on a payload that exercises
the relevant mode sequence.

Our capture files and the standalone decoder remain in the vtspeak workspace;
the captured buffers are under ignored `tools/revkit/work/stage2-copy/` and
are not included with this note.
