===== 0x100130e0 =====
Function: FUN_100130e0 @ 100130e0

undefined4 __cdecl FUN_100130e0(int param_1,int param_2)

{
  char cVar1;
  short *psVar2;
  int iVar3;
  char *pcVar4;
  undefined4 *puVar5;
  undefined1 *puVar6;
  uint uVar7;
  undefined4 uVar8;
  short *psVar9;
  int *piVar10;
  int iVar11;
  int *piVar12;
  byte *pbVar13;
  
  psVar2 = *(short **)(param_2 + 0x4c);
  iVar11 = *(int *)(param_2 + 0x1312bc);
  FUN_10012c70((int)psVar2);
  iVar3 = 0;
  if (0 < psVar2[1]) {
    pcVar4 = (char *)((int)psVar2 + 0xa09);
    do {
      iVar3 = iVar3 + 1;
      *pcVar4 = s___Z_Z_________10079a30[*(short *)(pcVar4 + -0x11)];
      pcVar4 = pcVar4 + 0x3c0;
    } while (iVar3 < psVar2[1]);
  }
  if (*(char *)(iVar3 * 0x3c0 + 0x649 + (int)psVar2) != '^') {
    *(undefined1 *)(iVar3 * 0x3c0 + 0x649 + (int)psVar2) = 0x5a;
  }
  if ((char)psVar2[0x23b85] == '\x06') {
    *(undefined1 *)(psVar2[1] * 0x3c0 + 0x649 + (int)psVar2) = 0x5e;
  }
  else {
    *(undefined1 *)(psVar2[1] * 0x3c0 + 0x649 + (int)psVar2) = 0x5a;
  }
  FUN_10012df0(psVar2);
  FUN_10012f00(param_1,param_2);
  iVar3 = 0;
  if (0 < psVar2[1]) {
    psVar9 = psVar2 + 0x504;
    do {
      *(undefined1 *)psVar9 = 0;
      iVar3 = iVar3 + 1;
      psVar9 = psVar9 + 0x1e0;
    } while (iVar3 < psVar2[1]);
  }
  iVar3 = 0;
  if (psVar2[1] != 1 && -1 < psVar2[1] + -1) {
    pcVar4 = (char *)((int)psVar2 + 0xa09);
    do {
      if (((pcVar4[-0xd] == '\b') || (pcVar4[0x3b3] != '\b')) ||
         ((*pcVar4 != ']' && (*pcVar4 != '\\')))) {
        if (*(short *)(pcVar4 + -0x11) == 0xc) goto LAB_100131e8;
      }
      else {
        *pcVar4 = '\\';
LAB_100131e8:
        pcVar4[-1] = '\x01';
      }
      iVar3 = iVar3 + 1;
      pcVar4 = pcVar4 + 0x3c0;
    } while (iVar3 < psVar2[1] + -1);
  }
  iVar3 = 0;
  if (psVar2[1] != 1 && -1 < psVar2[1] + -1) {
    piVar10 = (int *)(param_2 + 0x121a64);
    piVar12 = (int *)(iVar11 + 0x20);
    do {
      iVar11 = *piVar12;
      if (iVar11 == -2) {
        *piVar10 = 100;
      }
      else if (-1 < iVar11) {
        *piVar10 = iVar11;
      }
      iVar3 = iVar3 + 1;
      piVar10 = piVar10 + 1;
      piVar12 = piVar12 + 0x25;
    } while (iVar3 < psVar2[1] + -1);
  }
  iVar11 = 1;
  if (1 < psVar2[1]) {
    puVar5 = (undefined4 *)(param_2 + 0x121d80);
    do {
      if (-1 < (int)puVar5[-199]) {
        *puVar5 = 2;
      }
      iVar11 = iVar11 + 1;
      puVar5 = puVar5 + 1;
    } while (iVar11 < psVar2[1]);
  }
  iVar11 = 0;
  if (psVar2[1] != 1 && -1 < psVar2[1] + -1) {
    puVar6 = (undefined1 *)((int)psVar2 + 0xa09);
    puVar5 = (undefined4 *)(param_2 + 0x121d80);
    do {
      switch(*puVar5) {
      case 0:
        *puVar6 = 0x5d;
        puVar6[-1] = 1;
        break;
      case 1:
        *puVar6 = 0x5c;
        puVar6[-1] = 1;
        break;
      case 2:
        *puVar6 = 0x5b;
        break;
      case 3:
        *puVar6 = 0x5a;
      }
      iVar11 = iVar11 + 1;
      puVar5 = puVar5 + 1;
      puVar6 = puVar6 + 0x3c0;
    } while (iVar11 < psVar2[1] + -1);
  }
  iVar3 = *(int *)(param_2 + 0x121d80 + iVar11 * 4);
  if ((-1 < iVar3) && (iVar3 < 3)) {
    *(undefined1 *)(iVar11 * 0x3c0 + 0xa09 + (int)psVar2) = 0x5b;
  }
  uVar7 = 0;
  iVar11 = 0;
  param_2 = 0;
  if (0 < psVar2[1]) {
    pcVar4 = (char *)((int)psVar2 + 0xa09);
    pbVar13 = (byte *)((int)psVar2 + 0x6e1);
    do {
      iVar11 = iVar11 + (uint)*pbVar13 * 2;
      if (1000 < iVar11) {
        if (uVar7 == param_2) {
          return uVar7 & 0xffff0000;
        }
        uVar7 = uVar7 - 1;
        pbVar13 = pbVar13 + -0x3c0;
        pcVar4 = pcVar4 + -0x3c0;
        *pcVar4 = '[';
      }
      cVar1 = *pcVar4;
      if ((((cVar1 == '[') || (cVar1 == 'Z')) || (cVar1 == '^')) || (cVar1 == '`')) {
        param_2 = uVar7 + 1;
        iVar11 = 0;
      }
      uVar7 = uVar7 + 1;
      pbVar13 = pbVar13 + 0x3c0;
      pcVar4 = pcVar4 + 0x3c0;
    } while ((int)uVar7 < (int)psVar2[1]);
  }
  uVar8 = FUN_10012df0(psVar2);
  return CONCAT22((short)((uint)uVar8 >> 0x10),1);
}



