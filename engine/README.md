# Go engine

This module is the independent Go implementation. It reads 2013 M16 Paul unit
resources, decodes DAT units, and resolves supported ASCII dictionary words to
phone alternatives. It does not yet select units from those phones or render
text to speech.

## Packages

- `dat` validates the observed 2013 Paul unit-index layout, reads DAT and UPM
  spans, decodes a unit to signed 16-bit PCM, and writes a mono 16 kHz WAVE.
- `voice` opens one matching index/DAT/UPM bank read-only. `Bank.ReadRecord`
  returns indexed metadata and feature bytes without reading waveform data;
  `Bank.ReadUnit` also reads the UPM and DAT spans, decodes PCM, and verifies
  the cached first/shared/last UPM bytes plus the decoded sample count against
  the combined UPM vector.
  `OpenPaul2013` groups the four banks and decision trees in the local Paul
  M16 data root.
- `tree3` parses and evaluates the supported Paul decision trees and loads all
  nine duration and eight pitch trees after checking their observed shapes.
- `duration.OpenPaul2013` loads the shared English dictionary and Paul
  duration trees from their local roots. `Engine.Evaluate` accepts source text,
  resolves words with unique pronunciations, derives the ordinary-path vectors,
  dispatches each phone through the recovered phone-class selector, and
  returns raw duration-tree values for later timing work.
- `distance` validates and reads the little-endian triangular `cepdist.tbl`
  resource, with checked symmetric lookup by metric code and the observed
  packed-code mask. It builds the observed 256-bin triangular feature-distance
  table using the DLL's symmetric difference formula. `voice.OpenPaul2013`
  loads the raw 1,024-entry table and builds the feature table with the banks
  and trees.
- `text.LoadTXT2Tables` reads the nine shared classification and replacement
  tables with validated framing. `text.LoadEmbeddedDictionary` validates all
  indexed records in the shared embedded English lexicon, and
  `text.ParsePhonePayload` parses its direct-ID and alternative-path forms.
  `text.Paul2013PhoneIDCodebook` provides the documented 256-slot compact-ID
  mapping. `PhonePayload.ExpandPronunciations` expands IDs through a five-byte
  codebook into internal symbol bytes, and
  `EmbeddedDictionary.ResolvePaul2013SurfacePronunciation` joins those steps
  for one canonical token. `text.DecodeCMUPhones` labels the runtime-backed
  internal phone bytes and rejects structural or unknown values. The
  frontend now retains the four embedded pronunciation metadata flags on
  lexical tokens and carries them through sequence flattening; the legacy
  token-to-phone path copies these values, though their meanings remain
  unlabeled. They are distinct from the five duration-tree row values.
  `text.SelectLexicalPronunciations` flattens selected alternatives and retains
  each token's phone span and original separators. The engine can select
  automatically when every token has exactly one dictionary pronunciation via
  `LexiconFrontend.ResolveUniquePhoneSequence`; context-based pronunciation
  ranking remains unimplemented, so ambiguous entries fail closed.
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
  sentinels. It assumes the standard terminal `Z` marker; special marker
  transitions and the outer `FUN_10012c70` token-subdivision rules remain to be
  ported. `duration.EvaluatePaul2013SourceText` resolves text with the local
  dictionary, automatically uses unique pronunciations, derives the duration
  inputs, and needs no per-phone metadata from its caller.
  `duration.EvaluatePaul2013Text` dispatches each phone to the
  observed duration-tree family and evaluates its nine-short vector. The result
  is still a raw tree value; its time unit and timeline conversion are unknown.
  The older `text.BuildPaul2013TokenDurationInputs` remains a low-level trace
  helper.
  The legacy path is partly traced:
  `FUN_10012c70` groups phone records and calls `FUN_10013c00`, while
  `FUN_100135d0` derives neighbor sentinels and tree position from token and
  terminal markers, and reads the row and auxiliary values produced above.
  The ordinary path now evaluates duration trees, but tree outputs are not yet
  converted to a synthesis timeline.
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
  month-specific calendar validity is not checked. Other date and time forms,
  general normalization, and phone-context generation remain unsupported.
  `text.LexiconFrontend.ResolveText` tokenizes ASCII surfaces and returns
  dictionary alternatives with original model-coded bytes and CMU labels; it
  resolves expanded number words and rejects unsupported numeric separators;
  The unique-pronunciation path can also build one ordered phone sequence
  without a caller-supplied choice list. It is not yet connected to
  `pipeline.Pipeline` because model-context projection and ambiguous
  pronunciation ranking are still missing.
- `text` implements the observed seven-byte to five-byte class-key transform,
  two ten-byte feature views, and their weighted mismatch tables.
