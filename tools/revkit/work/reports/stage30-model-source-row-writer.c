/* Ghidra pseudocode excerpt for FUN_10044fd0 in vt_pau.dll.
 * Parameter and local names are decompiler-generated.
 */

undefined4 __cdecl FUN_10044fd0(
    short *rows, char *text, undefined4 first, undefined4 second,
    undefined1 type, undefined1 class, uint flag)
{
    /* The decompiler's length loop is strlen(text). */
    if (strlen(text) != 0) {
        if (99 < *rows)
            return 0;

        short *row = rows + *rows * 0x4a + 10;
        *(undefined4 *)row = first;
        *(undefined4 *)(row + 2) = second;
        *(uint *)(row + 4) = strlen(text);
        *(undefined1 *)((int)row + 0x23) = type;
        *(undefined1 *)(row + 0x12) = class;
        *(uint *)(row + 0x14) = flag & 0xff;
        *(undefined1 *)(row + 0x18) = 0xff;
        strcpy((char *)(row + 0x1a), text);
        ++*rows;
        *(undefined4 *)(rows + 2) = *(undefined4 *)(row + 2);
    }
    return 1;
}

/* Byte offsets relative to a row are 0, 4, 8, 0x23, 0x24, 0x28,
 * 0x30, and 0x34. Other fields, including +0x2c, are not written here.
 */

/* FUN_10045070 uses the same row layout as FUN_10044fd0 and additionally
 * copies its third argument to row +0x52 after copying the primary string.
 */
undefined4 __cdecl FUN_10045070(
    short *rows, char *text, char *auxiliary, undefined4 first,
    undefined4 second, undefined1 type, undefined1 class, uint flag)
{
    if (strlen(text) != 0) {
        if (99 < *rows)
            return 0;
        short *row = rows + *rows * 0x4a + 10;
        *(undefined4 *)row = first;
        *(undefined4 *)(row + 2) = second;
        *(uint *)(row + 4) = strlen(text);
        *(undefined1 *)((int)row + 0x23) = type;
        *(undefined1 *)(row + 0x12) = class;
        *(uint *)(row + 0x14) = flag & 0xff;
        *(undefined1 *)(row + 0x18) = 0xff;
        strcpy((char *)(row + 0x1a), text);
        strcpy((char *)(row + 0x29), auxiliary);
        ++*rows;
        *(undefined4 *)(rows + 2) = *(undefined4 *)(row + 2);
    }
    return 1;
}

/* FUN_100451e0 passes an unsplit string directly to FUN_10044fd0. If a
 * space or tab occurs, it walks the bytes, NUL-terminates each component in
 * a 32-byte local buffer, and calls FUN_10044fd0 for each component. A zero
 * result from a component append is returned immediately.
 */

/* FUN_100452c0 passes the third string directly to FUN_10045070 when it has
 * no backslash. Otherwise it splits that string at each backslash, including
 * empty components, copies each component into a 68-byte local buffer, and
 * calls FUN_10045070 with the same primary string and coordinates.
 */
