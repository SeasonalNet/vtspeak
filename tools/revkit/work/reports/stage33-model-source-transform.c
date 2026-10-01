/* Ghidra pseudocode excerpt from FUN_1000d640 @ 0x1000d640.
 * This records only its directly portable copy/split branch.
 */

/* If param_4 is not -1 and the input has more than two bytes, a byte at
 * param_1[strlen(param_1)-2] with no 0xc0 attribute bits causes a two-byte
 * suffix split. The prefix is copied to param_2 and param_1 is mutated to the
 * suffix beginning at strlen-2. */
if (param_4 != -1 && strlen(param_1) > 2 &&
    (DAT_1007e188[(byte)param_1[strlen(param_1)-2]] & 0xc0) == 0) {
    strncpy(param_2, param_1, strlen(param_1)-2);
    param_2[strlen(param_1)-2] = '\0';
    strcpy(param_1, param_1 + strlen(param_1)-2);
    return;
}

/* All paths that reach the common tail copy the remaining param_1 C string
 * into param_2 and then clear param_1. The param_4 == -1, empty-param_5,
 * zero-param_6 path enters the larger contraction scanner before that tail;
 * this excerpt does not claim to implement that scanner. */
strcpy(param_2, param_1);
param_1[0] = '\0';

/* FUN_1000d190 calls this helper with two 32-byte local buffers. The Go port
 * uses owned byte slices and leaves caller capacity enforcement explicit. */
