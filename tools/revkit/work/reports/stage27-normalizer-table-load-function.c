===== 0x100033e0 =====
Function: FUN_100033e0 @ 100033e0

int * __cdecl FUN_100033e0(int param_1,int param_2,short param_3,undefined1 param_4)

{
  int *piVar1;
  char *pcVar2;
  int iVar3;
  char *pcVar4;
  undefined1 local_208 [512];
  int local_8;
  
  FUN_10063ed6(local_208,(byte *)s__s__s_1007732c);
  piVar1 = (int *)FUN_1001d9c0(0x14);
  if (piVar1 == (int *)0x0) {
    return (int *)0x0;
  }
  pcVar2 = FUN_10010300(local_208,&param_2,1,&local_8);
  piVar1[2] = (int)pcVar2;
  if (pcVar2 != (char *)0x0) {
    pcVar2[param_2] = '\0';
    if (local_8 == 0) {
      return (int *)0x0;
    }
    param_1 = piVar1[2];
    *piVar1 = local_8;
    pcVar2 = (char *)(piVar1[2] + param_2);
    iVar3 = FUN_1001d9c0(*piVar1 << 2);
    piVar1[3] = iVar3;
    iVar3 = FUN_1001d9c0(*piVar1 << 2);
    piVar1[4] = iVar3;
    iVar3 = 0;
    if (0 < *piVar1) {
      do {
        pcVar4 = FUN_10003540(&param_1,pcVar2);
        if (pcVar4 == (char *)0x0) {
          return (int *)0x0;
        }
        if (param_3 == 0x4c) {
          *(char **)(piVar1[3] + iVar3 * 4) = pcVar4;
          *(char **)(piVar1[4] + iVar3 * 4) = pcVar4;
        }
        else if (param_3 == 0x4d) {
          *(char **)(piVar1[3] + iVar3 * 4) = pcVar4;
          pcVar4 = FUN_10003540(&param_1,pcVar2);
          if (pcVar4 == (char *)0x0) {
            return (int *)0x0;
          }
          *(char **)(piVar1[4] + iVar3 * 4) = pcVar4;
        }
        iVar3 = iVar3 + 1;
      } while (iVar3 < *piVar1);
    }
    *(undefined1 *)(piVar1 + 1) = param_4;
    return piVar1;
  }
  return (int *)0x0;
}



===== 0x10003660 =====
Function: FUN_10003660 @ 10003660

undefined4 __cdecl FUN_10003660(int param_1)

{
  int iVar1;
  undefined *puVar2;
  char *pcVar3;
  undefined4 *puVar4;
  undefined4 uVar5;
  char *pcVar6;
  int iVar7;
  undefined1 local_210 [512];
  int local_10;
  int local_c;
  int local_8;
  
  DAT_100fee00 = (undefined4 *)0x0;
  FUN_10063ed6(local_210,(byte *)s__s_exceptdict_10077334);
  puVar2 = FUN_10024d20(local_210,&DAT_100771a8);
  if (puVar2 == (undefined *)0x0) {
    return 0xffffffff;
  }
  iVar7 = *(int *)(puVar2 + 0x10);
  pcVar3 = (char *)FUN_1001d9c0(iVar7 + 1);
  if (pcVar3 == (char *)0x0) {
    return 0xffffffff;
  }
  FUN_10025440((int)puVar2,0,0,pcVar3,1,iVar7);
  FUN_100253d0(puVar2);
  DAT_100fee00 = (undefined4 *)FUN_1001d9c0(0xc);
  if (DAT_100fee00 == (undefined4 *)0x0) {
    return 0xffffffff;
  }
  FUN_100014f0(DAT_100fee00,(undefined4 *)pcVar3,4,1);
  FUN_100014f0(DAT_100fee00 + 1,(undefined4 *)(pcVar3 + 4),4,1);
  pcVar6 = pcVar3 + 8;
  puVar4 = FUN_1001da00(DAT_100fee00[1] + 1,0xc);
  DAT_100fee00[2] = puVar4;
  if (DAT_100fee00[2] != 0) {
    while( true ) {
      FUN_100014f0(&param_1,(undefined4 *)pcVar6,4,1);
      FUN_100014f0(&local_8,(undefined4 *)(pcVar6 + 4),4,1);
      pcVar6 = pcVar6 + 8;
      *(int *)(DAT_100fee00[2] + param_1 * 0xc) = local_8;
      uVar5 = FUN_1001d9c0(local_8 * 4);
      *(undefined4 *)(DAT_100fee00[2] + 4 + param_1 * 0xc) = uVar5;
      uVar5 = FUN_1001d9c0(local_8 * 4);
      *(undefined4 *)(DAT_100fee00[2] + 8 + param_1 * 0xc) = uVar5;
      iVar7 = DAT_100fee00[2] + param_1 * 0xc;
      if ((*(int *)(iVar7 + 4) == 0) || (*(int *)(iVar7 + 8) == 0)) break;
      iVar7 = 0;
      if (0 < local_8) {
        do {
          FUN_100014f0(&local_c,(undefined4 *)pcVar6,4,1);
          FUN_100014f0(&local_10,(undefined4 *)(pcVar6 + 4),4,1);
          pcVar6 = pcVar6 + 8;
          uVar5 = FUN_1001d9c0(local_c + 1);
          *(undefined4 *)(*(int *)(DAT_100fee00[2] + 4 + param_1 * 0xc) + iVar7 * 4) = uVar5;
          uVar5 = FUN_1001d9c0(local_10 + 1);
          *(undefined4 *)(*(int *)(DAT_100fee00[2] + 8 + param_1 * 0xc) + iVar7 * 4) = uVar5;
          iVar1 = DAT_100fee00[2] + param_1 * 0xc;
          puVar4 = *(undefined4 **)(*(int *)(iVar1 + 4) + iVar7 * 4);
          if (puVar4 == (undefined4 *)0x0) {
            return 0xffffffff;
          }
          if (*(int *)(*(int *)(iVar1 + 8) + iVar7 * 4) == 0) {
            return 0xffffffff;
          }
          FUN_100014f0(puVar4,(undefined4 *)pcVar6,1,local_c);
          *(undefined1 *)
           (*(int *)(*(int *)(DAT_100fee00[2] + 4 + param_1 * 0xc) + iVar7 * 4) + local_c) = 0;
          pcVar6 = pcVar6 + local_c;
          FUN_100014f0(*(undefined4 **)(*(int *)(DAT_100fee00[2] + 8 + param_1 * 0xc) + iVar7 * 4),
                       (undefined4 *)pcVar6,1,local_10);
          *(undefined1 *)
           (*(int *)(*(int *)(DAT_100fee00[2] + 8 + param_1 * 0xc) + iVar7 * 4) + local_10) = 0;
          pcVar6 = pcVar6 + local_10;
          iVar7 = iVar7 + 1;
        } while (iVar7 < local_8);
      }
      if (param_1 == DAT_100fee00[1]) {
        FUN_1001da30(pcVar3);
        return 1;
      }
    }
  }
  return 0xffffffff;
}



===== 0x10003a10 =====
Function: FUN_10003a10 @ 10003a10

int FUN_10003a10(void)

{
  undefined4 uVar1;
  char local_204 [512];
  
  FUN_10063ed6(local_204,(byte *)s__s_atmt_tree3_10077344);
  uVar1 = FUN_100019d0(0x100fee0c,local_204,0x1b);
  return (-(uint)((short)uVar1 != 0) & 2) - 1;
}



