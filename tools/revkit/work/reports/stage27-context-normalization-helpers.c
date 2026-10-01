===== 0x10003a70 =====
Function: FUN_10003a70 @ 10003a70

void __cdecl FUN_10003a70(char *param_1,byte *param_2,short param_3,int param_4)

{
  char cVar1;
  uint uVar2;
  uint uVar3;
  char *pcVar4;
  byte *pbVar5;
  char *pcVar6;
  byte *pbVar7;
  byte local_f0 [200];
  uint local_28;
  char local_24 [32];
  
  if (param_3 == 0) {
    FUN_1000fd70(param_1,(int)local_24);
  }
  else {
    uVar2 = 0xffffffff;
    do {
      pcVar4 = param_1;
      if (uVar2 == 0) break;
      uVar2 = uVar2 - 1;
      pcVar4 = param_1 + 1;
      cVar1 = *param_1;
      param_1 = pcVar4;
    } while (cVar1 != '\0');
    uVar2 = ~uVar2;
    pcVar4 = pcVar4 + -uVar2;
    pcVar6 = local_24;
    for (uVar3 = uVar2 >> 2; uVar3 != 0; uVar3 = uVar3 - 1) {
      *(undefined4 *)pcVar6 = *(undefined4 *)pcVar4;
      pcVar4 = pcVar4 + 4;
      pcVar6 = pcVar6 + 4;
    }
    for (uVar2 = uVar2 & 3; uVar2 != 0; uVar2 = uVar2 - 1) {
      *pcVar6 = *pcVar4;
      pcVar4 = pcVar4 + 1;
      pcVar6 = pcVar6 + 1;
    }
  }
  uVar2 = 0xffffffff;
  pcVar4 = local_24;
  do {
    pcVar6 = pcVar4;
    if (uVar2 == 0) break;
    uVar2 = uVar2 - 1;
    pcVar6 = pcVar4 + 1;
    cVar1 = *pcVar4;
    pcVar4 = pcVar6;
  } while (cVar1 != '\0');
  uVar2 = ~uVar2;
  pbVar5 = (byte *)(pcVar6 + -uVar2);
  pbVar7 = local_f0;
  for (uVar3 = uVar2 >> 2; uVar3 != 0; uVar3 = uVar3 - 1) {
    *(undefined4 *)pbVar7 = *(undefined4 *)pbVar5;
    pbVar5 = pbVar5 + 4;
    pbVar7 = pbVar7 + 4;
  }
  for (uVar2 = uVar2 & 3; uVar2 != 0; uVar2 = uVar2 - 1) {
    *pbVar7 = *pbVar5;
    pbVar5 = pbVar5 + 1;
    pbVar7 = pbVar7 + 1;
  }
  local_28 = 0xffffffff;
  pcVar4 = local_24;
  do {
    if (local_28 == 0) break;
    local_28 = local_28 - 1;
    cVar1 = *pcVar4;
    pcVar4 = pcVar4 + 1;
  } while (cVar1 != '\0');
  local_28 = ~local_28;
  FUN_10011820(local_f0,param_2,param_4,0);
  return;
}



===== 0x10003c50 =====
Function: FUN_10003c50 @ 10003c50

undefined4 __cdecl FUN_10003c50(undefined1 *param_1,byte *param_2)

