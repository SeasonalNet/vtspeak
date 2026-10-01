/* Static evidence for FUN_1000d190's punctuation-retention dictionary probe.
 * The caller invokes FUN_10003a70 with mode 0 on the accumulated row text,
 * then checks whether its output buffer is nonempty.
 */

/* FUN_10003a70 mode 0 applies FUN_1000fd70's embedded-key encoding and calls
 * FUN_10011820 with the caller's embedded dictionary state. Its output buffer
 * is later copied only when byte zero is nonzero. The Go port uses the loaded
 * canonical Paul embedded dictionary for this presence check; it does not
 * parse pronunciation payloads for this branch. */
