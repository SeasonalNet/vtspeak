===== 0x1003e240 =====
Function: FUN_1003e240 @ 1003e240

int __cdecl FUN_1003e240(short *param_1)

{
  char cVar1;
  int *piVar2;
  bool bVar3;
  short sVar4;
  int iVar5;
  char *pcVar6;
  uint uVar7;
  uint uVar8;
  char *_SubStr;
  short *_Str;
  short *psVar9;
  short *psVar10;
  char *pcVar11;
  int local_18;
  int local_14;
  
  piVar2 = *(int **)(param_1 + 8);
  bVar3 = false;
  iVar5 = FUN_1000e2f0(piVar2,(byte *)(param_1 + 10),(int)*param_1,1,(int)(param_1 + 0x1cf6));
  if (iVar5 < 0) {
    return CONCAT22((short)((uint)iVar5 >> 0x10),0xffff);
  }
  local_14 = 0;
  local_18 = 0;
  if (0 < *(short *)((int)piVar2 + 0x429a2)) {
    _SubStr = (char *)((int)piVar2 + 0x429ad);
    _Str = param_1 + -0x26;
    do {
      if ((local_18 < 1) || (*(short *)(_SubStr + -0x75) != *(short *)(_SubStr + -5))) {
        *(char *)(_Str + 0x48) = -1;
        if ((*(char *)((int)_Str + 0x83) == 'U') || ((char)_Str[0x42] != 'Y')) {
          uVar7 = 0xffffffff;
          pcVar6 = _SubStr + 0x1e;
          do {
            pcVar11 = pcVar6;
            if (uVar7 == 0) break;
            uVar7 = uVar7 - 1;
            pcVar11 = pcVar6 + 1;
            cVar1 = *pcVar6;
            pcVar6 = pcVar11;
          } while (cVar1 != '\0');
          uVar7 = ~uVar7;
          pcVar6 = pcVar11 + -uVar7;
          psVar9 = _Str + 0x59;
          for (uVar8 = uVar7 >> 2; uVar8 != 0; uVar8 = uVar8 - 1) {
            *(undefined4 *)psVar9 = *(undefined4 *)pcVar6;
            pcVar6 = pcVar6 + 4;
            psVar9 = psVar9 + 2;
          }
          for (uVar7 = uVar7 & 3; uVar7 != 0; uVar7 = uVar7 - 1) {
            *(char *)psVar9 = *pcVar6;
            pcVar6 = pcVar6 + 1;
            psVar9 = (short *)((int)psVar9 + 1);
          }
        }
        iVar5 = *(int *)(_Str + 0x38);
        *(int *)(_Str + 0x38) = iVar5 + 1;
        bVar3 = false;
        *(char *)((int)_Str + iVar5 + 0x91) = -1;
        *(char *)(_Str + 0x41) = _SubStr[-1];
        _Str[0x3e] = *(short *)(_SubStr + -7);
        local_14 = local_14 + 1;
        _Str = _Str + 0x4a;
      }
      else if ((*_SubStr == -0x5e) && (_SubStr[1] == -2)) {
        if (!bVar3) {
          pcVar6 = _strstr((char *)_Str,_SubStr);
          *pcVar6 = '\0';
          bVar3 = true;
          *(int *)(_Str + -0x16) = *(int *)(_Str + -0x16) + -2;
        }
      }
      else {
        if (bVar3) {
          uVar7 = 0xffffffff;
          pcVar6 = _SubStr;
          do {
            pcVar11 = pcVar6;
            if (uVar7 == 0) break;
            uVar7 = uVar7 - 1;
            pcVar11 = pcVar6 + 1;
            cVar1 = *pcVar6;
            pcVar6 = pcVar11;
          } while (cVar1 != '\0');
          uVar7 = ~uVar7;
          iVar5 = -1;
          psVar9 = _Str;
          do {
            psVar10 = psVar9;
            if (iVar5 == 0) break;
            iVar5 = iVar5 + -1;
            psVar10 = (short *)((int)psVar9 + 1);
            sVar4 = *psVar9;
            psVar9 = psVar10;
          } while ((char)sVar4 != '\0');
          pcVar6 = pcVar11 + -uVar7;
          pcVar11 = (char *)((int)psVar10 + -1);
          for (uVar8 = uVar7 >> 2; uVar8 != 0; uVar8 = uVar8 - 1) {
            *(undefined4 *)pcVar11 = *(undefined4 *)pcVar6;
            pcVar6 = pcVar6 + 4;
            pcVar11 = pcVar11 + 4;
          }
          for (uVar7 = uVar7 & 3; uVar7 != 0; uVar7 = uVar7 - 1) {
            *pcVar11 = *pcVar6;
            pcVar6 = pcVar6 + 1;
            pcVar11 = pcVar11 + 1;
          }
        }
        if (*(int *)(_Str + -0x12) < 5) {
          if (_SubStr[0x1e] != '\0') {
            uVar7 = 0xffffffff;
            pcVar6 = _SubStr + 0x1e;
            do {
              pcVar11 = pcVar6;
              if (uVar7 == 0) break;
              uVar7 = uVar7 - 1;
              pcVar11 = pcVar6 + 1;
              cVar1 = *pcVar6;
              pcVar6 = pcVar11;
            } while (cVar1 != '\0');
            uVar7 = ~uVar7;
            iVar5 = -1;
            psVar9 = _Str + 0xf;
            do {
              psVar10 = psVar9;
              if (iVar5 == 0) break;
              iVar5 = iVar5 + -1;
              psVar10 = (short *)((int)psVar9 + 1);
              sVar4 = *psVar9;
              psVar9 = psVar10;
            } while ((char)sVar4 != '\0');
            pcVar6 = pcVar11 + -uVar7;
            pcVar11 = (char *)((int)psVar10 + -1);
            for (uVar8 = uVar7 >> 2; uVar8 != 0; uVar8 = uVar8 - 1) {
              *(undefined4 *)pcVar11 = *(undefined4 *)pcVar6;
              pcVar6 = pcVar6 + 4;
              pcVar11 = pcVar11 + 4;
            }
            for (uVar7 = uVar7 & 3; uVar7 != 0; uVar7 = uVar7 - 1) {
              *pcVar11 = *pcVar6;
              pcVar6 = pcVar6 + 1;
              pcVar11 = pcVar11 + 1;
            }
          }
          *(char *)((int)_Str + *(int *)(_Str + -0x12) + -4) = -1;
          iVar5 = *(int *)(_Str + -0x12);
          *(int *)(_Str + -0x12) = iVar5 + 1;
          *(char *)((int)_Str + iVar5 + -3) = -1;
          _Str[-0xc] = _Str[-0xc] | *(ushort *)(_SubStr + -7);
        }
      }
      _SubStr = _SubStr + 0x70;
      local_18 = local_18 + 1;
    } while (local_18 < *(short *)((int)piVar2 + 0x429a2));
  }
  if (0 < *(short *)((int)piVar2 + 0x429a2)) {
    param_1[1] = (short)piVar2[0x10a69];
  }
  return (-(uint)(*param_1 != local_14) & 0xfffffffe) + 1;
}



