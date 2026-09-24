/* Ghidra decompiler output: approximate C pseudocode, not original source. */
/* Program: vt_pau.dll */

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



/* ===== FUN_10017510 ===== */
/* Entry: 10017510 */

/* WARNING: Type propagation algorithm not settling */

undefined4 __cdecl FUN_10017510(int param_1)

{
  char *pcVar1;
  char cVar2;
  short sVar3;
  short *psVar4;
  char cVar5;
  byte bVar6;
  int iVar7;
  int iVar8;
  int iVar9;
  short sVar10;
  short sVar11;
  int iVar12;
  uint uVar13;
  short sVar14;
  undefined4 *puVar15;
  int iVar16;
  char acStackY_11474 [1000];
  char acStackY_1108c [1000];
  char acStackY_10ca4 [1000];
  byte abStackY_108bc [1000];
  char acStackY_104d4 [30769];
  char acStackY_8ca3 [999];
  byte abStackY_88bc [1000];
  char acStackY_84d4 [28768];
  char local_1474 [1000];
  char local_108c [999];
  char cStackY_ca5;
  undefined4 local_ca4 [249];
  byte abStackY_8bd [999];
  char acStackY_4d6 [1002];
  short asStackY_ec [100];
  short *local_24;
  int local_20;

  FUN_10064270();
  sVar10 = 0;
  psVar4 = *(short **)(param_1 + 0x4c);
  if (0 < psVar4[1]) {
    do {
      sVar11 = 0;
      iVar7 = sVar10 * 0x3c0;
      if (*(char *)(iVar7 + 0x6e1 + (int)psVar4) != '\0') {
        do {
          iVar12 = (int)sVar11;
          sVar11 = sVar11 + 1;
          *(undefined1 *)(iVar12 + iVar7 + 0x8a9 + (int)psVar4) = 0;
        } while (sVar11 < (short)(ushort)*(byte *)(iVar7 + 0x6e1 + (int)psVar4));
      }
      sVar10 = sVar10 + 1;
    } while (sVar10 < psVar4[1]);
  }
  local_24 = psVar4;
  FUN_10017100(param_1);
  local_20 = 0;
  if (0 < *psVar4) {
    do {
      puVar15 = local_ca4;
      for (iVar7 = 0xfa; iVar7 != 0; iVar7 = iVar7 + -1) {
        *puVar15 = 0;
        puVar15 = puVar15 + 1;
      }
      iVar7 = *(int *)(local_24 + (short)local_20 * 8 + 10);
      cVar2 = *(char *)(local_24[(short)local_20 * 8 + 7] * 0x3c0 + -3 + iVar7);
      sVar10 = 0;
      iVar12 = 0;
      if (0 < local_24[(short)local_20 * 8 + 7]) {
        do {
          sVar11 = (short)iVar12;
          iVar8 = (int)sVar11;
          asStackY_ec[iVar8] = sVar10;
          if (sVar11 == 0) {
            cVar5 = '\x04';
          }
          else {
            cVar5 = *(char *)(iVar8 * 0x3c0 + -3 + iVar7);
            if (cVar5 == 'a') {
              cVar5 = '\x03';
            }
            else {
              cVar5 = (cVar5 != ']') + '\x01';
            }
          }
          iVar16 = (int)sVar10;
          *(char *)((int)local_ca4 + iVar16) = cVar5 << 3;
          iVar8 = iVar8 * 0x3c0;
          sVar11 = 0;
          if (*(char *)(iVar8 + 0x95 + iVar7) != '\0') {
            do {
              if ((0 < sVar11) && (*(char *)(sVar11 + iVar8 + 0x328 + iVar7) != '0')) {
                iVar9 = sVar11 + iVar16;
                *(char *)((int)local_ca4 + iVar9) = *(char *)((int)local_ca4 + iVar9) + '\b';
              }
              iVar9 = (int)sVar11;
              if ((iVar9 < (int)(*(byte *)(iVar8 + 0x95 + iVar7) - 1)) &&
                 (*(char *)(iVar9 + iVar8 + 0x329 + iVar7) != '0')) {
                pcVar1 = (char *)((int)local_ca4 + iVar16 + iVar9);
                *pcVar1 = *pcVar1 + '\x01';
              }
              local_108c[iVar9 + iVar16] = (char)iVar12;
              local_1474[iVar9 + iVar16] = (char)sVar11;
              sVar11 = sVar11 + 1;
            } while (sVar11 < (short)(ushort)*(byte *)(iVar8 + 0x95 + iVar7));
          }
          puVar15 = (undefined4 *)(iVar8 + 0x2e8 + iVar7);
          uVar13 = (uint)*(byte *)(iVar8 + 0x95 + iVar7);
          local_20 = 0x100176d7;
          FUN_10063f30((undefined4 *)(acStackY_4d6 + iVar16 + 2),puVar15,uVar13);
          cVar5 = *(char *)(iVar8 + 0x3bd + iVar7);
          sVar10 = sVar10 + (ushort)*(byte *)(iVar8 + 0x95 + iVar7);
          if (cVar5 == 'Z') {
            cVar5 = '\x05';
          }
          else if (cVar5 == '^') {
            cVar5 = '\x06';
          }
          else if (cVar5 == '[') {
            cVar5 = '\x04';
          }
          else if (cVar5 == 'a') {
            cVar5 = '\x03';
          }
          else {
            cVar5 = (cVar5 != ']') + '\x01';
          }
          *(char *)((int)local_ca4 + sVar10 + -1) = *(char *)((int)local_ca4 + sVar10 + -1) + cVar5;
          iVar12 = uVar13 + 1;
        } while ((short)iVar12 < *(short *)((int)puVar15 + 0xe));
      }
      puVar15 = local_ca4;
      uVar13 = (uint)sVar10;
      local_20 = 0x1001775c;
      FUN_10063f30((undefined4 *)((int)abStackY_8bd + 1),puVar15,uVar13);
      if (0 < sVar10) {
        iVar12 = 0;
        do {
          bVar6 = abStackY_8bd[iVar12 + 1];
          *(undefined1 *)((int)local_ca4 + iVar12) = 0;
          if ((bVar6 & 0x38) != 0) {
            *(undefined1 *)((int)local_ca4 + iVar12) = 0x80;
          }
          if ((bVar6 & 7) != 0) {
            *(char *)((int)local_ca4 + iVar12) = *(char *)((int)local_ca4 + iVar12) + '@';
          }
          iVar12 = iVar12 + 1;
          uVar13 = uVar13 - 1;
        } while (uVar13 != 0);
      }
      sVar14 = 1;
      sVar11 = sVar10;
      if (1 < sVar10) {
        do {
          bVar6 = abStackY_8bd[sVar14] >> 3 & 7;
          if (((bVar6 != 0) && ((&DAT_1007b9e0)[acStackY_4d6[sVar14 + 1]] == '\0')) &&
             (sVar14 < sVar10)) {
            do {
              cVar5 = acStackY_4d6[sVar14 + 2];
              abStackY_8bd[(int)sVar14 + 1] = (abStackY_8bd[(int)sVar14 + 1] & 7) + bVar6 * '\b';
              if ((&DAT_1007b9e0)[cVar5] != '\0') break;
              sVar14 = sVar14 + 1;
            } while (sVar14 < sVar10);
          }
          sVar14 = sVar14 + 1;
        } while (sVar14 < sVar10);
      }
joined_r0x10017808:
      sVar14 = sVar11;
      sVar11 = sVar14 + -1;
      if (0 < sVar11) {
        bVar6 = abStackY_8bd[(int)sVar11 + 1];
        if (((bVar6 & 7) != 0) && ((&DAT_1007b9e0)[acStackY_4d6[sVar11 + 2]] == '\0')) {
          for (sVar11 = sVar14 + -2; -1 < sVar11; sVar11 = sVar11 + -1) {
            if (((abStackY_8bd[(int)sVar11 + 1] & 7) != 0) ||
               (abStackY_8bd[(int)sVar11 + 1] = (abStackY_8bd[(int)sVar11 + 1] & 0x38) + (bVar6 & 7)
               , (&DAT_1007b9e0)[acStackY_4d6[sVar11 + 2]] != '\0')) break;
          }
        }
        goto joined_r0x10017808;
      }
      sVar11 = 0;
      if (0 < *(short *)((int)puVar15 + 0xe) + -1) {
        iVar12 = 0;
        do {
          iVar8 = iVar12 * 0x3c0;
          if (*(char *)(iVar8 + 0x3bc + iVar7) == '\x01') {
            sVar14 = *(byte *)(iVar8 + 0x95 + iVar7) - 1;
            if (-1 < sVar14) {
              sVar3 = asStackY_ec[iVar12];
              do {
                iVar16 = (int)sVar3 + (int)sVar14;
                if (((abStackY_8bd[iVar16 + 1] & 7) == 1) || ((abStackY_8bd[iVar16 + 1] & 7) == 2))
                {
                  iVar9 = sVar14 + iVar8;
                  *(char *)(iVar9 + 0x25d + iVar7) = *(char *)(iVar9 + 0x25d + iVar7) + '\x01';
                }
              } while (((&DAT_1007b9e0)[acStackY_4d6[iVar16 + 2]] == '\0') &&
                      (sVar14 = sVar14 + -1, -1 < sVar14));
            }
            sVar14 = 0;
            if (*(char *)(iVar8 + 0x455 + iVar7) != '\0') {
              sVar3 = asStackY_ec[iVar12 + 1];
              do {
                iVar12 = (int)sVar3 + (int)sVar14;
                bVar6 = abStackY_8bd[iVar12 + 1] >> 3 & 7;
                if ((bVar6 == 1) || (bVar6 == 2)) {
                  iVar16 = sVar14 + iVar8;
                  *(char *)(iVar16 + 0x61d + iVar7) = *(char *)(iVar16 + 0x61d + iVar7) + '\n';
                }
              } while (((&DAT_1007b9e0)[acStackY_4d6[iVar12 + 2]] == '\0') &&
                      (sVar14 = sVar14 + 1, sVar14 < (short)(ushort)*(byte *)(iVar8 + 0x455 + iVar7)
                      ));
            }
          }
          sVar11 = sVar11 + 1;
          iVar12 = (int)sVar11;
        } while (iVar12 < *(short *)((int)puVar15 + 0xe) + -1);
      }
      sVar11 = 0;
      if (0 < *(short *)((int)puVar15 + 0xe)) {
        do {
          sVar14 = 0;
          iVar12 = sVar11 * 0x3c0;
          bVar6 = *(byte *)(iVar12 + 0x95 + iVar7);
          if (bVar6 != 0) {
            do {
              iVar8 = sVar14 + iVar12 + iVar7;
              if (*(char *)(iVar8 + 0x329) == '2') {
                if (((int)sVar14 < (int)(bVar6 - 1)) && (*(byte *)(iVar8 + 0x25e) / 10 == 0)) {
                  *(byte *)(iVar8 + 0x25e) = *(byte *)(iVar8 + 0x25e) + 10;
                }
                if ((uint)*(byte *)(iVar8 + 0x25d) % 10 == 0) {
                  *(byte *)(iVar8 + 0x25d) = *(byte *)(iVar8 + 0x25d) + 1;
                }
              }
              bVar6 = *(byte *)(iVar12 + 0x95 + iVar7);
              sVar14 = sVar14 + 1;
            } while (sVar14 < (short)(ushort)bVar6);
          }
          sVar11 = sVar11 + 1;
        } while (sVar11 < *(short *)((int)puVar15 + 0xe));
      }
      sVar11 = 0;
      if (0 < sVar10) {
        do {
          if (sVar11 == 0) {
            iVar12 = iVar7 + local_108c[0] * 0x3c0 + local_1474[0] * 7;
            *(undefined1 *)(iVar12 + 0x96) = 0x5a;
            *(undefined1 *)(iVar12 + 0x97) = 0x5a;
          }
          else if (sVar11 == 1) {
            iVar12 = iVar7 + local_108c[1] * 0x3c0 + local_1474[1] * 7;
            *(undefined1 *)(iVar12 + 0x96) = 0x5a;
            *(char *)(iVar12 + 0x97) = acStackY_4d6[2];
          }
          else {
            iVar12 = local_1474[sVar11] * 7 + iVar7 + local_108c[sVar11] * 0x3c0;
            *(char *)(iVar12 + 0x96) = acStackY_4d6[sVar11];
            *(char *)(iVar12 + 0x97) = acStackY_4d6[sVar11 + 1];
          }
          iVar12 = (int)sVar11;
          iVar9 = (int)local_1474[iVar12];
          iVar8 = local_108c[iVar12] * 0x3c0;
          iVar16 = iVar9 * 7 + iVar8 + iVar7;
          *(char *)(iVar16 + 0x98) = acStackY_4d6[iVar12 + 2];
          cVar5 = cVar2;
          if (iVar12 == sVar10 + -1) {
            *(char *)(iVar16 + 0x99) = cVar2;
            iVar9 = (iVar9 + 0x16) * 7 + iVar8;
          }
          else {
            *(char *)(iVar16 + 0x99) = acStackY_4d6[iVar12 + 3];
            if (iVar12 == sVar10 + -2) {
              iVar9 = (iVar9 + 0x16) * 7 + iVar8;
            }
            else {
              iVar9 = (iVar9 + 0x16) * 7 + iVar8;
              cVar5 = acStackY_4d6[iVar12 + 4];
            }
          }
          *(char *)(iVar9 + iVar7) = cVar5;
          *(byte *)(iVar16 + 0x9b) = abStackY_8bd[iVar12 + 1] + *(char *)((int)local_ca4 + iVar12);
          *(undefined1 *)(iVar16 + 0x9c) = 0;
          if (*(char *)(iVar8 + 0x2df + iVar7) == '\f') {
            *(undefined1 *)(iVar16 + 0x9c) = 0x20;
          }
          sVar11 = sVar11 + 1;
        } while (sVar11 < sVar10);
      }
      local_20 = local_20 + 1;
    } while ((short)local_20 < *local_24);
  }
  return 1;
}



