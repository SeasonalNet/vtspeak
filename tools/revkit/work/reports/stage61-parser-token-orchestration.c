/* Static basis and implementation boundary for
 * text.BuildPaul2013ModelTokenRowsFromParserRows.
 *
 * FUN_1000d190 @ 0x1000d190 (Ghidra pseudocode excerpt in
 * stage29-parser-helpers.c) reads parser rows at +0x2c with stride 0x94,
 * calls FUN_1000d640 repeatedly until the mutated source string is empty,
 * classifies each transformed string with FUN_1000fd20, resolves it through
 * FUN_10003a70 mode 0, and writes a 0x554-byte token row through
 * FUN_1000d450. The prior token row's parser index and source mode byte select
 * FUN_1000d450's ordinary/final or conditional marker branch. The model token
 * count is capped at 200.
 *
 * The Go implementation composes the recovered string transform, embedded
 * key lookup, phone-payload expansion, and row writer. Dictionary misses are
 * represented by the observed zero-alternative row shape. The caller still
 * supplies initialized parser rows and the outer path-phone gate. Parser-row
 * production, the FUN_10007520 normalization/exception/TPP cascade, and the
 * native special pronunciation selector remain outside this function.
 *
 * This is a static port composition. It has not been compared row-by-row
 * against a runtime capture and does not establish end-to-end text parity.
 */
