/* Static evidence for repeated one-character quote fallback writes in the
 * token-scanner functions. Addresses below are in
 * vt_pau-objdump-disassembly.txt; surrounding function names remain stripped.
 */

/* Each branch checks a double quote and a native count of one, then compares
 * a caller coordinate against -1. When the signed comparison is <= -1 the
 * last source row's +0x2c dword receives 1; positive coordinates are stored
 * in the neighboring coordinate field instead. */
if (scanner_path_accepted && character == '"' && character_count == 1 &&
    coordinate <= -1 && source_row_count > 0) {
    last_source_row->type_2c = 1;
}

/* Matching gates and writes are visible at 0x1004f426, 0x1004f6c5,
 * 0x1004f867, 0x1004fa17, 0x1004fbb4, 0x10050073, 0x10050229,
 * 0x10050278, 0x1005084b, 0x10050b4e, and 0x100510df. Each has a distinct
 * preceding scanner path; the helper accepts that decision and the final
 * character/count/coordinate values explicitly. Static evidence only; no
 * runtime write-watch comparison is available for these paths.
 */
