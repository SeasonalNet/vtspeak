# Go engine

This module is the independent Go implementation. It reads 2013 M16 Paul
resources and contains evidence-backed text, pronunciation, duration/pitch,
candidate-selection, and audio-reconstruction mechanics. The source-text to
model-context path and the selection/rendering handoff are still incomplete,
so it does not yet synthesize arbitrary text end to end.

## Packages

- `api` exposes measured Paul 2013 export contracts for speaker names and
  metadata, database-size queries, per-speaker pitch/speed/volume/pause/
  emphasis settings, the five-case user-dictionary limit switch, the full
  byte-to-boolean highlight setting, process-global parenthesis and reading
  rule setters, and the unit-history setter's pre-load gate. Database load
  state and size values remain caller inputs; unresolved synthesis effects of
  the highlight and history settings are not claimed.
- `dat` validates the observed 2013 Paul unit-index layout, reads DAT and UPM
  spans, decodes a unit to signed 16-bit PCM, and writes a mono 16 kHz WAVE.
  `WriteWAV` streams decoded bytes, and `WriteWAVFromSamples` converts samples
  through a fixed-size buffer when callers do not need a complete in-memory
  WAVE value.
- `voice` opens one matching index/DAT/UPM bank read-only. `Bank.ReadRecord`
  returns indexed metadata and feature bytes without reading waveform data;
  `Bank.ReadUnit` also reads the UPM and DAT spans, decodes PCM, and verifies
  the cached first/shared/last UPM bytes plus the decoded sample count against
  the combined UPM vector.
  `OpenPaul2013` groups the four banks and decision trees in the local Paul
  M16 data root. `Paul2013.GlobalUnitOrdinal` maps a bank-local unit reference
  through the ordered `dblist.idx` bank ranges for cross-bank continuity tests.
- `tree3` parses and evaluates the supported Paul decision trees and loads all
  nine duration and eight pitch trees after checking their observed shapes.
  It also loads the shared 531-node, one-byte `engbi.tree3` pronunciation
  classifier from the dictionary root. `ParseAt` reads a tree from a
  concatenated buffer; `LoadPaul2013ATMTTrees` validates and loads the shared
  27-tree `atmt.tree3` container without modifying its vendor input.
