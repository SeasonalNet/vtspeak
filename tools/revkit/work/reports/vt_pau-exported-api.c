/* Ghidra decompiler output; reconstructed pseudocode, not original source. */
/* Program: vt_pau.dll */

/* ===== VT_CsvParser_Init_ENG @ 10016bc0 ===== */

undefined4 * VT_CsvParser_Init_ENG(void)

{
  undefined4 *puVar1;
  undefined4 uVar2;

  puVar1 = (undefined4 *)FUN_1001d9c0(0x18);
  if (puVar1 != (undefined4 *)0x0) {
    *puVar1 = 0;
    puVar1[1] = 0;
    puVar1[2] = 0;
    puVar1[3] = 100;
    puVar1[4] = 0;
    uVar2 = FUN_1001c440(PTR_DAT_1007b6bc);
    puVar1[5] = uVar2;
  }
  return puVar1;
}



/* ===== VT_CsvParser_Parsing_ENG @ 10016bd0 ===== */

void __cdecl VT_CsvParser_Parsing_ENG(undefined4 *param_1,char *param_2,int param_3)

{
                    /* 0x16bd0  14  VT_CsvParser_Parsing_ENG */
  FUN_100166b0(param_1,param_2,param_3);
  return;
}



/* ===== VT_CsvParser_Exit_ENG @ 10016bf0 ===== */

void __cdecl VT_CsvParser_Exit_ENG(undefined *param_1)

{
                    /* 0x16bf0  8  VT_CsvParser_Exit_ENG */
  FUN_10016840(param_1);
  return;
}



/* ===== VT_CsvParser_GetNfields_ENG @ 10016c10 ===== */

void __cdecl VT_CsvParser_GetNfields_ENG(int param_1)

{
                    /* 0x16c10  10  VT_CsvParser_GetNfields_ENG */
  FUN_100168b0(param_1);
  return;
}



/* ===== VT_CsvParser_GetField_ENG @ 10016c30 ===== */

void __cdecl VT_CsvParser_GetField_ENG(int param_1,int param_2)

{
                    /* 0x16c30  9  VT_CsvParser_GetField_ENG */
  FUN_100168f0(param_1,param_2);
  return;
}



/* ===== VT_CsvParser_MakeCsv_ENG @ 10016c50 ===== */

void __cdecl
VT_CsvParser_MakeCsv_ENG(undefined4 *param_1,int param_2,undefined4 *param_3,uint param_4)

{
                    /* 0x16c50  13  VT_CsvParser_MakeCsv_ENG */
  FUN_10016940(param_1,param_2,param_3,param_4);
  return;
}



/* ===== VT_CsvParser_IsCsv_ENG @ 10016c70 ===== */

void __cdecl VT_CsvParser_IsCsv_ENG(int param_1,int param_2,int param_3)

{
                    /* 0x16c70  12  VT_CsvParser_IsCsv_ENG */
  FUN_10016b30(param_1,param_2,param_3);
  return;
}



/* ===== VT_GetPathKey_ENG @ 1001c690 ===== */

undefined4 * __cdecl VT_GetPathKey_ENG(int param_1)

{
  char cVar1;
  byte bVar2;
  byte *pbVar3;
  int iVar4;
  uint uVar5;
  uint uVar6;
  byte *pbVar7;
  char *pcVar8;
  char *pcVar9;
  char *pcVar10;
  bool bVar11;

                    /* 0x1c690  24  VT_GetPathKey_ENG */
  pbVar7 = &DAT_1009f948;
  pbVar3 = (&PTR_DAT_1007c6b0)[param_1 * 6];
  do {
    bVar2 = *pbVar3;
    bVar11 = bVar2 < *pbVar7;
    if (bVar2 != *pbVar7) {
LAB_1001c6d2:
      iVar4 = (1 - (uint)bVar11) - (uint)(bVar11 != 0);
      goto LAB_1001c6d7;
    }
    if (bVar2 == 0) break;
    bVar2 = pbVar3[1];
    bVar11 = bVar2 < pbVar7[1];
    if (bVar2 != pbVar7[1]) goto LAB_1001c6d2;
    pbVar3 = pbVar3 + 2;
    pbVar7 = pbVar7 + 2;
  } while (bVar2 != 0);
  iVar4 = 0;
LAB_1001c6d7:
  if (iVar4 == 0) {
    return (undefined4 *)0x0;
  }
  FUN_10063ed6((undefined1 *)&DAT_100fe6e0,&DAT_10079404);
  uVar5 = 0xffffffff;
  pcVar8 = &DAT_1007c934;
  do {
    pcVar10 = pcVar8;
    if (uVar5 == 0) break;
    uVar5 = uVar5 - 1;
    pcVar10 = pcVar8 + 1;
    cVar1 = *pcVar8;
    pcVar8 = pcVar10;
  } while (cVar1 != '\0');
  uVar5 = ~uVar5;
  iVar4 = -1;
  pcVar8 = (char *)&DAT_100fe6e0;
  do {
    pcVar9 = pcVar8;
    if (iVar4 == 0) break;
    iVar4 = iVar4 + -1;
    pcVar9 = pcVar8 + 1;
    cVar1 = *pcVar8;
    pcVar8 = pcVar9;
  } while (cVar1 != '\0');
  pcVar8 = pcVar10 + -uVar5;
  pcVar10 = pcVar9 + -1;
  for (uVar6 = uVar5 >> 2; uVar6 != 0; uVar6 = uVar6 - 1) {
    *(undefined4 *)pcVar10 = *(undefined4 *)pcVar8;
    pcVar8 = pcVar8 + 4;
    pcVar10 = pcVar10 + 4;
  }
  for (uVar5 = uVar5 & 3; uVar5 != 0; uVar5 = uVar5 - 1) {
    *pcVar10 = *pcVar8;
    pcVar8 = pcVar8 + 1;
    pcVar10 = pcVar10 + 1;
  }
  uVar5 = 0xffffffff;
  pcVar8 = &DAT_1007c930;
  do {
    pcVar10 = pcVar8;
    if (uVar5 == 0) break;
    uVar5 = uVar5 - 1;
    pcVar10 = pcVar8 + 1;
    cVar1 = *pcVar8;
    pcVar8 = pcVar10;
  } while (cVar1 != '\0');
  uVar5 = ~uVar5;
  iVar4 = -1;
  pcVar8 = (char *)&DAT_100fe6e0;
  do {
    pcVar9 = pcVar8;
    if (iVar4 == 0) break;
    iVar4 = iVar4 + -1;
    pcVar9 = pcVar8 + 1;
    cVar1 = *pcVar8;
    pcVar8 = pcVar9;
  } while (cVar1 != '\0');
  pcVar8 = pcVar10 + -uVar5;
  pcVar10 = pcVar9 + -1;
  for (uVar6 = uVar5 >> 2; uVar6 != 0; uVar6 = uVar6 - 1) {
    *(undefined4 *)pcVar10 = *(undefined4 *)pcVar8;
    pcVar8 = pcVar8 + 4;
    pcVar10 = pcVar10 + 4;
  }
  for (uVar5 = uVar5 & 3; uVar5 != 0; uVar5 = uVar5 - 1) {
    *pcVar10 = *pcVar8;
    pcVar8 = pcVar8 + 1;
    pcVar10 = pcVar10 + 1;
  }
  return &DAT_100fe6e0;
}



/* ===== VT_TextToFile_ENG @ 1001da50 ===== */

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



/* ===== VT_TextToBuffer_ENG @ 1001dc40 ===== */

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



/* ===== VT_TextToPreprocessInfoFile_ENG @ 1001dd60 ===== */

void __cdecl
VT_TextToPreprocessInfoFile_ENG
          (byte *param_1,LPCSTR param_2,byte param_3,int param_4,int param_5,int param_6,int param_7
          ,undefined4 param_8,int param_9,int param_10)

{
                    /* 0x1dd60  56  VT_TextToPreprocessInfoFile_ENG */
  FUN_1001ef20(param_1,param_2,param_3,param_4,param_5,param_6,param_7,param_8,param_9,param_10);
  return;
}



/* ===== VT_TextToPcmBuffer_ProgressBar_ENG @ 1001dda0 ===== */

void __cdecl
VT_TextToPcmBuffer_ProgressBar_ENG
          (HWND param_1,UINT param_2,byte *param_3,undefined4 *param_4,int *param_5,int param_6,
          int param_7,int *param_8,int param_9,int param_10,int param_11,undefined4 param_12,
          int param_13,int param_14)

{
                    /* 0x1dda0  55  VT_TextToPcmBuffer_ProgressBar_ENG */
  FUN_10021230(param_1,param_2,param_3,param_4,param_5,param_6,param_7,param_8,param_9,param_10,
               param_11,param_12,param_13,param_14);
  return;
}



/* ===== VT_TextToBufferEX_ENG @ 1001ddf0 ===== */

int __cdecl
VT_TextToBufferEX_ENG
          (int param_1,byte *param_2,undefined4 *param_3,int *param_4,int param_5,int param_6,
          uint param_7,undefined4 *param_8,undefined4 *param_9,int param_10,int param_11,
          int param_12,undefined4 param_13,int *param_14,int *param_15)

