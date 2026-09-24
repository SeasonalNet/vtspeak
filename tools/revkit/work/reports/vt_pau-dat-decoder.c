/* Ghidra decompiler output: approximate C pseudocode, not original source. */

/* ===== 0x10001b30 ===== */
/* Entry: 10001b30 FUN_10001b30 */

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
      uVar2 = (*(code *)(&PTR_decode_10001c75_10001e70)[uVar1])();
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



/* ===== 0x10001ea0 ===== */
/* Entry: 10001ea0 FUN_10001ea0 */

uint __cdecl FUN_10001ea0(int param_1,int param_2)

{
  uint uVar1;
  undefined4 uVar2;
  int iVar3;
  uint uVar4;

  if (*(int *)(param_1 + 0xc3c) == 0) {
    uVar2 = FUN_10001fb0(param_1);
    *(undefined4 *)(param_1 + 0xc30) = uVar2;
    *(undefined4 *)(param_1 + 0xc3c) = 0x20;
  }
  uVar4 = 0;
  iVar3 = *(int *)(param_1 + 0xc3c) + -1;
  *(int *)(param_1 + 0xc3c) = iVar3;
  uVar1 = *(uint *)(param_1 + 0xc30) & 1 << ((byte)iVar3 & 0x1f);
  while (uVar1 == 0) {
    if (iVar3 == 0) {
      uVar2 = FUN_10001fb0(param_1);
      *(undefined4 *)(param_1 + 0xc30) = uVar2;
      *(undefined4 *)(param_1 + 0xc3c) = 0x20;
    }
    uVar4 = uVar4 + 1;
    iVar3 = *(int *)(param_1 + 0xc3c) + -1;
    *(int *)(param_1 + 0xc3c) = iVar3;
    uVar1 = *(uint *)(param_1 + 0xc30) & 1 << ((byte)iVar3 & 0x1f);
  }
  while( true ) {
    if (param_2 == 0) {
      return uVar4;
    }
    iVar3 = *(int *)(param_1 + 0xc3c);
    if (param_2 <= iVar3) break;
    uVar4 = uVar4 << ((byte)iVar3 & 0x1f) |
            *(uint *)(*(int *)(param_1 + 0x420) + iVar3 * 4) & *(uint *)(param_1 + 0xc30);
    uVar2 = FUN_10001fb0(param_1);
    *(undefined4 *)(param_1 + 0xc30) = uVar2;
    param_2 = param_2 - *(int *)(param_1 + 0xc3c);
    *(undefined4 *)(param_1 + 0xc3c) = 0x20;
  }
  iVar3 = *(int *)(param_1 + 0xc3c) - param_2;
  uVar1 = *(uint *)(*(int *)(param_1 + 0x420) + param_2 * 4);
  *(int *)(param_1 + 0xc3c) = iVar3;
  return uVar1 & *(uint *)(param_1 + 0xc30) >> ((byte)iVar3 & 0x1f) |
         uVar4 << ((byte)param_2 & 0x1f);
}



/* ===== 0x10002170 ===== */
/* Entry: 10002170 FUN_10002170 */

void __cdecl FUN_10002170(int param_1)

{
  uint uVar1;

  uVar1 = FUN_10001ea0(param_1,2);
  FUN_10001ea0(param_1,uVar1);
  return;
}



/* ===== 0x100020a0 ===== */
/* Entry: 100020a0 FUN_100020a0 */

uint __cdecl FUN_100020a0(int param_1,int param_2)

{
  uint uVar1;

  uVar1 = FUN_10001ea0(param_1,param_2 + 1);
  if ((uVar1 & 1) != 0) {
    return ~(uVar1 >> 1);
  }
  return uVar1 >> 1;
}



/* ===== 0x100020d0 ===== */
/* Entry: 100020d0 FUN_100020d0 */

void __cdecl FUN_100020d0(int *param_1,int param_2,int param_3)

{
  if ((param_3 != 0) && (0 < param_2)) {
    do {
      param_2 = param_2 + -1;
      *param_1 = *param_1 << ((byte)param_3 & 0x1f);
      param_1 = param_1 + 1;
    } while (param_2 != 0);
  }
  return;
}



/* ===== 0x10002100 ===== */
/* Entry: 10002100 FUN_10002100 */

void __cdecl FUN_10002100(int param_1,int *param_2,int param_3)

