===== 0x1001c900 =====
Function: FUN_1001c900 @ 1001c900

undefined4 __cdecl FUN_1001c900(byte *param_1)

{
  byte bVar1;
  int iVar2;
  
  iVar2 = FUN_1001c960((int)param_1,2);
  if (((iVar2 != 0) && (bVar1 = *param_1, (bVar1 & 0x80) != 0)) &&
     ((((0xa0 < bVar1 && (((bVar1 < 0xae && (0xa0 < param_1[1])) && (param_1[1] != 0xff)))) ||
       ((bVar1 == 0xfd && (param_1[1] == 0xfe)))) ||
      ((bVar1 == 0xae && ((0xa0 < param_1[1] && (param_1[1] < 0xc3)))))))) {
    return 1;
  }
  return 0;
}



===== 0x1002e990 =====
Function: FUN_1002e990 @ 1002e990

undefined4 __cdecl
FUN_1002e990(int param_1,int param_2,int *param_3,undefined4 param_4,undefined4 *param_5,
            undefined4 param_6,undefined4 param_7)

{
  int *piVar1;
  int iVar2;
  undefined4 uVar3;
  int *piVar4;
  int iVar5;
  int iVar6;
  int *piVar7;
  int local_8;
  
  piVar4 = param_3;
  iVar6 = 0;
  iVar2 = *param_3;
  param_3 = (int *)0x0;
  uVar3 = *param_5;
  local_8 = -1;
  piVar7 = (int *)(param_1 + 4);
  if (*piVar7 != 0) {
    do {
      iVar5 = FUN_1002ea80(*piVar4 + param_2,*piVar7);
      if (((iVar5 != 0) &&
          (iVar5 = FUN_1001c320((byte *)(*piVar4 + param_2),(byte *)piVar7[-1],*piVar7), iVar5 == 0)
          ) && ((int)param_3 < *piVar7)) {
        param_3 = (int *)*piVar7;
        local_8 = iVar6;
      }
      piVar1 = piVar7 + 6;
      piVar7 = piVar7 + 6;
      iVar6 = iVar6 + 1;
    } while (*piVar1 != 0);
    if (-1 < local_8) {
      iVar6 = param_1 + local_8 * 0x18;
      iVar6 = (**(code **)(iVar6 + 0x10))(param_2,piVar4,param_4,param_5,param_6,iVar6,param_7);
      if (iVar6 == 0) {
        *piVar4 = iVar2;
        *param_5 = uVar3;
        return 0;
      }
      return 1;
    }
  }
  *piVar4 = iVar2;
  *param_5 = uVar3;
  return 0;
}



===== 0x1002efd0 =====
Function: FUN_1002efd0 @ 1002efd0

undefined4 __cdecl
FUN_1002efd0(int param_1,int *param_2,int param_3,int *param_4,int param_5,int param_6,int param_7)

