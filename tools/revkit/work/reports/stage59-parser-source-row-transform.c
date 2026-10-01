/* Static evidence for the row-to-string slice of FUN_1000d190.
 * See stage29-parser-helpers.c at 0x1000d190 and
 * stage32/33/35/36 source-transform reports.
 */

/* For each 0x94-byte row, the native loop selects the raw row type at +0x2c,
 * reads the source C string at +0x34, and reads associated text at +0x52.
 * A marker is appended when the row text length at +0x08 is below 0x1d and
 * the row type is nonzero and mapped. FUN_1000d640 then receives the mode byte
 * at +0x30, the associated string, whether the row class at +0x24 is 'S',
 * and whether this is the final row in the input list. FUN_1000fd20 classifies
 * the transformed output before token lookup and row writing continue.
 */

/* This port composes only those observed fields/operations. The non-final
 * punctuation lookup through FUN_10003a70 remains a supplied callback, and
 * the following dictionary lookup, token-index association, token-row write,
 * and caller arena contract remain separate stages. */
