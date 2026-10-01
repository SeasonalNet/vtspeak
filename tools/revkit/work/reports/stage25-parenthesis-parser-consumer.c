===== 0x10055660 =====
Function: FUN_100544f0 @ 100544f0

/* WARNING: Type propagation algorithm not settling */

uint __cdecl FUN_100544f0(short *param_1,int *param_2,byte *param_3,byte *param_4,short param_5)

{
  char cVar1;
  byte bVar2;
  short sVar3;
  char *pcVar4;
  uint uVar5;
  int *piVar6;
  int iVar8;
  int iVar9;
  int iVar10;
  uint uVar11;
  short *psVar12;
  uint *puVar13;
  char *pcVar14;
  int *piVar15;
  uint *puVar16;
  bool bVar17;
  undefined4 uVar18;
  undefined1 uVar19;
  undefined1 uVar20;
  int local_2a0 [4];
  int local_290;
  int local_284;
  char local_280;
  int local_25c [4];
  int local_24c;
  int local_240;
  char local_23c;
  uint local_218;
  char local_214 [4];
  char local_210 [4];
  char local_20c;
  int local_118 [4];
  int local_108;
  int local_fc;
  char local_f8;
  byte local_d4 [60];
  int local_98 [4];
  int local_88;
  undefined4 local_7c;
  char local_78 [36];
  uint local_54;
  int local_50 [4];
  int local_40;
  int local_34;
  byte local_30 [36];
  int local_c;
  int local_8;
  undefined4 uVar7;
  
  local_54 = 0;
  local_d4[0] = 0;
  local_8 = (int)*param_1;
  local_218 = local_218 & 0xffffff00;
  FUN_1005a350(param_2,param_3,param_4,0x12,*(int *)(param_1 + 0x1cf4));
  iVar10 = param_2[7];
  if ((iVar10 == 4) || (iVar10 == 5)) goto LAB_1005469c;
  iVar8 = param_2[4];
  if ((iVar8 == 0) || ((iVar10 != 3 && (iVar10 != 6)))) {
    return 0xffffffff;
  }
  if (iVar8 != 1) {
    if (iVar8 == 2) {
      cVar1 = (char)param_2[8];
      bVar2 = *(byte *)((int)param_2 + 0x21);
      iVar9 = FUN_10056150(0x10095c48,(char *)(param_2 + 8),DAT_1009b088,0x53);
      if (iVar9 < 0) {
        iVar9 = FUN_10056150(0x10095e58,(char *)(param_2 + 8),DAT_1009b090,0x53);
        if (iVar9 < 0) {
          iVar9 = FUN_10056150(0x10095f18,(char *)(param_2 + 8),DAT_1009b094,0x53);
          if (iVar9 < 0) {
            iVar9 = FUN_10056150(0x10096878,(char *)(param_2 + 8),DAT_1009b098,0x53);
            if (iVar9 < 0) {
              if ((((cVar1 == -0x58) || (cVar1 == -0x57)) && (0xcc < bVar2)) && (bVar2 < 0xe7)) {
                if ((cVar1 == -0x58) && (bVar2 == 0xcf)) {
                  piVar6 = (int *)param_2[5];
                  iVar10 = FUN_1005a350(local_50,(byte *)((int)piVar6 + (int)param_3),
                                        (byte *)((int)piVar6 + (int)param_4),0x12,
                                        *(int *)(param_1 + 0x1cf4));
                  if (((local_34 == 2) &&
                      ((uVar5 = FUN_10061940((char *)local_30), (short)uVar5 != 0 &&
                       (local_50[0] < 3)))) && (iVar8 = FUN_10039410(local_30), -1 < iVar8)) {
                    uVar7 = FUN_100451e0(param_1,s_copyright_100804fc,param_2[2],param_2[3],0x53,
                                         0x44,0x13);
                    param_2 = piVar6;
                    if ((short)uVar7 == 0) {
                      return 0;
                    }
LAB_10055fb0:
                    FUN_10039d50(local_30,(char *)&local_218,local_30);
                    uVar7 = FUN_100451e0(param_1,(char *)&local_218,local_50[2],local_50[3],0x53,
                                         0x44,0x13);
                    sVar3 = (short)uVar7;
LAB_10055886:
                    if (sVar3 == 0) {
                      return 0;
                    }
                    return (int)param_2 + iVar10;
                  }
                }
                local_218._0_2_ = (ushort)(byte)(bVar2 + 0x74);
                uVar7 = FUN_10044fd0(param_1,(char *)&local_218,param_2[2],param_2[3],0x53,0x53,0x13
                                    );
                sVar3 = (short)uVar7;
                goto LAB_100552d1;
              }
              if ((char)param_2[8] == '+') {
                iVar10 = param_2[3];
                iVar8 = param_2[2];
                pcVar4 = s_plus_plus_10081520;
                goto LAB_100552cc;
              }
              bVar17 = (char)param_2[8] == '-';
LAB_100557c8:
              if (!bVar17) goto LAB_1005469c;
              goto LAB_100557ce;
            }
            iVar10 = param_2[3];
            iVar8 = param_2[2];
            pcVar4 = s_plus_or_minus_10096888 + iVar9 * 0x30;
          }
          else {
            iVar10 = param_2[3];
            iVar8 = param_2[2];
            pcVar4 = &DAT_10095f28 + iVar9 * 0x30;
          }
        }
        else {
          iVar10 = param_2[3];
          iVar8 = param_2[2];
          pcVar4 = s_sub_one_10095e68 + iVar9 * 0x30;
        }
      }
      else {
        iVar10 = param_2[3];
        iVar8 = param_2[2];
        pcVar4 = s_to_the_first_10095c58 + iVar9 * 0x30;
      }
    }
    else {
      if (iVar8 < 3) {
        return 0xffffffff;
      }
      if ((char)param_2[8] != '+') {
        if (((char)param_2[8] == '.') && (0 < local_8)) {
          (param_1 + local_8 * 0x4a + -0x2a)[0] = 1;
          (param_1 + local_8 * 0x4a + -0x2a)[1] = 0;
          return param_2[5];
        }
        goto LAB_1005469c;
      }
      iVar10 = param_2[3];
      iVar8 = param_2[2];
      pcVar4 = s_plus_plus_plus_10081510;
    }
LAB_100552cc:
    uVar7 = FUN_100451e0(param_1,pcVar4,iVar8,iVar10,0x53,0x44,0x13);
    sVar3 = (short)uVar7;
LAB_100552d1:
    if (sVar3 == 0) {
      return 0;
    }
    goto LAB_1005469c;
  }
  if ((((iVar10 == 3) && ((char)param_2[8] != '\0')) &&
      (pcVar4 = _strchr(&DAT_100813c4,(int)(char)param_2[8]), pcVar4 != (char *)0x0)) &&
     (param_2[4] == 1)) {
    piVar6 = local_25c;
    for (iVar10 = 0x11; iVar10 != 0; iVar10 = iVar10 + -1) {
      *piVar6 = 0;
      piVar6 = piVar6 + 1;
    }
    piVar6 = local_98;
    for (iVar10 = 0x11; iVar10 != 0; iVar10 = iVar10 + -1) {
      *piVar6 = 0;
      piVar6 = piVar6 + 1;
    }
    local_d4[0] = 0;
    uVar5 = FUN_10063300((char *)local_d4,(char *)(param_3 + -(int)param_4),(char *)param_3,2);
    if ((short)uVar5 != 0) {
      FUN_10063590(local_25c,local_98,local_d4,param_4 + -(int)(short)uVar5,0x12);
    }
  }
  if (param_2[7] == 3) {
    if (((char)param_2[8] == ',') && (param_2[4] == 1)) {
      if ((0 < local_8) &&
         ((iVar10 = FUN_1005a350(local_50,param_3 + param_2[5],param_4 + param_2[5],0x12,
                                 *(int *)(param_1 + 0x1cf4)), iVar10 != 0 && (local_34 != 9)))) {
        iVar10 = FUN_1001c2c0(local_30,&DAT_10078ea0);
        if (((iVar10 != 0) &&
            (iVar10 = FUN_1001c2c0(local_30,(byte *)s_either_1007902c), iVar10 != 0)) ||
           (1 < local_50[0])) {
          (param_1 + local_8 * 0x4a + -0x2a)[0] = 5;
          (param_1 + local_8 * 0x4a + -0x2a)[1] = 0;
          return param_2[5];
        }
        (param_1 + local_8 * 0x4a + -0x2a)[0] = 0xc;
        (param_1 + local_8 * 0x4a + -0x2a)[1] = 0;
      }
      goto LAB_1005469c;
    }
    if (((((char)param_2[8] == '\0') ||
         (pcVar4 = _strchr(&DAT_100813c4,(int)(char)param_2[8]), pcVar4 == (char *)0x0)) ||
        (param_2[4] != 1)) || ((*param_2 < 1 && (local_88 != 0)))) goto LAB_10054a2f;
    piVar6 = (int *)param_2[5];
    local_c = FUN_1005a350(local_50,(byte *)((int)piVar6 + (int)param_3),
                           (byte *)((int)piVar6 + (int)param_4),0x12,*(int *)(param_1 + 0x1cf4));
    if ((local_c == 0) || (local_34 == 9)) goto LAB_1005469c;
    if (local_34 == 3) {
      if (local_40 == 1) {
        iVar10 = FUN_10056150(0x10081b40,(char *)local_30,DAT_1009b080,0x49);
        if (iVar10 < 0) goto LAB_10054869;
        piVar6 = (int *)((int)piVar6 + local_c);
        local_c = FUN_1005a350(local_118,(byte *)((int)piVar6 + (int)param_3),
                               (byte *)((int)piVar6 + (int)param_4),0x12,*(int *)(param_1 + 0x1cf4))
        ;
        if (((local_c == 0) || (local_fc == 9)) ||
           ((local_108 != 1 ||
            (((local_fc != 3 || ((char)param_2[8] != local_f8)) || (local_50[0] != local_118[0])))))
           ) goto LAB_1005469c;
        pcVar4 = s_exclamation_mark_10081b50 + iVar10 * 0x30;
        param_2 = piVar6;
        goto LAB_1005480d;
      }
    }
    else {
LAB_10054869:
      if (((local_34 == 6) && (local_40 == 1)) &&
         (iVar10 = FUN_10056150(0x1009ac30,(char *)local_30,DAT_1009b120,0x49), -1 < iVar10)) {
        piVar6 = (int *)((int)piVar6 + local_c);
        local_c = FUN_1005a350(local_118,(byte *)((int)piVar6 + (int)param_3),
                               (byte *)((int)piVar6 + (int)param_4),0x12,*(int *)(param_1 + 0x1cf4))
        ;
        if (((((local_c == 0) || (local_fc == 9)) || (local_108 < 1)) ||
            ((local_fc != 3 || ((char)param_2[8] != local_f8)))) || (local_40 != 1))
        goto LAB_1005469c;
        pcVar4 = s_exclamation_mark_1009ac40 + iVar10 * 0x30;
        param_2 = piVar6;
LAB_1005480d:
        uVar5 = 0xffffffff;
        do {
          pcVar14 = pcVar4;
          if (uVar5 == 0) break;
          uVar5 = uVar5 - 1;
          pcVar14 = pcVar4 + 1;
          cVar1 = *pcVar4;
          pcVar4 = pcVar14;
        } while (cVar1 != '\0');
        uVar5 = ~uVar5;
        puVar13 = (uint *)(pcVar14 + -uVar5);
        puVar16 = &local_218;
        for (uVar11 = uVar5 >> 2; uVar11 != 0; uVar11 = uVar11 - 1) {
          *puVar16 = *puVar13;
          puVar13 = puVar13 + 1;
          puVar16 = puVar16 + 1;
        }
        for (uVar5 = uVar5 & 3; uVar5 != 0; uVar5 = uVar5 - 1) {
          *(char *)puVar16 = (char)*puVar13;
          puVar13 = (uint *)((int)puVar13 + 1);
          puVar16 = (uint *)((int)puVar16 + 1);
        }
        uVar7 = FUN_100451e0(param_1,(char *)&local_218,local_50[2],local_50[3],0x53,0x44,0x13);
        sVar3 = (short)uVar7;
        goto LAB_10054851;
      }
    }
    if ((1 < local_50[0]) ||
       (uVar5 = FUN_100560a0(&PTR_DAT_100780fc,(char *)local_30,DAT_10078248,0x49),
       0x7fffffff < uVar5)) {
      if (((char)param_2[8] == '\"') &&
         ((local_50[0] < 2 &&
          (iVar10 = FUN_100560a0(&PTR_DAT_1007844c,(char *)local_30,DAT_1007847c,0x49), -1 < iVar10)
          ))) {
        if (0 < local_8) {
          (param_1 + local_8 * 0x4a + -0x2a)[0] = 1;
          (param_1 + local_8 * 0x4a + -0x2a)[1] = 0;
          return param_2[5];
        }
      }
      else if (((((char)param_2[8] == '\'') && (local_34 == 1)) && (local_50[0] == 0)) &&
              (iVar10 = FUN_100560a0(&PTR_DAT_1009ade0,(char *)local_30,DAT_1009b124,0x49),
              -1 < iVar10)) {
        FUN_10063ed6((undefined1 *)&local_218,&DAT_10081564);
        uVar7 = FUN_100451e0(param_1,(char *)&local_218,param_2[2],local_50[3],0x41,0x44,0x15);
        sVar3 = (short)uVar7;
        param_2 = piVar6;
LAB_10054851:
        if (sVar3 == 0) {
          return 0;
        }
        return local_c + (int)param_2;
      }
      goto LAB_1005469c;
    }
  }
  else {
LAB_10054a2f:
    if (((param_2[7] == 3) && ((char)param_2[8] == '\'')) && ((param_2[4] == 1 && (*param_2 == 0))))
    {
      if ((0 < local_88) && (local_7c == 1)) {
LAB_10054ace:
        uVar5 = param_2[5];
        local_c = FUN_1005a350(local_50,param_3 + uVar5,param_4 + uVar5,0x18,
                               *(int *)(param_1 + 0x1cf4));
        if (local_8 < 1) goto LAB_1005469c;
        if (((local_40 < 1) || (local_34 != 1)) || (local_50[0] != 0)) {
          if (local_88 < 3) goto LAB_1005469c;
          iVar8 = FUN_1001c2c0((byte *)(local_78 + local_88 + -2),&DAT_10077688);
          iVar10 = local_8;
          if ((iVar8 != 0) ||
             (((local_34 != 3 && (local_34 != 9)) &&
              ((local_34 != 8 && ((local_34 != 7 && (local_50[0] < 1)))))))) {
            if ((2 < local_88) &&
               (((local_7c == 1 && (*(short *)(&DAT_1007e388 + local_78[local_88 + -1] * 2) == 0x73)
                 ) && (((local_34 == 9 ||
                        (((local_34 == 8 || (local_34 == 7)) || (0 < local_50[0])))) &&
                       (*(int *)(param_1 + local_8 * 0x4a + -0x3c) + 2 < 0x1d)))))) {
              iVar8 = FUN_10063ed6((undefined1 *)&local_218,&DAT_10081558);
              *(int *)(param_1 + iVar10 * 0x4a + -0x3e) = param_2[3];
              *param_1 = (short)local_8 + -1;
              if (*(int *)(param_1 + iVar10 * 0x4a + -0x2c) == 0x19) {
                if (*(int *)(param_1 + iVar10 * 0x4a + -0x36) == 0) {
                  *(undefined4 *)(param_1 + iVar10 * 0x4a + -0x36) =
                       *(undefined4 *)(param_1 + iVar10 * 0x4a + -0x3c);
                  *(int *)(param_1 + iVar10 * 0x4a + -0x34) = local_40 + 1;
                }
                else {
                  *(int *)(param_1 + iVar10 * 0x4a + -0x34) =
                       *(int *)(param_1 + iVar10 * 0x4a + -0x34) + local_40 + 1;
                }
                iVar8 = (int)*param_1;
                if (0 < iVar8) {
                  psVar12 = param_1 + iVar8 * 0x4a + 0xc;
                  do {
                    if (*(int *)(psVar12 + -0x4c) == *(int *)(psVar12 + -2)) {
                      *(undefined4 *)(psVar12 + -0x4a) = *(undefined4 *)psVar12;
                    }
                    psVar12 = psVar12 + -0x4a;
                    iVar8 = iVar8 + -1;
                  } while (iVar8 != 0);
                }
                uVar7 = *(undefined4 *)(param_1 + iVar10 * 0x4a + -0x3e);
                uVar18 = *(undefined4 *)(param_1 + iVar10 * 0x4a + -0x40);
                uVar11 = 0x19;
                uVar19 = 0x59;
                uVar20 = 0x55;
              }
              else {
                uVar11 = CONCAT31((int3)((uint)iVar8 >> 8),(char)param_1[iVar10 * 0x4a + -0x2c]);
                uVar19 = (undefined1)param_1[iVar10 * 0x4a + -0x2e];
                uVar7 = *(undefined4 *)(param_1 + iVar10 * 0x4a + -0x3e);
                uVar18 = *(undefined4 *)(param_1 + iVar10 * 0x4a + -0x40);
                uVar20 = 0x41;
              }
              uVar7 = FUN_100451e0(param_1,(char *)&local_218,uVar18,uVar7,uVar20,uVar19,uVar11);
              if ((short)uVar7 != 0) {
                return uVar5;
              }
              return 0;
            }
            goto LAB_1005469c;
          }
        }
        iVar10 = local_8;
        if (*(int *)(param_1 + local_8 * 0x4a + -0x3c) + 1 + local_40 < 0x1d) {
          iVar8 = local_8 * 0x4a;
          iVar9 = FUN_10062fe0(local_30);
          iVar8 = FUN_10062fe0((byte *)(param_1 + iVar8 + -0x26));
          if (iVar9 + iVar8 < 0x40) {
            if (0 < param_5) {
              *param_1 = *param_1 + 1;
              return uVar5;
            }
            if ((local_34 == 1) && (local_50[0] == 0)) {
              FUN_10063ed6((undefined1 *)&local_218,(byte *)s__s__s_1008155c);
              *(int *)(param_1 + iVar10 * 0x4a + -0x3e) = local_50[3];
              uVar5 = uVar5 + local_c;
            }
            else {
              FUN_10063ed6((undefined1 *)&local_218,&DAT_10081558);
              *(int *)(param_1 + iVar10 * 0x4a + -0x3e) = param_2[3];
            }
            *param_1 = (short)local_8 + -1;
            if (*(int *)(param_1 + iVar10 * 0x4a + -0x2c) == 0x19) {
              if (*(int *)(param_1 + iVar10 * 0x4a + -0x36) == 0) {
                *(undefined4 *)(param_1 + iVar10 * 0x4a + -0x36) =
                     *(undefined4 *)(param_1 + iVar10 * 0x4a + -0x3c);
                local_40 = local_40 + 1;
              }
              else {
                local_40 = *(int *)(param_1 + iVar10 * 0x4a + -0x34) + local_40 + 1;
              }
              *(int *)(param_1 + iVar10 * 0x4a + -0x34) = local_40;
              iVar8 = (int)*param_1;
              if (0 < iVar8) {
                psVar12 = param_1 + iVar8 * 0x4a + 0xc;
                do {
                  if (*(int *)(psVar12 + -0x4c) == *(int *)(psVar12 + -2)) {
                    *(undefined4 *)(psVar12 + -0x4a) = *(undefined4 *)psVar12;
                  }
                  psVar12 = psVar12 + -0x4a;
                  iVar8 = iVar8 + -1;
                } while (iVar8 != 0);
              }
              uVar7 = *(undefined4 *)(param_1 + iVar10 * 0x4a + -0x3e);
              uVar18 = *(undefined4 *)(param_1 + iVar10 * 0x4a + -0x40);
              uVar11 = 0x19;
              uVar20 = 0x59;
              uVar19 = 0x55;
            }
            else {
              uVar11 = CONCAT31((int3)((uint)*(int *)(param_1 + iVar10 * 0x4a + -0x2c) >> 8),
                                (char)param_1[iVar10 * 0x4a + -0x2c]);
              uVar7 = *(undefined4 *)(param_1 + iVar10 * 0x4a + -0x3e);
              uVar18 = *(undefined4 *)(param_1 + iVar10 * 0x4a + -0x40);
              uVar20 = 0x44;
              uVar19 = 0x41;
            }
            uVar7 = FUN_100451e0(param_1,(char *)&local_218,uVar18,uVar7,uVar19,uVar20,uVar11);
            if ((short)uVar7 != 0) {
              return uVar5;
            }
            return 0;
          }
        }
        goto LAB_1005469c;
      }
      if ((((local_24c < 1) ||
           ((local_8 < 1 || (*(int *)(param_1 + local_8 * 0x4a + -0x3e) != local_25c[3])))) ||
          ((*(int *)(param_1 + local_8 * 0x4a + -0x2c) != 0x16 &&
           (*(int *)(param_1 + local_8 * 0x4a + -0x2c) != 0x18)))) ||
         ((local_78[0] != '.' || (local_88 != 1)))) {
        if (0 < local_88) goto LAB_10054e24;
      }
      else {
        if (local_98[0] == 0) goto LAB_10054ace;
LAB_10054e24:
        if ((local_7c == 1) || (local_7c == 2)) goto LAB_1005469c;
      }
      iVar10 = param_2[5];
      iVar8 = FUN_1005a350(local_50,param_3 + iVar10,param_4 + iVar10,0x12,
                           *(int *)(param_1 + 0x1cf4));
      if ((iVar8 != 0) && (local_34 != 9)) {
        if (((local_34 == 1) && (local_50[0] == 0)) &&
           (iVar9 = FUN_100560a0(&PTR_DAT_1009ade0,(char *)local_30,DAT_1009b124,0x49), -1 < iVar9))
        {
          FUN_10063ed6((undefined1 *)&local_218,&DAT_10081564);
          uVar7 = FUN_100451e0(param_1,(char *)&local_218,param_2[2],local_50[3],0x41,0x44,0x15);
          if ((short)uVar7 == 0) {
            return 0;
          }
          return iVar10 + iVar8;
        }
        if (((local_50[0] < 2) &&
            (iVar10 = FUN_100560a0(&PTR_DAT_100780fc,(char *)local_30,DAT_10078248,0x49),
            -1 < iVar10)) && (0 < local_8)) {
          (param_1 + local_8 * 0x4a + -0x2a)[0] = 1;
          (param_1 + local_8 * 0x4a + -0x2a)[1] = 0;
          return param_2[5];
        }
      }
      goto LAB_1005469c;
    }
    if (((param_2[7] == 3) && ((char)param_2[8] != '\0')) &&
       ((pcVar4 = _strchr(&DAT_100813c4,(int)(char)param_2[8]), pcVar4 != (char *)0x0 &&
        (param_2[4] == 1)))) {
      iVar10 = FUN_1005a350(local_50,param_3 + param_2[5],param_4 + param_2[5],0x12,
                            *(int *)(param_1 + 0x1cf4));
      if ((iVar10 == 0) || (local_34 == 9)) goto LAB_1005469c;
      if ((1 < local_50[0]) ||
         (uVar5 = FUN_100560a0(&PTR_DAT_100780fc,(char *)local_30,DAT_10078248,0x49),
         0x7fffffff < uVar5)) {
        if (((((char)param_2[8] == '\"') && (local_50[0] < 2)) &&
            (iVar10 = FUN_100560a0(&PTR_DAT_1007844c,(char *)local_30,DAT_1007847c,0x49),
            -1 < iVar10)) && (0 < local_8)) {
          (param_1 + local_8 * 0x4a + -0x2a)[0] = 1;
          (param_1 + local_8 * 0x4a + -0x2a)[1] = 0;
          return param_2[5];
        }
        goto LAB_1005469c;
      }
      goto LAB_100557ce;
    }
    if ((((param_2[7] == 3) && ((char)param_2[8] != '\0')) &&
        (pcVar4 = _strchr(&DAT_10081554,(int)(char)param_2[8]), pcVar4 != (char *)0x0)) &&
       (param_2[4] == 1)) {
      piVar6 = local_98;
      for (iVar10 = 0x11; iVar10 != 0; iVar10 = iVar10 + -1) {
        *piVar6 = 0;
        piVar6 = piVar6 + 1;
      }
      local_d4[0] = 0;
      uVar5 = FUN_10063300((char *)local_d4,(char *)(param_3 + -(int)param_4),(char *)param_3,1);
      if ((short)uVar5 != 0) {
        FUN_10063510(local_98,local_d4,param_4 + -(int)(short)uVar5,0x12);
      }
      uVar5 = param_2[5];
      local_c = FUN_1005a350(local_2a0,param_3 + uVar5,param_4 + uVar5,0x12,
                             *(int *)(param_1 + 0x1cf4));
      if ((local_c == 0) || (local_284 == 9)) goto LAB_1005469c;
      if (((local_284 == 3) && ((local_280 == '=' && (local_290 == 1)))) && (local_2a0[0] == 0)) {
        uVar5 = uVar5 + local_c;
        local_c = FUN_1005a350(local_50,param_3 + uVar5,param_4 + uVar5,0x12,
                               *(int *)(param_1 + 0x1cf4));
        if ((local_c == 0) || (local_54 = uVar5, local_34 == 9)) goto LAB_1005469c;
      }
      else {
        piVar6 = local_2a0;
        piVar15 = local_50;
        for (iVar10 = 0x11; iVar10 != 0; iVar10 = iVar10 + -1) {
          *piVar15 = *piVar6;
          piVar6 = piVar6 + 1;
          piVar15 = piVar15 + 1;
        }
      }
      if (((local_40 < 1) ||
          (((local_34 != 3 || (local_30[0] == 0)) ||
           (pcVar4 = _strchr(&DAT_10081550,(int)(char)local_30[0]), pcVar4 == (char *)0x0)))) ||
         ((local_40 != 1 || (2 < local_50[0])))) {
        iVar10 = local_8;
        if (local_88 < 1) goto LAB_1005469c;
        if ((((*param_1 < 1) || (*(int *)(param_1 + local_8 * 0x4a + -0x3e) != local_98[3])) ||
            ((local_7c != 2 && (local_7c != 1)))) ||
           ((((1 < *param_2 || (local_40 < 1)) || (local_34 != 2)) || (1 < local_50[0])))) {
          if ((0 < *param_1) && (*(int *)(param_1 + local_8 * 0x4a + -0x3e) == local_98[3])) {
            if (local_7c == 1) {
              if (local_88 != 1) {
                return param_2[5];
              }
            }
            else if ((local_7c != 2) || (uVar5 = FUN_10061940(local_78), (short)uVar5 == 0))
            goto LAB_1005469c;
            if ((((*param_2 < 2) && (0 < local_40)) &&
                ((local_34 == 1 || ((local_34 == 2 || (local_34 == 3)))))) &&
               ((local_50[0] < 2 && (0 < iVar10)))) {
              (param_1 + iVar10 * 0x4a + -0x2a)[0] = 1;
              (param_1 + iVar10 * 0x4a + -0x2a)[1] = 0;
              return param_2[5];
            }
          }
          goto LAB_1005469c;
        }
        if (local_54 != 0) {
          if ((char)param_2[8] == '<') {
            puVar13 = (uint *)s_less_than_or_equal_to_10080438;
            puVar16 = &local_218;
            for (iVar10 = 5; iVar10 != 0; iVar10 = iVar10 + -1) {
              *puVar16 = *puVar13;
              puVar13 = puVar13 + 1;
              puVar16 = puVar16 + 1;
            }
            *(short *)puVar16 = (short)*puVar13;
          }
          else {
            puVar13 = (uint *)s_greater_than_or_equal_to_10080410;
            puVar16 = &local_218;
            for (iVar10 = 6; iVar10 != 0; iVar10 = iVar10 + -1) {
              *puVar16 = *puVar13;
              puVar13 = puVar13 + 1;
              puVar16 = puVar16 + 1;
            }
            *(char *)puVar16 = (char)*puVar13;
          }
          iVar10 = param_2[2];
          goto LAB_1005532f;
        }
        if ((char)param_2[8] == '<') {
          local_218._0_1_ = s_less_than_1008042c[0];
          local_218._1_1_ = s_less_than_1008042c[1];
          local_218._2_1_ = s_less_than_1008042c[2];
          local_218._3_1_ = s_less_than_1008042c[3];
          local_214[0] = s_less_than_1008042c[4];
          local_214[1] = s_less_than_1008042c[5];
          local_214[2] = s_less_than_1008042c[6];
          local_214[3] = s_less_than_1008042c[7];
          local_210[0] = s_less_than_1008042c[8];
          local_210[1] = s_less_than_1008042c[9];
        }
        else {
          local_218._0_1_ = s_greater_than_10080400[0];
          local_218._1_1_ = s_greater_than_10080400[1];
          local_218._2_1_ = s_greater_than_10080400[2];
          local_218._3_1_ = s_greater_than_10080400[3];
          local_214[0] = s_greater_than_10080400[4];
          local_214[1] = s_greater_than_10080400[5];
          local_214[2] = s_greater_than_10080400[6];
          local_214[3] = s_greater_than_10080400[7];
          local_210[0] = s_greater_than_10080400[8];
          local_210[1] = s_greater_than_10080400[9];
          local_210[2] = s_greater_than_10080400[10];
          local_210[3] = s_greater_than_10080400[0xb];
          local_20c = s_greater_than_10080400[0xc];
        }
        iVar10 = param_2[3];
        iVar8 = param_2[2];
        pcVar4 = (char *)&local_218;
      }
      else {
        iVar10 = FUN_1005a350(local_118,param_3 + uVar5 + local_c,param_4 + uVar5 + local_c,0x12,
                              *(int *)(param_1 + 0x1cf4));
        if ((((iVar10 == 0) || (local_fc == 9)) || (local_108 < 1)) ||
           (((local_fc != 2 || (2 < local_118[0])) ||
            ((local_88 != 0 && ((local_7c != 2 && (local_7c != 1)))))))) goto LAB_1005469c;
        if (local_54 != 0) {
          if ((char)param_2[8] == '<') {
            puVar13 = (uint *)s_less_than_or_equal_to_10080438;
            puVar16 = &local_218;
            for (iVar10 = 5; iVar10 != 0; iVar10 = iVar10 + -1) {
              *puVar16 = *puVar13;
              puVar13 = puVar13 + 1;
              puVar16 = puVar16 + 1;
            }
            *(short *)puVar16 = (short)*puVar13;
          }
          else {
            puVar13 = (uint *)s_greater_than_or_equal_to_10080410;
            puVar16 = &local_218;
            for (iVar10 = 6; iVar10 != 0; iVar10 = iVar10 + -1) {
              *puVar16 = *puVar13;
              puVar13 = puVar13 + 1;
              puVar16 = puVar16 + 1;
            }
            *(char *)puVar16 = (char)*puVar13;
          }
          iVar10 = param_2[2];
LAB_1005532f:
          uVar7 = FUN_100451e0(param_1,(char *)&local_218,iVar10,local_2a0[3],0x53,0x44,0x13);
          return -(uint)((short)uVar7 != 0) & local_54;
        }
        if ((char)param_2[8] == '<') {
          local_218._0_1_ = s_less_than_1008042c[0];
          local_218._1_1_ = s_less_than_1008042c[1];
          local_218._2_1_ = s_less_than_1008042c[2];
          local_218._3_1_ = s_less_than_1008042c[3];
          local_214[0] = s_less_than_1008042c[4];
          local_214[1] = s_less_than_1008042c[5];
          local_214[2] = s_less_than_1008042c[6];
          local_214[3] = s_less_than_1008042c[7];
          local_210[0] = s_less_than_1008042c[8];
          local_210[1] = s_less_than_1008042c[9];
        }
        else {
          local_218._0_1_ = s_greater_than_10080400[0];
          local_218._1_1_ = s_greater_than_10080400[1];
          local_218._2_1_ = s_greater_than_10080400[2];
          local_218._3_1_ = s_greater_than_10080400[3];
          local_214[0] = s_greater_than_10080400[4];
          local_214[1] = s_greater_than_10080400[5];
          local_214[2] = s_greater_than_10080400[6];
          local_214[3] = s_greater_than_10080400[7];
          local_210[0] = s_greater_than_10080400[8];
          local_210[1] = s_greater_than_10080400[9];
          local_210[2] = s_greater_than_10080400[10];
          local_210[3] = s_greater_than_10080400[0xb];
          local_20c = s_greater_than_10080400[0xc];
        }
        iVar10 = param_2[3];
        iVar8 = param_2[2];
        pcVar4 = (char *)&local_218;
      }
      goto LAB_100552cc;
    }
    if (param_2[7] != 3) goto LAB_1005469c;
    if (((char)param_2[8] != ']') || (param_2[4] != 1)) {
      if ((((char)param_2[8] == '(') || ((char)param_2[8] == '[')) && (param_2[4] == 1)) {
        uVar5 = param_2[5];
        iVar10 = FUN_10062ea0((char *)(param_3 + uVar5));
        if (iVar10 == 0) {
          return uVar5;
        }
        if ((*(int *)(param_1 + 4) != 0) && (*(int *)(param_1 + 4) <= iVar10)) {
          if (local_8 < 1) {
            return uVar5;
          }
          (param_1 + local_8 * 0x4a + -0x2a)[0] = 1;
          (param_1 + local_8 * 0x4a + -0x2a)[1] = 0;
          return uVar5;
        }
        return uVar5 + iVar10;
      }
      if ((((char)param_2[8] == ')') || ((char)param_2[8] == ']')) && (param_2[4] == 1)) {
        if (*param_2 < 2) {
          piVar6 = local_98;
          for (iVar10 = 0x11; iVar10 != 0; iVar10 = iVar10 + -1) {
            *piVar6 = 0;
            piVar6 = piVar6 + 1;
          }
          local_d4[0] = 0;
          uVar5 = FUN_10063300((char *)local_d4,(char *)(param_3 + -(int)param_4),(char *)param_3,1)
          ;
          if ((short)uVar5 != 0) {
            FUN_10063510(local_98,local_d4,param_4 + -(int)(short)uVar5,0x12);
          }
          if ((local_88 < 1) ||
             (((local_7c != 2 || (uVar5 = FUN_10061940(local_78), (short)uVar5 == 0)) &&
              (local_7c != 1)))) goto LAB_1005469c;
          iVar10 = FUN_1005a350(local_50,param_3 + param_2[5],param_4 + param_2[5],0x12,
                                *(int *)(param_1 + 0x1cf4));
        }
        else {
          iVar10 = FUN_1005a350(local_50,param_3 + param_2[5],param_4 + param_2[5],0x12,
                                *(int *)(param_1 + 0x1cf4));
        }
        if ((iVar10 == 0) || (local_34 == 9)) goto LAB_1005469c;
        goto joined_r0x100557be;
      }
      if (((char)param_2[8] == '.') && (param_2[4] == 1)) {
        piVar6 = (int *)param_2[5];
        iVar10 = FUN_1005a350(local_50,(byte *)((int)piVar6 + (int)param_3),
                              (byte *)((int)piVar6 + (int)param_4),0x12,*(int *)(param_1 + 0x1cf4));
        if ((((iVar10 == 0) || (local_34 == 9)) || (local_50[0] != 0)) ||
           (iVar8 = FUN_1001c2c0(local_30,&DAT_1008154c), iVar8 != 0)) goto LAB_1005469c;
        uVar7 = FUN_100451e0(param_1,s_dot_net_10081544,param_2[2],local_50[3],0x53,0x44,0x13);
        sVar3 = (short)uVar7;
        param_2 = piVar6;
        goto LAB_10055886;
      }
      if ((((char)param_2[8] == ':') && (param_2[4] == 1)) && (*param_2 < 2)) {
        piVar6 = local_25c;
        for (iVar10 = 0x11; iVar10 != 0; iVar10 = iVar10 + -1) {
          *piVar6 = 0;
          piVar6 = piVar6 + 1;
        }
        piVar6 = local_98;
        for (iVar10 = 0x11; iVar10 != 0; iVar10 = iVar10 + -1) {
          *piVar6 = 0;
          piVar6 = piVar6 + 1;
        }
        local_d4[0] = 0;
        uVar5 = FUN_10063300((char *)local_d4,(char *)(param_3 + -(int)param_4),(char *)param_3,2);
        if ((short)uVar5 != 0) {
          FUN_10063590(local_25c,local_98,local_d4,param_4 + -(int)(short)uVar5,0x12);
        }
        if ((((local_24c < 1) || (local_240 != 3)) || ((local_24c != 1 || (local_23c != ':')))) &&
           (0 < local_88)) {
          if ((local_7c == 2) && (uVar5 = FUN_10061940(local_78), (short)uVar5 != 0)) {
            iVar10 = param_2[5];
            local_c = FUN_1005a350(local_50,param_3 + iVar10,param_4 + iVar10,0x12,
                                   *(int *)(param_1 + 0x1cf4));
            if ((local_c != 0) && (local_34 != 9)) {
              if ((local_34 == 2) &&
                 ((uVar5 = FUN_10061940((char *)local_30), (short)uVar5 != 0 &&
                  (*param_2 == local_50[0])))) {
                FUN_1005a350(local_118,param_3 + iVar10 + local_c,param_4 + iVar10 + local_c,0x12,
                             *(int *)(param_1 + 0x1cf4));
                if (((0 < local_108) && (local_fc == 3)) && (local_108 == 1)) {
                  return param_2[5];
                }
              }
              else if (0 < local_8) {
                (param_1 + local_8 * 0x4a + -0x2a)[0] = 1;
                (param_1 + local_8 * 0x4a + -0x2a)[1] = 0;
                return param_2[5];
              }
            }
          }
          else if ((0 < local_88) && ((local_7c == 1 && (0 < local_8)))) {
            (param_1 + local_8 * 0x4a + -0x2a)[0] = 1;
            (param_1 + local_8 * 0x4a + -0x2a)[1] = 0;
            return param_2[5];
          }
        }
        goto LAB_1005469c;
      }
      if ((((char)param_2[8] == '~') && (param_2[4] == 1)) && (*param_2 < 2)) {
        piVar6 = local_98;
        for (iVar10 = 0x11; iVar10 != 0; iVar10 = iVar10 + -1) {
          *piVar6 = 0;
          piVar6 = piVar6 + 1;
        }
        local_d4[0] = 0;
        uVar5 = FUN_10063300((char *)local_d4,(char *)(param_3 + -(int)param_4),(char *)param_3,1);
        if ((short)uVar5 == 0) goto LAB_1005469c;
        FUN_10063510(local_98,local_d4,param_4 + -(int)(short)uVar5,0x12);
        iVar10 = FUN_1005a350(local_50,param_3 + param_2[5],param_4 + param_2[5],0x12,
                              *(int *)(param_1 + 0x1cf4));
        if (((iVar10 == 0) || (local_34 == 9)) || (local_88 < 1)) goto LAB_1005469c;
        if (local_7c == 1) {
          if (local_88 != 1) {
            return param_2[5];
          }
        }
        else if (local_7c != 2) goto LAB_1005469c;
        if (local_40 < 1) goto LAB_1005469c;
        if (local_34 == 1) {
          if (local_40 != 1) {
            return param_2[5];
          }
        }
        else if (local_34 != 2) goto LAB_1005469c;
        if (*param_2 != local_50[0]) goto LAB_1005469c;
        uVar7 = FUN_10044fd0(param_1,(char *)&DAT_10077380,param_2[2],param_2[3],0x53,0x44,0x13);
        sVar3 = (short)uVar7;
      }
      else if (((char)param_2[8] == '@') && (param_2[4] == 1)) {
        uVar7 = FUN_10044fd0(param_1,&DAT_10077694,param_2[2],param_2[3],0x53,0x44,0x13);
        sVar3 = (short)uVar7;
      }
      else {
        if (((char)param_2[8] == '\\') && (param_2[4] == 1)) {
          iVar10 = param_2[3];
          iVar8 = param_2[2];
          pcVar4 = s_back_slash_1008123c;
          goto LAB_100552cc;
        }
        if (((char)param_2[8] == '&') && (param_2[4] == 1)) {
          uVar7 = FUN_10044fd0(param_1,(char *)&DAT_10077400,param_2[2],param_2[3],0x53,0x44,0x13);
          sVar3 = (short)uVar7;
        }
        else {
          if (((char)param_2[8] != '#') || (param_2[4] != 1)) {
            if ((param_2[4] == 1) && ((char)param_2[8] == -0x57)) {
              piVar6 = (int *)param_2[5];
              iVar10 = FUN_1005a350(local_50,(byte *)((int)piVar6 + (int)param_3),
                                    (byte *)((int)piVar6 + (int)param_4),3,
                                    *(int *)(param_1 + 0x1cf4));
              if ((iVar10 == 0) ||
                 ((((local_34 == 9 || (local_34 != 2)) ||
                   (uVar5 = FUN_10061940((char *)local_30), (short)uVar5 == 0)) ||
                  ((2 < local_50[0] || (iVar8 = FUN_10039410(local_30), iVar8 < 0))))))
              goto LAB_1005469c;
              uVar7 = FUN_10044fd0(param_1,s_copyright_100804fc,param_2[2],param_2[3],0x53,0x44,0x13
                                  );
              param_2 = piVar6;
              if ((short)uVar7 == 0) {
                return 0;
              }
              goto LAB_10055fb0;
            }
            if ((param_2[4] != 1) || ((char)param_2[8] != -0x52)) goto LAB_1005469c;
            iVar10 = param_2[3];
            iVar8 = param_2[2];
            pcVar4 = s_registered_trademark_1008152c;
            goto LAB_100552cc;
          }
          iVar10 = FUN_1005a350(local_50,param_3 + param_2[5],param_4 + param_2[5],0x12,
                                *(int *)(param_1 + 0x1cf4));
          if ((((iVar10 == 0) || (local_34 == 9)) || (local_34 != 2)) ||
             ((uVar5 = FUN_10061940((char *)local_30), (short)uVar5 == 0 || (2 < local_50[0]))))
          goto LAB_1005469c;
          uVar7 = FUN_10044fd0(param_1,s_number_100775c4,param_2[2],param_2[3],0x53,0x44,0x13);
          sVar3 = (short)uVar7;
        }
      }
      goto LAB_100552d1;
    }
    piVar6 = local_98;
    for (iVar10 = 0x11; iVar10 != 0; iVar10 = iVar10 + -1) {
      *piVar6 = 0;
      piVar6 = piVar6 + 1;
    }
    local_d4[0] = 0;
    uVar5 = FUN_10063300((char *)local_d4,(char *)(param_3 + -(int)param_4),(char *)param_3,1);
    if ((short)uVar5 != 0) {
      FUN_10063510(local_98,local_d4,param_4 + -(int)(short)uVar5,0x12);
    }
    if ((local_88 < 1) ||
       ((((local_7c != 2 || (uVar5 = FUN_10061940(local_78), (short)uVar5 == 0)) &&
         ((local_7c != 1 || (local_88 != 1)))) ||
        (((iVar10 = FUN_1005a350(local_50,param_3 + param_2[5],param_4 + param_2[5],0x12,
                                 *(int *)(param_1 + 0x1cf4)), iVar10 == 0 || (local_34 == 9)) ||
         (1 < local_50[0])))))) goto LAB_1005469c;
joined_r0x100557be:
    if ((local_34 != 1) && (local_34 != 2)) {
      bVar17 = local_34 == 3;
      goto LAB_100557c8;
    }
  }
LAB_100557ce:
  if (0 < local_8) {
    (param_1 + local_8 * 0x4a + -0x2a)[0] = 1;
    (param_1 + local_8 * 0x4a + -0x2a)[1] = 0;
    return param_2[5];
  }
LAB_1005469c:
  return param_2[5];
}



