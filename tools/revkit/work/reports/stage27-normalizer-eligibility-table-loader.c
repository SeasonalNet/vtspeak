===== 0x100032e0 =====
Function: FUN_100032c0 @ 100032c0

int __cdecl FUN_100032c0(int param_1)

{
  int iVar1;
  undefined4 uVar2;
  undefined2 extraout_var;
  char local_104 [256];
  
  iVar1 = FUN_10003660(param_1);
  if (iVar1 == -1) {
    return -1;
  }
  VT_SetDecimal0Pron_ENG();
  DAT_100fee04 = FUN_100033e0(param_1,0x100772fc,0x4d,0x49);
  if (DAT_100fee04 == (int *)0x0) {
    return 0xffff;
  }
  DAT_100fee08 = FUN_100033e0(param_1,0x100772f0,0x4c,0x49);
  if (DAT_100fee08 == (int *)0x0) {
    return 0xffff;
  }
  VT_SetDecimal0Pron_ENG();
  iVar1 = FUN_10003a10();
  if ((short)iVar1 == -1) {
    return iVar1;
  }
  VT_SetDecimal0Pron_ENG();
  FUN_10063ed6(local_104,(byte *)s__s_poly_tree3_100772a0);
  uVar2 = FUN_10001900((int *)&DAT_100ff100,local_104);
  if ((short)uVar2 == 0) {
    return CONCAT22((short)((uint)uVar2 >> 0x10),0xffff);
  }
  VT_SetDecimal0Pron_ENG();
  return CONCAT22(extraout_var,1);
}



===== 0x10003390 =====
Function: FUN_100032c0 @ 100032c0

int __cdecl FUN_100032c0(int param_1)

{
  int iVar1;
  undefined4 uVar2;
  undefined2 extraout_var;
  char local_104 [256];
  
  iVar1 = FUN_10003660(param_1);
  if (iVar1 == -1) {
    return -1;
  }
  VT_SetDecimal0Pron_ENG();
  DAT_100fee04 = FUN_100033e0(param_1,0x100772fc,0x4d,0x49);
  if (DAT_100fee04 == (int *)0x0) {
    return 0xffff;
  }
  DAT_100fee08 = FUN_100033e0(param_1,0x100772f0,0x4c,0x49);
  if (DAT_100fee08 == (int *)0x0) {
    return 0xffff;
  }
  VT_SetDecimal0Pron_ENG();
  iVar1 = FUN_10003a10();
  if ((short)iVar1 == -1) {
    return iVar1;
  }
  VT_SetDecimal0Pron_ENG();
  FUN_10063ed6(local_104,(byte *)s__s_poly_tree3_100772a0);
  uVar2 = FUN_10001900((int *)&DAT_100ff100,local_104);
  if ((short)uVar2 == 0) {
    return CONCAT22((short)((uint)uVar2 >> 0x10),0xffff);
  }
  VT_SetDecimal0Pron_ENG();
  return CONCAT22(extraout_var,1);
}