{
  int iVar1;
  undefined2 *puVar2;
  int iVar3;

  puVar2 = (undefined2 *)(param_1 + 0x204);
  iVar3 = param_3;
  if (0 < param_3) {
    do {
      iVar1 = *param_2;
      if (0x7fff < iVar1) {
        iVar1 = 0x7fff;
      }
      *puVar2 = (short)iVar1;
      puVar2 = puVar2 + 1;
      param_2 = param_2 + 1;
      iVar3 = iVar3 + -1;
    } while (iVar3 != 0);
  }
  FUN_10063f30((undefined4 *)(*(int *)(param_1 + 0x410) + *(int *)(param_1 + 0x41c)),
               (undefined4 *)(param_1 + 0x204),param_3 * 2);
  *(int *)(param_1 + 0x41c) = *(int *)(param_1 + 0x41c) + param_3 * 2;
  return;
}



/* ===== 0x10001fb0 ===== */
/* Entry: 10001fb0 FUN_10001fb0 */

undefined4 __cdecl FUN_10001fb0(int param_1)

{
  undefined4 *puVar1;
  undefined1 uVar2;
  undefined1 uVar3;
  undefined1 uVar4;
  undefined1 uVar5;
  undefined1 *puVar6;
  int iVar7;

  if (*(int *)(param_1 + 0xc38) < 4) {
    iVar7 = *(int *)(param_1 + 0x418);
    puVar1 = (undefined4 *)(param_1 + 0x82c);
    if (*(int *)(param_1 + 0x414) < iVar7 + 0x400) {
      FUN_10063f30(puVar1,(undefined4 *)(*(int *)(param_1 + 0x40c) + iVar7),
                   *(int *)(param_1 + 0x414) - iVar7);
      iVar7 = *(int *)(param_1 + 0x414);
      *(int *)(param_1 + 0xc38) = *(int *)(param_1 + 0xc38) + (iVar7 - *(int *)(param_1 + 0x418));
    }
    else {
      FUN_10063f30(puVar1,(undefined4 *)(*(int *)(param_1 + 0x40c) + iVar7),0x400);
      *(int *)(param_1 + 0xc38) = *(int *)(param_1 + 0xc38) + 0x400;
      iVar7 = *(int *)(param_1 + 0x418) + 0x400;
    }
    *(undefined4 **)(param_1 + 0xc2c) = puVar1;
    *(int *)(param_1 + 0x418) = iVar7;
  }
  puVar6 = *(undefined1 **)(param_1 + 0xc2c);
  uVar2 = *puVar6;
  uVar3 = puVar6[2];
  uVar4 = puVar6[1];
  uVar5 = puVar6[3];
  *(undefined1 **)(param_1 + 0xc2c) = puVar6 + 4;
  *(int *)(param_1 + 0xc38) = *(int *)(param_1 + 0xc38) + -4;
  return CONCAT31(CONCAT21(CONCAT11(uVar2,uVar4),uVar3),uVar5);
}



/* ===== 0x10001c59 ===== */
/* Entry: 10001c59 decode_10001c59 */

/* WARNING (jumptable): Unable to track spacebase fully for stack */

undefined4 decode_10001c59(void)

