/* Static offset map for the full-model exception bridge.
 *
 * Sources are Ghidra pseudocode in:
 *   lead3-ax-metadata-helper-static-2026-09-25.txt (FUN_100091b0)
 *   stage29-parser-helpers.c (FUN_10007520 and FUN_1000ea20 call layout)
 *   stage6-exception-dictionary-expanded.txt (FUN_10008dc0)
 *   stage6-dictionary-transform-expanded.txt (FUN_1000ca30/FUN_1000ca50)
 *
 * FUN_10007520 uses a short-array base at model+0x429a2, reads a signed
 * context count there, and advances rows by 0x70. The FUN_1000ea20 context
 * row starts at model+0x429a8, so its parser-row index is +0x00, source
 * surface is +0x05, phone string is +0x23, and X retry class is +0x64.
 * The context row's adjacent state short starts two bytes before that row.
 *
 * For each row, FUN_100091b0 returns a low short derived from the linked
 * parser row's +0x23 kind, +0x24 form, and whether the previous context row
 * links to the same parser row. FUN_10007520 enters FUN_10008dc0 only when
 * that short is zero, parser-row +0x30 is -1 or the model phone string starts
 * with NUL, and parser-row +0x23 is 'U' or parser-row +0x52 starts with NUL.
 * The parser-status/phone-marker and kind/context-form tests are the outer
 * FUN_10007520 gate around both the inner normalization block and the later
 * post-handler rules. The FUN_100091b0 low short only closes the inner
 * exception/normalization block; when the outer tests pass, post-handlers
 * still run even if that short is nonzero. When the outer tests fail, both
 * normalization and post-handler processing are skipped.
 *
 * FUN_10008dc0 scans context rows contiguously from the gate row, derives
 * normalized lookup keys, uses +0x66 as its X-prefix retry marker, and returns
 * the matched phone-code string plus the number of rows included. A match
 * calls FUN_1000ca50 with the model pointer at context-row-start-2. The writer
 * emits bytes at that pointer+0x25 (context-row+0x23), consumes a 'd' boundary
 * when its count reaches FUN_1000ca30(surface), writes NUL at each destination
 * row, and stores the short value 2 at context-row-start-2.
 *
 * The Go bridge in text/model_context_normalization.go and
 * duration/exceptions_model.go ports this field mapping and writer behavior.
 * The FUN_100091b0 row mutations, FUN_10007520 outer iteration/index updates,
 * candidate-sequence producer, later special handlers, and runtime row
 * comparison remain separate work; this report makes no whole-cascade claim.
 */
