/* Static map of FUN_1000c710's trailing-apostrophe normalization branch.
 *
 * Source: stage6-deep-text-functions.txt, FUN_1000c710 @ 0x1000c710. It only
 * handles source strings whose last byte is apostrophe and copies them into a
 * 32-byte local buffer. When the final three source bytes compare equal to
 * DAT_100779a4 (the NUL-terminated bytes "in'" at vt_pau.dll file offset
 * 0x779a4), it first replaces the final apostrophe with 'G' if the preceding
 * character has attribute bit 7, otherwise with 'g', and calls FUN_1000cb30.
 * A match rewrites a final output '.' to '-' when present and ORs 9 into the
 * caller's row-state byte.
 *
 * If that lookup misses, the function truncates the local input at the
 * apostrophe and calls FUN_1000cb30 again. A match again ORs 9. Otherwise it
 * calls FUN_10002f10 directly with the caller's row-state byte; a handled
 * result ORs 8. If that also returns zero, FUN_10009cd0 encodes the truncated
 * source and a nonempty output ORs 8 (plus the encoder's own 0x20 update).
 *
 * The Go method composes existing mapped-string comparison, context lookup,
 * generic normalization, and ASCII context encoding behavior. The lookup is
 * caller supplied, with an embedded-dictionary convenience wrapper. The
 * unsupported-resource boundary and lack of row-for-row native runtime
 * comparison remain explicit.
 */