{
  undefined4 uVar1;
  int iVar2;
  undefined4 *puVar3;
  uint uVar4;
  int *piVar5;
  undefined4 *puVar6;
  int iVar7;
  int iVar8;
  int unaff_EBP;
  int unaff_EDI;

  if (0 < unaff_EDI) {
    puVar6 = (undefined4 *)(unaff_EBP + -0x414);
    for (; unaff_EDI != 0; unaff_EDI = unaff_EDI + -1) {
      *puVar6 = 0;
      puVar6 = puVar6 + 1;
    }
    unaff_EDI = *(int *)(unaff_EBP + 8);
  }
  iVar2 = unaff_EDI / 2;
  if (0 < unaff_EDI) {
    piVar5 = (int *)(unaff_EBP + -0x414);
    iVar7 = unaff_EDI;
    do {
      iVar8 = *piVar5;
      piVar5 = piVar5 + 1;
      iVar2 = iVar2 + iVar8;
      iVar7 = iVar7 + -1;
    } while (iVar7 != 0);
  }
  puVar6 = (undefined4 *)(unaff_EBP + -0x10);
  iVar7 = 3;
  do {
    puVar6[-1] = *puVar6;
    puVar6 = puVar6 + 1;
    iVar7 = iVar7 + -1;
  } while (iVar7 != 0);
  iVar7 = *(int *)(unaff_EBP + -4);
  iVar8 = 3;
  puVar6 = (undefined4 *)(unaff_EBP + -0x420);
  *(int *)(unaff_EBP + -8) = iVar2 / unaff_EDI << ((byte)iVar7 & 0x1f);
  puVar3 = (undefined4 *)(unaff_EBP + -0x420 + unaff_EDI * 4);
  do {
    uVar1 = *puVar3;
    puVar3 = puVar3 + 1;
    *puVar6 = uVar1;
    puVar6 = puVar6 + 1;
    iVar8 = iVar8 + -1;
  } while (iVar8 != 0);
  FUN_100020d0((int *)(unaff_EBP + -0x414),unaff_EDI,iVar7);
  FUN_10002100(unaff_EBP + -0x1060,(int *)(unaff_EBP + -0x414),unaff_EDI);
  while( true ) {
    uVar4 = FUN_10001ea0(unaff_EBP + -0x1060,2);
    *(uint *)(unaff_EBP + 0x10) = uVar4;
    if (uVar4 == 4) break;
    switch(*(int *)(unaff_EBP + 0x10)) {
    case 0:
    case 1:
    case 2:
    case 3:
    case 8:
      if (*(int *)(unaff_EBP + 0x10) != 8) {
        FUN_10001ea0(unaff_EBP + -0x1060,3);
      }
      iVar2 = 4;
      do {
        iVar2 = iVar2 + -1;
      } while (iVar2 != 0);
                    /* WARNING: Could not recover jumptable at 0x10001c52. Too many branches */
                    /* WARNING: Treating indirect jump as call */
      uVar1 = (*(code *)(&PTR_decode_10001c75_10001e70)[*(int *)(unaff_EBP + 0x10)])();
      return uVar1;
    default:
      return 0;
    case 5:
      uVar1 = FUN_10002170(unaff_EBP + -0x1060);
      *(undefined4 *)(unaff_EBP + 8) = uVar1;
      break;
    case 6:
      uVar4 = FUN_10001ea0(unaff_EBP + -0x1060,2);
      *(uint *)(unaff_EBP + -4) = uVar4;
    }
  }
  **(undefined4 **)(unaff_EBP + 0x14) = *(undefined4 *)(unaff_EBP + -0xc44);
  return 1;
}



/* ===== 0x10001c75 ===== */
/* Entry: 10001c75 decode_10001c75 */

/* WARNING (jumptable): Unable to track spacebase fully for stack */

undefined4 decode_10001c75(void)

