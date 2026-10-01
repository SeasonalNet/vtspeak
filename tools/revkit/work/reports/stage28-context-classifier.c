===== 0x1000e570 =====
Function: FUN_1000e570 @ 1000e570

short __cdecl FUN_1000e570(short *param_1,byte *param_2)

{
  bool bVar1;
  short sVar2;
  undefined3 extraout_var;
  undefined3 extraout_var_00;
  undefined3 extraout_var_01;
  int iVar3;
  undefined3 extraout_var_02;
  undefined4 uVar4;
  byte *pbVar5;
  byte *pbVar6;
  byte *pbVar7;
  byte *pbVar8;
  
  bVar1 = FUN_1000e540((char *)((int)param_1 + (int)param_2 * 0x70 + 0xb));
  if ((((((short)CONCAT31(extraout_var,bVar1) == 0) &&
        (bVar1 = FUN_1000e4e0((char *)((int)param_1 + (int)param_2 * 0x70 + 0xb)),
        (short)CONCAT31(extraout_var_00,bVar1) == 0)) &&
       (bVar1 = FUN_1000e4b0((char *)((int)param_1 + (int)param_2 * 0x70 + 0xb)),
       (short)CONCAT31(extraout_var_01,bVar1) == 0)) || (*param_1 + -1 <= (int)(param_2 + 1))) ||
     ((sVar2 = FUN_1000e570(param_1,param_2 + 1), sVar2 != 0x50 && (sVar2 != 0x4e)))) {
    iVar3 = FUN_1001c2c0((byte *)((int)param_1 + (int)param_2 * 0x70 + 0xb),&DAT_10078f90);
    if (iVar3 == 0) {
      return 0x50;
    }
    iVar3 = FUN_1001c2c0((byte *)((int)param_1 + (int)param_2 * 0x70 + 0xb),&DAT_10077374);
    if ((iVar3 == 0) && ((int)(param_2 + 1) < *param_1 + -1)) {
      pbVar6 = (byte *)((int)param_1 + (int)param_2 * 0x70 + 0x7b);
      iVar3 = FUN_1001c2c0(pbVar6,(byte *)s_about_10077538);
      if ((iVar3 == 0) ||
         (((iVar3 = FUN_1001c2c0(pbVar6,&DAT_10079270), iVar3 == 0 ||
           (iVar3 = FUN_1001c2c0(pbVar6,&DAT_1007760c), iVar3 == 0)) ||
          (iVar3 = FUN_1001c2c0(pbVar6,&DAT_10078ae8), iVar3 == 0)))) {
        return 0x50;
      }
    }
    iVar3 = FUN_1001c2c0((byte *)((int)param_1 + (int)param_2 * 0x70 + 0xb),&DAT_10079014);
    if (iVar3 == 0) {
      if ((int)(param_2 + 1) < *param_1 + -1) {
        pbVar6 = (byte *)((int)param_1 + (int)param_2 * 0x70 + 0x7b);
        iVar3 = FUN_1001c2c0(pbVar6,(byte *)s_about_10077538);
        if (((((iVar3 == 0) || (iVar3 = FUN_1001c2c0(pbVar6,&DAT_1007926c), iVar3 == 0)) ||
             ((iVar3 = FUN_1001c2c0(pbVar6,&DAT_10077688), iVar3 == 0 ||
              ((iVar3 = FUN_1001c2c0(pbVar6,(byte *)s_thouth_10079264), iVar3 == 0 ||
               (iVar3 = FUN_1001c2c0(pbVar6,&DAT_10078eb8), iVar3 == 0)))))) ||
            (iVar3 = FUN_1001c2c0(pbVar6,&DAT_1007925c), iVar3 == 0)) ||
           (iVar3 = FUN_1001c2c0(pbVar6,&DAT_10079254), iVar3 == 0)) {
          return 0x50;
        }
      }
      if ((((int)(param_2 + 2) < *param_1 + -1) &&
          (iVar3 = FUN_1001c2c0((byte *)((int)param_1 + (int)param_2 * 0x70 + 0x7b),&DAT_1007924c),
          iVar3 == 0)) &&
         (iVar3 = FUN_1001c2c0((byte *)((int)param_1 + (int)param_2 * 0x70 + 0xeb),&DAT_10077368),
         iVar3 == 0)) {
        return 0x50;
      }
    }
    pbVar6 = (byte *)((int)param_1 + (int)param_2 * 0x70 + 0xb);
    bVar1 = FUN_1000e510((char *)pbVar6);
    if ((short)CONCAT31(extraout_var_02,bVar1) != 0) {
      pbVar6 = param_2 + 1;
      if (((int)pbVar6 < *param_1 + -1) &&
         (((uVar4 = FUN_1000e430((byte *)((int)param_1 + (int)param_2 * 0x70 + 0x7b),(int)pbVar6),
           (short)uVar4 != 0 ||
           (uVar4 = FUN_1000e3a0((byte *)((int)param_1 + (int)param_2 * 0x70 + 0x7b),&DAT_1009f948),
           (short)uVar4 != 0)) ||
          (iVar3 = FUN_1001c2c0((byte *)((int)param_1 + (int)param_2 * 0x70 + 0x7b),
                                (byte *)&DAT_100776fc), iVar3 == 0)))) {
        return 0x50;
      }
      if (((int)(param_2 + 2) < *param_1 + -1) &&
         (uVar4 = FUN_1000e3a0((byte *)((int)param_1 + (int)param_2 * 0x70 + 0x7b),
                               (byte *)((int)param_1 + (int)param_2 * 0x70 + 0xeb)),
         (short)uVar4 != 0)) {
        return 0x50;
      }
      if (*param_1 + -1 <= (int)pbVar6) {
        return 0x50;
      }
      pbVar5 = (byte *)((int)pbVar6 * 0x70 + 0xb + (int)param_1);
      pbVar7 = pbVar6;
      do {
        iVar3 = FUN_1001c2c0(pbVar5,(byte *)&DAT_100776fc);
        if (iVar3 == 0) {
          if (*param_1 + -1 <= (int)pbVar6) {
            return 0x50;
          }
          pbVar5 = param_2 + 2;
          param_2 = (byte *)((int)pbVar6 * 0x70 + 0xb + (int)param_1);
          do {
            iVar3 = FUN_1001c2c0(param_2,(byte *)&DAT_100776fc);
            if ((iVar3 == 0) && ((int)pbVar5 < *param_1 + -1)) {
              pbVar7 = param_2 + 0x70;
              iVar3 = FUN_1001c2c0(pbVar7,&DAT_1007926c);
              if (iVar3 == 0) {
                pbVar8 = pbVar5;
                if (*param_1 + -1 <= (int)pbVar5) {
                  return 0x50;
                }
                while (iVar3 = FUN_1001c2c0(pbVar7,(byte *)&DAT_100776fc), iVar3 != 0) {
                  pbVar8 = pbVar8 + 1;
                  pbVar7 = pbVar7 + 0x70;
                  if (*param_1 + -1 <= (int)pbVar8) {
                    return 0x50;
                  }
                }
              }
              else {
                sVar2 = FUN_1000e570(param_1,pbVar5);
                if (sVar2 == 0x4e) {
                  return 0x4e;
                }
              }
            }
            pbVar6 = pbVar6 + 1;
            param_2 = param_2 + 0x70;
            pbVar5 = pbVar5 + 1;
          } while ((int)pbVar6 < *param_1 + -1);
          return 0x50;
        }
        pbVar7 = pbVar7 + 1;
        pbVar5 = pbVar5 + 0x70;
      } while ((int)pbVar7 < *param_1 + -1);
      return 0x50;
    }
    if ((((int)param_2 < *param_1 + -1) &&
        ((uVar4 = FUN_1000e430(pbVar6,(int)param_2), (short)uVar4 != 0 ||
         (uVar4 = FUN_1000e3a0(pbVar6,&DAT_1009f948), (short)uVar4 != 0)))) ||
       (((int)(param_2 + 1) < *param_1 + -1 &&
        (uVar4 = FUN_1000e3a0(pbVar6,(byte *)((int)param_1 + (int)param_2 * 0x70 + 0x7b)),
        (short)uVar4 != 0)))) {
      return 0x4e;
    }
    sVar2 = 0x44;
  }
  return sVar2;
}



