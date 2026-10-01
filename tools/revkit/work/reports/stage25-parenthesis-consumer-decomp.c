===== 0x1003e4a0 =====
Function: FUN_1003e4a0 @ 1003e4a0

uint __cdecl FUN_1003e4a0(short *param_1,uint *param_2,undefined **param_3,byte *param_4)

{
  byte bVar1;
  undefined **ppuVar2;
  char cVar3;
  char *pcVar4;
  uint *puVar5;
  undefined4 uVar6;
  byte *pbVar7;
  undefined4 uVar8;
  int iVar9;
  uint uVar10;
  uint uVar11;
  undefined **ppuVar12;
  undefined **ppuVar13;
  uint *puVar14;
  char *pcVar15;
  byte *pbVar16;
  bool bVar17;
  short sVar18;
  byte local_694 [12];
  int local_688;
  short local_67c;
  char local_542;
  uint local_414 [2];
  byte *local_40c;
  uint local_408;
  short local_400;
  short local_3fe;
  undefined4 local_194;
  uint *local_94;
  uint local_90 [4];
  uint local_80;
  int local_74;
  byte local_70 [35];
  char acStack_4d [33];
  undefined4 local_2c;
  undefined **local_28;
  byte *local_24;
  undefined4 local_20;
  byte local_1c [8];
  undefined4 local_14;
  uint local_10;
  uint *local_c;
  undefined **local_8;
  
  local_c = *(uint **)(param_1 + 0x1cf2);
  ppuVar13 = (undefined **)0x0;
  local_10 = 0;
  local_28 = param_3;
  *local_c = 0;
  local_c[1] = 0x14;
  puVar5 = local_c + 2;
  for (iVar9 = 0x640; iVar9 != 0; iVar9 = iVar9 + -1) {
    *puVar5 = 0;
    puVar5 = puVar5 + 1;
  }
  local_24 = param_4;
  local_194 = local_194 & 0xffffff00;
  pbVar16 = local_694;
  for (iVar9 = 0x140; iVar9 != 0; iVar9 = iVar9 + -1) {
    pbVar16[0] = 0;
    pbVar16[1] = 0;
    pbVar16[2] = 0;
    pbVar16[3] = 0;
    pbVar16 = pbVar16 + 4;
  }
  local_1c[0] = 0;
  acStack_4d[1] = 0;
  ppuVar12 = (undefined **)0x0;
  local_8 = (undefined **)0x0;
  local_2c = 0;
  local_14 = 0;
  local_20 = FUN_10040250((int *)local_694,param_1,(int *)param_2,(char *)param_3,(int)param_4);
  if (local_20 == 0) goto LAB_1003e5bc;
  if (local_20 == 0xffffffff) {
    bVar17 = local_542 == '$';
    if (bVar17) {
      acStack_4d[1] = 0x24;
      local_1c[0] = 0x24;
      local_8 = (undefined **)0x1;
    }
    ppuVar13 = (undefined **)(uint)bVar17;
    ppuVar12 = (undefined **)(uint)bVar17;
    local_20 = 0;
  }
  else {
    param_3 = (undefined **)((int)param_3 + local_20);
    param_4 = param_4 + local_20;
  }
  if (local_67c != 0) {
    local_14 = 1;
  }
  uVar10 = local_20;
  if (ppuVar12 < (undefined **)0x5) {
    do {
      if ((int)uVar10 < 1) {
LAB_1003e5bc:
        uVar10 = param_2[5];
        puVar5 = param_2;
        puVar14 = local_90;
        for (iVar9 = 0x11; iVar9 != 0; iVar9 = iVar9 + -1) {
          *puVar14 = *puVar5;
          puVar5 = puVar5 + 1;
          puVar14 = puVar14 + 1;
        }
        local_10 = uVar10;
        ppuVar13 = local_8;
        uVar10 = local_20;
      }
      else {
        local_10 = FUN_1005a350((int *)local_90,(byte *)param_3,param_4,0xb,
                                *(int *)(param_1 + 0x1cf4));
        if ((local_10 == 0) || (local_74 == 9)) break;
      }
      ppuVar2 = local_8;
      if ((local_90[1] == 0) && (2 < (int)local_90[0])) break;
      if ((((local_74 == 3) || (local_74 == 7)) && (local_80 == 1)) && (local_70[0] == 0x2f)) {
        local_1c[(int)ppuVar12] = 0x2f;
        pcVar4 = acStack_4d + (int)ppuVar2 + 1;
        *pcVar4 = '/';
        local_194 = DAT_10078b04;
        goto LAB_1003ec05;
      }
      if ((local_90[0] != 0) || (2 < (int)local_80)) goto LAB_1003e717;
      if (((local_70[0] == 0x2e) &&
          ((0 < (int)ppuVar12 && (*(char *)((int)&local_20 + 3 + (int)ppuVar12) == 'A')))) &&
         (acStack_4d[(int)ppuVar13] != '\0')) {
        pcVar4 = _strchr(&DAT_1008053c,(int)acStack_4d[(int)ppuVar13]);
        if (pcVar4 == (char *)0x0) {
          if (local_90[0] != 0) goto LAB_1003e717;
          goto LAB_1003e6b6;
        }
        acStack_4d[(int)ppuVar13 + 1] = local_70[0];
        param_3 = (undefined **)((int)param_3 + local_10);
        param_4 = param_4 + local_10;
        *(undefined2 *)((int)ppuVar12 * 0x140 + -0x122 + (int)local_c) = 1;
        local_8 = (undefined **)((int)ppuVar13 + 1);
      }
      else {
LAB_1003e6b6:
        uVar11 = local_10;
        if (((int)local_80 < 3) && (local_70[0] == 0x2e)) {
          if (0 < (int)ppuVar13) {
            if (acStack_4d[(int)ppuVar13] == 'W') {
              acStack_4d[(int)ppuVar13 + 1] = '.';
              param_4 = param_4 + uVar11;
              param_3 = (undefined **)((int)param_3 + uVar11);
              local_40c = local_40c + 1;
              local_408 = local_408 + 1;
              local_20 = uVar10 + uVar11;
              local_3fe = 1;
              local_8 = (undefined **)((int)ppuVar13 + 1);
              uVar10 = uVar10 + uVar11;
              goto LAB_1003ecd9;
            }
            goto LAB_1003e717;
          }
        }
        else {
LAB_1003e717:
          if (0 < (int)ppuVar13) {
            if (local_70[0] == 0x2a) {
              if (local_80 == 1) goto LAB_1003e7d2;
LAB_1003e73d:
              iVar9 = FUN_1001c1f0((char *)local_70,&DAT_100803b0);
              if (((iVar9 == 0) && (local_80 == 2)) ||
                 (((local_70[0] == 0xb7 || (local_70[0] == 0xd7)) && (local_80 == 1))))
              goto LAB_1003e7d2;
            }
            else {
              if ((local_70[0] != 0x78) || (local_80 != 1)) goto LAB_1003e73d;
LAB_1003e7d2:
              ppuVar2 = local_8;
              if ((int)local_90[0] < 2) {
                local_1c[(int)ppuVar12] = 0x4f;
                pcVar4 = acStack_4d + (int)ppuVar2 + 1;
                sVar18 = *(short *)(&DAT_1007e388 + (char)local_70[0] * 2);
                *pcVar4 = 'O';
                if (sVar18 == 0x78) {
                  local_194 = CONCAT22(local_194._2_2_,DAT_1007d344);
                }
                else {
                  local_194._0_3_ = CONCAT12(DAT_100776b6,DAT_100776b4);
                }
                goto LAB_1003ec05;
              }
            }
            iVar9 = FUN_1001c1f0((char *)local_70,&DAT_100803ac);
            ppuVar2 = local_8;
            if ((((iVar9 == 0) && (local_80 == 2)) || ((local_70[0] == 0xf7 && (local_80 == 1)))) &&
               ((int)local_90[0] < 2)) {
              local_1c[(int)ppuVar12] = 0x4f;
              local_194 = DAT_10078b04;
              pcVar4 = acStack_4d + (int)ppuVar2 + 1;
              *pcVar4 = 'O';
              goto LAB_1003ec05;
            }
          }
        }
        if (((((int)ppuVar12 < 1) || (*(char *)((int)&local_20 + 3 + (int)ppuVar12) != 'A')) ||
            ((local_80 != 1 && (local_80 != 2)))) ||
           (((local_74 != 3 || (local_90[0] != 0)) ||
            (iVar9 = FUN_10056150(0x10095c48,(char *)local_70,DAT_1009b088,0x49), puVar5 = local_c,
            uVar11 = local_10, iVar9 < 0)))) {
          if (local_408 == 0) {
            if (local_74 == 1) {
              iVar9 = FUN_10056150(0x10095d98,(char *)local_70,DAT_1009b08c,0x49);
              uVar11 = local_10;
              if (-1 < iVar9) {
                local_414[1] = local_90[2];
                local_414[0] = local_90[0];
                local_408 = local_80;
                param_4 = param_4 + local_10;
                acStack_4d[(int)ppuVar13 + 1] = 'W';
                local_400 = (short)iVar9 + 10;
                uVar10 = local_20 + uVar11;
                param_3 = (undefined **)((int)param_3 + uVar11);
                local_20 = uVar10;
                local_8 = (undefined **)((int)ppuVar13 + 1);
                local_40c = param_4;
                goto LAB_1003ecd9;
              }
              goto LAB_1003e968;
            }
LAB_1003e975:
            if ((local_74 == 3) || (local_74 == 6)) goto LAB_1003e983;
LAB_1003ea45:
            if (local_74 == 1) {
              uVar8 = FUN_10040ab0((char *)&local_194,DAT_100a7a60,(char *)local_70,(short)local_14)
              ;
              ppuVar2 = local_8;
              sVar18 = (short)local_14;
              if ((short)uVar8 == 0) {
                if ((local_74 != 1) ||
                   (uVar6 = FUN_10040ab0((char *)&local_194,DAT_100a7a64,(char *)local_70,
                                         (short)local_14), uVar8 = local_14, pbVar16 = local_24,
                   ppuVar2 = local_28, (short)uVar6 == 0)) goto LAB_1003eb54;
                iVar9 = uVar10 + local_10;
                pcVar4 = acStack_4d + (int)local_8 + 1;
                *pcVar4 = 'C';
                local_1c[(int)ppuVar12] = 0x41;
                sVar18 = (short)uVar8;
                FUN_1003fbf0((int)param_1,(char *)&local_194,sVar18,local_70,(char *)ppuVar2,
                             (int)pbVar16,(int *)0x42,iVar9);
              }
              else {
                local_1c[(int)ppuVar12] = 0x41;
                pcVar4 = acStack_4d + (int)ppuVar2 + 1;
                *pcVar4 = 'P';
              }
              if (sVar18 != 0) {
                local_14 = 0;
                *(undefined2 *)(local_c + (int)ppuVar12 * 0x50 + 8) = 1;
              }
              if ((0 < (int)local_408) &&
                 ((acStack_4d[(int)local_8] == 'W' ||
                  ((acStack_4d[(int)local_8] == '.' && (local_3fe != 0)))))) {
                puVar5 = local_c + (int)ppuVar12 * 0x50;
                goto LAB_1003ea11;
              }
            }
            else {
LAB_1003eb54:
              uVar8 = FUN_10040b60((char *)&local_194,(char *)local_70,(int)param_3,uVar10,
                                   (int *)&local_10,*(int *)(param_1 + 0x1cf4));
              ppuVar2 = local_8;
              if ((short)uVar8 == 0) {
                if ((((ppuVar12 != (undefined **)0x2) || (local_1c[0] != 0x24)) ||
                    (local_1c[1] != '/')) || ((local_74 != 1 || (local_c[0x52] != 0)))) break;
                local_1c[2] = 0x41;
                pcVar4 = acStack_4d + (int)local_8 + 1;
                *pcVar4 = 'a';
              }
              else {
                bVar17 = (short)local_14 != 0;
                local_1c[(int)ppuVar12] = 0x41;
                pcVar4 = acStack_4d + (int)ppuVar2 + 1;
                *pcVar4 = 'E';
                if (bVar17) {
                  local_14 = 0;
                  *(undefined2 *)(local_c + (int)ppuVar12 * 0x50 + 8) = 1;
                }
                if (0 < (int)local_408) {
                  cVar3 = acStack_4d[(int)ppuVar2];
                  goto LAB_1003e9ec;
                }
              }
            }
          }
          else {
LAB_1003e968:
            if (local_74 != 1) goto LAB_1003e975;
LAB_1003e983:
            uVar8 = FUN_10040ab0((char *)&local_194,DAT_100a7a5c,(char *)local_70,(short)local_14);
            ppuVar2 = local_8;
            if ((short)uVar8 == 0) goto LAB_1003ea45;
            bVar17 = (short)local_14 != 0;
            local_1c[(int)ppuVar12] = 0x41;
            pcVar4 = acStack_4d + (int)ppuVar2 + 1;
            *pcVar4 = 'E';
            if (bVar17) {
              local_14 = 0;
              *(undefined2 *)(local_c + (int)ppuVar12 * 0x50 + 8) = 1;
            }
            if ((int)local_408 < 1) goto LAB_1003ec05;
            cVar3 = acStack_4d[(int)local_8];
LAB_1003e9ec:
            if ((cVar3 == 'W') || ((cVar3 == '.' && (local_3fe != 0)))) {
              puVar5 = local_c + (int)ppuVar12 * 0x50;
LAB_1003ea11:
              puVar5[3] = local_414[1];
              puVar5[2] = local_414[0];
              *(short *)(puVar5 + 7) = local_400;
              puVar5 = local_414;
              for (iVar9 = 0x50; iVar9 != 0; iVar9 = iVar9 + -1) {
                *puVar5 = 0;
                puVar5 = puVar5 + 1;
              }
            }
          }
LAB_1003ec05:
          if (local_c[(int)ppuVar12 * 0x50 + 3] == 0) {
            local_c[(int)ppuVar12 * 0x50 + 3] = local_90[2];
            local_c[(int)ppuVar12 * 0x50 + 2] = local_90[0];
          }
          local_c[(int)ppuVar12 * 0x50 + 4] = (uint)(param_4 + local_10);
          local_c[(int)ppuVar12 * 0x50 + 5] = local_80;
          cVar3 = *pcVar4;
          *(undefined1 *)(local_c + (int)ppuVar12 * 0x50 + 6) = (undefined1)local_74;
          *(char *)((int)local_c + (int)ppuVar12 * 0x140 + 0x1a) = cVar3;
          local_c[(int)ppuVar12 * 0x50 + 9] = local_10 + (int)param_3;
          uVar10 = 0xffffffff;
          pbVar16 = local_70;
          do {
            pbVar7 = pbVar16;
            if (uVar10 == 0) break;
            uVar10 = uVar10 - 1;
            pbVar7 = pbVar16 + 1;
            bVar1 = *pbVar16;
            pbVar16 = pbVar7;
          } while (bVar1 != 0);
          uVar10 = ~uVar10;
          local_94 = local_c + (int)ppuVar12 * 0x50 + 10;
          puVar5 = (uint *)(pbVar7 + -uVar10);
          puVar14 = local_c + (int)ppuVar12 * 0x50 + 10;
          for (uVar11 = uVar10 >> 2; uVar11 != 0; uVar11 = uVar11 - 1) {
            *puVar14 = *puVar5;
            puVar5 = puVar5 + 1;
            puVar14 = puVar14 + 1;
          }
          for (uVar10 = uVar10 & 3; uVar10 != 0; uVar10 = uVar10 - 1) {
            *(char *)puVar14 = (char)*puVar5;
            puVar5 = (uint *)((int)puVar5 + 1);
            puVar14 = (uint *)((int)puVar14 + 1);
          }
          if ((char)local_194 != '\0') {
            uVar10 = 0xffffffff;
            pcVar4 = (char *)&local_194;
            do {
              pcVar15 = pcVar4;
              if (uVar10 == 0) break;
              uVar10 = uVar10 - 1;
              pcVar15 = pcVar4 + 1;
              cVar3 = *pcVar4;
              pcVar4 = pcVar15;
            } while (cVar3 != '\0');
            uVar10 = ~uVar10;
            pcVar4 = pcVar15 + -uVar10;
            pcVar15 = (char *)((int)local_c + (int)ppuVar12 * 0x140 + 0x46);
            for (uVar11 = uVar10 >> 2; uVar11 != 0; uVar11 = uVar11 - 1) {
              *(undefined4 *)pcVar15 = *(undefined4 *)pcVar4;
              pcVar4 = pcVar4 + 4;
              pcVar15 = pcVar15 + 4;
            }
            for (uVar10 = uVar10 & 3; uVar10 != 0; uVar10 = uVar10 - 1) {
              *pcVar15 = *pcVar4;
              pcVar4 = pcVar4 + 1;
              pcVar15 = pcVar15 + 1;
            }
            local_194 = local_194 & 0xffffff00;
          }
          local_20 = local_20 + local_10;
          param_4 = param_4 + local_10;
          param_3 = (undefined **)((int)param_3 + local_10);
          ppuVar12 = (undefined **)((int)ppuVar12 + 1);
          local_8 = (undefined **)((int)local_8 + 1);
          uVar10 = local_20;
        }
        else {
          param_3 = (undefined **)((int)param_3 + local_10);
          local_20 = uVar10 + local_10;
          param_4 = param_4 + local_10;
          acStack_4d[(int)ppuVar13 + 1] = 's';
          puVar5[(int)ppuVar12 * 0x50 + -0x4c] = puVar5[(int)ppuVar12 * 0x50 + -0x4c] + uVar11;
          puVar5[(int)ppuVar12 * 0x50 + -0x4b] = puVar5[(int)ppuVar12 * 0x50 + -0x4b] + local_10;
          *(short *)(puVar5 + (int)ppuVar12 * 0x50 + -0x49) = (short)iVar9 + 1;
          local_8 = (undefined **)((int)ppuVar13 + 1);
          uVar10 = local_20;
        }
      }
LAB_1003ecd9:
      ppuVar13 = local_8;
    } while ((int)ppuVar12 < 5);
  }
  local_1c[(int)ppuVar12] = 0;
  bVar1 = local_1c[0];
  bVar17 = local_1c[0] == 0x24;
  acStack_4d[(int)ppuVar13 + 1] = '\0';
  if (bVar17) {
    param_2 = (uint *)0x1;
  }
  else {
    if (bVar1 != 0x41) {
      return 0xffffffff;
    }
    param_2 = (uint *)0x2;
  }
  param_3 = (undefined **)0x0;
  param_4 = (byte *)0x0;
  local_8 = ppuVar12;
  do {
    if ((int)ppuVar12 < 1) {
      return 0xffffffff;
    }
    local_8 = ppuVar12;
    local_1c[(int)ppuVar12] = 0;
    puVar5 = local_c;
    if (param_2 == (uint *)0x1) {
      if (ppuVar12 == (undefined **)0x3) {
        param_4 = (byte *)0x1;
        ppuVar12 = (undefined **)0x0;
        goto LAB_1003ed9a;
      }
    }
    else {
      ppuVar13 = param_3;
      pbVar16 = param_4;
      if (param_2 == (uint *)0x2) {
        if (ppuVar12 == (undefined **)0x1) {
          param_4 = (byte *)0x2;
        }
        else if (ppuVar12 == (undefined **)0x2) {
          param_4 = (byte *)0x3;
        }
        else {
          if (ppuVar12 != (undefined **)0x3) {
            ppuVar13 = ppuVar12;
            pbVar16 = DAT_10080520;
            if (ppuVar12 != (undefined **)0x5) goto LAB_1003ef3a;
            goto LAB_1003ed88;
          }
          param_4 = (byte *)0x5;
        }
      }
      else {
LAB_1003ed88:
        param_4 = pbVar16;
        param_3 = ppuVar13;
        ppuVar12 = param_3;
        if ((int)param_4 <= (int)param_3) goto LAB_1003ef3a;
      }
LAB_1003ed9a:
      ppuVar13 = &PTR_DAT_10080508 + (int)ppuVar12;
      do {
        param_3 = ppuVar13;
        pbVar16 = local_1c;
        pbVar7 = *param_3;
        do {
          bVar1 = *pbVar7;
          bVar17 = bVar1 < *pbVar16;
          if (bVar1 != *pbVar16) {
LAB_1003ede8:
            iVar9 = (1 - (uint)bVar17) - (uint)(bVar17 != 0);
            goto LAB_1003eded;
          }
          if (bVar1 == 0) break;
          bVar1 = pbVar7[1];
          bVar17 = bVar1 < pbVar16[1];
          if (bVar1 != pbVar16[1]) goto LAB_1003ede8;
          pbVar7 = pbVar7 + 2;
          pbVar16 = pbVar16 + 2;
        } while (bVar1 != 0);
        iVar9 = 0;
LAB_1003eded:
        if (iVar9 == 0) {
          uVar10 = local_c[(int)local_8 * 0x50 + -0x4c] - (int)local_24;
          if (ppuVar12 == (undefined **)0x1) {
            if (0 < local_688) {
              if ((*(char *)((int)local_c + 0x1a) == 'C') &&
                 (uVar8 = FUN_1003fbf0((int)param_1,(char *)0x0,0,(byte *)(local_c + 10),
                                       (char *)local_28,(int)local_24,(int *)0x42,uVar10),
                 (short)uVar8 != 0)) {
                local_2c = 1;
              }
              else if ((0 < local_688) &&
                      (iVar9 = FUN_1001c2c0((byte *)(puVar5 + 10),&DAT_10077688), iVar9 == 0)) {
                puVar14 = puVar5 + 10;
                goto LAB_1003eef9;
              }
            }
          }
          else if (ppuVar12 == (undefined **)0x3) {
            if ((*(char *)((int)local_c + 0x1a) == 'C') &&
               (uVar8 = FUN_1003fbf0((int)param_1,(char *)0x0,0,(byte *)(local_c + 10),
                                     (char *)local_28,(int)local_24,(int *)0x42,uVar10),
               (short)uVar8 != 0)) {
              local_2c = 1;
            }
            if (*(char *)((int)puVar5 + 0x29a) == 'C') {
              puVar14 = puVar5 + 0xaa;
LAB_1003eef9:
              uVar8 = FUN_1003fbf0((int)param_1,(char *)0x0,0,(byte *)puVar14,(char *)local_28,
                                   (int)local_24,(int *)0x42,uVar10);
              if ((short)uVar8 != 0) {
                local_2c = 1;
              }
            }
          }
          *puVar5 = (uint)local_8;
          uVar8 = FUN_1003f300((int)param_1,(int)puVar5,(int)local_694,(int)ppuVar12,(short)local_2c
                              );
          param_3 = ppuVar12;
          if ((short)uVar8 != 0) {
            uVar11 = FUN_1003ef80(param_1,local_694);
            return -(uint)((short)uVar11 != 0) & uVar10;
          }
          break;
        }
        ppuVar12 = (undefined **)((int)ppuVar12 + 1);
        ppuVar13 = param_3 + 1;
        param_3 = ppuVar12;
      } while ((int)ppuVar12 < (int)param_4);
    }
LAB_1003ef3a:
    ppuVar12 = (undefined **)((int)local_8 - 1);
    local_8 = ppuVar12;
  } while( true );
}