- `duration.OpenPaul2013` loads the shared English dictionary and Paul
  duration trees plus the 27 shared ATMT trees from their local roots. It
  also retains all nine shared `.txt2` tables; its
  `ApplyPaul2013TPPNumericTextRows` method runs the direct F/G and WAB pass on
  caller-supplied normalized model rows.
  `Engine.EvaluateATMTTree` evaluates a caller-selected tree with explicit
  numeric inputs; it does not infer a tree role or feature meaning.
  `Engine.EvaluatePaul2013NormalizerCharacter` builds the recovered ten-short
  input for one supplied character/category position and evaluates its
  letter or apostrophe tree. `EvaluatePaul2013NormalizerCharacterSequence`
  scans supported ASCII letter/apostrophe tokens from right to left, feeds
  each signed result byte into the next tree input, applies the observed
  short-token internal-E override, then returns results in source order. The
  `text.ShapePaul2013NormalizerClassCodes` ports the following opaque-byte
  escape and punctuation shaping pass from `FUN_10002200`, while
  `duration.Engine.ClassifyPaul2013NormalizerToken` composes the character
  trace with `text.RemapPaul2013NormalizerClassCodes`, which applies the observed
  `FUN_100024f0` explicit remaps and table at DLL RVA `0x77f5a`. Its sentinel
  entries are covered by the native explicit remap cases. The CHC gate and
  supported inner `FUN_10002f10` dispatch now connect these stages to direct
  context-code fallback or opaque class-code output with its observed `0x10`
  row-flag mask. `ApplyPaul2013GenericNormalizerWrites` writes the selected
  bytes and NUL terminator into a copied destination and ORs the result mask
  into the separate caller flag byte. `NormalizePaul2013ContextToken` ports
  the outer `FUN_1000cb30` branch order around an explicit context-lookup
  callback. `NormalizePaul2013ContextTokenToBuffer` composes that dispatch
  with copied C-string writes, the initial no-match NUL, and the generic
  normalizer's local flag byte. Its embedded-dictionary variant uses the
  loaded lexicon for an already-normalized token.
  `text.ApplyPaul2013ModelSourceClassNormalization` ports
  `FUN_10009030`'s ordered S/C source-class branches against the full model
  arena, and `duration.Engine.NormalizePaul2013ModelContextRow` supplies its
  generic normalizer stage. `NormalizePaul2013ModelContextRows` composes those
  calls in an explicit caller-supplied order and returns the rows handled.
  Per-row preparation return and exception eligibility gates now also read
  directly from the full model arena. `RunPaul2013KnownModelContextPasses`
  now composes the supported parts of the counted-row iteration, including
  the exception-match index adjustment. Other first-pass row mutations, the
  generic dictionary/TPP and later special-handler chain, and runtime row
  comparison remain incomplete; see the Stage 63 and Stage 64 reports.
  `text.DerivePaul2013ModelContextPreparation` and
  `EvaluatePaul2013ModelContextExceptionGate` bridge the known low-short
  preparation return and caller gate to the full model arena.
  `text.ApplyPaul2013ModelContextFormY` also applies the form-Y parser-string
  copy, state bit, and adjacent duplicate-row advance. Its `WithTail` variant
  applies directly recovered literal suffix rows from an explicit auxiliary
  surface. Generic auxiliary-string dispatch remains unresolved; see the
  Stage 65 static report.
  `text.ApplyPaul2013ModelContextFormS` applies the form-S context-string
  encoding and row flag, and has the same explicit-tail literal variant.
  `text.ApplyPaul2013ModelContextSecondPass` ports the final counted-row pass
  for the literal previous-surface `the` and following-surface `de` cases,
  including its class gate and preceding-row phone-string/state writes.
  `text.InitializePaul2013ModelContextProcessingFlags` ports the initial
  counted-row state initialization in `FUN_10007520`; the combined
  build-and-normalize entry points now apply it before source-class
  normalization.
  `duration.Engine.RunPaul2013KnownModelContextPasses` now composes the
  supported Y/S mutations, exception lookup and native matched-row advance,
  then the native dispatch gate and `FUN_10009030` source-class normalizer
  before `FUN_100086c0` and the supported fallback chain, plus the final
  `the`/`de` pass in model-row order. A closed exception gate with an open
  outer gate reaches known post-handler rules; a closed outer gate skips both
  normalization and post-handlers and is listed in `OuterGateSkippedRows`.
  This distinction follows the caller pseudocode in Stage 64. It reports rows
  still requiring unsupported suffix branches and TPP row updates instead of
  treating those rows as complete.
  `BuildAndRunPaul2013KnownModelContextPassesFromParserRows` and
  `BuildAndRunPaul2013KnownModelContextPassesFromOrdinarySource` connect this
  pass to the existing token/context projection.
  The counted-row pass attempts supported A140/spelling, c040, B800, C3A0,
  trailing-apostrophe, and terminal generic stages in native order after each
  preceding handler misses, writing matched outputs into the copied model
  arena. Historical `With...Rows` entry points remain for compatibility and
  conflict checks or controlled mode overrides.
  `RunPaul2013KnownModelContextPassesWithCompoundContractionB800AndC3A0Rows`
  retains its historical row arguments for conflict checks.
  `BuildAndRunPaul2013KnownModelContextPassesWithCompoundContractionAndC3A0RowsFromOrdinarySegmentsInSharedArena`
  connects the same pass to shared-arena ordinary-source projection.
  `Paul2013CompoundOuterMode` derives the `FUN_10009dc0` mode from the
  preceding context row's parser-row index and parser type field. The
  counted-row pass now attempts `FUN_10009dc0` for each row that reaches its
  native branch and derives the mode automatically. Explicit mode maps remain
  available as overrides for controlled comparisons.
  The counted-row pass now runs under the `FUN_10007520` dispatch gate after
  exception misses: it applies `FUN_10009030`, then the parser-class-A
  `FUN_10002f10` precheck when gated. Both direct and hyphen-split
  `FUN_10003110` routes are composed before the supported `FUN_100086c0` path
  and the automatic fallback chain. That path reads source from
  `model + 0x429ad + row*0x70`; the counted-row surface begins at row offset
  `+5`. A nonzero result applies the
  recovered `FUN_10009cd0` encoding and bit-3 write, then runs the known
  post-handler rules. The nonzero phone-marker gate then blocks lower-priority
  fallback handlers. The hyphen handler preserves its 64-byte output cap and
  accumulated flags. Rows remain in `UnresolvedRows` because the remaining
  TPP updates and runtime parity are not established; unsupported gate cases
  return errors.
  Its `...B800AndC3A0RowsFromOrdinarySegmentsInSharedArena` variant also
  connects the pass to ordinary-segment projection. Unsupported suffix
  branches, remaining TPP handlers, and native row parity remain unresolved.
  The `FUN_10007520` caller listing establishes the stage order; the
  trailing-apostrophe route is documented in Stage 67 and the component and
  terminal routes in Stages 82 through 87. See also Stages 66, 81, and 111.
  `text.Paul2013ContextWordPairSpecial`, `text.Paul2013IsNumberWord`, and
  `duration.EvaluatePaul2013FUN10008580` port three counted-row helper
  predicates and their mapped-string/table conditions. Their outputs are
  available for caller composition; outer selection and runtime row parity
  remain separate. See Stage 112.
  `duration.ApplyPaul2013FUN10007520UseMarker` ports the subsequent `use`
  row rewrite gated by the third phone-code byte and `FUN_10008580`, including
  its row-state bit. The known counted-row pass applies it after supported
  per-row results; other preceding post-handler rules remain open. See Stage
  113.
  `duration.ApplyPaul2013FUN10007520UpMarker` ports the adjacent-row `UP`
  case and its `08 35 00` output. The counted-row pass composes it after the
  `use` rule; the remaining native post-handler rules are incomplete. See
  Stage 114.
  `duration.ApplyPaul2013FUN10007520Minute` and
  `duration.ApplyPaul2013FUN10007520CloseMarker` also port the direct `minute`
  number-word rewrite and gated `close` phone-code rewrite. The counted-row
  pass runs them for supported per-row results; other table-backed post-handler
  rules and full runtime parity remain open. See Stage 115.
  `duration.FindPaul2013ContextPhraseWindow` ports the bounded L/R/B neighbor
  selection and sorted phrase-table search of `FUN_1000cc60`. It accepts the
  table snapshot and native comparator as inputs; all four `FUN_10007520`
  call sites are connected through the bow, lead, mouth, and read rules.
  Runtime parity remains open. See Stage 116.
  `duration.ApplyPaul2013FUN10007520MouthMarker` also connects the recovered
  class-12 table to the direct two-row `mouth` rewrite. Other post-handler
  branches and runtime parity remain open. See Stage 117.
  `duration.ApplyPaul2013FUN10007520Bow` and
  `duration.ApplyPaul2013FUN10007520Lead` connect the recovered six-word and
  ten-word tables to their bounded bidirectional phrase lookups and phone-code
  writes. The `read` temporal table is covered in Stage 119. See Stage 118.
  `duration.ApplyPaul2013FUN10007520ReadMarker` ports the `read` apostrophe
  gate, its recovered table and temporal phrase predicates, and the `0x18`
  rewrite. Runtime row parity and decompiler accuracy remain unverified. See
  Stage 119.
  Form-Y/S model-context handlers also derive `FUN_100091b0`'s suffix input
  from the selected parser row's `+0x14` offset and `+0x18` length when no
  explicit override is supplied. The ordinary parser-row writer still does
  not produce that offset. See Stage 120.
  Their result now exposes the tail to the duration pass, which follows a
  direct-literal miss with embedded `FUN_1000cb30`, direct `FUN_10002f10`, and
  `FUN_10009cd0` in native order. It appends the `d` marker and preserves the
  nested flags before running post-handler rules. Form-Y post-handlers use the
  duplicate-advanced context row. This uses the existing ASCII normalizer and
  dictionary support; unresolved TPP/resource traversal, the parser offset
  producer, non-ASCII behavior, and runtime row parity remain open. See Stage
  121.
  `ApplyPaul2013FUN10007520ArticleMarker` adds the directly decidable `a`
  following-class and same-parser-row apostrophe-s rewrites to the post-handler
  pass, including the row-zero gate using recovered class-3 and class-4
  tables. Runtime row parity remains open. See Stage 122.
  `duration.Engine.ResolvePaul2013PronunciationExceptionFromModelRows` adds
  caller-ordered exception candidate extraction, X-prefix retry, and decoded
  phone output; a match also applies the recovered `FUN_1000ca50` delimiter
  splitting and phone-row/state writes to a copied model arena. Candidate rows
  must be contiguous from the gate row. Its
  `ResolvePaul2013PronunciationExceptionFromModelContextRows` entry point now
  derives that prefix directly from the model: it stops on an unnormalizable
  next surface or when the accumulated one-to-four component total is full.
  The preceding normalization mutations and native runtime row comparison
  remain open.
  `LookupPaul2013EmbeddedContextToken` supplies the embedded-key transform,
  payload parse, pronunciation count, first internal phone-symbol string, and
  payload bit-7 fallback gate;
  `NormalizePaul2013ContextTokenFromEmbeddedDictionary` composes that path.
  Its normalization result now preserves `FUN_1000cb30`'s native short return
  (`1` handled, `-1` miss) for callers that continue the ordered fallback.
  `TransformPaul2013ModelSourceFromEmbeddedDictionary` applies the decoded
  E/A result type to the source scanner's terminal-punctuation split.
  `NormalizePaul2013TrailingApostrophe` ports `FUN_1000c710`'s mapped `in'`
  retry with case-dependent `G/g`, apostrophe-stripped lookup, generic
  normalization, and context-string fallback with the observed row-flag
  updates. Its context lookup remains a resource input; an embedded-dictionary
  wrapper is available. The TPP-family lookup and other caller-side token
  transforms remain open, so this is not a complete text normalizer. The
  recovered `MC` branch classifies `MAC` plus the uppercased suffix when the
  suffix contains a table-recognized vowel; fallback context codes still
  derive from the original token.
  `NormalizePaul2013MappedSpellingSuffix` ports all nine recovered literal
  `FUN_1000b4d0` suffix rewrites and their `FUN_1000cb30` lookup/flag behavior,
  including minimum stem lengths. It must follow the earlier `FUN_10009dc0`
  and `FUN_1000a140` handlers in a caller's cascade. The static evidence is in
  `tools/revkit/work/reports/stage69-model-context-suffix-rewrites.txt`; no
  runtime row comparison has been made.
  `NormalizePaul2013AnceSuffixContext`, `NormalizePaul2013NessSuffixContext`,
  `NormalizePaul2013MentSuffixContext`, `NormalizePaul2013ShipSuffixContext`,
  and `NormalizePaul2013LessSuffixContext` port directly compared 7- and
  8-byte `ance`, `ness`, `ment`, `ship`, and `less` branches in
  `FUN_1000a140`, including the `FUN_1000cb30` stem lookup, branch-specific
  marker bytes, and row flag bit.
  `NormalizePaul2013EstSuffixContext` ports the 6- to 8-byte `est` branch,
  including its four-byte truncation and `0x07 0x2d 0x37` output tail; the
  adjacent `ers` handler now ports both ordered truncations and its table-gated
  `i`-to-`y` rewrite, with the native `D`, `0x1a D`, and `0x07 79` tails
  (Stage 99). `NormalizePaul2013IngSuffixContext` ports the 5- to 31-byte
  branch, including the short-input guard around the `length-7` and
  `length-6` table reads; the old long-input method delegates to it (Stage
  101). `NormalizePaul2013LySuffixContext` ports the recovered 5- to 31-byte
  branch while preserving the old long-input entry point (Stage 100).
  `NormalizePaul2013IstSuffixContext` applies the direct three-byte check,
  suffix-stripped lookup, and `#79` output for the same 5- to 31-byte range;
  the old long-input method remains as a wrapper (Stage 102).
  `NormalizePaul2013FulSuffixContext` ports the direct `ful` stem lookup and
  space/control/plus tail between `ist` and `ly` (Stage 103).
  `NormalizePaul2013EdSuffixContext` also ports the bounded long-input `ed`
  branch with its marker-aware lookup and raw DLL table checks.
  `NormalizePaul2013ErSuffixContext` ports the bounded long-input `er` branch,
  including its ordered truncation lookups and shared raw table predicates.
  The direct suffix inventory is summarized in Stage 107; branches outside
  that inventory and runtime row parity remain open. See the Stage 70 through
  Stage 76 reports for the early suffix findings.
  `NormalizePaul2013SupportedSuffixFallbackAfter09DC0` composes these
  recovered bounded branches before the later mapped spelling rewrites when
  the preceding `FUN_10009dc0` handler has already missed. Its scope and
  evidence boundary are recorded in the Stage 77 report.
  `duration.NormalizePaul2013CompoundContext` ports the scanner and output
  assembly from `FUN_10009dc0`, preserving accumulated output when the native
  63-byte gate stops later copies. Callers can provide the component chain
  directly or use the supported composition below; see Stage 78 for evidence.
  `NormalizePaul2013CompoundContextWithEmbeddedDictionary` adds the native
  component gates and `FUN_1000cb30` branch, then requires an explicit callback
  for the remaining ordered fallback chain. Its
  `NormalizePaul2013CompoundContextWithSupportedHandlers` variant composes the
  recovered A140, trailing-apostrophe, contraction-suffix, C3A0, and generic
  handlers. Its c040 prefix uses the supported dictionary/compound/A140/generic
  composition by default; callers can still override that prefix chain.
  `Engine.LookupPaul2013CompoundDictionaryGate` ports the embedded-payload
  flag check used before that chain, and `Paul2013CompoundCharacterGate`
  ports its adjacent ASCII character predicate. `NormalizePaul2013CompoundModelContextRow`
  applies a matched result to a counted row; the compound-and-C3A0 ordered
  pass and shared-arena segment wrapper accept explicit outer branch choices.
  Marker writes before an empty output buffer, outer branch selection, and
  runtime parity remain open; see Stages 79, 80, and 110.
  `text.Paul2013FUN10010010` and `text.Paul2013FUN100100D0` port two
  self-contained callees used by `FUN_100086c0`: the signed-byte attribute
  predicate and mapped-vowel counter. Stage 104 recovers their signed-index
  table prefixes: high-bit inputs have zero attributes and do not count as
  vowels. This establishes these helpers' byte behavior without extending the
  separate source-gate and normalization boundaries.
  `duration.Engine.EvaluatePaul2013FUN100086C0Precheck` carries their results
  through the parent's early status, phone, state, vowel, and generic-
  eligibility gates. It identifies the native neighbor-scan boundary
  explicitly. `ScanPaul2013FUN100086C0Neighbors` ports the nearest
  previous/following classed-surface scans and their `FUN_1000ffd0` gates.
  `CompletePaul2013FUN100086C0NeighborDecision` carries the status-specific
  post-scan state and length branches, including the mapped `Lexus` plus exact
  `IS` case. The six-string `DAT_100783d0` table (`ALL`, `BY`, `DEAR`, `DO`,
  `IS`, `NOT`) is recovered directly from the DLL and searched in native
  `strcmp` order.
  `Engine.EvaluatePaul2013FUN100086C0SupportedPath` also ports the outer
  `param_5` gate, including its one-byte/status exception, mapped `Wi` check,
  and compound character gate. It composes the recovered stages and returns
  the neighbor trace alongside the native low-short result.
  `Paul2013NativeCStringTable` ports the mode-`0x53` `FUN_100560a0` search over
  sorted strings and validates caller-supplied snapshots before search.
  `NormalizePaul2013C3A0ComponentSequenceWithFUN100086C0` connects the
  recovered helper to the supported A140/component cascade; the caller still
  provides the model-row index. The embedded-dictionary entry point now also
  projects `FUN_10003c50` fields into the raw `local_25c` layout, so callers
  no longer need to build that record themselves. `NormalizePaul2013C3A0ModelContextRow`
  applies a matched sequence to a counted model row, and
  `NormalizePaul2013C3A0ModelContextRows` carries ordered writes across an
  explicit row list. The composed counted-row pass attempts C3A0 after B800
  misses; the legacy `WithC3A0Rows` input is retained for compatibility.
  `BuildAndRunPaul2013KnownModelContextPassesWithC3A0RowsFromOrdinarySegmentsInSharedArena`
  connects that pass to shared-arena sentence segments.
  `RunPaul2013KnownModelContextPassesWithCompoundAndC3A0Rows` and its
  shared-arena wrapper retain compatibility inputs while the counted-row pass
  attempts compound and C3A0 in native order. Unresolved neighbor cases fail
  closed.
  `NormalizePaul2013ContractionSuffixContext` ports the direct `FUN_1000c040`
  `'n`, `'ee`, `'er`, and `'ers` marker tails, including uppercase endings
  and the `'n` context-class condition. Its
  `NormalizePaul2013ContractionSuffixWithSupportedPrefix` entry point
  composes the parsed-record, supported `FUN_100086c0`, compound, A140,
  generic, and context-encoder prefix branches. It requires the model arena
  and row index for the row-dependent decision; unresolved branches fail
  closed. `NormalizePaul2013ContractionModelContextRow` applies a match to a
  caller-selected counted row. Outer row dispatch and runtime comparison
  remain separate. See the Stage 81 report.
  `Engine.ApplyPaul2013B800JoinedDictionaryModelContextRow` ports the supported
  `FUN_1000b800` control flow: joined-surface dictionary lookup and group
  writes, then the literal and pointer-table suffix rules on a dictionary miss.
  It reports the matched suffix rule and returns a native-style miss when no
  branch applies. Static evidence does not establish runtime row parity. The
  counted-row pass attempts it after contraction misses. See the Stage 111
  report and pseudocode.
  `text.ScanPaul2013C3A0Component` ports `FUN_1000c860`'s bounded ASCII
  component scan, including its class transitions and apostrophe-`s` handling.
  The `FUN_1000c3a0` proper-name/TPP dispatch that calls it remains incomplete;
  see the Stage 82 report.
  `NormalizePaul2013ApostropheSSuffixFromEmbeddedDictionary` ports the later
  case-sensitive `'s` stem lookup and marker tail, using the loaded dictionary
  and generic normalizer. It remains one branch of the incomplete
  `FUN_1000c3a0`; see Stage 83.
  `NormalizePaul2013C3A0ComponentSequence` composes the repeated component
  scan and output assembly through the native handler callback, including the
  single-component bypass, `d` joins, capacity arithmetic, and final bit-3
  gate. The handler chain remains incomplete; see Stage 84.
  `text.Paul2013C3A0Gate` now ports `FUN_10008cc0`'s mapped prefix exclusions,
  case-class transition, and apostrophe checks; the source-driven sequence
  entry point applies it before scanning. Stage 105 recovers the gate's
  signed-index prefix: high-bit bytes produce an ineligible result. Later
  scanner and component processing retain separate boundaries; see Stages 85
  and 105.
  `NormalizePaul2013C3A0ComponentSequenceWithEmbeddedDictionary` runs the
  direct parsed-record branch: components at least two bytes long with a
  positive pronunciation count emit the first internal-symbol string and set
  flag bit 0. The ordered variant also runs the payload-bit-7 generic
  precheck; a successful precheck skips the earlier exception/
  `FUN_1000a140` callback instead of selecting the parsed string. It then
  applies the embedded apostrophe-`s` branch and generic/context fallback. See
  Stages 86 and 87; row-level runtime parity remains open.
  `NormalizePaul2013C3A0SupportedA140Component` adapts the directly ported
  bounded A140 suffix subset for a caller's earlier-handler chain, after its
  exception check misses. Its supported set includes the long `liness` and
  `ly` paths with their table-checked `i`-to-`y` retries. Stage 92 adds the
  long `ist` stem lookup and `#79` tail; Stage 103 adds the following `ful`
  lookup before `ly`. Stage 96 adds long `ing` stem selection and its
  `#.` tail; Stage 97 adds the 5- to 31-byte terminal-`s` ordered stems,
  protected output gate, marker tail, and bit-3 update. Stage 106 ports the
  full signed short-table memory window and removes the former lowercase-only
  index boundary. Stages 100 through 102 extend `ly`, `ing`, and `ist`
  handling through the short input range. Native exception row/resource
  inputs, A140 branches outside this direct suffix inventory, and runtime row
  parity remain open. Short `est`
  and `ers` handlers are recorded in Stages 98 and 99; Stages 100 through 103
  cover short `ly`, `ing`, `ist`, and `ful`. See Stages 88 through 92 and
  Stages 96 through 107. The branch inventory and its static-only boundary are
  recorded in `tools/revkit/work/reports/stage107-a140-static-branch-inventory.txt`.
  `NormalizePaul2013C3A0ComponentSequenceWithSupportedA140` composes these
  supported suffixes after a required caller handler for the preceding
  `FUN_100086c0`/exception path, then runs the existing apostrophe-`s` and
  generic/context stages. It leaves the earlier handler's resources and state
  explicit.
  `Engine.Evaluate` accepts source text,
  evaluates candidate pronunciation rows with the loaded `engbi` classifier,
  ranks path groups, derives the ordinary-path vectors, dispatches each phone
  through the recovered phone-class selector, and returns raw duration-tree
  values. It also derives lexical convenience pitch inputs per phone,
  evaluates the paired scalar and 12-value pitch trees, and carries their
  outputs alongside each duration result. It applies the six rounded
  seven-value windows from `FUN_100137c0`/`FUN_10013790` to adjacent vectors
  within each lexical token, while retaining the raw tree rows. The native
  smoothing routine walks a separate 16-byte descriptor table; its mapping
  to model-record groups remains unresolved. `SmoothPaul2013PitchDescriptorTableRuns`
  uses the signed table count at `+0x2` as the outer iteration bound, then
  reads each call's row count from the following descriptor's leading short.
  `ReadPaul2013PitchSmoothingTable` extracts those counts and current row
  pointers from a table snapshot and resolves vectors through an explicit
  address-space callback. Pointer-to-model-record mapping remains open. The
  `EvaluatePaul2013PitchTextWithMarkers` entry point accepts one explicit
  terminal marker per token and selects `bt/bf`, `sbt/sbf`, or `qbt/qbf` for
  `^`, `Z`, or `[`; other markers use `nbt/nbf`. Marker production from source
  text remains open. `text.ParsePaul2013CMUPronunciation` and
  `BuildPaul2013CMUPhoneSequence` accept an explicit `x-cmu` phone string;
  `text.ParsePaul2013VTMLCMUPhoneme` also reads the directly captured single
  forced-phoneme element and retains its source span;
  `Engine.EvaluateVTMLCMUPhoneme` connects it to both loaded tree families
  when the caller supplies markers and position states.
  `Engine.EvaluateSequenceWithPositionStates` evaluates its duration and pitch
  trees when the caller supplies the observed per-phone state vector. This
  does not produce those states or selected model units. The frontend
  parses pronunciation paths into the legacy groups, and
  `text.RankPaul2013PronunciationPathGroups` ports class-membership
  scoring and stable first-wins tie selection.
  `text.SelectPaul2013PronunciationByPathMarker` ports the marker lookup and
  fallback in `FUN_10003f10` and returns the phone string aligned with the
  selected path row. `Engine.SelectPaul2013ContextPronunciationFromEmbeddedDictionary`
  composes that selector with the loaded embedded dictionary and preserves
  `FUN_1000cbe0`'s empty-output and return-code behavior. The marker producer
  remains unresolved.
  `Engine.ResolvePronunciationSequence` exposes the selected phone sequence.
  `LexiconFrontend.AnalyzeText` returns token data with partial rows from
  `text.BuildPaul2013KnownPronunciationFeatures` and per-path-code candidate
  rows from `text.BuildPaul2013PronunciationPathCodeFeatures`. The token rows mark
  the four surrounding-word categories (positions 0–3) and first-byte class
  flag (position 14); candidate rows also mark the path-code class (position
  13). Triplet positions 4–7 are now fully produced from the sorted-word
  tables, literal context rules, digit scan, punctuation scan, and default
  class in `FUN_10006ce0`. Positions 8–12 now apply the case-insensitive
  word-shape classifier recovered from `FUN_10007160` to the two preceding
  words, center word, and two following words; absent words map to zero.
  Its complete ordered return cascade is ported: class-2 suffix and paired
  character-table gates, class-3 suffixes, then the `ly`, `ing`, `ist`/`ists`,
  `ed`, and final-`s` classes. Comparisons use the embedded byte map and
  preserve signed-byte indexing in `FUN_1000ffa0` and `FUN_1000ffd0`.
  Together with the path-code class at position 13, these produce complete
  15-value rows for each code in a nonempty path group. The engine evaluates
  rows in dictionary order, accumulates outputs by group, and selects the
  corresponding pronunciation alternative. The resolver requires the
  recovered one-path-group-per-alternative layout and rejects other records.
  A full local dictionary audit found one nonempty path group per alternative
  in all 829 alternatives across 414 ambiguous records. Feature production is
  statically derived and has not been compared row by row with runtime
  captures.
