/* Static evidence for two post-append type-8 writes.
 * See vt_pau-objdump-disassembly.txt at 0x10053206-0x10053231 and
 * 0x10053875-0x100538c8.
 */

/* Both paths append an A-state source row. If the input record's short at
 * +2 is zero, the newly appended row receives type 8. */
if (row_append_succeeded && input_record_flag_at_plus_2 == 0) {
    last_source_row->type_2c = 8;
}

/* The first append and input-short check are at 0x10053206-0x10053229; the
 * second append and equivalent check are at 0x10053875-0x100538c0. These
 * branches are static-only in the available evidence. Their preceding token
 * recognition paths remain caller-supplied.
 */
