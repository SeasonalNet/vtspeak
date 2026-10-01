/* Static composition boundary for
 * text.BuildPaul2013ModelPhoneRowsFromOrdinarySource.
 *
 * The ordinary offset/text producer is documented in
 * docs/reimplementation-plan.md and compared with the captured Stage 20
 * parser rows. The downstream model traversal follows the Ghidra pseudocode
 * for FUN_1000d190 and FUN_1000ea20 in stage29-parser-helpers.c and
 * stage29-model-text-parser.c; their individual layouts and gates are
 * summarized in stage61-parser-token-orchestration.c.
 *
 * This function composes the observed ordinary offsets and +0x34 text with
 * source-length values at +0x08 plus explicit per-row controls for +0x2c,
 * +0x30, and +0x52. The length follows the native FUN_10044fd0/FUN_10045070
 * row writer; the remaining controls are not inferred from source spelling.
 * It accepts a single sentence
 * segment and produces model token/context rows only. The FUN_10007520
 * normalization/exception/TPP cascade, finalizer arena construction, and
 * text-to-audio selection remain separate. No full native runtime comparison
 * has been performed for this composition.
 *
 * The ordinary parser-offset arena convenience builder initializes +0x30 to
 * 0xff, matching the captured untagged/unknown/function POS status for the
 * controlled "record" probe in docs/reverse-engineering/phone-symbol-codebook.md.
 * Other POS statuses remain explicit row controls.
 */
