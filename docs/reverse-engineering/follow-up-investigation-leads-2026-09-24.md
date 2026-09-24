# Follow-up investigation leads (2026-09-24)

This note records proposed research directions after the feasibility pass for
the local 2013 M16 Paul package. These are leads, not claims that the work has
started or that the current roadmap stages have changed status. The priority
order favors questions that can extend existing runtime evidence. Keep all
VoiceText binaries and model data read-only, and distinguish observed behavior
from semantic interpretation.

## 1. Label more decision-tree inputs with controlled contrasts

**Question.** Which tree input fields track phone identity, stress, neighboring
phones, part of speech, and phrase or word boundaries? What do those fields
change in the selected tree outputs and final synthesis?

**Starting evidence.** Stage 10 already holds the visible token constant while
varying forced CMU phones and vowel stress. It identifies the one-based phone
identity values in some inputs, several consonant groups that select different
onset trees, and the stress value in selected tree positions. The remaining
tree selectors and the wider prosody consequences are still numeric. See the
[phone codebook](phone-symbol-codebook.md#controlled-tree-input-labels) and
the Stage 10 [tree tracing notes](../../tools/revkit/work/stage10/README.md).

**Bounded probe.** Extend the controlled fixtures one factor at a time: vary
one left or right phone, then a boundary, part-of-speech tag, or pronunciation
alternative while holding the rest fixed. Capture the tree identity, full
input vector, result, and downstream timeline rows. Repeat any apparent rule
with a second phone or word before labeling it.

**Useful result.** A verified mapping from one input field to a controlled
property and a measured downstream effect, with counterexamples or limits
recorded. The binary may expose operational roles without revealing the
vendor's original field names.

## 2. Broaden full-synthesis parity coverage

**Question.** How far do the recovered text, selection, timing, and audio
paths reproduce original WAV output across different text classes and API
controls?

**Starting evidence.** DAT decode output matches the DLL for all 580,474 local
units by byte count and SHA-256. A smaller set of controlled synthesis runs
matches captured sample buffers to WAV data, including default and pitch
controls plus number and abbreviation fixtures. This is not whole-synthesis
corpus parity. See the [corpus parity and Stage 9 report](dat-corpus-parity-and-stage9-2026-09-24.md).

**Bounded probe.** Define a small, versioned golden input set covering ordinary
prose, names, numbers, dates/times, addresses, abbreviations, and short edge
cases. Capture deterministic full WAVs at default settings first, then vary
pitch, speed, and volume separately. Compare whole files and report the first
divergent boundary when bytes differ: token/context rows, selected unit IDs,
timeline, PCM blocks, or WAVE accounting.

**Useful result.** A reproducible pass/fail matrix with a clear coverage claim
and divergence localization. Do not extrapolate a small golden set to every
utterance or every engine/package version.

## 3. Investigate the unknown TPP code meanings

**Question.** What observable token or pronunciation behavior is associated
with TPP code families `A`–`G`, `AX`, and the numeric suffixes?

**Starting evidence.** All 31,550 TPP keys decode and re-encode; the typed
payload grammar, caller paths, and `F`/`G` suffix parsing are known. The code
families have corpus associations, but several suffix and component-bit
semantics remain unnamed. See the TPP discussion in the
[engine and model findings](voice-engine-and-model-formats.md#stage-6-text-and-pronunciation-resources)
and the [open-topics review](open-topics-review-2026-09-24.md).

**Bounded probe.** Select keys whose payloads differ in only one code or
suffix, where possible. Compare upstream token records, the returned
phone/context rows, selected pronunciation, and final surface/audio for
controlled neighboring text. Keep raw codes as the authoritative labels;
only assign a linguistic name when multiple independent observations support
it.

**Useful result.** A tested mapping from a code value to a particular
transformation or branch. If the available corpus has no discriminating
examples, record that limitation instead of guessing from key spelling.

## 4. Expand abbreviation probes by likely context

**Question.** Which of the 427 H/T/C abbreviation keys change with context,
case, punctuation, or neighboring token type?

**Starting evidence.** Every key has one runtime probe, and 40 targeted
contrasts cover selected honorific, company, date, measurement, address, and
repetition cases. The C sweep found 93 outputs different from its neutral
template, but a literal result in that template does not establish that a key
is inactive. See the [Stage 10 context report](abbreviation-context-contrasts-2026-09-24.md).

**Bounded probe.** Use the recorded one-context outputs to prioritize keys
with plausible alternate readings. For each selected key, compare a small
set of meaningful contexts (for example address versus ordinary prose, date
versus month name, or company suffix versus isolated token), plus case and
punctuation controls. Capture ordered token surfaces and phone/context rows.

**Useful result.** A per-key set of verified context contrasts, not a blanket
expansion table. Preserve unresolved keys and contexts explicitly.

## 5. Compare other voices and package versions

**Question.** Which parts of the Paul findings generalize to the other
supplied voices or to other VoiceText package versions?

**Starting evidence.** The four supplied voice DLLs have matching executable
`.text` bytes in the recorded comparison, while small `.data` and voice/model
inputs differ. The deep runtime and format analysis covers Paul M16; the other
voice data and older index variants are not characterized to the same depth.
See the [asset map](../assets.md) and [open-topics review](open-topics-review-2026-09-24.md).

**Bounded probe.** If matching data packages are available, first inventory
file sets, hashes, headers, and tree/index layouts without modifying inputs.
Then run a short shared set of API and text probes against each voice and
compare load behavior, selected structural formats, and output properties.
Treat a matching code section as evidence of shared implementation code, not
proof that all voice behavior or data formats match.

**Useful result.** A compatibility matrix stating which layouts and behaviors
were checked for each package, and which findings remain Paul-specific.

## 6. Probe API behavior beyond the main text-to-WAV path

**Question.** What behavior does a native replacement need to preserve for
other output selectors, buffer/thread modes, and error cases?

**Starting evidence.** The public header and the 2013 Paul text-to-WAV path
are mapped, but runtime behavior across every selector, mode, and failure
path is not. See the API boundary in the [open-topics review](open-topics-review-2026-09-24.md)
and the recovered [engine findings](voice-engine-and-model-formats.md).

**Bounded probe.** First select the compatibility surface required by a
future implementation. For each in-scope API mode, record arguments,
return/error codes, output format, and whether state persists across calls.
Avoid probing unrelated modes without an implementation or research need.

**Useful result.** A tested API contract for the chosen scope, including
explicitly unsupported or untested cases.

## Lower-priority leads and scope limits

- **DAT mode 8:** no mode-8 frame occurs in the local Paul corpus. Runtime
  validation needs a genuine payload from another package or an isolated
  decoder test that does not pretend to establish real-corpus compatibility.
- **Anonymous model feature labels:** more static inspection may refine field
  flow, but names are unlikely to emerge without discriminating controlled
  inputs, an independently labeled dataset, or a vendor schema.
- **Reuse rights:** technical feasibility does not establish rights to use or
  redistribute vendor code, voice data, or recordings. Any implementation
  plan needs its own licensing assumptions.

## Suggested first sequence

1. Continue the controlled tree-input contrasts, since fixtures and a tree
   trace already exist and the result can reduce the largest semantic gap.
2. In parallel at the research-planning level, define the small full-WAV
   golden corpus and its exact comparison boundaries.
3. Use the tree and golden-corpus observations to choose whether TPP semantics
   or abbreviation context breadth offers the next highest-value probe.
4. Broaden to other voices or API modes only when a concrete compatibility
   target calls for it.

This sequence is advisory. It does not start a new roadmap stage, authorize
changes to vendor inputs, or make a compatibility or licensing claim.
