/* Ghidra decompiler output: approximate C pseudocode, not original source. */
/* Program: vt_pau.dll */

/* ===== VT_LOADTTS_EXT_ENG ===== */
/* Entry: 10027af0 */

/* WARNING: Globals starting with '_' overlap smaller symbols at the same address */

short __cdecl
VT_LOADTTS_EXT_ENG(HWND param_1,int param_2,char *param_3,undefined4 param_4,undefined4 param_5,
                  char *param_6,char *param_7,uint param_8)

{
  char cVar1;
  byte bVar2;
  int *piVar3;
  undefined4 uVar4;
  uint uVar5;
  uint uVar6;
  int iVar7;
  int iVar8;
  byte *pbVar9;
  char *pcVar10;
  byte *pbVar11;
  undefined4 *puVar12;
  byte local_208 [514];
  short local_6;

                    /* 0x27af0  32  VT_LOADTTS_EXT_ENG */
  FUN_10025e60();
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  if (DAT_1007d6dc < 0) {
    DAT_1007d6dc = FUN_10025e60();
  }
  if ((param_2 < 0) || (iVar8 = param_2, 5 < param_2)) {
    iVar8 = 1;
  }
  if (DAT_1009fc50 == 0) {
    VT_SetDecimal0Pron_ENG();
  }
  if (param_3 == (char *)0x0) {
    FUN_10028620(&DAT_1009fa4c,iVar8);
  }
  else {
    iVar7 = -1;
    pcVar10 = param_3;
    do {
      if (iVar7 == 0) break;
      iVar7 = iVar7 + -1;
      cVar1 = *pcVar10;
      pcVar10 = pcVar10 + 1;
    } while (cVar1 != '\0');
    if (iVar7 == -2) {
      FUN_10028620(&DAT_1009fa4c,iVar8);
    }
    else {
      uVar5 = 0xffffffff;
      do {
        pcVar10 = param_3;
        if (uVar5 == 0) break;
        uVar5 = uVar5 - 1;
        pcVar10 = param_3 + 1;
        cVar1 = *param_3;
        param_3 = pcVar10;
      } while (cVar1 != '\0');
      uVar5 = ~uVar5;
      pbVar9 = (byte *)(pcVar10 + -uVar5);
      pbVar11 = local_208;
      for (uVar6 = uVar5 >> 2; uVar6 != 0; uVar6 = uVar6 - 1) {
        *(undefined4 *)pbVar11 = *(undefined4 *)pbVar9;
        pbVar9 = pbVar9 + 4;
        pbVar11 = pbVar11 + 4;
      }
      for (uVar5 = uVar5 & 3; uVar5 != 0; uVar5 = uVar5 - 1) {
        *pbVar11 = *pbVar9;
        pbVar9 = pbVar9 + 1;
        pbVar11 = pbVar11 + 1;
      }
      FUN_10027aa0((char *)local_208);
      if (DAT_100a045c == '\x01') {
        iVar7 = FUN_1001c2c0(&DAT_1009fa4c,local_208);
        if (iVar7 == 0) goto LAB_10027c40;
        if ((&DAT_100a7480)[param_2] == '\x01') {
          return 1;
        }
        uVar5 = 0xffffffff;
        pbVar9 = local_208;
        do {
          pbVar11 = pbVar9;
          if (uVar5 == 0) break;
          uVar5 = uVar5 - 1;
          pbVar11 = pbVar9 + 1;
          bVar2 = *pbVar9;
          pbVar9 = pbVar11;
        } while (bVar2 != 0);
        uVar5 = ~uVar5;
        pbVar9 = pbVar11 + -uVar5;
        pbVar11 = &DAT_1009fa4c;
        for (uVar6 = uVar5 >> 2; uVar6 != 0; uVar6 = uVar6 - 1) {
          *(undefined4 *)pbVar11 = *(undefined4 *)pbVar9;
          pbVar9 = pbVar9 + 4;
          pbVar11 = pbVar11 + 4;
        }
      }
      else {
        uVar5 = 0xffffffff;
        pbVar9 = local_208;
        do {
          pbVar11 = pbVar9;
          if (uVar5 == 0) break;
          uVar5 = uVar5 - 1;
          pbVar11 = pbVar9 + 1;
          bVar2 = *pbVar9;
          pbVar9 = pbVar11;
        } while (bVar2 != 0);
        uVar5 = ~uVar5;
        pbVar9 = pbVar11 + -uVar5;
        pbVar11 = &DAT_1009fa4c;
        for (uVar6 = uVar5 >> 2; uVar6 != 0; uVar6 = uVar6 - 1) {
          *(undefined4 *)pbVar11 = *(undefined4 *)pbVar9;
          pbVar9 = pbVar9 + 4;
          pbVar11 = pbVar11 + 4;
        }
      }
      for (uVar5 = uVar5 & 3; uVar5 != 0; uVar5 = uVar5 - 1) {
        *pbVar11 = *pbVar9;
        pbVar9 = pbVar9 + 1;
        pbVar11 = pbVar11 + 1;
      }
    }
  }
LAB_10027c40:
  FUN_10027aa0(&DAT_1009fa4c);
  if (DAT_100a045c != '\x01') {
    if (param_1 != (HWND)0x0) {
      PostMessageA(param_1,0x464,10,0);
    }
    FUN_100274e0(DAT_100a9d38);
    if (DAT_100a0458 == 0) {
      InitializeCriticalSection((LPCRITICAL_SECTION)&DAT_100a8400);
      DAT_100a0458 = 1;
    }
    if (DAT_1009fe54 == 0) {
      FUN_10024cc0();
      DAT_1009fe54 = 1;
    }
    VT_SetDecimal0Pron_ENG();
    FUN_1002a600();
    puVar12 = &DAT_100a647c;
    for (iVar7 = 0x400; iVar7 != 0; iVar7 = iVar7 + -1) {
      *puVar12 = 0;
      puVar12 = puVar12 + 1;
    }
    FUN_10025e70(param_1,&local_6);
    if (local_6 != 0) {
      if (DAT_100a0460 != (undefined *)0x0) {
        FUN_1001da30(DAT_100a0460);
      }
      DAT_100a0460 = (undefined *)0x0;
      return local_6;
    }
    puVar12 = &DAT_100a047c;
    for (iVar7 = 0x1800; iVar7 != 0; iVar7 = iVar7 + -1) {
      *puVar12 = 0;
      puVar12 = puVar12 + 1;
    }
    _DAT_100a7480 = 0;
    _DAT_100a7488 = 0;
    puVar12 = &DAT_100a0464;
    for (iVar7 = 6; iVar7 != 0; iVar7 = iVar7 + -1) {
      *puVar12 = 0;
      puVar12 = puVar12 + 1;
    }
    DAT_100a045c = '\x01';
    _DAT_100a748c = 0;
    _DAT_100a7484 = 0;
    puVar12 = &DAT_100d46e0;
    for (iVar7 = 0xa800; iVar7 != 0; iVar7 = iVar7 + -1) {
      *puVar12 = 0;
      puVar12 = puVar12 + 1;
    }
  }
  if ((&DAT_100a0464)[iVar8] == 0) {
    piVar3 = FUN_10025f10(param_1,&local_6,iVar8);
    (&DAT_100a0464)[iVar8] = piVar3;
    if (local_6 != 0) {
      return local_6;
    }
    piVar3[0x1342] = (int)DAT_100a0460;
    *(undefined4 *)((&DAT_100a0464)[iVar8] + 0x290) = 0;
    iVar7 = VT_CheckLicense_ENG(param_6,param_7,param_8,(&PTR_DAT_1007c6a8)[iVar8 * 6]);
    if (iVar7 == 0) {
      (&DAT_100a7488)[iVar8] = 1;
      uVar4 = FUN_10028c50(param_6,param_7,param_8);
      *(undefined4 *)((&DAT_100a0464)[iVar8] + 0x4d14) = uVar4;
    }
    else {
      iVar7 = (&DAT_100a0464)[iVar8];
      (&DAT_100a7488)[iVar8] = 0;
      *(undefined4 *)(iVar7 + 0x4d14) = 1;
    }
    VT_SetDecimal0Pron_ENG();
    iVar7 = *(int *)((&DAT_100a0464)[iVar8] + 0x4d14);
    if ((iVar7 < 1) || (0x400 < iVar7)) {
      *(undefined4 *)((&DAT_100a0464)[iVar8] + 0x4d14) = 1;
      (&DAT_100a7488)[iVar8] = 0;
    }
  }
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  return 0;
}



/* ===== VT_UNLOADTTS_EXT_ENG ===== */
/* Entry: 10027ea0 */

void __cdecl VT_UNLOADTTS_EXT_ENG(int param_1)

