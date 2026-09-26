# Lead 6: `vt_pau.dll` API coverage map (2026-09-25)

## Scope and coverage labels

This map inventories the named exports of the local 2013 Paul M16 DLL and
tracks evidence for each surface. `objdump -p binary/vt_pau.dll` reports the
PE internal name `vt_eng.dll`, 66 named exports, ordinals 1–66. Ordinals 1–61
are function exports; 62–66 are data exports. `include/vt_eng.h` declares 13
of the 61 functions. The remaining names are still part of the binary's
export surface even when the shipped header does not declare them.

Coverage terms:

- **Runtime:** called in the local Wine process and a result or output was
  captured. The exact input and host conditions remain part of each claim.
- **Static:** export and wrapper pseudocode (or a shared-address alias) were
  inspected; no direct API runtime result is claimed.
- **Inventory only:** the export was inventoried, but its underlying behavior
  or data meaning was not characterized in the reviewed evidence.

Ghidra output is pseudocode and can misstate wrapper prototypes. In particular,
the header declares short returns for loader and user-dictionary functions
whose thin Ghidra wrappers are rendered as `void`. Runtime probes have
captured selected user-dictionary short returns and staged loader results;
other branches still need direct return-boundary evidence.

## Header-declared API coverage

| API | Runtime evidence | Remaining coverage |
| --- | --- | --- |
| `VT_LOADTTS_ENG` | **Runtime:** existing Stage 1–2 and Stage 9 Wine traces load the local Paul model and reach synthesis; prior stages capture model-load failures. | Not an exhaustive input/error matrix. Non-null license variants, every speaker slot, concurrent/repeated load, and wrapper return ABI across all failures are not covered here. |
| `VT_UNLOADTTS_ENG` | **Runtime:** unloading speaker 1 in-process makes subsequent file, buffer, and play calls return their declared database-not-loaded codes. | Other speaker IDs, unload while active playback/streaming, repeated unload, and full last-speaker DLL teardown are not characterized. |
| `VT_LOAD_UserDict_ENG` | **Runtime, selected errors and candidate rows:** indexes `-1` and `1024` expose low AX `-1`; the existing WAVE path and candidate CSVs expose `-3`. Tested candidates include plain and quoted three-field `P` rows, a four-field `P` row, and a normalized `A` row. The wrapper is declared `void`; these are raw register observations, not declared C returns. Private helper probes accepted source `vtspeakprobe` (normalized length 12), target `hello`, and target phoneme `HH`; its converter returned 1 and produced a one-character quote string. These helper successes did not make the candidate rows load. | Successful load, busy index, unknown error, accepted row format, and dictionary effect on synthesis remain untested. |
| `VT_UNLOAD_UserDict_ENG` | **Runtime, selected errors:** low AX was `-2` for index `-1` and `-1` for empty slot 0. The wrapper is declared `void`, so these are raw register observations, not declared C return values. | Unload of a loaded dictionary, in-use/busy behavior, and unknown error remain untested. |
| `VT_PLAYTTS_ENG` | **Runtime:** null text (`-2`), empty text (`-3`), database unloaded (`-4`), device-init failure (`-5`), and success (`1`) through an ALSA null sink. On success a WinMM handle was open and play state was 1. | Thread-create and unknown errors, completion message delivery, audible device output, repeated/overlapping play, and speaker selection are unverified. |
| `VT_STOPTTS_ENG` | **Partial runtime:** called while a null-sink playback handle was active; the inferior exited normally during GDB evaluation of this void call. | No call-return marker or explicit post-stop state was captured; effective stop/reset/close behavior remains static-only. |
| `VT_PAUSETTS_ENG` | **Runtime:** call returned while null-sink playback had an active handle. | No sampled playback state or paused audio data was collected; behavior on a physical device remains open. |
| `VT_RESTARTTTS_ENG` | **Runtime:** call returned after pause while the null-sink handle was active. | Resume timing, completed/stopped state, and physical-device behavior remain open. |
| `VT_TextToFile_ENG` | **Runtime:** selectors 0–5 and 7–9 succeed; selector 6 and out-of-range 10 fail. Null/empty text, null path, database-unloaded, repeat, text-format, configuration, and one VTML substitution case are covered. | Non-null unwritable paths, create-thread and unknown errors, full markup/error grammar, all speaker/dictionary combinations, and arbitrary call interleavings are not covered. |
| `VT_TextToBuffer_ENG` | **Runtime:** formats 0–3, null/empty arguments, invalid format, database-unloaded, short-call length semantics, 60,000-byte chunk drain, busy/cancel/no-context transitions, and one file → buffer → file sequence are covered on thread ID 0. | Nonzero thread IDs fail in this Wine setup; other hosts, abnormal/unknown errors, physical-capacity failures, and all format-specific long streams remain open. |
| `VT_SetPitchSpeedVolumePause_ENG` | **Runtime:** setter bounds and getter values checked; isolated pitch, speed, and volume extrema change the fixed sample's WAVE bytes. Pause extrema leave that sample unchanged. | Other slots, negative-value no-op cases, all intermediate values, and pause-bearing text are not characterized. |
| `VT_SetCommaPause_ENG` | **Runtime:** getter values confirm clamp to 0–65,535; tested extrema leave the fixed sample's WAVE bytes unchanged. | Pause-bearing text, intermediate values, and slots other than speaker 1 remain open. |
| `VT_GetTTSInfo_ENG` | **Runtime:** declared requests 0–26, invalid request 27, null value, string-size boundaries, and the null-license/default lookup path are covered. | `VT_INFO_ERROR_NOT_SUPPORTED_REQUEST`, unknown error, supplied license file, and every request's exact destination-size edge remain open. |