- `voice.BuildPaul2013ClassCatalog` reconstructs sorted exact-key classes and
  their member-unit references from the four read-only unit indexes. Its
  `LookupContext` method applies the key transform to an existing context and
  returns that exact class when present.
  `selection.RankClassRecords` applies the observed count and population
  limits while retaining unit membership. The catalog has not yet been
  checked against a captured DLL candidate list. `selection.ExpandClassRecords`
  flattens ranked memberships in source order under an explicit unit cap; it
  does not construct contexts or candidate scoring metadata. The catalog
  does not score individual units. `selection.ExactContextUnitCandidates`
  maps an existing seven-byte context to its exact class and returns bounded
  unit references; it does not broaden missing matches.
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
  `ScoreTransitionLayer` evaluates every predecessor for each supplied
  candidate and records the minimum predecessor index;
  `ScoreAndPruneTransitionLayer` then applies the measured cutoff while
  preserving predecessor indexes.
  `InitializePathLayer` wraps caller-supplied finite local costs in the first
  path layer with no predecessor links; it does not produce those costs.
  Context/duration producers, weight-row selection, fallback candidate
  generation, per-unit scoring, and initial local-cost production remain
  missing, so text selection is not yet wired.
  `ScoreUnitCost` ports the feature-distance scaling and category adjustment
  from `FUN_100182e0`. `UnitRecordBytePenalty` implements the observed
  byte-1 through byte-3 rules using the two recovered DLL lookup windows and
  supplied context markers. `UnitByteFivePenalty` applies the two recovered
  byte-5 field tables and rejects unsupported high-field pairs. Cross-byte
  pair penalties are applied when mapped categories fit the recovered 3-by-3
  tables; sentinel/out-of-range categories fail closed. Context-array
  producers and candidate-state integration remain incomplete.
  `KnownUnitCategoricalPenalty` combines the recovered rules currently
  executable from these inputs.
- `text`, `selection`, and `synthesis` define separate stage contracts.
- `synthesis.NormalizePaul2013Controls` resolves negative defaults, maps speed
  zero to 50, and clamps effective pitch, speed, and volume to the observed
  ranges of 50–200, 50–400, and 0–500 before the selected-unit renderers
  consume them. Pause handling remains outside the renderer contract.
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
- `synthesis.BuildPaul2013UPMSegments` builds the observed five-word records
  between adjacent UPM periods, including doubled sample-grid positions and
  the recovered speed-ratio arithmetic. `SelectPaul2013UPMSegment` and
  `Paul2013ResampledPeriodLength` port the nearest-segment scan and captured
  period-length clamp. `PlanPaul2013UPMSegmentResampling` follows the moving
  segment scan and carries its cumulative output count, without touching audio.
  `LimitPaul2013PeriodLength` applies the cumulative
  sample-budget cap from `FUN_1002afb0`; `Paul2013UPMSegmentTargetLength`
  derives that budget from the first segment's `FirstPeriod` field. The caller
  supplies the cumulative prior output length. `ResamplePaul2013PeriodWithinBudget`
  composes the cap with isolated-period resampling. `ResamplePaul2013Period`
  uses linear interpolation; the interpolation is an
  approximation, and these helpers do not perform context selection or
  timeline joins.
  `synthesis.PeriodResamplingRenderer` consumes selected units, splits their
  decoded PCM at UPM boundaries, and applies that helper to each period. It
  supports pitch-control experiments with default speed; it does not match the
  DLL's neighboring-window reconstruction, context blending, or timeline
  output and makes no synthesis-parity claim.
- `synthesis.UPMSegmentPlanRenderer` executes the moving segment plan against
  each selected unit's own UPM spans and applies the cumulative length cap.
  It accepts pitch and speed controls, but omits neighbor-window reconstruction,
  and context blending. It copies the selected unit's raw second-period tail
  as a stand-in for reconstructed context audio; its output is experimental.
- `synthesis.Paul2013ContextMultipliers` applies the observed left/right gates
  from `FUN_1002d230` to caller-supplied seven-byte rows, mode bytes, and
  neighbor indexes. It reports side eligibility only; neighbor gathering and
  sample blending remain separate work.
- `synthesis.BuildPaul2013UPMEdgeWeights` builds the observed integer left,
  current, and right ramp arrays from gathered neighbor counts and gate
  multipliers. It applies the observed period-count and five-entry context
  caps. It does not sample or mix the corresponding audio windows.
- `synthesis.ApplyPaul2013BlendWindow` crops or zero-pads one neighbor window
  and applies the rising coefficient curve used by the traced helpers. The
  curve is calculated from an inferred analytic form of the DLL's 4096-entry
  table, so sample rounding may differ.
- `synthesis.AddPaul2013WeightedWindows` accumulates already-windowed samples
  on a caller-supplied output timeline using the traced integer multiply,
  divide, clip, and 16-bit accumulator behavior. It does not construct the
  windows, choose their offsets, or derive normalization.
- `synthesis.MixPaul2013UPMInterval` combines current, left, and right source
  windows for one UPM interval using the leading/trailing edge ramps. Source
  window selection and reconstruction remain caller work, and the analytic
  curve may differ from the DLL table by a sample.
- `synthesis.MixPaul2013UPMTimeline` places prepared periods consecutively in
  a unit's output span and validates all periods before accumulation. It does
  not join adjacent selected units or derive source windows.
- `synthesis.BuildPaul2013TimelineRows` derives each selected unit's sample
  count and doubled first/last UPM edge spans. `Paul2013TimelineOutputFrames`
  applies the captured non-final-row cursor advance, and
  `Paul2013SyntheticTimelineSamples` evaluates the known boundary-row integer
  expression from its two numeric inputs. The current capture set has a
  seven-frame join residual, so these helpers build/validate timing plans but
  do not claim exact rendered frame counts or PCM joins.
- `cmd/vtdecode` extracts one unit from an index and DAT bank into a WAV
  file. `cmd/vtlex` prints normalized tokens and their embedded-dictionary
  pronunciation alternatives, including phone labels and original internal
  symbol bytes. `cmd/vtconcat` writes a new WAVE from an explicitly supplied
  sequence of unit references. Its `-render-mode` flag exposes raw concatenation,
  isolated-period resampling, and the experimental UPM segment-plan renderer.
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