===== 0x10013170 =====
Function: FUN_100130e0 @ 100130e0

undefined4 __cdecl FUN_100130e0(int param_1,int param_2)

{
  char cVar1;
  short *psVar2;
  int iVar3;
  char *pcVar4;
  undefined4 *puVar5;
  undefined1 *puVar6;
  uint uVar7;
  undefined4 uVar8;
  short *psVar9;
  int *piVar10;
  int iVar11;
  int *piVar12;
  byte *pbVar13;
  
  psVar2 = *(short **)(param_2 + 0x4c);
  iVar11 = *(int *)(param_2 + 0x1312bc);
  FUN_10012c70((int)psVar2);
  iVar3 = 0;
  if (0 < psVar2[1]) {
    pcVar4 = (char *)((int)psVar2 + 0xa09);
    do {
      iVar3 = iVar3 + 1;
      *pcVar4 = s___Z_Z_________10079a30[*(short *)(pcVar4 + -0x11)];
      pcVar4 = pcVar4 + 0x3c0;
    } while (iVar3 < psVar2[1]);
  }
  if (*(char *)(iVar3 * 0x3c0 + 0x649 + (int)psVar2) != '^') {
    *(undefined1 *)(iVar3 * 0x3c0 + 0x649 + (int)psVar2) = 0x5a;
  }
  if ((char)psVar2[0x23b85] == '\x06') {
    *(undefined1 *)(psVar2[1] * 0x3c0 + 0x649 + (int)psVar2) = 0x5e;
  }
  else {
    *(undefined1 *)(psVar2[1] * 0x3c0 + 0x649 + (int)psVar2) = 0x5a;
  }
  FUN_10012df0(psVar2);
  FUN_10012f00(param_1,param_2);
  iVar3 = 0;
  if (0 < psVar2[1]) {
    psVar9 = psVar2 + 0x504;
    do {
      *(undefined1 *)psVar9 = 0;
      iVar3 = iVar3 + 1;
      psVar9 = psVar9 + 0x1e0;
    } while (iVar3 < psVar2[1]);
  }
  iVar3 = 0;
  if (psVar2[1] != 1 && -1 < psVar2[1] + -1) {
    pcVar4 = (char *)((int)psVar2 + 0xa09);
    do {
      if (((pcVar4[-0xd] == '\b') || (pcVar4[0x3b3] != '\b')) ||
         ((*pcVar4 != ']' && (*pcVar4 != '\\')))) {
        if (*(short *)(pcVar4 + -0x11) == 0xc) goto LAB_100131e8;
      }
      else {
        *pcVar4 = '\\';
LAB_100131e8:
        pcVar4[-1] = '\x01';
      }
      iVar3 = iVar3 + 1;
      pcVar4 = pcVar4 + 0x3c0;
    } while (iVar3 < psVar2[1] + -1);
  }
  iVar3 = 0;
  if (psVar2[1] != 1 && -1 < psVar2[1] + -1) {
    piVar10 = (int *)(param_2 + 0x121a64);
    piVar12 = (int *)(iVar11 + 0x20);
    do {
      iVar11 = *piVar12;
      if (iVar11 == -2) {
        *piVar10 = 100;
      }
      else if (-1 < iVar11) {
        *piVar10 = iVar11;
      }
      iVar3 = iVar3 + 1;
      piVar10 = piVar10 + 1;
      piVar12 = piVar12 + 0x25;
    } while (iVar3 < psVar2[1] + -1);
  }
  iVar11 = 1;
  if (1 < psVar2[1]) {
    puVar5 = (undefined4 *)(param_2 + 0x121d80);
    do {
      if (-1 < (int)puVar5[-199]) {
        *puVar5 = 2;
      }
      iVar11 = iVar11 + 1;
      puVar5 = puVar5 + 1;
    } while (iVar11 < psVar2[1]);
  }
  iVar11 = 0;
  if (psVar2[1] != 1 && -1 < psVar2[1] + -1) {
    puVar6 = (undefined1 *)((int)psVar2 + 0xa09);
    puVar5 = (undefined4 *)(param_2 + 0x121d80);
    do {
      switch(*puVar5) {
      case 0:
        *puVar6 = 0x5d;
        puVar6[-1] = 1;
        break;
      case 1:
        *puVar6 = 0x5c;
        puVar6[-1] = 1;
        break;
      case 2:
        *puVar6 = 0x5b;
        break;
      case 3:
        *puVar6 = 0x5a;
      }
      iVar11 = iVar11 + 1;
      puVar5 = puVar5 + 1;
      puVar6 = puVar6 + 0x3c0;
    } while (iVar11 < psVar2[1] + -1);
  }
  iVar3 = *(int *)(param_2 + 0x121d80 + iVar11 * 4);
  if ((-1 < iVar3) && (iVar3 < 3)) {
    *(undefined1 *)(iVar11 * 0x3c0 + 0xa09 + (int)psVar2) = 0x5b;
  }
  uVar7 = 0;
  iVar11 = 0;
  param_2 = 0;
  if (0 < psVar2[1]) {
    pcVar4 = (char *)((int)psVar2 + 0xa09);
    pbVar13 = (byte *)((int)psVar2 + 0x6e1);
    do {
      iVar11 = iVar11 + (uint)*pbVar13 * 2;
      if (1000 < iVar11) {
        if (uVar7 == param_2) {
          return uVar7 & 0xffff0000;
        }
        uVar7 = uVar7 - 1;
        pbVar13 = pbVar13 + -0x3c0;
        pcVar4 = pcVar4 + -0x3c0;
        *pcVar4 = '[';
      }
      cVar1 = *pcVar4;
      if ((((cVar1 == '[') || (cVar1 == 'Z')) || (cVar1 == '^')) || (cVar1 == '`')) {
        param_2 = uVar7 + 1;
        iVar11 = 0;
      }
      uVar7 = uVar7 + 1;
      pbVar13 = pbVar13 + 0x3c0;
      pcVar4 = pcVar4 + 0x3c0;
    } while ((int)uVar7 < (int)psVar2[1]);
  }
  uVar8 = FUN_10012df0(psVar2);
  return CONCAT22((short)((uint)uVar8 >> 0x10),1);
}



