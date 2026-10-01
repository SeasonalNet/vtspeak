/* Static evidence for type-7 writes around two source-row appends.
 * See vt_pau-objdump-disassembly.txt at 0x10057811-0x10057841 and
 * 0x1005784d-0x1005788b.
 */

/* The first branch requires positive dwords at parser-state +0x28c and
 * +0x3cc, then stores 7 at the pending source row's +0x2c before the generic
 * append. It copies strings from parser-state +0x2be and calls FUN_10044fd0. */
if (state_at_0x28c > 0 && state_at_0x3cc > 0) {
    pending_source_row->type_2c = 7;
    append_source_row(...);
}

/* After a successful FUN_100451e0 space/tab split, the function writes 7 to
 * the final appended row's +0x2c at 0x1005788b. The native split helper's
 * result is tested at 0x10057873-0x10057876. */
if (split_append_succeeded) {
    last_source_row->type_2c = 7;
}

/* These sites are static-only; the state dwords' meanings and preceding
 * parser path remain unresolved. */