{
  int *piVar1;
  int iVar2;
  bool bVar3;

                    /* 0x27ea0  58  VT_UNLOADTTS_EXT_ENG */
  if ((param_1 < 0) || (5 < param_1)) {
    param_1 = 1;
  }
  VT_STOPTTS_ENG();
  DAT_100a9d30 = param_1;
  if ((&DAT_100a0464)[param_1] != 0) {
    FUN_10027700(param_1);
  }
  piVar1 = &DAT_100a0464;
  do {
    if (*piVar1 != 0) {
      return;
    }
    piVar1 = piVar1 + 1;
  } while ((int)piVar1 < 0x100a047c);
  iVar2 = 0;
  do {
    VT_UNLOAD_UserDict_ENG(iVar2);
    iVar2 = iVar2 + 1;
  } while (iVar2 < 0x400);
  FUN_10027fa0();
  DestroyWindow(DAT_100a8418);
  if (DAT_100a0458 == 1) {
    DeleteCriticalSection((LPCRITICAL_SECTION)&DAT_100a8400);
    DAT_100a0458 = 0;
  }
  if (DAT_1009fe54 == 1) {
    VT_SetDecimal0Pron_ENG();
    DAT_1009fe54 = 0;
  }
  DAT_1009fa4c = DAT_1009f948;
  bVar3 = DAT_1009fc50 == 1;
  (&DAT_100a7488)[param_1] = 0;
  if (bVar3) {
    VT_SetDecimal0Pron_ENG();
  }
  VT_SetDecimal0Pron_ENG();
  return;
}



/* ===== VT_CheckLicense_ENG ===== */
/* Entry: 10029b80 */

int __cdecl VT_CheckLicense_ENG(char *param_1,char *param_2,uint param_3,byte *param_4)

{
  char cVar1;
  bool bVar2;
  char *pcVar3;
  char *_Str;
  char *pcVar4;
  undefined3 extraout_var;
  undefined3 extraout_var_00;
  undefined3 extraout_var_01;
  undefined3 extraout_var_02;
  undefined3 extraout_var_03;
  undefined3 extraout_var_04;
  undefined3 extraout_var_05;
  uint uVar5;
  int iVar6;
  char *pcVar7;
  DWORD local_44;
  int local_40;
  int local_3c;
  int local_38;
  char *local_c;

                    /* 0x29b80  3  VT_CheckLicense_ENG */
  FUN_10028ca0(&local_44);
  FUN_10028db0(param_1,param_2,param_3,&local_44);
  if (local_44 == 0) {
    pcVar3 = (char *)FUN_1001c440(local_c);
    _Str = _strstr(local_c,s_<vw_verify_1007d974);
    iVar6 = 0;
    if (_Str != (char *)0x0) {
      pcVar4 = _strstr(_Str,&DAT_1007d970);
      if (pcVar4 == (char *)0x0) {
        iVar6 = 0;
      }
      else {
        uVar5 = 0xffffffff;
        pcVar7 = s_<vw_verify_1007d974;
        do {
          if (uVar5 == 0) break;
          uVar5 = uVar5 - 1;
          cVar1 = *pcVar7;
          pcVar7 = pcVar7 + 1;
        } while (cVar1 != '\0');
        FUN_10063f30((undefined4 *)pcVar3,(undefined4 *)(_Str + (~uVar5 - 1)),
                     (uint)(pcVar4 + (-(int)_Str - (~uVar5 - 1))));
        pcVar3[(int)(pcVar4 + (-(int)_Str - (~uVar5 - 1)))] = '\0';
        bVar2 = FUN_10029d30(pcVar3,(byte *)&DAT_1007da50,param_4);
        if (CONCAT31(extraout_var,bVar2) == 0) {
          iVar6 = -5;
        }
        else {
          bVar2 = FUN_10029d30(pcVar3,(byte *)&DAT_1007da48,param_4);
          if (CONCAT31(extraout_var_00,bVar2) == 0) {
            iVar6 = -6;
          }
          else {
            bVar2 = FUN_10029d30(pcVar3,(byte *)s_speaker_1007da40,param_4);
            if (CONCAT31(extraout_var_01,bVar2) == 0) {
              iVar6 = -7;
            }
            else {
              bVar2 = FUN_10029d30(pcVar3,(byte *)s_version_1007da38,param_4);
              if (CONCAT31(extraout_var_02,bVar2) == 0) {
                iVar6 = -8;
              }
              else {
                bVar2 = FUN_10029d30(pcVar3,(byte *)s_dbaccess_1007da2c,param_4);
                if (CONCAT31(extraout_var_03,bVar2) == 0) {
                  iVar6 = -9;
                }
                else {
                  bVar2 = FUN_10029d30(pcVar3,(byte *)s_sampling_1007da20,param_4);
                  if (CONCAT31(extraout_var_04,bVar2) == 0) {
                    iVar6 = -10;
                  }
                  else {
                    bVar2 = FUN_10029d30(pcVar3,(byte *)s_dbsize_1007da18,param_4);
                    iVar6 = (-(uint)(CONCAT31(extraout_var_05,bVar2) != 0) & 0xb) - 0xb;
                  }
                }
              }
            }
          }
        }
      }
    }
    FUN_1001da30(pcVar3);
  }
  else if (local_40 < 0) {
    iVar6 = -2;
  }
  else if (local_3c < 0) {
    iVar6 = -3;
  }
  else {
    iVar6 = ((-1 < local_38) - 1 & 0xfffffffd) - 1;
  }
  FUN_10028ce0((int)&local_44);
  return iVar6;
}



/* ===== VT_TextToFile_ENG ===== */
/* Entry: 1001da50 */

undefined4 __cdecl
VT_TextToFile_ENG(undefined4 param_1,byte *param_2,int param_3,int param_4,int param_5,int param_6,
                 int param_7,undefined4 param_8,int param_9,int param_10)

{
  undefined4 uVar1;

                    /* 0x1da50  53  VT_TextToFile_ENG */
  switch(param_1) {
  case 0:
    uVar1 = FUN_1001e290(param_2,param_3,param_4,param_5,param_6,param_7,param_8,param_9,param_10);
    return uVar1;
  case 1:
    uVar1 = FUN_1001e430(param_2,param_3,param_4,param_5,param_6,param_7,param_8,param_9,param_10);
    return uVar1;
  case 2:
    uVar1 = FUN_1001e630(param_2,param_3,param_4,param_5,param_6,param_7,param_8,param_9,param_10);
    return uVar1;
  case 3:
    uVar1 = FUN_1001ed40(param_2,param_3,param_4,param_5,param_6,param_7,param_8,param_9,param_10);
    return uVar1;
  case 4:
    uVar1 = FUN_1001e830(param_2,param_3,param_4,param_5,param_6,param_7,param_8,param_9,param_10);
    return uVar1;
  case 5:
    uVar1 = FUN_1001e870(param_2,param_3,param_4,param_5,param_6,param_7,param_8,param_9,param_10);
    return uVar1;
  default:
    return CONCAT22((short)((uint)param_1 >> 0x10),0xffff);
  case 7:
    uVar1 = FUN_1001e8b0(param_2,param_3,param_4,param_5,param_6,param_7,param_8,param_9,param_10);
    return uVar1;
  case 8:
    uVar1 = FUN_1001e8f0(param_2,param_3,param_4,param_5,param_6,param_7,param_8,param_9,param_10);
    return uVar1;
  case 9:
    uVar1 = FUN_1001fc40(param_2,param_3,param_4,param_5,param_6,param_7,param_8,param_9,param_10);
    return uVar1;
  }
}



/* ===== VT_TextToBuffer_ENG ===== */
/* Entry: 1001dc40 */

undefined4 __cdecl
VT_TextToBuffer_ENG(undefined4 param_1,byte *param_2,undefined4 *param_3,int *param_4,int param_5,
                   int param_6,int param_7,int param_8,int param_9,int param_10,undefined4 param_11,
                   int param_12,int param_13)

{
  undefined4 uVar1;

                    /* 0x1dc40  52  VT_TextToBuffer_ENG */
  switch(param_1) {
  case 0:
    uVar1 = FUN_100200c0(param_2,param_3,param_4,param_5,param_6,param_7,param_8,param_9,param_10,
                         param_11,param_12,param_13);
    return uVar1;
  case 1:
    uVar1 = FUN_10020460(param_2,param_3,param_4,param_5,param_6,param_7,param_8,param_9,param_10,
                         param_11,param_12,param_13);
    return uVar1;
  case 2:
    uVar1 = FUN_10020930(param_2,param_3,param_4,param_5,param_6,param_7,param_8,param_9,param_10,
                         param_11,param_12,param_13);
    return uVar1;
  case 3:
    uVar1 = FUN_10020e00(param_2,param_3,param_4,param_5,param_6,param_7,param_8,param_9,param_10,
                         param_11,param_12,param_13);
    return uVar1;
  default:
    return 0xffffffff;
  }
}