{
  int iVar1;

                    /* 0x1ddf0  51  VT_TextToBufferEX_ENG */
  if ((-1 < param_1) && (param_1 < 3)) {
    iVar1 = FUN_10021bb0(param_1,param_2,param_3,param_4,param_5,param_6,param_7,param_8,param_9,
                         param_10,param_11,param_12,param_13,param_14,param_15);
    return iVar1;
  }
  return -1;
}



/* ===== VT_TextToLipSyncLog_ENG @ 1001de50 ===== */

void __cdecl
VT_TextToLipSyncLog_ENG
          (byte *param_1,byte *param_2,int param_3,int param_4,int param_5,int param_6,
          undefined4 param_7,int param_8,int param_9)

{
                    /* 0x1de50  54  VT_TextToLipSyncLog_ENG */
  FUN_10022200(param_1,param_2,param_3,param_4,param_5,param_6,param_7,param_8,param_9);
  return;
}



/* ===== VTDTTS_MakeInfo_ENG @ 1001de90 ===== */

void __cdecl
VTDTTS_MakeInfo_ENG(byte *param_1,undefined4 param_2,int param_3,int param_4,int param_5,int param_6
                   ,undefined4 param_7,int param_8,int param_9)

{
                    /* 0x1de90  1  VTDTTS_MakeInfo_ENG */
  FUN_10022350(param_1,param_2,param_3,param_4,param_5,param_6,param_7,param_8,param_9);
  return;
}



/* ===== VT_VerifyTTS_ENG @ 1001ded0 ===== */

void __cdecl VT_VerifyTTS_ENG(byte *param_1,int param_2,int param_3,int param_4)

{
                    /* 0x1ded0  61  VT_VerifyTTS_ENG */
  FUN_100226f0(param_1,param_2,param_3,param_4);
  return;
}



/* ===== VT_GetDefUserDictAbsoluteName_ENG @ 1001def0 ===== */

void VT_GetDefUserDictAbsoluteName_ENG(void)

{
                    /* 0x1def0  19  VT_GetDefUserDictAbsoluteName_ENG */
  FUN_100227f0(0);
  return;
}



/* ===== VT_GetDefUserDictRelativeName_ENG @ 1001df00 ===== */

void VT_GetDefUserDictRelativeName_ENG(void)

{
                    /* 0x1df00  20  VT_GetDefUserDictRelativeName_ENG */
  FUN_100227f0(1);
  return;
}



/* ===== VT_AllocSyncInfo_New_ENG @ 10026360 ===== */

int * VT_AllocSyncInfo_New_ENG(void)

{
  int *piVar1;
  int iVar2;
  undefined4 uVar3;
  int iVar4;

                    /* 0x26360  2  VT_AllocSyncInfo_New_ENG */
  piVar1 = (int *)FUN_1001d9c0(0x38);
  if (piVar1 == (int *)0x0) {
    return (int *)0x0;
  }
  piVar1[1] = 600;
  piVar1[2] = 0x41;
  iVar2 = FUN_1001d9c0(0x5460);
  *piVar1 = iVar2;
  if (iVar2 == 0) {
LAB_10026398:
    VT_FreeSyncInfo_New_ENG(piVar1);
    return (int *)0x0;
  }
  iVar2 = 0;
  if (0 < piVar1[1]) {
    iVar4 = 0;
    do {
      uVar3 = FUN_1001d9c0(piVar1[2] << 3);
      *(undefined4 *)(iVar4 + 4 + *piVar1) = uVar3;
      if (*(int *)(iVar4 + 4 + *piVar1) == 0) goto LAB_10026398;
      iVar2 = iVar2 + 1;
      iVar4 = iVar4 + 0x24;
    } while (iVar2 < piVar1[1]);
  }
  VT_InitSyncInfo_New_ENG(piVar1);
  return piVar1;
}



/* ===== VT_FreeSyncInfo_New_ENG @ 100263f0 ===== */

void __cdecl VT_FreeSyncInfo_New_ENG(int *param_1)

{
  undefined *puVar1;
  int iVar2;
  int iVar3;

                    /* 0x263f0  16  VT_FreeSyncInfo_New_ENG */
  if (param_1 != (int *)0x0) {
    if (*param_1 != 0) {
      iVar2 = 0;
      if (0 < param_1[1]) {
        iVar3 = 0;
        do {
          puVar1 = *(undefined **)(iVar3 + 4 + *param_1);
          if (puVar1 != (undefined *)0x0) {
            FUN_1001da30(puVar1);
          }
          iVar2 = iVar2 + 1;
          iVar3 = iVar3 + 0x24;
        } while (iVar2 < param_1[1]);
      }
      FUN_1001da30((undefined *)*param_1);
    }
    FUN_1001da30((undefined *)param_1);
  }
  return;
}



/* ===== VT_InitSyncInfo_New_ENG @ 10026450 ===== */

void __cdecl VT_InitSyncInfo_New_ENG(int *param_1)

{
  int iVar1;
  int iVar2;
  int iVar3;

                    /* 0x26450  30  VT_InitSyncInfo_New_ENG */
  if (param_1 != (int *)0x0) {
    iVar3 = 0;
    param_1[3] = 0;
    param_1[1] = 600;
    param_1[2] = 0x41;
    iVar1 = 0;
    do {
      *(undefined2 *)(iVar1 + *param_1) = 0;
      *(undefined4 *)(iVar1 + 8 + *param_1) = 0;
      *(undefined4 *)(iVar1 + 0xc + *param_1) = 0;
      *(undefined4 *)(iVar1 + 0x10 + *param_1) = 0;
      *(undefined4 *)(iVar1 + 0x14 + *param_1) = 0;
      *(undefined4 *)(iVar1 + 0x18 + *param_1) = 0;
      *(undefined4 *)(iVar1 + 0x1c + *param_1) = 0;
      *(undefined4 *)(iVar1 + 0x20 + *param_1) = 0;
      iVar2 = 0;
      if (0 < param_1[2]) {
        do {
          iVar2 = iVar2 + 1;
          *(undefined2 *)(*(int *)(iVar1 + 4 + *param_1) + -4 + iVar2 * 8) = 0;
          *(undefined4 *)(*(int *)(iVar1 + 4 + *param_1) + -8 + iVar2 * 8) = 0;
        } while (iVar2 < param_1[2]);
      }
      iVar3 = iVar3 + 1;
      iVar1 = iVar1 + 0x24;
    } while (iVar3 < param_1[1]);
    param_1[4] = -1;
    param_1[5] = 0;
    param_1[7] = 0;
    param_1[6] = -1;
    param_1[8] = -1;
    param_1[9] = -1;
    param_1[10] = 0;
    param_1[0xc] = 0;
    param_1[0xb] = -1;
    param_1[0xd] = -1;
  }
  return;
}



/* ===== VT_CopySyncInfo_New_ENG @ 10026500 ===== */

void __cdecl VT_CopySyncInfo_New_ENG(int *param_1,int *param_2)

{
  int *piVar1;
  int iVar2;
  int iVar3;

                    /* 0x26500  7  VT_CopySyncInfo_New_ENG */
  piVar1 = param_1;
  iVar2 = 0;
  if ((param_1 != (int *)0x0) && (param_2 != (int *)0x0)) {
    param_1[3] = param_2[3];
    param_1[1] = param_2[1];
    param_1[2] = param_2[2];
    param_1 = (int *)0x0;
    if (0 < param_2[1]) {
      do {
        *(undefined2 *)(iVar2 + *piVar1) = *(undefined2 *)(iVar2 + *param_2);
        *(undefined4 *)(iVar2 + 8 + *piVar1) = *(undefined4 *)(iVar2 + 8 + *param_2);
        *(undefined4 *)(iVar2 + 0xc + *piVar1) = *(undefined4 *)(iVar2 + 0xc + *param_2);
        *(undefined4 *)(iVar2 + 0x10 + *piVar1) = *(undefined4 *)(iVar2 + 0x10 + *param_2);
        *(undefined4 *)(iVar2 + 0x14 + *piVar1) = *(undefined4 *)(iVar2 + 0x14 + *param_2);
        *(undefined4 *)(iVar2 + 0x18 + *piVar1) = *(undefined4 *)(iVar2 + 0x18 + *param_2);
        *(undefined4 *)(iVar2 + 0x1c + *piVar1) = *(undefined4 *)(iVar2 + 0x1c + *param_2);
        *(undefined4 *)(iVar2 + 0x20 + *piVar1) = *(undefined4 *)(iVar2 + 0x20 + *param_2);
        iVar3 = 0;
        if (0 < param_2[2]) {
          do {
            iVar3 = iVar3 + 1;
            *(undefined2 *)(*(int *)(iVar2 + 4 + *piVar1) + -4 + iVar3 * 8) =
                 *(undefined2 *)(*(int *)(iVar2 + 4 + *param_2) + -4 + iVar3 * 8);
            *(undefined4 *)(*(int *)(iVar2 + 4 + *piVar1) + -8 + iVar3 * 8) =
                 *(undefined4 *)(*(int *)(iVar2 + 4 + *param_2) + -8 + iVar3 * 8);
          } while (iVar3 < param_2[2]);
        }
        param_1 = (int *)((int)param_1 + 1);
        iVar2 = iVar2 + 0x24;
      } while ((int)param_1 < param_2[1]);
    }
    piVar1[4] = param_2[4];
    piVar1[5] = param_2[5];
    piVar1[7] = param_2[7];
    piVar1[6] = param_2[6];
    piVar1[8] = param_2[8];
    piVar1[9] = param_2[9];
    piVar1[10] = param_2[10];
    piVar1[0xc] = param_2[0xc];
    piVar1[0xb] = param_2[0xb];
    piVar1[0xd] = param_2[0xd];
  }
  return;
}



