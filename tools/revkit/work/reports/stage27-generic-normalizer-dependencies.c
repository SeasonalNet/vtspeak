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



===== 0x1000fd70 =====
Function: FUN_1000fd70 @ 1000fd70

void __cdecl FUN_1000fd70(char *param_1,int param_2)

{
  char cVar1;
  int iVar2;
  int iVar3;
  int iVar4;
  uint uVar5;
  int iVar6;
  char *pcVar7;
  char local_a8 [152];
  int local_10;
  int local_c;
  int local_8;
  
  uVar5 = 0xffffffff;
  iVar2 = 0;
  iVar6 = 0;
  pcVar7 = param_1;
  do {
    if (uVar5 == 0) break;
    uVar5 = uVar5 - 1;
    cVar1 = *pcVar7;
    pcVar7 = pcVar7 + 1;
  } while (cVar1 != '\0');
  iVar3 = iVar2;
  if (0 < (int)(~uVar5 - 1)) {
    do {
      iVar2 = iVar3 + 1;
      local_a8[iVar3] = (&DAT_1007e388)[param_1[iVar3] * 2];
      iVar3 = iVar2;
    } while (iVar2 < (int)(~uVar5 - 1));
  }
  local_a8[iVar2] = '\0';
  uVar5 = 0xffffffff;
  pcVar7 = local_a8;
  do {
    if (uVar5 == 0) break;
    uVar5 = uVar5 - 1;
    cVar1 = *pcVar7;
    pcVar7 = pcVar7 + 1;
  } while (cVar1 != '\0');
  param_1 = (char *)0x0;
  local_10 = ~uVar5 - 2;
  if (0 < (int)(~uVar5 - 2)) {
    pcVar7 = local_a8;
    do {
      local_8 = 0;
      local_c = 0x7e;
      iVar2 = local_8;
      while (local_8 = iVar2, iVar2 = local_c, 1 < local_c - local_8) {
        iVar2 = (local_8 + local_c) / 2;
        iVar3 = _strncmp(pcVar7,&DAT_10077b00 + iVar2 * 3,2);
        if (iVar3 == 0) goto LAB_1000fe49;
        if (iVar3 < 0) {
          local_c = iVar2;
          iVar2 = local_8;
        }
      }
      iVar4 = _strncmp(pcVar7,&DAT_10077b00 + local_c * 3,2);
      iVar3 = local_8;
      if (iVar4 == 0) {
LAB_1000fe49:
        *(byte *)(iVar6 + param_2) = (char)iVar2 + 1U | 0x80;
        param_1 = param_1 + 1;
        pcVar7 = pcVar7 + 1;
      }
      else {
        iVar2 = _strncmp(pcVar7,&DAT_10077b00 + local_8 * 3,2);
        if (iVar2 == 0) {
          *(byte *)(iVar6 + param_2) = (char)iVar3 + 1U | 0x80;
          param_1 = param_1 + 1;
          pcVar7 = pcVar7 + 1;
        }
        else {
          *(char *)(iVar6 + param_2) = *pcVar7;
        }
      }
      iVar6 = iVar6 + 1;
      param_1 = param_1 + 1;
      pcVar7 = pcVar7 + 1;
    } while ((int)param_1 < local_10);
  }
  *(char *)(iVar6 + param_2) = local_a8[(int)param_1];
  *(undefined1 *)(iVar6 + 1 + param_2) = 0;
  return;
}



===== 0x10002c70 =====
Function: FUN_10002c70 @ 10002c70

uint __cdecl FUN_10002c70(char *param_1)