===== 0x10013250 =====
Function: FUN_100130e0 @ 100130e0

undefined4 __cdecl FUN_100130e0(int param_1,int param_2)

{
  char cVar1;
  short *psVar2;
  int iVar3;
  char *pcVar4;
  undefined4 *puVar5;
  undefined1 *puVar6;
  uint uVar7;
  undefined4 uVar8;
  short *psVar9;
  int *piVar10;
  int iVar11;
  int *piVar12;
  byte *pbVar13;
  
  psVar2 = *(short **)(param_2 + 0x4c);
  iVar11 = *(int *)(param_2 + 0x1312bc);
  FUN_10012c70((int)psVar2);
  iVar3 = 0;
  if (0 < psVar2[1]) {
    pcVar4 = (char *)((int)psVar2 + 0xa09);
    do {
      iVar3 = iVar3 + 1;
      *pcVar4 = s___Z_Z_________10079a30[*(short *)(pcVar4 + -0x11)];
      pcVar4 = pcVar4 + 0x3c0;
    } while (iVar3 < psVar2[1]);
  }
  if (*(char *)(iVar3 * 0x3c0 + 0x649 + (int)psVar2) != '^') {
    *(undefined1 *)(iVar3 * 0x3c0 + 0x649 + (int)psVar2) = 0x5a;
  }
  if ((char)psVar2[0x23b85] == '\x06') {
    *(undefined1 *)(psVar2[1] * 0x3c0 + 0x649 + (int)psVar2) = 0x5e;
  }
  else {
    *(undefined1 *)(psVar2[1] * 0x3c0 + 0x649 + (int)psVar2) = 0x5a;
  }
  FUN_10012df0(psVar2);
  FUN_10012f00(param_1,param_2);
  iVar3 = 0;
  if (0 < psVar2[1]) {
    psVar9 = psVar2 + 0x504;
    do {
      *(undefined1 *)psVar9 = 0;
      iVar3 = iVar3 + 1;
      psVar9 = psVar9 + 0x1e0;
    } while (iVar3 < psVar2[1]);
  }
  iVar3 = 0;
  if (psVar2[1] != 1 && -1 < psVar2[1] + -1) {
    pcVar4 = (char *)((int)psVar2 + 0xa09);
    do {
      if (((pcVar4[-0xd] == '\b') || (pcVar4[0x3b3] != '\b')) ||
         ((*pcVar4 != ']' && (*pcVar4 != '\\')))) {
        if (*(short *)(pcVar4 + -0x11) == 0xc) goto LAB_100131e8;
      }
      else {
        *pcVar4 = '\\';
LAB_100131e8:
        pcVar4[-1] = '\x01';
      }
      iVar3 = iVar3 + 1;
      pcVar4 = pcVar4 + 0x3c0;
    } while (iVar3 < psVar2[1] + -1);
  }
  iVar3 = 0;
  if (psVar2[1] != 1 && -1 < psVar2[1] + -1) {
    piVar10 = (int *)(param_2 + 0x121a64);
    piVar12 = (int *)(iVar11 + 0x20);
    do {
      iVar11 = *piVar12;
      if (iVar11 == -2) {
        *piVar10 = 100;
      }
      else if (-1 < iVar11) {
        *piVar10 = iVar11;
      }
      iVar3 = iVar3 + 1;
      piVar10 = piVar10 + 1;
      piVar12 = piVar12 + 0x25;
    } while (iVar3 < psVar2[1] + -1);
  }
  iVar11 = 1;
  if (1 < psVar2[1]) {
    puVar5 = (undefined4 *)(param_2 + 0x121d80);
    do {
      if (-1 < (int)puVar5[-199]) {
        *puVar5 = 2;
      }
      iVar11 = iVar11 + 1;
      puVar5 = puVar5 + 1;
    } while (iVar11 < psVar2[1]);
  }
  iVar11 = 0;
  if (psVar2[1] != 1 && -1 < psVar2[1] + -1) {
    puVar6 = (undefined1 *)((int)psVar2 + 0xa09);
    puVar5 = (undefined4 *)(param_2 + 0x121d80);
    do {
      switch(*puVar5) {
      case 0:
        *puVar6 = 0x5d;
        puVar6[-1] = 1;
        break;
      case 1:
        *puVar6 = 0x5c;
        puVar6[-1] = 1;
        break;
      case 2:
        *puVar6 = 0x5b;
        break;
      case 3:
        *puVar6 = 0x5a;
      }
      iVar11 = iVar11 + 1;
      puVar5 = puVar5 + 1;
      puVar6 = puVar6 + 0x3c0;
    } while (iVar11 < psVar2[1] + -1);
  }
  iVar3 = *(int *)(param_2 + 0x121d80 + iVar11 * 4);
  if ((-1 < iVar3) && (iVar3 < 3)) {
    *(undefined1 *)(iVar11 * 0x3c0 + 0xa09 + (int)psVar2) = 0x5b;
  }
  uVar7 = 0;
  iVar11 = 0;
  param_2 = 0;
  if (0 < psVar2[1]) {
    pcVar4 = (char *)((int)psVar2 + 0xa09);
    pbVar13 = (byte *)((int)psVar2 + 0x6e1);
    do {
      iVar11 = iVar11 + (uint)*pbVar13 * 2;
      if (1000 < iVar11) {
        if (uVar7 == param_2) {
          return uVar7 & 0xffff0000;
        }
        uVar7 = uVar7 - 1;
        pbVar13 = pbVar13 + -0x3c0;
        pcVar4 = pcVar4 + -0x3c0;
        *pcVar4 = '[';
      }
      cVar1 = *pcVar4;
      if ((((cVar1 == '[') || (cVar1 == 'Z')) || (cVar1 == '^')) || (cVar1 == '`')) {
        param_2 = uVar7 + 1;
        iVar11 = 0;
      }
      uVar7 = uVar7 + 1;
      pbVar13 = pbVar13 + 0x3c0;
      pcVar4 = pcVar4 + 0x3c0;
    } while ((int)uVar7 < (int)psVar2[1]);
  }
  uVar8 = FUN_10012df0(psVar2);
  return CONCAT22((short)((uint)uVar8 >> 0x10),1);
}



