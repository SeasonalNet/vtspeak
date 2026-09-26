# VoiceText reverse-engineering coverage (2026-09-25)

This page records which parts of the local 2013 M16 Paul VoiceText package
have been examined, what evidence supports that coverage, and which analysis
methods are useful for the remaining questions. Coverage here is reported by
evidence boundary; it is not a percentage of all possible execution paths.

## Techniques and when to use them

| Technique | Use it to answer | Current use and limits |
| --- | --- | --- |
| Ghidra disassembly, decompiler pseudocode, and cross-references | Find callers, data flow, branch conditions, and likely formats. | Ghidra 12.1.4 has pseudocode for all 59 public `VT_*`/`VTDTTS_*` exports and selected loader, license, text, synthesis, and output internals. Pseudocode is a hypothesis aid: names, types, and control flow may be wrong. The complete internal-function inventory and decompilation coverage have not been measured. See the [static analysis report](../../tools/revkit/work/reports/vt_pau-analysis.md). |
| `objdump -d` assembly listing | Check exact instruction bytes, instruction boundaries, branches, and calling details independently of Ghidra's pseudocode; provide a searchable address-based reference. | A broad listing of the shared VoiceText code is now preserved in [`vt_pau-objdump-disassembly.txt`](../../tools/revkit/work/reports/vt_pau-objdump-disassembly.txt). It covers the code sections `objdump` identifies as instructions, but does not assign semantic names or show which instructions execute. |
| `objdump -D` | Investigate code/data boundaries or a suspected code sequence in a section not selected by `-d`. | `-D` also decodes nonempty non-code sections. Much of that output is data rendered as meaningless instructions, so use it for a specific boundary question rather than as the main reference. GNU documents the distinction between [`-d` and `-D`](https://sourceware.org/binutils/docs/binutils/objdump.html). |
| Controlled Wine/GDB probes | Observe concrete arguments, registers, memory writes, return values, file reads, and output at a chosen boundary. | This is the main dynamic method used so far. It gives precise values for selected calls and inputs, not complete path coverage. Interventions should change one controlled value at a time and retain a matching control. |
| Intel PT | Observe executed branch paths with lower source-level intrusion, then correlate addresses against the static map. | Stage 13 captured a corrected control/intervention pair, but the paired branch comparison is unfinished. One earlier control decode yielded 591 branch records from 96 DLL source addresses; it had AUX loss and decode errors and is not a complete trace. See the [Stage 13 analysis](intel-pt-wine-capture-2026-09-25.md). |
| Pin or Frida Stalker | Collect compact basic-block or edge coverage when broad executed-code coverage is needed. | No checked-in run currently demonstrates either tool on this Wine target. A short feasibility pilot may be useful if PT decoding remains cumbersome; filter to `vt_pau.dll` and a known synthesis window, then compare its coverage with existing GDB markers and output. |
| Corpus scans and independent parsers | Check structure and broad data coverage across the local model corpus. | Strong coverage already exists for local DAT decoder output, index records, tree parsing, and dictionary inventories. Those checks establish the stated format/output boundaries; they do not name undocumented feature fields or prove whole-synthesis parity. |
| Controlled input/output comparisons | Determine which outputs change when a specific input, field, or selector changes. | Used throughout the stages, including selected WAV/PCM checks, abbreviation contexts, G83 interventions, and file-output API selectors. Results apply to the tested controls and cases; do not infer semantic labels from correlations alone. |
| Symbolic execution or emulation | Resolve a small branch-heavy function when concrete probes and static analysis leave competing explanations. | No whole-program symbolic pass is established in the current evidence. Restrict this to a bounded function with captured inputs and explicit stubs/models for external calls; whole-DLL Wine behavior would require substantial environment modeling. |
| Dynamic taint or broad data-flow instrumentation | Follow a value through many functions when its consumers are not apparent statically. | No broad taint run is established. Prefer a targeted register/memory trace for a named unresolved value, such as G83 score components, after narrowing the relevant path with GDB or branch coverage. |

## Current coverage map

| Area | Current evidence | Coverage boundary / open work |
| --- | --- | --- |
| DLL exports and top-level call path | Ghidra pseudocode for all 59 public exports; selected core pseudocode; `include/vt_eng.h` supplies public names and signatures. | Export inventory is not an inventory of every internal function. Not every selector, error, buffer size, threading mode, or exported API has runtime coverage. |
| Shared VoiceText engine code | The analysis report records identical executable `.text` bytes in `vt_pau.dll`, `vt_jam.dll`, `vt_jul.dll`, and `vt_kat.dll`; the Paul DLL is therefore the representative code image for this set. | Voice/model data differs. Results on Paul data do not establish synthesis behavior for other voices or package versions. |
| Other local voices and package versions | Lead 5 inventories all five local M16 models and attempts all four supplied host/DLL pairs. Paul and James each synthesize three shared fixtures; a package-supplied James DLL also synthesizes those fixtures with different PCM. Julie and Kate stop at the first missing `tree3` voice tree; Bridget has no local runtime pair. Stages 18–19 test Kate adapters without obtaining validated PCM. | This closes the checked-in pair comparison, not general cross-version support. The Julie legacy candidate stops at missing common data; the Bridget candidate lacks a matching host/package set. Whole-synthesis parity for other voices remains unestablished. See the [Lead 5 report](lead5-voice-package-runtime-2026-09-25.md). |
| Static machine-code reference | Selected disassembly excerpts already support specific findings, including the Stage 8 copy helper. The broad `vt_pau.dll` `-d` listing is now available as a cross-check source. | Disassembly helps verify bytes and control flow; it does not establish runtime execution or the meaning of proprietary data fields. |
| DAT decoder and local indexed corpus | All 580,474 local records were inventoried. The independent decoder was compared with the DLL at the decoder-output boundary by byte count and SHA-256; captured fixtures also have byte-for-byte PCM comparisons. | This is local-package decoder parity, not whole-synthesis corpus parity. DAT mode 8 is absent from the local corpus, so that static branch lacks a real runtime sample. See the [corpus parity and Stage 9 report](dat-corpus-parity-and-stage9-2026-09-24.md). |
| UPM, indexes, and `tree3` | Local index fields and UPM timing/join behavior are checked across the corpus; the parser covers all 17 Paul duration/pitch trees, with 307 runtime lookups compared. | Operational fields and tested tree inputs are mapped, while many acoustic feature names and physical interpretations remain unavailable. |
| Text, dictionaries, and pronunciation | The text/resource layouts, all 31,550 TPP keys, H/T/C inventories, and phone/context record layouts have broad static/corpus coverage. Stage 14 tests selected abbreviation keys; Stage 15 records 1,750 observations across all 427 H/T/C keys and selected contexts. | Corpus enumeration is not every possible linguistic context. Numeric TPP classes, some context conditions, and other semantic labels remain unresolved. See the [open-topics review](open-topics-review-2026-09-24.md) and [Lead 4 results](lead4-abbreviation-context-results-2026-09-25.md). |
| Candidate selection and synthesis | Stages 7–8 map candidate scoring, selection, timeline construction, reconstruction, and output for the documented feasibility cases. Stage 9 adds selected number/abbreviation sample-buffer-to-WAV checks. Stage 11 adds eight selected text-category captures with local PCM-to-WAV checks. | Not every context-gate combination or full-corpus synthesis output has been compared. The eight whole-file hashes differ from the peer manifest; peer runtime inputs are unavailable, so that mismatch is not localized. |
| G83 intervention and executed branches | Controlled G83-zeroing changes selected acoustic descriptors/scores and selected unit IDs in the documented probe. Stage 13 has a corrected PT capture pair and an older bounded control decode. | No paired PT edge comparison is complete; the decoded control is partial and from a different capture. G83 semantics and the score-component path remain open. |
| File and buffer output APIs | Stage 16 checks `VT_TextToFile_ENG` selectors 0–10, selected errors, buffer formats 0–3, short-call length behavior, a 289,636-byte chunk/drain/cancel lifecycle, file → buffer → file byte parity, declared text formats on ASCII plus one `<vtml_sub>` alias/control, user-dictionary selected errors, configuration boundaries, and playback with an ALSA null sink. The [API coverage map](lead6-api-coverage-map-2026-09-25.md) inventories all 66 named exports. | Nonzero thread IDs, other markup handlers and malformed markup, successful user-dictionary load/use, actual audible playback, broader call sequences, remaining error paths, and other voices remain unestablished. See [Lead 6 results](lead6-file-api-behavior-2026-09-25.md). |

## How to read this map

The strongest coverage is at specific structural and behavior boundaries: local
DAT decoder values, selected index/tree mechanics, mapped text and candidate
paths, selected synthesis outputs, and the tested file-output API cases. The
Ghidra export pass is broad at the public interface, but the repository does
not claim complete decompilation of every internal function or complete runtime
branch coverage. The PT captures are evidence that collection worked; until a
paired decode is completed, they do not provide a control-versus-intervention
branch map.

Use the listing as an address/byte reference alongside the Ghidra project and
the cited runtime probes. When adding findings, state the DLL/data package,
function or address, exact tested inputs, direct observation, and remaining
inference. Keep semantic labels and compatibility claims within the tested
evidence boundary.

## Assembly listing provenance

The text listing was generated from the read-only `binary/vt_pau.dll`
with host GNU Binutils `objdump` 2.46. Its input SHA-256 is
`200945bbb2853cc56b3eaa03770eedf26066b20bfc471ef5b0213387035b37ce`.
The output contains 160,474 lines and is 8,116,338 bytes (about 7.8 MiB).
Recreate it from the repository root with:

```sh
objdump -d -Mintel --show-raw-insn binary/vt_pau.dll \
  > tools/revkit/work/reports/vt_pau-objdump-disassembly.txt
```

`-d` selected the `.text` code section in this PE32 i386 DLL. The output keeps
the DLL virtual addresses and raw instruction bytes for comparison with
Ghidra's image addresses and focused runtime traces. It is a generated
disassembly of the proprietary input, not original source code.