/* ===== FUN_10013380 ===== */
/* Entry: 10013380 */

undefined4 __cdecl FUN_10013380(int param_1,int param_2)

{
  int iVar1;
  int iVar2;
  undefined2 uVar3;
  byte *pbVar4;
  int iVar5;
  short local_2c [10];
  byte *local_18;
  int local_14;
  byte *local_10;
  int local_c;
  int local_8;

  iVar2 = param_1;
  local_c = 0;
  iVar5 = *(int *)(param_2 + 0x4c);
  if (*(short *)(iVar5 + 2) < 1) {
    return CONCAT22((short)((uint)iVar5 >> 0x10),1);
  }
  pbVar4 = (byte *)(iVar5 + 0x6e0);
  local_14 = iVar5;
  do {
    param_2 = 0;
    param_1 = 0;
    if (*pbVar4 != 0) {
      local_8 = 0;
      do {
        iVar5 = 0;
        local_18 = pbVar4 + *(byte *)(local_8 + *(int *)(pbVar4 + -0x8c) + 0x1c) + 0x254;
        if (*(char *)(local_8 + *(int *)(pbVar4 + -0x8c) + 0x1d) != '\0') {
          local_10 = pbVar4 + param_2 * 2 + -0x84;
          do {
            FUN_100135d0(local_2c,local_14,local_c,param_1,iVar5);
            iVar1 = (char)(&DAT_1007b6c0)[(char)local_18[iVar5]] * 0xb;
            if ((&DAT_1007987b)[iVar1] == '\x01') {
              switch((&DAT_1007987c)[iVar1]) {
              case 1:
                uVar3 = FUN_10001670((int *)(iVar2 + 0xb4),(int)local_2c);
                break;
              case 2:
                uVar3 = FUN_10001670((int *)(iVar2 + 0xd0),(int)local_2c);
                break;
              case 3:
                uVar3 = FUN_10001670((int *)(iVar2 + 0xec),(int)local_2c);
                break;
              case 4:
                uVar3 = FUN_10001670((int *)(iVar2 + 0x108),(int)local_2c);
                break;
              default:
                uVar3 = 0x640;
              }
            }
            else {
              switch((&DAT_10079880)[iVar1]) {
              case 1:
                uVar3 = FUN_10001670((int *)(iVar2 + 0x124),(int)local_2c);
                break;
              case 2:
                uVar3 = FUN_10001670((int *)(iVar2 + 0x140),(int)local_2c);
                break;
              case 3:
                uVar3 = FUN_10001670((int *)(iVar2 + 0x15c),(int)local_2c);
                break;
              case 4:
              case 5:
                uVar3 = FUN_10001670((int *)(iVar2 + 0x178),(int)local_2c);
                break;
              case 6:
                uVar3 = FUN_10001670((int *)(iVar2 + 0x194),(int)local_2c);
                break;
              default:
                uVar3 = 500;
              }
            }
            param_2 = param_2 + 1;
            *(undefined2 *)local_10 = uVar3;
            local_10 = local_10 + 2;
            iVar5 = iVar5 + 1;
          } while (iVar5 < (int)(uint)*(byte *)(local_8 + 0x1d + *(int *)(pbVar4 + -0x8c)));
        }
        param_1 = param_1 + 1;
        local_8 = local_8 + 0x1e;
        iVar5 = local_14;
      } while (param_1 < (int)(uint)*pbVar4);
    }
    pbVar4 = pbVar4 + 0x3c0;
    local_c = local_c + 1;
  } while (local_c < *(short *)(iVar5 + 2));
  return CONCAT22((short)((uint)iVar5 >> 0x10),1);
}



/* ===== FUN_100138c0 ===== */
/* Entry: 100138c0 */

undefined4 __cdecl FUN_100138c0(int param_1,int param_2)

{
  byte bVar1;
  undefined2 uVar2;
  int *piVar3;
  int iVar4;
  byte *pbVar5;
  int iVar6;
  int *piVar7;
  short *psVar8;
  short local_28 [11];
  short local_12;
  short *local_10;
  int local_c;
  int local_8;

  local_c = 0;
  psVar8 = *(short **)(param_2 + 0x4c);
  iVar4 = param_2;
  local_10 = psVar8;
  if (0 < psVar8[1]) {
    pbVar5 = (byte *)(psVar8 + 0x370);
    do {
      iVar4 = 0;
      if (*pbVar5 != 0) {
        local_8 = 0;
        do {
          FUN_10013a20(local_28,(int)psVar8,local_c,iVar4);
          bVar1 = pbVar5[0x329];
          if ((bVar1 == 0x5e) && (iVar4 == *pbVar5 - 1)) {
            piVar3 = (int *)(param_1 + 0x204);
            piVar7 = (int *)(param_1 + 0x274);
          }
          else if ((bVar1 == 0x5a) && (iVar4 == *pbVar5 - 1)) {
            piVar3 = (int *)(param_1 + 0x1e8);
            piVar7 = (int *)(param_1 + 600);
          }
          else if ((bVar1 == 0x5b) && (iVar4 == *pbVar5 - 1)) {
            piVar3 = (int *)(param_1 + 0x1cc);
            piVar7 = (int *)(param_1 + 0x23c);
          }
          else {
            piVar3 = (int *)(param_1 + 0x1b0);
            piVar7 = (int *)(param_1 + 0x220);
          }
          uVar2 = FUN_10001670(piVar3,(int)local_28);
          *(char *)(local_8 + 3 + *(int *)(pbVar5 + -0x8c)) = (char)uVar2;
          local_12 = (short)*(char *)(local_8 + 3 + *(int *)(pbVar5 + -0x8c));
          FUN_100016a0(piVar7,(int)local_28,(undefined4 *)(local_8 + 4 + *(int *)(pbVar5 + -0x8c)));
          iVar4 = iVar4 + 1;
          local_8 = local_8 + 0x1e;
          psVar8 = local_10;
        } while (iVar4 < (int)(uint)*pbVar5);
      }
      pbVar5 = pbVar5 + 0x3c0;
      iVar4 = local_c + 1;
      local_c = iVar4;
    } while (iVar4 < psVar8[1]);
  }
  iVar6 = 0;
  if (0 < *psVar8) {
    do {
      FUN_100137c0(iVar6,param_2);
      iVar4 = (int)*psVar8;
      iVar6 = iVar6 + 1;
    } while (iVar6 < iVar4);
  }
  return CONCAT22((short)((uint)iVar4 >> 0x10),1);
}



/* ===== FUN_10022dc0 ===== */
/* Entry: 10022dc0 */

int __cdecl FUN_10022dc0(int param_1,int *param_2)

{
  short sVar1;
  int iVar2;
  byte *pbVar3;
  int iVar4;
  int iVar5;
  int *piVar6;
  byte *pbVar7;
  int *piVar8;
  int iVar9;

  piVar6 = param_2;
  iVar2 = param_2[0x13];
  pbVar3 = (byte *)param_2[0x4c4af];
  if (param_2[2] == 0) {
    return 0;
  }
  while( true ) {
    FUN_10022850(param_1,(int)piVar6);
    FUN_1003d350((undefined **)(piVar6[2] + piVar6[1]),pbVar3);
    iVar4 = *(int *)(pbVar3 + 4);
    param_2 = (int *)0x0;
    if (0 < *(short *)pbVar3) {
      piVar8 = (int *)(iVar2 + 0x650);
      pbVar7 = pbVar3 + 0x18;
      do {
        iVar9 = *(int *)(pbVar7 + -4);
        iVar5 = piVar6[1];
        piVar8[-1] = iVar5 + iVar9;
        if (*(int *)(pbVar7 + -4) < *(int *)pbVar7) {
          *piVar8 = *(int *)pbVar7 + -1 + piVar6[1];
        }
        else {
          *piVar8 = iVar5 + iVar9;
        }
        pbVar7 = pbVar7 + 0x94;
        param_2 = (int *)((int)param_2 + 1);
        piVar8 = piVar8 + 0xf0;
      } while ((int)param_2 < (int)*(short *)pbVar3);
    }
    sVar1 = *(short *)pbVar3;
    iVar9 = 0;
    *(short *)(iVar2 + 2) = sVar1;
    if (0 < sVar1) {
      piVar8 = piVar6 + 0x48508;
      do {
        piVar8[-200] = piVar6[0x4843f];
        *piVar8 = piVar6[0x4843d];
        piVar8[200] = piVar6[0x4843e];
        piVar8[400] = -1;
        piVar8[600] = -1;
        piVar8[800] = -1;
        iVar9 = iVar9 + 1;
        piVar8 = piVar8 + 1;
      } while (iVar9 < *(short *)(iVar2 + 2));
    }
    if (iVar4 == 0) break;
    if (0 < *(short *)pbVar3) {
      FUN_10022970(piVar6);
      if (*(char *)(DAT_100a0460 + 0x20424) == '\x01') {
        iVar9 = 0;
        if (0 < *(short *)(iVar2 + 2)) {
          piVar8 = (int *)(iVar2 + 0x650);
          do {
            iVar9 = iVar9 + 1;
            piVar8[-1] = *(int *)(piVar6[0x11de4] + *(int *)(piVar6[0x11de5] + piVar8[-1] * 4) * 4);
            *piVar8 = *(int *)(piVar6[0x11de4] + *(int *)(piVar6[0x11de6] + *piVar8 * 4) * 4);
            piVar8 = piVar8 + 0xf0;
          } while (iVar9 < *(short *)(iVar2 + 2));
        }
      }
      else {
        iVar9 = 0;
        if (0 < *(short *)(iVar2 + 2)) {
          piVar8 = (int *)(iVar2 + 0x650);
          do {
            if (*piVar6 <= piVar8[-1]) {
              piVar8[-1] = *piVar6 + -1;
            }
            if (*piVar6 <= *piVar8) {
              *piVar8 = *piVar6 + -1;
            }
            if (piVar8[-1] < 0) {
              piVar8[-1] = 0;
            }
            if (*piVar8 < 0) {
              *piVar8 = 0;
            }
            iVar9 = iVar9 + 1;
            piVar8[-1] = *(int *)(piVar6[0x11de5] + piVar8[-1] * 4);
            *piVar8 = *(int *)(piVar6[0x11de6] + *piVar8 * 4);
            piVar8 = piVar8 + 0xf0;
          } while (iVar9 < *(short *)(iVar2 + 2));
        }
      }
      iVar2 = piVar6[1];
      piVar6[1] = iVar2 + iVar4;
      return iVar2 + iVar4;
    }
    piVar6[1] = piVar6[1] + iVar4;
  }
  piVar6[0x11] = 1;
  return piVar6[1];
}



/* ===== FUN_10026750 ===== */
/* Entry: 10026750 */

undefined4 __cdecl FUN_10026750(int param_1,int param_2)

{
  uint *puVar1;
  int *piVar2;
  uint uVar3;

  puVar1 = *(uint **)(*(int *)(param_2 + 0x4c) + 0x4770c);
  if ((int)*puVar1 < 0) {
    if ((int)*(uint *)(param_2 + 0x1312c8) < 0) {
      uVar3 = *(uint *)(param_2 + 0x1312c4);
      if ((int)uVar3 < 0) {
        if (*(char *)(param_2 + 0x122850) != '\0') {
          *puVar1 = 0;
          goto LAB_100267a6;
        }
        uVar3 = *(uint *)(param_1 + 0x4d00);
      }
    }
    else {
      uVar3 = ~-(uint)(*(char *)(param_2 + 0x122850) != '\0') & *(uint *)(param_2 + 0x1312c8);
    }
    *puVar1 = uVar3;
  }
LAB_100267a6:
  piVar2 = *(int **)(*(int *)(param_2 + 0x4c) + 0x4770c);
  if (0xfffe < *piVar2) {
    *piVar2 = 0xffff;
  }
  return 1;
}



/* ===== FUN_10016c90 ===== */
/* Entry: 10016c90 */

undefined4 __cdecl FUN_10016c90(int param_1)