===== 0x10012f00 =====
Function: FUN_10012f00 @ 10012f00

undefined4 __cdecl FUN_10012f00(int param_1,int param_2)

{
  int *piVar1;
  int iVar2;
  short sVar3;
  short *psVar4;
  short *psVar5;
  short sVar6;
  int iVar7;
  short local_38;
  short local_36;
  short local_34;
  short local_32;
  short local_30;
  short local_2e;
  undefined2 local_2c;
  undefined2 local_2a;
  ushort local_28;
  ushort local_26;
  ushort local_24;
  undefined2 local_22;
  ushort local_20;
  ushort local_1e;
  ushort local_1c;
  short *local_18;
  int local_14;
  int local_10;
  int local_c;
  int local_8;
  
  psVar4 = *(short **)(param_2 + 0x4c);
  local_10 = 0;
  local_18 = psVar4;
  if (*psVar4 < 1) {
    return 1;
  }
  do {
    iVar2 = local_10;
    sVar6 = 0;
    local_c = 0;
    param_2 = 0;
    local_8 = 0;
    local_14 = psVar4[local_10 * 8 + 7] + -1;
    psVar5 = psVar4;
    if (0 < local_14) {
      iVar7 = local_14 * 0x3c0;
      do {
        sVar6 = sVar6 + 1;
        local_8 = local_8 + 1;
        piVar1 = (int *)(iVar7 + *(int *)(psVar4 + iVar2 * 8 + 10));
        local_c = local_c + (uint)*(byte *)(piVar1 + 0x25);
        param_2 = param_2 + (uint)*(byte *)(piVar1 + 0x25);
        if ((piVar1[-0xef] != piVar1[1]) && (piVar1[-0xef] + 1 != *piVar1)) {
          local_38 = psVar4[iVar2 * 8 + 7] - sVar6;
          local_34 = psVar4[iVar2 * 8 + 7];
          local_30 = (short)local_c;
          local_32 = local_18[(local_10 + 1) * 8] - local_30;
          local_2e = local_18[(local_10 + 1) * 8];
          local_2c = (undefined2)local_8;
          local_2a = (undefined2)param_2;
          local_28 = (ushort)*(byte *)(iVar7 + 0x3b3 + *(int *)(psVar4 + iVar2 * 8 + 10));
          local_26 = (ushort)*(byte *)(iVar7 + -0xe + *(int *)(psVar4 + iVar2 * 8 + 10));
          local_24 = (ushort)*(byte *)(iVar7 + 0x3b2 + *(int *)(psVar4 + iVar2 * 8 + 10));
          local_22 = *(undefined2 *)(iVar7 + -0x12 + *(int *)(psVar4 + iVar2 * 8 + 10));
          local_20 = (ushort)(*(int *)(iVar7 + -0x3bc + *(int *)(psVar4 + iVar2 * 8 + 10)) ==
                             *(int *)(iVar7 + *(int *)(psVar4 + iVar2 * 8 + 10) + 4));
          local_1e = (ushort)(byte)(&DAT_1007b9e0)
                                   [**(char **)(iVar7 + 0x2e4 + *(int *)(psVar4 + iVar2 * 8 + 10))];
          local_1c = (ushort)(byte)(&DAT_1007b9e0)
                                   [*(char *)((*(byte *)(iVar7 + -0x32b +
                                                        *(int *)(psVar4 + iVar2 * 8 + 10)) - 1) +
                                             *(int *)(iVar7 + -0xdc +
                                                     *(int *)(psVar4 + iVar2 * 8 + 10)))];
          local_36 = sVar6;
          sVar3 = FUN_10001670((int *)(*(int *)(param_1 + 0x4d08) + 0x200),(int)&local_38);
          if (500 < sVar3) {
            *(undefined1 *)(iVar7 + -3 + *(int *)(psVar4 + iVar2 * 8 + 10)) = 0x5c;
            param_2 = 0;
            local_8 = 0;
          }
        }
        iVar7 = iVar7 + -0x3c0;
        local_14 = local_14 + -1;
        psVar5 = local_18;
      } while (local_14 != 0);
    }
    local_10 = local_10 + 1;
    psVar4 = psVar5;
  } while (local_10 < *psVar5);
  return 1;
}



===== 0x10012df0 =====
Function: FUN_10012df0 @ 10012df0

undefined4 __cdecl FUN_10012df0(short *param_1)