Detailed observations and source traces are in the [Lead 6 API behavior
report](lead6-file-api-behavior-2026-09-25.md) and [Stage 16](../../tools/revkit/work/stage16/README.md).

## Complete named-export inventory

The table includes all 66 named PE exports. The header-declared surface is
marked **H**. All other function exports are **E** (exported, not declared in
`vt_eng.h`). Data exports are **D**. “Static” means wrapper pseudocode or a
shared-address alias was inspected, but no API-level runtime behavior is
claimed. Several E functions are likely internal-support interfaces; their
presence alone does not establish that they are supported application APIs.

| Ordinal | Export | Kind | Coverage |
| ---: | --- | :---: | --- |
| 1 | `VTDTTS_MakeInfo_ENG` | E | **Runtime:** with Stage 16 text, speaker 1, other scalar args -1, and a unique second-argument prefix, raw EAX was `1`; emitted a `.bin.dtt`/`.asc.dtt` pair. Headers identify `VTDTTS BINARY` and `VTDTTS ASCII`; the ASCII file lists phone records and fields such as file index, PCM/PM offsets and sizes, phone mode, pitch endpoints, and rates. The generated basename did not directly match the supplied prefix. Static path normalizes invalid speaker slots to slot 1 and checks loaded/nonempty text. Exact naming, field layout, and values beyond this input remain open. |
| 2 | `VT_AllocSyncInfo_New_ENG` | E | Runtime: allocated two non-null sync-info objects; copy and free calls completed. |
| 3 | `VT_CheckLicense_ENG` | E | Runtime: null-path call returned `0`; no supplied-file contents were captured. Static: `FUN_10028db0` selects an in-memory transform when input plus positive length are supplied, otherwise a file reader. For inputs at least 7 bytes, both helpers copy unchanged unless the first and last three bytes match; on a match they strip a 3-byte prefix, subtract the following seed byte from each payload byte, and discard the 3-byte suffix. The parser splits resulting text at semicolons into an internal list; positional numeric validators enforce observed upper bounds (`99,999,999` and `1024`), and another formats system date as `YYYYMMDD` for comparison with a parsed number. The checker then searches for `<vw_verify ` through the next `/>` and checks `os`, `lang`, `speaker`, `version`, `dbaccess`, `sampling`, and `dbsize`. Token meanings, any integrity/authentication role of the parsed fields, and supplied-file acceptance remain open. |
| 4 | `VT_CheckUserDict_SourceNorm_ENG` | E | Runtime: call with `example` completed and left the input unchanged; wrapper's normalized scratch output is not exposed. Separately, private normalizer `FUN_1005f2e0` returned length 12 and produced `vtspeakprobe` unchanged. |
| 5 | `VT_CheckUserDict_TargetNorm_ENG` | E | Runtime: inputs `example` and `hello` each returned `1`; semantics beyond these samples are open. |
| 6 | `VT_CheckUserDict_TargetPhon_ENG` | E | Runtime: `example` returned `-9`; `HH` returned `1`. Private converter `FUN_1005f710` returned 1 for `HH` and produced a one-character quote string. Semantics beyond these samples are open. |
| 7 | `VT_CopySyncInfo_New_ENG` | E | Runtime: copied scalar fields 4 and 6 (`0x13579bdf`, `0x2468ace0`) between API-allocated objects; the full structure contract remains open. |
| 8 | `VT_CsvParser_Exit_ENG` | E | Runtime: freed the non-null object returned by `VT_CsvParser_Init_ENG`. |
| 9 | `VT_CsvParser_GetField_ENG` | E | Runtime: returned three parsed fields; index 3 returned null. |
| 10 | `VT_CsvParser_GetNfields_ENG` | E | Runtime: returned field count 3 for the tested record. |
| 11 | `VT_CsvParser_Init_ENG` | E | Runtime: returned a non-null parser object, then paired with `VT_CsvParser_Exit_ENG`. |
| 12 | `VT_CsvParser_IsCsv_ENG` | E | Runtime: returned 1 for expected count 3 and 0 for expected count 4 on the tested record. |
| 13 | `VT_CsvParser_MakeCsv_ENG` | E | Runtime: serialized `A` and `b,c` as `"A","b,c"`; return 1. |
| 14 | `VT_CsvParser_Parsing_ENG` | E | Runtime: quoted-comma row returned 3 fields; an unclosed quote was accepted and stripped; an empty middle field was preserved; one embedded LF was absent from the returned field. |
| 15 | `VT_DestroyWindow_ENG` | E | Static; runtime unprobed |
| 16 | `VT_FreeSyncInfo_New_ENG` | E | Runtime: freed both objects returned by `VT_AllocSyncInfo_New_ENG`. |
| 17 | `VT_GetCommaPause_ENG` | E | Runtime via API probe; prototype is absent from header |
| 18 | `VT_GetDBSize_ENG` | E | Runtime for slots 0–5 after Paul M16 load: slot 1 returned `1` and `508121688`; other slots returned `-1` and left zero-initialized output unchanged. |
| 19 | `VT_GetDefUserDictAbsoluteName_ENG` | E | Runtime after load: `../data-common/userdict/userdict_eng.csv` (the export name is `Absolute`, but the observed path contains `../`). |
| 20 | `VT_GetDefUserDictRelativeName_ENG` | E | Runtime after load: `data-common/userdict/userdict_eng.csv`. |
| 21 | `VT_GetDefVersion_ENG` | E | Runtime: returned `Paul-M16-FileIO` in the loaded Paul process. |
| 22 | `VT_GetLicenseComment_ENG` | E | Runtime: null license path returned length `322`; output redacted because it contains host/license-specific fields. |
| 23 | `VT_GetLicenseInfo_ENG` | E | Runtime: null-path size queries for selectors 0/12 returned 4; selector 0 capacities 0/3 returned `-5`; capacity 4 with null output returned `-4`. Static: selector 3 copies raw `<vw_verify>` attributes; selectors 4–11/13–15 map to `os`, `lang`, `speaker`, `version`, `dbaccess`, `sampling`, `app`, `wavsave`, `bgaudio`, `dbsize`, `realtime`; 12 parses `savetime` to four bytes; 1/2 expose two unnamed parser strings. No license contents were captured. |
| 24 | `VT_GetPathKey_ENG` | E | Runtime: queried slots 0–5; returned `SOFTWARE\\VW\\VT\\<name>\\M16` keys. |
| 25 | `VT_GetPitchSpeedVolumePause_ENG` | E | Runtime via configuration probe; prototype is absent from header |
| 26 | `VT_GetSpeakerName_ENG` | E | Runtime: queried slots 0–5 and invalid slot -1; see the export-query capture. |
| 27 | `VT_GetTTSInfo_ENG` | H | Runtime; see header-declared coverage |
| 28 | `VT_GetUserDictLimit_ENG` | E | Runtime: selectors 0–5 return `30, 10, 50, 65, 65, -1` |
| 29 | `VT_INIT_ENG` | E | Runtime: direct call returned `-1` in the loaded Paul process. |
| 30 | `VT_InitSyncInfo_New_ENG` | E | Runtime on an API-allocated object: fields 1 and 2 became `600` and `65`; fields 4, 6, and 13 became `-1`. |
| 31 | `VT_LOADTTS_ENG` | H | Runtime; see header-declared coverage |
| 32 | `VT_LOADTTS_EXT_ENG` | E | Static loader pseudocode and wrapper path; direct call unprobed |
| 33 | `VT_LOAD_UserDict_ENG` | H | Runtime selected errors; see header-declared coverage |
| 34 | `VT_LOAD_UserDict_EXT_ENG` | E | Static: delegates parsing to `FUN_1005e330`; it accepts parser field counts 3 or 4, requires the third field to be one character `A` or `P`, normalizes the source, then uses target normalization for `A` or phoneme validation/conversion for `P`. Runtime candidates have not loaded successfully. |
| 35 | `VT_PAUSETTS_ENG` | H | Runtime call returned with active null-sink handle |
| 36 | `VT_PLAYTTS_ENG` | H | Runtime; see header-declared coverage |
| 37 | `VT_RESTARTTTS_ENG` | H | Runtime call returned with active null-sink handle |
| 38 | `VT_STOPTTS_ENG` | H | Partial runtime; see header-declared coverage |
| 39 | `VT_SetCommaPause_ENG` | H | Runtime; see header-declared coverage |
| 40 | `VT_SetDecimal0Pron_ENG` | E | Static shared-address alias to a no-op return; same implementation as ordinals 44 and 49. |
| 41 | `VT_SetEmphasisFactor_ENG` | E | Runtime: slot 1 clamps `200` to `95` and `-200` to `-95`; restored prior value. |
| 42 | `VT_SetEnglishReadingRule_KOR` | E | Runtime: state accepted a test of `-1` as `0`; restored prior value. |
| 43 | `VT_SetParenthesisCharNumber_ENG` | E | Runtime: state accepted a test of `-1` as `0`; restored prior value. |
| 44 | `VT_SetPhone0Pron_ENG` | E | Static shared-address alias to the same no-op return at ordinal 40. |
| 45 | `VT_SetPitchSpeedVolumePause_ENG` | H | Runtime; see header-declared coverage |
| 46 | `VT_SetSoundCardID_ENG` | E | Runtime: set a temporary value, returned prior sentinel `0xffffffff`, then restored it. |
| 47 | `VT_SetTextTypeForHighlight_ENG` | E | Runtime: nonzero input `7` stored as `1`; restored prior value. |
| 48 | `VT_SetUnitSelectHistoryMode_ENG` | E | Runtime: before load, calls with 0/1 stored 0/1; after load, both calls were ignored. Format-4 control synthesis returned `1` in both modes and produced identical WAVE SHA-256 `a9bb244d9d0cdb664a7a64d14eeb2acd0c45b22d19383d88ff157ba337dd1a69`. Static: flag gates extra per-item arrays/counters; other-input effects remain open. |
| 49 | `VT_SetVirtualTagMode_ENG` | E | Static shared-address alias to the same no-op return at ordinal 40. |
| 50 | `VT_SpeakersInfo_ENG` | E | Runtime: queried slots 0–5; returned `6` and copied the name and embedded database path for each slot. |
| 51 | `VT_TextToBufferEX_ENG` | E | Runtime: selectors 0–2 each produced bytes matching ordinary buffer formats 0–2 for the Stage 16 input. First call returned 0 with 23,606 bytes for selector 0 and 11,803 bytes for selectors 1/2; one flag-1 poll returned 1 with length 0. Guard bytes remained unchanged. Selectors -1 and 3 returned -1 with remaining arguments null/zero. Other pointer/scalar argument values, invalid inputs, multi-poll/chunk/cancel behavior, and extra-argument meanings remain open. |
| 52 | `VT_TextToBuffer_ENG` | H | Runtime; see header-declared coverage |
| 53 | `VT_TextToFile_ENG` | H | Runtime; see header-declared coverage |
| 54 | `VT_TextToLipSyncLog_ENG` | E | Runtime after model load, using the Stage 16 text, an empty second string, and remaining numeric arguments -1: raw EAX was 1 for slots 1 and -1 (fallback). Null text produced AX -3; empty text produced AX -4 (raw EAX `0x003efffc`). A nonempty path-like second argument caused the GDB-injected call to terminate the process with an access violation before a result was captured; this single probe does not establish normal API behavior. The decompiled wrapper is void. The meaning of the second string and output/log format remain unknown. |
| 55 | `VT_TextToPcmBuffer_ProgressBar_ENG` | E | Runtime after model load: with the Stage 16 text, thread 0, speaker 1, option values -1, null window/message, and a 60,000-byte guarded buffer, it wrote 23,606 bytes. Captured bytes match ordinary buffer format 0; guard bytes were unchanged. Raw EAX after the void wrapper was 1, not a declared return contract. Long-text progress notifications, callback delivery, errors, and other arguments remain unprobed. |
| 56 | `VT_TextToPreprocessInfoFile_ENG` | E | Runtime after model load: byte flag 0 with all other arguments null/zero completed normally; raw EAX after the void wrapper was 1. With byte flag 1, Stage 16 text, speaker 1, and second argument `Z:/work/stage16/preprocess-flag1-probe.out`, low AX after the void wrapper was `-6`; no file appeared at that path. Static code returns early for flag 0; flag 4 has a special four-resource loop and flag 6 references `test.pcm`. File format, second-argument role, and behavior of other flag values remain open. |
| 57 | `VT_UNLOADTTS_ENG` | H | Runtime; see header-declared coverage |
| 58 | `VT_UNLOADTTS_EXT_ENG` | E | Static: invalid speaker slots normalize to slot 1; stops playback and unloads that slot. If no slots remain loaded, it unloads dictionary slots 0–1023 and tears down shared window/critical-section state. Direct call unprobed. |
| 59 | `VT_UNLOAD_UserDict_ENG` | H | Runtime selected errors; see header-declared coverage |
| 60 | `VT_UNLOAD_UserDict_EXT_ENG` | E | Static: accepts dictionary indexes 0–1023, returns `-3` when its subsystem is uninitialized or the dictionary is in use by a loaded speaker, frees an idle dictionary and returns `1`, or returns `-1` for an empty slot. Header wrapper runtime-tested only on selected errors. |
| 61 | `VT_VerifyTTS_ENG` | E | Runtime after load, trailing args 0/0: valid text on loaded slot 1 and slot -1 returned raw `EAX=1` (invalid slot falls back to 1); valid text on unloaded slot 0 returned `AX=-4`; null text on slot 1 returned `AX=-2`. The `void` wrapper preserves internal status in EAX. Other inputs/results remain open. |
| 62 | `VT_gHeapStartAddress_ENG` | D | RVA `0x000ff11c`; 32-bit runtime reads were `0` before model load, after load, and after unloading speaker 1. Its purpose remains unknown. |
| 63 | `VT_gVersionFirst_ENG` | D | RVA `0x0007d6f8`; 32-bit runtime read: `3`; part of tuple `(3, 11, 7, 1)` matching PE FileVersion/ProductVersion `3.11.7.1`. |
| 64 | `VT_gVersionFourth_ENG` | D | RVA `0x0007d704`; 32-bit runtime read: `1`; part of the tuple matching PE FileVersion/ProductVersion `3.11.7.1`. |
| 65 | `VT_gVersionSecond_ENG` | D | RVA `0x0007d6fc`; 32-bit runtime read: `11`; part of the tuple matching PE FileVersion/ProductVersion `3.11.7.1`. |
| 66 | `VT_gVersionThird_ENG` | D | RVA `0x0007d700`; 32-bit runtime read: `7`; part of the tuple matching PE FileVersion/ProductVersion `3.11.7.1`. |

