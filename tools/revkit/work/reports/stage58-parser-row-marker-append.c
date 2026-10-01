/* Static evidence for the row-type marker append in FUN_1000d190.
 * See stage29-parser-helpers.c at 0x1000d190 and vt_pau.dll .data. */

/* When the native value at param_1[-9] is below 0x1d and the dword at
 * param_1 is nonzero, the function appends a marker to the current row text:
 *   type 2 -> DAT_10077700 ('.')
 *   type 3 -> DAT_100776f4 ('?')
 *   type 4 -> DAT_100776f8 ('!')
 *   type 5 or 12 -> DAT_100776fc (',')
 *   type 11 -> DAT_1007739c ('-')
 * The source C string is at row +0x34 and the marker strings are read at
 * 0x1000d214-0x1000d246. Their bytes are in the read-only analysis input's
 * .data section at file offsets 0x77700, 0x776f4, 0x776f8, 0x776fc, and
 * 0x7739c respectively.
 */

/* The function uses a 32-byte local string buffer. The Go port returns an
 * error when source plus marker would exceed 31 content bytes, avoiding the
 * native overrun behavior. HeaderValue remains an opaque caller input. */
