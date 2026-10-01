===== 0x1000d208 =====
Function: FUN_1000d190 @ 1000d190

int __cdecl FUN_1000d190(int *param_1,int param_2,int param_3,short param_4)

{
  byte bVar1;
  int *piVar2;
  char cVar3;
  short sVar4;
  undefined3 extraout_var;
  int iVar5;
  uint uVar6;
  uint uVar7;
  int iVar8;
  byte *pbVar9;
  int *piVar10;
  int *piVar11;
  char *pcVar12;
  char *pcVar13;
  char *pcVar14;
  byte *pbVar15;
  byte local_234 [200];
  byte local_16c [204];
  char local_a0 [68];
  char local_5c [32];
  char local_3c [32];
  undefined4 local_1c;
  uint local_18;
  int local_14;
  int local_10;
  undefined4 local_c;
  uint local_8;
  
  piVar2 = param_1;
  sVar4 = 1;
  local_c = 1;
  local_14 = 0;
  local_10 = 0;
  if (param_3 < 1) {
LAB_1000d415:
    return (int)sVar4;
  }
  param_1 = (int *)(param_2 + 0x2c);
  do {
    uVar6 = 0xffffffff;
    piVar10 = param_1 + 2;
    do {
      piVar11 = piVar10;
      if (uVar6 == 0) break;
      uVar6 = uVar6 - 1;
      piVar11 = (int *)((int)piVar10 + 1);
      iVar8 = *piVar10;
      piVar10 = piVar11;
    } while ((char)iVar8 != '\0');
    uVar6 = ~uVar6;
    pcVar12 = (char *)((int)piVar11 - uVar6);
    pcVar14 = local_3c;
    for (uVar7 = uVar6 >> 2; uVar7 != 0; uVar7 = uVar7 - 1) {
      *(undefined4 *)pcVar14 = *(undefined4 *)pcVar12;
      pcVar12 = pcVar12 + 4;
      pcVar14 = pcVar14 + 4;
    }
    for (uVar6 = uVar6 & 3; uVar6 != 0; uVar6 = uVar6 - 1) {
      *pcVar14 = *pcVar12;
      pcVar12 = pcVar12 + 1;
      pcVar14 = pcVar14 + 1;
    }
    local_8 = (uint)(local_10 == param_3 + -1);
    if ((param_1[-9] < 0x1d) && (*param_1 != 0)) {
      switch(*param_1) {
      case 2:
        pcVar12 = (char *)&DAT_10077700;
        break;
      case 3:
        pcVar12 = (char *)&DAT_100776f4;
        break;
      case 4:
        pcVar12 = (char *)&DAT_100776f8;
        break;
      case 5:
      case 0xc:
        pcVar12 = (char *)&DAT_100776fc;
        break;
      default:
        goto switchD_1000d214_caseD_6;
      case 0xb:
        pcVar12 = &DAT_1007739c;
      }
      uVar6 = 0xffffffff;
      do {
        pcVar14 = pcVar12;
        if (uVar6 == 0) break;
        uVar6 = uVar6 - 1;
        pcVar14 = pcVar12 + 1;
        cVar3 = *pcVar12;
        pcVar12 = pcVar14;
      } while (cVar3 != '\0');
      uVar6 = ~uVar6;
      iVar8 = -1;
      pcVar12 = local_3c;
      do {
        pcVar13 = pcVar12;
        if (iVar8 == 0) break;
        iVar8 = iVar8 + -1;
        pcVar13 = pcVar12 + 1;
        cVar3 = *pcVar12;
        pcVar12 = pcVar13;
      } while (cVar3 != '\0');
      pcVar12 = pcVar14 + -uVar6;
      pcVar14 = pcVar13 + -1;
      for (uVar7 = uVar6 >> 2; uVar7 != 0; uVar7 = uVar7 - 1) {
        *(undefined4 *)pcVar14 = *(undefined4 *)pcVar12;
        pcVar12 = pcVar12 + 4;
        pcVar14 = pcVar14 + 4;
      }
      for (uVar6 = uVar6 & 3; uVar6 != 0; uVar6 = uVar6 - 1) {
        *pcVar14 = *pcVar12;
        pcVar12 = pcVar12 + 1;
        pcVar14 = pcVar14 + 1;
      }
    }
switchD_1000d214_caseD_6:
    uVar6 = 0xffffffff;
    iVar8 = param_1[1];
    pcVar12 = (char *)((int)param_1 + 0x26);
    do {
      pcVar14 = pcVar12;
      if (uVar6 == 0) break;
      uVar6 = uVar6 - 1;
      pcVar14 = pcVar12 + 1;
      cVar3 = *pcVar12;
      pcVar12 = pcVar14;
    } while (cVar3 != '\0');
    uVar6 = ~uVar6;
    pcVar12 = pcVar14 + -uVar6;
    pcVar14 = local_a0;
    for (uVar7 = uVar6 >> 2; uVar7 != 0; uVar7 = uVar7 - 1) {
      *(undefined4 *)pcVar14 = *(undefined4 *)pcVar12;
      pcVar12 = pcVar12 + 4;
      pcVar14 = pcVar14 + 4;
    }
    iVar5 = param_1[-2];
    for (uVar6 = uVar6 & 3; uVar6 != 0; uVar6 = uVar6 - 1) {
      *pcVar14 = *pcVar12;
      pcVar12 = pcVar12 + 1;
      pcVar14 = pcVar14 + 1;
    }
    uVar6 = (uint)((char)iVar5 == 'S');
    local_18 = uVar6;
    do {
      FUN_1000d640(local_3c,local_5c,(short)local_8,(char)iVar8,local_a0,(short)uVar6,
                   (int)(piVar2 + 0x1204a));
      cVar3 = FUN_1000fd20(local_5c);
      local_1c = CONCAT31(extraout_var,cVar3);
      FUN_10003a70(local_5c,local_16c,0,(int)(piVar2 + 0x1204a));
      if ((&stack0x00000000 == (undefined1 *)0x16c) || (local_16c[0] == 0)) {
        local_234[0] = 0;
      }
      else {
        uVar7 = 0xffffffff;
        pbVar9 = local_16c;
        do {
          pbVar15 = pbVar9;
          if (uVar7 == 0) break;
          uVar7 = uVar7 - 1;
          pbVar15 = pbVar9 + 1;
          bVar1 = *pbVar9;
          pbVar9 = pbVar15;
        } while (bVar1 != 0);
        uVar7 = ~uVar7;
        pbVar9 = pbVar15 + -uVar7;
        pbVar15 = local_234;
        for (uVar6 = uVar7 >> 2; uVar6 != 0; uVar6 = uVar6 - 1) {
          *(undefined4 *)pbVar15 = *(undefined4 *)pbVar9;
          pbVar9 = pbVar9 + 4;
          pbVar15 = pbVar15 + 4;
        }
        for (uVar7 = uVar7 & 3; uVar6 = local_18, uVar7 != 0; uVar7 = uVar7 - 1) {
          *pbVar15 = *pbVar9;
          pbVar9 = pbVar9 + 1;
          pbVar15 = pbVar15 + 1;
        }
      }
      sVar4 = (short)*piVar2;
      if (199 < sVar4) {
        local_c = 0;
        break;
      }
      if (sVar4 < 1) {
        iVar5 = -1;
      }
      else {
        iVar5 = (int)(short)piVar2[sVar4 * 0x155 + -0x154];
      }
      FUN_1000d450((undefined2 *)((int)piVar2 + sVar4 * 0x554 + 2),local_5c,local_234,local_14,iVar5
                   ,0,local_1c,(char)iVar8,local_a0,param_4);
      *(short *)piVar2 = (short)*piVar2 + 1;
    } while (local_3c[0] != '\0');
    local_14 = local_14 + 1;
    sVar4 = 0;
    if ((short)local_c == 0) goto LAB_1000d415;
    local_10 = local_10 + 1;
    param_1 = param_1 + 0x25;
    if (param_3 <= local_10) {
      return (int)(short)local_c;
    }
  } while( true );
}