{
  char cVar1;
  short sVar2;
  int iVar3;
  int iVar4;
  char *pcVar5;
  int iVar6;
  uint uVar7;
  undefined2 *puVar8;
  undefined2 *puVar9;
  char *pcVar10;
  int iVar11;
  byte *pbVar12;
  int local_8;

  iVar3 = *(int *)(param_1 + 0x4c);
  iVar4 = *(int *)(param_1 + 0x1312bc);
  local_8 = 0;
  if (0 < *(short *)(iVar3 + 2)) {
    puVar8 = (undefined2 *)(iVar3 + 0xa00);
    puVar9 = (undefined2 *)(iVar4 + 0x24);
    do {
      *(undefined2 **)(puVar8 + -0x6a) = puVar9 + 0x12;
      *puVar8 = *puVar9;
      *(undefined2 **)(puVar8 + 2) = puVar9 + 0x10;
      *(undefined2 **)(puVar8 + -0x68) = puVar9 + 0x21;
      *(undefined1 *)((int)puVar8 + -3) = *(undefined1 *)((int)puVar9 + 0x13);
      puVar8[-4] = puVar9[0xe];
      *(undefined1 *)(puVar8 + -2) = *(undefined1 *)(puVar9 + 0xc);
      *(undefined1 *)(puVar8 + -1) = *(undefined1 *)((int)puVar9 + 0x15);
      *(undefined1 *)((int)puVar8 + -1) = *(undefined1 *)(puVar9 + 0xb);
      puVar8[-3] = puVar9[8];
      pcVar5 = *(char **)(puVar8 + -0x68);
      uVar7 = 0xffffffff;
      pcVar10 = pcVar5;
      do {
        if (uVar7 == 0) break;
        uVar7 = uVar7 - 1;
        cVar1 = *pcVar10;
        pcVar10 = pcVar10 + 1;
      } while (cVar1 != '\0');
      param_1 = ~uVar7 - 1;
      if (0x41 < param_1) {
        param_1 = 0x41;
      }
      iVar6 = 0;
      iVar11 = 0;
      if (param_1 < 1) {
        return 0;
      }
      do {
        if (pcVar5[iVar11] == 'd') {
          if (0 < iVar6) {
            *(undefined1 *)((int)puVar8 + iVar6 + -0x8c) = 0x31;
          }
        }
        else if (pcVar5[iVar11] == 'c') {
          if (0 < iVar6) {
            *(undefined1 *)((int)puVar8 + iVar6 + -0x8c) = 0x32;
          }
        }
        else {
          *(undefined1 *)((int)puVar8 + iVar6 + -0x116) = 0;
          cVar1 = pcVar5[iVar11];
          if (cVar1 == 'M') {
            *(undefined1 *)((int)puVar8 + iVar6 + -0x116) = 1;
          }
          else if ((cVar1 < '\x01') || ('E' < cVar1)) break;
          *(char *)((int)puVar8 + iVar6 + -0xcc) = pcVar5[iVar11];
          *(undefined1 *)((int)puVar8 + iVar6 + -0x8b) = 0x30;
          iVar6 = iVar6 + 1;
        }
        iVar11 = iVar11 + 1;
      } while (iVar11 < param_1);
      if (iVar6 == 0) {
        return 0;
      }
      *(char *)((int)puVar8 + -799) = (char)iVar6;
      puVar9 = puVar9 + 0x4a;
      local_8 = local_8 + 1;
      puVar8 = puVar8 + 0x1e0;
    } while (local_8 < *(short *)(iVar3 + 2));
  }
  iVar6 = 0;
  if (0 < *(short *)(iVar3 + 2)) {
    pbVar12 = (byte *)(iVar3 + 0x92b);
    do {
      iVar11 = iVar6 * 4;
      iVar6 = iVar6 + 1;
      *pbVar12 = (*(int *)(*(int *)(iVar3 + 0x4771c) + iVar11) == -1) - 1U & 0xc;
      pbVar12 = pbVar12 + 0x3c0;
    } while (iVar6 < *(short *)(iVar3 + 2));
  }
  iVar6 = *(short *)(iVar3 + 2) * 0x3c0 + iVar3;
  sVar2 = *(short *)(iVar6 + 0x638);
  if (sVar2 == 3) {
    if (*(short *)(iVar4 + 2) == 0) {
      *(undefined1 *)(iVar3 + 0x4770a) = 6;
      return 1;
    }
  }
  else if (sVar2 != 4) {
    if (*(char *)(iVar6 + 0x56b) == '\f') {
      *(undefined1 *)(iVar3 + 0x4770a) = 7;
      return 1;
    }
    *(undefined1 *)(iVar3 + 0x4770a) = 5;
    return 1;
  }
  *(undefined1 *)(iVar3 + 0x4770a) = 7;
  return 1;
}



/* ===== FUN_100130e0 ===== */
/* Entry: 100130e0 */

undefined4 __cdecl FUN_100130e0(int param_1,int param_2)

{
  char cVar1;
  short *psVar2;
  int iVar3;
  char *pcVar4;
  undefined4 *puVar5;
  undefined1 *puVar6;
  uint uVar7;
  undefined4 uVar8;
  short *psVar9;
  int *piVar10;
  int iVar11;
  int *piVar12;
  byte *pbVar13;

  psVar2 = *(short **)(param_2 + 0x4c);
  iVar11 = *(int *)(param_2 + 0x1312bc);
  FUN_10012c70((int)psVar2);
  iVar3 = 0;
  if (0 < psVar2[1]) {
    pcVar4 = (char *)((int)psVar2 + 0xa09);
    do {
      iVar3 = iVar3 + 1;
      *pcVar4 = s___Z_Z_________10079a30[*(short *)(pcVar4 + -0x11)];
      pcVar4 = pcVar4 + 0x3c0;
    } while (iVar3 < psVar2[1]);
  }
  if (*(char *)(iVar3 * 0x3c0 + 0x649 + (int)psVar2) != '^') {
    *(undefined1 *)(iVar3 * 0x3c0 + 0x649 + (int)psVar2) = 0x5a;
  }
  if ((char)psVar2[0x23b85] == '\x06') {
    *(undefined1 *)(psVar2[1] * 0x3c0 + 0x649 + (int)psVar2) = 0x5e;
  }
  else {
    *(undefined1 *)(psVar2[1] * 0x3c0 + 0x649 + (int)psVar2) = 0x5a;
  }
  FUN_10012df0(psVar2);
  FUN_10012f00(param_1,param_2);
  iVar3 = 0;
  if (0 < psVar2[1]) {
    psVar9 = psVar2 + 0x504;
    do {
      *(undefined1 *)psVar9 = 0;
      iVar3 = iVar3 + 1;
      psVar9 = psVar9 + 0x1e0;
    } while (iVar3 < psVar2[1]);
  }
  iVar3 = 0;
  if (psVar2[1] != 1 && -1 < psVar2[1] + -1) {
    pcVar4 = (char *)((int)psVar2 + 0xa09);
    do {
      if (((pcVar4[-0xd] == '\b') || (pcVar4[0x3b3] != '\b')) ||
         ((*pcVar4 != ']' && (*pcVar4 != '\\')))) {
        if (*(short *)(pcVar4 + -0x11) == 0xc) goto LAB_100131e8;
      }
      else {
        *pcVar4 = '\\';
LAB_100131e8:
        pcVar4[-1] = '\x01';
      }
      iVar3 = iVar3 + 1;
      pcVar4 = pcVar4 + 0x3c0;
    } while (iVar3 < psVar2[1] + -1);
  }
  iVar3 = 0;
  if (psVar2[1] != 1 && -1 < psVar2[1] + -1) {
    piVar10 = (int *)(param_2 + 0x121a64);
    piVar12 = (int *)(iVar11 + 0x20);
    do {
      iVar11 = *piVar12;
      if (iVar11 == -2) {
        *piVar10 = 100;
      }
      else if (-1 < iVar11) {
        *piVar10 = iVar11;
      }
      iVar3 = iVar3 + 1;
      piVar10 = piVar10 + 1;
      piVar12 = piVar12 + 0x25;
    } while (iVar3 < psVar2[1] + -1);
  }
  iVar11 = 1;
  if (1 < psVar2[1]) {
    puVar5 = (undefined4 *)(param_2 + 0x121d80);
    do {
      if (-1 < (int)puVar5[-199]) {
        *puVar5 = 2;
      }
      iVar11 = iVar11 + 1;
      puVar5 = puVar5 + 1;
    } while (iVar11 < psVar2[1]);
  }
  iVar11 = 0;
  if (psVar2[1] != 1 && -1 < psVar2[1] + -1) {
    puVar6 = (undefined1 *)((int)psVar2 + 0xa09);
    puVar5 = (undefined4 *)(param_2 + 0x121d80);
    do {
      switch(*puVar5) {
      case 0:
        *puVar6 = 0x5d;
        puVar6[-1] = 1;
        break;
      case 1:
        *puVar6 = 0x5c;
        puVar6[-1] = 1;
        break;
      case 2:
        *puVar6 = 0x5b;
        break;
      case 3:
        *puVar6 = 0x5a;
      }
      iVar11 = iVar11 + 1;
      puVar5 = puVar5 + 1;
      puVar6 = puVar6 + 0x3c0;
    } while (iVar11 < psVar2[1] + -1);
  }
  iVar3 = *(int *)(param_2 + 0x121d80 + iVar11 * 4);
  if ((-1 < iVar3) && (iVar3 < 3)) {
    *(undefined1 *)(iVar11 * 0x3c0 + 0xa09 + (int)psVar2) = 0x5b;
  }
  uVar7 = 0;
  iVar11 = 0;
  param_2 = 0;
  if (0 < psVar2[1]) {
    pcVar4 = (char *)((int)psVar2 + 0xa09);
    pbVar13 = (byte *)((int)psVar2 + 0x6e1);
    do {
      iVar11 = iVar11 + (uint)*pbVar13 * 2;
      if (1000 < iVar11) {
        if (uVar7 == param_2) {
          return uVar7 & 0xffff0000;
        }
        uVar7 = uVar7 - 1;
        pbVar13 = pbVar13 + -0x3c0;
        pcVar4 = pcVar4 + -0x3c0;
        *pcVar4 = '[';
      }
      cVar1 = *pcVar4;
      if ((((cVar1 == '[') || (cVar1 == 'Z')) || (cVar1 == '^')) || (cVar1 == '`')) {
        param_2 = uVar7 + 1;
        iVar11 = 0;
      }
      uVar7 = uVar7 + 1;
      pbVar13 = pbVar13 + 0x3c0;
      pcVar4 = pcVar4 + 0x3c0;
    } while ((int)uVar7 < (int)psVar2[1]);
  }
  uVar8 = FUN_10012df0(psVar2);
  return CONCAT22((short)((uint)uVar8 >> 0x10),1);
}



/* ===== FUN_10026630 ===== */
/* Entry: 10026630 */

int __cdecl FUN_10026630(int param_1,int *param_2)

{
  uint uVar1;
  undefined4 uVar2;
  int iVar3;

  *(undefined2 *)(param_2 + 0x10) = 0;
  *(undefined2 *)(param_2[0x13] + 0x47704) = 0;
  *(undefined2 *)(param_2[0x13] + 0x47706) = 0;
  param_2[0x3ef5b] = 0;
  do {
    uVar1 = FUN_10022dc0(param_1,param_2);
    param_2[1] = uVar1;
    if (((int)uVar1 < 0) || (param_2[0x11] != 0)) goto LAB_1002669f;
    FUN_10026750(param_1,(int)param_2);
  } while ((char)param_2[8] == '\0');
  uVar1 = FUN_10016c90((int)param_2);
  if (uVar1 == 0) {
    return 0xffff;
  }
  if ((char)param_2[8] != '\x01') {
    uVar2 = FUN_100130e0(param_1,(int)param_2);
    if ((short)uVar2 == 0) {
      return CONCAT22((short)((uint)uVar2 >> 0x10),0xffff);
    }
    iVar3 = FUN_10017510((int)param_2);
    if (iVar3 == 0) {
      return 0xffff;
    }
    if (*(char *)((int)param_2 + 0x21) == '\x03') {
      FUN_100180b0(param_1,(int)param_2);
    }
    if (*(char *)((int)param_2 + 0x21) == '\x05') {
      FUN_10017c30(param_1,(int)param_2);
    }
    if ((*(char *)((int)param_2 + 0x21) == '\x06') || (*(char *)((int)param_2 + 0x21) == '\a')) {
      FUN_10017e80(param_1,(int)param_2);
    }
    if (*(char *)((int)param_2 + 0x21) == '\n') {
      FUN_10017d80(param_1,(int)param_2);
    }
    uVar2 = FUN_10013380(param_1,(int)param_2);
    if ((short)uVar2 == 0) {
      return CONCAT22((short)((uint)uVar2 >> 0x10),0xffff);
    }
    uVar2 = FUN_100138c0(param_1,(int)param_2);
    return ((short)uVar2 != 0) - 1;
  }
LAB_1002669f:
  return uVar1 & 0xffff0000;
}



/* ===== FUN_100267d0 ===== */
/* Entry: 100267d0 */

short * __cdecl FUN_100267d0(int *param_1,int *param_2)

{
  short sVar1;
  short *psVar2;
  uint uVar3;
  short *psVar4;

  psVar4 = (short *)param_2[0x11];
  param_2[0xc] = 0;
  psVar2 = psVar4;
  while (psVar2 == (short *)0x0) {
    uVar3 = FUN_1002c9b0(param_1,param_2);
    sVar1 = (short)param_2[0x10];
    psVar4 = (short *)CONCAT22((short)(uVar3 >> 0x10),sVar1);
    if (sVar1 == 1) break;
    if (sVar1 == 0) {
      *(short *)(param_2[0x13] + 0x47704) = *(short *)(param_2[0x13] + 0x47704) + 1;
      *(undefined2 *)(param_2[0x13] + 0x47706) = 0;
      psVar4 = (short *)param_2[0x13];
      param_2[0x3ef5b] = 0;
      if ((psVar4[1] == 0) || (psVar4[0x23b82] == *psVar4)) {
        psVar4 = (short *)FUN_10026630((int)param_1,param_2);
        while (((short)psVar4 < 0 || (*(short *)(param_2[0x13] + 2) == 0))) {
          psVar4 = (short *)FUN_10026630((int)param_1,param_2);
          if (param_2[0x11] != 0) {
            return psVar4;
          }
        }
      }
    }
    psVar2 = (short *)param_2[0x11];
  }
  return (short *)((uint)psVar4 & 0xffff0000);
}



