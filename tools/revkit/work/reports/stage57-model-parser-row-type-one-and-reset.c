/* Static evidence for two source-row +0x2c mutations in FUN_1003d3d0.
 * See vt_pau-objdump-disassembly.txt at 0x1003d4aa-0x1003d4bc and
 * 0x1003da33-0x1003da88.
 */

/* After FUN_1005a350 returns 1, a nonempty row list is checked. The final
 * row receives type 1 only when its existing +0x2c dword is zero. The write
 * uses the pointer formed at 0x1003d4b0 and occurs at 0x1003d4b6. */
if (scanner_result == 1 && source_row_count > 0 &&
    last_source_row->type_2c == 0) {
    last_source_row->type_2c = 1;
}

/* A separate loop path checks a selected row's type equals 2, two row words
 * at +0x28 equal 0x12, and does not take the branch where both words at +0x22
 * equal 0x1f. It then clears that selected row's +0x2c dword. */
if (selected_row->type_2c == 2 && peer_a_plus_28 == 0x12 &&
    peer_b_plus_28 == 0x12 &&
    (current_plus_22 != 0x1f || peer_plus_22 != 0x1f)) {
    selected_row->type_2c = 0;
}

/* The values compared at +0x28 and +0x22 are retained as raw fields; their
 * meanings, caller scanner inputs, and runtime behavior are unresolved. */
