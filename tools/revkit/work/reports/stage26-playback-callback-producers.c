===== 0x10026b40 =====
Function: FUN_10026b40 @ 10026b40

/* WARNING: Globals starting with '_' overlap smaller symbols at the same address */

undefined4 __cdecl FUN_10026b40(int *param_1,int *param_2,char *param_3)

{
  char cVar1;
  char *pcVar2;
  int iVar3;
  int iVar4;
  MMRESULT MVar5;
  uint uVar6;
  
  uVar6 = 0xffffffff;
  param_2[1] = 0;
  param_2[2] = (int)param_3;
  do {
    if (uVar6 == 0) break;
    uVar6 = uVar6 - 1;
    cVar1 = *param_3;
    param_3 = param_3 + 1;
  } while (cVar1 != '\0');
  *param_2 = ~uVar6 - 1;
  param_2[0x11] = 0;
  iVar3 = FUN_10026630((int)param_1,param_2);
  iVar4 = param_2[0x11];
  while( true ) {
    if (iVar4 != 0) {
      return CONCAT22((short)((uint)iVar3 >> 0x10),0xfffe);
    }
    if ((-1 < (short)iVar3) && (*(short *)(param_2[0x13] + 2) != 0)) break;
    iVar3 = FUN_10026630((int)param_1,param_2);
    iVar4 = param_2[0x11];
  }
  DAT_100a7494 = (char *)param_2[0xe];
  iVar4 = FUN_10026870(param_1,param_2,*(int *)(DAT_100a7494 + 4));
  if (iVar4 < 1) {
    return CONCAT22((short)((uint)iVar4 >> 0x10),0xffff);
  }
  param_2[0x12] = 0;
  _DAT_100a9d20 = 1;
  _DAT_100a9d24 = 16000;
  _DAT_100a9d28 = 32000;
  _DAT_100a9d22 = 1;
  _DAT_100a9d2c = 2;
  _DAT_100a9d2e = 0x10;
  MVar5 = waveOutOpen((LPHWAVEOUT)&DAT_100a7490,DAT_1007d6e0,(LPCWAVEFORMATEX)&DAT_100a9d20,
                      DAT_100a8418,0,0x10000);
  if (MVar5 != 0) {
    return CONCAT22((short)(MVar5 >> 0x10),0xffff);
  }
  PostMessageA(DAT_100a9d3c,DAT_100a9d34,*(WPARAM *)(**(int **)(param_1[0xa4] + 0x47774) + 0xc),
               *(LPARAM *)(**(int **)(param_1[0xa4] + 0x47774) + 0x10));
  pcVar2 = DAT_100a7494;
  DAT_100a9d44 = 1;
  if (param_2[0x11] == 0) {
    iVar4 = param_2[0xc];
    DAT_100a7494 = *(char **)(DAT_100a7494 + 8);
    iVar3 = FUN_10026870(param_1,param_2,*(int *)(DAT_100a7494 + 4));
    if (iVar3 < 1) {
      return CONCAT22((short)((uint)iVar3 >> 0x10),0xffff);
    }
    VT_SetDecimal0Pron_ENG();
    FUN_10026d00(pcVar2,iVar4);
    uVar6 = FUN_10026d00(DAT_100a7494,param_2[0xc]);
    return uVar6 & 0xffff0000;
  }
  VT_SetDecimal0Pron_ENG();
  param_2[0x12] = 1;
  uVar6 = FUN_10026d00(DAT_100a7494,param_2[0xc]);
  return uVar6 & 0xffff0000;
}


===== 0x10026870 =====
Function: FUN_10026870 @ 10026870

int __cdecl FUN_10026870(int *param_1,int *param_2,int param_3)