/* ===== FUN_1002c9b0 ===== */
/* Entry: 1002c9b0 */

uint __cdecl FUN_1002c9b0(int *param_1,undefined4 *param_2)

{
  char cVar1;
  undefined4 in_EAX;
  uint uVar2;

  cVar1 = *(char *)(param_2 + 8);
  uVar2 = CONCAT31((int3)((uint)in_EAX >> 8),cVar1);
  if (cVar1 == '\x01') goto LAB_1002c9fc;
  if (*(short *)(param_2 + 0x10) == 0) {
    if (cVar1 == '\x02') goto LAB_1002c9fc;
    FUN_10024980(*(short *)(param_2[0x13] + 0x47704),param_2,param_1);
    uVar2 = param_2[0x3b189];
    if (uVar2 == 0) goto LAB_1002c9fc;
  }
  if (*(char *)(param_2 + 8) != '\x03') {
    uVar2 = FUN_1002bd90(param_2,param_1);
  }
LAB_1002c9fc:
  return uVar2 & 0xffff0000;
}



/* ===== FUN_1002bd90 ===== */
/* Entry: 1002bd90 */

uint __cdecl FUN_1002bd90(undefined4 *param_1,int *param_2)

{
  int iVar1;
  int *piVar2;
  int iVar3;
  undefined4 *puVar4;
  int *piVar5;
  int iVar6;
  int iVar7;
  uint uVar8;
  uint uVar9;
  uint uVar10;
  short sVar11;

  iVar1 = param_1[0x13];
  piVar2 = (int *)param_1[0x11ddd];
  iVar7 = CONCAT22((short)((uint)piVar2 >> 0x10),*(short *)(iVar1 + 0x47704));
  if ((*(short *)(iVar1 + 0x47706) == 0) && (*(short *)(param_1 + 0x10) == 0)) {
    iVar7 = *(short *)(iVar1 + 0x47704) * 0x10;
    FUN_1002c220((int)(param_1 + 0x3b77e),(int)param_1,*(ushort *)(iVar7 + 0xc + iVar1),
                 CONCAT22((short)((uint)iVar7 >> 0x10),*(undefined2 *)(iVar7 + 0xe + iVar1)),
                 (int)param_2);
    if (*(char *)((int)param_1 + 0x122851) == '\x01') {
      uVar8 = param_1[0xc] + 2;
      param_1[0xc] = uVar8;
      if (60000 < (int)uVar8) {
        param_1[0xc] = 0;
      }
      return uVar8 & 0xffff0000;
    }
    if (*(short *)(iVar1 + 0x47706) == 0) {
      *(undefined2 *)(iVar1 + 0x47708) = 0;
    }
    iVar7 = FUN_1002c530((int)param_1,(int)(param_1 + 0x3b77e));
  }
  sVar11 = *(short *)(iVar1 + 0x47706);
  if (sVar11 < *(short *)(param_1 + 0x3ef5a)) {
    do {
      iVar3 = (int)sVar11;
      iVar7 = param_1[0xc];
      if (*(char *)((int)param_1 + iVar3 * 0x34 + 0xede1f) == '\x01') {
        uVar8 = param_1[iVar3 * 0xd + 0x3b781] * 2;
        if (60000 < (int)(iVar7 + uVar8)) {
          uVar10 = 60000 - param_1[0xc];
          puVar4 = (undefined4 *)(param_1[0xd] + param_1[0xc]);
          for (uVar9 = uVar10 >> 2; uVar9 != 0; uVar9 = uVar9 - 1) {
            *puVar4 = 0;
            puVar4 = puVar4 + 1;
          }
          for (uVar9 = uVar10 & 3; uVar9 != 0; uVar9 = uVar9 - 1) {
            *(undefined1 *)puVar4 = 0;
            puVar4 = (undefined4 *)((int)puVar4 + 1);
          }
          iVar3 = (int)sVar11;
          param_1[0xc] = param_1[0xc] + uVar10;
          param_1[iVar3 * 0xd + 0x3b781] = (int)(uVar8 - uVar10) / 2;
          piVar5 = (int *)(*piVar2 + 8 + *(short *)((int)param_1 + iVar3 * 0x34 + 0xede26) * 0x24);
          *piVar5 = *piVar5 + (param_1[0xc] - iVar7 >> 1);
          iVar1 = *(int *)(*piVar2 + 4 + *(short *)((int)param_1 + iVar3 * 0x34 + 0xede26) * 0x24);
          piVar5 = (int *)(iVar1 + *(short *)(param_1 + iVar3 * 0xd + 0x3b78a) * 8);
          *piVar5 = *(int *)(iVar1 + *(short *)(param_1 + iVar3 * 0xd + 0x3b78a) * 8) +
                    (param_1[0xc] - iVar7 >> 1);
LAB_1002c0f7:
          *(undefined2 *)(param_1 + 0x10) = 1;
          return (uint)piVar5 & 0xffff0000;
        }
        puVar4 = (undefined4 *)(param_1[0xd] + iVar7);
        for (uVar9 = (param_1[iVar3 * 0xd + 0x3b781] & 0x7fffffff) >> 1; uVar9 != 0;
            uVar9 = uVar9 - 1) {
          *puVar4 = 0;
          puVar4 = puVar4 + 1;
        }
        for (uVar9 = uVar8 & 3; uVar9 != 0; uVar9 = uVar9 - 1) {
          *(undefined1 *)puVar4 = 0;
          puVar4 = (undefined4 *)((int)puVar4 + 1);
        }
        param_1[0xc] = param_1[0xc] + uVar8;
        *(undefined2 *)(iVar1 + 0x47708) = 0;
      }
      else {
        if (param_1[0x3ef5b] == 0) {
          puVar4 = FUN_1002d230(param_1,param_2,sVar11);
          param_1[0x3ef5b] = puVar4;
        }
        iVar6 = ((((int)param_1[iVar3 * 0xd + 0x3b77f] >> 1) + 10000) /
                (int)param_1[iVar3 * 0xd + 0x3b77f]) * param_1[iVar3 * 0xd + 0x3b781] + 0x32;
        piVar5 = (int *)-(iVar6 >> 0x1f);
        if (60000 < param_1[0xc] + 2000 + (iVar6 / 100) * 2) goto LAB_1002c0f7;
        if (param_1[0x3ef5b] == 0) {
          FUN_1002c8b0((short *)(param_1 + 0x43bbd),(int)(param_1 + iVar3 * 0xd + 0x3b77e),
                       (int)param_1,param_2);
        }
        if ((param_1[iVar3 * 0xd + 0x3b77e] == 100) &&
           ((((int)param_1[iVar3 * 0xd + 0x3b77f] >> 1) + 10000) /
            (int)param_1[iVar3 * 0xd + 0x3b77f] == 100)) {
          *(undefined2 *)((int)param_1 + 0x10b5e6) = 0;
        }
        else {
          if (param_1[0x3ef5b] == 0) {
            FUN_1002bbd0((int)(param_1 + 0x43b27),(int)(param_1 + iVar3 * 0xd + 0x3b77e),param_1,
                         (int)param_2);
          }
          iVar6 = FUN_1002bc60(param_1 + 0x42d7a,(short *)(param_1 + 0x43b27),
                               param_1 + iVar3 * 0xd + 0x3b77e);
          *(short *)((int)param_1 + 0x10b5e6) = (short)iVar6;
        }
        FUN_1002afb0((undefined2 *)param_1,sVar11);
        param_1[0x3ef5b] = 0;
      }
      piVar5 = (int *)(*piVar2 + 8 + *(short *)((int)param_1 + iVar3 * 0x34 + 0xede26) * 0x24);
      *piVar5 = *piVar5 + (param_1[0xc] - iVar7 >> 1);
      iVar6 = (int)sVar11;
      iVar3 = *(int *)(*piVar2 + 4 + *(short *)((int)param_1 + iVar3 * 0x34 + 0xede26) * 0x24);
      sVar11 = sVar11 + 1;
      *(int *)(iVar3 + *(short *)(param_1 + iVar6 * 0xd + 0x3b78a) * 8) =
           *(int *)(iVar3 + *(short *)(param_1 + iVar6 * 0xd + 0x3b78a) * 8) +
           (param_1[0xc] - iVar7 >> 1);
      *(short *)(iVar1 + 0x47706) = *(short *)(iVar1 + 0x47706) + 1;
      iVar7 = iVar1;
    } while (sVar11 < *(short *)(param_1 + 0x3ef5a));
  }
  if (*(short *)(param_1 + 0x3ef5a) == 0) {
    *(undefined2 *)(param_1 + 0x10) = 1;
    param_1[0x11] = 1;
    return 1;
  }
  *(undefined2 *)(param_1 + 0x10) = 0;
  return CONCAT22((short)((uint)iVar7 >> 0x10),1);
}



/* ===== FUN_1002d230 ===== */
/* Entry: 1002d230 */

/* WARNING: Removing unreachable block (ram,0x1002d663) */
/* WARNING: Removing unreachable block (ram,0x1002d6a8) */

undefined4 * __cdecl FUN_1002d230(undefined4 *param_1,int *param_2,short param_3)