/* ===== VT_SetSoundCardID_ENG @ 10026b20 ===== */

undefined4 __cdecl VT_SetSoundCardID_ENG(undefined4 param_1)

{
  undefined4 uVar1;

                    /* 0x26b20  46  VT_SetSoundCardID_ENG */
  uVar1 = DAT_1007d6e0;
  DAT_1007d6e0 = param_1;
  return uVar1;
}



/* ===== VT_PAUSETTS_ENG @ 10027010 ===== */

void VT_PAUSETTS_ENG(void)

{
                    /* 0x27010  35  VT_PAUSETTS_ENG */
  if (DAT_100a7490 != (HWAVEOUT)0x0) {
    waveOutPause(DAT_100a7490);
  }
  return;
}



/* ===== VT_RESTARTTTS_ENG @ 10027030 ===== */

void VT_RESTARTTTS_ENG(void)

{
                    /* 0x27030  37  VT_RESTARTTTS_ENG */
  if (DAT_100a7490 != (HWAVEOUT)0x0) {
    waveOutRestart(DAT_100a7490);
  }
  return;
}



/* ===== VT_PLAYTTS_ENG @ 10027050 ===== */

undefined4 __cdecl
VT_PLAYTTS_ENG(HWND param_1,UINT param_2,byte *param_3,int param_4,int param_5,int param_6,
              int param_7,undefined4 param_8,int param_9,int param_10)

{
  int *piVar1;
  int iVar2;
  undefined4 uVar3;
  char *pcVar4;
  undefined4 *puVar5;
  BOOL BVar6;
  int iVar7;
  short local_6;

                    /* 0x27050  36  VT_PLAYTTS_ENG */
  VT_SetDecimal0Pron_ENG();
  DAT_100a9d3c = param_1;
  DAT_100a9d34 = param_2;
  VT_STOPTTS_ENG();
  DAT_100a841c = 0;
  iVar2 = FUN_10025e60();
  if (iVar2 != 0) {
    return CONCAT22((short)((uint)iVar2 >> 0x10),0xffff);
  }
  if ((param_4 < 0) || (DAT_100a9d30 = param_4, 5 < param_4)) {
    DAT_100a9d30 = 1;
  }
  piVar1 = (int *)(&DAT_100a0464)[DAT_100a9d30];
  if (piVar1 == (int *)0x0) {
    uVar3 = FUN_10025e60();
    return CONCAT22((short)((uint)uVar3 >> 0x10),0xfffc);
  }
  if (param_3 == (byte *)0x0) {
    uVar3 = FUN_10025e60();
    return CONCAT22((short)((uint)uVar3 >> 0x10),0xfffe);
  }
  if (*param_3 == 0) {
    uVar3 = FUN_10025e60();
    return CONCAT22((short)((uint)uVar3 >> 0x10),0xfffd);
  }
  if ((param_4 < 0) || (iVar2 = param_4, 5 < param_4)) {
    iVar2 = 1;
  }
  iVar2 = FUN_10025fc0(&local_6,param_9,iVar2);
  piVar1[0xa4] = iVar2;
  if (local_6 != 1) {
    uVar3 = FUN_10025e60();
    return CONCAT22((short)((uint)uVar3 >> 0x10),0xffff);
  }
  iVar7 = 1;
  if ((-1 < param_4) && (param_4 < 6)) {
    iVar7 = param_4;
  }
  *(char *)(iVar2 + 0x47784) = (char)iVar7;
  FUN_100286c0(piVar1,piVar1[0xa4],param_5,param_6,param_7,param_8,param_10);
  pcVar4 = (char *)FUN_1001c990((byte *)piVar1[0xa4],param_3,0);
  if (pcVar4 == (char *)0x0) {
    FUN_10027750((undefined *)piVar1[0xa4]);
    uVar3 = FUN_10025e60();
    return CONCAT22((short)((uint)uVar3 >> 0x10),0xfffe);
  }
  puVar5 = &DAT_100a9d60;
  for (iVar2 = 600; iVar2 != 0; iVar2 = iVar2 + -1) {
    *puVar5 = 0;
    puVar5 = puVar5 + 1;
  }
  puVar5 = &DAT_100a842c;
  do {
    *puVar5 = 0;
    puVar5 = puVar5 + 8;
  } while ((int)puVar5 < 0x100a9d2c);
  DAT_100a7a84 = 0;
  DAT_100a9d44 = 0;
  DAT_100aa6c0 = 0;
  DAT_100a9d48 = 0;
  DAT_100a9d40 = 1;
  uVar3 = FUN_10026b40(piVar1,(int *)piVar1[0xa4],pcVar4);
  if ((short)uVar3 < 0) {
    DAT_100a7498 = 0;
    FUN_10025e60();
    VT_STOPTTS_ENG();
    BVar6 = PostMessageA(DAT_100a9d3c,DAT_100a9d34,0,-1);
    return CONCAT22((short)((uint)BVar6 >> 0x10),0xfffb);
  }
  uVar3 = FUN_10025e60();
  DAT_100a7498 = 1;
  return CONCAT22((short)((uint)uVar3 >> 0x10),1);
}



/* ===== VT_STOPTTS_ENG @ 10027560 ===== */

void VT_STOPTTS_ENG(void)

{
  int iVar1;
  int iVar2;
  HWAVEOUT hwo;
  int *piVar3;
  tagMSG local_24;
  int local_8;

                    /* 0x27560  38  VT_STOPTTS_ENG */
  DAT_100a841c = 1;
  iVar1 = (&DAT_100a0464)[DAT_100a9d30];
  local_8 = iVar1;
  iVar2 = FUN_10025e60();
  if (iVar2 == 0) {
    if (DAT_100a7490 != (HWAVEOUT)0x0) {
      waveOutReset(DAT_100a7490);
      piVar3 = &DAT_100a842c;
      hwo = DAT_100a7490;
      do {
        if (*piVar3 == 1) {
          if (hwo != (HWAVEOUT)0x0) {
            waveOutUnprepareHeader(hwo,(LPWAVEHDR)(piVar3 + -3),0x20);
            hwo = DAT_100a7490;
          }
          *piVar3 = 0;
        }
        piVar3 = piVar3 + 8;
      } while ((int)piVar3 < 0x100a90ac);
      piVar3 = &DAT_100a90ac;
      do {
        if (*piVar3 == 1) {
          if (hwo != (HWAVEOUT)0x0) {
            waveOutUnprepareHeader(hwo,(LPWAVEHDR)(piVar3 + -3),0x20);
            hwo = DAT_100a7490;
          }
          *piVar3 = 0;
        }
        piVar3 = piVar3 + 8;
      } while ((int)piVar3 < 0x100a9d2c);
      if (hwo != (HWAVEOUT)0x0) {
        waveOutClose(hwo);
      }
      DAT_100a7490 = (HWAVEOUT)0x0;
    }
    if ((iVar1 != 0) && (*(int *)(iVar1 + 0x290) != 0)) {
      iVar2 = PeekMessageA(&local_24,DAT_100a8418,0,0,1);
      while (iVar2 != 0) {
        TranslateMessage(&local_24);
        DispatchMessageA(&local_24);
        iVar2 = PeekMessageA(&local_24,DAT_100a8418,0,0,1);
        iVar1 = local_8;
      }
      FUN_10027750(*(undefined **)(iVar1 + 0x290));
      *(undefined4 *)(iVar1 + 0x290) = 0;
    }
    FUN_10025e60();
  }
  return;
}



/* ===== VT_LOAD_UserDict_EXT_ENG @ 10027880 ===== */

undefined4 __cdecl
VT_LOAD_UserDict_EXT_ENG(int param_1,int param_2,char *param_3,undefined *param_4)

{
  undefined4 in_EAX;
  undefined2 extraout_var;
  int *piVar1;
  undefined2 extraout_var_00;
  undefined2 extraout_var_01;

                    /* 0x27880  34  VT_LOAD_UserDict_EXT_ENG */
  if ((0x3ff < param_1) || (param_1 < 0)) {
    return CONCAT22((short)((uint)in_EAX >> 0x10),0xffff);
  }
  if (DAT_100a0458 == 0) {
    return 0xfffc;
  }
  EnterCriticalSection((LPCRITICAL_SECTION)&DAT_100a8400);
  if ((&DAT_100a647c)[param_1] != 0) {
    LeaveCriticalSection((LPCRITICAL_SECTION)&DAT_100a8400);
    return CONCAT22(extraout_var,0xfffe);
  }
  (&DAT_100a647c)[param_1] = 0xffffff01;
  LeaveCriticalSection((LPCRITICAL_SECTION)&DAT_100a8400);
  piVar1 = FUN_1005e330(param_2,param_3,param_4);
  (&DAT_100a647c)[param_1] = piVar1;
  EnterCriticalSection((LPCRITICAL_SECTION)&DAT_100a8400);
  if ((&DAT_100a647c)[param_1] == -0xff) {
    (&DAT_100a647c)[param_1] = 0;
  }
  if ((&DAT_100a647c)[param_1] == 0) {
    LeaveCriticalSection((LPCRITICAL_SECTION)&DAT_100a8400);
    return CONCAT22(extraout_var_00,0xfffd);
  }
  LeaveCriticalSection((LPCRITICAL_SECTION)&DAT_100a8400);
  return CONCAT22(extraout_var_01,1);
}