{
  short *psVar1;
  int iVar2;
  int iVar3;
  int iVar4;
  
  param_2[0xd] = param_3;
  FUN_100267d0(param_1,param_2);
  iVar2 = param_2[0xc];
  if (param_1[0x133f] != 0) {
    if (0 < param_1[0x133f]) {
      psVar1 = (short *)param_2[0xd];
      iVar2 = iVar2 / 2;
      iVar4 = 0;
      if (0 < iVar2) {
        do {
          if (iVar4 == 0) {
            iVar3 = ((int)(short)param_2[0x48435] * param_1[0x133f] + 0x32) / 100 + (int)*psVar1;
            if ((0x7ffe < iVar3) || (iVar3 < -0x7fff)) {
              iVar3 = ((iVar3 < 0x7fff) - 1 & 0xfffd) - 0x7fff;
            }
            *(short *)(param_2 + 0x48a17) = (short)iVar3;
          }
          else {
            iVar3 = (int)psVar1[iVar4] +
                    ((int)*(short *)((int)param_2 + iVar4 * 2 + 0x12285a) * param_1[0x133f] + 0x32)
                    / 100;
            if ((0x7ffe < iVar3) || (iVar3 < -0x7fff)) {
              iVar3 = ((iVar3 < 0x7fff) - 1 & 0xfffd) - 0x7fff;
            }
            *(short *)((int)param_2 + iVar4 * 2 + 0x12285c) = (short)iVar3;
          }
          iVar4 = iVar4 + 1;
        } while (iVar4 < iVar2);
      }
      *(undefined2 *)(param_2 + 0x48435) = *(undefined2 *)((int)param_2 + iVar4 * 2 + 0x12285a);
      FUN_10063f30((undefined4 *)psVar1,param_2 + 0x48a17,iVar2 * 2);
      return param_2[0xc];
    }
    psVar1 = (short *)param_2[0xd];
    iVar2 = iVar2 / 2;
    iVar4 = iVar2;
    while (iVar3 = iVar4 + -1, -1 < iVar3) {
      if (iVar3 == 0) {
        iVar4 = ((int)(short)param_2[0x48435] * param_1[0x133f] + 0x32) / 100 + (int)*psVar1;
        if ((0x7ffe < iVar4) || (iVar4 < -0x7fff)) {
          iVar4 = ((iVar4 < 0x7fff) - 1 & 0xfffd) - 0x7fff;
        }
        *(short *)(param_2 + 0x48a17) = (short)iVar4;
        iVar4 = iVar3;
      }
      else {
        iVar4 = (int)psVar1[iVar3] + ((int)psVar1[iVar4 + -2] * param_1[0x133f] + 0x32) / 100;
        if ((0x7ffe < iVar4) || (iVar4 < -0x7fff)) {
          iVar4 = ((iVar4 < 0x7fff) - 1 & 0xfffd) - 0x7fff;
        }
        *(short *)((int)param_2 + iVar3 * 2 + 0x12285c) = (short)iVar4;
        iVar4 = iVar3;
      }
    }
    *(short *)(param_2 + 0x48435) = psVar1[iVar2 + -1];
    FUN_10063f30((undefined4 *)psVar1,param_2 + 0x48a17,iVar2 * 2);
    iVar2 = param_2[0xc];
  }
  return iVar2;
}


===== 0x10026d00 =====
Function: FUN_10026d00 @ 10026d00

undefined4 __cdecl FUN_10026d00(char *param_1,int param_2)

