/* Static disassembly evidence for FUN_10054050's multi-source-row write.
 * Names and semantics of the surrounding parser state are not inferred.
 * See vt_pau-objdump-disassembly.txt at 0x10054131-0x1005416d.
 */

/* The generic branch calls FUN_100451e0 to append an A row. On successful
 * return it checks the count at the input source-row table; when that count
 * exceeds one, it writes 10 into the just-appended parser source row's +0x2c
 * dword. */
if (generic_row_append_succeeded && input_source_row_count > 1) {
    last_output_source_row->type_2c = 10;
}

/* Address evidence: 0x10054131-0x1005414c prepares and calls the generic
 * FUN_100451e0 append with type byte 0x41; a zero return exits through the
 * failure path. At 0x10054151 the input count is compared against one and
 * counts <= 1 skip the write. The output parser-row count is then read as a
 * signed short and used with the 0x94-byte stride; the dword assignment of
 * 0x0a is at 0x1005415f. This is a static observation, with no matching
 * runtime write-watch capture. The Go port represents append success and
 * input count as explicit gate inputs.
 */
