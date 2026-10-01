/* Static disassembly evidence for the accepted-path type write in
 * FUN_100544f0. See vt_pau-objdump-disassembly.txt at 0x10055a34-0x10055a59.
 */

/* On this scanner branch, a positive source-row count is the final predicate
 * before writing 1 into the last source row's +0x2c dword. Earlier scanner,
 * lookup, and mode gates enter this block and remain caller-supplied. */
if (scanner_path_accepted && source_row_count > 0) {
    last_source_row->type_2c = 1;
}

/* Address evidence: source-row count load at 0x10055a34, positive test and
 * failure branch at 0x10055a37-0x10055a3a, and dword assignment at
 * 0x10055a4a. This branch is static-only; no direct runtime write-watch
 * capture is available.
 */