## Practical meaning

The Lead 6 runtime matrix is broad across core authoring/playback calls and
selected export-only helpers, but it is not exhaustive behavior coverage for
every export. Successful user-dictionary load/use, broader extended-buffer
argument and error behavior, progress callbacks, nonzero-flag preprocessing
output, lip-sync output semantics, MakeInfo output-name mapping and full field
interpretation, and complete
argument/error coverage need separate bounded probes. Static license parsing
now establishes the in-memory/file input branches, conditional framing and
byte-offset transform, post-read semicolon/list parser, selected positional
numeric bounds, a system-date comparison, and the subsequent `<vw_verify>`
attribute checks. The parsed fields' meanings, any integrity/authentication
role, and acceptance of a supplied record remain unverified.
The four version data exports match the DLL's version resource; the heap-start
data export's purpose remains unresolved. Package-wide generalization is a
separate comparison task; this map is grounded in the local Paul M16 DLL.

## Export-query runtime evidence

The follow-up trace calls `VT_GetSpeakerName_ENG`, `VT_SpeakersInfo_ENG`,
`VT_GetDefVersion_ENG`, and `VT_GetDBSize_ENG` from the already-loaded Paul
process. It queries the six numeric speaker slots; `VT_GetSpeakerName_ENG`
also receives `-1`, which falls back to `Paul`. Results are captured in
[`export-queries-api.log`](../../tools/revkit/work/stage16/export-queries-api.log)
and reproduced by
[`run-export-queries.sh`](../../tools/revkit/work/stage16/run-export-queries.sh).

