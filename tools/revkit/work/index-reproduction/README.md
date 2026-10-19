# Index overlay reproduction

The Stage 19 converted indexes and Stage 20 candidate overlays are generated analysis outputs, not vendor inputs. The checked-in generators rebuild these outputs from the local, read-only Kate/Bridget model copies and `binary/vt_pau.dll`. The generator families include `convert_legacy_indexes.py`, `make-index-attribute-variants.py`, `make-index-reordered-variant.py`, `build-bridget-2005-index-adapter.py`, and the Stage 20 repack, attribute-transfer, duplicate, row-swap, and continuity-overlay scripts. Their exact variant parameters and source paths are recorded in the Stage 19/20 README files and scripts.

The accompanying `verification.json` records original SHA-256 and size for all 256 byte-identical outputs. The five `index-adapter-key-repacked-matched-context-attrb-transfer-metric-swap` indexes reproduce with `python3 tools/revkit/work/stage20/swap-matched-context-metric-rows.py all`; all five generated files matched the originals byte-for-byte.

The scratch rerun used the existing generators with outputs redirected outside the repository, preserving vendor inputs. To restore the verified files, rerun the corresponding Stage 19/20 generator commands described in `tools/revkit/work/stage19/README.md`, `tools/revkit/work/stage20/README.md`, and the script argument help. Re-run the byte audit before pruning a future regeneration.