- `distance` validates and reads the little-endian triangular `cepdist.tbl`
  resource, with checked symmetric lookup by metric code and the observed
  packed-code mask. It builds the observed 256-bin triangular feature-distance
  table using the DLL's symmetric difference formula. `voice.OpenPaul2013`
  loads the raw 1,024-entry table and builds the feature table with the banks
  and trees.
- `text.LoadTXT2Tables` reads the nine shared classification and replacement
  tables with validated framing. `text.LoadEmbeddedDictionary` validates all
  indexed records in the shared embedded English lexicon, and
  `text.LoadExceptionDictionary` parses all four sorted `exceptdict` groups;
  `text.NormalizePaul2013UserDictionarySource` ports the separate
  `VT_CheckUserDict_SourceNorm_ENG` helper, including its signed result,
  copied prefix, and unterminated pair-error state. The public export discards
  this scratch result; Stage 21 covers all two-byte inputs and the one-byte
  contexts around rejected pairs. See Stage 93.
  `text.NormalizePaul2013UserDictionaryTarget` ports the mutable
  `VT_CheckUserDict_TargetNorm_ENG` scan, including its exact `[SKIP]` and
  `[CI]` cases, internal hyphen rule, segment counters, and angle-bracket
  state. Numeric statuses stay opaque; broad sequence parity remains open.
  See Stage 94.
  `text.CheckPaul2013UserDictionaryTargetPhon` ports the writable-buffer
  `VT_CheckUserDict_TargetPhon_ENG` validator, including trimming, the 69
  uppercase phone spellings, special `#` token, terminal case-folded `[CI]`
  truncation, and native status values. Its adjacent converter maps those
  spellings to the recovered internal bytes and enforces the 65-phone limit.
  Marker meaning outside the tested dictionary cases is not inferred. See the
  Stage 21 target-phone captures and the lead-6 file/API behavior report.
  `duration.Engine.LookupPronunciationException` exposes their exact
  one-to-four-component category lookup for caller-normalized and
  caller-encoded keys. `text.NormalizePaul2013ExceptionSurface` and
  `LookupPronunciationExceptionSurface` also port the mapped-character,
  apostrophe/hyphen, and component-count behavior with explicit byte
  attributes. `Paul2013ExceptionCharacterAttributes` provisions the relevant
  ASCII mask bits recovered at DLL VA `0x1007e188`; the source-token exception
  dispatch cascade is not yet connected to text resolution.
  `text.BuildPaul2013NormalizerFeatureWindow` ports the ten-short classifier
  input layout from `FUN_10002db0`, including its seven-character window and
  signed trailing categories. `duration.Engine.ClassifyPaul2013NormalizerToken`
  produces reverse-scan character categories, applies `FUN_100024f0` remapping
  and `FUN_10002200` shaping, and retains the per-character trace.
  `LookupPaul2013CHCFlags` connects loaded `chc_sort.txt2` rows to
  `FUN_10002680`'s exact-key four-bit mask lookup. The `FUN_100026f0`
  substring matcher is also ported for exact matches, compound splits,
  trailing-S retry, and the observed long-token fallback.
  `duration.Engine.IsPaul2013NormalizerEligible` composes those lookups with
  `FUN_10002c70`'s ASCII letter/apostrophe scan, vowel boundaries, and
  run-specific masks. It rejects tokens longer than the native 31-byte local
  run buffer; its scan cases are synthetic and lack direct native comparison.
  Exact dictionary/TPP dispatch remains incomplete.
  `LookupPaul2013SurfaceSequence` joins caller-selected adjacent surfaces,
  accumulates category counts, and retains the longest exact prefix. The
  model-row exception resolver derives the bounded contiguous prefix and reads
  its `X` retry markers from the model context rows; source-row production and
  the preceding cascade still determine whether the native exception gate is
  entered.
  `text.LoadTPPDictionary`
  parses all sequential typed-text records and their observed one/two-atom
  payload grammar; `duration.Engine.LookupTypedText` applies the known key
  transform and returns opaque tag/suffix atoms. `TPPAtom.ComponentPattern`
  decodes A-E component counts and raw binary values without naming their
  meanings; `duration.Engine.LookupTypedComponentPattern` connects that
  decoder to loaded dictionary keys. `AX` remains a marker with unresolved
  purpose. `text.ParsePaul2013NativeInteger` ports `FUN_100645ba`'s
  attribute-driven leading skip, optional sign, digit accumulation, and
  32-bit wraparound. `TPPAtom.NumericByte` validates the observed F/G decimal
  suffix grammar, uses that parser, and narrows to one byte while retaining
  those tags' meanings as unknown.
  `Paul2013ProperNameTPPSelector` maps one through four components to the
  observed A-D selectors from `FUN_10034180`; the native count gate skips five
  or more components. Its duration-engine wrapper connects that choice to the
  exact TPP lookup. `text.AcceptPaul2013ProperNameTPPComponentPattern` ports
  the observed post-lookup gate: single-component hits check only the marker;
  multi-component hits reject zero bits paired with marker `d` or `A`.
  `duration.Engine.MatchProperNameTPPComponentPattern` composes lookup and
  gate when callers supply the per-component markers.
  `text.Paul2013ProperNameTPPComponentsFromArena` projects the native counted
  `0x140` component rows, including the final row's raw `+0x18` scanner-input
  dword. A zero input bypasses context checks; for nonzero input,
  `text.EvaluatePaul2013ProperNameTPPContextGate` evaluates the observed
  post-scanner decision tree from caller-supplied scanner and helper outcomes.
  `text.Paul2013NativeStringSearchTable` ports the mode-selected
  `FUN_100035d0` search over caller-supplied sorted snapshots, and the context
  table wrapper can derive its two lookup outcomes from scanner output.
  `text.Paul2013ProperNameContextLateLookup` also ports the bounded
  `FUN_10041df0` Roman-numeral and two-byte cases consumed by that gate.
  The table-backed context wrapper derives `FUN_10010120`'s scanner-word
  predicate directly from the scanner output as well.
  `text.ObservePaul2013ProperNameScannerInput` derives the leading whitespace
  count and status 8/9 early returns. `text.ObservePaul2013ProperNameASCIIToken`
  also ports bounded mode-0x15 ASCII digit runs and letter words with internal
  dot or hyphen separators, plus two exact one-digit mixed forms. It derives
  native `FUN_10062f50` weights from shared letter-expansion records and
  handles the mapped `United States` branch. On the letter-token path, it
  stops before ASCII punctuation and whitespace delimiters, returns the
  consumed-byte cursor, and treats an unjoined dot or hyphen as a boundary.
  A digit run also stops before one ASCII delimiter immediately followed by
  NUL; longer numeric lookahead remains unsupported. Its duration wrapper composes known
  output with the context gate. Other scanner paths remain explicit inputs; see
  Stages 140 and 141.
  `duration.Engine.MatchProperNameTPPComponentArena` connects that projection
  to exact lookup and the marker gate, then reports when zero at `+0x18`
  bypasses later context predicates. Nonzero values can be resolved when the
  caller supplies scanner/helper outcomes; unavailable outcomes leave the
  contextual path pending. Full scanner production and runtime context-table
  snapshots stay unresolved; see Stages 131, 132, 134, 136, 138, 139, 140,
  and 141.
  `duration.Engine.MatchAndAppendProperNameTPPComponentArena` composes the
  zero-input bypass with the native matched-component source-row follow-up;
  its `WithContext` variant also evaluates the modeled nonzero path. Pending
  or rejected candidates leave the source arena untouched; see Stages 133 and
  134.
  `text.ApplyPaul2013TPPNumericAtom` ports
  the F/G per-token suffix writes, preceding-row flag, and optional WAB class
  byte from the native consumer. `text.ApplyPaul2013TPPTokenRange` exposes the
  same generic `FUN_1000e0c0` row writes for any caller-supplied byte value and
  affected range, so selectors A–E can reuse the known writer without
  guessing their value conversion or range production. WAB class lookup uses
  the mapped binary search recovered for `wab.txt2`.
  `text.MatchPaul2013TPPWindow` ports the bounded row scan, case-mapped
  `dash` fast-skip, `A`-row checks, cumulative cutoff, hyphenated compound
  writes, and last-match behavior from `FUN_1000e160`.
  `TPPDictionary.LookupNumericText` extracts exact F/G decimal suffixes, and
  `duration.Engine.LookupTypedTextCode` exposes those bytes through the
  loaded engine. `TPPDictionary.LookupSelectedText` and
  `TPPDictionary.LookupSelectedAtom` expose A-G results, including secondary
  atoms and AX's literal X payload. The duration engine exposes the suffix
  through `LookupSelectedTypedText` and the tag with suffix through
  `LookupSelectedTypedAtom`.
  `text.MatchPaul2013TPPWindowWithDictionary` connects direct F-atom matches to
  the loaded shared dictionary at index zero. Nonzero indexes
  fail closed because their alternate lookup resources are not loaded.
  Component-bit row application and selector construction in other callers
  remain unresolved, as do other multi-token dispatch paths and code
  semantics.
  `text.ApplyPaul2013TPPNumericTextRows` joins the direct F and G numeric
  paths to row updates and the no-code WAB class fallback. It can use the
  loaded one-column WAB table directly; a callback remains available for
  controlled comparisons. `WritePaul2013TPPNumericTextUpdates` applies the
  recovered `+0x20`, `+0x25`, and `+0x26` fields to a copied caller-row arena;
  `Paul2013TPPWindowRowsFromParserArena` projects native fields at `+0x00`,
  `+0x04`, `+0x08`, `+0x23`, and `+0x34`, and
  `duration.Engine.ApplyPaul2013TPPNumericParserArena` composes the projection,
  routing, and writes for a copied parser arena.
  `BuildAndRunPaul2013KnownModelContextPassesWithTPPFromParserRows` applies
  those writes after the supported parser-to-model context stages, matching
  their relative order in `FUN_1000e2f0`, with the recovered
  `FUN_1000e990` context-table flag scan between them. Ordinary-source and
  shared-arena sentence-segment variants connect the same sequence while
  retaining explicit parser controls.
  `text.SplitPaul2013ExceptionPhoneCodes`
  ports the `d` delimiter retention and destination-row advance in
  `FUN_1000ca50` for explicit per-row delimiter counts, and
  `duration.Engine.ResolvePaul2013PronunciationExceptionSequence` connects
  exception lookup to decoded phone rows with those caller-supplied counts.
  `text.ParsePhonePayload`
  parses its direct-ID and alternative-path forms.
  `text.Paul2013PhoneIDCodebook` provides the documented 256-slot compact-ID
  mapping. `PhonePayload.ExpandPronunciations` expands IDs through a five-byte
  codebook into internal symbol bytes, and
  `EmbeddedDictionary.ResolvePaul2013SurfacePronunciation` joins those steps
  for one canonical token. `PhonePayload.BuildPaul2013DictionaryPhoneRowsForState`
  ports both row branches. `PhonePayload.SelectPaul2013PronunciationByPathMarker`
  returns the phone string paired with a caller-selected path row. Source
  marker production and the conditional branch's selector/gate inputs remain
  explicit.
  `text.DecodeCMUPhones` labels the runtime-backed
  internal phone bytes and rejects structural or unknown values. The
  frontend now retains the four embedded pronunciation metadata flags on
  lexical tokens and carries them through sequence flattening; the legacy
  token-to-phone path copies these values, though their meanings remain
  unlabeled. They are distinct from the five duration-tree row values.
  `text.SelectLexicalPronunciations` flattens selected alternatives and retains
  each token's phone span and original separators.
  `LexiconFrontend.ResolveUniquePhoneSequence` remains a low-level
  unique-alternative helper. The duration engine selects ambiguous
  pronunciations with its loaded classifier and ranker.
  `text.BuildPaul2013TokenPhoneNeighborhoods` now derives identity/stress
  features over the complete utterance and groups them back by token, so
  adjacent words supply each other's phone neighbors. Separators remain
  available for boundary and prosody processing.
  `text.BuildPaul2013PhoneGroups` ports the vowel-anchor and permitted-onset
  boundary scan from `FUN_10013f30`, using the observed Paul phoneme labels and
  onset-cluster table. Token results include these phone-group spans and
  per-phone context projections. The engine ports `FUN_10013c00`'s duration
  row fields (nucleus stress, onset/coda flags, first/middle/final row state,
  phone offset/count) and its 1/2/3 onset/nucleus/coda labels, including the
  no-vowel fallback row. `text.BuildPaul2013TokenDurationInputsFromText` joins
  those fields with utterance-wide neighboring phones and the ordinary edge
  sentinels. Its explicit-marker form ports final boundary identities 40/42
  and position states 3/2 for the marker cases in `FUN_100135d0`.
  `duration.Engine.EvaluateWithMarkers` passes caller-supplied markers through
  both duration and pitch evaluation. `text.SplitPaul2013TokenPhoneBlocks`
  ports the explicit per-phone delimiter scan from `FUN_10012c70`.
  `text.BuildPaul2013ModelRecordPhoneGroups` composes that subdivision with
  the native `+0x329` marker stream and recovered phone-group writer for a
  supplied 0x3c0-byte model-state record. Decoded identity, vowel, and stress
  values can be supplied explicitly, or derived from `+0x2e8` by
  `BuildPaul2013ModelRecordPhoneGroupsFromSymbols` through the static Paul CMU
  codebook. Stage 10 runtime captures cover all 69 observed symbol codes and
  all 39 identity ordinals, validating this symbol-to-feature path without a
  direct table dump. The result now includes emitted native `0x1e` group rows,
  the parallel per-phone 1/2/3 labels, per-phone contexts, and the native group
  count; the no-vowel fallback preserves its uncounted row write.
  `text.BuildPaul2013ModelRecordDurationInputs` composes those record-derived
  fields with the nine-short writer and derives the right boundary ordinal
  from the native terminal marker. It also derives per-phone position states,
  the left boundary ordinal, and group count from the group rows and ordered
  record boundary; callers supply the record index and preceding marker.
  Results follow counted-group order and retain physical phone indexes;
  uncounted zero-vowel fallback rows emit no duration vector. Neighbor
  sentinels follow native first/last counted-group positions, even when a
  fallback span sits outside the counted groups.
  `duration.EvaluatePaul2013ModelRecordDuration` evaluates these vectors with
  the loaded duration trees and returns each phone's identity and index, tree,
  leaf, raw value, and exact input vector. The result remains a raw tree
  output without established physical timing units.
  Its stream variant consumes packed 0x3c0-byte records and derives record
  indexes and preceding terminal markers from stream order.
  `duration.EvaluatePaul2013ModelStateArenaDuration` reads the native signed
  count and record base from the complete model-state arena, then evaluates
  that stream without requiring callers to extract the records first.
  `text.BuildPaul2013ModelRecordPitchInputs` follows the native pitch loop's
  counted-group iteration and emits one 11-short row per group. The duration
  package evaluates the selected scalar/vector pitch pair from a model record
  or full model-state arena; results are raw vectors before the separate
  adjacent-group smoothing pass. The descriptor-aware smoother accepts its
  counts and rows explicitly while native pointer-target extraction remains
  unresolved.
  `text.ApplyPaul2013PhoneBlockDurationLimit` also inserts the captured
  `[` boundary when doubled per-phone values exceed the 1,000-unit block limit.
  `text.BuildPaul2013MarkerTreeGroups` ports the separate `FUN_10012df0`
  grouping rule for `\` and `]` boundaries and sums explicit per-phone
  `+0x94` bytes within each group.
  `text.ApplyPaul2013MarkerTreeDecision` applies `FUN_10012f00`'s result
  branch: results above 500 write `\` at the preceding record's `+0x3bd` and reset both
  running accumulators. `text.BuildPaul2013MarkerTreeInput` assembles the
  15-short vector from adjacent native phone records plus explicit token and
  split counters. `text.EvaluatePaul2013MarkerTreeInput` runs a caller-selected
  scalar tree and applies that decision; `text.RunPaul2013MarkerTreeToken`
  carries those counters through a supplied native-order token scan. The two
  pointer-backed string bytes and native tree selection remain caller inputs
  to that isolated token helper; the arena adapter below supplies iteration.
  `text.BuildPaul2013TokenBoundaryMarkerTransitions` ports the subsequent
  adjacent-record slash and previous-state-flag rules from explicit row bytes.
  `text.ReadPaul2013TokenBoundaryArenaInputs` now reads the signed marker code
  at record `+0x3ac`, transition state byte at `+0x3b0`, and duration-limit
  byte at `+0x95` from the counted model-state arena. The
  `text.BuildPaul2013TokenBoundaryMarkersFromModelStateArena` entry point
  composes the marker-table lookup with the transition in the order observed
  in `FUN_100130e0`; the later state/value arrays and marker dispatch still
  need their external producers. `text.ApplyPaul2013TokenBoundaryDurationLimitFromModelStateArena`
  applies the cumulative limit to those arena bytes and caller markers,
  returning an error if the first value would require a preceding boundary
  that the input does not contain.
  `text.BuildPaul2013TokenBoundaryStateResultWithTransitions` composes the
  state-array, marker, and post-dispatch adjacent-token passes for staged
  callers.
  `text.InitializePaul2013TokenBoundaryArena` writes the initial table and
  mode-dependent final marker before the explicit marker-tree stage.
  `text.FinishPaul2013TokenBoundaryArena` then writes native flags and markers,
  preserving transition flags through dispatch, applies the cumulative limit
  across up to 100 records, and rebuilds group descriptors. A nonnegative
  next-token value sets the preceding token's state to 2.
  `Paul2013FinalizedModelState.FinishTokenBoundaries` reads parser outputs at
  finalized state `+0x20` with `0x94` strides. Caller state initialization
  and the intervening tree stage remain explicit; see Stage 158.
  `text.PopulatePaul2013ModelPhoneGroupArena` writes group pointers, counts,
  labels, and recovered group-row fields while preserving unwritten bytes
  and the native uncounted fallback tail. `PreparePaul2013TokenBoundaryArena`
  composes that pass with marker initialization and pre-tree descriptors.
  Native model-symbol decoding now includes structural `M` through its
  identity/stress mapping, while preserving original record bytes. See
  Stages 159 and 162.
  `text.RunPaul2013MarkerTreeArena` derives groups, performs the reverse scan,
  resolves eligible `+0x2e4` string reads, and writes preceding terminals.
  `RunPaul2013TokenBoundaryPipeline` composes preparation, tree evaluation,
  and final dispatch; `Paul2013FinalizedModelState.RunTokenBoundaries` also
  derives parser outputs. The scalar tree and caller state arrays remain
  supplied inputs. Stage 160 records the corrected write target and tests.
  `PopulatePaul2013RecordGroupDescriptorArena` writes the native descriptor
  header, spans, and pointers. The composed pipeline performs this write both
  before tree scanning and after final marker dispatch; see Stage 161.
  `Paul2013FinalizedModelState.RunTokenBoundariesWithWorkspace` selects the
  loaded shared `engbi.tree3` and reads/writes native workspace arrays.
  `duration.Engine.RunPaul2013ModelTokenBoundaries` supplies that resource
  from `OpenPaul2013`. `ApplyPaul2013PositionStateProgramToWorkspace` reads
  native defaults and writes the recovered state-program arrays, preserving
  the shifted boundary-value tail cell. Range/event inputs remain explicit;
  see Stage 163 for loader and workspace evidence.
  State-program event results now retain the shared native cursor; upper
  clamps advance it and terminal processing resumes from it. Workspace
  adapters read/write selector cursors at `+0x1223f4/+0x122404`.
  `ApplyPaul2013PositionStateProgramFromModelArena` derives native row keys,
  selector-specific inter-record gaps, and the terminal gate. Reversed gaps
  from overlapping records match no events. See Stage 164.
  `Paul2013FinalizedModelState.ApplyPositionStatesAndRunTokenBoundaries`
  connects state production, optional interval mapping, slot-7 record flags,
  final mode, and boundary writeback. The duration engine exposes this through
  `RunPaul2013PositionStateBoundaries`. Interval mapping requires the full
  enclosing arena and validates its record extent; see Stage 165.
  `ReadPaul2013PositionStateDescriptorProgram` now reads the six native
  count/cursor/boundary/value descriptors and the recovered DLL clamp limits.
  `RunPaul2013PositionDescriptorBoundaries` connects those inputs to the
  loaded boundary pipeline; see Stage 166. The adjacent
  `ReadPaul2013PositionIndexWorkspaceTables` derives mapping pointers, source
  length and the shared remap gate. `RunPaul2013NativePositionWorkspaceBoundaries`
  composes both workspace adapters; see Stage 167. Pointer resolution and the
  producers of the native snapshots remain explicit.
  Stage 168 adds `PreparePaul2013PositionSegmentWorkspace` for native
  override selection, clamping, default propagation and sentinel reset.
  `RunPaul2013PreparedPositionWorkspaceBoundaries` composes preparation with
  the loaded workspace path. Successful native workspace calls now advance
  source position by consumed parser bytes; the source-parser loop and its
  zero-byte end-of-source branch remain outside these adapters.
  `RunPaul2013PositionSegmentDriver` now ports that native driver around an
  explicit parser callback: null source, empty-record retries, zero-byte
  completion, interval/default writes and counted state/mapping passes.
  Position adapters no longer require control fields populated only later
  by the native caller; see Stage 169. The loaded
  `RunPaul2013PositionSegmentBoundaryDriver` joins successful segments to
  phone/control projection and boundary writeback without fabricating a
  finalizer result; see Stage 170. Source parsing and native pointer ownership
  remain explicit dependencies.
  `RunPaul2013AlternateModelParser` ports the alternate `FUN_1003e070`
  row loop, including ordinary rows, angle contents, bracket skipping,
  whitespace and native low-short return codes. The shared scanner remains
  explicit. `DispatchPaul2013ModelParserWithAlternateScanner` supplies this
  implementation to the existing dispatcher unless overridden; see Stage 171.
  `RunPaul2013PrimaryModelParser` now ports the primary caller's ordered
  handler cascade, candidate rollback, marked-row boundary handling and
  special merge. `DispatchPaul2013ModelParserWithHandlers` supplies both
  inner loops while retaining explicit overrides. Reached missing handlers
  fail explicitly; their bodies and the general scanner remain dependencies.
  See Stage 172 for the native call order and validation boundary.
  `NewPaul2013TerminalParserHandler` now supplies FUN_10051a00's mode-1
  rescan, terminal row types and native signed returns. Its deeper punctuation
  recognizer remains explicit; see Stage 173.
  `NewPaul2013ModelScannerWithTerminalPrefixes` handles the scanner's complete
  whitespace-only and multiline early-return fields before delegating ordinary
  tokens. Those terminal paths run through the primary parser without fixture
  scanner bodies; see Stage 174.
  `ClassifyPaul2013TerminalContext` ports FUN_10052000's complete ordered
  six-record classifier. `NewPaul2013TerminalContextRecognizer` builds those
  records using the recovered backward copy and rolling scanner, then scans
  the following context and dispatches the classifier. It binds the complete
  source utterance and preserves the native field projection. Ordinary scanning
  and the model key/table fallback remain explicit; see Stages 175 and 176.
  `Paul2013TerminalAbbreviationTable` ports the model's abbreviation lookup
  across adjacent case variants. The case/status, single-punctuation and
  internal-dot features are also recovered; see Stage 177.
  `LoadPaul2013TerminalModel` now loads the four required text tables and
  `sbd.tree3`, supplies the recovered following-word classifier and complete
  16-short key, and evaluates the scalar fallback. Pass `model.Lookup` to
  `NewPaul2013TerminalContextRecognizer` to connect this path to the terminal
  handler. Ordinary scanning is the remaining dependency on this path;
  Stage 178 records the implementation and evidence boundary.
  `NewPaul2013ModelScannerWithASCII` now supplies standalone mode-0/1 letter
  and punctuation paths with recovered length/cost limits and raw result
  fields, plus mode-0x12 punctuation. Nonzero model pointers require the
  explicit post-token lookup; mixed tokens, numbers and other modes delegate.
  The recovered scanner, punctuation model and terminal handler run together
  on supported text without scanner fixtures; see Stage 179.
  `NewPaul2013ScannerIndexedTokenLookup` supplies both post-token search
  bodies and their recovered boundary decisions. Native index-pointer
  projection and mode-0x17 scanning remain explicit; see Stage 180.
  `NewPaul2013ScannerArenaTokenLookup` now reads the native index objects
  through a memory/string resolver and supplies that projection automatically.
  It snapshots active key ranges without reading inactive gaps; see Stage 181.
  Mode `0x17` ASCII words, hyphen joins and recognized apostrophe suffixes
  now run through the scanner and index-boundary check without fixtures.
  Numeric/non-ASCII mode-0x17 paths remain delegated; see Stage 182.
  Modes 0/1/0x17 now supply digit runs, decimals, three-digit comma groups
  and native ordinal suffixes. Mode 0's `401K`/`401(K)` forms return word
  status. Leading decimals require a known preceding-byte boundary;
  see Stage 183 for these numeric branches and remaining scanner work.
  Modes 0/1 also join letter components separated by dots and the native
  `int'` prefix. Mode 8 supplies ASCII words with lookup-gated hyphen/dot
  joins, its lowercase-to-uppercase dot boundary, apostrophe suffixes and
  the special `wi'` form. A rejected model join returns before the ordinary
  post-token lookup, preserving the native result flag. Mode-8 nonletter
  starts and non-ASCII continuations still delegate; see Stage 184.
  Additional ASCII word modes now preserve their individual apostrophe,
  hyphen, uppercase/digit and `United States` joins. Mode `0x17` also takes
  the common dot/`int'` join before its suffix dispatch; Stage 185 records
  the x86 evidence correcting the earlier mode-0x17 dot boundary.
  The remaining ASCII word branches now include mode-5 `pa'anga`, mode-0xf
  `A.M.`/`P.M.`, mode-0x1b abbreviation and slash/ampersand rules, and
  mode-0x1d digit joins. Other mode values use the native common dot/`int'`
  branch. ASCII nonletter starts still require their individual mode ports;
  see Stage 186 for this word-dispatch coverage and its limits.
  `text.BuildPaul2013InitialTokenMarkers` ports the 13-entry marker lookup at
  `0x10079a30` for explicit signed state codes; the code producer remains
  unresolved.
  `duration.Engine.EvaluateWithPhoneMarkers` also accepts those bytes and
  rebuilds duration and pitch rows within each resulting block. The ordinary
  evaluator supplies zero markers, so it preserves one block per token;
  source text does not yet produce the state codes, special phone markers, or
  budget inputs.
  `text.ParsePaul2013ContextCodes` separately parses the bounded code string
  consumed by `FUN_10016c90`, retaining its opaque code bytes, `d`/`c`
  preceding-phone markers, and `M` flags. It does not derive the string or
  connect it to model-context construction. `text.SummarizePaul2013ContextCodeState`
  ports that function's explicit per-token state flags and final mode-byte
  dispatch.
  `text.PopulatePaul2013ModelStateRecordsFromParserRows` now projects its
  recovered writes into counted 0x3c0-byte records from caller-produced
  0x94-byte parser rows. It writes the three address fields through an
  explicit resolver and preserves all unobserved record bytes. The per-token
  state values and terminal pitch-value count remain caller inputs; the
  function derives the final mode byte at arena +0x4770a from the projected
  last record. This projection is based on Ghidra pseudocode and has no
  captured native-record parity claim.
  `text.InitializePaul2013PositionStateArrays` reproduces the per-row defaults
  from `FUN_10022dc0`, and `text.ApplyPaul2013PositionStateProgram` executes
  the recovered `FUN_10022970` range and event passes in native order when the
  caller supplies their tables and intervals. Its optional final mapping reads
  model-state row indexes at +0x64c/+0x650 and applies the supplied lookup
  tables. Parser-row production and those model-state inputs remain unresolved.
  `text.ApplyPaul2013PositionStateProgramFromParserState` now reads the raw
  `0x94`-stride parser-state rows and supplies their absolute intervals to the
  program's selector-3/4 event passes, rejecting row-count mismatches. The
  bounded ordinary-source adapter applies one supplied program per captured
  sentence segment and preserves source-byte origins across `.?!` resets. The
  synthesis adapters also seed the three initial arrays from independently
  selected and clamped pitch, speed, and volume controls, including a single
  call from bounded ordinary source. `text.NormalizePaul2013PositionStateProducerArrays`
  ports `FUN_1001d5d0`'s paired-boundary compaction and negative-value fill for
  modes 0-2 after their arrays are populated.
  `text.ApplyPaul2013PositionStateProgramWithProducerArrays` connects that
  step to the state-range consumer. Range/event values, gates, the matcher,
  and the raw parser-state producer remain explicit inputs; see Stage 142.
  `duration.EvaluatePaul2013SourceText` is a low-level
  unique-pronunciation helper. `Engine.Evaluate` handles ambiguous entries
  before deriving duration inputs; neither path needs per-phone metadata from
  its caller.
  `duration.EvaluatePaul2013Text` dispatches each phone to the
  observed duration-tree family and evaluates its nine-short vector. The result
  is still a raw tree value; its time unit and timeline conversion are unknown.
  The older `text.BuildPaul2013TokenDurationInputs` remains a low-level trace
  helper.
  The legacy path is partly traced:
  `FUN_10012c70` groups phone records and calls `FUN_10013c00`, while
  `FUN_100135d0` derives neighbor sentinels and tree position from token and
  terminal markers, and reads the row and auxiliary values produced above.
  The ordinary path now evaluates duration and pitch trees and derives the
  lexical-path adjacent-vector pitch averages. Their native descriptor/group
  mapping and physical timing use are not yet connected to a synthesis
  timeline.
  `CMUPhone.TreeFeatures` returns the original internal phone byte alongside
  the directly observed one-based identity ordinal and vowel stress value.
  These are distinct representations.
  `BuildPaul2013PhoneNeighborhoods` arranges those known values around each
  phone and marks sequence edges; it does not
  merge word spans or invent tree-specific feature positions.
  `ApplyObservedPaul2013OnsetIdentityFeature` writes the controlled consonant
  ordinal into position 0 of a caller-selected onset-scalar input; tree-family
  selection and the remaining fields are caller work.
  `ApplyObservedPaul2013InteriorDurationPhoneFeatures` fills the observed
  current/previous/next identities and current stress in positions 0–3 for an
  interior phone. `BuildObservedPaul2013DurationTreeInput` remains a low-level
  trace helper for explicit inputs; ordinary text uses
  `BuildPaul2013TokenDurationInputsFromText` to derive those fields. The
  interior convenience builder remains available for controlled tree probes.
  `ApplyObservedPaul2013PairedScalarFeatures` writes the captured
  second/first identity ordinals to positions 0/1 and fills position 3 only
  for a second phone labeled AH. `ApplyObservedPaul2013AHPairFeatures` retains
  the stricter AH-only contract. Both leave unresolved fields untouched.
  `ApplyObservedPaul2013AHStressFields` writes AH stress to position 2 of the
  traced third-scalar and 12-value vector inputs and copies the supplied
  preceding scalar result into vector position 11; it does not evaluate trees.
  `text.Paul2013EmbeddedKeyTables` provisions the documented key-transform
  tables, and `EmbeddedDictionary.ResolvePaul2013Surface` performs canonical
  token lookup without injected tables. `ExpandPaul2013Number` handles
  observed plain/signed integers, grouped commas, and decimal forms. It uses
  the captured four-digit year rule, cardinal groups through 15 digits, and
  digit-by-digit speech for longer or leading-zero values and decimal places.
  `ExpandPaul2013Ordinal` handles ordinal digit suffixes from `0th` through
  `31st`; the 1–31 spelling rules are inferred from English morphology, with
  only the second/fifth words directly represented in captured outputs.
  `ExpandPaul2013ClockTime` handles 12-hour `H:MM` values with nonzero minutes;
  only `3:45` is directly represented in captured output. The captured
  currency example `$5.00` and percentage example `25%` are now handled by
  dedicated normalizers; singular/plural and nonzero-cent wording plus the
  wider numeric grammar are implementation inferences. The captured
  `01/02/2024` form expands to `January second twenty twenty four`; the
  `MM/DD/YYYY` widths and component bounds beyond it are inferred, and
  month-specific calendar validity is not checked. Other date and time forms
  and general normalization remain unsupported. `ExpandPaul2013Telephone`
  handles the captured `555-1234` form as “five hundred fifty five to twelve
  thirty four”; other digit groups in this layout follow the cardinal
  expander as inference, and other telephone layouts remain unsupported.
  `LexiconFrontend.ResolveTextWithPaul2013InlinePauses` preserves captured
  `<vtml_pause time="N"/>` events and their `N * 16` frame durations.
  `text.Paul2013InlinePauseModelSourceRule` wires that exact form into the
  recovered longest-key dispatcher and returns each timing event to its
  caller; other attributes and VTML tag effects remain unsupported.
  `text.Paul2013InlineMarkModelSourceRule` also ports the captured unnamed
  and named self-closing marks, including the empty-name no-record case,
  case-insensitive `name`, and 512-byte truncation to 511 bytes. It emits a
  source-offset event without changing lexical text; output-frame placement
  remains caller-owned.
  `LexiconFrontend.ResolveTextWithPaul2013InlineMarks` resolves text around
  those tags, retains exact mark source offsets and record fields, and carries
  them through explicit pronunciation selection. A token span may enclose a
  mark tag when the tag falls inside that lexical token; audio-frame placement
  remains unresolved.
  `text.ExpandPaul2013CapturedVTMLSubstitutions` replaces the body in the
  directly captured `<vtml_sub alias="...">surface</vtml_sub>` form with its
  ASCII alias; by itself it returns no source-span mapping. Escapes and general
  VTML parsing remain unsupported.
  `LexiconFrontend.ResolveTextWithPaul2013CapturedVTMLSubstitutions` resolves
  that form through the embedded dictionary and maps token spans to the alias
  bytes in the original input, rejecting tokens that cross a rewritten span.
  `duration.Engine.EvaluateWithInlinePauses` carries those events through
  dictionary pronunciation selection and both duration/pitch tree families,
  returning the sequence and token results with pause boundaries intact.
  `LexiconFrontend.ResolveTextWithPaul2013CapturedVTML` composes the captured
  pause, mark, and substitution forms in one source pass, retaining original
  token spans and pause/mark records. `text.SelectPaul2013CapturedVTMLPronunciations`
  carries both event streams through selected pronunciation assembly.
  `duration.Engine.EvaluateWithCapturedVTML`
  carries the parsed text through automatic pronunciation selection and the
  duration/pitch trees. Pause frames are still not inserted into rendered PCM,
  and marks remain source-positioned because their frame mapping is unresolved.
  `synthesis.InsertPaul2013InlinePauses` inserts zero-valued PCM at explicit
  token-boundary frame offsets. The renderer does not yet produce those
  offsets from selected text rows, so pause events are not wired into
  `pipeline.Pipeline`. `InsertPaul2013InlinePausesAfterTokenChunks` derives
  offsets when input PCM is already separated into one chunk per lexical token.
  `text.LexiconFrontend.ResolveText` tokenizes ASCII surfaces and returns
  dictionary alternatives with original model-coded bytes and CMU labels. It
  preserves inclusive source-byte spans through pronunciation selection;
  expanded number words inherit the source token span. These provenance spans
  are distinct from parser rows. `text.BuildPaul2013OrdinaryParserOffsetSegments`
  derives exclusive offsets for ASCII-letter words and unsigned cardinals;
  `BuildPaul2013OrdinaryParserOffsetRows` also carries the captured opaque
  `A`/`D`/`S` parser-type byte and populated `+0x34` text for supported
  ordinary and numeric forms. Numeric row strings come from the bounded number
  normalizers, matching the Stage 20 examples while leaving broader grammar
  behavior as inference. In captured numeric rows, the leading sign, final
  currency row, and telephone
  hyphen use `S`; percent suffixes use `A`. These
  offsets retain commas inside a segment and reset after `.?!`, matching
  captured `FUN_1003e240` rows. Numeric multiplicity is directly captured for
  `1`, `12`, `123`, `1234`, and `2024`; other unsigned-cardinal counts follow
  the number expander as inference. Captured offsets also cover `+12`, `1.25`,
  `$5.00`, `25%`, `21st`, `01/02/2024`, `3:45 PM`, and `555-1234`. Telephone
  rows split into cardinalized groups with a separate hyphen row. The frontend
  now speaks the captured form using those group words joined by “to”. General
  abbreviation, exception, TPP, and full parser-state production remain
  unresolved. The `FUN_1000e570` model-context classifier now uses the
  extracted sorted tables and mapped-string rules by default. The fixed
  `FUN_1000ea20` projection is ported by `text.NormalizePaul2013ModelPhoneRows`;
  `text.SelectPaul2013ModelPronunciationAlternative` supplies the ordinary
  tree-scored `FUN_100068b0` path when given `engbi.tree3`; the model text
  parser wires it when `PronunciationTree` is provided. Its short-input fast
  path is ported from the mapped `a` exception. The shared token-row reader
  also supports thirty-nine direct `FUN_100049b0` word branches: `anti`, `conflicts with`,
  class-gated `supply`, `does` after `much` or before `not`, the contextual
  `converts` and `minute` rules,
  `wound up/down/through`, the `number` selector rule, `perfect` before punctuation,
  `polish Jewish`, terminal/punctuation-followed `learned`, and `present` with
  `to` two rows later, the recovered `can` context cases, selected `close`
  contexts including `at hand`, `it up`, and the preceding pair `to a`,
  selected `closer` contexts,
  `mouth to mouth`, `brain reading`, `record(s)` with a number or `it`/`to`
  context (or `%` after the ten-word subject set), `refuse(s)` before `to`,
  `import` after `export`/`the`/`an` or before `bank`, `increase(s)` after
  `the`/`an` or before a class-3 surface, `resume(s)` after two rows before
  sentence punctuation when preceding rows pass class-3/class-11 gates,
  `transform(s)` before `into` two or
  three rows later, `use` after `still`, `wind` before `up`/`down`/`through`,
  `de` beside an uppercase-initial surface, `job hunting`, all recovered
  `laden` selector conditions,
  `house` beside uppercase-initial surfaces or before `arrest`/`number(s)`,
  `dogged` with its recovered phrase-window rules, `invalid` after the recovered
  class-3 and negation contexts, `lied` with its recovered context gates,
  `dove into`, `elaborate into`, and selected `live/lives`
  contexts (`lives on`, `live on TV`, `live(s) it up`, `live(s) up to`,
  `live(s) to` a surface longer than `self` ending in `self`, and the
  following-literal/feature-4 selectors and the recovered previous-surface
  exception set),
  `separate(s)` after `be` or a class-3 predecessor, or before its recovered
  forward-context set, `subject to` with its
  class-3 contexts, `contest(s)` before class-3 surfaces, and uppercase-initial
  `interstate` before another uppercase-initial surface or `interstate(s) to`.
  For other words, the selector ports additional `FUN_10004300` gates in
  native order: before `of` (except after `nice`) it tries 0x13, 0x16, 0x14,
  0x15; before `'s` it tries 0x14, 0x15, 0x13, 0x16; compound prefixes
  (`de`, `inter`, `re`) before a same-index hyphen row request 0x0d. After
  `less`, `more`,
  `so`, or `very`, it chooses the first class-2 path alternative, then class 8.
  After the recovered 20-word adverb set, it tries `(` when the preceding
  triplet has class 3, then `%`, `*`, and `&`. A following `it`, class-12 or
  class-7 word, or one of `him`, `me`, `them`, `us` activates the `%` gate
  only when class-10 paths occur in the selected class-9/12 alternative or
  are absent. Previous words accepted by `FUN_10010690` select the first
  class-4/7 path, then 0x13, unless a class-5/6 path occurs in another
  alternative. The `have`/`has`/`had` branch tries `(`, and the recovered
  `FUN_10010760` verb set before a surface longer than five bytes ending in
  `ing` selects the apostrophe path. `FUN_10010840`'s word-pair gate is also
  applied at these branches; a matching previous pair selects `%`. A `how`
  context immediately before the current
  row or two rows earlier selects 0x0e when another row follows and the
  current surface is class 8. A preceding triplet classified as 7 uses the
  same guarded class-4/7-or-0x13 selection as `FUN_10010690`.
  This covers the recovered `FUN_10004300` fallback cascade. Its call path has
  static support but no direct runtime row/choice comparison; the remaining
  name/context branches belong to `FUN_100049b0`. Evidence is in
  `tools/revkit/work/reports/stage6-phone-record-consumers.txt`.
  Direct word cases that do not select a path now enter this generic fallback,
  matching `FUN_100049b0`'s final call to `FUN_10004300`.
  `text.FindPaul2013PronunciationForPathCode` ports its reusable
  `FUN_10010640` helper: it matches the requested path code against each
  alternative through the DLL class table and returns the first aligned phone
  string.
  `text.WritePaul2013ParserOffsetRows` writes the known start,
  end, and +0x24 discriminator fields into initialized `0x94` parser rows while preserving
  unported fields. Its optional `Text` and `AuxiliaryText` values write the
  C strings at `+0x34` and `+0x52` consumed by `FUN_1000d190`. Stage 20 traces
  show expanded numeric lexemes at `+0x34` across rows sharing the original
  numeric source span. The ordinary row builder uses the bounded number
  normalizers to populate `Text` with expanded lexemes; broader numeric
  grammars remain inference. Auxiliary text remains caller-owned.
  `WritePaul2013ParserOffsetRowsWithTypeWrites` applies sparse caller-supplied
  writes to `+0x2c` in order, without deriving values from the discriminator
  or offsets.
  `WritePaul2013ParserOffsetRowsWithTypeGates` composes the captured terminal
  and comma writers plus the static dot writer with explicitly row-indexed
  scanner gates; it preserves gate order and does not infer scanner
  acceptance. The dot branch has no direct runtime observation. Other
  row-type producers can be composed with
  `WritePaul2013ParserOffsetRowsWithTypeApplications`, which invokes the
  corresponding existing row-type helper for an indexed row. Applications
  must still provide any other source fields and native gate values read by
  their helper.
  `BuildPaul2013OrdinaryParserOffsetRowArena` composes that
  writer with the captured offset builder for a single sentence segment, and
  `RunPaul2013ModelTextParser` uses it when callers omit parser offset rows.
  `BuildPaul2013OrdinaryParserOffsetRowArenaWithTypeGates` and its
  `WithTypeApplications` variant also compose explicit row-type gates or
  existing row-type helpers into that source path.
  Multi-segment input remains fail-closed because each captured segment resets
  its coordinate origin. `ApplyPaul2013ParserTerminalRowType` ports the direct
  punctuation writes 2/3/4 from `FUN_10051a00`, and
  `Paul2013ParserPreviousRowTerminalPunctuation` ports the `FUN_10062f10`
  membership check against `.?!;`; `Paul2013ParserTerminalGateFromPreviousRow`
  derives its zero-result gate, and its whitespace-prefix variant derives the
  two-LF status 8 path. Other scanner statuses and the alternate
  `FUN_10051cc0` result remain explicit inputs.
  `ScanPaul2013ModelParserWhitespacePrefix` ports only `FUN_1005a350`'s
  leading whitespace loop, including its LF count, single-LF marker, and
  two-LF status 8 path; it does not parse tokens.
  `ApplyPaul2013ParserCommaRowType` ports its comma-path writes 5/12 from
  `FUN_100544f0`; it applies the DLL's mapped comparisons against `too` and
  `either`. `ApplyPaul2013ParserDotRowType` ports its statically recovered dot
  write of 1, which has no direct runtime observation yet.
  `ApplyPaul2013ParserQuoteRowType` also ports a static double-quote write of
  1; its preceding scanner path remains an explicit gate. Code 12 is also
  static-only. `ApplyPaul2013ParserMultiSourceRowType` ports a static code-10
  write after a successful generic row append when the input has multiple
  source rows. `ApplyPaul2013ParserUnmatchedRowType` ports the static lookup-
  miss write at `0x10054f28` against a separate table.
  `ApplyPaul2013ParserCloseDelimiterRowType` ports the static close-delimiter
  write at `0x100557df`. `ApplyPaul2013ParserModeOneRowType` ports the static
  mode-1 write at `0x10055a8e`. `ApplyPaul2013ParserAcceptedPathRowType`
  ports the final count check and write at `0x10055a4a`; earlier scanner
  conditions remain caller inputs. `ApplyPaul2013ParserBoundedScanRowType`
  ports the final boundary check and write at `0x10055677`, with raw values
  explicit. `ApplyPaul2013ParserPositiveScanDotRowType` ports the additional
  dot write at `0x10056077`. `ApplyPaul2013ParserHyphenFallbackRowType` ports
  a separate static type-7 write at `0x1004eb4d`. Stage 20 write watches
  confirmed the 2/3/4/5 assignments. `ApplyPaul2013ParserQuoteCoordinateRowType`
  ports repeated quote/count/negative-coordinate type-1 fallbacks across the
  scanner routines. `ApplyPaul2013ParserInputFlagRowType` ports a separate
  post-append type-5 write at `0x100445b1`. Two post-append type-7 paths at
  `0x10057722` and `0x100577db` are exposed as separate helpers; the latter
  requires another input row. `ApplyPaul2013ParserZeroInputFlagRowType` ports
  two post-append type-8 writes when the input short at `+2` is zero.
  `ApplyPaul2013ParserHyphenModeRowType` ports a post-append type-11 write
  when the native mode is 2 and the following byte is `-`.
  `ApplyPaul2013ParserExactInputFlagRowType` ports a separate type-5 write
  when the input short at `+2` equals exactly one.
  `ApplyPaul2013ParserContextualRowTypeOne` ports the context-gated type-1
  assignment at `0x10042e8f`, preserving the compared state values as raw
  numeric inputs. `ApplyPaul2013ParserModeOnePriorKeyRowTypeOne` ports a
  separate type-1 write at `0x100430ce` after a row-key match, mode 1, and two
  positive state values. `ApplyPaul2013ParserPreAppendRowType7` and
  `ApplyPaul2013ParserSplitAppendRowType7` port the pending and final-row
  type-7 stores at `0x1005782b` and `0x1005788b`.
  `ApplyPaul2013ParserAcceptedScannerRowTypeOne` handles the zero-field write
  at `0x1003d4b6`, and `ApplyPaul2013ParserTypeTwoReset` handles the gated
  clear at `0x1003da7a`. Other candidate type stores
  across the parser are
  inventoried in `tools/revkit/work/reports/stage46-parser-row-type-write-inventory.c`;
  only the byte-stride candidate at `0x1005551d` remains unattributed, and
  incoming predicates for static stores remain caller inputs.
  `text.AppendPaul2013ModelSourceRow` ports the bounded append writes in
  `FUN_10044fd0` and the optional auxiliary-string write from `FUN_10045070`.
  `AppendPaul2013ModelSourceRowsSplit` also ports `FUN_100451e0`'s space/tab
  splitting, and `AppendPaul2013ModelSourceRowsBackslash` ports
  `FUN_100452c0`'s auxiliary-string splitting. These preserve fields those
  helpers leave untouched, including `+0x2c`; they do not construct that
  value. `AppendPaul2013MatchedComponentSourceRows` composes the accepted
  `FUN_10034180` candidate follow-up from `FUN_10034110`, including its raw
  append constants and last-row `+0x1e` marker write. Candidate production and
  runtime row parity remain open; see the Stage 130 report.
  `AppendPaul2013ParserRowTypeMarker` applies the static punctuation
  suffix mapping used by `FUN_1000d190` for selected nonzero row types under
  the native header-value gate; it bounds the 32-byte local string buffer.
  `TransformPaul2013ModelParserSourceRow` composes the recovered fields,
  marker append, supported source transform, and source-type classifier for
  one 0x94-byte row. The non-final punctuation lookup remains an explicit
  callback; `TransformPaul2013ModelParserSourceRowWithEmbeddedDictionary`
  connects it to `EmbeddedDictionary.ContainsPaul2013Surface`.
  `text.DispatchPaul2013ModelSourceRule` ports the generic
  `FUN_1002e990` longest-key selection and restores both cursors when the
  selected handler declines. Its rule table and handlers remain caller inputs;
  see Stage 143.
  `text.BuildPaul2013ModelTokenRowsFromParserRows` traverses initialized rows,
  repeats the source transform while a remainder remains, performs exact
  embedded-dictionary lookup, and writes token rows using the preceding-row
  and mode gates. Parser-row production and the earlier normalization/TPP
  cascade remain upstream requirements; this composition has not had a
  row-by-row runtime comparison.
  `text.BuildPaul2013ModelPhoneRowsFromParserRows` adds the recovered fixed
  token-to-context-row projection, ordinary multi-pronunciation selector,
  thirty-nine direct native name/context word branches, and the early
  `August` shortcut after `an`/`the`, beside comma-initial rows, at the final
  row, before terminal punctuation in the penultimate row or a class-4
  surface, or beside one-byte rows with attribute mask `0xc0` (path byte
  0x1e). It also ports the recovered class-2/3, word-pair, and character-table
  gates from `FUN_10003ff0`, selecting 0x07 when no 0x1e gate matches. This
  shortcut is static-only; runtime choice parity remains unverified. The generic
  `FUN_10004300` fallback cascade is statically ported, including its
  cross-alternative class guards; no direct runtime row/choice comparison is
  available. Additional direct word/context cases in `FUN_100049b0` remain
  unsupported. The `close` branch also
  covers `at hand` (0x0e) and `it up` ('%').
  `text.BuildPaul2013ModelPhoneRowsFromOrdinarySource` connects the bounded
  word/number offset-and-text producer to that path, deriving `+0x08` from
  each source string and requiring caller-supplied `+0x2c`, `+0x30`, and
  `+0x52` controls. Multiple sentence segments remain unsupported by this call.
  `BuildPaul2013ModelPhoneRowsFromOrdinarySourceWithTypeApplications` also
  applies the existing row-type helpers in order over the supplied initial
  `+0x2c` values before projecting tokens and contexts; mode and auxiliary text
  remain caller controls. The
  `BuildPaul2013ModelPhoneRowsFromOrdinarySourceWithDefaultParserControls`
  variant uses mode `0xff` and an empty auxiliary string for the bounded
  ordinary inputs, while keeping initial row types and type-helper gates
  explicit.
  `text.BuildPaul2013ModelPhoneRowsFromOrdinarySegments` returns one projected
  result per segment, preserving the captured offset reset; its results remain
  segment-local. The
  `BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArena` entry point
  preserves those offset resets while flattening rows in order and projecting
  them once through the shared model counters (Stage 109).
  Its `WithTypeApplications` variant applies each segment's recovered row-type
  helpers using local row indexes and row counts, then maps successful writes
  into the flattened arena. The
  `WithDefaultParserControls` variant applies mode `0xff` and an empty
  auxiliary string to those rows while retaining explicit initial row types.
  `BuildAndRunPaul2013KnownModelContextPassesFromOrdinarySegmentsInSharedArena`
  carries that model through the known ordered passes and returns rows still
  requiring unported handlers in `UnresolvedRows`.
  Its `WithTypeApplications` variant runs the same passes after composing the
  per-segment row type helpers.
  `BuildAndRunPaul2013KnownModelContextPassesWithTPPFromOrdinarySegmentsInSharedArenaWithTypeApplications`
  continues through the supported numeric F/G TPP and WAB updates in the
  native order; other component-based TPP effects remain open.
  Its `WithDefaultParserControls` variant supplies the captured mode `0xff`
  and empty auxiliary string while keeping initial row types and type-helper
  inputs explicit, then runs the same ordered model-context and TPP passes.
  `duration.Engine.BuildPaul2013ModelPhoneRowsFromParserRows` and its ordinary
  source variant provide these compositions with the engine's loaded
  dictionary and pronunciation classifier; they do not invent the unresolved
  parser-row controls or run the upstream normalization/TPP cascade.
  Their result retains the parser-row arena for downstream processing.
  `BuildAndNormalizePaul2013ModelPhoneRowsFromParserRows` and its ordinary
  source variant compose the projection with `FUN_10009030` over explicit
  context-row indexes; the caller still supplies the native eligibility/order
  result, and the earlier exception/TPP stages remain outside this path.
  `text.RunPaul2013ModelTextParser` also ports the
  `FUN_1000e2f0` orchestration and error returns; its phone-row parser,
  row processor, and TPP parser remain explicit callbacks. The normalizer uses
  the built-in row projection with supplied parser rows or derives them for a
  supported single sentence segment.
  `text.UpdatePaul2013ModelContextTableFlag` ports `FUN_1000e990`'s gate,
  ordered classifier calls, P/N result handling, and output short, using the
  built-in classifier unless a comparison override is supplied.
  `text.RunPaul2013ModelParserFinalizer` connects that wrapper to the
  `FUN_1003e240` caller-arena writes and adapts the result to
  `DispatchPaul2013ModelParser`.
  `text.FinalizePaul2013ModelParserRows` ports the post-parser grouping and
  memory writes in `FUN_1003e240` against an explicit caller arena and the
  model context-table bytes. Its helper implementations, source-row
  producers, and full-arena runtime comparison remain open.
  `text.FinalizePaul2013ModelParserRowsAndBuildModelStateRecords` connects
  the finalized +0x66 phone strings to record construction, transferring the
  count and inclusive intervals before optional index remapping. Per-token
  state values and pointer address mapping remain explicit; see Stage 157.
  Resolved tokens retain
  the ordinary/final dictionary phone-row projection from `FUN_1000d450`, and
  selected token spans carry an independent copy for later parser stages.
  `text.WritePaul2013DictionaryPhoneRow` writes its established `0x554` row
  fields and leaves native-unwritten bytes untouched; token indexes, surfaces,
  and source-marker production remain caller inputs.
  `Paul2013DictionaryPhoneRows.WithPaul2013ModelSourceType` applies
  `FUN_1000fd20`'s result-byte classifier when the caller supplies its
  transformed source string. Its convenience wrapper is ASCII-only; callers
  can use the full recovered table through `ClassifyPaul2013ModelSourceType`.
  `TransformAndClassifyPaul2013ModelSource` composes the direct
  `FUN_1000d640` copy/trailing-pair split paths with that classification. Its
  scanner path also handles unpunctuated ASCII strings by trimming leading
  spaces, copying the remaining bytes, and clearing the consumed source.
  It also splits terminal punctuation and terminal hyphens according to the
  final-row gate. Callers can supply an embedded-token lookup for the non-final
  punctuation decision; the ASCII helper fails closed when that lookup is
  needed but absent. `Paul2013UnsignedCharacterAttributeTable` exposes the raw
  256 table bytes; signed-char string paths use
  `Paul2013ExceptionCharacterAttributes`, which maps high-bit bytes through
  the verified zero prefix. Stage 108 ports the scanner's ordinary high-byte
  copy behavior and its directly compared `0xa2 0xfe` split. This scanner subset is
  static-only; its control-flow summary is
  `tools/revkit/work/reports/stage35-model-source-scan.c`,
  `tools/revkit/work/reports/stage36-model-source-terminal-splits.c`, and
  `tools/revkit/work/reports/stage37-character-attributes.c`. The related
  pseudocode excerpts are `tools/revkit/work/reports/stage32-model-source-type.c`
  and `tools/revkit/work/reports/stage33-model-source-transform.c`.
  `text.BuildPaul2013PhoneContextRows` projects an already selected lexical
  sequence into the known fields of the counted `0x70` phone/context rows.
  `BuildPaul2013PhoneContextRowsFromParserRows` derives the native `X` marker
  from parser-row byte `+0x24`; token-to-parser-row indexes remain explicit.
  The literal
  contraction branches of `FUN_100091b0` are ported, while its generic
  dictionary/TPP normalization path and remaining row mutations are open.
  `text.DerivePaul2013PhoneContextPreparationFromRows` derives every
  `FUN_100091b0` low-short return from the source kind, subtype, and adjacent
  source indexes. `text.ExtractPaul2013ExceptionRowInputs` reads surfaces and
  X retry markers from caller-selected rows, and
  `ResolvePaul2013PronunciationExceptionFromPhoneContextRows` composes those
  fields with the counted-row gate and exception lookup. Row eligibility,
  ordering, and production remain upstream inputs.
  `text.EncodePaul2013ContextString` ports the ASCII letter expansion and `d`
  delimiter writer from `FUN_10009cd0`, including its 64-byte output cap and
  `+0x20` row-flag update. Its 26 expansion records are copied from DLL table
  `0x10077e58`. `text.BuildPaul2013ContextCodesFromSourceString` feeds that
  result into `FUN_10016c90`'s per-phone marker scan. Non-ASCII mappings,
  token-state production, and the rest of the `FUN_100091b0` cascade are still
  open.
  `text.ApplyPaul2013LiteralContextHandler` also ports its direct contraction
  and suffix writes for the observed literal surfaces; other strings still
  depend on the unported generic dictionary/TPP path.
  `text.BuildPaul2013ExceptionDispatchGateFromRows` extracts the exception
  gate fields from `0x94` parser rows and `0x70` phone/context rows, and the
  duration engine can compose that gate with its loaded exception lookup.
  Unsupported numeric
  separators still fail closed.
  The duration engine now builds a selected phone sequence without a
  caller-supplied choice list. It is not yet connected to `pipeline.Pipeline`
  because model-context projection and unit selection are still missing.