{
  LPWAVEHDR pwVar1;
  int iVar2;
  HWAVEOUT pHVar3;
  int iVar4;
  int iVar5;
  uint uVar6;
  uint uVar7;
  int *piVar8;
  undefined4 *puVar9;
  int iVar10;
  bool bVar11;
  int local_1a0 [100];
  int local_10;
  int *local_c;
  int *local_8;
  
  uVar7 = 0;
  local_10 = 0;
  local_8 = (int *)0x0;
  local_c = *(int **)(*(int *)((&DAT_100a0464)[DAT_100a9d30] + 0x290) + 0x47774);
  piVar8 = local_1a0;
  for (iVar4 = 100; iVar4 != 0; iVar4 = iVar4 + -1) {
    *piVar8 = 0;
    piVar8 = piVar8 + 1;
  }
  uVar6 = param_2 >> 1;
  iVar4 = (-(uint)(DAT_100a9d40 != 1) & 0xfffffe70) + 500;
  piVar8 = (int *)&stack0xfffffe5c;
  do {
    iVar10 = iVar4 * 0x24;
    if (0 < *(int *)(iVar10 + 8 + *local_c)) {
      if (uVar7 == 0) {
        local_8 = (int *)iVar4;
      }
      iVar5 = (int)DAT_100a7a84;
      bVar11 = DAT_100a7a84 == 599;
      (&DAT_100a7aa0)[iVar5] = DAT_100aa6c0;
      (&DAT_100a9d60)[iVar5] = iVar4;
      if (bVar11) {
        DAT_100a7a84 = 0;
      }
      else {
        DAT_100a7a84 = DAT_100a7a84 + 1;
      }
      piVar8 = piVar8 + 1;
      iVar5 = *local_c;
      iVar2 = *(int *)(iVar10 + 8 + iVar5);
      uVar7 = uVar7 + iVar2;
      *piVar8 = iVar2;
      *(undefined4 *)(iVar10 + 8 + iVar5) = 0;
      if (uVar6 < uVar7) {
        iVar5 = uVar7 - uVar6;
        uVar7 = uVar7 - iVar5;
        *piVar8 = *piVar8 - iVar5;
        *(int *)(iVar10 + 8 + *local_c) = iVar5;
      }
    }
    if (iVar4 == 599) {
      iVar4 = 0;
    }
    else {
      iVar4 = iVar4 + 1;
    }
  } while (uVar7 < uVar6);
  DAT_100a9d40 = (uint)(iVar4 < 0x65);
  if ((iVar4 < 100) && (500 < (int)local_8)) {
    iVar4 = (iVar4 - (int)local_8) + 600;
  }
  else {
    iVar4 = iVar4 - (int)local_8;
  }
  pHVar3 = (HWAVEOUT)((int)DAT_100aa6c0 + 1);
  DAT_100aa6c0 = pHVar3;
  if (pHVar3 == (HWAVEOUT)0x7ffffffe) {
    DAT_100aa6c0 = (HWAVEOUT)0x0;
  }
  if (iVar4 != 0) {
    pHVar3 = DAT_100a7490;
    if (*param_1 == '\0') {
      piVar8 = &DAT_100a842c;
      do {
        if (*piVar8 == 1) {
          if (pHVar3 != (HWAVEOUT)0x0) {
            waveOutUnprepareHeader(pHVar3,(LPWAVEHDR)(piVar8 + -3),0x20);
            pHVar3 = DAT_100a7490;
          }
          *piVar8 = 0;
        }
        piVar8 = piVar8 + 8;
      } while ((int)piVar8 < 0x100a90ac);
      local_c = (int *)0x0;
      if (0 < iVar4) {
        local_8 = local_1a0;
        puVar9 = &DAT_100a8430;
        while (iVar10 = *local_8, iVar10 != 0) {
          pwVar1 = (LPWAVEHDR)(puVar9 + -4);
          pwVar1->lpData = (LPSTR)(*(int *)(param_1 + 4) + local_10 * 2);
          puVar9[-3] = iVar10 * 2;
          *puVar9 = 0;
          puVar9[1] = 1;
          puVar9[-1] = 1;
          if ((pHVar3 != (HWAVEOUT)0x0) &&
             (waveOutPrepareHeader(pHVar3,pwVar1,0x20), pHVar3 = DAT_100a7490,
             DAT_100a7490 != (HWAVEOUT)0x0)) {
            waveOutWrite(DAT_100a7490,pwVar1,0x20);
            pHVar3 = DAT_100a7490;
          }
          local_10 = local_10 + iVar10;
          local_c = (int *)((int)local_c + 1);
          puVar9 = puVar9 + 8;
          local_8 = local_8 + 1;
          if (iVar4 <= (int)local_c) {
            return CONCAT22((short)((uint)pHVar3 >> 0x10),1);
          }
        }
      }
    }
    else {
      piVar8 = &DAT_100a90ac;
      do {
        if (*piVar8 == 1) {
          if (pHVar3 != (HWAVEOUT)0x0) {
            waveOutUnprepareHeader(pHVar3,(LPWAVEHDR)(piVar8 + -3),0x20);
            pHVar3 = DAT_100a7490;
          }
          *piVar8 = 0;
        }
        piVar8 = piVar8 + 8;
      } while ((int)piVar8 < 0x100a9d2c);
      local_c = (int *)0x0;
      if (0 < iVar4) {
        local_8 = local_1a0;
        puVar9 = &DAT_100a90b0;
        do {
          iVar10 = *local_8;
          if (iVar10 == 0) break;
          pwVar1 = (LPWAVEHDR)(puVar9 + -4);
          pwVar1->lpData = (LPSTR)(*(int *)(param_1 + 4) + local_10 * 2);
          puVar9[-3] = iVar10 * 2;
          *puVar9 = 0;
          puVar9[1] = 1;
          puVar9[-1] = 1;
          if ((pHVar3 != (HWAVEOUT)0x0) &&
             (waveOutPrepareHeader(pHVar3,pwVar1,0x20), pHVar3 = DAT_100a7490,
             DAT_100a7490 != (HWAVEOUT)0x0)) {
            waveOutWrite(DAT_100a7490,pwVar1,0x20);
            pHVar3 = DAT_100a7490;
          }
          local_10 = local_10 + iVar10;
          local_c = (int *)((int)local_c + 1);
          puVar9 = puVar9 + 8;
          local_8 = local_8 + 1;
        } while ((int)local_c < iVar4);
      }
    }
  }
  return CONCAT22((short)((uint)pHVar3 >> 0x10),1);
}