{
  char cVar1;
  int *piVar2;
  int iVar3;
  uint uVar4;
  int iVar5;
  int iVar6;
  char *pcVar7;
  char local_c04 [1536];
  char local_604 [1536];
  
  piVar2 = *(int **)(param_7 + 0x122448);
  iVar3 = FUN_1002eab0((undefined4 *)local_c04,param_1,param_2,param_6);
  if (iVar3 != 0) {
    uVar4 = 0xffffffff;
    pcVar7 = local_c04;
    do {
      if (uVar4 == 0) break;
      uVar4 = uVar4 - 1;
      cVar1 = *pcVar7;
      pcVar7 = pcVar7 + 1;
    } while (cVar1 != '\0');
    iVar6 = (~uVar4 - 1) + *(int *)(param_6 + 0xc) + *(int *)(param_6 + 4);
    iVar3 = FUN_10016150(local_c04,*(int **)(DAT_100a0460 + 0x20428),0);
    if (iVar3 != 0) {
      if (param_5 != 0) {
        *(undefined1 *)(*param_4 + param_3) = 0xa2;
        *(undefined1 *)(*param_4 + 1 + param_3) = 0xfe;
        iVar3 = *param_2;
        if (iVar3 < *(int *)(param_7 + 0x477a8)) {
          param_6 = *(int *)(*(int *)(param_7 + 0x4779c) + iVar3 * 4);
          param_5 = *(int *)(*(int *)(param_7 + 0x477a0) + -4 + (iVar3 + iVar6) * 4);
          iVar3 = 0;
        }
        else {
          param_5 = *(int *)(*(int *)(param_7 + 0x4779c) + -4 + *(int *)(param_7 + 0x477a8) * 4) + 1
          ;
          iVar3 = 0;
          param_6 = param_5;
        }
        do {
          *(int *)(*(int *)(param_7 + 0x47794) + (*param_4 + iVar3) * 4) = param_6;
          iVar5 = *param_4 + iVar3;
          iVar3 = iVar3 + 1;
          *(int *)(*(int *)(param_7 + 0x47798) + iVar5 * 4) = param_5;
        } while (iVar3 < 2);
        *(int *)(*piVar2 * 0x210 + 8 + piVar2[1]) = *param_4 + 1;
        *(undefined1 *)(*piVar2 * 0x210 + 0xc + piVar2[1]) = DAT_1009f948;
        *(undefined1 *)(*piVar2 * 0x210 + 0x20c + piVar2[1]) = 2;
      }
      *piVar2 = *piVar2 + 1;
      *param_4 = *param_4 + 2;
      *param_2 = *param_2 + iVar6;
      return 1;
    }
    iVar3 = FUN_10016150(local_c04,*(int **)(DAT_100a0460 + 0x20428),4);
    if (iVar3 != 0) {
      local_604[0] = '\0';
      iVar3 = FUN_1002ebf0(local_604,local_c04,&DAT_1007fdc8,*(char **)(param_6 + 0x14));
      if ((iVar3 != 0) ||
         (iVar3 = FUN_1002ebf0(local_604,local_c04,&DAT_1007fdc0,*(char **)(param_6 + 0x14)),
         iVar3 != 0)) {
        if (param_5 != 0) {
          *(undefined1 *)(*param_4 + param_3) = 0xa2;
          *(undefined1 *)(*param_4 + 1 + param_3) = 0xfe;
          iVar3 = *param_2;
          if (iVar3 < *(int *)(param_7 + 0x477a8)) {
            param_6 = *(int *)(*(int *)(param_7 + 0x4779c) + iVar3 * 4);
            param_5 = *(int *)(*(int *)(param_7 + 0x477a0) + -4 + (iVar3 + iVar6) * 4);
            iVar3 = 0;
          }
          else {
            param_5 = *(int *)(*(int *)(param_7 + 0x4779c) + -4 + *(int *)(param_7 + 0x477a8) * 4) +
                      1;
            iVar3 = 0;
            param_6 = param_5;
          }
          do {
            *(int *)(*(int *)(param_7 + 0x47794) + (*param_4 + iVar3) * 4) = param_6;
            iVar5 = *param_4 + iVar3;
            iVar3 = iVar3 + 1;
            *(int *)(*(int *)(param_7 + 0x47798) + iVar5 * 4) = param_5;
          } while (iVar3 < 2);
          *(int *)(*piVar2 * 0x210 + 8 + piVar2[1]) = *param_4 + 1;
          uVar4 = 0xffffffff;
          pcVar7 = local_604;
          do {
            if (uVar4 == 0) break;
            uVar4 = uVar4 - 1;
            cVar1 = *pcVar7;
            pcVar7 = pcVar7 + 1;
          } while (cVar1 != '\0');
          if (~uVar4 - 1 < 0x200) {
            uVar4 = 0xffffffff;
            pcVar7 = local_604;
            do {
              if (uVar4 == 0) break;
              uVar4 = uVar4 - 1;
              cVar1 = *pcVar7;
              pcVar7 = pcVar7 + 1;
            } while (cVar1 != '\0');
            FUN_10063f30((undefined4 *)(*piVar2 * 0x210 + 0xc + piVar2[1]),(undefined4 *)local_604,
                         ~uVar4 - 1);
            uVar4 = 0xffffffff;
            pcVar7 = local_604;
            do {
              if (uVar4 == 0) break;
              uVar4 = uVar4 - 1;
              cVar1 = *pcVar7;
              pcVar7 = pcVar7 + 1;
            } while (cVar1 != '\0');
            *(undefined1 *)(~uVar4 + *piVar2 * 0x210 + 0xb + piVar2[1]) = 0;
            *(undefined1 *)(*piVar2 * 0x210 + 0x20c + piVar2[1]) = 1;
          }
          else {
            FUN_10063f30((undefined4 *)(*piVar2 * 0x210 + 0xc + piVar2[1]),(undefined4 *)local_604,
                         0x1ff);
            *(undefined1 *)(*piVar2 * 0x210 + 0x20b + piVar2[1]) = 0;
            *(undefined1 *)(*piVar2 * 0x210 + 0x20c + piVar2[1]) = 3;
          }
        }
        *piVar2 = *piVar2 + 1;
        *param_4 = *param_4 + 2;
        *param_2 = *param_2 + iVar6;
        return 1;
      }
    }
  }
  return 0;
}



