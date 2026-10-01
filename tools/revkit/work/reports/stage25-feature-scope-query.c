===== 0x10023f90 =====
Function: FUN_10023f90 @ 10023f90

uint __cdecl FUN_10023f90(int param_1,uint param_2,int param_3,int param_4,short param_5)

{
  int *piVar1;
  uint uVar2;
  
  piVar1 = (int *)(param_3 + 0xec170 + *(short *)(param_3 + 0xec620) * 4);
  if (*(char *)(param_3 + 0xec62c + param_2 * 6) == '\0') {
    uVar2 = FUN_10023dc0(param_1,piVar1,param_3,param_4);
    *(short *)(param_3 + 0xec620) = *(short *)(param_3 + 0xec620) + (short)uVar2;
    return uVar2;
  }
  if (0x121 < *(short *)(param_3 + 0xec620)) {
    return param_2 & 0xffff0000;
  }
  uVar2 = FUN_10023e70(param_1,piVar1,param_2,param_3,param_4,param_5);
  *(short *)(param_3 + 0xec620) = *(short *)(param_3 + 0xec620) + (short)uVar2;
  return uVar2;
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



===== 0x10023e70 =====
Function: FUN_10023e70 @ 10023e70

uint __cdecl
FUN_10023e70(int param_1,int *param_2,int param_3,int param_4,int param_5,short param_6)

{
  int iVar1;
  short sVar2;
  uint uVar3;
  uint uVar4;
  int iVar5;
  int local_10 [2];
  char local_8;
  
  iVar5 = param_3 * 3 + 0x76314;
  iVar1 = param_4 + iVar5 * 2;
  local_8 = '\n';
  uVar3 = FUN_10019570(local_10,param_1,(int)*(char *)(param_4 + 4 + iVar5 * 2),param_5);
  sVar2 = (short)uVar3;
  while (sVar2 == 0) {
    if ((param_6 != 1) && (local_8 < '\x06')) goto LAB_10023f7b;
    local_8 = local_8 + -1;
    uVar3 = FUN_10019570(local_10,param_1,(int)*(char *)(iVar1 + 4),param_5);
    sVar2 = (short)uVar3;
  }
  uVar4 = FUN_10023c70(10,10000,param_2,param_1,(short)*(char *)(iVar1 + 4),local_10,param_4,param_5
                      );
  iVar5 = FUN_10023060(param_2,(short)uVar4,param_5);
  while( true ) {
    if ((9 < iVar5) || (uVar3 = CONCAT31((int3)((uint)iVar5 >> 8),local_8), local_8 < '\x01')) {
      return CONCAT22((short)((uint)iVar5 >> 0x10),(short)uVar4);
    }
    if ((param_6 != 1) && (local_8 < '\x06')) break;
    local_8 = local_8 + -1;
    FUN_10019570(local_10,param_1,(int)*(char *)(iVar1 + 4),param_5);
    uVar4 = FUN_10023c70(10,10000,param_2,param_1,(short)*(char *)(iVar1 + 4),local_10,param_4,
                         param_5);
    iVar5 = FUN_10023060(param_2,(short)uVar4,param_5);
  }
LAB_10023f7b:
  return uVar3 & 0xffff0000;
}



===== 0x10023c70 =====
Function: FUN_10023c70 @ 10023c70

uint __cdecl
FUN_10023c70(int param_1,int param_2,int *param_3,int param_4,short param_5,int *param_6,int param_7
            ,int param_8)

{
  int iVar1;
  int iVar2;
  undefined4 *puVar3;
  int iVar4;
  short sVar5;
  byte local_c [8];
  
  iVar1 = param_7 + 0x477ac;
  iVar2 = FUN_10023060((int *)*param_6,(short)param_6[1],param_8);
  if ((param_6[1] <= param_1) && (iVar2 <= param_2)) {
    puVar3 = FUN_10063f30(param_3,(undefined4 *)*param_6,param_6[1] << 2);
    return CONCAT22((short)((uint)puVar3 >> 0x10),(short)param_6[1]);
  }
  FUN_10016ea0(param_4,local_c);
  sVar5 = 0;
  iVar2 = 0;
  if (0 < param_6[1]) {
    do {
      *(undefined4 *)(*(int *)(iVar1 + sVar5 * 4) + 8) = *(undefined4 *)(*param_6 + iVar2 * 4);
      iVar4 = FUN_10023a70(local_c,(byte *)(*(int *)(*param_6 + iVar2 * 4) * 5 +
                                           *(int *)(param_8 + 0x90)),param_5);
      sVar5 = sVar5 + 1;
      *(float *)(*(int *)(iVar1 + iVar2 * 4) + 4) = (float)iVar4;
      if (9999 < sVar5) break;
      iVar2 = iVar2 + 1;
    } while (iVar2 < param_6[1]);
  }
  FUN_1001b5f0((int)sVar5,iVar1);
  if (param_1 < sVar5) {
    sVar5 = (short)param_1;
  }
  iVar4 = 0;
  param_1 = 0;
  iVar2 = 0;
  if (0 < sVar5) {
    iVar1 = iVar1 - (int)param_3;
    do {
      iVar2 = *(int *)(*(int *)(iVar1 + (int)param_3) + 8);
      *param_3 = iVar2;
      iVar4 = iVar4 + (uint)*(ushort *)(*(int *)(param_8 + 0x8c) + iVar2 * 2);
      iVar2 = param_1;
      if (param_2 < iVar4) break;
      iVar2 = param_1 + 1;
      param_3 = param_3 + 1;
      param_1 = iVar2;
    } while (iVar2 < sVar5);
  }
  if (sVar5 <= iVar2) {
    return CONCAT22((short)((uint)iVar2 >> 0x10),sVar5);
  }
  return iVar2 + 1;
}



===== 0x10024060 =====
Function: FUN_10024060 @ 10024060

bool __cdecl FUN_10024060(uint param_1,undefined4 *param_2,int param_3)

{
  byte *pbVar1;
  undefined4 *puVar2;
  byte bVar3;
  char cVar4;
  ushort uVar5;
  ushort uVar6;
  undefined4 *puVar7;
  short sVar8;
  short sVar9;
  int iVar10;
  uint uVar11;
  int iVar12;
  undefined4 local_18;
  byte local_12;
  ushort *local_10;
  int local_c;
  ushort *local_8;
  
  puVar7 = param_2;
  iVar12 = param_2[0x13];
  pbVar1 = (byte *)((int)param_2 + (param_1 * 3 + 0x76314) * 2);
  iVar10 = (uint)*pbVar1 * 0x3c0;
  bVar3 = *(byte *)(iVar10 + 0x6de + iVar12);
  pbVar1[4] = 0;
  pbVar1[3] = bVar3;
  cVar4 = *(char *)(iVar10 + 0x92b + iVar12);
  puVar2 = (undefined4 *)((uint)pbVar1[2] * 7 + iVar10 + 0x6e2 + iVar12);
  *(undefined2 *)(param_2 + 0x3b188) = 0;
  *(undefined2 *)(param_2 + 0x3b05b) = 0;
  if (cVar4 == '\f') {
    FUN_10063f30(&local_18,puVar2,7);
    local_12 = local_12 & 0xdf;
    FUN_10018770(&local_18,param_1,(int)param_2,param_3,1);
    if (*(short *)(param_2 + 0x3b188) == 0) {
      return false;
    }
    sVar8 = 2;
  }
  else {
    sVar8 = 1;
  }
  FUN_10018770(puVar2,param_1,(int)param_2,param_3,sVar8);
  if (*(short *)(param_2 + 0x3b188) == 0) {
    return false;
  }
  sVar8 = FUN_10024010(*(short *)(param_2 + 0x3b188),param_2 + 0x3b05c);
  *(short *)(param_2 + 0x3b188) = sVar8;
  uVar11 = FUN_10023af0(0x1e,10000,param_2 + 0x3b05c,(int)puVar2,param_2,param_3);
  *(short *)(param_2 + 0x3b188) = (short)uVar11;
  sVar8 = FUN_10024010(*(short *)(param_2 + 0x3b05b),param_2 + 0x3b03d);
  *(short *)(param_2 + 0x3b05b) = sVar8;
  local_8 = (ushort *)(param_2 + param_1 * 0x3f + 0x2ba25);
  local_10 = local_8 + 2;
  FUN_10063f30((undefined4 *)local_10,param_2 + 0x3b03d,(int)sVar8 << 2);
  uVar5 = *(ushort *)(param_2 + 0x3b188);
  uVar6 = *(ushort *)(param_2 + 0x3b05b);
  if ((short)uVar6 < (short)uVar5) {
    local_c = 0;
    param_2 = (undefined4 *)(uint)uVar6;
    if (0 < (short)uVar5) {
      do {
        sVar8 = *(short *)(puVar7 + 0x3b05b);
        sVar9 = 0;
        if (0 < sVar8) {
          do {
            if (puVar7[(short)local_c + 0x3b05c] == puVar7[sVar9 + 0x3b03d]) break;
            sVar9 = sVar9 + 1;
          } while (sVar9 < sVar8);
        }
        if (sVar9 == sVar8) {
          puVar7[param_1 * 0x3f + (int)(short)(ushort)param_2 + 0x2ba26] =
               puVar7[(short)local_c + 0x3b05c];
          param_2 = (undefined4 *)((int)param_2 + 1);
        }
        local_c = local_c + 1;
        uVar6 = (ushort)param_2;
      } while ((short)local_c < *(short *)(puVar7 + 0x3b188));
    }
    param_2._0_2_ = uVar6;
    *local_8 = (ushort)param_2;
    local_8[1] = *(ushort *)(puVar7 + 0x3b05b);
  }
  else {
    *local_8 = uVar5;
    local_8[1] = *(ushort *)(param_2 + 0x3b05b);
  }
  iVar12 = FUN_10023060((int *)local_10,*local_8,param_3);
  if ((cVar4 == '\b') && (0 < iVar12)) {
    return true;
  }
  return 9 < iVar12;
}



