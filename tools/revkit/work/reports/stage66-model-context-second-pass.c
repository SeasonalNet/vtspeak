/* Static offset map for the final counted-row pass in FUN_100091b0.
 *
 * Source: Ghidra pseudocode at 0x100091b0 in
 * lead3-ax-metadata-helper-static-2026-09-25.txt, also reproduced in
 * stage29-parser-helpers.c. The pass reads the signed context count through
 * short *psVar1 where psVar1 = model + 0x429a2, then scans row indexes 1..n-1.
 *
 * For row i, the current phone-code byte is model + 0x429a2 + i*0x70 + 0x29
 * (context row +0x23). The preceding surface begins at
 * model + 0x429a2 + i*0x70 - 0x65 (previous row +0x05); the following surface
 * begins at model + 0x429a2 + i*0x70 + 0x7b (next row +0x05). The pass checks
 * that the code byte is positive, below 'F', and has a nonzero 16-bit class
 * entry in DAT_10077d94; one of those two surfaces must compare equal to
 * DAT_1007736c ("the") or DAT_1007761c ("de"), respectively; and the code
 * byte must not be 'C'.
 *
 * When the gate passes, the three direct writes target the preceding context
 * row (i-1): phone byte +1 = 0x26, phone byte +2 = 0, and state byte |= 8.
 * This follows from the observed row base at model + 0x429a8 + i*0x70, phone
 * field at +0x23, and adjacent state byte at row-start -2. The native
 * following-surface pointer is read even on the last counted row. The Go port
 * accepts the full model arena and requires the complete bounded C-string
 * field to be present before comparing it; the previous-surface match retains
 * the native short-circuit behavior and does not require that following read.
 *
 * The Go implementation maps the observed class ranges from
 * DAT_10077d94 through the existing paul2013ContextCodeClass helper. This is
 * static pseudocode coverage only; no row-for-row native runtime comparison
 * has been made. Generic first-pass dictionary/TPP branches, remaining
 * FUN_100091b0 mutations, FUN_10007520 iteration/index changes, and later
 * special handlers remain outside this slice.
 */