/* ===== VT_LOAD_UserDict_ENG @ 10027960 ===== */

void __cdecl VT_LOAD_UserDict_ENG(int param_1,int param_2)

{
                    /* 0x27960  33  VT_LOAD_UserDict_ENG */
  VT_LOAD_UserDict_EXT_ENG(param_1,param_2,(char *)0x0,(undefined *)0xffffffff);
  return;
}



/* ===== VT_UNLOAD_UserDict_EXT_ENG @ 10027980 ===== */

undefined4 __cdecl VT_UNLOAD_UserDict_EXT_ENG(int param_1)

{
  int *piVar1;
  int iVar2;
  undefined4 in_EAX;
  undefined2 extraout_var;
  undefined2 extraout_var_00;
  undefined2 extraout_var_01;
  int iVar3;
  int *piVar4;
  int *piVar5;
  int *local_8;

                    /* 0x27980  60  VT_UNLOAD_UserDict_EXT_ENG */
  if ((0x3ff < param_1) || (param_1 < 0)) {
    return CONCAT22((short)((uint)in_EAX >> 0x10),0xfffe);
  }
  if (DAT_100a0458 == 0) {
    return 0xfffd;
  }
  EnterCriticalSection((LPCRITICAL_SECTION)&DAT_100a8400);
  piVar1 = (int *)(&DAT_100a647c)[param_1];
  if (piVar1 != (int *)0x0) {
    local_8 = &DAT_100a0464;
    piVar5 = &DAT_100a047c;
    do {
      if (*local_8 != 0) {
        iVar2 = *(int *)(*local_8 + 0x4d14);
        iVar3 = 0;
        piVar4 = piVar5;
        if (0 < iVar2) {
          do {
            if ((*piVar4 != 0) && (*(int **)(*piVar4 + 0x1312c0) == piVar1)) {
              LeaveCriticalSection((LPCRITICAL_SECTION)&DAT_100a8400);
              return CONCAT22(extraout_var_01,0xfffd);
            }
            iVar3 = iVar3 + 1;
            piVar4 = piVar4 + 1;
          } while (iVar3 < iVar2);
        }
      }
      piVar5 = piVar5 + 0x400;
      local_8 = local_8 + 1;
      if (0x100a647b < (int)piVar5) {
        FUN_1005ef40(piVar1);
        (&DAT_100a647c)[param_1] = 0;
        LeaveCriticalSection((LPCRITICAL_SECTION)&DAT_100a8400);
        return CONCAT22(extraout_var_00,1);
      }
    } while( true );
  }
  LeaveCriticalSection((LPCRITICAL_SECTION)&DAT_100a8400);
  return CONCAT22(extraout_var,0xffff);
}



/* ===== VT_UNLOAD_UserDict_ENG @ 10027a80 ===== */

void __cdecl VT_UNLOAD_UserDict_ENG(int param_1)

{
                    /* 0x27a80  59  VT_UNLOAD_UserDict_ENG */
  VT_UNLOAD_UserDict_EXT_ENG(param_1);
  return;
}



/* ===== VT_LOADTTS_EXT_ENG @ 10027af0 ===== */

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



/* ===== VT_LOADTTS_ENG @ 10027e70 ===== */

void __cdecl VT_LOADTTS_ENG(HWND param_1,int param_2,char *param_3,char *param_4)

{
                    /* 0x27e70  31  VT_LOADTTS_ENG */
  VT_LOADTTS_EXT_ENG(param_1,param_2,param_3,0,0xffffffff,param_4,(char *)0x0,0xffffffff);
  return;
}



/* ===== VT_UNLOADTTS_EXT_ENG @ 10027ea0 ===== */

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



/* ===== VT_UNLOADTTS_ENG @ 10027f80 ===== */

void __cdecl VT_UNLOADTTS_ENG(int param_1)

{
                    /* 0x27f80  57  VT_UNLOADTTS_ENG */
  VT_UNLOADTTS_EXT_ENG(param_1);
  return;
}



/* ===== VT_DestroyWindow_ENG @ 10027fd0 ===== */

void VT_DestroyWindow_ENG(void)

{
                    /* 0x27fd0  15  VT_DestroyWindow_ENG */
  DestroyWindow(DAT_100a8418);
  return;
}



/* ===== VT_SetPitchSpeedVolumePause_ENG @ 10027fe0 ===== */

void __cdecl
VT_SetPitchSpeedVolumePause_ENG(int param_1,int param_2,int param_3,int param_4,int param_5)

{
  int iVar1;

                    /* 0x27fe0  45  VT_SetPitchSpeedVolumePause_ENG */
  if ((param_5 < 0) || (5 < param_5)) {
    param_5 = 1;
  }
  iVar1 = (&DAT_100a0464)[param_5];
  if (iVar1 != 0) {
    if (-1 < param_1) {
      *(int *)(iVar1 + 0x4cf4) = param_1;
      if (200 < param_1) {
        *(undefined4 *)(iVar1 + 0x4cf4) = 200;
      }
      if (*(int *)(iVar1 + 0x4cf4) < 0x32) {
        *(undefined4 *)(iVar1 + 0x4cf4) = 0x32;
      }
    }
    if (-1 < param_2) {
      *(int *)(iVar1 + 0x4cf0) = param_2;
      if (400 < param_2) {
        *(undefined4 *)(iVar1 + 0x4cf0) = 400;
      }
      if (*(int *)(iVar1 + 0x4cf0) < 0x32) {
        *(undefined4 *)(iVar1 + 0x4cf0) = 0x32;
      }
    }
    if (-1 < param_3) {
      *(int *)(iVar1 + 0x4cf8) = param_3;
      if (500 < param_3) {
        *(undefined4 *)(iVar1 + 0x4cf8) = 500;
      }
      if (*(int *)(iVar1 + 0x4cf8) < 0) {
        *(undefined4 *)(iVar1 + 0x4cf8) = 0;
      }
    }
    if ((-1 < param_4) && (*(int *)(iVar1 + 0x4d00) = param_4, 0xffff < param_4)) {
      *(undefined4 *)(iVar1 + 0x4d00) = 0xffff;
    }
  }
  return;
}



/* ===== VT_GetPitchSpeedVolumePause_ENG @ 100280c0 ===== */

undefined4 __cdecl
VT_GetPitchSpeedVolumePause_ENG
          (undefined4 *param_1,undefined4 *param_2,undefined4 *param_3,undefined4 *param_4,
          int param_5)

{
  int iVar1;

                    /* 0x280c0  25  VT_GetPitchSpeedVolumePause_ENG */
  if ((param_5 < 0) || (5 < param_5)) {
    param_5 = 1;
  }
  iVar1 = (&DAT_100a0464)[param_5];
  if (iVar1 == 0) {
    return 0xffffffff;
  }
  if (param_1 != (undefined4 *)0x0) {
    *param_1 = *(undefined4 *)(iVar1 + 0x4cf4);
  }
  if (param_2 != (undefined4 *)0x0) {
    *param_2 = *(undefined4 *)(iVar1 + 0x4cf0);
  }
  if (param_3 != (undefined4 *)0x0) {
    *param_3 = *(undefined4 *)(iVar1 + 0x4cf8);
  }
  if (param_4 != (undefined4 *)0x0) {
    *param_4 = *(undefined4 *)(iVar1 + 0x4d00);
  }
  return 1;
}



/* ===== VT_GetDBSize_ENG @ 10028130 ===== */

undefined4 __cdecl VT_GetDBSize_ENG(int *param_1,int param_2)

{
  int iVar1;

                    /* 0x28130  18  VT_GetDBSize_ENG */
  if ((param_2 < 0) || (iVar1 = param_2, 5 < param_2)) {
    iVar1 = 1;
  }
  if ((&DAT_100a0464)[iVar1] == 0) {
    return 0xffffffff;
  }
  if (DAT_100a747c == '\0') {
    return 0xffffffff;
  }
  if ((param_2 < 0) || (iVar1 = param_2, 5 < param_2)) {
    iVar1 = 1;
  }
  if ((&DAT_100a7480)[iVar1] == '\0') {
    return 0xffffffff;
  }
  if ((param_2 < 0) || (5 < param_2)) {
    param_2 = 1;
  }
  *param_1 = *(int *)(&DAT_100a74a0 + param_2 * 4) + DAT_100a749c;
  return 1;
}



/* ===== VT_SetCommaPause_ENG @ 100281b0 ===== */

void __cdecl VT_SetCommaPause_ENG(int param_1,int param_2)

{
  int iVar1;

                    /* 0x281b0  39  VT_SetCommaPause_ENG */
  if ((param_2 < 0) || (5 < param_2)) {
    param_2 = 1;
  }
  iVar1 = (&DAT_100a0464)[param_2];
  if (((iVar1 != 0) && (-1 < param_1)) && (*(int *)(iVar1 + 0x4d04) = param_1, 0xffff < param_1)) {
    *(undefined4 *)(iVar1 + 0x4d04) = 0xffff;
  }
  return;
}



/* ===== VT_GetCommaPause_ENG @ 100281f0 ===== */

undefined4 __cdecl VT_GetCommaPause_ENG(undefined4 *param_1,int param_2)