{
  char cVar1;
  uint uVar2;
  uint uVar3;
  int iVar4;
  int iVar5;
  int iVar6;
  int iVar7;
  char *pcVar8;
  char local_28 [32];
  undefined4 local_8;
  
  uVar3 = 0xffffffff;
  uVar2 = 0;
  iVar7 = 0;
  pcVar8 = param_1;
  do {
    if (uVar3 == 0) break;
    uVar3 = uVar3 - 1;
    cVar1 = *pcVar8;
    pcVar8 = pcVar8 + 1;
  } while (cVar1 != '\0');
  iVar4 = ~uVar3 - 1;
  local_8 = 0;
  iVar5 = 0;
  if (iVar4 < 1) goto LAB_10002d4e;
LAB_10002c97:
  cVar1 = param_1[iVar7];
  uVar2 = CONCAT31((int3)((uint)param_1 >> 8),cVar1);
  if ((((&DAT_1007e188)[cVar1] & 0xc0) == 0) && (cVar1 != '\'')) goto LAB_10002d4e;
  if (((&DAT_1007e188)[cVar1] & 0xc0) != 0) {
    if (*(short *)(&DAT_10077d5e + *(short *)(&DAT_1007e388 + cVar1 * 2) * 2) == 0) {
      local_28[iVar5] = cVar1;
      iVar6 = iVar5 + 1;
      if ((iVar7 == ~uVar3 - 2) && (1 < iVar6)) {
        local_28[iVar5 + 1] = '\0';
        uVar2 = FUN_100026f0(local_28,iVar6,iVar4,(-(ushort)(iVar7 != iVar5) & 0xfff9) + 8);
        if ((short)uVar2 != 0) goto LAB_10002d3c;
        goto LAB_10002d4e;
      }
      goto LAB_10002d3e;
    }
    local_8 = 1;
    if (iVar5 < 2) goto LAB_10002d3c;
    local_28[iVar5] = '\0';
    uVar2 = FUN_100026f0(local_28,iVar5,iVar4,(-(ushort)(iVar7 != iVar5) & 0xfffe) + 4);
    if ((short)uVar2 == 0) goto LAB_10002d4e;
  }
LAB_10002d3c:
  iVar6 = 0;
LAB_10002d3e:
  iVar7 = iVar7 + 1;
  iVar5 = iVar6;
  if (iVar4 <= iVar7) goto code_r0x10002d47;
  goto LAB_10002c97;
code_r0x10002d47:
  if ((short)local_8 != 0) {
    return CONCAT22((short)(uVar2 >> 0x10),1);
  }
LAB_10002d4e:
  return uVar2 & 0xffff0000;
}



===== 0x1000fed0 =====
Function: FUN_1000fed0 @ 1000fed0

void __cdecl FUN_1000fed0(char *param_1,char *param_2)

{
  char cVar1;
  
  cVar1 = *param_2;
  if (cVar1 == '\0') {
    *param_1 = '\0';
    return;
  }
  do {
    if (((&DAT_1007e188)[cVar1] & 0x40) != 0) {
      cVar1 = (char)*(undefined2 *)(&DAT_1007e688 + cVar1 * 2);
    }
    param_2 = param_2 + 1;
    *param_1 = cVar1;
    param_1 = param_1 + 1;
    cVar1 = *param_2;
  } while (cVar1 != '\0');
  *param_1 = '\0';
  return;
}



===== 0x1000ff10 =====
Function: FUN_1000ff10 @ 1000ff10

uint __cdecl FUN_1000ff10(char *param_1)

{
  int iVar1;
  char cVar2;
  uint in_EAX;
  int iVar3;
  
  if (param_1 != (char *)0x0) {
    cVar2 = *param_1;
    in_EAX = CONCAT31((int3)(in_EAX >> 8),cVar2);
    while (cVar2 != '\0') {
      iVar3 = (int)(char)in_EAX;
      if (((&DAT_1007e188)[iVar3] & 0xc0) != 0) {
        iVar1 = iVar3 * 2;
        iVar3 = (int)*(short *)(&DAT_1007e388 + iVar1);
        if (*(short *)(&DAT_10077d5e + iVar3 * 2) != 0) {
          return CONCAT22(*(short *)(&DAT_1007e388 + iVar1) >> 0xf,1);
        }
      }
      cVar2 = param_1[1];
      in_EAX = CONCAT31((int3)((uint)iVar3 >> 8),cVar2);
      param_1 = param_1 + 1;
    }
  }
  return in_EAX & 0xffff0000;
}



===== 0x10002db0 =====
Function: FUN_10002db0 @ 10002db0

void __cdecl FUN_10002db0(undefined4 *param_1,int param_2,int param_3,int param_4,int param_5)

