/* Static evidence for a post-append type-11 write at 0x10045a90.
 * See vt_pau-objdump-disassembly.txt at 0x10045a6f-0x10045a98.
 */

/* After FUN_10044fd0 succeeds, native mode 2 checks the byte at the selected
 * source coordinate. A hyphen assigns type 11 to the last source row. */
if (row_append_succeeded && input_mode == 2 && following_byte == '-') {
    last_source_row->type_2c = 11;
}

/* Append success is tested at 0x10045a6f-0x10045a74; mode is compared at
 * 0x10045a77-0x10045a7c; the byte comparison is at 0x10045a7e-0x10045a85;
 * and the dword write is at 0x10045a90. The mode's meaning remains unknown.
 * This path is static-only in the available evidence.
 */
