# Recovered intervention recordings

Two Bridget recordings reproduce exactly when the actual text is `Hello world.\n`. Their named historical fixture says `Hello.\n`; the original structural-classids log nevertheless contains the longer phone-context stream (`22 17 2b 30 42 1b 2b 15`) and 26 selected units. The historical shared Stage 5 fixture is the likely cause of this mismatch. The fresh runtime tests establish the required literal text independently of that causal inference.

`verified-recipes.json` contains complete proposed runtime batches, with a read-only input mount overriding `/probe/bridget-hello.txt`. Each generated complete WAVE matched the retained original byte for byte:

- structural-classids: 28,820 bytes, SHA-256 `db849c1ef7cf9dc0af65052cb6afd01e60c20260db99c5cc27070704efb9942c`.
- terminal-marker-duplicate-contextkeys-v2: 27,888 bytes, SHA-256 `c6bccacc755dab6c429578e6a8002566df06b30d9ff6aebfcaaf75c5f07ae33c`.

Both complete RIFF sizes agree with file sizes; no trailing-byte reconstruction is involved. The recordings were synthesized in isolated scratch directories, never read as generation inputs. Scratch results are under `/tmp/intervention-probes-live/bridget-{structural,terminal}-Hello_worlddot/work/stage20/probe.wav`. Historical command provenance: September 27 session `01a0e11d-8cd6-7970-a43a-617207e88167`, zero-based JSONL lines 57825 and 59059.

The key-repack-key5-b-c-copy-attr40 recording also reproduces exactly (19,120 bytes, SHA-256 `b117aefbeb1f65b4d83ca190bc17a797c7bf0af87f4ed31db96f263139ed0696`) from the concurrent Stage 16 mark-position fixture `xx<vtml_mark name="..."/>z`, preserved in `key-mark-position-2.txt` with a final newline. The key experiment launched at 2026-09-27T22:04:45Z; another session launched `run-buffer-ex-guard-size.sh 0 16` at 22:04:49Z, which temporarily installed `stage16/input-ex-records-vtml-mark-position-2.txt` into shared Stage 5. This explains the distinct original duration/selected-unit sequence and is strong evidence of historical shared-fixture interference. Runtime output in `/tmp/intervention-probes-third/key-6ce35496bd20/work/corpus-parity/stage20/probe/adapted-pah0.wav` matches the retained original.

All three supplied recipes passed a second fresh scratch run using the unmodified `prepare`, `docker_command`, and `verify_batch` functions. Generated model copies and Wine prefixes were then discarded; generated audio and logs remain in the cited scratch directories.