{
  undefined4 uVar1;
  uint uVar2;
  int iVar3;
  undefined4 *puVar4;
  int *piVar5;
  undefined4 *puVar6;
  int iVar7;
  int iVar8;
  int unaff_EBX;
  int unaff_EBP;
  int unaff_ESI;
  int unaff_EDI;

  if (0 < unaff_EDI) {
    piVar5 = (int *)(unaff_EBP + -0x414);
    *(undefined4 *)(unaff_EBP + 0x10) = *(undefined4 *)(unaff_EBP + 8);
    do {
      uVar2 = FUN_100020a0(unaff_EBP + -0x1060,unaff_EBX);
      *piVar5 = uVar2 + unaff_ESI;
      piVar5 = piVar5 + 1;
      iVar3 = *(int *)(unaff_EBP + 0x10) + -1;
      *(int *)(unaff_EBP + 0x10) = iVar3;
    } while (iVar3 != 0);
    unaff_EDI = *(int *)(unaff_EBP + 8);
  }
  iVar3 = unaff_EDI / 2;
  if (0 < unaff_EDI) {
    piVar5 = (int *)(unaff_EBP + -0x414);
    iVar7 = unaff_EDI;
    do {
      iVar8 = *piVar5;
      piVar5 = piVar5 + 1;
      iVar3 = iVar3 + iVar8;
      iVar7 = iVar7 + -1;
    } while (iVar7 != 0);
  }
  puVar6 = (undefined4 *)(unaff_EBP + -0x10);
  iVar7 = 3;
  do {
    puVar6[-1] = *puVar6;
    puVar6 = puVar6 + 1;
    iVar7 = iVar7 + -1;
  } while (iVar7 != 0);
  iVar7 = *(int *)(unaff_EBP + -4);
  iVar8 = 3;
  puVar6 = (undefined4 *)(unaff_EBP + -0x420);
  *(int *)(unaff_EBP + -8) = iVar3 / unaff_EDI << ((byte)iVar7 & 0x1f);
  puVar4 = (undefined4 *)(unaff_EBP + -0x420 + unaff_EDI * 4);
  do {
    uVar1 = *puVar4;
    puVar4 = puVar4 + 1;
    *puVar6 = uVar1;
    puVar6 = puVar6 + 1;
    iVar8 = iVar8 + -1;
  } while (iVar8 != 0);
  FUN_100020d0((int *)(unaff_EBP + -0x414),unaff_EDI,iVar7);
  FUN_10002100(unaff_EBP + -0x1060,(int *)(unaff_EBP + -0x414),unaff_EDI);
  while( true ) {
    uVar2 = FUN_10001ea0(unaff_EBP + -0x1060,2);
    *(uint *)(unaff_EBP + 0x10) = uVar2;
    if (uVar2 == 4) break;
    switch(*(int *)(unaff_EBP + 0x10)) {
    case 0:
    case 1:
    case 2:
    case 3:
    case 8:
      if (*(int *)(unaff_EBP + 0x10) != 8) {
        FUN_10001ea0(unaff_EBP + -0x1060,3);
      }
      iVar3 = 4;
      do {
        iVar3 = iVar3 + -1;
      } while (iVar3 != 0);
                    /* WARNING: Could not recover jumptable at 0x10001c52. Too many branches */
                    /* WARNING: Treating indirect jump as call */
      uVar1 = (*(code *)(&PTR_decode_10001c75_10001e70)[*(int *)(unaff_EBP + 0x10)])();
      return uVar1;
    default:
      return 0;
    case 5:
      uVar1 = FUN_10002170(unaff_EBP + -0x1060);
      *(undefined4 *)(unaff_EBP + 8) = uVar1;
      break;
    case 6:
      uVar2 = FUN_10001ea0(unaff_EBP + -0x1060,2);
      *(uint *)(unaff_EBP + -4) = uVar2;
    }
  }
  **(undefined4 **)(unaff_EBP + 0x14) = *(undefined4 *)(unaff_EBP + -0xc44);
  return 1;
}



/* ===== 0x10001cb1 ===== */
/* Entry: 10001cb1 decode_10001cb1 */

/* WARNING (jumptable): Unable to track spacebase fully for stack */

undefined4 decode_10001cb1(void)