{
  short *psVar1;
  short sVar2;
  short *psVar3;
  short *psVar4;
  short *psVar5;
  int iVar6;
  char *pcVar7;
  int iVar8;
  int iVar9;
  
  psVar1 = param_1;
  iVar9 = 0;
  sVar2 = param_1[1];
  if (sVar2 != 0) {
    *(short **)(param_1 + 10) = param_1 + 0x326;
    psVar3 = param_1 + 7;
    param_1 = (short *)0x0;
    psVar1[6] = 0;
    *psVar3 = 0;
    psVar5 = (short *)0x0;
    if (sVar2 != 1 && -1 < sVar2 + -1) {
      pcVar7 = (char *)((int)psVar1 + 0xa09);
      psVar4 = psVar3;
      do {
        *psVar3 = *psVar3 + 1;
        psVar5 = psVar4;
        if ((*pcVar7 != '\\') && (*pcVar7 != ']')) {
          param_1 = (short *)((int)param_1 + 1);
          psVar5 = psVar4 + 8;
          *(char **)(psVar4 + 0xb) = pcVar7 + 3;
          psVar4[7] = (short)iVar9 + 1;
          *psVar5 = 0;
          psVar3 = psVar5;
        }
        iVar9 = iVar9 + 1;
        pcVar7 = pcVar7 + 0x3c0;
        psVar4 = psVar5;
        psVar5 = param_1;
      } while (iVar9 < psVar1[1] + -1);
    }
    iVar9 = 0;
    psVar1[(int)psVar5 * 8 + 7] = psVar1[(int)psVar5 * 8 + 7] + 1;
    param_1 = (short *)0x0;
    sVar2 = (short)psVar5 + 1;
    *psVar1 = sVar2;
    if (0 < sVar2) {
      psVar3 = psVar1 + 7;
      do {
        *(short **)(psVar3 + 5) = psVar1 + iVar9 * 0xf + 0xbea6;
        psVar3[1] = 0;
        iVar6 = 0;
        iVar8 = 0;
        if (0 < *psVar3) {
          do {
            psVar3[1] = psVar3[1] + (ushort)*(byte *)(iVar8 + 0x94 + *(int *)(psVar3 + 3));
            iVar6 = iVar6 + 1;
            iVar8 = iVar8 + 0x3c0;
          } while (iVar6 < *psVar3);
        }
        iVar9 = iVar9 + psVar3[1];
        param_1 = (short *)((int)param_1 + 1);
        psVar3 = psVar3 + 8;
      } while ((int)param_1 < (int)*psVar1);
    }
    return 1;
  }
  *param_1 = 0;
  return 1;
}



===== 0x10012c70 =====
Function: FUN_10012c70 @ 10012c70

undefined4 __cdecl FUN_10012c70(int param_1)

{
  int iVar1;
  byte bVar2;
  char cVar3;
  char extraout_AL;
  int iVar4;
  short sVar5;
  short sVar6;
  undefined1 local_98 [68];
  char local_54 [68];
  int local_10;
  int local_c;
  char *local_8;
  
  *(undefined2 *)(param_1 + 4) = 0;
  local_c = 0;
  if (*(short *)(param_1 + 2) < 1) {
    return 1;
  }
  do {
    iVar1 = (short)local_c * 0x3c0 + 0x64c + param_1;
    bVar2 = *(byte *)(iVar1 + 0x95);
    *(int *)(iVar1 + 8) = param_1 + 0x17d4c + *(short *)(param_1 + 4) * 0x1e;
    if (bVar2 != 0) {
      local_10 = (int)(short)(ushort)bVar2;
      local_8 = (char *)(iVar1 + 0x2e8);
      iVar4 = 0;
      do {
        cVar3 = *local_8;
        local_98[iVar4] = (&DAT_1007b6c0)[cVar3];
        local_54[iVar4] = (&DAT_1007baa8)[cVar3];
        local_8 = local_8 + 1;
        local_10 = local_10 + -1;
        iVar4 = iVar4 + 1;
      } while (local_10 != 0);
      local_10 = 0;
    }
    sVar5 = 0;
    sVar6 = 0;
    *(undefined1 *)(iVar1 + 0x94) = 0;
    if (bVar2 != 0) {
      do {
        if ((*(char *)(sVar6 + 0x329 + iVar1) != '0') || ((int)sVar6 == bVar2 - 1)) {
          *(char *)(*(int *)(iVar1 + 8) + 0x1c + (uint)*(byte *)(iVar1 + 0x94) * 0x1e) = (char)sVar5
          ;
          iVar4 = (int)sVar5;
          FUN_10013c00((int)(local_98 + iVar4),local_54 + iVar4,(sVar6 - sVar5) + 1,
                       (undefined1 *)(*(int *)(iVar1 + 8) + (uint)*(byte *)(iVar1 + 0x94) * 0x1e),
                       iVar4 + 0x2e8 + iVar1,iVar4 + 0x36a + iVar1);
          *(char *)(iVar1 + 0x94) = *(char *)(iVar1 + 0x94) + extraout_AL;
          sVar5 = sVar6 + 1;
        }
        bVar2 = *(byte *)(iVar1 + 0x95);
        sVar6 = sVar6 + 1;
      } while (sVar6 < (short)(ushort)bVar2);
    }
    *(short *)(param_1 + 4) = *(short *)(param_1 + 4) + (ushort)*(byte *)(iVar1 + 0x94);
    local_c = local_c + 1;
  } while ((short)local_c < *(short *)(param_1 + 2));
  return CONCAT22((short)((uint)local_c >> 0x10),1);
}



===== 0x10013da0 =====
Function: FUN_10013da0 @ 10013da0

undefined4 __cdecl FUN_10013da0(int param_1,int param_2,int param_3)

