===== 0x1000e2f0 =====
Function: FUN_1000e2f0 @ 1000e2f0

int __cdecl FUN_1000e2f0(int *param_1,byte *param_2,int param_3,short param_4,int param_5)

{
  int iVar1;
  int iVar2;
  
  if (param_3 == 0) {
    *(undefined2 *)((int)param_1 + 0x429a2) = 0;
    return 0;
  }
  *(undefined2 *)param_1 = 0;
  *(short *)((int)param_1 + 0x429a2) = 0;
  iVar1 = FUN_1000d190(param_1,(int)param_2,param_3,param_4);
  if (iVar1 < 0) {
    return -1;
  }
  if (0 < (short)*param_1) {
    FUN_1000ea20((short *)param_1,(int)param_2,param_3);
  }
  if (param_4 != 0) {
    FUN_10007520((int)param_1,(int)param_2);
  }
  FUN_1000e990((short *)((int)param_1 + 0x429a2));
  iVar2 = FUN_1000dfc0(param_2,param_3,param_5);
  if (iVar2 < 0) {
    return -1;
  }
  return iVar1;
}