===== 0x1000eb5c =====
Function: FUN_1000ea20 @ 1000ea20

void __cdecl FUN_1000ea20(short *param_1,int param_2,int param_3)

{
  char cVar1;
  int iVar2;
  short *psVar3;
  uint uVar4;
  uint uVar5;
  char *pcVar6;
  char *pcVar7;
  int local_c;
  byte *local_8;
  
  param_1[0x214d1] = 0;
  local_8 = (byte *)0x0;
  if (0 < *param_1) {
    psVar3 = param_1 + 1;
    local_c = 0;
    do {
      *(short *)(local_c + 0x429a8 + (int)param_1) = psVar3[1];
      uVar4 = 0xffffffff;
      pcVar6 = (char *)((int)psVar3 + 5);
      do {
        pcVar7 = pcVar6;
        if (uVar4 == 0) break;
        uVar4 = uVar4 - 1;
        pcVar7 = pcVar6 + 1;
        cVar1 = *pcVar6;
        pcVar6 = pcVar7;
      } while (cVar1 != '\0');
      uVar4 = ~uVar4;
      pcVar6 = pcVar7 + -uVar4;
      pcVar7 = (char *)(local_c + 0x429ad + (int)param_1);
      for (uVar5 = uVar4 >> 2; uVar5 != 0; uVar5 = uVar5 - 1) {
        *(undefined4 *)pcVar7 = *(undefined4 *)pcVar6;
        pcVar6 = pcVar6 + 4;
        pcVar7 = pcVar7 + 4;
      }
      for (uVar4 = uVar4 & 3; uVar4 != 0; uVar4 = uVar4 - 1) {
        *pcVar7 = *pcVar6;
        pcVar6 = pcVar6 + 1;
        pcVar7 = pcVar7 + 1;
      }
      if (*psVar3 < 2) {
        uVar4 = 0xffffffff;
        pcVar6 = (char *)((int)psVar3 + 0x37);
        do {
          pcVar7 = pcVar6;
          if (uVar4 == 0) break;
          uVar4 = uVar4 - 1;
          pcVar7 = pcVar6 + 1;
          cVar1 = *pcVar6;
          pcVar6 = pcVar7;
        } while (cVar1 != '\0');
        uVar4 = ~uVar4;
        pcVar6 = pcVar7 + -uVar4;
        pcVar7 = (char *)(local_c + 0x429cb + (int)param_1);
        for (uVar5 = uVar4 >> 2; uVar5 != 0; uVar5 = uVar5 - 1) {
          *(undefined4 *)pcVar7 = *(undefined4 *)pcVar6;
          pcVar6 = pcVar6 + 4;
          pcVar7 = pcVar7 + 4;
        }
        for (uVar4 = uVar4 & 3; uVar4 != 0; uVar4 = uVar4 - 1) {
          *pcVar7 = *pcVar6;
          pcVar6 = pcVar6 + 1;
          pcVar7 = pcVar7 + 1;
        }
      }
      else {
        FUN_100068b0((byte *)(local_c + 0x429cb + (int)param_1),(int)(param_1 + 1),
                     (byte *)(int)*param_1,local_8,param_3);
      }
      *(char *)(local_c + 0x429ac + (int)param_1) = (char)psVar3[2];
      if (*(char *)(param_2 + 0x24 + psVar3[1] * 0x94) == 'X') {
        *(undefined1 *)(local_c + 0x42a0c + (int)param_1) = 0x58;
      }
      else {
        *(undefined1 *)(local_c + 0x42a0c + (int)param_1) = 0x30;
      }
      *(short *)(local_c + 0x42a0e + (int)param_1) = psVar3[0x2a6];
      *(short *)(local_c + 0x42a10 + (int)param_1) = psVar3[0x2a7];
      *(short *)(local_c + 0x42a12 + (int)param_1) = psVar3[0x2a8];
      *(short *)(local_c + 0x42a14 + (int)param_1) = psVar3[0x2a9];
      iVar2 = *(int *)(param_2 + 0x2c + psVar3[1] * 0x94);
      if (iVar2 == 3) {
        *(undefined2 *)(local_c + 0x429aa + (int)param_1) = 3;
      }
      else if (iVar2 == 0) {
        *(undefined2 *)(local_c + 0x429aa + (int)param_1) = 0;
      }
      else {
        *(undefined2 *)(local_c + 0x429aa + (int)param_1) = 1;
      }
      param_1[0x214d1] = param_1[0x214d1] + 1;
      local_8 = local_8 + 1;
      psVar3 = psVar3 + 0x2aa;
      local_c = local_c + 0x70;
    } while ((int)local_8 < (int)*param_1);
  }
  return;
}



