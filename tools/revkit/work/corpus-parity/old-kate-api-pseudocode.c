/* Ghidra pseudocode for the extracted 2006 MSI DLL; names/control flow may be imperfect. */
/* Program: vt_eng.dll */

/* ===== VT_GetTTSInfo_ENG ===== */
/* Entry: 10022c40 */

undefined4 __cdecl VT_GetTTSInfo_ENG(undefined4 param_1,char *param_2,DWORD *param_3,int param_4)

{
  char cVar1;
  BYTE BVar2;
  DWORD DVar3;
  undefined *puVar4;
  UINT UVar5;
  int iVar6;
  uint uVar7;
  uint uVar8;
  DWORD *pDVar9;
  char *pcVar10;
  char *pcVar11;
  BYTE *pBVar12;
  BYTE *pBVar13;
  BYTE local_200 [512];
  
                    /* 0x22c40  21  VT_GetTTSInfo_ENG */
  if (DAT_10070974 < 0) {
    DAT_10070974 = FUN_1001f8b0();
  }
  if (DAT_1008ceb0 == 0) {
    FUN_1000e890();
  }
  if (DAT_100942cc == 0) {
    InitializeCriticalSection((LPCRITICAL_SECTION)&DAT_10096a00);
    DAT_100942cc = 1;
  }
  if (DAT_100940b4 == 0) {
    FUN_1001efd0();
    DAT_100940b4 = 1;
  }
  if (param_3 == (DWORD *)0x0) {
    return 3;
  }
  switch(param_1) {
  case 0:
    uVar7 = 0xffffffff;
    pcVar10 = PTR_DAT_10070988;
    do {
      if (uVar7 == 0) break;
      uVar7 = uVar7 - 1;
      cVar1 = *pcVar10;
      pcVar10 = pcVar10 + 1;
    } while (cVar1 != '\0');
    pcVar10 = PTR_DAT_10070988;
    if ((int)(~uVar7 - 1) < param_4) {
LAB_10022d85:
      uVar7 = 0xffffffff;
      do {
        pcVar11 = pcVar10;
        if (uVar7 == 0) break;
        uVar7 = uVar7 - 1;
        pcVar11 = pcVar10 + 1;
        cVar1 = *pcVar10;
        pcVar10 = pcVar11;
      } while (cVar1 != '\0');
      uVar7 = ~uVar7;
      pDVar9 = (DWORD *)(pcVar11 + -uVar7);
      for (uVar8 = uVar7 >> 2; uVar8 != 0; uVar8 = uVar8 - 1) {
        *param_3 = *pDVar9;
        pDVar9 = pDVar9 + 1;
        param_3 = param_3 + 1;
      }
      for (uVar7 = uVar7 & 3; uVar7 != 0; uVar7 = uVar7 - 1) {
        *(char *)param_3 = (char)*pDVar9;
        pDVar9 = (DWORD *)((int)pDVar9 + 1);
        param_3 = (DWORD *)((int)param_3 + 1);
      }
      return 0;
    }
    break;
  case 1:
    DVar3 = FUN_10022430(param_2);
    *param_3 = DVar3;
    return 0;
  case 2:
    DVar3 = FUN_10022430(param_2);
    if (DVar3 != 0) {
      *param_3 = 1;
      return 0;
    }
    puVar4 = FUN_10022010();
    *param_3 = (DWORD)puVar4;
    return 0;
  case 3:
    iVar6 = -1;
    pcVar10 = (char *)&DAT_1008ccac;
    do {
      if (iVar6 == 0) break;
      iVar6 = iVar6 + -1;
      cVar1 = *pcVar10;
      pcVar10 = pcVar10 + 1;
    } while (cVar1 != '\0');
    if (iVar6 == -2) {
      FUN_10021e80(local_200);
      uVar7 = 0xffffffff;
      pBVar12 = local_200;
      do {
        if (uVar7 == 0) break;
        uVar7 = uVar7 - 1;
        BVar2 = *pBVar12;
        pBVar12 = pBVar12 + 1;
      } while (BVar2 != '\0');
      if ((int)(~uVar7 - 1) < param_4) {
        uVar7 = 0xffffffff;
        pBVar12 = local_200;
        do {
          pBVar13 = pBVar12;
          if (uVar7 == 0) break;
          uVar7 = uVar7 - 1;
          pBVar13 = pBVar12 + 1;
          BVar2 = *pBVar12;
          pBVar12 = pBVar13;
        } while (BVar2 != '\0');
        uVar7 = ~uVar7;
        pDVar9 = (DWORD *)(pBVar13 + -uVar7);
        for (uVar8 = uVar7 >> 2; uVar8 != 0; uVar8 = uVar8 - 1) {
          *param_3 = *pDVar9;
          pDVar9 = pDVar9 + 1;
          param_3 = param_3 + 1;
        }
        for (uVar7 = uVar7 & 3; uVar7 != 0; uVar7 = uVar7 - 1) {
          *(BYTE *)param_3 = (BYTE)*pDVar9;
          pDVar9 = (DWORD *)((int)pDVar9 + 1);
          param_3 = (DWORD *)((int)param_3 + 1);
        }
        return 0;
      }
    }
    else {
      uVar7 = 0xffffffff;
      pcVar10 = (char *)&DAT_1008ccac;
      do {
        if (uVar7 == 0) break;
        uVar7 = uVar7 - 1;
        cVar1 = *pcVar10;
        pcVar10 = pcVar10 + 1;
      } while (cVar1 != '\0');
      if ((int)(~uVar7 - 1) < param_4) {
        pcVar10 = (char *)&DAT_1008ccac;
        goto LAB_10022d85;
      }
    }
    break;
  case 4:
  case 6:
  case 8:
  case 9:
  case 0x13:
  case 0x16:
    *param_3 = 0;
    return 0;
  case 5:
    *param_3 = 2;
    return 0;
  case 7:
    UVar5 = GetACP();
    *param_3 = UVar5;
    return 0;
  case 10:
    *param_3 = 16000;
    return 0;
  case 0xb:
    *param_3 = 200;
    return 0;
  case 0xc:
  case 0xf:
  case 0x12:
    *param_3 = 100;
    return 0;
  case 0xd:
  case 0x10:
    *param_3 = 0x32;
    return 0;
  case 0xe:
    *param_3 = 400;
    return 0;
  case 0x11:
    *param_3 = 500;
    return 0;
  case 0x14:
    *param_3 = 0x10000;
    return 0;
  case 0x15:
    *param_3 = 0x2af;
    return 0;
  default:
    return 2;
  }
  return 4;
}



/* ===== VT_LOADTTS_ENG ===== */
/* Entry: 10021570 */

/* WARNING: Globals starting with '_' overlap smaller symbols at the same address */

short __cdecl VT_LOADTTS_ENG(HWND param_1,int param_2,char *param_3,char *param_4)

{
  char cVar1;
  byte bVar2;
  DWORD DVar3;
  undefined *puVar4;
  int *piVar5;
  uint uVar6;
  uint uVar7;
  int iVar8;
  int iVar9;
  byte *pbVar10;
  char *pcVar11;
  byte *pbVar12;
  undefined4 *puVar13;
  short local_206 [3];
  byte local_200 [512];
  
                    /* 0x21570  25  VT_LOADTTS_ENG */
  FUN_1000e890();
  FUN_1000e890();
  if (DAT_10070974 < 0) {
    DAT_10070974 = FUN_1001f8b0();
  }
  if ((param_2 < 0) || (iVar9 = param_2, 1 < param_2)) {
    iVar9 = 0;
  }
  if (DAT_1008ceb0 == 0) {
    FUN_1000e890();
  }
  if (param_3 == (char *)0x0) {
    FUN_10021e80((BYTE *)&DAT_1008ccac);
  }
  else {
    iVar8 = -1;
    pcVar11 = param_3;
    do {
      if (iVar8 == 0) break;
      iVar8 = iVar8 + -1;
      cVar1 = *pcVar11;
      pcVar11 = pcVar11 + 1;
    } while (cVar1 != '\0');
    if (iVar8 == -2) {
      FUN_10021e80((BYTE *)&DAT_1008ccac);
    }
    else {
      uVar6 = 0xffffffff;
      do {
        pcVar11 = param_3;
        if (uVar6 == 0) break;
        uVar6 = uVar6 - 1;
        pcVar11 = param_3 + 1;
        cVar1 = *param_3;
        param_3 = pcVar11;
      } while (cVar1 != '\0');
      uVar6 = ~uVar6;
      pbVar10 = (byte *)(pcVar11 + -uVar6);
      pbVar12 = local_200;
      for (uVar7 = uVar6 >> 2; uVar7 != 0; uVar7 = uVar7 - 1) {
        *(undefined4 *)pbVar12 = *(undefined4 *)pbVar10;
        pbVar10 = pbVar10 + 4;
        pbVar12 = pbVar12 + 4;
      }
      for (uVar6 = uVar6 & 3; uVar6 != 0; uVar6 = uVar6 - 1) {
        *pbVar12 = *pbVar10;
        pbVar10 = pbVar10 + 1;
        pbVar12 = pbVar12 + 1;
      }
      FUN_10021530((char *)local_200);
      if (DAT_100942d0 == '\x01') {
        iVar8 = FUN_10014e70((byte *)&DAT_1008ccac,local_200);
        if (iVar8 == 0) goto LAB_100216b5;
        if ((&DAT_100962e4)[param_2] == '\x01') {
          return 1;
        }
        uVar6 = 0xffffffff;
        pbVar10 = local_200;
        do {
          pbVar12 = pbVar10;
          if (uVar6 == 0) break;
          uVar6 = uVar6 - 1;
          pbVar12 = pbVar10 + 1;
          bVar2 = *pbVar10;
          pbVar10 = pbVar12;
        } while (bVar2 != 0);
        uVar6 = ~uVar6;
        pbVar10 = pbVar12 + -uVar6;
        pbVar12 = (byte *)&DAT_1008ccac;
        for (uVar7 = uVar6 >> 2; uVar7 != 0; uVar7 = uVar7 - 1) {
          *(undefined4 *)pbVar12 = *(undefined4 *)pbVar10;
          pbVar10 = pbVar10 + 4;
          pbVar12 = pbVar12 + 4;
        }
      }
      else {
        uVar6 = 0xffffffff;
        pbVar10 = local_200;
        do {
          pbVar12 = pbVar10;
          if (uVar6 == 0) break;
          uVar6 = uVar6 - 1;
          pbVar12 = pbVar10 + 1;
          bVar2 = *pbVar10;
          pbVar10 = pbVar12;
        } while (bVar2 != 0);
        uVar6 = ~uVar6;
        pbVar10 = pbVar12 + -uVar6;
        pbVar12 = (byte *)&DAT_1008ccac;
        for (uVar7 = uVar6 >> 2; uVar7 != 0; uVar7 = uVar7 - 1) {
          *(undefined4 *)pbVar12 = *(undefined4 *)pbVar10;
          pbVar10 = pbVar10 + 4;
          pbVar12 = pbVar12 + 4;
        }
      }
      for (uVar6 = uVar6 & 3; uVar6 != 0; uVar6 = uVar6 - 1) {
        *pbVar12 = *pbVar10;
        pbVar10 = pbVar10 + 1;
        pbVar12 = pbVar12 + 1;
      }
    }
  }
LAB_100216b5:
  FUN_10021530((char *)&DAT_1008ccac);
  if ((DAT_100962e6 == '\0') && (DVar3 = FUN_10022430(param_4), DVar3 == 0)) {
    DAT_100962e6 = '\x01';
  }
  if (DAT_100942d0 != '\x01') {
    if (param_1 != (HWND)0x0) {
      PostMessageA(param_1,0x464,10,0);
    }
    FUN_10021080(DAT_10097438);
    if (DAT_100942cc == 0) {
      InitializeCriticalSection((LPCRITICAL_SECTION)&DAT_10096a00);
      DAT_100942cc = 1;
    }
    if (DAT_100940b4 == 0) {
      FUN_1001efd0();
      DAT_100940b4 = 1;
    }
    FUN_1000e890();
    puVar13 = &DAT_100952e0;
    for (iVar8 = 0x400; iVar8 != 0; iVar8 = iVar8 + -1) {
      *puVar13 = 0;
      puVar13 = puVar13 + 1;
    }
    FUN_1001f8c0(param_1,local_206);
    if (local_206[0] != 0) {
      if (DAT_100942d4 != (LPVOID)0x0) {
        FUN_10016500(DAT_100942d4);
      }
      DAT_100942d4 = (LPVOID)0x0;
      return local_206[0];
    }
    DAT_100942d0 = '\x01';
    if (DAT_100962e6 == '\0') {
      *(undefined2 *)((int)DAT_100942d4 + 0x2048c) = 1;
    }
    else {
      puVar4 = FUN_10022010();
      *(short *)((int)DAT_100942d4 + 0x2048c) = (short)puVar4;
    }
    DAT_100942d8 = 0;
    puVar13 = &DAT_100942e0;
    for (iVar8 = 0x400; iVar8 != 0; iVar8 = iVar8 + -1) {
      *puVar13 = 0;
      puVar13 = puVar13 + 1;
    }
    DAT_100942dc = 0;
    _DAT_100962e4 = 0;
    puVar13 = &DAT_100977a0;
    for (iVar8 = 0x1c00; iVar8 != 0; iVar8 = iVar8 + -1) {
      *puVar13 = 0;
      puVar13 = puVar13 + 1;
    }
  }
  if ((&DAT_100942d8)[iVar9] == 0) {
    piVar5 = FUN_1001f950(param_1,local_206,iVar9);
    (&DAT_100942d8)[iVar9] = piVar5;
    if (local_206[0] != 0) {
      return local_206[0];
    }
    piVar5[0xe11] = (int)DAT_100942d4;
    *(undefined4 *)((&DAT_100942d8)[iVar9] + 0x294) = 0;
  }
  FUN_1000e890();
  FUN_1000e890();
  FUN_1000e890();
  return local_206[0];
}



/* ===== VT_TextToFile_ENG ===== */
/* Entry: 10016510 */

undefined4 __cdecl
VT_TextToFile_ENG(undefined4 param_1,byte *param_2,int param_3,int *param_4,int param_5,int param_6,
                 int param_7,undefined4 param_8,int param_9,int param_10)

{
  undefined4 uVar1;
  
                    /* 0x16510  37  VT_TextToFile_ENG */
  switch(param_1) {
  case 0:
    uVar1 = FUN_10016dd0(param_2,param_3,(int)param_4,param_5,param_6,param_7,param_8,param_9,
                         param_10);
    return uVar1;
  case 1:
    uVar1 = FUN_10016f60(param_2,param_3,(int)param_4,param_5,param_6,param_7,param_8,param_9,
                         param_10);
    return uVar1;
  case 2:
    uVar1 = FUN_10017140(param_2,param_3,(int)param_4,param_5,param_6,param_7,param_8,param_9,
                         param_10);
    return uVar1;
  case 3:
    uVar1 = FUN_10017820(param_2,param_3,(int)param_4,param_5,param_6,param_7,param_8,param_9,
                         param_10);
    return uVar1;
  case 4:
    uVar1 = FUN_10017320(param_2,param_3,param_4,param_5,param_6,param_7,param_8,param_9,param_10);
    return uVar1;
  case 5:
    uVar1 = FUN_10017360(param_2,param_3,param_4,param_5,param_6,param_7,param_8,param_9,param_10);
    return uVar1;
  default:
    return CONCAT22((short)((uint)param_1 >> 0x10),0xffff);
  case 7:
    uVar1 = FUN_100173a0(param_2,param_3,param_4,param_5,param_6,param_7,param_8,param_9,param_10);
    return uVar1;
  case 8:
    uVar1 = FUN_100173e0(param_2,param_3,param_4,param_5,param_6,param_7,param_8,param_9,param_10);
    return uVar1;
  case 9:
    uVar1 = FUN_10018710(param_2,param_3,(int)param_4,param_5,param_6,param_7,param_8,param_9,
                         param_10);
    return uVar1;
  }
}



/* ===== VT_UNLOADTTS_ENG ===== */
/* Entry: 10021890 */

void __cdecl VT_UNLOADTTS_ENG(int param_1)

{
  int *piVar1;
  int iVar2;
  
                    /* 0x21890  42  VT_UNLOADTTS_ENG */
  if ((param_1 < 0) || (1 < param_1)) {
    param_1 = 0;
  }
  VT_STOPTTS_ENG();
  DAT_10097430 = param_1;
  if ((&DAT_100942d8)[param_1] != 0) {
    FUN_10021290(param_1);
  }
  piVar1 = &DAT_100942d8;
  do {
    if (*piVar1 != 0) {
      return;
    }
    piVar1 = piVar1 + 1;
  } while ((int)piVar1 < 0x100942e0);
  iVar2 = 0;
  do {
    VT_UNLOAD_UserDict_ENG(iVar2);
    iVar2 = iVar2 + 1;
  } while (iVar2 < 0x400);
  FUN_10021970();
  DestroyWindow(DAT_10096a18);
  if (DAT_100942cc == 1) {
    DeleteCriticalSection((LPCRITICAL_SECTION)&DAT_10096a00);
    DAT_100942cc = 0;
  }
  if (DAT_100940b4 == 1) {
    FUN_1000e890();
    DAT_100940b4 = 0;
  }
  FUN_10022380();
  DAT_100962e6 = 0;
  DAT_1008ccac._0_1_ = DAT_1008cbe8;
  if (DAT_1008ceb0 == 1) {
    FUN_1000e890();
  }
  FUN_1000e890();
  return;
}



/* ===== FUN_10017320 ===== */
/* Entry: 10017320 */

void __cdecl
FUN_10017320(byte *param_1,int param_2,int *param_3,int param_4,int param_5,int param_6,
            undefined4 param_7,int param_8,int param_9)

{
  FUN_10017420(param_1,param_2,1,param_3,param_4,param_5,param_6,param_7,param_8,param_9);
  return;
}



/* ===== FUN_10016dd0 ===== */
/* Entry: 10016dd0 */

undefined4 __cdecl
FUN_10016dd0(byte *param_1,undefined4 param_2,int param_3,int param_4,int param_5,int param_6,
            undefined4 param_7,int param_8,int param_9)

{
  int *piVar1;
  int *piVar2;
  int *piVar3;
  int iVar4;
  undefined4 *puVar5;
  char *pcVar6;
  undefined4 uVar7;
  int iVar8;
  
  iVar4 = param_3;
  if ((param_3 < 0) || (iVar8 = param_3, 1 < param_3)) {
    iVar8 = 0;
  }
  piVar1 = (int *)(&DAT_100942d8)[iVar8];
  if (piVar1 == (int *)0x0) {
    return CONCAT22((short)((uint)iVar8 >> 0x10),0xfffb);
  }
  if (param_1 == (byte *)0x0) {
    return 0xfffd;
  }
  if (*param_1 == 0) {
    return CONCAT22((short)((uint)param_1 >> 0x10),0xfffc);
  }
  piVar2 = FUN_1001efe0(param_2,&DAT_1006f110);
  if (piVar2 == (int *)0x0) {
    return 0xfffa;
  }
  iVar8 = 0;
  piVar3 = (int *)FUN_1001f9f0((undefined2 *)&param_3,param_8);
  if ((short)param_3 != 1) {
    return CONCAT22((short)((uint)piVar3 >> 0x10),0xfffe);
  }
  if ((iVar4 < 0) || (1 < iVar4)) {
    iVar4 = 0;
  }
  *(char *)(piVar3 + 0x14a38) = (char)iVar4;
  FUN_10021f10(piVar1,(int)piVar3,param_4,param_5,param_6,param_7,param_9);
  puVar5 = FUN_10016a50((byte *)0x0);
  piVar3[0xb] = (int)puVar5;
  pcVar6 = (char *)FUN_10015370((int)piVar3,param_1,0);
  if (pcVar6 == (char *)0x0) {
    FUN_10016bb0((LPVOID)piVar3[0xb]);
    uVar7 = FUN_100212e0(piVar3);
    return CONCAT22((short)((uint)uVar7 >> 0x10),0xfffd);
  }
  FUN_10020750((int)piVar1,piVar3,pcVar6);
  iVar4 = piVar3[0x13];
  while (iVar4 == 0) {
    FUN_10020510(piVar1,piVar3,*(int *)(piVar3[0x10] + 4));
    FUN_10016c00((int)piVar3);
    FUN_1001f660(piVar2,iVar8,0,*(char **)(piVar3[0x10] + 4),1,piVar3[0xe]);
    iVar8 = iVar8 + piVar3[0xe];
    iVar4 = piVar3[0x13];
  }
  FUN_1001f550(piVar2);
  FUN_10016bb0((LPVOID)piVar3[0xb]);
  uVar7 = FUN_100212e0(piVar3);
  return CONCAT22((short)((uint)uVar7 >> 0x10),1);
}



/* ===== FUN_10016f60 ===== */
/* Entry: 10016f60 */

undefined4 __cdecl
FUN_10016f60(byte *param_1,undefined4 param_2,int param_3,int param_4,int param_5,int param_6,
            undefined4 param_7,int param_8,int param_9)

{
  int *piVar1;
  int *piVar2;
  int *piVar3;
  undefined4 *puVar4;
  char *pcVar5;
  undefined4 uVar6;
  int iVar7;
  int iVar8;
  
  iVar7 = param_3;
  if ((param_3 < 0) || (iVar8 = param_3, 1 < param_3)) {
    iVar8 = 0;
  }
  piVar1 = (int *)(&DAT_100942d8)[iVar8];
  if (piVar1 == (int *)0x0) {
    return CONCAT22((short)((uint)iVar8 >> 0x10),0xfffb);
  }
  if (param_1 == (byte *)0x0) {
    return 0xfffd;
  }
  if (*param_1 == 0) {
    return CONCAT22((short)((uint)param_1 >> 0x10),0xfffc);
  }
  piVar2 = FUN_1001efe0(param_2,&DAT_1006f110);
  if (piVar2 == (int *)0x0) {
    return 0xfffa;
  }
  iVar8 = 0;
  piVar3 = (int *)FUN_1001f9f0((undefined2 *)&param_3,param_8);
  if ((short)param_3 != 1) {
    return CONCAT22((short)((uint)piVar3 >> 0x10),0xfffe);
  }
  if ((iVar7 < 0) || (1 < iVar7)) {
    iVar7 = 0;
  }
  *(char *)(piVar3 + 0x14a38) = (char)iVar7;
  FUN_10021f10(piVar1,(int)piVar3,param_4,param_5,param_6,param_7,param_9);
  puVar4 = FUN_10016a50((byte *)0x0);
  piVar3[0xb] = (int)puVar4;
  pcVar5 = (char *)FUN_10015370((int)piVar3,param_1,0);
  if (pcVar5 == (char *)0x0) {
    FUN_10016bb0((LPVOID)piVar3[0xb]);
    uVar6 = FUN_100212e0(piVar3);
    return CONCAT22((short)((uint)uVar6 >> 0x10),0xfffd);
  }
  FUN_10020750((int)piVar1,piVar3,pcVar5);
  iVar7 = piVar3[0x13];
  while (iVar7 == 0) {
    FUN_10020510(piVar1,piVar3,*(int *)(piVar3[0x10] + 4));
    FUN_10016c00((int)piVar3);
    iVar7 = 0;
    if (0 < piVar3[0xe] / 2) {
      do {
        iVar7 = iVar7 + 1;
        *(undefined1 *)(*(int *)(piVar3[0x11] + 4) + -1 + iVar7) =
             *(undefined1 *)
              (*(short *)(*(int *)(piVar3[0x10] + 4) + -2 + iVar7 * 2) + 0x1848c + piVar1[0xe11]);
      } while (iVar7 < piVar3[0xe] / 2);
    }
    FUN_1001f660(piVar2,iVar8,0,*(char **)(piVar3[0x11] + 4),1,piVar3[0xe] / 2);
    iVar8 = iVar8 + piVar3[0xe] / 2;
    iVar7 = piVar3[0x13];
  }
  FUN_1001f550(piVar2);
  FUN_10016bb0((LPVOID)piVar3[0xb]);
  uVar6 = FUN_100212e0(piVar3);
  return CONCAT22((short)((uint)uVar6 >> 0x10),1);
}



/* ===== FUN_10017140 ===== */
/* Entry: 10017140 */

undefined4 __cdecl
FUN_10017140(byte *param_1,undefined4 param_2,int param_3,int param_4,int param_5,int param_6,
            undefined4 param_7,int param_8,int param_9)

{
  int *piVar1;
  int *piVar2;
  int *piVar3;
  undefined4 *puVar4;
  char *pcVar5;
  undefined4 uVar6;
  int iVar7;
  int iVar8;
  
  iVar7 = param_3;
  if ((param_3 < 0) || (iVar8 = param_3, 1 < param_3)) {
    iVar8 = 0;
  }
  piVar1 = (int *)(&DAT_100942d8)[iVar8];
  if (piVar1 == (int *)0x0) {
    return CONCAT22((short)((uint)iVar8 >> 0x10),0xfffb);
  }
  if (param_1 == (byte *)0x0) {
    return 0xfffd;
  }
  if (*param_1 == 0) {
    return CONCAT22((short)((uint)param_1 >> 0x10),0xfffc);
  }
  piVar2 = FUN_1001efe0(param_2,&DAT_1006f110);
  if (piVar2 == (int *)0x0) {
    return 0xfffa;
  }
  iVar8 = 0;
  piVar3 = (int *)FUN_1001f9f0((undefined2 *)&param_3,param_8);
  if ((short)param_3 != 1) {
    return CONCAT22((short)((uint)piVar3 >> 0x10),0xfffe);
  }
  if ((iVar7 < 0) || (1 < iVar7)) {
    iVar7 = 0;
  }
  *(char *)(piVar3 + 0x14a38) = (char)iVar7;
  FUN_10021f10(piVar1,(int)piVar3,param_4,param_5,param_6,param_7,param_9);
  puVar4 = FUN_10016a50((byte *)0x0);
  piVar3[0xb] = (int)puVar4;
  pcVar5 = (char *)FUN_10015370((int)piVar3,param_1,0);
  if (pcVar5 == (char *)0x0) {
    FUN_10016bb0((LPVOID)piVar3[0xb]);
    uVar6 = FUN_100212e0(piVar3);
    return CONCAT22((short)((uint)uVar6 >> 0x10),0xfffd);
  }
  FUN_10020750((int)piVar1,piVar3,pcVar5);
  iVar7 = piVar3[0x13];
  while (iVar7 == 0) {
    FUN_10020510(piVar1,piVar3,*(int *)(piVar3[0x10] + 4));
    FUN_10016c00((int)piVar3);
    iVar7 = 0;
    if (0 < piVar3[0xe] / 2) {
      do {
        iVar7 = iVar7 + 1;
        *(undefined1 *)(*(int *)(piVar3[0x11] + 4) + -1 + iVar7) =
             *(undefined1 *)
              (*(short *)(*(int *)(piVar3[0x10] + 4) + -2 + iVar7 * 2) + 0x848c + piVar1[0xe11]);
      } while (iVar7 < piVar3[0xe] / 2);
    }
    FUN_1001f660(piVar2,iVar8,0,*(char **)(piVar3[0x11] + 4),1,piVar3[0xe] / 2);
    iVar8 = iVar8 + piVar3[0xe] / 2;
    iVar7 = piVar3[0x13];
  }
  FUN_1001f550(piVar2);
  FUN_10016bb0((LPVOID)piVar3[0xb]);
  uVar6 = FUN_100212e0(piVar3);
  return CONCAT22((short)((uint)uVar6 >> 0x10),1);
}



/* ===== FUN_10017820 ===== */
/* Entry: 10017820 */

undefined4 __cdecl
FUN_10017820(byte *param_1,undefined4 param_2,int param_3,int param_4,int param_5,int param_6,
            undefined4 param_7,int param_8,int param_9)

{
  int *piVar1;
  int *piVar2;
  int *piVar3;
  int iVar4;
  undefined4 *puVar5;
  char *pcVar6;
  undefined4 uVar7;
  int iVar8;
  int iVar9;
  
  iVar4 = param_3;
  if ((param_3 < 0) || (iVar8 = param_3, 1 < param_3)) {
    iVar8 = 0;
  }
  piVar1 = (int *)(&DAT_100942d8)[iVar8];
  if (piVar1 == (int *)0x0) {
    return CONCAT22((short)((uint)iVar8 >> 0x10),0xfffb);
  }
  if (param_1 == (byte *)0x0) {
    return 0xfffd;
  }
  if (*param_1 == 0) {
    return CONCAT22((short)((uint)param_1 >> 0x10),0xfffc);
  }
  piVar2 = FUN_1001efe0(param_2,&DAT_1006f110);
  if (piVar2 == (int *)0x0) {
    return 0xfffa;
  }
  iVar8 = 0;
  piVar3 = (int *)FUN_1001f9f0((undefined2 *)&param_3,param_8);
  if ((short)param_3 != 1) {
    return CONCAT22((short)((uint)piVar3 >> 0x10),0xfffe);
  }
  if ((iVar4 < 0) || (1 < iVar4)) {
    iVar4 = 0;
  }
  *(char *)(piVar3 + 0x14a38) = (char)iVar4;
  FUN_10021f10(piVar1,(int)piVar3,param_4,param_5,param_6,param_7,param_9);
  puVar5 = FUN_10016a50((byte *)0x0);
  piVar3[0xb] = (int)puVar5;
  pcVar6 = (char *)FUN_10015370((int)piVar3,param_1,0);
  if (pcVar6 == (char *)0x0) {
    FUN_10016bb0((LPVOID)piVar3[0xb]);
    uVar7 = FUN_100212e0(piVar3);
    return CONCAT22((short)((uint)uVar7 >> 0x10),0xfffd);
  }
  FUN_10020750((int)piVar1,piVar3,pcVar6);
  iVar9 = 0;
  iVar4 = piVar3[0x13];
  while (iVar4 == 0) {
    FUN_10020510(piVar1,piVar3,*(int *)(piVar3[0x10] + 4));
    FUN_10016c00((int)piVar3);
    FUN_10024ed0(*(short **)(piVar3[0x10] + 4),*(int *)(piVar3[0x11] + 4),piVar3[0xe] / 2,iVar9,
                 (int)piVar3);
    FUN_1001f660(piVar2,iVar8,0,*(char **)(piVar3[0x11] + 4),1,
                 (int)(piVar3[0xe] + (piVar3[0xe] >> 0x1f & 3U)) >> 2);
    iVar8 = iVar8 + ((int)(piVar3[0xe] + (piVar3[0xe] >> 0x1f & 3U)) >> 2);
    iVar9 = iVar9 + 1;
    iVar4 = piVar3[0x13];
  }
  FUN_1001f550(piVar2);
  FUN_10016bb0((LPVOID)piVar3[0xb]);
  uVar7 = FUN_100212e0(piVar3);
  return CONCAT22((short)((uint)uVar7 >> 0x10),1);
}



/* ===== FUN_10017420 ===== */
/* Entry: 10017420 */

undefined4 __cdecl
FUN_10017420(byte *param_1,int param_2,uint param_3,int *param_4,int param_5,int param_6,int param_7
            ,undefined4 param_8,int param_9,int param_10)

{
  byte *pbVar1;
  undefined2 uVar10;
  int *piVar2;
  int *piVar3;
  int *piVar4;
  undefined4 *puVar5;
  char *pcVar6;
  undefined4 uVar7;
  int iVar8;
  uint uVar9;
  int iVar11;
  uint uVar12;
  
  piVar4 = param_4;
  pbVar1 = param_1;
  if (((int)param_4 < 0) || (1 < (int)param_4)) {
    param_4 = (int *)0x0;
  }
  param_4 = (int *)(&DAT_100942d8)[(int)param_4];
  uVar10 = (undefined2)((uint)param_4 >> 0x10);
  if (param_4 == (int *)0x0) {
    return 0xfffb;
  }
  if (param_1 == (byte *)0x0) {
    return CONCAT22(uVar10,0xfffd);
  }
  if (*param_1 == 0) {
    return CONCAT22(uVar10,0xfffc);
  }
  piVar2 = FUN_1001efe0(param_2,&DAT_1006ffc8);
  if (piVar2 == (int *)0x0) {
    return 0xfffa;
  }
  piVar3 = (int *)FUN_1001f9f0((undefined2 *)&param_1,param_9);
  if ((short)param_1 != 1) {
    return CONCAT22((short)((uint)piVar3 >> 0x10),0xfffe);
  }
  if (((int)piVar4 < 0) || (1 < (int)piVar4)) {
    piVar4 = (int *)0x0;
  }
  *(char *)(piVar3 + 0x14a38) = (char)piVar4;
  FUN_10021f10(param_4,(int)piVar3,param_5,param_6,param_7,param_8,param_10);
  puVar5 = FUN_10016a50((byte *)0x0);
  uVar9 = param_3;
  piVar3[0xb] = (int)puVar5;
  *(short *)((int)piVar3 + 0x114eb2) = (short)param_3;
  FUN_10017ce0((int)piVar3);
  pcVar6 = (char *)FUN_10015370((int)piVar3,pbVar1,0);
  piVar4 = param_4;
  if (pcVar6 == (char *)0x0) {
    FUN_10016bb0((LPVOID)piVar3[0xb]);
    uVar7 = FUN_100212e0(piVar3);
    return CONCAT22((short)((uint)uVar7 >> 0x10),0xfffd);
  }
  FUN_10020750((int)param_4,piVar3,pcVar6);
  FUN_100180f0(piVar3,piVar2);
  param_2 = 0;
  if (piVar3[0x13] == 0) {
    param_9 = (uVar9 & 0xffff) - 6;
    do {
      iVar11 = param_9;
      FUN_10020510(piVar4,piVar3,*(int *)(piVar3[0x10] + 4));
      FUN_10016c00((int)piVar3);
      FUN_100180f0(piVar3,piVar2);
      switch(iVar11) {
      case 0:
        iVar11 = 0;
        if (0 < piVar3[0xe] / 2) {
          do {
            iVar11 = iVar11 + 1;
            *(undefined1 *)(*(int *)(piVar3[0x11] + 4) + -1 + iVar11) =
                 *(undefined1 *)
                  (*(short *)(*(int *)(piVar3[0x10] + 4) + -2 + iVar11 * 2) + 0x1848c +
                  piVar4[0xe11]);
          } while (iVar11 < piVar3[0xe] / 2);
        }
        goto LAB_100176e4;
      case 1:
        iVar11 = 0;
        if (0 < piVar3[0xe] / 2) {
          do {
            iVar11 = iVar11 + 1;
            *(undefined1 *)(*(int *)(piVar3[0x11] + 4) + -1 + iVar11) =
                 *(undefined1 *)
                  (*(short *)(*(int *)(piVar3[0x10] + 4) + -2 + iVar11 * 2) + 0x848c + piVar4[0xe11]
                  );
          } while (iVar11 < piVar3[0xe] / 2);
        }
LAB_100176e4:
        pcVar6 = *(char **)(piVar3[0x11] + 4);
        uVar9 = piVar3[0xe] / 2;
        uVar12 = 1;
LAB_1001772d:
        FUN_1001f660(piVar2,0,2,pcVar6,uVar12,uVar9);
        break;
      default:
        if (DAT_10070974 == 0) {
          uVar9 = piVar3[0xe] / 2;
          pcVar6 = *(char **)(piVar3[0x10] + 4);
          uVar12 = 2;
          goto LAB_1001772d;
        }
        FUN_1001f700(piVar2,0,2,*(char **)(piVar3[0x10] + 4),2,piVar3[0xe] / 2);
        break;
      case 0xb:
        FUN_10015320(*(short **)(piVar3[0x10] + 4),*(byte **)(piVar3[0x11] + 4),piVar3[0xe] / 2,
                     param_2,(int)piVar3);
        FUN_1001f660(piVar2,0,2,*(char **)(piVar3[0x11] + 4),1,
                     (int)(piVar3[0xe] + (piVar3[0xe] >> 0x1f & 3U)) >> 2);
        break;
      case 0x93:
        iVar11 = 0;
        if (0 < piVar3[0xe] / 2) {
          do {
            iVar8 = (int)*(short *)(*(int *)(piVar3[0x10] + 4) + iVar11 * 2);
            param_4 = (int *)CONCAT31(param_4._1_3_,
                                      (char)(iVar8 + (iVar8 >> 0x1f & 0xffU) >> 8) + -0x80);
            FUN_10058950((undefined4 *)(*(int *)(piVar3[0x11] + 4) + iVar11),&param_4,1);
            iVar11 = iVar11 + 1;
          } while (iVar11 < piVar3[0xe] / 2);
        }
        uVar9 = piVar3[0xe] / 2;
        uVar12 = 1;
        pcVar6 = *(char **)(piVar3[0x11] + 4);
        goto LAB_1001772d;
      }
      param_2 = param_2 + 1;
    } while (piVar3[0x13] == 0);
  }
  FUN_1001f550(piVar2);
  FUN_10016bb0((LPVOID)piVar3[0xb]);
  uVar7 = FUN_100212e0(piVar3);
  return CONCAT22((short)((uint)uVar7 >> 0x10),1);
}