{
  undefined4 *puVar1;
  short *psVar2;
  char cVar3;
  int iVar4;
  undefined4 *puVar5;
  short *psVar6;
  int iVar7;
  undefined2 extraout_var;
  uint uVar8;
  short *psVar9;
  short sVar10;
  undefined4 *extraout_EDX;
  undefined4 *extraout_EDX_00;
  undefined4 *extraout_EDX_01;
  undefined4 *puVar11;
  byte bVar12;
  short sVar13;
  short *psVar14;
  int iVar15;
  short sVar16;
  int iVar17;
  byte bStack0000000f;
  undefined1 local_f4 [8];
  undefined4 local_ec;
  byte local_d4;
  char local_ce;
  undefined1 local_c0 [8];
  undefined4 local_b8;
  uint local_a0;
  char local_9a;
  undefined4 local_8c [8];
  uint local_6c;
  int local_58;
  short *local_54;
  short *local_50;
  short *local_4c;
  int local_48;
  undefined4 *local_44;
  int local_40;
  short *local_3c;
  undefined4 *local_38;
  undefined4 *local_34;
  undefined4 *local_30;
  undefined4 *local_2c;
  undefined4 *local_28;
  undefined4 *local_24;
  short *local_20;
  short *local_1c;
  int local_18;
  int local_14;
  int local_10;
  int local_c;
  char local_6;
  char local_5;

  local_54 = (short *)(param_1 + 0x3ef5c);
  local_38 = param_1 + 0x429f4;
  local_20 = (short *)(param_1 + 0x42bb6);
  local_50 = (short *)(param_1 + 0x4166c);
  local_30 = param_1 + 0x42a8a;
  puVar1 = param_1 + 0x4354b;
  local_34 = param_1 + 0x42b20;
  local_40 = (int)param_1 + 0x10b132;
  local_1c = (short *)(param_1 + 0x42ce3);
  local_24 = param_1 + 0x43bbd;
  local_2c = param_1 + 0x43b27;
  local_14 = param_2[0x19];
  local_6 = '\x02';
  local_44 = param_1 + param_3 * 0xd + 0x3b77e;
  local_5 = '\x02';
  iVar4 = (int)*(short *)(param_1 + param_3 * 0xd + 0x3b782);
  iVar7 = iVar4 * 0x3f;
  local_c = param_1[*(short *)((int)param_1 + iVar4 * 2 + 0xae0c4) + iVar7 + 0x2ba44];
  iVar17 = local_c * 7 + local_14;
  if (*(short *)(param_1 + param_3 * 0xd + 0x3b782) == 0) {
    local_6 = '\0';
  }
  else if (((*(byte *)(iVar17 + 5) & 0x80) == 0x80) && (0x17 < (*(byte *)(iVar17 + 5) & 0x38))) {
    local_6 = '\0';
  }
  else {
    local_18 = param_1[*(short *)((int)param_1 + iVar4 * 2 + 0xae0c2) + iVar7 + 0x2ba05];
    local_10 = local_18 * 7 + local_14;
    cVar3 = *(char *)((int)param_1 + iVar4 * 6 + 0xec62c);
    if ((cVar3 == '\x02') && (local_18 == local_c)) {
      local_6 = '\0';
    }
    else if (local_18 == param_2[0x1b] + -1) {
      local_6 = '\0';
    }
    else if (cVar3 != '\x02') {
      if ((local_18 + 1 == local_c) && ((*(byte *)(local_10 + 6) & 0x80) != 0)) {
        local_6 = '\0';
      }
      else if ((*(byte *)(local_10 + 6) & 0x80) == 0) {
        local_6 = '\0';
      }
      else if ((*(byte *)(iVar17 + -1) & 0x80) == 0) {
        local_6 = '\x01';
      }
    }
  }
  if (iVar4 == param_1[0x3b189] + -1) {
    local_5 = '\0';
  }
  else if (((*(byte *)(iVar17 + 5) & 0x40) == 0x40) && (2 < (*(byte *)(iVar17 + 5) & 7))) {
    local_5 = '\0';
  }
  else {
    local_10 = param_1[*(short *)((int)param_1 + iVar4 * 2 + 0xae0c6) + iVar7 + 0x2ba83];
    cVar3 = *(char *)((int)param_1 + iVar4 * 6 + 0xec62c);
    if ((cVar3 == '\x01') && (local_c == local_10)) {
      local_5 = '\0';
    }
    else if (local_10 == 0) {
      local_5 = '\0';
    }
    else if (cVar3 != '\x01') {
      if ((local_c + 1 == local_10) && ((*(byte *)(iVar17 + 6) & 0x80) != 0)) {
        local_5 = '\0';
      }
      else if ((*(byte *)(local_10 * 7 + local_14 + -1) & 0x80) == 0) {
        local_5 = '\0';
      }
      else if ((*(byte *)(iVar17 + 6) & 0x80) == 0) {
        local_5 = '\x01';
      }
    }
  }
  if ((local_6 == '\0') && (local_5 == '\0')) {
    return (undefined4 *)0x0;
  }
  puVar5 = local_44;
  puVar11 = local_8c;
  local_4c = (short *)(param_1 + 0x402e4);
  for (iVar7 = 0xd; iVar7 != 0; iVar7 = iVar7 + -1) {
    *puVar11 = *puVar5;
    puVar5 = puVar5 + 1;
    puVar11 = puVar11 + 1;
  }
  FUN_1002c8b0((short *)(param_1 + 0x402e4),(int)local_8c,(int)param_1,param_2);
  FUN_1002bbd0((int)local_30,(int)local_8c,param_1,(int)param_2);
  if (local_6 == '\0') {
    local_c = 1;
    bStack0000000f = 0;
  }
  else {
    local_ce = *(char *)((int)param_1 + *(short *)(local_44 + 4) * 6 + 0xec62c);
    iVar7 = local_18;
    if (local_ce != '\x02') {
      iVar7 = local_18 + 1;
    }
    local_ec = local_8c[2];
    FUN_1002d1e0((int)local_f4,local_54,(int)local_38,iVar7,param_1,param_2);
    bStack0000000f = local_d4;
    if ((byte)local_6c < local_d4) {
      bStack0000000f = (byte)local_6c;
    }
    if (5 < bStack0000000f) {
      bStack0000000f = 5;
    }
    local_c = CONCAT22((short)((uint)puVar11 >> 0x10),(ushort)bStack0000000f) *
              CONCAT22(extraout_var,(short)local_6);
  }
  iVar7 = local_c;
  if (local_5 == '\0') {
    bVar12 = 0;
    sVar13 = 1;
    local_28 = (undefined4 *)0x0;
  }
  else {
    local_9a = *(char *)((int)param_1 + *(short *)(local_44 + 4) * 6 + 0xec62c);
    iVar4 = local_10;
    if (local_9a != '\x01') {
      iVar4 = local_10 + -1;
    }
    local_b8 = local_8c[2];
    FUN_1002d1e0((int)local_c0,local_50,(int)local_34,iVar4,param_1,param_2);
    bVar12 = (byte)local_a0;
    uVar8 = local_a0;
    if ((byte)local_6c < (byte)local_a0) {
      bVar12 = (byte)local_6c;
      uVar8 = local_6c;
    }
    local_28 = (undefined4 *)(uVar8 & 0xff);
    if (5 < bVar12) {
      bVar12 = 5;
      local_28 = (undefined4 *)0x5;
    }
    sVar13 = (ushort)bVar12 * (short)local_5;
  }
  param_2._0_2_ = sVar13 * (short)iVar7;
  puVar5 = local_2c;
  for (uVar8 = (local_6c & 0xff) >> 1; uVar8 != 0; uVar8 = uVar8 - 1) {
    *puVar5 = 0;
    puVar5 = puVar5 + 1;
  }
  sVar10 = 0;
  for (iVar7 = (local_6c & 1) << 1; iVar7 != 0; iVar7 = iVar7 + -1) {
    *(undefined1 *)puVar5 = 0;
    puVar5 = (undefined4 *)((int)puVar5 + 1);
  }
  iVar7 = (uint)bVar12 * local_c;
  sVar16 = (ushort)bStack0000000f * sVar13;
  do {
    if (sVar10 < (short)(ushort)bStack0000000f) {
      local_20[sVar10] = sVar16;
    }
    else {
      local_20[sVar10] = 0;
    }
    if (sVar10 < (short)(ushort)bVar12) {
      local_1c[(short)((ushort)(byte)local_6c - sVar10)] = (short)iVar7;
    }
    else {
      local_1c[(short)((ushort)(byte)local_6c - sVar10)] = 0;
    }
    iVar7 = iVar7 - local_c;
    sVar10 = sVar10 + 1;
    sVar16 = sVar16 - sVar13;
  } while (sVar10 <= (short)(ushort)(byte)local_6c);
  iVar7 = (short)(ushort)(byte)local_6c + 1;
  psVar6 = local_20;
  do {
    psVar9 = (short *)(((int)local_1c - (int)local_20) + (int)psVar6);
    sVar13 = *psVar6;
    psVar6 = psVar6 + 1;
    iVar7 = iVar7 + -1;
    *(short *)((local_40 - (int)local_20) + -2 + (int)psVar6) = ((short)param_2 - *psVar9) - sVar13;
  } while (iVar7 != 0);
  if ((byte)local_6c != 0) {
    local_3c = (short *)((int)local_20 - (int)local_1c);
    local_58 = (int)local_38 - (int)local_20;
    local_14 = 0;
    local_48 = (int)local_2c - (int)local_20;
    local_c = (int)(short)(ushort)(byte)local_6c;
    psVar9 = local_20;
    psVar6 = local_1c;
    do {
      psVar14 = psVar6 + 1;
      if ((*psVar9 == 0) && (*psVar14 == 0)) {
        *(undefined2 *)(local_48 + (int)psVar9) =
             *(undefined2 *)(((int)local_30 - (int)local_20) + (int)psVar9);
      }
      else {
        psVar2 = (short *)(((int)local_30 - (int)local_20) + (int)psVar9);
        iVar7 = ((int)*(short *)((int)local_3c + (int)psVar14 + (local_40 - (int)local_20)) +
                (int)*(short *)((int)psVar2 + (local_40 - (int)local_30))) * (int)*psVar2;
        if (0 < *psVar9) {
          iVar7 = iVar7 + ((int)*psVar9 + (int)*(short *)((int)local_3c + (int)psVar14)) *
                          (int)*(short *)(local_58 + (int)psVar9);
        }
        if (0 < *psVar14) {
          iVar7 = iVar7 + (int)*(short *)((int)local_34 +
                                         (((local_a0 & 0xff) - (local_6c & 0xff)) + local_14) * 2) *
                          ((int)*psVar6 + (int)*psVar14);
        }
        *(short *)(local_48 + (int)psVar9) =
             (short)((iVar7 + (short)param_2) / ((short)param_2 * 2));
      }
      psVar9 = psVar9 + 1;
      local_14 = local_14 + 1;
      local_c = local_c + -1;
      psVar6 = psVar14;
    } while (local_c != 0);
  }
  puVar5 = local_24;
  for (iVar7 = 5000; iVar7 != 0; iVar7 = iVar7 + -1) {
    *puVar5 = 0;
    puVar5 = puVar5 + 1;
  }
  iVar7 = 0;
  local_18 = 0;
  local_c = 0;
  local_14 = 0;
  _param_3 = (undefined4 *)0x0;
  if (local_5 == '\0') {
    local_a0 = local_a0 & 0xffffff00;
    iVar4 = 0;
    local_28 = (undefined4 *)0x0;
  }
  else {
    iVar17 = (local_a0 & 0xff) - (int)local_28;
    iVar4 = 0;
    if (0 < iVar17) {
      local_28 = (undefined4 *)0x0;
      do {
        iVar7 = iVar7 + CONCAT22((short)((uint)local_28 >> 0x10),
                                 *(undefined2 *)((int)local_34 + (int)local_28 * 2));
        iVar4 = iVar4 + 1;
        local_28 = (undefined4 *)(int)(short)iVar4;
        local_18 = iVar7;
      } while ((int)local_28 < iVar17);
    }
  }
  iVar7 = CONCAT22((short)((uint)iVar4 >> 0x10),(ushort)(byte)local_6c);
  local_10 = 0;
  if ((byte)local_6c != 0) {
    do {
      iVar17 = local_c;
      iVar4 = local_10;
      puVar5 = local_2c;
      iVar15 = (int)(short)local_10;
      local_28 = (undefined4 *)
                 ((CONCAT22((short)((uint)local_28 >> 0x10),(ushort)(byte)local_a0) - iVar7) +
                 local_10);
      if ((local_20[iVar15] == 0) && (local_1c[iVar15 + 1] == 0)) {
        puVar5 = FUN_10063f30((undefined4 *)((int)local_24 + param_3 * 2),
                              (undefined4 *)(local_4c + (short)local_c),
                              (int)*(short *)((int)local_2c + iVar15 * 2) << 1);
        local_c = iVar17 + CONCAT22((short)((uint)puVar5 >> 0x10),
                                    *(undefined2 *)((int)local_30 + iVar15 * 2));
        local_28 = (undefined4 *)
                   ((int)_param_3 +
                   CONCAT22((short)((uint)local_30 >> 0x10),
                            *(undefined2 *)((int)local_2c + iVar15 * 2)));
        puVar5 = local_2c;
        local_10 = iVar4;
        _param_3 = local_28;
      }
      else {
        local_3c = local_4c + (short)local_c;
        FUN_1002d010((int)puVar1,(int)local_3c,*(short *)((int)local_30 + iVar15 * 2),
                     *(short *)((int)local_2c + iVar15 * 2));
        sVar13 = 0;
        if (0 < *(short *)((int)puVar5 + iVar15 * 2)) {
          do {
            iVar7 = ((int)*(short *)((int)puVar1 + sVar13 * 2) *
                    (int)*(short *)(local_40 + iVar15 * 2)) / (int)(short)param_2;
            if ((0x7ffe < iVar7) || (iVar7 < -0x7fff)) {
              iVar7 = ((iVar7 < 0x7fff) - 1 & 0xfffd) - 0x7fff;
            }
            psVar6 = (short *)((int)local_24 + ((int)sVar13 + (int)param_3) * 2);
            *psVar6 = *psVar6 + (short)iVar7;
            sVar13 = sVar13 + 1;
          } while (sVar13 < *(short *)((int)puVar5 + iVar15 * 2));
        }
        if (0 < local_20[iVar15]) {
          FUN_1002d010((int)puVar1,(int)(local_54 + (short)local_14),
                       *(short *)((int)local_38 + iVar15 * 2),*(short *)((int)puVar5 + iVar15 * 2));
          sVar13 = 0;
          if (0 < *(short *)((int)puVar5 + iVar15 * 2)) {
            do {
              iVar7 = ((int)*(short *)((int)puVar1 + sVar13 * 2) * (int)local_20[iVar15]) /
                      (int)(short)param_2;
              if ((0x7ffe < iVar7) || (iVar7 < -0x7fff)) {
                iVar7 = ((iVar7 < 0x7fff) - 1 & 0xfffd) - 0x7fff;
              }
              psVar6 = (short *)((int)local_24 + ((int)sVar13 + (int)param_3) * 2);
              *psVar6 = *psVar6 + (short)iVar7;
              sVar13 = sVar13 + 1;
            } while (sVar13 < *(short *)((int)puVar5 + iVar15 * 2));
          }
        }
        if (0 < local_1c[iVar15]) {
          FUN_1002d010((int)puVar1,(int)(local_50 + (short)local_18),
                       *(short *)((int)local_34 + (short)local_28 * 2),
                       *(short *)((int)puVar5 + iVar15 * 2));
          sVar13 = 0;
          if (0 < *(short *)((int)puVar5 + iVar15 * 2)) {
            do {
              iVar7 = ((int)*(short *)((int)puVar1 + sVar13 * 2) * (int)local_1c[iVar15]) /
                      (int)(short)param_2;
              if ((0x7ffe < iVar7) || (iVar7 < -0x7fff)) {
                iVar7 = ((iVar7 < 0x7fff) - 1 & 0xfffd) - 0x7fff;
              }
              psVar6 = (short *)((int)local_24 + ((int)sVar13 + (int)param_3) * 2);
              *psVar6 = *psVar6 + (short)iVar7;
              sVar13 = sVar13 + 1;
            } while (sVar13 < *(short *)((int)puVar5 + iVar15 * 2));
          }
        }
        FUN_1002cdf0(puVar1,local_3c,*(short *)((int)local_30 + iVar15 * 2),
                     *(short *)((int)puVar5 + iVar15 * 2));
        sVar13 = 0;
        puVar11 = extraout_EDX;
        if (0 < *(short *)((int)puVar5 + iVar15 * 2)) {
          do {
            iVar7 = ((int)*(short *)((int)puVar1 + sVar13 * 2) *
                    (int)*(short *)(local_40 + 2 + iVar15 * 2)) / (int)(short)param_2;
            if ((0x7ffe < iVar7) || (iVar7 < -0x7fff)) {
              iVar7 = ((iVar7 < 0x7fff) - 1 & 0xfffd) - 0x7fff;
            }
            psVar6 = (short *)((int)local_24 + ((int)sVar13 + (int)param_3) * 2);
            *psVar6 = *psVar6 + (short)iVar7;
            sVar13 = sVar13 + 1;
            puVar11 = local_24;
          } while (sVar13 < *(short *)((int)puVar5 + iVar15 * 2));
        }
        if (0 < local_20[iVar15 + 1]) {
          FUN_1002cdf0(puVar1,local_54 + (short)local_14,*(short *)((int)local_38 + iVar15 * 2),
                       *(short *)((int)puVar5 + iVar15 * 2));
          sVar13 = 0;
          puVar11 = extraout_EDX_00;
          if (0 < *(short *)((int)puVar5 + iVar15 * 2)) {
            local_3c = (short *)(int)(short)param_2;
            local_48 = (int)param_3;
            do {
              iVar7 = ((int)*(short *)((int)puVar1 + sVar13 * 2) * (int)local_20[iVar15 + 1]) /
                      (int)local_3c;
              if ((0x7ffe < iVar7) || (iVar7 < -0x7fff)) {
                iVar7 = ((iVar7 < 0x7fff) - 1 & 0xfffd) - 0x7fff;
              }
              psVar6 = (short *)((int)local_24 + (sVar13 + local_48) * 2);
              *psVar6 = *psVar6 + (short)iVar7;
              sVar13 = sVar13 + 1;
              puVar11 = local_24;
            } while (sVar13 < *(short *)((int)puVar5 + iVar15 * 2));
          }
        }
        if (0 < local_1c[iVar15 + 1]) {
          FUN_1002cdf0(puVar1,local_50 + (short)local_18,
                       *(short *)((int)local_34 + (short)local_28 * 2),
                       *(short *)((int)puVar5 + iVar15 * 2));
          sVar13 = 0;
          puVar11 = extraout_EDX_01;
          if (0 < *(short *)((int)puVar5 + iVar15 * 2)) {
            do {
              iVar7 = ((int)*(short *)((int)puVar1 + sVar13 * 2) * (int)local_1c[iVar15 + 1]) /
                      (int)(short)param_2;
              if ((0x7ffe < iVar7) || (iVar7 < -0x7fff)) {
                iVar7 = ((iVar7 < 0x7fff) - 1 & 0xfffd) - 0x7fff;
              }
              psVar6 = (short *)((int)local_24 + ((int)sVar13 + (int)param_3) * 2);
              *psVar6 = *psVar6 + (short)iVar7;
              sVar13 = sVar13 + 1;
              puVar11 = local_24;
            } while (sVar13 < *(short *)((int)puVar5 + iVar15 * 2));
          }
        }
        if (0 < local_20[iVar15]) {
          puVar11 = (undefined4 *)
                    CONCAT22((short)((uint)puVar11 >> 0x10),
                             *(undefined2 *)((int)local_38 + iVar15 * 2));
          local_14 = local_14 + (int)puVar11;
        }
        psVar6 = local_1c;
        if (0 < local_1c[iVar15 + 1]) {
          psVar6 = (short *)CONCAT22((short)((uint)local_1c >> 0x10),
                                     *(undefined2 *)((int)local_34 + (short)local_28 * 2));
          local_18 = local_18 + (int)psVar6;
          puVar11 = local_34;
        }
        local_28 = (undefined4 *)
                   CONCAT22((short)((uint)psVar6 >> 0x10),*(undefined2 *)((int)puVar5 + iVar15 * 2))
        ;
        local_c = local_c + CONCAT22((short)((uint)puVar11 >> 0x10),
                                     *(undefined2 *)((int)local_30 + iVar15 * 2));
        puVar5 = (undefined4 *)((int)_param_3 + (int)local_28);
        _param_3 = puVar5;
      }
      iVar7 = CONCAT22((short)((uint)puVar5 >> 0x10),(ushort)(byte)local_6c);
      local_10 = local_10 + 1;
    } while ((short)local_10 < (short)(ushort)(byte)local_6c);
  }
  local_44[3] = (int)param_3;
  *(char *)(local_44 + 8) = (char)local_6c;
  *(undefined2 *)((int)local_44 + 0x22) = *(undefined2 *)local_2c;
  *(undefined2 *)(local_44 + 9) = *(undefined2 *)((int)local_2c + (local_6c & 0xff) * 2 + -2);
  return local_44;
}