===== 0x1002eab0 =====
Function: FUN_1002eab0 @ 1002eab0

undefined4 __cdecl FUN_1002eab0(undefined4 *param_1,int param_2,int *param_3,int param_4)

{
  char *pcVar1;
  
  pcVar1 = _strstr((char *)(*param_3 + param_2),*(char **)(param_4 + 8));
  if (pcVar1 == (char *)0x0) {
    *(undefined1 *)param_1 = DAT_1009f948;
    return 0;
  }
  pcVar1 = pcVar1 + ((-*param_3 - *(int *)(param_4 + 4)) - param_2);
  if (0x5ff < (int)pcVar1) {
    *(undefined1 *)param_1 = DAT_1009f948;
    return 0;
  }
  FUN_10063f30(param_1,(undefined4 *)(*(int *)(param_4 + 4) + *param_3 + param_2),(uint)pcVar1);
  pcVar1[(int)param_1] = '\0';
  return 1;
}



===== 0x1002ebf0 =====
Function: FUN_1002ebf0 @ 1002ebf0

undefined4 __cdecl FUN_1002ebf0(char *param_1,char *param_2,byte *param_3,char *param_4)

{
  char cVar1;
  byte bVar2;
  byte *pbVar3;
  char *pcVar4;
  byte *pbVar5;
  int iVar6;
  char *pcVar7;
  uint uVar8;
  uint uVar9;
  char *pcVar10;
  byte *pbVar11;
  char local_604 [1536];
  
  pcVar7 = param_4;
  pbVar3 = param_3;
  uVar8 = 0xffffffff;
  pcVar4 = param_2;
  do {
    pcVar10 = pcVar4;
    if (uVar8 == 0) break;
    uVar8 = uVar8 - 1;
    pcVar10 = pcVar4 + 1;
    cVar1 = *pcVar4;
    pcVar4 = pcVar10;
  } while (cVar1 != '\0');
  uVar8 = ~uVar8;
  pcVar4 = pcVar10 + -uVar8;
  pcVar10 = local_604;
  for (uVar9 = uVar8 >> 2; uVar9 != 0; uVar9 = uVar9 - 1) {
    *(undefined4 *)pcVar10 = *(undefined4 *)pcVar4;
    pcVar4 = pcVar4 + 4;
    pcVar10 = pcVar10 + 4;
  }
  for (uVar8 = uVar8 & 3; uVar8 != 0; uVar8 = uVar8 - 1) {
    *pcVar10 = *pcVar4;
    pcVar4 = pcVar4 + 1;
    pcVar10 = pcVar10 + 1;
  }
  pcVar4 = FUN_1002eb30(local_604,param_3,param_4);
  if (pcVar4 == (char *)0x0) {
    return 0;
  }
  pbVar5 = (byte *)FUN_10024b50(pcVar4,pcVar7,&param_2);
  uVar8 = 0xffffffff;
  pbVar11 = pbVar3;
  do {
    if (uVar8 == 0) break;
    uVar8 = uVar8 - 1;
    bVar2 = *pbVar11;
    pbVar11 = pbVar11 + 1;
  } while (bVar2 != 0);
  iVar6 = FUN_1001c320(pbVar5,pbVar3,~uVar8 - 1);
  if (iVar6 != 0) {
    return 0;
  }
  pcVar7 = FUN_1002ece0(param_2,pcVar7);
  if (pcVar7 == (char *)0x0) {
    return 0;
  }
  pcVar7 = FUN_10024b50(pcVar7,&DAT_1007d4f8,&param_2);
  if (pcVar7 == (char *)0x0) {
    return 0;
  }
  uVar8 = 0xffffffff;
  do {
    pcVar4 = pcVar7;
    if (uVar8 == 0) break;
    uVar8 = uVar8 - 1;
    pcVar4 = pcVar7 + 1;
    cVar1 = *pcVar7;
    pcVar7 = pcVar4;
  } while (cVar1 != '\0');
  uVar8 = ~uVar8;
  pcVar7 = pcVar4 + -uVar8;
  for (uVar9 = uVar8 >> 2; uVar9 != 0; uVar9 = uVar9 - 1) {
    *(undefined4 *)param_1 = *(undefined4 *)pcVar7;
    pcVar7 = pcVar7 + 4;
    param_1 = param_1 + 4;
  }
  for (uVar8 = uVar8 & 3; uVar8 != 0; uVar8 = uVar8 - 1) {
    *param_1 = *pcVar7;
    pcVar7 = pcVar7 + 1;
    param_1 = param_1 + 1;
  }
  return 1;
}