- `text` implements the observed seven-byte to five-byte class-key transform,
  two ten-byte feature views, and their weighted mismatch tables.
  `text.PreparePaul2013SelectionContextsFromPackedContextPhoneGroup` now
  composes the recovered packed-category pass and extracts its per-phone
  seven-byte rows as selection contexts in record/phone order. It requires an
  already-selected native record group; source record production and group
  dispatch remain upstream inputs.
- `voice.BuildPaul2013ClassCatalog` reconstructs sorted exact-key classes and
  their member-unit references from the four read-only unit indexes. Its
  `LookupContext` method applies the key transform to an existing context and
  returns that exact class when present; `LookupPrefix` reproduces the native
  one-to-five-byte class-key range search.
  `selection.RankClassRecords` applies the observed count and population
  limits while retaining unit membership, including the native 10,000-class
  scoring bound used when its source-order fast path does not apply. Feature
  views rank only the first 10,000 entries; oversized key-only pools fail
  closed because the DLL sort consumes an uncharacterized scratch tail. The
  recovered insertion path also determines equal-score order for lists up to
  17; the larger-list implementation now matches a natural 75-entry native
  partition permutation and a controlled 75-entry heap-fallback permutation.
  The second vector uses process-local score replacement, and does not show
  that ordinary workloads take the fallback. Exhaustive tie-order parity across
  call sites remains unverified. Catalog construction matches the native
  61,566-key table, three captured key-to-ID lookups, and
  their captured unit expansions. `selection.ExpandClassRecords`
  flattens ranked memberships in source order under an explicit unit cap; it
  does not construct contexts or candidate scoring metadata. The catalog
  does not score individual units. `selection.ExactContextUnitCandidates`
  maps an existing seven-byte context to its exact class and returns bounded
  unit references; it does not broaden missing matches.
  `selection.LookupPaul2013PackedContextCandidates` composes the packed
  per-phone context pass with exact catalog lookup and unit expansion, keeping
  misses as empty pools. It does not run native whole-position queries or
  fallback rows.
  `selection.LookupPaul2013PhoneMarkerRecordCandidates` composes the ordered
  record-group producer, marker prepass, packed-context preparation, and exact
  lookup for a full supplied native record stream. The text-level
  `PreparePaul2013PhoneMarkerRecordGroups` result now also exposes each group's
  extracted seven-byte `SelectionContexts` alongside its prepared record rows.
  It preserves groups and
  exact misses; raw record production and pointer-backed key bytes remain
  caller inputs, and whole-position broadening and scoring remain separate.
  `selection.LookupPaul2013PhoneMarkerArenaCandidates` includes the shared
  phone-state reset for a contiguous native record arena.
  `selection.LookupPaul2013PhoneMarkerModelStateArenaCandidatesWithPointerResolver`
  reads the signed count at +2 and record base at +0x64c from the full model
  state arena, then connects it to that same marker/group/exact-lookup path.
  The marker-string address space, row producer, and all later query-state
  inputs remain caller supplied.
  `selection.RankPaul2013PhoneMarkerRecordContinuity` feeds those exact pools
  into the recovered continuity pass and restores the native group alignment.
  It requires caller-produced context modes and fails on exact misses until
  native fallback rows have been applied.
  `selection.Paul2013ContextModesFromTransitionStates` reads mode byte `+4`
  from the compact six-byte state rows (`0` for accepted whole-position rows,
  `1`/`2` for fallback rows). Its `Rank...FromTransitionStates` and
  `Score...FromTransitionStates` compositions connect those rows to continuity
  ranking and local scoring when the state order matches candidate-position
  order.
  `selection.ScorePaul2013PhoneMarkerRecordCandidates` then applies the
  bounded local-score stage, including protected-prefix expansion for
  oversized pools. State-derived feature contexts and candidate-key maps stay
  caller supplied. `selection.FinalizePaul2013PhoneMarkerRecordShortlists`
  carries those rows through native tail ordering and protected-candidate
  restoration while preserving record groups; tail starts and class rows stay
  caller supplied. `selection.SelectPaul2013PhoneMarkerRecordPath` composes
  the exact-record route through shortlist finalization and model-backed path
  scoring. `selection.SelectPaul2013PhoneMarkerRecordPathFromTransitionStates`
  derives its per-phone modes from aligned six-byte transition rows. Score
  state and per-edge transition state remain explicit; exact misses still
  require fallback first.
  `selection.SelectPaul2013PhoneMarkerArenaPath` starts from the contiguous
  record arena and includes the shared phone-state reset and packed-context
  exact lookup before those stages. Its `FromTransitionStates` variant also
  derives per-phone modes from aligned six-byte native state rows; both the
  resolved-key and pointer-resolved arena routes are available.
  `selection.ReadCandidateRecords` reads those candidates' indexed signatures
  and feature bytes in order without loading or decoding waveforms.
  `selection.PruneCandidates` implements the observed cumulative-cost cutoff
  and flagged-candidate retention as a standalone step;
  `selection.Backtrack` follows predecessor links to choose one unit per
  context. `selection.ModeTwoCategoricalPenalty` implements one directly
  observed branch of adjacent-context scoring.
  `ModeOtherCategoricalPenalty` ports the other branch's byte-table penalties;
  `ScoreTransition` looks up the masked raw distance and both generated
  feature distances, then combines the categorical branches, DLL-derived
  coefficient rows, context predicates, cumulative cost, and duration term.
  `Paul2013TransitionWeightRow` derives the coefficient row from two observed
  six-byte state records, and `ScoreTransitionFromContextStates` uses that
  producer directly. `Paul2013PhoneRowHasTransitionFlag` ports
  `FUN_10017010`'s byte-`+2` lookup through the exact 256-byte table at RVA
  `0x7ba44`; `ApplyPaul2013TransitionStatePhoneFields` applies the byte-0,
  byte-1 period lookup, byte-2, and byte-5 writes made by `FUN_10024680` when
  supplied the record's resolved pointer target. `ReadPaul2013TransitionPeriodPhoneSpans`
  extracts each row's byte `+0x1d` at the native `0x1e` stride. The later
  query/fallback writes and compact one-row/two-row advancement are composed by
  `BuildPaul2013TransitionContextStates`, with acceptance linked to
  `QueryPaul2013Position` results. `ReadPaul2013TransitionEligibilityShorts`
  reads the native signed-short gate at workspace offset `+0x65c`, using the
  descriptor group ordinal and `0x1e0`-short stride; its input adapter applies
  the positive-short gate to ordered phone rows.
  `BuildPaul2013TransitionContextStatesFromEligibilityArena` composes that
  read with the native phone-field writes and compact query/fallback advance.
  The arena writer and mapping from source text remain unresolved. The caller
  still maps those phone states to candidate
  positions; full text-to-context production is unresolved.
  `BuildPaul2013TransitionScoreInput` maps the current context mode to the
  observed current/previous metric and two feature-code groups in each parsed
  unit record. It fills the signatures and coefficient row too. Duration and
  other context predicates remain explicit because their state producers are
  not recovered.
  It also applies the Stage 20 zero-subtotal shortcut for same-unit mode-2
  edges and marked consecutive edges. When `TransitionScoreInput.Model` is
  set, `voice.Paul2013.GlobalUnitOrdinal` maps bank-local IDs through the
  ordered `dblist.idx` ranges, allowing that check across bank boundaries.
  Without a model, the scorer uses the established same-bank adjacency rule.
  Cross-bank runtime confirmation is from the matched-index Kate captures;
  direct Paul cross-bank parity remains untested.
  `RankPaul2013CandidatesByContinuity` ports the span metadata pass in
  `FUN_100230a0`: it reads signature marker bit `0x80`, matches adjacent
  global IDs against supplied candidate layers, computes right-match count,
  total span, weighted side coverage, and ascending initial key order. Equal
  key order is not guaranteed by the native larger-list sort. It accepts
  the expanded candidate layers and per-position context modes; it does not
  generate those candidates or recover their context modes.
  `Paul2013CandidateMembershipFlags` ports the first oversized-list marker
  pass in `FUN_10023350`: for 31–10,000 continuity-ranked candidates, it sets
  the protected flag when a caller-mapped candidate identity appears in the
  whole-position key list. The unit-to-key mapping and later flag additions
  remain caller work. `ScorePaul2013CandidateLocalPrefix` connects that flag
  pass and the native prefix boundary to candidate-record reads,
  `FUN_100182e0` input construction, and continuity-scaled local costs. It
  returns scored-prefix markers alongside the untouched tail; context values
  and final tail ordering still come from their separate caller stages.
  `EvaluatePaul2013WholePosition` ports the deterministic post-lookup portion
  of `FUN_10024060`: it numerically sorts IDs and removes duplicates as
  `FUN_10024010` does, then applies the 30-class/10,000-unit ranked shortlist,
  member-count metric summation, and strict `>9` acceptance rule (row class 8
  accepts any positive subtotal). Runtime comparison established that native
  class weights equal class populations. Query-variant generation and
  state-dependent class lookup remain caller work.
  `LookupPaul2013WholePosition` now connects ordered, caller-supplied query
  variants to `voice.ClassCatalog`; found records are then sorted and
  deduplicated in native numeric-ID order before shortlist evaluation. It does
  not generate those variants or build the two-row fallback path.
  `RunPaul2013CatalogQuerySequence` connects the query-control-flow port to
  whole-position catalog dispatch and preserves candidate lists across passes.
  `QueryPaul2013WholePosition` executes the ordinary mode-1 pass and the
  model-class-12 masked mode-1/original mode-2 sequence, then applies the
  native post-lookup shortlist. `BuildPaul2013WholePositionRecordInput` reads
  its seven-byte target signature, row-class byte, and model-class byte from
  the observed three-record arena block selected by the six-byte state row;
  `QueryPaul2013WholePositionFromRecordState` connects that projection to the
  query. `QueryPaul2013WholePositionFromTransitionState` accepts the state
  row produced by the `FUN_10024680` port; both functions address the same
  six-byte workspace base. `QueryPaul2013PositionFromTransitionState` also
  composes the conditional fallback with that record-derived base row and
  model class, deriving its signed row-class delta from record `+0x6de`.
  `QueryPaul2013PositionFromStateArenas` reads the prior-metric and live
  context-gate bytes through the state pointer's `+0x4c` target. Producing
  that state arena and eligibility ordering remain upstream inputs.
  `LookupPaul2013WholePositionByPrefix` ports the `FUN_10023dc0` prefix
  broadening loop, 300-class cap, and cumulative 30-member stop for one
  already-produced key, optionally continuing a caller-supplied candidate
  list. `LookupPaul2013QueryDispatch` now joins this path for zero row class
  and routes nonzero row classes through the feature-view helper, including
  the native 290-class early return. `RunPaul2013CatalogQuerySequence` now
  obtains nonzero-row ranges directly from the reconstructed class catalog;
  its optional scope callback is retained for controlled comparisons.
  `LookupPaul2013WholePositionByFeatureViewWithQueryMode` applies the native
  shrinking view prefix to a caller-supplied candidate scope: query mode 1
  relaxes from ten bytes to one, while modes 0 and 2 stop at six. It then
  applies the 10-class/10,000-member ranker and its `>=10` coverage stop. It
  preserves caller scope order on the native fast path and limits scored
  candidates to the first 10,000 matching classes. The older helper remains
  as a compatibility wrapper. A swapped three-bit field in the feature views
  was corrected; the full catalog now matches captured P half-key range
  counts for two queries. Broader parity, native equal-view tie order, and
  row-state/query-mode production remain unresolved.
  `BuildPaul2013FallbackRows` ports the observed
  `FUN_100242a0` row copy: it emits row classes 1 and 2, advances byte +3 by
  the supplied signed model-table delta and row ordinal, and for model class
  12 creates two tree targets (byte +6 masked with `0x80` in mode 1, then the
  original signature in mode 2). `QueryPaul2013FallbackRows` connects both
  generated rows to the catalog query sequence and ranks their deduplicated
  class records with the row's feature view before assembling the native
  ranked-prefix/raw-sidecar layout. The signed model-table delta, prior metric,
  and live context gate remain explicit caller inputs; feature queries use
  the reconstructed catalog index by default. A positive mode-1 query from an
  empty list snapshots the first 30 raw IDs for the sidecar. Ranked IDs retain
  score order; raw IDs absent from them follow in ascending order, with the native
  63-entry capacity enforced. `QueryPaul2013Position` joins whole-position
  acceptance to the conditional fallback call and reports the native row
  advance (one or two). `ExpandPaul2013PositionQueryCandidates` carries an
  accepted shortlist or both fallback ID rows into ordered unit-reference
  pools under a caller-supplied unit cap. It resolves fallback IDs only from
  class records retained by those query results and fails on missing records.
  `RankPaul2013QueriedCandidatePositions` then connects those pools to the
  native continuity metadata pass when one explicit context mode is supplied
  for every expanded position. `RankAndScorePaul2013QueriedCandidatePositions`
  carries the ranked rows through bounded local scoring, including the
  oversized protected-prefix path; score contexts and candidate-key mappings
  remain caller-produced and must align with each ranked row.
  `FinalizePaul2013QueriedCandidateShortlist` applies the secondary tail score,
  native tail sort, 30-unit cap, and protected-candidate replacements when all
  required local costs and caller-derived class values are present.
  `SelectPaul2013QueriedCandidatePath` composes the queried positions through
  shortlist finalization and model-backed path selection when all required
  per-edge state is supplied. `SelectPaul2013QueriedCandidatePathFromTransitionStates`
  derives one mode per accepted query or two per fallback query from the
  aligned six-byte transition states before running that path.
  `Paul2013ContinuityFeatureScale` and `ScorePaul2013CandidateLocalCost` then
  apply the recovered multiplier `1 / (weighted side + max(total span / 2, 1))`
  to the feature term while leaving categorical penalties unchanged. Its
  float32 constants were read from local DLL `.rdata` at `0x1006d170` and
  `0x1006d174`. `Paul2013CandidateLocalScorePrefix` ports the corresponding
  `FUN_10023350` scoring-prefix choice: prefer the leading full-span group;
  otherwise score small lists in full and bound large lists at the first
  strict key increase after 30 entries, retaining boundary ties.
  `ScorePaul2013CandidateLocalPrefix` now runs that boundary against aligned
  caller-produced score contexts and loads only records inside the prefix.
  It also returns the first oversized-list protection marker pass. The
  marker pass is limited to the initial scored prefix;
  `ExpandPaul2013CandidateProtectionPrefix` ports the following large-pool
  scan, including copied candidate fields, destination ordering-key retention,
  and the tenth-protected-candidate ordering boundary. Candidate identity
  keys, whole-position pool size, and aligned score contexts remain explicit
  caller inputs. `ScorePaul2013ExpandedCandidateLocalPrefix` carries the
  copied candidates through record reads and local scoring, using its source
  index map to realign those contexts.
  `Paul2013TailOrderingScore` ports the two recovered category multipliers;
  `SortPaul2013CandidateTailAndCap` sorts only the supplied tail and applies
  the 30-candidate cap. `FinalizePaul2013CandidateShortlist` also restores
  flagged candidates from the sorted remainder into trailing unprotected
  slots, following the native five-protected-candidate limit. Tail start,
  secondary scores, and the complete protected count still come from
  context-conditioned caller state; equal-score ordering may differ from the
  DLL's sort helper.
  `ScoreTransitionLayer` evaluates every predecessor for each supplied
  candidate and records the minimum predecessor index;
  `ScoreAndPruneTransitionLayer` then applies the measured cutoff while
  preserving predecessor indexes.
  `InitializePathLayer` wraps caller-supplied finite local costs in the first
  path layer with no predecessor links; it does not produce those costs.
  `SelectPaul2013MinimumCostPath` connects initialization, all-predecessor
  scoring, per-layer pruning at the native multiplier, and final backtracking
  when callers supply candidate layers, local costs, context gates, and an
  edge scorer. `SelectPaul2013ModelMinimumCostPath` reads indexed records,
  normalizes each raw candidate local score by the observed divisor of two,
  constructs mode-specific transition inputs, scores edges with loaded metric
  and feature tables, prunes cumulative layers, and backtracks the path.
  Candidate layers, pruning gates, and per-edge state records remain caller
  inputs. `SelectPaul2013PathFromLocalScorePasses` adapts continuity/local-score
  results directly into that search while preserving unit order and protected
  flags; it rejects any candidate without a recovered path cost.
  Context/duration producers, fallback state inputs, and production of the
  initial candidate layer remain missing, so text selection is not yet wired.
  `BuildPaul2013UnitScoreInput` maps each candidate's parsed feature row to
  the primary and mode-selected additional feature codes consumed by
  `FUN_100182e0`, and derives its signature-based scaling/equality branches.
  It still requires context bytes, marker values, and the context-mode producer.
  `ScoreUnitCost` applies the feature-distance scaling and category adjustment
  from `FUN_100182e0`. `UnitRecordBytePenalty` implements the observed
  byte-1 through byte-3 rules using the two recovered DLL lookup windows and
  supplied context markers. `UnitByteFivePenalty` applies the two recovered
  byte-5 field tables and rejects unsupported high-field pairs. Cross-byte
  pair penalties are applied when mapped categories fit the recovered 3-by-3
  tables; sentinel/out-of-range categories fail closed. Context-array
  producers and candidate-state integration remain incomplete.
  `KnownUnitCategoricalPenalty` combines the recovered rules currently
  executable from these inputs.
  `selection.UniqueExactSelector` connects exact class lookup to the pipeline
  only when each supplied context has exactly one indexed unit. Missing or
  multi-unit classes fail closed; this does not replace native candidate
  broadening, scoring, or path selection.