{
                    /* 0x281f0  17  VT_GetCommaPause_ENG */
  if ((param_2 < 0) || (5 < param_2)) {
    param_2 = 1;
  }
  if ((&DAT_100a0464)[param_2] == 0) {
    return 0xffffffff;
  }
  if (param_1 != (undefined4 *)0x0) {
    *param_1 = *(undefined4 *)((&DAT_100a0464)[param_2] + 0x4d04);
  }
  return 1;
}



/* ===== VT_GetSpeakerName_ENG @ 10028230 ===== */

undefined * __cdecl VT_GetSpeakerName_ENG(int param_1)

{
                    /* 0x28230  26  VT_GetSpeakerName_ENG */
  if ((-1 < param_1) && (param_1 < 6)) {
    return (&PTR_DAT_1007c6b0)[param_1 * 6];
  }
  return PTR_DAT_1007c6c8;
}



/* ===== VT_GetDefVersion_ENG @ 10028260 ===== */

undefined4 * VT_GetDefVersion_ENG(void)

{
  char cVar1;
  uint uVar2;
  int iVar3;
  uint uVar4;
  char *pcVar5;
  char *pcVar6;
  char *pcVar7;

                    /* 0x28260  21  VT_GetDefVersion_ENG */
  FUN_10063ed6((undefined1 *)&DAT_100a0258,&DAT_1007bf98);
  uVar2 = 0xffffffff;
  pcVar5 = &DAT_1007d940;
  do {
    pcVar7 = pcVar5;
    if (uVar2 == 0) break;
    uVar2 = uVar2 - 1;
    pcVar7 = pcVar5 + 1;
    cVar1 = *pcVar5;
    pcVar5 = pcVar7;
  } while (cVar1 != '\0');
  uVar2 = ~uVar2;
  iVar3 = -1;
  pcVar5 = (char *)&DAT_100a0258;
  do {
    pcVar6 = pcVar5;
    if (iVar3 == 0) break;
    iVar3 = iVar3 + -1;
    pcVar6 = pcVar5 + 1;
    cVar1 = *pcVar5;
    pcVar5 = pcVar6;
  } while (cVar1 != '\0');
  pcVar5 = pcVar7 + -uVar2;
  pcVar7 = pcVar6 + -1;
  for (uVar4 = uVar2 >> 2; uVar4 != 0; uVar4 = uVar4 - 1) {
    *(undefined4 *)pcVar7 = *(undefined4 *)pcVar5;
    pcVar5 = pcVar5 + 4;
    pcVar7 = pcVar7 + 4;
  }
  for (uVar2 = uVar2 & 3; uVar2 != 0; uVar2 = uVar2 - 1) {
    *pcVar7 = *pcVar5;
    pcVar5 = pcVar5 + 1;
    pcVar7 = pcVar7 + 1;
  }
  uVar2 = 0xffffffff;
  pcVar5 = &DAT_1007c930;
  do {
    pcVar7 = pcVar5;
    if (uVar2 == 0) break;
    uVar2 = uVar2 - 1;
    pcVar7 = pcVar5 + 1;
    cVar1 = *pcVar5;
    pcVar5 = pcVar7;
  } while (cVar1 != '\0');
  uVar2 = ~uVar2;
  iVar3 = -1;
  pcVar5 = (char *)&DAT_100a0258;
  do {
    pcVar6 = pcVar5;
    if (iVar3 == 0) break;
    iVar3 = iVar3 + -1;
    pcVar6 = pcVar5 + 1;
    cVar1 = *pcVar5;
    pcVar5 = pcVar6;
  } while (cVar1 != '\0');
  pcVar5 = pcVar7 + -uVar2;
  pcVar7 = pcVar6 + -1;
  for (uVar4 = uVar2 >> 2; uVar4 != 0; uVar4 = uVar4 - 1) {
    *(undefined4 *)pcVar7 = *(undefined4 *)pcVar5;
    pcVar5 = pcVar5 + 4;
    pcVar7 = pcVar7 + 4;
  }
  for (uVar2 = uVar2 & 3; uVar2 != 0; uVar2 = uVar2 - 1) {
    *pcVar7 = *pcVar5;
    pcVar5 = pcVar5 + 1;
    pcVar7 = pcVar7 + 1;
  }
  uVar2 = 0xffffffff;
  pcVar5 = s__FileIO_1007d938;
  do {
    pcVar7 = pcVar5;
    if (uVar2 == 0) break;
    uVar2 = uVar2 - 1;
    pcVar7 = pcVar5 + 1;
    cVar1 = *pcVar5;
    pcVar5 = pcVar7;
  } while (cVar1 != '\0');
  uVar2 = ~uVar2;
  iVar3 = -1;
  pcVar5 = (char *)&DAT_100a0258;
  do {
    pcVar6 = pcVar5;
    if (iVar3 == 0) break;
    iVar3 = iVar3 + -1;
    pcVar6 = pcVar5 + 1;
    cVar1 = *pcVar5;
    pcVar5 = pcVar6;
  } while (cVar1 != '\0');
  pcVar5 = pcVar7 + -uVar2;
  pcVar7 = pcVar6 + -1;
  for (uVar4 = uVar2 >> 2; uVar4 != 0; uVar4 = uVar4 - 1) {
    *(undefined4 *)pcVar7 = *(undefined4 *)pcVar5;
    pcVar5 = pcVar5 + 4;
    pcVar7 = pcVar7 + 4;
  }
  for (uVar2 = uVar2 & 3; uVar2 != 0; uVar2 = uVar2 - 1) {
    *pcVar7 = *pcVar5;
    pcVar5 = pcVar5 + 1;
    pcVar7 = pcVar7 + 1;
  }
  return &DAT_100a0258;
}



/* ===== VT_SetEmphasisFactor_ENG @ 10028310 ===== */

void __cdecl VT_SetEmphasisFactor_ENG(int param_1,int param_2)

{
  int iVar1;

                    /* 0x28310  41  VT_SetEmphasisFactor_ENG */
  if ((param_2 < 0) || (5 < param_2)) {
    param_2 = 1;
  }
  iVar1 = (&DAT_100a0464)[param_2];
  if (iVar1 != 0) {
    *(int *)(iVar1 + 0x4cfc) = param_1;
    if (param_1 < 0) {
      if (param_1 < -0x5f) {
        *(undefined4 *)(iVar1 + 0x4cfc) = 0xffffffa1;
      }
    }
    else if (0x5f < param_1) {
      *(undefined4 *)(iVar1 + 0x4cfc) = 0x5f;
      return;
    }
  }
  return;
}



/* ===== VT_SetTextTypeForHighlight_ENG @ 10028360 ===== */

void __cdecl VT_SetTextTypeForHighlight_ENG(undefined1 param_1)

{
                    /* 0x28360  47  VT_SetTextTypeForHighlight_ENG */
  if (DAT_100a0460 != 0) {
    *(undefined1 *)(DAT_100a0460 + 0x20424) = param_1;
    if (*(char *)(DAT_100a0460 + 0x20424) != '\0') {
      *(undefined1 *)(DAT_100a0460 + 0x20424) = 1;
    }
  }
  return;
}



/* ===== VT_SetParenthesisCharNumber_ENG @ 10028390 ===== */

void __cdecl VT_SetParenthesisCharNumber_ENG(undefined4 param_1)

{
                    /* 0x28390  43  VT_SetParenthesisCharNumber_ENG */
  if (DAT_100a0460 != 0) {
    *(undefined4 *)(DAT_100a0460 + 0x2041c) = param_1;
    if (*(int *)(DAT_100a0460 + 0x2041c) < 0) {
      *(undefined4 *)(DAT_100a0460 + 0x2041c) = 0;
    }
  }
  return;
}



/* ===== VT_SetEnglishReadingRule_KOR @ 100283c0 ===== */

void __cdecl VT_SetEnglishReadingRule_KOR(undefined4 param_1)

{
                    /* 0x283c0  42  VT_SetEnglishReadingRule_KOR */
  if (DAT_100a0460 != 0) {
    *(undefined4 *)(DAT_100a0460 + 0x20420) = param_1;
    if (*(int *)(DAT_100a0460 + 0x20420) < 0) {
      *(undefined4 *)(DAT_100a0460 + 0x20420) = 0;
    }
  }
  return;
}



/* ===== VT_SetUnitSelectHistoryMode_ENG @ 100283f0 ===== */

void __cdecl VT_SetUnitSelectHistoryMode_ENG(char param_1)

{
                    /* 0x283f0  48  VT_SetUnitSelectHistoryMode_ENG */
  if (DAT_100a747c == '\0') {
    DAT_100a045d = param_1;
    if (param_1 != '\x01') {
      DAT_100a045d = 0;
      return;
    }
    VT_SetDecimal0Pron_ENG();
  }
  return;
}



/* ===== VT_SetDecimal0Pron_ENG @ 10028420 ===== */

void VT_SetDecimal0Pron_ENG(void)

{
                    /* 0x28420  40  VT_SetDecimal0Pron_ENG
                       0x28420  44  VT_SetPhone0Pron_ENG
                       0x28420  49  VT_SetVirtualTagMode_ENG */
  return;
}



/* ===== VT_GetLicenseComment_ENG @ 10029a00 ===== */

uint __cdecl VT_GetLicenseComment_ENG(char *param_1,char *param_2,int param_3)

