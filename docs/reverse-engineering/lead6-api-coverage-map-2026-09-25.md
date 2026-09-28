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
| `VT_LOADTTS_ENG` | **Runtime:** existing Stage 1–2 and Stage 9 Wine traces load the local Paul model and reach synthesis. Direct tracing of the extended implementation confirms the null-path default, a successful explicit absolute base path, and loading with the supplied license record by file path and memory buffer. | Not an exhaustive input/error matrix. Other license contents/lengths, every speaker slot, concurrent/repeated load, and wrapper return ABI across all failures are not covered here. |
| `VT_UNLOADTTS_ENG` | **Runtime:** unloading speaker 1 in-process makes subsequent file, buffer, and play calls return their declared database-not-loaded codes. | Other speaker IDs, unload while active playback/streaming, repeated unload, and full last-speaker DLL teardown are not characterized. |
| `VT_LOAD_UserDict_ENG` | **Runtime:** indexes `-1` and `1024` expose low AX `-1`. Heap-backed path calls load `hello,HH,P` and `hello,world,A`; a 4th field (`extra` or empty), double-quoted fields, CRLF, an empty source, an unmatched quote in the final type field, an embedded quote in an unquoted source, and a doubled-quote escaped source are accepted in the tested rows. Two/five fields, empty file, empty targets, invalid type `X`, and type `PP` return `-3`; afterward unload returns `-1` for the empty slot and a valid reload at that same index returns `1`. Duplicate load of an occupied slot returns `-2`. With the subsystem flag forced to zero, load returns `-4` and unload returns `-3`; the flag is restored afterward. Static parser code retrieves only fields 0–2; `extra` and empty fourth-field rows produce identical tested synthesis output. Fully quoted ordinary fields and a final unmatched type-field quote both produce the same tested WAV as the plain `P` row. An embedded source quote leaves `hello` unchanged but changes `he"llo` to the same tested WAV as the plain `hello` mapping; the doubled-quote escaped source row produces the same output for `he"llo`. The accepted empty-source row `,HH,P` produced byte-identical WAVs to control on `hello` and `.` with the dictionary gate forced on. Unsupported `P` targets `example` and `NOTAPHONE` terminate the process with access violation rather than returning a loader status. The default dictionary remains empty. | Whether an empty source can match any other token boundary/input, quote-dependent effects on other inputs, other quote placements, true concurrent contention, other invalid phoneme forms, and unknown error coverage remain open. Runtime synthesis effect is conditional on a per-speaker gate; in this Paul process, the status-only license check returned `-11` at its final `dbsize` predicate, leaving capacity 1 and gate 0. The supplied `dbsize=300` is below the established minimum 452 for the 496,212 KiB Paul database. |
| `VT_UNLOAD_UserDict_ENG` | **Runtime:** low AX `-2` for index `-1`, `-1` for empty slot 0, and `1` after unloading each successfully loaded dictionary at indexes 27/28. The extended unloader's forced-uninitialized and injected-reference paths return `-3` (see ordinal 60). | Natural active-use behavior through the simple API and unknown errors remain unprobed. |
| `VT_PLAYTTS_ENG` | **Runtime:** null text (`-2`), empty text (`-3`), database unloaded (`-4`), device-init failure (`-5`), and success (`1`) through an ALSA null sink. On success a WinMM handle was open and play state was 1. | Thread-create and unknown errors, completion message delivery, audible device output, repeated/overlapping play, and speaker selection are unverified. |
| `VT_STOPTTS_ENG` | **Partial runtime:** called while a null-sink playback handle was active; the inferior exited normally during GDB evaluation of this void call. | No call-return marker or explicit post-stop state was captured; effective stop/reset/close behavior remains static-only. |
| `VT_PAUSETTS_ENG` | **Runtime:** call returned while null-sink playback had an active handle. | No sampled playback state or paused audio data was collected; behavior on a physical device remains open. |
| `VT_RESTARTTTS_ENG` | **Runtime:** call returned after pause while the null-sink handle was active. | Resume timing, completed/stopped state, and physical-device behavior remain open. |
| `VT_TextToFile_ENG` | **Runtime:** selectors 0–5 and 7–9 succeed; selector 6 and out-of-range 10 fail. Null/empty text, null path, database-unloaded, repeat, text-format, configuration, and one VTML substitution case are covered. | Non-null unwritable paths, create-thread and unknown errors, full markup/error grammar, all speaker/dictionary combinations, and arbitrary call interleavings are not covered. |
| `VT_TextToBuffer_ENG` | **Runtime:** formats 0–3, null/empty arguments, invalid format, database-unloaded, short-call length semantics, 60,000-byte chunk drain, busy/cancel/no-context transitions, and one file → buffer → file sequence are covered on thread ID 0. | Nonzero thread IDs fail in this Wine setup; other hosts, abnormal/unknown errors, physical-capacity failures, and all format-specific long streams remain open. |
| `VT_SetPitchSpeedVolumePause_ENG` | **Runtime:** setter bounds and getter values checked; isolated pitch, speed, and volume extrema change the fixed sample's WAVE bytes. Pause extrema leave that sample unchanged. | Other slots, negative-value no-op cases, intermediate scalar pause values, and scalar pause effects on synthesis remain open. Inline VTML pause `time` is separately confirmed in Item 51. |
| `VT_SetCommaPause_ENG` | **Runtime:** getter values confirm clamp to 0–65,535; tested extrema leave the fixed sample's WAVE bytes unchanged. | Comma-pause effects on synthesis, intermediate values, and slots other than speaker 1 remain open. |
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
| 2 | `VT_AllocSyncInfo_New_ENG` | E | **Runtime plus static layout and live cursor coverage:** allocates a 14-dword/56-byte header, 600 rows of 36 bytes, and a distinct 520-byte nested allocation per row (65 eight-byte entries). Header dwords map to row-array pointer, row/nested capacities, row write cursor, then output-slice cursors: global start frame; start row, row-frame offset, nested-entry index, nested-frame offset; inclusive global end frame; end row, row-frame offset, nested-entry index, nested-frame offset. `FUN_10021640` advances the start cursor from the prior endpoint; `FUN_100216d0` computes the endpoint from output length and nested frame totals; `FUN_100217e0` consumes the range for EX marker mapping. Live values cross-check the global and nested coordinates. Producer source chains are mapped for row `+20/+24/+28/+32`: `+20/+24` are source-text byte-coordinate fields used inclusively by the formatter, `+28` is the active parser-output record count, and `+32` copies synthesis-record `+0x28`, which indexes parallel synthesis-state arrays; the separate token-loop ordinal is record `+0x10`. The live `(0,10)`, `(13,23)`, `(26,36)`, and `(39,49)` pairs match successive ASCII `Hello world` spans. A new row begins on the first synthesis record or when adjacent records change group index (`+0x28`) or record kind (`+0x27`); the raw class byte at `+0x32` is not that second boundary test. Synthesis-record `+0x2e` is assigned by the producer as the current circular SyncInfo row-slot tag (0–599); on same-kind reuse it compares the current slot with the previous record's tag and refreshes source bounds only when the slot changes. A zero cursor maps to wrapped previous slot 599. The full `DAT_1007daa8` table maps 0x60 raw class-byte indices to nested selectors; ordinary selectors index a 39-entry phone-label table (`AA` through `ZH`), while `0x5a..0x5e` map to the separate silence selector `0x28`. The consumer formats nested durations and row source/frame totals as readable records. Remaining questions are upstream schema/domain names for the group index, kind, and raw class bytes, and any finer phonetic definitions beyond the emitted phone mnemonics. Static loader/unloader pseudocode stores this object in per-state offset `+0x47774` and frees it at teardown. See the behavior report and consumer traces. |
| 3 | `VT_CheckLicense_ENG` | E | **Runtime:** the supplied 468-byte record returned `0` through file and memory readers; equivalent framed input also passed both paths. Wrong frame seed/suffix, short input, and attempted field/tag mutations returned `-2`. `FUN_10028db0` uses the memory transform when nonnull input and positive length are supplied, otherwise reads the file. The framed reader strips the 3-byte wrapper and subtracts its seed; plain input passes unchanged. Parsing splits semicolon records and key/value tokens, and the `License` value is resolved by a second exact-key lookup. Seven fields map for this sample to host ID, `VW_VTAPI`, expiry token, channel, user, OS, and verification/comment field. Static conversion requires a 96-character token, splits it into two 48-character halves, and converts selected hex pairs into derived values whose meanings remain unknown. The checker queries `os`, `lang`, `speaker`, `version`, `dbaccess`, `sampling`, and `dbsize`, not `app` or `expdate`; direct helper tests rejected invalid `os`, `version`, `dbaccess`, and `sampling`, while absent keys return `-1` and are treated as true. The `speaker` list admits slots 0, 1, 3, and 4 (Kate, Paul, Julie, James) and rejects 2 and 5 (em001, Ashley). With Paul loaded, `dbsize=300` returned `-11`; values 300/450/451 fail and 452/500 pass. With James slot 4 loaded, the same file and memory forms both returned `0` at `dbsize=300`: James is 247,848 KiB, under the 330,000-KiB allowance, while Paul is 496,212 KiB. Disassembly computes 1,100 KiB per `dbsize` unit. The positional date conversion is not conclusively mapped to XML `expdate`, and `VW_VTAPI` remains unexplained. |
| 4 | `VT_CheckUserDict_SourceNorm_ENG` | E | Runtime: call with `example` completed and left the input unchanged; wrapper's normalized scratch output is not exposed. Separately, private normalizer `FUN_1005f2e0` returned length 12 and produced `vtspeakprobe` unchanged. |
| 5 | `VT_CheckUserDict_TargetNorm_ENG` | E | Runtime: inputs `example` and `hello` each returned `1`; semantics beyond these samples are open. |
| 6 | `VT_CheckUserDict_TargetPhon_ENG` | E | **Runtime inventory and path coverage:** accepted single-byte tokens are `B D F G K L M N P R S T V W Y Z`; accepted uppercase pairs are `CH DH HH JH NG SH TH ZH`; accepted vowel stems with suffixes `0`–`2` are `AA AE AH AO AW AY EH ER EY IH IY OW OY UH UW`. Sweeps covered 100 selected single-byte inputs, all 676 uppercase pairs, all 2,028 uppercase pairs with suffixes 0–2, and digits 0–9 on each of the 15 accepted vowel stems; every other tested candidate returned `-9`. `#` alone returns 1 as a special token. Null, empty, and whitespace-only inputs return `-1`; valid separated sequences return 1; `HH EXAMPLE` returns `-9`. The exact case-folded `[CI]` marker after a valid phone returns 2 (ASCII case variants, adjacent/separated forms, and trailing whitespace were checked); it is trimmed/mutated in the writable caller buffer. Isolated `[` or `[CI]` returns `-2`; a nonmatching bracket suffix after a valid phone returns `-12`. The post-trim ASCII byte counter accepts 259 bytes and returns `-5` at 261 bytes for the tested valid repeated-token sequence; the earlier apparent outlier came from reusing the mutated buffer and is resolved by a fresh-buffer rerun. Private converter `FUN_1005f710` was called for all 69 accepted spellings plus `#`; it returned 1, its complete byte mapping is in the behavior note, and its 65-token maximum is runtime-checked (66 returns 0 and clears output byte 0). Numeric values are not acoustic labels, and `[CI]` semantics remain unknown. See the detailed validator section and Stage 16 captures. |
| 7 | `VT_CopySyncInfo_New_ENG` | E | **Runtime plus static copy coverage:** called with native dimensions 600×65. All 11 header dwords at `+0xc..+34`, 4,800 row values, and 78,000 nested values compared equal; the two dimensions at `+4/+8` were independently confirmed copied by shape probes. Static pseudocode copies all 13 scalar header dwords (`+4..+34`) while retaining the destination row-array pointer (`+0`); all 600 destination nested pointers and the row-array pointer remained distinct. Runtime edges: `(NULL,NULL)`, `(src,NULL)`, and `(NULL,src)` returned without a fault; self-copy is an identity operation for a valid object; 2×3, 0×3, 2×0, negative-row, and negative-width sources produced source-driven destination dimensions/loops. When destination metadata began at 1×1 and physical allocation remained native-sized, copying 2×3 replaced metadata with 2×3 and copied both rows and all three nested entries. Negative source rows skip row copies; positive rows with negative width copy fixed row data but no nested values. Guard-page probes now confirm a 2×3 source faults on the first write beyond a physically one-row destination (write to `0x1771000`, instruction `0x10026546`) after row 0 is copied and dimensions become 2×3; a 1×2 source faults on the second nested entry beyond an 8-byte destination nested area (write to `0x1771004`, instruction `0x100265b9`) after entry 0 is copied. Static and runtime evidence agree: source dimensions control the loops without destination-capacity checks. See the isolated Stage 16 guard-page captures. |
| 8 | `VT_CsvParser_Exit_ENG` | E | Runtime: freed the non-null object returned by `VT_CsvParser_Init_ENG`. |
| 9 | `VT_CsvParser_GetField_ENG` | E | Runtime: returned the quoted-comma and unclosed-quote fields; index 3 returned null on the three-field sample. Preserved the empty middle field in `a,,c`, omitted the empty final field in `a,`, returned `a"b` for both tested quote forms, and retained LF within the `b<LF>c` field of `a,b<LF>c,d`. |
| 10 | `VT_CsvParser_GetNfields_ENG` | E | Runtime: count 3 for the standard quoted-comma row, `a,,c`, and `a,<LF>b`; count 2 for the unclosed quote, doubled-quote, quote-in-unquoted, and CRLF cases; count 1 for `a,`; count 3 for `a,b<LF>c,d`. |
| 11 | `VT_CsvParser_Init_ENG` | E | Runtime: returned a non-null parser object, then paired with `VT_CsvParser_Exit_ENG`. |
| 12 | `VT_CsvParser_IsCsv_ENG` | E | Runtime: returned 1 for expected count 3 and 0 for expected count 4 on the quoted-comma sample. |
| 13 | `VT_CsvParser_MakeCsv_ENG` | E | Runtime with fields `A` and `b,c`: every capacity 0–64 was swept; 10 is the first successful size, 1–9 return raw low AX -1 with partial/empty NUL-terminated output, and 10–64 return raw low AX 1. Capacities 65, 128, 1024, 4096, 65,536, 131,072, 1,048,576, 4,194,304, 16,777,216, and 67,108,864 returned 1 with physically matched allocations; prefix/post guards remained unchanged and the final advertised byte was NUL. Capacity 0 still writes the full CSV and underwrites `buffer-1`. With `a"b` and `plain`, output doubled the embedded quote: `"a""b","plain"`. Capacity initialization leaves `0x31` (`'1'`) in unused bytes and writes a NUL at `buffer+capacity-1`. Null field-array/count 0 returned low AX 1 and an empty string. Counts -1 and `INT_MIN` with a null field-array also returned low AX 1 and wrote an empty string; only capacity initialization and terminators changed the output. In separate processes, null output, null field-array/count 1, and one null field element/count 1 each terminated with Windows `0xc0000005` before API return. Static flow places these faults in buffer initialization, field-array dereference, and input-string scan respectively. The `void` wrapper makes low AX a register observation, not a declared C result. Negative capacities, positive capacities beyond 64 MiB, overflow, and allocation-failure behavior remain open. |
| 14 | `VT_CsvParser_Parsing_ENG` | E | Runtime: quoted-comma row returned 3 fields; `a,"unterminated` returned 2 fields with the opening quote removed; `a,,c` preserved an empty middle field; `a,` returned only one field; doubled quotes in `"a""b",c` and a quote in unquoted `a"b,c` both produced `a"b`,`c`. `a,\r\nb` returned `a`,`b`; `a,b<LF>c,d` returned `a`,`b<LF>c`,`d`. Static code compares the third argument only with `1`: exactly 1 parses by modifying the caller's input buffer; all other values parse a private copy. Runtime flags INT_MIN, -1, 0, 1, 2, 3, 255, and INT_MAX all returned 1/count 3 with the same fields for `a,"b,c",d`; only flag 1 changed the source string to `a` by overwriting the first delimiter. CP1252 `é/ï` bytes and UTF-8 `é/ï` sequences were preserved byte-for-byte in their fields while ASCII comma split them into two fields. A 41-case quote/empty-field/CR-LF matrix is documented in the behavior report; malformed cases outside that matrix, alternate delimiter configuration, and encodings beyond CP1252/UTF-8 remain open. |
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
| 30 | `VT_InitSyncInfo_New_ENG` | E | **Runtime plus static initialization coverage:** sets header field 3 to 0, fields 1/2 to 600/65, fields 4–13 to `[-1,0,-1,0,-1,-1,0,-1,0,-1]`; clears row fields at offsets 0 and 8–32 and both value words of each nested entry while preserving all tested allocation pointers. Runtime sampled rows 0, 1, and 599 and confirmed pointers stayed unchanged. Live EX records map row offset 8 as audio frames, offsets 12/16 as inclusive source-text positions, and nested entries as frame components plus a 16-bit formatter selector (`0x28` selects silence). Header fields 4–13 are now mapped as start/end timeline cursors by `FUN_10021640`/`FUN_100216d0`; frame-position words initialize to -1 and row/entry index words to 0. State-derived row fields at offsets 20–32 remain opaque beyond their direct source/control roles. |
| 31 | `VT_LOADTTS_ENG` | H | Runtime; see header-declared coverage |
| 32 | `VT_LOADTTS_EXT_ENG` | E | **Runtime + static:** the wrapper enters with slot `-1`, null database path, and extended arguments `(0, 0xffffffff, NULL, 0, 0xffffffff)`; low AX is `0`, and `VT_GetDBSize_ENG` reports `1 / 508,121,688` for Paul. Default path is `../`; explicit `Z:\work\` normalizes to `Z:/work/`. `Z:\voice-data\paul\M16` and its parent return AX `3` (tagger error), so this argument is the resource base, not the voice leaf. With the supplied 468-byte record, file argument 6 and memory argument 7 both load successfully, but per-voice checks differ: Paul slot 1 returns checker `-11`, gate 0, capacity 1; James slot 4 returns checker `0`, gate 1, capacity 6. Both loader calls return AX `0`: the checker result controls the gate/capacity but is not propagated as the loader error. Paul is 496,212 KiB, above the record's `dbsize=300` allowance (330,000 KiB); James is 247,848 KiB, below it. Capacity 6 matches this record's `channel=6`, an inference from one successful record. Static code normalizes invalid slots to 1, maps null/empty paths to the default helper, normalizes nonempty paths, and compares the base with the shared path. Fresh-process first calls for `-2`, `-1`, and `6` each returned `0`, allocating only slot 1 (`0x02`). Repeating each of `-2`, `-1`, `6`, and `1000` with a different base returned `0`; each left Paul queryable and then loaded James successfully from the original base (`occupancy=0x12`). Static code indexes a six-entry per-slot flag array using the original argument on this changed-base branch before using normalized slot 1; observed return/state behavior does not remove the out-of-range read for invalid values. Identical slot-1 reload returned `0` with unchanged state; unload cleared the mask and reload succeeded. A deliberately nonexistent different base with a valid loaded slot returns `1` for Paul alone and for either Paul or James when both are loaded, preserving both states and sizes; the original base then returns `0` for both. With correct read-only roots and aliases, slots 1 (Paul) and 4 (James) load with sizes 508,121,688 and 253,797,273 bytes. Loading both succeeds in either order (`occupancy=0x12`); unloading one preserves the other. Slots 0 and 3 fail with prosody-database error 8 on missing first trees `data-kate/.../tree3/pitch/nbt.tree3` and `data-julie/.../tree3/pitch/nbt.tree3`; their local packages have `tree2`. Slots 2 and 5 request absent trees below `data-em001` and `data-ashley`. Missing-resource runtime probes map errors 3–10. For error 2, GDB-forced NULL results at both the shared 0x2042c-byte and per-speaker 0x4d18-byte allocation sites independently return AX 2; these force null branches rather than natural memory exhaustion. Error 11 is declared by the header but no assignment was found in the recovered loader error-propagation chain. Extended arguments 4/5 are unused; arguments 6–8 feed the license checker and size helper as file path, memory text, and text length. The 514-byte path buffer has no visible copy bound check; long paths remain unprobed. Concurrent threads, malformed paths, and natural allocation failures remain open. Evidence includes `load-ext-resource-error-*-api.log`, `load-ext-error2-static-v3-force-allocation-api.log`, `load-ext-error2-static-v5-force-speaker-api.log`, `load-ext-slot-matrix-alias-api.log`, `load-ext-multi-api.log`, `load-ext-multi-path-switch-api.log`, `load-ext-invalid-switch-recovery-*-api.log`, `load-ext-james-license-api.log`, `load-ext-slot4-filetrace-v3-api.log`, and per-slot traces 0/2/3/5. See also `load-ext-wrapper-api.log`, `load-ext-default-path-api.log`, `load-ext-explicit-runtime-root-api.log`, `load-ext-explicit-db-path-api.log`, `load-ext-explicit-parent-path-api.log`, `load-ext-license-file-state-api.log`, `load-ext-memory-license-state-api.log`, `load-ext-state-api.log`, and `load-ext-path-switch-api.log`. |
| 33 | `VT_LOAD_UserDict_ENG` | H | Runtime: invalid indexes return low AX `-1`; valid three- and four-field `P`/`A` samples, an empty-source row, unmatched final-type quote, embedded source quote, and doubled-quote escaped source load. Forced-gate synthesis shows the unmatched-type form matches the plain `P` row; the embedded/escaped source forms affect `he"llo` but leave `hello` unchanged. The empty-source row has no WAV effect on `hello` or `.`. Malformed field counts, empty targets, and invalid type return `-3` and leave the slot empty; unsupported `P` target samples crash the process. Its target validation is separately mapped under ordinal 6, including unsupported-target crashes in the tested parser path. See the detailed user-dictionary section and Stage 16 captures. |
| 34 | `VT_LOAD_UserDict_EXT_ENG` | E | Static wrapper contract plus direct runtime calls: `(index, filename, buffer, length)`. A nonnull buffer with positive length selects memory parsing even when a valid filename is also supplied; otherwise it uses the filename path. A null filename plus the 10-byte `hello,HH,P` buffer loads, and lengths 10 and 11 succeed; 1/9-byte truncations return `-3`. The occupied-slot guard returns `-2`. Forcing the subsystem flag at `0x100a0458` to zero produces `-4`, after which the original value is restored. `FUN_1005e330` accepts 3/4 fields but calls `VT_CsvParser_GetField_ENG` only for indexes 0–2; with the dictionary gate forced on, no-fourth, `extra`, empty-fourth, and fully quoted three-field `P` rows produce identical WAV hashes. Concurrent reservation races and unknown error cases remain unprobed. Unsupported `P` targets rejected by the exported phoneme checker trigger an access violation in the tested parser path. |
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
| 51 | `VT_TextToBufferEX_ENG` | E | Runtime/static: selectors 0–2 match ordinary buffer formats; 78-byte and long-mark utterances drain through fixed-max chunks (60,000 bytes selector 0; 30,000 selectors 1/2). Negative flag returns that maximum before synthesis validation. `*length = 1` is not capacity; selector 2's first capture was 8,984 bytes, then 30,000 on repeats. Guard pages with 1 writable byte and 4,096 bytes (all selectors), plus 16 bytes (selector 0), all fault on the next protected-page write in the `rep movs` helper while incoming `*length` is 1. Busy (`-7`), cancel (flag 2 → `1`), poll-after-cancel/no-state (`-2`), null/empty text (`-3`/`-4`), null buffer (`-5`), unloaded database (`-6`), invalid selector (`-1`), and nonzero thread (`-2` in Wine) are captured. Optional outputs at positions 8–9 accept null; address 1 at either independently causes `STATUS_ACCESS_VIOLATION` (`0xc0000005`) before return. Position 9 points to a 16-byte descriptor with `0x210`-stride rows. Plain text and `<vtml_sub>` give count 0. `<vtml_mark/>` yields kind-2 rows with empty payload; a nonempty `name`/`NAME` yields kind 1 and a 512-byte name yields kind 3 truncated to 511 bytes. Explicitly empty name emits no row; a one-space name is retained. These forms were checked across selectors 0–2. The row `+8` marker position is mapped through per-output-character arrays. For positions 1/2/5, descriptor `+0/+4` are `0xefa`/`0x23a4`/`0x5d54` and match the selected SyncInfo row `+8`; all SyncInfo row `+8` values sum to 9,191/14,281/29,132 frames, matching PCM16 byte lengths divided by 2 and A-law/μ-law byte lengths directly. A long marker after 51 source bytes yields `+0/+4 = 0x19fac/0x401c` on poll 2 with nonzero header `+0x10 = 0x15f90`, directly confirming `+0 = +4 + +0x10`; the same base advances 30,000 frames per full poll under all selectors. Header start index 10 selects SyncInfo rows 10–11 for the poll-2 marker; `0x16b1 + 0x39d0 - header +0x18 (0x1065) = +4 (0x401c)`, and adding header +0x10 (0x15f90) gives utterance-global +0 (0x19fac). This resolves the nonzero partial-row subtraction. `-8` is a defensive branch when the per-call copy cursor is already at/past its fixed output limit; the ordinary crossing path splits carry-over and loops only below the limit. A 78-byte VTML call completed nine data/terminal polls under each selector without `-8`; unbroken `a` tokens of 1,024/4,096/16,384 bytes drained under selector 1 in 4/13/47 polls with no `-8` (84,789/343,768/1,377,407 bytes total). A debugger-forced conditional-branch probe returned `-8` across selectors 0–2; the original branch bytes were restored, flag-2 cancellation returned `1`, and the next poll returned `-2`, confirming cleanup. The forced branch establishes return/recovery behavior, not natural reachability of the cursor comparison; no valid input has triggered it, and the carry-over unit maximum is not statically bounded. `-9` is absent from this handler (the license export's separate `dbaccess` failure is unrelated). Positions 10–13 are pitch/speed/volume/pause; 14 is dictionary index; 15 is byte text type. Static clamps are pitch 50–200, speed 50–400, volume 0–500; text types 4/6 reset the four preceding options. Pitch/speed and volume affect output; pause 0/250/65535, dictionary 1023, and types 0–7 are unchanged on plain and comma-punctuation fixtures. On the record-bearing VTML fixture, pitch/speed alter marker coordinates/output lengths, while volume 0, pause 250, dictionary 1023, and combined extreme settings with type 4 preserve the observed rows across selectors. Complete drains of identical VTML pause utterances at N=0/200/1000 return 17,762/20,962/33,762 frames under selectors 1/2; selector 0 returns twice the byte count. The 200 and 1,000 ms values add exactly 3,200 and 16,000 frames over N=0, confirming milliseconds at 16 kHz. The supplied callers do not reveal how consumers interpret the descriptor rows. Populated alternate-dictionary effects, scalar pause-argument placement on other text categories, full option cross-products, and natural reachability of the `-8` cursor condition remain open; no valid probe reached the condition, and this handler has no `-9` return. See [detailed EX behavior](lead6-file-api-behavior-2026-09-25.md#vt_texttobufferex_eng). |
| 52 | `VT_TextToBuffer_ENG` | H | Runtime; see header-declared coverage |
| 53 | `VT_TextToFile_ENG` | H | Runtime; see header-declared coverage |
| 54 | `VT_TextToLipSyncLog_ENG` | E | Runtime after model load on `Hello world.`: bytewise heap-backed calls return raw EAX 1 and create reports at exact tested paths: `vtspeak-lipname.txt`, `./vtspeak-lipdot`, `lipsync-probe/rel.txt`, `lipsync-probe/report.txt`, `lipsync-probe/vtspeak-liprelative.txt`, and `Z:/work/stage16/lipsync-filename-bytewise-absolute-z-output.txt`. The supplied path is used as the output path for these extension, relative, subdirectory, and absolute Z-drive cases. Static `FUN_1001df10` builds a candidate `length-sync-%s-%s.txt` name, then returns a null context when its path argument exactly matches the empty global at `0x1009f948`; the runtime sentinel call confirms that behavior and creates no file. For a stable argument, the later branch selecting the already-formatted candidate also requires that same path to equal the empty global, so the earlier return makes the candidate-name branch unreachable through the public API's ordinary call path; nonempty paths select the caller-supplied path, as confirmed by runtime output. This resolves the candidate fallback's reachability for ordinary calls, though the intended use of the dead-looking candidate branch is unknown. Null and empty second strings return raw EAX 1 but create no report in the Stage 5 working directory. At speed 100, eight phone lengths and two word totals sum to 11,803 frames; `VT_TextToBufferEX_ENG` format 0 returns 23,606 PCM bytes = 11,803 16-bit frames. At speed 200, report totals are 5,682 frames, matching the paired buffer call's 11,364 bytes / 2. This establishes frame units and speed sensitivity for this fixture. The older `vtspeak-lead6` report content was retained without its original path. An earlier absolute-path fault came from invalid GDB string setup and is not API evidence. Null text produced AX -3 and empty text AX -4 (`EAX=0x003efffc`). The `void` wrapper makes these register values, not declared C returns. UNC paths and other pathname forms; the invalid pointer value and why the writer receives it on a missing-parent path; graceful write-error handling; and report/audio effects on other text, speakers, and options remain open. Evidence: [`lipsync-filename-bytewise-extension.log`](../../tools/revkit/work/stage16/lipsync-filename-bytewise-extension.log), [`lipsync-filename-bytewise-dot-relative.log`](../../tools/revkit/work/stage16/lipsync-filename-bytewise-dot-relative.log), [`lipsync-filename-bytewise-subdir-relative-verified.log`](../../tools/revkit/work/stage16/lipsync-filename-bytewise-subdir-relative-verified.log), [`lipsync-filename-bytewise-absolute-z.log`](../../tools/revkit/work/stage16/lipsync-filename-bytewise-absolute-z.log), [`lipsync-speed-100.log`](../../tools/revkit/work/stage16/lipsync-speed-100.log), [`lipsync-speed-200.log`](../../tools/revkit/work/stage16/lipsync-speed-200.log), [`lipsync-null-path.log`](../../tools/revkit/work/stage16/lipsync-null-path.log), [`lipsync-sentinel.log`](../../tools/revkit/work/stage16/lipsync-sentinel.log), [`lipsync-filename-bytewise-missing-parent-stop.log`](../../tools/revkit/work/stage16/lipsync-filename-bytewise-missing-parent-stop.log), and [`buffer-ex-0.log`](../../tools/revkit/work/stage16/buffer-ex-0.log). |
| 55 | `VT_TextToPcmBuffer_ProgressBar_ENG` | E | Runtime after model load: with the Stage 16 text, thread 0, speaker 1, option values -1, null window/message, and a 60,000-byte guarded buffer, it wrote 23,606 bytes. Captured bytes match ordinary buffer format 0; guard bytes were unchanged. Raw EAX after the void wrapper was 1, not a declared return contract. Long-text progress notifications, callback delivery, errors, and other arguments remain unprobed. |
| 56 | `VT_TextToPreprocessInfoFile_ENG` | E | The void wrapper delegates to `FUN_1001ef20` at `0x1001ef20`; flag 0 returns before pointer checks. With a valid loaded-Paul tuple, byte flags 1–10, 11, 254, and 255 completed with raw low AX 1 (register observation, not a declared return). Corrected heap-backed calls prove `param2` is used literally for flags 1–5, 7–11, 254, and 255; flag 4 appends `.0`–`.3`, while flag 6 ignores `param2` and writes fixed `test.pcm`. The prior mangled names came from GDB string setup corrupting the pointer before API entry. Outputs now have field-level descriptions: flag 3 emits a source-length dot ruler plus decimal internal phone/control bytes; flag 5 emits their CMU-style phone labels; flag 7 emits source text, per-word surface, structural boundary code, phones, bitmask letters, and inclusive ASCII source spans; flag 10 separates tested punctuation with a space. Flag 6 emits 23,606 bytes of PCM for the baseline sample. Its flag-7 prefix is selected from eight fixed period strings by a `GetTickCount`-seeded generator; identical-input fresh runs varied, so it is not a text alignment measure. Why the filler is included remains unknown. Boundary-class names, metadata-letter meanings, and broader punctuation/text behavior remain open. A one-process sweep called every value 12–253 on the fixed fixture; all returned raw EAX 1 through the void wrapper and wrote empty files. Other texts and fresh-process behavior remain untested. See the behavior report and Stage 16 captures. |
| 57 | `VT_UNLOADTTS_ENG` | H | Runtime; see header-declared coverage |
| 58 | `VT_UNLOADTTS_EXT_ENG` | E | **Runtime + static:** direct call for loaded speaker 1 changes `VT_GetDBSize_ENG(1)` from `1 / 508121688` to `-1` (output remains its zero sentinel); `VT_GetSpeakerName_ENG(1)` still returns `Paul`. Subsequent valid-text calls return `-5` (`VT_TextToFile_ENG`) and `-6` (`VT_TextToBuffer_ENG`, length 0). A separate Paul/James probe loads slots 1 and 4 in both orders (`occupancy=0x12`); unloading either preserves the other slot and size query, and unloading the final slot clears both. Static pseudocode normalizes invalid speaker slots to 1, stops playback, and only when no slots remain loaded unloads dictionary slots 0–1023 and tears down shared window/critical-section state. Invalid-slot runtime behavior remains untested; the void export's return value is not claimed. See `unload-ext-loaded-api.log`, `load-ext-multi-api.log`, and `run-unload-ext-loaded.sh`. |
| 59 | `VT_UNLOAD_UserDict_ENG` | H | Runtime: invalid/empty errors plus successful unloads of loaded slots 27/28; see header-declared coverage. |
| 60 | `VT_UNLOAD_UserDict_EXT_ENG` | E | Static: accepts indexes 0–1023, returns `-3` if uninitialized or if a loaded speaker references the dictionary, frees an idle dictionary and returns `1`, or returns `-1` for an empty slot. Runtime confirms idle success and forced-uninitialized `-3`. A controlled reference-table test inserted a scratch synthesis context whose `+0x1312c0` field equals the loaded dictionary pointer; unload returned `-3`, then returned `1` after the injected table pointer was restored/cleared. This confirms the in-use guard's comparison/return branch, but is not a naturally active synthesis context. Its scan checks each loaded speaker's per-state dictionary-index range `[0, capacity)`; James capacity 6 is runtime-established, while its relationship to `channel=6` remains inferred. |
| 61 | `VT_VerifyTTS_ENG` | E | Static: negative and >=6 speaker slots normalize to slot 1; null/empty text return -2/-3. Runtime after load: valid text on slot 1 and slot -1 returned `EAX=1`; unloaded in-range slot 0 returned `AX=-4`; null text returned `AX=-2`; nonnull empty text had low `AX=-3` (upper EAX was stale); the one-character malformed input `<` returned low `AX=-5`, reaching the final-helper error path. `param3` selects a user-dictionary index: -1 and out-of-range indexes use the default entry; in-range empty indexes fall back to entry 0. Tested indexes -2, -1, 0, 1, 1023, and 1024 all returned 1 for the Stage 16 text. `param4` is the text-type byte also set by the EX buffer API; exhaustive 0–255 calls with `param3=-1` all returned 1 for this text. Because the wrapper is `void`, report low AX for error results; upper EAX bits may be stale. Broader malformed markup, active nonempty dictionaries, and other voices remain open. |
| 62 | `VT_gHeapStartAddress_ENG` | D | RVA `0x000ff11c`; 32-bit runtime reads were `0` before model load, after load, and after unloading speaker 1. Its purpose remains unknown. |
| 63 | `VT_gVersionFirst_ENG` | D | RVA `0x0007d6f8`; 32-bit runtime read: `3`; part of tuple `(3, 11, 7, 1)` matching PE FileVersion/ProductVersion `3.11.7.1`. |
| 64 | `VT_gVersionFourth_ENG` | D | RVA `0x0007d704`; 32-bit runtime read: `1`; part of the tuple matching PE FileVersion/ProductVersion `3.11.7.1`. |
| 65 | `VT_gVersionSecond_ENG` | D | RVA `0x0007d6fc`; 32-bit runtime read: `11`; part of the tuple matching PE FileVersion/ProductVersion `3.11.7.1`. |
| 66 | `VT_gVersionThird_ENG` | D | RVA `0x0007d700`; 32-bit runtime read: `7`; part of the tuple matching PE FileVersion/ProductVersion `3.11.7.1`. |

### Extended loader resource-error probes

The following results come from fresh processes loading Paul slot 1. Each
resource omission used a read-only overlay made of symlinks to the untouched
voice package, with the named file omitted. These results map missing-resource
returns for this DLL and package; they do not establish every malformed-file
case.

| Resource omitted | `VT_LOADTTS_EXT_ENG` low AX | Header code |
| --- | ---: | --- |
| `data-common/dict-eng/atmt.tree3` | 3 | Tagger |
| `data-common/dict-eng/engbi.tree3` | 4 | Break index |
| `data-common/dict-eng/sbd.tree3`, `tppdict_eng`, or `hashidx_eng_tpp` | 5 | TPP dictionary |
| `data-paul/M16/ttsdata/dist_tbl/cepdist.tbl` | 6 | Table |
| `data-paul/M16/dblist.idx` or `mc_idx_tbl/unit-gen.idx` | 7 | Unit index |
| Paul prosody tree (and missing first prosody-tree requests for slots 0, 2, 3, and 5) | 8 | Prosody database |
| `data-paul/M16/dat/merged-gen.dat` | 9 | PCM database |
| `data-paul/M16/dat/merged-gen.upm` | 10 | PM database |

Captures use the name pattern `load-ext-resource-error-<case>-api.log` under
Stage 16; see the [reproduction instructions](../../tools/revkit/work/stage16/README.md).

An empty `dblist.idx` also returns 7. An empty `merged-gen.dat` instead allows
the voice to load: database size falls from 508,121,688 to 161,530,436 bytes,
a decrease of 346,591,252 bytes, exactly the original `merged-gen.dat` file
length. This indicates that the empty file is accepted and its payload size is
excluded from the reported database size. Both error-2 allocation branches
were confirmed by GDB probes. Replacing the shared 0x2042c-byte allocation's
successful return pointer with NULL made the export return low AX 2; the
unmodified control returned 0. Replacing the per-speaker 0x4d18-byte
allocation's successful return pointer with NULL also returned AX 2. These
probes force the post-allocation null condition; they do not simulate natural
allocator failure. Error 11 (`UNKNOWN`) is declared by the public header, but
the recovered extended-loader path passes error values only from its two
initialization helpers; their mapped assignments are allocation error 2 and
resource errors 3–10, with no assignment of 11 found in that chain. This is
scoped to this DLL's recovered loader path, not every build or export. The
separate license-checker result `-11` is not loader error 11: for the supplied
Paul record the checker fails its database-size predicate, the loader leaves
the license gate disabled, and the loader still returns AX 0. Captures and
reproduction instructions are listed in the Stage 16 evidence README.

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
dimensions and smaller source-declared dimensions. Null arguments are
null-guarded, self-copy preserves values, and source dimensions are copied to
the destination. A 2×3 source copies two rows and three nested entries per row;
a 0×3 source copies no rows; and a 2×0 source copies its two fixed-size rows
but no nested entries. Unselected destination tail storage stays untouched.
Internal pseudocode uses header field 3 as a 600-row write cursor and maps the
remaining header scalars to global frame, row, and nested-entry start/end
coordinates. Row field 0 is the nested-item count; nested values are frame
lengths, and `0x28` selects the silence formatting branch. Row `+20/+24`
are source-buffer byte-coordinate fields; `+24` is an inclusive endpoint in
the formatter, corroborated by the safe runtime capture. The row-start test
is mapped: first record, changed adjacent group/index `+0x28`, or changed
adjacent kind byte `+0x27`; it does not compare raw class `+0x32`. Ordinary
selectors index the consumer's 39 phone mnemonics (`AA` through `ZH`), and
selector `0x28` formats as silence. A normal-application trace also captured
41 kind-2 synthesis records and validated live raw-class lookup examples. The
earlier trace read the wrong base for state `+2`; a corrected sequence capture
shows active counts 9 and 6 for the 41- and 25-record batches. Remaining gaps
are the upstream names/domains for group index `+0x28` and kind/class schema.
The runtime fault/corruption outcome when source dimensions exceed physical
destination capacity remains unprobed; static copy loops lack a capacity guard.

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
count 3 but not 4. The MakeCsv probe used the same two fields and swept every
capacity 0–64, then sampled 65, 128, 1,024, 4,096, and 65,536 with guarded
allocations sized to the requested capacity. An additional trace covered
131,072, 1,048,576, 4,194,304, 16,777,216, and 67,108,864 bytes with matched
physical allocations. All returned raw low AX 1 with intact prefix/post guards and a
NUL at the advertised final byte. Capacity 10 was the first
to return raw low AX 1 with the full nine-byte serialization plus NUL; 1–9
returned -1 with partial or empty output, and 10–64 returned 1. Capacity 0
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
three fields under flags 0, 1, and 2. Static code compares the flag only with
1; runtime values INT_MIN, -1, 0, 1, 2, 3, 255, and INT_MAX all returned the
same fields, with only 1 modifying the caller buffer in place. CP1252 and UTF-8
samples preserved their high bytes and split at ASCII comma. Other quote and
line-position patterns, broader encoding behavior, and `MakeCsv` null-pointer
combinations remain open. Negative counts -1 and `INT_MIN` are runtime
confirmed to skip field traversal and produce an empty string. Exact calls are in
[`csv-parser-edges-api.log`](../../tools/revkit/work/stage16/csv-parser-edges-api.log)
and [`trace-csv-parser-edges.gdb`](../../tools/revkit/work/stage16/trace-csv-parser-edges.gdb).
The flag and encoding captures are
[`csv-flag-matrix-api.log`](../../tools/revkit/work/stage16/csv-flag-matrix-api.log)
and [`csv-encoding-edges-api.log`](../../tools/revkit/work/stage16/csv-encoding-edges-api.log).
The large-capacity trace and capture are
[`trace-csv-capacity-large.gdb`](../../tools/revkit/work/stage16/trace-csv-capacity-large.gdb)
and [`csv-capacity-large-api.log`](../../tools/revkit/work/stage16/csv-capacity-large-api.log).
`MakeCsv` error behavior for unprobed combinations remains open. Its original
trace is in
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

- CSV remains open for malformed quote sequences beyond the tested matrix,
  alternate delimiter configuration, and encodings beyond the tested CP1252
  and UTF-8 byte sequences. `MakeCsv` remains uncharacterized for allocation
  failure, arithmetic overflow, and capacities near 32-bit limits. Sync-info
  copy behavior is now directly captured through first-write faults for both
  undersized row arrays and nested-entry arrays; earlier copied data and
  destination dimensions persist before the fault. Header fields map to global
  frame and row/nested-entry cursor
  coordinates. Producer source chains, the row-start predicate, adjacent
  `+0x2e` row-slot assignment/source-span refresh predicate, and complete
  raw-class-to-phone-label lookup are mapped. The upstream schema/domain names
  for the `+0x28` group index, kind, and raw class remain open. The emitted
  phone labels are mapped, though finer phonetic definitions are not. Null/self-copy calls,
  smaller/zero/negative source dimensions, and smaller destination metadata
  have direct runtime coverage.
- `VTDTTS_MakeInfo_ENG` still has open boundaries: the binary header's leading `3`, whether `merged-alp` can be selected in this build, the exact physical interpretation of mode-2 `Shift Size`, punctuation-to-silence duration rules beyond the three cases, options/text/settings beyond the controlled examples, and error behavior. No supplied executable imports the export and no local DTT reader was found; the intended external consumer therefore remains unidentified.
- `VT_TextToPreprocessInfoFile_ENG` higher-level names for flag 7's structural
  boundary classes and metadata letters, the purpose of its tick-count-selected
  leading dot run, flag 3 values outside the matched phone/control subset,
  punctuation and markup beyond the controlled ASCII cases, changed synthesis
  settings, malformed nonempty text, and invalid scalar combinations. The
  fixed-fixture sweep called every byte value 12–253; all 242 wrote empty files
  and left raw EAX 1. Other texts and fresh-process behavior remain untested.
  Corrected filename behavior and the directly observed record formats
  are documented in the behavior report. `VT_TextToLipSyncLog_ENG` has verified
  frame units, literal relative/backslash paths, and absolute paths on mapped
  C: and Z: drives for tested forms. The static path from this export through
  `FUN_10022200` to `FUN_1001e0c0` confirms its output contains phone/silence
  durations and per-row source text/spans with total frames; the emitted phone
  labels are `AA` through `ZH`.
  On `Hello world.`, text types 0–7 and volume 0, pause 250, and dictionary 0
  leave the report byte-identical to the all-default call; pitch 50/200 changes
  frame totals to 12,850/11,998 from 11,803. Types 4/6 reset the tested pitch
  50 when combined with those options. A nonexistent relative parent triggers
  target exit `0xc0000005` during the API call, with no raw EAX captured. Still
  open: other path grammar and write-error handling, the empty-sentinel/default
  name path, and option effects on other text, speaker, and output content. See
  the [detailed option/path probes](lead6-file-api-behavior-2026-09-25.md#vt_texttolipsynclog_eng-option-and-path-probes).
- License derived-value meanings and any integrity/authentication role of the
  96-character token; relationship between the positional expiry token and
  XML `expdate`; exact meaning of `VW_VTAPI`; and cross-record validation.
- Broader user-dictionary validation (including malformed quote/phoneme forms),
  arbitrary buffer-boundary safety, true concurrent reservation contention,
  and natural active-context unload timing. The fourth CSV field is accepted
  but ignored by the recovered loader path for the tested rows; forced
  in-use/uninitialized extended-unloader branches return `-3`. The
  ordinary Paul gate is now explained: the checker returned `-11` for the
  final `dbsize` predicate, leaving capacity 1 and gate 0 because the supplied
  value 300 is below the measured minimum 452. A disposable debugger run with
  the gate forced to 1 made the known `P` and `A` outputs differ from control.
  The repository's default English CSV is zero bytes.
- Broader malformed-markup behavior, active nonempty user-dictionary effects, and other voices for `VT_VerifyTTS_ENG`; the external consumer behavior for the now-decoded EX descriptor rows; and natural reachability of the `-8` cursor comparison from valid or corrupted carry-over. The `VerifyTTS` final-helper `-5` path is reachable with the one-character input `<`. The EX forced-branch probe established its `-8` return and cancellation cleanup; `-9` is absent from the EX handler path.
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
