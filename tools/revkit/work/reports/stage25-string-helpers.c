===== 0x1001c2c0 =====
Function: FUN_1001c2c0 @ 1001c2c0

int __cdecl FUN_1001c2c0(byte *param_1,byte *param_2)

{
  byte bVar1;
  int iVar2;
  
  if (param_1 == param_2) {
    return 0;
  }
  iVar2 = (int)*(short *)(&DAT_1007e388 + (uint)*param_1 * 2) -
          (int)*(short *)(&DAT_1007e388 + (uint)*param_2 * 2);
  while( true ) {
    if (iVar2 != 0) {
      return iVar2;
    }
    param_2 = param_2 + 1;
    bVar1 = *param_1;
    param_1 = param_1 + 1;
    if (bVar1 == 0) break;
    iVar2 = (int)*(short *)(&DAT_1007e388 + (uint)*param_1 * 2) -
            (int)*(short *)(&DAT_1007e388 + (uint)*param_2 * 2);
  }
  return 0;
}



===== 0x10062400 =====
Function: FUN_10062400 @ 10062400

char * __cdecl FUN_10062400(char *param_1,char *param_2)

{
  char cVar1;
  uint uVar2;
  uint uVar3;
  char *pcVar4;
  char *pcVar5;
  
  if (*param_1 != '\0') {
    uVar2 = 0xffffffff;
    pcVar4 = param_1;
    do {
      if (uVar2 == 0) break;
      uVar2 = uVar2 - 1;
      cVar1 = *pcVar4;
      pcVar4 = pcVar4 + 1;
    } while (cVar1 != '\0');
    FUN_10063ed6(param_1 + (~uVar2 - 1),&DAT_1009c420);
    return param_1;
  }
  uVar2 = 0xffffffff;
  do {
    pcVar4 = param_2;
    if (uVar2 == 0) break;
    uVar2 = uVar2 - 1;
    pcVar4 = param_2 + 1;
    cVar1 = *param_2;
    param_2 = pcVar4;
  } while (cVar1 != '\0');
  uVar2 = ~uVar2;
  pcVar4 = pcVar4 + -uVar2;
  pcVar5 = param_1;
  for (uVar3 = uVar2 >> 2; uVar3 != 0; uVar3 = uVar3 - 1) {
    *(undefined4 *)pcVar5 = *(undefined4 *)pcVar4;
    pcVar4 = pcVar4 + 4;
    pcVar5 = pcVar5 + 4;
  }
  for (uVar2 = uVar2 & 3; uVar2 != 0; uVar2 = uVar2 - 1) {
    *pcVar5 = *pcVar4;
    pcVar4 = pcVar4 + 1;
    pcVar5 = pcVar5 + 1;
  }
  return param_1;
}