- `text`, `selection`, and `synthesis` define separate stage contracts.
- `synthesis.BuildPaul2013TimelinePCMRow` extracts the captured row window:
  row type 1 starts at sample zero, while type 2 starts at the decoded segment
  length minus the row's leading span. It validates the copied range against
  the supplied backing buffer, which may include scratch samples after the
  decoded segment. Text-to-row production remains a separate pipeline step.
- `synthesis.ScalePaul2013TimelinePCMRow` applies the row's integer percentage
  gain with truncation toward zero and the observed signed 16-bit saturation
  limits before joining.
- `synthesis.NormalizePaul2013Controls` resolves negative defaults, maps speed
  zero to 50, and clamps effective pitch, speed, and volume to the observed
  ranges of 50–200, 50–400, and 0–500 before the selected-unit renderers
  consume them. `synthesis.ResolvePaul2013ControlOverrides` also ports
  `FUN_10022850`'s independent per-control override switches and clamps over
  already-selected engine-level values. Pause handling remains outside the
  renderer contract.
- `pipeline` orchestrates those contracts and returns `ErrStageUnavailable`
  when a stage is absent.
- `pipeline.SynthesizeWAV` writes the returned sample values directly into the
  final WAVE allocation, with RIFF-size and host-allocation bounds checked.