===== 0x10016cee =====
Function: FUN_10016c90 @ 10016c90

undefined4 __cdecl FUN_10016c90(int param_1)

{
  char cVar1;
  short sVar2;
  int iVar3;
  int iVar4;
  char *pcVar5;
  int iVar6;
  uint uVar7;
  undefined2 *puVar8;
  undefined2 *puVar9;
  char *pcVar10;
  int iVar11;
  byte *pbVar12;
  int local_8;
  
  iVar3 = *(int *)(param_1 + 0x4c);
  iVar4 = *(int *)(param_1 + 0x1312bc);
  local_8 = 0;
  if (0 < *(short *)(iVar3 + 2)) {
    puVar8 = (undefined2 *)(iVar3 + 0xa00);
    puVar9 = (undefined2 *)(iVar4 + 0x24);
    do {
      *(undefined2 **)(puVar8 + -0x6a) = puVar9 + 0x12;
      *puVar8 = *puVar9;
      *(undefined2 **)(puVar8 + 2) = puVar9 + 0x10;
      *(undefined2 **)(puVar8 + -0x68) = puVar9 + 0x21;
      *(undefined1 *)((int)puVar8 + -3) = *(undefined1 *)((int)puVar9 + 0x13);
      puVar8[-4] = puVar9[0xe];
      *(undefined1 *)(puVar8 + -2) = *(undefined1 *)(puVar9 + 0xc);
      *(undefined1 *)(puVar8 + -1) = *(undefined1 *)((int)puVar9 + 0x15);
      *(undefined1 *)((int)puVar8 + -1) = *(undefined1 *)(puVar9 + 0xb);
      puVar8[-3] = puVar9[8];
      pcVar5 = *(char **)(puVar8 + -0x68);
      uVar7 = 0xffffffff;
      pcVar10 = pcVar5;
      do {
        if (uVar7 == 0) break;
        uVar7 = uVar7 - 1;
        cVar1 = *pcVar10;
        pcVar10 = pcVar10 + 1;
      } while (cVar1 != '\0');
      param_1 = ~uVar7 - 1;
      if (0x41 < param_1) {
        param_1 = 0x41;
      }
      iVar6 = 0;
      iVar11 = 0;
      if (param_1 < 1) {
        return 0;
      }
      do {
        if (pcVar5[iVar11] == 'd') {
          if (0 < iVar6) {
            *(undefined1 *)((int)puVar8 + iVar6 + -0x8c) = 0x31;
          }
        }
        else if (pcVar5[iVar11] == 'c') {
          if (0 < iVar6) {
            *(undefined1 *)((int)puVar8 + iVar6 + -0x8c) = 0x32;
          }
        }
        else {
          *(undefined1 *)((int)puVar8 + iVar6 + -0x116) = 0;
          cVar1 = pcVar5[iVar11];
          if (cVar1 == 'M') {
            *(undefined1 *)((int)puVar8 + iVar6 + -0x116) = 1;
          }
          else if ((cVar1 < '\x01') || ('E' < cVar1)) break;
          *(char *)((int)puVar8 + iVar6 + -0xcc) = pcVar5[iVar11];
          *(undefined1 *)((int)puVar8 + iVar6 + -0x8b) = 0x30;
          iVar6 = iVar6 + 1;
        }
        iVar11 = iVar11 + 1;
      } while (iVar11 < param_1);
      if (iVar6 == 0) {
        return 0;
      }
      *(char *)((int)puVar8 + -799) = (char)iVar6;
      puVar9 = puVar9 + 0x4a;
      local_8 = local_8 + 1;
      puVar8 = puVar8 + 0x1e0;
    } while (local_8 < *(short *)(iVar3 + 2));
  }
  iVar6 = 0;
  if (0 < *(short *)(iVar3 + 2)) {
    pbVar12 = (byte *)(iVar3 + 0x92b);
    do {
      iVar11 = iVar6 * 4;
      iVar6 = iVar6 + 1;
      *pbVar12 = (*(int *)(*(int *)(iVar3 + 0x4771c) + iVar11) == -1) - 1U & 0xc;
      pbVar12 = pbVar12 + 0x3c0;
    } while (iVar6 < *(short *)(iVar3 + 2));
  }
  iVar6 = *(short *)(iVar3 + 2) * 0x3c0 + iVar3;
  sVar2 = *(short *)(iVar6 + 0x638);
  if (sVar2 == 3) {
    if (*(short *)(iVar4 + 2) == 0) {
      *(undefined1 *)(iVar3 + 0x4770a) = 6;
      return 1;
    }
  }
  else if (sVar2 != 4) {
    if (*(char *)(iVar6 + 0x56b) == '\f') {
      *(undefined1 *)(iVar3 + 0x4770a) = 7;
      return 1;
    }
    *(undefined1 *)(iVar3 + 0x4770a) = 5;
    return 1;
  }
  *(undefined1 *)(iVar3 + 0x4770a) = 7;
  return 1;
}



===== 0x1003e238 =====
Function: FUN_1003e210 @ 1003e210

void __cdecl FUN_1003e210(undefined2 *param_1)

{
  int iVar1;
  undefined4 *puVar2;
  
  *(undefined4 *)(param_1 + 2) = 0;
  *param_1 = 0;
  param_1[1] = 0;
  if (param_1[6] != 1) {
    param_1[6] = 0;
  }
  puVar2 = (undefined4 *)(param_1 + 10);
  for (iVar1 = 0xe74; iVar1 != 0; iVar1 = iVar1 + -1) {
    *puVar2 = 0;
    puVar2 = puVar2 + 1;
  }
  return;
}



