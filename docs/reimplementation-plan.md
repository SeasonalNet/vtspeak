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
   span fields, the opaque 21-byte feature row, its 7-byte signature, and the
   combined UPM span. `engine/voice` opens matching index/DAT/UPM files;
   `Bank.ReadRecord` reads metadata without touching waveform data, while
   `Bank.ReadUnit` returns decoded units after checking cached UPM edge bytes
   and the sample-count/UPM-sum invariant established across the Paul corpus.
   `engine/tree3` parses, evaluates,
   and loads the 17
   observed Paul duration/pitch trees. `engine/distance` parses the raw
   triangular `cepdist.tbl`; three documented runtime lookups match its
   symmetric indexing. `distance.SymmetricFeatureDifference` implements the
   observed feature-distance formula and zero-denominator behavior.
   `voice.OpenPaul2013` groups all four banks, the tree catalog, and the
   validated 1,024-entry metric table plus the generated 256-bin feature
   distance table from a local Paul M16 data root.
   `text.LoadTXT2Tables` reads the nine mapped shared `.txt2` tables.
   `text.LoadEmbeddedDictionary` validates the indexed embedded lexicon, and
   `text.ParsePhonePayload` parses its direct-ID and alternative-path forms.
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
   nine inputs for the ordinary single-utterance marker path. The outer
   `FUN_10012c70` token subdivision and special marker transitions remain
   pending. `duration.OpenPaul2013` loads the shared embedded dictionary and
   Paul duration trees from local roots. `Engine.Evaluate` resolves uniquely
   pronounced text and derives duration inputs without caller-supplied
   per-phone metadata.
   `duration.EvaluatePaul2013Text` ports the observed phone-class selector
   tables and evaluates each assembled vector against its duration tree.
   Results are raw 16-bit tree outputs; their physical unit and timeline
   conversion remain unknown. `ApplyObservedPaul2013OnsetIdentityFeature`
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
   Call-specific tree-vector construction, general text normalization, and
   phone-context generation remain unimplemented.
   Compare the resource parsers and transforms with runtime traces.
3. **Text to selected units — lexical lookup and class catalog started;
   context generation and unit selection pending.**
   `engine/text.LexiconFrontend` tokenizes ASCII surface runs, expands the
   observed plain/signed integer and decimal forms plus inferred numeric
   ordinal spellings from 0 through 31, resolves the embedded
   dictionary, and returns CMU-labeled alternatives with their original
   internal symbol bytes. It retains the four embedded-payload metadata flags
   on each token and through phone-sequence flattening; their semantics remain
   unlabeled, and they are not the five duration-tree row values. It preserves
   separators and rejects unsupported numeric separators. Ordinal
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
   `voice.BuildPaul2013ClassCatalog` reconstructs sorted exact-key classes and
   unit membership from the local four-bank indexes, and its exact lookup
   consumes an already-built context; `text` and `selection`
   implement the observed key transform, feature views, weighted mismatch
   calculation, and class-count/population limit behavior. The catalog has
   not yet been checked against a captured DLL candidate list.
   `selection.PruneCandidates` implements the measured post-transition
   cutoff as an isolated operation, including its context gate, ten-candidate
   minimum, and node-flag bypass. `selection.Backtrack` follows the saved
   minimum-cost predecessor chain from the final context.
   `selection.ExpandClassRecords` now flattens retained class members in
   ranked order under an explicit unit cap. It does not generate
   context-conditioned classes, node flags, duration data, or scoring inputs.
   `selection.ExactContextUnitCandidates` connects a supplied seven-byte
   context to its exact class and returns bounded unit references; it does not
   implement the legacy context-broadening path.
   `selection.ReadCandidateRecords` loads their indexed unit records without
   reading waveform payloads, preserving candidate order.
   `selection.ModeTwoCategoricalPenalty` implements the mode-2 branch's
   observed same-unit and signature-byte-2 costs.
   `selection.ModeOtherCategoricalPenalty` implements the DLL's other
   signature-table branch. `selection.ScoreTransition` now looks up raw
   and generated feature distances and combines the recovered category rules,
   coefficient row, mode divisor, context predicates, cumulative cost, and
   duration term. `selection.ScoreUnitCost` ports the feature-distance scaling
   branches of `FUN_100182e0` when given resolved feature pairs and its
   categorical penalty. `selection.UnitRecordBytePenalty` ports the direct
   byte-1 through byte-3 mismatch branches using the recovered DLL maps and
   supplied state markers. `selection.UnitByteFivePenalty` applies the
   recovered high- and low-field cost tables, with unsupported high-field
   values rejected; `selection.KnownUnitCategoricalPenalty` combines these
   recovered byte costs with in-range cross-byte table penalties. It rejects
   mapped categories outside the recovered 3-by-3 tables. Context-array
   producers, duration-term producers, weight-row selection,
   fallback candidate generation, per-unit candidate scoring, and
   initial local-cost production remain unimplemented.
   `selection.ScoreTransitionLayer` now evaluates all predecessor candidates
   for each supplied current candidate and records the winning index for
   backtracking. `selection.ScoreAndPruneTransitionLayer` then applies the
   measured cutoff while preserving predecessor references.
   `selection.InitializePathLayer` wraps caller-supplied local costs for the
   first layer; local-cost production remains unimplemented. Implement
   normalization, phone/context generation, context-conditioned class
   expansion, and remaining unit and adjacent-transition behavior. Then
   construct the initial candidate layer and compare intermediate records and
   selected IDs with controlled captures.
