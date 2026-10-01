===== 0x1003a570 =====
Function: FUN_1003a570 @ 1003a570

void __cdecl
FUN_1003a570(char *param_1,byte *param_2,short param_3,char param_4,undefined4 param_5,int param_6)

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
  
  if (param_3 == 0x45) {
    FUN_1003a610(param_1,(int)local_24,param_4);
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
  FUN_10011820(local_f0,param_2,param_6,1);
  return;
}



===== 0x10024b50 =====
Function: FUN_10024b50 @ 10024b50

char * __cdecl FUN_10024b50(char *param_1,char *param_2,undefined4 *param_3)

{
  int iVar1;
  char *pcVar2;
  char *pcVar3;
  
  if (param_1 == (char *)0x0) {
    param_1 = (char *)*param_3;
  }
  iVar1 = FUN_10024c30(param_1,param_2);
  pcVar3 = param_1 + iVar1;
  if (*pcVar3 == '\0') {
    *param_3 = pcVar3;
    return (char *)0x0;
  }
  pcVar2 = FUN_10024c80(pcVar3,param_2);
  if (pcVar2 == (char *)0x0) {
    pcVar2 = FUN_10024bb0(pcVar3,'\0');
    *param_3 = pcVar2;
    return pcVar3;
  }
  *pcVar2 = '\0';
  *param_3 = pcVar2 + 1;
  return pcVar3;
}



===== 0x1001c230 =====
Function: FUN_1001c230 @ 1001c230

char * __cdecl FUN_1001c230(char *param_1,char *param_2)

{
  char cVar1;
  char *pcVar2;
  
  pcVar2 = param_1;
  cVar1 = *param_2;
  *param_1 = cVar1;
  while (cVar1 != '\0') {
    param_2 = param_2 + 1;
    param_1 = param_1 + 1;
    cVar1 = *param_2;
    *param_1 = cVar1;
  }
  return pcVar2;
}



===== 0x10063ed6 =====
Function: FUN_10063ed6 @ 10063ed6

int __cdecl FUN_10063ed6(undefined1 *param_1,byte *param_2)

{
  int iVar1;
  undefined1 *local_24;
  int local_20;
  undefined1 *local_1c;
  undefined4 local_18;
  
  local_1c = param_1;
  local_24 = param_1;
  local_18 = 0x42;
  local_20 = 0x7fffffff;
  iVar1 = FUN_100653b4((int *)&local_24,param_2,(undefined4 *)&stack0x0000000c);
  local_20 = local_20 + -1;
  if (local_20 < 0) {
    FUN_1006529c(0,(int *)&local_24);
  }
  else {
    *local_24 = 0;
  }
  return iVar1;
}



