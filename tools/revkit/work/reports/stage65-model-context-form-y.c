/* Static offset map for the bounded Y/S portions of FUN_100091b0.
 *
 * Source: lead3-ax-metadata-helper-static-2026-09-25.txt,
 * FUN_100091b0 @ 0x100091b0. Ghidra pseudocode shows the handler first
 * requiring parser-row kind 'U' and rejecting a context row that repeats its
 * predecessor's parser-row index. For parser-row form 'Y', it copies the
 * NUL-terminated parser string at +0x52 to the model phone-string field, ORs
 * bit 2 into the adjacent context state byte, and advances the caller's row
 * index across following context rows that carry the same parser-row index.
 * For form 'S', it encodes the existing model surface into the phone field
 * through FUN_10009cd0 and preserves its row-flag update. Both branches then
 * read an auxiliary parser string and enter the shared suffix/code cascade.
 * The Go WithTail variants accept that auxiliary string explicitly and apply
 * only the separately recovered literal suffix cases; other tail dispatch
 * and mutations remain outside these helpers.
 *
 * FUN_1000ea20 places context row i at model+0x429a8+i*0x70. The copied
 * phone string starts at context-row+0x23. The state byte starts at
 * context-row-2, matching FUN_100091b0's pointer arithmetic. The row index
 * advance compares each adjacent context row's signed parser index at +0x00.
 *
 * The Go implementation applies these evidence-backed prefix writes and
 * optionally composes the known literal suffix cases. It does not claim full
 * FUN_100091b0 behavior or runtime parity.
 */
