/* Static evidence for a second post-append type-5 assignment.
 * See vt_pau-objdump-disassembly.txt at 0x10059f35-0x10059f92.
 */

/* After a successful A-row append, the native short at input record +2 must
 * equal exactly one before type 5 is written to the final output source row. */
if (row_append_succeeded && input_record_flag_at_plus_2 == 1) {
    last_source_row->type_2c = 5;
}

/* Append success is checked at 0x10059f35-0x10059f3a; the short comparison
 * occurs at 0x10059f7b-0x10059f80; the dword assignment is at 0x10059f8b.
 * This branch is static-only in the available evidence.
 */