===== 0x1002c530 =====
Function: FUN_1002c530 @ 1002c530

undefined4 __cdecl FUN_1002c530(int param_1,int param_2)

{
  short sVar1;
  int *piVar2;
  short *psVar3;
  int iVar4;
  int iVar5;
  char *local_10;
  int local_8;
  
  local_8 = 0;
  iVar4 = *(int *)(param_1 + 0x4c);
  piVar2 = *(int **)(param_1 + 0x47774);
  if (0 < *(short *)(param_2 + 0xdf70)) {
    local_10 = (char *)(param_1 + 0xede1f);
    psVar3 = (short *)(param_2 + 0x2e);
    do {
      if (((local_8 < 1) || (psVar3[-0x1d] != psVar3[-3])) ||
         (*(char *)((int)psVar3 + -0x3b) != *(char *)((int)psVar3 + -7))) {
        sVar1 = (short)piVar2[3];
        *psVar3 = sVar1;
        if (*local_10 == '\x01') {
          *(undefined4 *)(*piVar2 + 0xc + sVar1 * 0x24) = *(undefined4 *)(param_1 + 0x4777c);
          *(undefined4 *)(*piVar2 + 0x10 + *psVar3 * 0x24) = *(undefined4 *)(param_1 + 0x47780);
          if (piVar2[3] == 0) {
            iVar5 = 599;
          }
          else {
            iVar5 = piVar2[3] + -1;
          }
          iVar5 = iVar5 * 0x24;
          *(undefined4 *)(*piVar2 + 0x14 + *psVar3 * 0x24) = *(undefined4 *)(iVar5 + 0x14 + *piVar2)
          ;
          *(undefined4 *)(*piVar2 + 0x18 + *psVar3 * 0x24) = *(undefined4 *)(iVar5 + 0x18 + *piVar2)
          ;
          *(undefined4 *)(*piVar2 + 0x1c + *psVar3 * 0x24) = *(undefined4 *)(iVar5 + 0x1c + *piVar2)
          ;
          *(undefined4 *)(*piVar2 + 0x20 + *psVar3 * 0x24) = *(undefined4 *)(iVar5 + 0x20 + *piVar2)
          ;
        }
        else {
          *(undefined4 *)(*piVar2 + 0xc + sVar1 * 0x24) =
               *(undefined4 *)(psVar3[-3] * 0x3c0 + 0x64c + iVar4);
          *(undefined4 *)(*piVar2 + 0x10 + *psVar3 * 0x24) =
               *(undefined4 *)(psVar3[-3] * 0x3c0 + 0x650 + iVar4);
          *(undefined4 *)(param_1 + 0x4777c) = *(undefined4 *)(*piVar2 + 0xc + *psVar3 * 0x24);
          *(undefined4 *)(param_1 + 0x47780) = *(undefined4 *)(*piVar2 + 0x10 + *psVar3 * 0x24);
          *(undefined4 *)(*piVar2 + 0x14 + *psVar3 * 0x24) = *(undefined4 *)(iVar4 + 0x64c);
          *(undefined4 *)(*piVar2 + 0x18 + *psVar3 * 0x24) =
               *(undefined4 *)(*(short *)(iVar4 + 2) * 0x3c0 + 0x290 + iVar4);
          *(int *)(*piVar2 + 0x1c + *psVar3 * 0x24) = (int)*(short *)(iVar4 + 2);
          *(int *)(*piVar2 + 0x20 + *psVar3 * 0x24) = (int)psVar3[-3];
        }
        *(undefined4 *)(*piVar2 + 8 + *psVar3 * 0x24) = 0;
        *(undefined2 *)(*piVar2 + *psVar3 * 0x24) = 0;
        iVar5 = piVar2[3];
        piVar2[3] = iVar5 + 1;
        if (iVar5 + 1 == 600) {
          piVar2[3] = 0;
        }
      }
      else {
        if (piVar2[3] == 0) {
          *psVar3 = 599;
        }
        else {
          *psVar3 = (short)piVar2[3] + -1;
        }
        if (psVar3[-0x1a] != *psVar3) {
          *(undefined4 *)(*piVar2 + 0xc + *psVar3 * 0x24) =
               *(undefined4 *)(psVar3[-3] * 0x3c0 + 0x64c + iVar4);
          *(undefined4 *)(*piVar2 + 0x10 + *psVar3 * 0x24) =
               *(undefined4 *)(psVar3[-3] * 0x3c0 + 0x650 + iVar4);
          *(undefined4 *)(*piVar2 + 0x14 + *psVar3 * 0x24) = *(undefined4 *)(iVar4 + 0x64c);
          *(undefined4 *)(*piVar2 + 0x18 + *psVar3 * 0x24) =
               *(undefined4 *)(*(short *)(iVar4 + 2) * 0x3c0 + 0x290 + iVar4);
          *(int *)(*piVar2 + 0x1c + *psVar3 * 0x24) = (int)*(short *)(iVar4 + 2);
          *(int *)(*piVar2 + 0x20 + *psVar3 * 0x24) = (int)psVar3[-3];
        }
      }
      local_10 = local_10 + 0x34;
      local_8 = local_8 + 1;
      psVar3 = psVar3 + 0x1a;
    } while (local_8 < *(short *)(param_2 + 0xdf70));
  }
  local_8 = 0;
  if (0 < *(short *)(param_2 + 0xdf70)) {
    psVar3 = (short *)(param_2 + -0xc);
    do {
      if (((local_8 < 1) || (*psVar3 != psVar3[0x1a])) ||
         ((psVar3[2] != psVar3[0x1c] ||
          (*(char *)((int)psVar3 + -1) != *(char *)((int)psVar3 + 0x33))))) {
        iVar4 = psVar3[0x1d] * 0x24;
        psVar3[0x1e] = *(short *)(iVar4 + *piVar2);
        iVar5 = (int)*(short *)(*piVar2 + iVar4);
        *(undefined4 *)(*(int *)((short *)(*piVar2 + iVar4) + 2) + iVar5 * 8) = 0;
        *(ushort *)(*(int *)(iVar4 + 4 + *piVar2) + 4 + iVar5 * 8) =
             (ushort)(byte)(&DAT_1007daa8)[*(byte *)(psVar3 + 0x1f)];
        *(short *)(iVar4 + *piVar2) = *(short *)(iVar4 + *piVar2) + 1;
      }
      else {
        psVar3[0x1e] = *(short *)(*piVar2 + psVar3[0x1d] * 0x24) + -1;
      }
      local_8 = local_8 + 1;
      psVar3 = psVar3 + 0x1a;
    } while (local_8 < *(short *)(param_2 + 0xdf70));
  }
  return 1;
}