/* ===== FUN_10020750 ===== */
/* Entry: 10020750 */

undefined4 __cdecl FUN_10020750(int param_1,int *param_2,char *param_3)

{
  char cVar1;
  int iVar2;
  undefined4 uVar3;
  uint uVar4;
  
  uVar4 = 0xffffffff;
  param_2[1] = 0;
  param_2[2] = (int)param_3;
  do {
    if (uVar4 == 0) break;
    uVar4 = uVar4 - 1;
    cVar1 = *param_3;
    param_3 = param_3 + 1;
  } while (cVar1 != '\0');
  *param_2 = ~uVar4 - 1;
  param_2[0x13] = 0;
  uVar3 = FUN_100202f0(param_1,param_2);
  iVar2 = param_2[0x13];
  while( true ) {
    if (iVar2 != 0) {
      return CONCAT22((short)((uint)uVar3 >> 0x10),0xffff);
    }
    if ((-1 < (short)uVar3) && (*(short *)(param_2[0x15] + 2) != 0)) break;
    uVar3 = FUN_100202f0(param_1,param_2);
    iVar2 = param_2[0x13];
  }
  return param_2[0x15] & 0xffff0000;
}



/* ===== FUN_100202f0 ===== */
/* Entry: 100202f0 */

undefined4 __cdecl FUN_100202f0(int param_1,int *param_2)

{
  uint *puVar1;
  uint uVar2;
  int iVar3;
  undefined4 uVar4;
  
  do {
    uVar2 = FUN_1001c180(param_1,param_2);
    param_2[1] = uVar2;
    if (((int)uVar2 < 0) || (param_2[0x13] != 0)) goto LAB_1002044f;
    puVar1 = *(uint **)(param_2[0x15] + 0x52870);
    if ((int)*puVar1 < 0) {
      if (param_2[0x4b2c5] < 0) {
        uVar2 = param_2[0x4b2c4];
        if ((int)uVar2 < 0) {
          if ((char)param_2[0x47829] != '\0') {
            *puVar1 = 0;
            goto LAB_10020362;
          }
          uVar2 = *(uint *)(param_1 + 0x3840);
        }
        *puVar1 = uVar2;
      }
      else {
        *puVar1 = ~-(uint)((char)param_2[0x47829] != '\0') & param_2[0x4b2c5];
      }
    }
LAB_10020362:
    if (0xffff < **(int **)(param_2[0x15] + 0x52870)) {
      **(int **)(param_2[0x15] + 0x52870) = 0x10000;
    }
  } while ((char)param_2[8] == '\0');
  iVar3 = FUN_10010620((int)param_2);
  if (iVar3 == 0) {
    return 0xffff;
  }
  if ((char)param_2[8] == '\x01') {
    uVar2 = param_2[0x15];
    *(undefined2 *)(param_2 + 0x12) = 0;
    *(undefined2 *)(uVar2 + 0x5285c) = 0;
    return uVar2 & 0xffff0000;
  }
  uVar2 = FUN_1000f3e0(param_1,(int)param_2);
  if ((short)uVar2 == 0) {
    return CONCAT22((short)(uVar2 >> 0x10),0xffff);
  }
  iVar3 = FUN_1001b800((int)param_2);
  if (iVar3 == 0) {
    return 0xffff;
  }
  if (*(char *)((int)param_2 + 0x21) == '\x03') {
    FUN_1001bcc0();
  }
  if (*(char *)((int)param_2 + 0x21) == '\x05') {
    FUN_1001be90(param_1,(int)param_2);
  }
  if ((*(char *)((int)param_2 + 0x21) == '\x06') || (*(char *)((int)param_2 + 0x21) == '\a')) {
    FUN_1001bf60(param_1,(int)param_2);
  }
  uVar4 = FUN_1000f5a0(param_1,(int)param_2);
  if ((short)uVar4 == 0) {
    return CONCAT22((short)((uint)uVar4 >> 0x10),0xffff);
  }
  uVar2 = FUN_1000fd30(param_1,(int)param_2);
  if ((short)uVar2 == 0) {
    return CONCAT22((short)(uVar2 >> 0x10),0xffff);
  }
  *(undefined2 *)(param_2 + 0x12) = 0;
  *(undefined2 *)(param_2[0x15] + 0x5285c) = 0;
LAB_1002044f:
  return uVar2 & 0xffff0000;
}



/* ===== FUN_1001f8c0 ===== */
/* Entry: 1001f8c0 */

uint __cdecl FUN_1001f8c0(HWND param_1,undefined2 *param_2)

{
  uint uVar1;
  undefined2 extraout_var;
  
  DAT_100942d4 = (undefined1 *)FUN_10016490(0x2049c);
  if (DAT_100942d4 == (undefined1 *)0x0) {
    *param_2 = 2;
    return (uint)param_2 & 0xffff0000;
  }
  DAT_100942d4[0x20494] = 0;
  *(undefined4 *)(DAT_100942d4 + 0x20490) = 1;
  FUN_10014f70(DAT_100942d4);
  uVar1 = FUN_1000e930((int)DAT_100942d4,param_2,param_1);
  if ((short)uVar1 < 0) {
    return uVar1 & 0xffff0000;
  }
  *param_2 = 0;
  DAT_100962e0 = 1;
  FUN_1000e890();
  return CONCAT22(extraout_var,1);
}



/* ===== FUN_1001f8b0 ===== */
/* Entry: 1001f8b0 */

undefined4 FUN_1001f8b0(void)

{
  return 0;
}



/* ===== FUN_1000e930 ===== */
/* Entry: 1000e930 */

undefined4 __cdecl FUN_1000e930(int param_1,undefined2 *param_2,HWND param_3)

{
  short sVar1;
  undefined4 uVar2;
  undefined2 extraout_var;
  uint uVar3;
  int *piVar4;
  undefined2 extraout_var_00;
  undefined4 extraout_ECX;
  undefined4 extraout_EDX;
  undefined1 local_200 [512];
  
  FUN_100588f6(local_200,(byte *)s__sdict_eng_1006ce4c);
  FUN_1000e890();
  FUN_10021c60(1,param_3,0);
  uVar2 = FUN_1000e820(param_1,param_2);
  if ((short)uVar2 < 0) {
    return CONCAT22((short)((uint)uVar2 >> 0x10),0xffff);
  }
  FUN_1000e890();
  FUN_1000e890();
  FUN_1000e890();
  sVar1 = FUN_1000e880();
  if (sVar1 < 0) {
    return CONCAT22(extraout_var,0xffff);
  }
  FUN_1000e890();
  FUN_1000e890();
  FUN_1000e890();
  FUN_10021c60(1,param_3,10);
  uVar3 = FUN_1003cab0(local_200);
  if ((short)uVar3 == 0) {
    *param_2 = 5;
    return CONCAT22((short)(uVar3 >> 0x10),0xffff);
  }
  FUN_1000e890();
  FUN_1000e890();
  FUN_1000e890();
  FUN_10021c60(1,param_3,0x14);
  FUN_1000e8a0();
  FUN_1000e890();
  FUN_1000e890();
  FUN_1000e890();
  FUN_10021c60(1,param_3,0x1e);
  FUN_1000e890();
  FUN_1000e890();
  FUN_10021c60(1,param_3,0x23);
  FUN_10025060(param_1);
  FUN_100250c0(extraout_ECX,extraout_EDX,param_1);
  FUN_10025090(param_1);
  piVar4 = FUN_1002b1e0(&PTR_DAT_10071cf8);
  *(int **)(param_1 + 0x20498) = piVar4;
  FUN_1000e890();
  FUN_1000e890();
  FUN_1000e890();
  FUN_10021c60(1,param_3,0x28);
  return CONCAT22(extraout_var_00,1);
}



/* ===== FUN_1000e820 ===== */
/* Entry: 1000e820 */

undefined4 __cdecl FUN_1000e820(undefined4 param_1,undefined2 *param_2)

{
  int iVar1;
  undefined1 local_200 [512];
  
  FUN_100588f6(local_200,(byte *)s__sdict_eng_1006ce4c);
  iVar1 = FUN_10003550(local_200);
  if (iVar1 < 0) {
    *param_2 = 3;
    return CONCAT22((short)((uint)param_2 >> 0x10),0xffff);
  }
  return CONCAT22((short)((uint)iVar1 >> 0x10),1);
}



/* ===== FUN_1000e880 ===== */
/* Entry: 1000e880 */

undefined2 FUN_1000e880(void)

{
  return 1;
}



/* ===== FUN_1000e8a0 ===== */
/* Entry: 1000e8a0 */

int FUN_1000e8a0(void)

{
  undefined1 local_200 [512];
  
  if (DAT_100952e0 != (uint *)0x0) {
    return CONCAT22((short)((uint)DAT_100952e0 >> 0x10),1);
  }
  FUN_100588f6(local_200,&DAT_1006ce58);
  DAT_100952e0 = FUN_1003bd50(local_200);
  return (-(uint)(DAT_100952e0 != (uint *)0x0) & 2) - 1;
}



/* ===== FUN_10003550 ===== */
/* Entry: 10003550 */

int __cdecl FUN_10003550(undefined4 param_1)

{
  int iVar1;
  
  iVar1 = FUN_100035a0(param_1);
  return ((-1 < iVar1) - 1 & 0xfffffffe) + 1;
}



/* ===== FUN_10003710 ===== */
/* Entry: 10003710 */

int * __cdecl FUN_10003710(undefined4 param_1,undefined4 param_2,short param_3,undefined1 param_4)

{
  int *piVar1;
  char *pcVar2;
  int iVar3;
  int iVar4;
  uint uVar5;
  int local_208;
  uint local_204;
  undefined1 local_200 [512];
  
  FUN_100588f6(local_200,(byte *)s__s__s_1006b31c);
  piVar1 = (int *)FUN_10016490(0x14);
  if (piVar1 == (int *)0x0) {
    return (int *)0x0;
  }
  pcVar2 = FUN_10024b90(local_200,&local_204,1);
  piVar1[2] = (int)pcVar2;
  if (pcVar2 == (char *)0x0) {
    return (int *)0x0;
  }
  pcVar2[local_204] = '\0';
  iVar3 = FUN_1000d070((char *)piVar1[2]);
  if (iVar3 == 0) {
    return (int *)0x0;
  }
  local_208 = piVar1[2];
  *piVar1 = iVar3;
  iVar3 = piVar1[2];
  iVar4 = FUN_10016490(*piVar1 << 2);
  piVar1[3] = iVar4;
  iVar4 = FUN_10016490(*piVar1 << 2);
  piVar1[4] = iVar4;
  iVar4 = 0;
  if (0 < *piVar1) {
    do {
      pcVar2 = FUN_10003890(&local_208,(char *)(iVar3 + local_204));
      if (pcVar2 == (char *)0x0) {
        return (int *)0x0;
      }
      if (param_3 == 0x4c) {
        *(char **)(piVar1[3] + iVar4 * 4) = pcVar2;
LAB_1000381c:
        *(char **)(piVar1[4] + iVar4 * 4) = pcVar2;
      }
      else if (param_3 == 0x4d) {
        *(char **)(piVar1[3] + iVar4 * 4) = pcVar2;
        pcVar2 = FUN_10003890(&local_208,(char *)(iVar3 + local_204));
        if (pcVar2 == (char *)0x0) {
          return (int *)0x0;
        }
        goto LAB_1000381c;
      }
      iVar4 = iVar4 + 1;
    } while (iVar4 < *piVar1);
  }
  *(undefined1 *)(piVar1 + 1) = param_4;
  uVar5 = FUN_10003920(piVar1);
  if ((short)uVar5 == 0) {
    FUN_10003970((undefined *)piVar1);
  }
  return piVar1;
}



/* ===== FUN_1003cab0 ===== */
/* Entry: 1003cab0 */

uint __cdecl FUN_1003cab0(undefined4 param_1)

{
  uint uVar1;
  undefined4 uVar2;
  undefined1 local_200 [512];
  
  DAT_10096680 = FUN_10003710(param_1,s_streeta_txt2_100750b8,0x4d,0x49);
  if (DAT_10096680 == (int *)0x0) {
    return 0;
  }
  DAT_10096684 = FUN_10003710(param_1,s_streetf_txt2_100750a8,0x4c,0x49);
  if (DAT_10096684 == (int *)0x0) {
    return 0;
  }
  DAT_10096688 = FUN_10003710(param_1,s_city1_txt2_1007509c,0x4d,0x49);
  if (DAT_10096688 == (int *)0x0) {
    return 0;
  }
  DAT_1009668c = FUN_10003710(param_1,s_city2_txt2_10075090,0x4d,0x49);
  if (DAT_1009668c == (int *)0x0) {
    return 0;
  }
  DAT_10096690 = FUN_10003710(param_1,s_city3_txt2_10075084,0x4d,0x49);
  if (DAT_10096690 == (int *)0x0) {
    return 0;
  }
  DAT_10096694 = FUN_10003710(param_1,s_city4_txt2_10075078,0x4d,0x49);
  if (DAT_10096694 == (int *)0x0) {
    return 0;
  }
  DAT_10096698 = FUN_10003710(param_1,s_citya_txt2_1007506c,0x4d,0x49);
  if (DAT_10096698 == (int *)0x0) {
    return 0;
  }
  DAT_1009669c = FUN_10003710(param_1,s_sbdw_txt2_10075060,0x4d,0x49);
  if (DAT_1009669c == (int *)0x0) {
    return 0;
  }
  DAT_100966a0 = FUN_10003710(param_1,s_abbrh_txt2_10075054,0x4d,0x49);
  if (DAT_100966a0 == (int *)0x0) {
    return 0;
  }
  DAT_100966a4 = FUN_10003710(param_1,s_abbrt_txt2_10075048,0x4d,0x49);
  if (DAT_100966a4 == (int *)0x0) {
    return 0;
  }
  DAT_100966a8 = FUN_10003710(param_1,s_abbrc_txt2_1007503c,0x4d,0x49);
  if (DAT_100966a8 == (int *)0x0) {
    return 0;
  }
  FUN_100588f6(local_200,(byte *)s__s_sbd_tree2_1007502c);
  uVar1 = FUN_100017d0((int *)&DAT_100966ac,local_200);
  if ((short)uVar1 == 0) {
    return uVar1;
  }
  uVar2 = FUN_1003c3a0(0x10096680);
  return (uint)((short)uVar2 != 0);
}



/* ===== FUN_1003c3a0 ===== */
/* Entry: 1003c3a0 */

undefined4 __cdecl FUN_1003c3a0(int param_1)

{
  int iVar1;
  undefined4 *puVar2;
  undefined4 uVar3;
  
  puVar2 = (undefined4 *)FUN_10016490(0x14);
  *(undefined4 **)(param_1 + 0x48) = puVar2;
  if (puVar2 == (undefined4 *)0x0) {
    return 0;
  }
  *puVar2 = 8;
  uVar3 = FUN_10016490(0x20);
  *(undefined4 *)(*(int *)(param_1 + 0x48) + 0x10) = uVar3;
  uVar3 = FUN_10016490(0x20);
  *(undefined4 *)(*(int *)(param_1 + 0x48) + 4) = uVar3;
  uVar3 = FUN_10016490(0x10);
  *(undefined4 *)(*(int *)(param_1 + 0x48) + 8) = uVar3;
  uVar3 = FUN_10016490(0x10);
  *(undefined4 *)(*(int *)(param_1 + 0x48) + 0xc) = uVar3;
  **(undefined4 **)(*(int *)(param_1 + 0x48) + 0x10) = &DAT_1007afa0;
  *(undefined **)(*(int *)(*(int *)(param_1 + 0x48) + 0x10) + 4) = &DAT_1007d280;
  *(undefined **)(*(int *)(*(int *)(param_1 + 0x48) + 0x10) + 8) = &DAT_1007d9a0;
  *(undefined **)(*(int *)(*(int *)(param_1 + 0x48) + 0x10) + 0xc) = &DAT_1007e600;
  *(undefined **)(*(int *)(*(int *)(param_1 + 0x48) + 0x10) + 0x10) = &DAT_10080100;
  *(undefined **)(*(int *)(*(int *)(param_1 + 0x48) + 0x10) + 0x14) = &DAT_10080220;
  *(undefined **)(*(int *)(*(int *)(param_1 + 0x48) + 0x10) + 0x18) = &DAT_100807c0;
  *(undefined **)(*(int *)(*(int *)(param_1 + 0x48) + 0x10) + 0x1c) = &DAT_10080ac0;
  **(undefined4 **)(*(int *)(param_1 + 0x48) + 4) = DAT_1008aac8;
  *(undefined4 *)(*(int *)(*(int *)(param_1 + 0x48) + 4) + 4) = DAT_1008aacc;
  *(undefined4 *)(*(int *)(*(int *)(param_1 + 0x48) + 4) + 8) = DAT_1008aad0;
  *(undefined4 *)(*(int *)(*(int *)(param_1 + 0x48) + 4) + 0xc) = DAT_1008aad4;
  *(undefined4 *)(*(int *)(*(int *)(param_1 + 0x48) + 4) + 0x10) = DAT_1008aad8;
  *(undefined4 *)(*(int *)(*(int *)(param_1 + 0x48) + 4) + 0x14) = DAT_1008aadc;
  *(undefined4 *)(*(int *)(*(int *)(param_1 + 0x48) + 4) + 0x18) = DAT_1008aae0;
  *(undefined4 *)(*(int *)(*(int *)(param_1 + 0x48) + 4) + 0x1c) = DAT_1008aae4;
  **(undefined2 **)(*(int *)(param_1 + 0x48) + 8) = 0x49;
  *(undefined2 *)(*(int *)(*(int *)(param_1 + 0x48) + 8) + 2) = 0x49;
  *(undefined2 *)(*(int *)(*(int *)(param_1 + 0x48) + 8) + 4) = 0x53;
  *(undefined2 *)(*(int *)(*(int *)(param_1 + 0x48) + 8) + 6) = 0x49;
  *(undefined2 *)(*(int *)(*(int *)(param_1 + 0x48) + 8) + 8) = 0x53;
  *(undefined2 *)(*(int *)(*(int *)(param_1 + 0x48) + 8) + 10) = 0x49;
  *(undefined2 *)(*(int *)(*(int *)(param_1 + 0x48) + 8) + 0xc) = 0x53;
  *(undefined2 *)(*(int *)(*(int *)(param_1 + 0x48) + 8) + 0xe) = 0x4d;
  puVar2 = *(undefined4 **)(*(int *)(param_1 + 0x48) + 0xc);
  *puVar2 = 0;
  puVar2[1] = 0;
  puVar2[2] = 0;
  puVar2[3] = 0;
  puVar2 = (undefined4 *)FUN_10016490(0x14);
  *(undefined4 **)(param_1 + 0x4c) = puVar2;
  if (puVar2 == (undefined4 *)0x0) {
    return 0;
  }
  *puVar2 = 4;
  uVar3 = FUN_10016490(0x10);
  *(undefined4 *)(*(int *)(param_1 + 0x4c) + 0x10) = uVar3;
  uVar3 = FUN_10016490(0x10);
  *(undefined4 *)(*(int *)(param_1 + 0x4c) + 4) = uVar3;
  uVar3 = FUN_10016490(8);
  *(undefined4 *)(*(int *)(param_1 + 0x4c) + 8) = uVar3;
  uVar3 = FUN_10016490(8);
  *(undefined4 *)(*(int *)(param_1 + 0x4c) + 0xc) = uVar3;
  **(undefined4 **)(*(int *)(param_1 + 0x4c) + 0x10) = &DAT_1007afa0;
  *(undefined **)(*(int *)(*(int *)(param_1 + 0x4c) + 0x10) + 4) = &DAT_1007d280;
  *(undefined **)(*(int *)(*(int *)(param_1 + 0x4c) + 0x10) + 8) = &DAT_1007d9a0;
  *(undefined **)(*(int *)(*(int *)(param_1 + 0x4c) + 0x10) + 0xc) = &DAT_10080ac0;
  **(undefined4 **)(*(int *)(param_1 + 0x4c) + 4) = DAT_1008aac8;
  *(undefined4 *)(*(int *)(*(int *)(param_1 + 0x4c) + 4) + 4) = DAT_1008aacc;
  *(undefined4 *)(*(int *)(*(int *)(param_1 + 0x4c) + 4) + 8) = DAT_1008aad0;
  *(undefined4 *)(*(int *)(*(int *)(param_1 + 0x4c) + 4) + 0xc) = DAT_1008aae4;
  **(undefined2 **)(*(int *)(param_1 + 0x4c) + 8) = 0x49;
  *(undefined2 *)(*(int *)(*(int *)(param_1 + 0x4c) + 8) + 2) = 0x49;
  *(undefined2 *)(*(int *)(*(int *)(param_1 + 0x4c) + 8) + 4) = 0x53;
  *(undefined2 *)(*(int *)(*(int *)(param_1 + 0x4c) + 8) + 6) = 0x4d;
  puVar2 = *(undefined4 **)(*(int *)(param_1 + 0x4c) + 0xc);
  *puVar2 = 0;
  puVar2[1] = 0;
  puVar2 = (undefined4 *)FUN_10016490(0x14);
  *(undefined4 **)(param_1 + 0x50) = puVar2;
  if (puVar2 == (undefined4 *)0x0) {
    return 0;
  }
  *puVar2 = 3;
  uVar3 = FUN_10016490(0xc);
  *(undefined4 *)(*(int *)(param_1 + 0x50) + 0x10) = uVar3;
  uVar3 = FUN_10016490(0xc);
  *(undefined4 *)(*(int *)(param_1 + 0x50) + 4) = uVar3;
  uVar3 = FUN_10016490(6);
  *(undefined4 *)(*(int *)(param_1 + 0x50) + 8) = uVar3;
  uVar3 = FUN_10016490(6);
  *(undefined4 *)(*(int *)(param_1 + 0x50) + 0xc) = uVar3;
  **(undefined4 **)(*(int *)(param_1 + 0x50) + 0x10) = &DAT_1007e600;
  *(undefined **)(*(int *)(*(int *)(param_1 + 0x50) + 0x10) + 4) = &DAT_10080100;
  *(undefined **)(*(int *)(*(int *)(param_1 + 0x50) + 0x10) + 8) = &DAT_10080ac0;
  **(undefined4 **)(*(int *)(param_1 + 0x50) + 4) = DAT_1008aad4;
  *(undefined4 *)(*(int *)(*(int *)(param_1 + 0x50) + 4) + 4) = DAT_1008aad8;
  *(undefined4 *)(*(int *)(*(int *)(param_1 + 0x50) + 4) + 8) = DAT_1008aae4;
  **(undefined2 **)(*(int *)(param_1 + 0x50) + 8) = 0x49;
  *(undefined2 *)(*(int *)(*(int *)(param_1 + 0x50) + 8) + 2) = 0x53;
  *(undefined2 *)(*(int *)(*(int *)(param_1 + 0x50) + 8) + 4) = 0x4d;
  puVar2 = *(undefined4 **)(*(int *)(param_1 + 0x50) + 0xc);
  *puVar2 = 0;
  *(undefined2 *)(puVar2 + 1) = 0;
  puVar2 = (undefined4 *)FUN_10016490(0x14);
  *(undefined4 **)(param_1 + 0x54) = puVar2;
  if (puVar2 == (undefined4 *)0x0) {
    return 0;
  }
  *puVar2 = 3;
  uVar3 = FUN_10016490(0xc);
  *(undefined4 *)(*(int *)(param_1 + 0x54) + 0x10) = uVar3;
  uVar3 = FUN_10016490(0xc);
  *(undefined4 *)(*(int *)(param_1 + 0x54) + 4) = uVar3;
  uVar3 = FUN_10016490(6);
  *(undefined4 *)(*(int *)(param_1 + 0x54) + 8) = uVar3;
  uVar3 = FUN_10016490(6);
  *(undefined4 *)(*(int *)(param_1 + 0x54) + 0xc) = uVar3;
  **(undefined4 **)(*(int *)(param_1 + 0x54) + 0x10) = &DAT_10080220;
  *(undefined **)(*(int *)(*(int *)(param_1 + 0x54) + 0x10) + 4) = &DAT_100807c0;
  *(undefined **)(*(int *)(*(int *)(param_1 + 0x54) + 0x10) + 8) = &DAT_10080ac0;
  **(undefined4 **)(*(int *)(param_1 + 0x54) + 4) = DAT_1008aadc;
  *(undefined4 *)(*(int *)(*(int *)(param_1 + 0x54) + 4) + 4) = DAT_1008aae0;
  *(undefined4 *)(*(int *)(*(int *)(param_1 + 0x54) + 4) + 8) = DAT_1008aae4;
  **(undefined2 **)(*(int *)(param_1 + 0x54) + 8) = 0x49;
  *(undefined2 *)(*(int *)(*(int *)(param_1 + 0x54) + 8) + 2) = 0x53;
  *(undefined2 *)(*(int *)(*(int *)(param_1 + 0x54) + 8) + 4) = 0x4d;
  puVar2 = *(undefined4 **)(*(int *)(param_1 + 0x54) + 0xc);
  *puVar2 = 0;
  *(undefined2 *)(puVar2 + 1) = 0;
  puVar2 = (undefined4 *)FUN_10016490(0x14);
  *(undefined4 **)(param_1 + 0x58) = puVar2;
  if (puVar2 == (undefined4 *)0x0) {
    return 0;
  }
  *puVar2 = 3;
  uVar3 = FUN_10016490(0xc);
  *(undefined4 *)(*(int *)(param_1 + 0x58) + 0x10) = uVar3;
  uVar3 = FUN_10016490(0xc);
  *(undefined4 *)(*(int *)(param_1 + 0x58) + 4) = uVar3;
  uVar3 = FUN_10016490(6);
  *(undefined4 *)(*(int *)(param_1 + 0x58) + 8) = uVar3;
  *(undefined4 *)(*(int *)(param_1 + 0x58) + 0xc) = 0;
  **(undefined4 **)(*(int *)(param_1 + 0x58) + 0x10) = &DAT_10088fe8;
  *(undefined **)(*(int *)(*(int *)(param_1 + 0x58) + 0x10) + 4) = &DAT_10089200;
  *(undefined **)(*(int *)(*(int *)(param_1 + 0x58) + 0x10) + 8) = &DAT_100893a8;
  **(undefined4 **)(*(int *)(param_1 + 0x58) + 4) = DAT_1008ab10;
  *(undefined4 *)(*(int *)(*(int *)(param_1 + 0x58) + 4) + 4) = DAT_1008ab18;
  *(undefined4 *)(*(int *)(*(int *)(param_1 + 0x58) + 4) + 8) = DAT_1008ab24;
  **(undefined2 **)(*(int *)(param_1 + 0x58) + 8) = 0x49;
  iVar1 = *(int *)(*(int *)(param_1 + 0x58) + 8);
  *(undefined2 *)(iVar1 + 2) = 0x49;
  *(undefined2 *)(*(int *)(*(int *)(param_1 + 0x58) + 8) + 4) = 0x49;
  return CONCAT22((short)((uint)iVar1 >> 0x10),1);
}



/* ===== FUN_10025060 ===== */
/* Entry: 10025060 */

undefined4 __cdecl FUN_10025060(int param_1)

{
  undefined2 uVar1;
  undefined2 extraout_var;
  uint uVar2;
  undefined2 *puVar3;
  
  uVar2 = 0;
  puVar3 = (undefined2 *)(param_1 + 0x28c);
  do {
    uVar1 = FUN_10024f70(uVar2);
    *puVar3 = uVar1;
    uVar2 = uVar2 + 1;
    puVar3 = puVar3 + 1;
  } while ((int)uVar2 < 0x100);
  return CONCAT22(extraout_var,1);
}



/* ===== FUN_100250c0 ===== */
/* Entry: 100250c0 */

undefined4 __fastcall FUN_100250c0(undefined4 param_1,undefined4 param_2,int param_3)

{
  uint uVar1;
  undefined4 extraout_ECX;
  undefined4 extraout_EDX;
  int iVar2;
  
  iVar2 = -0x8000;
  do {
    uVar1 = FUN_10024fe0(param_1,param_2,(short)iVar2);
    *(char *)(param_3 + 0x1848c + iVar2) = (char)uVar1;
    iVar2 = iVar2 + 1;
    param_1 = extraout_ECX;
    param_2 = extraout_EDX;
  } while (iVar2 < 0x8000);
  return CONCAT22((short)(uVar1 >> 0x10),1);
}



/* ===== FUN_10025090 ===== */
/* Entry: 10025090 */

undefined4 __cdecl FUN_10025090(int param_1)

{
  undefined4 uVar1;
  int iVar2;
  
  iVar2 = -0x8000;
  do {
    uVar1 = FUN_10024f90(iVar2);
    *(char *)(param_1 + 0x848c + iVar2) = (char)uVar1;
    iVar2 = iVar2 + 1;
  } while (iVar2 < 0x8000);
  return CONCAT22((short)((uint)uVar1 >> 0x10),1);
}



/* ===== FUN_10014f70 ===== */
/* Entry: 10014f70 */

void __cdecl FUN_10014f70(undefined1 *param_1)

{
  FUN_100588f6(param_1,(byte *)s__sdata_common__1006f2c4);
  return;
}



/* ===== FUN_10022430 ===== */
/* Entry: 10022430 */

DWORD __cdecl FUN_10022430(char *param_1)

{
  char cVar1;
  byte bVar2;
  byte *pbVar3;
  DWORD *pDVar4;
  int iVar5;
  byte *pbVar6;
  char *pcVar7;
  bool bVar8;
  char local_200 [512];
  
  if (DAT_10070974 < 0) {
    DAT_10070974 = FUN_1001f8b0();
  }
  if (DAT_1008ceb0 == 0) {
    FUN_1000e890();
  }
  if (DAT_100942cc == 0) {
    InitializeCriticalSection((LPCRITICAL_SECTION)&DAT_10096a00);
    DAT_100942cc = 1;
  }
  if (DAT_100940b4 == 0) {
    FUN_1001efd0();
    DAT_100940b4 = 1;
  }
  if (param_1 != (char *)0x0) {
    iVar5 = -1;
    pcVar7 = param_1;
    do {
      if (iVar5 == 0) break;
      iVar5 = iVar5 + -1;
      cVar1 = *pcVar7;
      pcVar7 = pcVar7 + 1;
    } while (cVar1 != '\0');
    if (iVar5 != -2) goto LAB_1002252c;
  }
  pbVar6 = &DAT_1008cbe8;
  pbVar3 = (byte *)&DAT_1008ccac;
  do {
    bVar2 = *pbVar3;
    bVar8 = bVar2 < *pbVar6;
    if (bVar2 != *pbVar6) {
LAB_100224e1:
      iVar5 = (1 - (uint)bVar8) - (uint)(bVar8 != 0);
      goto LAB_100224e6;
    }
    if (bVar2 == 0) break;
    bVar2 = pbVar3[1];
    bVar8 = bVar2 < pbVar6[1];
    if (bVar2 != pbVar6[1]) goto LAB_100224e1;
    pbVar3 = pbVar3 + 2;
    pbVar6 = pbVar6 + 2;
  } while (bVar2 != 0);
  iVar5 = 0;
LAB_100224e6:
  if (iVar5 == 0) {
    FUN_10021e80((BYTE *)&DAT_1008ccac);
    FUN_10021530((char *)&DAT_1008ccac);
  }
  FUN_10014f70(local_200);
  FUN_100588f6(local_200,&DAT_1006ce58);
  param_1 = local_200;
LAB_1002252c:
  iVar5 = FUN_10022560(s_VoiceText_10070a38,param_1,(byte *)s_License_10070a54);
  if (iVar5 < 0) {
    pDVar4 = FUN_10059b03();
    return *pDVar4;
  }
  return 0;
}



/* ===== FUN_10022560 ===== */
/* Entry: 10022560 */

undefined4 __cdecl FUN_10022560(char *param_1,char *param_2,byte *param_3)

