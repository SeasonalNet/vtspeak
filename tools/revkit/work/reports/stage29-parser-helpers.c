===== 0x1000d190 =====
Function: FUN_1000d190 @ 1000d190

int __cdecl FUN_1000d190(int *param_1,int param_2,int param_3,short param_4)

{
  byte bVar1;
  int *piVar2;
  char cVar3;
  short sVar4;
  undefined3 extraout_var;
  int iVar5;
  uint uVar6;
  uint uVar7;
  int iVar8;
  byte *pbVar9;
  int *piVar10;
  int *piVar11;
  char *pcVar12;
  char *pcVar13;
  char *pcVar14;
  byte *pbVar15;
  byte local_234 [200];
  byte local_16c [204];
  char local_a0 [68];
  char local_5c [32];
  char local_3c [32];
  undefined4 local_1c;
  uint local_18;
  int local_14;
  int local_10;
  undefined4 local_c;
  uint local_8;
  
  piVar2 = param_1;
  sVar4 = 1;
  local_c = 1;
  local_14 = 0;
  local_10 = 0;
  if (param_3 < 1) {
LAB_1000d415:
    return (int)sVar4;
  }
  param_1 = (int *)(param_2 + 0x2c);
  do {
    uVar6 = 0xffffffff;
    piVar10 = param_1 + 2;
    do {
      piVar11 = piVar10;
      if (uVar6 == 0) break;
      uVar6 = uVar6 - 1;
      piVar11 = (int *)((int)piVar10 + 1);
      iVar8 = *piVar10;
      piVar10 = piVar11;
    } while ((char)iVar8 != '\0');
    uVar6 = ~uVar6;
    pcVar12 = (char *)((int)piVar11 - uVar6);
    pcVar14 = local_3c;
    for (uVar7 = uVar6 >> 2; uVar7 != 0; uVar7 = uVar7 - 1) {
      *(undefined4 *)pcVar14 = *(undefined4 *)pcVar12;
      pcVar12 = pcVar12 + 4;
      pcVar14 = pcVar14 + 4;
    }
    for (uVar6 = uVar6 & 3; uVar6 != 0; uVar6 = uVar6 - 1) {
      *pcVar14 = *pcVar12;
      pcVar12 = pcVar12 + 1;
      pcVar14 = pcVar14 + 1;
    }
    local_8 = (uint)(local_10 == param_3 + -1);
    if ((param_1[-9] < 0x1d) && (*param_1 != 0)) {
      switch(*param_1) {
      case 2:
        pcVar12 = (char *)&DAT_10077700;
        break;
      case 3:
        pcVar12 = (char *)&DAT_100776f4;
        break;
      case 4:
        pcVar12 = (char *)&DAT_100776f8;
        break;
      case 5:
      case 0xc:
        pcVar12 = (char *)&DAT_100776fc;
        break;
      default:
        goto switchD_1000d214_caseD_6;
      case 0xb:
        pcVar12 = &DAT_1007739c;
      }
      uVar6 = 0xffffffff;
      do {
        pcVar14 = pcVar12;
        if (uVar6 == 0) break;
        uVar6 = uVar6 - 1;
        pcVar14 = pcVar12 + 1;
        cVar3 = *pcVar12;
        pcVar12 = pcVar14;
      } while (cVar3 != '\0');
      uVar6 = ~uVar6;
      iVar8 = -1;
      pcVar12 = local_3c;
      do {
        pcVar13 = pcVar12;
        if (iVar8 == 0) break;
        iVar8 = iVar8 + -1;
        pcVar13 = pcVar12 + 1;
        cVar3 = *pcVar12;
        pcVar12 = pcVar13;
      } while (cVar3 != '\0');
      pcVar12 = pcVar14 + -uVar6;
      pcVar14 = pcVar13 + -1;
      for (uVar7 = uVar6 >> 2; uVar7 != 0; uVar7 = uVar7 - 1) {
        *(undefined4 *)pcVar14 = *(undefined4 *)pcVar12;
        pcVar12 = pcVar12 + 4;
        pcVar14 = pcVar14 + 4;
      }
      for (uVar6 = uVar6 & 3; uVar6 != 0; uVar6 = uVar6 - 1) {
        *pcVar14 = *pcVar12;
        pcVar12 = pcVar12 + 1;
        pcVar14 = pcVar14 + 1;
      }
    }
switchD_1000d214_caseD_6:
    uVar6 = 0xffffffff;
    iVar8 = param_1[1];
    pcVar12 = (char *)((int)param_1 + 0x26);
    do {
      pcVar14 = pcVar12;
      if (uVar6 == 0) break;
      uVar6 = uVar6 - 1;
      pcVar14 = pcVar12 + 1;
      cVar3 = *pcVar12;
      pcVar12 = pcVar14;
    } while (cVar3 != '\0');
    uVar6 = ~uVar6;
    pcVar12 = pcVar14 + -uVar6;
    pcVar14 = local_a0;
    for (uVar7 = uVar6 >> 2; uVar7 != 0; uVar7 = uVar7 - 1) {
      *(undefined4 *)pcVar14 = *(undefined4 *)pcVar12;
      pcVar12 = pcVar12 + 4;
      pcVar14 = pcVar14 + 4;
    }
    iVar5 = param_1[-2];
    for (uVar6 = uVar6 & 3; uVar6 != 0; uVar6 = uVar6 - 1) {
      *pcVar14 = *pcVar12;
      pcVar12 = pcVar12 + 1;
      pcVar14 = pcVar14 + 1;
    }
    uVar6 = (uint)((char)iVar5 == 'S');
    local_18 = uVar6;
    do {
      FUN_1000d640(local_3c,local_5c,(short)local_8,(char)iVar8,local_a0,(short)uVar6,
                   (int)(piVar2 + 0x1204a));
      cVar3 = FUN_1000fd20(local_5c);
      local_1c = CONCAT31(extraout_var,cVar3);
      FUN_10003a70(local_5c,local_16c,0,(int)(piVar2 + 0x1204a));
      if ((&stack0x00000000 == (undefined1 *)0x16c) || (local_16c[0] == 0)) {
        local_234[0] = 0;
      }
      else {
        uVar7 = 0xffffffff;
        pbVar9 = local_16c;
        do {
          pbVar15 = pbVar9;
          if (uVar7 == 0) break;
          uVar7 = uVar7 - 1;
          pbVar15 = pbVar9 + 1;
          bVar1 = *pbVar9;
          pbVar9 = pbVar15;
        } while (bVar1 != 0);
        uVar7 = ~uVar7;
        pbVar9 = pbVar15 + -uVar7;
        pbVar15 = local_234;
        for (uVar6 = uVar7 >> 2; uVar6 != 0; uVar6 = uVar6 - 1) {
          *(undefined4 *)pbVar15 = *(undefined4 *)pbVar9;
          pbVar9 = pbVar9 + 4;
          pbVar15 = pbVar15 + 4;
        }
        for (uVar7 = uVar7 & 3; uVar6 = local_18, uVar7 != 0; uVar7 = uVar7 - 1) {
          *pbVar15 = *pbVar9;
          pbVar9 = pbVar9 + 1;
          pbVar15 = pbVar15 + 1;
        }
      }
      sVar4 = (short)*piVar2;
      if (199 < sVar4) {
        local_c = 0;
        break;
      }
      if (sVar4 < 1) {
        iVar5 = -1;
      }
      else {
        iVar5 = (int)(short)piVar2[sVar4 * 0x155 + -0x154];
      }
      FUN_1000d450((undefined2 *)((int)piVar2 + sVar4 * 0x554 + 2),local_5c,local_234,local_14,iVar5
                   ,0,local_1c,(char)iVar8,local_a0,param_4);
      *(short *)piVar2 = (short)*piVar2 + 1;
    } while (local_3c[0] != '\0');
    local_14 = local_14 + 1;
    sVar4 = 0;
    if ((short)local_c == 0) goto LAB_1000d415;
    local_10 = local_10 + 1;
    param_1 = param_1 + 0x25;
    if (param_3 <= local_10) {
      return (int)(short)local_c;
    }
  } while( true );
}



