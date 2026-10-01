===== 0x10011820 =====
Function: FUN_10011820 @ 10011820

undefined4 __cdecl FUN_10011820(byte *param_1,byte *param_2,int param_3,undefined4 param_4)

{
  char cVar1;
  byte bVar2;
  uint *puVar3;
  byte *pbVar4;
  int iVar5;
  uint uVar6;
  uint uVar7;
  int iVar8;
  char *pcVar9;
  char *pcVar10;
  byte *pbVar11;
  bool bVar12;
  byte local_dc [200];
  int local_14;
  int local_10;
  int local_c;
  int local_8;
  
  iVar5 = param_3;
  if ((char)param_4 == '\0') {
    local_c = DAT_1009f9a8;
    local_8 = DAT_1009f950;
    local_10 = DAT_1009f9c0;
    iVar8 = DAT_1009f9a0;
    uVar6 = DAT_1009f990;
  }
  else {
    local_c = DAT_1009f968;
    local_8 = DAT_1009f954;
    local_10 = DAT_1009f94c;
    iVar8 = DAT_1009f960;
    uVar6 = DAT_1009f994;
  }
  if (param_3 == 0) {
    uVar7 = FUN_10012030(param_1,*(int *)(param_1 + 200) - 1,(char)param_4);
    if (uVar7 < uVar6) {
      FUN_10025440(iVar8,uVar7 * 2,0,(char *)&param_3,1,2);
      if ((*(int *)(param_1 + 200) + -1 == (int)(char)param_3) && (*param_1 == param_3._1_1_)) {
        FUN_10025440(local_c,uVar7 * 4 + 4,0,(char *)&local_14,4,1);
        FUN_10025440(local_8,local_14,0,(char *)local_dc,1,local_10);
        pbVar4 = local_dc;
        do {
          bVar2 = *pbVar4;
          bVar12 = bVar2 < *param_1;
          if (bVar2 != *param_1) {
LAB_100119d3:
            iVar5 = (1 - (uint)bVar12) - (uint)(bVar12 != 0);
            goto LAB_100119d8;
          }
          if (bVar2 == 0) break;
          bVar2 = pbVar4[1];
          bVar12 = bVar2 < param_1[1];
          if (bVar2 != param_1[1]) goto LAB_100119d3;
          pbVar4 = pbVar4 + 2;
          param_1 = param_1 + 2;
        } while (bVar2 != 0);
        iVar5 = 0;
LAB_100119d8:
        if (iVar5 == 0) {
          uVar6 = 0xffffffff;
          pbVar4 = local_dc;
          do {
            if (uVar6 == 0) break;
            uVar6 = uVar6 - 1;
            bVar2 = *pbVar4;
            pbVar4 = pbVar4 + 1;
          } while (bVar2 != 0);
          uVar7 = 0xffffffff;
          pbVar4 = local_dc + ~uVar6;
          do {
            pbVar11 = pbVar4;
            if (uVar7 == 0) break;
            uVar7 = uVar7 - 1;
            pbVar11 = pbVar4 + 1;
            bVar2 = *pbVar4;
            pbVar4 = pbVar11;
          } while (bVar2 != 0);
          uVar7 = ~uVar7;
          pbVar4 = pbVar11 + -uVar7;
          pbVar11 = param_2;
          for (uVar6 = uVar7 >> 2; uVar6 != 0; uVar6 = uVar6 - 1) {
            *(undefined4 *)pbVar11 = *(undefined4 *)pbVar4;
            pbVar4 = pbVar4 + 4;
            pbVar11 = pbVar11 + 4;
          }
          for (uVar7 = uVar7 & 3; uVar7 != 0; uVar7 = uVar7 - 1) {
            *pbVar11 = *pbVar4;
            pbVar4 = pbVar4 + 1;
            pbVar11 = pbVar11 + 1;
          }
          uVar6 = 0xffffffff;
          pbVar4 = param_2;
          do {
            if (uVar6 == 0) break;
            uVar6 = uVar6 - 1;
            bVar2 = *pbVar4;
            pbVar4 = pbVar4 + 1;
          } while (bVar2 != 0);
          *(uint *)(param_2 + 200) = ~uVar6 - 1;
          return 1;
        }
      }
      param_2[200] = 0;
      param_2[0xc9] = 0;
      param_2[0xca] = 0;
      param_2[0xcb] = 0;
      *param_2 = 0;
      return 0xffffffff;
    }
  }
  else {
    puVar3 = (uint *)FUN_10011a50(param_3,param_1);
    puVar3 = FUN_10011ac0(iVar5,puVar3,param_1,iVar8,local_c,local_8,uVar6,local_10,param_4);
    if (*(char *)puVar3[2] != '\0') {
      uVar6 = 0xffffffff;
      pcVar10 = (char *)puVar3[2];
      do {
        pcVar9 = pcVar10;
        if (uVar6 == 0) break;
        uVar6 = uVar6 - 1;
        pcVar9 = pcVar10 + 1;
        cVar1 = *pcVar10;
        pcVar10 = pcVar9;
      } while (cVar1 != '\0');
      uVar6 = ~uVar6;
      pbVar4 = (byte *)(pcVar9 + -uVar6);
      pbVar11 = param_2;
      for (uVar7 = uVar6 >> 2; uVar7 != 0; uVar7 = uVar7 - 1) {
        *(undefined4 *)pbVar11 = *(undefined4 *)pbVar4;
        pbVar4 = pbVar4 + 4;
        pbVar11 = pbVar11 + 4;
      }
      for (uVar6 = uVar6 & 3; uVar6 != 0; uVar6 = uVar6 - 1) {
        *pbVar11 = *pbVar4;
        pbVar4 = pbVar4 + 1;
        pbVar11 = pbVar11 + 1;
      }
      uVar6 = 0xffffffff;
      pcVar10 = (char *)puVar3[2];
      do {
        if (uVar6 == 0) break;
        uVar6 = uVar6 - 1;
        cVar1 = *pcVar10;
        pcVar10 = pcVar10 + 1;
      } while (cVar1 != '\0');
      *(uint *)(param_2 + 200) = ~uVar6 - 1;
      return 1;
    }
  }
  *param_2 = 0;
  param_2[200] = 0;
  param_2[0xc9] = 0;
  param_2[0xca] = 0;
  param_2[0xcb] = 0;
  return 0xffffffff;
}