/* ===== FUN_1001e830 ===== */
/* Entry: 1001e830 */

void __cdecl
FUN_1001e830(byte *param_1,int param_2,int param_3,int param_4,int param_5,int param_6,
            undefined4 param_7,int param_8,int param_9)

{
  FUN_1001e930(param_1,param_2,1,param_3,param_4,param_5,param_6,param_7,param_8,param_9);
  return;
}



/* ===== FUN_1001e930 ===== */
/* Entry: 1001e930 */

undefined4 __cdecl
FUN_1001e930(byte *param_1,int param_2,undefined2 param_3,int param_4,int param_5,int param_6,
            int param_7,undefined4 param_8,int param_9,int param_10)

{
  int *piVar1;
  int iVar2;
  char *pcVar3;
  int iVar4;
  byte *pbVar5;
  undefined4 *puVar6;
  char *pcVar7;
  undefined4 uVar8;
  uint uVar9;

  iVar4 = param_4;
  if ((param_4 < 0) || (iVar2 = param_4, 5 < param_4)) {
    iVar2 = 1;
  }
  piVar1 = (int *)(&DAT_100a0464)[iVar2];
  if (piVar1 == (int *)0x0) {
    return 0xfffb;
  }
  if (param_1 == (byte *)0x0) {
    return 0xfffd;
  }
  if (*param_1 == 0) {
    return CONCAT22((short)((uint)param_1 >> 0x10),0xfffc);
  }
  pcVar3 = FUN_10024d20(param_2,&DAT_1007d60c);
  if (pcVar3 == (char *)0x0) {
    return 0xfffa;
  }
  if ((iVar4 < 0) || (5 < iVar4)) {
    iVar4 = 1;
  }
  pbVar5 = (byte *)FUN_10025fc0((undefined2 *)((int)&param_4 + 2),param_9,iVar4);
  if (param_4._2_2_ != 1) {
    return CONCAT22((short)((uint)pbVar5 >> 0x10),0xfffe);
  }
  FUN_100286c0(piVar1,(int)pbVar5,param_5,param_6,param_7,param_8,param_10);
  puVar6 = FUN_1001df10((byte *)0x0);
  *(undefined4 **)(pbVar5 + 0x2c) = puVar6;
  *(undefined2 *)(pbVar5 + 0x11b244) = param_3;
  FUN_1001f1f0((int)pbVar5);
  pcVar7 = (char *)FUN_1001c990(pbVar5,param_1,0);
  if (pcVar7 == (char *)0x0) {
    FUN_1001e070(*(undefined **)(pbVar5 + 0x2c));
    uVar8 = FUN_10027750(pbVar5);
    return CONCAT22((short)((uint)uVar8 >> 0x10),0xfffd);
  }
  FUN_10026ab0((int)piVar1,(int *)pbVar5,pcVar7);
  FUN_1001f600((int *)pbVar5,pcVar3);
  iVar4 = *(int *)(pbVar5 + 0x44);
  param_2 = 0;
  do {
    if ((iVar4 != 0) ||
       (iVar4 = FUN_10026870(piVar1,(int *)pbVar5,*(int *)(*(int *)(pbVar5 + 0x38) + 4)), iVar4 < 1)
       ) {
      FUN_100253d0(pcVar3);
      FUN_1001e070(*(undefined **)(pbVar5 + 0x2c));
      uVar8 = FUN_10027750(pbVar5);
      return CONCAT22((short)((uint)uVar8 >> 0x10),1);
    }
    FUN_1001e0c0((int)pbVar5,(int)param_1);
    FUN_1001f600((int *)pbVar5,pcVar3);
    switch(param_3) {
    case 6:
      iVar4 = 0;
      if (0 < *(int *)(pbVar5 + 0x30) / 2) {
        do {
          iVar4 = iVar4 + 1;
          *(undefined1 *)(*(int *)(*(int *)(pbVar5 + 0x3c) + 4) + -1 + iVar4) =
               *(undefined1 *)
                (*(short *)(*(int *)(*(int *)(pbVar5 + 0x38) + 4) + -2 + iVar4 * 2) + 0x1841c +
                piVar1[0x1342]);
        } while (iVar4 < *(int *)(pbVar5 + 0x30) / 2);
      }
      goto LAB_1001ec02;
    case 7:
      iVar4 = 0;
      if (0 < *(int *)(pbVar5 + 0x30) / 2) {
        do {
          iVar4 = iVar4 + 1;
          *(undefined1 *)(*(int *)(*(int *)(pbVar5 + 0x3c) + 4) + -1 + iVar4) =
               *(undefined1 *)
                (*(short *)(*(int *)(*(int *)(pbVar5 + 0x38) + 4) + -2 + iVar4 * 2) + 0x841c +
                piVar1[0x1342]);
        } while (iVar4 < *(int *)(pbVar5 + 0x30) / 2);
      }
LAB_1001ec02:
      iVar2 = *(int *)(pbVar5 + 0x30);
      uVar9 = 1;
      pcVar7 = *(char **)(*(int *)(pbVar5 + 0x3c) + 4);
LAB_1001ec4b:
      FUN_10025500((int)pcVar3,0,2,pcVar7,uVar9,iVar2 / 2);
      break;
    default:
      iVar2 = *(int *)(pbVar5 + 0x30);
      if (DAT_1007d6dc == 0) {
        iVar4 = *(int *)(pbVar5 + 0x38);
        uVar9 = 2;
        goto LAB_1001ec47;
      }
      FUN_100255a0((int)pcVar3,0,2,*(char **)(*(int *)(pbVar5 + 0x38) + 4),2,iVar2 / 2);
      break;
    case 0x11:
      FUN_1001c8b0(*(short **)(*(int *)(pbVar5 + 0x38) + 4),*(byte **)(*(int *)(pbVar5 + 0x3c) + 4),
                   *(int *)(pbVar5 + 0x30) / 2,param_2,(int)pbVar5);
      FUN_10025500((int)pcVar3,0,2,*(char **)(*(int *)(pbVar5 + 0x3c) + 4),1,
                   (int)(*(int *)(pbVar5 + 0x30) + (*(int *)(pbVar5 + 0x30) >> 0x1f & 3U)) >> 2);
      break;
    case 0x99:
      iVar4 = 0;
      if (0 < *(int *)(pbVar5 + 0x30) / 2) {
        do {
          iVar2 = (int)*(short *)(*(int *)(*(int *)(pbVar5 + 0x38) + 4) + iVar4 * 2);
          param_10 = CONCAT13((char)(iVar2 + (iVar2 >> 0x1f & 0xffU) >> 8) + -0x80,
                              (undefined3)param_10);
          FUN_10063f30((undefined4 *)(*(int *)(*(int *)(pbVar5 + 0x3c) + 4) + iVar4),
                       (undefined4 *)((int)&param_10 + 3),1);
          iVar4 = iVar4 + 1;
        } while (iVar4 < *(int *)(pbVar5 + 0x30) / 2);
      }
      iVar2 = *(int *)(pbVar5 + 0x30);
      iVar4 = *(int *)(pbVar5 + 0x3c);
      uVar9 = 1;
LAB_1001ec47:
      pcVar7 = *(char **)(iVar4 + 4);
      goto LAB_1001ec4b;
    }
    iVar4 = *(int *)(pbVar5 + 0x44);
    param_2 = param_2 + 1;
  } while( true );
}



/* ===== FUN_1001f1f0 ===== */
/* Entry: 1001f1f0 */

undefined4 __cdecl FUN_1001f1f0(int param_1)

