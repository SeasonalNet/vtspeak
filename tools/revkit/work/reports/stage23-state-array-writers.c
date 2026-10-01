===== 0x10022970 =====
Function: FUN_10022970 @ 10022970

undefined4 __cdecl FUN_10022970(int *param_1)

{
  int iVar1;
  int iVar2;
  int iVar3;
  int iVar4;
  int iVar5;
  int iVar6;
  bool bVar7;
  int iVar8;
  int *piVar9;
  int iVar10;
  int iVar11;
  int iVar12;
  int iVar13;
  int iVar14;
  int iVar15;
  int local_20;
  int *local_c;
  int local_8;
  
  iVar1 = param_1[0x13];
  local_8 = 0;
  iVar2 = *(int *)(param_1[0x4c4af] + 4);
  iVar8 = (int)*(short *)(iVar1 + 2);
  do {
    if ((((local_8 == 0) || (local_8 == 1)) || (local_8 == 2)) || (local_8 == 7)) {
      iVar3 = param_1[(local_8 + 0x1223c) * 4];
      iVar4 = param_1[local_8 * 4 + 0x488f2];
      iVar5 = param_1[local_8 * 4 + 0x488f3];
      if ((iVar4 != 0) && (iVar5 != 0)) {
        iVar12 = *(int *)(&DAT_1007d66c + local_8 * 4);
        iVar15 = *(int *)(&DAT_1007d690 + local_8 * 4);
        switch(local_8) {
        case 0:
          local_c = param_1 + 0x48440;
          break;
        case 1:
          local_c = param_1 + 0x48508;
          break;
        case 2:
          local_c = param_1 + 0x485d0;
          break;
        case 7:
          local_c = param_1 + 0x48828;
        }
        iVar11 = 0;
        iVar14 = 0;
        local_20 = 0;
        if (0 < iVar3) {
          do {
            if (iVar11 < iVar8) {
              piVar9 = (int *)(iVar11 * 0x3c0 + 0x64c + iVar1);
              do {
                if (*(int *)(iVar4 + local_20 * 4) <= *piVar9) break;
                iVar11 = iVar11 + 1;
                piVar9 = piVar9 + 0xf0;
              } while (iVar11 < iVar8);
            }
            if (local_20 == 0) {
              if (iVar14 < iVar11) {
                iVar14 = iVar11;
              }
            }
            else {
              iVar6 = *(int *)(iVar5 + -4 + local_20 * 4);
              if (iVar14 < iVar11) {
                piVar9 = local_c + iVar14;
                iVar13 = iVar11 - iVar14;
                iVar14 = iVar14 + iVar13;
                do {
                  iVar10 = iVar6;
                  if (iVar6 < 0) {
                    iVar10 = *piVar9;
                  }
                  if ((local_8 != 7) || (iVar10 != -1)) {
                    if (iVar10 < iVar15) {
                      iVar10 = iVar15;
                    }
                    if (iVar12 < iVar10) {
                      iVar10 = iVar12;
                    }
                  }
                  *piVar9 = iVar10;
                  piVar9 = piVar9 + 1;
                  iVar13 = iVar13 + -1;
                } while (iVar13 != 0);
              }
            }
            if ((local_20 == iVar3 + -1) && (iVar14 < iVar8)) {
              iVar14 = iVar8;
            }
            local_20 = local_20 + 1;
          } while (local_20 < iVar3);
        }
      }
    }
    local_8 = local_8 + 1;
  } while (local_8 < 8);
  local_8 = 0;
  do {
    if ((local_8 == 3) || (local_8 == 4)) {
      iVar3 = param_1[(local_8 + 0x1223c) * 4];
      iVar4 = param_1[local_8 * 4 + 0x488f2];
      iVar5 = param_1[local_8 * 4 + 0x488f3];
      if ((iVar4 != 0) && (0 < iVar3)) {
        if (local_8 == 3) {
          local_c = param_1 + 0x48698;
        }
        else if (local_8 == 4) {
          local_c = param_1 + 0x48760;
        }
        iVar12 = 0;
        if (0 < iVar8) {
LAB_10022b88:
          if (local_8 == 3) {
            if (iVar12 == 0) {
              iVar15 = *(int *)(iVar1 + 0x64c);
              local_20 = 0;
            }
            else {
              iVar15 = *(int *)(iVar12 * 0x3c0 + 0x64c + iVar1);
              local_20 = *(int *)(iVar12 * 0x3c0 + 0x290 + iVar1);
            }
          }
          else if (iVar12 == iVar8 + -1) {
            local_20 = *(int *)(iVar12 * 0x3c0 + 0x650 + iVar1);
            iVar15 = param_1[1] + iVar2;
          }
          else {
            iVar15 = *(int *)(iVar12 * 0x3c0 + 0xa0c + iVar1);
            local_20 = *(int *)(iVar12 * 0x3c0 + 0x650 + iVar1);
          }
          bVar7 = false;
          iVar14 = param_1[local_8 * 4 + 0x488f1];
LAB_10022c05:
          do {
            iVar11 = iVar14;
            if (iVar14 < iVar3) {
              piVar9 = (int *)(iVar4 + iVar14 * 4);
              do {
                if ((local_20 <= *piVar9) && (iVar11 = iVar14, *piVar9 <= iVar15)) break;
                iVar14 = iVar14 + 1;
                piVar9 = piVar9 + 1;
                iVar11 = iVar14;
              } while (iVar14 < iVar3);
            }
            if (iVar11 == iVar3) goto LAB_10022c9f;
            if (local_8 == 3) {
              if (bVar7) {
                local_c[iVar12] = local_c[iVar12] + *(int *)(iVar5 + iVar11 * 4);
              }
              else {
                local_c[iVar12] = *(int *)(iVar5 + iVar11 * 4);
              }
              iVar14 = iVar11 + 1;
              if (local_c[iVar12] < DAT_1007d69c) {
                local_c[iVar12] = DAT_1007d69c;
              }
              if (DAT_1007d678 < local_c[iVar12]) {
                local_c[iVar12] = DAT_1007d678;
                bVar7 = true;
                param_1[0x488fd] = iVar14;
                goto LAB_10022c05;
              }
            }
            else {
              local_c[iVar12] = *(int *)(iVar5 + iVar11 * 4);
            }
            iVar14 = iVar11 + 1;
            bVar7 = true;
            param_1[local_8 * 4 + 0x488f1] = iVar14;
          } while( true );
        }
LAB_10022cb4:
        if (local_8 == 3) {
          iVar12 = param_1[1] + iVar2;
          param_1[0x48910] = -1;
          if (iVar12 == *param_1) {
            bVar7 = false;
            iVar15 = *(int *)(iVar8 * 0x3c0 + 0x290 + iVar1);
            iVar14 = param_1[0x488fd];
            do {
              if (iVar14 < iVar3) {
                piVar9 = (int *)(iVar4 + iVar14 * 4);
                do {
                  if ((iVar15 <= *piVar9) && (*piVar9 <= iVar12)) break;
                  iVar14 = iVar14 + 1;
                  piVar9 = piVar9 + 1;
                } while (iVar14 < iVar3);
              }
              if (iVar14 == iVar3) break;
              iVar11 = *(int *)(iVar5 + iVar14 * 4);
              if (bVar7) {
                param_1[0x48910] = param_1[0x48910] + iVar11;
              }
              else {
                param_1[0x48910] = iVar11;
              }
              iVar14 = iVar14 + 1;
              if (param_1[0x48910] < DAT_1007d69c) {
                param_1[0x48910] = DAT_1007d69c;
              }
              if (DAT_1007d678 < param_1[0x48910]) {
                param_1[0x48910] = DAT_1007d678;
              }
              param_1[0x488fd] = iVar14;
              bVar7 = true;
            } while( true );
          }
        }
      }
    }
    local_8 = local_8 + 1;
    if (7 < local_8) {
      return 1;
    }
  } while( true );
LAB_10022c9f:
  iVar12 = iVar12 + 1;
  if (iVar8 <= iVar12) goto LAB_10022cb4;
  goto LAB_10022b88;
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



===== 0x10025fc0 =====
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



===== 0x10016c90 =====
Function: FUN_10016c90 @ 10016c90

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



===== 0x10017e80 =====
Function: FUN_10017e80 @ 10017e80

undefined4 __cdecl FUN_10017e80(undefined4 param_1,int param_2)

{
  int iVar1;
  undefined2 *puVar2;
  int iVar3;
  int iVar4;
  int iVar5;
  undefined4 *puVar6;
  byte *pbVar7;
  undefined1 local_14 [12];
  undefined2 *local_8;
  
  iVar1 = param_2;
  iVar4 = *(int *)(param_2 + 0x1312bc);
  if (*(char *)(param_2 + 0x21) == '\a') {
    FUN_10025e40(*(undefined4 **)(param_2 + 0x1c),&DAT_1007bfc4);
  }
  iVar3 = *(int *)(param_2 + 0x4c);
  iVar5 = 0;
  param_2 = 0;
  if (0 < *(short *)(iVar3 + 2)) {
    local_8 = (undefined2 *)(iVar4 + 0x30);
    do {
      switch(*(undefined1 *)(iVar5 + 0xa09 + iVar3)) {
      case 0x5b:
        break;
      case 0x5c:
        break;
      case 0x5d:
        break;
      case 0x5e:
      }
      FUN_10025e40(*(undefined4 **)(iVar1 + 0x1c),(byte *)s__s__d____1007bfb8);
      iVar4 = 0;
      if (*(char *)(iVar5 + 0x6e1 + *(int *)(iVar1 + 0x4c)) != '\0') {
        do {
          FUN_10025e40(*(undefined4 **)(iVar1 + 0x1c),&DAT_1007bf90);
          iVar3 = iVar5 + *(int *)(iVar1 + 0x4c);
          if (*(char *)(iVar3 + 0x975 + iVar4) == '1') {
            puVar6 = *(undefined4 **)(iVar1 + 0x1c);
            pbVar7 = &DAT_1007bfb4;
LAB_10017f7f:
            FUN_10025e40(puVar6,pbVar7);
          }
          else if (*(char *)(iVar3 + 0x975 + iVar4) == '2') {
            puVar6 = *(undefined4 **)(iVar1 + 0x1c);
            pbVar7 = &DAT_1007bf8c;
            goto LAB_10017f7f;
          }
          iVar4 = iVar4 + 1;
        } while (iVar4 < (int)(uint)*(byte *)(iVar5 + 0x6e1 + *(int *)(iVar1 + 0x4c)));
      }
      puVar2 = local_8;
      FUN_10018020((byte)*local_8,local_14);
      FUN_10025e40(*(undefined4 **)(iVar1 + 0x1c),(byte *)s_____s__1007bfa8);
      FUN_10025e40(*(undefined4 **)(iVar1 + 0x1c),(byte *)s___d__d__1007bf9c);
      iVar3 = *(int *)(iVar1 + 0x4c);
      param_2 = param_2 + 1;
      local_8 = puVar2 + 0x4a;
      iVar5 = iVar5 + 0x3c0;
    } while (param_2 < *(short *)(iVar3 + 2));
  }
  return CONCAT22((short)((uint)iVar3 >> 0x10),1);
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



