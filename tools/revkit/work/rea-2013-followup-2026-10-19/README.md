# 2013 engine REA follow-up

Target: `binary/vt_pau.dll`, x86 PE32, SHA-256
`200945bbb2853cc56b3eaa03770eedf26066b20bfc471ef5b0213387035b37ce`.

- `evidence.json` contains the PE resource tree, its 988-byte `RT_VERSION`
  payload read, three principal API call graphs, selected G83 scoring and
  selection value-flow traces, and focused function dossiers.
- `evidence-second-session.json` contains expanded call graphs for playback,
  user-dictionary, metadata, and configuration APIs, selected-type inspection,
  and 2006/2013 function comparisons.

The curated findings and evidence boundaries are in
[`docs/reverse-engineering/rea-2013-followup-extractions-2026-10-19.md`](../../../../docs/reverse-engineering/rea-2013-followup-extractions-2026-10-19.md).