/* ===== FUN_10024680 ===== */
/* Entry: 10024680 */

uint __cdecl FUN_10024680(short param_1,undefined4 *param_2,int param_3)

{
  int iVar1;
  byte *pbVar2;
  short *psVar3;
  short sVar4;
  int iVar5;
  bool bVar6;
  short sVar7;
  short sVar8;
  int iVar9;
  undefined3 extraout_var;
  undefined4 uVar10;
  byte *pbVar11;
  uint uVar12;
  byte bStack00000007;
  int local_20;
  uint local_18;
  int local_14;
  uint local_10;

  iVar5 = param_2[0x13];
  uVar12 = 0;
  iVar9 = param_1 * 0x10 + iVar5;
  local_18 = 0;
  local_14 = 0;
  sVar4 = *(short *)(iVar9 + 0xc);
  iVar9 = (int)*(short *)(iVar9 + 0xe);
  if (0 < iVar9) {
    local_20 = sVar4 * 0x1e0;
    pbVar11 = (byte *)(sVar4 * 0x3c0 + 0x6e1 + iVar5);
    do {
      pbVar11[-0xffffffff00000003] = 0;
      pbVar11[-0xffffffff00000002] = 0;
      *(uint *)(pbVar11 + -0x89) = (int)param_2 + uVar12 * 6 + 0xec628;
      local_10 = 0;
      if (*pbVar11 != 0) {
        bStack00000007 = (char)local_14 + (char)sVar4;
        do {
          pbVar2 = (byte *)((int)param_2 + uVar12 * 6 + 0xec628);
          *pbVar2 = bStack00000007;
          pbVar2[2] = (byte)local_10;
          bVar6 = FUN_10017010((local_10 & 0xff) * 7 + 0x6e2 + iVar5 + (uint)*pbVar2 * 0x3c0);
          pbVar2[5] = bVar6;
          sVar7 = 0;
          sVar8 = 0;
          iVar1 = (uint)*pbVar2 * 0x3c0 + iVar5;
          if (*(char *)((uint)*pbVar2 * 0x3c0 + 0x6e0 + iVar5) != '\0') {
            do {
              if ((sVar8 <= (short)(ushort)pbVar2[2]) &&
                 (uVar12 = local_18,
                 (int)(uint)pbVar2[2] <
                 (int)((uint)*(byte *)(*(int *)(iVar1 + 0x654) + 0x1d + sVar7 * 0x1e) + (int)sVar8))
                 ) break;
              sVar8 = sVar8 + (ushort)*(byte *)(*(int *)(iVar1 + 0x654) + 0x1d + sVar7 * 0x1e);
              sVar7 = sVar7 + 1;
            } while (sVar7 < (short)(ushort)*(byte *)(iVar1 + 0x6e0));
          }
          pbVar2[1] = (byte)sVar7;
          if (0 < *(short *)(iVar5 + 0x65c + (local_20 + local_10) * 2)) {
            bVar6 = FUN_10024060(uVar12,param_2,param_3);
            sVar8 = (short)CONCAT31(extraout_var,bVar6);
            if (sVar8 < 1) {
              uVar10 = FUN_100242a0(uVar12,(int)param_2,param_3);
              sVar8 = (short)uVar10;
            }
            psVar3 = (short *)((uint)*pbVar2 * 0x3c0 + 0x6de + iVar5);
            *psVar3 = *psVar3 + sVar8;
            uVar12 = uVar12 + (int)sVar8;
            local_18 = uVar12;
          }
          local_10 = local_10 + 1;
        } while ((int)local_10 < (int)(uint)*pbVar11);
      }
      pbVar11 = pbVar11 + 0x3c0;
      local_14 = local_14 + 1;
      local_20 = local_20 + 0x1e0;
    } while (local_14 < iVar9);
  }
  param_2[0x3b189] = uVar12;
  return uVar12;
}



/* ===== FUN_10024060 ===== */
/* Entry: 10024060 */

bool __cdecl FUN_10024060(uint param_1,undefined4 *param_2,int param_3)

{
  byte *pbVar1;
  undefined4 *puVar2;
  byte bVar3;
  char cVar4;
  ushort uVar5;
  ushort uVar6;
  undefined4 *puVar7;
  short sVar8;
  short sVar9;
  int iVar10;
  uint uVar11;
  int iVar12;
  undefined4 local_18;
  byte local_12;
  ushort *local_10;
  int local_c;
  ushort *local_8;

  puVar7 = param_2;
  iVar12 = param_2[0x13];
  pbVar1 = (byte *)((int)param_2 + (param_1 * 3 + 0x76314) * 2);
  iVar10 = (uint)*pbVar1 * 0x3c0;
  bVar3 = *(byte *)(iVar10 + 0x6de + iVar12);
  pbVar1[4] = 0;
  pbVar1[3] = bVar3;
  cVar4 = *(char *)(iVar10 + 0x92b + iVar12);
  puVar2 = (undefined4 *)((uint)pbVar1[2] * 7 + iVar10 + 0x6e2 + iVar12);
  *(undefined2 *)(param_2 + 0x3b188) = 0;
  *(undefined2 *)(param_2 + 0x3b05b) = 0;
  if (cVar4 == '\f') {
    FUN_10063f30(&local_18,puVar2,7);
    local_12 = local_12 & 0xdf;
    FUN_10018770(&local_18,param_1,(int)param_2,param_3,1);
    if (*(short *)(param_2 + 0x3b188) == 0) {
      return false;
    }
    sVar8 = 2;
  }
  else {
    sVar8 = 1;
  }
  FUN_10018770(puVar2,param_1,(int)param_2,param_3,sVar8);
  if (*(short *)(param_2 + 0x3b188) == 0) {
    return false;
  }
  sVar8 = FUN_10024010(*(short *)(param_2 + 0x3b188),param_2 + 0x3b05c);
  *(short *)(param_2 + 0x3b188) = sVar8;
  uVar11 = FUN_10023af0(0x1e,10000,param_2 + 0x3b05c,(int)puVar2,param_2,param_3);
  *(short *)(param_2 + 0x3b188) = (short)uVar11;
  sVar8 = FUN_10024010(*(short *)(param_2 + 0x3b05b),param_2 + 0x3b03d);
  *(short *)(param_2 + 0x3b05b) = sVar8;
  local_8 = (ushort *)(param_2 + param_1 * 0x3f + 0x2ba25);
  local_10 = local_8 + 2;
  FUN_10063f30((undefined4 *)local_10,param_2 + 0x3b03d,(int)sVar8 << 2);
  uVar5 = *(ushort *)(param_2 + 0x3b188);
  uVar6 = *(ushort *)(param_2 + 0x3b05b);
  if ((short)uVar6 < (short)uVar5) {
    local_c = 0;
    param_2 = (undefined4 *)(uint)uVar6;
    if (0 < (short)uVar5) {
      do {
        sVar8 = *(short *)(puVar7 + 0x3b05b);
        sVar9 = 0;
        if (0 < sVar8) {
          do {
            if (puVar7[(short)local_c + 0x3b05c] == puVar7[sVar9 + 0x3b03d]) break;
            sVar9 = sVar9 + 1;
          } while (sVar9 < sVar8);
        }
        if (sVar9 == sVar8) {
          puVar7[param_1 * 0x3f + (int)(short)(ushort)param_2 + 0x2ba26] =
               puVar7[(short)local_c + 0x3b05c];
          param_2 = (undefined4 *)((int)param_2 + 1);
        }
        local_c = local_c + 1;
        uVar6 = (ushort)param_2;
      } while ((short)local_c < *(short *)(puVar7 + 0x3b188));
    }
    param_2._0_2_ = uVar6;
    *local_8 = (ushort)param_2;
    local_8[1] = *(ushort *)(puVar7 + 0x3b05b);
  }
  else {
    *local_8 = uVar5;
    local_8[1] = *(ushort *)(param_2 + 0x3b05b);
  }
  iVar12 = FUN_10023060((int *)local_10,*local_8,param_3);
  if ((cVar4 == '\b') && (0 < iVar12)) {
    return true;
  }
  return 9 < iVar12;
}



/* ===== FUN_100242a0 ===== */
/* Entry: 100242a0 */

undefined4 __cdecl FUN_100242a0(int param_1,int param_2,int param_3)

