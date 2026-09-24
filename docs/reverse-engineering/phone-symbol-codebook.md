# Compact phone-symbol ID table

`FUN_10003c50` expands dictionary phone IDs through a 256-entry table initialized by `FUN_1000ebc0` at DLL address `0x100fe900`. Each ID has a five-byte slot; the slot is a NUL-terminated sequence of up to four internal symbol bytes. The table below was transcribed from the initializer; `—` means no nonzero byte was assigned to that slot. The 256 IDs can represent sequences; they are not themselves single phonemes.

## Runtime-backed labels for internal symbol bytes

The [VoiceText VTML guide](https://static.carahsoft.com/concrete/files/1615/2520/8261/Voice-Text_Markup_Language.pdf)
defines an English `x-cmu` phoneme input and its symbolic phonetic table. That
guide alone does not identify this DLL's private byte values. Stage 10 sent
each listed consonant followed by `AH0`, and each of the 15 listed vowels at
stress levels `0`, `1`, and `2` preceded by `T`, through the original 2013 Paul
engine. At `FUN_10007520` entry, the upstream token record's pronunciation
field (`0x94`-byte stride, `+0x52`) contained the two expected internal bytes
in input order. The [probe generator](../../tools/revkit/scripts/generate_cmu_probes.py),
[Wine captures](../../tools/revkit/work/stage10/README.md), and
[verifier](../../tools/revkit/work/scripts/verify_cmu_codebook.py) cover all
69 labels below. They form a one-to-one mapping onto every byte `0x01`–`0x45`,
with no gaps or collisions. The guide's `tomato` example independently yields
`T AH0 M EY1 T OW0` → `39 07 2c 1e 39 2f`; a three-tag control and separate
`AH0`/`AH1`/`AH2` contrast agree with the mapping.

For each vowel row, the three byte values correspond to stress `0`, `1`, and
`2` in that order:

| CMU vowel | Stress 0 | Stress 1 | Stress 2 | CMU vowel | Stress 0 | Stress 1 | Stress 2 |
| --- | --- | --- | --- | --- | --- | --- | --- |
| AA | `01` | `02` | `03` | AE | `04` | `05` | `06` |
| AH | `07` | `08` | `09` | AO | `0a` | `0b` | `0c` |
| AW | `0d` | `0e` | `0f` | AY | `10` | `11` | `12` |
| EH | `17` | `18` | `19` | ER | `1a` | `1b` | `1c` |
| EY | `1d` | `1e` | `1f` | IH | `23` | `24` | `25` |
| IY | `26` | `27` | `28` | OW | `2f` | `30` | `31` |
| OY | `32` | `33` | `34` | UH | `3b` | `3c` | `3d` |
| UW | `3e` | `3f` | `40` | | | | |

| CMU consonant | Byte | CMU consonant | Byte | CMU consonant | Byte | CMU consonant | Byte |
| --- | --- | --- | --- | --- | --- | --- | --- |
| B | `13` | CH | `14` | D | `15` | DH | `16` |
| F | `20` | G | `21` | HH | `22` | JH | `29` |
| K | `2a` | L | `2b` | M | `2c` | N | `2d` |
| NG | `2e` | P | `35` | R | `36` | S | `37` |
| SH | `38` | T | `39` | TH | `3a` | V | `41` |
| W | `42` | Y | `43` | Z | `44` | ZH | `45` |

These are supported **CMU input labels for byte values**, not an IPA
transcription of every dictionary entry or proof that every context produces
the same acoustic realization. `0x00` and the special `0x63`/`0x64` bytes
retain their separate control/record roles. The 256-entry ID table below can
now be read as sequences of the labeled bytes.

## Controlled part-of-speech field contrast

The same [VTML guide](https://static.carahsoft.com/concrete/files/1615/2520/8261/Voice-Text_Markup_Language.pdf)
documents `<vtml_partofsp>` values. Stage 10 applied each documented value to
the same isolated word, `record`. Its upstream `0x94`-byte token field at
`+0x30` changed as follows before `FUN_10007520`; the phone bytes below are
from the returned `0x70`-byte row at `+0x25` and decoded with the mapping
above:

| Tag | Upstream `+0x30` | Returned phone sequence |
| --- | --- | --- |
| None, `unknown`, `function` | `ff` | R IH0 K AO1 R D |
| `noun` | `13` | R EH1 K ER0 D |
| `modifier` | `0e` | R EH1 K ER0 D |
| `verb` | `25` | R IH0 K AO1 R D |
| `interjection` | `24` | R IH0 K AO1 R D |

This identifies `+0x30` as a **tag-responsive lexical-class constraint** in
this path. It does not establish a global enum for every upstream status
value: `function` and `unknown` are indistinguishable on this word. All seven
runs had the same four trailing dictionary metadata words, `0, 1, 0, 0`.
Noun/modifier and verb/interjection therefore provide pairs with identical
phone sequences but different status bytes. Their captured duration/pitch
tree input sequences also match exactly within each pair, so this status
distinction does not appear as a changed tree feature in those runs. These
contrasts leave the four metadata-word meanings and other packed model
features open.

The full metadata trace across the 40 linguistic/context inputs and nine
part-of-speech controls covers 227 output rows and ten distinct four-word
patterns. Examples include `Apartment` as `0,1,1,1`, the unexpanded second
`Apt` as `0,0,0,1`, `Hello` as `0,1,0,0`, and the tested `January` rows as
`1,0,1,1`. The seven isolated part-of-speech variants of `record` keep the
same `0,1,0,0` tuple. These are copied Boolean dictionary payload fields from
bits 6, 4, 5, and 7; the captures do not justify semantic names for the
individual bits.

## Controlled tree input labels

Thirty-nine one-token Stage 10 runs forced each of the 24 CMU consonants
followed by `AH0`, and each of the 15 CMU vowels after `T`, onto the same
visible token `probe`. Five additional controls compared `P AH0`, `B AH0`,
`T AH0`, `P AH1`, and `P AH2`. Each single-phone run reached the same four
ordered duration/pitch calls, though the first selected tree varies with the
consonant. The [probe generator](../../tools/revkit/scripts/generate_cmu_probes.py),
[Wine captures](../../tools/revkit/work/stage10/README.md), and
[analyzer](../../tools/revkit/work/scripts/analyze_stage10_tree_features.py)
preserve and check the raw inputs.

The [tree trace](../../tools/revkit/work/stage10/trace-pos-trees.gdb) prints
16 signed halfwords per call; positions below are zero-based within that
printed vector (the last two calls begin printing at the input pointer plus
four bytes). The analyzer checks all 39 symbol vectors. It finds that the
identity values equal the one-based alphabetical rank of each CMU base
symbol. For example AA=`1`, CH=`8`, IH=`17`, P=`27`, and ZH=`39`. The CMU
phone codes are separate (`P=0x35`, `B=0x13`, `T=0x39`).

| Tree input | Directly observed response | First-tree groups for consonants |
| --- | --- | --- |
| Selected onset scalar, printed position 0 | CMU identity ordinal of the controlled consonant | `0x14bf434`: P, B, T, D, K, G |
| Same onset scalar family | Same consonant groups select these trees | `0x14bf46c`: CH, JH; `0x14bf450`: F, V, TH, DH, S, Z, SH, ZH, HH; `0x14bf488`: M, N, NG, L; `0x14bf4a4`: R, W, Y |
| Paired scalar (`0x14bf3c4` or `0x14bf3e0`) | Printed positions 0/1 carry the controlled second/first phone ordinals; position 3 changes with AH stress | P/B/T ordinal `27`/`7`/`31`; AH0/AH1/AH2 stress `0`/`1`/`2` |
| `0x14bf4f8`, third scalar | Printed position 2 follows AH stress | `0`/`1`/`2` |
| `0x14bf568`, vector | Printed position 2 follows AH stress; position 11 matches the preceding scalar result | stress `0`/`1`/`2`; position 11 `0`/`4`/`0` |

These controls name the base-phone identity value and show that consonant
tree selection tracks distinct consonant groups. The groups align with
consonant manner categories, with L grouped with M/N/NG; that family label is
an interpretation of the grouping, not a name recovered from the binary.
The last vector's position 11 is not an independent stress encoding. Other
tree inputs, family roles, and global enum semantics remain open.

| ID | Bytes | ID | Bytes | ID | Bytes | ID | Bytes |
| ---: | --- | ---: | --- | ---: | --- | ---: | --- |
| `0x00` | `—` | `0x01` | `07 13 07 2b` | `0x02` | `15 07 2d` | `0x03` | `2d 23 37` |
| `0x04` | `02 2b` | `0x05` | `13 02` | `0x06` | `13 26` | `0x07` | `22 02 36` |
| `0x08` | `2a 0b 36` | `0x09` | `2c 44` | `0x0a` | `07 35` | `0x0b` | `36 2f` |
| `0x0c` | `37 11` | `0x0d` | `1f 39` | `0x0e` | `2a 39` | `0x0f` | `2c 26` |
| `0x10` | `36 3f` | `0x11` | `37 08` | `0x12` | `2c 07 2a` | `0x13` | `02 35` |
| `0x14` | `05 2b` | `0x15` | `15 08` | `0x16` | `18 39` | `0x17` | `20 0b 36` |
| `0x18` | `22 11` | `0x19` | `25 44 07 2c` | `0x1a` | `29 26` | `0x1b` | `2c 02 36` |
| `0x1c` | `2c 08` | `0x1d` | `2c 18` | `0x1e` | `35 05` | `0x1f` | `43 3b` |
| `0x20` | `07 20` | `0x21` | `0b 2b` | `0x22` | `13 08` | `0x23` | `13 1b` |
| `0x24` | `15 24` | `0x25` | `2b 08` | `0x26` | `2b 30` | `0x27` | `2c 13` |
| `0x28` | `42 11` | `0x29` | `2a 07 2d` | `0x2a` | `07 2d 39` | `0x2b` | `23 37 39` |
| `0x2c` | `2b 23 2e` | `0x2d` | `37 39 1a` | `0x2e` | `07 15` | `0x2f` | `09 2d` |
| `0x30` | `0a 36` | `0x31` | `13 05` | `0x32` | `15 44` | `0x33` | `18 37 39` |
| `0x34` | `20 1a` | `0x35` | `2a 08` | `0x36` | `2b 11` | `0x37` | `2b 1e` |
| `0x38` | `2b 3f` | `0x39` | `36 08` | `0x3a` | `36 11` | `0x3b` | `36 1e` |
| `0x3c` | `37 1a` | `0x3d` | `37 2a 26` | `0x3e` | `39 23 2e` | `0x3f` | `03 36` |
| `0x40` | `0c 36` | `0x41` | `15 1a` | `0x42` | `15 26` | `0x43` | `21 1a` |
| `0x44` | `22 05` | `0x45` | `23 44` | `0x46` | `24 36` | `0x47` | `27 2d` |
| `0x48` | `2a 26` | `0x49` | `2a 2f` | `0x4a` | `2a 30` | `0x4b` | `2b 02` |
| `0x4c` | `2b 27` | `0x4d` | `2b 44` | `0x4e` | `2c 24` | `0x4f` | `2e 21` |
| `0x50` | `30 2b` | `0x51` | `36 02` | `0x52` | `36 28` | `0x53` | `36 30` |
| `0x54` | `37 24` | `0x55` | `37 26` | `0x56` | `37 27` | `0x57` | `39 2f` |
| `0x58` | `42 24` | `0x59` | `44 07 2c` | `0x5a` | `13 07 2b` | `0x5b` | `37 07 2d` |
| `0x5c` | `39 07 2d` | `0x5d` | `39 23 15` | `0x5e` | `01 36` | `0x5f` | `05 2a` |
| `0x60` | `12 44` | `0x61` | `18 2a` | `0x62` | `1a 15` | `0x63` | `1e 38 07 2d` |
| `0x64` | `20 2b` | `0x65` | `2a 05` | `0x66` | `2a 1a` | `0x67` | `2a 42` |
| `0x68` | `2b 05` | `0x69` | `2b 18` | `0x6a` | `2b 1a` | `0x6b` | `2b 24` |
| `0x6c` | `2c 05` | `0x6d` | `2c 1a` | `0x6e` | `2d 1a` | `0x6f` | `2e 2a` |
| `0x70` | `35 1a` | `0x71` | `36 27` | `0x72` | `41 1a` | `0x73` | `43 3f` |
| `0x74` | `2c 07 2d` | `0x75` | `38 07 2d` | `0x76` | `13 1a` | `0x77` | `2c 35` |
| `0x78` | `02 2d` | `0x79` | `05 2d` | `0x7a` | `07 13` | `0x7b` | `07 2a` |
| `0x7c` | `07 2b` | `0x7d` | `07 2c` | `0x7e` | `07 2d` | `0x7f` | `07 37` |
| `0x80` | `07 39` | `0x81` | `0b 36` | `0x82` | `13 07` | `0x83` | `13 36` |
| `0x84` | `15 07` | `0x85` | `15 23` | `0x86` | `18 2b` | `0x87` | `18 2d` |
| `0x88` | `18 36` | `0x89` | `1a 44` | `0x8a` | `21 36` | `0x8b` | `23 15` |
| `0x8c` | `23 2a` | `0x8d` | `23 2d` | `0x8e` | `23 2e` | `0x8f` | `23 37` |
| `0x90` | `23 39` | `0x91` | `24 2b` | `0x92` | `24 2d` | `0x93` | `26 07` |
| `0x94` | `2a 02` | `0x95` | `2a 07` | `0x96` | `2a 2b` | `0x97` | `2a 36` |
| `0x98` | `2a 37` | `0x99` | `2b 07` | `0x9a` | `2b 15` | `0x9b` | `2b 23` |
| `0x9c` | `2b 26` | `0x9d` | `2c 07` | `0x9e` | `2d 07` | `0x9f` | `2d 15` |
| `0xa0` | `2d 23` | `0xa1` | `2d 26` | `0xa2` | `2d 37` | `0xa3` | `2d 39` |
| `0xa4` | `2d 44` | `0xa5` | `35 36` | `0xa6` | `36 05` | `0xa7` | `36 07` |
| `0xa8` | `36 18` | `0xa9` | `36 23` | `0xaa` | `36 24` | `0xab` | `36 26` |
| `0xac` | `37 07` | `0xad` | `37 23` | `0xae` | `37 2a` | `0xaf` | `37 35` |
| `0xb0` | `37 39` | `0xb1` | `38 07` | `0xb2` | `39 07` | `0xb3` | `39 1a` |
| `0xb4` | `39 23` | `0xb5` | `39 26` | `0xb6` | `39 36` | `0xb7` | `39 37` |
| `0xb8` | `01` | `0xb9` | `02` | `0xba` | `03` | `0xbb` | `04` |
| `0xbc` | `05` | `0xbd` | `06` | `0xbe` | `07` | `0xbf` | `08` |
| `0xc0` | `09` | `0xc1` | `0a` | `0xc2` | `0b` | `0xc3` | `0c` |
| `0xc4` | `0d` | `0xc5` | `0e` | `0xc6` | `0f` | `0xc7` | `10` |
| `0xc8` | `11` | `0xc9` | `12` | `0xca` | `13` | `0xcb` | `14` |
| `0xcc` | `15` | `0xcd` | `16` | `0xce` | `17` | `0xcf` | `18` |
| `0xd0` | `19` | `0xd1` | `1a` | `0xd2` | `1b` | `0xd3` | `1c` |
| `0xd4` | `1d` | `0xd5` | `1e` | `0xd6` | `1f` | `0xd7` | `20` |
| `0xd8` | `21` | `0xd9` | `22` | `0xda` | `23` | `0xdb` | `24` |
| `0xdc` | `25` | `0xdd` | `26` | `0xde` | `27` | `0xdf` | `28` |
| `0xe0` | `29` | `0xe1` | `2a` | `0xe2` | `2b` | `0xe3` | `2c` |
| `0xe4` | `2d` | `0xe5` | `2e` | `0xe6` | `2f` | `0xe7` | `30` |
| `0xe8` | `31` | `0xe9` | `32` | `0xea` | `33` | `0xeb` | `34` |
| `0xec` | `35` | `0xed` | `36` | `0xee` | `37` | `0xef` | `38` |
| `0xf0` | `39` | `0xf1` | `3a` | `0xf2` | `3b` | `0xf3` | `3c` |
| `0xf4` | `3d` | `0xf5` | `3e` | `0xf6` | `3f` | `0xf7` | `40` |
| `0xf8` | `41` | `0xf9` | `42` | `0xfa` | `43` | `0xfb` | `44` |
| `0xfc` | `45` | `0xfd` | `64` | `0xfe` | `63` | `0xff` | `—` |

The `0x00` slot is zero and is also the direct-ID list terminator. No nonzero assignment was found for `0xff`. Payloads in the repeated-alternative form use `0xff` as their alternative terminator before the ID sequence is expanded. Code bytes including `0x63` and `0x64` are retained literally in slots; downstream code also uses `d` (`0x64`) as a phone/context component separator.
