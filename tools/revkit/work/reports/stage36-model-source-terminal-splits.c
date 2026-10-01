/*
 * Static control-flow note for FUN_1000d640 in
 * stage6-deep-text-functions.txt (PE VA 0x1000d640).
 *
 * The scanner branch is gated by param_4 == -1, empty param_5, and param_6
 * == 0. Native argument param_3 is tested at the final punctuation split:
 * when nonzero, output is terminated before the mark and param_1 is left
 * beginning at that mark. When zero, FUN_10003a70 checks the accumulated
 * output unless the preceding copied byte was '.', or the lookup result has
 * nonzero type at local_d4[0]. A positive decision keeps the punctuation in
 * output and consumes through it; otherwise output is truncated and param_1
 * remains at the punctuation. This yields the caller-visible FinalRow
 * argument in model_source_transform.go.
 *
 * At the same scanner position, terminal '-' truncates output and leaves the
 * remainder at the hyphen. Earlier decimal/comma separators are copied with
 * adjacent bytes only when the character attribute mask has 0x10 on both
 * digits; this branch does not set the preceding-period flag. The portable
 * implementation only applies this rule to the recovered ASCII digit table.
 *
 * Evidence is the cited Ghidra pseudocode and its cross-check at
 * stage6-deep-text-functions.txt; this is static analysis, not an independent
 * runtime trace. Callback semantics and all non-ASCII attribute entries stay
 * caller-supplied or unsupported.
 */
