===== 0x10016960 =====
Function: FUN_10016960 @ 10016960

undefined4 __cdecl FUN_10016960(undefined4 *param_1,int param_2,undefined4 *param_3,uint param_4)

{
  byte *pbVar1;
  undefined4 *puVar2;
  uint uVar3;
  int iVar4;
  int iVar5;
  uint uVar6;
  int iVar7;
  undefined4 *puVar8;
  
  uVar3 = param_4;
  puVar2 = param_3;
  puVar8 = param_3;
  for (uVar6 = param_4 >> 2; uVar6 != 0; uVar6 = uVar6 - 1) {
    *puVar8 = 0x31313131;
    puVar8 = puVar8 + 1;
  }
  iVar7 = 0;
  uVar6 = param_4 & 3;
  param_4 = 0;
  for (; uVar6 != 0; uVar6 = uVar6 - 1) {
    *(undefined1 *)puVar8 = 0x31;
    puVar8 = (undefined4 *)((int)puVar8 + 1);
  }
  *(undefined1 *)((int)param_3 + (uVar3 - 1)) = 0;
  if (0 < param_2) {
    param_3 = param_1;
    do {
      pbVar1 = (byte *)*param_3;
      iVar4 = FUN_1002ea80(iVar7 + (int)puVar2,1);
      if (iVar4 == 0) {
        return 0xffffffff;
      }
      *(undefined1 *)(iVar7 + (int)puVar2) = 0x22;
      iVar4 = FUN_10016a40(pbVar1,(undefined4 *)(iVar7 + 1 + (int)puVar2));
      if (iVar4 < 1) {
        return 0xffffffff;
      }
      iVar4 = iVar7 + 1 + iVar4;
      iVar7 = FUN_1002ea80(iVar4 + (int)puVar2,1);
      if (iVar7 == 0) {
        return 0xffffffff;
      }
      iVar7 = iVar4 + 1;
      *(undefined1 *)(iVar4 + (int)puVar2) = 0x22;
      if (param_4 != param_2 - 1U) {
        iVar5 = FUN_1002ea80(iVar7 + (int)puVar2,1);
        if (iVar5 == 0) {
          return 0xffffffff;
        }
        *(undefined1 *)(iVar7 + (int)puVar2) = 0x2c;
        iVar7 = iVar4 + 2;
      }
      param_4 = param_4 + 1;
      param_3 = param_3 + 1;
    } while ((int)param_4 < param_2);
  }
  *(undefined1 *)(iVar7 + (int)puVar2) = 0;
  return 1;
}



===== 0x10016a40 =====
Function: FUN_10016a40 @ 10016a40

int __cdecl FUN_10016a40(byte *param_1,undefined4 *param_2)

{
  byte bVar1;
  int iVar2;
  uint uVar3;
  int iVar4;
  undefined4 *puVar5;
  byte *pbVar6;
  
  pbVar6 = param_1;
  uVar3 = 0xffffffff;
  iVar4 = 0;
  do {
    if (uVar3 == 0) break;
    uVar3 = uVar3 - 1;
    bVar1 = *param_1;
    param_1 = param_1 + 1;
  } while (bVar1 != 0);
  param_1 = (byte *)0x0;
  if (0 < (int)(~uVar3 - 1)) {
    do {
      iVar2 = FUN_1002ea80((int)pbVar6,2);
      bVar1 = *pbVar6;
      if (iVar2 == 0) {
        if (bVar1 != 0x22) {
          iVar2 = FUN_1002ea80((int)param_2,1);
          goto joined_r0x10016af4;
        }
LAB_10016acc:
        iVar2 = FUN_1002ea80((int)param_2,2);
        if (iVar2 == 0) {
          return -1;
        }
        *(undefined1 *)param_2 = 0x22;
        iVar4 = iVar4 + 2;
        *(undefined1 *)((int)param_2 + 1) = 0x22;
        puVar5 = (undefined4 *)((int)param_2 + 2);
LAB_10016b06:
        param_1 = param_1 + 1;
        pbVar6 = pbVar6 + 1;
      }
      else {
        if ((bVar1 & 0x80) == 0) {
          if (bVar1 == 0x22) goto LAB_10016acc;
          iVar2 = FUN_1002ea80((int)param_2,1);
joined_r0x10016af4:
          if (iVar2 == 0) {
            return -1;
          }
          iVar4 = iVar4 + 1;
          puVar5 = (undefined4 *)((int)param_2 + 1);
          FUN_10063f30(param_2,(undefined4 *)pbVar6,1);
          goto LAB_10016b06;
        }
        iVar2 = FUN_1002ea80((int)param_2,2);
        if (iVar2 == 0) {
          return -1;
        }
        FUN_10063f30(param_2,(undefined4 *)pbVar6,2);
        param_1 = param_1 + 2;
        pbVar6 = pbVar6 + 2;
        iVar4 = iVar4 + 2;
        puVar5 = (undefined4 *)((int)param_2 + 2);
      }
      param_2 = puVar5;
    } while ((int)param_1 < (int)(~uVar3 - 1));
  }
  return iVar4;
}



===== 0x1002ea80 =====
Function: FUN_1002ea80 @ 1002ea80

undefined4 __cdecl FUN_1002ea80(int param_1,int param_2)

{
  int iVar1;
  
  iVar1 = 0;
  if (0 < param_2) {
    do {
      if (*(char *)(iVar1 + param_1) == '\0') {
        return 0;
      }
      iVar1 = iVar1 + 1;
    } while (iVar1 < param_2);
  }
  return 1;
}



===== 0x1002e1b0 =====
Function: FUN_1002e1b0 @ 1002e1b0

char * __cdecl FUN_1002e1b0(undefined4 *param_1,char *param_2)

{
  char *pcVar1;
  char cVar2;
  char *pcVar3;
  char *pcVar4;
  
  pcVar4 = (char *)*param_1;
  if ((pcVar4 != (char *)0x0) && (pcVar4 <= param_2)) {
    while( true ) {
      if ((((&DAT_1007e188)[*pcVar4] & 6) == 0) && (*pcVar4 != '\0')) {
        pcVar3 = pcVar4;
        cVar2 = *pcVar4;
        while (((cVar2 != '\0' && (cVar2 != '\n')) && (cVar2 != '\r'))) {
          if (param_2 <= pcVar4) {
            return (char *)0x0;
          }
          pcVar1 = pcVar4 + 1;
          pcVar4 = pcVar4 + 1;
          cVar2 = *pcVar1;
        }
        *pcVar4 = '\0';
        *param_1 = pcVar4 + 1;
        return pcVar3;
      }
      if (param_2 <= pcVar4) break;
      pcVar4 = pcVar4 + 1;
    }
  }
  return (char *)0x0;
}



