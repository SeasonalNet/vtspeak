/* Static evidence and field map for the Go FUN_10008dc0 candidate-prefix port.
 *
 * Source: FUN_10008dc0 in stage6-exception-dictionary-expanded.txt. The
 * function scans context rows from its starting index up to the supplied row
 * count. FUN_1000c9c0 normalizes each row surface and returns one plus its
 * literal-hyphen count. The caller accumulates that category and checks the
 * exception dictionary's observed inclusive category range, 1 through 4,
 * before attempting a lookup. A category above four terminates the scan.
 *
 * On the first row, an unnormalizable surface returns without a match. On a
 * later row, it ends the scan while preserving the longest earlier match.
 * In-range surfaces are joined by literal hyphens and looked up as exact
 * compressed keys. If the joined lookup misses and the current context row's
 * +0x66 marker is 'X', the current normalized surface is retried with the
 * literal h' prefix and the accumulated category. A successful match passes
 * its phone-code string and matched row count to FUN_1000ca50.
 *
 * The Go port is text.BuildPaul2013NativeExceptionCandidateRows plus the
 * existing sequence lookup and the duration-engine model-context wrapper.
 * The producer reads model context
 * surfaces at row +0x05 and the X marker at +0x64. It preserves contiguous
 * order, stops before an invalid or over-limit row, and leaves the outer
 * FUN_10007520 entry gate and preceding FUN_100091b0 mutations separate.
 * This is static pseudocode coverage; no row-by-row native runtime comparison
 * has been made.
 */
