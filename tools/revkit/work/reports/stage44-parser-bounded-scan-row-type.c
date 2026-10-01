/* Static disassembly evidence for the last unmodeled type-1 assignment in
 * FUN_100544f0. See vt_pau-objdump-disassembly.txt at 0x1005565d-0x1005567e.
 */

/* The dword at [edi+8] must be nonzero and no greater than eax under a signed
 * comparison. The source-row count must then be positive before the last
 * source row's +0x2c field is assigned 1. */
if (scanner_path_accepted && compared_value != 0 &&
    compared_value <= limit_value && source_row_count > 0) {
    last_source_row->type_2c = 1;
}

/* Address evidence: compared-value load/test at 0x1005565d-0x10055662;
 * signed comparison and greater-than branch at 0x10055664-0x10055666; source
 * row count test at 0x10055668-0x1005566d; dword assignment at 0x10055677.
 * The compared values' semantic roles and earlier scanner path remain
 * unresolved and are explicit Go inputs. There is no runtime write-watch
 * capture for this assignment.
 */
