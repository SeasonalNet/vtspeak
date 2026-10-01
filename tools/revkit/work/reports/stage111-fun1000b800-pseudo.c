===== 0x1000b800 =====
Function: FUN_1000b800 @ 1000b800

uint __cdecl FUN_1000b800(int param_1,int param_2)

{
  char cVar1;
  byte bVar2;
  short sVar3;
  undefined1 uVar4;
  undefined1 uVar5;
  uint in_EAX;
  char *pcVar6;
  short *psVar7;
  int iVar8;
  byte *pbVar9;
  uint uVar10;
  uint uVar11;
  byte *pbVar12;
  int iVar13;
  int iVar14;
  char *pcVar15;
  char *pcVar16;
  byte *pbVar17;
  bool bVar18;
  char local_7c [68];
  char local_38 [32];
  char *local_18;
  char *local_14;
  byte *local_10;
  uint local_c;
  char local_5;
  
  iVar13 = param_1;
  if (param_2 < 1) goto LAB_1000c031;
  pcVar15 = (char *)(param_2 * 0x70);
  sVar3 = *(short *)((int)pcVar15 + 0x42938 + param_1);
  in_EAX = CONCAT22((short)(in_EAX >> 0x10),sVar3);
  if (sVar3 != *(short *)((int)pcVar15 + 0x429a8 + param_1)) goto LAB_1000c031;
  local_10 = (byte *)((int)pcVar15 + 0x429ad + param_1);
  local_18 = pcVar15;
  pcVar6 = _strchr(&DAT_10077988,(int)*(char *)((int)pcVar15 + 0x429ad + param_1));
  if ((pcVar6 != (char *)0x0) &&
     (cVar1 = *(char *)((int)pcVar15 + 0x429ae + param_1),
     in_EAX = CONCAT31((int3)((uint)pcVar6 >> 8),cVar1), cVar1 == '\0')) goto LAB_1000c031;
  iVar14 = param_2 + -1;
  if (0 < iVar14) {
    psVar7 = (short *)(iVar14 * 0x70 + 0x429a8 + param_1);
    do {
      if (psVar7[-0x38] != *psVar7) break;
      iVar14 = iVar14 + -1;
      psVar7 = psVar7 + -0x38;
    } while (0 < iVar14);
  }
  local_38[0] = '\0';
  local_c = iVar14;
  if (iVar14 <= param_2) {
    local_14 = (char *)(iVar14 * 0x70 + 0x429ad + param_1);
    do {
      if (local_c == iVar14) {
        uVar10 = 0xffffffff;
        pcVar15 = local_14;
        do {
          pcVar6 = pcVar15;
          if (uVar10 == 0) break;
          uVar10 = uVar10 - 1;
          pcVar6 = pcVar15 + 1;
          cVar1 = *pcVar15;
          pcVar15 = pcVar6;
        } while (cVar1 != '\0');
        uVar10 = ~uVar10;
        pcVar15 = pcVar6 + -uVar10;
        pcVar6 = local_38;
        for (uVar11 = uVar10 >> 2; uVar11 != 0; uVar11 = uVar11 - 1) {
          *(undefined4 *)pcVar6 = *(undefined4 *)pcVar15;
          pcVar15 = pcVar15 + 4;
          pcVar6 = pcVar6 + 4;
        }
        for (uVar10 = uVar10 & 3; uVar10 != 0; uVar10 = uVar10 - 1) {
          *pcVar6 = *pcVar15;
          pcVar15 = pcVar15 + 1;
          pcVar6 = pcVar6 + 1;
        }
      }
      else {
        uVar10 = 0xffffffff;
        pcVar15 = local_14;
        do {
          pcVar6 = pcVar15;
          if (uVar10 == 0) break;
          uVar10 = uVar10 - 1;
          pcVar6 = pcVar15 + 1;
          cVar1 = *pcVar15;
          pcVar15 = pcVar6;
        } while (cVar1 != '\0');
        uVar10 = ~uVar10;
        iVar8 = -1;
        pcVar15 = local_38;
        do {
          pcVar16 = pcVar15;
          if (iVar8 == 0) break;
          iVar8 = iVar8 + -1;
          pcVar16 = pcVar15 + 1;
          cVar1 = *pcVar15;
          pcVar15 = pcVar16;
        } while (cVar1 != '\0');
        pcVar15 = pcVar6 + -uVar10;
        pcVar6 = pcVar16 + -1;
        for (uVar11 = uVar10 >> 2; uVar11 != 0; uVar11 = uVar11 - 1) {
          *(undefined4 *)pcVar6 = *(undefined4 *)pcVar15;
          pcVar15 = pcVar15 + 4;
          pcVar6 = pcVar6 + 4;
        }
        for (uVar10 = uVar10 & 3; uVar10 != 0; uVar10 = uVar10 - 1) {
          *pcVar6 = *pcVar15;
          pcVar15 = pcVar15 + 1;
          pcVar6 = pcVar6 + 1;
        }
      }
      local_c = local_c + 1;
      local_14 = local_14 + 0x70;
      pcVar15 = local_18;
    } while ((int)local_c <= param_2);
  }
  iVar8 = FUN_1000cb30(local_38,local_7c,param_1 + 0x48128);
  if (iVar8 == 1) {
    iVar13 = iVar14 * 0x70 + 0x429a2 + param_1;
    uVar10 = 0xffffffff;
    pcVar15 = local_7c;
    do {
      pcVar6 = pcVar15;
      if (uVar10 == 0) break;
      uVar10 = uVar10 - 1;
      pcVar6 = pcVar15 + 1;
      cVar1 = *pcVar15;
      pcVar15 = pcVar6;
    } while (cVar1 != '\0');
    pbVar12 = (byte *)~uVar10;
    local_18 = (char *)(iVar13 + 0x29);
    pcVar15 = pcVar6 + -(int)pbVar12;
    for (uVar10 = (uint)pbVar12 >> 2; uVar10 != 0; uVar10 = uVar10 - 1) {
      *(undefined4 *)local_18 = *(undefined4 *)pcVar15;
      pcVar15 = pcVar15 + 4;
      local_18 = local_18 + 4;
    }
    for (uVar10 = (uint)pbVar12 & 3; uVar10 != 0; uVar10 = uVar10 - 1) {
      *local_18 = *pcVar15;
      pcVar15 = pcVar15 + 1;
      local_18 = local_18 + 1;
    }
    pbVar9 = (byte *)(iVar13 + 4);
    *pbVar9 = *pbVar9 | 8;
    iVar14 = iVar14 + 1;
    if (iVar14 <= param_2) {
      iVar13 = (param_2 - iVar14) + 1;
      pbVar12 = (byte *)(iVar14 * 0x70 + 0x429a6 + param_1);
      do {
        pbVar12[0x25] = 0;
        *pbVar12 = *pbVar12 | 8;
        pbVar12 = pbVar12 + 0x70;
        iVar13 = iVar13 + -1;
      } while (iVar13 != 0);
    }
    return CONCAT22((short)((uint)pbVar12 >> 0x10),1);
  }
  uVar10 = 0xffffffff;
  local_c = 0;
  pcVar6 = (char *)((int)pcVar15 + 0x4295b + param_1);
  do {
    if (uVar10 == 0) break;
    uVar10 = uVar10 - 1;
    cVar1 = *pcVar6;
    pcVar6 = pcVar6 + 1;
  } while (cVar1 != '\0');
  iVar14 = ~uVar10 - 1;
  if (iVar14 < 1) {
    param_1._3_1_ = '\0';
  }
  else {
    param_1._3_1_ = *(char *)((int)pcVar15 + iVar14 + 0x4295a + param_1);
  }
  if (iVar14 < 2) {
    local_5 = '\0';
  }
  else {
    local_5 = *(char *)(iVar14 + (int)pcVar15 + 0x42959 + iVar13);
  }
  iVar14 = FUN_1001c2c0(local_10,&DAT_100773a0);
  pbVar12 = local_10;
  if (iVar14 == 0) {
    uVar4 = 0x44;
    uVar5 = uVar4;
    if (param_1._3_1_ < '\x01') {
LAB_1000bd81:
      *(undefined1 *)((int)pcVar15 + 0x429cb + iVar13) = uVar5;
      local_c = 1;
    }
    else if ((((param_1._3_1_ == ' ') || (param_1._3_1_ == '*')) || (param_1._3_1_ == '5')) ||
            ((param_1._3_1_ == '9' || (param_1._3_1_ == ':')))) {
      *(undefined1 *)((int)pcVar15 + 0x429cb + iVar13) = 0x37;
      local_c = 1;
    }
    else {
      if (((param_1._3_1_ != '\x14') && (param_1._3_1_ != ')')) &&
         ((param_1._3_1_ != '7' &&
          (((param_1._3_1_ != '8' && (param_1._3_1_ != 'D')) && (param_1._3_1_ != 'E'))))))
      goto LAB_1000bd81;
LAB_1000ba88:
      *(undefined1 *)((int)pcVar15 + 0x429cb + iVar13) = 0x23;
      *(undefined1 *)((int)pcVar15 + 0x429cc + iVar13) = uVar4;
      local_c = 2;
    }
  }
  else {
    pbVar17 = &DAT_10077568;
    pbVar9 = local_10;
    do {
      bVar2 = *pbVar9;
      bVar18 = bVar2 < *pbVar17;
      if (bVar2 != *pbVar17) {
LAB_1000bacf:
        iVar14 = (1 - (uint)bVar18) - (uint)(bVar18 != 0);
        goto LAB_1000bad4;
      }
      if (bVar2 == 0) break;
      bVar2 = pbVar9[1];
      bVar18 = bVar2 < pbVar17[1];
      if (bVar2 != pbVar17[1]) goto LAB_1000bacf;
      pbVar9 = pbVar9 + 2;
      pbVar17 = pbVar17 + 2;
    } while (bVar2 != 0);
    iVar14 = 0;
LAB_1000bad4:
    if (iVar14 != 0) {
      pbVar17 = &DAT_10077984;
      pbVar9 = local_10;
      do {
        bVar2 = *pbVar9;
        bVar18 = bVar2 < *pbVar17;
        if (bVar2 != *pbVar17) {
LAB_1000bb08:
          iVar14 = (1 - (uint)bVar18) - (uint)(bVar18 != 0);
          goto LAB_1000bb0d;
        }
        if (bVar2 == 0) break;
        bVar2 = pbVar9[1];
        bVar18 = bVar2 < pbVar17[1];
        if (bVar2 != pbVar17[1]) goto LAB_1000bb08;
        pbVar9 = pbVar9 + 2;
        pbVar17 = pbVar17 + 2;
      } while (bVar2 != 0);
      iVar14 = 0;
LAB_1000bb0d:
      if (iVar14 != 0) {
        pbVar17 = &DAT_100778cc;
        pbVar9 = local_10;
        do {
          bVar2 = *pbVar9;
          bVar18 = bVar2 < *pbVar17;
          if (bVar2 != *pbVar17) {
LAB_1000bb41:
            iVar14 = (1 - (uint)bVar18) - (uint)(bVar18 != 0);
            goto LAB_1000bb46;
          }
          if (bVar2 == 0) break;
          bVar2 = pbVar9[1];
          bVar18 = bVar2 < pbVar17[1];
          if (bVar2 != pbVar17[1]) goto LAB_1000bb41;
          pbVar9 = pbVar9 + 2;
          pbVar17 = pbVar17 + 2;
        } while (bVar2 != 0);
        iVar14 = 0;
LAB_1000bb46:
        if (iVar14 != 0) {
          pbVar17 = &DAT_100778c8;
          pbVar9 = local_10;
          do {
            bVar2 = *pbVar9;
            bVar18 = bVar2 < *pbVar17;
            if (bVar2 != *pbVar17) {
LAB_1000bb7a:
              iVar14 = (1 - (uint)bVar18) - (uint)(bVar18 != 0);
              goto LAB_1000bb7f;
            }
            if (bVar2 == 0) break;
            bVar2 = pbVar9[1];
            bVar18 = bVar2 < pbVar17[1];
            if (bVar2 != pbVar17[1]) goto LAB_1000bb7a;
            pbVar9 = pbVar9 + 2;
            pbVar17 = pbVar17 + 2;
          } while (bVar2 != 0);
          iVar14 = 0;
LAB_1000bb7f:
          if (iVar14 != 0) {
            pbVar17 = &DAT_100778c4;
            pbVar9 = local_10;
            do {
              bVar2 = *pbVar9;
              bVar18 = bVar2 < *pbVar17;
              if (bVar2 != *pbVar17) {
LAB_1000bbb3:
                iVar14 = (1 - (uint)bVar18) - (uint)(bVar18 != 0);
                goto LAB_1000bbb8;
              }
              if (bVar2 == 0) break;
              bVar2 = pbVar9[1];
              bVar18 = bVar2 < pbVar17[1];
              if (bVar2 != pbVar17[1]) goto LAB_1000bbb3;
              pbVar9 = pbVar9 + 2;
              pbVar17 = pbVar17 + 2;
            } while (bVar2 != 0);
            iVar14 = 0;
LAB_1000bbb8:
            if (iVar14 != 0) {
              pbVar17 = &DAT_100778c0;
              pbVar9 = local_10;
              do {
                bVar2 = *pbVar9;
                bVar18 = bVar2 < *pbVar17;
                if (bVar2 != *pbVar17) {
LAB_1000bbec:
                  iVar14 = (1 - (uint)bVar18) - (uint)(bVar18 != 0);
                  goto LAB_1000bbf1;
                }
                if (bVar2 == 0) break;
                bVar2 = pbVar9[1];
                bVar18 = bVar2 < pbVar17[1];
                if (bVar2 != pbVar17[1]) goto LAB_1000bbec;
                pbVar9 = pbVar9 + 2;
                pbVar17 = pbVar17 + 2;
              } while (bVar2 != 0);
              iVar14 = 0;
LAB_1000bbf1:
              if (iVar14 != 0) {
                pbVar17 = &DAT_100778bc;
                pbVar9 = local_10;
                do {
                  bVar2 = *pbVar9;
                  bVar18 = bVar2 < *pbVar17;
                  if (bVar2 != *pbVar17) {
LAB_1000bc25:
                    iVar14 = (1 - (uint)bVar18) - (uint)(bVar18 != 0);
                    goto LAB_1000bc2a;
                  }
                  if (bVar2 == 0) break;
                  bVar2 = pbVar9[1];
                  bVar18 = bVar2 < pbVar17[1];
                  if (bVar2 != pbVar17[1]) goto LAB_1000bc25;
                  pbVar9 = pbVar9 + 2;
                  pbVar17 = pbVar17 + 2;
                } while (bVar2 != 0);
                iVar14 = 0;
LAB_1000bc2a:
                if (iVar14 != 0) {
                  pbVar17 = &DAT_100778b8;
                  pbVar9 = local_10;
                  do {
                    bVar2 = *pbVar9;
                    bVar18 = bVar2 < *pbVar17;
                    if (bVar2 != *pbVar17) {
LAB_1000bc5e:
                      iVar14 = (1 - (uint)bVar18) - (uint)(bVar18 != 0);
                      goto LAB_1000bc63;
                    }
                    if (bVar2 == 0) break;
                    bVar2 = pbVar9[1];
                    bVar18 = bVar2 < pbVar17[1];
                    if (bVar2 != pbVar17[1]) goto LAB_1000bc5e;
                    pbVar9 = pbVar9 + 2;
                    pbVar17 = pbVar17 + 2;
                  } while (bVar2 != 0);
                  iVar14 = 0;
LAB_1000bc63:
                  if (iVar14 != 0) {
                    pbVar17 = &DAT_100778b4;
                    pbVar9 = local_10;
                    do {
                      bVar2 = *pbVar9;
                      bVar18 = bVar2 < *pbVar17;
                      if (bVar2 != *pbVar17) {
LAB_1000bc97:
                        iVar14 = (1 - (uint)bVar18) - (uint)(bVar18 != 0);
                        goto LAB_1000bc9c;
                      }
                      if (bVar2 == 0) break;
                      bVar2 = pbVar9[1];
                      bVar18 = bVar2 < pbVar17[1];
                      if (bVar2 != pbVar17[1]) goto LAB_1000bc97;
                      pbVar9 = pbVar9 + 2;
                      pbVar17 = pbVar17 + 2;
                    } while (bVar2 != 0);
                    iVar14 = 0;
LAB_1000bc9c:
                    if (iVar14 != 0) {
                      pbVar17 = &DAT_100778b0;
                      pbVar9 = local_10;
                      do {
                        bVar2 = *pbVar9;
                        bVar18 = bVar2 < *pbVar17;
                        if (bVar2 != *pbVar17) {
LAB_1000bcd0:
                          iVar14 = (1 - (uint)bVar18) - (uint)(bVar18 != 0);
                          goto LAB_1000bcd5;
                        }
                        if (bVar2 == 0) break;
                        bVar2 = pbVar9[1];
                        bVar18 = bVar2 < pbVar17[1];
                        if (bVar2 != pbVar17[1]) goto LAB_1000bcd0;
                        pbVar9 = pbVar9 + 2;
                        pbVar17 = pbVar17 + 2;
                      } while (bVar2 != 0);
                      iVar14 = 0;
LAB_1000bcd5:
                      if (iVar14 != 0) {
                        iVar14 = FUN_1001c2c0(local_10,&DAT_100778ac);
                        if (iVar14 == 0) {
                          if (0 < param_2) {
                            pcVar6 = (char *)((int)pcVar15 + 0x4293d + iVar13);
                            iVar14 = FUN_100560a0(&PTR_s_anybody_100786f0,pcVar6,DAT_10078718,0x49);
                            if ((iVar14 == -1) &&
                               (iVar14 = FUN_100560a0(&PTR_DAT_1007871c,pcVar6,DAT_10078764,0x49),
                               iVar14 == -1)) {
                              if (((((param_1._3_1_ != '*') &&
                                    ((param_1._3_1_ != '5' && (param_1._3_1_ != '7')))) &&
                                   (param_1._3_1_ != ' ')) &&
                                  (((param_1._3_1_ != '\"' && (param_1._3_1_ != '8')) &&
                                   (param_1._3_1_ != ':')))) && (param_1._3_1_ != '\x14')) {
                                uVar5 = 0x15;
                                uVar4 = 0x15;
                                if ((param_1._3_1_ == '9') || (param_1._3_1_ == '\x15'))
                                goto LAB_1000ba88;
                                goto LAB_1000bd81;
                              }
                              *(undefined1 *)((int)pcVar15 + 0x429cb + iVar13) = 0x39;
                              local_c = 1;
                              goto LAB_1000c005;
                            }
                          }
                          if ((((param_1._3_1_ < '\x01') || (param_1._3_1_ < '\x01')) ||
                              (('E' < param_1._3_1_ ||
                               (*(short *)(&DAT_10077d94 + param_1._3_1_ * 2) == 0)))) &&
                             ((((local_5 < '\x01' || (local_5 < '\x01')) || ('E' < local_5)) ||
                              ((*(short *)(&DAT_10077d94 + local_5 * 2) == 0 ||
                               (param_1._3_1_ != '6')))))) {
                            *(undefined1 *)((int)pcVar15 + 0x429cb + iVar13) = 7;
                            *(undefined1 *)((int)pcVar15 + 0x429cc + iVar13) = 0x15;
                            local_c = 2;
                          }
                          else {
                            *(undefined1 *)((int)pcVar15 + 0x429cb + iVar13) = 0x15;
                            local_c = 1;
                          }
                        }
                        else {
                          iVar14 = FUN_1001c2c0(pbVar12,&DAT_100778a8);
                          if (iVar14 == 0) {
                            if ((((param_1._3_1_ < '\x01') || (param_1._3_1_ < '\x01')) ||
                                ('E' < param_1._3_1_)) ||
                               (*(short *)(&DAT_10077d94 + param_1._3_1_ * 2) == 0))
                            goto LAB_1000bf26;
                            *(undefined1 *)((int)pcVar15 + 0x429cb + iVar13) = 0x2c;
                            local_c = 1;
                          }
                        }
                        goto LAB_1000c005;
                      }
                    }
                    if ((((param_1._3_1_ < '\x01') || (param_1._3_1_ < '\x01')) ||
                        (('E' < param_1._3_1_ ||
                         (*(short *)(&DAT_10077d94 + param_1._3_1_ * 2) == 0)))) &&
                       ((((local_5 < '\x01' || (local_5 < '\x01')) || ('E' < local_5)) ||
                        ((*(short *)(&DAT_10077d94 + local_5 * 2) == 0 || (param_1._3_1_ != '6')))))
                       ) {
                      *(undefined1 *)((int)pcVar15 + 0x429cb + iVar13) = 7;
                      *(undefined1 *)((int)pcVar15 + 0x429cc + iVar13) = 0x2b;
                      local_c = 2;
                    }
                    else {
                      *(undefined1 *)((int)pcVar15 + 0x429cb + iVar13) = 0x2b;
                      local_c = 1;
                    }
                    goto LAB_1000c005;
                  }
                }
                if (((param_1._3_1_ < '\x01') || (param_1._3_1_ < '\x01')) ||
                   (('E' < param_1._3_1_ || (*(short *)(&DAT_10077d94 + param_1._3_1_ * 2) == 0))))
                {
                  *(undefined1 *)((int)pcVar15 + 0x429cb + iVar13) = 0x1a;
                  local_c = 1;
                }
                else {
                  *(undefined1 *)((int)pcVar15 + 0x429cb + iVar13) = 0x36;
                  local_c = 1;
                }
                goto LAB_1000c005;
              }
            }
LAB_1000bf26:
            *(undefined1 *)((int)pcVar15 + 0x429cb + iVar13) = 7;
            *(undefined1 *)((int)pcVar15 + 0x429cc + iVar13) = 0x2c;
            local_c = 2;
            goto LAB_1000c005;
          }
        }
        if ((((param_1._3_1_ < '\x01') || (param_1._3_1_ < '\x01')) || ('E' < param_1._3_1_)) ||
           (*(short *)(&DAT_10077d94 + param_1._3_1_ * 2) == 0)) {
          *(undefined1 *)((int)pcVar15 + 0x429cb + iVar13) = 7;
          *(undefined1 *)((int)pcVar15 + 0x429cc + iVar13) = 0x41;
          local_c = 2;
        }
        else {
          *(undefined1 *)((int)pcVar15 + 0x429cb + iVar13) = 0x41;
          local_c = 1;
        }
        goto LAB_1000c005;
      }
    }
    if ((((param_1._3_1_ < '\x01') || (param_1._3_1_ < '\x01')) ||
        (('E' < param_1._3_1_ || (*(short *)(&DAT_10077d94 + param_1._3_1_ * 2) == 0)))) &&
       ((((local_5 < '\x01' || (local_5 < '\x01')) || ('E' < local_5)) ||
        ((*(short *)(&DAT_10077d94 + local_5 * 2) == 0 || (param_1._3_1_ != '6')))))) {
      *(undefined1 *)((int)pcVar15 + 0x429cb + iVar13) = 7;
      *(undefined1 *)((int)pcVar15 + 0x429cc + iVar13) = 0x2d;
      *(undefined1 *)((int)pcVar15 + 0x429cd + iVar13) = 0x39;
      local_c = 3;
    }
    else {
      *(undefined1 *)((int)pcVar15 + 0x429cb + iVar13) = 0x2d;
      *(undefined1 *)((int)pcVar15 + 0x429cc + iVar13) = 0x39;
      local_c = 2;
    }
  }
LAB_1000c005:
  *(undefined1 *)((int)pcVar15 + local_c + 0x429cb + iVar13) = 0;
  in_EAX = local_c;
  if (0 < (int)local_c) {
    pbVar12 = (byte *)((int)pcVar15 + 0x429a6 + iVar13);
    *pbVar12 = *pbVar12 | 8;
    return CONCAT22((short)((uint)((int)pcVar15 + 0x429a6 + iVar13) >> 0x10),1);
  }
LAB_1000c031:
  return in_EAX & 0xffff0000;
}



