# MakeInfo DTT reproduction audit

The exact-output audit covers all 391,702 Stage 21 and Stage 16 MakeInfo
`.dtt` files (476,366,397 bytes). Every output was regenerated in isolated
scratch with the local VoiceText Paul runtime and read-only voice inputs, then
compared byte for byte before its tracked copy was removed. No `.dtt` files
remain in `tools/revkit/work`.

`verification.json` records each removed filename, original size, and SHA-256,
along with its generator and comparison method. The source runners remain in
`stage16/`, `stage21/`, and `trace-reproduction/generators/`; their exact inputs
and settings are in those scripts and the corresponding Stage 16/21 notes.
The output batches covered are:

| Family | Files | Original bytes |
| --- | ---: | ---: |
| `mi-ws12` | 18,570 | 22,096,770 |
| `mi-word-space` | 13,824 | 16,104,856 |
| `mi-interword-punct` | 16,224 | 13,850,504 |
| `mi-word-interword` | 3,456 | 4,017,115 |
| `mi-word-wsgrid` | 6,912 | 7,521,028 |
| `mi-word-wsclass12` | 2,166 | 1,897,878 |
| `mi-word-wsclass` | 1,152 | 1,252,629 |
| `mi-wsfull12` | 124,416 | 143,349,348 |
| `mi-letter-a-prefix` | 62,190 | 95,806,599 |
| `mi-letter-triple-a` | 64,896 | 75,492,716 |
| `mi-letter-vary-prefix-fixed-aa` | 43,264 | 62,243,219 |
| `mi-letter-third-sweep` | 384 | 558,000 |
| `mi-letter-fourth-sweep` | 800 | 1,305,938 |
| `mi-letter-four-context` | 128 | 194,106 |
| `mi-letter-pair` | 27,040 | 23,587,431 |
| `invalid-case-map` | 5,408 | 6,592,051 |
| `mi-alphabet-small-matrices` | 572 | 273,739 |
| `mi-pause-silence` | 216 | 145,296 |
| `mi-pause-edge` | 54 | 34,725 |
| `stage16-makeinfo` | 30 | 42,449 |

`invalid-case-map` reproduces the historical malformed middle-case mapping
exactly; it is recorded as an old result, not as correct case-conversion
behavior. All 256 reproducible Stage 19/20 `.idx` overlays, including the five
metric-row swaps, are audited separately in
[`index-reproduction`](../index-reproduction/README.md).