{
  char cVar1;
  uint uVar2;
  uint uVar3;
  uint uVar4;
  char *pcVar5;
  DWORD local_44 [14];
  char *local_c;

                    /* 0x29a00  22  VT_GetLicenseComment_ENG */
  FUN_10028ca0(local_44);
  FUN_10028db0(param_1,(char *)0x0,0xffffffff,local_44);
  if (local_44[0] == 0) {
    uVar4 = 0xffffffff;
    pcVar5 = local_c;
    do {
      if (uVar4 == 0) break;
      uVar4 = uVar4 - 1;
      cVar1 = *pcVar5;
      pcVar5 = pcVar5 + 1;
    } while (cVar1 != '\0');
    uVar4 = ~uVar4;
    if ((int)(uVar4 - 1) < 0) {
      uVar4 = 0xfffffffe;
    }
    else if (-1 < param_3) {
      if (local_c == (char *)0x0) {
        uVar4 = 0xfffffffd;
      }
      else if (param_3 < (int)uVar4) {
        uVar4 = 0xfffffffc;
      }
      else {
        uVar2 = 0xffffffff;
        do {
          pcVar5 = local_c;
          if (uVar2 == 0) break;
          uVar2 = uVar2 - 1;
          pcVar5 = local_c + 1;
          cVar1 = *local_c;
          local_c = pcVar5;
        } while (cVar1 != '\0');
        uVar2 = ~uVar2;
        pcVar5 = pcVar5 + -uVar2;
        for (uVar3 = uVar2 >> 2; uVar3 != 0; uVar3 = uVar3 - 1) {
          *(undefined4 *)param_2 = *(undefined4 *)pcVar5;
          pcVar5 = pcVar5 + 4;
          param_2 = param_2 + 4;
        }
        for (uVar2 = uVar2 & 3; uVar2 != 0; uVar2 = uVar2 - 1) {
          *param_2 = *pcVar5;
          pcVar5 = pcVar5 + 1;
          param_2 = param_2 + 1;
        }
      }
    }
  }
  else {
    uVar4 = 0xffffffff;
  }
  FUN_10028ce0((int)local_44);
  return uVar4;
}



/* ===== VT_CheckLicense_ENG @ 10029b80 ===== */

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



/* ===== VT_GetLicenseInfo_ENG @ 1002a110 ===== */

int __cdecl
VT_GetLicenseInfo_ENG
          (int param_1,char *param_2,char *param_3,uint param_4,byte *param_5,int param_6)

{
  byte bVar1;
  char cVar2;
  char *_Str;
  char *pcVar3;
  byte *pbVar4;
  undefined4 uVar5;
  uint uVar6;
  int iVar7;
  void *extraout_ECX;
  void *extraout_ECX_00;
  void *extraout_ECX_01;
  void *this;
  uint uVar8;
  void *this_00;
  byte *pbVar9;
  char *pcVar10;
  byte *pbVar11;
  DWORD local_58;
  int local_54;
  int local_50;
  int local_4c;
  byte *local_38;
  byte *local_30;
  byte *local_2c;
  byte *local_20;
  undefined2 local_18;
  undefined1 uStack_16;
  char cStack_15;
  char local_14 [2];
  char cStack_12;
  char cStack_11;
  char local_10;

                    /* 0x2a110  23  VT_GetLicenseInfo_ENG */
  FUN_10028ca0(&local_58);
  FUN_10028db0(param_2,param_3,param_4,&local_58);
  if (local_58 != 0) {
    if (local_54 < 0) {
      FUN_10028ce0((int)&local_58);
      return -3;
    }
    if ((-1 < local_50) && (-1 < local_4c)) {
      FUN_10028ce0((int)&local_58);
      return -2;
    }
  }
  switch(param_1) {
  case 0:
  case 0xc:
    iVar7 = 4;
    goto LAB_1002a1a2;
  case 1:
    pbVar4 = local_30;
    break;
  case 2:
    pbVar4 = local_38;
    break;
  case 3:
  case 4:
  case 5:
  case 6:
  case 7:
  case 8:
  case 9:
  case 10:
  case 0xb:
  case 0xd:
  case 0xe:
  case 0xf:
    pbVar4 = local_20;
    break;
  default:
    goto switchD_1002a196_default;
  }
  uVar6 = 0xffffffff;
  do {
    if (uVar6 == 0) break;
    uVar6 = uVar6 - 1;
    bVar1 = *pbVar4;
    pbVar4 = pbVar4 + 1;
  } while (bVar1 != 0);
  iVar7 = ~uVar6 - 1;
  if (iVar7 < 0) {
switchD_1002a196_default:
    FUN_10028ce0((int)&local_58);
    return -1;
  }
LAB_1002a1a2:
  if (param_6 < 0) {
    FUN_10028ce0((int)&local_58);
    return iVar7;
  }
  if ((param_1 == 0) || (param_1 == 0xc)) {
    if (param_6 < iVar7) {
      FUN_10028ce0((int)&local_58);
      return -5;
    }
  }
  else if (param_6 <= iVar7) {
    FUN_10028ce0((int)&local_58);
    return -5;
  }
  if (param_5 == (byte *)0x0) {
    FUN_10028ce0((int)&local_58);
    return -4;
  }
  _Str = _strstr((char *)local_20,s_<vw_verify_1007d974);
  pbVar4 = (byte *)0x0;
  this = extraout_ECX;
  if (_Str != (char *)0x0) {
    pcVar3 = _strstr(_Str,&DAT_1007d970);
    if (pcVar3 == (char *)0x0) {
      pbVar4 = (byte *)0x0;
      this = extraout_ECX_00;
    }
    else {
      uVar6 = 0xffffffff;
      pcVar10 = s_<vw_verify_1007d974;
      do {
        if (uVar6 == 0) break;
        uVar6 = uVar6 - 1;
        cVar2 = *pcVar10;
        pcVar10 = pcVar10 + 1;
      } while (cVar2 != '\0');
      pcVar3 = pcVar3 + (-(int)_Str - (~uVar6 - 1));
      pbVar4 = (byte *)FUN_1001d9c0((size_t)(pcVar3 + 1));
      FUN_10063f30((undefined4 *)pbVar4,(undefined4 *)(_Str + (~uVar6 - 1)),(uint)pcVar3);
      pbVar4[(int)pcVar3] = 0;
      this = extraout_ECX_01;
    }
  }
  if (param_1 == 0) {
    uVar5 = FUN_10064645(this,local_2c);
    *(undefined4 *)param_5 = uVar5;
    goto LAB_1002a512;
  }
  if (param_1 == 1) {
LAB_1002a2db:
    uVar6 = 0xffffffff;
    do {
      pbVar9 = local_30;
      if (uVar6 == 0) break;
      uVar6 = uVar6 - 1;
      pbVar9 = local_30 + 1;
      bVar1 = *local_30;
      local_30 = pbVar9;
    } while (bVar1 != 0);
    uVar6 = ~uVar6;
    pbVar9 = pbVar9 + -uVar6;
    for (uVar8 = uVar6 >> 2; uVar8 != 0; uVar8 = uVar8 - 1) {
      *(undefined4 *)param_5 = *(undefined4 *)pbVar9;
      pbVar9 = pbVar9 + 4;
      param_5 = param_5 + 4;
    }
    for (uVar6 = uVar6 & 3; uVar6 != 0; uVar6 = uVar6 - 1) {
      *param_5 = *pbVar9;
      pbVar9 = pbVar9 + 1;
      param_5 = param_5 + 1;
    }
  }
  else {
    if (param_1 == 2) {
      uVar6 = 0xffffffff;
      do {
        pbVar9 = local_38;
        if (uVar6 == 0) break;
        uVar6 = uVar6 - 1;
        pbVar9 = local_38 + 1;
        bVar1 = *local_38;
        local_38 = pbVar9;
      } while (bVar1 != 0);
      uVar6 = ~uVar6;
      pbVar9 = pbVar9 + -uVar6;
      for (uVar8 = uVar6 >> 2; uVar8 != 0; uVar8 = uVar8 - 1) {
        *(undefined4 *)param_5 = *(undefined4 *)pbVar9;
        pbVar9 = pbVar9 + 4;
        param_5 = param_5 + 4;
      }
      for (uVar6 = uVar6 & 3; uVar6 != 0; uVar6 = uVar6 - 1) {
        *param_5 = *pbVar9;
        pbVar9 = pbVar9 + 1;
        param_5 = param_5 + 1;
      }
      goto LAB_1002a512;
    }
    if (pbVar4 == (byte *)0x0) {
      if (param_1 == 0xc) {
        param_5[0] = 0xff;
        param_5[1] = 0xff;
        param_5[2] = 0xff;
        param_5[3] = 0xff;
        goto LAB_1002a512;
      }
    }
    else {
      local_30 = pbVar4;
      if (param_1 == 3) goto LAB_1002a2db;
      if (param_1 == 4) {
        uStack_16 = DAT_1007da52;
        local_18._0_1_ = (undefined1)DAT_1007da50;
        local_18._1_1_ = DAT_1007da50._1_1_;
      }
      else if (param_1 == 5) {
        local_18._0_1_ = (char)DAT_1007da48;
        local_18._1_1_ = DAT_1007da48._1_1_;
        uStack_16 = DAT_1007da48._2_1_;
        cStack_15 = DAT_1007da48._3_1_;
        local_14[0] = DAT_1007da4c;
      }
      else if (param_1 == 6) {
        local_18._0_1_ = s_speaker_1007da40[0];
        local_18._1_1_ = s_speaker_1007da40[1];
        uStack_16 = s_speaker_1007da40[2];
        cStack_15 = s_speaker_1007da40[3];
        local_14[0] = s_speaker_1007da40[4];
        local_14[1] = s_speaker_1007da40[5];
        cStack_12 = s_speaker_1007da40[6];
        cStack_11 = s_speaker_1007da40[7];
      }
      else if (param_1 == 7) {
        local_18._0_1_ = s_version_1007da38[0];
        local_18._1_1_ = s_version_1007da38[1];
        uStack_16 = s_version_1007da38[2];
        cStack_15 = s_version_1007da38[3];
        local_14[0] = s_version_1007da38[4];
        local_14[1] = s_version_1007da38[5];
        cStack_12 = s_version_1007da38[6];
        cStack_11 = s_version_1007da38[7];
      }
      else if (param_1 == 8) {
        local_18._0_1_ = s_dbaccess_1007da2c[0];
        local_18._1_1_ = s_dbaccess_1007da2c[1];
        uStack_16 = s_dbaccess_1007da2c[2];
        cStack_15 = s_dbaccess_1007da2c[3];
        local_14[0] = s_dbaccess_1007da2c[4];
        local_14[1] = s_dbaccess_1007da2c[5];
        cStack_12 = s_dbaccess_1007da2c[6];
        cStack_11 = s_dbaccess_1007da2c[7];
        local_10 = s_dbaccess_1007da2c[8];
      }
      else if (param_1 == 9) {
        local_18._0_1_ = s_sampling_1007da20[0];
        local_18._1_1_ = s_sampling_1007da20[1];
        uStack_16 = s_sampling_1007da20[2];
        cStack_15 = s_sampling_1007da20[3];
        local_14[0] = s_sampling_1007da20[4];
        local_14[1] = s_sampling_1007da20[5];
        cStack_12 = s_sampling_1007da20[6];
        cStack_11 = s_sampling_1007da20[7];
        local_10 = s_sampling_1007da20[8];
      }
      else if (param_1 == 10) {
        local_18._0_1_ = (char)DAT_1007da94;
        local_18._1_1_ = DAT_1007da94._1_1_;
        uStack_16 = DAT_1007da94._2_1_;
        cStack_15 = DAT_1007da94._3_1_;
      }
      else if (param_1 == 0xc) {
        local_18._0_1_ = s_savetime_1007da88[0];
        local_18._1_1_ = s_savetime_1007da88[1];
        uStack_16 = s_savetime_1007da88[2];
        cStack_15 = s_savetime_1007da88[3];
        local_14[0] = s_savetime_1007da88[4];
        local_14[1] = s_savetime_1007da88[5];
        cStack_12 = s_savetime_1007da88[6];
        cStack_11 = s_savetime_1007da88[7];
        local_10 = s_savetime_1007da88[8];
      }
      else if (param_1 == 0xb) {
        local_18._0_1_ = (char)DAT_1007da80;
        local_18._1_1_ = DAT_1007da80._1_1_;
        uStack_16 = DAT_1007da80._2_1_;
        cStack_15 = DAT_1007da80._3_1_;
        local_14[0] = (char)DAT_1007da84;
        local_14[1] = DAT_1007da84._1_1_;
        cStack_12 = DAT_1007da84._2_1_;
        cStack_11 = DAT_1007da84._3_1_;
      }
      else if (param_1 == 0xd) {
        local_18._0_1_ = (char)DAT_1007da78;
        local_18._1_1_ = DAT_1007da78._1_1_;
        uStack_16 = DAT_1007da78._2_1_;
        cStack_15 = DAT_1007da78._3_1_;
        local_14[0] = (char)DAT_1007da7c;
        local_14[1] = DAT_1007da7c._1_1_;
        cStack_12 = DAT_1007da7c._2_1_;
        cStack_11 = DAT_1007da7c._3_1_;
      }
      else if (param_1 == 0xe) {
        local_18._0_1_ = s_dbsize_1007da18[0];
        local_18._1_1_ = s_dbsize_1007da18[1];
        uStack_16 = s_dbsize_1007da18[2];
        cStack_15 = s_dbsize_1007da18[3];
        cStack_12 = s_dbsize_1007da18[6];
        local_14[0] = s_dbsize_1007da18[4];
        local_14[1] = s_dbsize_1007da18[5];
      }
      else if (param_1 == 0xf) {
        local_18._0_1_ = s_realtime_1007da6c[0];
        local_18._1_1_ = s_realtime_1007da6c[1];
        uStack_16 = s_realtime_1007da6c[2];
        cStack_15 = s_realtime_1007da6c[3];
        local_14[0] = s_realtime_1007da6c[4];
        local_14[1] = s_realtime_1007da6c[5];
        cStack_12 = s_realtime_1007da6c[6];
        cStack_11 = s_realtime_1007da6c[7];
        local_10 = s_realtime_1007da6c[8];
      }
      else {
        local_18._0_1_ = DAT_1009f948;
      }
      iVar7 = FUN_1002ebf0((char *)pbVar4,(char *)pbVar4,(byte *)&local_18,&DAT_1007d964);
      if (iVar7 != 0) {
        if (param_1 == 0xc) {
          uVar5 = FUN_10064645(this_00,pbVar4);
          *(undefined4 *)param_5 = uVar5;
        }
        else {
          uVar6 = 0xffffffff;
          pbVar9 = pbVar4;
          do {
            pbVar11 = pbVar9;
            if (uVar6 == 0) break;
            uVar6 = uVar6 - 1;
            pbVar11 = pbVar9 + 1;
            bVar1 = *pbVar9;
            pbVar9 = pbVar11;
          } while (bVar1 != 0);
          uVar6 = ~uVar6;
          pbVar9 = pbVar11 + -uVar6;
          for (uVar8 = uVar6 >> 2; uVar8 != 0; uVar8 = uVar8 - 1) {
            *(undefined4 *)param_5 = *(undefined4 *)pbVar9;
            pbVar9 = pbVar9 + 4;
            param_5 = param_5 + 4;
          }
          for (uVar6 = uVar6 & 3; uVar6 != 0; uVar6 = uVar6 - 1) {
            *param_5 = *pbVar9;
            pbVar9 = pbVar9 + 1;
            param_5 = param_5 + 1;
          }
        }
        goto LAB_1002a512;
      }
      if (param_1 == 0xc) {
        param_5[0] = 0xff;
        param_5[1] = 0xff;
        param_5[2] = 0xff;
        param_5[3] = 0xff;
        goto LAB_1002a512;
      }
    }
    *param_5 = DAT_1009f948;
  }
LAB_1002a512:
  FUN_1001da30(pbVar4);
  FUN_10028ce0((int)&local_58);
  return 0;
}