{
  undefined2 uVar1;
  short sVar2;
  uint uVar3;
  
  *param_1 = 0xffffffff;
  param_1[1] = 0xffffffff;
  param_1[2] = 0xffffffff;
  param_1[3] = 0xffffffff;
  param_1[4] = 0xffffffff;
  if (param_4 + -3 < 0) {
    uVar1 = 0;
  }
  else {
    uVar3 = FUN_10002d70(*(char *)(param_2 + -3 + param_4));
    uVar1 = (undefined2)uVar3;
  }
  *(undefined2 *)param_1 = uVar1;
  if (param_4 + -2 < 0) {
    uVar1 = 0;
  }
  else {
    uVar3 = FUN_10002d70(*(char *)(param_2 + -2 + param_4));
    uVar1 = (undefined2)uVar3;
  }
  *(undefined2 *)((int)param_1 + 2) = uVar1;
  if (param_4 + -1 < 0) {
    uVar1 = 0;
  }
  else {
    uVar3 = FUN_10002d70(*(char *)(param_2 + -1 + param_4));
    uVar1 = (undefined2)uVar3;
  }
  *(undefined2 *)(param_1 + 1) = uVar1;
  uVar3 = FUN_10002d70(*(char *)(param_2 + param_4));
  *(short *)((int)param_1 + 6) = (short)uVar3;
  if (param_4 + 1 < param_3) {
    uVar3 = FUN_10002d70(*(char *)(param_2 + 1 + param_4));
    uVar1 = (undefined2)uVar3;
  }
  else {
    uVar1 = 0;
  }
  *(undefined2 *)(param_1 + 2) = uVar1;
  if (param_4 + 2 < param_3) {
    uVar3 = FUN_10002d70(*(char *)(param_2 + 2 + param_4));
    uVar1 = (undefined2)uVar3;
  }
  else {
    uVar1 = 0;
  }
  *(undefined2 *)((int)param_1 + 10) = uVar1;
  if (param_4 + 3 < param_3) {
    uVar3 = FUN_10002d70(*(char *)(param_2 + 3 + param_4));
    uVar1 = (undefined2)uVar3;
  }
  else {
    uVar1 = 0;
  }
  *(undefined2 *)(param_1 + 3) = uVar1;
  if (param_4 + 1 < param_3) {
    sVar2 = (short)*(char *)((param_3 - param_4) + -2 + param_5);
  }
  else {
    sVar2 = 0;
  }
  *(short *)((int)param_1 + 0xe) = sVar2;
  if (param_4 + 2 < param_3) {
    sVar2 = (short)*(char *)((param_3 - param_4) + -3 + param_5);
  }
  else {
    sVar2 = 0;
  }
  *(short *)(param_1 + 4) = sVar2;
  if (param_4 + 3 < param_3) {
    *(short *)((int)param_1 + 0x12) = (short)*(char *)((param_3 - param_4) + -4 + param_5);
    return;
  }
  *(undefined2 *)((int)param_1 + 0x12) = 0;
  return;
}



===== 0x10010050 =====
Function: FUN_10010050 @ 10010050

undefined4 __cdecl FUN_10010050(char *param_1,int param_2,int param_3)

{
  char cVar1;
  uint in_EAX;
  int iVar2;
  char *pcVar3;
  
  if (-1 < param_2) {
    iVar2 = -1;
    pcVar3 = param_1;
    do {
      if (iVar2 == 0) break;
      iVar2 = iVar2 + -1;
      cVar1 = *pcVar3;
      pcVar3 = pcVar3 + 1;
    } while (cVar1 != '\0');
    in_EAX = param_3;
    if (param_3 <= (char)(~(byte)iVar2 - 1)) {
      if (param_2 < param_3) {
        do {
          cVar1 = param_1[param_2];
          param_1[param_2] = param_1[param_3];
          param_1[param_3] = cVar1;
          param_2 = param_2 + 1;
          param_3 = param_3 + -1;
        } while (param_2 < param_3);
      }
      return CONCAT22((short)((uint)param_3 >> 0x10),1);
    }
  }
  return in_EAX & 0xffff0000;
}



===== 0x100024f0 =====
Function: FUN_100024f0 @ 100024f0

uint __cdecl FUN_100024f0(char *param_1,char *param_2)

