/* Static disassembly evidence for the additional dot write in
 * FUN_100544f0. See vt_pau-objdump-disassembly.txt at 0x1005603c-0x10056088.
 */

/* This branch is reached with a scanner result in eax. Positive status is
 * required; the current character must be '.'; and a source row must exist. */
if (scanner_status > 0 && current_character == '.' && source_row_count > 0) {
    last_source_row->type_2c = 1;
}

/* Address evidence: positive-status branch at 0x1005603c; plus-token route at
 * 0x1005603e-0x10056059; dot comparison and rejection at 0x1005605e-0x10056060;
 * source-row count check at 0x10056066-0x1005606d; dword assignment at
 * 0x10056077. This path is static-only and lacks runtime write-watch evidence.
 */