- `synthesis.ConcatenatingRenderer` assembles explicit selected-unit PCM in
  sequence with default controls for diagnostics. Its optional
  `ApplyObservedGain` setting applies the normalized 0–500% gain with 16-bit
  saturation; the default remains raw concatenation at volume 200 for unit
  inspection. It does not implement timeline clipping,
  UPM joins, pitch/speed controls, or context reconstruction.
- `synthesis.CursorTimelineRenderer` emits the selected units at the observed
  normal-row cursor positions by omitting each non-final row's trailing UPM
  span. It is a default-control timing experiment: it does not reconstruct or
  blend omitted boundary audio, insert synthetic rows, or claim parity.
- `synthesis.BuildPaul2013UPMSegments` builds the observed five-word records
  between adjacent UPM periods, including doubled sample-grid positions and
  the recovered speed-ratio arithmetic. `SelectPaul2013UPMSegment` and
  `Paul2013ResampledPeriodLength` port the nearest-segment scan and captured
  period-length clamp. `PlanPaul2013UPMSegmentResampling` follows the moving
  segment scan, permits a segment to be selected again while its output budget
  remains, and carries the cumulative output count without touching audio.
  `LimitPaul2013PeriodLength` applies the cumulative
  sample-budget cap from `FUN_1002afb0`; `Paul2013UPMSegmentTargetLength`
  derives that budget from the first segment's `FirstPeriod` field. The plan
  stops when that budget is exhausted, including the captured integer case
  where the next target coordinate would otherwise repeat. The caller supplies
  the cumulative prior output length. `ResamplePaul2013PeriodWithinBudget`
  composes the cap with isolated-period resampling. `ResamplePaul2013Period`
  uses linear interpolation; the interpolation is an
  approximation, and these helpers do not perform context selection or
  timeline joins.
  `synthesis.PeriodResamplingRenderer` consumes selected units, splits their
  decoded PCM at UPM boundaries, and applies that helper to each period. It
  supports pitch-control experiments with default speed; it does not match the
  DLL's neighboring-window reconstruction, context blending, or timeline
  output and makes no synthesis-parity claim.