===== 0x1000e540 =====
Function: FUN_1000e540 @ 1000e540

bool __cdecl FUN_1000e540(char *param_1)

{
  uint uVar1;
  
  uVar1 = FUN_100560a0(&PTR_s_about_10078140,param_1,DAT_1007824c,0x49);
  return uVar1 < 0x80000000;
}



===== 0x1000e4e0 =====
Function: FUN_1000e4e0 @ 1000e4e0

bool __cdecl FUN_1000e4e0(char *param_1)

{
  uint uVar1;
  
  uVar1 = FUN_100560a0(&PTR_DAT_10078038,param_1,DAT_100780c0,0x49);
  return uVar1 < 0x80000000;
}



===== 0x1000e4b0 =====
Function: FUN_1000e4b0 @ 1000e4b0

bool __cdecl FUN_1000e4b0(char *param_1)

{
  uint uVar1;
  
  uVar1 = FUN_100560a0(&PTR_DAT_100780e0,param_1,DAT_10078244,0x49);
  return uVar1 < 0x80000000;
}



===== 0x1000e510 =====
Function: FUN_1000e510 @ 1000e510

bool __cdecl FUN_1000e510(char *param_1)

{
  uint uVar1;
  
  uVar1 = FUN_100560a0(&PTR_DAT_100780fc,param_1,DAT_10078248,0x49);
  return uVar1 < 0x80000000;
}



