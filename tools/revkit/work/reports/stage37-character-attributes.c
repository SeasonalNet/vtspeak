/*
 * Direct 256-byte table read from binary/vt_pau.dll at file offset 0x7e188,
 * which maps VA 0x1007e188 in this image. The vendor input SHA-256 was
 * 200945bbb2853cc56b3eaa03770eedf26066b20bfc471ef5b0213387035b37ce.
 * Observed ranges:
 *   00-07: 01; 08-0d: 02; 0e-1f: 01; 20: 04
 *   21-2f: 08; 30-39: 30; 3a-40: 08
 *   41-46: a0; 47-5a: 80; 5b-60: 08
 *   61-66: 60; 67-7a: 40; 7b-7e: 08; 7f: 01
 *   80-82: 08; 83: 40; 84-89: 08; 8a: 80; 8b: 08; 8c: 80
 *   8d-99: 08; 9a: 40; 9b: 08; 9c: 40; 9d-9e: 08; 9f: 80
 *   a0-bf: 08; c0-d6: 80; d7: 08; d8-de: 80; df: 40
 *   e0-e6: 40; e7: 08; e8-f6: 40; f7: 08; f8-ff: 40
 *
 * Signed-char string consumers use negative offsets for bytes 0x80-0xff. The
 * prefix at file offset 0x7e108 (VA 0x1007e108) is 128 zero bytes, so the
 * effective signed-char attribute map has zero for those input bytes even
 * though the raw table above contains nonzero high-byte entries. The Go code
 * keeps both views separate. Consumer masks include 0xc0 for class checks,
 * 0x10 for decimal grouping, and 0x08 for scanner checks.
 *
 * These are opaque table flags, not linguistic labels. The bytes remain
 * proprietary vendor-derived runtime data; this report does not change their
 * licensing status. No voice or pronunciation data is included.
 */