{
  ushort uVar1;
  ushort uVar2;
  int iVar3;
  undefined2 *puVar4;
  int iVar5;
  undefined2 *puVar6;
  int iVar7;
  undefined2 *puVar8;
  int iVar9;
  short sVar10;
  uint uVar11;
  uint uVar12;
  uint local_8;

  iVar3 = param_1;
  uVar1 = *(ushort *)(param_1 + 0x11b244);
  uVar11 = (uint)uVar1;
  iVar9 = 0;
  uVar2 = 0;
  switch(uVar11) {
  default:
    param_1 = 0x10;
    break;
  case 2:
  case 6:
  case 7:
  case 0x11:
  case 0x31:
  case 0x99:
    param_1 = 8;
  }
  local_8 = 1;
  uVar12 = 0;
  switch(uVar11) {
  case 1:
  case 0x99:
    uVar11 = 1;
    iVar9 = param_1 + 7 >> 3;
    goto LAB_1001f27b;
  case 2:
    iVar9 = 0x80;
    param_1._0_2_ = 4;
    uVar2 = 0x20;
    local_8 = 0xf4;
    break;
  default:
    goto switchD_1001f259_caseD_3;
  case 6:
  case 7:
    iVar9 = 1;
    break;
  case 0x11:
    iVar9 = 0x100;
    param_1._0_2_ = 4;
    uVar2 = 2;
    local_8 = 0x1f9;
    break;
  case 0x31:
    iVar9 = 0x41;
    param_1._0_2_ = 0;
    uVar2 = 2;
    local_8 = 0x140;
  }
  uVar12 = uVar11;
  if ((uVar1 == 0x31) || (uVar1 != 1)) {
switchD_1001f259_caseD_3:
    iVar5 = uVar2 + 0x12;
    uVar11 = uVar12;
  }
  else {
LAB_1001f27b:
    iVar5 = 0x10;
  }
  sVar10 = (short)uVar11;
  iVar7 = iVar5 + 0x1c;
  if (sVar10 != 1) {
    iVar7 = iVar5 + 0x28;
  }
  *(int *)(iVar3 + 0xedda8) = iVar5;
  *(undefined1 *)(iVar3 + 0xedd98) = 0x52;
  *(undefined1 *)(iVar3 + 0xedd99) = 0x49;
  *(undefined1 *)(iVar3 + 0xedd9a) = 0x46;
  *(undefined1 *)(iVar3 + 0xedd9b) = 0x46;
  *(int *)(iVar3 + 0xedd9c) = iVar7;
  *(undefined1 *)(iVar3 + 0xedda0) = 0x57;
  *(undefined1 *)(iVar3 + 0xedda1) = 0x41;
  *(undefined1 *)(iVar3 + 0xedda2) = 0x56;
  *(undefined1 *)(iVar3 + 0xedda3) = 0x45;
  *(undefined1 *)(iVar3 + 0xedda4) = 0x66;
  *(undefined1 *)(iVar3 + 0xedda5) = 0x6d;
  *(undefined1 *)(iVar3 + 0xedda6) = 0x74;
  *(undefined1 *)(iVar3 + 0xedda7) = 0x20;
  *(short *)(iVar3 + 0xeddac) = sVar10;
  *(undefined2 *)(iVar3 + 0xeddae) = 1;
  *(undefined4 *)(iVar3 + 0xeddb0) = 16000;
  *(short *)(iVar3 + 0xeddb8) = (short)iVar9;
  *(uint *)(iVar3 + 0xeddb4) = (iVar9 * 16000 + (local_8 >> 1)) / local_8;
  *(undefined2 *)(iVar3 + 0xeddba) = (undefined2)param_1;
  if (sVar10 != 1) {
    *(ushort *)(iVar3 + 0xeddbc) = uVar2;
  }
  if (uVar11 == 2) {
    *(undefined2 *)(iVar3 + 0xeddc0) = 7;
    *(undefined2 *)(iVar3 + 0xeddbe) = (undefined2)local_8;
    puVar8 = (undefined2 *)(iVar3 + 0xeddc4);
    puVar6 = &DAT_1007d51a;
    do {
      puVar4 = puVar6 + 2;
      puVar8[-1] = puVar6[-1];
      *puVar8 = *puVar6;
      puVar8 = puVar8 + 2;
      puVar6 = puVar4;
    } while ((int)puVar4 < 0x1007d536);
  }
  else if ((uVar11 == 0x11) || (puVar4 = (undefined2 *)(uVar11 - 0x31), puVar4 == (undefined2 *)0x0)
          ) {
    puVar4 = (undefined2 *)0x0;
    *(undefined2 *)(iVar3 + 0xeddbe) = (undefined2)local_8;
  }
  if (sVar10 != 1) {
    *(undefined1 *)(iVar3 + 0xeddde) = 0x66;
    *(undefined1 *)(iVar3 + 0xedddf) = 0x61;
    *(undefined1 *)(iVar3 + 0xedde0) = 99;
    *(undefined1 *)(iVar3 + 0xedde1) = 0x74;
    *(undefined4 *)(iVar3 + 0xedde4) = 4;
    *(undefined4 *)(iVar3 + 0xedde8) = 0;
  }
  *(undefined1 *)(iVar3 + 0xedded) = 0x61;
  *(undefined1 *)(iVar3 + 0xeddef) = 0x61;
  *(undefined1 *)(iVar3 + 0xeddec) = 100;
  *(undefined1 *)(iVar3 + 0xeddee) = 0x74;
  *(undefined4 *)(iVar3 + 0xeddf0) = 0;
  return CONCAT22((short)((uint)puVar4 >> 0x10),1);
}



/* ===== FUN_1001f600 ===== */
/* Entry: 1001f600 */

undefined4 __cdecl FUN_1001f600(int *param_1,char *param_2)

{
  int *piVar1;
  short sVar2;
  int *piVar3;
  char *pcVar4;
  uint uVar5;
  int iVar6;
  int local_8;

  pcVar4 = param_2;
  piVar3 = param_1;
  switch((short)param_1[0x46c91]) {
  case 1:
    param_1[0x3b77c] = param_1[0x3b77c] + param_1[0xc];
    goto LAB_1001f6a0;
  case 2:
  case 0x11:
  case 0x31:
    uVar5 = param_1[0xc];
    iVar6 = param_1[0x3b77c] + (uVar5 >> 2);
    break;
  default:
    uVar5 = param_1[0xc];
    iVar6 = param_1[0x3b77c] + uVar5;
    break;
  case 6:
  case 7:
  case 0x99:
    uVar5 = param_1[0xc];
    iVar6 = param_1[0x3b77c] + (uVar5 >> 1);
  }
  param_1[0x3b77c] = iVar6;
  param_1[0x3b77a] = param_1[0x3b77a] + (int)uVar5 / 2;
LAB_1001f6a0:
  param_1 = param_1 + 0x3b77c;
  piVar1 = piVar3 + 0x3b767;
  iVar6 = *param_1 + 0x14 + piVar3[0x3b76a];
  *piVar1 = iVar6;
  if ((short)piVar3[0x3b76b] != 1) {
    *piVar1 = piVar3[0x3b779] + 8 + iVar6;
  }
  if (DAT_1007d6dc != 0) {
    FUN_100255a0((int)param_2,0,0,(char *)(piVar3 + 0x3b766),1,4);
    FUN_100255a0((int)param_2,4,0,(char *)piVar1,4,1);
    FUN_100255a0((int)param_2,8,0,(char *)(piVar3 + 0x3b768),1,4);
    FUN_100255a0((int)param_2,0xc,0,(char *)(piVar3 + 0x3b769),1,4);
    FUN_100255a0((int)param_2,0x10,0,(char *)(piVar3 + 0x3b76a),4,1);
    FUN_100255a0((int)param_2,0x14,0,(char *)(piVar3 + 0x3b76b),2,1);
    FUN_100255a0((int)param_2,0x16,0,(char *)((int)piVar3 + 0xeddae),2,1);
    FUN_100255a0((int)param_2,0x18,0,(char *)(piVar3 + 0x3b76c),4,1);
    FUN_100255a0((int)param_2,0x1c,0,(char *)(piVar3 + 0x3b76d),4,1);
    FUN_100255a0((int)param_2,0x20,0,(char *)(piVar3 + 0x3b76e),2,1);
    FUN_100255a0((int)param_2,0x22,0,(char *)((int)piVar3 + 0xeddba),2,1);
    iVar6 = 0x24;
    if ((short)piVar3[0x3b76b] != 1) {
      FUN_100255a0((int)param_2,0x24,0,(char *)(piVar3 + 0x3b76f),2,1);
      iVar6 = 0x26;
    }
    sVar2 = (short)piVar3[0x3b76b];
    if (sVar2 == 2) {
      FUN_100255a0((int)param_2,iVar6,0,(char *)((int)piVar3 + 0xeddbe),2,1);
      FUN_100255a0((int)param_2,iVar6 + 2,0,(char *)(piVar3 + 0x3b770),2,1);
      param_2 = (char *)(piVar3 + 0x3b771);
      iVar6 = iVar6 + 4;
      local_8 = 7;
      do {
        FUN_100255a0((int)pcVar4,iVar6,0,param_2 + -2,2,1);
        FUN_100255a0((int)pcVar4,iVar6 + 2,0,param_2,2,1);
        iVar6 = iVar6 + 4;
        param_2 = param_2 + 4;
        local_8 = local_8 + -1;
      } while (local_8 != 0);
    }
    else if ((sVar2 == 0x11) || (sVar2 == 0x31)) {
      FUN_100255a0((int)param_2,iVar6,0,(char *)((int)piVar3 + 0xeddbe),2,1);
      iVar6 = iVar6 + 2;
    }
    if ((short)piVar3[0x3b76b] != 1) {
      FUN_100255a0((int)pcVar4,iVar6,0,(char *)((int)piVar3 + 0xeddde),1,4);
      FUN_100255a0((int)pcVar4,iVar6 + 4,0,(char *)(piVar3 + 0x3b779),4,1);
      FUN_100255a0((int)pcVar4,iVar6 + 8,0,(char *)(piVar3 + 0x3b77a),4,1);
      iVar6 = iVar6 + 0xc;
    }
    FUN_100255a0((int)pcVar4,iVar6,0,(char *)(piVar3 + 0x3b77b),1,4);
    uVar5 = FUN_100255a0((int)pcVar4,iVar6 + 4,0,(char *)param_1,4,1);
    return CONCAT22((short)(uVar5 >> 0x10),1);
  }
  FUN_10025500((int)param_2,0,0,(char *)(piVar3 + 0x3b766),1,4);
  FUN_10025500((int)param_2,4,0,(char *)piVar1,4,1);
  FUN_10025500((int)param_2,8,0,(char *)(piVar3 + 0x3b768),1,4);
  FUN_10025500((int)param_2,0xc,0,(char *)(piVar3 + 0x3b769),1,4);
  FUN_10025500((int)param_2,0x10,0,(char *)(piVar3 + 0x3b76a),4,1);
  FUN_10025500((int)param_2,0x14,0,(char *)(piVar3 + 0x3b76b),2,1);
  FUN_10025500((int)param_2,0x16,0,(char *)((int)piVar3 + 0xeddae),2,1);
  FUN_10025500((int)param_2,0x18,0,(char *)(piVar3 + 0x3b76c),4,1);
  FUN_10025500((int)param_2,0x1c,0,(char *)(piVar3 + 0x3b76d),4,1);
  FUN_10025500((int)param_2,0x20,0,(char *)(piVar3 + 0x3b76e),2,1);
  FUN_10025500((int)param_2,0x22,0,(char *)((int)piVar3 + 0xeddba),2,1);
  iVar6 = 0x24;
  if ((short)piVar3[0x3b76b] != 1) {
    FUN_10025500((int)param_2,0x24,0,(char *)(piVar3 + 0x3b76f),2,1);
    iVar6 = 0x26;
  }
  sVar2 = (short)piVar3[0x3b76b];
  if (sVar2 == 2) {
    FUN_10025500((int)param_2,iVar6,0,(char *)((int)piVar3 + 0xeddbe),2,1);
    FUN_10025500((int)param_2,iVar6 + 2,0,(char *)(piVar3 + 0x3b770),2,1);
    param_2 = (char *)(piVar3 + 0x3b771);
    iVar6 = iVar6 + 4;
    local_8 = 7;
    do {
      FUN_10025500((int)pcVar4,iVar6,0,param_2 + -2,2,1);
      FUN_10025500((int)pcVar4,iVar6 + 2,0,param_2,2,1);
      iVar6 = iVar6 + 4;
      param_2 = param_2 + 4;
      local_8 = local_8 + -1;
    } while (local_8 != 0);
  }
  else if ((sVar2 == 0x11) || (sVar2 == 0x31)) {
    FUN_10025500((int)param_2,iVar6,0,(char *)((int)piVar3 + 0xeddbe),2,1);
    iVar6 = iVar6 + 2;
  }
  if ((short)piVar3[0x3b76b] != 1) {
    FUN_10025500((int)pcVar4,iVar6,0,(char *)((int)piVar3 + 0xeddde),1,4);
    FUN_10025500((int)pcVar4,iVar6 + 4,0,(char *)(piVar3 + 0x3b779),4,1);
    FUN_10025500((int)pcVar4,iVar6 + 8,0,(char *)(piVar3 + 0x3b77a),4,1);
    iVar6 = iVar6 + 0xc;
  }
  FUN_10025500((int)pcVar4,iVar6,0,(char *)(piVar3 + 0x3b77b),1,4);
  uVar5 = FUN_10025500((int)pcVar4,iVar6 + 4,0,(char *)param_1,4,1);
  return CONCAT22((short)(uVar5 >> 0x10),1);
}



