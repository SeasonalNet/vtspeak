/* Static disassembly evidence for the sign-scanner type-7 write.
 * See vt_pau-objdump-disassembly.txt at 0x1004eb1f-0x1004eb55.
 */

/* After a recognized one-character '+' or '-' token, the '-' branch checks
 * the caller coordinate. A value greater than -1 is stored in a neighboring
 * row field; a value <= -1 instead assigns 7 to the last row's +0x2c dword. */
if (scanner_path_accepted && character == '-' && character_count == 1 &&
    coordinate <= -1 && source_row_count > 0) {
    last_source_row->type_2c = 7;
}

/* Address evidence: sign-byte and count checks at 0x1004eb1f-0x1004eb2a;
 * coordinate comparison at 0x1004eb30-0x1004eb39; alternate coordinate
 * write on the positive branch at 0x1004eb3b-0x1004eb45; type-7 assignment
 * on the nonpositive branch at 0x1004eb47-0x1004eb54. Earlier scanner
 * predicates remain caller inputs. This branch is static-only in the
 * available evidence.
 */
