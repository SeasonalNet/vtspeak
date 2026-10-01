===== 0x10026630 =====
Function: FUN_10026630 @ 10026630

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



===== 0x100260ec =====
Function: FUN_10025fc0 @ 10025fc0

int __cdecl FUN_10025fc0(undefined2 *param_1,int param_2,int param_3)

{
  int iVar1;
  undefined4 *puVar2;
  int *piVar3;
  undefined4 uVar4;
  int iVar5;
  int iVar6;
  
  VT_SetDecimal0Pron_ENG();
  EnterCriticalSection((LPCRITICAL_SECTION)&DAT_100a8400);
  iVar6 = 0;
  if (0 < *(int *)((&DAT_100a0464)[param_3] + 0x4d14)) {
    piVar3 = &DAT_100a047c + param_3 * 0x400;
    do {
      if (*piVar3 == 0) {
        (&DAT_100a047c)[param_3 * 0x400 + iVar6] = -0xff;
        iVar1 = FUN_1001d9c0(0x1312e0);
        (&DAT_100a047c)[param_3 * 0x400 + iVar6] = iVar1;
        if ((iVar1 == 0) || (iVar1 == -0xff)) {
          *param_1 = 0;
          LeaveCriticalSection((LPCRITICAL_SECTION)&DAT_100a8400);
          return 0;
        }
        *(int *)(iVar1 + 0x11b248) = iVar6;
        *(undefined1 *)(iVar1 + 0x47784) = (undefined1)param_3;
        if (iVar1 != 0) {
          LeaveCriticalSection((LPCRITICAL_SECTION)&DAT_100a8400);
          VT_SetDecimal0Pron_ENG();
          if ((param_2 < 0x400) && (-1 < param_2)) {
            if ((&DAT_100a647c)[param_2] != 0) {
              *(int *)(iVar1 + 0x1312c0) = (&DAT_100a647c)[param_2];
              goto LAB_100260ec;
            }
            if (DAT_100a647c != 0) {
              *(int *)(iVar1 + 0x1312c0) = DAT_100a647c;
              goto LAB_100260ec;
            }
          }
          else if (DAT_100a647c != 0) {
            *(int *)(iVar1 + 0x1312c0) = DAT_100a647c;
            goto LAB_100260ec;
          }
          *(undefined4 *)(iVar1 + 0x1312c0) = 0;
LAB_100260ec:
          VT_SetDecimal0Pron_ENG();
          *(undefined4 *)(iVar1 + 8) = 0;
          *(undefined4 *)(iVar1 + 0x34) = 0;
          *(undefined4 *)(iVar1 + 0x30) = 0;
          *(undefined4 *)(iVar1 + 0x44) = 0;
          *(undefined2 *)(iVar1 + 0x1210d4) = 0;
          *(undefined1 *)(iVar1 + 0x122850) = 1;
          *(undefined4 *)(iVar1 + 0x1312c8) = 0xffffffff;
          *(undefined4 *)(iVar1 + 0x47788) = 0;
          *(undefined4 *)(iVar1 + 0x4778c) = 0;
          *(undefined4 *)(iVar1 + 0x1312c4) = 0xffffffff;
          *(undefined1 *)(iVar1 + 0x122851) = 0;
          *(undefined4 *)(iVar1 + 0x1312cc) = 0xffffffff;
          *(undefined4 *)(iVar1 + 0x1312d0) = 0;
          *(undefined4 *)(iVar1 + 0x122444) = 0xffffffff;
          *(undefined1 *)(iVar1 + 0x20) = 4;
          *(undefined1 *)(iVar1 + 0x21) = 0;
          *(undefined4 *)(iVar1 + 0x2c) = 0;
          *(undefined4 *)(iVar1 + 0x47794) = 0;
          *(undefined4 *)(iVar1 + 0x47798) = 0;
          *(undefined4 *)(iVar1 + 0x47790) = 0;
          *(undefined4 *)(iVar1 + 0x477a4) = 0xffffffff;
          *(int *)(iVar1 + 0x4c) = iVar1 + 0x50;
          *(undefined4 *)(iVar1 + 0x122440) = 0xffffffff;
          *(undefined1 *)(iVar1 + 0x1312d4) = 0;
          *(undefined1 *)(iVar1 + 0x1312d6) = 0;
          *(undefined1 *)(iVar1 + 0x1312d7) = 0;
          *(undefined1 *)(iVar1 + 0x1312d8) = 0;
          *(undefined1 *)(iVar1 + 0x1312d9) = 0;
          *(undefined1 *)(iVar1 + 0x1312da) = 0;
          *(undefined1 *)(iVar1 + 0x1312dc) = 0;
          *(undefined1 *)(iVar1 + 0x1312dd) = 0;
          *(undefined1 *)(iVar1 + 0x1312db) = 0;
          *(undefined1 *)(iVar1 + 0x1312d5) = 0;
          iVar6 = iVar1 + 0x1223c0;
          param_3 = 8;
          do {
            FUN_1001cbc0(iVar6);
            iVar6 = iVar6 + 0x10;
            param_3 = param_3 + -1;
          } while (param_3 != 0);
          FUN_1001cc70(iVar1);
          piVar3 = (int *)(iVar1 + 0x477ac);
          iVar6 = iVar1 + 0x513ec;
          iVar5 = 10000;
          do {
            *piVar3 = iVar6;
            iVar6 = iVar6 + 0x14;
            piVar3 = piVar3 + 1;
            iVar5 = iVar5 + -1;
          } while (iVar5 != 0);
          VT_SetDecimal0Pron_ENG();
          puVar2 = FUN_10040df0();
          *(undefined4 **)(iVar1 + 0x1312bc) = puVar2;
          if ((&DAT_100a7488)[*(char *)(iVar1 + 0x47784)] == '\0') {
            *(undefined4 *)(*(int *)(iVar1 + 0x1312bc) + 0x39e8) = 0;
          }
          else {
            *(undefined4 *)(*(int *)(iVar1 + 0x1312bc) + 0x39e8) = *(undefined4 *)(iVar1 + 0x1312c0)
            ;
          }
          FUN_1003e470(*(int *)(iVar1 + 0x1312bc),*(undefined4 *)(DAT_100a0460 + 0x2041c));
          FUN_1003e480(*(int *)(iVar1 + 0x1312bc));
          VT_SetDecimal0Pron_ENG();
          *(undefined4 *)(iVar1 + 0x47778) = 0;
          *(undefined4 *)(iVar1 + 0x4777c) = 0;
          *(undefined4 *)(iVar1 + 0x47780) = 0;
          piVar3 = VT_AllocSyncInfo_New_ENG();
          *(int **)(iVar1 + 0x47774) = piVar3;
          if (piVar3 != (int *)0x0) {
            VT_InitSyncInfo_New_ENG(*(int **)(iVar1 + 0x47774));
            VT_SetDecimal0Pron_ENG();
            param_3 = 0;
            piVar3 = (int *)(iVar1 + 0x38);
            while( true ) {
              iVar6 = FUN_1001d9c0(0xc);
              *piVar3 = iVar6;
              if (iVar6 == 0) {
                *param_1 = 0;
                return 0;
              }
              uVar4 = FUN_1001d9c0(60000);
              *(undefined4 *)(*piVar3 + 4) = uVar4;
              if (*(int *)(*piVar3 + 4) == 0) break;
              piVar3 = piVar3 + 1;
              param_3 = param_3 + 1;
              if (1 < param_3) {
                *(undefined4 *)(*(int *)(iVar1 + 0x38) + 8) = *(undefined4 *)(iVar1 + 0x3c);
                *(undefined4 *)(*(int *)(iVar1 + 0x3c) + 8) = *(undefined4 *)(iVar1 + 0x38);
                **(undefined1 **)(iVar1 + 0x38) = 0;
                **(undefined1 **)(iVar1 + 0x3c) = 1;
                VT_SetDecimal0Pron_ENG();
                *param_1 = 1;
                VT_SetDecimal0Pron_ENG();
                return iVar1;
              }
            }
            *param_1 = 0;
            return 0;
          }
          *param_1 = 0;
          return 0;
        }
        break;
      }
      iVar6 = iVar6 + 1;
      piVar3 = piVar3 + 1;
    } while (iVar6 < *(int *)((&DAT_100a0464)[param_3] + 0x4d14));
  }
  *param_1 = 0;
  LeaveCriticalSection((LPCRITICAL_SECTION)&DAT_100a8400);
  return 0;
}