{
  char cVar1;
  char *pcVar2;
  uint uVar3;
  
  cVar1 = *param_2;
  pcVar2 = param_1;
  do {
    if (cVar1 == '\0') {
      *pcVar2 = '\0';
      uVar3 = FUN_10002200(param_1);
      if ((short)uVar3 != 0) {
        return CONCAT22((short)(uVar3 >> 0x10),1);
      }
      *param_1 = (char)uVar3;
      return uVar3;
    }
    if (cVar1 != '\x01') {
      if (cVar1 == '\t') {
        *pcVar2 = '\a';
        pcVar2 = pcVar2 + 1;
        *pcVar2 = '+';
      }
      else if (cVar1 == '\n') {
        *pcVar2 = '\a';
        pcVar2 = pcVar2 + 1;
        *pcVar2 = ',';
      }
      else if (cVar1 == '\v') {
        *pcVar2 = '\a';
        pcVar2 = pcVar2 + 1;
        *pcVar2 = '-';
      }
      else if (cVar1 == '\f') {
        *pcVar2 = '\x1a';
      }
      else if (cVar1 == '\'') {
        *pcVar2 = '!';
        pcVar2 = pcVar2 + 1;
        *pcVar2 = 'D';
      }
      else if (cVar1 == '(') {
        *pcVar2 = '!';
        pcVar2 = pcVar2 + 1;
        *pcVar2 = 'E';
      }
      else if (cVar1 == '*') {
        *pcVar2 = '\"';
        pcVar2 = pcVar2 + 1;
        *pcVar2 = 'B';
      }
      else if (cVar1 == '3') {
        *pcVar2 = '*';
        pcVar2 = pcVar2 + 1;
        *pcVar2 = '7';
      }
      else if (cVar1 == '4') {
        *pcVar2 = '*';
        pcVar2 = pcVar2 + 1;
        *pcVar2 = '8';
      }
      else if (cVar1 == 'E') {
        *pcVar2 = '9';
        pcVar2 = pcVar2 + 1;
        *pcVar2 = '7';
      }
      else if (cVar1 == 'N') {
        *pcVar2 = 'B';
        pcVar2 = pcVar2 + 1;
        *pcVar2 = '\b';
      }
      else if (cVar1 == 'O') {
        *pcVar2 = 'B';
        pcVar2 = pcVar2 + 1;
        *pcVar2 = '\t';
      }
      else if (cVar1 == 'Q') {
        *pcVar2 = 'C';
        pcVar2 = pcVar2 + 1;
        *pcVar2 = '\a';
      }
      else if (cVar1 == 'R') {
        *pcVar2 = 'C';
        pcVar2 = pcVar2 + 1;
        *pcVar2 = ';';
      }
      else if (cVar1 == 'S') {
        *pcVar2 = 'C';
        pcVar2 = pcVar2 + 1;
        *pcVar2 = '<';
      }
      else if (cVar1 == 'T') {
        *pcVar2 = 'C';
        pcVar2 = pcVar2 + 1;
        *pcVar2 = '=';
      }
      else if (cVar1 == 'U') {
        *pcVar2 = 'C';
        pcVar2 = pcVar2 + 1;
        *pcVar2 = '>';
      }
      else if (cVar1 == 'V') {
        *pcVar2 = 'C';
        pcVar2 = pcVar2 + 1;
        *pcVar2 = '?';
      }
      else if (cVar1 == 'W') {
        *pcVar2 = 'C';
        pcVar2 = pcVar2 + 1;
        *pcVar2 = '@';
      }
      else {
        if ((cVar1 < '\x01') || ('Y' < cVar1)) {
          *pcVar2 = '\0';
          return (uint)pcVar2 & 0xffff0000;
        }
        *pcVar2 = (&DAT_10077f5a)[cVar1 * 2];
      }
      pcVar2 = pcVar2 + 1;
    }
    cVar1 = param_2[1];
    param_2 = param_2 + 1;
  } while( true );
}



===== 0x1000ff60 =====
Function: FUN_1000ff60 @ 1000ff60

uint __cdecl FUN_1000ff60(char *param_1)

{
  char cVar1;
  uint in_EAX;
  
  if (param_1 != (char *)0x0) {
    cVar1 = *param_1;
    in_EAX = CONCAT31((int3)(in_EAX >> 8),cVar1);
    while (cVar1 != '\0') {
      cVar1 = (char)in_EAX;
      if ((('\0' < cVar1) && (cVar1 < 'F')) &&
         (in_EAX = (uint)cVar1, *(short *)(&DAT_10077d94 + in_EAX * 2) != 0)) {
        return CONCAT22(cVar1 >> 7,1);
      }
      cVar1 = param_1[1];
      in_EAX = CONCAT31((int3)(in_EAX >> 8),cVar1);
      param_1 = param_1 + 1;
    }
  }
  return in_EAX & 0xffff0000;
}



