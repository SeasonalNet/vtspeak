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
whose thin Ghidra wrappers are rendered as `void`. Runtime probes capture the
low AX status at those wrapper calls; other branches still need direct
return-boundary evidence.

## Header-declared API coverage

| API | Runtime evidence | Remaining coverage |
| --- | --- | --- |
| `VT_LOADTTS_ENG` | **Runtime:** existing Stage 1–2 and Stage 9 Wine traces load the local Paul model and reach synthesis; prior stages capture model-load failures. | Not an exhaustive input/error matrix. Non-null license variants, every speaker slot, concurrent/repeated load, and wrapper return ABI across all failures are not covered here. |
| `VT_UNLOADTTS_ENG` | **Runtime:** unloading speaker 1 in-process makes subsequent file, buffer, and play calls return their declared database-not-loaded codes. | Other speaker IDs, unload while active playback/streaming, repeated unload, and full last-speaker DLL teardown are not characterized. |
| `VT_LOAD_UserDict_ENG` | **Runtime:** indexes `-1` and `1024` expose low AX `-1`. With engine-heap-backed strings, plain three-field rows `hello,HH,P` and `hello,world,A` both load at indexes 27/28 with low AX `1`. Earlier `-3` candidate attempts used suspect literal-string pointers and do not establish row rejection. Static parser accepts 3 or 4 fields, requires field 2 to be one-character `A` or `P`, normalizes field 0, then checks/converts field 1 according to the type. The default `data-common/userdict/userdict_eng.csv` remains empty. | Busy/uninitialized states, malformed-row boundaries, four-field acceptance, and unknown errors remain unprobed. Runtime synthesis effect is conditional on a per-speaker gate; in this Paul process, the status-only license check returned `-11` at its final `dbsize` predicate, leaving capacity 1 and gate 0. The supplied `dbsize=300` is below the established minimum 452 for the 496,212 KiB Paul database. |
| `VT_UNLOAD_UserDict_ENG` | **Runtime:** low AX `-2` for index `-1`, `-1` for empty slot 0, and `1` after unloading each successfully loaded dictionary at indexes 27/28. Static extended-unloader code returns `-3` for uninitialized/in-use dictionaries. | In-use/busy and uninitialized runtime cases, plus unknown errors, remain unprobed. |
| `VT_PLAYTTS_ENG` | **Runtime:** null text (`-2`), empty text (`-3`), database unloaded (`-4`), device-init failure (`-5`), and success (`1`) through an ALSA null sink. On success a WinMM handle was open and play state was 1. | Thread-create and unknown errors, completion message delivery, audible device output, repeated/overlapping play, and speaker selection are unverified. |
| `VT_STOPTTS_ENG` | **Partial runtime:** called while a null-sink playback handle was active; the inferior exited normally during GDB evaluation of this void call. | No call-return marker or explicit post-stop state was captured; effective stop/reset/close behavior remains static-only. |
| `VT_PAUSETTS_ENG` | **Runtime:** call returned while null-sink playback had an active handle. | No sampled playback state or paused audio data was collected; behavior on a physical device remains open. |
| `VT_RESTARTTTS_ENG` | **Runtime:** call returned after pause while the null-sink handle was active. | Resume timing, completed/stopped state, and physical-device behavior remain open. |
| `VT_TextToFile_ENG` | **Runtime:** selectors 0–5 and 7–9 succeed; selector 6 and out-of-range 10 fail. Null/empty text, null path, database-unloaded, repeat, text-format, configuration, and one VTML substitution case are covered. | Non-null unwritable paths, create-thread and unknown errors, full markup/error grammar, all speaker/dictionary combinations, and arbitrary call interleavings are not covered. |
| `VT_TextToBuffer_ENG` | **Runtime:** formats 0–3, null/empty arguments, invalid format, database-unloaded, short-call length semantics, 60,000-byte chunk drain, busy/cancel/no-context transitions, and one file → buffer → file sequence are covered on thread ID 0. | Nonzero thread IDs fail in this Wine setup; other hosts, abnormal/unknown errors, physical-capacity failures, and all format-specific long streams remain open. |
| `VT_SetPitchSpeedVolumePause_ENG` | **Runtime:** setter bounds and getter values checked; isolated pitch, speed, and volume extrema change the fixed sample's WAVE bytes. Pause extrema leave that sample unchanged. | Other slots, negative-value no-op cases, all intermediate values, and pause-bearing text are not characterized. |
| `VT_SetCommaPause_ENG` | **Runtime:** getter values confirm clamp to 0–65,535; tested extrema leave the fixed sample's WAVE bytes unchanged. | Pause-bearing text, intermediate values, and slots other than speaker 1 remain open. |
| `VT_GetTTSInfo_ENG` | **Runtime:** declared requests 0–26, invalid request 27, null value, string-size boundaries, and the null-license/default path are covered. Requests 1 and 2 also ran with the supplied verification file: API return 0 and output values 0 (license check) and 6 (max channel). With invalid path `License`, the API still returned 0 while request 1 wrote -1 and request 2 fell back to 1. | `VT_INFO_ERROR_NOT_SUPPORTED_REQUEST`, unknown error, and exact destination-size edges for requests beyond the existing string cases remain open. |

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
| 1 | `VTDTTS_MakeInfo_ENG` | E | **Runtime + cross-checked layout and controlled variation:** the export emits `.bin.dtt`/`.asc.dtt`; the second argument is the literal prefix and the DLL appends those suffixes. Phone rows map to `unit-*.idx` and their `.dat`/`.upm` spans: mode 0 selects a full decoded unit and combined UPM span; modes 1 and 2 select the indexed first and second side spans. Mode 2 alone carries `Shift Size`; its writer computes a 16-bit subtraction from a selected-unit descriptor field (`+0x0c`) and timeline field (`-0x05`), but the physical meaning remains unknown. `TypeFlag=2` is the detailed phone/unit row. `TypeFlag=1` is now observed between phone rows for `Hello, world.` (`Size=3200`) and `Hello. World.`/`Hello... World.` (`Size=14800`); it is emitted as a size-only row, and static synthesis handling zero-fills that interval. This establishes its role as an inserted silence interval for these punctuation cases. The synthesis loop treats `Size` as a PCM16 sample count and zeroes twice that many bytes. `File Index` is zero-based in header order for tested rows: 0=`merged-gen`, 1=`merged-num`, 2=`merged-etc`; index 3 (`merged-alp`) has not been observed in a row. Rate captures map `Pitch_rate` directly to pitch, `Volume_rate` directly to volume, and `Duration_rate` to integer inverse speed `((speed >> 1)+10000)/speed` (120 gives 83); the tested pause setting did not change them. Text and rate changes retained the same field structure in tested cases, except punctuation adds the silence row. The output is a selected synthesis-timeline manifest; no local downstream reader was found, and no consumer purpose beyond that serialized role is established. The preceding header value `3` remains unlabeled. Invalid speaker slots normalize to slot 1; loaded/nonempty-text checks are static. See the behavior report and Stage 16 captures. |
| 2 | `VT_AllocSyncInfo_New_ENG` | E | **Runtime plus static layout:** allocates a 14-dword/56-byte header, 600 rows of 36 bytes, and a distinct 520-byte nested allocation per row (65 eight-byte entries). Initialization sets header fields 1/2 to 600/65 and clears the sampled row and nested values; full loop bounds and offsets are from Ghidra pseudocode. Static loader/unloader pseudocode stores this object in per-state offset `+0x47774` and frees it at teardown. |
| 3 | `VT_CheckLicense_ENG` | E | **Runtime:** the supplied 468-byte read-only record returned `0` through both explicit file-path and memory-buffer calls; a framed equivalent also returned `0` through both reader paths. Wrong frame seed/suffix, short input, and each attempted field/tag mutation returned `-2`. After loading Paul, the supplied file with its actual per-voice context returned `-11` at `dbsize`. The comment getter matched the final nonempty colon field exactly. `FUN_10028db0` chooses memory transform when input plus positive length are supplied, otherwise the file reader. Matching first/last three bytes strips the 3-byte wrapper and subtracts the seed from each payload byte; plain inputs pass unchanged. The parser splits semicolon records, scans key/value tokens, and links the `License` value to a second exact-key lookup. Its seven fields match, for this sample, host ID, `VW_VTAPI`, expiry token, channel, user, OS, and verification/comment field. Selector and XML comparisons establish these mappings for this record. Static conversion requires a 96-character token, splits it into two 48-character halves, and converts selected two-character hexadecimal slices into derived comparison values; its algorithm and output meanings remain unknown. The checker queries `os`, `lang`, `speaker`, `version`, `dbaccess`, `sampling`, and `dbsize`, not `app` or `expdate`. Direct helper tests rejected invalid `os`, `version`, `dbaccess`, and `sampling`; absent keys return `-1`, which the caller treats as true. The `speaker` list admits the DLL voice IDs in slots 0, 1, 3, and 4 (Kate, Paul, Julie, James) and rejects slots 2 and 5 (em001, Ashley). For loaded slot 1, `dbsize` values 300/450/451 fail and 452/500 pass; the exported database size is 508,121,688 bytes. Disassembly computes a loaded-size limit of 1,100 KiB per `dbsize` unit. The positional date conversion is not conclusively mapped to XML `expdate`, and `VW_VTAPI` remains unexplained. |
| 4 | `VT_CheckUserDict_SourceNorm_ENG` | E | Runtime: call with `example` completed and left the input unchanged; wrapper's normalized scratch output is not exposed. Separately, private normalizer `FUN_1005f2e0` returned length 12 and produced `vtspeakprobe` unchanged. |
| 5 | `VT_CheckUserDict_TargetNorm_ENG` | E | Runtime: inputs `example` and `hello` each returned `1`; semantics beyond these samples are open. |
| 6 | `VT_CheckUserDict_TargetPhon_ENG` | E | Runtime: `example` returned `-9`; `HH` returned `1`. Private converter `FUN_1005f710` returned 1 for `HH` and produced a one-character quote string. Semantics beyond these samples are open. |
| 7 | `VT_CopySyncInfo_New_ENG` | E | **Runtime plus static copy coverage:** called with native dimensions 600×65. All 11 copied header fields, 4,800 row values, and 78,000 nested values compared equal; all 600 destination nested pointers and the row-array pointer remained distinct. Ghidra pseudocode confirms the corresponding loops and leaves destination allocations in place. |
| 8 | `VT_CsvParser_Exit_ENG` | E | Runtime: freed the non-null object returned by `VT_CsvParser_Init_ENG`. |
| 9 | `VT_CsvParser_GetField_ENG` | E | Runtime: returned the quoted-comma and unclosed-quote fields; index 3 returned null on the three-field sample. Preserved the empty middle field in `a,,c`, omitted the empty final field in `a,`, returned `a"b` for both tested quote forms, and retained LF within the `b<LF>c` field of `a,b<LF>c,d`. |
| 10 | `VT_CsvParser_GetNfields_ENG` | E | Runtime: count 3 for the standard quoted-comma row, `a,,c`, and `a,<LF>b`; count 2 for the unclosed quote, doubled-quote, quote-in-unquoted, and CRLF cases; count 1 for `a,`; count 3 for `a,b<LF>c,d`. |
| 11 | `VT_CsvParser_Init_ENG` | E | Runtime: returned a non-null parser object, then paired with `VT_CsvParser_Exit_ENG`. |
| 12 | `VT_CsvParser_IsCsv_ENG` | E | Runtime: returned 1 for expected count 3 and 0 for expected count 4 on the quoted-comma sample. |
| 13 | `VT_CsvParser_MakeCsv_ENG` | E | Runtime with fields `A` and `b,c`: every capacity 0–64 was swept; 10 is the first successful size, 1–9 return raw low AX -1 with partial/empty NUL-terminated output, and 10–64 return raw low AX 1. Capacity 0 still writes the full CSV and underwrites `buffer-1`. With `a"b` and `plain`, output doubled the embedded quote: `"a""b","plain"`. Capacity initialization leaves `0x31` (`'1'`) in unused bytes and writes a NUL at `buffer+capacity-1`. Null field-array with count 0 returned low AX 1 and an empty string. The `void` wrapper makes low AX a register observation, not a declared C result. Null pointers in other argument combinations, capacities above 64, and negative count remain unprobed. |
| 14 | `VT_CsvParser_Parsing_ENG` | E | Runtime: quoted-comma row returned 3 fields; `a,"unterminated` returned 2 fields with the opening quote removed; `a,,c` preserved an empty middle field; `a,` returned only one field; doubled quotes in `"a""b",c` and a quote in unquoted `a"b,c` both produced `a"b`,`c`. `a,\r\nb` returned `a`,`b`; `a,b<LF>c,d` returned `a`,`b<LF>c`,`d`. The same quoted-comma row returned identical fields with flags 0, 1, and 2. These are individual ASCII observations, not a complete parser contract. |
| 15 | `VT_DestroyWindow_ENG` | E | Static; runtime unprobed |
| 16 | `VT_FreeSyncInfo_New_ENG` | E | Runtime: freed both objects returned by `VT_AllocSyncInfo_New_ENG`; the probe restored the allocator's 600-row/65-entry dimensions before freeing. Static pseudocode frees each row's nested allocation, then the row array and header. |
| 17 | `VT_GetCommaPause_ENG` | E | Runtime via API probe; prototype is absent from header |
| 18 | `VT_GetDBSize_ENG` | E | Runtime for slots 0–5 after Paul M16 load: slot 1 returned `1` and `508121688`; other slots returned `-1` and left zero-initialized output unchanged. |
| 19 | `VT_GetDefUserDictAbsoluteName_ENG` | E | Runtime after load: `../data-common/userdict/userdict_eng.csv` (the export name is `Absolute`, but the observed path contains `../`). |
| 20 | `VT_GetDefUserDictRelativeName_ENG` | E | Runtime after load: `data-common/userdict/userdict_eng.csv`. |
| 21 | `VT_GetDefVersion_ENG` | E | Runtime: returned `Paul-M16-FileIO` in the loaded Paul process. |
| 22 | `VT_GetLicenseComment_ENG` | E | Runtime: supplied file path returned 322 bytes including NUL; its 321 payload bytes matched the final nonempty colon field exactly. Framed-file output matched too. A 0–2048 capacity sweep through default and explicit-file paths first succeeded at 322; capacity 321 returned `-4`. Guard bytes stayed unchanged; record contents are omitted from the capture. |
| 23 | `VT_GetLicenseInfo_ENG` | E | Runtime: all selectors 0–15 succeeded for both supplied-file and in-memory inputs, with byte-for-byte equal outputs. Selector 0 is a four-byte integer matching channel; selector 1 returns the positional `0` string matching XML `expdate`; selector 2 returns host ID; selector 3 returns only tag attribute text. Selectors 4–11 and 13–15 matched their named XML attributes; 12 parses `savetime` as a four-byte integer. Capacity sweeps through both paths: selectors 0/12 first succeed at 4 bytes, selector 1 at 2, selector 2 at 13, and selectors 3–15 at 322; capacity 321 returns `-5` for selectors 3–15. A negative-capacity size query returns 321 for selectors 3–15, versus selected-output lengths for 0–2 and 12. Null output returns `-4` for all selectors, invalid selector `-1`, and guard bytes remain unchanged. The unusual shared requirement matches the comment-field payload length; record values are omitted from the capture. |
| 24 | `VT_GetPathKey_ENG` | E | Runtime: queried slots 0–5; returned `SOFTWARE\\VW\\VT\\<name>\\M16` keys. |
| 25 | `VT_GetPitchSpeedVolumePause_ENG` | E | Runtime via configuration probe; prototype is absent from header |
| 26 | `VT_GetSpeakerName_ENG` | E | Runtime: queried slots 0–5 and invalid slot -1; see the export-query capture. |
| 27 | `VT_GetTTSInfo_ENG` | H | Runtime; see header-declared coverage |
| 28 | `VT_GetUserDictLimit_ENG` | E | Runtime: selectors 0–5 return `30, 10, 50, 65, 65, -1` |
| 29 | `VT_INIT_ENG` | E | Runtime: direct call returned `-1` in the loaded Paul process. |
| 30 | `VT_InitSyncInfo_New_ENG` | E | **Runtime plus static initialization coverage:** sets header field 3 to 0, fields 1/2 to 600/65, fields 4–13 to `[-1,0,-1,0,-1,-1,0,-1,0,-1]`; clears row fields at offsets 0 and 8–32 and both value words of each nested entry while preserving all tested allocation pointers. Runtime sampled rows 0, 1, and 599 and confirmed pointers stayed unchanged. Field meanings are unknown. |
| 31 | `VT_LOADTTS_ENG` | H | Runtime; see header-declared coverage |
| 32 | `VT_LOADTTS_EXT_ENG` | E | Static loader pseudocode and wrapper path; direct call unprobed |
| 33 | `VT_LOAD_UserDict_ENG` | H | Runtime: invalid indexes return low AX `-1`; valid 3-field `P` and `A` files loaded at 27/28 with low AX `1`; see header-declared coverage and Stage 16 lifecycle captures. |
| 34 | `VT_LOAD_UserDict_EXT_ENG` | E | Static wrapper contract plus runtime success through the simple wrapper: validates index range 0–1023, checks subsystem initialization and slot occupancy under a critical section, reserves the slot during parsing, then returns `-3` for an empty parse result or `1` when a dictionary pointer remains. `FUN_1005e330` accepts parser field counts 3 or 4, reads fields 0–2, requires field 2 to be one character `A` or `P`, normalizes the source, then uses target normalization for `A` or phoneme validation/conversion for `P`. Direct extended-buffer loading and state errors remain unprobed. |
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
| 54 | `VT_TextToLipSyncLog_ENG` | E | Runtime after model load, using the Stage 16 text and simple second string `vtspeak-lead6`, raw EAX was 1 and one ASCII report was produced. It contains per-phone labels and integer length values, word-index/text spans, and matching total word/phone lengths (11,803 for this sample). Static code routes the second string through the log-file setup (`FUN_1001df10`); static data contains `length-sync-%s-%s.txt`, and the runtime produced a matching `length-sync-*` file. This supports, but does not prove, treating the string as a naming/tag input. Empty second string also returned raw EAX 1; null text produced AX -3 and empty text AX -4 (`EAX=0x003efffc`). A separate path-like string caused an access violation. The `void` wrapper makes these register values, not declared C returns. Accepted string grammar, exact name mapping, measure units, path-like failure, and option effects remain open. |
| 55 | `VT_TextToPcmBuffer_ProgressBar_ENG` | E | Runtime after model load: with the Stage 16 text, thread 0, speaker 1, option values -1, null window/message, and a 60,000-byte guarded buffer, it wrote 23,606 bytes. Captured bytes match ordinary buffer format 0; guard bytes were unchanged. Raw EAX after the void wrapper was 1, not a declared return contract. Long-text progress notifications, callback delivery, errors, and other arguments remain unprobed. |
| 56 | `VT_TextToPreprocessInfoFile_ENG` | E | The void wrapper delegates to `FUN_1001ef20` at `0x1001ef20`; flag 0 returns before pointer checks. With a valid loaded-Paul tuple, byte flags 1–10, 11, 254, and 255 completed with raw low AX 1 (register observation, not a declared return). Corrected heap-backed calls prove `param2` is used literally for flags 1–5, 7–11, 254, and 255; flag 4 appends `.0`–`.3`, while flag 6 ignores `param2` and writes fixed `test.pcm`. The prior mangled names came from GDB string setup corrupting the pointer before API entry. Outputs now have field-level descriptions: flag 3 emits a source-length dot ruler plus decimal internal phone/control bytes; flag 5 emits their CMU-style phone labels; flag 7 emits source text, per-word surface, structural boundary code, phones, bitmask letters, and inclusive ASCII source spans; flag 10 separates tested punctuation with a space. Flag 6 emits 23,606 bytes of PCM for the baseline sample. Its flag-7 prefix is selected from eight fixed period strings by a `GetTickCount`-seeded generator; identical-input fresh runs varied, so it is not a text alignment measure. Why the filler is included remains unknown. Boundary-class names, metadata-letter meanings, and broader punctuation/text behavior remain open. Byte values 12–253 share the static default branch but were not individually runtime-tested. See the behavior report and Stage 16 captures. |
| 57 | `VT_UNLOADTTS_ENG` | H | Runtime; see header-declared coverage |
| 58 | `VT_UNLOADTTS_EXT_ENG` | E | Static: invalid speaker slots normalize to slot 1; stops playback and unloads that slot. If no slots remain loaded, it unloads dictionary slots 0–1023 and tears down shared window/critical-section state. Direct call unprobed. |
| 59 | `VT_UNLOAD_UserDict_ENG` | H | Runtime: invalid/empty errors plus successful unloads of loaded slots 27/28; see header-declared coverage. |
| 60 | `VT_UNLOAD_UserDict_EXT_ENG` | E | Static: accepts dictionary indexes 0–1023, returns `-3` when its subsystem is uninitialized or the dictionary is in use by a loaded speaker, frees an idle dictionary and returns `1`, or returns `-1` for an empty slot. Idle-success path is also observed through the wrapper; in-use and uninitialized errors remain static-only. |
| 61 | `VT_VerifyTTS_ENG` | E | Runtime after load, trailing args 0/0: valid text on loaded slot 1 and slot -1 returned raw `EAX=1` (invalid slot falls back to 1); valid text on unloaded slot 0 returned `AX=-4`; null text on slot 1 returned `AX=-2`. The `void` wrapper preserves internal status in EAX. Other inputs/results remain open. |
| 62 | `VT_gHeapStartAddress_ENG` | D | RVA `0x000ff11c`; 32-bit runtime reads were `0` before model load, after load, and after unloading speaker 1. Its purpose remains unknown. |
| 63 | `VT_gVersionFirst_ENG` | D | RVA `0x0007d6f8`; 32-bit runtime read: `3`; part of tuple `(3, 11, 7, 1)` matching PE FileVersion/ProductVersion `3.11.7.1`. |
| 64 | `VT_gVersionFourth_ENG` | D | RVA `0x0007d704`; 32-bit runtime read: `1`; part of the tuple matching PE FileVersion/ProductVersion `3.11.7.1`. |
| 65 | `VT_gVersionSecond_ENG` | D | RVA `0x0007d6fc`; 32-bit runtime read: `11`; part of the tuple matching PE FileVersion/ProductVersion `3.11.7.1`. |
| 66 | `VT_gVersionThird_ENG` | D | RVA `0x0007d700`; 32-bit runtime read: `7`; part of the tuple matching PE FileVersion/ProductVersion `3.11.7.1`. |