{
  undefined4 uVar1;
  uint uVar2;
  int iVar3;
  undefined4 *puVar4;
  int *piVar5;
  undefined4 *puVar6;
  int iVar7;
  int iVar8;
  int unaff_EBX;
  int unaff_EBP;
  int unaff_EDI;

  if (0 < unaff_EDI) {
    piVar5 = (int *)(unaff_EBP + -0x414);
    *(int *)(unaff_EBP + 0x10) = unaff_EDI;
    do {
      uVar2 = FUN_100020a0(unaff_EBP + -0x1060,unaff_EBX);
      *piVar5 = uVar2 + piVar5[-1];
      piVar5 = piVar5 + 1;
      iVar3 = *(int *)(unaff_EBP + 0x10) + -1;
      *(int *)(unaff_EBP + 0x10) = iVar3;
    } while (iVar3 != 0);
  }
  iVar3 = unaff_EDI / 2;
  if (0 < unaff_EDI) {
    piVar5 = (int *)(unaff_EBP + -0x414);
    iVar7 = unaff_EDI;
    do {
      iVar8 = *piVar5;
      piVar5 = piVar5 + 1;
      iVar3 = iVar3 + iVar8;
      iVar7 = iVar7 + -1;
    } while (iVar7 != 0);
  }
  puVar6 = (undefined4 *)(unaff_EBP + -0x10);
  iVar7 = 3;
  do {
    puVar6[-1] = *puVar6;
    puVar6 = puVar6 + 1;
    iVar7 = iVar7 + -1;
  } while (iVar7 != 0);
  iVar7 = *(int *)(unaff_EBP + -4);
  iVar8 = 3;
  puVar6 = (undefined4 *)(unaff_EBP + -0x420);
  *(int *)(unaff_EBP + -8) = iVar3 / unaff_EDI << ((byte)iVar7 & 0x1f);
  puVar4 = (undefined4 *)(unaff_EBP + -0x420 + unaff_EDI * 4);
  do {
    uVar1 = *puVar4;
    puVar4 = puVar4 + 1;
    *puVar6 = uVar1;
    puVar6 = puVar6 + 1;
    iVar8 = iVar8 + -1;
  } while (iVar8 != 0);
  FUN_100020d0((int *)(unaff_EBP + -0x414),unaff_EDI,iVar7);
  FUN_10002100(unaff_EBP + -0x1060,(int *)(unaff_EBP + -0x414),unaff_EDI);
  while( true ) {
    uVar2 = FUN_10001ea0(unaff_EBP + -0x1060,2);
    *(uint *)(unaff_EBP + 0x10) = uVar2;
    if (uVar2 == 4) break;
    switch(*(int *)(unaff_EBP + 0x10)) {
    case 0:
    case 1:
    case 2:
    case 3:
    case 8:
      if (*(int *)(unaff_EBP + 0x10) != 8) {
        FUN_10001ea0(unaff_EBP + -0x1060,3);
      }
      iVar3 = 4;
      do {
        iVar3 = iVar3 + -1;
      } while (iVar3 != 0);
                    /* WARNING: Could not recover jumptable at 0x10001c52. Too many branches */
                    /* WARNING: Treating indirect jump as call */
      uVar1 = (*(code *)(&PTR_decode_10001c75_10001e70)[*(int *)(unaff_EBP + 0x10)])();
      return uVar1;
    default:
      return 0;
    case 5:
      uVar1 = FUN_10002170(unaff_EBP + -0x1060);
      *(undefined4 *)(unaff_EBP + 8) = uVar1;
      break;
    case 6:
      uVar2 = FUN_10001ea0(unaff_EBP + -0x1060,2);
      *(uint *)(unaff_EBP + -4) = uVar2;
    }
  }
  **(undefined4 **)(unaff_EBP + 0x14) = *(undefined4 *)(unaff_EBP + -0xc44);
  return 1;
}



/* ===== 0x10001ce7 ===== */
/* Entry: 10001ce7 decode_10001ce7 */

/* WARNING (jumptable): Unable to track spacebase fully for stack */

undefined4 decode_10001ce7(void)

{
  undefined4 uVar1;
  uint uVar2;
  int iVar3;
  undefined4 *puVar4;
  int *piVar5;
  undefined4 *puVar6;
  int iVar7;
  int iVar8;
  int unaff_EBX;
  int unaff_EBP;
  int unaff_EDI;

  if (0 < unaff_EDI) {
    piVar5 = (int *)(unaff_EBP + -0x41c);
    *(int *)(unaff_EBP + 0x10) = unaff_EDI;
    do {
      uVar2 = FUN_100020a0(unaff_EBP + -0x1060,unaff_EBX);
      piVar5[2] = uVar2 + (piVar5[1] * 2 - *piVar5);
      piVar5 = piVar5 + 1;
      iVar3 = *(int *)(unaff_EBP + 0x10) + -1;
      *(int *)(unaff_EBP + 0x10) = iVar3;
    } while (iVar3 != 0);
  }
  iVar3 = unaff_EDI / 2;
  if (0 < unaff_EDI) {
    piVar5 = (int *)(unaff_EBP + -0x414);
    iVar7 = unaff_EDI;
    do {
      iVar8 = *piVar5;
      piVar5 = piVar5 + 1;
      iVar3 = iVar3 + iVar8;
      iVar7 = iVar7 + -1;
    } while (iVar7 != 0);
  }
  puVar6 = (undefined4 *)(unaff_EBP + -0x10);
  iVar7 = 3;
  do {
    puVar6[-1] = *puVar6;
    puVar6 = puVar6 + 1;
    iVar7 = iVar7 + -1;
  } while (iVar7 != 0);
  iVar7 = *(int *)(unaff_EBP + -4);
  iVar8 = 3;
  puVar6 = (undefined4 *)(unaff_EBP + -0x420);
  *(int *)(unaff_EBP + -8) = iVar3 / unaff_EDI << ((byte)iVar7 & 0x1f);
  puVar4 = (undefined4 *)(unaff_EBP + -0x420 + unaff_EDI * 4);
  do {
    uVar1 = *puVar4;
    puVar4 = puVar4 + 1;
    *puVar6 = uVar1;
    puVar6 = puVar6 + 1;
    iVar8 = iVar8 + -1;
  } while (iVar8 != 0);
  FUN_100020d0((int *)(unaff_EBP + -0x414),unaff_EDI,iVar7);
  FUN_10002100(unaff_EBP + -0x1060,(int *)(unaff_EBP + -0x414),unaff_EDI);
  while( true ) {
    uVar2 = FUN_10001ea0(unaff_EBP + -0x1060,2);
    *(uint *)(unaff_EBP + 0x10) = uVar2;
    if (uVar2 == 4) break;
    switch(*(int *)(unaff_EBP + 0x10)) {
    case 0:
    case 1:
    case 2:
    case 3:
    case 8:
      if (*(int *)(unaff_EBP + 0x10) != 8) {
        FUN_10001ea0(unaff_EBP + -0x1060,3);
      }
      iVar3 = 4;
      do {
        iVar3 = iVar3 + -1;
      } while (iVar3 != 0);
                    /* WARNING: Could not recover jumptable at 0x10001c52. Too many branches */
                    /* WARNING: Treating indirect jump as call */
      uVar1 = (*(code *)(&PTR_decode_10001c75_10001e70)[*(int *)(unaff_EBP + 0x10)])();
      return uVar1;
    default:
      return 0;
    case 5:
      uVar1 = FUN_10002170(unaff_EBP + -0x1060);
      *(undefined4 *)(unaff_EBP + 8) = uVar1;
      break;
    case 6:
      uVar2 = FUN_10001ea0(unaff_EBP + -0x1060,2);
      *(uint *)(unaff_EBP + -4) = uVar2;
    }
  }
  **(undefined4 **)(unaff_EBP + 0x14) = *(undefined4 *)(unaff_EBP + -0xc44);
  return 1;
}



