/* Static disassembly evidence for one FUN_100544f0 source-row type path.
 * This is an address-oriented reading of objdump, not recovered source.
 * See vt_pau-objdump-disassembly.txt at 0x10054930-0x100549a1.
 */

/* A preceding scanner branch reaches 0x10054930. The local count at
 * [ebp-0x4c] must be below two. FUN_100560a0 searches 0x49 entries rooted at
 * 0x1007844c; result <= -1 returns to the scanner loop. */
if (local_count < 2 &&
    FindSortedTable(table_0x1007844c, 0x49, lookup_key) > -1 &&
    source_row_count > 0) {
    last_source_row->type_2c = 1;
    return scanner_advance;
}

/* Observed instructions: cmp [ebp-0x4c],2 at 0x10054930; table call at
 * 0x1005497b and result test at 0x10054983; source-row count check at
 * 0x10054988-0x1005498d; dword write at 0x10054999. The destination is the
 * final source row's +0x2c field using the 0x94-byte row stride.
 *
 * A second quote branch at 0x10054fd9 uses the same character and count
 * conditions, searches the same 0x49-entry table, and writes 1 at
 * 0x1005501f. Its caller-visible scanner path differs and is not collapsed
 * into the preceding branch by this report. The Go helper accepts the
 * preceding scanner acceptance and native lookup result as explicit inputs.
 * No runtime write-watch evidence for either quote branch is available.
 *
 * The other statically visible writes at 0x10055677, 0x10055a4a, and
 * 0x10055a8e are detailed in stage44-parser-bounded-scan-row-type.c,
 * stage43-parser-accepted-path-row-type.c, and
 * stage42-parser-mode-one-row-type.c.
 * The lookup-miss write at 0x10054f28 and closing-delimiter write at
 * 0x100557df are detailed in stage40-parser-unmatched-row-type.c and
 * stage41-parser-close-delimiter-row-type.c. The single-quote path begins at
 * 0x100549ab.
 * Code 5 at 0x100546af and code 12 at 0x10054694 are covered separately by
 * stage31-model-parser-row-type-producers.c.
 */
