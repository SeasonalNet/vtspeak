/* Static evidence for a context-gated type-1 assignment.
 * See vt_pau-objdump-disassembly.txt at 0x10042e50-0x10042e96.
 */

/* The branch compares a prior row key against the current key, checks a
 * caller state value equals 1, checks two values in [1, 2], and accepts a
 * mode value of 1 or 2. It then stores dword 1 at the current row's +0x2c.
 * Other earlier function predicates are intentionally represented as an
 * explicit gate by ApplyPaul2013ParserContextualRowTypeOne. */
if (prior_row_key == current_key && predicate_value == 1 &&
    scan_value >= 1 && scan_value <= 2 &&
    (mode_value == 1 || mode_value == 2) &&
    secondary_value >= 1 && secondary_value <= 2) {
    current_source_row->type_2c = 1;
}

/* The comparisons are at 0x10042e54-0x10042e8d and the write is at
 * 0x10042e8f. The containing routine's earlier scanner path and the semantic
 * meaning of its stack values are unresolved. Static-only; no runtime
 * write-watch is available for this branch. */