===== 0x1000ea20 =====
Function: FUN_1000ea20 @ 1000ea20

void __cdecl FUN_1000ea20(short *param_1,int param_2,int param_3)

{
  char cVar1;
  int iVar2;
  short *psVar3;
  uint uVar4;
  uint uVar5;
  char *pcVar6;
  char *pcVar7;
  int local_c;
  byte *local_8;
  
  param_1[0x214d1] = 0;
  local_8 = (byte *)0x0;
  if (0 < *param_1) {
    psVar3 = param_1 + 1;
    local_c = 0;
    do {
      *(short *)(local_c + 0x429a8 + (int)param_1) = psVar3[1];
      uVar4 = 0xffffffff;
      pcVar6 = (char *)((int)psVar3 + 5);
      do {
        pcVar7 = pcVar6;
        if (uVar4 == 0) break;
        uVar4 = uVar4 - 1;
        pcVar7 = pcVar6 + 1;
        cVar1 = *pcVar6;
        pcVar6 = pcVar7;
      } while (cVar1 != '\0');
      uVar4 = ~uVar4;
      pcVar6 = pcVar7 + -uVar4;
      pcVar7 = (char *)(local_c + 0x429ad + (int)param_1);
      for (uVar5 = uVar4 >> 2; uVar5 != 0; uVar5 = uVar5 - 1) {
        *(undefined4 *)pcVar7 = *(undefined4 *)pcVar6;
        pcVar6 = pcVar6 + 4;
        pcVar7 = pcVar7 + 4;
      }
      for (uVar4 = uVar4 & 3; uVar4 != 0; uVar4 = uVar4 - 1) {
        *pcVar7 = *pcVar6;
        pcVar6 = pcVar6 + 1;
        pcVar7 = pcVar7 + 1;
      }
      if (*psVar3 < 2) {
        uVar4 = 0xffffffff;
        pcVar6 = (char *)((int)psVar3 + 0x37);
        do {
          pcVar7 = pcVar6;
          if (uVar4 == 0) break;
          uVar4 = uVar4 - 1;
          pcVar7 = pcVar6 + 1;
          cVar1 = *pcVar6;
          pcVar6 = pcVar7;
        } while (cVar1 != '\0');
        uVar4 = ~uVar4;
        pcVar6 = pcVar7 + -uVar4;
        pcVar7 = (char *)(local_c + 0x429cb + (int)param_1);
        for (uVar5 = uVar4 >> 2; uVar5 != 0; uVar5 = uVar5 - 1) {
          *(undefined4 *)pcVar7 = *(undefined4 *)pcVar6;
          pcVar6 = pcVar6 + 4;
          pcVar7 = pcVar7 + 4;
        }
        for (uVar4 = uVar4 & 3; uVar4 != 0; uVar4 = uVar4 - 1) {
          *pcVar7 = *pcVar6;
          pcVar6 = pcVar6 + 1;
          pcVar7 = pcVar7 + 1;
        }
      }
      else {
        FUN_100068b0((byte *)(local_c + 0x429cb + (int)param_1),(int)(param_1 + 1),
                     (byte *)(int)*param_1,local_8,param_3);
      }
      *(char *)(local_c + 0x429ac + (int)param_1) = (char)psVar3[2];
      if (*(char *)(param_2 + 0x24 + psVar3[1] * 0x94) == 'X') {
        *(undefined1 *)(local_c + 0x42a0c + (int)param_1) = 0x58;
      }
      else {
        *(undefined1 *)(local_c + 0x42a0c + (int)param_1) = 0x30;
      }
      *(short *)(local_c + 0x42a0e + (int)param_1) = psVar3[0x2a6];
      *(short *)(local_c + 0x42a10 + (int)param_1) = psVar3[0x2a7];
      *(short *)(local_c + 0x42a12 + (int)param_1) = psVar3[0x2a8];
      *(short *)(local_c + 0x42a14 + (int)param_1) = psVar3[0x2a9];
      iVar2 = *(int *)(param_2 + 0x2c + psVar3[1] * 0x94);
      if (iVar2 == 3) {
        *(undefined2 *)(local_c + 0x429aa + (int)param_1) = 3;
      }
      else if (iVar2 == 0) {
        *(undefined2 *)(local_c + 0x429aa + (int)param_1) = 0;
      }
      else {
        *(undefined2 *)(local_c + 0x429aa + (int)param_1) = 1;
      }
      param_1[0x214d1] = param_1[0x214d1] + 1;
      local_8 = local_8 + 1;
      psVar3 = psVar3 + 0x2aa;
      local_c = local_c + 0x70;
    } while ((int)local_8 < (int)*param_1);
  }
  return;
}