4. **Selected units to speech — diagnostic concatenation and isolated period
   resampling implemented; legacy rendering pending.**
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
   minimum/maximum clamp. `PlanPaul2013UPMSegmentResampling` follows the
   moving segment scan and carries the cumulative output count without reading
   audio. `LimitPaul2013PeriodLength` applies the cumulative
   context-budget cap from `FUN_1002afb0`, and
   `Paul2013UPMSegmentTargetLength` derives the budget from the first segment's
   `FirstPeriod` field. The caller supplies the running prior-output count;
   `ResamplePaul2013PeriodWithinBudget` composes the cap with the isolated
   resampler. `NormalizePaul2013Controls` resolves the observed negative
   defaults, maps speed zero to 50, and clamps effective pitch, speed, and
   volume to 50–200, 50–400, and 0–500. The raw-unit, isolated-period, and
   segment-plan renderers use that result; when gain is enabled they apply
   volume/100 with the observed 16-bit saturation. Pause handling remains
   unsupported.
   Linear interpolation remains an approximation. The
   opt-in `PeriodResamplingRenderer` splits decoded units at UPM boundaries
   and applies the isolated interpolation with default speed; it is an
   experimental renderer and makes no DLL-parity claim. Separately,
   `UPMSegmentPlanRenderer` uses the moving segment plan over each selected
   unit's own UPM spans and accepts pitch and speed controls, but omits
   neighbor-window reconstruction and context blending. It copies raw
   second-period tails from the selected unit as a stand-in for reconstructed
   context audio. It is experimental and does not claim synthesis parity.
   `Paul2013ContextMultipliers` ports the observed left/right context
   gates over caller-supplied model rows and neighbor indexes. The gate does
   not gather neighbors. `BuildPaul2013UPMEdgeWeights` constructs the observed
   integer edge ramps from gathered counts after applying the observed
   period-count and five-entry caps, together with gate multipliers.
   `ApplyPaul2013BlendWindow` crops or zero-pads one supplied window and
   applies the rising curve used by the traced helpers; its analytic form is
   inferred from sampled values in the DLL's 4096-entry table, so final sample
   rounding can differ. `AddPaul2013WeightedWindows` applies the traced
   integer multiply/divide, per-contribution clipping, and 16-bit accumulation
   to caller-prepared windows; callers still supply their offsets and
   normalization. `MixPaul2013UPMInterval` connects the window preparation,
   edge weights, and accumulation for one interval when callers supply the
   selected source windows. `MixPaul2013UPMTimeline` places prepared periods
   consecutively within one unit and validates them before accumulation.
   `BuildPaul2013TimelineRows` derives normal row sample counts and the first/
   last UPM edge spans from selected units. `Paul2013TimelineOutputFrames`
   applies the measured non-final-row cursor advance, and
   `Paul2013SyntheticTimelineSamples` ports the captured sentence-boundary
   integer expression when its two numeric inputs are available. The exact
   2013 join still differs from this cursor prediction by seven frames in one
   stable capture, so this is duration-to-timeline planning, not exact PCM
   assembly.
   Model-window reconstruction and selection, integration of the segment plan
   with audio reads, inter-unit joins, output clipping, and prosody are not
   implemented.
   Port those behaviors and compare PCM/WAVE output with the Stage 8 captures
   before claiming synthesis parity.