- `synthesis.SplitPaul2013UPMSamples` exposes each exact decoded-sample span
  as an independent copied window with its sample offset. It requires positive
  UPM periods that cover the PCM exactly. The context-row renderer below uses
  explicit selected rows to prepare these windows for the Stage 8 blend.
- `synthesis.UPMSegmentPlanRenderer` executes the moving segment plan against
  each selected unit's own UPM spans and applies the cumulative length cap.
  At default pitch and speed it preserves the direct samples. For adjusted
  controls, `BlendPaul2013UPMSegmentWindows` applies the captured falling
  window to the current output context and the rising window to the
  tail-aligned selected source period, with the observed crop, zero-pad,
  truncate, and 16-bit-add behavior. It then appends the following source
  period and carries that span as context for the next step. The first context
  window is zero-filled because upstream timeline scratch state is not wired
  in. Text-driven row selection and full timeline integration remain
  incomplete.
- `synthesis.UPMTimelineJoinRenderer` joins complete selected-unit payloads
  with equal or unequal UPM edge spans using the direct 2013 timeline path.
  Complete payloads stand in for native selected windows, so this mode
  supports controlled audio experiments without claiming native timeline
  parity. `UPMEqualSpanJoinRenderer` remains available as a strict variant.
- `synthesis.Paul2013ContextMultipliers` applies the observed left/right gates
  from `FUN_1002d230` to caller-supplied seven-byte rows, mode bytes, and
  neighbor indexes. It reports side eligibility; the row-level context mixer
  below applies those gates to explicit selected rows.