| Export | Observed result | Boundary |
| --- | --- | --- |
| `VT_GetSpeakerName_ENG` | Slots 0–5 returned `Kate`, `Paul`, `em001`, `Julie`, `James`, and `Ashley`; slot -1 returned `Paul`. | Names are compiled into this DLL; this does not show that the corresponding model packages load. |
| `VT_SpeakersInfo_ENG` | Slots 0–5 returned `6` and copied lowercase IDs plus `d:/eng/db/.../pcm/` path strings. | Those strings are DLL metadata, not paths that were opened in this run. |
| `VT_GetDefVersion_ENG` | Returned `Paul-M16-FileIO`. | This is the function's runtime string in this DLL/package only. |
| `VT_GetDBSize_ENG` | Slot 1 returned `1` and `508121688`; slots 0 and 2–5 returned `-1` and left a zero-initialized output value unchanged. | Only the loaded slot is established; no interpretation of the reported byte count beyond the function name is inferred. |

The same Stage 16 helper trace queried `VT_GetPathKey_ENG` for slots 0–5 and
called both default user-dictionary name exports. Their exact returned strings
are recorded in `helper-exports-api.log`; the absolute-name export returned a
path retaining a `../` prefix. It also parsed
`alpha,"beta,gamma",delta` as three fields, serialized `A` and `b,c` as
`"A","b,c"`, and classified the same parsed input as CSV for expected field
count 3 but not 4. This covers quoted-comma handling and selected helper
contracts only; malformed, multiline, encoding, and capacity boundaries remain
open.