===== 0x1003a610 =====
Function: FUN_1003a610 @ 1003a610

void __cdecl FUN_1003a610(char *param_1,int param_2,char param_3)

{
  char cVar1;
  byte bVar2;
  int iVar3;
  uint uVar4;
  uint uVar5;
  char *pcVar6;
  byte *pbVar7;
  byte *pbVar8;
  char *pcVar9;
  int iVar10;
  byte local_a4 [152];
  int local_c;
  int local_8;
  
  local_8 = 0;
  if (param_3 == 'S') {
    uVar4 = 0xffffffff;
    do {
      pcVar9 = param_1;
      if (uVar4 == 0) break;
      uVar4 = uVar4 - 1;
      pcVar9 = param_1 + 1;
      cVar1 = *param_1;
      param_1 = pcVar9;
    } while (cVar1 != '\0');
    uVar4 = ~uVar4;
    pbVar7 = (byte *)(pcVar9 + -uVar4);
    pbVar8 = local_a4;
    for (uVar5 = uVar4 >> 2; uVar5 != 0; uVar5 = uVar5 - 1) {
      *(undefined4 *)pbVar8 = *(undefined4 *)pbVar7;
      pbVar7 = pbVar7 + 4;
      pbVar8 = pbVar8 + 4;
    }
    for (uVar4 = uVar4 & 3; uVar4 != 0; uVar4 = uVar4 - 1) {
      *pbVar8 = *pbVar7;
      pbVar7 = pbVar7 + 1;
      pbVar8 = pbVar8 + 1;
    }
  }
  else {
    uVar4 = 0xffffffff;
    iVar3 = 0;
    pcVar9 = param_1;
    do {
      if (uVar4 == 0) break;
      uVar4 = uVar4 - 1;
      cVar1 = *pcVar9;
      pcVar9 = pcVar9 + 1;
    } while (cVar1 != '\0');
    iVar10 = iVar3;
    if (0 < (int)(~uVar4 - 1)) {
      do {
        iVar3 = iVar10 + 1;
        local_a4[iVar10] = (&DAT_1007e388)[param_1[iVar10] * 2];
        iVar10 = iVar3;
      } while (iVar3 < (int)(~uVar4 - 1));
    }
    local_a4[iVar3] = 0;
  }
  uVar4 = 0xffffffff;
  pbVar7 = local_a4;
  do {
    if (uVar4 == 0) break;
    uVar4 = uVar4 - 1;
    bVar2 = *pbVar7;
    pbVar7 = pbVar7 + 1;
  } while (bVar2 != 0);
  iVar10 = 0;
  local_c = ~uVar4 - 2;
  iVar3 = 0;
  if (0 < (int)(~uVar4 - 2)) {
    pbVar7 = local_a4;
    do {
      param_1 = (char *)0x0;
      _param_3 = (char *)0x7e;
      pcVar9 = param_1;
      while (param_1 = pcVar9, 1 < (int)_param_3 - (int)param_1) {
        pcVar6 = (char *)((int)(param_1 + (int)_param_3) / 2);
        iVar3 = FUN_1001c280(pbVar7,&DAT_10081568 + (int)pcVar6 * 3,2);
        if (iVar3 == 0) {
          *(byte *)(local_8 + param_2) = (char)pcVar6 + 1U | 0x80;
          iVar10 = iVar10 + 1;
          local_8 = local_8 + 1;
          pbVar7 = pbVar7 + 1;
          goto LAB_1003a781;
        }
        pcVar9 = pcVar6;
        if (iVar3 < 0) {
          pcVar9 = param_1;
          _param_3 = pcVar6;
        }
      }
      iVar3 = FUN_1001c280(pbVar7,&DAT_10081568 + (int)_param_3 * 3,2);
      if (iVar3 == 0) {
        *(byte *)(local_8 + param_2) = (char)_param_3 + 1U | 0x80;
        iVar10 = iVar10 + 1;
        local_8 = local_8 + 1;
        pbVar7 = pbVar7 + 1;
      }
      else {
        iVar3 = FUN_1001c280(pbVar7,&DAT_10081568 + (int)param_1 * 3,2);
        if (iVar3 == 0) {
          *(byte *)(local_8 + param_2) = (char)param_1 + 1U | 0x80;
          iVar10 = iVar10 + 1;
          local_8 = local_8 + 1;
          pbVar7 = pbVar7 + 1;
        }
        else {
          *(byte *)(local_8 + param_2) = *pbVar7;
          local_8 = local_8 + 1;
        }
      }
LAB_1003a781:
      iVar10 = iVar10 + 1;
      pbVar7 = pbVar7 + 1;
      iVar3 = local_8;
    } while (iVar10 < local_c);
  }
  *(byte *)(iVar3 + param_2) = local_a4[iVar10];
  *(undefined1 *)(iVar3 + 1 + param_2) = 0;
  return;
}



===== 0x10064645 =====
Function: FUN_10064645 @ 10064645

void __thiscall FUN_10064645(void *this,byte *param_1)

{
  FUN_100645ba(this,param_1);
  return;
}