/* ===== VT_CheckUserDict_SourceNorm_ENG @ 1002a550 ===== */

void __cdecl VT_CheckUserDict_SourceNorm_ENG(byte *param_1)

{
  undefined1 local_38 [52];

                    /* 0x2a550  4  VT_CheckUserDict_SourceNorm_ENG */
  local_38[0] = 0;
  FUN_1005f2e0(local_38,param_1);
  return;
}



/* ===== VT_CheckUserDict_TargetNorm_ENG @ 1002a570 ===== */

int __cdecl VT_CheckUserDict_TargetNorm_ENG(byte *param_1)

{
  int iVar1;

                    /* 0x2a570  5  VT_CheckUserDict_TargetNorm_ENG */
  iVar1 = FUN_1005f3b0(param_1);
  return (int)(short)iVar1;
}



/* ===== VT_CheckUserDict_TargetPhon_ENG @ 1002a590 ===== */

int __cdecl VT_CheckUserDict_TargetPhon_ENG(byte *param_1)

{
  int iVar1;

                    /* 0x2a590  6  VT_CheckUserDict_TargetPhon_ENG */
  iVar1 = FUN_1005f5c0(param_1);
  return (int)(short)iVar1;
}



/* ===== VT_GetUserDictLimit_ENG @ 1002a5b0 ===== */

undefined4 __cdecl VT_GetUserDictLimit_ENG(undefined4 param_1)

{
  undefined4 uVar1;

                    /* 0x2a5b0  28  VT_GetUserDictLimit_ENG */
  uVar1 = 0xffffffff;
  switch(param_1) {
  case 0:
    uVar1 = 0x1e;
    break;
  case 1:
    return 10;
  case 2:
    return 0x32;
  case 3:
  case 4:
    return 0x41;
  }
  return uVar1;
}



/* ===== VT_GetTTSInfo_ENG @ 1002a690 ===== */

undefined4 __cdecl VT_GetTTSInfo_ENG(undefined4 param_1,char *param_2,UINT *param_3,int param_4)