/* ===== FUN_1001c990 ===== */
/* Entry: 1001c990 */

undefined4 __cdecl FUN_1001c990(byte *param_1,byte *param_2,int param_3)

{
  size_t sVar1;
  byte bVar2;
  int *piVar3;
  byte *pbVar4;
  undefined4 uVar5;
  int iVar6;
  int iVar7;
  byte *pbVar8;
  char *pcVar9;
  uint uVar10;
  int iVar11;
  uint uVar12;
  byte *pbVar13;
  byte *pbVar14;
  int iVar15;

  pbVar4 = param_1;
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
    bVar2 = *pbVar8;
    pbVar8 = pbVar8 + 1;
  } while (bVar2 != 0);
  iVar11 = ~uVar10 - 1;
  if (param_3 == 0) {
    pbVar8 = (byte *)FUN_1001d9c0(~uVar10);
    if (pbVar8 == (byte *)0x0) {
      return 0;
    }
    uVar10 = 0xffffffff;
    do {
      pbVar13 = param_2;
      if (uVar10 == 0) break;
      uVar10 = uVar10 - 1;
      pbVar13 = param_2 + 1;
      bVar2 = *param_2;
      param_2 = pbVar13;
    } while (bVar2 != 0);
    uVar10 = ~uVar10;
    pbVar13 = pbVar13 + -uVar10;
    pbVar14 = pbVar8;
    for (uVar12 = uVar10 >> 2; uVar12 != 0; uVar12 = uVar12 - 1) {
      *(undefined4 *)pbVar14 = *(undefined4 *)pbVar13;
      pbVar13 = pbVar13 + 4;
      pbVar14 = pbVar14 + 4;
    }
    for (uVar10 = uVar10 & 3; param_2 = pbVar8, uVar10 != 0; uVar10 = uVar10 - 1) {
      *pbVar14 = *pbVar13;
      pbVar13 = pbVar13 + 1;
      pbVar14 = pbVar14 + 1;
    }
  }
  if (*(undefined **)(param_1 + 0x47794) != (undefined *)0x0) {
    FUN_1001da30(*(undefined **)(param_1 + 0x47794));
  }
  if (*(undefined **)(param_1 + 0x47798) != (undefined *)0x0) {
    FUN_1001da30(*(undefined **)(param_1 + 0x47798));
  }
  if (*(undefined **)(param_1 + 0x47790) != (undefined *)0x0) {
    FUN_1001da30(*(undefined **)(param_1 + 0x47790));
  }
  sVar1 = iVar11 * 4;
  uVar5 = FUN_1001d9c0(sVar1);
  *(undefined4 *)(param_1 + 0x47794) = uVar5;
  uVar5 = FUN_1001d9c0(sVar1);
  *(undefined4 *)(param_1 + 0x47798) = uVar5;
  uVar5 = FUN_1001d9c0(sVar1);
  *(undefined4 *)(param_1 + 0x47790) = uVar5;
  iVar6 = 0;
  if (0 < iVar11) {
    do {
      *(int *)(*(int *)(param_1 + 0x47794) + iVar6 * 4) = iVar6;
      *(int *)(*(int *)(param_1 + 0x47798) + iVar6 * 4) = iVar6;
      iVar6 = iVar6 + 1;
    } while (iVar6 < iVar11);
  }
  iVar15 = 0;
  iVar6 = 0;
  if (0 < iVar11) {
    param_1 = param_2;
    do {
      iVar7 = FUN_1001c900(param_1);
      if (iVar7 == 0) {
        *(int *)(*(int *)(pbVar4 + 0x47790) + iVar15 * 4) = iVar6;
        iVar15 = iVar15 + 1;
        param_1 = param_1 + 1;
      }
      else {
        iVar15 = iVar15 + 2;
        *(int *)(*(int *)(pbVar4 + 0x47790) + -8 + iVar15 * 4) = iVar6;
        *(int *)(*(int *)(pbVar4 + 0x47790) + -4 + iVar15 * 4) = iVar6;
        param_1 = param_1 + 2;
      }
      iVar6 = iVar6 + 1;
    } while (iVar15 < iVar11);
  }
  pbVar8 = FUN_1001ce10((int)pbVar4,param_2,param_3);
  pcVar9 = (char *)FUN_1001d250((int)pbVar4,(char *)pbVar8,FUN_1001d370);
  uVar5 = FUN_1001d250((int)pbVar4,pcVar9,FUN_1001d5d0);
  piVar3 = *(int **)(pbVar4 + 0x122448);
  param_3 = *piVar3;
  if (0 < param_3) {
    if (*(char *)(DAT_100a0460 + 0x20424) == '\x01') {
      if (0 < param_3) {
        iVar11 = 0;
        do {
          iVar6 = iVar11 + 8;
          iVar15 = iVar11 + 8;
          iVar11 = iVar11 + 0x210;
          *(undefined4 *)(iVar15 + piVar3[1]) =
               *(undefined4 *)
                (*(int *)(pbVar4 + 0x47790) +
                *(int *)(*(int *)(pbVar4 + 0x47794) + *(int *)(iVar6 + piVar3[1]) * 4) * 4);
          param_3 = param_3 + -1;
        } while (param_3 != 0);
      }
    }
    else if (0 < param_3) {
      iVar11 = 0;
      do {
        iVar6 = iVar11 + 8;
        iVar15 = iVar11 + 8;
        iVar11 = iVar11 + 0x210;
        *(undefined4 *)(iVar15 + piVar3[1]) =
             *(undefined4 *)(*(int *)(pbVar4 + 0x47794) + *(int *)(iVar6 + piVar3[1]) * 4);
        param_3 = param_3 + -1;
      } while (param_3 != 0);
    }
  }
  VT_SetDecimal0Pron_ENG();
  return uVar5;
}



