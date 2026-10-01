/* Static evidence for two post-append type-7 assignments.
 * See vt_pau-objdump-disassembly.txt at 0x10057716-0x10057729 and
 * 0x100577c8-0x100577e3.
 */

/* In the first path a successful FUN_100451e0 A-row append is followed
 * directly by a type-7 write to the newly appended source row. */
if (first_a_row_append_succeeded) {
    last_source_row->type_2c = 7; /* 0x10057722 */
}

/* In the repeated-row path, append success is followed by an input-index
 * increment and comparison with the input row count. A type-7 write occurs
 * only when the loop has another input row. */
if (repeated_a_row_append_succeeded && more_input_rows_remain) {
    last_source_row->type_2c = 7; /* 0x100577db */
}

/* Append call success is checked at 0x1005770a-0x10057715 and
 * 0x100577bf-0x100577c6. The second path's increment/bounds check is at
 * 0x100577c8-0x100577d0; dword stores are at 0x10057722 and 0x100577db.
 * This is static-only evidence without matching runtime write watches.
 */
