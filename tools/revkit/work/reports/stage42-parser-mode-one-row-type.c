/* Static disassembly evidence for FUN_100544f0's mode-one type write.
 * See vt_pau-objdump-disassembly.txt at 0x10055a5a-0x10055a99.
 */

/* This scanner branch requires a positive preceding scan count, mode 1, and
 * at least one source row before assigning 1 to the last row's +0x2c dword. */
if (scanner_path_accepted && prior_scan_count > 0 && mode == 1 &&
    source_row_count > 0) {
    last_source_row->type_2c = 1;
}

/* Address evidence: branch entry and positive scan-count check at
 * 0x10055a5a-0x10055a62; mode comparison at 0x10055a68-0x10055a72; source
 * row count at 0x10055a78-0x10055a7d; dword assignment at 0x10055a8e.
 * Prior scanner gates remain caller inputs, and there is no runtime
 * write-watch capture for this path.
 */