## Static export coverage and unresolved runtime surfaces

The checked-in Ghidra pseudocode report provides wrapper-level static
characterization for the CSV parser family (`Init`, `Parsing`, `Exit`,
`GetNfields`, `GetField`, `MakeCsv`, `IsCsv`), synchronization-info helpers,
license and loader paths, user-dictionary wrappers, file/buffer variants,
playback controls, configuration setters, and information helpers. Static
pseudocode covers all 61 function exports, including shared-address aliases.
The five data exports are now read at runtime; only the heap-start variable's
purpose remains unresolved. Static review is not
equivalent to runtime behavior coverage. Several wrappers only delegate
to internal functions, and their full contracts, pointer requirements,
error handling, or concurrency behavior remain unverified.

Remaining coverage work:

- CSV malformed quoting, multiline records, encodings, and output capacities;
  other sync-info fields and full copy semantics.
- `VTDTTS_MakeInfo_ENG` output-name mapping and complete file-layout decoding.
- `VT_TextToPreprocessInfoFile_ENG` successful nonzero-flag behavior, file
  semantics, and other flag values; `VT_TextToLipSyncLog_ENG` second-argument
  meaning and output semantics.
- License positional token meanings, any integrity/authentication role of the
  parsed fields, and acceptance of a supplied record.
