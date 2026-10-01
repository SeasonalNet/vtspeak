/* Ghidra pseudocode excerpt for FUN_1000fd20, called from FUN_1000d190.
 * The names and types are decompiler-generated, not recovered source.
 */

/* Each byte uses DAT_1007e188. Bytes without either 0xc0 attribute bit do
 * not affect the result. A first 0x40-only character yields 1. A 0x80
 * character yields 2 when no prior class was seen, or 3 after one; subsequent
 * 0x80 characters keep the result at 3. */
short state = 0;
char result = 0;
for (char *cursor = source; cursor != NULL && *cursor != '\0'; cursor++) {
    byte attributes = DAT_1007e188[(byte)*cursor];
    if ((attributes & 0xc0) != 0) {
        if ((attributes & 0x80) == 0) {
            if (result == 0) {
                state = 1;
                result = 1;
            }
        } else {
            bool hadClass = state != 0;
            state = hadClass + 2;
            result = hadClass + 2;
        }
    }
}
return result;

/* FUN_1000d190 passes this byte to FUN_1000d450. Its input is the transformed
 * local string after FUN_1000d640; that transform is a separate stage. */
