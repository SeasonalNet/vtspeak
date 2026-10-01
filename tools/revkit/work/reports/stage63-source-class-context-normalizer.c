/* Static evidence and offset map for text.ApplyPaul2013ModelSourceClassNormalization.
 *
 * FUN_10007520 @ 0x10007520 calls FUN_10009030 with:
 *   param_1 = model + 0x429a2 + row_index * 0x70 + 4
 *   param_2 = signed-byte model context row + 4
 *   param_3 = parser/source row + 0x23
 *   param_4 = parser/source row + 0x24
 *
 * FUN_1000ea20 places the corresponding model row at model + 0x429a8 +
 * row_index * 0x70. Therefore param_1 is two bytes before that row: param_1+7
 * is the row surface at +5, param_1+0x25 is its phone string at +0x23,
 * param_1+0x6e is the metadata word at +0x6c, and *param_1 is the adjacent
 * flag byte. FUN_10009030 checks source class S/C, the I/M/E source-kind
 * branches, E/A phone types, existing phone strings, the metadata-gated
 * FUN_10002f10 call, and FUN_10009cd0 context-encoding fallback in that order.
 *
 * The Go function ports that branch order over the model arena, with the
 * duration engine supplying the recovered generic normalizer. Unsupported
 * generic calls fail closed when no callback is supplied. Ghidra evidence is
 * in lead3-ax-rowflag-static-2026-09-25.txt and the callers are in
 * stage29-parser-helpers.c. This composition has not had a native runtime
 * row-by-row comparison.
 */