/* ===== FUN_10026ab0 ===== */
/* Entry: 10026ab0 */

undefined4 __cdecl FUN_10026ab0(int param_1,int *param_2,char *param_3)

{
  char cVar1;
  int iVar2;
  int iVar3;
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
  param_2[0x11] = 0;
  iVar3 = FUN_10026630(param_1,param_2);
  iVar2 = param_2[0x11];
  while( true ) {
    if (iVar2 != 0) {
      return CONCAT22((short)((uint)iVar3 >> 0x10),0xffff);
    }
    if ((-1 < (short)iVar3) && (*(short *)(param_2[0x13] + 2) != 0)) break;
    iVar3 = FUN_10026630(param_1,param_2);
    iVar2 = param_2[0x11];
  }
  return param_2[0x13] & 0xffff0000;
}



/* ===== FUN_10026870 ===== */
/* Entry: 10026870 */

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



/* ===== FUN_1001e0c0 ===== */
/* Entry: 1001e0c0 */

undefined4 __cdecl FUN_1001e0c0(int param_1,int param_2)

{
  int *piVar1;
  int iVar2;
  undefined4 in_EAX;
  short *psVar3;
  int iVar4;
  undefined4 local_614;
  int local_14;
  int local_10;
  int *local_c;
  int local_8;

  piVar1 = *(int **)(param_1 + 0x2c);
  local_c = *(int **)(param_1 + 0x47774);
  if (piVar1 == (int *)0x0) {
    return CONCAT22((short)((uint)in_EAX >> 0x10),1);
  }
  if (*(int *)(param_1 + 0x44) == 0) {
    if (piVar1[2] == local_c[3]) goto LAB_1001e278;
  }
  else {
    piVar1[2] = local_c[3];
  }
  local_10 = piVar1[2];
  if (local_10 < *(int *)(param_1 + 0x47778)) {
    local_10 = local_10 + 600;
  }
  local_14 = *(int *)(param_1 + 0x47778);
  if (local_14 < local_10) {
    do {
      local_8 = 0;
      iVar4 = (local_14 % 600) * 0x24;
      psVar3 = (short *)(*local_c + iVar4);
      if (0 < *(short *)(*local_c + iVar4)) {
        do {
          iVar2 = *(int *)(*(int *)(psVar3 + 2) + local_8 * 8);
          if (*(short *)(*(int *)(psVar3 + 2) + local_8 * 8 + 4) == 0x28) {
            FUN_10025e40((undefined4 *)piVar1[4],(byte *)s__2d__sil_____d_1007d5f8);
          }
          else {
            FUN_10025e40((undefined4 *)piVar1[4],(byte *)s__2d___3s_____d_1007d5e4);
          }
          piVar1[1] = piVar1[1] + iVar2;
          psVar3 = (short *)(*local_c + iVar4);
          local_8 = local_8 + 1;
        } while (local_8 < *psVar3);
      }
      local_8 = *(int *)(*local_c + 8 + iVar4);
      iVar4 = *local_c + iVar4;
      iVar2 = *(int *)(iVar4 + 0xc);
      iVar4 = *(int *)(iVar4 + 0x10) - iVar2;
      FUN_10063f30(&local_614,(undefined4 *)(iVar2 + param_2),iVar4 + 1);
      *(undefined1 *)((int)&local_614 + iVar4 + 1) = 0;
      FUN_10025e40((undefined4 *)piVar1[4],(byte *)s_________________________________1007d5c0);
      iVar4 = local_8;
      FUN_10025e40((undefined4 *)piVar1[4],(byte *)s__d____d__s____d_1007d5ac);
      FUN_10025e40((undefined4 *)piVar1[4],(byte *)s_________________________________1007d588);
      local_14 = local_14 + 1;
      *piVar1 = *piVar1 + iVar4;
    } while (local_14 < local_10);
  }
  *(int *)(param_1 + 0x47778) = local_10;
  if (599 < local_10) {
    *(int *)(param_1 + 0x47778) = local_10 + -600;
  }
LAB_1001e278:
  iVar4 = local_c[3];
  piVar1[2] = iVar4;
  return CONCAT22((short)((uint)iVar4 >> 0x10),1);
}



/* ===== FUN_10024d20 ===== */
/* Entry: 10024d20 */

undefined * __cdecl FUN_10024d20(undefined4 param_1,byte *param_2)

{
  byte bVar1;
  undefined *puVar2;
  undefined4 *puVar3;
  LPCSTR pCVar4;
  int iVar5;
  uint uVar6;
  uint uVar7;
  byte *pbVar8;
  byte *pbVar9;
  uint local_c;
  byte *local_8;

  puVar2 = (undefined *)FUN_1001d9c0(0x420);
  if (puVar2 != (undefined *)0x0) {
    uVar6 = 0xffffffff;
    local_8 = puVar2 + 0x14;
    pbVar8 = param_2;
    do {
      pbVar9 = pbVar8;
      if (uVar6 == 0) break;
      uVar6 = uVar6 - 1;
      pbVar9 = pbVar8 + 1;
      bVar1 = *pbVar8;
      pbVar8 = pbVar9;
    } while (bVar1 != 0);
    uVar6 = ~uVar6;
    pbVar8 = pbVar9 + -uVar6;
    pbVar9 = local_8;
    for (uVar7 = uVar6 >> 2; uVar7 != 0; uVar7 = uVar7 - 1) {
      *(undefined4 *)pbVar9 = *(undefined4 *)pbVar8;
      pbVar8 = pbVar8 + 4;
      pbVar9 = pbVar9 + 4;
    }
    for (uVar6 = uVar6 & 3; uVar6 != 0; uVar6 = uVar6 - 1) {
      *pbVar9 = *pbVar8;
      pbVar8 = pbVar8 + 1;
      pbVar9 = pbVar9 + 1;
    }
    puVar3 = FUN_10024e20(puVar2,&local_c,param_2);
    if (puVar3 == (undefined4 *)0x0) {
      pCVar4 = (LPCSTR)FUN_10025360((int)puVar2,param_1);
      if ((pCVar4 != (LPCSTR)0x0) && (*pCVar4 != '\0')) {
        iVar5 = FUN_10025370((int)puVar2,pCVar4,local_c,0x80);
        if (-1 < iVar5) {
          FUN_10025400((int)puVar2,0,0);
          FUN_10025400((int)puVar2,0,2);
          *(undefined4 *)(puVar2 + 0x10) = *(undefined4 *)(puVar2 + 0xc);
          FUN_10025400((int)puVar2,0,0);
          if (*param_2 == 0x72) {
            iVar5 = FUN_10024cd0((int)puVar2);
            if (iVar5 == 0) {
              DAT_100a74b8 = DAT_100a74b8 + *(int *)(puVar2 + 0x10);
            }
          }
          return puVar2;
        }
      }
    }
    FUN_1001da30(puVar2);
  }
  return (undefined *)0x0;
}



/* ===== FUN_10024e20 ===== */
/* Entry: 10024e20 */

undefined4 * __cdecl FUN_10024e20(undefined4 param_1,undefined4 *param_2,byte *param_3)