## Practical meaning

The Lead 6 runtime matrix is broad across core authoring/playback calls and
selected export-only helpers, but it is not exhaustive behavior coverage for
every export. Broader user-dictionary row validation, direct extended-buffer
argument and error behavior, progress callbacks, nonzero-flag preprocessing
output, lip-sync report variation across text/options, MakeInfo's unlabeled
header value, mode-2 shift-field units, broader punctuation-to-silence rules,
and complete argument/error coverage need separate bounded probes. Static license parsing
now establishes the in-memory/file input branches, conditional framing and
byte-offset transform, post-read semicolon/list parser, selected positional
numeric bounds, a system-date comparison, and the subsequent `<vw_verify>`
attribute checks. The supplied record now passes the exported checker through
both file and memory readers; its seven positional fields are mapped for that
sample to host ID, API/product marker, expiry token, channel, user, OS, and the
verification/comment field. The speaker allow-list and loaded database-size
limit were also exercised with the DLL's actual voice slots. Disassembly shows
the 96-character token is split and converted through a custom consistency
path, but does not assign business meanings to its derived values. The date
conversion mapping, cross-record schema, and exact meaning of the `VW_VTAPI`
marker remain open.
The four version data exports match the DLL's version resource; the heap-start
data export's purpose remains unresolved. Package-wide generalization is a
separate comparison task; this map is grounded in the local Paul M16 DLL.
Sync-info structure and copy dimensions are now characterized, including a
runtime comparison of all copied header, row, and nested values at native
dimensions. Internal pseudocode uses header field 3 as a 600-row write cursor,
row field 0 as the nested-item count, and nested value `0x28` as a formatting
branch selector. The semantic units of the accumulated values, remaining row
fields, null/self-copy behavior, and mismatched-dimension behavior remain
unknown.

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
count 3 but not 4. The MakeCsv capacity probe used the same two fields and
swept every capacity 0–16 with a 64-byte allocation. Capacity 10 was the first
to return raw low AX 1 with the full nine-byte serialization plus NUL; 1–9
returned -1 with partial or empty output, and 10–16 returned 1. Capacity 0
also emitted the full value despite declaring no output space. The disassembly
at `0x10016960` fills the advertised region with byte `0x31`, then writes a NUL
at `buffer + capacity - 1`, which is `buffer - 1` at capacity zero. A
sacrificial prefix byte set to `0x5a` changed to `0` at capacity zero and
remained unchanged at positive capacities. A separate serialization of
`a"b` and `plain` emitted `"a""b","plain"`, confirming that embedded quotes
are doubled. The wrapper is `void`; raw AX is not a declared C return.

