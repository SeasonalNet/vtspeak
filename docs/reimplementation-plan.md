# Native engine implementation plan

**Status:** Started. Go is the current implementation language; see the
[language decision](reimplementation-language-decision.md). This plan can be
revised as reverse-engineering evidence arrives.

## Scope and evidence boundary

Build an independent native engine that reads supported VoiceText resources and
eventually synthesizes speech. Begin with the documented 2013 M16 Paul package.
Keep vendor binaries and voice data read-only and outside Git. An engine
implementation does not grant distribution rights to those inputs.

The [research roadmap](reverse-engineering/roadmap.md) establishes feasibility
for Paul and identifies behavior to implement. Each implementation milestone
requires its own comparison against captured original-engine output. Existing
Python research scripts remain evidence tools, not runtime dependencies.

## Milestones

1. **Unit DAT boundary — in progress.** Read a versioned Paul unit index's
   little-endian DAT offset and length; decode controls, predictors, history,
   and sample conversion into 16-bit PCM; wrap that PCM in a 16 kHz mono WAV
   for inspection. Compare the Go output with the tracked Stage 2 DLL captures
   across the observed modes 0–3. Reject malformed and unsupported input.
   The `engine/dat` package and `vtdecode` command implement this slice.
2. **Model resources — scaffold started.** `engine/dat` reads versioned index
   span fields, including both UPM-side sample counts, the opaque 21-byte
   feature row, its 7-byte signature, and the combined UPM span. Its
   `UnitRecord.UPMSides` view preserves the shared boundary period on both
   sides. `synthesis.SplitPaul2013UnitUPMSides` maps both views onto decoded
   PCM, retains their original combined-vector offsets, and checks the exact
   side sample spans. `voice.Bank.ReadUnit` validates those index/UPM side
   relationships alongside the combined decoded sample count.
   `engine/voice` opens matching index/DAT/UPM files;
   `Bank.ReadRecord` reads metadata without touching waveform data, while
   `Bank.ReadUnit` returns decoded units after checking cached UPM edge bytes
   and the sample-count/UPM-sum invariant established across the Paul corpus.
   `engine/tree3` parses, evaluates, and loads the 17 observed Paul
   duration/pitch trees. Its `ParseAt` also supports concatenated tree records,
   and `LoadPaul2013ATMTTrees` validates the 27-tree shared English container
   against its count and EOF. `duration.OpenPaul2013` now loads those trees,
   and `Engine.EvaluateATMTTree` exposes selected-tree evaluation over explicit
   numeric inputs without assigning unresolved roles. `engine/distance`
   parses the raw triangular `cepdist.tbl`; three documented runtime lookups match its
   symmetric indexing. `distance.SymmetricFeatureDifference` implements the
   observed feature-distance formula and zero-denominator behavior.
   `voice.OpenPaul2013` groups all four banks, the tree catalog, and the
   validated 1,024-entry metric table plus the generated 256-bin feature
   distance table from a local Paul M16 data root.
   `text.LoadTXT2Tables` reads the nine mapped shared `.txt2` tables, which
   `duration.OpenPaul2013` now retains for native text-table consumers.
   `text.LoadEmbeddedDictionary` validates the indexed embedded lexicon, and
   `text.ParsePhonePayload` parses its direct-ID and alternative-path forms.
   `PhonePayload.BuildPaul2013DictionaryPhoneRowsForState` ports both
   `FUN_1000d450` branches over the parsed payload: it preserves result
   metadata, expands phone rows, builds the 0xff-terminated `d`-joined
   path-control stream, and handles the conditional non-final marker row.
   `PhonePayload.SelectPaul2013PronunciationByPathMarker` joins that payload
   representation to `FUN_10003f10`'s path-row search and aligned phone-string
   return. `text.WritePaul2013DictionaryPhoneRow` writes the recovered
   `0x554` token-result fields, including the conditional marker/sentinel and
   selected-phone fields, while preserving bytes the native writer leaves
   untouched. Token indexes and surfaces plus source-marker and selector/gate
   production remain caller inputs. `duration.Engine.SelectPaul2013ContextPronunciationFromEmbeddedDictionary`
   composes the loaded embedded lookup, payload expansion, and marker
   selection used by `FUN_1000cbe0`, preserving its empty C-string output and
   `-1`/`1` return distinction.
   `text.LoadExceptionDictionary` parses the four sorted category groups in
   `exceptdict`; `duration.Engine.LookupPronunciationException` exposes exact
   lookup for caller-normalized and caller-encoded keys.
   `text.NormalizePaul2013ExceptionSurface` ports the recovered character map,
   apostrophe/hyphen handling, and component counter. The relevant ASCII
   character-attribute mask is provisioned from DLL table `0x1007e188`.
   `text.EncodePaul2013ContextString` ports the ASCII letter expansion,
   delimiter insertion, 64-byte cap, and row-flag bit from `FUN_10009cd0`,
   using the 26 observed expansion records at DLL VA `0x10077e58`;
   `ApplyPaul2013ContextStringToPhoneRow` connects it to the `0x70` row, and
   `BuildPaul2013ContextCodesFromSourceString` composes the result with
   `FUN_10016c90`'s code and per-phone-marker scan. The separate token-state
   arrays and final mode code remain caller inputs.
   Non-ASCII character-table behavior remains unsupported. The decompiler
   evidence is in `tools/revkit/work/reports/lead3-ax-rowflag-static-2026-09-25.txt`.
   `text.ApplyPaul2013LiteralContextHandler` also ports the directly compared
   `FUN_100091b0` contraction and suffix rows (`'s`, `'ve`, `'re`, `'em`,
   `'ll`, `'m`, `'n`, `'d`, `'er`, and `'ee`) with the observed trailing-code
   class table at `0x10077d94`; unlisted strings still require the generic
   dictionary/TPP branch.
   `text.BuildPaul2013ExceptionDispatchGateFromRows` extracts the direct gate
   fields from supplied `0x94` source/parser rows and `0x70` phone-context
   rows. `text.DerivePaul2013PhoneContextPreparationFromRows` also ports the
   complete `FUN_100091b0` low-short return: early-return and non-`Y`/`S` tail
   paths return zero, while `Y`/`S` paths return one after setting `local_c`.
   The counted-row exception resolver derives this gate input for every valid
   row even where cascade row mutations remain unported. These branches and
   return conditions are taken from the Ghidra pseudocode in
   `tools/revkit/work/reports/lead3-ax-metadata-helper-static-2026-09-25.txt`.
   `text.ShouldDispatchPaul2013PronunciationException` and the duration-engine
   wrapper port the `FUN_10007520` gate into `FUN_10008dc0`; the row-based
   `ResolvePaul2013PronunciationExceptionFromRows` joins gate, sequence lookup,
   optional `h'` retry, row-count derivation, and phone-code decoding.
   `text.ApplyPaul2013ModelSourceClassNormalization` ports the ordered
   `FUN_10009030` S/C source-class branches against the actual `FUN_1000ea20`
   model-row offsets, including existing-code flags, metadata-gated generic
   normalization, and context-string fallback. Its duration-engine wrapper
   supplies the recovered `FUN_10002f10` implementation;
   `duration.Engine.NormalizePaul2013ModelContextRows` applies it in an
   explicit caller-supplied order and returns the handled indexes.
   `duration.Engine.RunPaul2013KnownModelContextPasses` composes the supported
   row iteration and matched-exception index adjustment. Other
   `FUN_100091b0` mutations, the generic dictionary/TPP and later
   special-handler chain, and runtime row comparison remain open;
   the field mapping is recorded in
   `tools/revkit/work/reports/stage63-source-class-context-normalizer.c`.
   `text.DerivePaul2013ModelContextPreparation` and
   `EvaluatePaul2013ModelContextExceptionGate` now derive the known
   `FUN_100091b0` low-short result and `FUN_10007520` gate from the full model
   arena. `text.BuildPaul2013NativeExceptionCandidateRows` ports
   `FUN_10008dc0`'s contiguous row scan, stopping at an unnormalizable surface
   or when the accumulated one-to-four component total exceeds the exception
   dictionary range. `duration.Engine.ResolvePaul2013PronunciationExceptionFromModelContextRows`
   composes that producer with the gate, X-prefix retry, exception lookup,
   phone decoding, and `FUN_1000ca50` delimiter-split writes into a copied
   model arena. The preceding normalization mutations and native runtime
   row comparison remain open. The scan evidence is recorded in
   `tools/revkit/work/reports/stage68-exception-candidate-scan.c`; the
   full-model offset
   and write map is documented in
   `tools/revkit/work/reports/stage64-model-context-exception-dispatch.c`;
   the implementation remains static-only without row-for-row runtime parity.
   `text.ApplyPaul2013ModelContextFormY` ports the form-Y parser text copy,
   row-state bit, and adjacent duplicate-parser-row advance at the start of
   `FUN_100091b0`; `text.ApplyPaul2013ModelContextFormS` ports the form-S
   context-string encoding and row flag. Their auxiliary-string suffix-handler
   tail accepts the direct literal cases through the `WithTail` variants when
   callers supply the auxiliary surface; generic suffix dispatch remains open.
   See
   `tools/revkit/work/reports/stage65-model-context-form-y.c`.
   `text.ApplyPaul2013ModelContextSecondPass` also ports the final counted-row
   pass in `FUN_100091b0`: it checks the phone-code class and the directly
   observed previous `the` / following `de` surfaces, then writes `0x26`, NUL,
   and state bit 3 into the preceding context row's phone string and state
   byte. The native following-surface read is preserved against the supplied
   arena, including for the final counted row. The generic
   dictionary/TPP cascade, remaining first-pass mutations, unsupported
   `FUN_10007520` branches, later special handlers, and runtime row comparison
   remain open; see
   `tools/revkit/work/reports/stage66-model-context-second-pass.c`.
   `text.InitializePaul2013ModelContextProcessingFlags` ports the initial
   `FUN_10007520` counted-row loop, setting the 16-bit per-row state to 1 when
   the first phone byte is nonzero and 0 otherwise. The combined model-text
   build-and-normalize entry points now run this pass before source-class
   normalization. It does not imply that the intervening exception and
   special-handler cascade is complete.
   `duration.Engine.RunPaul2013KnownModelContextPasses` composes the known
   model-row order: Y/S writes and duplicate-row skipping, exception lookup
   with matched-token row advance, the native dispatch gate,
   `FUN_10009030` before `FUN_100086c0` and the supported fallback chain, and
   the final counted-row `the`/`de` mutations. An open outer gate with a
   closed exception gate still reaches known post-handler rules; a closed
   outer gate skips both and is reported in `OuterGateSkippedRows`. It accepts
   literal-tail strings as explicit inputs
   and returns rows that still need unsupported suffix branches and TPP row
   updates in `UnresolvedRows`. The
   `BuildAndRunPaul2013KnownModelContextPassesFromParserRows` and
   `BuildAndRunPaul2013KnownModelContextPassesFromOrdinarySource` wrappers
   connect this pass to the existing token/context projection. The row
   `RunPaul2013KnownModelContextPassesWithCompoundContractionAndC3A0Rows`
   now uses the recovered ordered suffix/contraction/B800/C3A0 handlers,
   trailing-apostrophe path, and terminal generic fallback automatically after
   each preceding miss. Historical row-selection arguments remain for
   compatibility and controlled comparisons. The
   `BuildAndRunPaul2013KnownModelContextPassesWithCompoundContractionAndC3A0RowsFromOrdinarySegmentsInSharedArena`
   entry point connects the same pass to shared-arena ordinary-source rows.
   The `...WithCompoundContractionB800AndC3A0Rows...` variant now carries
   compatibility conflict checks through the automatic B800 and
   C3A0 stages.
   `Paul2013CompoundOuterMode` derives FUN_10007520's short for
   FUN_10009dc0 from the preceding counted context row's parser index and
   parser-row type at +0x2c. The counted-row pass now attempts the handler for
   each row that reaches that native branch and derives the mode automatically;
   explicit mode maps remain available for controlled comparisons. Later
   supported A140, contraction, B800, C3A0, trailing-apostrophe, and generic
   stages now run in native order. Unsupported suffix branches, remaining TPP
   behavior, and runtime row parity remain unresolved.
   `ApplyPaul2013FUN10007520ArticleMarker` also ports the `a` post-handler's
   following-class and same-parser-row apostrophe-s cases, including its
   initial-row gate through the recovered class-3 and class-4 tables. Runtime
   row parity remains open; see Stage 122.
   `ApplyPaul2013FUN10007520TheMarker` ports the `the` post-handler's
   initial-row, following-class, and same-parser-row predicates and its
   `0x07`/`0x1e` marker writes. The counted-row pass runs it before the `use`
   rule, matching the recovered order. Marker meaning and runtime row parity
   remain unverified.
   The counted-row runner now honors the `FUN_10007520` dispatch gate and runs
   `FUN_10009030`, parser-class-A `FUN_10002f10` precheck when gated, and the
   `FUN_100086c0` path before the automatic fallback chain. Both direct and
   hyphen-split `FUN_10003110` generic routes are composed. The latter reads
   from `model + 0x429ad + row*0x70`
   (the counted-row surface starts at row offset `+5`); a nonzero native short
   applies the recovered `FUN_10009cd0` encoding and bit-3 write
   before known post-handler rules. The following nonzero phone-marker gate
   also blocks lower-priority fallback handlers. The hyphen handler preserves
   its 64-byte output cap and accumulated flags. These rows remain in
   `UnresolvedRows` because remaining TPP updates and runtime parity are not
   established; unsupported generic-token paths fail with an error. The
   static basis is `tools/revkit/work/reports/stage27-context-normalization-helpers.c`
   and the `FUN_10007520` caller listing in
   `tools/revkit/work/reports/stage6-phone-dictionary-expanded.txt`; the
   split loop and capacity arithmetic are read from address `0x10003110` in
   `tools/revkit/work/reports/vt_pau-objdump-disassembly.txt`. Runtime row
   parity for this composed branch remains unverified.
   Remaining TPP handlers, the row producer, and row-for-row runtime
   comparison remain open. Historical B800/C3A0 row arguments are retained
   for compatibility; the composed pass attempts those handlers after their
   predecessors miss.
   `text.ExtractPaul2013ExceptionRowInputs` and
   `ResolvePaul2013PronunciationExceptionFromPhoneContextRows` read normalized
   surfaces and X retry markers from caller-selected `0x70` rows, then compose
   the counted-row gate and exception resolver. The source of the row
   eligibility/order list and the producer of those rows remain open. The
   outer `FUN_1000cb30` branch order is now ported: a positive context lookup
   count returns its text, while a nonpositive count enters the generic
   normalizer only when its recovered short gate is set. The context lookup
   producer and its dictionary/TPP resource traversal remain caller supplied;
   table-backed literal contraction handlers are ported separately.
   `text.Paul2013FUN10010010` and `text.Paul2013FUN100100D0` port two
   self-contained `FUN_100086c0` callees: the first-byte-table predicate over
   `DAT_1007e188` and the mapped-vowel count over the u16 table at
   `0x1007e388`. Stage 104 ports the signed-index prefixes: `FUN_10010010`
   reads zero attributes for high-bit bytes, while `FUN_100100d0` counts none
   as vowels. `duration.Engine.EvaluatePaul2013FUN100086C0Precheck` now ports
   the caller-gated prefix through the `E`/`A` status branches, existing phone
   and row-state checks, mapped-vowel count, generic eligibility short-circuit,
   and directly recovered alternate-record state checks. It returns an
   explicit `NeedsNeighborScan` disposition where the native path enters its
   previous/next-row search. `duration.ScanPaul2013FUN100086C0Neighbors` now
   scans each direction to its nearest first/second-byte class candidate and
   applies the recovered `FUN_1000ffd0` gate.
   `duration.CompletePaul2013FUN100086C0NeighborDecision` carries the
   post-scan status/state/length branches, including the mapped `Lexus` /
   exact `IS` special case. The strings are at VA/file offsets
   `0x10077874`/`0x77874` and `0x10077870`/`0x77870`; the comparator bodies are
   at `0x1001c2c0` (mapped equality) and `0x1001c1f0` (byte-string equality).
   `DAT_100783d0` holds six 32-bit little-endian pointers (VA/file offset
   `0x100783d0`/`0x783d0`); `DAT_100783e8` is the following 32-bit count
   (VA/file offset `0x100783e8`/`0x783e8`), whose value is 6. In pointer order,
   targets are `ALL` (`0x10078c0c`), `BY` (`0x10078c08`), `DEAR`
   (`0x10078c00`), `DO` (`0x10078bfc`), `IS` (`0x10077870`), and `NOT`
   (`0x10078bf8`). The mode-`0x53` search now uses that recovered table
   directly. The composed entry point also ports the
   `param_5` caller gate: its one-byte/non-`E` exception, mapped `Wi` check at
   `0x1007787c`, and compound character gate. Production of the model arena
   and local `local_25c` bytes, plus row-for-row runtime parity, remain open;
   this is not an end-to-end normalization claim.
   `Engine.EvaluatePaul2013FUN100086C0SupportedPath`
   composes the evidence-backed stages in native order and returns the
   neighbor trace with the low-short disposition. `Paul2013NativeCStringTable`
   implements `FUN_100560a0`'s mode-`0x53` sorted pointer-array search using
   bytewise `strcmp` order; `FUN_100086c0` now uses the directly recovered
   six-string table at `DAT_100783d0`.
   `text.BuildPaul2013NormalizerFeatureWindow` ports the ten-short input layout
   in `FUN_10002db0`, including its seven-character window, source-edge zero
   fill, and signed trailing category bytes. `FUN_10002d70`'s lowercase table
   at DLL RVA `0x7e688` was confirmed to map lowercase ASCII to uppercase
   before subtracting `0x40`. Feature-window construction is available, and
   `duration.Engine.EvaluatePaul2013NormalizerCharacter` connects one supplied
   character/category window to its selected shared ATMT tree.
   `EvaluatePaul2013NormalizerCharacterSequence` performs the supported ASCII
   letter/apostrophe token pass from right to left, feeds each signed category
   byte into the next character window, and applies the observed short-token
   internal-E override to `0x1c`, then restores result order to source order.
   `text.ShapePaul2013NormalizerClassCodes` ports the subsequent
   `FUN_10002200` opaque-byte shaping pass, including class-escape decoding,
   direct punctuation mappings, and comma separator insertion. It uses the
   recovered `0x10077d94` class-membership ranges and is covered by synthetic
   byte-level cases, not a row-for-row native capture.
   `text.RemapPaul2013NormalizerClassCodes` also ports the preceding
   `FUN_100024f0` explicit code cases and low-byte lookup table at DLL RVA
   `0x10077f5a`; its sentinel entries are covered by the native explicit cases.
   `duration.Engine.ClassifyPaul2013NormalizerToken` composes the character
   results with this remapping and shaping pass, while retaining the trace.
   `text.Table.LookupPaul2013CHCFlags` and the duration-engine wrapper connect
   `chc_sort.txt2` to `FUN_10002680`'s exact-key search and four-bit mask gate.
   `MatchPaul2013CHCSubstring` ports `FUN_100026f0`'s exact query, observed
   compound prefix/suffix checks, trailing-S retry, and long-token/short-
   substring fallback against that loaded table. `IsPaul2013NormalizerEligible`
   composes those lookups with `FUN_10002c70`'s ASCII letter/apostrophe scan,
   vowel boundaries, and run-specific masks. Its 31-byte limit follows the
   decompiled local run buffer; no row-for-row native eligibility comparison
   has been made. `NormalizePaul2013GenericToken` now ports the supported
   `FUN_10002f10` branch control after eligibility: one-byte vowel tokens use
   the direct context encoder, while other eligible ASCII tokens produce
   shaped class codes or take the direct context fallback. A result containing
   an observed context-class byte sets the native `0x10` row-flag mask; a
   nonempty direct context fallback carries `0x20`.
   `ApplyPaul2013GenericNormalizerWrites` applies the selected class/context
   output to a copied caller buffer, terminates it with NUL, and ORs the
   returned mask into the separate native flag byte. The `MC` prefix branch is
   also recovered: when the suffix contains a table-recognized vowel, it
   classifies `MAC` plus the uppercased suffix; a failed class lookup still
   encodes the original source token. `duration.Engine.LookupPaul2013EmbeddedContextToken`
   now performs the loaded embedded-key transform, dictionary lookup, payload
   parse, and phone-ID expansion consumed by `FUN_1000cb30`; it supplies the
   parsed pronunciation count, first internal-symbol string, and metadata
   bit-7 fallback gate. `NormalizePaul2013ContextTokenFromEmbeddedDictionary`
   connects this lookup to the outer branch order.
   `Paul2013ContextNormalizationResult.NativeReturnCode` preserves
   `FUN_1000cb30`'s direct short result (`1` handled, `-1` miss), so ordered
   callers can distinguish a generic/dictionary hit from a fallback miss.
   `TransformPaul2013ModelSourceFromEmbeddedDictionary` also feeds the E/A
   result type into `FUN_1000d640`'s terminal punctuation split. The TPP-family lookup,
   broader caller-side token transforms, and source-row mutations remain
   separate and unresolved. `NormalizePaul2013ContextTokenToBuffer` composes
   the supported lookup/generic branch with the native output-buffer writes,
   including the initial no-match NUL, dictionary-text copy, and generic local
   flag byte. `NormalizePaul2013ContextTokenFromEmbeddedDictionaryToBuffer`
   connects those writes to the loaded embedded lexicon for an already-
   normalized token.
   `duration.Engine.NormalizePaul2013TrailingApostrophe` ports
   `FUN_1000c710`'s ordered `in'` rewrite/lookup, apostrophe-stripped lookup,
   generic normalizer, and direct context-encoder fallback, including its
   caller row-flag masks and trailing-period-to-hyphen write. The lookup
   resource remains caller supplied; a wrapper uses the loaded embedded
   dictionary. Static pseudocode coverage is recorded in
   `tools/revkit/work/reports/stage67-trailing-apostrophe-context.c`; no
   row-for-row runtime comparison has been made.
   `duration.Engine.NormalizePaul2013MappedSpellingSuffix` ports all nine
   directly recovered suffix/replacement pairs from `FUN_1000b4d0`, including
   its minimum stem lengths, and calls the loaded `FUN_1000cb30` path, applying
   flag mask 9 only on a nonempty result. Use it only after the earlier
   `FUN_10009dc0` and `FUN_1000a140` handlers miss. Static evidence and the
   `.data` string bytes are recorded in
   `tools/revkit/work/reports/stage69-model-context-suffix-rewrites.txt`;
   runtime row comparison remains open.
   `duration.Engine.NormalizePaul2013AnceSuffixContext`,
   `NormalizePaul2013NessSuffixContext`, and
   `NormalizePaul2013MentSuffixContext` port the 7- and 8-byte `ance`,
   `ness`, and `ment` branches in `FUN_1000a140`. Each normalizes the
   four-byte-truncated stem and sets row flag bit 3 only when the native
   result is nonempty. The `ance` branch appends `0x07 0x2d 0x23 0x37`; the
   `ness` and `ment` branches append `0x2c 0x07 0x2d 0x39`. Other function
   branches remain open; evidence is in the Stage 70 through Stage 73
   reports. `NormalizePaul2013ShipSuffixContext` and
   `NormalizePaul2013LessSuffixContext` also port their directly compared
   7- and 8-byte cases, including the `0x07 0x2d 0x37` and `0x38 0x25 0x35`
   marker bytes respectively. Stage 98 adds the 6- to 8-byte mapped `est`
   branch with its four-byte stem truncation and `0x07 0x2d 0x37` tail. Stage
  99 ports short `ers`, including its ordered truncations, table-gated
  `i`-to-`y` rewrite, and `D`, `0x1a D`, and `0x07 79` output tails. Stage 100
  extends the mapped `ly` handler to 5- to 31-byte inputs, and Stage 101 ports
  the short `ing` path with its guarded table reads, marker-selected lookup,
  fallback, and `#.` tail. Stage 102 extends `ist` to the 5- to 31-byte range
  using its direct three-byte check, suffix-stripped lookup, and `#79` tail.
  Stage 103 ports the intervening `ful` handler and its space/control/plus
  tail. Stage 106 captures the full signed 256-index short-table window used
  by the A140 predicates, including adjacent data reached through negative
  indexes. Stage 107 cross-checks every directly compared suffix family
  against a Go handler. Remaining A140 branches outside those suffix
  comparisons are table driven or state dependent; their complete coverage
  and native runtime row comparison remain open.
   `NormalizePaul2013EdSuffixContext`
   now ports the direct long-input `ed` branch, including marker-aware `%`
   lookup, fallback dictionary normalization, the inspected raw character
   table predicates, and its output marker. Stage 106 captures the signed
   table window used by these predicates; original-runtime row parity remains
   open. See the Stage 75 and Stage 106 reports.
   `NormalizePaul2013ErSuffixContext` ports the bounded long-input `er`
   branch, including its ordered truncation lookups and shared raw table
   predicates. Other cases and original-runtime row parity remain open; see
   the Stage 76 report.
   `NormalizePaul2013SupportedSuffixFallbackAfter09DC0` connects these
   bounded `FUN_1000a140` handlers to the later `FUN_1000b4d0` mapped spelling
   rewrites in recovered order, for callers that have already observed a
   miss from `FUN_10009dc0`. This does not implement that earlier handler or
   the unported `FUN_1000a140` branches. Evidence and the composition
   boundary are recorded in
   `tools/revkit/work/reports/stage77-supported-suffix-fallback-composition.txt`.
   `duration.NormalizePaul2013CompoundContext` ports the `FUN_10009dc0`
   component scanner, explicit single-letter `a` branch, output `d` joins,
   overflow behavior, and final row flag. Component normalization is an
   explicit callback because its model-state dictionary gates and
   `FUN_1000c040`/`FUN_1000c3a0` handlers are not complete. The evidence and
   boundary are recorded in
   `tools/revkit/work/reports/stage78-compound-context-scanner.txt`.
   `Engine.LookupPaul2013CompoundDictionaryGate` ports the preceding
   `FUN_10003b10` payload-flag predicate over the loaded embedded dictionary,
   retaining raw flag and lookup presence. `Paul2013CompoundCharacterGate`
   ports the adjacent ASCII predicate from `FUN_1000ffd0`; the compound
   scanner passes both results to its callback. `NormalizePaul2013CompoundContextWithEmbeddedDictionary`
   now follows those gates with `FUN_1000cb30`, preserving its native return
   code and bit-0 update. `NormalizePaul2013CompoundContextWithSupportedHandlers`
   composes the recovered A140, trailing-apostrophe, contraction-suffix, C3A0,
   and generic branches in order; c040 uses its supported prefix composition
   by default, with an optional caller override.
   `NormalizePaul2013CompoundModelContextRow` writes a match into a counted
   model row. `RunPaul2013KnownModelContextPassesWithCompoundAndC3A0Rows` and
   the shared-arena segment wrapper place explicitly selected compound rows
   after form/exception handling and before C3A0. The `local_14` prior-output
   marker is applied when an accumulated code byte exists; the write before
   an empty output buffer, outer branch selection, and runtime parity remain
   open; see Stages 79, 80, and 110.
   `NormalizePaul2013ContractionSuffixContext` ports the direct `FUN_1000c040`
   endings `'n`, `'ee`, `'er`, and `'ers`, including uppercase forms, their
   output markers, the context-class table condition for `'n`, and row flag
   bit 3. `NormalizePaul2013ContractionSuffixWithSupportedPrefix` composes
   the parsed-record, supported `FUN_100086c0`, compound, A140, generic, and
   context-encoder prefix branches with explicit model-row inputs.
   `NormalizePaul2013ContractionModelContextRow` applies a matched result to
   a caller-selected counted row. Outer row dispatch, unresolved selected
   branches, and runtime parity remain separate; see the Stage 81 report.
   `Engine.ApplyPaul2013B800JoinedDictionaryModelContextRow` ports the
   supported `FUN_1000b800` control flow: joined-surface embedded-dictionary
   lookup and group writes, followed on a miss by every recovered literal and
   table-backed suffix branch. An unmatched surface returns a native-style
   miss; malformed auxiliary row strings fail closed. The decompilation is
   static evidence only; runtime row comparison remains open. See the Stage
   111 report and its pseudocode attachment. The composed pass attempts this
   handler after contraction misses.
   `text.ScanPaul2013C3A0Component` ports the byte scanner in the adjacent
   `FUN_1000c860`: it tracks source offsets, splits on the two recovered
   character-class boundaries, and preserves both apostrophe-`s` cases.
   It is a scanner primitive, not the complete proper-name/TPP handler;
   non-ASCII behavior and runtime row parity remain open. See Stage 82.
   `text.Paul2013ContextWordPairSpecial`, `text.Paul2013IsNumberWord`, and
   `duration.EvaluatePaul2013FUN10008580` port the directly called
   `FUN_10008440`, `FUN_10008550`, and `FUN_10008580` predicates, including
   their sorted word-table lookups and model-relative string checks. The outer
   caller's row selection and native row parity remain open; see Stage 112.
   `duration.ApplyPaul2013FUN10007520UseMarker` also ports the `use` surface,
   phone-code, and `FUN_10008580` gate before its `D` rewrite and row-state
   bit-3 write. The known counted-row pass runs it after supported per-row
   handler results. Stages 114–119 and 122 port the other directly recovered
   `FUN_10007520` post-handler rules; runtime row parity and interactions
   across unresolved upstream handlers remain open. See Stage 113.
   `duration.ApplyPaul2013FUN10007520UpMarker` also ports the direct `UP`
   neighbor rule and its `08 35 00` phone-code write, and the counted-row pass
   runs it after the `use` rule. Runtime row parity and interactions across
   unresolved upstream handlers remain open. See Stage 114.
   `duration.ApplyPaul2013FUN10007520Minute` also ports the number-word-gated
   `minute` code write, and `duration.ApplyPaul2013FUN10007520CloseMarker`
   ports the `close` rewrite gated by `FUN_10008580`. The counted-row pass runs
   both for supported per-row results. Stages 116–119 port the context-window
   helper and the remaining direct table-backed rules. Runtime row parity and
   interactions across unresolved upstream handlers remain open. See Stage
   115.
   `duration.FindPaul2013ContextPhraseWindow` also ports `FUN_1000cc60`'s
   bounded L/R/B neighbor window, sorted lookup, and phrase continuation for
   an explicit table snapshot and comparator. All four calls from
   `FUN_10007520` are connected by Stages 117 through 119.
   See Stage 116.
   `duration.ApplyPaul2013FUN10007520MouthMarker` connects the existing
   17-entry class-12 table to the `mouth` rule's two-row `R` lookup and output
   write. Stages 118, 119, and 122 cover the remaining directly recovered
   table and literal rules in this post-handler pass. Runtime row parity and
   ordering with unresolved upstream handlers remain open. See Stage 117.
   `duration.ApplyPaul2013FUN10007520Bow` and
   `duration.ApplyPaul2013FUN10007520Lead` connect their recovered six- and
   ten-entry tables to the `B` lookup, output bytes, and row flag. The read
   temporal phrase table is covered in Stage 119. See Stage 118.
   `duration.ApplyPaul2013FUN10007520ReadMarker` also ports the `read`
   apostrophe gate, recovered class and phrase predicates, and code-byte
   rewrite. Static evidence does not prove runtime row parity. See Stage 119.
   Form-Y/S handlers now derive the `FUN_100091b0` suffix input from parser-row
   `+0x14` and `+0x18` fields when no override is supplied. The known
   `FUN_10044fd0`/`FUN_10045070` appenders and `FUN_100457b0` caller do not
   produce those fields; the remaining producer/later updater is still open.
   Do not infer the tail from the distinct `+0x28` flag. See Stage 120.
   The Y/S results expose that selected tail to the duration pass. On a miss
   from the direct contraction literals, `duration.Engine.ApplyPaul2013ModelContextTail`
   now runs the embedded `FUN_1000cb30`, direct `FUN_10002f10`, and
   `FUN_10009cd0` fallback chain in native order, appending the recovered `d`
   marker and row-state updates. The pass applies this before the post-handler
   rules and uses the duplicate-advanced context row for form-Y post-handlers.
   This composes existing ASCII and embedded-dictionary support; unresolved
   TPP/resource traversal, the ordinary parser-row offset producer, non-ASCII
   behavior, and runtime row parity remain open. See Stage 121.
   `NormalizePaul2013ApostropheSSuffixContext` also ports the direct
   case-sensitive `'s` stem lookup and its three output-tail cases from
   `FUN_1000c3a0`; the embedded-dictionary wrapper supplies the recovered
   `FUN_1000cb30` stem normalization. This suffix is reached only after earlier
   handlers miss and does not complete the proper-name/TPP path. See Stage 83.
   `NormalizePaul2013C3A0ComponentSequence` now composes the repeated scanner
   calls, first-single-component bypass, explicit component-handler boundary,
   `d` joining, 63-byte check, and final success/flag gate. Per-component
   dictionary, exception, TPP, and normalization choices remain callback
   inputs. See Stage 84.
   `text.Paul2013C3A0Gate` ports the source gate in `FUN_10008cc0`, including
   its mapped `Mac`/`Mc` prefix exclusions, observed character-class
   transition, and apostrophe-`s` checks. The source-driven sequence entry
   point now derives this gate directly. Stage 105 recovers the signed-index
   prefix: high-bit bytes produce the native ineligible result. Downstream
   scanner non-ASCII handling and native runtime comparison remain open; see
   Stages 85 and 105.
   `NormalizePaul2013C3A0EmbeddedRecordComponent` ports the direct
   two-byte-or-longer, positive-pronunciation-count branch using the existing
   `FUN_10003a70`/`FUN_10003c50` embedded lookup result. It copies the first
   internal-symbol string and sets row-flag bit 0. Stage 87 also ports the
   payload-bit-7-gated generic precheck: its successful return skips exception
   and `FUN_1000a140` handlers; it does not select the parsed string. The
   earlier handlers remain callback work, and row-for-row runtime parity is
   open; see Stages 86 and 87.
   `NormalizePaul2013C3A0ComponentSequenceWithSupportedHandlers` now composes
   the parsed-record branch, the bit-7-gated generic precheck, conditionally
   supplied earlier exception/`FUN_1000a140` handlers, the embedded-dictionary
   apostrophe-`s` rewrite, and the supported `FUN_10002f10`/`FUN_10009cd0`
   final output choice in recovered order. The earlier handlers keep their
   unresolved resource and row-state inputs explicit; see Stage 87.
   `NormalizePaul2013C3A0SupportedA140Component` adapts the bounded,
   directly-portable `FUN_1000a140` suffix cases for use after a caller's
   exception handler misses. Stage 89 adds the 9- to 31-byte mapped `liness`
   path, including its table-gated `i`-to-`y` retry, ahead of the existing
   long `ed` and `er` branches. Stage 90 adds the mapped long `ly` branch,
   its own table-gated retry, and the `&`/control/plus marker behavior. Stage
   106 captures the full signed short-table index window, removing the former
   lowercase-only table-index boundary. Stage 91 marks
   mapped 9- to 31-byte terminal `s` and `ing` inputs as A140 branches.
   Stage 96 ports the `ing` branch's ordered stem lookups, ASCII table
   predicates, and `#.` tail. Stage 97 ports the 5- to 31-byte terminal-`s`
   branch, including mapped fixed-width comparisons, ordered `s`/`es`/`ies`
   stem choices, the protected output gate, marker tail, and row-flag bit 3.
   Stage 92 ports the long mapped
   `ist` path, including stem lookup, `#79` marker tail, and branch ordering
   before `ly`. Stage 100 broadens `ly` to 5- to 31-byte inputs; Stage 101
   ports short `ing` with the pseudocode's `length < 7` guard around earlier
   table indexes; Stage 102 extends `ist` over the same input range; Stage 103
   ports mapped `ful` between `ist` and `ly`. Stage 107 records the directly
   compared suffix families represented by these handlers. Table-driven and
   state-dependent branches, exception resource/state production, and runtime
   row parity remain open. See Stages 88 through 92 and Stages 96 through 107.
   `NormalizePaul2013C3A0ComponentSequenceWithSupportedA140` now composes
   that subset after a required caller handler for the preceding
   `FUN_100086c0`/exception path, then runs the existing apostrophe-`s` and
   generic/context fallback stages. The required callback preserves the
   unported earlier stage instead of silently skipping it; row-level runtime
   parity and that callback's resource/state production remain open.
   `NormalizePaul2013C3A0ComponentSequenceWithFUN100086C0` now supplies that
   earlier handler from the full recovered `FUN_100086c0` path, including its
   caller gate, neighbor scan, static six-string table, and output context
   encoder. `PhonePayload.BuildPaul2013ParsedDictionaryRecord` now projects
   the loaded embedded payload into the bounded `FUN_10003c50` record,
   including metadata words, expanded phone slots, and path rows.
   `BuildPaul2013C3A0AlternateRecord` connects that projection to the engine,
   and `NormalizePaul2013C3A0ComponentSequenceWithEmbeddedFUN100086C0`
   removes the raw-record producer callback. The caller still supplies the
   model-row index; runtime row comparison remains open. The direct
   `FUN_10009cd0` fallback also propagates its `0x20` row-flag bit when output
   is nonempty. `Engine.NormalizePaul2013C3A0ModelContextRow` now reads the
   counted row surface/state, runs that composition, and copies matched output
   into the row phone field. `NormalizePaul2013C3A0ModelContextRows` carries
   writes across a caller-ordered row list. `RunPaul2013KnownModelContextPassesWithC3A0Rows`
   retains its row list for compatibility while the composed pass attempts
   C3A0 automatically after B800 misses. The shared-arena segment wrapper
   connects it to ordinary sentences. Unresolved neighbor cases fail closed.
   Row-for-row runtime parity remains open.
   recorded in `tools/revkit/work/reports/stage95-c3a0-local-record-projection.txt`.
   This is not a complete text normalizer. The short-E neighbor map is
   implemented only for the accepted ASCII letter/apostrophe path.
   Ghidra pseudocode is recorded in
   `tools/revkit/work/reports/stage27-generic-normalizer.c`,
   `tools/revkit/work/reports/stage27-generic-normalizer-character-map.c` and
   `tools/revkit/work/reports/stage27-generic-normalizer-dependencies.c`;
   the additional lookup-shaper pseudocode is in
   `tools/revkit/work/reports/stage27-generic-normalizer-lookup-followup.c` and
   `tools/revkit/work/reports/stage27-normalizer-key-shaping-dependencies.c`.
   `LookupPaul2013SurfaceSequence` composes caller-selected adjacent surfaces
   into exact hyphen-joined prefix lookups. The explicit-marker variant also
   ports the native `X`-row retry as an exact `h'`-prefixed lookup of the
   current normalized surface. The gate remains derived from supplied rows;
   row eligibility/order and the producer of the X marker remain explicit
   upstream inputs.
   `text.SplitPaul2013ExceptionPhoneCodes` ports the `d` delimiter retention
   and destination-row advance in `FUN_1000ca50` given explicit row delimiter
   counts. `text.Paul2013ExceptionSurfaceDelimiterCount` ports the
   `FUN_1000ca30` literal-hyphen count plus one, and exception matches can now
   derive those destination counts from the matched row surfaces.
   `text.DecodePaul2013ExceptionPhoneRows` composes the split with the
   runtime-backed internal phone-symbol labels. A local `San-Francisco`
   record decodes into two nonempty destination rows when supplied the
   observed one-component counts for its two source rows; selecting eligible
   surfaces and producing the retry markers remain caller work.
   `duration.Engine.ResolvePaul2013PronunciationExceptionSequence` connects
   exact-prefix lookup through destination-row decoding for explicit counts;
   `ResolvePaul2013PronunciationExceptionSurfaces` derives them from the row
   surfaces and composes the same decoding path. Eligibility and retry-marker
   production remain caller inputs.
   `text.LoadTPPDictionary` validates the sequential TPP record framing and
   all observed one/two-atom payload forms. `duration.Engine.LookupTypedText`
   applies the known key transform and returns opaque typed atoms.
   `TPPAtom.ComponentPattern` decodes the directly observed A-E component
   counts and raw 0/1 values, and `duration.Engine.LookupTypedComponentPattern`
   connects that decoder to exact shared-dictionary keys. The values remain
   unnamed; AX is retained as a one-component marker without a bit value.
   `Paul2013ProperNameTPPSelector` ports `FUN_10034180`'s one-through-four
   component A-D selector mapping; the captured x86 count gate skips selector
   construction for five or more components, despite E-family records being
   directly retrievable through the shared helper. The engine wrapper composes
   this selector with exact-key lookup. `AcceptPaul2013ProperNameTPPComponentPattern`
   ports the post-lookup gate: a one-component hit checks only its marker,
   while a multi-component hit rejects literal zero bits paired with marker
   `d` or `A`. `duration.Engine.MatchProperNameTPPComponentPattern` composes
   lookup and gate when callers supply the component markers. The caller
   pseudocode and instruction addresses are recorded in
   `tools/revkit/work/reports/stage6-tpp-callers.txt` and
   `docs/reverse-engineering/lead3-tpp-typed-code-findings-2026-09-24.md`.
   `text.BuildPaul2013ProperNameTPPSurface` now ports the preceding native
   component-row surface join, hyphen separator, stored-length accumulation,
   and 29-byte limit. `duration.Engine.MatchProperNameTPPComponents` composes
   those row fields with A-D selection, exact lookup, and the existing
   marker/value gate. `text.Paul2013ProperNameTPPComponentsFromArena` reads the
   counted `0x140` rows and projects their stored length, marker, C-string,
   and final candidate row's raw `+0x18` scanner-input dword;
   `duration.Engine.MatchProperNameTPPComponentArena` composes that projection
   with the supported lookup gate and reports when zero at `+0x18` bypasses
   all later context predicates. A nonzero or unavailable value leaves those
   predicates pending. `text.EvaluatePaul2013ProperNameTPPContextGate` ports
   the post-scanner Boolean decision tree for the nonzero branch when callers
   supply scanner counters/status and observed helper/table outcomes. It
   preserves the native low-word success/rejection paths without assigning
   meanings to pointer-backed inputs. Scanner production, candidate-row
   production, runtime context-table snapshots, and runtime parity remain
   open; the generic table search is ported in
   `tools/revkit/work/reports/stage136-fun100035d0-sorted-context-table-search.txt`.
   The one-component `FUN_10041df0` Roman-numeral and two-byte lookup is
   derived from scanner output by the table-backed gate; its bounded match
   cases are recorded in
   `tools/revkit/work/reports/stage138-proper-name-late-context-lookup.txt`.
   The direct `FUN_10010120` scanner-word predicate is also derived from
   scanner output; its source branch map is recorded in
   `tools/revkit/work/reports/stage139-proper-name-scanner-word-predicate.txt`.
   `text.ObservePaul2013ProperNameScannerInput` additionally ports the
   scanner's leading whitespace counters and status 8/9 early returns; its
   other ordinary token output and status remain explicit inputs. See
   `tools/revkit/work/reports/stage140-proper-name-scanner-prefix-and-terminal-paths.txt`.
   `text.ObservePaul2013ProperNameASCIIToken` now ports bounded mode-0x15
   paths for NUL-terminated ASCII digit runs, decimal values, comma groups
   with exactly three digits and a nonzero leading byte, mapped ASCII ordinal
   suffixes (`st`, `nd`, `rd`, and `th`) with the 11th–13th override, and ASCII-letter
   words with internal dot or hyphen separators. It also handles the two
   directly observed one-digit uppercase-letter patterns and the mapped
   two-component `United States` path. It derives `FUN_10062f50` weights from
   the shared letter-expansion records to enforce the native scan limit. Its
   duration wrapper composes known scanner results with the existing TPP
   context gate; unsupported scan shapes retain pending context. The bounded
   letter-token observer now stops before ASCII whitespace delimiters,
   preserves the native consumed-byte cursor, and drops an unjoined trailing
   dot or hyphen. ASCII punctuation also ends a letter token when it does not
   join an accepted internal dot/hyphen pair. Numeric tokens support one
   optional decimal after valid comma groups and one ASCII delimiter
   immediately before NUL; longer numeric lookahead, non-ASCII mapped ordinal
   comparisons, and other numeric forms remain open.
   Other digit combinations, other punctuation forms, multibyte input, and
   the rest of the scanner remain open; see
   `tools/revkit/work/reports/stage141-proper-name-ascii-letter-scanner.txt`.
   See also
   `tools/revkit/work/reports/stage129-proper-name-tpp-surface-assembly.txt`
   and `tools/revkit/work/reports/stage131-proper-name-tpp-component-arena-projection.txt`,
   with the bypass gate recorded in
   `tools/revkit/work/reports/stage132-proper-name-tpp-context-bypass.txt` and
   the evaluated decision tree in
   `tools/revkit/work/reports/stage134-proper-name-tpp-context-decision.txt`.
   `duration.Engine.MatchAndAppendProperNameTPPComponentArena` now connects
   the zero-input bypass and evaluated nonzero context gate to
   `FUN_10034110`'s source-row append. Pending or rejected context does not
   mutate the source arena. Candidate arena production, scanner production,
    runtime context-table snapshots, and runtime row parity remain open; see
    `tools/revkit/work/reports/stage133-proper-name-tpp-followup-composition.txt`.
   `MatchAndAppendProperNameTPPComponentArenaWithASCIIScannerInputAndContextTables`
   also composes the bounded mode-0x15 ASCII scanner with that append: an
   unsupported scan preserves the native zero-pointer bypass, while a
   nonzero pointer leaves context pending and does not append. This closes the
   static composition for those scan cases; scanner coverage and native row
   parity remain limited as described in Stage 141.
   `TPPDictionary.LookupNumericText` and `duration.Engine.LookupTypedTextCode`
   select an exact F or G suffix after lookup, preserving its decimal text.
   `TPPDictionary.LookupSelectedText` and
   `TPPDictionary.LookupSelectedAtom` expose the selected atom for every A-G
   family; `duration.Engine.LookupSelectedTypedAtom` and
   `LookupSelectedTypedText` expose its tag and suffix. This includes secondary
   atoms such as G95 in `A0 G95`; AX returns its literal X suffix without
   assigning it semantics.
   `MatchPaul2013TPPWindowWithDictionary`
   connects direct F-atom matches to the loaded shared dictionary at native
   index zero and rejects nonzero indexes whose alternate lookup resources
   are not loaded. Marker production and selector construction in other
   callers remain unresolved.
   `text.ParsePaul2013NativeInteger` ports `FUN_100645ba`'s attribute-driven
   leading skip, optional sign, digit accumulation, and 32-bit wraparound.
   `TPPAtom.NumericByte` applies the F/G dictionary's validated decimal suffix
   grammar through that helper and narrows to one byte.
   `text.ApplyPaul2013TPPTokenRange` ports `FUN_1000e0c0`'s generic writes for
   a supplied byte value and inclusive row range; F/G numeric atoms compose
   that writer after parsing. It also exposes the known writer for A–E when
   their caller provides a narrowed value and range. Those selector-specific
   conversions and ranges remain unresolved.
   `text.ApplyPaul2013TPPNumericTextRows` and
   `duration.Engine.ApplyPaul2013TPPNumericTextRows` compose the observed `FUN_1000dfc0`
   route for direct F compound codes, direct G row codes, and the no-code WAB
   class-byte fallback. It derives the affected row range and processing
   flags. `text.WritePaul2013TPPNumericTextUpdates` applies the exact
   `+0x20`, `+0x25`, and `+0x26` writes to a copied 0x94-byte caller arena;
   `Paul2013TPPWindowRowsFromParserArena` projects the native `+0x00`, `+0x04`,
   `+0x08`, `+0x23`, and `+0x34` fields directly from those rows, and
   `duration.Engine.ApplyPaul2013TPPNumericParserArena` composes that
   projection, routing, and copied-arena writes. The lower-level
   `ApplyPaul2013TPPNumericTextRowsToParserArena` remains available for
   caller-supplied typed rows. The
   `BuildAndRunPaul2013KnownModelContextPassesWithTPPFromParserRows` entry
   point runs model-context projection/normalization before the TPP row pass,
   and applies the recovered `FUN_1000e990` context-table flag scan between
   them, matching that order in `FUN_1000e2f0`. Selector extraction across the
   observed A-G payload families is now available; remaining caller-specific
   component row effects and alternate dictionary indexes are unresolved.
   The ordinary-source and shared-arena sentence-segment variants compose
   the same ordered stages
   from their bounded source-row producers while retaining explicit row
   controls.
   `Table.LookupPaul2013WABClass` ports the
   sorted-table binary search through the native mapped-character comparator.
   `LookupSelectedAtom` reproduces the selected discriminator result over the
   validated one/two-atom payload grammar, including secondary selectors;
   alternate dictionary indexes and caller-specific non-F/G row effects remain
   unresolved.
   `text.MatchPaul2013TPPWindow` ports `FUN_1000e160`'s bounded following-row
   scan, case-mapped `dash`-row fast-skip, `A`-row eligibility, cumulative
   cutoff, five-row limit, hyphenated compound assembly, and
   last-successful-match behavior. The caller passes the observed character
   map. Ghidra pseudocode and call-site disassembly are preserved in
   `tools/revkit/work/reports/stage25-string-helpers.c`,
   `tools/revkit/work/reports/stage25-tpp-lookup-chain.c`,
   `tools/revkit/work/reports/stage25-tpp-index-lookup.c`, and
   `tools/revkit/work/reports/vt_pau-objdump-disassembly.txt`.
   `text.EncodeEmbeddedKey` ports the observed character-map and greedy
   pair-compression algorithm. The documented key tables and 256-slot Paul
   2013 phone-ID codebook are now provisioned in the engine, so a canonical
   token can resolve through the local shared dictionary into internal phone
   symbols. `CMUPhone.TreeFeatures` maps each labeled phone to its original
   internal symbol byte and the observed one-based identity ordinal and stress
   value. The internal byte and tree ordinal are distinct representations.
   `BuildPaul2013PhoneNeighborhoods`
   arranges those values for a phone span and marks missing edge neighbors.
   `BuildPaul2013TokenPhoneNeighborhoods` applies it to the complete flattened
   utterance and groups the result by token, deriving cross-word phone
   neighbors automatically. Neither helper infers tree-specific vector
   positions. `BuildPaul2013PhoneGroups` ports the vowel-anchor and allowed
   onset-cluster boundary scan from `FUN_10013f30`. Token results now include
   the `FUN_10013c00` row bytes and onset/nucleus/coda labels, including its
   no-vowel fallback. `BuildPaul2013TokenDurationInputsFromText` assembles all
   nine inputs for the ordinary single-utterance marker path.
   `text.SplitPaul2013TokenPhoneBlocks` ports `FUN_10012c70`'s delimiter rule
   for supplied per-phone marker bytes: non-`'0'` markers and the final phone
   end a block. `text.ApplyPaul2013PhoneBlockDurationLimit` ports the following
   cumulative doubled-value limit and inserts `[` after the preceding phone
   when the running total exceeds 1,000. `text.BuildPaul2013MarkerTreeGroups`
   ports the separate `FUN_10012df0` rule that keeps `\` and `]` boundaries
   within a group and sums the explicit per-phone `+0x94` values. These byte
   value producers remain caller work. `text.ApplyPaul2013MarkerTreeDecision`
   ports the `FUN_10012f00` result branch: a tree result above 500 writes `\\`
   at the preceding record's `+0x3bd` and resets the two scan accumulators. The 15-short
   input producer is ported by `text.BuildPaul2013MarkerTreeInput` from
   adjacent 0x3c0-byte records and explicit token/split counters. It derives
   the record fields and applies `DAT_1007b9e0`; callers still resolve its two
   pointer-backed string bytes and select the native marker tree.
   `text.EvaluatePaul2013MarkerTreeInput` evaluates that scalar tree and applies
   the marker decision. `text.RunPaul2013MarkerTreeToken` composes those steps
   over a caller-provided native-order token scan and carries the split
   accumulators forward. The two string bytes are restricted to ASCII because
   high-bit values index before `DAT_1007b9e0` under native signed-char rules.
   Full token-group arena iteration and string-pointer resolution remain
   outside this helper.
   `BuildPaul2013TokenPhoneNeighborhoodsWithPhoneMarkers` rebuilds
   `FUN_10013f30`/`FUN_10013c00` groups inside those blocks, and the explicit
   duration and pitch evaluators consume the resulting rows. The ordinary
   evaluator supplies one `'0'` marker per phone, preserving one block per
   token. `BuildPaul2013TokenBoundaryStateResult` ports `FUN_100130e0`'s
   per-token value normalization and state-2 propagation from supplied arrays;
   `BuildPaul2013TokenBoundaryMarkers` ports its marker mapping and final-token
   override. `BuildPaul2013InitialTokenMarkers` ports the 13-entry marker table
   at `0x10079a30` for explicit signed state codes; the code producer and
   state arrays remain unresolved.
   `duration.OpenPaul2013` loads the shared
   embedded dictionary, its 531-node `engbi.tree3` ambiguity classifier, and
   Paul duration trees from
   local roots. The engine generates complete 15-short rows for each path
   code, evaluates that classifier, and selects the winning path group and
   matching pronunciation alternative. The
   `text.ParsePaul2013PronunciationPath` now ports the group delimiters and
   duplicate control-marker handling from `FUN_100105d0`/`FUN_10010550`, and
   `text.RankPaul2013PronunciationPathGroups` ports classifier-output
   membership scoring plus the first-wins maximum using the recovered
   46-entry class map. `text.SelectPaul2013PronunciationByPathMarker` ports
   `FUN_10003f10`'s ordered selector-byte search, `)`/`%` match, `0x17`
   fallback, and row-zero default, returning the aligned phone string. The
   selector producer remains unresolved. `PhonePayload.BuildPaul2013DictionaryPhoneRowsForState`
   ports the ordinary/final and conditional marker branches of
   `FUN_1000d450`; `LexiconFrontend.ResolveText` retains the ordinary/final
   row projection on each token, and pronunciation flattening carries an
   independent copy into each selected token span. `text.WritePaul2013DictionaryPhoneRow`
   writes the established row fields from explicit token-index and surface
   inputs; source marker and selector-gate production remain explicit.
   `text.BuildPaul2013PhoneContextRows` projects selected spans into the
   counted `0x70` rows from `FUN_1000ea20`, carrying the selected phone string,
   surface, result type, marker class, and metadata.
   `BuildPaul2013PhoneContextRowsFromParserRows` derives the native `X` source
   marker from parser-row byte `+0x24`; mapping selected lexical spans to
   parser-row indexes remains explicit. The initial processing bit is ported.
   Selected `FUN_100091b0` branches are implemented, including its final
   counted-row pass for the literal `the`/`de` cases. Its generic
   dictionary/TPP normalization, remaining first-pass row updates, and
   `FUN_10007520` outer row mutations remain unimplemented. The complete
   low-short return is derived from source kind, form, and adjacent source-row
   indexes.
   `LexiconFrontend.AnalyzeText`
   returns resolved tokens alongside the partial rows from
   `BuildPaul2013KnownPronunciationFeatures`
   and candidate rows from `BuildPaul2013PronunciationPathCodeFeatures`. These
   mark positions 0–3 and 14 from token surfaces, position 13 for each path
   code using the recovered `DAT_100783ec` map, and values for triplet
   positions 4–7 using empty-input handling, literal neighbor rules, ordered
   sorted-word tables, then the character-attribute digit scan, punctuation
   scan, and default class. The DLL's `DAT_1007e188` bit `0x10` is set only for
   ASCII digits at the signed-byte values scanned here; the final set at
   `0x10077724` contains the punctuation bytes. Positions 8–12 are generated
   by the recovered ordered word-shape checks in `FUN_10007160`, applied to
   the two preceding words, center word, and two following words; missing
   words return zero. The classifier preserves the native ordered `ing`,
   `ist`/`ists`, `ed`, and final-`s` return branches and uses the DLL byte map.
   The static feature builder now produces all 15 values
   for each path-code candidate. Classifier outputs are accumulated by path
   group; the winning group selects the corresponding lexical alternative.
   The resolver enforces the observed one-group-per-alternative layout. The
   full local dictionary audit confirms that layout for 829 alternatives
   across 414 ambiguous records. The
   static feature values have not yet been compared row-for-row against
   runtime captures. `Engine.Evaluate` selects ambiguous pronunciations with
   the loaded classifier and ranker, then derives duration inputs without
   caller-supplied per-phone metadata.
   `duration.EvaluatePaul2013Text` ports the observed phone-class selector
   tables and evaluates each assembled vector against its duration tree. The
   lexical convenience path derives 11 pitch inputs from phone/group
   projections, evaluates the scalar tree, forwards its stored signed byte as
   vector input 11, and evaluates the paired 12-value tree. The native-record
   path follows `FUN_100138c0`'s counted-group loop:
   `text.BuildPaul2013ModelRecordPitchInputs` emits one input per counted
   group, and `duration.EvaluatePaul2013ModelStateArenaPitch` evaluates the
   selected scalar/vector pair from the model-state arena. These are raw group
   vectors before the separate adjacent-group smoothing pass.
   `EvaluatePaul2013PitchTextWithMarkers`
   also ports the `FUN_100138c0` terminal dispatch: `^` selects `bt/bf`, `Z`
   selects `sbt/sbf`, `[` selects `qbt/qbf`, and other marker bytes use
   `nbt/nbf`. The ordinary text entry point still supplies `Z` only at the end
   of the utterance. `FUN_100135d0`'s recognized terminal markers now select
   duration boundary identity 40 and position state 3; other marker bytes use
   identity 42 and state 2. `duration.Engine.EvaluateWithMarkers` passes the
   same supplied markers to duration and pitch. Source-text marker production
   remains open. Explicit per-phone block subdivision is available in `text`
   and the marker-aware evaluator; the ordinary evaluator treats each lexical
   token as one block. `text.ParsePaul2013ContextCodes` ports the bounded
   `FUN_10016c90` code-string scan, including its preceding-phone `d`/`c`
   markers and separate `M` flag. `text.SummarizePaul2013ContextCodeState`
   ports its explicit per-token -1/12 state mapping and final 5/6/7 dispatch;
   the opaque code meanings and caller state fields remain unresolved. The
   `text.PopulatePaul2013ModelStateRecordsFromParserRows` projection now writes
   the observed phone-code, marker, M-flag, row-control, state-flag, and pointer
   fields into counted 0x3c0-byte records from already-produced 0x94-byte
   parser rows. Pointer values use a caller resolver, and unspecified record
   bytes are preserved. Per-token state values and terminal pitch-value count
   remain caller inputs; the projection derives the final mode byte at arena
   +0x4770a from the last record. This static-derived projection has no
   captured native-record parity claim. The
   `FUN_100130e0` state-array and marker decompiles are
   preserved in `tools/revkit/work/reports/stage22-context-marker-producers.c`.
   `BuildPaul2013TokenBoundaryMarkerTransitions` also ports its adjacent-record
   slash conversion and previous-state flag gate. `ReadPaul2013TokenBoundaryArenaInputs`
   reads the signed marker code at record `+0x3ac`, transition state byte at
   `+0x3b0`, and duration-limit byte at `+0x95` from the counted model-state
   arena. `BuildPaul2013TokenBoundaryMarkersFromModelStateArena` composes the
   initial marker lookup and transition in the observed `FUN_100130e0` order.
   `ApplyPaul2013TokenBoundaryDurationLimitFromModelStateArena` applies the
   recovered cumulative limit using the arena's `+0x95` bytes and explicit
   caller markers.
   The later state/value arrays and their producer remain unresolved;
   `BuildPaul2013TokenBoundaryStateResultWithTransitions` remains a staged
   post-dispatch composition for explicit callers.
   `InitializePaul2013TokenBoundaryArena` and
   `FinishPaul2013TokenBoundaryArena` now expose native arena writeback in
   order around the explicitly supplied marker-tree stage. The latter resets
   flags, preserves transition flags through dispatch, inserts boundaries
   with native overflow rescanning for up to 100 records, and rebuilds group
   descriptors. The x86-confirmed propagation sets state[i-1] when value[i]
   is nonnegative, correcting the previous same-index implementation.
   `Paul2013FinalizedModelState.FinishTokenBoundaries` derives the nonfinal
   outputs from finalized parser state +0x20 at 0x94 strides. Caller value/
   state initialization and marker-tree orchestration remain open; Stage 158
   records evidence and limits.
   `PopulatePaul2013ModelPhoneGroupArena` now writes FUN_10012c70's group
   pointers, group counts, per-phone labels, and the five recovered group-row
   fields into the caller arena, preserving other bytes. It retains the
   uncounted fallback tail and its possible overwrite by the next record.
   `PreparePaul2013TokenBoundaryArena` composes that writer with initial
   marker writes and pre-tree record descriptors. Its native address base
   is explicit. Stage 159 records the original implementation boundary;
   Stage 162 extends model-symbol decoding through native identity/stress
   tables, including structural M, and corrects the swapped B9E0/BAA8 windows.
   `RunPaul2013MarkerTreeArena` now derives those descriptors and scans each
   group in reverse, resolving eligible current/previous +0x2e4 string reads
   and applying decisions to the preceding record's terminal byte. The
   scalar tree remains supplied; the native pointer is model +0x200. The
   feature writer's historical NextTokenCumulative field is supplied with
   this group's +0x94 total, not the following group's span.
   `RunPaul2013TokenBoundaryPipeline` composes preparation, tree scan, and
   dispatch; `Paul2013FinalizedModelState.RunTokenBoundaries` also extracts
   parser outputs from retained finalized state. Caller array initialization
   remains explicit. Stage 162 supports structural M through the native
   identity/stress tables. See Stage 160 for scan evidence.
   `PopulatePaul2013RecordGroupDescriptorArena` now writes FUN_10012df0's
   descriptor counts, starts, spans, and native pointers, preserving padding.
   Short-offset value views are converted to byte offsets for native pointers.
   The composed pipeline applies this writer before marker-tree scanning and
   after final dispatch, so changed boundaries rebuild the native header and
   descriptor stream as well as Go value views. See Stage 161.
   Stage 163 establishes the marker-tree resource: FUN_10012280 formats
   `%sdict-eng/engbi.tree3` from DLL VA 0x100793ec and FUN_10001000 loads it
   into the same +0x200 object later read by FUN_10012f00. The workspace
   pipeline uses the shared catalog's existing Pronunciation tree, and the
   duration engine supplies its loaded resource through
   `RunPaul2013ModelTokenBoundaries`.
   `ApplyPaul2013PositionStateProgramToWorkspace` reads the native scalar
   defaults and writes the six recovered state arrays. Boundary values begin
   at +0x121a64, one cell into slot 3, while boundary states use slot 4 at
   +0x121d80. The final value's retained tail cell is preserved rather than
   initialized. `RunTokenBoundariesWithWorkspace` derives these arrays and
   writes the final updates back. Range/event table producers, pointer
   address mapping, and full native arena comparison remain open.
   Stage 164 corrects event cursor propagation: upper-clamp and ordinary
   matches write the same selector cursor, and the terminal pass starts from
   the final consumed index. OverflowBoundary remains a compatibility
   diagnostic, not a second native cursor. Program results retain event
   results and workspace adapters read/write +0x1223f4/+0x122404, including
   terminal continuation and the +0x122440 sentinel reset.
   `ApplyPaul2013PositionStateProgramFromModelArena` derives range row keys
   from record starts, selector-3 gaps [previous end, current start] with
   first minimum zero, selector-4 gaps [current end, next start] with final
   maximum at segment end, and the terminal source-end gate. Reversed gaps
   from overlapping records are empty under native comparisons. Existing
   parser-span convenience adapters remain available but do not substitute
   for these native inter-record ranges. Event values and clamps remain inputs.
   Stage 165 connects that position program to the finalized model-state
   boundary pipeline. Optional mapped intervals are written into records
   after event processing, slot-7 results replace record +0x2df flags, and
   the final mode is recomputed before initial marker selection. The duration
   engine supplies its loaded shared tree through
   `RunPaul2013PositionStateBoundaries`. The full arena is now supplied to
   `MapPaul2013PositionStateRecordIndexes`; its bounds check includes the
   +0x64c record base, preventing short-arena slices and out-of-bounds reads.
   Stage 166 replaces explicit range/event slices and clamp limits with the
   recovered workspace descriptors at +0x1223c0 + selector*0x10. It follows
   native count/pointer gates, resolves and copies needed int32 cells, retains
   event cursors, and uses the DLL tables at VA 0x1007d66c/0x1007d690.
   `RunPaul2013PositionDescriptorBoundaries` exposes the loaded composed path.
   Stage 167 derives post-state mapping pointers from +0x47790/+0x47794/
   +0x47798, the ordinary clamp bound from workspace +0, and the exact remap
   gate (shared object +0x20424 == 1). It resolves only the table prefixes
   needed by record indexes and subsequently used final-table indexes.
   `RunPaul2013NativePositionWorkspaceBoundaries` combines those mapping
   inputs with the descriptors and boundary processing. Pointer ownership,
   source-table/descriptor producers, native snapshot production, segment
   advancement and full original-engine arena parity remain open.
   Stage 168 closes segment-default preparation and successful segment
   advancement: FUN_10022850 selects exact-1 overrides or engine-context
   defaults, clamps the three signed scalars, propagates active/default fields
   and resets +0x1312c4. The native workspace pipeline advances +4 using the
   parser consumed-byte count with x86 32-bit wrapping. Its prepared wrapper
   composes both operations. Zero-byte completion, empty-record retries and
   the earlier source-parser segment loop remain unported here.
   Stage 169 ports FUN_10022dc0's driver around an explicit FUN_1003d350
   parser callback: null source returns zero without preparation, empty or
   negative row counts with nonzero consumed bytes advance and retry, and
   zero consumed bytes set workspace +0x44 to 1 after row/default writes.
   Successful counted segments run native state/mapping and advance +4.
   Position/mapping readers now validate count/storage independently of
   phone/control fields populated only later by FUN_10016c90.
   Stage 170 connects this driver to phone/control record population using
   its produced slot-7 values and retained parser state, then to the loaded
   boundary tree and workspace writeback. Completion branches skip these
   dependencies. The duration engine exposes the composed driver. Parser
   implementation, table producers, pointer ownership and original-engine
   arena/whole-synthesis comparison remain open.
   Stage 171 ports the alternate inner parser FUN_1003e070 around an explicit
   FUN_1005a350 scanner. Native status-1 tokens append A/D rows, angle contents
   append A/S rows, and bracket contents are skipped; scanner coordinates and
   flag 0x15 flow through the existing native row writer. Whitespace, terminal
   row-count gates and failure/full-row low-short returns are preserved.
   Its dispatcher adapter fills the alternate callback only when absent.
   General scanner and primary-parser reconstruction remain incomplete.
   Stage 172 ports FUN_1003d3d0's primary caller control flow: mode handler,
   scanner-prefix row type, status-gated handler, terminal/fallback handlers,
   native cascade order, repeated candidate rollback, marked-row boundaries
   and FUN_10033760's special two-row merge. Its handler map retains native
   addresses and fails on a missing reached handler. Scanner fields +4 and
   +0x14 remain explicit raw values. The dispatcher now supplies both inner
   loops unless overridden; scanner and individual handler bodies remain
   separate implementation work. No source-text parity is established.
   Stage 173 ports the complete FUN_10051a00 terminal-handler body around
   the scanner and FUN_10051cc0 context recognizer. It rescans in mode 1,
   writes qualifying row types, handles native status 8/9 and structural
   A2 FD tokens, and returns signed advances from scanner +0x14 rather than
   the scanner return. Prior-row fields and punctuation eligibility are derived.
   Stage 174 supplies FUN_1005a350's full whitespace-only/multiline early
   result fields, including coordinates, status, prefix flag and advance.
   These paths compose with the recovered terminal handler and primary loop
   without fixture scanning. Ordinary scanning and the punctuation recognizer
   remain explicit dependencies; no whole-text or synthesis parity is claimed.
   Stage 175 ports FUN_10052000's full ordered six-record N/P/A classifier,
   including signed-character attributes, case/spacing predicates, digit and
   letter strings, quoted following tokens and mapped literal comparisons.
   Stage 176 ports FUN_10063300's backward context copy, FUN_10063660's
   three-record rolling scanner and FUN_10051cc0's window assembly/dispatch.
   The recognizer binds the complete NUL-terminated source, projects exactly
   the copied scanner fields and tests the fallback result's low byte.
   Ordinary scanning and FUN_10052520/FUN_10001670 key/table binding remain
   explicit dependencies. Tests cover recovered caller behavior using labeled
   scanner/lookup fixtures; no new original-engine or whole-text parity claim.
   Stage 177 ports FUN_10052d00's abbreviation lookup using the existing
   native mapped binary search, exact/mapped case flags and adjacent variant
   scans. Qualifying variants preserve the original search index. It also
   ports FUN_100526b0's three-table membership mask, FUN_10052710's counted
   case/status feature, FUN_100527c0's punctuation literal/jump-table feature
   and FUN_10052c90's first-dot predicate. All local abbreviation resources
   are covered by row/case-variant tests. Full key production, following-word
   classification and the sbd.tree3 evaluation adapter remain implementation work.
   Stage 178 ports the complete FUN_100528b0 following-word classifier:
   sbdw_sort.txt2 membership, the DLL's 48-short class translation, counted
   digits/uppercase, punctuation groups, mapped suffixes and recursive
   contraction prefixes. FUN_10052520 now produces all 16 key shorts from
   the recovered scanner window, including native -1 empty-neighbor fields.
   `LoadPaul2013TerminalModel` loads only the four relevant tables and the
   shared scalar sbd.tree3; its Lookup method supplies the recognizer's model
   fallback. The local tree has 120 nodes and 121 scalar leaves. Tests compare
   the class translation against the read-only DLL and exercise every local
   sbdw row, full key layouts, recursive branches and recognizer composition.
   General scanning and original-engine differential coverage remain open.
   Stage 179 supplies standalone ASCII letter and punctuation scanner paths
   in modes 0/1, mode-0x12 punctuation, ignored-byte skipping and complete raw
   scanner fields. Letter limits use FUN_10062f50's recovered 26-entry cost
   table and the native 29-byte maximum. Punctuation run statuses/limits and
   single-byte backslash behavior are preserved. Raw field +0x18 now represents
   the post-token lookup short at +0x1a; nonzero pointers require explicit
   FUN_1005e010/FUN_1005e1a0 callbacks. Supported text runs through the real
   scanner, loaded punctuation model and terminal handler without fixture
   scanners. Mixed joins, numeric and non-ASCII paths remain delegated work.
   Stage 180 ports FUN_1005e010/FUN_1005e1a0 over explicit native index
   projections: longest-first length buckets, inclusive endpoint/midpoint
   searches, exact/mapped prefix comparison and repeated boundary-gate calls.
   It also ports FUN_1005f840's full boundary return decisions and
   FUN_1005f9b0's apostrophe suffix predicate. The scanner adapter resolves
   the model address to both indexes and writes the recovered +0x1a flag.
   Native pointer projection and mode-0x17 scanning remain explicit work.
   Stage 181 supplies that native pointer projection via bounded memory and
   C-string readers. Model +0xc/+0x10 select the two index objects; their
   signed count, bound arrays and 20-byte row key pointers are decoded and
   copied. Negative/inactive ranges are skipped and active keys are stored
   sparsely, preserving row indexes without reading gaps. The arena scanner
   lookup adapter reads fresh snapshots for each call. Mode-0x17 scanning
   and the producers that build these native indexes remain open work.
   Stage 182 ports FUN_100621f0's apostrophe suffix scanner and mode-0x17
   ASCII word/hyphen/contraction paths. Suffix bytes retain input case; the
   plural-s apostrophe requires following whitespace, and n't recognition
   requires preceding n. Dot-apostrophe forms retain the dot. Cost-limited
   hyphenated words backtrack to the last delimiter as the native scanner does.
   The index-boundary integration now uses this scanner without fixtures.
   Numeric/non-ASCII paths and native index construction remain open.
   Stage 183 ports the common ASCII numeric loop for modes 0/1/0x17,
   including 29-byte truncation, decimals, comma groups with exactly three
   following digits, native ordinal/teen suffix rules and the exact mixed-case
   Rd exclusion. Mode 0 401K/401(K) returns word status rather than numeric
   status. Leading decimal scanning is supplied when its preceding byte is
   known to permit it; unknown predecessor state delegates explicitly.
   Numeric boundary lookups now use the recovered scanner without fixtures.
   Other modes, mixed words, non-ASCII and native index construction remain.
   Stage 184 ports mode-0/1 dot and mapped `int'` joins, plus mode-8
   ASCII word joins. Mode 8 gates hyphens and eligible dots with exact then
   mapped model searches over the following source. A match stops before
   the delimiter and bypasses post-token lookup. Without a prior dot,
   lowercase-first words stop before a dot followed by an uppercase letter.
   Recognized apostrophe suffixes can continue into a hyphenated component
   without an internal lookup; `wi'` is the separate mapped special case.
   Cost/length overflow rolls back to the last dot/hyphen for all supported
   word modes. Mode 0/1 does not consume ordinary hyphen/contraction joins;
   x86 branches at 1005acb6/1005af62 establish that mode distinction.
   Mode-8 nonletter starts, other modes, non-ASCII and native index builders
   remain open, as does original-engine differential scanner coverage.
   Stage 185 supplies ASCII word branches for modes 9/0xc/0x11/0x12/0x15/
   0x16/0x18/0x19/0x1a/0x1c. Modes 0xc/0x1a join unrestricted apostrophe
   components; mode 0x16 also joins hyphens and common dots. Modes 0x11/0x15
   add uppercase-letter/digit/uppercase-letter and mapped `United States`
   joins. Mode 0x19 consumes a digit after one uppercase letter without
   requiring another uppercase letter. Other listed modes stop at delimiters.
   The common dot/`int'` dispatch also precedes mode-0x17 suffix handling;
   the prior Stage 182 test expectation that `word.next` stopped at `word`
   was corrected from x86 at 1005ac34-1005acb9. All supported word branches
   share cost accumulation, last-delimiter rollback and native result fields.
   Numeric/nonletter starts in these additional modes still delegate.
   Stage 186 ports the remaining special ASCII word branches: mode 5 joins
   mapped `pa` followed by apostrophe and an `anga` prefix; mode 0xf consumes
   the trailing dot of mapped `A.M`/`P.M`; mode 0x1b applies case-sensitive
   abbreviation endings, bounded slash endings and one-letter slash/uppercase
   ampersand joins. Mode 0x1d joins a digit after a letter run; consecutive
   digits stop after the first joined digit. Other mode values now take the
   native default dot/`int'` branch, including signed mode values. This closes
   the static ASCII-letter dispatch inventory at 1005ac34-1005b503 while
   retaining explicit fallback for non-ASCII and unrecovered nonletter starts.
   Original-engine differential coverage and scanner index construction remain.
   Both paths return raw tree values; their physical units and timeline
   conversion remain unknown.
   `ApplyObservedPaul2013OnsetIdentityFeature`
   writes the controlled consonant identity ordinal to position 0 of an
   explicitly selected onset-scalar input; it does not select the tree family.
   `ApplyObservedPaul2013InteriorDurationPhoneFeatures` fills the current,
   previous, and next phone identities at duration-input positions 0–2 and
   current stress at position 3 for an interior phone.
   `BuildObservedPaul2013DurationTreeInput` remains a low-level helper that
   accepts explicit neighbors and metadata. Its concrete-phone
   `BuildObservedPaul2013InteriorDurationTreeInput` remains available for
   controlled tree probes.
   `ApplyObservedPaul2013PairedScalarFeatures`
   writes the observed second/first phone identity ordinals into positions 0
   and 1, and fills position 3 only for a second phone labeled AH. The
   AH-specific helper remains as a stricter compatibility entry point. Other
   slots stay caller-supplied because their producers are unresolved.
   `ApplyObservedPaul2013AHStressFields` fills position 2 in the traced third
   scalar and 12-value vector inputs, and copies the supplied preceding scalar
   result into vector position 11. It does not evaluate either tree.
   `BuildObservedPaul2013DurationTreeInput` now has a row-order regression
   check against the first two nine-short inputs in the captured `P AH0` VTML
   phoneme probe (`tools/revkit/work/stage10/feature-p-ah0-trees.log`).
   `BuildPaul2013TokenDurationInputsWithPositionStates` composes phone-group
   row production and vector assembly when the caller supplies the observed
   per-phone states. `BuildPaul2013TokenPitchInputsWithPositionStates` applies
   the same states to the pitch inputs, and
   `Engine.EvaluateSequenceWithPositionStates` evaluates both tree families
   with that shared state vector. Tests supply state 1 for both captured
   phones and evaluate the loaded duration/pitch trees; they validate the
   captured inputs and downstream tree calls, not the source-to-state producer
   for this VTML path. The convenience text builder produces terminal state 3
   for its ordinary `Z` path, while this VTML capture uses state 1; it does
   not claim parity for the forced-phoneme path.
   Call-specific tree-vector construction and general text normalization
   remain unimplemented.
   Compare the resource parsers and transforms with runtime traces.
3. **Text to selected units — lexical lookup and class catalog started;
   context generation and unit selection pending.**
   `engine/text.LexiconFrontend` tokenizes ASCII surface runs, expands the
   observed plain/signed integer and decimal forms plus inferred numeric
   ordinal spellings from 0 through 31, resolves the embedded
   dictionary, and returns CMU-labeled alternatives with their original
   internal symbol bytes. It retains the four embedded-payload metadata flags
   and the recovered dictionary phone-row projection on each token and through
   phone-sequence flattening; their semantics remain unlabeled, and the flags
   are not the five duration-tree row values.
   `ParsePaul2013CMUPronunciation` and
   `BuildPaul2013CMUPhoneSequence` also parse the captured `x-cmu` label form
   into a caller-selected phone sequence and recovered internal symbol bytes.
   `ParsePaul2013VTMLCMUPhoneme` now extracts that pronunciation and source
   span from the directly captured single
   `<vtml_phoneme alphabet="x-cmu" ph="...">surface</vtml_phoneme>.` form.
   `duration.Engine.EvaluateVTMLCMUPhoneme` connects that source form to both
   loaded tree families when the caller supplies the capture's marker and
   position-state arrays. Other VTML forms and source position-state
   production remain unsupported.
   The ordinary lexical path preserves separators and
   spells the captured `555-1234` telephone form as cardinal groups joined by
   “to”, matching the Stage 20 parser text, and rejects unsupported numeric
   separators. Other digit groups within this layout follow cardinal
   normalization as inference; other telephone layouts are
   not generalized from that single capture. Ordinal
   normalization uses inferred English spelling
   rules; only `second` and `fifth` have direct runtime output evidence.
   The captured `3:45 PM` form is normalized to `three forty five PM`;
   additional 12-hour `H:MM` values use explicit formatting bounds and are
   inference. The captured `01/02/2024` slash date expands to
   `January second twenty twenty four`; MM/DD/YYYY widths and month/day bounds
   beyond this case are inferred, and month-specific calendar validity is not
   checked. Currency amounts with optional two-digit cents and integer/decimal
   percentages are now normalized; only `$5.00` and `25%` are directly
   captured, so broader wording is an implementation inference. Other date
   and time forms and number grammars remain unimplemented. Its output is not
   the seven-byte model context consumed by selection.
   `LexiconFrontend.ResolveTextWithPaul2013InlinePauses` also recognizes the
   captured self-closing `<vtml_pause time="N"/>` form, retains its position
   between resolved tokens, and converts milliseconds to 16 kHz frames. The
   Stage 16 200 ms and 1,000 ms cases map to 3,200 and 16,000 frames. This
   event metadata is carried through explicit pronunciation flattening.
   `duration.Engine.EvaluateWithInlinePauses` now resolves those tokens,
   chooses pronunciations through the loaded classifier, evaluates the
   duration and pitch trees, and returns both tree results and pause events.
   `synthesis.InsertPaul2013InlinePauses` inserts the corresponding zero PCM
   frames when callers provide each token boundary's output-frame offset. The
   text-to-timeline mapping does not yet produce those offsets, so this is not
   wired into ordinary synthesis; other markup forms remain unsupported.
   `voice.BuildPaul2013ClassCatalog` reconstructs sorted exact-key classes and
   unit membership from the local four-bank indexes, and its exact lookup
   consumes an already-built context; `text` and `selection`
   implement the observed key transform, feature views, weighted mismatch
   calculation, and class-count/population limit behavior. Catalog construction
   matches the native 61,566-key table, three captured key-to-ID lookups, and
   their captured unit expansions; broader candidate-list parity remains open.
   `selection.UniqueExactSelector` connects the exact lookup to the pipeline
   when every context has one indexed unit. It returns an error for a missing
   or multi-unit class; native candidate broadening, ranking, and path-state
   production remain required for the general case.
   `ClassCatalog.LookupPrefix` ports the variable-width class-key comparator
   in `FUN_10019450`, and `selection.LookupPaul2013WholePositionByPrefix`
   ports the `FUN_10023dc0` five-to-one-byte broadening loop, 300-class cap,
   and cumulative 30-member coverage stop. It accepts classes already
   produced by an earlier state-generated query. `selection.LookupPaul2013QueryDispatch`
   now ports the `FUN_10023f90` zero-row route by appending prefix matches and
   the nonzero-row route over the reconstructed full-catalog feature-view
   index, with the native 290-class early return. An explicit scope callback
   remains available for controlled comparisons. Row-state-to-view/query
   production remains missing, and reconstructed equal-view ordering cannot
   be verified without the absent `sclass.idx` resource.
   Nested exact-prefix matches remain duplicated through coverage
   calculation and are canonicalized by the downstream shortlist.
   `selection.LookupPaul2013WholePositionByFeatureView` applies the
   10-to-1-byte `FUN_1002df50` prefix relaxation,
   10-class/10,000-member ranking, and the feature path's `>=10` coverage stop
   to a caller-supplied candidate scope. The explicit-query-mode entry point
   now follows `FUN_10023e70`: query mode 1 may relax to one byte, while modes
   0 and 2 stop at six; the earlier wrapper remains for compatibility. It
   preserves scope order on the native fast path and limits scoring to the
   first 10,000 matching classes.
   Captured P half-key probes now match reconstructed full-catalog range
   counts at widths 10, 9, and 8 in both feature-view modes for two queries.
   They exposed and corrected a swap between the low and high three-bit fields
   at byte 8 of the two feature views; broader query coverage and native
   ordering remain unverified.
   `selection.PruneCandidates` implements the measured post-transition
   cutoff, including its context gate, ten-candidate minimum, and node-flag
   bypass. `selection.SelectPaul2013ModelMinimumCostPath` now connects local
   candidate costs to record-backed `FUN_10018c80` transition scoring, native
   cumulative pruning, and `selection.Backtrack`. The higher-level
   `selection.SelectPaul2013QueriedCandidatePath` composes accepted/fallback
   class queries, unit expansion, continuity ranking, bounded local scoring,
   shortlist finalization, transition scoring, and path selection. It still
   requires caller-produced position-query results, state-derived local-score
   contexts and key mappings, per-edge transition states, and pruning gates.
   `SelectPaul2013QueriedCandidatePathFromTransitionStates` derives the one-/two-row
   mode sequence from the corresponding six-byte states before invoking it.
   The current text frontend does not produce those inputs, and
   reconstructed class-query ordering lacks verification against `sclass.idx`;
   full text-driven candidate-list parity remains open. See the Stage 7
   postselection and transition pseudocode reports.
   `selection.ExpandClassRecords` now flattens retained class members in
   ranked order under an explicit unit cap. It does not generate
   context-conditioned classes, node flags, duration data, or scoring inputs.
   `selection.ExactContextUnitCandidates` connects a supplied seven-byte
   context to its exact class and returns bounded unit references.
   `selection.LookupPaul2013PackedContextCandidates` composes that exact-key
   path from prepared native phone records and preserves an empty pool for
   each catalog miss. The five-to-one-byte broadening of an already-produced key is now available
   through `LookupPaul2013WholePositionByPrefix`. Text-derived feature-query
   generation and row-state-to-view selection remain pending.
   `selection.ReadCandidateRecords` loads their indexed unit records without
   reading waveform payloads, preserving candidate order.
   `selection.ModeTwoCategoricalPenalty` implements the mode-2 branch's
   observed same-unit and signature-byte-2 costs.
   `selection.ModeOtherCategoricalPenalty` implements the DLL's other
   signature-table branch. `selection.ScoreTransition` now looks up raw
   and generated feature distances and combines the recovered category rules,
   coefficient row, mode divisor, context predicates, cumulative cost, and
   duration term. `selection.Paul2013TransitionWeightRow` derives the local
   coefficient row from two six-byte transition state records, and
   `ScoreTransitionFromContextStates` feeds that result into the scorer.
   `selection.Paul2013PhoneRowHasTransitionFlag` ports `FUN_10017010`'s
   byte-`+2` lookup through the 256-byte DLL table at RVA `0x7ba44`;
   `selection.ApplyPaul2013TransitionStatePhoneFields` applies
   `FUN_10024680`'s byte-0, byte-1 period scan, byte-2, and byte-5 writes when
   given the record's pointer target. `selection.ReadPaul2013TransitionPeriodPhoneSpans`
   extracts byte `+0x1d` from its `0x1e`-byte rows, preserving fields written
   later by query/fallback calls. The lookup is visible in the direct
   disassembly at `0x10017010`; the call and surrounding writes are in the
   `FUN_10024680` pseudocode in `tools/revkit/work/reports/stage7-backtrack.txt`.
   `selection.BuildPaul2013TransitionContextStates` composes these fields with
   the observed query/fallback writes and compact one-row/two-row advancement.
   `selection.ReadPaul2013TransitionEligibilityShorts` reads the native
   signed-short gate from workspace offset `+0x65c`, indexed by the explicit
   descriptor group ordinal and phone index at a `0x1e0`-short group stride;
   its adapter applies the positive-short condition to ordered phone inputs.
   `selection.BuildPaul2013TransitionContextStatesFromEligibilityArena`
   composes that read with the phone-state writes and compact query/fallback
   advancement when record bytes, pointer targets, and query results are
   supplied.
   The arena writer and mapping from text to these shorts remain unresolved.
   An eligible phone consumes `QueryPaul2013Position` directly for the
   accepted/fallback decision. Mapping from compact states to candidate edges
   also remains unresolved.
   `selection.BuildPaul2013TransitionScoreInput` selects the current and
   previous metric groups and their two feature-code pairs from parsed unit
   records for modes 0–2, and fills signatures plus the coefficient row from
   explicit state records. Duration, context predicates, resolving the
   pointer-backed period spans, and mapping phone-row writes to per-edge state
   rows remain unresolved. The Stage 20
   zero-subtotal shortcut is also applied for same-unit mode-2 edges and
   marked consecutive edges. `voice.Paul2013.GlobalUnitOrdinal` maps bank-local
   references through the ordered `dblist.idx` ranges when a model is attached
   to `TransitionScoreInput`, enabling cross-bank adjacency. The mapping is
   supported by matched-index Kate runtime ranges and the documented Paul
   bank order; direct Paul cross-bank parity remains untested.
   `selection.RankPaul2013CandidatesByContinuity` ports the metadata pass in
   `FUN_100230a0`: it applies mode-dependent same-ID/adjacent-ID traversal,
   signature marker bit `0x80`, candidate-layer membership, weighted side
   coverage, and the initial sort key. It requires already-expanded candidate
   layers and per-position modes, so it does not yet produce those inputs or
   implement the subsequent context-conditioned local-score pass.
   `selection.Paul2013CandidateMembershipFlags` ports the first `+0x12`
   protection-flag pass in `FUN_10023350` over the initially scored prefix.
   `selection.ExpandPaul2013CandidateProtectionPrefix` ports the following
   large-whole-position scan: it copies matching tail candidate identity and
   span fields into the expanded score prefix, retains the destination slot's
   ordering key, and stops at the first ordering-key change after ten protected
   candidates. Unit-to-key mapping and the whole-position pool count remain
   caller inputs.
   `selection.EvaluatePaul2013WholePosition` ports the post-lookup part of
   `FUN_10024060`: it sorts class IDs numerically and removes duplicates via
   `FUN_10024010`, then applies the 30-class and
   10,000-unit shortlist through the recovered ranker, sums member counts as
   the class metric, and applies the strict `>9` acceptance rule with the
   row-class-8 positive-subtotal exception. Captured class weights match their
   source member counts. `selection.LookupPaul2013WholePosition` still accepts
   an ordered list of supplied query variants.
   `selection.QueryPaul2013WholePosition` now connects the evidence-backed
   query producer and catalog lookup for one
   whole-position row, including model-class-12's masked mode-1 pass followed
   by the original mode-2 pass, candidate accumulation, and the post-lookup
   shortlist. `selection.BuildPaul2013WholePositionRecordInput` now reads the
   target signature at three-record-block `+0x6e2 + phoneIndex*7`, row class
   at `+0x6de`, and model class at `+0x92b` from the contiguous arena selected
   by the existing six-byte workspace row; it also applies the native writes
   to row bytes `+3` and `+4`. The record-field projection is
   connected to the query by `QueryPaul2013WholePositionFromRecordState`.
   `QueryPaul2013WholePositionFromTransitionState` now accepts the six-byte
   row built by `FUN_10024680`'s port: the backtrack pseudocode writes those
   rows at `+0xec628`, the exact base used by `FUN_10024060`'s
   `(position*3 + 0x76314)*2` expression. Its row byte `+0` block ordinal and
   byte `+2` phone ordinal therefore flow directly into the record-backed
   query. `QueryPaul2013PositionFromTransitionState` now carries that same
   projection through `QueryPaul2013Position`, so a rejected whole-position
   result enters the native two-row fallback with its record-derived base row
   and model class. The fallback row-class delta is derived from the signed
   byte at record `+0x6de`, the same field written to query-row byte `+3`.
   `QueryPaul2013PositionFromStateArenas` also reads the prior-metric byte and
   live context-gate byte using the selectors at the native pointer target
   passed through state `+0x4c`. The pointer target's contents and
   eligibility-driven row ordering still need upstream production; text-derived
   state remains open.
   `selection.Paul2013QueryRulePairs` selects the captured
   `FUN_10018770` left/right rule lists and applies their count-remainder and
   context-row gates for explicit signature/state inputs.
   `selection.Paul2013QueryBaseAttempts` applies the first seven-byte
   signature mutations and packed rule code. The DLL category maps now drive
   `selection.Paul2013QueryCategoryFollowups`, including the source-restored
   third lookup branch and the state-byte-gated middle-category branch.
   `selection.Paul2013QueryRecursiveSignature` applies the recursive target
   map and mode/local-gate conditions. Callers still supply lookup results and
   the state byte at `+0x8ea`; `selection.RunPaul2013QuerySequence` now joins
   these mechanics through a caller lookup callback, tracks signed-short
   subtotals and the function return gate, and runs the single recursive
   mode-zero pass. `selection.RunPaul2013CatalogQuerySequence` connects that
   control flow to catalog lookups and carries candidates across query passes.
   `voice.ClassCatalog` now reconstructs both full class-ID vectors sorted by
   the two 10-byte feature views, and `LookupFeatureViewPrefix` performs the
   native-width binary range lookup from `FUN_1002df50`. The native reader
   loads these vectors from `sclass.idx`; the local Paul package does not
   contain that resource, so the Go reconstruction orders equal views by
   ascending class ID. Native ordering remains unverified because the local
   Paul package lacks `sclass.idx`; the two captured P half-key range-count
   discrepancies are resolved. `selection.LookupPaul2013WholePositionByFeatureViewCatalog`
   now consumes these indexed ranges through the native prefix relaxation and
   shortlist limits. Per-position state and query-mode production remain
   unimplemented.
   `text.ApplyPaul2013PositionStateRanges` ports the sorted-boundary sweep in
   `FUN_10022970` for state arrays 0, 1, 2, and 7: negative interval values
   carry the destination's prior value, ordinary modes clamp to supplied
   bounds, and mode 7 preserves its `-1` sentinel. The caller's row-key,
   boundary, and value producers remain explicit inputs.
   `text.NormalizePaul2013PositionStateProducerSeries` also ports
   `FUN_1001d5d0`'s mode 0-2 paired-boundary compaction and negative-value
   fill after those arrays have been populated. The table matcher and source
   coordinate producers remain explicit. The
   `text.ApplyPaul2013PositionStateProgramWithProducerArrays` entry point
   connects this cleanup to the existing `FUN_10022970` range passes; see
   `tools/revkit/work/reports/stage142-position-state-producer-pairing.txt`.
   `text.BuildPaul2013PositionIntervals` also ports `FUN_10022dc0`'s
   base-coordinate conversion from parser start/end offsets into inclusive
   event intervals. `text.BuildPaul2013PositionIntervalsFromParserRows`
   composes that conversion with the ordinary parser-row projection while
   retaining the caller's segment base and the native 100-row limit. The
   upstream source-row control and range/event table producers remain open.
   `text.MapPaul2013PositionStateIndexes` ports its separate
   start/end table lookups, ordinary index clamps, and optional final remap.
   `text.ResolvePaul2013RepeatedStateValue` ports the following
   `FUN_10026750` scalar transition: it resolves a negative current sentinel
   through the primary/secondary values, zero flag, or engine default, then
   clamps positive values above 65534 to 65535. These values retain their
   observed branch roles without inferred semantic labels. The parser's row,
   boundary, and table producers remain caller inputs, so this does not close
   text-derived position-state or query production.
   `text.InitializePaul2013ModelParserState` ports `FUN_1003e210`'s bounded
   state clearing, while `text.Paul2013ModeProbeHasWord` ports the null and
   pointed-word check in `FUN_1005f2c0`. `text.DispatchPaul2013ModelParser`
   joins them to `FUN_1003d350`'s 100-row sentinel initialization, mode
   dispatch, and boolean return conversion. Ghidra pseudocode for
   `FUN_1003e240` shows `-1` on its initial error path and when its
    processed-row count differs from the expected count; `FUN_1003d350`
   converts any nonzero low short to true. The Go `Success` field preserves
   that native conversion and does not imply parse validity. The two inner
   parsers still require injected callbacks. `text.RollbackPaul2013ModelParserCandidateRow`
   ports the repeated `FUN_1003d3d0` handler-miss restore: it resets the
   committed row count and source cursor, clears the rejected 0x94-byte row,
   and restores its `+0x0c` sentinel. The primary handler cascade remains
   incomplete; the repeated writes are documented in
   `tools/revkit/work/reports/stage6-tpp-callers.txt`.
   `text.FinalizePaul2013ModelParserRows`
   now ports the post-parser row grouping and caller-arena writes in
   `FUN_1003e240`, including the model-table boundary comparisons, delimiter
   truncation, surface and phone-string copies, counters, flags, and return
   value. It requires the complete surrounding caller arena and the model
   context-table bytes beginning at `model+0x429a2`; the initial
   `FUN_1000e2f0` call and source-row production remain outside this helper.
   Runtime captures currently expose selected parser-row fields, not the full
   caller arena, so those output writes have static disassembly support but no
   full-arena runtime comparison. Dispatch still accepts a finalizer callback.
   `text.FinalizePaul2013ModelParserRowsAndBuildModelStateRecords` connects
   these grouped +0x66 phone strings to `FUN_10016c90` record construction.
   Its record-source view begins at the native parser-state base, rather than
   the earlier source-row slice. It transfers the signed count and inclusive
   intervals from +0x14/+0x18 before optional position-index remapping and
   derives the terminal pitch count from finalized state +2. Per-token state
   values and pointer addresses remain caller inputs. The finalizer's U/Y
   copy predicate now matches the x86 OR branch; see Stage 157.
   `text.RunPaul2013ModelTextParser` now ports `FUN_1000e2f0`'s wrapper order:
   empty-input reset, phone-row parse, conditional normalization and row
   processing, the `FUN_1000e990` context-table flag scan, and final TPP parse.
   Negative phone-parser and TPP results follow the native `-1` return paths.
   The `FUN_1000e570` classifier now ports its ordered recursion, literal
   comparisons, membership branches, and P/N/D result flow. Its six sorted
   rule tables, 22 comparison literals, and 256 signed-short weights were
   extracted from the read-only DLL; the weights match the already-provisioned
   low-byte character map for all 256 entries. `text.UpdatePaul2013ModelContextTableFlag`
   uses this classifier by default while retaining an optional override for
   controlled comparisons. `text.NormalizePaul2013ModelPhoneRows` now ports
   `FUN_1000ea20`'s `0x554` token-row to `0x70` context-row writes, including
   parser-row class and type, metadata, and single-component phone-string
   fields. `RunPaul2013ModelTextParser` uses this projection with supplied
   parser rows or derives rows for a supported single sentence segment.
   `text.SelectPaul2013ModelPronunciationAlternative`
   ports the tree-scored path-group selection when supplied with the loaded
   `engbi.tree3` through `PronunciationTree`; a custom selector may override
   it. It ports the short-input fast path: for fewer than two source bytes,
   surfaces other than mapped string `a` copy the first phone alternative;
   `a` continues through contextual scoring. `DAT_10077718` was read as `a`
   from DLL `.data` at RVA/file offset `0x77718`.
   `text.ReadPaul2013ModelPronunciationRows` now validates and snapshots the
   surface, alternative count, terminated path-control stream, and aligned
   phone strings from the native 0x554-byte rows. The tree scorer consumes
   this shared row projection, which also supplies the exact fields read by
   `FUN_100049b0`. `text.SelectPaul2013ModelNameContextPronunciation` ports
   the recovered direct word branches, including `anti`, plus the `August`
   shortcut after `an`/`the`, beside a comma-initial row, at the final row,
   before terminal punctuation in the penultimate row, before a class-4 surface, or beside a
   one-byte surface with attribute mask `0xc0` (path byte 0x1e). For article
   `a`, the observed early direct exits in `FUN_10003ff0` now return the
   one-byte 0x1e phone for a preceding `an`/`the`, a following class-4 row,
   the recovered class-2/3 and word-pair gates, the comma-pair condition, the
   previous-row character-attribute predicate, adjacent one-byte rows with
   attribute mask `0xc0`, the final row, or a penultimate row before `.?!,`.
   It returns the one-byte 0x07 phone when none match. The control flow and
   direct stores are in `FUN_10003ff0` at `0x10003ff0`, particularly the
   return paths at `0x100042e4` and `0x100042d5` in
   `tools/revkit/work/reports/vt_pau-objdump-disassembly.txt`. This static port
   has no runtime choice comparison yet. `anti` selects path byte 0x0e when
   followed by a class-4 surface and 0x13 otherwise; `conflicts` followed by `with` selects
   `'*'`; `supply` followed by a class-3 surface selects 0x13; `does` after
   `much` or before `not` selects `*`; `converts`
   selects 0x16 before `from` or after `the`/`a`, and `'*'` otherwise;
   `wound` followed by `up`, `down`, or `through` selects `&`; and `minute`
   selects 0x0e before `amount` or 0x13 in `the`/`a` and `minute by minute`
   contexts. `number` selects 0x0f when preceded by
   `less`/`more` or followed by `with`, and 0x13 otherwise; `perfect` before
   sentence punctuation selects 0x0e; `polish Jewish` selects 0x0e; and
   terminal or punctuation-followed `learned` selects `&`; `present` with
   `to` two rows later selects `%` when the native extra-row bound is met.
   `can` selects 0x13 after `steel`, 0x12 before `be`, 0x13 for its recovered
   class-12/context-table neighbor and punctuation gates, and 0x12 otherwise.
   The supported `close` cases select `%` after `I`/`you`/`we`/`they` or the
   `to a` sequence, 0x0e after `feel` or before `than`/`by`/`on`/`to`/`upon`/
   `call`/`calls`, and 0x13 after `fisherman`/`fishermans`. `close at hand`
   and `close in with` select 0x0e; `close it up` and `close down on` select
   `%`. Other `in`/`up` contexts continue to the generic fallback, following
   the ordered bounds and comparisons in the decompilation. The `closer`
   branch selects 0x0f before
   `than`/`by`/`on`/`to`/`upon`/`call`/`calls`/`look`/`looks`,
   for `at hand` or `in with`, and when the following surface returns class 2
   from `FUN_10007160`. Those are the direct predicates in the recovered
   `closer` branch; misses fall through to the shared `FUN_10004300` path.
   Runtime selector parity remains unverified. The classifier's
   paired `FUN_1000ffa0`/`FUN_1000ffd0` character-table gates are ported with
   the native signed-byte table view.
   The `mouth to mouth` branch selects 0x13; other mouth contexts use the
   generic fallback. `brain reading` selects 0x0e; other reading contexts use
   the generic fallback. `record`/`records` select 0x16/0x13 when preceded by
   a number word or followed by `it`/`to`; preceding the ten-entry
   `FUN_100786f0` set selects `%`. `dogged` searches up to three preceding
   rows for the ten recovered `'ve`/`had`/`has`/`have` phrase forms and selects
   0x0e or `(` using the match-position and preceding class-12/`that` gate.
   It is static-only; see
   `tools/revkit/work/reports/stage122-dogged-pronunciation-selector.txt`.
   `invalid` selects 0x0e after class 3, or after class 3 followed by `n't` or
   `not`; other contexts use the generic fallback. See
   `tools/revkit/work/reports/stage123-invalid-pronunciation-selector.txt`.
   `lied` selects `&` when lowercase `l` follows uppercase, or after its
   following-word/`FUN_10010730` gates; `have`/`had`/`been` predecessors select
   `(` at the remaining gate.
   See `tools/revkit/work/reports/stage124-lied-pronunciation-selector.txt`.
   `refuse`/`refuses` before `to` select
   `*`/`%`. `import` after `export`/`the`/`an` or before `bank` selects 0x13.
   `increase`/`increases` after `the`/`an` or before a class-3 surface select
   0x16/0x13 respectively; other contexts use the generic fallback.
   `transform`/`transforms` before `into` two or three rows later select `*`/`%`;
   `use` after `still` selects `%`; `wind` before `up`/`down`/`through` selects
   `%`; `de` beside an uppercase-initial surface selects 0x0c; `house` beside
   an uppercase-initial surface or before `arrest`/`number(s)`, and `job hunting`
   select 0x13; `laden` selects 0x14
   before `'s`, after `bin`, or at the start with an uppercase initial, and
   0x0e otherwise; `dove into` and `elaborate into` select `&`; see
   `tools/revkit/work/reports/stage125-elaborate-pronunciation-selector.txt`.
   `lives on` selects `*` except
   after `put`/`lay`, where it selects 0x16; `live on TV` selects 0x0e.
   `live(s) it up`, `live(s) up to`, and `live(s) to` a following surface
   longer than `self` that ends in `self` select `%`/`*` by singular/plural form.
   The recovered following-literal set and a following `FUN_10007160`
   feature-4 surface select `%`/`*` for `live`/`lives`; the recovered previous-surface
   exception set selects 0x0e/0x16. Other live contexts use the generic
   fallback. See
   `tools/revkit/work/reports/stage128-live-pronunciation-selector.txt`.
   `resume`/`resumes` after two rows and before sentence punctuation select
   `*`/`%` only when the previous row is outside class 3 and the row two
   positions back is class 11. See
   `tools/revkit/work/reports/stage126-resume-pronunciation-selector.txt`.
   `separate`/`separates` after `be` or a predecessor whose triplet feature is
   3 select 0x0e; before `from`/`into`/`out`/`up`, or those first two plus
   `and` two rows later, they select `*`/`%`. See
   `tools/revkit/work/reports/stage127-separate-pronunciation-selector.txt`.
   `subject to` selects 0x0e after a class-3 predecessor or after a class-3
   predecessor two rows earlier with `job`/`not` between; otherwise it selects
   `%`. `contest`/`contests` followed by a class-3 surface select 0x16/
   0x13 respectively; other contest contexts use the generic fallback.
   `interstate` selects 0x0e before another surface whose initial byte passes
   the native high-bit continuity-table predicate; `interstate to` and
   `interstates to` select `%` and `*`, respectively.
   Additional generic `FUN_10004300` gates are ported from its static
   pseudocode: before `of` (unless the current surface is `nice`), try 0x13,
   0x16, 0x14, then 0x15; before `'s`, try 0x14, 0x15, 0x13, then 0x16;
   then `de`/`inter`/`re` before a same-index hyphen row select 0x0d. After `less`,
   `more`, `so`, or `very`, prefer the first class-2 path alternative, then
   class 8. After the recovered 20-word adverb set, try `(` when the preceding
   triplet returns class 3, then `%`, `*`, and `&`. A following `it`, class-12
   or class-7 word, or one of `him`, `me`, `them`, `us` activates the `%` gate
   when the class-10 cross-alternative guard from `FUN_10010900` passes.
   `FUN_10010690` preceding-word matches select the first class-4/7 path, then
   0x13, subject to the class-5/6 cross-alternative guard. Other direct
   `FUN_100049b0` branches and runtime choice parity remain open. The `have`/`has`/`had`
   branch tries `(`; `FUN_10010760` verbs before a surface longer than five
   bytes ending in `ing` select the apostrophe path. Its surrounding pair gate
   uses the recovered `FUN_10010840` predicate; a matching prior pair selects
   `%`. A `how` context immediately
   before the current row or two rows earlier selects 0x0e when another row
   follows and the current surface is class 8. A class-7 preceding triplet
   takes the same guarded class-4/7-or-0x13 path as the `FUN_10010690` gate.
   Evidence:
   `FUN_10004300`, `FUN_10010690`, `FUN_10010700`, `FUN_10010730`,
   `FUN_10010760`, `FUN_10010790`, `FUN_10010840`, and `FUN_10010900` in
   `tools/revkit/work/reports/stage6-phone-record-consumers.txt` and
   `tools/revkit/work/reports/vt_pau-objdump-disassembly.txt`.
   The direct `close` selector returns `%` after the preceding pair `to a`,
   from the branch in `FUN_100049b0`'s Stage 6 pseudocode and the string bytes
   at DLL VAs `0x10077380` and `0x10077718`; this case is static-only. Other
   `close` cases fall through to the generic fallback. The scorer tries
   direct branches first; when their context gate or selected path does not
   match, native control flow calls `FUN_10004300`. Other `FUN_100049b0`
   name/context branches remain open. The
   statically recovered `FUN_10004300` fallback cascade is now represented,
   including the ordered path selectors and cross-alternative guards, but has
   no end-to-end runtime choice comparison.
   `text.FindPaul2013PronunciationForPathCode`
   ports the alternative scan in `FUN_10010640`: each path byte and requested
   code are mapped through `DAT_100783ec`, with first-match ordering and the
   native comparison of mapped `-1` entries preserved. It returns the aligned
   phone string and alternative index. Thirty-nine direct `FUN_100049b0` word
   branches now call this helper; remaining direct word/context cases in
   `FUN_100049b0` are unported. The generic `FUN_10004300` fallback cascade is
   statically represented, but no direct runtime row/choice comparison has
   been made. The phone-row
   parser, row processor, and TPP
   parser remain callbacks, so source parsing is still incomplete. The extraction script and
   JSON evidence are `tools/revkit/scripts/extract_paul2013_context_classifier_tables.py`
   and `tools/revkit/work/reports/stage28-context-classifier-tables.json`;
   the decompiler pseudocode is in
   `tools/revkit/work/reports/stage28-context-classifier.c`. No direct
   row-by-row runtime comparison has been made for this Go classifier.
   `text.RunPaul2013ModelParserFinalizer` composes that wrapper with the row
   finalizer using the state `+0`, `+0x14`, and `+0x39ec` fields and the model
   table at `+0x429a2`; its `Callback` method adapts the complete
   `FUN_1003e240` stage to `DispatchPaul2013ModelParser`. It returns `-1`
   without finalizing when the inner wrapper fails. Its remaining phone-row
   parser, row processor, and TPP parser dependencies are callbacks; the
   normalizer uses the built-in projection when caller-arena parser-offset
   rows are available. The classifier override is optional, and the outer
   selected parser callbacks remain caller supplied.
   `text.BuildPaul2013ModelTokenRows` projects resolved lexical tokens through
   the recovered `FUN_1000d450` row writer and validates their explicit
   parser-row indexes. `BuildPaul2013ModelTokenRowsFromParserRows` now covers
   the source-row traversal and token-index association for initialized
   `0x94` rows, while the parser-row producers and helper cascade remain
   upstream requirements.
   `text.BuildPaul2013ModelPhoneRows` composes that writer with the fixed
   `FUN_1000ea20` projection and the loaded-tree pronunciation selector for
   already-resolved tokens. `BuildPaul2013ModelPhoneRowsFromParserRows`
   derives the token-to-parser-row association from the native traversal and
   projects its result; neither entry point runs the normalization/TPP cascade.
   `Paul2013DictionaryPhoneRows.WithPaul2013ModelSourceType` ports
   `FUN_1000fd20`'s result-byte classification for an explicitly supplied
   transformed source string. `TransformAndClassifyPaul2013ModelSource`
   composes it with the copy and trailing-pair split branches of
   `FUN_1000d640`. The source-type convenience wrapper remains ASCII-only;
   `TransformAndClassifyPaul2013ModelSource` uses the full table on supported
   non-scanner byte paths. The scanner gate
   (`param_4 == -1`, empty-associated-text, zero-count) now supports the
   unpunctuated ASCII path: it trims leading spaces, copies source bytes in
   order, and clears the consumed source. Terminal punctuation and terminal
   hyphens split at the native final-row boundary; non-final punctuation uses
   an optional embedded-token lookup callback and fails closed when that
   callback is required but absent. `Paul2013UnsignedCharacterAttributeTable`
   exposes the raw 256 bytes at `DAT_1007e188`. Signed-char consumers use
   `Paul2013ExceptionCharacterAttributes`, which maps high-bit bytes through
   the verified zero prefix before the table. Scanner contraction handling
   now copies ordinary high-bit bytes and preserves the directly compared
   `0xa2 0xfe` split (Stage 108). The ASCII convenience wrapper remains
   strict. This behavior is static-only and has no row-for-row runtime
   comparison. Evidence is in
   `tools/revkit/work/reports/stage108-d640-high-byte-scanner.txt`, plus
   `tools/revkit/work/reports/stage37-character-attributes.c`. The
   pseudocode excerpts are
   `tools/revkit/work/reports/stage32-model-source-type.c`,
   `tools/revkit/work/reports/stage33-model-source-transform.c`, and
   `tools/revkit/work/reports/stage35-model-source-scan.c`, and
   `tools/revkit/work/reports/stage36-model-source-terminal-splits.c`.
   `ComparePaul2013MappedCString` and
   `FindPaul2013SortedCString` are the reusable `FUN_1001c2c0`/
   `FUN_100560a0` primitives used by the classifier implementation.
   `text.DispatchPaul2013ModelSourceRule` ports `FUN_1002e990`'s longest
   matching-key dispatch, mapped-byte comparison, selected callback, and
   cursor restore on miss. `text.Paul2013ModelSourceRuleDescriptors` now
   exposes the 31 ordered keys and observed handler/auxiliary addresses from
   the table at DLL VA `0x1007e888`; `BindPaul2013ModelSourceRules` combines
   those entries with explicit implementations and fails closed when any
   handler is absent. The generic table is therefore available, while most
   tag handlers remain unported. `text.Paul2013InlinePauseModelSourceRule`
   connects the table's `<vtml_pause` entry to the exact supported
   `<vtml_pause time="N"/>` parser, emits a lexical separator, and reports
   the timing event through a callback. The disassembly shows that native
   handler `0x1002fb10` forwards to the shared tag routine with selector 6;
   the Go adapter covers only the separately observed narrow lexical form,
   not that routine's settings, marker, or parser-state effects. Other tag
   handlers and general dispatcher-to-parser integration beyond this pause
   adapter remain unported. `text.Paul2013InlineMarkModelSourceRule` now
   binds the captured `<vtml_mark` handler and emits kind-1 named, kind-2
   unnamed, or kind-3 truncated-name events. It preserves case-insensitive
   `name`/`NAME`, the observed empty-name no-record case, and the 511-byte
   inline-name boundary while leaving frame positions to the caller. Other
   mark attributes, non-ASCII names, and dispatcher-to-audio integration
   remain unsupported. See Stages 143 through 145 and Stage 150:
   `tools/revkit/work/reports/stage143-model-source-longest-rule-dispatch.txt`,
   `tools/revkit/work/reports/stage144-model-source-rule-table.txt`, and
   `tools/revkit/work/reports/stage145-inline-pause-rule-adapter.txt`, and
   `tools/revkit/work/reports/stage150-inline-mark-rule-adapter.txt`.
   `text.ExpandPaul2013CapturedVTMLSubstitutions` expands the exact captured
   `<vtml_sub alias="...">surface</vtml_sub>` shape by replacing its body with
   the ASCII alias. Stage 16 showed this fixture synthesizing byte-identically
   to the alias text and differently from the body alone. Other attributes,
   escaped/non-ASCII aliases, nested markup, malformed forms, and original
   source-span mapping are unsupported; see Stage 151.
   `LexiconFrontend.ResolveTextWithPaul2013CapturedVTMLSubstitutions` connects
   this expansion to embedded-dictionary lookup and remaps token spans to the
   contiguous alias bytes in the original source. It declines any token whose
   expanded span crosses a rewritten boundary; other markup and alias-source
   semantics remain unsupported. See Stage 152.
   `LexiconFrontend.ResolveTextWithPaul2013InlineMarks` connects the captured
   mark forms to embedded-dictionary resolution and carries their kind, name,
   and source offset through explicit pronunciation selection. Tokens can
   enclose an inline mark in their original-source span; output-frame
   placement and audio integration remain open. See Stage 153.
   `LexiconFrontend.ResolveTextWithPaul2013CapturedVTML` now composes the
   captured pause, mark, and substitution forms in one bounded source pass.
   It resolves rewritten text through the embedded dictionary, maps token
   spans back to original bytes, preserves pause boundaries and mark events,
   and rejects tokens crossing substitution boundaries. The pause path retains
   its observed separator insertion. `duration.Engine.EvaluateWithCapturedVTML`
   carries those tokens through the loaded pronunciation classifier and
   duration/pitch trees. This composes the independently captured source forms
   but does not establish mark-to-frame mapping or insert pause frames into
   rendered speech. The combined parser has local dictionary/runtime-resource
   tests; see Stage 154.
   A follow-up decompilation of the remaining `FUN_1000e2f0` callees shows why
   those remaining callbacks are not a raw-string-only gap: `FUN_1000d190` consumes
   a caller-arena parser-row array at `source+0x2c` with a `0x94` stride and
   writes `0x554`-byte model token rows; `FUN_1000ea20` projects those rows
   into the counted `0x70` context table and calls `FUN_100068b0` for
   multi-alternative rows (the fixed projection and tree-scored ordinary path
   selection are ported, with native special cases still open); `FUN_10007520` runs the ordered normalization,
   exception, and TPP cascade; `FUN_1000dfc0` drives the TPP window matcher and
   row updates. Their pseudocode is preserved in
   `tools/revkit/work/reports/stage29-model-text-parser.c` and
   `tools/revkit/work/reports/stage29-parser-helpers.c`. The row-to-string
   slice is composed in `TransformPaul2013ModelParserSourceRow`, and the
   token traversal/dictionary-to-context projection are available for
   initialized parser rows. A complete `LexiconFrontend` path still needs the
   upstream parser producers and row-control fields, the ordered
   `FUN_10007520` cascade, the TPP row updates, and full caller-arena wiring.
   `text.BuildPaul2013OrdinaryParserOffsetSegments`
   derives exclusive offsets for ASCII words and captured number forms.
   `text.BuildPaul2013OrdinaryParserOffsetRows` adds the captured raw row-type
   discriminator for ASCII words (`A`), normalized numeric rows (`D`), the
   leading sign, final currency, and telephone hyphen rows (`S`), and percent
   suffix rows (`A`). It also populates `+0x34` with ASCII token surfaces and
   expanded words from the bounded numeric normalizers, matching the Stage 20
   runtime examples. Broader numeric grammar behavior remains inference. The
   discriminator behavior is bounded to captured classes; it does not
   construct full `0x94` parser rows or finalizer state.
   `text.WritePaul2013ParserOffsetRows` now projects those offsets and the
   discriminator into fields +0, +4, and +0x24 of initialized `0x94` parser
   rows, preserving other controls. When `Text` is supplied, it also writes
   `strlen(Text)` at +0x08 and the NUL-terminated string at +0x34; optional
   `AuxiliaryText` is written at +0x52. These match the native source-row
   append layout and strings consumed by `FUN_1000d190`. Stage 20 traces show
   expanded numeric lexemes at +0x34 (`one`, `point`, `two`, `five`, `dollars`,
   `percent`, and `PM`) while several rows can retain the same source span.
   The ordinary row builder uses the bounded number normalizers to populate
   `Text` with expanded lexemes; broader number grammars remain inference. It
   does not produce the auxiliary string at +0x52. The native reads are
   recorded in
   `tools/revkit/work/reports/stage34-model-parser-source-strings.c`.
   Caller-supplied `Text` values over 29 bytes fail closed because they would
   overlap +0x52.
   `text.WritePaul2013ParserOffsetRowsWithTypeWrites` composes that writer
   with ordered sparse writes to parser-row `+0x2c`. It requires each row type
   to be supplied explicitly and does not infer the value from offsets or the
   `+0x24` discriminator; unresolved type producers therefore remain open.
   `text.WritePaul2013ParserOffsetRowsWithTypeGates` now composes the captured
   terminal and comma helpers plus the static dot and scan-fallback helpers
   over explicitly row-indexed scanner gates. It retains successful writes in
   input order, including repeated writes to a row, and does not infer scanner
   acceptance. The dot and scan-fallback branches have no direct runtime
   write observations; other type-producing
   branches can use `text.WritePaul2013ParserOffsetRowsWithTypeApplications`
   to invoke their existing row-type helper on an indexed row. Those
   applications must still provide source fields and native gate values read
   by the helper.
   `text.BuildPaul2013OrdinaryParserOffsetRowArena`
   composes the captured offset builder and writer for one sentence segment;
   `RunPaul2013ModelTextParser` derives that arena when callers omit it. The
   `WithTypeGates` and `WithTypeApplications` arena builders now compose
   explicit row-type producers directly into that ordinary-source path. The
   arena initializes parser-row `+0x30` to the captured untagged `0xff` mode;
   nondefault part-of-speech modes still require an explicit row control.
   `BuildPaul2013ModelPhoneRowsFromOrdinarySegments` retains the segment-local
   projection. Stage 109 adds
   `BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArena`, which
   flattens captured `.?!` segments in order, preserves each segment's reset
   offsets, and projects them once through the shared model counters. It
   requires row controls per generated source row and a fresh model arena.
   Its `WithTypeApplications` variant applies ordered row-type helpers using
   each segment's local row indexes/counts, then translates writes into the
   shared flattened parser arena. Mode bytes and auxiliary text remain caller
   controls.
   `BuildAndRunPaul2013KnownModelContextPassesFromOrdinarySegmentsInSharedArena`
   connects that projection to the known counted-row passes and reports rows
   that still need unported handlers. Its `WithTypeApplications` variant now
   accepts the per-segment ordered row type helpers before those passes.
   The corresponding `WithTPPFromOrdinarySegmentsInSharedArenaWithTypeApplications`
   entry point continues through the supported numeric F/G and WAB updates;
   component-based TPP effects remain unresolved. The
   `WithDefaultParserControls` variant also composes the captured `0xff` mode
   and empty auxiliary string with explicit initial row types and type
   applications, then runs the same counted-row and numeric TPP passes.
   Parser-control production, earlier normalization/TPP work, remaining row
   mutations, and multi-segment runtime comparison remain unresolved; see
   `tools/revkit/work/reports/stage109-model-parser-shared-segment-arena.txt`.
   `text.AppendPaul2013ModelSourceRow` ports the successful row append in
   `FUN_10044fd0`: the 100-row bound, C-string length/copy, coordinate and
   discriminator writes, low-byte flag, `0xff` marker, count increment, and
   last-coordinate update. It preserves unwritten fields, including `+0x2c`;
   source-row construction and that field's producer remain separate inputs.
   Its optional auxiliary-string form ports `FUN_10045070`'s additional copy
   at row `+0x52`. `AppendPaul2013ModelSourceRowsSplit` ports
   `FUN_100451e0`'s space/tab splitting and repeated append order; components
   longer than its 32-byte native local buffer fail closed.
   `AppendPaul2013ModelSourceRowsBackslash` ports `FUN_100452c0`'s auxiliary
   string splitting, preserving empty components and append order while
   bounding the native 68-byte temporary buffer and destination row.
   `text.AppendPaul2013MatchedComponentSourceRows` composes the accepted
   candidate follow-up in `FUN_10034110`: it reads component rows at `+0x0c`
   with `0x140` stride, applies the literal `0x41`/`0x44`/`0x12` append
   arguments, and writes the native `0x1f` byte at source-row `+0x1e` after
   each successful split. Candidate production and runtime row parity remain
   open; see `tools/revkit/work/reports/stage130-matched-component-source-row-followup.txt`.
   `text.AppendPaul2013ParserRowTypeMarker` ports `FUN_1000d190`'s fixed
   punctuation suffixes for row types 2, 3, 4, 5, 12, and 11 under its
   nonzero-type and header-value gate. The opaque header value remains an
   explicit input, and strings that exceed the native 32-byte local buffer
   fail closed. Evidence is in
   `tools/revkit/work/reports/stage58-parser-row-marker-append.c`.
   `text.TransformPaul2013ModelParserSourceRow` composes that marker with the
   observed `+0x08`, `+0x24`, `+0x2c`, `+0x30`, `+0x34`, and `+0x52` row
   fields, the supported `FUN_1000d640` transform, and `FUN_1000fd20` source
   classification. It bounds the local strings and leaves the embedded
   punctuation lookup callback explicit. The callback can now use
   `EmbeddedDictionary.ContainsPaul2013Surface`, which applies the recovered
   key transform and exact record lookup used by `FUN_10003a70` mode 0. Token
   lookup, previous-row association, and `FUN_1000d450` production are now
   composed by the Stage 61 parser-row entry point when initialized rows are
   supplied. See
   `tools/revkit/work/reports/stage59-parser-source-row-transform.c`.
   `TransformPaul2013ModelParserSourceRowWithEmbeddedDictionary` composes this
   callback with the row transform; mode-0 lookup evidence is in
   `tools/revkit/work/reports/stage60-parser-punctuation-dictionary-lookup.c`.
   `text.BuildPaul2013ModelTokenRowsFromParserRows` now traverses initialized
   `0x94` rows, repeats `FUN_1000d640` while its source remainder is nonempty,
   resolves each transformed string through the embedded dictionary, and
   writes the `0x554` result row with the observed preceding-parser-index and
   mode-byte gates. It also preserves the native empty-lookup row shape and
   200-token stop. Parser-row production and the earlier
   `FUN_10007520`/TPP cascade are still upstream requirements; no row-by-row
   runtime comparison has been made. The adjacent
   `text.BuildPaul2013ModelPhoneRowsFromParserRows` composes this traversal
   with `FUN_1000ea20`'s context-row projection, the recovered ordinary
   multi-pronunciation selector and the recovered direct `FUN_100049b0` word
   branches. The `August` shortcut selects 0x1e after
   `an`/`the`, beside a comma-initial row, at the final row, before terminal
   punctuation or a class-4 surface, and beside one-byte rows classified by attribute mask
   `0xc0`; the article `a` uses the direct `FUN_10003ff0` gates listed above
   and returns 0x07 otherwise. It also composes the statically
   recovered `FUN_10004300` fallback cascade. Additional direct name/context
   branches and runtime pronunciation-choice comparison remain open. See
   `tools/revkit/work/reports/stage61-parser-token-orchestration.c`.
   `text.BuildPaul2013ModelPhoneRowsFromOrdinarySource` connects the captured
   ordinary word/number offset-and-text producer to that composition. It
   derives `+0x08` from each generated source string and requires explicit
   per-row `+0x2c`, `+0x30`, and `+0x52` controls,
   rejects multiple sentence segments, and does not claim general text-parser
   coverage. Its evidence boundary is in
   `tools/revkit/work/reports/stage62-ordinary-source-to-model-context.c`.
   `text.BuildPaul2013ModelPhoneRowsFromOrdinarySourceWithTypeApplications`
   now composes ordered recovered row-type helpers over the supplied initial
   `+0x2c` values before token/context projection. Mode and auxiliary text
   remain explicit caller controls in that API. The new
   `text.BuildPaul2013ModelPhoneRowsFromOrdinarySourceWithDefaultParserControls`
   initializes mode `+0x30` to the captured `0xff` value and leaves auxiliary
   text `+0x52` empty for the bounded ordinary word, punctuation, and numeric
   input classes. Initial row types and recovered type-helper gates remain
   explicit; Stage 20 observed empty auxiliary strings on its inspected
   numeric and punctuation rows. Its shared-arena counterpart applies those
   defaults across `.?!` segments while preserving their native offsets and
   counters. These conveniences do not supply the unresolved `+0x2c` producer
   or upstream normalization/TPP cascade.
   `duration.Engine.BuildPaul2013ModelPhoneRowsFromParserRows` and its
   ordinary-source variant bind the loaded dictionary and pronunciation tree
   to these compositions; callers still provide the parser-row controls and
   upstream normalization/TPP results. The orchestration result retains the
   parser-row arena. `BuildAndNormalizePaul2013ModelPhoneRowsFromParserRows`
   and its ordinary-source variant compose projection with the supported
   `FUN_10009030` branches over explicit context-row indexes; the caller still
   supplies the eligibility/order list, and the earlier exception/TPP stages
   remain open.
   The Ghidra pseudocode excerpt is
   `tools/revkit/work/reports/stage30-model-source-row-writer.c`.
   `text.ApplyPaul2013ParserTerminalRowType` ports the direct `+0x2c` writes
   2/3/4 in `FUN_10051a00`, and `ApplyPaul2013ParserCommaRowType` ports its
   comma-path writes 5/12 in `FUN_100544f0`; 5 was observed at runtime and 12
   is static-only. `ApplyPaul2013ParserDotRowType` ports that function's
   separate static dot branch, which writes 1; this branch has not been
   observed at runtime. `ApplyPaul2013ParserQuoteRowType` ports the
   statically recovered double-quote branch that writes 1 after a successful
   lookup in the table at `0x1007844c`, a local row count below two, and a
   nonempty source-row list. Its preceding scanner path remains an explicit
   gate and the branch has no runtime observation. The single-quote path
   beginning at `0x100549ab` appends a source row but has no direct `+0x2c`
   write. `ApplyPaul2013ParserUnmatchedRowType` ports the separate lookup-miss
   write at `0x10054f28`, using the table at `0x100780fc`; the preceding
   scanner path and native lookup result remain explicit. Other static
   `ApplyPaul2013ParserCloseDelimiterRowType` ports the static closing
   parenthesis/bracket write at `0x100557df` for accepted final scanner
   statuses 1 through 3. `ApplyPaul2013ParserModeOneRowType` ports the static
   write at `0x10055a8e` for a positive prior scan count and mode 1; its
   preceding scanner path remains explicit. `ApplyPaul2013ParserAcceptedPathRowType`
   ports the final count gate and write at `0x10055a4a`; all earlier scanner
   and table conditions remain caller-supplied. `ApplyPaul2013ParserBoundedScanRowType`
   ports the final nonzero-and-within-limit comparison and write at
   `0x10055677`; the compared values' meanings and earlier scanner path remain
   unresolved. The `Paul2013ParserPreviousRowTerminalPunctuation` helper
   ports the
   `FUN_10062f10` check against the DLL string at `0x100813bc` (`.?!;`), and
   `Paul2013ParserTerminalGateFromPreviousRow` derives the zero-result gate
   used by `FUN_10051a00`. Its whitespace-prefix variant composes the
   two-or-more-LF status 8 from `FUN_1005a350` with that predicate. Other
   scanner statuses and the alternate accepted `FUN_10051cc0` result remain
   explicit inputs. The native disassembly and
   string bytes establish this predicate; they do not recover the scanner
   status producer. See
   `tools/revkit/work/reports/stage31-model-parser-row-type-producers.c`.
   For the 0x94-stride stores attributed to `FUN_100544f0`, the direct
   assignment sites have final-write helpers, but their incoming
   branch predicates are not fully recovered. A nearby byte-stride write at
   `0x1005551d` still needs destination attribution. Other parser routines
   also write the same candidate field. `ApplyPaul2013ParserQuoteCoordinateRowType`
   ports the repeated one-character quote, count-one, coordinate-sentinel
   fallback at `0x1004f426`, `0x1004f6c5`, `0x1004f867`, `0x1004fa17`,
   `0x1004fbb4`, `0x10050073`, `0x10050229`, `0x10050278`, `0x1005084b`,
   `0x10050b4e`, and `0x100510df`; their preceding scanner paths remain
   explicit. See `tools/revkit/work/reports/stage48-parser-quote-coordinate-row-type.c`.
   `ApplyPaul2013ParserInputFlagRowType` ports the separate type-5 write at
   `0x100445b1` after a successful append when the input short at `+2` is
   nonzero; see `tools/revkit/work/reports/stage49-parser-input-flag-row-type.c`.
   `ApplyPaul2013ParserAcceptedAppendRowType7` and
   `ApplyPaul2013ParserIntermediateAppendRowType7` port the post-append stores
   at `0x10057722` and `0x100577db`; the latter requires another input row to
   remain. See `tools/revkit/work/reports/stage50-parser-post-append-row-type-seven.c`.
   `ApplyPaul2013ParserHyphenModeRowType` ports the type-11 store at
   `0x10045a90` after an append when mode is 2 and the following byte is `-`;
   the mode remains opaque. See
   `tools/revkit/work/reports/stage52-parser-hyphen-mode-row-type-eleven.c`.
   `ApplyPaul2013ParserExactInputFlagRowType` ports the separate type-5 write
   at `0x10059f8b` when the input short at `+2` equals exactly one; evidence is
   in `tools/revkit/work/reports/stage53-parser-exact-input-flag-row-type-five.c`.
   `ApplyPaul2013ParserContextualRowTypeOne` ports the context-gated assignment
   at `0x10042e8f` while retaining its stack values as raw numeric inputs;
   semantic labels and the earlier scanner path remain unresolved. See
   `tools/revkit/work/reports/stage54-contextual-row-type-one.c`.
   `ApplyPaul2013ParserModeOnePriorKeyRowTypeOne` ports the separate
   `0x100430ce` assignment after a prior-row key match, mode 1, and two
   positive raw state values; semantic labels and earlier scanner conditions
   remain unresolved. See
   `tools/revkit/work/reports/stage55-mode-one-prior-key-row-type-one.c`.
   `ApplyPaul2013ParserPreAppendRowType7` and
   `ApplyPaul2013ParserSplitAppendRowType7` port the pending-row and
   post-split writes at `0x1005782b` and `0x1005788b`; the state dword meanings
   remain unresolved. See
   `tools/revkit/work/reports/stage56-split-append-row-type-seven.c`.
   A broader disassembly pass found two more `FUN_1003d3d0` mutations outside
   the immediate `-0x54` store search. `ApplyPaul2013ParserAcceptedScannerRowTypeOne`
   ports the accepted scanner's zero-field type-1 write at `0x1003d4b6`.
   `ApplyPaul2013ParserTypeTwoReset` ports the selected type-2 row clear at
   `0x1003da7a`, with its peer-row words and sentinels supplied explicitly.
   See `tools/revkit/work/reports/stage57-model-parser-row-type-one-and-reset.c`;
   semantic names and runtime observations remain open.
   `ApplyPaul2013ParserZeroInputFlagRowType` ports
   the type-8 writes at `0x10053229` and `0x100538c0` after successful
   appends when the input short at `+2` is zero; see
   `tools/revkit/work/reports/stage51-parser-zero-flag-row-type-eight.c`.
   The 0x94-stride stores identified by the displacement audit now each have
   a final-write helper, but most incoming predicates are still explicit
   caller inputs and these branches lack runtime write watches. The separate
   byte-stride candidate at `0x1005551d` remains unattributed. The audit is
   recorded in `tools/revkit/work/reports/stage46-parser-row-type-write-inventory.c`.
   `ApplyPaul2013ParserPositiveScanDotRowType` ports the extra positive-scan
   dot assignment at `0x10056077`. `ApplyPaul2013ParserScanFallbackRowType`
   maps the byte-indexed assignment at `0x1005551d` to the preceding source
   row's `+0x2c` type field using the native `0x94` stride and `+0x14` row
   base. Its scanner and lookup results remain explicit caller inputs; this
   static branch has no runtime write watch. See
   `tools/revkit/work/reports/stage58-parser-scan-fallback-row-type-one.c`.
   The cross-function address search is recorded in
   `tools/revkit/work/reports/stage46-parser-row-type-write-inventory.c`.
   `ApplyPaul2013ParserMultiSourceRowType`
   ports the separate static `FUN_10054050` write of 10 after a successful
   generic row append when the input source-row count exceeds one; the append
   result remains an explicit gate. See
   `tools/revkit/work/reports/stage39-parser-multi-source-row-type.c`; lookup-
   miss and closing-delimiter evidence are in
   `tools/revkit/work/reports/stage40-parser-unmatched-row-type.c` and
   `tools/revkit/work/reports/stage41-parser-close-delimiter-row-type.c`;
   mode-one evidence is in
   `tools/revkit/work/reports/stage42-parser-mode-one-row-type.c`, and the
   accepted-path write is in
   `tools/revkit/work/reports/stage43-parser-accepted-path-row-type.c`, and
   the boundary comparison is in
   `tools/revkit/work/reports/stage44-parser-bounded-scan-row-type.c`.
   The distinct sign-scanner write is documented in
   `tools/revkit/work/reports/stage47-parser-hyphen-fallback-row-type.c`.
   Additional `FUN_100544f0` branch evidence is in
   `tools/revkit/work/reports/stage45-parser-positive-dot-row-type.c`.
   `text.ScanPaul2013ModelParserWhitespacePrefix` ports the entry loop in
   `FUN_1005a350` that skips leading space, tab, CR, and LF, sets its single-LF
   marker, and returns status 8 after two or more LFs. This does not implement
   token scanning or derive the statuses used by the row-type branches.
   Scanner classification, previous-row eligibility, and following-token
   scanner outputs remain explicit; the comma branch now
   applies its mapped `too`/`either` comparisons itself. The static excerpts
   and direct write-watch evidence
   are in `tools/revkit/work/reports/stage31-model-parser-row-type-producers.c`
   and `tools/revkit/work/stage20/trace-model-parser-row-type-writes.gdb`.
   The Stage 20 trace also records the `+0x2c` dword consumed by
   `FUN_1000ea20` for ordinary, punctuation, and captured numeric rows.
   Values vary across otherwise ordinary `A` rows (including 0, 2, 3, 4, and
   5), while the captured numeric expansions commonly use 0 on intermediate
   rows and 4 on the final row. The observations are preserved in
   `tools/revkit/work/stage20/model-parser-row-types.tsv`, and the trace is
   reproducible with `trace-model-parser-row-type.gdb`. The full producer is
   still incomplete: incoming `+0x2c` branch predicates and type producers
   outside the enumerated sites remain explicit caller inputs and must not be inferred
   from offsets or the `+0x24` discriminator.
   Runtime rows match repeated words, comma retention, `.?!` resets, and
   source spans for the captured `+12` form, decimals, currency, percentages,
   ordinals, one slash date, a clock time, and the seven-digit telephone form.
   Row counts outside directly captured values follow the existing normalizers
   as inference. Other abbreviation, exception, and TPP parser outputs, raw
   parser-state construction, callbacks, and source lookup tables remain
   unported. `text.BuildPaul2013OrdinaryParserPositionIntervals` now composes
   those offset rows with the recovered inclusive base-coordinate conversion,
   returning one event-interval slice per sentence segment. It advances each
   segment base by that segment's source-byte origin, matching the native
   caller's consumed-length coordinate update. The caller's initial base and
   later event/table values remain explicit. See
   `tools/revkit/work/reports/stage25-parser-finalizer.c`.
   `text.BuildPaul2013PositionIntervalsFromParserState` adapts the raw parser
   state rows consumed by `FUN_10022dc0` (signed count at +0, rows at +0x14
   with 0x94 stride, signed start/end offsets at row offsets +0/+4) into the
   already-ported inclusive absolute intervals. It validates the 100-row
   capacity and byte bounds; it does not produce parser rows.
   `text.InitializePaul2013PositionStateArrays` ports the same caller's
   per-row initialization before `FUN_10022970`: arrays 0, 1, and 2 copy their
   three supplied defaults, while arrays 3, 4, and 7 start at -1. Slots 5 and 6
   are not initialized or consumed by the recovered state pass.
   `text.ApplyPaul2013PositionStateProgram` composes initialization with the
   mode 0/1/2/7 range passes and mode 3/4 event passes in native order,
   including the optional terminal mode-3 accumulator. Its optional final map
   reads signed indexes at model-state row offsets +0x64c/+0x650 with the
   native 0x3c0 stride and applies the explicit start/end/final tables.
   Parser-row, range, event, and gate producers remain explicit.
   `text.RunPaul2013ModelParserLoop` now ports `FUN_100267d0`'s repeated
   `FUN_1002c9b0` advance and `FUN_10026630` row-processing control flow,
   including its signed-short retry condition, counter writes, stop handle,
   ordinary high-word return, and full callback result on its inner stop-handle
   return. The callbacks still own text parsing and row production, and the
   state fields preserve native offsets where meaning is unknown.
   `text.PropagatePaul2013PackedContextCategories` ports the forward high-field
   and backward low-field propagation loops inside `FUN_10017510`, using the
   same opaque `DAT_1007b9e0` table already used by unit scoring. It accepts
   the native parallel byte streams explicitly.
   `text.BuildPaul2013PackedContextCategories` ports the preceding per-record
   packed-byte construction and returns the flattened phone-symbol stream and
   record/phone index map. `text.ApplyPaul2013PhoneMarkerPrepass` ports
   `FUN_10017100`'s ordered per-record code rewrites and terminal-boundary
   rewrites, using `FUN_10017030`/`FUN_10017090` and their recovered table
   windows. `text.PreparePaul2013PhoneMarkerContextGroup` composes that pass
   with packed-context construction when callers supply the selected record
   indexes. `text.Paul2013SpecialTerminalKeyMatch` ports the ten-entry
   `DAT_1007beec`/`FUN_100560a0` membership check; bytes reached through its
   pointer field and the character map remain caller inputs.
   `text.ApplyPaul2013PhoneMarkerPrepassWithPointerResolver` now reads the
   little-endian 32-bit address stored at record `+0x3b8` and invokes an
   explicit C-string address resolver only when the native terminal lookup
   branch is reached. This composes the recovered pointer field with the
   already-ported membership check while leaving address-space bounds and
   target-string ownership with the caller. The resolver now flows through
   `PreparePaul2013PhoneMarkerRecordGroupsFromArenaWithPointerResolver`,
   `LookupPaul2013PhoneMarkerArenaCandidatesWithPointerResolver`, and
   `SelectPaul2013PhoneMarkerArenaPathWithPointerResolver`, so the complete
   arena-to-exact-candidate path can use native pointer fields without a
   separate key list. Candidate-state, score, transition, and upstream source
   record producers remain explicit. `FUN_10017100`
   mutates records; its caller `FUN_10012df0` separately produces the group
   descriptors. `text.BuildPaul2013RecordGroupDescriptors` ports that producer's
   grouping over ordered 0x3c0-byte records: `\\` and `]` continue the current
   group, other terminal bytes close it after the current record, and the final
   record always belongs to a group. The descriptor retains each group's start
   and count, its native +0xe start-boundary index, the sum of each record's
   byte at +0x94, and the output-short offset derived from base 0xbea6 plus the
   preceding row-span totals times 0x0f. The +0x94 field's meaning remains
   opaque. The decompiler pseudocode is in
   `tools/revkit/work/reports/stage22-context-marker-producers.c`.
   `text.PreparePaul2013PhoneMarkerRecordGroups` connects these descriptors to
   the marker prepass and the per-group packed-context helper; the pointer-backed
   key bytes remain explicit caller inputs. This composed helper receives record
   slices rather than the native contiguous arena, so it does not apply the
   shared reset; `PreparePaul2013PhoneMarkerRecordGroupsFromArena` applies that
   reset and retains its result. The prepared group result now also carries
   each group's extracted seven-byte selection contexts in native record/phone
   order, ready for exact catalog lookup.
   `text.WritePaul2013PackedContextPhoneRows` ports the final seven-byte
   per-phone writes at `+0x96..+0x9c` using the flattened map, propagated and
   initial packed bytes, and explicit native records. `text.ApplyPaul2013PackedContextPhoneStateAdjustments`
   ports the backward `+0x25d` and next-record `+0x61d` category adjustments,
   gated by the current record's `+0x3bc` and bounded by the aliased next
   record phone count, followed by the marker-`'2'` `+0x25e/+0x25d` pass.
   `text.ResetPaul2013PackedContextSharedPhoneState` ports the initial clear at
   record offset `+0x25d` (enclosing group-block offset `+0x8a9`) when given the
   complete contiguous native record arena and record count.
   `text.PreparePaul2013PackedContextPhoneGroup` composes the packed-category
   passes over caller-selected raw records while retaining both category
   stages and updated record copies. The slice-based group helper cannot apply
   the shared reset because its input does not establish the complete native
   record stream; the contiguous-arena entry point retains the reset copy.
   `text.ExtractPaul2013SelectionContexts` reads the seven-byte rows back in
   record/phone order, and `text.PreparePaul2013SelectionContextsFromPackedContextPhoneGroup`
   connects the marker pass to the `text.Context` values consumed by catalog
   lookup. `selection.LookupPaul2013PackedContextCandidates` composes one
   prepared group with exact class lookup and bounded unit expansion. Its
   `LookupPaul2013PhoneMarkerRecordCandidates` variant also accepts the full
   ordered record stream, derives its native groups, runs the marker prepass,
   and returns exact candidate pools aligned within each group. Both retain
   exact misses as empty pools. Raw record production and pointer-backed key
   bytes remain upstream inputs. `LookupPaul2013PhoneMarkerArenaCandidates`
   additionally applies the shared phone-state reset to a contiguous native
   record arena before preparing groups. The
   `LookupPaul2013PhoneMarkerModelStateArenaCandidatesWithPointerResolver`
   entry point reads the signed record count at +2 and slices records from
   +0x64c in the full model-state arena before applying the same reset and
   exact lookup. Parser-row production, pointer address mapping, eligibility
   shorts, query state, and later score/path state remain explicit or
   unresolved. `RankPaul2013PhoneMarkerRecordContinuity`
   connects nonempty exact pools to the native continuity pass while retaining
   group alignment; it requires explicit context modes and leaves misses to the
   fallback stage. Whole-position query broadening and fallback-row production
   remain separate. `ScorePaul2013PhoneMarkerRecordCandidates` composes the
   continuity and bounded local-score passes when the caller supplies their
   state-derived contexts and candidate-key maps. `Paul2013ContextModesFromTransitionStates`
   reads mode byte `+4` from the compact six-byte transition rows, and the
   `...FromTransitionStates` rank/score entry points connect those modes when
   the rows align with candidate positions. Score contexts and key maps remain
   explicit.
   `selection.BuildPaul2013TransitionContextStatesFromRecordGroup` now maps
   the native per-record phone loop into compact states: it adds each
   record-loop ordinal to the byte-converted descriptor `+0x0c` base,
   resolves each record's `+0x08` period-row pointer through a caller
   address-space callback, and reads the corresponding eligibility block at
   one `0x1e0`-short stride per record. Position-query results remain aligned
   caller inputs, with nil required for ineligible phones. This composes the
   recovered phone-state writer, query/fallback state advancement, and compact
   output order without claiming to produce the descriptor, pointer map,
   eligibility shorts, or position-query inputs.
   `FinalizePaul2013PhoneMarkerRecordShortlists` connects those score passes to
   native tail ordering, cap, and protected-candidate restoration while
   preserving record groups. `SelectPaul2013PhoneMarkerRecordPath` carries exact
   pools through continuity, local scoring, shortlist finalization, and path
   scoring. It still requires caller-produced transition states or modes,
   score/tail state, and
   per-edge transition state; exact misses must pass through fallback first.
   `SelectPaul2013PhoneMarkerRecordPathFromTransitionStates` extracts the
   continuity modes from aligned six-byte transition rows before running that
   composed path.
   `SelectPaul2013PhoneMarkerArenaPath` includes the contiguous-arena shared
   reset, grouping, packed-context exact lookup, and those downstream stages.
   `SelectPaul2013PhoneMarkerArenaPathFromTransitionStates` composes that
   route with mode extraction from aligned six-byte state rows, including the
   pointer-resolved key variant. Score, shortlist, and edge-state inputs remain
   explicit.
   The
   `FUN_10017030`, `FUN_10017090`,
   `FUN_10017100`, and `FUN_100560a0` pseudocode is preserved in
   `tools/revkit/work/reports/stage27-record-marker-passes.c`.
   `ApplyPaul2013PositionEventValues`
   ports the per-row selector 3 sum/clamp and selector 4 last-match assignment;
   it preserves the persistent matched-event cursor and selector 3's separate
   upper-clamp cursor. `EvaluatePaul2013PositionTerminalAccumulator` ports the
   final selector-3 scan when the caller supplies its gate and interval. The
   interval and event-table producers remain unresolved.
   `selection.BuildPaul2013FallbackRows` ports `FUN_100242a0`'s two-row
   fallback byte construction and its class-12 dual tree-target dispatch.
   Runtime watchpoints directly confirm fallback row byte +4 values 1 and 2;
   static decompilation supports the byte +3 increment and target-byte mask.
   `selection.QueryPaul2013FallbackRows` connects the two constructed rows to
   their ordered tree-query sequences, numeric-ID deduplication, row-view
   shortlist ranking, and final ranked-prefix/raw-sidecar assembly. Ranked IDs
   retain scoring order; a positive mode-1 lookup from an empty candidate list
   snapshots the first 30 raw IDs, which are appended by ascending ID. The
   model-table increment at `+0x6de`, per-position prior metric, and live gate
   remain caller inputs; this does not recover outer per-position state
   production.
   `selection.QueryPaul2013Position` composes this fallback with the
   whole-position acceptance gate from the native caller and reports its
   one-row or two-row advance. `selection.ExpandPaul2013PositionQueryCandidates`
   expands the accepted shortlist or both ordered fallback ID rows into
   bounded unit-reference pools. It resolves fallback IDs only against class
   records carried through those query passes and fails on missing records.
   `selection.RankPaul2013QueriedCandidatePositions` connects the expanded
   pools to the native continuity metadata pass using caller-produced modes
   aligned to the one- or two-row query advance.
   `selection.RankAndScorePaul2013QueriedCandidatePositions` carries those
   ranked rows through bounded local scoring, including protected-prefix
   expansion for oversized pools. Caller-produced score contexts and identity
   keys must align to continuity-ranked order; tail ordering and text
   production of those values remain separate stages.
   `selection.FinalizePaul2013QueriedCandidateShortlist` connects the scored
   rows to the secondary tail multiplier, native tail sort, 30-unit cap, and
   protected-candidate replacements. It requires local costs for every
   candidate plus caller-produced candidate/context classes for oversized
   rows, so an unscored tail still fails closed.
   `selection.SelectPaul2013QueriedCandidatePath` composes those queries,
   score passes, and shortlists through model-backed transition scoring,
   pruning, and backtracking. It still requires text-derived query results,
   score contexts, class mappings, and per-edge state from the caller.
   `SelectPaul2013QueriedCandidatePathFromTransitionStates` derives modes from
   aligned six-byte whole-position/fallback transition rows; query production
   and the other text/state inputs remain upstream.
   `selection.Paul2013ContinuityFeatureScale` derives the multiplier passed to
   the `FUN_100182e0` feature term, and `ScorePaul2013CandidateLocalCost`
   applies it without scaling categorical costs. It uses float32 constants
   1.0 and 2.0 read from the local DLL at `0x1006d170` and `0x1006d174`.
   `selection.Paul2013CandidateLocalScorePrefix` implements the recovered
   full-span-first rule and the large-list score boundary, including ties.
   `selection.ScorePaul2013CandidateLocalPrefix` connects that boundary and
   oversized-list protection markers to candidate record reads,
   `FUN_100182e0` input construction, and continuity-scaled costs for the
   scored prefix. It marks only the initial native prefix; callers that need
   the subsequent protection expansion can apply the explicit expansion
   helper and realign their score contexts through its source-index map.
   `selection.ScorePaul2013ExpandedCandidateLocalPrefix` composes the
   expansion, context realignment, record reads, and local scoring. Candidate
   identity keys and score-context values remain explicit inputs; final tail
   ordering remains a separate stage.
   `selection.Paul2013TailOrderingScore` applies the category multipliers, and
   `SortPaul2013CandidateTailAndCap` sorts the supplied suffix and retains 30.
   `FinalizePaul2013CandidateShortlist` completes the oversized-list path by
   sorting the recovered suffix, capping at 30, and restoring flagged tail
   candidates into trailing unprotected slots, with the native limit of five
   protected entries. The caller still supplies the suffix boundary,
   secondary scores, and final protected count, including any later flags.
   `sortPaul2013Native` ports `FUN_1001b5f0`'s 17-entry insertion cutoff,
   median-of-three partition, strict scans, and 32:1 imbalance fallback to
   `FUN_1001b400`'s heap path; class ranking, continuity ranking, and oversized
   tail ordering now use it. A natural 75-entry local-score tail from the
   repeated-Hello run matches the complete native permutation with its exact
   float32 score bits and input order. A second 75-entry trace replaces scores
   in process memory with a pattern that forces the 32:1 heap fallback; that
   complete native permutation also matches. The vectors are
   `tools/revkit/work/stage20/native-tail-sort-order.tsv` and
   `tools/revkit/work/stage20/native-heap-tail-sort-order.tsv`; their capture
   scripts are `trace-native-tail-sort.gdb` and `trace-native-heap-sort.gdb`.
   The forced vector checks the fallback for that constructed partition only;
   it does not show that a natural workload takes this path or establish
   parity for every range size. The caller
   must still derive the suffix boundary and secondary local scores from
   context-conditioned state.
   `selection.BuildPaul2013UnitScoreInput` maps candidate feature codes from
   the reconstructed 21-byte index row and selects the columns consumed by
   `FUN_100182e0` for context modes 0–2. It also derives the signature-based
   scale and equality branches. Context feature values, context signatures,
   marker values, and the mode producer still must be supplied.
   `selection.ScoreUnitCost` applies the feature-distance scaling branches of
   `FUN_100182e0` when given that explicit context and its categorical
   penalty. `selection.UnitRecordBytePenalty` ports the direct
   byte-1 through byte-3 mismatch branches using the recovered DLL maps and
   supplied state markers. `selection.UnitByteFivePenalty` applies the
   recovered high- and low-field cost tables, with unsupported high-field
   values rejected; `selection.KnownUnitCategoricalPenalty` combines these
   recovered byte costs with in-range cross-byte table penalties. It rejects
   mapped categories outside the recovered 3-by-3 tables. Context-array and
   duration-term producers, context-state assembly, and text-derived local-cost
   production remain unimplemented; complete per-unit scoring is still not
   connected to generated text contexts. With explicit context values,
   `selection.ScorePaul2013CandidateLocalPrefix` reads candidate records,
   applies the recovered score boundary and membership flags, and computes
   local costs for the scored prefix.
   `selection.ScoreTransitionLayer` now evaluates all predecessor candidates
   for each supplied current candidate and records the winning index for
   backtracking. `selection.ScoreAndPruneTransitionLayer` then applies the
   measured cutoff while preserving predecessor references.
   `selection.InitializePathLayer` wraps caller-supplied local costs for the
   first layer. `selection.SelectPaul2013MinimumCostPath` now connects the
   initial layer, per-position all-predecessor scoring, native multiplier-1
   pruning, and final backtracking when callers provide candidate layers,
   local costs, context gates, and the edge scorer.
   `selection.SelectPaul2013ModelMinimumCostPath` connects model record reads,
   mode-specific transition construction, scoring, pruning, and backtracking.
   It applies the native `/2` local-score normalization to each candidate in
   the first and later layers, using the current normalized node score as the
   transition's local term. Candidate layers, context-state rows, and gates
   still require caller production. `selection.SelectPaul2013PathFromLocalScorePasses`
   adapts continuity/local-score output into those layers while preserving
   order and protection flags; it rejects candidates without a recovered path
   cost instead of inventing one.
   The selection package now also ports the deterministic query expansion
   around `FUN_10018770`: packed left/right rule selection, context-row gates,
   category followups, recursive lookup, and the ordered catalog-backed
   `FUN_10023f90` dispatch. This can generate and accumulate class candidates
   when the target signature, six-byte context row, prior metric, context gate,
   and model catalog are supplied. It does not produce those source/model
   inputs. The remaining selection handoff is therefore the producer chain:
   normalize source rows, build phone/context rows and seven-byte target
   signatures, derive per-position context state and duration terms, then feed
   the resulting candidate layers into the already-ported path scorer. Capture
   comparison is still needed for that composed source-to-selected-unit path.
4. **Selected units to speech — diagnostic concatenation, isolated period
   resampling, and adjacent-window reconstruction implemented; full legacy
   rendering pending.**
   `engine/synthesis.ConcatenatingRenderer` reads and
   joins complete decoded units into 16 kHz PCM under default controls. Its
   opt-in `ApplyObservedGain` flag applies the normalized 0–500% gain with
   16-bit saturation per sample; raw concatenation remains available at the
   default volume for unit inspection. Pitch/speed changes are unsupported.
   `BuildPaul2013UPMSegments` now
   constructs the captured five-word records between adjacent UPM periods,
   including the doubled sample positions and recovered speed-ratio
   arithmetic. `SelectPaul2013UPMSegment` ports the closest-segment scan, and
   `Paul2013ResampledPeriodLength` applies the observed pitch ratio and
   minimum/maximum clamp. `PlanPaul2013UPMSegmentResampling` follows the moving
   segment scan, including repeated selection while the native output budget
   remains, and carries the cumulative output count without reading audio.
   `LimitPaul2013PeriodLength` applies the cumulative
   context-budget cap from `FUN_1002afb0`, and
   `Paul2013UPMSegmentTargetLength` derives the budget from the first segment's
   `FirstPeriod` field. The caller supplies the running prior-output count;
   `ResamplePaul2013PeriodWithinBudget` composes the cap with the isolated
   resampler. `ResolvePaul2013ControlOverrides` ports `FUN_10022850`'s
   independent pitch/speed/volume source switches and clamps the selected
   engine-level values. `NormalizePaul2013Controls` separately resolves
   negative API defaults, maps speed zero to 50, and clamps effective pitch,
   speed, and volume to 50–200, 50–400, and 0–500. The
   raw-unit, isolated-period, and
   segment-plan renderers use that result; when gain is enabled they apply
   volume/100 with the observed 16-bit saturation. The engine-level pause
   control remains unsupported; inline VTML pause tags are handled separately
   through explicit token-boundary frame offsets.
   Linear interpolation remains only in the isolated-period experiment.
   `PeriodResamplingRenderer` splits decoded units at UPM boundaries and
   applies that interpolation with default speed; it is an experimental
   renderer and makes no DLL-parity claim.
   `SplitPaul2013UPMSamples` exposes those exact UPM-to-PCM source spans as
   independent sample windows with offsets, rejecting zero periods and any
   vector that fails to cover the decoded PCM exactly. Separately,
   `SplitPaul2013UnitUPMSides` maps the record's two overlapping period sides
   to those windows and checks both declared side sample spans. The
   `PeriodResamplingRenderer` consumes the side views and emits their shared
   boundary only once. The
   `UPMSegmentPlanRenderer` uses the moving segment plan over each selected
   unit's own UPM spans and accepts pitch and speed controls. At default pitch
   and speed it preserves direct samples. For adjusted controls,
   `BlendPaul2013UPMSegmentWindows` ports the current output-context window
   and tail-aligned selected source-period window in `FUN_1002afb0`: the
   falling and rising halves of the embedded table, aligned crop/zero-pad
   behavior, per-lane truncation and saturation, and the final 16-bit sum.
   The segment renderer then appends the following source period and carries
   it as context for the next step. Its initial context is zero-filled because
   upstream timeline scratch state is not connected. Text-driven segment
   state, shared-buffer scratch tails, row and timeline integration, and full
   synthesis parity remain incomplete.
   `Paul2013ContextMultipliers` ports the observed left/right context
   gates over caller-supplied model rows and neighbor indexes. The gate does
   not gather neighbors. `BuildPaul2013UPMEdgeWeights` constructs the observed
   integer edge ramps from gathered counts after applying the observed
   period-count and five-entry caps, together with gate multipliers.
   `BuildPaul2013ContextEdgePlan` connects the supplied context rows and
   gathered neighbor counts through the gate, caps, and ramp construction;
   `MixPaul2013ContextUPMTimeline` carries the resulting plan through the
   timeline accumulator for supplied period windows. The row-level API below
   reconstructs windows from explicit selected rows; text-driven neighbor
   gathering and production of native shared-buffer scratch tails remain
   unresolved.
   `ApplyPaul2013BlendWindow` ports the aligned source/output placement and
   phase length in `FUN_1002d010`/`FUN_1002cdf0`, cropping or zero-padding one
   supplied window. `FUN_1002d010` indexes the rising half of the 8,192-entry
   float table at `0x1006d1b0`; `FUN_1002cdf0` starts at index 4,096 for the
   falling half. Both Go paths now use the exact 8,192-entry float32 slice
   extracted from the read-only DLL; the analytic curve remains a comparison
   model and differs from some stored coefficients.
   Weighted products clamp to 32,766/-32,767 and use the DLL conversion helper's
   x87 truncate-toward-zero behavior before storing each int16 sample.
   `AddPaul2013WeightedWindows` applies the traced
   integer multiply/divide, per-contribution clipping, and 16-bit accumulation
   to caller-prepared windows; callers still supply their offsets and
   normalization. `MixPaul2013UPMInterval` connects the window preparation,
   edge weights, and accumulation for one interval when callers supply the
   selected source windows. `MixPaul2013UPMTimeline` places prepared periods
   consecutively within one unit and validates them before accumulation.
   `BuildPaul2013TimelineUnitView` ports `FUN_1002c120`'s combined, first-side,
   and second-side DAT/UPM offsets, sample counts, period counts, and edge
   values from a decoded unit and explicit view mode. Its row-mode producer
   remains unresolved. `BuildPaul2013NormalTimelineRow` adds the directly
   observed normal-row kind, primary index, three carried control words, and
   doubled edge spans while representing descriptor pointers as unit
   references. `Paul2013TimelineBankFileIndex` applies the observed fixed
   `gen`/`num`/`etc`/`alp` file ordering, and
   `BuildPaul2013NormalTimelineRowsFromUnitRefs` reads selected units and
   composes normal-row construction using that mapping. The model-state mode,
   primary sequence index, and carried control words remain explicit.
   `BuildPaul2013SyntheticTimelineRow` fills the observed boundary
   row kind, duration, opaque index, sentinels, and marker from explicit lookup
   and state values. `BuildPaul2013TimelineRows` composes the combined view
   for selected units. `Paul2013TimelineOutputFrames`
   applies the measured non-final-row cursor advance, and
   `Paul2013SyntheticTimelineSamples` ports the captured sentence-boundary
   integer expression when its two numeric inputs are available. The exact
   `Paul2013TimelineOutputFramesWithSynthetic` accounts for ordered
   normal and synthetic rows with their distinct length rules; production of
   synthetic rows from source text remains open. The exact 2013 join still
   differs from this cursor prediction by seven frames in one
   stable capture, so this is duration-to-timeline planning, not exact PCM
   assembly.
   The segment plan now runs over audio read from each selected unit, validates
   that UPM spans cover its decoded PCM, and stops at the observed cumulative
   output budget. It still copies the source unit's second-period tail as a
   stand-in for reconstructed context audio. `MixPaul2013Timeline`
   now joins already-prepared rows when adjacent edge spans are equal or
   unequal, using the active-span coefficient schedules and cursor accounting
   in the direct `FUN_1002aac0` path with the exact 8,192-entry float32 table
   extracted from the read-only DLL at RVA `0x6d1b0`.
   `engine/synthesis/stage8_capture_test.go` also verifies a byte-identical
   default-pitch WAVE for the captured two-unit `Hi.` case, where neither row
   returns a context row. This local-model check does not cover the eligible
   neighbor path or text-derived row modes.
   `ApplyPaul2013BlendWindow` and `MixPaul2013UPMInterval` use the same table.
   The strict `MixPaul2013EqualSpanTimeline` wrapper retains its explicit
   equality check. The reproducible extractor maps the PE section and writes
   only the coefficient slice; the vendor DLL remains unchanged. The mixer
   still requires prepared source windows. `MixPaul2013ContextTimelineRows`
   prepares these from explicit rows, while text-driven selection and scratch
   tail production remain unresolved. `UPMTimelineJoinRenderer` and
   `vtconcat -render-mode join` connect the general join path to complete
   selected-unit payloads for diagnostics; native text usually builds selected
   subwindows instead. Text-driven model-window selection, output clipping,
   and prosody remain unimplemented.
   `BuildPaul2013TimelinePCMRow` applies the captured view-mode source offset
   in `FUN_1002c8b0`: modes 0 and 1 start at decoded sample zero, and mode 2
   starts at decoded sample count minus the row's leading span. The normal-row
   renderer now carries mode `+0x26` separately from row kind `+0x27`, and
   checks the copy against the supplied buffer extent, which may include
   scratch samples after the decoded segment.
   `ScalePaul2013TimelinePCMRow` applies
   the captured per-row percentage gain with integer truncation and native
   saturation before joining. `RenderPaul2013TimelineRows` reads each selected
   unit, applies those row operations from explicit metadata, and connects the
   prepared windows to the general mixer. Its strict
   `RenderPaul2013EqualSpanTimelineRows` variant checks equal adjacent spans.
   `BuildPaul2013TimelinePCMInputs` maps typed normal rows directly into the
   renderer's row inputs, and `RenderPaul2013NormalTimelineRows` composes that
   mapping with unit extraction and joining. The row renderer accepts an
   explicit scratch PCM tail for type-2 source windows that extend beyond the
   selected unit's decoded span. `MixPaul2013ContextTimelineRows` now
   connects explicit current/left/right selected rows to the Stage 8 UPM
   source windows: it loads each unit, applies its supplied view mode and
   gain, aligns the capped left prefix and right suffix to the current row,
   checks the row's sample count and doubled edge spans against the decoded
   unit view, and invokes the captured context mixer. It leaves output
   untouched when both native context gates are closed. Context rows, modes,
   selected unit references, and any shared-buffer scratch tail remain caller
   inputs; text-to-row and text-driven neighbor-selection producers remain
   missing.
   `RenderPaul2013ContextTimelineRows` composes that step over explicit
   per-primary-row inputs and passes the resulting rows through the direct
   default-pitch timeline join. It remains an explicit-input renderer, not a
   complete text-to-speech path.
   `BuildPaul2013ContextTimelineRowInputs` maps already-built normal timeline
   rows into that renderer's context inputs, preserving explicit scratch tails;
   `RenderPaul2013NormalContextTimelineRows` composes the adapter and renderer.
   `BuildPaul2013NormalContextTimelineRowInput` ports the observed
   mode-dependent left `base+1` and right `base-1` row adjustments and attaches
   the selected rows and scratch tails. Base indexes, gate rows, and modes
   remain explicit upstream inputs because their producers are not recovered.
   `BuildPaul2013NormalContextTimelineInputs` and
   `RenderPaul2013NormalContextTimelineFromRows` apply that adjustment across
   the supplied gate sequence before rendering.
   Compare PCM/WAVE output with the eligible-neighbor Stage 8 captures before
   claiming synthesis parity. Text-derived gate/index production, shared
   scratch-tail production, and complete synthesis parity remain open.
5. **Additional packages and interfaces.** Add another voice or public API
   mode only after its format and original-engine behavior have separate
   evidence and acceptance criteria.

## Current use

From `engine/`, run `go run ./cmd/vtdecode -index INDEX -dat BANK -unit 0
-output /tmp/unit.wav`, using matching, local 2013 Paul files. This extracts
one unit; it does not turn text into speech. `go test ./...` checks the decoder
against tracked Stage 2 DLL PCM captures. `voice.OpenBank` reads the unit's
index, DAT, and UPM data. The embedded lexicon parser validates every local
record and its observed pronunciation payload grammar.
`text.NormalizePaul2013UserDictionarySource` separately ports the
user-dictionary source-normalizer helper, retaining its signed return code and
partial-output termination state; the public wrapper discards that result.
Its two-byte input domain and one-byte contexts around rejected pairs have
direct runtime coverage, while longer sequences follow the recovered static
loop. `text.NormalizePaul2013UserDictionaryTarget` ports the disassembled
target-normalizer scan and its observed mutation paths, retaining native
numeric statuses. Its one-byte input domain and selected multi-byte patterns
have original-runtime captures; broader sequence parity remains open. Key
encoding, its small Paul 2013 tables, and compact phone-ID expansion are
provisioned.
`text.CheckPaul2013UserDictionaryTargetPhon` ports the writable-buffer
`VT_CheckUserDict_TargetPhon_ENG` validator, including its 69 accepted
uppercase phone spellings, special `#` token, terminal case-folded `[CI]`
truncation, trimming, 260-byte rejection boundary, and native result codes.
`text.ConvertPaul2013UserDictionaryTargetPhon` maps validated phone tokens to
the recovered internal bytes and enforces the 65-token limit. The marker's
broader target-dictionary effect remains unassigned beyond the tested
case-insensitive source matching; see the Stage 21 target-phone captures.
The `api` package also ports Stage 21's six-slot speaker-name and metadata
tables and `VT_GetDBSize_ENG` selector/load/output-write contract, plus the
five-case `VT_GetUserDictLimit_ENG` selector and the complete byte normalization observed
for `VT_SetTextTypeForHighlight_ENG`, and the pre-load gate and byte mapping
observed for `VT_SetUnitSelectHistoryMode_ENG`. It models
`VT_SetPitchSpeedVolumePause_ENG`, its pointer-selective getter,
`VT_SetCommaPause_ENG` and its nullable-output getter, and
`VT_SetEmphasisFactor_ENG` over caller-provided six-slot state, preserving
selector fallback, loaded-slot gates, negative-value no-ops, and native clamps.
It also ports the initialized-object gates and negative-to-zero clamps for
`VT_SetParenthesisCharNumber_ENG` and `VT_SetEnglishReadingRule_KOR`, and
applies the initialized-object gate to the highlight setter.
Database size/load state stays caller-provided. The highlight flag's effect on
auxiliary output and the history setting's synthesis effect remain unresolved;
captures and reproduction steps are in `tools/revkit/work/stage21/README.md`.
`text.LexiconFrontend.ResolveText` resolves ASCII dictionary surfaces into
phone alternatives while preserving the original internal symbol bytes and
CMU labels. It also carries inclusive original-source byte spans through
lexical resolution and pronunciation selection; expanded words retain their
source number token's span. `text.BuildPaul2013OrdinaryParserOffsetSegments`
separately reproduces the captured exclusive parser offsets for ASCII-letter
words and captured numeric forms: signed integers, decimals, `$5.00`, `25%`,
`21st`, `01/02/2024`, `3:45 PM`, and `555-1234`. The captured telephone rows
are grouped cardinal spans separated by a hyphen row. The frontend's spoken
form now matches the captured row text, using “to” for the hyphen. Other row
counts use existing normalizers as inference.
Abbreviation, exception, and TPP parser branches, raw row construction, and
source lookup tables remain unported. It expands observed
plain and signed integers, grouped comma
integers, and decimal forms, including the four-digit year rule. The captured
`3:45 PM` clock example is also expanded, with only the captured
`H:MM` formatting range supported by inference. Numeric ordinal suffixes from
zero through thirty-one are expanded; those spellings are inferred except for
the captured second/fifth word outputs. The captured `01/02/2024` date is
expanded as `January second twenty twenty four`; broader MM/DD/YYYY component
bounds are inferred, with no month-specific day validation. Currency and
percentage forms use separate normalizers; only `$5.00` and `25%` have direct
output captures.
The frontend still cannot produce model contexts from arbitrary text or
produce audio end to end.
`EmbeddedDictionary.ResolvePaul2013SurfacePronunciation` remains
available when callers supply key tables and a codebook directly. The 69
controlled internal phone bytes have CMU labels; structural and unknown values
fail closed. `voice.BuildPaul2013ClassCatalog` derives the exact-key class
membership index from the four Paul unit banks, and `selection.RankClassRecords`
implements the measured class shortlist limits when given a candidate list.
`selection.ExpandClassRecords` flattens retained memberships under an
explicit unit cap, covering only the final class-to-member step of legacy
candidate expansion. `selection.ExactContextUnitCandidates` returns bounded
unit references from an exact context-key match without candidate broadening
or per-unit scoring. Catalog construction matches the native 61,566-key table,
three captured key-to-ID lookups, and their captured unit expansions; broader
candidate-list parity remains open. Transition scoring is now connected by
`selection.SelectPaul2013ModelMinimumCostPath`, and
`selection.SelectPaul2013QueriedCandidatePath` composes the query-to-backtrack
stages when callers provide their state and score inputs. The frontend still
does not produce those inputs, so this does not establish text-driven
candidate-list parity.
`pipeline.Pipeline` joins
injected frontend, selector, and renderer implementations, but the concrete
text-to-audio stages are not implemented.
`go run ./cmd/vtconcat -data-root ../data-paul/M16 -units gen:12,gen:13
-output /tmp/units.wav` writes a diagnostic WAVE from explicit unit IDs; the
command refuses to overwrite an existing output. Its `-render-mode` option
also exposes isolated-period and experimental UPM segment-plan renderers for
known unit sequences. These modes do not add text-based unit selection or
whole-synthesis parity.
`go run ./cmd/vtlex -dictionary-root ../data-common/dict-eng -text "Hello."`
prints the normalized lexical tokens, pronunciation alternatives, CMU labels,
and original internal symbol bytes. Its optional `-choices` list and
`text.SelectLexicalPronunciations` flatten caller-chosen alternatives into
phone order while preserving token spans and separators. The engine also
provides `LexiconFrontend.ResolveUniquePhoneSequence` as a low-level
unique-alternative helper. `duration.Engine.ResolvePronunciationSequence` and
`Engine.Evaluate` produce complete 15-value rows for each path code, evaluate
them with `engbi.tree3`, apply class-membership scoring with stable first-wins
ties, and select the matching lexical alternative. The resolver currently
requires one path group per alternative, as constructed by the recovered
dictionary producer, and rejects other layouts. Its feature rows have not yet
been compared row-for-row with runtime captures.
`text.BuildPaul2013TokenPhoneNeighborhoods`
derives phone neighbors across the full utterance and groups features back by
token. It also derives group count and the observed single/non-final/final
position state per phone from the `FUN_10013f30` groups, including the no-vowel
fallback observed in `FUN_10013c00`. It also ports the duration-row fields and
per-phone auxiliary labels from those groups. The new
`text.BuildPaul2013TokenDurationInputsFromText` joins them with utterance-wide
phone neighbors and ordinary boundary sentinels without caller metadata. This
path assumes the standard terminal `Z` marker. The `FUN_10012c70`
subdivision is now also exposed by `text.BuildPaul2013ModelRecordPhoneGroups`,
which reads each supplied native record's `+0x329` marker bytes and composes
the existing `FUN_10013f30`/`FUN_10013c00` group projection.
`BuildPaul2013ModelRecordPhoneGroupsFromSymbols` derives decoded phone
identity, vowel, and stress values from record bytes `+0x2e8` using the static
Paul CMU codebook. Stage 10 runtime captures cover all 69 observed symbol
codes and all 39 identity ordinals; this validates the symbol-to-feature path
without a direct dump of `DAT_1007b6c0`/`DAT_1007baa8`. The verifier and
captures are `tools/revkit/work/scripts/verify_cmu_codebook.py` and
`tools/revkit/work/stage10/`. Its result now carries the emitted native
`0x1e`-byte group rows, parallel per-phone labels, and the uncounted tail row
written by the zero-vowel fallback. `FUN_100130e0` state-to-marker mapping is
available when callers supply its state inputs; production of those inputs
from source text and model state remains unported.
`text.BuildPaul2013ModelRecordDurationInputs` joins the record-derived group
fields and terminal marker to the existing nine-short duration input writer.
It derives the nine input fields, group count, per-phone position states, and
left/right boundary identities from the emitted group rows and terminal
markers. The enclosing-stream record index and preceding marker remain
explicit because they are stored outside the 0x3c0-byte record.
Its outputs follow counted-group order, retain physical phone indexes, and
omit uncounted zero-vowel fallback rows. Neighbor sentinels follow the first
and last counted-group positions even when fallback phones lie outside them.
`duration.EvaluatePaul2013ModelRecordDuration` now evaluates those vectors
with the loaded duration trees and returns each phone's identity, tree, leaf,
raw value, and input vector. These remain raw tree outputs without established
physical timing units.
`duration.EvaluatePaul2013ModelRecordStreamDuration` consumes packed records
and derives record indexes and preceding terminal markers from stream order.
`duration.EvaluatePaul2013ModelStateArenaDuration` reads the signed native
record count and record base at +2 and +0x64c from the full model-state arena
before evaluating that stream.
`text.BuildPaul2013ModelRecordPitchInputs` reads the same emitted group rows
and produces one 11-short input per counted group, matching the loop bound in
`FUN_100138c0`. `duration.EvaluatePaul2013ModelRecordPitch` evaluates its
marker-selected scalar/vector trees, and the stream and model-state-arena
entry points derive the preceding marker and native record sequence. Their
outputs are raw group vectors before boundary smoothing.
The one-group `P AH0` input in `tools/revkit/work/stage10/feature-p-ah0-trees.log`
is `[1,3,0,3,3,1,0,4,1,1,1]`, consistent with the recovered group-row field
mapping; broader record-row differential coverage remains open.
The native `FUN_100138c0` caller reads its signed outer iteration count at
descriptor-table offset `+0x2`. Each `FUN_100137c0` call then reads the
following 16-byte descriptor's leading count and the current row-array pointer
at table offset `+0x18` plus the descriptor index stride. The existing
`SmoothPaul2013PitchVectors` matches the window arithmetic, while
`SmoothPaul2013PitchDescriptorTableRuns` applies the table-level bound and
following-descriptor count to caller-projected vector rows. Native
`ReadPaul2013PitchSmoothingTable` reads the signed table count at `+0x2`, the
16-byte-stride descriptor counts, and current row pointers at `+0x18` plus
that stride; it resolves those pointers through an explicit caller callback.
Pointer-to-model-record mapping remains unresolved; lexical-token pairing does
not prove those scopes are equivalent.
`duration.EvaluatePaul2013Text` selects one of the nine duration trees using the
DLL's phone-class selector tables and evaluates each vector; outputs are not
yet converted to timing. Unit selection remains unwired from text.
The package currently does not claim parity for mode 8, other voice
generations, or whole-synthesis output.