===== 0x1000e430 =====
Function: FUN_1000e430 @ 1000e430

undefined4 __cdecl FUN_1000e430(byte *param_1,int param_2)

{
  uint in_EAX;
  int iVar1;
  
  if (0 < param_2) {
    in_EAX = FUN_100560a0(&PTR_DAT_100780c8,(char *)param_1,DAT_10078240,0x49);
    if (-1 < (int)in_EAX) {
      return CONCAT22((short)(in_EAX >> 0x10),1);
    }
  }
  if (param_2 == 0) {
    iVar1 = FUN_1001c2c0(param_1,&DAT_100773a0);
    in_EAX = 0;
    if (iVar1 != 0) {
      in_EAX = FUN_100560a0(&PTR_DAT_100780c8,(char *)param_1,DAT_10078240,0x49);
      if (-1 < (int)in_EAX) {
        return CONCAT22((short)(in_EAX >> 0x10),1);
      }
    }
  }
  return in_EAX & 0xffff0000;
}



===== 0x1000e3a0 =====
Function: FUN_1000e3a0 @ 1000e3a0

undefined4 __cdecl FUN_1000e3a0(byte *param_1,byte *param_2)

{
  int iVar1;
  uint uVar2;
  
  iVar1 = FUN_100560a0(&PTR_DAT_10078044,(char *)param_1,DAT_100780c4,0x49);
  if (-1 < iVar1) {
    return CONCAT22((short)((uint)iVar1 >> 0x10),1);
  }
  iVar1 = FUN_1001c2c0(param_1,&DAT_10077764);
  if (iVar1 != 0) {
    iVar1 = FUN_1001c2c0(param_1,&DAT_10077768);
    if (iVar1 != 0) {
      uVar2 = FUN_1001c2c0(param_1,&DAT_1007776c);
      if (uVar2 != 0) goto LAB_1000e41b;
    }
  }
  uVar2 = FUN_1001c2c0(param_2,&DAT_10077568);
  if (uVar2 == 0) {
    return 1;
  }
LAB_1000e41b:
  return uVar2 & 0xffff0000;
}



===== 0x1001c2c0 =====
Function: FUN_1001c2c0 @ 1001c2c0

int __cdecl FUN_1001c2c0(byte *param_1,byte *param_2)

{
  byte bVar1;
  int iVar2;
  
  if (param_1 == param_2) {
    return 0;
  }
  iVar2 = (int)*(short *)(&DAT_1007e388 + (uint)*param_1 * 2) -
          (int)*(short *)(&DAT_1007e388 + (uint)*param_2 * 2);
  while( true ) {
    if (iVar2 != 0) {
      return iVar2;
    }
    param_2 = param_2 + 1;
    bVar1 = *param_1;
    param_1 = param_1 + 1;
    if (bVar1 == 0) break;
    iVar2 = (int)*(short *)(&DAT_1007e388 + (uint)*param_1 * 2) -
            (int)*(short *)(&DAT_1007e388 + (uint)*param_2 * 2);
  }
  return 0;
}