{
  BYTE BVar1;
  char cVar2;
  UINT UVar3;
  char *pcVar4;
  int iVar5;
  uint uVar6;
  uint uVar7;
  UINT *pUVar8;
  BYTE *pBVar9;
  BYTE *pBVar10;
  char *pcVar11;
  BYTE local_208 [512];
  undefined4 local_8;

                    /* 0x2a690  27  VT_GetTTSInfo_ENG */
  local_8 = 0;
  if (DAT_1007d6dc < 0) {
    DAT_1007d6dc = FUN_10025e60();
  }
  if (DAT_1009fc50 == 0) {
    VT_SetDecimal0Pron_ENG();
  }
  if (DAT_100a0458 == 0) {
    InitializeCriticalSection((LPCRITICAL_SECTION)&DAT_100a8400);
    DAT_100a0458 = 1;
  }
  if (DAT_1009fe54 == 0) {
    FUN_10024cc0();
    DAT_1009fe54 = 1;
  }
  if (param_3 == (UINT *)0x0) {
    return 3;
  }
  switch(param_1) {
  case 0:
    uVar6 = 0xffffffff;
    pcVar11 = PTR_DAT_1007d6f4;
    do {
      pcVar4 = PTR_DAT_1007d6f4;
      if (uVar6 == 0) break;
      uVar6 = uVar6 - 1;
      cVar2 = *pcVar11;
      pcVar11 = pcVar11 + 1;
    } while (cVar2 != '\0');
    break;
  case 1:
    UVar3 = VT_CheckLicense_ENG(param_2,(char *)0x0,0xffffffff,(byte *)0x0);
    *param_3 = UVar3;
    return local_8;
  case 2:
    iVar5 = VT_CheckLicense_ENG(param_2,(char *)0x0,0xffffffff,(byte *)0x0);
    if (iVar5 == 0) {
      UVar3 = FUN_10028c50(param_2,(char *)0x0,0xffffffff);
      *param_3 = UVar3;
      return local_8;
    }
  case 6:
    *param_3 = 1;
    return local_8;
  case 3:
    iVar5 = -1;
    pcVar4 = &DAT_1009fa4c;
    do {
      if (iVar5 == 0) break;
      iVar5 = iVar5 + -1;
      cVar2 = *pcVar4;
      pcVar4 = pcVar4 + 1;
    } while (cVar2 != '\0');
    if (iVar5 == -2) {
      FUN_10028620(local_208,-1);
      uVar6 = 0xffffffff;
      pBVar9 = local_208;
      do {
        if (uVar6 == 0) break;
        uVar6 = uVar6 - 1;
        BVar1 = *pBVar9;
        pBVar9 = pBVar9 + 1;
      } while (BVar1 != '\0');
      if (param_4 <= (int)(~uVar6 - 1)) {
        return 4;
      }
      pBVar9 = local_208;
    }
    else {
      uVar6 = 0xffffffff;
      pcVar4 = &DAT_1009fa4c;
      do {
        if (uVar6 == 0) break;
        uVar6 = uVar6 - 1;
        cVar2 = *pcVar4;
        pcVar4 = pcVar4 + 1;
      } while (cVar2 != '\0');
      if (param_4 <= (int)(~uVar6 - 1)) {
        return 4;
      }
      pBVar9 = &DAT_1009fa4c;
    }
    uVar6 = 0xffffffff;
    do {
      pBVar10 = pBVar9;
      if (uVar6 == 0) break;
      uVar6 = uVar6 - 1;
      pBVar10 = pBVar9 + 1;
      BVar1 = *pBVar9;
      pBVar9 = pBVar10;
    } while (BVar1 != '\0');
    uVar6 = ~uVar6;
    pUVar8 = (UINT *)(pBVar10 + -uVar6);
    for (uVar7 = uVar6 >> 2; uVar7 != 0; uVar7 = uVar7 - 1) {
      *param_3 = *pUVar8;
      pUVar8 = pUVar8 + 1;
      param_3 = param_3 + 1;
    }
    for (uVar6 = uVar6 & 3; uVar6 != 0; uVar6 = uVar6 - 1) {
      *(BYTE *)param_3 = (BYTE)*pUVar8;
      pUVar8 = (UINT *)((int)pUVar8 + 1);
      param_3 = (UINT *)((int)param_3 + 1);
    }
    return local_8;
  case 4:
  case 8:
  case 9:
  case 0x13:
  case 0x16:
  case 0x1a:
    *param_3 = 0;
    return local_8;
  case 5:
    *param_3 = 6;
    return local_8;
  case 7:
    UVar3 = GetACP();
    *param_3 = UVar3;
    return local_8;
  case 10:
    *param_3 = 16000;
    return local_8;
  case 0xb:
  case 0x19:
    *param_3 = 200;
    return local_8;
  case 0xc:
  case 0xf:
  case 0x12:
    *param_3 = 100;
    return local_8;
  case 0xd:
  case 0x10:
    *param_3 = 0x32;
    return local_8;
  case 0xe:
    *param_3 = 400;
    return local_8;
  case 0x11:
    *param_3 = 500;
    return local_8;
  case 0x14:
  case 0x18:
    *param_3 = 0xffff;
    return local_8;
  case 0x15:
    *param_3 = 0x2af;
    return local_8;
  case 0x17:
    iVar5 = -1;
    pcVar4 = &DAT_1009fa4c;
    do {
      if (iVar5 == 0) break;
      iVar5 = iVar5 + -1;
      cVar2 = *pcVar4;
      pcVar4 = pcVar4 + 1;
    } while (cVar2 != '\0');
    if (iVar5 == -2) {
      FUN_10028620(local_208,-1);
      pcVar4 = FUN_1002a600();
    }
    else {
      pcVar4 = FUN_1002a600();
    }
    if (pcVar4 == (char *)0x0) {
      return local_8;
    }
    uVar6 = 0xffffffff;
    pcVar11 = pcVar4;
    do {
      if (uVar6 == 0) break;
      uVar6 = uVar6 - 1;
      cVar2 = *pcVar11;
      pcVar11 = pcVar11 + 1;
    } while (cVar2 != '\0');
    break;
  default:
    return 2;
  case 0x65:
    return DAT_100a7498;
  }
  if (param_4 <= (int)(~uVar6 - 1)) {
    return 4;
  }
  uVar6 = 0xffffffff;
  do {
    pcVar11 = pcVar4;
    if (uVar6 == 0) break;
    uVar6 = uVar6 - 1;
    pcVar11 = pcVar4 + 1;
    cVar2 = *pcVar4;
    pcVar4 = pcVar11;
  } while (cVar2 != '\0');
  uVar6 = ~uVar6;
  pUVar8 = (UINT *)(pcVar11 + -uVar6);
  for (uVar7 = uVar6 >> 2; uVar7 != 0; uVar7 = uVar7 - 1) {
    *param_3 = *pUVar8;
    pUVar8 = pUVar8 + 1;
    param_3 = param_3 + 1;
  }
  for (uVar6 = uVar6 & 3; uVar6 != 0; uVar6 = uVar6 - 1) {
    *(char *)param_3 = (char)*pUVar8;
    pUVar8 = (UINT *)((int)pUVar8 + 1);
    param_3 = (UINT *)((int)param_3 + 1);
  }
  return local_8;
}



/* ===== VT_INIT_ENG @ 1002aa50 ===== */

undefined2 VT_INIT_ENG(void)

{
                    /* 0x2aa50  29  VT_INIT_ENG */
  return 0xffff;
}



/* ===== VT_SpeakersInfo_ENG @ 1002aa60 ===== */

undefined4 __cdecl VT_SpeakersInfo_ENG(int param_1,char *param_2,char *param_3)

{
  char cVar1;
  uint uVar2;
  uint uVar3;
  char *pcVar4;
  char *pcVar5;

                    /* 0x2aa60  50  VT_SpeakersInfo_ENG */
  uVar2 = 0xffffffff;
  pcVar4 = (&PTR_DAT_1007c6a8)[param_1 * 6];
  do {
    pcVar5 = pcVar4;
    if (uVar2 == 0) break;
    uVar2 = uVar2 - 1;
    pcVar5 = pcVar4 + 1;
    cVar1 = *pcVar4;
    pcVar4 = pcVar5;
  } while (cVar1 != '\0');
  uVar2 = ~uVar2;
  pcVar4 = pcVar5 + -uVar2;
  for (uVar3 = uVar2 >> 2; uVar3 != 0; uVar3 = uVar3 - 1) {
    *(undefined4 *)param_2 = *(undefined4 *)pcVar4;
    pcVar4 = pcVar4 + 4;
    param_2 = param_2 + 4;
  }
  for (uVar2 = uVar2 & 3; uVar2 != 0; uVar2 = uVar2 - 1) {
    *param_2 = *pcVar4;
    pcVar4 = pcVar4 + 1;
    param_2 = param_2 + 1;
  }
  uVar2 = 0xffffffff;
  pcVar4 = (&PTR_s_d__eng_db_susan_pcm__1007c6ac)[param_1 * 6];
  do {
    pcVar5 = pcVar4;
    if (uVar2 == 0) break;
    uVar2 = uVar2 - 1;
    pcVar5 = pcVar4 + 1;
    cVar1 = *pcVar4;
    pcVar4 = pcVar5;
  } while (cVar1 != '\0');
  uVar2 = ~uVar2;
  pcVar4 = pcVar5 + -uVar2;
  for (uVar3 = uVar2 >> 2; uVar3 != 0; uVar3 = uVar3 - 1) {
    *(undefined4 *)param_3 = *(undefined4 *)pcVar4;
    pcVar4 = pcVar4 + 4;
    param_3 = param_3 + 4;
  }
  for (uVar2 = uVar2 & 3; uVar2 != 0; uVar2 = uVar2 - 1) {
    *param_3 = *pcVar4;
    pcVar4 = pcVar4 + 1;
    param_3 = param_3 + 1;
  }
  return 6;
}