/* ===== 0x10001d20 ===== */
/* Entry: 10001d20 decode_10001d20 */

/* WARNING (jumptable): Unable to track spacebase fully for stack */

undefined4 decode_10001d20(void)

{
  undefined4 uVar1;
  uint uVar2;
  int iVar3;
  undefined4 *puVar4;
  int *piVar5;
  undefined4 *puVar6;
  int iVar7;
  int iVar8;
  int unaff_EBX;
  int unaff_EBP;
  int unaff_EDI;

  if (0 < unaff_EDI) {
    *(int *)(unaff_EBP + 0x10) = unaff_EDI;
    piVar5 = (int *)(unaff_EBP + -0x41c);
    do {
      uVar2 = FUN_100020a0(unaff_EBP + -0x1060,unaff_EBX);
      iVar3 = *(int *)(unaff_EBP + 0x10) + -1;
      piVar5[2] = uVar2 + (piVar5[1] - *piVar5) * 3 + piVar5[-1];
      *(int *)(unaff_EBP + 0x10) = iVar3;
      piVar5 = piVar5 + 1;
    } while (iVar3 != 0);
  }
  iVar3 = unaff_EDI / 2;
  if (0 < unaff_EDI) {
    piVar5 = (int *)(unaff_EBP + -0x414);
    iVar7 = unaff_EDI;
    do {
      iVar8 = *piVar5;
      piVar5 = piVar5 + 1;
      iVar3 = iVar3 + iVar8;
      iVar7 = iVar7 + -1;
    } while (iVar7 != 0);
  }
  puVar6 = (undefined4 *)(unaff_EBP + -0x10);
  iVar7 = 3;
  do {
    puVar6[-1] = *puVar6;
    puVar6 = puVar6 + 1;
    iVar7 = iVar7 + -1;
  } while (iVar7 != 0);
  iVar7 = *(int *)(unaff_EBP + -4);
  iVar8 = 3;
  puVar6 = (undefined4 *)(unaff_EBP + -0x420);
  *(int *)(unaff_EBP + -8) = iVar3 / unaff_EDI << ((byte)iVar7 & 0x1f);
  puVar4 = (undefined4 *)(unaff_EBP + -0x420 + unaff_EDI * 4);
  do {
    uVar1 = *puVar4;
    puVar4 = puVar4 + 1;
    *puVar6 = uVar1;
    puVar6 = puVar6 + 1;
    iVar8 = iVar8 + -1;
  } while (iVar8 != 0);
  FUN_100020d0((int *)(unaff_EBP + -0x414),unaff_EDI,iVar7);
  FUN_10002100(unaff_EBP + -0x1060,(int *)(unaff_EBP + -0x414),unaff_EDI);
  while( true ) {
    uVar2 = FUN_10001ea0(unaff_EBP + -0x1060,2);
    *(uint *)(unaff_EBP + 0x10) = uVar2;
    if (uVar2 == 4) break;
    switch(*(int *)(unaff_EBP + 0x10)) {
    case 0:
    case 1:
    case 2:
    case 3:
    case 8:
      if (*(int *)(unaff_EBP + 0x10) != 8) {
        FUN_10001ea0(unaff_EBP + -0x1060,3);
      }
      iVar3 = 4;
      do {
        iVar3 = iVar3 + -1;
      } while (iVar3 != 0);
                    /* WARNING: Could not recover jumptable at 0x10001c52. Too many branches */
                    /* WARNING: Treating indirect jump as call */
      uVar1 = (*(code *)(&PTR_decode_10001c75_10001e70)[*(int *)(unaff_EBP + 0x10)])();
      return uVar1;
    default:
      return 0;
    case 5:
      uVar1 = FUN_10002170(unaff_EBP + -0x1060);
      *(undefined4 *)(unaff_EBP + 8) = uVar1;
      break;
    case 6:
      uVar2 = FUN_10001ea0(unaff_EBP + -0x1060,2);
      *(uint *)(unaff_EBP + -4) = uVar2;
    }
  }
  **(undefined4 **)(unaff_EBP + 0x14) = *(undefined4 *)(unaff_EBP + -0xc44);
  return 1;
}