{
  char cVar1;
  byte bVar2;
  byte *pbVar3;
  int iVar4;
  uint uVar5;
  uint uVar6;
  int *piVar7;
  byte *pbVar8;
  undefined4 uVar9;
  char *pcVar10;
  char *pcVar11;
  undefined **ppuVar12;
  bool bVar13;
  byte local_18 [20];
  
  uVar9 = 0;
  if ((((0 < param_2) && (0 < param_3)) && (iVar4 = (param_3 - param_2) + 1, 0 < iVar4)) &&
     (iVar4 < 4)) {
    if (iVar4 == 1) {
      iVar4 = *(char *)(param_1 + param_2 * 4) * 0xb;
      if ((&DAT_1007987b)[iVar4] == '\0') {
        pbVar8 = &DAT_10079194;
        pbVar3 = &DAT_10079878 + iVar4;
        do {
          bVar2 = *pbVar3;
          bVar13 = bVar2 < *pbVar8;
          if (bVar2 != *pbVar8) {
LAB_10013e2b:
            iVar4 = (1 - (uint)bVar13) - (uint)(bVar13 != 0);
            goto LAB_10013e30;
          }
          if (bVar2 == 0) break;
          bVar2 = pbVar3[1];
          bVar13 = bVar2 < pbVar8[1];
          if (bVar2 != pbVar8[1]) goto LAB_10013e2b;
          pbVar3 = pbVar3 + 2;
          pbVar8 = pbVar8 + 2;
        } while (bVar2 != 0);
        iVar4 = 0;
LAB_10013e30:
        if (iVar4 != 0) {
          uVar9 = 1;
        }
      }
      return uVar9;
    }
    local_18[0] = 0;
    if (param_2 <= param_3) {
      piVar7 = (int *)(param_1 + param_2 * 4);
      param_1 = iVar4;
      do {
        iVar4 = *piVar7;
        piVar7 = piVar7 + 1;
        uVar5 = 0xffffffff;
        pcVar10 = &DAT_10079878 + iVar4 * 0xb;
        do {
          pcVar11 = pcVar10;
          if (uVar5 == 0) break;
          uVar5 = uVar5 - 1;
          pcVar11 = pcVar10 + 1;
          cVar1 = *pcVar10;
          pcVar10 = pcVar11;
        } while (cVar1 != '\0');
        uVar5 = ~uVar5;
        iVar4 = -1;
        pbVar3 = local_18;
        do {
          pbVar8 = pbVar3;
          if (iVar4 == 0) break;
          iVar4 = iVar4 + -1;
          pbVar8 = pbVar3 + 1;
          bVar2 = *pbVar3;
          pbVar3 = pbVar8;
        } while (bVar2 != 0);
        pbVar3 = (byte *)(pcVar11 + -uVar5);
        pbVar8 = pbVar8 + -1;
        for (uVar6 = uVar5 >> 2; uVar6 != 0; uVar6 = uVar6 - 1) {
          *(undefined4 *)pbVar8 = *(undefined4 *)pbVar3;
          pbVar3 = pbVar3 + 4;
          pbVar8 = pbVar8 + 4;
        }
        for (uVar5 = uVar5 & 3; uVar5 != 0; uVar5 = uVar5 - 1) {
          *pbVar8 = *pbVar3;
          pbVar3 = pbVar3 + 1;
          pbVar8 = pbVar8 + 1;
        }
        uVar5 = 0xffffffff;
        pcVar10 = &DAT_1007b3cc;
        do {
          pcVar11 = pcVar10;
          if (uVar5 == 0) break;
          uVar5 = uVar5 - 1;
          pcVar11 = pcVar10 + 1;
          cVar1 = *pcVar10;
          pcVar10 = pcVar11;
        } while (cVar1 != '\0');
        uVar5 = ~uVar5;
        iVar4 = -1;
        pbVar3 = local_18;
        do {
          pbVar8 = pbVar3;
          if (iVar4 == 0) break;
          iVar4 = iVar4 + -1;
          pbVar8 = pbVar3 + 1;
          bVar2 = *pbVar3;
          pbVar3 = pbVar8;
        } while (bVar2 != 0);
        pbVar3 = (byte *)(pcVar11 + -uVar5);
        pbVar8 = pbVar8 + -1;
        for (uVar6 = uVar5 >> 2; uVar6 != 0; uVar6 = uVar6 - 1) {
          *(undefined4 *)pbVar8 = *(undefined4 *)pbVar3;
          pbVar3 = pbVar3 + 4;
          pbVar8 = pbVar8 + 4;
        }
        param_1 = param_1 + -1;
        for (uVar5 = uVar5 & 3; uVar5 != 0; uVar5 = uVar5 - 1) {
          *pbVar8 = *pbVar3;
          pbVar3 = pbVar3 + 1;
          pbVar8 = pbVar8 + 1;
        }
      } while (param_1 != 0);
    }
    uVar5 = 0xffffffff;
    pbVar3 = local_18;
    do {
      if (uVar5 == 0) break;
      uVar5 = uVar5 - 1;
      bVar2 = *pbVar3;
      pbVar3 = pbVar3 + 1;
    } while (bVar2 != 0);
    (&stack0xffffffe6)[~uVar5] = 0;
    for (ppuVar12 = &PTR_DAT_10079a4c; pbVar3 = *ppuVar12, *pbVar3 != 0; ppuVar12 = ppuVar12 + 1) {
      pbVar8 = local_18;
      do {
        bVar2 = *pbVar8;
        bVar13 = bVar2 < *pbVar3;
        if (bVar2 != *pbVar3) {
LAB_10013f06:
          iVar4 = (1 - (uint)bVar13) - (uint)(bVar13 != 0);
          goto LAB_10013f0b;
        }
        if (bVar2 == 0) break;
        bVar2 = pbVar8[1];
        bVar13 = bVar2 < pbVar3[1];
        if (bVar2 != pbVar3[1]) goto LAB_10013f06;
        pbVar8 = pbVar8 + 2;
        pbVar3 = pbVar3 + 2;
      } while (bVar2 != 0);
      iVar4 = 0;
LAB_10013f0b:
      if (iVar4 == 0) {
        return 1;
      }
    }
  }
  return 0;
}



===== 0x10013f30 =====
Function: FUN_10013f30 @ 10013f30

uint __cdecl FUN_10013f30(int *param_1,int param_2,char *param_3,undefined4 param_4,int param_5)

