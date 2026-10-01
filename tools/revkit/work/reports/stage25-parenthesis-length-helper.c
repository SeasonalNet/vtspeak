===== 0x10062ea0 =====
Function: FUN_10062ea0 @ 10062ea0

int __cdecl FUN_10062ea0(char *param_1)

{
  char cVar1;
  int iVar2;
  int iVar3;
  
  iVar2 = 0;
  iVar3 = 0;
  while( true ) {
    cVar1 = param_1[iVar2];
    if ((cVar1 == '\0') || (2999 < iVar3)) {
      return 0;
    }
    if (cVar1 == ')') {
      return iVar2;
    }
    if (cVar1 == ']') break;
    if (*param_1 == ' ') {
      iVar3 = iVar3 + 1;
    }
    iVar2 = iVar2 + 1;
  }
  return iVar2;
}