{
  byte *pbVar1;
  int iVar2;
  byte bVar3;
  char cVar4;
  byte *pbVar5;
  uint uVar6;
  uint uVar7;
  int iVar8;
  byte *pbVar9;
  byte *pbVar10;
  char *pcVar11;
  char *pcVar12;
  char *pcVar13;
  byte local_cc;
  byte local_cb;
  byte local_ca [198];
  
  *(undefined2 *)(param_1 + 2) = 0;
  *(undefined2 *)(param_1 + 4) = 0;
  *(undefined2 *)(param_1 + 6) = 0;
  *(undefined2 *)(param_1 + 8) = 0;
  *(undefined4 *)(param_1 + 0xc) = 0;
  *param_1 = 0;
  if ((param_2 == (byte *)0x0) || (*param_2 == 0)) {
    return 1;
  }
  uVar6 = 0xffffffff;
  do {
    pbVar5 = param_2;
    if (uVar6 == 0) break;
    uVar6 = uVar6 - 1;
    pbVar5 = param_2 + 1;
    bVar3 = *param_2;
    param_2 = pbVar5;
  } while (bVar3 != 0);
  uVar6 = ~uVar6;
  pbVar5 = pbVar5 + -uVar6;
  pbVar10 = &local_cc;
  for (uVar7 = uVar6 >> 2; uVar7 != 0; uVar7 = uVar7 - 1) {
    *(undefined4 *)pbVar10 = *(undefined4 *)pbVar5;
    pbVar5 = pbVar5 + 4;
    pbVar10 = pbVar10 + 4;
  }
  for (uVar7 = uVar6 & 3; uVar7 != 0; uVar7 = uVar7 - 1) {
    *pbVar10 = *pbVar5;
    pbVar5 = pbVar5 + 1;
    pbVar10 = pbVar10 + 1;
  }
  if ((local_cc & 4) == 4) {
    *param_1 = 0x45;
  }
  if ((local_cc & 8) == 8) {
    *param_1 = 0x41;
  }
  if ((local_cc & 0x40) == 0x40) {
    *(undefined2 *)(param_1 + 2) = 1;
  }
  if ((local_cc & 0x10) == 0x10) {
    *(undefined2 *)(param_1 + 4) = 1;
  }
  if ((local_cc & 0x20) == 0x20) {
    *(undefined2 *)(param_1 + 6) = 1;
  }
  if ((local_cc & 0x80) == 0x80) {
    *(undefined2 *)(param_1 + 8) = 1;
  }
  if ((local_cc & 1) == 1) {
    pbVar5 = &local_cb;
    param_1[0x10] = 0;
    while (local_cb != 0) {
      uVar6 = 0xffffffff;
      pcVar11 = &DAT_100fe900 + (uint)local_cb * 5;
      do {
        pcVar13 = pcVar11;
        if (uVar6 == 0) break;
        uVar6 = uVar6 - 1;
        pcVar13 = pcVar11 + 1;
        cVar4 = *pcVar11;
        pcVar11 = pcVar13;
      } while (cVar4 != '\0');
      uVar6 = ~uVar6;
      iVar8 = -1;
      pcVar11 = param_1 + 0x10;
      do {
        pcVar12 = pcVar11;
        if (iVar8 == 0) break;
        iVar8 = iVar8 + -1;
        pcVar12 = pcVar11 + 1;
        cVar4 = *pcVar11;
        pcVar11 = pcVar12;
      } while (cVar4 != '\0');
      pcVar11 = pcVar13 + -uVar6;
      pcVar13 = pcVar12 + -1;
      for (uVar7 = uVar6 >> 2; uVar7 != 0; uVar7 = uVar7 - 1) {
        *(undefined4 *)pcVar13 = *(undefined4 *)pcVar11;
        pcVar11 = pcVar11 + 4;
        pcVar13 = pcVar13 + 4;
      }
      pbVar5 = pbVar5 + 1;
      local_cb = *pbVar5;
      for (uVar6 = uVar6 & 3; uVar6 != 0; uVar6 = uVar6 - 1) {
        *pcVar13 = *pcVar11;
        pcVar11 = pcVar11 + 1;
        pcVar13 = pcVar13 + 1;
      }
    }
    param_1[0x155] = 0xff;
    uVar6 = *(int *)(param_1 + 0xc) + 1;
    *(uint *)(param_1 + 0xc) = uVar6;
  }
  else {
    uVar6 = CONCAT31((int3)(uVar6 >> 8),local_cc) & 0xffffff02;
    if ((char)uVar6 == '\x02') {
      param_2 = &local_cb;
      pbVar5 = &local_cb;
      while( true ) {
        bVar3 = *pbVar5;
        uVar6 = CONCAT31((int3)(uVar6 >> 8),bVar3);
        pbVar10 = pbVar5;
        if (bVar3 == 0) break;
        while ((bVar3 != 0xff && (param_2 = pbVar10, (char)uVar6 != '\0'))) {
          bVar3 = pbVar10[1];
          param_2 = pbVar10 + 1;
          pbVar10 = param_2;
          uVar6 = (uint)bVar3;
        }
        if (*pbVar10 == 0xff) {
          *pbVar10 = 0;
          pbVar10 = pbVar10 + 1;
          param_2 = pbVar10;
        }
        bVar3 = *pbVar5;
        pbVar9 = pbVar5;
        while (bVar3 != 0x7c) {
          pbVar1 = pbVar9 + 1;
          pbVar9 = pbVar9 + 1;
          bVar3 = *pbVar1;
        }
        *pbVar9 = 0;
        pbVar9 = pbVar9 + 1;
        iVar8 = 0;
        bVar3 = *pbVar5;
        while (bVar3 != 0) {
          iVar2 = iVar8 + *(int *)(param_1 + 0xc) * 0x14;
          iVar8 = iVar8 + 1;
          pbVar5 = pbVar5 + 1;
          param_1[iVar2 + 0x155] = bVar3 - 1;
          pbVar10 = param_2;
          bVar3 = *pbVar5;
        }
        param_1[iVar8 + *(int *)(param_1 + 0xc) * 0x14 + 0x155] = 0xff;
        param_1[*(int *)(param_1 + 0xc) * 0x41 + 0x10] = 0;
        bVar3 = *pbVar9;
        pbVar5 = pbVar10;
        while (bVar3 != 0) {
          uVar6 = 0xffffffff;
          pcVar11 = &DAT_100fe900 + (uint)bVar3 * 5;
          do {
            pcVar13 = pcVar11;
            if (uVar6 == 0) break;
            uVar6 = uVar6 - 1;
            pcVar13 = pcVar11 + 1;
            cVar4 = *pcVar11;
            pcVar11 = pcVar13;
          } while (cVar4 != '\0');
          uVar6 = ~uVar6;
          iVar8 = -1;
          pcVar11 = param_1 + *(int *)(param_1 + 0xc) * 0x41 + 0x10;
          do {
            pcVar12 = pcVar11;
            if (iVar8 == 0) break;
            iVar8 = iVar8 + -1;
            pcVar12 = pcVar11 + 1;
            cVar4 = *pcVar11;
            pcVar11 = pcVar12;
          } while (cVar4 != '\0');
          pcVar11 = pcVar13 + -uVar6;
          pcVar13 = pcVar12 + -1;
          for (uVar7 = uVar6 >> 2; uVar7 != 0; uVar7 = uVar7 - 1) {
            *(undefined4 *)pcVar13 = *(undefined4 *)pcVar11;
            pcVar11 = pcVar11 + 4;
            pcVar13 = pcVar13 + 4;
          }
          pbVar9 = pbVar9 + 1;
          bVar3 = *pbVar9;
          for (uVar6 = uVar6 & 3; pbVar5 = param_2, uVar6 != 0; uVar6 = uVar6 - 1) {
            *pcVar13 = *pcVar11;
            pcVar11 = pcVar11 + 1;
            pcVar13 = pcVar13 + 1;
          }
        }
        uVar6 = *(int *)(param_1 + 0xc) + 1;
        *(uint *)(param_1 + 0xc) = uVar6;
      }
    }
  }
  return CONCAT22((short)(uVar6 >> 0x10),1);
}



