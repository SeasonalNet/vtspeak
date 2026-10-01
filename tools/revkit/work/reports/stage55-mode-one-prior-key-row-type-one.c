/* Static evidence for the additional type-1 write at 0x100430ce.
 * See vt_pau-objdump-disassembly.txt at 0x1004307a-0x100430d1.
 */

/* The preceding instructions compute a 0x94-byte row address from the signed
 * source-row count and compare the row key at +4 with the current key. The
 * branch then requires a positive count, mode value 1, and two positive
 * caller state values before storing 1 at that row's +0x2c. */
if (prior_row_key == current_key && mode_value == 1 &&
    state_a > 0 && state_b > 0) {
    prior_source_row->type_2c = 1;
}

/* The native stack variables' semantic labels and earlier scanner path are
 * unresolved. This branch is static-only in the available evidence. */