- Broader user-dictionary validation, one successful load/use lifecycle, and
  resulting synthesis behavior.
- Other `VT_VerifyTTS_ENG` arguments/results; additional
  `VT_TextToBufferEX_ENG` arguments, errors, chunking, and cancellation.
- Whether unit-selection history affects other text, settings, or
  repeatability; the purpose of `VT_gHeapStartAddress_ENG` (it read zero before
  load, after load, and after unloading speaker 1).

The version data exports match
the PE version resource, but no internal code references to these data exports
were found in the reviewed pseudocode/disassembly, so external-consumer
semantics remain unverified. `VT_SetDecimal0Pron_ENG`,
`VT_SetPhone0Pron_ENG`, and
`VT_SetVirtualTagMode_ENG` resolve to the same no-op routine in the PE export
table and reviewed pseudocode.

“Exhaustive” here means every named export is inventoried and each row states
the evidence level. It does not claim exhaustive branch, argument,
concurrency, or compatibility coverage.

## Sources

- `binary/vt_pau.dll` PE export directory, read with `objdump -p`.
- `include/vt_eng.h` public declarations and constants.
- `tools/revkit/work/reports/vt_pau-exported-api.c` Ghidra pseudocode for
  exported functions. Treat this as decompiler output, not recovered source.
- [Lead 6 API behavior](lead6-file-api-behavior-2026-09-25.md) and Stage 16
  captures for direct runtime evidence.
- `tools/revkit/work/stage16/export-queries-api.log` for the additional
  metadata and database-size query calls.
- `tools/revkit/work/stage16/helper-exports-api.log` for safe helper, setter,
  allocation-lifecycle, and exported-data observations. License output is
  intentionally omitted from the log.