===== 0x10002f10 =====
Function: FUN_10002f10 @ 10002f10

undefined4 __cdecl FUN_10002f10(char *param_1,char *param_2,byte *param_3)

{
  char cVar1;
  undefined2 uVar2;
  undefined4 uVar3;
  undefined2 extraout_var;
  uint uVar4;
  undefined2 extraout_var_00;
  uint uVar5;
  uint uVar6;
  int iVar7;
  int iVar8;
  char *pcVar9;
  char local_7c [32];
  char local_5c [31];
  char acStack_3d [33];
  undefined4 local_1c [5];
  int local_8;
  
  uVar3 = FUN_10002c70(param_2);
  if ((short)uVar3 == 0) {
    return uVar3;
  }
  uVar5 = 0xffffffff;
  pcVar9 = param_2;
  do {
    if (uVar5 == 0) break;
    uVar5 = uVar5 - 1;
    cVar1 = *pcVar9;
    pcVar9 = pcVar9 + 1;
  } while (cVar1 != '\0');
  uVar5 = ~uVar5;
  uVar6 = uVar5 - 1;
  if (uVar6 == 1) {
    if (*(short *)(&DAT_10077d5e + *(short *)(&DAT_1007e388 + *param_2 * 2) * 2) != 0) {
      FUN_10009cd0((int)param_1,param_2,param_3);
      return CONCAT22(extraout_var,1);
    }
  }
  else if (((2 < (int)uVar6) && (*(short *)(&DAT_1007e688 + *param_2 * 2) == 0x4d)) &&
          (*(short *)(&DAT_1007e688 + param_2[1] * 2) == 0x43)) {
    uVar4 = FUN_1000ff10(param_2 + 2);
    if ((short)uVar4 != 0) {
      FUN_1000fed0(local_7c,param_2 + 2);
      FUN_10063ed6(acStack_3d + 1,(byte *)s_MAC_s_10077230);
      goto LAB_10002fd7;
    }
  }
  FUN_1000fed0(acStack_3d + 1,param_2);
  uVar5 = uVar6;
LAB_10002fd7:
  uVar6 = uVar5 - 1;
  iVar7 = 0;
  if (-1 < (int)uVar6) {
    local_8 = 1 - (int)(acStack_3d + 1);
    iVar8 = iVar7;
    do {
      if ((((&DAT_1007e188)[acStack_3d[uVar6 + 1]] & 0xc0) != 0) ||
         (iVar7 = iVar8, acStack_3d[uVar6 + 1] == '\'')) {
        FUN_10002db0(local_1c,(int)(acStack_3d + 1),uVar5,uVar6,(int)local_5c);
        if (((&DAT_1007e188)[acStack_3d[uVar6 + 1]] & 0xc0) == 0) {
          iVar7 = 0x1a;
        }
        else {
          iVar7 = acStack_3d[uVar6 + 1] + -0x41;
        }
        uVar2 = FUN_10001670((int *)(&DAT_100fee0c + iVar7 * 0x1c),(int)local_1c);
        local_5c[iVar8] = (char)uVar2;
        iVar7 = iVar8 + 1;
        if (((acStack_3d[uVar6 + 1] == 'E') && (local_5c[iVar8] == '\x01')) &&
           ((1 < (int)uVar5 &&
            (((((int)uVar5 < 4 && (0 < (int)uVar6)) &&
              ((int)(acStack_3d + local_8 + uVar6 + 1) < (int)uVar5)) &&
             ((*(short *)(&DAT_10077d9e + acStack_3d[uVar6] * 2) == 0 &&
              (*(short *)(&DAT_10077d9e + acStack_3d[uVar6 + 2] * 2) == 0)))))))) {
          local_5c[iVar8] = '\x1c';
        }
      }
      uVar6 = uVar6 - 1;
      iVar8 = iVar7;
    } while (uVar6 < 0x80000000);
  }
  local_5c[iVar7] = '\0';
  FUN_10010050(local_5c,0,iVar7 + -1);
  FUN_100024f0(param_1,local_5c);
  uVar5 = FUN_1000ff60(param_1);
  if ((short)uVar5 == 0) {
    FUN_10009cd0((int)param_1,param_2,param_3);
    return CONCAT22(extraout_var_00,1);
  }
  *param_3 = *param_3 | 0x10;
  return CONCAT22((short)((uint)param_3 >> 0x10),1);
}