- `synthesis.BuildPaul2013UPMEdgeWeights` builds the observed integer left,
  current, and right ramp arrays from gathered neighbor counts and gate
  multipliers. It applies the observed period-count and five-entry context
  caps. It does not sample or mix the corresponding audio windows.
- `synthesis.BuildPaul2013ContextEdgePlan` composes the context-side gates,
  neighbor-count caps, and UPM edge ramps. The caller supplies gathered
  neighbor counts.
- `synthesis.MixPaul2013ContextUPMTimeline` carries that plan through the
  timeline accumulator for supplied period windows.
- `synthesis.MixPaul2013ContextTimelineRows` prepares those windows from
  explicit selected current/left/right timeline rows. It applies each supplied
  row mode and gain, aligns the capped left prefix and right suffix by UPM
  period, checks row sample counts and doubled edge spans against the decoded
  unit views, then invokes the captured blend path. It skips mixing when both
  native context gates are closed. Context selection from text and production
  of any shared-buffer scratch tail remain unsupported.
- `synthesis.RenderPaul2013ContextTimelineRows` composes that per-row context
  reconstruction with the direct default-pitch timeline join. It requires
  already-produced primary rows, gates, and selected side rows; it does not
  claim text-to-audio selection parity.
- `synthesis.BuildPaul2013ContextTimelineRowInputs` carries typed normal
  timeline rows into the context renderer and preserves explicit scratch PCM
  tails. `BuildPaul2013NormalContextTimelineRowInput` applies the observed
  mode-dependent left `base+1` and right `base-1` row adjustments from
  `FUN_1002d230`, then attaches those rows and their scratch tails. The caller
  still supplies the base indexes, gate rows and modes; their upstream
  production remains unresolved. `RenderPaul2013NormalContextTimelineRows`
  composes explicit row inputs and the renderer. `BuildPaul2013NormalContextTimelineInputs`
  and `RenderPaul2013NormalContextTimelineFromRows` apply the row adjustment
  for each supplied gate before rendering.
- `synthesis.ApplyPaul2013BlendWindow` crops or zero-pads one neighbor window
  at the aligned edge and applies the rising curve for `FUN_1002d010` or the
  falling curve for `FUN_1002cdf0`. Both use the smaller source/output length
  for phase. The engine embeds the 8,192 exact float32 entries from DLL RVA
  `0x6d1b0`; the right helper starts at index 4,096. From the repository root,
  recreate the tracked table from the read-only DLL with
  `python3 tools/revkit/scripts/extract_paul2013_window_table.py binary/vt_pau.dll engine/synthesis/window_table.bin`.
- `synthesis.AddPaul2013WeightedWindows` accumulates already-windowed samples
  on a caller-supplied output timeline using the traced integer multiply,
  divide, clip, and 16-bit accumulator behavior. It does not construct the
  windows, choose their offsets, or derive normalization.
- `synthesis.MixPaul2013UPMInterval` combines current, left, and right source
  windows for one UPM interval using the leading/trailing edge ramps. Source
  window selection and reconstruction remain caller work; it uses the exact
  embedded coefficient table for the window calculation.
- `synthesis.MixPaul2013UPMTimeline` places prepared periods consecutively in
  a unit's output span and validates all periods before accumulation. It does
  not join adjacent selected units or derive source windows.
- `synthesis.MixPaul2013Timeline` joins prepared normal-row windows with equal
  or unequal adjacent edge spans, using the active-span coefficient schedules
  from `FUN_1002aac0` and the exact embedded 2013 coefficient table. The strict
  `MixPaul2013EqualSpanTimeline` wrapper retains its explicit equality check.
  Neither function extracts row windows from unit data.
- `synthesis.RenderPaul2013TimelineRows` connects explicit row metadata to
  type-dependent source-window extraction, per-row gain, and the general join
  path. `RenderPaul2013EqualSpanTimelineRows` preserves a strict equal-span
  variant. Both read selected units and return mono 16 kHz PCM; row production
  from text remains unsupported. `BuildPaul2013TimelinePCMInputs` maps typed
  normal rows into the renderer's row inputs, and
  `RenderPaul2013NormalTimelineRows` carries them through without requiring
  callers to copy the row kind, sample span, edge spans, and gain by hand.
  Captured type-2 rows can read beyond the decoded unit into the native shared
  PCM buffer; explicit renderer inputs can carry that scratch tail separately.
- `synthesis.BuildPaul2013TimelineUnitView` ports the selected descriptor's
  combined, first-side, or second-side DAT/UPM offset, count, sample span, and
  cached edge selection from `FUN_1002c120`; its view-mode producer remains
  caller work. `synthesis.BuildPaul2013NormalTimelineRow` adds the observed
  row kind, primary index, three carried control words, and doubled edge
  spans while representing descriptor pointers as unit references.
  `synthesis.Paul2013TimelineBankFileIndex` maps the four Paul banks to the
  captured file-index order, and
  `BuildPaul2013NormalTimelineRowsFromUnitRefs` resolves selected units and
  builds those rows. It still takes the model-state mode, primary sequence
  index, and carried controls as explicit per-row inputs.
  `synthesis.BuildPaul2013SyntheticTimelineRow` fills the observed boundary
  row kind, duration, opaque index, sentinels, and marker when its lookup and
  state values are supplied.
  `synthesis.BuildPaul2013TimelineRows` composes the combined
  view for each selected unit. `Paul2013TimelineOutputFrames`
  applies the captured non-final-row cursor advance, and
  `Paul2013SyntheticTimelineSamples` evaluates the known boundary-row integer
  expression from its two numeric inputs. `Paul2013TimelineOutputFramesWithSynthetic`
  combines ordered normal and synthetic rows while preserving their separate
  length rules. Synthetic-row source values and all mixed-boundary ordering
  have not yet been derived from text. The current capture set has a
  seven-frame join residual, so these helpers build/validate timing plans but
  do not claim exact rendered frame counts or PCM joins.
- `cmd/vtdecode` extracts one unit from an index and DAT bank into a WAV
  file. `cmd/vtlex` prints normalized tokens and their embedded-dictionary
  pronunciation alternatives, including phone labels and original internal
  symbol bytes. `cmd/vtconcat` writes a new WAVE from an explicitly supplied
  sequence of unit references. Its `-render-mode` flag exposes raw
  concatenation, the cursor-timeline experiment, general and strict equal-span
  UPM joins, isolated-period resampling, and the experimental UPM segment-plan
  renderer.
  The latter modes consume existing renderer implementations; they do not
  select units from text or claim whole-synthesis parity. Pitch, speed, and
  volume controls accept `-1` for their observed defaults. The opt-in
  `-apply-observed-gain` flag enables the measured volume scaling.

For example, render a known unit sequence with the experimental segment plan:

```sh
go run ./cmd/vtconcat -data-root ../data-paul/M16 \
  -units gen:12,gen:13 -render-mode segments -pitch 120 -speed 100 \
  -output /tmp/segments.wav
```

Inspect the current text-to-pronunciation frontend with the local shared
dictionary:

```sh
go run ./cmd/vtlex -dictionary-root ../data-common/dict-eng -text "Hello."
```

To assemble one explicitly chosen alternative per normalized token, pass a
comma-separated zero-based `-choices` list, such as `-choices 0,1,0`. The
command prints the flattened phones and symbols plus each selected token's
phone span; it does not rank alternatives.

## Development

Run the module checks from this directory:

```sh
go test ./...
go vet ./...
```

The DAT tests compare all tracked Stage 2 payload captures with the original
DLL PCM. The text tests validate every local embedded dictionary record and
pronunciation payload when those shared assets are present. The voice test also
checks a local Paul unit against the combined UPM span. None of these checks
claim whole-synthesis parity.

See [the implementation plan](../docs/reimplementation-plan.md) for stage
boundaries, evidence requirements, and unsupported behavior.
