===== 0x10002200 =====
Function: FUN_10002200 @ 10002200

uint __cdecl FUN_10002200(char *param_1)

{
  char cVar1;
  bool bVar2;
  char cVar3;
  uint in_EAX;
  undefined4 uVar4;
  uint uVar5;
  int iVar6;
  uint uVar7;
  int iVar8;
  char *pcVar9;
  char *pcVar10;
  char *pcVar11;
  char local_4c [68];
  int local_8;
  
  pcVar11 = param_1;
  if ((param_1 == (char *)0x0) || (*param_1 == '\0')) {
    return in_EAX & 0xffff0000;
  }
  uVar5 = 0xffffffff;
  bVar2 = false;
  pcVar9 = param_1;
  do {
    if (uVar5 == 0) break;
    uVar5 = uVar5 - 1;
    cVar3 = *pcVar9;
    pcVar9 = pcVar9 + 1;
  } while (cVar3 != '\0');
  local_8 = ~uVar5 - 1;
  iVar8 = 0;
  iVar6 = 0;
  if (0 < local_8) {
    do {
      if (0x40 < iVar6) break;
      if (((iVar8 < 1) || (cVar3 = param_1[iVar8], cVar3 < '\x1a')) || ('\x1c' < cVar3)) {
        local_4c[iVar6] = param_1[iVar8];
      }
      else {
        cVar1 = param_1[iVar8 + -1];
        if (((((cVar1 < '\x01') || ('\x03' < cVar1)) && ((cVar1 < '\x04' || ('\x06' < cVar1)))) &&
            (((cVar1 < '\n' || ('\f' < cVar1)) && ((cVar1 < '\x17' || ('\x19' < cVar1)))))) &&
           (((cVar1 < '#' || ('%' < cVar1)) && ((cVar1 < ';' || ('=' < cVar1)))))) {
          local_4c[iVar6] = cVar3;
        }
        else {
          local_4c[iVar6] = '6';
          bVar2 = true;
        }
      }
      iVar6 = iVar6 + 1;
      iVar8 = iVar8 + 1;
    } while (iVar8 < local_8);
  }
  local_4c[iVar6] = '\0';
  if (bVar2) {
    uVar5 = 0xffffffff;
    pcVar9 = local_4c;
    do {
      pcVar10 = pcVar9;
      if (uVar5 == 0) break;
      uVar5 = uVar5 - 1;
      pcVar10 = pcVar9 + 1;
      cVar3 = *pcVar9;
      pcVar9 = pcVar10;
    } while (cVar3 != '\0');
    uVar5 = ~uVar5;
    pcVar9 = pcVar10 + -uVar5;
    pcVar10 = param_1;
    for (uVar7 = uVar5 >> 2; uVar7 != 0; uVar7 = uVar7 - 1) {
      *(undefined4 *)pcVar10 = *(undefined4 *)pcVar9;
      pcVar9 = pcVar9 + 4;
      pcVar10 = pcVar10 + 4;
    }
    for (uVar5 = uVar5 & 3; uVar5 != 0; uVar5 = uVar5 - 1) {
      *pcVar10 = *pcVar9;
      pcVar9 = pcVar9 + 1;
      pcVar10 = pcVar10 + 1;
    }
    uVar5 = 0xffffffff;
    do {
      if (uVar5 == 0) break;
      uVar5 = uVar5 - 1;
      cVar3 = *param_1;
      param_1 = param_1 + 1;
    } while (cVar3 != '\0');
    local_8 = ~uVar5 - 1;
  }
  iVar6 = 0;
  iVar8 = 0;
  if (0 < local_8) {
    param_1 = (char *)0xffffffff;
    do {
      if (0x40 < iVar6) break;
      cVar3 = pcVar11[iVar8];
      if ((((cVar3 < '#') || ('%' < cVar3)) && ((cVar3 < ';' || ('=' < cVar3)))) ||
         ((iVar8 != local_8 + -1 &&
          ((local_8 + -1 <= iVar8 || (*(short *)(&DAT_10077d94 + pcVar11[iVar8 + 1] * 2) != 1))))))
      {
        if (((cVar3 == 'B') || (cVar3 == 'C')) && (iVar8 == local_8 + -1)) {
          if (cVar3 == 'B') {
            local_4c[iVar6] = '>';
          }
          else {
            if (cVar3 != 'C') goto LAB_1000249e;
            local_4c[iVar6] = '&';
          }
          goto LAB_1000249d;
        }
        if (((local_8 + -1 <= iVar8) || (cVar3 < '\a')) || ('\t' < cVar3)) {
          if ((((int)param_1 < 0) || (cVar3 < '+')) ||
             (('-' < cVar3 || (uVar5 = FUN_100021d0((int)pcVar11,(uint)param_1), (short)uVar5 != 0))
             )) {
            local_4c[iVar6] = pcVar11[iVar8];
          }
          else {
            if ((((0 < (int)param_1) &&
                 ((pcVar11[iVar8 + -1] == '+' || (pcVar11[iVar8 + -1] == '6')))) &&
                (*(short *)(&DAT_10077d94 + pcVar11[iVar8 + -2] * 2) != 0)) ||
               ((iVar8 < local_8 + -1 && (uVar4 = FUN_10002190(pcVar11), (short)uVar4 == 0)))) {
              cVar3 = pcVar11[iVar8];
              goto LAB_10002469;
            }
            cVar3 = pcVar11[iVar8];
            local_4c[iVar6] = '\a';
            local_4c[iVar6 + 1] = cVar3;
            iVar6 = iVar6 + 1;
          }
          goto LAB_1000249d;
        }
        cVar1 = pcVar11[iVar8 + 1];
        if (cVar1 == '6') {
          if (cVar3 == '\a') {
            local_4c[iVar6] = '\x1a';
            iVar6 = iVar6 + 1;
            iVar8 = iVar8 + 1;
            param_1 = param_1 + 1;
          }
          else if (cVar3 == '\b') {
            local_4c[iVar6] = '\x1b';
            iVar6 = iVar6 + 1;
            iVar8 = iVar8 + 1;
            param_1 = param_1 + 1;
          }
          else {
            if (cVar3 == '\t') {
              local_4c[iVar6] = '\x1c';
              iVar6 = iVar6 + 1;
            }
            iVar8 = iVar8 + 1;
            param_1 = param_1 + 1;
          }
        }
        else {
          if ((cVar1 < '\x1a') || ('\x1c' < cVar1)) {
LAB_10002469:
            local_4c[iVar6] = cVar3;
            goto LAB_1000249d;
          }
          local_4c[iVar6] = cVar1;
          iVar6 = iVar6 + 1;
          iVar8 = iVar8 + 1;
          param_1 = param_1 + 1;
        }
      }
      else {
        if (cVar3 == '#') {
          local_4c[iVar6] = '&';
        }
        else if (cVar3 == '$') {
          local_4c[iVar6] = '\'';
        }
        else if (cVar3 == '%') {
          local_4c[iVar6] = '(';
        }
        else if (cVar3 == ';') {
          local_4c[iVar6] = '>';
        }
        else if (cVar3 == '<') {
          local_4c[iVar6] = '?';
        }
        else {
          if (cVar3 != '=') goto LAB_1000249e;
          local_4c[iVar6] = '@';
        }
LAB_1000249d:
        iVar6 = iVar6 + 1;
      }
LAB_1000249e:
      iVar8 = iVar8 + 1;
      param_1 = param_1 + 1;
    } while (iVar8 < local_8);
  }
  uVar5 = 0xffffffff;
  local_4c[iVar6] = '\0';
  pcVar9 = local_4c;
  do {
    pcVar10 = pcVar9;
    if (uVar5 == 0) break;
    uVar5 = uVar5 - 1;
    pcVar10 = pcVar9 + 1;
    cVar3 = *pcVar9;
    pcVar9 = pcVar10;
  } while (cVar3 != '\0');
  uVar5 = ~uVar5;
  pcVar9 = pcVar10 + -uVar5;
  for (uVar7 = uVar5 >> 2; uVar7 != 0; uVar7 = uVar7 - 1) {
    *(undefined4 *)pcVar11 = *(undefined4 *)pcVar9;
    pcVar9 = pcVar9 + 4;
    pcVar11 = pcVar11 + 4;
  }
  for (uVar5 = uVar5 & 3; uVar5 != 0; uVar5 = uVar5 - 1) {
    *pcVar11 = *pcVar9;
    pcVar9 = pcVar9 + 1;
    pcVar11 = pcVar11 + 1;
  }
  return 1;
}



===== 0x10002680 =====
Function: FUN_10002680 @ 10002680

int __cdecl FUN_10002680(undefined4 param_1,ushort param_2)

{
  char *pcVar1;
  int iVar2;
  ushort uVar3;
  
  iVar2 = FUN_100035d0(DAT_100fee04,param_1);
  if (-1 < iVar2) {
    uVar3 = 0;
    pcVar1 = *(char **)(*(int *)(DAT_100fee04 + 0x10) + iVar2 * 4);
    if (*pcVar1 == '1') {
      uVar3 = 8;
    }
    if (pcVar1[1] == '1') {
      uVar3 = uVar3 | 4;
    }
    if (pcVar1[2] == '1') {
      uVar3 = uVar3 | 2;
    }
    if (pcVar1[3] == '1') {
      uVar3 = uVar3 | 1;
    }
    if ((uVar3 & param_2) == param_2) {
      return iVar2;
    }
  }
  return -1;
}



===== 0x10002180 =====
Function: FUN_10002170 @ 10002170

void __cdecl FUN_10002170(int param_1)

{
  uint uVar1;
  
  uVar1 = FUN_10001ea0(param_1,2);
  FUN_10001ea0(param_1,uVar1);
  return;
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



