/* Bounded control-flow summary from Ghidra pseudocode for
 * FUN_1000d640 @ 0x1000d640. The decompiler output is not original source. */

/* Under param_4 == -1, empty param_5, and zero param_6, the helper skips
 * leading ASCII spaces, then scans the remaining string. Ordinary bytes and
 * the recognized apostrophe suffixes are copied to param_2 in source order.
 * If the scan reaches the source NUL, it terminates param_2 and copies the
 * empty suffix back to the original param_1 buffer. */

/* Semicolon, period, comma, question-mark, and exclamation-mark cases have
 * additional param_3 and FUN_10003a70 behavior. The Go helper does not model
 * those branches. Its supported subset also rejects the penultimate-hyphen
 * case and non-ASCII bytes rather than guessing their table-dependent result. */