{
  int *piVar1;
  int *piVar2;
  int iVar3;
  int iVar4;
  int iVar5;
  int iVar6;
  int local_31c [66];
  int aiStack_214 [66];
  int local_10c [66];
  
  if (0x41 < param_5) {
    param_1[0x14a] = 0;
    return (uint)param_1 & 0xffff0000;
  }
  piVar1 = (int *)0x0;
  local_31c[0] = param_5;
  local_10c[0] = param_5;
  if (0 < param_5) {
    iVar5 = param_2 - (int)param_3;
    do {
      local_10c[(int)piVar1 + 1] = (int)param_3[iVar5];
      piVar1 = (int *)((int)piVar1 + 1);
      local_31c[(int)piVar1] = (int)*param_3;
      param_3 = param_3 + 1;
    } while ((int)piVar1 < param_5);
  }
  iVar5 = 0;
  if (0 < param_5) {
    piVar1 = (int *)0x0;
    piVar2 = aiStack_214;
    for (iVar4 = param_5; piVar2 = piVar2 + 1, iVar4 != 0; iVar4 = iVar4 + -1) {
      *piVar2 = 0;
    }
  }
  param_3 = (char *)0x1;
  if (0 < param_5) {
    iVar4 = 0;
    piVar2 = param_1 + 0x43;
    do {
      piVar1 = piVar2;
      if (((&DAT_1007987b)[*(int *)((int)local_10c + iVar4 + 4) * 0xb] == '\x01') && (iVar5 < 0x41))
      {
        iVar5 = iVar5 + 1;
        piVar1 = piVar2 + 4;
        *(int *)((int)aiStack_214 + iVar4 + 4) = iVar5;
        piVar2[3] = (int)param_3;
        iVar3 = *(int *)((int)local_31c + iVar4 + 4);
        *piVar1 = 0;
        piVar2[5] = 0;
        piVar2[6] = iVar3;
      }
      iVar4 = iVar4 + 4;
      param_3 = param_3 + 1;
      piVar2 = piVar1;
    } while ((int)param_3 <= param_5);
    if (iVar5 != 0) {
      param_2 = 1;
      if (0 < iVar5) {
        piVar1 = param_1 + 0x46;
        do {
          iVar6 = *piVar1 + -1;
          iVar3 = FUN_10013da0((int)local_10c,iVar6,iVar6);
          iVar4 = iVar6;
          while ((iVar3 != 0 && ((iVar4 = iVar4 + -1, param_2 < 2 || (piVar1[-4] < iVar4))))) {
            iVar3 = FUN_10013da0((int)local_10c,iVar4,iVar6);
          }
          iVar4 = iVar4 + 1;
          if (iVar6 < iVar4) {
            piVar1[1] = *piVar1;
          }
          else {
            piVar1[1] = iVar4;
            piVar2 = aiStack_214 + iVar4;
            for (iVar3 = (iVar6 - iVar4) + 1; iVar3 != 0; iVar3 = iVar3 + -1) {
              *piVar2 = param_2;
              piVar2 = piVar2 + 1;
            }
          }
          param_2 = param_2 + 1;
          piVar1 = piVar1 + 4;
        } while (param_2 <= iVar5);
      }
      iVar4 = 1;
      if (0 < iVar5) {
        piVar1 = param_1 + 0x48;
        do {
          iVar3 = piVar1[-2] + 1;
          for (iVar6 = iVar3; (iVar6 <= param_5 && (aiStack_214[iVar6] == 0)); iVar6 = iVar6 + 1) {
          }
          iVar6 = iVar6 + -1;
          if (iVar6 < iVar3) {
            *piVar1 = piVar1[-2];
          }
          else {
            *piVar1 = iVar6;
            piVar2 = aiStack_214 + iVar3;
            for (iVar6 = (iVar6 - iVar3) + 1; iVar6 != 0; iVar6 = iVar6 + -1) {
              *piVar2 = iVar4;
              piVar2 = piVar2 + 1;
            }
          }
          iVar4 = iVar4 + 1;
          piVar1 = piVar1 + 4;
        } while (iVar4 <= iVar5);
      }
      if (param_1[0x47] != 1) {
        param_1[0x47] = 1;
      }
      if (param_1[(iVar5 + 0x11) * 4] != param_5) {
        param_1[0x47] = param_5;
      }
      param_1[0x14a] = iVar5;
      if (-1 < param_5) {
        piVar1 = local_10c;
        for (iVar4 = param_5 + 1; iVar4 != 0; iVar4 = iVar4 + -1) {
          *param_1 = *piVar1;
          piVar1 = piVar1 + 1;
          param_1 = param_1 + 1;
        }
      }
      return CONCAT22((short)((uint)((iVar5 + 0x11) * 0x10) >> 0x10),1);
    }
  }
  param_1[0x14a] = 0;
  return CONCAT22((short)((uint)piVar1 >> 0x10),1);
}



===== 0x10013c00 =====
Function: FUN_10013c00 @ 10013c00

int __cdecl
FUN_10013c00(int param_1,char *param_2,short param_3,undefined1 *param_4,undefined4 param_5,
            int param_6)

