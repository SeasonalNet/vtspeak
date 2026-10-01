===== 0x1001cbc0 =====
Function: FUN_1001cbc0 @ 1001cbc0

undefined4 __cdecl FUN_1001cbc0(int param_1)

{
  *(undefined4 *)(param_1 + 8) = 0;
  *(undefined4 *)(param_1 + 0xc) = 0;
  *(undefined4 *)(param_1 + 4) = 0;
  return 1;
}



===== 0x1001cbe0 =====
Function: FUN_1001cbe0 @ 1001cbe0

undefined4 __cdecl FUN_1001cbe0(int *param_1,int param_2)

{
  int iVar1;
  
  FUN_1001cc30((int)param_1);
  *param_1 = param_2;
  if (0 < param_2) {
    iVar1 = FUN_1001d9c0(param_2 * 4);
    param_1[2] = iVar1;
    iVar1 = FUN_1001d9c0(param_2 * 4);
    param_1[3] = iVar1;
    return 1;
  }
  return 1;
}



===== 0x1001cc30 =====
Function: FUN_1001cc30 @ 1001cc30

undefined4 __cdecl FUN_1001cc30(int param_1)

{
  if (*(undefined **)(param_1 + 8) != (undefined *)0x0) {
    FUN_1001da30(*(undefined **)(param_1 + 8));
  }
  if (*(undefined **)(param_1 + 0xc) != (undefined *)0x0) {
    FUN_1001da30(*(undefined **)(param_1 + 0xc));
  }
  FUN_1001cbc0(param_1);
  return 1;
}



===== 0x1001cc70 =====
Function: FUN_1001cc70 @ 1001cc70

undefined4 __cdecl FUN_1001cc70(int param_1)

{
  undefined4 *puVar1;
  
  puVar1 = (undefined4 *)FUN_1001d9c0(0x10);
  if (puVar1 != (undefined4 *)0x0) {
    *puVar1 = 0;
    puVar1[2] = 0xffffffff;
    puVar1[3] = 0xffffffff;
    puVar1[1] = 0;
    *(undefined4 **)(param_1 + 0x122448) = puVar1;
    return 1;
  }
  *(undefined4 *)(param_1 + 0x122448) = 0;
  return 1;
}



===== 0x1001d9c0 =====
Function: FUN_1001d9c0 @ 1001d9c0

void __cdecl FUN_1001d9c0(size_t param_1)

{
  void *pvVar1;
  
  if (param_1 == 0) {
    param_1 = 2;
  }
  pvVar1 = _malloc(param_1);
  while (pvVar1 == (void *)0x0) {
    Sleep(10);
    pvVar1 = _malloc(param_1);
  }
  return;
}


===== 0x1001ccc0 =====
Function: FUN_1001ccc0 @ 1001ccc0

This helper allocates the descriptor row array referenced by descriptor
header `context[+0x122448]`. It stores `count` at header +0, requests
`count * 0x210` bytes through `FUN_1001d9c0` (malloc with a 10 ms retry loop),
and stores the returned pointer at header +4. For each row it initializes
three dwords only:

  row +0 = 0
  row +4 = 0
  row +8 = -1

The loop does not clear bytes `+0x0c..+0x20f`, which include the 512-byte
inline field `+0x0c..+0x20b`, kind byte `+0x20c`, and three trailing bytes.
The mark writer `FUN_1002efd0` subsequently sets row +8, writes either the
name string and NUL or only a NUL at +0x0c, and writes the kind byte. No
initialization of the bytes after that NUL is established by these routines.
Therefore raw post-NUL bytes in a returned descriptor are observable but
cannot be assigned API meaning from the recovered writer path; they may be
allocator residue. This is based on Ghidra pseudocode plus the corresponding
disassembly at 0x1001ccc0..0x1001cd22, not original source.