/* ===== 0x1002c120 ===== */
/* Entry: 1002c120 FUN_1002c120 */

void __cdecl FUN_1002c120(int param_1,int *param_2)

{
  char cVar1;
  byte bVar2;
  ushort uVar3;
  int iVar4;

  *(ushort *)(param_1 + 0x12) = (ushort)*(byte *)(param_2 + 5);
  cVar1 = *(char *)(param_1 + 0x26);
  *(int **)(param_1 + 0x18) = param_2;
  bVar2 = *(byte *)((int)param_2 + 9);
  uVar3 = (ushort)bVar2;
  if (cVar1 == '\x01') {
    *(ushort *)(param_1 + 0x22) = (ushort)*(byte *)(param_2 + 2);
    bVar2 = *(byte *)((int)param_2 + 9);
  }
  else {
    if (cVar1 != '\x02') {
      bVar2 = *(byte *)(param_2 + 2);
    }
    *(ushort *)(param_1 + 0x22) = (ushort)bVar2;
    bVar2 = *(byte *)((int)param_2 + 10);
  }
  *(ushort *)(param_1 + 0x24) = (ushort)bVar2;
  *(short *)(param_1 + 0x22) = *(short *)(param_1 + 0x22) << 1;
  *(short *)(param_1 + 0x24) = *(short *)(param_1 + 0x24) << 1;
  if (cVar1 == '\x01') {
    *(int *)(param_1 + 0x14) = *param_2;
    iVar4 = (int)(short)param_2[3];
  }
  else {
    if (cVar1 == '\x02') {
      *(int *)(param_1 + 0xc) = (int)*(short *)((int)param_2 + 0xe);
      *(int *)(param_1 + 0x14) = (int)(short)((short)param_2[3] + uVar3 * -2) + *param_2;
      goto LAB_1002c1cb;
    }
    *(int *)(param_1 + 0x14) = *param_2;
    iVar4 = ((int)*(short *)((int)param_2 + 0xe) + (int)(short)param_2[3]) - (int)(short)(uVar3 * 2)
    ;
  }
  *(int *)(param_1 + 0xc) = iVar4;
LAB_1002c1cb:
  if (cVar1 == '\x01') {
    *(int *)(param_1 + 0x1c) = param_2[1];
    *(undefined1 *)(param_1 + 0x20) = *(undefined1 *)((int)param_2 + 0x12);
    return;
  }
  if (cVar1 == '\x02') {
    *(uint *)(param_1 + 0x1c) = (*(byte *)((int)param_2 + 0x12) - 1) + param_2[1];
    *(undefined1 *)(param_1 + 0x20) = *(undefined1 *)((int)param_2 + 0x13);
    return;
  }
  *(int *)(param_1 + 0x1c) = param_2[1];
  *(char *)(param_1 + 0x20) = *(char *)((int)param_2 + 0x13) + *(char *)((int)param_2 + 0x12) + -1;
  return;
}