Parser edge probes found that doubled quotes in `"a""b",c` and an embedded
quote in unquoted `a"b,c` both yielded fields `a"b` and `c`. `a,` yielded one
field, unlike `a,,c`, which retained an empty middle field. With `a,\r\nb`,
the CRLF at the start of the second field was removed and the result was
`a`,`b`; with `a,b<LF>c,d`, the newline inside the second field remained,
producing `a`,`b<LF>c`,`d`. The identical row `a,"b,c",d` produced the same
three fields under flags 0, 1, and 2. These observations are limited to the
captured inputs. Other quote patterns, line positions, flag values, and
non-ASCII encodings remain open. `MakeCsv` capacities above 16 and null
pointers remain open. Exact calls are in
[`csv-parser-edges-api.log`](../../tools/revkit/work/stage16/csv-parser-edges-api.log)
and [`trace-csv-parser-edges.gdb`](../../tools/revkit/work/stage16/trace-csv-parser-edges.gdb).
`MakeCsv` behavior beyond the tested capacity range also remains open. Its
exact calls are in
[`csv-capacity-api.log`](../../tools/revkit/work/stage16/csv-capacity-api.log)
and [`trace-csv-capacity.gdb`](../../tools/revkit/work/stage16/trace-csv-capacity.gdb).

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

