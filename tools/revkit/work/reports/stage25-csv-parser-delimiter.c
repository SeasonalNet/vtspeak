===== 0x10016670 =====
Function: FUN_10016670 @ 10016670

undefined4 * FUN_10016670(void)

{
  undefined4 *puVar1;
  undefined4 uVar2;
  
  puVar1 = (undefined4 *)FUN_1001d9c0(0x18);
  if (puVar1 != (undefined4 *)0x0) {
    *puVar1 = 0;
    puVar1[1] = 0;
    puVar1[2] = 0;
    puVar1[3] = 100;
    puVar1[4] = 0;
    uVar2 = FUN_1001c440(PTR_DAT_1007b6bc);
    puVar1[5] = uVar2;
  }
  return puVar1;
}



===== 0x10016720 =====
Function: FUN_100166d0 @ 100166d0

undefined4 __cdecl FUN_100166d0(undefined4 *param_1,char *param_2,int param_3)

{
  byte bVar1;
  undefined4 uVar2;
  int iVar3;
  byte *pbVar4;
  byte *pbVar5;
  
  if (param_1 == (undefined4 *)0x0) {
    return 0xffffffff;
  }
  if (param_2 == (char *)0x0) {
    return 0xfffffffe;
  }
  *param_1 = param_2;
  if ((undefined *)param_1[1] != (undefined *)0x0) {
    FUN_1001da30((undefined *)param_1[1]);
  }
  uVar2 = FUN_1001c440(param_2);
  param_1[1] = uVar2;
  if ((undefined *)param_1[2] != (undefined *)0x0) {
    FUN_1001da30((undefined *)param_1[2]);
  }
  iVar3 = FUN_1001d9c0(param_1[3] << 2);
  param_1[2] = iVar3;
  if (iVar3 == 0) {
    return 0xfffffffd;
  }
  iVar3 = 0;
  if (0 < (int)param_1[3]) {
    do {
      iVar3 = iVar3 + 1;
      *(undefined4 *)(param_1[2] + -4 + iVar3 * 4) = 0;
    } while (iVar3 < (int)param_1[3]);
  }
  param_1[4] = 0;
  if (param_3 == 1) {
    pbVar5 = (byte *)*param_1;
  }
  else {
    pbVar5 = (byte *)param_1[1];
  }
  if (0 < (int)param_1[3]) {
    do {
      for (; (((bVar1 = *pbVar5, bVar1 == 9 || (bVar1 == 10)) || (bVar1 == 0xd)) || (bVar1 == 0x20))
          ; pbVar5 = pbVar5 + 1) {
      }
      if (*pbVar5 == 0x22) {
        pbVar5 = pbVar5 + 1;
        pbVar4 = (byte *)FUN_10016810((char *)pbVar5);
      }
      else {
        iVar3 = FUN_10064a30(pbVar5,(byte *)param_1[5]);
        pbVar4 = pbVar5 + iVar3;
      }
      for (; ((bVar1 = *pbVar4, bVar1 == 9 || (bVar1 == 10)) || ((bVar1 == 0xd || (bVar1 == 0x20))))
          ; pbVar4 = pbVar4 + 1) {
      }
      bVar1 = *pbVar4;
      *pbVar4 = 0;
      *(byte **)(param_1[2] + param_1[4] * 4) = pbVar5;
      iVar3 = param_1[4];
      param_1[4] = iVar3 + 1;
      if (bVar1 != 0) {
        pbVar4 = pbVar4 + 1;
      }
      if (*pbVar4 == 0) {
        return 1;
      }
      pbVar5 = pbVar4;
    } while (iVar3 + 1 < (int)param_1[3]);
  }
  return 0xfffffffc;
}



===== 0x10064a30 =====
Function: FUN_10064a30 @ 10064a30

int __cdecl FUN_10064a30(byte *param_1,byte *param_2)

{
  byte bVar1;
  int iVar2;
  byte abStack_28 [32];
  
  abStack_28[0x1c] = 0;
  abStack_28[0x1d] = 0;
  abStack_28[0x1e] = 0;
  abStack_28[0x1f] = 0;
  abStack_28[0x18] = 0;
  abStack_28[0x19] = 0;
  abStack_28[0x1a] = 0;
  abStack_28[0x1b] = 0;
  abStack_28[0x14] = 0;
  abStack_28[0x15] = 0;
  abStack_28[0x16] = 0;
  abStack_28[0x17] = 0;
  abStack_28[0x10] = 0;
  abStack_28[0x11] = 0;
  abStack_28[0x12] = 0;
  abStack_28[0x13] = 0;
  abStack_28[0xc] = 0;
  abStack_28[0xd] = 0;
  abStack_28[0xe] = 0;
  abStack_28[0xf] = 0;
  abStack_28[8] = 0;
  abStack_28[9] = 0;
  abStack_28[10] = 0;
  abStack_28[0xb] = 0;
  abStack_28[4] = 0;
  abStack_28[5] = 0;
  abStack_28[6] = 0;
  abStack_28[7] = 0;
  abStack_28[0] = 0;
  abStack_28[1] = 0;
  abStack_28[2] = 0;
  abStack_28[3] = 0;
  while( true ) {
    bVar1 = *param_2;
    if (bVar1 == 0) break;
    param_2 = param_2 + 1;
    abStack_28[(int)(uint)bVar1 >> 3] = abStack_28[(int)(uint)bVar1 >> 3] | '\x01' << (bVar1 & 7);
  }
  iVar2 = -1;
  do {
    iVar2 = iVar2 + 1;
    bVar1 = *param_1;
    if (bVar1 == 0) {
      return iVar2;
    }
    param_1 = param_1 + 1;
  } while ((abStack_28[(int)(uint)bVar1 >> 3] >> (bVar1 & 7) & 1) == 0);
  return iVar2;
}



