/* Static search inventory for writes ending at a displacement of -0x54 in
 * vt_pau-objdump-disassembly.txt. A matching displacement alone does not
 * establish a parser source-row +0x2c write: the base pointer and index scale
 * must also be traced. This report separates attributed FUN_100544f0 stores
 * from additional candidates that still require function-boundary review.
 */

/* FUN_100544f0 0x94-stride stores already represented by Go final-write
 * helpers or the previously documented comma/dot helpers:
 *   0x10054694 = 12, 0x100546af = 5
 *   0x10054999 = 1, 0x10054f28 = 1, 0x1005501f = 1
 *   0x10055677 = 1, 0x100557df = 1, 0x10055a4a = 1, 0x10055a8e = 1
 *   0x10056077 = 1
 *
 * 0x1005551d is also inside the function but uses a byte-scale index
 * ([esi+eax-0x54]); stage58-parser-scan-fallback-row-type-one.c traces it to
 * the preceding source row's +0x2c field. Incoming scanner and table
 * predicates for the listed stores remain partly caller supplied and the
 * static-only branches lack runtime watches.
 */

/* Other candidate stores found by the same displacement search:
 *   0x10042e8f = 1 (handled by stage54-contextual-row-type-one.c)
 *   0x100430ce = register (edx=1; handled by stage55-mode-one-prior-key-row-type-one.c)
 *   0x100445b1 = 5, 0x10045a90 = 11
 *   0x1004eb4d = 7
 *   0x10053229 = 8, 0x100538c0 = 8
 *   0x1005782b = register (esi=7), 0x1005788b = register (esi=7)
 *     (handled by stage56-split-append-row-type-seven.c)
 *   0x10059f8b = 5
 *
 * The quote-coordinate type-1 writes at 0x1004f426, 0x1004f6c5,
 * 0x1004f867, 0x1004fa17, 0x1004fbb4, 0x10050073, 0x10050229,
 * 0x10050278, 0x1005084b, 0x10050b4e, and 0x100510df are handled by
 * stage48-parser-quote-coordinate-row-type.c. The 0x1004eb4d branch is
 * handled by stage47-parser-hyphen-fallback-row-type.c.
 * The post-append type-5 write at 0x100445b1 is handled by
 * stage49-parser-input-flag-row-type.c.
 * The exact-input-short type-5 write at 0x10059f8b is handled by
 * stage53-parser-exact-input-flag-row-type-five.c.
 * The post-append type-7 writes at 0x10057722 and 0x100577db are handled by
 * stage50-parser-post-append-row-type-seven.c.
 * The type-8 stores at 0x10053229 and 0x100538c0 are handled by
 * stage51-parser-zero-flag-row-type-eight.c.
 * The type-11 store at 0x10045a90 is handled by
 * stage52-parser-hyphen-mode-row-type-eleven.c.
 *
 * All listed 0x94-stride candidates now have a Go final-write helper, though
 * several incoming predicates remain caller-supplied and the branches remain
 * static-only. The byte-indexed in-function candidate at 0x1005551d is
 * attributed to the preceding source-row field by the stride/offset
 * calculation in stage58; scanner and lookup outputs remain caller supplied.
 *
 * Follow-up instruction-flow review additionally attributes indirect writes
 * at 0x1003d4b6 and 0x1003da7a; see
 * stage57-model-parser-row-type-one-and-reset.c. This demonstrates why the
 * displacement grep alone cannot be treated as exhaustive.
 */
