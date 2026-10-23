# 2013 engine REA decompilation

Target: `binary/vt_pau.dll` (x86 PE), SHA-256
`200945bbb2853cc56b3eaa03770eedf26066b20bfc471ef5b0213387035b37ce`.

- `evidence.json` is the complete exported REA evidence bundle: overview,
  procedure inventory, and all batch-decompilation results.
- `vt_pau-all-procedures.c` contains the 1,098 unique successful Ghidra
  pseudocode bodies, sorted by entry address. These are not original source.
- `function-inventory.tsv` records each unique procedure address, label,
  status, and pseudocode length.
- `batch-evidence.tsv` maps the decompilation batches to retained Evidence IDs.
- `extract_pseudocode.py` reproduces the C-like listing from `evidence.json`.

The complete evidence bundle records one 10-function pilot batch in addition
to the 11 batches covering the full 1,098-procedure set. Those pilot procedures
are deduplicated in the listing and inventory.
