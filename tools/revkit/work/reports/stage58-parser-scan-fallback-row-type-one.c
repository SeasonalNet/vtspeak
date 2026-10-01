/* Static disassembly evidence for FUN_100544f0's byte-indexed type-1 store.
 * Source: vt_pau-objdump-disassembly.txt, 0x10055348-0x1005552d.
 *
 * At 0x1005536e-0x10055374, ESI is formed as EDI*0x94. The preceding-row
 * comparison reads [param_1 + ESI - 0x7c]. The final store at 0x1005551d is
 * [param_1 + ESI - 0x54] = 1. Since source rows begin at +0x14 and have a
 * 0x94-byte stride, EDI is the one-based row index and the target is that
 * row's predecessor at +0x2c:
 *   +0x14 + (EDI-1)*0x94 + 0x2c == EDI*0x94 - 0x54.
 *
 * The final path also checks a positive source-row count and key match,
 * scanner state 1 or state 2 with FUN_10061940 acceptance, parameter zero
 * below 2, positive scan length, parser mode 1..3, candidate count below 2,
 * a zero local flag, and a positive preceding scan result. The Go helper
 * exposes these scanner and lookup outputs in
 * Paul2013ParserScanFallbackRowTypeGate; it does not claim to produce them.
 * This is static evidence only; no runtime write watch is recorded here.
 */
