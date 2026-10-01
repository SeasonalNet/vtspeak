===== 0x10027290 =====
Function: FUN_10027290 @ 10027290

LRESULT FUN_10027290(HWND param_1,UINT param_2,WPARAM param_3,LPARAM param_4)

{
  int *piVar1;
  int *piVar2;
  int iVar3;
  LRESULT LVar4;
  int iVar5;
  int iVar6;
  
  piVar1 = (int *)(&DAT_100a0464)[DAT_100a9d30];
  if (param_2 != 0x3bd) {
    LVar4 = DefWindowProcA(param_1,param_2,param_3,param_4);
    return LVar4;
  }
  iVar5 = (int)DAT_100a9d44;
  if (DAT_100a9d48 == (&DAT_100a7aa0)[iVar5]) {
    if (DAT_100a841c != 1) {
      piVar2 = *(int **)(piVar1[0xa4] + 0x47774);
      if (*(int *)(piVar1[0xa4] + 0x48) == 1) {
        iVar3 = *piVar2;
        iVar6 = DAT_100a9d60;
        if (DAT_100a9d44 != 599) {
          iVar6 = (&DAT_100a9d64)[iVar5];
        }
        if ((*(int *)(iVar3 + 0x10 + (&DAT_100a9d60)[iVar5] * 0x24) ==
             *(int *)(iVar3 + 0x10 + iVar6 * 0x24)) && (DAT_100a7a84 <= DAT_100a9d44)) {
          DAT_100a7498 = 0;
          PostMessageA(DAT_100a9d3c,DAT_100a9d34,0,-1);
          VT_STOPTTS_ENG();
          DAT_100a9d44 = DAT_100a9d44 + 1;
          return 0;
        }
      }
      iVar3 = *piVar2;
      PostMessageA(DAT_100a9d3c,DAT_100a9d34,
                   *(WPARAM *)(iVar3 + (&DAT_100a9d60)[iVar5] * 0x24 + 0xc),
                   *(LPARAM *)(iVar3 + 0x10 + (&DAT_100a9d60)[iVar5] * 0x24));
      if (DAT_100a9d44 == 599) {
        DAT_100a9d44 = -1;
      }
    }
  }
  else if (DAT_100a841c != 1) {
    DAT_100a9d48 = DAT_100a9d48 + 1;
    if (DAT_100a9d48 == 0x7ffffffe) {
      DAT_100a9d48 = 0;
    }
    PostMessageA(DAT_100a9d3c,DAT_100a9d34,
                 *(WPARAM *)
                  (**(int **)(piVar1[0xa4] + 0x47774) + (&DAT_100a9d60)[iVar5] * 0x24 + 0xc),
                 *(LPARAM *)
                  (**(int **)(piVar1[0xa4] + 0x47774) + 0x10 + (&DAT_100a9d60)[iVar5] * 0x24));
    if (DAT_100a9d44 == 599) {
      DAT_100a9d44 = -1;
    }
    iVar5 = piVar1[0xa4];
    if (*(int *)(iVar5 + 0x48) != 0) {
      VT_STOPTTS_ENG();
      DAT_100a7498 = 0;
      PostMessageA(DAT_100a9d3c,DAT_100a9d34,0,-1);
      DAT_100a9d44 = DAT_100a9d44 + 1;
      return 0;
    }
    if (*(int *)(iVar5 + 0x44) != 0) {
      *(undefined4 *)(iVar5 + 0x48) = 1;
      DAT_100a9d44 = DAT_100a9d44 + 1;
      return 0;
    }
    DAT_100a7494 = *(char **)((int)DAT_100a7494 + 8);
    iVar5 = FUN_10026870(piVar1,(int *)piVar1[0xa4],*(int *)(DAT_100a7494 + 4));
    if (iVar5 < 1) {
      return 0;
    }
    FUN_10026d00(DAT_100a7494,*(int *)(piVar1[0xa4] + 0x30));
    DAT_100a9d44 = DAT_100a9d44 + 1;
    return 0;
  }
  DAT_100a9d44 = DAT_100a9d44 + 1;
  return 0;
}


