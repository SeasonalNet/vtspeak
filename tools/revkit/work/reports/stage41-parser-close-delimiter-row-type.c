/* Static disassembly evidence for FUN_100544f0's closing-delimiter type write.
 * See vt_pau-objdump-disassembly.txt at 0x10055691-0x100557f0.
 */

/* The scanner state must first select status 3 and a ')' or ']'; the native
 * branch checks the following span and performs a row scan. Its accepted
 * final scanner statuses are 1, 2, or 3. If the source-row count is positive,
 * the final row's +0x2c dword receives 1. */
if (scanner_path_accepted && (character == ')' || character == ']') &&
    final_scanner_status >= 1 && final_scanner_status <= 3 &&
    source_row_count > 0) {
    last_source_row->type_2c = 1;
}

/* Address evidence: initial status and delimiter checks at 0x10055691-
 * 0x100556ab; final scan result gates at 0x1005576d-0x100557c8; source-row
 * count check at 0x100557ce-0x100557d3; dword assignment at 0x100557df.
 * The Go helper takes the earlier scanner-path decision as an explicit gate.
 * This path is static-only; no runtime write-watch capture is available.
 */