===== 0x10007520 =====
Function: FUN_10007520 @ 10007520

void __cdecl FUN_10007520(int param_1,int param_2)

{
  short *psVar1;
  byte bVar2;
  char cVar3;
  short sVar4;
  uint uVar5;
  undefined4 uVar6;
  char *pcVar7;
  undefined4 uVar8;
  byte *pbVar9;
  undefined3 extraout_var_00;
  int iVar10;
  byte *pbVar11;
  bool bVar12;
  int local_c;
  int *local_8;
  undefined3 extraout_var;
  
  psVar1 = (short *)(param_1 + 0x429a2);
  local_8 = (int *)0x0;
  if (0 < *psVar1) {
    do {
      if (*(char *)((int)local_8 * 0x70 + 0x29 + (int)psVar1) == '\0') {
        psVar1[(int)local_8 * 0x38 + 2] = 0;
      }
      else {
        psVar1[(int)local_8 * 0x38 + 2] = 1;
      }
      uVar5 = FUN_100091b0(param_1,(int *)&local_8,param_2);
      iVar10 = param_2 + psVar1[(int)local_8 * 0x38 + 3] * 0x94;
      if (((*(char *)(param_2 + 0x30 + psVar1[(int)local_8 * 0x38 + 3] * 0x94) == -1) ||
          (*(char *)((int)psVar1 + (int)local_8 * 0x70 + 0x29) == '\0')) &&
         ((*(char *)(iVar10 + 0x23) == 'U' || (*(char *)(iVar10 + 0x52) == '\0')))) {
        if ((short)uVar5 == 0) {
          uVar6 = FUN_10008dc0(param_1 + 0x429a6,(int)local_8,(int)*psVar1,&local_c);
          if ((short)uVar6 != 0) {
            local_8 = (int *)((int)local_8 + local_c + -1);
            goto LAB_10008340;
          }
          uVar6 = FUN_10009030((byte *)(psVar1 + (int)local_8 * 0x38 + 2),
                               (char)psVar1[(int)local_8 * 0x38 + 5],
                               *(char *)(param_2 + psVar1[(int)local_8 * 0x38 + 3] * 0x94 + 0x23),
                               *(char *)(param_2 + 0x24 + psVar1[(int)local_8 * 0x38 + 3] * 0x94));
          if (((short)uVar6 == 0) &&
             ((*(char *)(param_2 + 0x24 + psVar1[(int)local_8 * 0x38 + 3] * 0x94) != 'A' ||
              (uVar6 = FUN_10002f10((char *)((int)psVar1 + (int)local_8 * 0x70 + 0x29),
                                    (char *)((int)psVar1 + (int)local_8 * 0x70 + 0xb),
                                    (byte *)(psVar1 + (int)local_8 * 0x38 + 2)), (short)uVar6 == 0))
             )) {
            uVar6 = FUN_100086c0((byte *)((int)local_8 * 0x70 + 0xb + (int)psVar1),psVar1,
                                 (char *)local_8,(char *)0x0,0);
            if ((short)uVar6 == 0) {
              if (*(char *)((int)psVar1 + (int)local_8 * 0x70 + 0x29) == '\0') {
                if (((int)local_8 < 1) ||
                   (*(int *)(param_2 + 0x2c + psVar1[(int)local_8 * 0x38 + -0x35] * 0x94) != 0xb)) {
                  uVar6 = 0;
                }
                else {
                  uVar6 = 1;
                }
                if (psVar1[(int)local_8 * 0x38 + 0x39] != 0) {
                  pcVar7 = _strchr((char *)((int)psVar1 + (int)local_8 * 0x70 + 0xb),0x2d);
                  if (pcVar7 == (char *)0x0) {
                    uVar8 = FUN_10002f10((char *)((int)psVar1 + (int)local_8 * 0x70 + 0x29),
                                         (char *)((int)psVar1 + (int)local_8 * 0x70 + 0xb),
                                         (byte *)(psVar1 + (int)local_8 * 0x38 + 2));
                    sVar4 = (short)uVar8;
                  }
                  else {
                    bVar12 = FUN_10003110((char *)((int)psVar1 + (int)local_8 * 0x70 + 0x29),
                                          (char *)((int)psVar1 + (int)local_8 * 0x70 + 0xb),
                                          (byte *)(psVar1 + (int)local_8 * 0x38 + 2));
                    sVar4 = (short)CONCAT31(extraout_var,bVar12);
                  }
                  if (sVar4 != 0) {
                    *(byte *)(psVar1 + (int)local_8 * 0x38 + 2) =
                         *(byte *)(psVar1 + (int)local_8 * 0x38 + 2) | 0x40;
                    goto LAB_1000792a;
                  }
                }
                uVar8 = FUN_10009dc0((char *)((int)psVar1 + (int)local_8 * 0x70 + 0x29),
                                     (byte *)((int)psVar1 + (int)local_8 * 0x70 + 0xb),(short)uVar6,
                                     (byte *)(psVar1 + (int)local_8 * 0x38 + 2),param_1,
                                     (char *)local_8);
                if (((((short)uVar8 == 0) &&
                     (uVar8 = FUN_1000a140((char *)((int)psVar1 + (int)local_8 * 0x70 + 0x29),
                                           (char *)((int)psVar1 + (int)local_8 * 0x70 + 0xb),
                                           (byte *)(psVar1 + (int)local_8 * 0x38 + 2),param_1),
                     (short)uVar8 == 0)) &&
                    (uVar8 = FUN_1000b4d0((char *)((int)psVar1 + (int)local_8 * 0x70 + 0x29),
                                          (char *)((int)psVar1 + (int)local_8 * 0x70 + 0xb),
                                          (byte *)(psVar1 + (int)local_8 * 0x38 + 2),param_1),
                    (short)uVar8 == 0)) &&
                   (((uVar6 = FUN_1000c040((char *)((int)psVar1 + (int)local_8 * 0x70 + 0x29),
                                           (char *)((int)psVar1 + (int)local_8 * 0x70 + 0xb),uVar6,
                                           (byte *)(psVar1 + (int)local_8 * 0x38 + 2),param_1,
                                           (char *)local_8), (short)uVar6 == 0 &&
                     (uVar5 = FUN_1000b800(param_1,(int)local_8), (short)uVar5 == 0)) &&
                    ((uVar6 = FUN_1000c3a0((char *)((int)psVar1 + (int)local_8 * 0x70 + 0x29),
                                           (byte *)((int)psVar1 + (int)local_8 * 0x70 + 0xb),
                                           (byte *)(psVar1 + (int)local_8 * 0x38 + 2),param_1,
                                           (char *)local_8), (short)uVar6 == 0 &&
                     ((uVar6 = FUN_1000c710((char *)((int)psVar1 + (int)local_8 * 0x70 + 0x29),
                                            (char *)((int)psVar1 + (int)local_8 * 0x70 + 0xb),
                                            (byte *)(psVar1 + (int)local_8 * 0x38 + 2),param_1),
                      (short)uVar6 == 0 &&
                      (uVar6 = FUN_10002f10((char *)((int)psVar1 + (int)local_8 * 0x70 + 0x29),
                                            (char *)((int)psVar1 + (int)local_8 * 0x70 + 0xb),
                                            (byte *)(psVar1 + (int)local_8 * 0x38 + 2)),
                      (short)uVar6 == 0)))))))) {
                  *(undefined1 *)(psVar1 + (int)local_8 * 0x38 + 5) = 0x45;
                  FUN_10009cd0((int)psVar1 + (int)local_8 * 0x70 + 0x29,
                               (char *)((int)psVar1 + (int)local_8 * 0x70 + 0xb),
                               (byte *)(psVar1 + (int)local_8 * 0x38 + 2));
                }
              }
            }
            else {
              FUN_10009cd0((int)psVar1 + (int)local_8 * 0x70 + 0x29,
                           (char *)((int)psVar1 + (int)local_8 * 0x70 + 0xb),
                           (byte *)(psVar1 + (int)local_8 * 0x38 + 2));
            }
          }
        }
LAB_1000792a:
        iVar10 = FUN_1001c2c0((byte *)((int)local_8 * 0x70 + 0xb + (int)psVar1),&DAT_10077718);
        if ((iVar10 == 0) &&
           (iVar10 = (int)local_8 * 0x70,
           *(char *)(param_2 + 0x24 + psVar1[(int)local_8 * 0x38 + 3] * 0x94) != 'S')) {
          if ((local_8 == (int *)0x0) &&
             (((((&DAT_1007e188)[*(char *)(param_1 + 0x429ad)] & 0x80) != 0 &&
               ('\x1c' < *(char *)(param_1 + 0x429cb))) && (*(char *)(param_1 + 0x429cb) < ' ')))) {
            if ((((*psVar1 < 2) ||
                 (iVar10 = FUN_100560a0(&PTR_DAT_1007844c,(char *)(param_1 + 0x42a1d),DAT_1007847c,
                                        0x49), iVar10 != -1)) ||
                (iVar10 = FUN_100560a0(&PTR_DAT_10078010,
                                       (char *)((int)local_8 * 0x70 + 0x7b + (int)psVar1),
                                       DAT_100780bc,0x49), iVar10 != -1)) ||
               (((&DAT_1007e188)[*(char *)((int)psVar1 + (int)local_8 * 0x70 + 0x7b)] & 0xc0) == 0))
            {
LAB_10007a7d:
              if ((((int)*psVar1 <= (int)((int)local_8 + 1)) ||
                  (psVar1[(int)local_8 * 0x38 + 3] != psVar1[(int)local_8 * 0x38 + 0x3b])) ||
                 (iVar10 = FUN_1001c2c0((byte *)((int)psVar1 + (int)local_8 * 0x70 + 0x7b),
                                        &DAT_100773a0), iVar10 != 0)) goto LAB_10007ae1;
              *(undefined1 *)((int)local_8 * 0x70 + 0x29 + (int)psVar1) = 0x1e;
            }
            else {
              *(undefined1 *)((int)psVar1 + (int)local_8 * 0x70 + 0x29) = 7;
            }
          }
          else {
            if ((((((&DAT_1007e188)[*(char *)((int)psVar1 + iVar10 + 0xb)] & 0x40) == 0) ||
                 (cVar3 = *(char *)((int)psVar1 + iVar10 + 0x29), cVar3 < '\x1d')) ||
                ('\x1f' < cVar3)) ||
               ((((int)*psVar1 <= (int)((int)local_8 + 1) ||
                 (((&DAT_1007e188)[*(char *)((int)psVar1 + iVar10 + 0x7b)] & 0xc0) == 0)) ||
                (((&DAT_1007e188)[*(char *)((int)psVar1 + iVar10 + 0x7b)] & 0x80) == 0))))
            goto LAB_10007a7d;
            *(undefined1 *)((int)psVar1 + iVar10 + 0x29) = 7;
          }
          *(byte *)(psVar1 + (int)local_8 * 0x38 + 2) =
               *(byte *)(psVar1 + (int)local_8 * 0x38 + 2) | 8;
        }
LAB_10007ae1:
        iVar10 = FUN_1001c2c0((byte *)((int)local_8 * 0x70 + 0xb + (int)psVar1),&DAT_10077424);
        if (((iVar10 == 0) && (*(char *)((int)local_8 * 0x70 + 0x2b + (int)psVar1) == '7')) &&
           (uVar6 = FUN_10008580((int)psVar1,(int)local_8), (short)uVar6 != 0)) {
          *(undefined1 *)((int)local_8 * 0x70 + 0x2b + (int)psVar1) = 0x44;
          *(byte *)(psVar1 + (int)local_8 * 0x38 + 2) =
               *(byte *)(psVar1 + (int)local_8 * 0x38 + 2) | 8;
        }
        if ((int)((int)local_8 + 1) < (int)*psVar1) {
          pbVar11 = &DAT_10077868;
          pbVar9 = (byte *)((int)psVar1 + (int)local_8 * 0x70 + 0xb);
          do {
            bVar2 = *pbVar9;
            bVar12 = bVar2 < *pbVar11;
            if (bVar2 != *pbVar11) {
LAB_10007ba2:
              iVar10 = (1 - (uint)bVar12) - (uint)(bVar12 != 0);
              goto LAB_10007ba7;
            }
            if (bVar2 == 0) break;
            bVar2 = pbVar9[1];
            bVar12 = bVar2 < pbVar11[1];
            if (bVar2 != pbVar11[1]) goto LAB_10007ba2;
            pbVar9 = pbVar9 + 2;
            pbVar11 = pbVar11 + 2;
          } while (bVar2 != 0);
          iVar10 = 0;
LAB_10007ba7:
          if ((iVar10 == 0) &&
             ((iVar10 = FUN_1001c2c0((byte *)((int)psVar1 + (int)local_8 * 0x70 + 0x7b),
                                     &DAT_10077718), iVar10 == 0 ||
              (iVar10 = FUN_1001c2c0((byte *)((int)local_8 * 0x70 + 0x7b + (int)psVar1),
                                     &DAT_1007736c), iVar10 == 0)))) {
            *(undefined1 *)((int)local_8 * 0x70 + 0x29 + (int)psVar1) = 8;
            *(undefined1 *)(psVar1 + (int)local_8 * 0x38 + 0x15) = 0x35;
            *(undefined1 *)((int)local_8 * 0x70 + 0x2b + (int)psVar1) = 0;
            *(byte *)(psVar1 + (int)local_8 * 0x38 + 2) =
                 *(byte *)(psVar1 + (int)local_8 * 0x38 + 2) | 8;
          }
        }
        iVar10 = FUN_1001c2c0((byte *)((int)local_8 * 0x70 + 0xb + (int)psVar1),&DAT_10077860);
        if ((iVar10 == 0) && ((char)psVar1[(int)local_8 * 0x38 + 0x15] == '\'')) {
          if (0 < (int)local_8) {
            uVar5 = FUN_100560a0(&PTR_DAT_10078010,
                                 (char *)((int)psVar1 + (int)local_8 * 0x70 + -0x65),DAT_100780bc,
                                 0x49);
            if ((((uVar5 < 0x80000000) ||
                 (((1 < (int)local_8 &&
                   (iVar10 = FUN_100560a0(&PTR_DAT_10078010,
                                          (char *)((int)local_8 * 0x70 + -0xd5 + (int)psVar1),
                                          DAT_100780bc,0x49), -1 < iVar10)) &&
                  (iVar10 = FUN_1001c2c0((byte *)((int)local_8 * 0x70 + -0x65 + (int)psVar1),
                                         &DAT_10077608), iVar10 == 0)))) ||
                ((0 < (int)local_8 &&
                 (uVar5 = FUN_100560a0(&PTR_DAT_100783a4,
                                       (char *)((int)local_8 * 0x70 + -0x65 + (int)psVar1),
                                       DAT_100783cc,0x49), uVar5 < 0x80000000)))) ||
               ((1 < (int)local_8 &&
                ((iVar10 = FUN_100560a0(&PTR_DAT_100783a4,
                                        (char *)((int)local_8 * 0x70 + -0xd5 + (int)psVar1),
                                        DAT_100783cc,0x49), -1 < iVar10 &&
                 (iVar10 = FUN_1001c2c0((byte *)((int)local_8 * 0x70 + -0x65 + (int)psVar1),
                                        &DAT_10077608), iVar10 == 0)))))) {
LAB_1000801b:
              *(undefined1 *)(psVar1 + (int)local_8 * 0x38 + 0x15) = 0x18;
            }
            else {
              if ((int)local_8 < 1) goto LAB_10007ee8;
              iVar10 = FUN_1001c2c0((byte *)((int)local_8 * 0x70 + -0x65 + (int)psVar1),
                                    &DAT_1007785c);
              if ((iVar10 == 0) ||
                 (iVar10 = FUN_1001c2c0((byte *)((int)local_8 * 0x70 + -0x65 + (int)psVar1),
                                        &DAT_10077858), iVar10 == 0)) goto LAB_1000801b;
              if ((int)local_8 < 1) {
LAB_10007ee8:
                if (((((int)local_8 < 3) ||
                     (((iVar10 = FUN_1001c2c0((byte *)((int)local_8 * 0x70 + -0x145 + (int)psVar1),
                                              &DAT_10077850), iVar10 != 0 &&
                       (iVar10 = FUN_1001c2c0((byte *)((int)local_8 * 0x70 + -0x145 + (int)psVar1),
                                              &DAT_1007784c), iVar10 != 0)) &&
                      (iVar10 = FUN_1001c2c0((byte *)((int)local_8 * 0x70 + -0x145 + (int)psVar1),
                                             &DAT_10077844), iVar10 != 0)))) &&
                    (((int)local_8 < 2 ||
                     (((iVar10 = FUN_1001c2c0((byte *)((int)local_8 * 0x70 + -0xd5 + (int)psVar1),
                                              &DAT_10077850), iVar10 != 0 &&
                       (iVar10 = FUN_1001c2c0((byte *)((int)local_8 * 0x70 + -0xd5 + (int)psVar1),
                                              &DAT_1007784c), iVar10 != 0)) &&
                      (iVar10 = FUN_1001c2c0((byte *)((int)local_8 * 0x70 + -0xd5 + (int)psVar1),
                                             &DAT_10077844), iVar10 != 0)))))) ||
                   (iVar10 = FUN_1001c2c0((byte *)((int)local_8 * 0x70 + -0x65 + (int)psVar1),
                                          (byte *)&DAT_10077380), iVar10 == 0)) goto LAB_10008044;
                goto LAB_1000801b;
              }
              iVar10 = FUN_1001c2c0(&DAT_10077854,
                                    (byte *)((int)local_8 * 0x70 + -0x65 + (int)psVar1));
              if (((iVar10 != 0) &&
                  (iVar10 = FUN_100560a0(&PTR_DAT_1007868c,
                                         (char *)((int)local_8 * 0x70 + -0x65 + (int)psVar1),
                                         DAT_100786ec,0x49), iVar10 < 0)) ||
                 (iVar10 = FUN_1000cc60((undefined *)psVar1,local_8,0x42,0x14,0x1007837c,
                                        DAT_100783a0,0x49), iVar10 < 0)) {
                if ((0 < (int)local_8) &&
                   (uVar5 = FUN_100560a0(&PTR_DAT_1007871c,
                                         (char *)((int)local_8 * 0x70 + -0x65 + (int)psVar1),
                                         DAT_10078764,0x49), uVar5 < 0x80000000)) goto LAB_1000801b;
                goto LAB_10007ee8;
              }
              *(undefined1 *)(psVar1 + (int)local_8 * 0x38 + 0x15) = 0x18;
            }
            *(byte *)(psVar1 + (int)local_8 * 0x38 + 2) =
                 *(byte *)(psVar1 + (int)local_8 * 0x38 + 2) | 8;
            goto LAB_10008044;
          }
        }
        else {
LAB_10008044:
          if (0 < (int)local_8) {
            bVar12 = FUN_10008550((char *)((int)local_8 * 0x70 + -0x65 + (int)psVar1));
            if (((short)CONCAT31(extraout_var_00,bVar12) != 0) &&
               (iVar10 = FUN_1001c2c0((byte *)((int)local_8 * 0x70 + 0xb + (int)psVar1),
                                      (byte *)s_minute_100774e8), iVar10 == 0)) {
              *(undefined1 *)((int)local_8 * 0x70 + 0x29 + (int)psVar1) = 0x2c;
              *(undefined1 *)(psVar1 + (int)local_8 * 0x38 + 0x15) = 0x24;
              *(undefined1 *)((int)local_8 * 0x70 + 0x2b + (int)psVar1) = 0x2d;
              *(undefined1 *)(psVar1 + (int)local_8 * 0x38 + 0x16) = 0x23;
              *(undefined1 *)((int)local_8 * 0x70 + 0x2d + (int)psVar1) = 0x39;
              *(undefined1 *)(psVar1 + (int)local_8 * 0x38 + 0x17) = 0;
              *(byte *)(psVar1 + (int)local_8 * 0x38 + 2) =
                   *(byte *)(psVar1 + (int)local_8 * 0x38 + 2) | 8;
            }
            if (((0 < (int)local_8) &&
                (iVar10 = FUN_1001c2c0((byte *)((int)local_8 * 0x70 + 0xb + (int)psVar1),
                                       (byte *)s_close_100776ec), iVar10 == 0)) &&
               (((char)psVar1[(int)local_8 * 0x38 + 0x16] == '7' &&
                (uVar6 = FUN_10008580((int)psVar1,(int)local_8), (short)uVar6 != 0)))) {
              *(undefined1 *)(psVar1 + (int)local_8 * 0x38 + 0x16) = 0x44;
              *(byte *)(psVar1 + (int)local_8 * 0x38 + 2) =
                   *(byte *)(psVar1 + (int)local_8 * 0x38 + 2) | 8;
            }
          }
        }
        iVar10 = FUN_1001c2c0((byte *)((int)local_8 * 0x70 + 0xb + (int)psVar1),&DAT_10077840);
        if (iVar10 == 0) {
          iVar10 = FUN_1000cc60((undefined *)psVar1,local_8,0x42,5,0x10078334,DAT_1007834c,0x49);
          if (iVar10 < 0) {
            *(undefined1 *)(psVar1 + (int)local_8 * 0x38 + 0x15) = 0xe;
          }
          else {
            *(undefined1 *)(psVar1 + (int)local_8 * 0x38 + 0x15) = 0x30;
          }
          *(byte *)(psVar1 + (int)local_8 * 0x38 + 2) =
               *(byte *)(psVar1 + (int)local_8 * 0x38 + 2) | 8;
        }
        iVar10 = FUN_1001c2c0((byte *)((int)local_8 * 0x70 + 0xb + (int)psVar1),&DAT_10077838);
        if (((iVar10 == 0) && ((char)psVar1[(int)local_8 * 0x38 + 0x15] == '\'')) &&
           (iVar10 = FUN_1000cc60((undefined *)psVar1,local_8,0x42,10,0x10078350,DAT_10078378,0x49),
           -1 < iVar10)) {
          *(undefined1 *)(psVar1 + (int)local_8 * 0x38 + 0x15) = 0x18;
          *(byte *)(psVar1 + (int)local_8 * 0x38 + 2) =
               *(byte *)(psVar1 + (int)local_8 * 0x38 + 2) | 8;
        }
        iVar10 = FUN_1001c2c0((byte *)((int)local_8 * 0x70 + 0xb + (int)psVar1),
                              (byte *)s_mouth_100774d8);
        if (((iVar10 == 0) && (*(char *)((int)local_8 * 0x70 + 0x2b + (int)psVar1) == '\x16')) &&
           ((1 < (int)local_8 &&
            (iVar10 = FUN_1000cc60((undefined *)psVar1,local_8,0x52,2,0x10078480,DAT_100784c4,0x49),
            -1 < iVar10)))) {
          *(undefined1 *)((int)local_8 * 0x70 + 0x2b + (int)psVar1) = 0x3a;
          *(byte *)(psVar1 + (int)local_8 * 0x38 + 2) =
               *(byte *)(psVar1 + (int)local_8 * 0x38 + 2) | 8;
        }
      }
LAB_10008340:
      local_8 = (int *)((int)local_8 + 1);
    } while ((int)local_8 < (int)*psVar1);
  }
  local_8 = (int *)0x0;
  if (0 < *psVar1) {
    do {
      if (((0 < (int)local_8) &&
          (cVar3 = *(char *)((int)psVar1 + (int)local_8 * 0x70 + 0x29), '\0' < cVar3)) &&
         (((cVar3 < 'F' && (*(short *)(&DAT_10077d94 + cVar3 * 2) != 0)) &&
          (((iVar10 = FUN_1001c2c0((byte *)((int)psVar1 + (int)local_8 * 0x70 + -0x65),&DAT_1007736c
                                  ), iVar10 == 0 ||
            (iVar10 = FUN_1001c2c0((byte *)((int)local_8 * 0x70 + -0x65 + (int)psVar1),&DAT_1007761c
                                  ), iVar10 == 0)) &&
           (*(char *)((int)psVar1 + (int)local_8 * 0x70 + 0x29) != 'C')))))) {
        *(undefined1 *)(psVar1 + (int)local_8 * 0x38 + -0x23) = 0x26;
        *(undefined1 *)((int)local_8 * 0x70 + -0x45 + (int)psVar1) = 0;
        *(byte *)(psVar1 + (int)local_8 * 0x38 + -0x36) =
             *(byte *)(psVar1 + (int)local_8 * 0x38 + -0x36) | 8;
      }
      local_8 = (int *)((int)local_8 + 1);
    } while ((int)local_8 < (int)*psVar1);
  }
  return;
}



