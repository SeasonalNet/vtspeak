# Stage 9: extended synthesis boundary captures

The two inputs come from the existing Stage 6 numbers and abbreviations
fixtures. `run.sh` runs the original Paul executable through the isolated
Wine/GDB setup, captures every `FUN_10026870` returned PCM block, and saves
the resulting WAVE files. The script restores the Stage 5 input and output
files after each run.

The GDB script also logs `FUN_1002c220` timeline row counts. Compare each
`*-blocks.pcm` file with the `data` chunk of its matching WAVE file. These
captures test the complete sample-buffer-to-WAVE boundary for two additional
text categories; they do not identify every packed model feature.

The number input produced 108,601 frames from four PCM blocks and timeline
row counts 22 and 78. The abbreviation input produced 116,689 frames from
five PCM blocks and row counts 41, 23, and 17. The checked-in
`work/scripts/verify_stage9.py` compares each concatenated PCM capture with
its WAVE data and confirms both WAVE files match their Stage 6 references.