{
  short sVar1;
  short sVar2;
  int iVar3;
  short sVar4;
  int iVar5;
  int local_580 [66];
  int iStack_478;
  short asStack_474 [2];
  int iStack_470;
  undefined1 auStack_46c [1044];
  char *in_stack_ffffffa8;
  undefined1 local_18 [20];
  
  FUN_10013f30(local_580,param_1,param_2,local_18,(int)param_3);
  if ((int)in_stack_ffffffa8 < 0x42) {
    if (in_stack_ffffffa8 == (char *)0x0) {
      param_4[0x1d] = (char)param_3;
      if ((char)param_3 != '\0') {
        iVar3 = 0;
        do {
          *(undefined1 *)(iVar3 + param_6) = 1;
          iVar3 = iVar3 + 1;
        } while (iVar3 < (int)(uint)(byte)param_4[0x1d]);
      }
      *param_4 = 0;
      param_4[1] = 4;
      param_4[2] = 1;
      return 0;
    }
  }
  else {
    in_stack_ffffffa8 = (char *)0x41;
  }
  iVar3 = 0;
  sVar1 = 1;
  if (0 < (int)in_stack_ffffffa8) {
    param_2 = (char *)0x1;
    do {
      param_4[(int)param_2 * 0x1e + -2] = param_4[0x1c] + (char)iVar3;
      sVar2 = asStack_474[(int)param_2 * 8];
      sVar4 = 0;
      if ((int)sVar2 < (&iStack_478)[(int)param_2 * 4]) {
        do {
          iVar5 = (int)sVar4;
          sVar4 = sVar4 + 1;
          *(undefined1 *)(iVar5 + iVar3 + param_6) = 1;
        } while ((int)(short)(sVar2 + sVar4) < (&iStack_478)[(int)param_2 * 4]);
      }
      sVar2 = sVar4 + 1;
      *(undefined1 *)(sVar4 + iVar3 + param_6) = 2;
      sVar4 = (short)(&iStack_478)[(int)param_2 * 4] + 1;
      if ((int)sVar4 <= (&iStack_470)[(int)param_2 * 4]) {
        do {
          sVar4 = sVar4 + 1;
          iVar5 = (int)sVar2;
          sVar2 = sVar2 + 1;
          *(undefined1 *)(iVar5 + iVar3 + param_6) = 3;
        } while ((int)sVar4 <= (&iStack_470)[(int)param_2 * 4]);
      }
      param_4[(int)param_2 * 0x1e + -0x1e] = auStack_46c[(int)param_2 * 0x10];
      iVar5 = (int)sVar2;
      if (*(char *)(iVar3 + param_6) == '\x02') {
        if (*(char *)(iVar3 + iVar5 + -1 + param_6) == '\x02') {
          param_4[(int)param_2 * 0x1e + -0x1d] = 0;
        }
        else {
          param_4[(int)param_2 * 0x1e + -0x1d] = 2;
        }
      }
      else if (*(char *)(iVar3 + iVar5 + -1 + param_6) == '\x02') {
        param_4[(int)param_2 * 0x1e + -0x1d] = 1;
      }
      else {
        param_4[(int)param_2 * 0x1e + -0x1d] = 3;
      }
      iVar3 = iVar3 + iVar5;
      param_4[(int)param_2 * 0x1e + -1] = (char)sVar2;
      if (sVar1 == 1) {
        param_4[2] = 1;
      }
      else if (param_2 == in_stack_ffffffa8) {
        param_4[(int)param_2 * 0x1e + -0x1c] = 3;
      }
      else {
        param_4[(int)param_2 * 0x1e + -0x1c] = 2;
      }
      sVar1 = sVar1 + 1;
      param_2 = (char *)(int)sVar1;
    } while ((int)param_2 <= (int)in_stack_ffffffa8);
  }
  return (int)in_stack_ffffffa8;
}



===== 0x100135d0 =====
Function: FUN_100135d0 @ 100135d0

undefined4 __cdecl FUN_100135d0(short *param_1,int param_2,int param_3,int param_4,int param_5)

{
  char cVar1;
  char cVar2;
  int iVar3;
  int iVar4;
  int iVar5;
  int iVar6;
  
  iVar3 = param_3 * 0x3c0 + param_2;
  iVar4 = iVar3 + 0x64c;
  iVar5 = param_4 * 0x1e;
  iVar6 = param_5 + (uint)*(byte *)(*(int *)(iVar3 + 0x654) + 0x1c + iVar5) + iVar4;
  cVar1 = *(char *)(iVar6 + 0x36a);
  *param_1 = (short)(char)(&DAT_1007b6c0)[*(char *)(iVar6 + 0x2e8)];
  if ((param_4 == 0) && (param_5 == 0)) {
    if ((param_3 == 0) || (*(char *)(iVar3 + 0x649) == '[')) {
      param_1[1] = 0x28;
    }
    else {
      param_1[1] = 0x2a;
    }
  }
  else {
    param_1[1] = (ushort)(byte)(&DAT_1007b6c0)
                               [*(char *)((uint)*(byte *)(*(int *)(iVar3 + 0x654) + 0x1c + iVar5) +
                                          iVar4 + 0x2e7 + param_5)];
  }
  if ((param_4 == *(byte *)(iVar3 + 0x6e0) - 1) &&
     (param_5 == *(byte *)(*(int *)(iVar3 + 0x654) + 0x1d + iVar5) - 1)) {
    cVar2 = *(char *)(iVar3 + 0xa09);
    if ((((cVar2 == '[') || (cVar2 == 'Z')) || (cVar2 == '^')) || (cVar2 == '`')) {
      param_1[2] = 0x28;
    }
    else {
      param_1[2] = 0x2a;
    }
  }
  else {
    param_1[2] = (ushort)(byte)(&DAT_1007b6c0)
                               [*(char *)((uint)*(byte *)(*(int *)(iVar3 + 0x654) + 0x1c + iVar5) +
                                          iVar4 + 0x2e9 + param_5)];
  }
  param_1[3] = (ushort)(byte)(&DAT_1007baa8)
                             [*(char *)((uint)*(byte *)(*(int *)(iVar3 + 0x654) + 0x1c + iVar5) +
                                        iVar4 + 0x2e8 + param_5)];
  param_1[4] = (short)*(char *)(*(int *)(iVar3 + 0x654) + 1 + iVar5);
  param_1[5] = (short)*(char *)(*(int *)(iVar3 + 0x654) + 2 + iVar5);
  if ((param_4 == 0) && ((param_3 == 0 || (*(char *)(iVar3 + 0x649) == '[')))) {
    param_1[6] = 1;
  }
  else if ((param_4 == *(byte *)(iVar3 + 0x6e0) - 1) &&
          (((cVar2 = *(char *)(iVar3 + 0xa09), cVar2 == '[' || (cVar2 == 'Z')) ||
           ((cVar2 == '^' || (cVar2 == '`')))))) {
    param_1[6] = 3;
  }
  else {
    param_1[6] = 2;
  }
  param_1[7] = (ushort)*(byte *)(iVar3 + 0x6e0);
  param_1[8] = (short)cVar1;
  return CONCAT22((short)((uint)iVar4 >> 0x10),1);
}