===== 0x10022850 =====
Function: FUN_10022850 @ 10022850

void __cdecl FUN_10022850(int param_1,int param_2)

{
  undefined4 uVar1;
  
  if (*(char *)(param_2 + 0x1210d7) == '\x01') {
    uVar1 = *(undefined4 *)(param_2 + 0x1210dc);
  }
  else {
    uVar1 = *(undefined4 *)(param_1 + 0x4cf4);
  }
  *(undefined4 *)(param_2 + 0x1210e8) = uVar1;
  if (*(char *)(param_2 + 0x1210d8) == '\x01') {
    uVar1 = *(undefined4 *)(param_2 + 0x1210e0);
  }
  else {
    uVar1 = *(undefined4 *)(param_1 + 0x4cf0);
  }
  *(undefined4 *)(param_2 + 0x1210ec) = uVar1;
  if (*(char *)(param_2 + 0x1210d9) == '\x01') {
    uVar1 = *(undefined4 *)(param_2 + 0x1210e4);
  }
  else {
    uVar1 = *(undefined4 *)(param_1 + 0x4cf8);
  }
  *(undefined4 *)(param_2 + 0x1210f0) = uVar1;
  if (200 < *(int *)(param_2 + 0x1210e8)) {
    *(undefined4 *)(param_2 + 0x1210e8) = 200;
  }
  if (*(int *)(param_2 + 0x1210e8) < 0x32) {
    *(undefined4 *)(param_2 + 0x1210e8) = 0x32;
  }
  if (400 < *(int *)(param_2 + 0x1210ec)) {
    *(undefined4 *)(param_2 + 0x1210ec) = 400;
  }
  if (*(int *)(param_2 + 0x1210ec) < 0x32) {
    *(undefined4 *)(param_2 + 0x1210ec) = 0x32;
  }
  if (500 < *(int *)(param_2 + 0x1210f0)) {
    *(undefined4 *)(param_2 + 0x1210f0) = 500;
  }
  if (*(int *)(param_2 + 0x1210f0) < 0) {
    *(undefined4 *)(param_2 + 0x1210f0) = 0;
  }
  *(undefined4 *)(param_2 + 0x1312c4) = 0xffffffff;
  *(undefined4 *)(param_2 + 0x1210fc) = *(undefined4 *)(param_2 + 0x1210e8);
  *(undefined4 *)(param_2 + 0x1210f4) = *(undefined4 *)(param_2 + 0x1210ec);
  *(undefined4 *)(param_2 + 0x1210f8) = *(undefined4 *)(param_2 + 0x1210f0);
  return;
}



===== 0x10022dc0 =====
Function: FUN_10022dc0 @ 10022dc0

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



===== 0x100130e0 =====
Function: FUN_100130e0 @ 100130e0

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