{
  byte *pbVar1;
  undefined4 *puVar2;
  short sVar3;
  short sVar4;
  uint uVar5;
  int iVar6;
  int local_2c;
  int local_28;
  undefined4 local_20;
  byte local_1a;
  undefined2 *local_18;
  byte *local_14;
  int local_10;
  int local_c;
  byte *local_8;

  iVar6 = param_1 * 3 + 0x76314;
  local_10 = *(int *)(param_2 + 0x4c);
  local_14 = (byte *)(param_2 + iVar6 * 2);
  puVar2 = (undefined4 *)
           ((uint)*local_14 * 0x3c0 + 0x6e2 +
           local_10 + (uint)*(byte *)(param_2 + 2 + iVar6 * 2) * 7);
  local_2c = param_2 + 0xec170;
  local_c = 0;
  do {
    sVar3 = (short)local_c;
    if (sVar3 == 0) {
      local_14[4] = 1;
    }
    else {
      local_8 = local_14 + sVar3 * 6;
      *(undefined4 *)local_8 = *(undefined4 *)local_14;
      *(undefined2 *)(local_8 + 4) = *(undefined2 *)(local_14 + 4);
      local_8[4] = 2;
    }
    iVar6 = (int)sVar3;
    pbVar1 = local_14 + iVar6 * 6;
    pbVar1[3] = *(char *)((uint)*pbVar1 * 0x3c0 + 0x6de + local_10) + (char)local_c;
    *(undefined2 *)(param_2 + 0xec620) = 0;
    *(undefined2 *)(param_2 + 0xec16c) = 0;
    if (*(char *)((uint)*local_14 * 0x3c0 + 0x92b + local_10) == '\f') {
      FUN_10063f30(&local_20,puVar2,7);
      local_1a = local_1a & 0x80;
      FUN_10018770(&local_20,iVar6 + param_1,param_2,param_3,1);
      FUN_10018770(puVar2,iVar6 + param_1,param_2,param_3,2);
    }
    else {
      FUN_10018770(puVar2,iVar6 + param_1,param_2,param_3,1);
    }
    sVar3 = FUN_10024010(*(short *)(param_2 + 0xec620),(int *)(param_2 + 0xec170));
    *(short *)(param_2 + 0xec620) = sVar3;
    local_28 = (int)sVar3;
    uVar5 = FUN_10023c70(0x1e,10000,(int *)(param_2 + 0xec170),(int)puVar2,(short)(char)pbVar1[4],
                         &local_2c,param_2,param_3);
    *(short *)(param_2 + 0xec620) = (short)uVar5;
    sVar3 = FUN_10024010(*(short *)(param_2 + 0xec16c),(int *)(param_2 + 0xec0f4));
    *(short *)(param_2 + 0xec16c) = sVar3;
    local_18 = (undefined2 *)(param_2 + 0xae894 + (iVar6 + param_1) * 0xfc);
    FUN_10063f30((undefined4 *)(local_18 + 2),(int *)(param_2 + 0xec0f4),(int)sVar3 << 2);
    sVar3 = 0;
    local_8 = (byte *)CONCAT22(local_8._2_2_,*(undefined2 *)(param_2 + 0xec16c));
    if (0 < *(short *)(param_2 + 0xec620)) {
      do {
        sVar4 = 0;
        if (0 < *(short *)(param_2 + 0xec16c)) {
          do {
            if (*(int *)(param_2 + 0xec170 + sVar3 * 4) == *(int *)(param_2 + 0xec0f4 + sVar4 * 4))
            break;
            sVar4 = sVar4 + 1;
          } while (sVar4 < *(short *)(param_2 + 0xec16c));
        }
        if (sVar4 == *(short *)(param_2 + 0xec16c)) {
          *(undefined4 *)(param_2 + 0xae898 + ((iVar6 + param_1) * 0x3f + (int)(short)local_8) * 4)
               = *(undefined4 *)(param_2 + 0xec170 + sVar3 * 4);
          local_8 = local_8 + 1;
        }
        sVar3 = sVar3 + 1;
      } while (sVar3 < *(short *)(param_2 + 0xec620));
    }
    *local_18 = (short)local_8;
    local_18[1] = *(undefined2 *)(param_2 + 0xec16c);
    local_c = local_c + 1;
    if (1 < (short)local_c) {
      return CONCAT22((short)((uint)local_c >> 0x10),2);
    }
  } while( true );
}



/* ===== FUN_10018770 ===== */
/* Entry: 10018770 */

undefined4 __cdecl
FUN_10018770(undefined4 *param_1,uint param_2,int param_3,int param_4,short param_5)

{
  byte bVar1;
  char cVar2;
  int iVar3;
  int iVar4;
  int iVar5;
  uint uVar6;
  uint uVar7;
  uint uVar8;
  byte bVar9;
  byte bVar10;
  int iVar11;
  short sVar12;
  char local_58;
  char local_57;
  char local_56;
  char local_55;
  char local_54;
  char local_53;
  int local_50;
  int local_4c;
  char *local_48;
  uint local_44;
  uint local_40;
  char *local_3c;
  short local_38;
  uint local_34;
  short local_30;
  uint local_2c;
  undefined4 local_28;
  int local_24;
  int local_20;
  byte *local_1c;
  uint local_18;
  uint local_14;
  byte local_10;
  undefined3 uStack_f;
  uint local_c;
  char local_7;
  char local_6;
  char local_5;

  iVar3 = param_3;
  local_28 = 0;
  local_38 = *(short *)(param_3 + 0xec620);
  local_1c = (byte *)(param_3 + (param_2 * 3 + 0x76314) * 2);
  bVar1 = *(byte *)((int)param_1 + 5);
  _local_10 = CONCAT31(uStack_f,bVar1 >> 7);
  local_c = CONCAT31(local_c._1_3_,bVar1 >> 6) & 0xffffff01;
  local_2c = CONCAT31(local_2c._1_3_,bVar1) & 0xffffff07;
  bVar9 = bVar1 >> 3;
  bVar10 = bVar9 & 7;
  local_34 = CONCAT31(local_34._1_3_,bVar9) & 0xffffff07;
  local_3c = (&PTR_DAT_1007c4dc)[(uint)(bVar1 >> 7) * 5 + (bVar9 & 7)];
  local_50 = (int)(short)*local_3c;
  local_48 = (&PTR_DAT_1007c504)[(bVar1 & 7) + (bVar1 >> 6 & 1) * 8];
  local_24 = 0;
  local_30 = (short)*local_48;
  if (0 < local_50) {
    local_44 = (uint)*(byte *)((uint)*local_1c * 0x3c0 + (uint)local_1c[2] + 0x8a9 +
                              *(int *)(param_3 + 0x4c));
    local_40 = local_44 / 10;
    do {
      cVar2 = local_3c[local_24 + 1];
      local_14 = CONCAT31(local_14._1_3_,cVar2);
      if (cVar2 == '\0') {
        param_3._3_1_ = '\0';
LAB_10018876:
        iVar11 = 1;
      }
      else {
        param_3._3_1_ = (char)_local_10;
        if (cVar2 == '\x01') goto LAB_10018876;
        iVar11 = 0;
      }
      if ((bVar10 == 1) || (bVar10 == 0)) {
        iVar4 = 1;
      }
      else {
        iVar4 = 0;
      }
      local_6 = iVar11 == iVar4;
      if ((((local_40 == 0) || ((bool)local_6)) || (local_1c[4] == 2)) &&
         (local_20 = 0, 0 < local_30)) {
        local_4c = (int)local_44 % 10;
        do {
          iVar11 = local_20;
          cVar2 = local_48[local_20 + 1];
          local_18 = CONCAT31(local_18._1_3_,cVar2);
          if (cVar2 == '\0') {
            local_5 = '\0';
LAB_100188eb:
            iVar4 = 1;
          }
          else {
            local_5 = (char)local_c;
            if (cVar2 == '\x01') goto LAB_100188eb;
            iVar4 = 0;
          }
          if (((char)local_2c == '\x01') || ((char)local_2c == '\0')) {
            iVar5 = 1;
          }
          else {
            iVar5 = 0;
          }
          local_7 = iVar4 == iVar5;
          if (((local_4c == 0) || ((bool)local_7)) || (local_1c[4] == 1)) {
            FUN_10063f30((undefined4 *)&local_58,param_1,7);
            if (local_24 == 0) {
              if (iVar11 != 0) goto LAB_100189f5;
              uVar6 = FUN_10023f90((int)&local_58,param_2,iVar3,param_4,param_5);
              sVar12 = (short)uVar6;
              if (((param_5 == 1) && (local_38 == 0)) && (0 < sVar12)) {
                if (*(short *)(iVar3 + 0xec620) < 0x1f) {
                  *(short *)(iVar3 + 0xec16c) = *(short *)(iVar3 + 0xec620);
                }
                else {
                  *(undefined2 *)(iVar3 + 0xec16c) = 0x1e;
                }
                FUN_10063f30((undefined4 *)(iVar3 + 0xec0f4),(undefined4 *)(iVar3 + 0xec170),
                             (int)*(short *)(iVar3 + 0xec16c) << 2);
              }
            }
            else {
              if (((0 < local_24) && (local_10 == '\x01')) && (2 < (byte)local_14)) {
                local_58 = (&DAT_1007c4a4)[local_14 & 0xff];
                local_57 = local_58;
              }
LAB_100189f5:
              if (((0 < iVar11) && ((char)local_c == '\x01')) && (2 < (byte)local_18)) {
                local_55 = (&DAT_1007c4a4)[local_18 & 0xff];
                local_54 = local_55;
              }
              local_53 = ((param_3._3_1_ * '\x02' + local_5) * '\b' + (byte)local_14) * '\b' +
                         (byte)local_18;
              uVar6 = FUN_10023f90((int)&local_58,param_2,iVar3,param_4,0);
              sVar12 = (short)uVar6;
            }
            if (((sVar12 != 0) || (local_10 != '\0')) || ((char)local_c != '\0')) {
              if ((&DAT_1007c2c4)[*(byte *)((int)param_1 + 2)] == '\0') {
                cVar2 = (&DAT_1007c2c4)[*(byte *)((int)param_1 + 1)];
                if (cVar2 == '\0') {
LAB_10018b96:
                  if ((&DAT_1007c2c4)[*(byte *)((int)param_1 + 3)] != '\0') {
                    local_55 = (&DAT_1007c2c4)[*(byte *)((int)param_1 + 3)];
                    uVar6 = FUN_10023f90((int)&local_58,param_2,iVar3,param_4,0);
                    sVar12 = sVar12 + (short)uVar6;
                  }
                }
                else if ((&DAT_1007c2c4)[*(byte *)((int)param_1 + 3)] == '\0') {
                  if (cVar2 == '\0') goto LAB_10018b96;
                  local_57 = cVar2;
                  uVar6 = FUN_10023f90((int)&local_58,param_2,iVar3,param_4,0);
                  sVar12 = sVar12 + (short)uVar6;
                }
                else {
                  local_57 = cVar2;
                  uVar6 = FUN_10023f90((int)&local_58,param_2,iVar3,param_4,0);
                  local_55 = (&DAT_1007c2c4)[*(byte *)((int)param_1 + 3)];
                  uVar7 = FUN_10023f90((int)&local_58,param_2,iVar3,param_4,0);
                  FUN_10063f30((undefined4 *)&local_58,param_1,7);
                  local_55 = (&DAT_1007c2c4)[*(byte *)((int)param_1 + 3)];
                  uVar8 = FUN_10023f90((int)&local_58,param_2,iVar3,param_4,0);
                  sVar12 = sVar12 + (short)uVar6 + (short)uVar7 + (short)uVar8;
                }
              }
              else if (((((char)local_c != '\0') || (*(byte *)((int)param_1 + 2) != 0x39)) ||
                       ((char)(&DAT_1007bb70)[*(byte *)((int)param_1 + 3)] < '\x01')) &&
                      (*(char *)((uint)*local_1c * 0x3c0 + (uint)local_1c[2] + 0x8ea +
                                *(int *)(iVar3 + 0x4c)) == '\0')) {
                local_56 = (&DAT_1007c2c4)[*(byte *)((int)param_1 + 2)];
                uVar6 = FUN_10023f90((int)&local_58,param_2,iVar3,param_4,0);
                sVar12 = sVar12 + (short)uVar6;
              }
              if ((((short)local_28 == 0) && (0 < sVar12)) &&
                 ((local_6 == '\x01' && (local_7 == '\x01')))) {
                local_28 = 1;
              }
            }
          }
          local_20 = local_20 + 1;
        } while (local_20 < local_30);
        bVar10 = (byte)local_34;
      }
      local_24 = local_24 + 1;
    } while (local_24 < local_50);
  }
  if ((param_5 != 2) && ((short)local_28 == 0)) {
    *(short *)(iVar3 + 0xec620) = local_38;
    return 0;
  }
  if ((param_5 != 0) && ((&DAT_1007c324)[*(byte *)((int)param_1 + 2)] != '\0')) {
    FUN_10063f30((undefined4 *)&local_58,param_1,7);
    local_56 = (&DAT_1007c324)[*(byte *)((int)param_1 + 2)];
    FUN_10018770((undefined4 *)&local_58,param_2,iVar3,param_4,0);
  }
  return 1;
}



/* ===== FUN_10023af0 ===== */
/* Entry: 10023af0 */

uint __cdecl
FUN_10023af0(uint param_1,int param_2,int *param_3,int param_4,undefined4 *param_5,int param_6)

{
  ushort uVar1;
  int iVar2;
  int iVar3;
  undefined4 *puVar4;
  int iVar5;
  int iVar6;
  short sVar7;
  int *piVar9;
  byte local_c [8];
  uint uVar8;

  puVar4 = param_5;
  iVar3 = FUN_10023060(param_5 + 0x3b05c,*(short *)(param_5 + 0x3b188),param_6);
  if (((int)*(short *)(param_5 + 0x3b188) <= (int)param_1) && (iVar3 <= param_2)) {
    puVar4 = FUN_10063f30(param_3,param_5 + 0x3b05c,*(short *)(param_5 + 0x3b188) * 4);
    return CONCAT22((short)((uint)puVar4 >> 0x10),*(undefined2 *)(param_5 + 0x3b188));
  }
  FUN_10016ea0(param_4,local_c);
  iVar3 = 0;
  if (0 < *(short *)(param_5 + 0x3b188)) {
    piVar9 = param_5 + 0x11deb;
    param_5 = param_5 + 0x3b05c;
    do {
      if (9999 < iVar3) break;
      *(undefined4 *)(*piVar9 + 8) = *param_5;
      iVar5 = FUN_10023a70(local_c,(byte *)(*(int *)(*piVar9 + 8) * 5 + *(int *)(param_6 + 0x90)),0)
      ;
      iVar3 = iVar3 + 1;
      param_5 = param_5 + 1;
      *(float *)(*piVar9 + 4) = (float)iVar5;
      piVar9 = piVar9 + 1;
    } while (iVar3 < *(short *)(puVar4 + 0x3b188));
  }
  uVar1 = *(ushort *)(puVar4 + 0x3b188);
  FUN_1001b5f0((int)(short)uVar1,(int)(puVar4 + 0x11deb));
  uVar8 = (uint)uVar1;
  if ((int)param_1 < (int)*(short *)(puVar4 + 0x3b188)) {
    uVar8 = param_1;
  }
  sVar7 = (short)uVar8;
  iVar3 = 0;
  iVar5 = 0;
  if (0 < sVar7) {
    iVar6 = (int)(puVar4 + 0x11deb) - (int)param_3;
    do {
      iVar2 = *(int *)(*(int *)(iVar6 + (int)param_3) + 8);
      *param_3 = iVar2;
      iVar3 = iVar3 + (uint)*(ushort *)(*(int *)(param_6 + 0x8c) + iVar2 * 2);
      if (param_2 < iVar3) break;
      iVar5 = iVar5 + 1;
      param_3 = param_3 + 1;
    } while (iVar5 < sVar7);
  }
  if (sVar7 <= iVar5) {
    return (int)sVar7;
  }
  return iVar5 + 1;
}



