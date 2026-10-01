/* Ghidra pseudocode excerpts for the observed parser-row +0x2c writers.
 * Names and types are decompiler-generated; these are not recovered source.
 */

/* FUN_10051a00: when the scanner's status is 8 and the prior-row predicate
 * succeeds, the last row's +0x2c dword is written only when it is zero.
 */
if (param_2[7] == 8 && previous_row_predicate_is_zero) {
    if ((char)param_2[8] == '?' && last_row_type == 0) {
        last_row_type = 3;
    } else if ((char)param_2[8] == '!' && last_row_type == 0) {
        last_row_type = 4;
    } else if (last_row_type == 0) {
        last_row_type = 2;
    }
    return param_2[5];
}

/* FUN_10051a00 also reaches the same assignments when FUN_10051cc0 returns
 * 0x50. Its scanner/token relationship checks are left as explicit caller
 * inputs by the Go helper.
 */

/* DAT_10078ea0 is the C string "too" at VA 0x10078ea0 (file offset
 * 0x78ea0 in .data). FUN_100544f0 also compares against "either" at
 * VA 0x1007902c.
 *
 * FUN_100544f0: a comma token with status 3 and one-character length scans
 * the following token. If it advances and its status is not 9, the last row
 * receives 5 when the following token differs from both "too" and "either"
 * under FUN_1001c2c0's mapped comparison, or the scanned
 * row count exceeds one; otherwise it receives 12.
 */
if (param_2[7] == 3 && (char)param_2[8] == ',' && param_2[4] == 1) {
    if (last_row_exists && following_token_advance != 0 && following_status != 9) {
        if ((first_compare != 0 && second_compare != 0) || following_row_count > 1)
            last_row_type = 5;
        else
            last_row_type = 12;
    }
}

/* FUN_100544f0's separate dot branch writes 1 when status is 3 or 6, its
 * character count is at least three, and at least one source row exists.
 */
if ((param_2[7] == 3 || param_2[7] == 6) && param_2[4] >= 3 &&
    (char)param_2[8] == '.' && *param_1 > 0) {
    last_row_type = 1;
}

/* FUN_100544f0's double-quote path has distinct scanner predecessors. This
 * is the common final gate represented by the Go helper; the caller supplies
 * whether its preceding scanner path was accepted. */
if (scanner_path_accepted && (char)param_2[8] == '"' && local_count < 2 &&
    FindSortedTable(quote_table, 0x49, lookup_key) > -1 &&
    *source_row_count > 0) {
    last_row_type = 1;
}

/* The boundary-write at 0x10055677 is detailed in
 * stage44-parser-bounded-scan-row-type.c. The accepted-path write at
 * 0x10055a4a is detailed in stage43-parser-accepted-path-row-type.c.
 * The lookup-miss
 * write at 0x10054f28 and closing-delimiter write at 0x100557df are detailed
 * in stage40-parser-unmatched-row-type.c and
 * stage41-parser-close-delimiter-row-type.c. The mode-one write at
 * 0x10055a8e is detailed in stage42-parser-mode-one-row-type.c. The
 * single-quote path beginning at 0x100549ab appends a source row but does not
 * directly assign +0x2c. */

/* Direct runtime watchpoints on the first four rows observed writes 0->2 at
 * 0x10051bc3 for "Hello hello.", 0->3 at 0x10051b7f for "Hello?", 0->4 at
 * 0x10051ba3 for the row before "!", and 0->5 at 0x100546b7 for the row
 * before a comma. FUN_1003e240 later writes zero to the observed fields.
 * Double-quote writes are static-only; see
 * stage38-parser-quote-row-type.c. This excerpt does not recover the full
 * +0x2c producer set.
 */
