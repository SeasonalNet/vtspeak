===== 0x1000cb30 =====
Function: FUN_1000cb30 @ 1000cb30

undefined4 __cdecl FUN_1000cb30(char *param_1,char *param_2,int param_3)

{
  char cVar1;
  undefined4 uVar2;
  uint uVar3;
  uint uVar4;
  char *pcVar5;
  char *pcVar6;
  undefined1 local_290 [8];
  short local_288;
  int local_284;
  char local_280 [428];
  byte local_d4 [204];
  byte local_8 [4];
  
  local_8[0] = 0;
  local_8[1] = 0;
  local_8[2] = 0;
  local_8[3] = 0;
  *param_2 = '\0';
  FUN_10003a70(param_1,local_d4,0,param_3);
  FUN_10003c50(local_290,local_d4);
  if (local_284 < 1) {
    if (local_288 != 0) {
      uVar2 = FUN_10002f10(param_2,param_1,local_8);
      if ((short)uVar2 != 0) {
        return 1;
      }
    }
    return 0xffffffff;
  }
  uVar3 = 0xffffffff;
  pcVar5 = local_280;
  do {
    pcVar6 = pcVar5;
    if (uVar3 == 0) break;
    uVar3 = uVar3 - 1;
    pcVar6 = pcVar5 + 1;
    cVar1 = *pcVar5;
    pcVar5 = pcVar6;
  } while (cVar1 != '\0');
  uVar3 = ~uVar3;
  pcVar5 = pcVar6 + -uVar3;
  for (uVar4 = uVar3 >> 2; uVar4 != 0; uVar4 = uVar4 - 1) {
    *(undefined4 *)param_2 = *(undefined4 *)pcVar5;
    pcVar5 = pcVar5 + 4;
    param_2 = param_2 + 4;
  }
  for (uVar3 = uVar3 & 3; uVar3 != 0; uVar3 = uVar3 - 1) {
    *param_2 = *pcVar5;
    pcVar5 = pcVar5 + 1;
    param_2 = param_2 + 1;
  }
  return 1;
}



