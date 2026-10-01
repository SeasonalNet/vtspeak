/* Static disassembly evidence for a FUN_100544f0 lookup-miss type write.
 * This is an address-oriented reading of objdump, not recovered source.
 * See vt_pau-objdump-disassembly.txt at 0x10054eeb-0x10054f30.
 */

/* The local count at [ebp-0x4c] must be less than two. FUN_100560a0 scans
 * 0x49 entries rooted at 0x100780fc; a result greater than -1 transfers to
 * another parser path, while a miss requires an existing source row and
 * writes 1 into that row's +0x2c field. */
if (scanner_path_accepted && local_count < 2 &&
    FindSortedTable(table_0x100780fc, 0x49, lookup_key) <= -1 &&
    source_row_count > 0) {
    last_source_row->type_2c = 1;
}

/* Evidence addresses: local-count compare at 0x10054eeb; table lookup call
 * at 0x10054f06; signed result comparison at 0x10054f0e-0x10054f11; source
 * row count check at 0x10054f17-0x10054f1c; dword write at 0x10054f28.
 * No runtime write-watch capture is available. The Go helper accepts the
 * preceding scanner path, native lookup result, and local count explicitly.
 */
