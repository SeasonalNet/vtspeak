#!/usr/bin/env python3
"""Generate small VTML x-cmu controls for the local VoiceText phone codebook."""

from pathlib import Path

OUTPUT = Path("tools/revkit/work/stage10/inputs")
CONSONANTS = (
    "P", "B", "T", "D", "CH", "JH", "K", "G", "F", "V", "TH", "DH",
    "S", "Z", "SH", "ZH", "HH", "M", "N", "NG", "L", "R", "W", "Y",
)
VOWELS = (
    "AA", "AE", "AH", "AO", "AW", "AY", "EH", "ER", "EY", "IH",
    "IY", "OW", "OY", "UH", "UW",
)


def tag(pronunciation: str, surface: str) -> str:
    return (
        f'<vtml_phoneme alphabet="x-cmu" ph="{pronunciation}">'
        f"{surface}</vtml_phoneme>"
    )


def main() -> None:
    OUTPUT.mkdir(parents=True, exist_ok=True)
    for symbol in CONSONANTS:
        name = f"tree-phone-consonant-{symbol.lower()}.txt"
        (OUTPUT / name).write_text(tag(f"{symbol} AH0", "probe") + ".\n")
    for symbol in VOWELS:
        name = f"tree-phone-vowel-{symbol.lower()}.txt"
        (OUTPUT / name).write_text(tag(f"T {symbol}0", "probe") + ".\n")
    for kind, symbols, group_size in (
        ("consonants", CONSONANTS, 8),
        ("vowels", VOWELS, 8),
        ("vowel-stress", tuple(f"{symbol}{stress}" for symbol in VOWELS for stress in range(3)), 8),
    ):
        for start in range(0, len(symbols), group_size):
            group = symbols[start : start + group_size]
            name = f"tree-{kind}-{start // group_size + 1}.txt"
            if kind == "consonants":
                pronunciations = tuple(f"{symbol} AH0" for symbol in group)
            elif kind == "vowels":
                pronunciations = tuple(f"T {symbol}0" for symbol in group)
            else:
                pronunciations = tuple(f"T {symbol}" for symbol in group)
            text = ". ".join(
                tag(pronunciation, f"probe{start + index + 1}")
                for index, pronunciation in enumerate(pronunciations)
            )
            (OUTPUT / name).write_text(text + ".\n")
    for start in range(0, len(CONSONANTS), 8):
        group = CONSONANTS[start : start + 8]
        name = f"cmu-consonants-{start // 8 + 1}.txt"
        text = " ".join(tag(f"{symbol} AH0", f"c{start + index + 1}") for index, symbol in enumerate(group))
        (OUTPUT / name).write_text(text + "\n")
    for start in range(0, len(VOWELS), 8):
        group = VOWELS[start : start + 8]
        name = f"cmu-vowels-{start // 8 + 1}.txt"
        text = " ".join(tag(f"T {symbol}1", f"v{start + index + 1}") for index, symbol in enumerate(group))
        (OUTPUT / name).write_text(text + "\n")
    stressed_vowels = tuple(
        f"{symbol}{stress}" for symbol in VOWELS for stress in range(3)
    )
    for start in range(0, len(stressed_vowels), 8):
        group = stressed_vowels[start : start + 8]
        name = f"cmu-vowel-stress-{start // 8 + 1}.txt"
        text = " ".join(
            tag(f"T {symbol}", f"s{start + index + 1}")
            for index, symbol in enumerate(group)
        )
        (OUTPUT / name).write_text(text + "\n")
    (OUTPUT / "cmu-stress.txt").write_text(
        " ".join(tag(f"T AH{stress}", f"stress{stress}") for stress in range(3)) + "\n"
    )


if __name__ == "__main__":
    main()