- Other CSV quote patterns and line positions, parser flag values beyond 0–2,
  non-ASCII encodings, `MakeCsv` capacities above 64 and other null-pointer
  combinations. Remaining sync-info questions are the accumulated values' units,
  other row-field meanings, null/self-copy behavior, and mismatched-dimension
  behavior.
- `VTDTTS_MakeInfo_ENG` still has open boundaries: the binary header's leading `3`, whether `merged-alp` can be selected in this build, the exact physical interpretation of mode-2 `Shift Size`, punctuation-to-silence duration rules beyond the three cases, options/text/settings beyond the controlled examples, and error behavior. No supplied executable imports the export and no local DTT reader was found; the intended external consumer therefore remains unidentified.
- `VT_TextToPreprocessInfoFile_ENG` higher-level names for flag 7's structural
  boundary classes and metadata letters, the purpose of its tick-count-selected
  leading dot run, flag 3 values outside the matched phone/control subset,
  punctuation and markup beyond the controlled ASCII cases, changed synthesis
  settings, malformed nonempty text, invalid scalar combinations, and runtime
  calls for byte values 12–253 (static dispatch routes them to the same default
  mode). Corrected filename behavior and the directly observed record formats
  are documented in the behavior report. `VT_TextToLipSyncLog_ENG` naming rules,
  measure units, path-like failure, and option effects.
- License derived-value meanings and any integrity/authentication role of the
  96-character token; relationship between the positional expiry token and
  XML `expdate`; exact meaning of `VW_VTAPI`; and cross-record validation.
- Broader user-dictionary validation (including malformed/four-field rows),
  direct extended-buffer loading, and live in-use/uninitialized errors. The
  ordinary Paul gate is now explained: the checker returned `-11` for the
  final `dbsize` predicate, leaving capacity 1 and gate 0 because the supplied
  value 300 is below the measured minimum 452. A disposable debugger run with
  the gate forced to 1 made the known `P` and `A` outputs differ from control.
  The repository's default English CSV is zero bytes.
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
