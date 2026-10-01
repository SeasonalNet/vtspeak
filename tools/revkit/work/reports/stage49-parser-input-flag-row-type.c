/* Static evidence for a post-append type-5 assignment at 0x100445b1.
 * See vt_pau-objdump-disassembly.txt at 0x100445a1-0x100445b8.
 */

/* After one of the row append paths succeeds, the short input flag at
 * [esi+2] is checked. A nonzero value marks the last output source row with
 * type 5; zero leaves the field unchanged. */
if (row_append_succeeded && input_record_flag_at_plus_2 != 0) {
    last_source_row->type_2c = 5;
}

/* Address evidence: short comparison at 0x100445a1-0x100445a6; row-index
 * calculation at 0x100445a8-0x100445ae; assignment at 0x100445b1. The writer
 * is statically visible but has no matching runtime write-watch capture.
 */