{
  byte bVar1;
  byte *pbVar2;
  int iVar3;
  byte *pbVar4;
  bool bVar5;

  bVar1 = *param_3;
  if (bVar1 == 0x61) {
    pbVar4 = &DAT_10077718;
    pbVar2 = param_3;
    do {
      bVar1 = *pbVar2;
      bVar5 = bVar1 < *pbVar4;
      if (bVar1 != *pbVar4) {
LAB_100251d5:
        iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
        goto LAB_100251da;
      }
      if (bVar1 == 0) break;
      bVar1 = pbVar2[1];
      bVar5 = bVar1 < pbVar4[1];
      if (bVar1 != pbVar4[1]) goto LAB_100251d5;
      pbVar2 = pbVar2 + 2;
      pbVar4 = pbVar4 + 2;
    } while (bVar1 != 0);
    iVar3 = 0;
LAB_100251da:
    if (iVar3 != 0) {
      pbVar4 = &DAT_10077694;
      pbVar2 = param_3;
      do {
        bVar1 = *pbVar2;
        bVar5 = bVar1 < *pbVar4;
        if (bVar1 != *pbVar4) {
LAB_10025211:
          iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
          goto LAB_10025216;
        }
        if (bVar1 == 0) break;
        bVar1 = pbVar2[1];
        bVar5 = bVar1 < pbVar4[1];
        if (bVar1 != pbVar4[1]) goto LAB_10025211;
        pbVar2 = pbVar2 + 2;
        pbVar4 = pbVar4 + 2;
      } while (bVar1 != 0);
      iVar3 = 0;
LAB_10025216:
      if (iVar3 != 0) {
        pbVar4 = &DAT_1007d6c0;
        pbVar2 = param_3;
        do {
          bVar1 = *pbVar2;
          bVar5 = bVar1 < *pbVar4;
          if (bVar1 != *pbVar4) {
LAB_1002524d:
            iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
            goto LAB_10025252;
          }
          if (bVar1 == 0) break;
          bVar1 = pbVar2[1];
          bVar5 = bVar1 < pbVar4[1];
          if (bVar1 != pbVar4[1]) goto LAB_1002524d;
          pbVar2 = pbVar2 + 2;
          pbVar4 = pbVar4 + 2;
        } while (bVar1 != 0);
        iVar3 = 0;
LAB_10025252:
        if (iVar3 != 0) {
          pbVar4 = &DAT_1007d6bc;
          pbVar2 = param_3;
          do {
            bVar1 = *pbVar2;
            bVar5 = bVar1 < *pbVar4;
            if (bVar1 != *pbVar4) {
LAB_10025289:
              iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
              goto LAB_1002528e;
            }
            if (bVar1 == 0) break;
            bVar1 = pbVar2[1];
            bVar5 = bVar1 < pbVar4[1];
            if (bVar1 != pbVar4[1]) goto LAB_10025289;
            pbVar2 = pbVar2 + 2;
            pbVar4 = pbVar4 + 2;
          } while (bVar1 != 0);
          iVar3 = 0;
LAB_1002528e:
          if (iVar3 != 0) {
            pbVar4 = &DAT_1007d6b8;
            pbVar2 = param_3;
            do {
              bVar1 = *pbVar2;
              bVar5 = bVar1 < *pbVar4;
              if (bVar1 != *pbVar4) {
LAB_100252c5:
                iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
                goto LAB_100252ca;
              }
              if (bVar1 == 0) break;
              bVar1 = pbVar2[1];
              bVar5 = bVar1 < pbVar4[1];
              if (bVar1 != pbVar4[1]) goto LAB_100252c5;
              pbVar2 = pbVar2 + 2;
              pbVar4 = pbVar4 + 2;
            } while (bVar1 != 0);
            iVar3 = 0;
LAB_100252ca:
            if (iVar3 == 0) {
              *param_2 = 0x8109;
              return (undefined4 *)0x0;
            }
            pbVar2 = &DAT_1007d6b4;
            do {
              bVar1 = *param_3;
              bVar5 = bVar1 < *pbVar2;
              if (bVar1 != *pbVar2) {
LAB_1002530d:
                iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
                goto LAB_10025312;
              }
              if (bVar1 == 0) break;
              bVar1 = param_3[1];
              bVar5 = bVar1 < pbVar2[1];
              if (bVar1 != pbVar2[1]) goto LAB_1002530d;
              param_3 = param_3 + 2;
              pbVar2 = pbVar2 + 2;
            } while (bVar1 != 0);
            iVar3 = 0;
LAB_10025312:
            if (iVar3 != 0) {
              return param_2;
            }
            *param_2 = 0x810a;
            return (undefined4 *)0x0;
          }
        }
        *param_2 = 0x410a;
        return (undefined4 *)0x1;
      }
    }
    *param_2 = 0x4109;
    return (undefined4 *)0x1;
  }
  if (bVar1 != 0x72) {
    if (bVar1 != 0x77) {
      *param_2 = 0x4000;
      return (undefined4 *)0x1;
    }
    pbVar4 = &DAT_1007d348;
    pbVar2 = param_3;
    do {
      bVar1 = *pbVar2;
      bVar5 = bVar1 < *pbVar4;
      if (bVar1 != *pbVar4) {
LAB_10024e81:
        iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
        goto LAB_10024e86;
      }
      if (bVar1 == 0) break;
      bVar1 = pbVar2[1];
      bVar5 = bVar1 < pbVar4[1];
      if (bVar1 != pbVar4[1]) goto LAB_10024e81;
      pbVar2 = pbVar2 + 2;
      pbVar4 = pbVar4 + 2;
    } while (bVar1 != 0);
    iVar3 = 0;
LAB_10024e86:
    if (iVar3 != 0) {
      pbVar4 = &DAT_1007d534;
      pbVar2 = param_3;
      do {
        bVar1 = *pbVar2;
        bVar5 = bVar1 < *pbVar4;
        if (bVar1 != *pbVar4) {
LAB_10024ebd:
          iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
          goto LAB_10024ec2;
        }
        if (bVar1 == 0) break;
        bVar1 = pbVar2[1];
        bVar5 = bVar1 < pbVar4[1];
        if (bVar1 != pbVar4[1]) goto LAB_10024ebd;
        pbVar2 = pbVar2 + 2;
        pbVar4 = pbVar4 + 2;
      } while (bVar1 != 0);
      iVar3 = 0;
LAB_10024ec2:
      if (iVar3 != 0) {
        pbVar4 = &DAT_1007d6d8;
        pbVar2 = param_3;
        do {
          bVar1 = *pbVar2;
          bVar5 = bVar1 < *pbVar4;
          if (bVar1 != *pbVar4) {
LAB_10024ef9:
            iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
            goto LAB_10024efe;
          }
          if (bVar1 == 0) break;
          bVar1 = pbVar2[1];
          bVar5 = bVar1 < pbVar4[1];
          if (bVar1 != pbVar4[1]) goto LAB_10024ef9;
          pbVar2 = pbVar2 + 2;
          pbVar4 = pbVar4 + 2;
        } while (bVar1 != 0);
        iVar3 = 0;
LAB_10024efe:
        if (iVar3 != 0) {
          pbVar4 = &DAT_1007d6d4;
          pbVar2 = param_3;
          do {
            bVar1 = *pbVar2;
            bVar5 = bVar1 < *pbVar4;
            if (bVar1 != *pbVar4) {
LAB_10024f35:
              iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
              goto LAB_10024f3a;
            }
            if (bVar1 == 0) break;
            bVar1 = pbVar2[1];
            bVar5 = bVar1 < pbVar4[1];
            if (bVar1 != pbVar4[1]) goto LAB_10024f35;
            pbVar2 = pbVar2 + 2;
            pbVar4 = pbVar4 + 2;
          } while (bVar1 != 0);
          iVar3 = 0;
LAB_10024f3a:
          if (iVar3 != 0) {
            pbVar4 = &DAT_10079290;
            pbVar2 = param_3;
            do {
              bVar1 = *pbVar2;
              bVar5 = bVar1 < *pbVar4;
              if (bVar1 != *pbVar4) {
LAB_10024f71:
                iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
                goto LAB_10024f76;
              }
              if (bVar1 == 0) break;
              bVar1 = pbVar2[1];
              bVar5 = bVar1 < pbVar4[1];
              if (bVar1 != pbVar4[1]) goto LAB_10024f71;
              pbVar2 = pbVar2 + 2;
              pbVar4 = pbVar4 + 2;
            } while (bVar1 != 0);
            iVar3 = 0;
LAB_10024f76:
            if (iVar3 == 0) {
              *param_2 = 0x8301;
              return (undefined4 *)0x0;
            }
            pbVar2 = &DAT_1007d60c;
            do {
              bVar1 = *param_3;
              bVar5 = bVar1 < *pbVar2;
              if (bVar1 != *pbVar2) {
LAB_10024fb9:
                iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
                goto LAB_10024fbe;
              }
              if (bVar1 == 0) break;
              bVar1 = param_3[1];
              bVar5 = bVar1 < pbVar2[1];
              if (bVar1 != pbVar2[1]) goto LAB_10024fb9;
              param_3 = param_3 + 2;
              pbVar2 = pbVar2 + 2;
            } while (bVar1 != 0);
            iVar3 = 0;
LAB_10024fbe:
            if (iVar3 != 0) {
              return param_2;
            }
            *param_2 = 0x8302;
            return (undefined4 *)0x0;
          }
        }
        *param_2 = 0x4302;
        return (undefined4 *)0x1;
      }
    }
    *param_2 = 0x4301;
    return (undefined4 *)0x1;
  }
  pbVar4 = &DAT_1007d358;
  pbVar2 = param_3;
  do {
    bVar1 = *pbVar2;
    bVar5 = bVar1 < *pbVar4;
    if (bVar1 != *pbVar4) {
LAB_1002502b:
      iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
      goto LAB_10025030;
    }
    if (bVar1 == 0) break;
    bVar1 = pbVar2[1];
    bVar5 = bVar1 < pbVar4[1];
    if (bVar1 != pbVar4[1]) goto LAB_1002502b;
    pbVar2 = pbVar2 + 2;
    pbVar4 = pbVar4 + 2;
  } while (bVar1 != 0);
  iVar3 = 0;
LAB_10025030:
  if (iVar3 != 0) {
    pbVar4 = &DAT_1007d6d0;
    pbVar2 = param_3;
    do {
      bVar1 = *pbVar2;
      bVar5 = bVar1 < *pbVar4;
      if (bVar1 != *pbVar4) {
LAB_10025067:
        iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
        goto LAB_1002506c;
      }
      if (bVar1 == 0) break;
      bVar1 = pbVar2[1];
      bVar5 = bVar1 < pbVar4[1];
      if (bVar1 != pbVar4[1]) goto LAB_10025067;
      pbVar2 = pbVar2 + 2;
      pbVar4 = pbVar4 + 2;
    } while (bVar1 != 0);
    iVar3 = 0;
LAB_1002506c:
    if (iVar3 != 0) {
      pbVar4 = &DAT_1007d6cc;
      pbVar2 = param_3;
      do {
        bVar1 = *pbVar2;
        bVar5 = bVar1 < *pbVar4;
        if (bVar1 != *pbVar4) {
LAB_100250a3:
          iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
          goto LAB_100250a8;
        }
        if (bVar1 == 0) break;
        bVar1 = pbVar2[1];
        bVar5 = bVar1 < pbVar4[1];
        if (bVar1 != pbVar4[1]) goto LAB_100250a3;
        pbVar2 = pbVar2 + 2;
        pbVar4 = pbVar4 + 2;
      } while (bVar1 != 0);
      iVar3 = 0;
LAB_100250a8:
      if (iVar3 != 0) {
        pbVar4 = &DAT_1007d6c8;
        pbVar2 = param_3;
        do {
          bVar1 = *pbVar2;
          bVar5 = bVar1 < *pbVar4;
          if (bVar1 != *pbVar4) {
LAB_100250df:
            iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
            goto LAB_100250e4;
          }
          if (bVar1 == 0) break;
          bVar1 = pbVar2[1];
          bVar5 = bVar1 < pbVar4[1];
          if (bVar1 != pbVar4[1]) goto LAB_100250df;
          pbVar2 = pbVar2 + 2;
          pbVar4 = pbVar4 + 2;
        } while (bVar1 != 0);
        iVar3 = 0;
LAB_100250e4:
        if (iVar3 != 0) {
          pbVar4 = &DAT_100771a8;
          pbVar2 = param_3;
          do {
            bVar1 = *pbVar2;
            bVar5 = bVar1 < *pbVar4;
            if (bVar1 != *pbVar4) {
LAB_1002511b:
              iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
              goto LAB_10025120;
            }
            if (bVar1 == 0) break;
            bVar1 = pbVar2[1];
            bVar5 = bVar1 < pbVar4[1];
            if (bVar1 != pbVar4[1]) goto LAB_1002511b;
            pbVar2 = pbVar2 + 2;
            pbVar4 = pbVar4 + 2;
          } while (bVar1 != 0);
          iVar3 = 0;
LAB_10025120:
          if (iVar3 == 0) {
            *param_2 = 0x8000;
            return (undefined4 *)0x0;
          }
          pbVar2 = &DAT_1007d6c4;
          do {
            bVar1 = *param_3;
            bVar5 = bVar1 < *pbVar2;
            if (bVar1 != *pbVar2) {
LAB_10025163:
              iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
              goto LAB_10025168;
            }
            if (bVar1 == 0) break;
            bVar1 = param_3[1];
            bVar5 = bVar1 < pbVar2[1];
            if (bVar1 != pbVar2[1]) goto LAB_10025163;
            param_3 = param_3 + 2;
            pbVar2 = pbVar2 + 2;
          } while (bVar1 != 0);
          iVar3 = 0;
LAB_10025168:
          if (iVar3 != 0) {
            return param_2;
          }
          *param_2 = 0x8002;
          return (undefined4 *)0x0;
        }
      }
      *param_2 = 0x4002;
      return (undefined4 *)0x1;
    }
  }
  *param_2 = 0x4000;
  return (undefined4 *)0x1;
}