{
  DWORD DVar1;
  char *pcVar2;
  int iVar3;
  DWORD *pDVar4;
  byte *pbVar5;
  undefined *puVar6;
  undefined4 *local_4;
  
  local_4 = (undefined4 *)0x0;
  pbVar5 = param_3;
  if (param_3 == (byte *)0x0) {
    pbVar5 = (byte *)s_License_10070a54;
  }
  if (DAT_100962f4 != (LPVOID)0x0) {
    FUN_10016500(DAT_100962f4);
    DAT_100962f4 = (LPVOID)0x0;
  }
  if (param_1 == (char *)0x0) {
    param_1 = s_VoiceText_10070a38;
  }
  DAT_100962f4 = (LPVOID)FUN_10014f40(param_1);
  if (param_2 == (char *)0x0) {
    param_2 = s_c__voicetext_eng_2004_data_commo_10070a74;
  }
  pcVar2 = FUN_10024b90(param_2,(uint *)&param_3,1);
  if (pcVar2 != (char *)0x0) {
    param_3[(int)pcVar2] = 0;
    iVar3 = FUN_10022810(pcVar2,(int)param_3,(int *)&local_4);
    if (iVar3 < 0) {
      pDVar4 = FUN_10059b03();
      DVar1 = *pDVar4;
      FUN_10022a20(local_4);
      pDVar4 = FUN_10059b03();
      *pDVar4 = DVar1;
      return 0xffffffff;
    }
    pDVar4 = FUN_10059b03();
    DVar1 = *pDVar4;
    if (DAT_100962f8 != (LPVOID)0x0) {
      FUN_10016500(DAT_100962f8);
    }
    if (DAT_100962fc != (byte *)0x0) {
      FUN_10016500(DAT_100962fc);
    }
    DAT_100962f8 = (char *)FUN_100229e0(local_4,pbVar5);
    if ((DAT_100962f8 != (char *)0x0) &&
       (DAT_100962fc = (byte *)_strchr(DAT_100962f8,0x3a), DAT_100962fc != (byte *)0x0)) {
      *DAT_100962fc = 0;
      DAT_100962fc = DAT_100962fc + 1;
      pbVar5 = (byte *)FUN_100229e0(local_4,DAT_100962fc);
      if (pbVar5 != (byte *)0x0) {
        DAT_10096300 = 1;
        DAT_100962fc = pbVar5;
      }
    }
    puVar6 = FUN_10022010();
    if ((-1 < (int)puVar6) && (iVar3 = FUN_10022710(), -1 < iVar3)) {
      DAT_100962f8 = (char *)FUN_10014f40(DAT_100962f8);
      DAT_100962fc = (byte *)FUN_10014f40((char *)DAT_100962fc);
      FUN_10022a20(local_4);
      return 0;
    }
    FUN_10022a20(local_4);
    DAT_100962f8 = (LPVOID)0x0;
    DAT_100962fc = (byte *)0x0;
    if (DVar1 != 0) {
      pDVar4 = FUN_10059b03();
      *pDVar4 = DVar1;
    }
  }
  return 0xffffffff;
}



/* ===== FUN_1001c180 ===== */
/* Entry: 1001c180 */

int __cdecl FUN_1001c180(int param_1,int *param_2)

{
  undefined4 *puVar1;
  int *piVar2;
  int iVar3;
  int iVar4;
  int *piVar5;
  int iVar6;
  int iVar7;
  int iVar8;
  
  piVar2 = (int *)param_2[0x4b2c2];
  if (param_2[2] == 0) {
    return 0;
  }
  while( true ) {
    FUN_1001c740(param_1,(int)param_2);
    FUN_100546d0(piVar2,(char *)(param_2[2] + param_2[1]));
    piVar5 = (int *)*piVar2;
    iVar8 = 0;
    iVar3 = *piVar5;
    if (0 < piVar5[1]) {
      iVar6 = 0;
      iVar7 = 0;
      do {
        *(int *)(param_2[0x15] + 0x728 + iVar6) = *(int *)(iVar7 + 0x14 + (int)piVar5) + param_2[1];
        iVar4 = *(int *)(*piVar2 + iVar7 + 0x18);
        if (*(int *)(*piVar2 + iVar7 + 0x14) < iVar4) {
          *(int *)(param_2[0x15] + 0x72c + iVar6) = iVar4 + -1 + param_2[1];
        }
        else {
          *(undefined4 *)(param_2[0x15] + iVar6 + 0x72c) =
               *(undefined4 *)(param_2[0x15] + iVar6 + 0x728);
        }
        piVar5 = (int *)*piVar2;
        iVar8 = iVar8 + 1;
        iVar7 = iVar7 + 0x88;
        iVar6 = iVar6 + 0x2fc;
      } while (iVar8 < piVar5[1]);
    }
    iVar8 = 0;
    *(undefined2 *)(param_2[0x15] + 2) = *(undefined2 *)(*piVar2 + 4);
    if (0 < *(short *)(param_2[0x15] + 2)) {
      piVar5 = param_2 + 0x473f8;
      do {
        piVar5[-200] = param_2[0x4732a];
        *piVar5 = param_2[0x47328];
        piVar5[200] = param_2[0x47329];
        piVar5[400] = -1;
        piVar5[600] = -1;
        iVar8 = iVar8 + 1;
        piVar5 = piVar5 + 1;
      } while (iVar8 < *(short *)(param_2[0x15] + 2));
    }
    if (iVar3 == 0) break;
    if (0 < *(int *)(*piVar2 + 4)) {
      FUN_1001c440((int)param_2);
      *(undefined4 *)(param_2[0x15] + 0x52864) = *(undefined4 *)(param_2[0x15] + 0x728);
      iVar8 = param_2[0x15];
      *(undefined4 *)(iVar8 + 0x52868) =
           *(undefined4 *)(iVar8 + 0x430 + *(short *)(iVar8 + 2) * 0x2fc);
      if (*(char *)(DAT_100942d4 + 0x20494) == '\x01') {
        iVar8 = param_2[0x15];
        iVar6 = 0;
        if (0 < *(short *)(iVar8 + 2)) {
          iVar7 = 0;
          do {
            iVar6 = iVar6 + 1;
            *(undefined4 *)(iVar7 + 0x728 + iVar8) =
                 *(undefined4 *)
                  (param_2[0x14a39] +
                  *(int *)(param_2[0x14a3a] + *(int *)(iVar7 + 0x728 + iVar8) * 4) * 4);
            piVar2 = (int *)(param_2[0x15] + 0x72c + iVar7);
            puVar1 = (undefined4 *)(param_2[0x15] + 0x72c + iVar7);
            iVar7 = iVar7 + 0x2fc;
            *puVar1 = *(undefined4 *)
                       (param_2[0x14a39] + *(int *)(param_2[0x14a3b] + *piVar2 * 4) * 4);
            iVar8 = param_2[0x15];
          } while (iVar6 < *(short *)(iVar8 + 2));
          iVar8 = param_2[1];
          param_2[1] = iVar8 + iVar3;
          return iVar8 + iVar3;
        }
      }
      else {
        iVar8 = param_2[0x15];
        iVar6 = 0;
        if (0 < *(short *)(iVar8 + 2)) {
          iVar7 = 0;
          do {
            if (*param_2 <= *(int *)(iVar7 + 0x728 + iVar8)) {
              *(int *)(iVar7 + 0x728 + iVar8) = *param_2 + -1;
            }
            if (*param_2 <= *(int *)(iVar7 + 0x72c + param_2[0x15])) {
              *(int *)(iVar7 + 0x72c + param_2[0x15]) = *param_2 + -1;
            }
            if (*(int *)(iVar7 + 0x728 + param_2[0x15]) < 0) {
              *(undefined4 *)(iVar7 + 0x728 + param_2[0x15]) = 0;
            }
            if (*(int *)(iVar7 + 0x72c + param_2[0x15]) < 0) {
              *(undefined4 *)(iVar7 + 0x72c + param_2[0x15]) = 0;
            }
            iVar6 = iVar6 + 1;
            *(undefined4 *)(iVar7 + 0x728 + param_2[0x15]) =
                 *(undefined4 *)(param_2[0x14a3a] + *(int *)(iVar7 + 0x728 + param_2[0x15]) * 4);
            iVar8 = iVar7 + 0x72c;
            iVar4 = iVar7 + 0x72c;
            iVar7 = iVar7 + 0x2fc;
            *(undefined4 *)(iVar4 + param_2[0x15]) =
                 *(undefined4 *)(param_2[0x14a3b] + *(int *)(iVar8 + param_2[0x15]) * 4);
            iVar8 = param_2[0x15];
          } while (iVar6 < *(short *)(iVar8 + 2));
        }
      }
      iVar8 = param_2[1];
      param_2[1] = iVar8 + iVar3;
      return iVar8 + iVar3;
    }
    param_2[1] = param_2[1] + iVar3;
  }
  param_2[0x13] = 1;
  return param_2[1];
}



/* ===== FUN_1001c440 ===== */
/* Entry: 1001c440 */

undefined4 __cdecl FUN_1001c440(int param_1)

{
  int iVar1;
  int iVar2;
  int iVar3;
  int iVar4;
  bool bVar5;
  int *piVar6;
  int iVar7;
  int iVar8;
  int iVar9;
  int iVar10;
  int iVar11;
  int iVar12;
  int iVar13;
  int local_24;
  int local_18;
  int local_c;
  int local_8;
  
  local_18 = 0;
  iVar9 = (int)*(short *)(*(int *)(param_1 + 0x54) + 2);
  do {
    if (((local_18 == 0) || (local_18 == 1)) || (local_18 == 2)) {
      iVar1 = *(int *)(param_1 + 0x11dc64 + local_18 * 0xc);
      iVar2 = *(int *)(param_1 + (local_18 * 3 + 0x47718) * 4);
      iVar3 = *(int *)(param_1 + local_18 * 0xc + 0x11dc68);
      if ((iVar1 != 0) && (iVar3 != 0)) {
        iVar11 = (&DAT_10070630)[local_18];
        iVar8 = (&DAT_10070618)[local_18];
        if (local_18 == 0) {
          local_24 = param_1 + 0x11ccc0;
        }
        else if (local_18 == 1) {
          local_24 = param_1 + 0x11cfe0;
        }
        else if (local_18 == 2) {
          local_24 = param_1 + 0x11d300;
        }
        iVar10 = 0;
        iVar13 = 0;
        local_c = 0;
        if (0 < iVar2) {
          do {
            if (iVar10 < iVar9) {
              piVar6 = (int *)(*(int *)(param_1 + 0x54) + 0x728 + iVar10 * 0x2fc);
              do {
                if (*(int *)(iVar1 + local_c * 4) <= *piVar6) break;
                iVar10 = iVar10 + 1;
                piVar6 = piVar6 + 0xbf;
              } while (iVar10 < iVar9);
            }
            if (local_c == 0) {
              if (iVar13 < iVar10) {
                iVar13 = iVar10;
              }
            }
            else {
              iVar4 = *(int *)(iVar3 + -4 + local_c * 4);
              if (iVar13 < iVar10) {
                piVar6 = (int *)(local_24 + iVar13 * 4);
                iVar12 = iVar10 - iVar13;
                iVar13 = iVar13 + iVar12;
                do {
                  iVar7 = iVar4;
                  if (iVar4 == 100) {
                    iVar7 = *piVar6;
                  }
                  if (iVar7 < iVar11) {
                    iVar7 = iVar11;
                  }
                  if (iVar8 < iVar7) {
                    iVar7 = iVar8;
                  }
                  *piVar6 = iVar7;
                  piVar6 = piVar6 + 1;
                  iVar12 = iVar12 + -1;
                } while (iVar12 != 0);
              }
            }
            if ((local_c == iVar2 + -1) && (iVar13 < iVar9)) {
              iVar13 = iVar9;
            }
            local_c = local_c + 1;
          } while (local_c < iVar2);
        }
      }
    }
    local_18 = local_18 + 1;
  } while (local_18 < 5);
  local_18 = 0;
  do {
    if ((local_18 == 3) || (local_18 == 4)) {
      iVar1 = *(int *)(param_1 + (local_18 * 3 + 0x47718) * 4);
      iVar2 = *(int *)(param_1 + 0x11dc64 + local_18 * 0xc);
      iVar3 = *(int *)(param_1 + 0x11dc68 + local_18 * 0xc);
      if ((iVar2 != 0) && (iVar3 != 0)) {
        if (local_18 == 3) {
          local_24 = param_1 + 0x11d620;
        }
        else if (local_18 == 4) {
          local_24 = param_1 + 0x11d940;
        }
        iVar11 = 0;
        if (0 < iVar9) {
          do {
            if (iVar11 == 0) {
              local_c = *(int *)(param_1 + 4);
              local_8 = *(int *)(*(int *)(param_1 + 0x54) + 0x728);
            }
            else {
              iVar8 = *(int *)(param_1 + 0x54) + iVar11 * 0x2fc;
              local_c = *(int *)(iVar8 + 0x430);
              local_8 = *(int *)(iVar8 + 0x728);
            }
            iVar8 = 0;
            bVar5 = false;
LAB_1001c682:
            iVar13 = iVar8;
            if (iVar8 < iVar1) {
              piVar6 = (int *)(iVar2 + iVar8 * 4);
              do {
                if ((local_c <= *piVar6) && (iVar13 = iVar8, *piVar6 <= local_8)) break;
                iVar8 = iVar8 + 1;
                piVar6 = piVar6 + 1;
                iVar13 = iVar8;
              } while (iVar8 < iVar1);
            }
            if (iVar13 != iVar1) {
              if (local_18 == 3) {
                iVar8 = *(int *)(iVar3 + iVar13 * 4);
                if (bVar5) {
                  piVar6 = (int *)(local_24 + iVar11 * 4);
                  *piVar6 = *piVar6 + iVar8;
                }
                else {
                  *(int *)(local_24 + iVar11 * 4) = iVar8;
                }
                iVar8 = iVar13 + 1;
                if (*(int *)(local_24 + iVar11 * 4) < DAT_1007063c) {
                  *(int *)(local_24 + iVar11 * 4) = DAT_1007063c;
                }
                if (*(int *)(local_24 + iVar11 * 4) <= DAT_10070624) goto LAB_1001c6ff;
                *(int *)(local_24 + iVar11 * 4) = DAT_10070624;
                bVar5 = true;
              }
              else {
                *(undefined4 *)(local_24 + iVar11 * 4) = *(undefined4 *)(iVar3 + iVar13 * 4);
LAB_1001c6ff:
                iVar8 = iVar13 + 1;
                bVar5 = true;
              }
              goto LAB_1001c682;
            }
            iVar11 = iVar11 + 1;
          } while (iVar11 < iVar9);
        }
      }
    }
    local_18 = local_18 + 1;
    if (4 < local_18) {
      return 1;
    }
  } while( true );
}



/* ===== FUN_1001c740 ===== */
/* Entry: 1001c740 */

void __cdecl FUN_1001c740(int param_1,int param_2)

{
  undefined4 uVar1;
  
  if (*(char *)(param_2 + 0x11cc83) == '\x01') {
    uVar1 = *(undefined4 *)(param_2 + 0x11cc88);
  }
  else {
    uVar1 = *(undefined4 *)(param_1 + 0x3834);
  }
  *(undefined4 *)(param_2 + 0x11cc94) = uVar1;
  if (*(char *)(param_2 + 0x11cc84) == '\x01') {
    uVar1 = *(undefined4 *)(param_2 + 0x11cc8c);
  }
  else {
    uVar1 = *(undefined4 *)(param_1 + 0x3830);
  }
  *(undefined4 *)(param_2 + 0x11cc98) = uVar1;
  if (*(char *)(param_2 + 0x11cc85) == '\x01') {
    uVar1 = *(undefined4 *)(param_2 + 0x11cc90);
  }
  else {
    uVar1 = *(undefined4 *)(param_1 + 0x3838);
  }
  *(undefined4 *)(param_2 + 0x11cc9c) = uVar1;
  if (200 < *(int *)(param_2 + 0x11cc94)) {
    *(undefined4 *)(param_2 + 0x11cc94) = 200;
  }
  if (*(int *)(param_2 + 0x11cc94) < 0x32) {
    *(undefined4 *)(param_2 + 0x11cc94) = 0x32;
  }
  if (200 < *(int *)(param_2 + 0x11cc98)) {
    *(undefined4 *)(param_2 + 0x11cc98) = 200;
  }
  if (*(int *)(param_2 + 0x11cc98) < 0x19) {
    *(undefined4 *)(param_2 + 0x11cc98) = 0x19;
  }
  if (500 < *(int *)(param_2 + 0x11cc9c)) {
    *(undefined4 *)(param_2 + 0x11cc9c) = 500;
  }
  if (*(int *)(param_2 + 0x11cc9c) < 0) {
    *(undefined4 *)(param_2 + 0x11cc9c) = 0;
  }
  *(undefined4 *)(param_2 + 0x11cca8) = *(undefined4 *)(param_2 + 0x11cc94);
  *(undefined4 *)(param_2 + 0x11cca0) = *(undefined4 *)(param_2 + 0x11cc98);
  *(undefined4 *)(param_2 + 0x11cca4) = *(undefined4 *)(param_2 + 0x11cc9c);
  *(undefined4 *)(param_2 + 0x12cb10) = 0xffffffff;
  *(undefined4 *)(param_2 + 0x11d620) = 0xffffffff;
  return;
}



/* ===== FUN_100546d0 ===== */
/* Entry: 100546d0 */

void __cdecl FUN_100546d0(int *param_1,char *param_2)

{
  int iVar1;
  
  FUN_10054700(param_1,&param_2);
  iVar1 = FUN_100547d0(param_2,param_1);
  FUN_10055050(*param_1,(short)iVar1);
  return;
}



/* ===== FUN_10054700 ===== */
/* Entry: 10054700 */

void __cdecl FUN_10054700(undefined4 *param_1,undefined4 *param_2)

{
  char cVar1;
  int *piVar2;
  int iVar3;
  uint uVar4;
  int iVar5;
  int *piVar6;
  int *piVar7;
  undefined4 *puVar8;
  
  piVar7 = (int *)*param_1;
  piVar2 = (int *)param_1[1];
  iVar5 = 0;
  *piVar7 = 0;
  piVar7[1] = 0;
  piVar6 = piVar7 + 3;
  for (iVar3 = 0xd48; iVar3 != 0; iVar3 = iVar3 + -1) {
    *piVar6 = 0;
    piVar6 = piVar6 + 1;
  }
  while( true ) {
    cVar1 = *(char *)*param_2;
    if ((((cVar1 != ' ') && (cVar1 != '\t')) && (cVar1 != '\n')) && (cVar1 != '\r')) break;
    iVar5 = iVar5 + 1;
    *param_2 = (char *)*param_2 + 1;
  }
  *piVar7 = iVar5;
  piVar7 = piVar2;
  for (uVar4 = (uint)(*piVar2 * 0x38) >> 2; piVar7 = piVar7 + 1, uVar4 != 0; uVar4 = uVar4 - 1) {
    *piVar7 = 0;
  }
  for (iVar3 = 0; iVar3 != 0; iVar3 = iVar3 + -1) {
    *(undefined1 *)piVar7 = 0;
    piVar7 = (int *)((int)piVar7 + 1);
  }
  *piVar2 = 0;
  *(undefined4 *)param_1[2] = 0;
  puVar8 = (undefined4 *)param_1[2];
  for (iVar3 = 0x19; puVar8 = puVar8 + 1, iVar3 != 0; iVar3 = iVar3 + -1) {
    *puVar8 = 0xffffffff;
  }
  puVar8 = (undefined4 *)param_1[2];
  for (iVar3 = 0x578; iVar3 != 0; iVar3 = iVar3 + -1) {
    *puVar8 = 0;
    puVar8 = puVar8 + 1;
  }
  *(undefined4 *)param_1[3] = 0;
  *(undefined4 *)(param_1[3] + 4) = 0x14;
  puVar8 = (undefined4 *)(param_1[3] + 8);
  for (uVar4 = *(int *)(param_1[3] + 4) * 0x43 & 0x3fffffff; uVar4 != 0; uVar4 = uVar4 - 1) {
    *puVar8 = 0;
    puVar8 = puVar8 + 1;
  }
  for (iVar3 = 0; iVar3 != 0; iVar3 = iVar3 + -1) {
    *(undefined1 *)puVar8 = 0;
    puVar8 = (undefined4 *)((int)puVar8 + 1);
  }
  return;
}



/* ===== FUN_100547d0 ===== */
/* Entry: 100547d0 */

int __cdecl FUN_100547d0(char *param_1,int *param_2)

