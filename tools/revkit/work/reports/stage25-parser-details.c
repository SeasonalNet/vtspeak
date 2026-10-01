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



===== 0x1003a7b0 =====
Function: FUN_1003a7b0 @ 1003a7b0

undefined4 __cdecl
FUN_1003a7b0(char *param_1,char *param_2,uint param_3,undefined4 param_4,int param_5)

{
  char *pcVar1;
  uint uVar2;
  char *pcVar3;
  char *pcVar4;
  byte local_d4;
  char local_d3 [203];
  char local_8 [4];
  
  uVar2 = param_3;
  pcVar1 = param_1;
  local_8[0] = '\0';
  *param_1 = '\0';
  pcVar3 = param_2;
  if ((param_2 != (char *)0x0) && (*param_2 != '\0')) {
    pcVar3 = (char *)FUN_1003a570(param_2,&local_d4,0x45,(char)param_4,param_3,param_5);
    if (-1 < (int)pcVar3) {
      if ((int)(char)local_d4 == (uVar2 & 0xff)) {
        pcVar4 = local_d3;
      }
      else {
        FUN_10063ed6(local_8,&DAT_10080384);
        pcVar4 = FUN_1001c390((char *)&local_d4,local_8);
        pcVar3 = (char *)0x0;
        if (pcVar4 == (char *)0x0) goto LAB_1003a86c;
        pcVar4 = pcVar4 + 2;
      }
      pcVar3 = FUN_10024b50(pcVar4,&DAT_1007b3cc,&param_1);
      pcVar3 = FUN_1001c230(pcVar1,pcVar3);
      return CONCAT22((short)((uint)pcVar3 >> 0x10),1);
    }
  }
LAB_1003a86c:
  return (uint)pcVar3 & 0xffff0000;
}