/* ===== FUN_10023c70 ===== */
/* Entry: 10023c70 */

uint __cdecl
FUN_10023c70(int param_1,int param_2,int *param_3,int param_4,short param_5,int *param_6,int param_7
            ,int param_8)

{
  int iVar1;
  int iVar2;
  undefined4 *puVar3;
  int iVar4;
  short sVar5;
  byte local_c [8];

  iVar1 = param_7 + 0x477ac;
  iVar2 = FUN_10023060((int *)*param_6,(short)param_6[1],param_8);
  if ((param_6[1] <= param_1) && (iVar2 <= param_2)) {
    puVar3 = FUN_10063f30(param_3,(undefined4 *)*param_6,param_6[1] << 2);
    return CONCAT22((short)((uint)puVar3 >> 0x10),(short)param_6[1]);
  }
  FUN_10016ea0(param_4,local_c);
  sVar5 = 0;
  iVar2 = 0;
  if (0 < param_6[1]) {
    do {
      *(undefined4 *)(*(int *)(iVar1 + sVar5 * 4) + 8) = *(undefined4 *)(*param_6 + iVar2 * 4);
      iVar4 = FUN_10023a70(local_c,(byte *)(*(int *)(*param_6 + iVar2 * 4) * 5 +
                                           *(int *)(param_8 + 0x90)),param_5);
      sVar5 = sVar5 + 1;
      *(float *)(*(int *)(iVar1 + iVar2 * 4) + 4) = (float)iVar4;
      if (9999 < sVar5) break;
      iVar2 = iVar2 + 1;
    } while (iVar2 < param_6[1]);
  }
  FUN_1001b5f0((int)sVar5,iVar1);
  if (param_1 < sVar5) {
    sVar5 = (short)param_1;
  }
  iVar4 = 0;
  param_1 = 0;
  iVar2 = 0;
  if (0 < sVar5) {
    iVar1 = iVar1 - (int)param_3;
    do {
      iVar2 = *(int *)(*(int *)(iVar1 + (int)param_3) + 8);
      *param_3 = iVar2;
      iVar4 = iVar4 + (uint)*(ushort *)(*(int *)(param_8 + 0x8c) + iVar2 * 2);
      iVar2 = param_1;
      if (param_2 < iVar4) break;
      iVar2 = param_1 + 1;
      param_3 = param_3 + 1;
      param_1 = iVar2;
    } while (iVar2 < sVar5);
  }
  if (sVar5 <= iVar2) {
    return CONCAT22((short)((uint)iVar2 >> 0x10),sVar5);
  }
  return iVar2 + 1;
}



/* ===== FUN_1002c8b0 ===== */
/* Entry: 1002c8b0 */

undefined4 __cdecl FUN_1002c8b0(short *param_1,int param_2,int param_3,int *param_4)

{
  undefined4 *puVar1;
  int iVar2;
  int iVar3;

  iVar3 = *(int *)(param_2 + 0xc);
  FUN_10025440(*(int *)(*param_4 + *(short *)(param_2 + 0x12) * 4),**(int **)(param_2 + 0x18),0,
               (char *)(param_3 + 0x118b34),1,(int)(short)(*(int **)(param_2 + 0x18))[4]);
  puVar1 = (undefined4 *)(param_3 + 0x113d14);
  FUN_10001b30((char *)(param_3 + 0x118b34),puVar1,(int)*(short *)(*(int *)(param_2 + 0x18) + 0x10),
               &param_4);
  if (*(char *)(param_2 + 0x26) == '\x02') {
    puVar1 = (undefined4 *)
             (param_3 + 0x113d14 +
             (short)(*(short *)(*(int *)(param_2 + 0x18) + 0xc) - *(short *)(param_2 + 0x22)) * 2);
  }
  FUN_10063f30((undefined4 *)param_1,puVar1,iVar3 * 2);
  if ((*(int *)(param_2 + 8) != 100) && (0 < iVar3)) {
    do {
      iVar2 = ((int)*param_1 * *(int *)(param_2 + 8)) / 100;
      if ((0x7ffe < iVar2) || (iVar2 < -0x7fff)) {
        iVar2 = ((iVar2 < 0x7fff) - 1 & 0xfffd) - 0x7fff;
      }
      *param_1 = (short)iVar2;
      param_1 = param_1 + 1;
      iVar3 = iVar3 + -1;
    } while (iVar3 != 0);
  }
  return 1;
}



/* ===== FUN_1002bbd0 ===== */
/* Entry: 1002bbd0 */

void __cdecl FUN_1002bbd0(int param_1,int param_2,undefined4 param_3,int param_4)

{
  short sVar1;
  int iVar2;
  byte abStackY_8130 [32736];
  byte local_130 [300];

  FUN_10025440(*(int *)(*(int *)(param_4 + 4) + *(short *)(param_2 + 0x12) * 4),
               *(int *)(param_2 + 0x1c),0,(char *)local_130,1,(uint)*(byte *)(param_2 + 0x20));
  sVar1 = 0;
  if (*(char *)(param_2 + 0x20) != '\0') {
    do {
      iVar2 = (int)sVar1;
      sVar1 = sVar1 + 1;
      *(ushort *)(param_1 + iVar2 * 2) = (ushort)local_130[iVar2];
    } while (sVar1 < (short)(ushort)*(byte *)(param_2 + 0x20));
  }
  sVar1 = 0;
  if (*(char *)(param_2 + 0x20) != '\0') {
    do {
      *(short *)(param_1 + sVar1 * 2) = *(short *)(param_1 + sVar1 * 2) << 1;
      sVar1 = sVar1 + 1;
    } while (sVar1 < (short)(ushort)*(byte *)(param_2 + 0x20));
  }
  return;
}



/* ===== FUN_10001b30 ===== */
/* Entry: 10001b30 */

undefined4 __cdecl
FUN_10001b30(undefined4 param_1,undefined4 param_2,undefined4 param_3,undefined4 *param_4)

{
  uint uVar1;
  undefined4 uVar2;
  int iVar3;
  undefined1 local_1064 [1032];
  undefined4 local_c5c;
  undefined4 local_c58;
  undefined4 local_c54;
  undefined4 local_c50;
  undefined4 local_c4c;
  undefined4 local_c48;
  undefined *local_c44;
  undefined1 local_838 [1024];
  undefined1 *local_438;
  undefined4 local_434;
  undefined4 local_42c;
  undefined4 local_428;
  undefined4 local_424;
  undefined4 local_420;
  undefined4 local_41c;
  undefined4 uStackY_1c;

  FUN_10064270();
  local_c54 = param_2;
  local_c58 = param_1;
  local_c50 = param_3;
  local_438 = local_838;
  local_424 = 0;
  local_420 = 0;
  local_c5c = 1;
  local_c4c = 0;
  local_c48 = 0;
  local_c44 = &DAT_100771ac;
  local_42c = 0;
  local_434 = 0;
  local_428 = 0;
  local_41c = 0;
  uStackY_1c = 0x10001bd6;
  uVar1 = FUN_10001ea0((int)local_1064,2);
  while (uVar1 != 4) {
    switch(uVar1) {
    case 0:
    case 1:
    case 2:
    case 3:
    case 8:
      if (uVar1 != 8) {
        uStackY_1c = 0x10001c0f;
        FUN_10001ea0((int)local_1064,3);
      }
      iVar3 = 4;
      do {
        iVar3 = iVar3 + -1;
      } while (iVar3 != 0);
                    /* WARNING: Could not recover jumptable at 0x10001c52. Too many branches */
                    /* WARNING: Treating indirect jump as call */
      uVar2 = (*(code *)(&PTR_LAB_10001e70)[uVar1])();
      return uVar2;
    default:
      return 0;
    case 5:
      FUN_10002170((int)local_1064);
      break;
    case 6:
      uStackY_1c = 0x10001e08;
      FUN_10001ea0((int)local_1064,2);
    }
    uStackY_1c = 0x10001e1c;
    uVar1 = FUN_10001ea0((int)local_1064,2);
  }
  *param_4 = local_c48;
  return 1;
}



/* ===== FUN_1002d010 ===== */
/* Entry: 1002d010 */

/* WARNING: Globals starting with '_' overlap smaller symbols at the same address */

void __cdecl FUN_1002d010(int param_1,int param_2,short param_3,short param_4)

{
  uint uVar1;
  uint uVar2;
  int iVar3;
  undefined4 *puVar4;
  longlong lVar5;
  int local_10;

  if (param_3 < param_4) {
    if (0 < param_3) {
      local_10 = (int)param_3;
      iVar3 = param_1 - param_2;
      do {
        lVar5 = __ftol();
        *(short *)(iVar3 + param_2) = (short)lVar5;
        param_2 = param_2 + 2;
        local_10 = local_10 + -1;
      } while (local_10 != 0);
    }
    uVar1 = (int)param_4 - (int)param_3;
    puVar4 = (undefined4 *)(param_1 + param_3 * 2);
    for (uVar2 = (uVar1 & 0x7fffffff) >> 1; uVar2 != 0; uVar2 = uVar2 - 1) {
      *puVar4 = 0;
      puVar4 = puVar4 + 1;
    }
    for (uVar1 = uVar1 * 2 & 3; uVar1 != 0; uVar1 = uVar1 - 1) {
      *(undefined1 *)puVar4 = 0;
      puVar4 = (undefined4 *)((int)puVar4 + 1);
    }
    return;
  }
  if (0 < param_4) {
    _param_4 = (int)param_4;
    iVar3 = param_1 - param_2;
    do {
      lVar5 = __ftol();
      *(short *)(param_2 + iVar3) = (short)lVar5;
      param_2 = param_2 + 2;
      _param_4 = _param_4 + -1;
    } while (_param_4 != 0);
  }
  return;
}



/* ===== FUN_1002cdf0 ===== */
/* Entry: 1002cdf0 */

/* WARNING: Globals starting with '_' overlap smaller symbols at the same address */

void __cdecl FUN_1002cdf0(undefined4 *param_1,int param_2,short param_3,short param_4)

{
  int iVar1;
  uint uVar2;
  uint uVar3;
  int iVar4;
  int iVar5;
  undefined4 *puVar6;
  longlong lVar7;
  int local_10;
  int local_8;

  if (param_3 < param_4) {
    iVar4 = (int)param_3;
    uVar2 = param_4 - iVar4;
    puVar6 = param_1;
    for (uVar3 = (uVar2 & 0x7fffffff) >> 1; uVar3 != 0; uVar3 = uVar3 - 1) {
      *puVar6 = 0;
      puVar6 = puVar6 + 1;
    }
    for (uVar2 = uVar2 * 2 & 3; uVar2 != 0; uVar2 = uVar2 - 1) {
      *(undefined1 *)puVar6 = 0;
      puVar6 = (undefined4 *)((int)puVar6 + 1);
    }
    if (0 < param_3) {
      iVar5 = 0;
      local_8 = iVar4;
      do {
        iVar1 = iVar5 + iVar4 * -2;
        lVar7 = __ftol();
        iVar5 = iVar5 + 2;
        *(short *)(iVar1 + param_4 * 2 + (int)param_1) = (short)lVar7;
        local_8 = local_8 + -1;
      } while (local_8 != 0);
      return;
    }
  }
  else if (0 < param_4) {
    local_10 = (int)param_4;
    iVar4 = 0;
    do {
      lVar7 = __ftol();
      *(short *)(iVar4 + (int)param_1) = (short)lVar7;
      iVar4 = iVar4 + 2;
      local_10 = local_10 + -1;
    } while (local_10 != 0);
  }
  return;
}



/* ===== FUN_1001b200 ===== */
/* Entry: 1001b200 */

void __cdecl FUN_1001b200(undefined4 *param_1,int param_2,undefined4 *param_3,int param_4)

{
  if (*(char *)(param_4 + 0x82) == '\0') {
    FUN_1001af50(param_1,param_2,param_3,param_4);
    return;
  }
  FUN_1001b0d0(param_1,param_2,(int)param_3,param_4);
  return;
}



/* ===== FUN_1002d1e0 ===== */
/* Entry: 1002d1e0 */

void __cdecl
FUN_1002d1e0(int param_1,short *param_2,int param_3,int param_4,undefined4 *param_5,int *param_6)

{
  int local_1c [6];

  FUN_1001b200(local_1c,param_4,param_5,(int)param_6);
  FUN_1002c120(param_1,local_1c);
  FUN_1002c8b0(param_2,param_1,(int)param_5,param_6);
  FUN_1002bbd0(param_3,param_1,param_5,(int)param_6);
  return;
}
