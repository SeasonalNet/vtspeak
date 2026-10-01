===== 0x1002df50 =====
Function: FUN_1002df50 @ 1002df50

undefined4 __cdecl FUN_1002df50(int *param_1,byte *param_2,int param_3)

{
  int *piVar1;
  int iVar2;
  int iVar3;
  int iVar4;
  int iVar5;
  int iVar6;
  int iVar7;
  byte local_1c [12];
  int local_10;
  int local_c;
  uint local_8;
  
  iVar2 = param_3;
  iVar4 = *param_2 + 0xd;
  local_8 = (*param_2 != 0) + 1;
  param_3 = 0;
  local_c = (int)(char)param_1[2];
  piVar1 = (int *)(iVar2 + iVar4 * 0xc);
  iVar7 = param_3;
  local_10 = *(int *)(iVar2 + 4 + iVar4 * 0xc) + -1;
  while( true ) {
    param_3 = iVar7;
    iVar4 = local_10;
    FUN_10016ef0((byte *)(*(int *)(*piVar1 + param_3 * 4) * 5 + *(int *)(iVar2 + 0x90)),local_1c,
                 local_8);
    iVar5 = FUN_10019350((int)param_2,local_1c,local_c);
    iVar6 = param_3;
    iVar3 = param_3;
    iVar7 = param_3;
    if (iVar5 == 0) break;
    FUN_10016ef0((byte *)(*(int *)(*piVar1 + iVar4 * 4) * 5 + *(int *)(iVar2 + 0x90)),local_1c,
                 local_8);
    iVar5 = FUN_10019350((int)param_2,local_1c,local_c);
    iVar6 = iVar4;
    iVar3 = iVar4;
    iVar7 = iVar4;
    if (iVar5 == 0) break;
    if (iVar4 - param_3 < 2) {
      param_1[1] = 0;
      return (uint)param_1 & 0xffff0000;
    }
    iVar6 = (iVar4 - param_3) / 2 + param_3;
    FUN_10016ef0((byte *)(*(int *)(*piVar1 + iVar6 * 4) * 5 + *(int *)(iVar2 + 0x90)),local_1c,
                 local_8);
    iVar4 = FUN_10019350((int)param_2,local_1c,local_c);
    iVar3 = iVar6;
    iVar7 = iVar6;
    if (iVar4 == 0) break;
    if (iVar4 < 0) {
      iVar7 = param_3;
      local_10 = iVar6;
    }
  }
  do {
    local_10 = iVar3;
    param_3 = iVar6;
    if (param_3 < 1) break;
    FUN_10016ef0((byte *)(*(int *)(*piVar1 + -4 + param_3 * 4) * 5 + *(int *)(iVar2 + 0x90)),
                 local_1c,local_8);
    iVar4 = FUN_10019350((int)param_2,local_1c,local_c);
    if (iVar4 != 0) break;
    iVar6 = param_3 + -1;
    iVar3 = local_10;
  } while( true );
  iVar7 = iVar7 * 4;
  for (; iVar7 = iVar7 + 4, local_10 < piVar1[1] + -1; local_10 = local_10 + 1) {
    FUN_10016ef0((byte *)(*(int *)(*piVar1 + iVar7) * 5 + *(int *)(iVar2 + 0x90)),local_1c,local_8);
    iVar4 = FUN_10019350((int)param_2,local_1c,local_c);
    if (iVar4 != 0) break;
  }
  *param_1 = *piVar1 + param_3 * 4;
  param_1[1] = (local_10 - param_3) + 1;
  return CONCAT22((short)((uint)param_1 >> 0x10),1);
}



===== 0x10016ef0 =====
Function: FUN_10016ef0 @ 10016ef0

uint __cdecl FUN_10016ef0(byte *param_1,undefined1 *param_2,uint param_3)

{
  byte bVar1;
  char cVar2;
  char cVar3;
  byte bVar4;
  byte bVar5;
  
  bVar1 = param_1[3];
  bVar5 = bVar1 & 7;
  bVar4 = bVar1 >> 3 & 7;
  if (((bVar1 & 0x80) != 0x80) || (cVar2 = '\x01', bVar4 < 2)) {
    cVar2 = '\0';
  }
  if (((bVar1 >> 6 & 1) != 1) || (cVar3 = '\x01', bVar5 < 2)) {
    cVar3 = '\0';
  }
  if (param_3 == 1) {
    *param_2 = 0;
    param_2[1] = param_1[4];
    param_2[2] = param_1[1];
    param_2[3] = (&DAT_1007b97c)[*param_1];
    param_2[4] = cVar2;
    param_2[5] = (&DAT_1007b8b4)[*param_1];
    param_2[6] = bVar4 * '\n' + cVar3;
    param_2[7] = *param_1;
    param_2[8] = bVar5;
    param_2[9] = (&DAT_1007b97c)[param_1[2]];
    return 1;
  }
  if (param_3 == 2) {
    *param_2 = 1;
    param_2[1] = param_1[4];
    param_2[2] = param_1[1];
    param_2[3] = (&DAT_1007b97c)[param_1[2]];
    param_2[4] = cVar3;
    param_2[5] = (&DAT_1007b918)[param_1[2]];
    param_2[6] = bVar5 * '\n' + cVar2;
    param_2[7] = param_1[2];
    param_2[8] = bVar4;
    param_2[9] = (&DAT_1007b97c)[*param_1];
    return 1;
  }
  return param_3 & 0xffff0000;
}



===== 0x10019350 =====
Function: FUN_10019350 @ 10019350

undefined4 __cdecl FUN_10019350(int param_1,byte *param_2,int param_3)

{
  int iVar1;
  int iVar2;
  
  iVar2 = 0;
  if (0 < param_3) {
    iVar1 = param_1 - (int)param_2;
    do {
      if (param_2[iVar1] < *param_2) {
        return 0xffffffff;
      }
      if (param_2[iVar1] != *param_2) {
        return 1;
      }
      iVar2 = iVar2 + 1;
      param_2 = param_2 + 1;
    } while (iVar2 < param_3);
  }
  return 0;
}



===== 0x10019570 =====
Function: FUN_10019570 @ 10019570

void __cdecl FUN_10019570(int *param_1,int param_2,uint param_3,int param_4)

{
  byte local_18 [12];
  byte local_c [8];
  
  FUN_10016ea0(param_2,local_c);
  FUN_10016ef0(local_c,local_18,param_3);
  FUN_1002df50(param_1,local_18,param_4);
  return;
}