===== 0x1000dfc0 =====
Function: FUN_1000dfc0 @ 1000dfc0

undefined4 __cdecl FUN_1000dfc0(byte *param_1,int param_2,int param_3)

{
  bool bVar1;
  undefined1 uVar2;
  undefined3 extraout_var;
  undefined4 uVar3;
  void *this;
  int *piVar4;
  int *piVar5;
  int iVar6;
  byte local_c [4];
  int *local_8;
  
  piVar4 = (int *)0x0;
  local_8 = (int *)0xffffffff;
  if (0 < param_2) {
    do {
      bVar1 = FUN_1000e160((char *)local_c,param_1,param_2,piVar4,(int *)&local_8,param_3);
      if ((short)CONCAT31(extraout_var,bVar1) == 0) {
        uVar3 = FUN_1003a7b0((char *)local_c,(char *)(param_1 + (int)piVar4 * 0x94 + 0x34),0x47,0x49
                             ,param_3);
        if ((short)uVar3 == 0) {
          iVar6 = FUN_100035d0(DAT_100fee08,param_1 + (int)piVar4 * 0x94 + 0x34);
          if (-1 < iVar6) {
            param_1[(int)piVar4 * 0x94 + 0x26] = (char)iVar6 + 1;
          }
        }
        else {
          iVar6 = 0;
          piVar5 = piVar4;
          uVar2 = FUN_10064645(this,local_c);
          uVar3 = FUN_1000e0c0((int)param_1,uVar2,(int)piVar5,iVar6);
          if ((short)uVar3 == 0) {
            return 0xffffffff;
          }
        }
      }
      else {
        piVar5 = local_8;
        uVar2 = FUN_10064645(local_c,local_c);
        uVar3 = FUN_1000e0c0((int)param_1,uVar2,(int)piVar4,(int)piVar5);
        piVar4 = local_8;
        if ((short)uVar3 == 0) {
          return 0xffffffff;
        }
        local_8 = (int *)0xffffffff;
      }
      piVar4 = (int *)((int)piVar4 + 1);
    } while ((int)piVar4 < param_2);
  }
  return 1;
}