{
  byte bVar1;
  int *piVar2;
  uint *puVar3;
  uint *puVar4;
  int *piVar5;
  int iVar6;
  short sVar7;
  char *pcVar8;
  undefined4 uVar9;
  uint uVar10;
  int *piVar11;
  char *pcVar12;
  byte *pbVar13;
  int *piVar14;
  int iVar15;
  uint *puVar16;
  byte *pbVar17;
  bool bVar18;
  char cVar19;
  
  piVar2 = (int *)*param_2;
  puVar3 = (uint *)param_2[2];
  puVar4 = (uint *)param_2[3];
  piVar5 = (int *)param_2[4];
  pcVar8 = param_1 + -*piVar2;
  param_2 = (int *)param_2[1];
  do {
    while (uVar9 = FUN_10058300((int)piVar5), (short)uVar9 != 0) {
      uVar10 = FUN_10055e60(piVar2,param_2,piVar5,(int *)&param_1,pcVar8);
      sVar7 = (short)uVar10;
      if ((((sVar7 == 2) || (sVar7 == 3)) || (sVar7 == 4)) || (sVar7 == -1)) goto LAB_1005503b;
      if (sVar7 != 1) break;
    }
    *puVar3 = 0;
    puVar16 = puVar3;
    for (iVar15 = 0x19; puVar16 = puVar16 + 1, iVar15 != 0; iVar15 = iVar15 + -1) {
      *puVar16 = 0xffffffff;
    }
    puVar16 = puVar3;
    for (iVar15 = 0x578; iVar15 != 0; iVar15 = iVar15 + -1) {
      *puVar16 = 0;
      puVar16 = puVar16 + 1;
    }
    *puVar4 = 0;
    puVar4[1] = 0x14;
    puVar16 = puVar4 + 2;
    for (iVar15 = 0x53c; iVar15 != 0; iVar15 = iVar15 + -1) {
      *puVar16 = 0;
      puVar16 = puVar16 + 1;
    }
    piVar11 = FUN_10054420((int *)puVar3,(int)piVar5,(int *)&param_1,pcVar8,'D');
    if (*piVar11 == 8) {
      iVar15 = piVar2[1];
      if ((0 < iVar15) &&
         (uVar9 = FUN_10057a70((char *)((int)piVar2 + iVar15 * 0x88 + -0x56),
                               piVar2[iVar15 * 0x22 + -0x1b]), (short)uVar9 == 0)) {
        if ((char)piVar11[6] == '?') {
          FUN_100579e0(piVar2,piVar2[1] + -1,0x3f,piVar11[4]);
        }
        else {
          FUN_100579e0(piVar2,piVar2[1] + -1,0x2e,piVar11[4]);
        }
      }
      *piVar2 = piVar11[4];
      if (piVar2[1] != 0) {
LAB_10055036:
        sVar7 = 1;
LAB_1005503b:
        return (int)sVar7;
      }
    }
    else {
      if (*piVar11 == 9) {
        iVar15 = piVar2[1];
        if ((0 < iVar15) &&
           (uVar9 = FUN_10057a70((char *)((int)piVar2 + iVar15 * 0x88 + -0x56),
                                 piVar2[iVar15 * 0x22 + -0x1b]), (short)uVar9 == 0)) {
          FUN_100579e0(piVar2,piVar2[1] + -1,0x2e,piVar2[piVar2[1] * 0x22 + -0x1c]);
        }
        *piVar2 = piVar11[4];
        return 0;
      }
      iVar15 = -1;
      piVar14 = piVar11 + 6;
      do {
        if (iVar15 == 0) break;
        iVar15 = iVar15 + -1;
        iVar6 = *piVar14;
        piVar14 = (int *)((int)piVar14 + 1);
      } while ((char)iVar6 != '\0');
      if (((iVar15 == -3) && (cVar19 = (char)piVar11[6], cVar19 != '\0')) &&
         (pcVar12 = _strchr(&DAT_10075500,(int)cVar19), pcVar12 != (char *)0x0)) {
LAB_100549b8:
        iVar15 = FUN_1004edc0(param_2,(int *)puVar3,(int)piVar5,piVar11,param_1);
        if ((short)iVar15 == 0x50) {
          uVar9 = FUN_10054220(param_2,piVar11,0);
          sVar7 = (short)uVar9;
          if (sVar7 != 1) goto LAB_1005503b;
          if (0 < piVar2[1]) {
            cVar19 = (char)piVar11[6];
            if ((cVar19 == '.') || (cVar19 == '?')) {
              iVar15 = piVar11[4];
            }
            else {
              iVar15 = piVar11[4];
              cVar19 = '.';
            }
            FUN_100579e0(piVar2,piVar2[1] + -1,cVar19,iVar15);
            goto LAB_10055036;
          }
        }
        else {
          piVar14 = FUN_10057930(param_2,(int *)puVar3,*puVar3 - 1,0xffffffff);
          if ((((piVar11[2] != 0) || (2 < piVar11[5])) || ((char)piVar11[6] != '.')) &&
             ((piVar14 == (int *)0x0 || (piVar14[1] == 0)))) goto LAB_10054a74;
          uVar9 = FUN_10054220(param_2,piVar11,0);
          sVar7 = (short)uVar9;
          if (((sVar7 == 2) || ((sVar7 == 3 || (sVar7 == 4)))) || (sVar7 == -1)) goto LAB_1005503b;
        }
      }
      else {
        iVar15 = -1;
        piVar14 = piVar11 + 6;
        do {
          if (iVar15 == 0) break;
          iVar15 = iVar15 + -1;
          iVar6 = *piVar14;
          piVar14 = (int *)((int)piVar14 + 1);
        } while ((char)iVar6 != '\0');
        if (iVar15 == -5) {
          pbVar17 = &DAT_100754fc;
          pbVar13 = (byte *)(piVar11 + 6);
          do {
            bVar1 = *pbVar17;
            bVar18 = bVar1 < *pbVar13;
            if (bVar1 != *pbVar13) {
LAB_100549ab:
              iVar15 = (1 - (uint)bVar18) - (uint)(bVar18 != 0);
              goto LAB_100549b0;
            }
            if (bVar1 == 0) break;
            bVar1 = pbVar17[1];
            bVar18 = bVar1 < pbVar13[1];
            if (bVar1 != pbVar13[1]) goto LAB_100549ab;
            pbVar17 = pbVar17 + 2;
            pbVar13 = pbVar13 + 2;
          } while (bVar1 != 0);
          iVar15 = 0;
LAB_100549b0:
          if (iVar15 == 0) goto LAB_100549b8;
        }
LAB_10054a74:
        uVar10 = FUN_1004a0b0(piVar2,param_2,puVar3,puVar4,piVar5,(uint *)&param_1,pcVar8);
        sVar7 = (short)uVar10;
        if (sVar7 != 1) {
          if ((((sVar7 == 2) || (sVar7 == 3)) || (sVar7 == 4)) || (sVar7 == -1)) goto LAB_1005503b;
          uVar9 = FUN_10050290(piVar2,param_2,(int *)puVar3,puVar4,(int)piVar5,(int *)&param_1,
                               pcVar8);
          sVar7 = (short)uVar9;
          if (sVar7 != 1) {
            if (((sVar7 == 2) || (sVar7 == 3)) || ((sVar7 == 4 || (sVar7 == -1))))
            goto LAB_1005503b;
            uVar9 = FUN_100344a0(piVar2,param_2,(int *)puVar3,puVar4,(int)piVar5,(int *)&param_1,
                                 pcVar8,0);
            sVar7 = (short)uVar9;
            if (sVar7 != 1) {
              if (((sVar7 == 2) || (sVar7 == 3)) || ((sVar7 == 4 || (sVar7 == -1))))
              goto LAB_1005503b;
              uVar10 = FUN_100552e0(piVar2,param_2,(int *)puVar3,puVar4,piVar5,(int *)&param_1,
                                    pcVar8);
              sVar7 = (short)uVar10;
              if (sVar7 != 1) {
                if ((((sVar7 == 2) || (sVar7 == 3)) || (sVar7 == 4)) || (sVar7 == -1))
                goto LAB_1005503b;
                uVar9 = FUN_100512a0(piVar2,param_2,(int *)puVar3,puVar4,(int)piVar5,(int *)&param_1
                                     ,pcVar8,0);
                sVar7 = (short)uVar9;
                if (sVar7 != 1) {
                  if (((sVar7 == 2) || (sVar7 == 3)) || ((sVar7 == 4 || (sVar7 == -1))))
                  goto LAB_1005503b;
                  uVar10 = FUN_10035d30(piVar2,param_2,(int *)puVar3,puVar4,(int)piVar5,
                                        (int *)&param_1,pcVar8,0);
                  sVar7 = (short)uVar10;
                  if (sVar7 != 1) {
                    if (((sVar7 == 2) || (sVar7 == 3)) || ((sVar7 == 4 || (sVar7 == -1))))
                    goto LAB_1005503b;
                    sVar7 = FUN_10040b50(piVar2,param_2,(int *)puVar3,puVar4,(int)piVar5,
                                         (int *)&param_1,pcVar8,0);
                    if (sVar7 != 1) {
                      if ((((sVar7 == 2) || (sVar7 == 3)) || (sVar7 == 4)) || (sVar7 == -1))
                      goto LAB_1005503b;
                      uVar9 = FUN_10042b70(piVar2,param_2,(int *)puVar3,puVar4,(int)piVar5,
                                           (int *)&param_1,pcVar8);
                      sVar7 = (short)uVar9;
                      if (sVar7 != 1) {
                        if (((sVar7 == 2) || (sVar7 == 3)) || ((sVar7 == 4 || (sVar7 == -1))))
                        goto LAB_1005503b;
                        uVar9 = FUN_1002f790(piVar2,param_2,puVar3,puVar4,(int)piVar5,
                                             (uint *)&param_1,pcVar8);
                        sVar7 = (short)uVar9;
                        if (sVar7 != 1) {
                          if (((sVar7 == 2) || (sVar7 == 3)) || ((sVar7 == 4 || (sVar7 == -1))))
                          goto LAB_1005503b;
                          uVar9 = FUN_1003d200(piVar2,param_2,(int *)puVar3,puVar4,(int)piVar5,
                                               (int *)&param_1,pcVar8,0);
                          sVar7 = (short)uVar9;
                          if (sVar7 != 1) {
                            if ((((sVar7 == 2) || (sVar7 == 3)) || (sVar7 == 4)) || (sVar7 == -1))
                            goto LAB_1005503b;
                            uVar9 = FUN_10040960(piVar2,param_2,(int *)puVar3,puVar4,(int *)&param_1
                                                 ,pcVar8,0);
                            sVar7 = (short)uVar9;
                            if (sVar7 != 1) {
                              if (((sVar7 == 2) || (sVar7 == 3)) || ((sVar7 == 4 || (sVar7 == -1))))
                              goto LAB_1005503b;
                              uVar9 = FUN_10043330(piVar2,param_2,(int *)puVar3,puVar4,(int)piVar5,
                                                   (int *)&param_1,pcVar8);
                              sVar7 = (short)uVar9;
                              if (sVar7 != 1) {
                                if (((sVar7 == 2) || (sVar7 == 3)) ||
                                   ((sVar7 == 4 || (sVar7 == -1)))) goto LAB_1005503b;
                                uVar9 = FUN_1003fa90(piVar2,param_2,(int *)puVar3,puVar4,(int)piVar5
                                                     ,(int *)&param_1,pcVar8);
                                sVar7 = (short)uVar9;
                                if (sVar7 != 1) {
                                  if ((((sVar7 == 2) || (sVar7 == 3)) || (sVar7 == 4)) ||
                                     (sVar7 == -1)) goto LAB_1005503b;
                                  uVar9 = FUN_1004f0c0(piVar2,param_2,(int *)puVar3,puVar4,
                                                       (int)piVar5,(int *)&param_1,pcVar8);
                                  sVar7 = (short)uVar9;
                                  if (sVar7 != 1) {
                                    if (((sVar7 == 2) || (sVar7 == 3)) ||
                                       ((sVar7 == 4 || (sVar7 == -1)))) goto LAB_1005503b;
                                    uVar9 = FUN_1002c320(piVar2,param_2,(int *)puVar3,puVar4,
                                                         (int)piVar5,(int *)&param_1,pcVar8);
                                    sVar7 = (short)uVar9;
                                    if (sVar7 != 1) {
                                      if (((sVar7 == 2) || (sVar7 == 3)) ||
                                         ((sVar7 == 4 || (sVar7 == -1)))) goto LAB_1005503b;
                                      uVar9 = FUN_1002c4f0(piVar2,param_2,(int *)puVar3,puVar4,
                                                           (int)piVar5,(int *)&param_1,pcVar8);
                                      sVar7 = (short)uVar9;
                                      if (sVar7 != 1) {
                                        if ((((sVar7 == 2) || (sVar7 == 3)) || (sVar7 == 4)) ||
                                           (sVar7 == -1)) goto LAB_1005503b;
                                        if (0x62 < piVar2[1]) goto LAB_10055036;
                                        FUN_10054220(param_2,piVar11,0);
                                        FUN_100479d0(piVar2,piVar11);
                                      }
                                    }
                                  }
                                }
                              }
                            }
                          }
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  } while( true );
}



/* ===== FUN_10055050 ===== */
/* Entry: 10055050 */

short __cdecl FUN_10055050(int param_1,short param_2)

{
  byte bVar1;
  char cVar2;
  byte *pbVar3;
  bool bVar4;
  int iVar5;
  char *pcVar6;
  uint uVar7;
  uint uVar8;
  byte *_SubStr;
  char *_Str;
  byte *pbVar9;
  char *pcVar10;
  byte *pbVar11;
  int local_c;
  int local_8;
  
  pbVar3 = *(byte **)(param_1 + 0x352c);
  bVar4 = false;
  iVar5 = FUN_1000b340(pbVar3,param_1 + 0xc,*(int *)(param_1 + 4),1);
  if (iVar5 < 0) {
    return -1;
  }
  local_c = 0;
  local_8 = 0;
  if (0 < *(short *)(pbVar3 + 0x429a2)) {
    _Str = (char *)(param_1 + -0x56);
    _SubStr = pbVar3 + 0x429ab;
    do {
      if ((local_8 < 1) || (*(short *)(_SubStr + -0x6f) != *(short *)(_SubStr + -3))) {
        _Str[0x84] = -1;
        if ((_Str[0x82] == 'U') || (_Str[0x83] != 'Y')) {
          uVar7 = 0xffffffff;
          pbVar9 = _SubStr + 0x1e;
          do {
            pbVar11 = pbVar9;
            if (uVar7 == 0) break;
            uVar7 = uVar7 - 1;
            pbVar11 = pbVar9 + 1;
            bVar1 = *pbVar9;
            pbVar9 = pbVar11;
          } while (bVar1 != 0);
          uVar7 = ~uVar7;
          pbVar9 = pbVar11 + -uVar7;
          pbVar11 = (byte *)(_Str + 0xa6);
          for (uVar8 = uVar7 >> 2; uVar8 != 0; uVar8 = uVar8 - 1) {
            *(undefined4 *)pbVar11 = *(undefined4 *)pbVar9;
            pbVar9 = pbVar9 + 4;
            pbVar11 = pbVar11 + 4;
          }
          for (uVar7 = uVar7 & 3; uVar7 != 0; uVar7 = uVar7 - 1) {
            *pbVar11 = *pbVar9;
            pbVar9 = pbVar9 + 1;
            pbVar11 = pbVar11 + 1;
          }
        }
        iVar5 = *(int *)(_Str + 0x76);
        bVar4 = false;
        *(int *)(_Str + 0x76) = iVar5 + 1;
        _Str[iVar5 + 0x85] = -1;
        _Str[0x66] = _SubStr[-1];
        local_c = local_c + 1;
        *(undefined2 *)(_Str + 100) = *(undefined2 *)(_SubStr + -5);
        _Str = _Str + 0x88;
      }
      else if ((*_SubStr == 0xa2) && (_SubStr[1] == 0xfe)) {
        if (!bVar4) {
          pcVar6 = _strstr(_Str,(char *)_SubStr);
          *pcVar6 = '\0';
          *(int *)(_Str + -0x16) = *(int *)(_Str + -0x16) + -2;
          bVar4 = true;
        }
      }
      else {
        if (bVar4) {
          uVar7 = 0xffffffff;
          pbVar9 = _SubStr;
          do {
            pbVar11 = pbVar9;
            if (uVar7 == 0) break;
            uVar7 = uVar7 - 1;
            pbVar11 = pbVar9 + 1;
            bVar1 = *pbVar9;
            pbVar9 = pbVar11;
          } while (bVar1 != 0);
          uVar7 = ~uVar7;
          iVar5 = -1;
          pcVar6 = _Str;
          do {
            pcVar10 = pcVar6;
            if (iVar5 == 0) break;
            iVar5 = iVar5 + -1;
            pcVar10 = pcVar6 + 1;
            cVar2 = *pcVar6;
            pcVar6 = pcVar10;
          } while (cVar2 != '\0');
          pbVar9 = pbVar11 + -uVar7;
          pbVar11 = (byte *)(pcVar10 + -1);
          for (uVar8 = uVar7 >> 2; uVar8 != 0; uVar8 = uVar8 - 1) {
            *(undefined4 *)pbVar11 = *(undefined4 *)pbVar9;
            pbVar9 = pbVar9 + 4;
            pbVar11 = pbVar11 + 4;
          }
          for (uVar7 = uVar7 & 3; uVar7 != 0; uVar7 = uVar7 - 1) {
            *pbVar11 = *pbVar9;
            pbVar9 = pbVar9 + 1;
            pbVar11 = pbVar11 + 1;
          }
        }
        if (*(int *)(_Str + -0x12) < 5) {
          if (_SubStr[0x1e] != 0) {
            uVar7 = 0xffffffff;
            pbVar9 = _SubStr + 0x1e;
            do {
              pbVar11 = pbVar9;
              if (uVar7 == 0) break;
              uVar7 = uVar7 - 1;
              pbVar11 = pbVar9 + 1;
              bVar1 = *pbVar9;
              pbVar9 = pbVar11;
            } while (bVar1 != 0);
            uVar7 = ~uVar7;
            iVar5 = -1;
            pcVar6 = _Str + 0x1e;
            do {
              pcVar10 = pcVar6;
              if (iVar5 == 0) break;
              iVar5 = iVar5 + -1;
              pcVar10 = pcVar6 + 1;
              cVar2 = *pcVar6;
              pcVar6 = pcVar10;
            } while (cVar2 != '\0');
            pbVar9 = pbVar11 + -uVar7;
            pbVar11 = (byte *)(pcVar10 + -1);
            for (uVar8 = uVar7 >> 2; uVar8 != 0; uVar8 = uVar8 - 1) {
              *(undefined4 *)pbVar11 = *(undefined4 *)pbVar9;
              pbVar9 = pbVar9 + 4;
              pbVar11 = pbVar11 + 4;
            }
            for (uVar7 = uVar7 & 3; uVar7 != 0; uVar7 = uVar7 - 1) {
              *pbVar11 = *pbVar9;
              pbVar9 = pbVar9 + 1;
              pbVar11 = pbVar11 + 1;
            }
          }
          _Str[*(int *)(_Str + -0x12) + -4] = -1;
          iVar5 = *(int *)(_Str + -0x12);
          *(int *)(_Str + -0x12) = iVar5 + 1;
          _Str[iVar5 + -3] = -1;
          *(ushort *)(_Str + -0x24) = *(ushort *)(_Str + -0x24) | *(ushort *)(_SubStr + -5);
        }
      }
      _SubStr = _SubStr + 0x6c;
      local_8 = local_8 + 1;
    } while (local_8 < *(short *)(pbVar3 + 0x429a2));
  }
  if (0 < *(short *)(pbVar3 + 0x429a2)) {
    *(undefined2 *)(param_1 + 8) = *(undefined2 *)(pbVar3 + 0x429a4);
  }
  if (*(int *)(param_1 + 4) == local_c) {
    if (param_2 != 0) {
      return (-(ushort)(param_2 != -1) & 2) - 1;
    }
    return 0;
  }
  return -1;
}



/* ===== FUN_1004edc0 ===== */
/* Entry: 1004edc0 */

int __cdecl FUN_1004edc0(int *param_1,int *param_2,int param_3,int *param_4,char *param_5)

{
  undefined2 uVar1;
  int iVar2;
  undefined4 local_20 [8];
  
  iVar2 = FUN_1004ee30(param_1,param_2,param_3,param_4,param_5);
  if ((short)iVar2 == 0x41) {
    FUN_1004ebc0(local_20,param_1,param_2,param_3,param_5);
    uVar1 = FUN_10001570(DAT_100966b0,(int)local_20);
    iVar2 = (-(uint)((char)uVar1 != '\0') & 2) + 0x4e;
  }
  return iVar2;
}



/* ===== FUN_1004a0b0 ===== */
/* Entry: 1004a0b0 */

uint __cdecl
FUN_1004a0b0(int *param_1,int *param_2,uint *param_3,uint *param_4,int *param_5,uint *param_6,
            char *param_7)

{
  char cVar1;
  undefined2 uVar2;
  uint uVar3;
  undefined2 extraout_var;
  
  if (param_3[*param_3 * 0xe + 0x57] == 4) {
    cVar1 = (char)(param_3 + *param_3 * 0xe + 0x57)[6];
    if ((cVar1 == 'A') || (cVar1 == 'f')) {
      uVar3 = FUN_1004a390(param_1,param_2,(int *)param_3,(int)param_5,(int *)param_6,param_7);
    }
    else if ((cVar1 == 'Y') || (cVar1 == 'Z')) {
      uVar3 = FUN_1004aaf0(param_1,param_2,(int *)param_3,param_5,(int *)param_6,param_7);
    }
    else if ((cVar1 < 'a') || ('e' < cVar1)) {
      if ((cVar1 < 'P') || ('V' < cVar1)) {
        if (cVar1 == 'l') {
          uVar3 = FUN_1004c3c0(param_1,param_2,(int *)param_3,(int)param_5,(int *)param_6,param_7);
        }
        else if (cVar1 == 'D') {
          uVar3 = FUN_100344a0(param_1,param_2,(int *)param_3,param_4,(int)param_5,(int *)param_6,
                               param_7,0x44);
        }
        else if (cVar1 == 'g') {
          uVar2 = FUN_10040b50(param_1,param_2,(int *)param_3,param_4,(int)param_5,(int *)param_6,
                               param_7,0x67);
          uVar3 = CONCAT22(extraout_var,uVar2);
        }
        else if ((cVar1 < 'h') || ('k' < cVar1)) {
          if (cVar1 == 'W') {
            uVar3 = FUN_1004c060(param_1,param_2,(int *)param_3,param_4,(int)param_5,(int *)param_6,
                                 param_7);
          }
          else if ((cVar1 < 'E') || ('O' < cVar1)) {
            if ((cVar1 < 'B') || ('C' < cVar1)) {
              if (cVar1 == 'X') {
                uVar3 = FUN_1004da00(param_1,param_2,(int *)param_3,(int)param_5,(int *)param_6,
                                     param_7);
              }
              else if (cVar1 == 'n') {
                uVar3 = FUN_1004e800(param_1,param_2,(int *)param_3,(int)param_5,(int *)param_6,
                                     param_7);
              }
              else {
                if (cVar1 != 'm') {
                  return (-(uint)(cVar1 != 'o') & 4) + 1;
                }
                uVar3 = FUN_1004e370(param_1,param_2,(int *)param_3,(int)param_5,(int *)param_6,
                                     param_7);
              }
            }
            else {
              uVar3 = FUN_1004c530(param_1,param_2,param_3,param_4,(int)param_5,param_6,param_7);
            }
          }
          else {
            uVar3 = FUN_10035d30(param_1,param_2,(int *)param_3,param_4,(int)param_5,(int *)param_6,
                                 param_7,(short)cVar1);
          }
        }
        else {
          uVar3 = FUN_100512a0(param_1,param_2,(int *)param_3,param_4,(int)param_5,(int *)param_6,
                               param_7,(short)cVar1);
        }
      }
      else {
        uVar3 = FUN_1004b690(param_1,param_2,(int *)param_3,(int)param_5,(int *)param_6,param_7);
      }
    }
    else {
      uVar3 = FUN_1004acf0(param_1,param_2,(int *)param_3,(int)param_5,(int *)param_6,param_7);
    }
    if ((short)uVar3 == 5) {
      return CONCAT22((short)(uVar3 >> 0x10),1);
    }
  }
  else {
    uVar3 = CONCAT22((short)((uint)(param_3 + *param_3 * 0xe + 0x57) >> 0x10),5);
  }
  return uVar3;
}



/* ===== FUN_10050290 ===== */
/* Entry: 10050290 */

undefined4 __cdecl
FUN_10050290(int *param_1,int *param_2,int *param_3,undefined4 param_4,int param_5,int *param_6,
            char *param_7)

{
  char cVar1;
  int iVar2;
  int iVar3;
  int iVar4;
  short sVar5;
  short sVar6;
  int iVar7;
  int *piVar8;
  undefined4 uVar9;
  uint uVar10;
  uint uVar11;
  int *piVar12;
  int *piVar13;
  char *pcVar14;
  char *pcVar15;
  int iVar16;
  byte *pbVar17;
  undefined2 local_100;
  undefined1 local_fe;
  
  iVar2 = *param_3;
  iVar7 = param_3[iVar2 * 0xe + 0x57];
  if (iVar7 != 6) goto LAB_100506bf;
  iVar3 = param_1[1];
  iVar4 = *param_2;
  cVar1 = (char)param_3[iVar2 * 0xe + 0x5d];
  if (cVar1 == -0x5c) {
    iVar7 = param_3[iVar2 * 0xe + 0x5b];
    iVar16 = param_3[iVar2 * 0xe + 0x5a];
    pcVar15 = s_currency_10074a34;
LAB_100502fa:
    iVar7 = FUN_10047b90(param_1,pcVar15,iVar16,iVar7,0x53,0x44);
LAB_10050300:
    sVar6 = (short)iVar7;
  }
  else {
    if (cVar1 == -0x57) {
      iVar7 = param_3[iVar2 * 0xe + 0x5b];
      iVar16 = param_3[iVar2 * 0xe + 0x5a];
      pcVar15 = s_Copyright_10075bc4;
      goto LAB_100502fa;
    }
    if (cVar1 == -0x4e) {
      iVar7 = FUN_10047ef0(param_1,s_superscript_two_10075cf0,param_3[iVar2 * 0xe + 0x5a],
                           param_3[iVar2 * 0xe + 0x5b],0x53,0x44);
      goto LAB_10050300;
    }
    if (cVar1 == -0x4d) {
      iVar7 = FUN_10047ef0(param_1,s_superscript_three_10075cdc,param_3[iVar2 * 0xe + 0x5a],
                           param_3[iVar2 * 0xe + 0x5b],0x53,0x44);
      goto LAB_10050300;
    }
    if (cVar1 == -0x4b) {
      iVar7 = param_3[iVar2 * 0xe + 0x5b];
      iVar16 = param_3[iVar2 * 0xe + 0x5a];
      pcVar15 = s_micro_10075cd4;
      goto LAB_100502fa;
    }
    if (cVar1 == -0x4a) {
      iVar7 = param_3[iVar2 * 0xe + 0x5b];
      iVar16 = param_3[iVar2 * 0xe + 0x5a];
      pcVar15 = s_Paragraph_10075cc8;
      goto LAB_100502fa;
    }
    if (cVar1 == -0x47) {
      iVar7 = FUN_10047ef0(param_1,s_superscript_one_10075cb8,param_3[iVar2 * 0xe + 0x5a],
                           param_3[iVar2 * 0xe + 0x5b],0x53,0x44);
      goto LAB_10050300;
    }
    if (cVar1 == -0x29) {
      iVar7 = param_3[iVar2 * 0xe + 0x5b];
      iVar16 = param_3[iVar2 * 0xe + 0x5a];
      pcVar15 = s_times_10075cb0;
      goto LAB_100502fa;
    }
    if ((((cVar1 != -0x37) && (cVar1 != -0x17)) && (cVar1 != -0x2f)) && (cVar1 != -0xf))
    goto LAB_10050637;
    piVar8 = FUN_10057930(param_2,param_3,iVar2 + -1,0xffffffff);
    if ((((piVar8 == (int *)0x0) || (param_3[iVar2 * 0xe + 0x59] != 0)) ||
        ((*piVar8 != 1 ||
         ((piVar8[1] != 0 || (iVar7 = FUN_100578e0((int)param_1,(int)piVar8), iVar7 < 0)))))) ||
       (0x1d < piVar8[5] + 3)) {
LAB_10050475:
      if (((char)param_3[iVar2 * 0xe + 0x5d] == -0x37) ||
         ((char)param_3[iVar2 * 0xe + 0x5d] == -0x17)) {
        local_100 = DAT_1006b6a4;
        local_fe = DAT_1006b6a6;
        sVar6 = (short)param_3[iVar2 * 0xe + 0x5a];
      }
      else {
        local_100 = DAT_10075c9c;
        local_fe = DAT_10075c9e;
        sVar6 = (short)param_3[iVar2 * 0xe + 0x5a];
      }
    }
    else {
      if (((char)param_3[iVar2 * 0xe + 0x5d] == -0x37) ||
         ((char)param_3[iVar2 * 0xe + 0x5d] == -0x17)) {
        pbVar17 = &DAT_10075ca0;
      }
      else {
        pbVar17 = &DAT_10075ca8;
      }
      FUN_100588f6((undefined1 *)&local_100,pbVar17);
      sVar6 = (short)piVar8[3];
      param_1[1] = iVar7;
      piVar8[1] = 0x1b;
      if (sVar6 == -1) goto LAB_10050475;
    }
    sVar5 = (short)param_3[iVar2 * 0xe + 0x5b];
    cVar1 = *(char *)*param_6;
joined_r0x100504d9:
    if ((((cVar1 == '\0') || (cVar1 == ' ')) || (cVar1 == '\t')) ||
       ((cVar1 == '\n' || (cVar1 == '\r')))) goto LAB_10050611;
    piVar8 = FUN_10054420(param_3,param_5,param_6,param_7,'D');
    if (piVar8 == (int *)0x0) {
      *param_2 = iVar4;
      *param_3 = iVar2;
      *param_6 = param_3[iVar2 + 1];
      param_1[1] = iVar3;
      if (iVar3 < 1) {
        iVar2 = *param_1;
        *param_1 = iVar2;
        return CONCAT22((short)((uint)iVar2 >> 0x10),4);
      }
      iVar2 = param_1[iVar3 * 0x22 + -0x1c];
      *param_1 = iVar2;
      return CONCAT22((short)((uint)iVar2 >> 0x10),4);
    }
    if (*piVar8 == 1) {
      piVar12 = piVar8 + 6;
      goto LAB_10050566;
    }
    cVar1 = (char)piVar8[6];
    if (cVar1 == -0x37) {
LAB_1005055d:
      piVar12 = (int *)&DAT_1006b6a4;
LAB_10050566:
      uVar10 = 0xffffffff;
      do {
        piVar13 = piVar12;
        if (uVar10 == 0) break;
        uVar10 = uVar10 - 1;
        piVar13 = (int *)((int)piVar12 + 1);
        iVar7 = *piVar12;
        piVar12 = piVar13;
      } while ((char)iVar7 != '\0');
      uVar10 = ~uVar10;
      iVar7 = -1;
      pcVar15 = (char *)&local_100;
      do {
        pcVar14 = pcVar15;
        if (iVar7 == 0) break;
        iVar7 = iVar7 + -1;
        pcVar14 = pcVar15 + 1;
        cVar1 = *pcVar15;
        pcVar15 = pcVar14;
      } while (cVar1 != '\0');
      pcVar15 = (char *)((int)piVar13 - uVar10);
      pcVar14 = pcVar14 + -1;
      for (uVar11 = uVar10 >> 2; uVar11 != 0; uVar11 = uVar11 - 1) {
        *(undefined4 *)pcVar14 = *(undefined4 *)pcVar15;
        pcVar15 = pcVar15 + 4;
        pcVar14 = pcVar14 + 4;
      }
      for (uVar10 = uVar10 & 3; uVar10 != 0; uVar10 = uVar10 - 1) {
        *pcVar14 = *pcVar15;
        pcVar15 = pcVar15 + 1;
        pcVar14 = pcVar14 + 1;
      }
      sVar5 = (short)piVar8[4];
      cVar1 = *(char *)*param_6;
      goto joined_r0x100504d9;
    }
    if (((cVar1 == -0x17) || (cVar1 == -0x2f)) || (cVar1 == -0xf)) {
      if ((cVar1 != -0x37) && (cVar1 != -0x17)) {
        piVar12 = (int *)&DAT_10075c9c;
        goto LAB_10050566;
      }
      goto LAB_1005055d;
    }
    FUN_100544b0(param_3,param_6);
LAB_10050611:
    iVar7 = FUN_10047b90(param_1,(char *)&local_100,(int)sVar6,(int)sVar5,0x41,0x41);
    sVar6 = (short)iVar7;
  }
  if (sVar6 == 1) {
LAB_10050637:
    uVar9 = FUN_10054280(param_2,param_3,0x1b);
    if ((short)uVar9 == 1) {
      return uVar9;
    }
    *param_3 = iVar2;
    *param_6 = param_3[iVar2 + 1];
    *param_2 = iVar4;
    param_1[1] = iVar3;
    if (iVar3 < 1) {
      *param_1 = *param_1;
      return uVar9;
    }
    *param_1 = param_1[iVar3 * 0x22 + -0x1c];
    return uVar9;
  }
  if (*param_3 != iVar2) {
    *param_3 = iVar2;
    iVar7 = param_3[iVar2 + 1];
    *param_6 = iVar7;
  }
LAB_100506bf:
  return CONCAT22((short)((uint)iVar7 >> 0x10),5);
}



/* ===== FUN_100344a0 ===== */
/* Entry: 100344a0 */

undefined4 __cdecl
FUN_100344a0(int *param_1,int *param_2,int *param_3,uint *param_4,int param_5,int *param_6,
            char *param_7,short param_8)

{
  int *piVar1;
  int iVar2;
  short sVar3;
  char cVar4;
  int *piVar5;
  undefined2 uVar11;
  int *piVar6;
  int iVar7;
  int *piVar8;
  char *pcVar9;
  undefined4 uVar10;
  int iVar12;
  bool bVar13;
  bool bVar14;
  bool bVar15;
  bool bVar16;
  bool bVar17;
  undefined2 in_stack_00000022;
  uint local_20;
  uint local_1c;
  uint local_18;
  byte *local_14;
  uint local_10;
  short local_8;
  uint local_4;
  
  piVar5 = param_3;
  sVar3 = 0;
  local_8 = 0;
  if (param_8 == 0) {
    piVar6 = (int *)*param_3;
    local_4 = param_3[(int)piVar6 * 0xe + 0x5a];
    piVar8 = param_3 + (int)piVar6 * 0xe + 0x57;
    param_3 = piVar6;
  }
  else {
    piVar6 = _param_8;
    if (param_8 != 0x44) goto LAB_100355ab;
    *param_3 = 0;
    piVar6 = param_3;
    for (iVar12 = 0x19; piVar6 = piVar6 + 1, iVar12 != 0; iVar12 = iVar12 + -1) {
      *piVar6 = -1;
    }
    piVar6 = param_3;
    for (iVar12 = 0x578; iVar12 != 0; iVar12 = iVar12 + -1) {
      *piVar6 = 0;
      piVar6 = piVar6 + 1;
    }
    piVar1 = (int *)*param_3;
    piVar8 = FUN_10054420(param_3,param_5,param_6,param_7,'D');
    piVar6 = (int *)0x0;
    if (piVar8 == (int *)0x0) goto LAB_100355ab;
    local_4 = piVar8[3];
    param_3 = piVar1;
  }
  iVar12 = param_1[1];
  iVar2 = *param_2;
  local_18 = 0xffffffff;
  local_10 = 0xffffffff;
  local_20 = 0xffffffff;
  local_1c = 0xffffffff;
  bVar13 = false;
  bVar15 = false;
  local_14 = (byte *)0x0;
  uVar11 = (undefined2)((uint)param_1 >> 0x10);
  if ((*piVar8 == 3) &&
     (local_1c = FUN_10050ac0(0x10075f10,(char *)(piVar8 + 6),DAT_1008aa6c,0x53), -1 < (int)local_1c
     )) {
    piVar6 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
    if (piVar6 == (int *)0x0) {
      *param_2 = iVar2;
      *piVar5 = (int)param_3;
      *param_6 = piVar5[(int)param_3 + 1];
      param_1[1] = iVar12;
      if (0 < iVar12) {
        *param_1 = param_1[iVar12 * 0x22 + -0x1c];
        return CONCAT22(uVar11,4);
      }
      goto LAB_100354ba;
    }
    if ((piVar6[2] < 2) && (*piVar6 == 2)) {
      local_14 = (byte *)(piVar6 + 6);
      piVar6 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
      if (piVar6 == (int *)0x0) {
        *param_2 = iVar2;
        *piVar5 = (int)param_3;
        *param_6 = piVar5[(int)param_3 + 1];
        param_1[1] = iVar12;
        if (0 < iVar12) {
          *param_1 = param_1[iVar12 * 0x22 + -0x1c];
          return CONCAT22(uVar11,4);
        }
        goto LAB_100354ba;
      }
      if ((piVar6[2] < 2) &&
         ((((*piVar6 == 3 || (*piVar6 == 7)) && (piVar6[5] == 1)) && ((char)piVar6[6] == '-')))) {
        bVar13 = true;
      }
      else {
        bVar13 = false;
        FUN_100544b0(piVar5,param_6);
      }
      piVar6 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
      if (piVar6 == (int *)0x0) {
        *param_2 = iVar2;
        *piVar5 = (int)param_3;
        *param_6 = piVar5[(int)param_3 + 1];
        goto LAB_10035368;
      }
      if (((piVar6[2] < 2) && (*piVar6 == 1)) &&
         (local_18 = FUN_10050a40(0x10075d00,(char *)(piVar6 + 6),DAT_1008aa68,0x49),
         -1 < (int)local_18)) {
        bVar15 = true;
        piVar6 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
        if (piVar6 == (int *)0x0) {
          *param_2 = iVar2;
          *piVar5 = (int)param_3;
          *param_6 = piVar5[(int)param_3 + 1];
          goto LAB_10035368;
        }
        if (((piVar6[2] != 0) || (2 < piVar6[5])) || ((char)piVar6[6] != '.')) {
          if ((((*piVar6 == 3) || (*piVar6 == 7)) && (piVar6[5] == 1)) && ((char)piVar6[6] == '-'))
          {
            piVar6 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
            iVar7 = _strncmp((char *)(piVar6 + 6),&DAT_1006b4e8,1);
            if ((((iVar7 == 0) ||
                 (iVar7 = _strncmp((char *)(piVar6 + 6),&DAT_1006c7d4,3), iVar7 == 0)) &&
                ((*(char *)((int)piVar6 + 0x19) == '-' || (*(char *)((int)piVar6 + 0x1b) == '-'))))
               && (piVar6[2] == 0)) {
              sVar3 = 1;
            }
            FUN_100544b0(piVar5,param_6);
          }
          goto LAB_10034809;
        }
      }
      else {
        FUN_100544b0(piVar5,param_6);
        if (bVar13) {
LAB_10034809:
          FUN_100544b0(piVar5,param_6);
        }
      }
      piVar8 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
      if (piVar8 == (int *)0x0) {
        *param_2 = iVar2;
        *piVar5 = (int)param_3;
        *param_6 = piVar5[(int)param_3 + 1];
        param_1[1] = iVar12;
        if (0 < iVar12) {
          *param_1 = param_1[iVar12 * 0x22 + -0x1c];
          return CONCAT22(uVar11,4);
        }
        goto LAB_100354ba;
      }
      if ((*piVar8 == 1) &&
         ((((1 < piVar8[5] && (piVar8[2] < 2)) || ((piVar8[5] == 1 && (piVar8[2] == 0)))) &&
          (local_20 = FUN_10050a40(0x10076210,(char *)(piVar8 + 6),DAT_1008aa70,0x49),
          -1 < (int)local_20)))) {
        bVar15 = true;
        local_8 = sVar3;
      }
      else {
LAB_100353be:
        piVar8 = FUN_100544b0(piVar5,param_6);
        local_8 = sVar3;
      }
    }
    else {
      if ((*piVar6 == 1) && (piVar6[2] < 2)) {
        local_10 = FUN_10050ac0(0x100766c0,(char *)(piVar6 + 6),DAT_1008aa74,0x49);
        if ((local_10 < 0x80000000) ||
           (((1 < piVar6[5] || ((piVar6[5] == 1 && (piVar6[2] == 0)))) &&
            (local_20 = FUN_10050a40(0x10076210,(char *)(piVar6 + 6),DAT_1008aa70,0x49),
            -1 < (int)local_20)))) {
          piVar6 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
          if (piVar6 != (int *)0x0) {
            if ((1 < piVar6[2]) || (*piVar6 != 2)) goto LAB_10034ae8;
            local_14 = (byte *)(piVar6 + 6);
            bVar15 = true;
            piVar6 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
            if (piVar6 != (int *)0x0) {
              if (((piVar6[2] < 2) && (((*piVar6 == 3 || (*piVar6 == 7)) && (piVar6[5] == 1)))) &&
                 ((char)piVar6[6] == '-')) {
                bVar13 = true;
              }
              else {
                bVar13 = false;
                FUN_100544b0(piVar5,param_6);
              }
              piVar6 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
              if (piVar6 != (int *)0x0) {
                if (((piVar6[2] < 2) && (*piVar6 == 1)) &&
                   (local_18 = FUN_10050a40(0x10075d00,(char *)(piVar6 + 6),DAT_1008aa68,0x49),
                   -1 < (int)local_18)) {
                  piVar8 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
                  if (piVar8 != (int *)0x0) goto LAB_10034a8f;
LAB_10034a50:
                  *param_2 = iVar2;
                  *piVar5 = (int)param_3;
                  *param_6 = piVar5[(int)param_3 + 1];
                  param_1[1] = iVar12;
                  if (0 < iVar12) {
                    *param_1 = param_1[iVar12 * 0x22 + -0x1c];
                    return CONCAT22(uVar11,4);
                  }
                  goto LAB_100354ba;
                }
                piVar8 = FUN_100544b0(piVar5,param_6);
                local_8 = sVar3;
                if (bVar13) {
                  piVar8 = FUN_100544b0(piVar5,param_6);
                }
                goto LAB_100353ec;
              }
            }
          }
LAB_10035351:
          *param_2 = iVar2;
          *piVar5 = (int)param_3;
          *param_6 = piVar5[(int)param_3 + 1];
LAB_10035368:
          param_1[1] = iVar12;
          if (0 < iVar12) {
            *param_1 = param_1[iVar12 * 0x22 + -0x1c];
            return CONCAT22(uVar11,4);
          }
          goto LAB_100354ba;
        }
      }
LAB_10034ae8:
      piVar8 = FUN_100544b0(piVar5,param_6);
      local_8 = sVar3;
    }
LAB_100353ec:
    bVar13 = true;
    if (param_8 != 0) {
      bVar13 = false;
      piVar8 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
      if (piVar8 == (int *)0x0) {
        *param_2 = iVar2;
        *piVar5 = (int)param_3;
        *param_6 = piVar5[(int)param_3 + 1];
        param_1[1] = iVar12;
        if (0 < iVar12) {
          *param_1 = param_1[iVar12 * 0x22 + -0x1c];
          return CONCAT22(uVar11,4);
        }
LAB_100354ba:
        *param_1 = *param_1;
        return CONCAT22(uVar11,4);
      }
      if (((*piVar8 == 8) && ((char)piVar8[6] == '\0')) &&
         (piVar8 = FUN_10054420(piVar5,param_5,param_6,param_7,'D'), piVar8 == (int *)0x0)) {
        *param_2 = iVar2;
        *piVar5 = (int)param_3;
        *param_6 = piVar5[(int)param_3 + 1];
        param_1[1] = iVar12;
        if (0 < iVar12) {
          *param_1 = param_1[iVar12 * 0x22 + -0x1c];
          return CONCAT22(uVar11,4);
        }
        goto LAB_100354ba;
      }
      if ((*piVar8 == 5) || (*piVar8 == 9)) {
        bVar13 = true;
      }
    }
  }
  else if (*piVar8 == 1) {
    local_20 = FUN_10050a40(0x10076210,(char *)(piVar8 + 6),DAT_1008aa70,0x49);
    if ((int)local_20 < 0) {
      if ((*piVar8 != 1) ||
         (local_10 = FUN_10050ac0(0x100766c0,(char *)(piVar8 + 6),DAT_1008aa74,0x49),
         (int)local_10 < 0)) goto LAB_10034ebf;
      piVar6 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
      if (piVar6 == (int *)0x0) {
        *param_2 = iVar2;
        *piVar5 = (int)param_3;
        *param_6 = piVar5[(int)param_3 + 1];
        goto LAB_10035368;
      }
      if ((1 < piVar6[2]) || (*piVar6 != 2)) goto LAB_10034eae;
      local_14 = (byte *)(piVar6 + 6);
      piVar6 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
      if (piVar6 == (int *)0x0) {
        *param_2 = iVar2;
        *piVar5 = (int)param_3;
        *param_6 = piVar5[(int)param_3 + 1];
        goto LAB_10035368;
      }
      if ((piVar6[2] < 2) &&
         ((((*piVar6 == 3 || (*piVar6 == 7)) && (piVar6[5] == 1)) && ((char)piVar6[6] == '-')))) {
        bVar13 = true;
      }
      else {
        bVar13 = false;
        FUN_100544b0(piVar5,param_6);
      }
      piVar6 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
      if (piVar6 == (int *)0x0) {
        *param_2 = iVar2;
        *piVar5 = (int)param_3;
        *param_6 = piVar5[(int)param_3 + 1];
        goto LAB_10035368;
      }
      if (((piVar6[2] < 2) && (*piVar6 == 1)) &&
         (local_18 = FUN_10050a40(0x10075d00,(char *)(piVar6 + 6),DAT_1008aa68,0x49),
         -1 < (int)local_18)) {
        piVar8 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
        if (piVar8 == (int *)0x0) goto LAB_10034a50;
LAB_10034a8f:
        if ((piVar8[2] != 0) || (2 < piVar8[5])) goto LAB_100353be;
        bVar13 = (char)piVar8[6] == '.';
      }
      else {
LAB_10034e97:
        piVar8 = FUN_100544b0(piVar5,param_6);
        bVar13 = !bVar13;
      }
      local_8 = sVar3;
      if (!bVar13) {
        piVar8 = FUN_100544b0(piVar5,param_6);
      }
      goto LAB_100353ec;
    }
    piVar6 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
    if (piVar6 == (int *)0x0) {
      *param_2 = iVar2;
      *piVar5 = (int)param_3;
      *param_6 = piVar5[(int)param_3 + 1];
      goto LAB_10035368;
    }
    if (*piVar6 == 3) {
      iVar7 = piVar8[5];
      bVar16 = SBORROW4(iVar7,1);
      bVar14 = iVar7 == 1;
      if (bVar14) {
        if (piVar6[2] != 0) {
          bVar16 = false;
          bVar14 = true;
          goto LAB_10034b8a;
        }
      }
      else {
LAB_10034b8a:
        if ((bVar14 || bVar16 != iVar7 + -1 < 0) || (1 < piVar6[2])) goto LAB_10034eae;
      }
      local_1c = FUN_10050ac0(0x10075f10,(char *)(piVar6 + 6),DAT_1008aa6c,0x53);
      if (-1 < (int)local_1c) {
        piVar6 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
        if (piVar6 != (int *)0x0) {
          if ((1 < piVar6[2]) || (*piVar6 != 2)) goto LAB_10034ae8;
          local_14 = (byte *)(piVar6 + 6);
          bVar15 = true;
          piVar6 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
          if (piVar6 != (int *)0x0) {
            if (((piVar6[2] < 2) && (((*piVar6 == 3 || (*piVar6 == 7)) && (piVar6[5] == 1)))) &&
               ((char)piVar6[6] == '-')) {
              bVar13 = true;
            }
            else {
              bVar13 = false;
              FUN_100544b0(piVar5,param_6);
            }
            piVar6 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
            if (piVar6 != (int *)0x0) {
              if (((1 < piVar6[2]) || (*piVar6 != 1)) ||
                 (local_18 = FUN_10050a40(0x10075d00,(char *)(piVar6 + 6),DAT_1008aa68,0x49),
                 (int)local_18 < 0)) goto LAB_10034e97;
              piVar8 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
              if (piVar8 == (int *)0x0) {
                *param_2 = iVar2;
                *piVar5 = (int)param_3;
                *param_6 = piVar5[(int)param_3 + 1];
                param_1[1] = iVar12;
                if (0 < iVar12) {
                  *param_1 = param_1[iVar12 * 0x22 + -0x1c];
                  return CONCAT22(uVar11,4);
                }
                goto LAB_100354ba;
              }
              goto LAB_10034a8f;
            }
          }
        }
        goto LAB_10035351;
      }
    }
LAB_10034eae:
    piVar8 = piVar6;
    if (param_8 != 0) goto LAB_100353d8;
  }
  else {
LAB_10034ebf:
    if (*piVar8 == 2) {
      local_14 = (byte *)(piVar8 + 6);
      piVar6 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
      if (piVar6 == (int *)0x0) {
        *param_2 = iVar2;
        *piVar5 = (int)param_3;
        *param_6 = piVar5[(int)param_3 + 1];
        param_1[1] = iVar12;
        if (0 < iVar12) {
          *param_1 = param_1[iVar12 * 0x22 + -0x1c];
          return CONCAT22(uVar11,4);
        }
        goto LAB_100354ba;
      }
      if ((piVar6[2] < 2) &&
         ((((*piVar6 == 3 || (*piVar6 == 7)) && (piVar6[5] == 1)) && ((char)piVar6[6] == '-')))) {
        bVar14 = true;
      }
      else {
        bVar14 = false;
        FUN_100544b0(piVar5,param_6);
      }
      piVar6 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
      if (piVar6 == (int *)0x0) {
        *param_2 = iVar2;
        *piVar5 = (int)param_3;
        *param_6 = piVar5[(int)param_3 + 1];
        goto LAB_10035368;
      }
      if (((piVar6[2] < 2) && (*piVar6 == 1)) &&
         (local_18 = FUN_10050a40(0x10075d00,(char *)(piVar6 + 6),DAT_1008aa68,0x49),
         -1 < (int)local_18)) {
        piVar6 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
        if (piVar6 == (int *)0x0) {
          *param_2 = iVar2;
          *piVar5 = (int)param_3;
          *param_6 = piVar5[(int)param_3 + 1];
          goto LAB_10035368;
        }
        if (((piVar6[2] != 0) || (2 < piVar6[5])) || ((char)piVar6[6] != '.')) goto LAB_10035021;
      }
      else {
        FUN_100544b0(piVar5,param_6);
        if (bVar14) {
LAB_10035021:
          FUN_100544b0(piVar5,param_6);
        }
      }
      piVar6 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
      if (piVar6 == (int *)0x0) goto LAB_10035351;
      if ((((piVar6[2] < 2) && (*piVar6 == 3)) && (*(char *)((int)piVar6 + 0x19) == '\0')) &&
         (((char)piVar6[6] != '\0' &&
          (pcVar9 = _strchr(&DAT_10074fbc,(int)(char)piVar6[6]), pcVar9 != (char *)0x0)))) {
        bVar14 = true;
      }
      else {
        bVar14 = false;
        FUN_100544b0(piVar5,param_6);
      }
      piVar8 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
      if (piVar8 == (int *)0x0) {
        *param_2 = iVar2;
        *piVar5 = (int)param_3;
        *param_6 = piVar5[(int)param_3 + 1];
        param_1[1] = iVar12;
        if (0 < iVar12) {
          *param_1 = param_1[iVar12 * 0x22 + -0x1c];
          return CONCAT22(uVar11,4);
        }
        goto LAB_100354ba;
      }
      if (piVar8[2] < 2) {
        if ((*piVar8 == 1) &&
           (local_10 = FUN_10050ac0(0x100766c0,(char *)(piVar8 + 6),DAT_1008aa74,0x49),
           -1 < (int)local_10)) {
          local_8 = sVar3;
          if (!bVar14) goto LAB_100353ec;
          piVar8 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
          if (piVar8 == (int *)0x0) {
            *param_2 = iVar2;
            *piVar5 = (int)param_3;
            *param_6 = piVar5[(int)param_3 + 1];
            goto LAB_10035368;
          }
          if (((1 < piVar8[2]) || (*piVar8 != 3)) || (*(char *)((int)piVar8 + 0x19) != '\0'))
          goto LAB_100353be;
          cVar4 = (char)piVar8[6];
        }
        else {
          if ((1 < piVar8[2]) ||
             (((*piVar8 != 3 ||
               (local_1c = FUN_10050ac0(0x10075f10,(char *)(piVar8 + 6),DAT_1008aa6c,0x53),
               0x7fffffff < local_1c)) &&
              ((*piVar8 != 1 ||
               (local_20 = FUN_10050a40(0x10076210,(char *)(piVar8 + 6),DAT_1008aa70,0x49),
               (int)local_20 < 0)))))) goto LAB_100353cc;
          piVar6 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
          if (piVar6 == (int *)0x0) {
            *param_2 = iVar2;
            *piVar5 = (int)param_3;
            *param_6 = piVar5[(int)param_3 + 1];
            param_1[1] = iVar12;
            if (0 < iVar12) {
              *param_1 = param_1[iVar12 * 0x22 + -0x1c];
              return CONCAT22(uVar11,4);
            }
            goto LAB_100354ba;
          }
          if ((int)local_1c < 0) {
            if (-1 < (int)local_20) {
              iVar7 = piVar8[5];
              bVar17 = SBORROW4(iVar7,1);
              bVar16 = iVar7 == 1;
              if (bVar16) {
                if (piVar6[2] != 0) {
                  bVar17 = false;
                  bVar16 = true;
                  goto LAB_100352db;
                }
              }
              else {
LAB_100352db:
                if ((bVar16 || bVar17 != iVar7 + -1 < 0) || (1 < piVar6[2])) goto LAB_1003530b;
              }
              if ((*piVar6 == 3) &&
                 (local_1c = FUN_10050ac0(0x10075f10,(char *)(piVar6 + 6),DAT_1008aa6c,0x53),
                 local_1c < 0x80000000)) goto LAB_10035329;
            }
LAB_1003530b:
            piVar8 = FUN_100544b0(piVar5,param_6);
            goto LAB_100354de;
          }
          if (*piVar6 == 1) {
            iVar7 = piVar6[5];
            bVar13 = SBORROW4(iVar7,1);
            bVar15 = iVar7 == 1;
            if (bVar15) {
              if (piVar6[2] != 0) {
                bVar13 = false;
                bVar15 = true;
                goto LAB_1003527a;
              }
            }
            else {
LAB_1003527a:
              if ((bVar15 || bVar13 != iVar7 + -1 < 0) || (1 < piVar6[2])) goto LAB_100352a5;
            }
            local_20 = FUN_10050a40(0x10076210,(char *)(piVar6 + 6),DAT_1008aa70,0x49);
            if (0x7fffffff < local_20) goto LAB_100352a5;
          }
          else {
LAB_100352a5:
            piVar6 = FUN_100544b0(piVar5,param_6);
          }
LAB_10035329:
          bVar15 = true;
          piVar8 = piVar6;
          local_8 = sVar3;
          if (!bVar14) goto LAB_100353ec;
          piVar8 = FUN_10054420(piVar5,param_5,param_6,param_7,'D');
          if (piVar8 == (int *)0x0) goto LAB_10035351;
          if (((1 < piVar8[2]) || (*piVar8 != 3)) || (*(char *)((int)piVar8 + 0x19) != '\0'))
          goto LAB_100353be;
          cVar4 = (char)piVar8[6];
        }
        if ((cVar4 == '\0') ||
           (pcVar9 = _strchr(&DAT_10074fb4,(int)cVar4), local_8 = sVar3, pcVar9 == (char *)0x0))
        goto LAB_100353be;
        goto LAB_100353ec;
      }
LAB_100353cc:
      if (param_8 == 0) goto LAB_100354de;
LAB_100353d8:
      piVar8 = FUN_100544b0(piVar5,param_6);
      local_8 = sVar3;
      goto LAB_100353ec;
    }
  }
LAB_100354de:
  if (((local_1c != 1) || (bVar15)) && (bVar13)) {
    uVar10 = FUN_10054280(param_2,piVar5,7);
    if ((short)uVar10 == 1) {
      FUN_100355c0(param_4,local_4,piVar8[4],local_14,local_1c,local_18,local_20,local_10,param_8,
                   local_8);
      uVar10 = FUN_100483c0(param_1,(int *)param_4);
      if ((short)uVar10 == 1) {
        return uVar10;
      }
    }
    *param_2 = iVar2;
    *piVar5 = (int)param_3;
    *param_6 = piVar5[(int)param_3 + 1];
    param_1[1] = iVar12;
    if (iVar12 < 1) {
      iVar12 = *param_1;
    }
    else {
      iVar12 = param_1[iVar12 * 0x22 + -0x1c];
    }
    *param_1 = iVar12;
  }
  piVar6 = param_3;
  if ((int *)*piVar5 != param_3) {
    *piVar5 = (int)param_3;
    piVar6 = (int *)piVar5[(int)param_3 + 1];
    *param_6 = (int)piVar6;
  }
LAB_100355ab:
  return CONCAT22((short)((uint)piVar6 >> 0x10),5);
}



/* ===== FUN_100552e0 ===== */
/* Entry: 100552e0 */

uint __cdecl
FUN_100552e0(int *param_1,int *param_2,int *param_3,undefined4 param_4,int *param_5,int *param_6,
            char *param_7)

{
  char cVar1;
  int iVar2;
  int *piVar3;
  int *piVar4;
  bool bVar5;
  int iVar6;
  int *piVar7;
  uint uVar8;
  char *pcVar9;
  int *piVar10;
  int iVar11;
  undefined3 extraout_var;
  undefined3 extraout_var_00;
  int iVar12;
  uint uVar13;
  int iVar14;
  int *piVar15;
  int *piVar16;
  char *pcVar17;
  int local_28;
  char local_20 [32];
  
  piVar4 = param_3;
  piVar3 = param_1;
  iVar11 = *param_3;
  iVar14 = param_1[1];
  iVar2 = *param_2;
  piVar10 = param_3 + iVar11 * 0xe + 0x57;
  if ((*piVar10 == 1) &&
     (iVar6 = FUN_10050830(0x10077d40,(char *)(piVar10 + 6),DAT_1008aa78,0x49), -1 < iVar6)) {
    iVar6 = piVar10[3];
    piVar7 = FUN_10054420(param_3,(int)param_5,param_6,param_7,'R');
    if (piVar7 == (int *)0x0) {
LAB_1005536d:
      *param_2 = iVar2;
      *piVar4 = iVar11;
      *param_6 = piVar4[iVar11 + 1];
      param_1[1] = iVar14;
      if (0 < iVar14) {
        iVar11 = param_1[iVar14 * 0x22 + -0x1c];
        *param_1 = iVar11;
        return CONCAT22((short)((uint)iVar11 >> 0x10),4);
      }
LAB_10055bfd:
      iVar11 = *piVar3;
      *piVar3 = iVar11;
      return CONCAT22((short)((uint)iVar11 >> 0x10),4);
    }
    if ((piVar7[2] != 0) || (piVar7[5] < 5)) {
      FUN_100544b0(param_3,param_6);
      goto LAB_100555c2;
    }
    iVar12 = param_1[1];
    FUN_10056980(local_20,(char *)(piVar10 + 6));
    uVar8 = FUN_10047b90(param_1,local_20,iVar6,piVar7[4],0x52,0x44);
    if (((short)uVar8 == 1) &&
       (uVar8 = FUN_10047ff0(param_1,param_5,(int)(piVar7 + 6),piVar7[5],iVar6,piVar7[4],'R'),
       (short)uVar8 == 1)) {
      param_3 = (int *)0xffffffff;
LAB_10055432:
      cVar1 = *(char *)*param_6;
      if (((((&DAT_10071088)[cVar1] & 0xd0) == 0) &&
          ((cVar1 == '\0' ||
           (pcVar9 = _strchr(s____________________1008ba28,(int)cVar1), pcVar9 == (char *)0x0)))) ||
         ((((*(char *)*param_6 != '\0' &&
            (pcVar9 = _strchr(&DAT_10075500,(int)*(char *)*param_6), pcVar9 != (char *)0x0)) ||
           (*(char *)*param_6 == ',')) &&
          ((((cVar1 = *(char *)(*param_6 + 1), cVar1 == ' ' || (cVar1 == '\t')) || (cVar1 == '\n'))
           || (((cVar1 == '\r' || (cVar1 == '\0')) ||
               (pcVar9 = _strchr(&DAT_10075500,(int)cVar1), pcVar9 != (char *)0x0)))))))) {
        if ((-1 < (int)param_3) && (iVar12 < param_1[1])) {
          piVar10 = param_1 + iVar12 * 0x22 + 6;
          do {
            *piVar10 = (int)param_3;
            iVar12 = iVar12 + 1;
            piVar10 = piVar10 + 0x22;
          } while (iVar12 < param_1[1]);
        }
        uVar8 = FUN_10054280(param_2,piVar4,0xc);
        if ((short)uVar8 == 1) {
          return uVar8;
        }
        goto LAB_1005557d;
      }
      piVar10 = FUN_10054420(piVar4,(int)param_5,param_6,param_7,'R');
      if (piVar10 == (int *)0x0) goto LAB_1005536d;
      uVar8 = FUN_10047ff0(param_1,param_5,(int)(piVar10 + 6),piVar10[5],iVar6,piVar10[4],'R');
      if ((short)uVar8 == 1) goto code_r0x10055511;
      *piVar4 = iVar11;
      *param_6 = piVar4[iVar11 + 1];
    }
    else {
LAB_1005557d:
      *piVar4 = iVar11;
      *param_6 = piVar4[iVar11 + 1];
    }
    *param_2 = iVar2;
    param_1[1] = iVar14;
    if (0 < iVar14) {
      *param_1 = param_1[iVar14 * 0x22 + -0x1c];
      return uVar8;
    }
  }
  else {
    if ((piVar10[2] == 0) &&
       ((((*piVar10 == 3 || (*piVar10 == 7)) && (piVar10[5] == 1)) && ((char)piVar10[6] == '@')))) {
      piVar10 = (int *)piVar10[3];
      piVar7 = FUN_10054420(param_3,(int)param_5,param_6,param_7,'R');
      if (piVar7 == (int *)0x0) {
        *param_2 = iVar2;
        *param_3 = iVar11;
        *param_6 = param_3[iVar11 + 1];
        param_1[1] = iVar14;
        if (0 < iVar14) {
          iVar11 = param_1[iVar14 * 0x22 + -0x1c];
          *param_1 = iVar11;
          return CONCAT22((short)((uint)iVar11 >> 0x10),4);
        }
        goto LAB_10055bfd;
      }
      if (piVar7[2] != 0) {
LAB_100555c2:
        if (*param_3 != iVar11) {
          *param_3 = iVar11;
          iVar11 = param_3[iVar11 + 1];
          *param_6 = iVar11;
        }
        return CONCAT22((short)((uint)iVar11 >> 0x10),5);
      }
      bVar5 = FUN_10055cb0((byte *)(piVar7 + 6),piVar7[5]);
      if (((short)CONCAT31(extraout_var,bVar5) == 0) ||
         (iVar6 = FUN_10014ed0((byte *)(piVar7 + 6),&DAT_1008ba3c,4), iVar6 == 0))
      goto LAB_100555c2;
      if ((param_3[0x67] == 0) &&
         (param_3 = (int *)FUN_100577a0(param_2,(int)piVar10), -1 < (int)param_3)) {
        iVar6 = FUN_10057850((int)param_1,param_2[(int)param_3 * 0xe + 4]);
        if (-1 < iVar6) {
          piVar10 = (int *)param_2[(int)param_3 * 0xe + 4];
          param_1[1] = iVar6;
          if ((int)param_3 < *param_2) {
            do {
              uVar8 = 0xffffffff;
              param_2[(int)param_3 * 0xe + 2] = 0xd;
              piVar15 = param_2 + (int)param_3 * 0xe + 7;
              do {
                piVar16 = piVar15;
                if (uVar8 == 0) break;
                uVar8 = uVar8 - 1;
                piVar16 = (int *)((int)piVar15 + 1);
                iVar6 = *piVar15;
                piVar15 = piVar16;
              } while ((char)iVar6 != '\0');
              uVar8 = ~uVar8;
              pcVar9 = (char *)((int)piVar16 - uVar8);
              pcVar17 = local_20;
              for (uVar13 = uVar8 >> 2; uVar13 != 0; uVar13 = uVar13 - 1) {
                *(undefined4 *)pcVar17 = *(undefined4 *)pcVar9;
                pcVar9 = pcVar9 + 4;
                pcVar17 = pcVar17 + 4;
              }
              for (uVar8 = uVar8 & 3; uVar8 != 0; uVar8 = uVar8 - 1) {
                *pcVar17 = *pcVar9;
                pcVar9 = pcVar9 + 1;
                pcVar17 = pcVar17 + 1;
              }
              iVar6 = param_2[(int)param_3 * 0xe + 6];
              param_7 = (char *)((int)param_3 + 2);
              if ((int)param_7 < *param_2) {
                piVar15 = param_2 + (int)param_3 * 0xe + 0xf;
                do {
                  if (((piVar15[-0xe] != 1) ||
                      (((*piVar15 != 3 && (*piVar15 != 7)) || (piVar15[5] != 1)))) ||
                     (((char)piVar15[6] != '\'' || (piVar15[0xe] != 1)))) break;
                  FUN_100588f6(local_20,(byte *)s__s__s_10074d6c);
                  param_3 = (int *)((int)param_3 + 2);
                  iVar6 = iVar6 + 1 + piVar15[0x13];
                  param_7 = param_7 + 2;
                  piVar15 = piVar15 + 0x1c;
                } while ((int)param_7 < *param_2);
              }
              uVar8 = FUN_10047ff0(param_1,param_5,(int)local_20,iVar6,(int)piVar10,piVar7[4],'I');
              if ((short)uVar8 != 1) {
                *param_2 = iVar2;
                param_1[1] = iVar14;
                if (iVar14 < 1) {
                  iVar14 = *param_1;
                }
                else {
                  iVar14 = param_1[iVar14 * 0x22 + -0x1c];
                }
                *param_1 = iVar14;
                *piVar4 = iVar11;
                *param_6 = piVar4[iVar11 + 1];
                return uVar8;
              }
              param_3 = (int *)((int)param_3 + 1);
            } while ((int)param_3 < *param_2);
          }
        }
      }
      param_1 = piVar10;
      param_3 = (int *)0x0;
      if (0 < *piVar4) {
        piVar10 = piVar4 + 0x6b;
        do {
          uVar8 = FUN_10047ff0(piVar3,param_5,(int)piVar10,piVar10[-1],(int)param_1,piVar7[4],'R');
          if ((short)uVar8 != 1) {
            *param_2 = iVar2;
            *piVar4 = iVar11;
            *param_6 = piVar4[iVar11 + 1];
            goto LAB_10055c76;
          }
          param_3 = (int *)((int)param_3 + 1);
          piVar10 = piVar10 + 0xe;
        } while ((int)param_3 < *piVar4);
      }
      uVar8 = FUN_10054280(param_2,piVar4,0xd);
      if ((short)uVar8 == 1) {
        return uVar8;
      }
LAB_10055c60:
      *piVar4 = iVar11;
      *param_6 = piVar4[iVar11 + 1];
      *param_2 = iVar2;
    }
    else {
      bVar5 = FUN_10055cb0((byte *)(piVar10 + 6),piVar10[5]);
      if ((short)CONCAT31(extraout_var_00,bVar5) == 0) goto LAB_100555c2;
      FUN_100544b0(param_3,param_6);
      piVar10 = FUN_10054420(param_3,(int)param_5,param_6,param_7,'R');
      if (piVar10 == (int *)0x0) {
        *param_2 = iVar2;
        *param_3 = iVar11;
        *param_6 = param_3[iVar11 + 1];
        param_1[1] = iVar14;
        if (0 < iVar14) {
          iVar11 = param_1[iVar14 * 0x22 + -0x1c];
          *param_1 = iVar11;
          return CONCAT22((short)((uint)iVar11 >> 0x10),4);
        }
        goto LAB_10055bfd;
      }
      piVar7 = (int *)piVar10[3];
      local_28 = -1;
      if ((piVar10[2] == 0) &&
         (param_3 = (int *)FUN_100577a0(param_2,(int)piVar7), -1 < (int)param_3)) {
        local_28 = FUN_10057850((int)param_1,param_2[(int)param_3 * 0xe + 4]);
        if (-1 < local_28) {
          piVar7 = (int *)param_2[(int)param_3 * 0xe + 4];
          param_1[1] = local_28;
          if ((int)param_3 < *param_2) {
            piVar15 = param_2 + (int)param_3 * 0xe + 6;
            do {
              piVar15[-4] = 0xc;
              uVar8 = FUN_10047ff0(param_1,param_5,(int)(piVar15 + 1),*piVar15,(int)piVar7,
                                   piVar10[4],'R');
              if ((short)uVar8 != 1) goto LAB_10055ba3;
              param_3 = (int *)((int)param_3 + 1);
              piVar15 = piVar15 + 0xe;
            } while ((int)param_3 < *param_2);
          }
        }
      }
      param_1 = piVar7;
      param_3 = (int *)0x0;
      if (0 < *piVar4) {
        piVar7 = piVar4 + 0x6b;
        do {
          uVar8 = FUN_10047ff0(piVar3,param_5,(int)piVar7,piVar7[-1],(int)param_1,piVar10[4],'R');
          if ((short)uVar8 != 1) goto LAB_10055ba3;
          param_3 = (int *)((int)param_3 + 1);
          piVar7 = piVar7 + 0xe;
        } while ((int)param_3 < *piVar4);
      }
      param_3 = (int *)0xffffffff;
LAB_10055ab4:
      cVar1 = *(char *)*param_6;
      if (((((&DAT_10071088)[cVar1] & 0xd0) == 0) &&
          ((cVar1 == '\0' ||
           (pcVar9 = _strchr(s____________________1008ba28,(int)cVar1), pcVar9 == (char *)0x0)))) ||
         ((((*(char *)*param_6 != '\0' &&
            (pcVar9 = _strchr(&DAT_10075500,(int)*(char *)*param_6), pcVar9 != (char *)0x0)) ||
           (*(char *)*param_6 == ',')) &&
          ((((cVar1 = *(char *)(*param_6 + 1), cVar1 == ' ' || (cVar1 == '\t')) || (cVar1 == '\n'))
           || (((cVar1 == '\r' || (cVar1 == '\0')) ||
               (pcVar9 = _strchr(&DAT_10075500,(int)cVar1), pcVar9 != (char *)0x0)))))))) {
        if (-1 < (int)param_3) {
          if (local_28 == -1) {
            local_28 = 0;
          }
          if (local_28 < piVar3[1]) {
            piVar10 = piVar3 + local_28 * 0x22 + 6;
            do {
              *piVar10 = (int)param_3;
              local_28 = local_28 + 1;
              piVar10 = piVar10 + 0x22;
            } while (local_28 < piVar3[1]);
          }
        }
        uVar8 = FUN_10054280(param_2,piVar4,0xc);
        if ((short)uVar8 == 1) {
          return uVar8;
        }
        goto LAB_10055c60;
      }
      piVar10 = FUN_10054420(piVar4,(int)param_5,param_6,param_7,'R');
      if (piVar10 == (int *)0x0) {
        *param_2 = iVar2;
        *piVar4 = iVar11;
        *param_6 = piVar4[iVar11 + 1];
        piVar3[1] = iVar14;
        if (0 < iVar14) {
          iVar11 = piVar3[iVar14 * 0x22 + -0x1c];
          *piVar3 = iVar11;
          return CONCAT22((short)((uint)iVar11 >> 0x10),4);
        }
        goto LAB_10055bfd;
      }
      uVar8 = FUN_10047ff0(piVar3,param_5,(int)(piVar10 + 6),piVar10[5],(int)param_1,piVar10[4],'R')
      ;
      if ((short)uVar8 == 1) goto code_r0x10055b97;
LAB_10055ba3:
      *piVar4 = iVar11;
      *param_6 = piVar4[iVar11 + 1];
      *param_2 = iVar2;
    }
LAB_10055c76:
    piVar3[1] = iVar14;
    if (0 < iVar14) {
      *piVar3 = piVar3[iVar14 * 0x22 + -0x1c];
      return uVar8;
    }
  }
  *piVar3 = *piVar3;
  return uVar8;
code_r0x10055511:
  param_3 = (int *)piVar10[4];
  goto LAB_10055432;
code_r0x10055b97:
  param_3 = (int *)piVar10[4];
  goto LAB_10055ab4;
}



/* ===== FUN_100512a0 ===== */
/* Entry: 100512a0 */

undefined4 __cdecl
FUN_100512a0(int *param_1,int *param_2,int *param_3,uint *param_4,int param_5,int *param_6,
            char *param_7,short param_8)

{
  byte bVar1;
  short sVar2;
  byte *pbVar3;
  int iVar4;
  undefined4 uVar5;
  undefined3 extraout_var;
  int iVar6;
  int *extraout_ECX;
  int *piVar7;
  void *this;
  void *extraout_ECX_00;
  void *extraout_ECX_01;
  void *this_00;
  void *this_01;
  bool bVar8;
  undefined2 in_stack_00000022;
  int *local_134;
  int local_130;
  int local_12c;
  int local_128;
  byte *local_124;
  byte *local_120;
  int *local_11c;
  undefined4 local_118;
  int *local_114;
  int *local_110;
  int local_10c;
  undefined4 local_108;
  undefined4 local_104;
  char local_100 [4];
  char local_fc [252];
  
  local_104 = 1;
  local_108 = 0;
  if (param_8 == 0) {
    local_130 = *param_3;
    local_134 = param_3 + local_130 * 0xe + 0x57;
    iVar6 = local_130;
    if (*local_134 != 2) goto LAB_10052435;
    local_10c = local_134[3];
  }
  else {
    iVar6 = _param_8;
    if ((param_8 < 0x68) || (0x6b < param_8)) goto LAB_10052435;
    *param_3 = 0;
    piVar7 = param_3;
    for (iVar6 = 0x19; piVar7 = piVar7 + 1, iVar6 != 0; iVar6 = iVar6 + -1) {
      *piVar7 = -1;
    }
    piVar7 = param_3;
    for (iVar6 = 0x578; iVar6 != 0; iVar6 = iVar6 + -1) {
      *piVar7 = 0;
      piVar7 = piVar7 + 1;
    }
    local_130 = *param_3;
    local_134 = FUN_10054420(param_3,param_5,param_6,param_7,'D');
    iVar6 = 0;
    if (local_134 == (int *)0x0) goto LAB_10052435;
    local_10c = local_134[3];
  }
  local_12c = *param_2;
  iVar6 = param_1[1];
  local_118 = 0;
  local_11c = (int *)0x0;
  local_114 = (int *)0x0;
  local_110 = (int *)0x0;
  local_128 = 0;
  local_120 = (byte *)FUN_10057930(param_2,param_3,*param_3 + -1,0xffffffff);
  if (local_120 == (byte *)0x0) {
LAB_10051491:
    uVar5 = FUN_10052450((byte *)(local_134 + 6),'H',param_8);
    if (((short)uVar5 != 0) && (local_120 != (byte *)0x0)) {
      local_124 = (byte *)&DAT_1006fe60;
      pbVar3 = (byte *)((int)local_120 + 0x18);
      do {
        bVar1 = *pbVar3;
        bVar8 = bVar1 < *local_124;
        if (bVar1 != *local_124) {
LAB_100514fb:
          iVar4 = (1 - (uint)bVar8) - (uint)(bVar8 != 0);
          goto LAB_10051500;
        }
        if (bVar1 == 0) break;
        bVar1 = pbVar3[1];
        bVar8 = bVar1 < local_124[1];
        if (bVar1 != local_124[1]) goto LAB_100514fb;
        pbVar3 = pbVar3 + 2;
        local_124 = local_124 + 2;
      } while (bVar1 != 0);
      iVar4 = 0;
LAB_10051500:
      if (iVar4 != 0) goto LAB_10051533;
    }
    uVar5 = FUN_10052450((byte *)(local_134 + 6),'H',param_8);
    if (((short)uVar5 != 0) && (local_120 == (byte *)0x0)) goto LAB_10051533;
    uVar5 = FUN_10052450((byte *)(local_134 + 6),'M',param_8);
    if ((short)uVar5 == 0) goto LAB_10051b43;
    local_114 = local_134;
    local_134 = FUN_10054420(param_3,param_5,param_6,param_7,'D');
    if (local_134 == (int *)0x0) {
      *param_2 = local_12c;
      *param_3 = local_130;
      *param_6 = param_3[local_130 + 1];
      param_1[1] = iVar6;
      if (0 < iVar6) {
        iVar6 = param_1[iVar6 * 0x22 + -0x1c];
        *param_1 = iVar6;
        return CONCAT22((short)((uint)iVar6 >> 0x10),4);
      }
      goto LAB_1005207f;
    }
    if ((2 < local_134[2]) ||
       ((((*local_134 != 3 && (*local_134 != 7)) || (local_134[5] != 1)) ||
        ((char)local_134[6] != ':')))) goto LAB_10051b43;
    local_134 = FUN_10054420(param_3,param_5,param_6,param_7,'D');
    if (local_134 == (int *)0x0) {
      *param_2 = local_12c;
      *param_3 = local_130;
      *param_6 = param_3[local_130 + 1];
      param_1[1] = iVar6;
      if (0 < iVar6) {
        iVar6 = param_1[iVar6 * 0x22 + -0x1c];
        *param_1 = iVar6;
        return CONCAT22((short)((uint)iVar6 >> 0x10),4);
      }
      goto LAB_1005207f;
    }
    if ((2 < local_134[2]) ||
       (uVar5 = FUN_10052450((byte *)(local_134 + 6),'S',param_8), (short)uVar5 == 0))
    goto LAB_10051b43;
    local_11c = local_134;
    local_134 = FUN_10054420(param_3,param_5,param_6,param_7,'D');
    if (local_134 == (int *)0x0) {
      *param_2 = local_12c;
      *param_3 = local_130;
      *param_6 = param_3[local_130 + 1];
      param_1[1] = iVar6;
      if (0 < iVar6) {
        iVar6 = param_1[iVar6 * 0x22 + -0x1c];
        *param_1 = iVar6;
        return CONCAT22((short)((uint)iVar6 >> 0x10),4);
      }
      goto LAB_1005207f;
    }
    if ((((local_134[2] < 2) && ((*local_134 == 3 || (*local_134 == 7)))) && (local_134[5] == 1)) &&
       ((char)local_134[6] == ':')) goto LAB_10051b43;
LAB_10051b35:
    local_128 = 1;
  }
  else {
    local_124 = &DAT_1006b4c0;
    pbVar3 = (byte *)((int)local_120 + 0x18);
    do {
      bVar1 = *pbVar3;
      bVar8 = bVar1 < *local_124;
      if (bVar1 != *local_124) {
LAB_10051400:
        iVar4 = (1 - (uint)bVar8) - (uint)(bVar8 != 0);
        goto LAB_10051405;
      }
      if (bVar1 == 0) break;
      bVar1 = pbVar3[1];
      bVar8 = bVar1 < local_124[1];
      if (bVar1 != local_124[1]) goto LAB_10051400;
      pbVar3 = pbVar3 + 2;
      local_124 = local_124 + 2;
    } while (bVar1 != 0);
    iVar4 = 0;
LAB_10051405:
    if (((iVar4 != 0) || (*local_134 != 2)) ||
       (uVar5 = FUN_10052450((byte *)(local_134 + 6),'H',param_8), (short)uVar5 != 0))
    goto LAB_10051491;
    local_124 = (byte *)&DAT_1006fe60;
    _param_8 = 1;
    param_8 = 1;
    local_108 = 1;
    pbVar3 = (byte *)((int)local_120 + 0x18);
    do {
      bVar1 = *pbVar3;
      bVar8 = bVar1 < *local_124;
      if (bVar1 != *local_124) {
LAB_10051480:
        iVar4 = (1 - (uint)bVar8) - (uint)(bVar8 != 0);
        goto LAB_10051485;
      }
      if (bVar1 == 0) break;
      bVar1 = pbVar3[1];
      bVar8 = bVar1 < local_124[1];
      if (bVar1 != local_124[1]) goto LAB_10051480;
      pbVar3 = pbVar3 + 2;
      local_124 = local_124 + 2;
    } while (bVar1 != 0);
    iVar4 = 0;
LAB_10051485:
    if (iVar4 == 0) goto LAB_10051491;
LAB_10051533:
    local_110 = local_134;
    if ((short)local_108 == 0) {
      if (param_8 != 0x69) goto LAB_1005154e;
      local_128 = 1;
      local_134 = FUN_10054420(param_3,param_5,param_6,param_7,'D');
      if (local_134 == (int *)0x0) {
        *param_2 = local_12c;
        *param_3 = local_130;
        *param_6 = param_3[local_130 + 1];
        param_1[1] = iVar6;
        if (0 < iVar6) {
          iVar6 = param_1[iVar6 * 0x22 + -0x1c];
          *param_1 = iVar6;
          return CONCAT22((short)((uint)iVar6 >> 0x10),4);
        }
        goto LAB_1005207f;
      }
    }
    else {
      param_8 = 0;
LAB_1005154e:
      local_134 = FUN_10054420(param_3,param_5,param_6,param_7,'D');
      if (local_134 == (int *)0x0) {
        *param_2 = local_12c;
        *param_3 = local_130;
        *param_6 = param_3[local_130 + 1];
        param_1[1] = iVar6;
        if (0 < iVar6) {
          iVar6 = param_1[iVar6 * 0x22 + -0x1c];
          *param_1 = iVar6;
          return CONCAT22((short)((uint)iVar6 >> 0x10),4);
        }
        goto LAB_1005207f;
      }
      if (((2 < local_134[2]) || (((*local_134 != 3 && (*local_134 != 7)) || (local_134[5] != 1))))
         || ((char)local_134[6] != ':')) goto LAB_10051b35;
      local_134 = FUN_10054420(param_3,param_5,param_6,param_7,'D');
      if (local_134 == (int *)0x0) {
        *param_2 = local_12c;
        *param_3 = local_130;
        *param_6 = param_3[local_130 + 1];
        param_1[1] = iVar6;
        if (0 < iVar6) {
          iVar6 = param_1[iVar6 * 0x22 + -0x1c];
          *param_1 = iVar6;
          return CONCAT22((short)((uint)iVar6 >> 0x10),4);
        }
        goto LAB_1005207f;
      }
      if ((local_134[2] < 3) &&
         (uVar5 = FUN_10052450((byte *)(local_134 + 6),'M',param_8), (short)uVar5 != 0)) {
        local_114 = local_134;
        if (param_8 == 0x6a) {
          local_128 = 1;
          local_134 = FUN_10054420(param_3,param_5,param_6,param_7,'D');
          if (local_134 == (int *)0x0) {
            *param_2 = local_12c;
            *param_3 = local_130;
            *param_6 = param_3[local_130 + 1];
            param_1[1] = iVar6;
            if (0 < iVar6) {
              iVar6 = param_1[iVar6 * 0x22 + -0x1c];
              *param_1 = iVar6;
              return CONCAT22((short)((uint)iVar6 >> 0x10),4);
            }
            goto LAB_1005207f;
          }
        }
        else {
          local_134 = FUN_10054420(param_3,param_5,param_6,param_7,'D');
          if (local_134 == (int *)0x0) {
            *param_2 = local_12c;
            *param_3 = local_130;
            *param_6 = param_3[local_130 + 1];
            param_1[1] = iVar6;
            if (0 < iVar6) {
              iVar6 = param_1[iVar6 * 0x22 + -0x1c];
              *param_1 = iVar6;
              return CONCAT22((short)((uint)iVar6 >> 0x10),4);
            }
            goto LAB_1005207f;
          }
          if (((local_134[2] < 3) && ((*local_134 == 3 || (*local_134 == 7)))) &&
             ((local_134[5] == 1 && ((char)local_134[6] == ':')))) {
            local_134 = FUN_10054420(param_3,param_5,param_6,param_7,'D');
            if (local_134 == (int *)0x0) {
              *param_2 = local_12c;
              *param_3 = local_130;
              *param_6 = param_3[local_130 + 1];
              param_1[1] = iVar6;
              if (0 < iVar6) {
                iVar6 = param_1[iVar6 * 0x22 + -0x1c];
                *param_1 = iVar6;
                return CONCAT22((short)((uint)iVar6 >> 0x10),4);
              }
              goto LAB_1005207f;
            }
            if ((local_134[2] < 3) &&
               (uVar5 = FUN_10052450((byte *)(local_134 + 6),'S',param_8), (short)uVar5 != 0)) {
              local_11c = local_134;
              local_134 = FUN_10054420(param_3,param_5,param_6,param_7,'D');
              if (local_134 == (int *)0x0) {
                *param_2 = local_12c;
                *param_3 = local_130;
                *param_6 = param_3[local_130 + 1];
                param_1[1] = iVar6;
                if (0 < iVar6) {
                  iVar6 = param_1[iVar6 * 0x22 + -0x1c];
                  *param_1 = iVar6;
                  return CONCAT22((short)((uint)iVar6 >> 0x10),4);
                }
                goto LAB_1005207f;
              }
              if (((1 < local_134[2]) ||
                  (((*local_134 != 3 && (*local_134 != 7)) || (local_134[5] != 1)))) ||
                 ((char)local_134[6] != ':')) {
                local_128 = 1;
                local_118 = 1;
              }
            }
          }
          else if ((param_8 == 0) || (param_8 == 0x68)) {
            local_128 = 1;
          }
        }
      }
    }
  }
LAB_10051b43:
  local_120 = (byte *)0x0;
  if ((short)local_128 == 0) {
LAB_10051f37:
    if ((short)local_118 != 0) goto LAB_10051f43;
LAB_10052422:
    iVar6 = local_130;
    if (*param_3 != local_130) {
      *param_3 = local_130;
      iVar6 = param_3[local_130 + 1];
      *param_6 = iVar6;
    }
LAB_10052435:
    return CONCAT22((short)((uint)iVar6 >> 0x10),5);
  }
  if (local_134[2] < 3) {
    if ((*local_134 == 1) &&
       (uVar5 = FUN_100524f0(param_4,(int *)&local_134,param_3,param_5,param_6,param_7,param_8),
       (short)uVar5 != 0)) {
      local_118 = 1;
      local_134 = FUN_10054420(param_3,param_5,param_6,param_7,'D');
      if (local_134 == (int *)0x0) {
        *param_2 = local_12c;
        *param_3 = local_130;
        *param_6 = param_3[local_130 + 1];
        param_1[1] = iVar6;
        if (0 < iVar6) {
          iVar6 = param_1[iVar6 * 0x22 + -0x1c];
          *param_1 = iVar6;
          return CONCAT22((short)((uint)iVar6 >> 0x10),4);
        }
        goto LAB_1005207f;
      }
      if (local_134[2] < 2) {
        iVar4 = *local_134;
joined_r0x10051e4c:
        if ((iVar4 == 1) &&
           (uVar5 = FUN_10052690(param_4,(int *)&local_134,param_3,param_5,param_6,param_7,0x44,
                                 param_8), (short)uVar5 != 0)) goto LAB_10051f43;
      }
LAB_10051e88:
      local_134 = FUN_100544b0(param_3,param_6);
    }
    else {
      if ((2 < local_134[2]) ||
         ((*local_134 != 1 ||
          (uVar5 = FUN_10052690(param_4,(int *)&local_134,param_3,param_5,param_6,param_7,0x49,
                                param_8), (short)uVar5 == 0)))) goto LAB_10051c7d;
      local_118 = 1;
    }
  }
  else {
LAB_10051c7d:
    if (param_8 == 0) {
LAB_10051f02:
      local_134 = FUN_100544b0(param_3,param_6);
      if ((param_8 == 0) &&
         (((piVar7 = local_114, local_110 == (int *)0x0 &&
           (piVar7 = local_11c, local_114 == (int *)0x0)) ||
          ((piVar7 == (int *)0x0 ||
           (bVar8 = FUN_10052800(param_2,param_3,param_5,*param_6),
           (short)CONCAT31(extraout_var,bVar8) == 0)))))) goto LAB_10051f37;
      local_118 = 1;
    }
    else {
      if (*local_134 != 5) {
        if ((*local_134 != 8) || ((char)local_134[6] != '\0')) goto LAB_10051f02;
        local_134 = FUN_10054420(param_3,param_5,param_6,param_7,'D');
        if (local_134 == (int *)0x0) {
          *param_2 = local_12c;
          *param_3 = local_130;
          *param_6 = param_3[local_130 + 1];
          param_1[1] = iVar6;
          if (0 < iVar6) {
            iVar6 = param_1[iVar6 * 0x22 + -0x1c];
            *param_1 = iVar6;
            return CONCAT22((short)((uint)iVar6 >> 0x10),4);
          }
          goto LAB_1005207f;
        }
        if (*local_134 != 5) goto LAB_10051f37;
      }
      local_120 = (byte *)0x1;
      local_134 = FUN_10054420(param_3,param_5,param_6,param_7,'D');
      if (local_134 == (int *)0x0) {
        *param_2 = local_12c;
        *param_3 = local_130;
        *param_6 = param_3[local_130 + 1];
        param_1[1] = iVar6;
        if (0 < iVar6) {
          iVar6 = param_1[iVar6 * 0x22 + -0x1c];
          *param_1 = iVar6;
          return CONCAT22((short)((uint)iVar6 >> 0x10),4);
        }
        goto LAB_1005207f;
      }
      if (local_134[2] < 3) {
        if ((*local_134 == 1) &&
           (uVar5 = FUN_100524f0(param_4,(int *)&local_134,param_3,param_5,param_6,param_7,param_8),
           (short)uVar5 != 0)) {
          local_118 = 1;
          local_134 = FUN_10054420(param_3,param_5,param_6,param_7,'D');
          if (local_134 == (int *)0x0) {
            *param_2 = local_12c;
            *param_3 = local_130;
            *param_6 = param_3[local_130 + 1];
            param_1[1] = iVar6;
            if (0 < iVar6) {
              iVar6 = param_1[iVar6 * 0x22 + -0x1c];
              *param_1 = iVar6;
              return CONCAT22((short)((uint)iVar6 >> 0x10),4);
            }
            goto LAB_1005207f;
          }
          if (local_134[2] < 2) {
            iVar4 = *local_134;
            goto joined_r0x10051e4c;
          }
          goto LAB_10051e88;
        }
        if (((local_134[2] < 3) && (*local_134 == 1)) &&
           (uVar5 = FUN_10052690(param_4,(int *)&local_134,param_3,param_5,param_6,param_7,0x49,
                                 param_8), (short)uVar5 != 0)) {
          local_118 = 1;
          goto LAB_10051f43;
        }
      }
      local_134 = FUN_100544b0(param_3,param_6);
      local_118 = 1;
    }
  }
LAB_10051f43:
  if ((param_8 == 0) || ((short)local_120 != 0)) {
    if ((short)local_118 == 0) goto LAB_10052422;
  }
  else {
    local_134 = FUN_10054420(param_3,param_5,param_6,param_7,'D');
    if (local_134 == (int *)0x0) {
      *param_2 = local_12c;
      *param_3 = local_130;
      *param_6 = param_3[local_130 + 1];
      param_1[1] = iVar6;
      if (0 < iVar6) {
        iVar6 = param_1[iVar6 * 0x22 + -0x1c];
        *param_1 = iVar6;
        return CONCAT22((short)((uint)iVar6 >> 0x10),4);
      }
LAB_1005207f:
      iVar6 = *param_1;
      *param_1 = iVar6;
      return CONCAT22((short)((uint)iVar6 >> 0x10),4);
    }
    if (((*local_134 == 8) && ((char)local_134[6] == '\0')) &&
       (local_134 = FUN_10054420(param_3,param_5,param_6,param_7,'D'), local_134 == (int *)0x0)) {
      *param_2 = local_12c;
      *param_3 = local_130;
      *param_6 = param_3[local_130 + 1];
      param_1[1] = iVar6;
      if (0 < iVar6) {
        iVar6 = param_1[iVar6 * 0x22 + -0x1c];
        *param_1 = iVar6;
        return CONCAT22((short)((uint)iVar6 >> 0x10),4);
      }
      goto LAB_1005207f;
    }
    if ((*local_134 != 5) && (*local_134 != 9)) goto LAB_10052422;
  }
  uVar5 = FUN_10054280(param_2,param_3,0xb);
  if ((short)uVar5 != 1) {
    *param_3 = local_130;
    *param_6 = param_3[local_130 + 1];
    *param_2 = local_12c;
    return uVar5;
  }
  local_128 = -1;
  if ((local_11c == (int *)0x0) || (local_128 = local_11c[4], piVar7 = local_11c, local_128 < 0)) {
    if (local_114 != (int *)0x0) {
      local_128 = local_114[4];
    }
    piVar7 = local_114;
    if (-1 < local_128) goto LAB_1005212a;
    if (local_110 != (int *)0x0) {
      local_128 = local_110[4];
      goto LAB_1005212a;
    }
  }
  else {
LAB_1005212a:
    if (local_110 != (int *)0x0) {
      FUN_10044d30(local_100,(byte *)(local_110 + 6),1);
      uVar5 = FUN_10047ef0(param_1,local_100,local_10c,local_128,0x44,0x44);
      piVar7 = extraout_ECX;
      if ((short)uVar5 != 1) {
        *param_3 = local_130;
        *param_6 = param_3[local_130 + 1];
        *param_2 = local_12c;
        param_1[1] = iVar6;
        if (0 < iVar6) {
          *param_1 = param_1[iVar6 * 0x22 + -0x1c];
          return uVar5;
        }
        goto LAB_10052404;
      }
    }
  }
  if (local_114 == (int *)0x0) {
LAB_10052336:
    sVar2 = (short)local_104;
  }
  else {
    local_120 = (byte *)(local_114 + 6);
    iVar4 = FUN_10059429(piVar7,local_120);
    if ((0 < iVar4) ||
       (((iVar4 = FUN_10059429(this,local_120), this_00 = extraout_ECX_00, iVar4 == 0 &&
         (local_11c != (int *)0x0)) &&
        (iVar4 = FUN_10059429(local_11c + 6,(byte *)(local_11c + 6)), this_00 = extraout_ECX_01,
        iVar4 != 0)))) {
      FUN_10044d30(local_100,local_120,1);
      uVar5 = FUN_10047ef0(param_1,local_100,local_10c,local_128,0x44,0x44);
      if ((short)uVar5 != 1) {
        *param_3 = local_130;
        *param_6 = param_3[local_130 + 1];
        *param_2 = local_12c;
        param_1[1] = iVar6;
        if (0 < iVar6) {
          *param_1 = param_1[iVar6 * 0x22 + -0x1c];
          return uVar5;
        }
        goto LAB_10052404;
      }
      goto LAB_10052336;
    }
    iVar4 = FUN_10059429(this_00,local_120);
    if (((iVar4 != 0) || (local_11c != (int *)0x0)) &&
       ((iVar4 = FUN_10059429(this_01,local_120), iVar4 != 0 ||
        (iVar4 = FUN_10059429(local_11c + 6,(byte *)(local_11c + 6)), iVar4 != 0))))
    goto LAB_10052336;
    if ((*param_4 == 0) && ((short)local_108 == 0)) {
      local_100[0] = s_o_clock_1008b9bc[0];
      local_100[1] = s_o_clock_1008b9bc[1];
      local_100[2] = s_o_clock_1008b9bc[2];
      local_100[3] = s_o_clock_1008b9bc[3];
      local_fc[0] = s_o_clock_1008b9bc[4];
      local_fc[1] = s_o_clock_1008b9bc[5];
      local_fc[2] = s_o_clock_1008b9bc[6];
      local_fc[3] = s_o_clock_1008b9bc[7];
      uVar5 = FUN_10047ef0(param_1,local_100,local_10c,local_128,0x44,0x44);
      if ((short)uVar5 != 1) {
        *param_3 = local_130;
        *param_6 = param_3[local_130 + 1];
        *param_2 = local_12c;
        param_1[1] = iVar6;
        if (0 < iVar6) {
          *param_1 = param_1[iVar6 * 0x22 + -0x1c];
          return uVar5;
        }
        goto LAB_10052404;
      }
    }
    sVar2 = 0;
  }
  if ((local_11c != (int *)0x0) && (sVar2 != 0)) {
    FUN_10044d30(local_100,(byte *)(local_11c + 6),1);
    uVar5 = FUN_10047ef0(param_1,local_100,local_10c,local_128,0x44,0x44);
    if ((short)uVar5 != 1) {
      *param_3 = local_130;
      *param_6 = param_3[local_130 + 1];
      *param_2 = local_12c;
      param_1[1] = iVar6;
      if (0 < iVar6) {
        *param_1 = param_1[iVar6 * 0x22 + -0x1c];
        return uVar5;
      }
      goto LAB_10052404;
    }
  }
  uVar5 = FUN_100483c0(param_1,(int *)param_4);
  if ((short)uVar5 == 1) {
    return CONCAT22((short)((uint)uVar5 >> 0x10),1);
  }
  *param_3 = local_130;
  *param_6 = param_3[local_130 + 1];
  *param_2 = local_12c;
  param_1[1] = iVar6;
  if (0 < iVar6) {
    *param_1 = param_1[iVar6 * 0x22 + -0x1c];
    return uVar5;
  }
LAB_10052404:
  *param_1 = *param_1;
  return uVar5;
}



/* ===== FUN_10054420 ===== */
/* Entry: 10054420 */

int * __cdecl FUN_10054420(int *param_1,int param_2,int *param_3,char *param_4,char param_5)

{
  int iVar1;
  
  if (99 < *param_1) {
    return (int *)0x0;
  }
  param_1[*param_1 + 1] = *param_3;
  if (param_5 == 'R') {
    FUN_10053c50(param_1 + *param_1 * 0xe + 0x65,param_3,param_4,0);
  }
  else {
    FUN_100528b0(param_1 + *param_1 * 0xe + 0x65,param_2,param_3,param_4);
  }
  iVar1 = *param_1;
  *param_1 = iVar1 + 1;
  param_1[iVar1 + 2] = *param_3;
  return param_1 + *param_1 * 0xe + 0x57;
}



/* ===== FUN_10055e60 ===== */
/* Entry: 10055e60 */

uint __cdecl FUN_10055e60(int *param_1,int *param_2,int *param_3,int *param_4,char *param_5)

{
  int iVar1;
  char cVar2;
  int iVar3;
  int iVar4;
  int iVar5;
  int iVar6;
  uint uVar7;
  uint uVar8;
  short sVar9;
  char *pcVar10;
  char *pcVar11;
  undefined1 local_94;
  char *local_84;
  int local_80;
  int local_7c [3];
  int local_70;
  int local_6c;
  undefined4 local_68;
  char local_64;
  char local_63 [31];
  char local_44 [68];
  
  iVar5 = FUN_10058300((int)param_3);
  if ((short)iVar5 != 0) {
    iVar3 = *param_2;
    iVar4 = param_1[1];
    iVar6 = FUN_100562e0(param_4,param_5,local_7c,param_3,0);
    iVar5 = iVar6;
    if (-1 < iVar6) {
      iVar5 = FUN_10054220(param_2,local_7c,0x19);
      if ((short)iVar5 == 1) {
        local_80 = local_70;
        iVar1 = param_3[4] + iVar6 * 8;
        if (*(char *)(param_3[4] + iVar6 * 8) != 'A') {
          uVar7 = FUN_10047cc0(param_1,&local_64,*(char **)(iVar1 + 4),local_70,local_6c,0x55,0x59);
          if ((short)uVar7 == 1) {
            return uVar7;
          }
          *param_2 = iVar3;
          param_1[1] = iVar4;
          if (0 < iVar4) {
            *param_1 = param_1[iVar4 * 0x22 + -0x1c];
            return uVar7;
          }
          *param_1 = *param_1;
          return uVar7;
        }
        uVar7 = 0xffffffff;
        pcVar10 = *(char **)(iVar1 + 4);
        do {
          pcVar11 = pcVar10;
          if (uVar7 == 0) break;
          uVar7 = uVar7 - 1;
          pcVar11 = pcVar10 + 1;
          cVar2 = *pcVar10;
          pcVar10 = pcVar11;
        } while (cVar2 != '\0');
        uVar7 = ~uVar7;
        pcVar10 = pcVar11 + -uVar7;
        pcVar11 = local_44;
        for (uVar8 = uVar7 >> 2; uVar8 != 0; uVar8 = uVar8 - 1) {
          *(undefined4 *)pcVar11 = *(undefined4 *)pcVar10;
          pcVar10 = pcVar10 + 4;
          pcVar11 = pcVar11 + 4;
        }
        for (uVar7 = uVar7 & 3; uVar7 != 0; uVar7 = uVar7 - 1) {
          *pcVar11 = *pcVar10;
          pcVar10 = pcVar10 + 1;
          pcVar11 = pcVar11 + 1;
        }
        local_84 = local_44;
        do {
          sVar9 = (short)iVar5;
          FUN_10053c50(local_7c,(int *)&local_84,local_44,1);
          if (local_7c[0] == 9) goto LAB_1005601a;
          if ((local_64 == '<') && (*(char *)((int)&local_68 + local_68 + 3) == '>')) {
            uVar7 = 0xffffffff;
            pcVar10 = local_63;
            do {
              pcVar11 = pcVar10;
              if (uVar7 == 0) break;
              uVar7 = uVar7 - 1;
              pcVar11 = pcVar10 + 1;
              cVar2 = *pcVar10;
              pcVar10 = pcVar11;
            } while (cVar2 != '\0');
            uVar7 = ~uVar7;
            local_94 = 0x53;
            pcVar10 = pcVar11 + -uVar7;
            pcVar11 = &local_64;
            for (uVar8 = uVar7 >> 2; uVar8 != 0; uVar8 = uVar8 - 1) {
              *(undefined4 *)pcVar11 = *(undefined4 *)pcVar10;
              pcVar10 = pcVar10 + 4;
              pcVar11 = pcVar11 + 4;
            }
            for (uVar7 = uVar7 & 3; uVar7 != 0; uVar7 = uVar7 - 1) {
              *pcVar11 = *pcVar10;
              pcVar10 = pcVar10 + 1;
              pcVar11 = pcVar11 + 1;
            }
            *(undefined1 *)((int)&local_68 + local_68 + 2) = 0;
            local_68 = local_68 + -2;
          }
          else {
            local_94 = 0x44;
          }
          iVar5 = FUN_10047b90(param_1,&local_64,local_80,local_6c,0x55,local_94);
          sVar9 = (short)iVar5;
        } while (sVar9 == 1);
        *param_2 = iVar3;
        param_1[1] = iVar4;
        if (0 < iVar4) {
          iVar5 = param_1[iVar4 * 0x22 + -0x1c];
          *param_1 = iVar5;
          return CONCAT22((short)((uint)iVar5 >> 0x10),sVar9);
        }
        local_7c[0] = *param_1;
        *param_1 = local_7c[0];
LAB_1005601a:
        return CONCAT22((short)((uint)local_7c[0] >> 0x10),sVar9);
      }
    }
  }
  return CONCAT22((short)((uint)iVar5 >> 0x10),5);
}



/* ===== FUN_10058300 ===== */
/* Entry: 10058300 */

undefined4 __cdecl FUN_10058300(int param_1)

{
  ushort uVar1;
  
  uVar1 = (ushort)((uint)param_1 >> 0x10);
  if ((param_1 != 0) && (*(short *)(param_1 + 0x14) != 0)) {
    return CONCAT22(uVar1,1);
  }
  return (uint)uVar1 << 0x10;
}



/* ===== FUN_10057930 ===== */
/* Entry: 10057930 */

int * __cdecl FUN_10057930(int *param_1,int *param_2,int param_3,uint param_4)

{
  int iVar1;
  int iVar2;
  
  if ((0x7fffffff < param_4) && (-1 < param_3)) {
    if ((param_3 < *param_2) && (iVar2 = param_3 + param_4, iVar2 < *param_2)) {
      if (-1 < iVar2) {
        return param_2 + iVar2 * 0xe + 0x65;
      }
      iVar2 = *param_1 + 1 + iVar2;
      iVar1 = iVar2 + -1;
      if ((-1 < iVar1) && (iVar1 <= *param_1 + -1)) {
        return param_1 + iVar2 * 0xe + -0xd;
      }
    }
  }
  return (int *)0x0;
}



/* ===== FUN_10054220 ===== */
/* Entry: 10054220 */

undefined4 __cdecl FUN_10054220(int *param_1,int *param_2,int param_3)

{
  undefined4 in_EAX;
  undefined2 uVar1;
  int iVar2;
  int *piVar3;
  
  if ((char)param_2[6] == '\0') {
    return CONCAT22((short)((uint)in_EAX >> 0x10),1);
  }
  uVar1 = (undefined2)((uint)param_1 >> 0x10);
  if (99 < *param_1) {
    return CONCAT22(uVar1,4);
  }
  piVar3 = param_1 + *param_1 * 0xe + 1;
  for (iVar2 = 0xe; iVar2 != 0; iVar2 = iVar2 + -1) {
    *piVar3 = *param_2;
    param_2 = param_2 + 1;
    piVar3 = piVar3 + 1;
  }
  param_1[*param_1 * 0xe + 2] = param_3;
  *param_1 = *param_1 + 1;
  return CONCAT22(uVar1,1);
}



/* ===== FUN_10020510 ===== */
/* Entry: 10020510 */

int __cdecl FUN_10020510(int *param_1,int *param_2,int param_3)

{
  short *psVar1;
  int iVar2;
  int iVar3;
  int iVar4;
  
  param_2[0xf] = param_3;
  FUN_10020460(param_1,param_2);
  iVar2 = param_2[0xe];
  if (param_1[0xe0f] != 0) {
    if (0 < param_1[0xe0f]) {
      psVar1 = (short *)param_2[0xf];
      iVar2 = iVar2 / 2;
      iVar4 = 0;
      if (0 < iVar2) {
        do {
          if (iVar4 == 0) {
            iVar3 = (int)*psVar1 + ((int)(short)param_2[0x47320] * param_1[0xe0f] + 0x32) / 100;
            if ((0x7ffe < iVar3) || (iVar3 < -0x7fff)) {
              iVar3 = ((iVar3 < 0x7fff) - 1 & 0xfffd) - 0x7fff;
            }
            *(short *)((int)param_2 + 0x11e0a6) = (short)iVar3;
          }
          else {
            iVar3 = ((int)*(short *)((int)param_2 + iVar4 * 2 + 0x11e0a4) * param_1[0xe0f] + 0x32) /
                    100 + (int)psVar1[iVar4];
            if ((0x7ffe < iVar3) || (iVar3 < -0x7fff)) {
              iVar3 = ((iVar3 < 0x7fff) - 1 & 0xfffd) - 0x7fff;
            }
            *(short *)((int)param_2 + iVar4 * 2 + 0x11e0a6) = (short)iVar3;
          }
          iVar4 = iVar4 + 1;
        } while (iVar4 < iVar2);
      }
      *(undefined2 *)(param_2 + 0x47320) = *(undefined2 *)((int)param_2 + iVar4 * 2 + 0x11e0a4);
      FUN_10058950((undefined4 *)psVar1,(undefined4 *)((int)param_2 + 0x11e0a6),iVar2 * 2);
      return param_2[0xe];
    }
    psVar1 = (short *)param_2[0xf];
    iVar2 = iVar2 / 2;
    iVar4 = iVar2;
    while (iVar3 = iVar4 + -1, -1 < iVar3) {
      if (iVar3 == 0) {
        iVar4 = (int)*psVar1 + ((int)(short)param_2[0x47320] * param_1[0xe0f] + 0x32) / 100;
        if ((0x7ffe < iVar4) || (iVar4 < -0x7fff)) {
          iVar4 = ((iVar4 < 0x7fff) - 1 & 0xfffd) - 0x7fff;
        }
        *(short *)((int)param_2 + 0x11e0a6) = (short)iVar4;
        iVar4 = iVar3;
      }
      else {
        iVar4 = ((int)psVar1[iVar4 + -2] * param_1[0xe0f] + 0x32) / 100 + (int)psVar1[iVar3];
        if ((0x7ffe < iVar4) || (iVar4 < -0x7fff)) {
          iVar4 = ((iVar4 < 0x7fff) - 1 & 0xfffd) - 0x7fff;
        }
        *(short *)((int)param_2 + iVar3 * 2 + 0x11e0a6) = (short)iVar4;
        iVar4 = iVar3;
      }
    }
    *(short *)(param_2 + 0x47320) = psVar1[iVar2 + -1];
    FUN_10058950((undefined4 *)psVar1,(undefined4 *)((int)param_2 + 0x11e0a6),iVar2 * 2);
    iVar2 = param_2[0xe];
  }
  return iVar2;
}



/* ===== FUN_10020460 ===== */
/* Entry: 10020460 */

uint __cdecl FUN_10020460(int *param_1,int *param_2)

{
  short sVar1;
  short *psVar2;
  uint uVar3;
  
  uVar3 = param_2[0x13];
  param_2[0xe] = 0;
  while (uVar3 == 0) {
    if ((short)param_2[0x12] == 0) {
      *(undefined2 *)(param_2[0x15] + 0x5285e) = 0;
    }
    uVar3 = FUN_10022f60(param_1,(float)param_2);
    sVar1 = (short)param_2[0x12];
    uVar3 = CONCAT22((short)(uVar3 >> 0x10),sVar1);
    if (sVar1 == 1) break;
    if (sVar1 == 0) {
      *(short *)(param_2[0x15] + 0x5285c) = *(short *)(param_2[0x15] + 0x5285c) + 1;
      psVar2 = (short *)param_2[0x15];
      if ((psVar2[1] == 0) || (psVar2[0x2942e] == *psVar2)) {
        uVar3 = FUN_100202f0((int)param_1,param_2);
        while (((short)uVar3 < 0 || (*(short *)(param_2[0x15] + 2) == 0))) {
          uVar3 = FUN_100202f0((int)param_1,param_2);
          if (param_2[0x13] != 0) {
            return uVar3;
          }
        }
      }
    }
    uVar3 = param_2[0x13];
  }
  return uVar3 & 0xffff0000;
}



/* ===== FUN_10015370 ===== */
/* Entry: 10015370 */

undefined4 __cdecl FUN_10015370(int param_1,byte *param_2,int param_3)

{
  size_t sVar1;
  int *piVar2;
  undefined4 *puVar3;
  byte bVar4;
  int *piVar5;
  undefined4 uVar6;
  int iVar7;
  byte *pbVar8;
  char *pcVar9;
  uint uVar10;
  int iVar11;
  uint uVar12;
  uint uVar13;
  int iVar14;
  byte *pbVar15;
  byte *pbVar16;
  
  if (param_2 == (byte *)0x0) {
    return 0;
  }
  if (*param_2 == 0) {
    return 0;
  }
  uVar10 = 0xffffffff;
  pbVar8 = param_2;
  do {
    if (uVar10 == 0) break;
    uVar10 = uVar10 - 1;
    bVar4 = *pbVar8;
    pbVar8 = pbVar8 + 1;
  } while (bVar4 != 0);
  uVar10 = ~uVar10;
  iVar11 = uVar10 - 1;
  if (param_3 == 0) {
    pbVar8 = (byte *)FUN_10016490(uVar10);
    if (pbVar8 == (byte *)0x0) {
      return 0;
    }
    uVar12 = 0xffffffff;
    do {
      pbVar15 = param_2;
      if (uVar12 == 0) break;
      uVar12 = uVar12 - 1;
      pbVar15 = param_2 + 1;
      bVar4 = *param_2;
      param_2 = pbVar15;
    } while (bVar4 != 0);
    uVar12 = ~uVar12;
    pbVar15 = pbVar15 + -uVar12;
    pbVar16 = pbVar8;
    for (uVar13 = uVar12 >> 2; uVar13 != 0; uVar13 = uVar13 - 1) {
      *(undefined4 *)pbVar16 = *(undefined4 *)pbVar15;
      pbVar15 = pbVar15 + 4;
      pbVar16 = pbVar16 + 4;
    }
    for (uVar12 = uVar12 & 3; param_2 = pbVar8, uVar12 != 0; uVar12 = uVar12 - 1) {
      *pbVar16 = *pbVar15;
      pbVar15 = pbVar15 + 1;
      pbVar16 = pbVar16 + 1;
    }
  }
  if (*(LPVOID *)(param_1 + 0x528e8) != (LPVOID)0x0) {
    FUN_10016500(*(LPVOID *)(param_1 + 0x528e8));
  }
  if (*(LPVOID *)(param_1 + 0x528ec) != (LPVOID)0x0) {
    FUN_10016500(*(LPVOID *)(param_1 + 0x528ec));
  }
  if (*(LPVOID *)(param_1 + 0x528e4) != (LPVOID)0x0) {
    FUN_10016500(*(LPVOID *)(param_1 + 0x528e4));
  }
  sVar1 = iVar11 * 4;
  uVar6 = FUN_10016490(sVar1);
  *(undefined4 *)(param_1 + 0x528e8) = uVar6;
  uVar6 = FUN_10016490(sVar1);
  *(undefined4 *)(param_1 + 0x528ec) = uVar6;
  uVar6 = FUN_10016490(sVar1);
  *(undefined4 *)(param_1 + 0x528e4) = uVar6;
  iVar7 = 0;
  if (0 < iVar11) {
    do {
      *(int *)(*(int *)(param_1 + 0x528e8) + iVar7 * 4) = iVar7;
      *(int *)(*(int *)(param_1 + 0x528ec) + iVar7 * 4) = iVar7;
      iVar7 = iVar7 + 1;
    } while (iVar7 < iVar11);
  }
  iVar7 = 0;
  iVar14 = 0;
  if (0 < iVar11) {
    do {
      bVar4 = param_2[iVar7];
      *(int *)(*(int *)(param_1 + 0x528e4) + iVar7 * 4) = iVar14;
      if (((bVar4 & 0x80) == 0) || (iVar7 == uVar10 - 2)) {
        iVar7 = iVar7 + 1;
      }
      else {
        iVar7 = iVar7 + 2;
        *(int *)(*(int *)(param_1 + 0x528e4) + -4 + iVar7 * 4) = iVar14;
      }
      iVar14 = iVar14 + 1;
    } while (iVar7 < iVar11);
  }
  pbVar8 = FUN_100157b0(param_1,param_2,param_3);
  pcVar9 = (char *)FUN_10015b50(param_1,(char *)pbVar8,&LAB_10015c70);
  if (*(char *)(param_1 + 0x11cc82) != '\0') {
    pcVar9 = (char *)FUN_10015b50(param_1,pcVar9,&LAB_100251e0);
  }
  uVar6 = FUN_10015b50(param_1,pcVar9,&LAB_10015f70);
  piVar5 = *(int **)(param_1 + 0x11dc9c);
  iVar11 = *piVar5;
  if (0 < iVar11) {
    if (*(char *)(DAT_100942d4 + 0x20494) == '\x01') {
      if (0 < iVar11) {
        iVar7 = 0;
        do {
          piVar2 = (int *)(piVar5[1] + 8 + iVar7);
          puVar3 = (undefined4 *)(piVar5[1] + 8 + iVar7);
          iVar7 = iVar7 + 0x110;
          iVar11 = iVar11 + -1;
          *puVar3 = *(undefined4 *)
                     (*(int *)(param_1 + 0x528e4) +
                     *(int *)(*(int *)(param_1 + 0x528e8) + *piVar2 * 4) * 4);
        } while (iVar11 != 0);
        return uVar6;
      }
    }
    else if (0 < iVar11) {
      iVar7 = 0;
      do {
        piVar2 = (int *)(piVar5[1] + 8 + iVar7);
        iVar7 = iVar7 + 0x110;
        iVar11 = iVar11 + -1;
        *piVar2 = *(int *)(*(int *)(param_1 + 0x528e8) + *piVar2 * 4);
      } while (iVar11 != 0);
    }
  }
  return uVar6;
}



/* ===== FUN_100157b0 ===== */
/* Entry: 100157b0 */

byte * __cdecl FUN_100157b0(int param_1,byte *param_2,int param_3)

{
  int iVar1;
  size_t sVar2;
  char cVar3;
  byte bVar4;
  int iVar5;
  undefined4 uVar6;
  byte *pbVar7;
  uint uVar8;
  uint uVar9;
  uint uVar10;
  int iVar11;
  byte *pbVar12;
  char *pcVar13;
  bool bVar14;
  
  pbVar7 = param_2;
  pbVar12 = PTR_s_BaekJongKwanLeeYunKeunLeeJoonWoo_1006f988;
  if (param_3 != 0) {
    return param_2;
  }
  do {
    bVar4 = *pbVar7;
    bVar14 = bVar4 < *pbVar12;
    if (bVar4 != *pbVar12) {
LAB_100157f9:
      iVar5 = (1 - (uint)bVar14) - (uint)(bVar14 != 0);
      goto LAB_100157fe;
    }
    if (bVar4 == 0) break;
    bVar4 = pbVar7[1];
    bVar14 = bVar4 < pbVar12[1];
    if (bVar4 != pbVar12[1]) goto LAB_100157f9;
    pbVar7 = pbVar7 + 2;
    pbVar12 = pbVar12 + 2;
  } while (bVar4 != 0);
  iVar5 = 0;
LAB_100157fe:
  if (iVar5 == 0) {
    uVar8 = 0xffffffff;
    pcVar13 = PTR_s_This_is_a_demo_of_the_VoiceText_E_1006f984;
    do {
      if (uVar8 == 0) break;
      uVar8 = uVar8 - 1;
      cVar3 = *pcVar13;
      pcVar13 = pcVar13 + 1;
    } while (cVar3 != '\0');
    *(uint *)(param_1 + 0x528f8) = ~uVar8 - 1;
    uVar8 = 0xffffffff;
    pbVar7 = param_2;
    do {
      if (uVar8 == 0) break;
      uVar8 = uVar8 - 1;
      bVar4 = *pbVar7;
      pbVar7 = pbVar7 + 1;
    } while (bVar4 != 0);
    sVar2 = (~uVar8 - 1) * 4;
    uVar6 = FUN_10016490(sVar2);
    *(undefined4 *)(param_1 + 0x528f0) = uVar6;
    uVar6 = FUN_10016490(sVar2);
    *(undefined4 *)(param_1 + 0x528f4) = uVar6;
    FUN_10058950(*(undefined4 **)(param_1 + 0x528f0),*(undefined4 **)(param_1 + 0x528e8),sVar2);
    FUN_10058950(*(undefined4 **)(param_1 + 0x528f4),*(undefined4 **)(param_1 + 0x528ec),sVar2);
    FUN_10016500(*(LPVOID *)(param_1 + 0x528e8));
    FUN_10016500(*(LPVOID *)(param_1 + 0x528ec));
    uVar6 = FUN_10016490(*(int *)(param_1 + 0x528f8) << 2);
    *(undefined4 *)(param_1 + 0x528e8) = uVar6;
    uVar6 = FUN_10016490(*(int *)(param_1 + 0x528f8) << 2);
    *(undefined4 *)(param_1 + 0x528ec) = uVar6;
    pbVar7 = (byte *)FUN_10016490(*(int *)(param_1 + 0x528f8) + 1);
    FUN_10058950((undefined4 *)pbVar7,(undefined4 *)PTR_s_This_is_a_demo_of_the_VoiceText_E_1006f984
                 ,*(uint *)(param_1 + 0x528f8));
    iVar5 = 0;
    pbVar7[*(int *)(param_1 + 0x528f8)] = 0;
    if (0 < *(int *)(param_1 + 0x528f8)) {
      do {
        iVar5 = iVar5 + 1;
        *(undefined4 *)(*(int *)(param_1 + 0x528e8) + -4 + iVar5 * 4) =
             **(undefined4 **)(param_1 + 0x528f0);
        *(undefined4 *)(*(int *)(param_1 + 0x528ec) + -4 + iVar5 * 4) =
             *(undefined4 *)(*(int *)(param_1 + 0x528f4) + -4 + sVar2);
      } while (iVar5 < *(int *)(param_1 + 0x528f8));
    }
    FUN_10016500(*(LPVOID *)(param_1 + 0x528f0));
    FUN_10016500(*(LPVOID *)(param_1 + 0x528f4));
    FUN_10016500(param_2);
    *(undefined4 *)(param_1 + 0x528f8) = 0xffffffff;
    return pbVar7;
  }
  pbVar7 = param_2;
  if (DAT_100962e6 == '\0') {
    param_3 = FUN_10015710(0,7);
    if (param_3 < 0) {
      param_3 = 0;
    }
    else if (7 < param_3) {
      param_3 = 7;
    }
    uVar8 = 0xffffffff;
    pcVar13 = (&PTR_s_You_do_not_have_a_valid_verifica_1006f964)[param_3];
    do {
      if (uVar8 == 0) break;
      uVar8 = uVar8 - 1;
      cVar3 = *pcVar13;
      pcVar13 = pcVar13 + 1;
    } while (cVar3 != '\0');
    uVar8 = ~uVar8;
    uVar9 = uVar8 - 1;
    uVar10 = 0xffffffff;
    do {
      if (uVar10 == 0) break;
      uVar10 = uVar10 - 1;
      bVar4 = *pbVar7;
      pbVar7 = pbVar7 + 1;
    } while (bVar4 != 0);
    uVar10 = ~uVar10 - 1;
    sVar2 = uVar10 * 4;
    *(uint *)(param_1 + 0x528f8) = uVar8 + uVar10;
    uVar6 = FUN_10016490(sVar2);
    *(undefined4 *)(param_1 + 0x528f0) = uVar6;
    uVar6 = FUN_10016490(sVar2);
    *(undefined4 *)(param_1 + 0x528f4) = uVar6;
    FUN_10058950(*(undefined4 **)(param_1 + 0x528f0),*(undefined4 **)(param_1 + 0x528e8),sVar2);
    FUN_10058950(*(undefined4 **)(param_1 + 0x528f4),*(undefined4 **)(param_1 + 0x528ec),sVar2);
    FUN_10016500(*(LPVOID *)(param_1 + 0x528e8));
    FUN_10016500(*(LPVOID *)(param_1 + 0x528ec));
    uVar6 = FUN_10016490(*(int *)(param_1 + 0x528f8) << 2);
    *(undefined4 *)(param_1 + 0x528e8) = uVar6;
    uVar6 = FUN_10016490(*(int *)(param_1 + 0x528f8) << 2);
    *(undefined4 *)(param_1 + 0x528ec) = uVar6;
    pbVar7 = (byte *)FUN_10016490(*(int *)(param_1 + 0x528f8) + 1);
    FUN_10058950((undefined4 *)pbVar7,
                 (undefined4 *)(&PTR_s_You_do_not_have_a_valid_verifica_1006f964)[param_3],uVar9);
    pbVar7[uVar9] = 0x20;
    FUN_10058950((undefined4 *)(pbVar7 + uVar8),(undefined4 *)param_2,uVar10);
    iVar5 = 0;
    pbVar7[*(int *)(param_1 + 0x528f8)] = 0;
    if (0 < (int)uVar9) {
      do {
        iVar5 = iVar5 + 1;
        *(undefined4 *)(*(int *)(param_1 + 0x528e8) + -4 + iVar5 * 4) =
             **(undefined4 **)(param_1 + 0x528f0);
        *(undefined4 *)(*(int *)(param_1 + 0x528ec) + -4 + iVar5 * 4) =
             **(undefined4 **)(param_1 + 0x528f4);
      } while (iVar5 < (int)uVar9);
    }
    *(undefined4 *)(*(int *)(param_1 + 0x528e8) + uVar9 * 4) = **(undefined4 **)(param_1 + 0x528f0);
    *(undefined4 *)(*(int *)(param_1 + 0x528ec) + uVar9 * 4) = **(undefined4 **)(param_1 + 0x528f4);
    iVar5 = 0;
    if (0 < (int)uVar10) {
      iVar11 = uVar9 * 4;
      do {
        iVar11 = iVar11 + 4;
        iVar1 = iVar5 * 4;
        iVar5 = iVar5 + 1;
        *(undefined4 *)(iVar11 + *(int *)(param_1 + 0x528e8)) =
             *(undefined4 *)(*(int *)(param_1 + 0x528f0) + iVar1);
        *(undefined4 *)(iVar11 + *(int *)(param_1 + 0x528ec)) =
             *(undefined4 *)(*(int *)(param_1 + 0x528f4) + -4 + iVar5 * 4);
      } while (iVar5 < (int)uVar10);
    }
    FUN_10016500(*(LPVOID *)(param_1 + 0x528f0));
    FUN_10016500(*(LPVOID *)(param_1 + 0x528f4));
    FUN_10016500(param_2);
    *(undefined4 *)(param_1 + 0x528f8) = 0xffffffff;
  }
  return pbVar7;
}



/* ===== FUN_10015b50 ===== */
/* Entry: 10015b50 */

undefined4 __cdecl FUN_10015b50(int param_1,char *param_2,undefined *param_3)

{
  size_t sVar1;
  char cVar2;
  undefined4 uVar3;
  uint uVar4;
  char *pcVar5;
  
  *(undefined4 *)(param_1 + 0x528f8) = 0xffffffff;
  (*(code *)param_3)(param_1,param_2,0);
  uVar4 = 0xffffffff;
  pcVar5 = param_2;
  do {
    if (uVar4 == 0) break;
    uVar4 = uVar4 - 1;
    cVar2 = *pcVar5;
    pcVar5 = pcVar5 + 1;
  } while (cVar2 != '\0');
  sVar1 = (~uVar4 - 1) * 4;
  uVar3 = FUN_10016490(sVar1);
  *(undefined4 *)(param_1 + 0x528f0) = uVar3;
  uVar3 = FUN_10016490(sVar1);
  *(undefined4 *)(param_1 + 0x528f4) = uVar3;
  FUN_10058950(*(undefined4 **)(param_1 + 0x528f0),*(undefined4 **)(param_1 + 0x528e8),sVar1);
  FUN_10058950(*(undefined4 **)(param_1 + 0x528f4),*(undefined4 **)(param_1 + 0x528ec),sVar1);
  FUN_10016500(*(LPVOID *)(param_1 + 0x528e8));
  FUN_10016500(*(LPVOID *)(param_1 + 0x528ec));
  uVar3 = FUN_10016490(*(int *)(param_1 + 0x528f8) << 2);
  *(undefined4 *)(param_1 + 0x528e8) = uVar3;
  uVar3 = FUN_10016490(*(int *)(param_1 + 0x528f8) << 2);
  *(undefined4 *)(param_1 + 0x528ec) = uVar3;
  uVar3 = FUN_10016490(*(int *)(param_1 + 0x528f8) + 1);
  (*(code *)param_3)(param_1,param_2,uVar3);
  FUN_10016500(*(LPVOID *)(param_1 + 0x528f0));
  FUN_10016500(*(LPVOID *)(param_1 + 0x528f4));
  FUN_10016500(param_2);
  *(undefined4 *)(param_1 + 0x528f8) = 0xffffffff;
  return uVar3;
}



/* ===== FUN_10016c00 ===== */
/* Entry: 10016c00 */

undefined4 __cdecl FUN_10016c00(int param_1)

{
  int *piVar1;
  int *piVar2;
  int iVar3;
  int iVar4;
  short sVar5;
  int iVar6;
  int iVar7;
  int iVar8;
  int iVar9;
  undefined4 local_100;
  
  piVar1 = *(int **)(param_1 + 0x2c);
  piVar2 = *(int **)(param_1 + 0x528d8);
  if (piVar1 == (int *)0x0) {
    return CONCAT22((short)((uint)param_1 >> 0x10),1);
  }
  if (*(int *)(param_1 + 0x4c) == 0) {
    if (piVar1[2] == (int)(short)piVar2[6]) goto LAB_10016dad;
  }
  else {
    piVar1[2] = (int)(short)piVar2[6];
  }
  iVar6 = piVar1[2];
  if (iVar6 < *(short *)((int)piVar2 + 0x1a)) {
    iVar6 = iVar6 + 200;
  }
  for (iVar4 = (int)*(short *)((int)piVar2 + 0x1a); iVar4 < iVar6; iVar4 = iVar4 + 1) {
    iVar8 = 0;
    iVar7 = (iVar4 % 200) * 0x20;
    iVar9 = *piVar2;
    if (0 < *(short *)(iVar9 + 0x10 + iVar7)) {
      do {
        iVar3 = *(int *)(*(int *)(iVar9 + iVar7) + iVar8 * 8);
        if (*(short *)(*(int *)(iVar9 + iVar7) + iVar8 * 8 + 4) == 0x5a) {
          FUN_1001f890((undefined4 *)piVar1[4],(byte *)s_0__sil_____d_1006ffb8);
        }
        else {
          FUN_1001f890((undefined4 *)piVar1[4],(byte *)s__2d___3s_____d_1006ffa4);
        }
        iVar8 = iVar8 + 1;
        piVar1[1] = piVar1[1] + iVar3;
        iVar9 = *piVar2;
      } while (iVar8 < *(short *)(iVar9 + 0x10 + iVar7));
    }
    iVar9 = *piVar2;
    iVar8 = *(int *)(iVar9 + 8 + iVar7);
    iVar3 = *(int *)(iVar9 + 4 + iVar7);
    iVar9 = *(int *)(iVar9 + 0xc + iVar7) - iVar8;
    FUN_10058950(&local_100,(undefined4 *)(*(int *)(param_1 + 8) + iVar8),iVar9 + 1);
    *(undefined1 *)((int)&local_100 + iVar9 + 1) = 0;
    FUN_1001f890((undefined4 *)piVar1[4],(byte *)s_________________________________1006ff80);
    FUN_1001f890((undefined4 *)piVar1[4],(byte *)s__d__d___s_____d_1006ff6c);
    FUN_1001f890((undefined4 *)piVar1[4],(byte *)s_________________________________1006ff48);
    *piVar1 = *piVar1 + iVar3;
  }
  sVar5 = (short)iVar6;
  *(short *)((int)piVar2 + 0x1a) = sVar5;
  if (199 < sVar5) {
    *(short *)((int)piVar2 + 0x1a) = sVar5 + -200;
  }
LAB_10016dad:
  iVar6 = piVar2[6];
  piVar1[2] = (int)(short)iVar6;
  return CONCAT22((short)iVar6 >> 0xf,1);
}



/* ===== FUN_1001f950 ===== */
/* Entry: 1001f950 */

int * __cdecl FUN_1001f950(HWND param_1,undefined2 *param_2,int param_3)

{
  int *piVar1;
  undefined4 uVar2;
  
  piVar1 = (int *)FUN_10016490(0x3850);
  if (piVar1 == (int *)0x0) {
    *param_2 = 2;
    return (int *)0x0;
  }
  piVar1[0xe0d] = 100;
  piVar1[0xe0b] = 100;
  piVar1[0xe0c] = 100;
  piVar1[0xe10] = 0x2af;
  piVar1[0xe0e] = 100;
  piVar1[0xe0f] = 0;
  FUN_10014f90((int)piVar1,param_3);
  uVar2 = FUN_1000eff0(piVar1,param_2,param_3,param_1);
  if ((short)uVar2 < 0) {
    return (int *)0x0;
  }
  *param_2 = 0;
  (&DAT_100962e4)[param_3] = 1;
  return piVar1;
}



/* ===== FUN_1000eff0 ===== */
/* Entry: 1000eff0 */

undefined4 __cdecl FUN_1000eff0(int *param_1,undefined2 *param_2,undefined4 param_3,HWND param_4)

{
  undefined4 uVar1;
  undefined2 extraout_var;
  
  FUN_10021c60(1,param_4,0x32);
  FUN_1000e890();
  uVar1 = FUN_1000eb20((int)param_1,param_2);
  if ((short)uVar1 < 0) {
    return CONCAT22((short)((uint)uVar1 >> 0x10),0xffff);
  }
  FUN_1000e890();
  FUN_1000e890();
  FUN_1000e890();
  FUN_10021c60(1,param_4,0x3c);
  uVar1 = FUN_1000ec70((int)param_1,param_2);
  if ((short)uVar1 < 0) {
    return CONCAT22((short)((uint)uVar1 >> 0x10),0xffff);
  }
  FUN_1000e890();
  FUN_1000e890();
  FUN_1000e890();
  FUN_10021c60(1,param_4,0x46);
  uVar1 = FUN_1000ed30((int)param_1,param_2);
  if ((short)uVar1 < 0) {
    return CONCAT22((short)((uint)uVar1 >> 0x10),0xffff);
  }
  FUN_1000e890();
  FUN_1000e890();
  FUN_1000e890();
  if (DAT_100942d1 != '\0') {
    uVar1 = FUN_1000eb60((int)param_1);
    if ((short)uVar1 < 0) {
      return CONCAT22((short)((uint)uVar1 >> 0x10),0xffff);
    }
  }
  FUN_10021c60(1,param_4,0x50);
  uVar1 = FUN_1000ee40(param_1,param_2);
  if ((short)uVar1 < 0) {
    return CONCAT22((short)((uint)uVar1 >> 0x10),0xffff);
  }
  FUN_1000e890();
  FUN_1000e890();
  FUN_1000e890();
  FUN_10021c60(1,param_4,100);
  FUN_1000e890();
  FUN_1000e890();
  FUN_1000e890();
  return CONCAT22(extraout_var,1);
}



/* ===== FUN_10001000 ===== */
/* Entry: 10001000 */

int __cdecl FUN_10001000(undefined4 param_1,int param_2)

{
  undefined4 uVar1;
  undefined1 local_200 [512];
  
  FUN_100588f6(local_200,(byte *)s__spitch_nbt_tree2_1006b194);
  uVar1 = FUN_100017d0((int *)(param_2 + 0x1b4),local_200);
  if ((short)uVar1 == 0) {
    return -1;
  }
  FUN_100588f6(local_200,(byte *)s__spitch_nbf_tree2_1006b180);
  uVar1 = FUN_10001c80((int *)(param_2 + 0x224),local_200);
  if ((short)uVar1 == 0) {
    return -1;
  }
  FUN_100588f6(local_200,(byte *)s__spitch_bt_tree2_1006b16c);
  uVar1 = FUN_100017d0((int *)(param_2 + 0x1d0),local_200);
  if ((short)uVar1 == 0) {
    return -1;
  }
  FUN_100588f6(local_200,(byte *)s__spitch_bf_tree2_1006b158);
  uVar1 = FUN_10001c80((int *)(param_2 + 0x240),local_200);
  if ((short)uVar1 == 0) {
    return -1;
  }
  FUN_100588f6(local_200,(byte *)s__spitch_qbt_tree2_1006b144);
  uVar1 = FUN_100017d0((int *)(param_2 + 0x208),local_200);
  if ((short)uVar1 == 0) {
    return -1;
  }
  FUN_100588f6(local_200,(byte *)s__spitch_qbf_tree2_1006b130);
  uVar1 = FUN_10001c80((int *)(param_2 + 0x278),local_200);
  if ((short)uVar1 == 0) {
    return -1;
  }
  FUN_100588f6(local_200,(byte *)s__spitch_sbt_tree2_1006b11c);
  uVar1 = FUN_100017d0((int *)(param_2 + 0x1ec),local_200);
  if ((short)uVar1 == 0) {
    return -1;
  }
  FUN_100588f6(local_200,(byte *)s__spitch_sbf_tree2_1006b108);
  uVar1 = FUN_10001c80((int *)(param_2 + 0x25c),local_200);
  if ((short)uVar1 == 0) {
    return -1;
  }
  FUN_100588f6(local_200,(byte *)s__sduration_vshort_tree2_1006b0f0);
  uVar1 = FUN_100017d0((int *)(param_2 + 0xb8),local_200);
  if ((short)uVar1 == 0) {
    return -1;
  }
  FUN_100588f6(local_200,(byte *)s__sduration_vlong_tree2_1006b0d8);
  uVar1 = FUN_100017d0((int *)(param_2 + 0xd4),local_200);
  if ((short)uVar1 == 0) {
    return -1;
  }
  FUN_100588f6(local_200,(byte *)s__sduration_vdi_tree2_1006b0c0);
  uVar1 = FUN_100017d0((int *)(param_2 + 0xf0),local_200);
  if ((short)uVar1 == 0) {
    return -1;
  }
  FUN_100588f6(local_200,(byte *)s__sduration_vsch_tree2_1006b0a8);
  uVar1 = FUN_100017d0((int *)(param_2 + 0x10c),local_200);
  if ((short)uVar1 == 0) {
    return -1;
  }
  FUN_100588f6(local_200,(byte *)s__sduration_cstop_tree2_1006b090);
  uVar1 = FUN_100017d0((int *)(param_2 + 0x128),local_200);
  if ((short)uVar1 == 0) {
    return -1;
  }
  FUN_100588f6(local_200,(byte *)s__sduration_cfri_tree2_1006b078);
  uVar1 = FUN_100017d0((int *)(param_2 + 0x144),local_200);
  if ((short)uVar1 == 0) {
    return -1;
  }
  FUN_100588f6(local_200,(byte *)s__sduration_caff_tree2_1006b060);
  uVar1 = FUN_100017d0((int *)(param_2 + 0x160),local_200);
  if ((short)uVar1 == 0) {
    return -1;
  }
  FUN_100588f6(local_200,(byte *)s__sduration_cnas_tree2_1006b048);
  uVar1 = FUN_100017d0((int *)(param_2 + 0x17c),local_200);
  if ((short)uVar1 == 0) {
    return -1;
  }
  FUN_100588f6(local_200,(byte *)s__sduration_capp_tree2_1006b030);
  uVar1 = FUN_100017d0((int *)(param_2 + 0x198),local_200);
  return (-(uint)((short)uVar1 != 0) & 2) - 1;
}



/* ===== FUN_10013820 ===== */
/* Entry: 10013820 */

undefined4 __cdecl FUN_10013820(short *param_1,undefined4 param_2)

{
  int iVar1;
  size_t sVar2;
  size_t sVar3;
  int *piVar4;
  char *pcVar5;
  undefined2 extraout_var;
  undefined4 uVar6;
  undefined4 *puVar7;
  undefined2 extraout_var_00;
  int iVar8;
  int iVar9;
  
  piVar4 = FUN_1001efe0(param_2,&DAT_1006b1a8);
  if (piVar4 == (int *)0x0) {
    return 0xffff;
  }
  sVar3 = piVar4[4];
  pcVar5 = (char *)FUN_10016490(sVar3);
  if (pcVar5 == (char *)0x0) {
    return 0;
  }
  FUN_1001f5a0(piVar4,0,0,pcVar5,1,sVar3);
  FUN_1001f550(piVar4);
  FUN_10001480((undefined4 *)param_1,(undefined4 *)pcVar5,2,1);
  iVar8 = (int)*param_1;
  iVar9 = (iVar8 / 2) * (iVar8 + 1);
  sVar2 = iVar9 * 4;
  if ((int)sVar3 < (int)(sVar2 + 2)) {
    FUN_10016500(pcVar5);
    return CONCAT22(extraout_var,0xffff);
  }
  uVar6 = FUN_10016490(iVar8 << 2);
  *(undefined4 *)(param_1 + 2) = uVar6;
  puVar7 = (undefined4 *)FUN_10016490(sVar2);
  *(undefined4 **)(param_1 + 4) = puVar7;
  FUN_10001480(puVar7,(undefined4 *)(pcVar5 + 2),4,iVar9);
  iVar9 = 0;
  iVar8 = 0;
  if (0 < *param_1) {
    do {
      iVar1 = iVar8 * 4;
      iVar8 = iVar8 + 1 + iVar9;
      *(int *)(*(int *)(param_1 + 2) + iVar9 * 4) = *(int *)(param_1 + 4) + iVar1;
      iVar9 = iVar9 + 1;
    } while (iVar9 < *param_1);
  }
  FUN_10016500(pcVar5);
  return CONCAT22(extraout_var_00,1);
}



/* ===== FUN_10013970 ===== */
/* Entry: 10013970 */

undefined4 __cdecl FUN_10013970(int param_1)

{
  int iVar1;
  int iVar2;
  undefined4 uVar3;
  int iVar4;
  int iVar5;
  float10 fVar6;
  undefined4 local_4;
  
  iVar2 = param_1;
  *(undefined2 *)(param_1 + 0x14) = 0x100;
  uVar3 = FUN_10016490(0x400);
  *(undefined4 *)(param_1 + 0x18) = uVar3;
  uVar3 = FUN_10016490(0x20200);
  *(undefined4 *)(param_1 + 0x1c) = uVar3;
  iVar4 = 0;
  iVar5 = 0;
  if (0 < *(short *)(param_1 + 0x14)) {
    do {
      iVar1 = iVar5 * 4;
      iVar5 = iVar5 + 1 + iVar4;
      *(int *)(*(int *)(param_1 + 0x18) + iVar4 * 4) = *(int *)(param_1 + 0x1c) + iVar1;
      iVar4 = iVar4 + 1;
    } while (iVar4 < *(short *)(param_1 + 0x14));
  }
  local_4 = 0;
  if (0 < *(short *)(param_1 + 0x14)) {
    iVar5 = 0;
    do {
      param_1 = 0;
      if (-1 < iVar5) {
        do {
          fVar6 = FUN_10013920((float)local_4,(float)param_1);
          param_1 = param_1 + 1;
          *(float *)(*(int *)(iVar5 + *(int *)(iVar2 + 0x18)) + -4 + param_1 * 4) = (float)fVar6;
        } while (param_1 <= local_4);
      }
      local_4 = local_4 + 1;
      iVar5 = iVar5 + 4;
    } while (local_4 < *(short *)(iVar2 + 0x14));
    return 1;
  }
  return 1;
}



/* ===== FUN_100129f0 ===== */
/* Entry: 100129f0 */

undefined4 __cdecl FUN_100129f0(int param_1,undefined4 param_2)

{
  char cVar1;
  byte bVar2;
  short sVar3;
  int *piVar4;
  char *pcVar5;
  byte *pbVar6;
  char *pcVar7;
  byte *pbVar8;
  int iVar9;
  short *psVar10;
  int iVar11;
  int iVar12;
  undefined4 uVar13;
  undefined4 uVar14;
  void *this;
  uint uVar15;
  uint uVar16;
  short *psVar17;
  char *pcVar18;
  byte *pbVar19;
  char *pcVar20;
  char *pcVar21;
  bool bVar22;
  char local_239;
  byte *local_238;
  char *local_234;
  int local_230;
  int local_22c;
  int local_228;
  char *local_224;
  int *local_220;
  LPVOID local_21c;
  undefined4 local_218;
  undefined4 local_214;
  undefined4 local_210;
  int local_20c;
  int local_208;
  int *local_204;
  undefined1 local_200 [512];
  
  piVar4 = FUN_1001efe0(param_2,&DAT_1006b1a8);
  if (piVar4 == (int *)0x0) {
    return 0;
  }
  uVar15 = piVar4[4];
  pcVar5 = (char *)FUN_10016490(uVar15 + 1);
  FUN_1001f5a0(piVar4,0,0,pcVar5,1,uVar15);
  FUN_1001f550(piVar4);
  pcVar20 = pcVar5 + uVar15;
  *pcVar20 = '\0';
  local_224 = pcVar5;
  pbVar6 = (byte *)FUN_10012900(&local_224,pcVar20);
  if (pbVar6 == (byte *)0x0) {
    return 0;
  }
  sVar3 = FUN_10059429(this,pbVar6);
  *(short *)(param_1 + 0xe98) = sVar3;
  if (0x14 < sVar3) {
    return 0;
  }
  local_22c = 0;
  if (0 < sVar3) {
    local_234 = (char *)(param_1 + 0x369a);
    local_238 = (byte *)0x0;
    pcVar18 = (char *)(param_1 + 0xe9a);
    do {
      pcVar7 = FUN_10012900(&local_224,pcVar20);
      if (pcVar7 == (char *)0x0) {
        return 0;
      }
      uVar15 = 0xffffffff;
      do {
        pcVar21 = pcVar7;
        if (uVar15 == 0) break;
        uVar15 = uVar15 - 1;
        pcVar21 = pcVar7 + 1;
        cVar1 = *pcVar7;
        pcVar7 = pcVar21;
      } while (cVar1 != '\0');
      uVar15 = ~uVar15;
      pcVar7 = pcVar21 + -uVar15;
      pcVar21 = pcVar18;
      for (uVar16 = uVar15 >> 2; uVar16 != 0; uVar16 = uVar16 - 1) {
        *(undefined4 *)pcVar21 = *(undefined4 *)pcVar7;
        pcVar7 = pcVar7 + 4;
        pcVar21 = pcVar21 + 4;
      }
      for (uVar15 = uVar15 & 3; uVar15 != 0; uVar15 = uVar15 - 1) {
        *pcVar21 = *pcVar7;
        pcVar7 = pcVar7 + 1;
        pcVar21 = pcVar21 + 1;
      }
      pcVar7 = FUN_10012900(&local_224,pcVar20);
      if (pcVar7 == (char *)0x0) {
        return 0;
      }
      uVar15 = 0xffffffff;
      do {
        pcVar21 = pcVar7;
        if (uVar15 == 0) break;
        uVar15 = uVar15 - 1;
        pcVar21 = pcVar7 + 1;
        cVar1 = *pcVar7;
        pcVar7 = pcVar21;
      } while (cVar1 != '\0');
      uVar15 = ~uVar15;
      pcVar7 = pcVar21 + -uVar15;
      pcVar21 = local_234;
      for (uVar16 = uVar15 >> 2; uVar16 != 0; uVar16 = uVar16 - 1) {
        *(undefined4 *)pcVar21 = *(undefined4 *)pcVar7;
        pcVar7 = pcVar7 + 4;
        pcVar21 = pcVar21 + 4;
      }
      for (uVar15 = uVar15 & 3; uVar15 != 0; uVar15 = uVar15 - 1) {
        *pcVar21 = *pcVar7;
        pcVar7 = pcVar7 + 1;
        pcVar21 = pcVar21 + 1;
      }
      uVar15 = 0xffffffff;
      pcVar7 = pcVar18;
      do {
        if (uVar15 == 0) break;
        uVar15 = uVar15 - 1;
        cVar1 = *pcVar7;
        pcVar7 = pcVar7 + 1;
      } while (cVar1 != '\0');
      pbVar19 = &DAT_1006f108;
      pbVar6 = (byte *)((int)local_238 + ~uVar15 + 0xe95 + param_1);
      pbVar8 = pbVar6;
      do {
        bVar2 = *pbVar8;
        bVar22 = bVar2 < *pbVar19;
        if (bVar2 != *pbVar19) {
LAB_10012b74:
          iVar9 = (1 - (uint)bVar22) - (uint)(bVar22 != 0);
          goto LAB_10012b79;
        }
        if (bVar2 == 0) break;
        bVar2 = pbVar8[1];
        bVar22 = bVar2 < pbVar19[1];
        if (bVar2 != pbVar19[1]) goto LAB_10012b74;
        pbVar8 = pbVar8 + 2;
        pbVar19 = pbVar19 + 2;
      } while (bVar2 != 0);
      iVar9 = 0;
LAB_10012b79:
      if (iVar9 == 0) {
        *pbVar6 = 0;
      }
      local_234 = local_234 + 0x14;
      local_22c = local_22c + 1;
      local_238 = (byte *)((int)local_238 + 0x200);
      pcVar18 = pcVar18 + 0x200;
    } while (local_22c < *(short *)(param_1 + 0xe98));
  }
  FUN_10016500(pcVar5);
  local_214 = 1000000;
  local_21c = (LPVOID)FUN_10016490(1000000);
  psVar10 = (short *)FUN_10016490(*(short *)(param_1 + 0xe98) * 0x68);
  *(short **)(param_1 + 0x88) = psVar10;
  *(undefined4 *)(param_1 + 0x70) = 0;
  iVar9 = 0;
  if (0 < *(short *)(param_1 + 0xe98)) {
    do {
      FUN_100588f6(local_200,(byte *)s__s_s_idx_1006f0fc);
      piVar4 = FUN_1001efe0(local_200,&DAT_1006b1a8);
      if (piVar4 == (int *)0x0) {
        return 0;
      }
      iVar12 = piVar4[4];
      psVar17 = psVar10 + iVar9 * 0x34;
      iVar11 = FUN_10012960(piVar4,(int)psVar17);
      if (iVar11 == 0) {
        return 0;
      }
      local_210 = 0;
      local_218 = 0;
      local_20c = (int)psVar17[0x33];
      local_220 = piVar4;
      local_208 = iVar12;
      FUN_10012140(psVar17,&local_220,0);
      FUN_10012030((undefined4 *)(psVar17 + 0x28),&local_220,4,1);
      FUN_1001f550(piVar4);
      *(int *)(param_1 + 0x70) = *(int *)(param_1 + 0x70) + *(int *)(psVar17 + 0x28);
      if (iVar9 == 0) {
        psVar10[0x30] = 0;
        psVar10[0x31] = 0;
      }
      else {
        *(int *)(psVar17 + 0x30) = *(int *)(psVar17 + -0xc) + *(int *)(psVar17 + -4);
      }
      iVar9 = iVar9 + 1;
    } while (iVar9 < *(short *)(param_1 + 0xe98));
  }
  FUN_100124d0((undefined4 *)(param_1 + 0x30),*(size_t *)(param_1 + 0x70));
  local_228 = 0;
  local_239 = '\0';
  local_22c = 0;
  if (0 < *(short *)(param_1 + 0xe98)) {
    psVar17 = psVar10 + 0x2b;
    do {
      FUN_100588f6(local_200,(byte *)s__s_s_idx_1006f0fc);
      local_220 = FUN_1001efe0(local_200,&DAT_1006b1a8);
      if (local_220 == (int *)0x0) {
        return 0;
      }
      iVar9 = local_220[4];
      local_210 = 0;
      local_218 = 0;
      sVar3 = psVar17[8];
      local_20c = (int)sVar3;
      local_208 = iVar9;
      local_204 = local_220;
      iVar12 = FUN_10012140(psVar17 + -0x2b,&local_220,1);
      piVar4 = (int *)(psVar17 + -3);
      iVar11 = FUN_10012030(piVar4,&local_220,4,1);
      sVar3 = sVar3 + (short)iVar12 + (short)iVar11;
      if ((char)psVar17[7] == '\0') {
        psVar17[-1] = sVar3;
        *psVar17 = (short)((iVar9 - sVar3) / *piVar4);
      }
      else {
        iVar9 = FUN_10012030((undefined4 *)psVar17,&local_220,2,1);
        psVar17[-1] = (short)iVar9 + sVar3;
      }
      FUN_10012580((int *)(psVar17 + -0x23),(int *)(param_1 + 0x30),local_228);
      iVar9 = FUN_10012450((int)(psVar17 + -0x2b),&local_220);
      if (iVar9 == 0) {
        return 0;
      }
      iVar9 = 0;
      *(char *)(psVar17 + -0x25) = local_239;
      if (0 < *piVar4) {
        do {
          *(char *)(*(int *)(psVar17 + -5) + iVar9) =
               *(char *)(*(int *)(psVar17 + -5) + iVar9) + local_239;
          iVar9 = iVar9 + 1;
        } while (iVar9 < *piVar4);
      }
      local_228 = local_228 + *piVar4;
      local_239 = local_239 + (char)psVar17[-0x2b];
      *(int **)(psVar17 + 3) = local_204;
      local_22c = local_22c + 1;
      psVar17 = psVar17 + 0x34;
    } while (local_22c < *(short *)(param_1 + 0xe98));
  }
  *(undefined2 *)(param_1 + 0x20) = 0;
  iVar9 = 0;
  psVar17 = psVar10;
  if (0 < *(short *)(param_1 + 0xe98)) {
    do {
      *(short *)(param_1 + 0x20) = *(short *)(param_1 + 0x20) + *psVar17;
      iVar9 = iVar9 + 1;
      psVar17 = psVar17 + 0x34;
    } while (iVar9 < *(short *)(param_1 + 0xe98));
  }
  uVar13 = FUN_10016490((int)*(short *)(param_1 + 0x20) << 2);
  *(undefined4 *)(param_1 + 0x24) = uVar13;
  uVar13 = FUN_10016490((int)*(short *)(param_1 + 0x20));
  iVar9 = 0;
  *(undefined4 *)(param_1 + 0x28) = uVar13;
  local_22c = 0;
  if (0 < *(short *)(param_1 + 0xe98)) {
    local_238 = (byte *)(param_1 + 0x369a);
    psVar17 = psVar10;
    do {
      uVar13 = FUN_10012730(local_238);
      local_230 = 0;
      if (0 < *psVar17) {
        do {
          uVar15 = 0xffffffff;
          pcVar5 = *(char **)(*(int *)(psVar17 + 2) + local_230 * 4);
          do {
            if (uVar15 == 0) break;
            uVar15 = uVar15 - 1;
            cVar1 = *pcVar5;
            pcVar5 = pcVar5 + 1;
          } while (cVar1 != '\0');
          uVar14 = FUN_10016490(~uVar15);
          *(undefined4 *)(*(int *)(param_1 + 0x24) + iVar9 * 4) = uVar14;
          FUN_10058950(*(undefined4 **)(*(int *)(param_1 + 0x24) + iVar9 * 4),
                       *(undefined4 **)(*(int *)(psVar17 + 2) + local_230 * 4),~uVar15);
          *(char *)(iVar9 + *(int *)(param_1 + 0x28)) = (char)uVar13;
          iVar9 = iVar9 + 1;
          local_230 = local_230 + 1;
        } while (local_230 < *psVar17);
      }
      FUN_10012480(psVar17);
      *(int *)(psVar17 + 2) = *(int *)(param_1 + 0x24) + (iVar9 - *psVar17) * 4;
      *(int *)(psVar17 + 4) = (iVar9 - *psVar17) + *(int *)(param_1 + 0x28);
      local_22c = local_22c + 1;
      local_238 = local_238 + 0x14;
      psVar17 = psVar17 + 0x34;
    } while (local_22c < *(short *)(param_1 + 0xe98));
  }
  FUN_10016500(local_21c);
  iVar9 = 1;
  if (1 < *(short *)(param_1 + 0xe98)) {
    psVar17 = psVar10 + 0x66;
    do {
      if ((char)*psVar17 != (char)psVar10[0x32]) {
        return 0;
      }
      iVar9 = iVar9 + 1;
      psVar17 = psVar17 + 0x34;
    } while (iVar9 < *(short *)(param_1 + 0xe98));
  }
  sVar3 = psVar10[0x32];
  *(char *)(param_1 + 0x84) = (char)sVar3;
  if ((char)sVar3 == '\0') {
    *(undefined2 *)(param_1 + 0x78) = 6;
    *(undefined2 *)(param_1 + 0x7a) = 4;
  }
  sVar3 = psVar10[0x2b];
  *(short *)(param_1 + 0x76) = sVar3;
  if (sVar3 < 0xc9) {
    iVar9 = 0;
    if (0 < *(int *)(param_1 + 0x70)) {
      do {
        if (*(short *)(param_1 + 0x20) <= (short)(ushort)*(byte *)(iVar9 + *(int *)(param_1 + 0x6c))
           ) {
          return 0;
        }
        iVar9 = iVar9 + 1;
      } while (iVar9 < *(int *)(param_1 + 0x70));
    }
    return 1;
  }
  return 0;
}



/* ===== FUN_100137a0 ===== */
/* Entry: 100137a0 */

bool __cdecl FUN_100137a0(int param_1,int param_2)

{
  int iVar1;
  
  iVar1 = FUN_10013380(param_1,param_2);
  if (iVar1 == 0) {
    iVar1 = FUN_10013030((int *)(param_1 + 0x8c),param_1 + 0x20);
    if (iVar1 == 0) {
      return false;
    }
  }
  iVar1 = FUN_10013320((int *)(param_1 + 0x8c),param_1 + 0x20);
  return iVar1 != 0;
}



/* ===== FUN_10013760 ===== */
/* Entry: 10013760 */

undefined4 __cdecl FUN_10013760(int param_1,int param_2)

{
  int iVar1;
  uint uVar2;
  
  iVar1 = FUN_10013570(param_1,param_2);
  if (iVar1 == 0) {
    uVar2 = FUN_10024850(param_1 + 0xa0,param_1 + 0x8c);
    if ((short)uVar2 == 0) {
      return 0;
    }
  }
  return 1;
}



/* ===== FUN_1000eb20 ===== */
/* Entry: 1000eb20 */

undefined4 __cdecl FUN_1000eb20(int param_1,undefined2 *param_2)

{
  int iVar1;
  
  iVar1 = FUN_10001000(param_1 + 0x498,param_1);
  if (iVar1 < 0) {
    *param_2 = 8;
    return CONCAT22((short)((uint)param_2 >> 0x10),0xffff);
  }
  return CONCAT22((short)((uint)iVar1 >> 0x10),1);
}



/* ===== FUN_1000ec70 ===== */
/* Entry: 1000ec70 */

undefined4 __cdecl FUN_1000ec70(int param_1,undefined2 *param_2)

{
  undefined4 uVar1;
  int iVar2;
  undefined2 extraout_var;
  undefined1 local_200 [512];
  
  FUN_100588f6(local_200,(byte *)s__sdist_tbl_cepdist_tbl_1006cf74);
  uVar1 = FUN_10013820((short *)(param_1 + 8),local_200);
  if ((short)uVar1 < 0) {
    *param_2 = 6;
    return CONCAT22((short)((uint)uVar1 >> 0x10),0xffff);
  }
  FUN_1000e890();
  iVar2 = FUN_10013970(param_1);
  if (iVar2 < 0) {
    *param_2 = 6;
    return CONCAT22((short)((uint)iVar2 >> 0x10),0xffff);
  }
  FUN_1000e890();
  return CONCAT22(extraout_var,1);
}



/* ===== FUN_1000ed30 ===== */
/* Entry: 1000ed30 */

undefined4 __cdecl FUN_1000ed30(int param_1,undefined2 *param_2)

{
  bool bVar1;
  int iVar2;
  undefined3 extraout_var;
  undefined2 extraout_var_00;
  undefined1 local_200 [512];
  
  FUN_100588f6(local_200,(byte *)s__sdblist_idx_1006d010);
  iVar2 = FUN_100129f0(param_1,local_200);
  if (iVar2 == 0) {
    *param_2 = 7;
    return 0xffff;
  }
  FUN_1000e890();
  FUN_100588f6(local_200,(byte *)s__sclass_idx_1006cfe0);
  bVar1 = FUN_100137a0(param_1,(int)local_200);
  if (CONCAT31(extraout_var,bVar1) == 0) {
    *param_2 = 7;
    return CONCAT22((short)((uint3)extraout_var >> 8),0xffff);
  }
  FUN_1000e890();
  FUN_100588f6(local_200,(byte *)s__sclasshp_idx_1006cfb0);
  iVar2 = FUN_10013760(param_1,(int)local_200);
  if (iVar2 == 0) {
    *param_2 = 7;
    return 0xffff;
  }
  FUN_1000e890();
  return CONCAT22(extraout_var_00,1);
}



/* ===== FUN_1000ee40 ===== */
/* Entry: 1000ee40 */

undefined4 __cdecl FUN_1000ee40(int *param_1,undefined2 *param_2)

{
  int iVar1;
  int *piVar2;
  int iVar3;
  undefined1 local_200 [512];
  
  iVar1 = FUN_10016490((int)(short)param_1[8] << 2);
  *param_1 = iVar1;
  piVar2 = (int *)0x0;
  if (iVar1 != 0) {
    iVar1 = 0;
    if (0 < (short)param_1[8]) {
      do {
        FUN_100588f6(local_200,(byte *)s__s_s_dat_1006d02c);
        piVar2 = FUN_1001efe0(local_200,&DAT_1006b1a8);
        *(int **)(*param_1 + iVar1 * 4) = piVar2;
        if (*(int *)(*param_1 + iVar1 * 4) == 0) goto LAB_1000ee66;
        iVar1 = iVar1 + 1;
      } while (iVar1 < (short)param_1[8]);
    }
    iVar1 = FUN_10016490((int)(short)param_1[8] << 2);
    param_1[1] = iVar1;
    if (iVar1 == 0) {
      *param_2 = 10;
      return CONCAT22((short)((uint)param_2 >> 0x10),0xffff);
    }
    iVar3 = 0;
    if (0 < (short)param_1[8]) {
      do {
        FUN_100588f6(local_200,(byte *)s__s_s_upm_1006d020);
        piVar2 = FUN_1001efe0(local_200,&DAT_1006b1a8);
        *(int **)(param_1[1] + iVar3 * 4) = piVar2;
        iVar1 = param_1[1];
        if (*(int *)(iVar1 + iVar3 * 4) == 0) {
          *param_2 = 10;
          return CONCAT22((short)((uint)iVar1 >> 0x10),0xffff);
        }
        iVar3 = iVar3 + 1;
      } while (iVar3 < (short)param_1[8]);
    }
    return CONCAT22((short)((uint)iVar1 >> 0x10),1);
  }
LAB_1000ee66:
  *param_2 = 9;
  return CONCAT22((short)((uint)piVar2 >> 0x10),0xffff);
}



/* ===== FUN_1000f3e0 ===== */
/* Entry: 1000f3e0 */

uint __cdecl FUN_1000f3e0(undefined4 param_1,int param_2)

{
  char cVar1;
  short *psVar2;
  int iVar3;
  char *pcVar4;
  undefined4 *puVar5;
  uint uVar6;
  int iVar7;
  short *psVar8;
  undefined1 *puVar9;
  uint uVar10;
  int *piVar11;
  
  iVar7 = 0;
  psVar2 = *(short **)(param_2 + 0x54);
  if (psVar2[1] != 1 && -1 < psVar2[1] + -1) {
    pcVar4 = (char *)((int)psVar2 + 0x945);
    do {
      iVar7 = iVar7 + 1;
      *pcVar4 = (pcVar4[-3] != '\x01') + '[';
      pcVar4 = pcVar4 + 0x2fc;
    } while (iVar7 < psVar2[1] + -1);
  }
  if ((char)psVar2[0x29436] == '\0') {
    *(undefined1 *)((int)psVar2 + psVar2[1] * 0x2fc + 0x649) = 0x5a;
  }
  else {
    *(undefined1 *)((int)psVar2 + psVar2[1] * 0x2fc + 0x649) = 0x5e;
  }
  iVar7 = 0;
  if (0 < psVar2[1]) {
    puVar5 = (undefined4 *)(param_2 + 0x11d940);
    do {
      if (-1 < (int)puVar5[-200]) {
        *puVar5 = 2;
      }
      iVar7 = iVar7 + 1;
      puVar5 = puVar5 + 1;
    } while (iVar7 < psVar2[1]);
  }
  iVar7 = 1;
  if (1 < psVar2[1]) {
    puVar9 = (undefined1 *)((int)psVar2 + 0x945);
    piVar11 = (int *)(param_2 + 0x11d944);
    do {
      if (*piVar11 == 2) {
        *puVar9 = 0x5b;
      }
      iVar7 = iVar7 + 1;
      piVar11 = piVar11 + 1;
      puVar9 = puVar9 + 0x2fc;
    } while (iVar7 < psVar2[1]);
  }
  uVar6 = 0;
  iVar7 = 0;
  uVar10 = 0;
  if (0 < psVar2[1]) {
    pcVar4 = (char *)((int)psVar2 + 0x945);
    psVar8 = psVar2 + 0x38f;
    do {
      iVar7 = iVar7 + *psVar8 * 2;
      if (1000 < iVar7) {
        if (uVar6 == uVar10) {
          return uVar6 & 0xffff0000;
        }
        uVar6 = uVar6 - 1;
        psVar8 = psVar8 + -0x17e;
        pcVar4 = pcVar4 + -0x2fc;
        *pcVar4 = '[';
      }
      cVar1 = *pcVar4;
      if (((cVar1 == '[') || (cVar1 == 'Z')) || (cVar1 == '^')) {
        iVar7 = 0;
        uVar10 = uVar6 + 1;
      }
      uVar6 = uVar6 + 1;
      psVar8 = psVar8 + 0x17e;
      pcVar4 = pcVar4 + 0x2fc;
    } while ((int)uVar6 < (int)psVar2[1]);
  }
  FUN_1000f1c0((int)psVar2);
  FUN_1000f2f0(psVar2);
  iVar7 = 0;
  if (0 < psVar2[1]) {
    psVar8 = psVar2 + 0x4a2;
    do {
      *(undefined1 *)psVar8 = 0;
      iVar7 = iVar7 + 1;
      psVar8 = psVar8 + 0x17e;
    } while (iVar7 < psVar2[1]);
  }
  iVar7 = 0;
  puVar9 = (undefined1 *)(psVar2[1] + -1);
  if (0 < (int)puVar9) {
    puVar9 = (undefined1 *)((int)psVar2 + 0x945);
    piVar11 = (int *)(param_2 + 0x11d944);
    do {
      iVar3 = *piVar11;
      if (iVar3 == 0) {
        *puVar9 = 0x5d;
LAB_1000f578:
        puVar9[-1] = 1;
      }
      else {
        if (iVar3 == 1) {
          *puVar9 = 0x5c;
          goto LAB_1000f578;
        }
        if (iVar3 == 2) {
          *puVar9 = 0x5b;
        }
      }
      iVar7 = iVar7 + 1;
      piVar11 = piVar11 + 1;
      puVar9 = puVar9 + 0x2fc;
    } while (iVar7 < psVar2[1] + -1);
  }
  return CONCAT22((short)((uint)puVar9 >> 0x10),1);
}



/* ===== FUN_1000f5a0 ===== */
/* Entry: 1000f5a0 */

undefined4 __cdecl FUN_1000f5a0(int param_1,int param_2)

{
  int iVar1;
  int iVar2;
  int iVar3;
  short sVar4;
  undefined2 uVar5;
  int iVar6;
  int iVar7;
  short *psVar8;
  int local_20;
  short *local_1c;
  short local_14 [10];
  
  iVar3 = param_1;
  iVar1 = *(int *)(param_2 + 0x54);
  local_20 = 0;
  uVar5 = (undefined2)((uint)iVar1 >> 0x10);
  if (*(short *)(iVar1 + 2) < 1) {
    return CONCAT22(uVar5,1);
  }
  psVar8 = (short *)(iVar1 + 0x71c);
  do {
    iVar6 = 0;
    param_2 = 0;
    param_1 = 0;
    if (0 < *psVar8) {
      do {
        iVar7 = 0;
        if (0 < *(short *)(*(int *)(psVar8 + -0x68) + 0x24 + iVar6)) {
          local_1c = psVar8 + param_2 + 0xf;
          do {
            FUN_1000f7f0(local_14,iVar1,local_20,param_1,iVar7);
            iVar2 = (char)(&DAT_1006fff0)
                          [*(char *)(*(int *)(*(int *)(psVar8 + -0x68) + 0x1c + iVar6) + iVar7)] *
                    0xb;
            if ((&DAT_1006d0f3)[iVar2] == '\x01') {
              switch((&DAT_1006d0f4)[iVar2]) {
              case 1:
                sVar4 = FUN_10001570(*(int **)(iVar3 + 0xbc),(int)local_14);
                break;
              case 2:
                sVar4 = FUN_10001570(*(int **)(iVar3 + 0xd8),(int)local_14);
                break;
              case 3:
                sVar4 = FUN_10001570(*(int **)(iVar3 + 0xf4),(int)local_14);
                break;
              case 4:
                sVar4 = FUN_10001570(*(int **)(iVar3 + 0x110),(int)local_14);
                break;
              default:
                sVar4 = 0x640;
              }
            }
            else {
              switch((&DAT_1006d0f8)[iVar2]) {
              case 1:
                sVar4 = FUN_10001570(*(int **)(iVar3 + 300),(int)local_14);
                break;
              case 2:
                sVar4 = FUN_10001570(*(int **)(iVar3 + 0x148),(int)local_14);
                break;
              case 3:
                sVar4 = FUN_10001570(*(int **)(iVar3 + 0x164),(int)local_14);
                break;
              case 4:
              case 5:
                sVar4 = FUN_10001570(*(int **)(iVar3 + 0x180),(int)local_14);
                break;
              case 6:
                sVar4 = FUN_10001570(*(int **)(iVar3 + 0x19c),(int)local_14);
                break;
              default:
                sVar4 = 500;
              }
            }
            *local_1c = sVar4;
            param_2 = param_2 + 1;
            local_1c = local_1c + 1;
            iVar7 = iVar7 + 1;
          } while (iVar7 < *(short *)(*(int *)(psVar8 + -0x68) + 0x24 + iVar6));
        }
        iVar6 = iVar6 + 0x28;
        param_1 = param_1 + 1;
      } while (param_1 < *psVar8);
    }
    psVar8 = psVar8 + 0x17e;
    local_20 = local_20 + 1;
  } while (local_20 < *(short *)(iVar1 + 2));
  return CONCAT22(uVar5,1);
}



/* ===== FUN_1000fd30 ===== */
/* Entry: 1000fd30 */

undefined4 __cdecl FUN_1000fd30(int param_1,int param_2)

{
  char cVar1;
  short *psVar2;
  undefined2 uVar3;
  int in_EAX;
  int *piVar4;
  undefined4 uVar5;
  int iVar6;
  int iVar7;
  short *psVar8;
  int *piVar9;
  int local_20;
  short local_18 [11];
  short local_2;
  
  psVar2 = *(short **)(param_2 + 0x54);
  local_20 = 0;
  if (0 < psVar2[1]) {
    psVar8 = psVar2 + 0x38e;
    do {
      iVar6 = 0;
      if (0 < *psVar8) {
        iVar7 = 0;
        do {
          FUN_1000fe90(local_18,(int)psVar2,local_20,iVar6);
          cVar1 = *(char *)((int)psVar8 + 0x229);
          if ((cVar1 == '^') && (iVar6 == *psVar8 + -1)) {
            piVar4 = *(int **)(param_1 + 0x20c);
            piVar9 = *(int **)(param_1 + 0x27c);
          }
          else if ((cVar1 == 'Z') && (iVar6 == *psVar8 + -1)) {
            piVar4 = *(int **)(param_1 + 0x1f0);
            piVar9 = *(int **)(param_1 + 0x260);
          }
          else if ((cVar1 == '[') && (iVar6 == *psVar8 + -1)) {
            piVar4 = *(int **)(param_1 + 0x1d4);
            piVar9 = *(int **)(param_1 + 0x244);
          }
          else {
            piVar4 = *(int **)(param_1 + 0x1b8);
            piVar9 = *(int **)(param_1 + 0x228);
          }
          uVar3 = FUN_10001570(piVar4,(int)local_18);
          *(char *)(*(int *)(psVar8 + -0x68) + 3 + iVar7) = (char)uVar3;
          local_2 = (short)*(char *)(*(int *)(psVar8 + -0x68) + 3 + iVar7);
          FUN_100019b0(piVar9,(int)local_18,(undefined4 *)(*(int *)(psVar8 + -0x68) + 4 + iVar7));
          iVar6 = iVar6 + 1;
          iVar7 = iVar7 + 0x28;
        } while (iVar6 < *psVar8);
      }
      psVar8 = psVar8 + 0x17e;
      in_EAX = local_20 + 1;
      local_20 = in_EAX;
    } while (in_EAX < psVar2[1]);
  }
  uVar3 = (undefined2)((uint)in_EAX >> 0x10);
  iVar6 = 0;
  if (0 < *psVar2) {
    do {
      uVar5 = FUN_1000fa70(iVar6,param_2);
      uVar3 = (undefined2)((uint)uVar5 >> 0x10);
      iVar6 = iVar6 + 1;
    } while (iVar6 < *psVar2);
  }
  return CONCAT22(uVar3,1);
}



/* ===== FUN_1001b800 ===== */
/* Entry: 1001b800 */

undefined4 __cdecl FUN_1001b800(int param_1)

{
  int iVar1;
  uint uVar2;
  char cVar3;
  short sVar4;
  int iVar5;
  undefined4 *puVar6;
  short sVar7;
  int iVar8;
  byte *pbVar9;
  byte bVar10;
  undefined1 uVar11;
  short sVar12;
  char *pcVar13;
  int iVar14;
  undefined1 local_4a;
  char local_49;
  undefined1 local_44;
  undefined4 local_43 [16];
  
  iVar1 = *(int *)(param_1 + 0x54);
  FUN_1001b5d0(param_1);
  sVar12 = 0;
  if (0 < *(short *)(iVar1 + 2)) {
    do {
      sVar4 = sVar12 + 1;
      sVar7 = sVar12 + -1;
      if (sVar7 < 0) {
LAB_1001b86d:
        local_4a = 0x5a;
        local_49 = '\x03';
      }
      else {
        do {
          if ((*(char *)(iVar1 + 0x945 + sVar7 * 0x2fc) == '[') ||
             (0 < *(short *)(iVar1 + sVar7 * 0x2fc + 0x71e))) {
            if (sVar7 < 0) goto LAB_1001b86d;
            iVar5 = sVar7 * 0x2fc;
            cVar3 = *(char *)(iVar5 + 0x945 + iVar1);
            if (cVar3 == '[') {
              local_4a = 0x5a;
              local_49 = '\x03';
            }
            else {
              local_4a = *(undefined1 *)(*(short *)(iVar5 + 0x71e + iVar1) + iVar5 + 0x657 + iVar1);
              local_49 = (cVar3 != ']') + '\x01';
            }
            goto LAB_1001b8bf;
          }
          sVar7 = sVar7 + -1;
        } while (-1 < sVar7);
        local_4a = 0x5a;
        local_49 = '\x03';
      }
LAB_1001b8bf:
      iVar5 = (int)sVar4;
      cVar3 = *(char *)(iVar1 + 0x649 + iVar5 * 0x2fc);
      sVar7 = sVar4;
      while ((((cVar3 != 'Z' && (cVar3 != '^')) && (cVar3 != '[')) &&
             (*(short *)(iVar1 + iVar5 * 0x2fc + 0x71e) < 1))) {
        sVar7 = sVar7 + 1;
        iVar5 = (int)sVar7;
        cVar3 = *(char *)(iVar1 + 0x649 + iVar5 * 0x2fc);
      }
      cVar3 = *(char *)(iVar1 + 0x649 + sVar7 * 0x2fc);
      if (cVar3 == 'Z') {
        uVar11 = 0x5a;
        param_1._0_1_ = '\x04';
      }
      else if (cVar3 == '^') {
        uVar11 = 0x5e;
        param_1._0_1_ = '\x05';
      }
      else if (cVar3 == '[') {
        uVar11 = 0x5b;
        param_1._0_1_ = '\x03';
      }
      else {
        uVar11 = *(undefined1 *)(iVar1 + sVar7 * 0x2fc + 0x658);
        param_1._0_1_ = (cVar3 != ']') + '\x01';
      }
      local_44 = local_4a;
      iVar5 = sVar12 * 0x2fc;
      FUN_10058950(local_43,(undefined4 *)(iVar5 + 0x658 + iVar1),
                   (int)*(short *)(iVar5 + 0x71e + iVar1));
      *(undefined1 *)((int)local_43 + (int)*(short *)(iVar5 + 0x71e + iVar1)) = uVar11;
      sVar12 = 0;
      if (0 < *(short *)(iVar5 + 0x71e + iVar1)) {
        do {
          iVar14 = (int)sVar12;
          puVar6 = (undefined4 *)(iVar5 + (iVar14 + 0x18c) * 5 + iVar1);
          *(undefined2 *)puVar6 = 0;
          *(undefined1 *)((int)puVar6 + 2) = 0;
          FUN_10058950(puVar6,(undefined4 *)(&local_44 + iVar14),3);
          cVar3 = '\0';
          if (sVar12 == 0) {
            cVar3 = local_49 * '\n';
          }
          iVar8 = *(short *)(iVar5 + 0x71e + iVar1) + -1;
          if (iVar14 == iVar8) {
            cVar3 = cVar3 + (char)param_1;
          }
          if ((0 < sVar12) && (*(char *)(iVar5 + iVar14 + 0x698 + iVar1) == '1')) {
            cVar3 = cVar3 + '\n';
          }
          if ((iVar14 < iVar8) && (*(char *)(iVar5 + iVar14 + 0x699 + iVar1) == '1')) {
            cVar3 = cVar3 + '\x01';
          }
          sVar12 = sVar12 + 1;
          iVar14 = iVar1 + iVar5 + iVar14 * 5;
          *(char *)(iVar14 + 0x7bf) = cVar3;
          *(char *)(iVar14 + 0x7c0) = cVar3;
        } while (sVar12 < *(short *)(iVar5 + 0x71e + iVar1));
      }
      bVar10 = *(byte *)(iVar5 + 0x7bf + iVar1) / 10;
      if ((((2 < bVar10) && (sVar12 = *(short *)(iVar5 + 0x71e + iVar1), 1 < sVar12)) &&
          ((&DAT_100703b0)[*(byte *)(iVar5 + 0x7bd + iVar1)] == '\0')) && (sVar7 = 1, 1 < sVar12)) {
        do {
          iVar14 = iVar5 + sVar7 * 5 + iVar1;
          *(byte *)(iVar14 + 0x7c0) = *(byte *)(iVar14 + 0x7c0) % 10 + bVar10 * '\n';
          if ((&DAT_100703b0)[*(byte *)(iVar14 + 0x7bd)] != '\0') break;
          sVar7 = sVar7 + 1;
        } while (sVar7 < *(short *)(iVar5 + 0x71e + iVar1));
      }
      sVar12 = *(short *)(iVar5 + 0x71e + iVar1);
      sVar7 = sVar12 + -1;
      iVar14 = iVar5 + sVar7 * 5 + iVar1;
      bVar10 = *(byte *)(iVar14 + 0x7bf) % 10;
      if (((0 < sVar7) && (2 < bVar10)) && ((&DAT_100703b0)[*(byte *)(iVar14 + 0x7bd)] == '\0')) {
        sVar12 = sVar12 + -2;
        while ((-1 < sVar12 &&
               (iVar14 = iVar5 + sVar12 * 5, iVar8 = iVar14 + iVar1,
               *(byte *)(iVar8 + 0x7c0) = (*(byte *)(iVar14 + 0x7c0 + iVar1) / 10) * '\n' + bVar10,
               (&DAT_100703b0)[*(byte *)(iVar8 + 0x7bd)] == '\0'))) {
          sVar12 = sVar12 + -1;
        }
      }
      sVar12 = sVar4;
    } while (sVar4 < *(short *)(iVar1 + 2));
  }
  sVar12 = 0;
  if (0 < *(short *)(iVar1 + 2)) {
    do {
      sVar4 = 0;
      iVar5 = sVar12 * 0x2fc;
      if (0 < *(short *)(iVar5 + 0x71e + iVar1)) {
        do {
          iVar14 = (int)sVar4;
          sVar4 = sVar4 + 1;
          *(undefined1 *)(iVar14 + iVar5 + 0x901 + iVar1) = 0;
        } while (sVar4 < *(short *)(iVar5 + 0x71e + iVar1));
      }
      sVar12 = sVar12 + 1;
    } while (sVar12 < *(short *)(iVar1 + 2));
  }
  sVar12 = 0;
  if (*(short *)(iVar1 + 2) != 1 && -1 < *(short *)(iVar1 + 2) + -1) {
    iVar5 = 0;
    do {
      iVar5 = iVar5 * 0x2fc;
      if (*(char *)(iVar5 + 0x944 + iVar1) == '\x01') {
        sVar4 = *(short *)(iVar5 + 0x71e + iVar1) + -1;
        if (-1 < sVar4) {
          iVar14 = (int)sVar4;
          pcVar13 = (char *)(iVar14 + iVar5 + 0x901 + iVar1);
          param_1 = iVar14 + 1;
          pbVar9 = (byte *)(iVar5 + iVar14 * 5 + 0x7c0 + iVar1);
          do {
            if (((uint)*pbVar9 % 10 == 1) || ((uint)*pbVar9 % 10 == 2)) {
              *pcVar13 = *pcVar13 + '\x01';
            }
            pbVar9 = pbVar9 + -5;
            pcVar13 = pcVar13 + -1;
            param_1 = param_1 + -1;
          } while (param_1 != 0);
        }
        sVar4 = 0;
        if (0 < *(short *)(iVar5 + 0xa1a + iVar1)) {
          do {
            uVar2 = *(byte *)(iVar5 + sVar4 * 5 + 0xabc + iVar1) / 10;
            if ((uVar2 == 1) || (uVar2 == 2)) {
              iVar14 = sVar4 + iVar5;
              *(char *)(iVar14 + 0xbfd + iVar1) = *(char *)(iVar14 + 0xbfd + iVar1) + '\n';
            }
            sVar4 = sVar4 + 1;
          } while (sVar4 < *(short *)(iVar5 + 0xa1a + iVar1));
        }
      }
      sVar12 = sVar12 + 1;
      iVar5 = (int)sVar12;
    } while (iVar5 < *(short *)(iVar1 + 2) + -1);
  }
  return 1;
}



/* ===== FUN_10001570 ===== */
/* Entry: 10001570 */

undefined2 __cdecl FUN_10001570(int *param_1,int param_2)

{
  int iVar1;
  int *piVar2;
  
LAB_1000157a:
  while( true ) {
    piVar2 = param_1;
    if ((int *)*piVar2 == (int *)0x0) {
      return (short)piVar2[3];
    }
    if (*(char *)((int)piVar2 + 0x11) != 'D') goto LAB_100015ae;
    iVar1 = FUN_10001540(*(short *)(param_2 + (char)piVar2[4] * 2),(short *)piVar2[2],
                         (uint)*(byte *)((int)piVar2 + 0x12));
    if (iVar1 == 0) break;
    param_1 = (int *)*piVar2;
  }
  goto LAB_100015a9;
LAB_100015ae:
  param_1 = (int *)*piVar2;
  if (*(short *)((int)piVar2 + 0xe) < *(short *)(param_2 + (char)piVar2[4] * 2)) {
LAB_100015a9:
    param_1 = (int *)piVar2[1];
  }
  goto LAB_1000157a;
}



/* ===== FUN_100019b0 ===== */
/* Entry: 100019b0 */

undefined4 __cdecl FUN_100019b0(int *param_1,int param_2,undefined4 *param_3)

{
  int *piVar1;
  int iVar2;
  undefined4 *puVar3;
  int *piVar4;
  
  piVar1 = param_1;
LAB_100019bd:
  while( true ) {
    piVar4 = piVar1;
    if ((int *)*piVar4 == (int *)0x0) {
      puVar3 = FUN_10058950(param_3,(undefined4 *)piVar4[3],(uint)*(byte *)(param_1 + 5) << 1);
      return CONCAT22((short)((uint)puVar3 >> 0x10),1);
    }
    if (*(char *)((int)piVar4 + 0x13) != 'D') goto LAB_100019f1;
    iVar2 = FUN_10001540(*(short *)(param_2 + *(char *)((int)piVar4 + 0x12) * 2),(short *)piVar4[2],
                         (uint)*(byte *)((int)piVar4 + 0x15));
    if (iVar2 == 0) break;
    piVar1 = (int *)*piVar4;
  }
  goto LAB_100019ec;
LAB_100019f1:
  piVar1 = (int *)*piVar4;
  if ((short)piVar4[4] < *(short *)(param_2 + *(char *)((int)piVar4 + 0x12) * 2)) {
LAB_100019ec:
    piVar1 = (int *)piVar4[1];
  }
  goto LAB_100019bd;
}



/* ===== FUN_10022f60 ===== */
/* Entry: 10022f60 */

uint __cdecl FUN_10022f60(int *param_1,float param_2)

{
  char cVar1;
  undefined4 in_EAX;
  uint uVar2;
  
  cVar1 = *(char *)((int)param_2 + 0x20);
  uVar2 = CONCAT31((int3)((uint)in_EAX >> 8),cVar1);
  if (cVar1 != '\x01') {
    if (*(short *)((int)param_2 + 0x48) == 0) {
      if (cVar1 == '\x02') goto LAB_10022f97;
      uVar2 = FUN_1001e470(*(short *)(*(int *)((int)param_2 + 0x54) + 0x5285c),param_2,param_1);
    }
    if (*(char *)((int)param_2 + 0x20) != '\x03') {
      uVar2 = FUN_10022fd0(*(short *)(*(int *)((int)param_2 + 0x54) + 0x5285c),(int)param_2,param_1)
      ;
      goto LAB_10022fc2;
    }
  }
LAB_10022f97:
  if (60000 < *(int *)((int)param_2 + 0x38)) {
    *(undefined4 *)((int)param_2 + 0x38) = 0;
    return uVar2 & 0xffff0000;
  }
LAB_10022fc2:
  return uVar2 & 0xffff0000;
}



/* ===== FUN_1000b340 ===== */
/* Entry: 1000b340 */

undefined4 __cdecl FUN_1000b340(byte *param_1,int param_2,int param_3,short param_4)

{
  byte *pbVar1;
  short sVar2;
  short extraout_var;
  
  if (param_3 == 0) {
    param_1[0x429a2] = 0;
    param_1[0x429a3] = 0;
    return 0;
  }
  pbVar1 = param_1 + 0x429a2;
  param_1[0] = 0;
  param_1[1] = 0;
  pbVar1[0] = 0;
  pbVar1[1] = 0;
  sVar2 = FUN_1000a4e0((short *)param_1,param_2,param_3,param_4);
  if (extraout_var < 0) {
    return 0xffffffff;
  }
  if (0 < *(short *)param_1) {
    FUN_1000ba60(param_1);
  }
  if (param_4 != 0) {
    FUN_10006410((int *)pbVar1,param_2,param_3);
  }
  FUN_1000b960((short *)pbVar1);
  return CONCAT22(extraout_var,sVar2);
}



/* ===== FUN_10021f10 ===== */
/* Entry: 10021f10 */

undefined4 __cdecl
FUN_10021f10(undefined4 param_1,int param_2,int param_3,int param_4,int param_5,undefined4 param_6,
            int param_7)

{
  if (param_7 < 0) {
    *(undefined1 *)(param_2 + 0x11cc82) = 0;
  }
  else {
    *(char *)(param_2 + 0x11cc82) = (char)param_7;
  }
  *(int *)(param_2 + 0x11cc88) = param_3;
  *(int *)(param_2 + 0x11cc8c) = param_4;
  *(int *)(param_2 + 0x11cc90) = param_5;
  *(undefined1 *)(param_2 + 0x11ccaf) = 0;
  *(undefined1 *)(param_2 + 0x11ccb0) = 0;
  *(undefined1 *)(param_2 + 0x11ccb1) = 0;
  if (param_3 < 0) {
    *(undefined1 *)(param_2 + 0x11cc83) = 0;
  }
  else {
    *(undefined1 *)(param_2 + 0x11cc83) = 1;
  }
  if (param_4 < 0) {
    *(undefined1 *)(param_2 + 0x11cc84) = 0;
  }
  else {
    *(undefined1 *)(param_2 + 0x11cc84) = 1;
    if (param_4 == 0) {
      *(undefined4 *)(param_2 + 0x11cc8c) = 200;
    }
    else {
      *(int *)(param_2 + 0x11cc8c) = (int)(10000 / (longlong)param_4);
    }
  }
  if (param_5 < 0) {
    *(undefined1 *)(param_2 + 0x11cc85) = 0;
  }
  else {
    *(undefined1 *)(param_2 + 0x11cc85) = 1;
  }
  *(undefined4 *)(param_2 + 0x12cb14) = param_6;
  *(int *)(*(int *)(param_2 + 0x54) + 0x5287c) = param_2 + 0x11ccc0;
  *(int *)(*(int *)(param_2 + 0x54) + 0x52874) = param_2 + 0x11cfe0;
  *(int *)(*(int *)(param_2 + 0x54) + 0x52878) = param_2 + 0x11d300;
  *(int *)(*(int *)(param_2 + 0x54) + 0x52870) = param_2 + 0x11d620;
  return 1;
}