5. **Additional packages and interfaces.** Add another voice or public API
   mode only after its format and original-engine behavior have separate
   evidence and acceptance criteria.

## Current use

From `engine/`, run `go run ./cmd/vtdecode -index INDEX -dat BANK -unit 0
-output /tmp/unit.wav`, using matching, local 2013 Paul files. This extracts
one unit; it does not turn text into speech. `go test ./...` checks the decoder
against tracked Stage 2 DLL PCM captures. `voice.OpenBank` reads the unit's
index, DAT, and UPM data. The embedded lexicon parser validates every local
record and its observed pronunciation payload grammar. Key encoding, its small
Paul 2013 tables, and compact phone-ID expansion are provisioned.
`text.LexiconFrontend.ResolveText` resolves ASCII dictionary surfaces into
phone alternatives while preserving the original internal symbol bytes and
CMU labels. It expands observed plain and signed integers, grouped comma
integers, and decimal forms, including the four-digit year rule. The captured
`3:45 PM` clock example is also expanded, with only the captured
`H:MM` formatting range supported by inference. Numeric ordinal suffixes from
zero through thirty-one are expanded; those spellings are inferred except for
the captured second/fifth word outputs. The captured `01/02/2024` date is
expanded as `January second twenty twenty four`; broader MM/DD/YYYY component
bounds are inferred, with no month-specific day validation. Currency and
percentage forms use separate normalizers; only `$5.00` and `25%` have direct
output captures.
The frontend does not yet produce model contexts or audio.
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
or per-unit scoring. These paths have not yet been compared with a captured DLL
candidate list. `selection.PruneCandidates` also implements the measured
filtering rule, but has no transition-score producer to consume yet.
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
provides `LexiconFrontend.ResolveUniquePhoneSequence`, which resolves and
selects the only pronunciation when every token has exactly one alternative;
it rejects ambiguous tokens because the legacy context-based pronunciation
ranker is not implemented. `text.BuildPaul2013TokenPhoneNeighborhoods`
derives phone neighbors across the full utterance and groups features back by
token. It also derives group count and the observed single/non-final/final
position state per phone from the `FUN_10013f30` groups, including the no-vowel
fallback observed in `FUN_10013c00`. It also ports the duration-row fields and
per-phone auxiliary labels from those groups. The new
`text.BuildPaul2013TokenDurationInputsFromText` joins them with utterance-wide
phone neighbors and ordinary boundary sentinels without caller metadata. This
path assumes the standard terminal `Z` marker. The outer `FUN_10012c70`
token-subdivision rules and special state-marker transitions remain unported.
`duration.EvaluatePaul2013Text` selects one of the nine duration trees using the
DLL's phone-class selector tables and evaluates each vector; outputs are not
yet converted to timing. Unit selection remains unwired from text.
The package currently does not claim parity for mode 8, other voice
generations, or whole-synthesis output.