/* ===== FUN_10025370 ===== */
/* Entry: 10025370 */

int __cdecl FUN_10025370(int param_1,LPCSTR param_2,uint param_3,uint param_4)

{
  undefined4 uVar1;
  uint unaff_ESI;

  if (*(char *)(param_1 + 0x14) == 'r') {
    uVar1 = FUN_1006a39e(param_2,param_3,unaff_ESI);
    *(undefined4 *)(param_1 + 4) = uVar1;
    return ((-1 < *(int *)(param_1 + 4)) - 1 & 0xfffffffe) + 1;
  }
  uVar1 = FUN_1006a39e(param_2,param_3,param_4);
  *(undefined4 *)(param_1 + 4) = uVar1;
  return ((-1 < *(int *)(param_1 + 4)) - 1 & 0xfffffffe) + 1;
}



/* ===== FUN_10025400 ===== */
/* Entry: 10025400 */

int __cdecl FUN_10025400(int param_1,LONG param_2,DWORD param_3)

{
  DWORD DVar1;
  DWORD DVar2;

  DVar1 = FUN_1006906d(*(uint *)(param_1 + 4),param_2,param_3);
  DVar2 = FUN_1006906d(*(uint *)(param_1 + 4),0,1);
  *(DWORD *)(param_1 + 0xc) = DVar2;
  return (-1 < (int)DVar1) - 1;
}



/* ===== FUN_10025500 ===== */
/* Entry: 10025500 */

uint __cdecl
FUN_10025500(int param_1,int param_2,DWORD param_3,char *param_4,uint param_5,int param_6)

{
  uint uVar1;

  EnterCriticalSection((LPCRITICAL_SECTION)&DAT_100a8400);
  if (param_3 == 0) {
    if (param_2 == *(int *)(param_1 + 0xc)) goto LAB_10025537;
    param_3 = 0;
  }
  FUN_10025400(param_1,param_2,param_3);
LAB_10025537:
  uVar1 = FUN_10025560(param_4,param_5,param_6,param_1);
  LeaveCriticalSection((LPCRITICAL_SECTION)&DAT_100a8400);
  return uVar1;
}



/* ===== FUN_100255a0 ===== */
/* Entry: 100255a0 */

uint __cdecl
FUN_100255a0(int param_1,int param_2,DWORD param_3,char *param_4,uint param_5,int param_6)

{
  uint uVar1;

  EnterCriticalSection((LPCRITICAL_SECTION)&DAT_100a8400);
  if (param_3 == 0) {
    if (param_2 == *(int *)(param_1 + 0xc)) goto LAB_100255d7;
    param_3 = 0;
  }
  FUN_10025400(param_1,param_2,param_3);
LAB_100255d7:
  uVar1 = FUN_10025600(param_4,param_5,param_6,param_1);
  LeaveCriticalSection((LPCRITICAL_SECTION)&DAT_100a8400);
  return uVar1;
}



/* ===== FUN_10025560 ===== */
/* Entry: 10025560 */

uint __cdecl FUN_10025560(char *param_1,uint param_2,int param_3,int param_4)

{
  uint uVar1;

  uVar1 = FUN_10069145(*(undefined **)(param_4 + 4),param_1,param_2 * param_3);
  *(uint *)(param_4 + 0xc) = *(int *)(param_4 + 0xc) + param_2 * param_3;
  return uVar1 / param_2;
}



/* ===== FUN_10025600 ===== */
/* Entry: 10025600 */

/* WARNING: Type propagation algorithm not settling */

uint __cdecl FUN_10025600(char *param_1,uint param_2,int param_3,int param_4)

{
  uint uVar1;
  char *pcVar2;
  uint uVar3;
  char local_24 [15];
  char acStack_15 [17];

  if ((param_2 != 1) && (param_2 < 9)) {
    uVar3 = 0;
    for (; param_3 != 0; param_3 = param_3 + -1) {
      FUN_10063f30((undefined4 *)((int)acStack_15 + 1),(undefined4 *)param_1,param_2);
      uVar1 = 0;
      if (param_2 != 0) {
        pcVar2 = acStack_15 + param_2;
        do {
          local_24[uVar1] = *pcVar2;
          uVar1 = uVar1 + 1;
          pcVar2 = pcVar2 + -1;
        } while (uVar1 < param_2);
      }
      uVar1 = FUN_10025560(local_24,1,param_2,param_4);
      uVar3 = uVar3 + uVar1;
      param_1 = param_1 + param_2;
    }
    return uVar3 / param_2;
  }
  uVar3 = FUN_10025560(param_1,param_2,param_3,param_4);
  return uVar3;
}
