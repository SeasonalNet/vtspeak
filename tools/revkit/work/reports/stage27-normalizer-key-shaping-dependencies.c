===== 0x100021d0 =====
Function: FUN_100021d0 @ 100021d0

uint __cdecl FUN_100021d0(int param_1,uint param_2)

{
  uint in_EAX;
  
  if (((param_1 != 0) && (in_EAX = param_2, -1 < (int)param_2)) && ((int)param_2 < 0x46)) {
    in_EAX = (uint)*(char *)(param_1 + param_2);
    if (*(short *)(&DAT_10077d94 + in_EAX * 2) != 0) {
      return CONCAT22(*(char *)(param_1 + param_2) >> 7,1);
    }
  }
  return in_EAX & 0xffff0000;
}



===== 0x10002190 =====
Function: FUN_10002190 @ 10002190

undefined4 __cdecl FUN_10002190(char *param_1)

{
  char cVar1;
  int iVar2;
  uint uVar3;
  char *pcVar4;
  
  uVar3 = 0xffffffff;
  iVar2 = 0;
  pcVar4 = param_1;
  do {
    if (uVar3 == 0) break;
    uVar3 = uVar3 - 1;
    cVar1 = *pcVar4;
    pcVar4 = pcVar4 + 1;
  } while (cVar1 != '\0');
  if (0 < (int)(~uVar3 - 1)) {
    do {
      if (*(short *)(&DAT_10077d94 + param_1[iVar2] * 2) != 0) {
        return 0;
      }
      iVar2 = iVar2 + 1;
    } while (iVar2 < (int)(~uVar3 - 1));
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



===== 0x100035d0 =====
Function: FUN_100035d0 @ 100035d0

int __cdecl FUN_100035d0(undefined *param_1,undefined4 param_2)

{
  undefined *puVar1;
  int iVar2;
  int iVar3;
  int iVar4;
  int local_8;
  
  puVar1 = param_1;
  if ((param_1 != (undefined *)0x0) && (0 < *(int *)param_1)) {
    iVar4 = *(int *)param_1 + -1;
    local_8 = 0;
    param_1 = FUN_1001c2c0;
    if (puVar1[4] != 'I') {
      param_1 = _strcmp;
    }
    iVar3 = iVar4;
    if (-1 < iVar4) {
      do {
        iVar3 = iVar3 / 2;
        iVar2 = (*(code *)param_1)(*(undefined4 *)(*(int *)(puVar1 + 0xc) + iVar3 * 4),param_2);
        if (iVar2 == 0) {
          return iVar3;
        }
        if (iVar2 < 0) {
          local_8 = iVar3 + 1;
        }
        else {
          iVar4 = iVar3 + -1;
        }
        iVar3 = iVar4 + local_8;
      } while (local_8 <= iVar4);
    }
  }
  return -1;
}



