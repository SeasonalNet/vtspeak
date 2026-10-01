===== 0x1001c990 =====
Function: FUN_1001c990 @ 1001c990

undefined4 __cdecl FUN_1001c990(byte *param_1,byte *param_2,int param_3)

{
  size_t sVar1;
  byte bVar2;
  int *piVar3;
  byte *pbVar4;
  undefined4 uVar5;
  int iVar6;
  int iVar7;
  byte *pbVar8;
  char *pcVar9;
  uint uVar10;
  int iVar11;
  uint uVar12;
  byte *pbVar13;
  byte *pbVar14;
  int iVar15;
  
  pbVar4 = param_1;
  if (param_2 == (byte *)0x0) {
    return 0;
  }
  if (*param_2 == 0) {
    return 0;
  }
  uVar10 = 0xffffffff;
  pbVar8 = param_2;
  do {
    if (uVar10 == 0) break;
    uVar10 = uVar10 - 1;
    bVar2 = *pbVar8;
    pbVar8 = pbVar8 + 1;
  } while (bVar2 != 0);
  iVar11 = ~uVar10 - 1;
  if (param_3 == 0) {
    pbVar8 = (byte *)FUN_1001d9c0(~uVar10);
    if (pbVar8 == (byte *)0x0) {
      return 0;
    }
    uVar10 = 0xffffffff;
    do {
      pbVar13 = param_2;
      if (uVar10 == 0) break;
      uVar10 = uVar10 - 1;
      pbVar13 = param_2 + 1;
      bVar2 = *param_2;
      param_2 = pbVar13;
    } while (bVar2 != 0);
    uVar10 = ~uVar10;
    pbVar13 = pbVar13 + -uVar10;
    pbVar14 = pbVar8;
    for (uVar12 = uVar10 >> 2; uVar12 != 0; uVar12 = uVar12 - 1) {
      *(undefined4 *)pbVar14 = *(undefined4 *)pbVar13;
      pbVar13 = pbVar13 + 4;
      pbVar14 = pbVar14 + 4;
    }
    for (uVar10 = uVar10 & 3; param_2 = pbVar8, uVar10 != 0; uVar10 = uVar10 - 1) {
      *pbVar14 = *pbVar13;
      pbVar13 = pbVar13 + 1;
      pbVar14 = pbVar14 + 1;
    }
  }
  if (*(undefined **)(param_1 + 0x47794) != (undefined *)0x0) {
    FUN_1001da30(*(undefined **)(param_1 + 0x47794));
  }
  if (*(undefined **)(param_1 + 0x47798) != (undefined *)0x0) {
    FUN_1001da30(*(undefined **)(param_1 + 0x47798));
  }
  if (*(undefined **)(param_1 + 0x47790) != (undefined *)0x0) {
    FUN_1001da30(*(undefined **)(param_1 + 0x47790));
  }
  sVar1 = iVar11 * 4;
  uVar5 = FUN_1001d9c0(sVar1);
  *(undefined4 *)(param_1 + 0x47794) = uVar5;
  uVar5 = FUN_1001d9c0(sVar1);
  *(undefined4 *)(param_1 + 0x47798) = uVar5;
  uVar5 = FUN_1001d9c0(sVar1);
  *(undefined4 *)(param_1 + 0x47790) = uVar5;
  iVar6 = 0;
  if (0 < iVar11) {
    do {
      *(int *)(*(int *)(param_1 + 0x47794) + iVar6 * 4) = iVar6;
      *(int *)(*(int *)(param_1 + 0x47798) + iVar6 * 4) = iVar6;
      iVar6 = iVar6 + 1;
    } while (iVar6 < iVar11);
  }
  iVar15 = 0;
  iVar6 = 0;
  if (0 < iVar11) {
    param_1 = param_2;
    do {
      iVar7 = FUN_1001c900(param_1);
      if (iVar7 == 0) {
        *(int *)(*(int *)(pbVar4 + 0x47790) + iVar15 * 4) = iVar6;
        iVar15 = iVar15 + 1;
        param_1 = param_1 + 1;
      }
      else {
        iVar15 = iVar15 + 2;
        *(int *)(*(int *)(pbVar4 + 0x47790) + -8 + iVar15 * 4) = iVar6;
        *(int *)(*(int *)(pbVar4 + 0x47790) + -4 + iVar15 * 4) = iVar6;
        param_1 = param_1 + 2;
      }
      iVar6 = iVar6 + 1;
    } while (iVar15 < iVar11);
  }
  pbVar8 = FUN_1001ce10((int)pbVar4,param_2,param_3);
  pcVar9 = (char *)FUN_1001d250((int)pbVar4,(char *)pbVar8,FUN_1001d370);
  uVar5 = FUN_1001d250((int)pbVar4,pcVar9,FUN_1001d5d0);
  piVar3 = *(int **)(pbVar4 + 0x122448);
  param_3 = *piVar3;
  if (0 < param_3) {
    if (*(char *)(DAT_100a0460 + 0x20424) == '\x01') {
      if (0 < param_3) {
        iVar11 = 0;
        do {
          iVar6 = iVar11 + 8;
          iVar15 = iVar11 + 8;
          iVar11 = iVar11 + 0x210;
          *(undefined4 *)(iVar15 + piVar3[1]) =
               *(undefined4 *)
                (*(int *)(pbVar4 + 0x47790) +
                *(int *)(*(int *)(pbVar4 + 0x47794) + *(int *)(iVar6 + piVar3[1]) * 4) * 4);
          param_3 = param_3 + -1;
        } while (param_3 != 0);
      }
    }
    else if (0 < param_3) {
      iVar11 = 0;
      do {
        iVar6 = iVar11 + 8;
        iVar15 = iVar11 + 8;
        iVar11 = iVar11 + 0x210;
        *(undefined4 *)(iVar15 + piVar3[1]) =
             *(undefined4 *)(*(int *)(pbVar4 + 0x47794) + *(int *)(iVar6 + piVar3[1]) * 4);
        param_3 = param_3 + -1;
      } while (param_3 != 0);
    }
  }
  VT_SetDecimal0Pron_ENG();
  return uVar5;
}



