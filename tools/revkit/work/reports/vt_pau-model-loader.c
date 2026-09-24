/* Ghidra decompiler output: approximate C pseudocode, not original source. */
/* Program: vt_pau.dll */

/* ===== FUN_1001c510 ===== */
/* Entry: 1001c510 */

void __cdecl FUN_1001c510(int param_1,undefined4 param_2)

{
  undefined1 local_c04 [512];
  undefined1 local_a04 [512];
  undefined1 local_804 [512];
  undefined1 local_604 [512];
  undefined1 local_404 [512];
  undefined1 local_204 [512];

  *(undefined4 *)(param_1 + 0x4d10) = param_2;
  FUN_10063ed6(local_204,(byte *)s__sdata__s__1007c924);
  FUN_10063ed6(local_604,(byte *)s__s__s__1007c90c);
  FUN_10063ed6(local_a04,(byte *)s__s_dat__1007c904);
  FUN_10063ed6(local_804,(byte *)s__s_dat__1007c904);
  FUN_10063ed6(local_c04,&DAT_1007c900);
  FUN_10063ed6(local_404,(byte *)s__s_ttsdata__1007c8f4);
  FUN_10063ed6((undefined1 *)(param_1 + 0x294),&DAT_10079404);
  FUN_10063ed6((undefined1 *)(param_1 + 0x494),(byte *)s__s_stree3__1007c8e8);
  FUN_10063ed6((undefined1 *)(param_1 + 0x694),&DAT_10079404);
  FUN_10063ed6((undefined1 *)(param_1 + 0x894),&DAT_10079404);
  FUN_10063ed6((undefined1 *)(param_1 + 0xa94),&DAT_10079404);
  FUN_10063ed6((undefined1 *)(param_1 + 0xc94),&DAT_10079404);
  return;
}



/* ===== FUN_10028620 ===== */
/* Entry: 10028620 */

void __cdecl FUN_10028620(BYTE *param_1,int param_2)

{
  BYTE BVar1;
  LPCSTR pCVar2;
  int iVar3;
  uint uVar4;
  uint uVar5;
  BYTE *pBVar6;
  undefined4 *puVar7;
  BYTE *pBVar8;
  BYTE local_404;
  undefined4 local_403;

  local_404 = '\0';
  puVar7 = &local_403;
  for (iVar3 = 0xff; iVar3 != 0; iVar3 = iVar3 + -1) {
    *puVar7 = 0;
    puVar7 = puVar7 + 1;
  }
  *(undefined2 *)puVar7 = 0;
  *(undefined1 *)((int)puVar7 + 2) = 0;
  if ((param_2 < 0) || (5 < param_2)) {
    param_2 = 1;
  }
  pCVar2 = (LPCSTR)VT_GetPathKey_ENG(param_2);
  FUN_10028480(&local_404,pCVar2);
  iVar3 = -1;
  pBVar6 = &local_404;
  do {
    if (iVar3 == 0) break;
    iVar3 = iVar3 + -1;
    BVar1 = *pBVar6;
    pBVar6 = pBVar6 + 1;
  } while (BVar1 != '\0');
  if (iVar3 == -2) {
    FUN_10063ed6(param_1,&DAT_1007d950);
    return;
  }
  uVar4 = 0xffffffff;
  pBVar6 = &local_404;
  do {
    pBVar8 = pBVar6;
    if (uVar4 == 0) break;
    uVar4 = uVar4 - 1;
    pBVar8 = pBVar6 + 1;
    BVar1 = *pBVar6;
    pBVar6 = pBVar8;
  } while (BVar1 != '\0');
  uVar4 = ~uVar4;
  pBVar6 = pBVar8 + -uVar4;
  for (uVar5 = uVar4 >> 2; uVar5 != 0; uVar5 = uVar5 - 1) {
    *(undefined4 *)param_1 = *(undefined4 *)pBVar6;
    pBVar6 = pBVar6 + 4;
    param_1 = param_1 + 4;
  }
  for (uVar4 = uVar4 & 3; uVar4 != 0; uVar4 = uVar4 - 1) {
    *param_1 = *pBVar6;
    pBVar6 = pBVar6 + 1;
    param_1 = param_1 + 1;
  }
  return;
}



/* ===== FUN_10012380 ===== */
/* Entry: 10012380 */

undefined4 __cdecl FUN_10012380(int param_1,undefined2 *param_2,HWND param_3)

{
  undefined4 uVar1;
  uint uVar2;
  int *piVar3;
  undefined2 extraout_var;
  undefined4 extraout_ECX;
  undefined4 extraout_EDX;
  undefined1 local_204 [512];

  FUN_10063ed6(local_204,(byte *)s__sdict_eng_100793e0);
  VT_SetDecimal0Pron_ENG();
  DAT_100a74b8 = 0;
  DAT_100a749c = 0;
  DAT_100a74bc = 0;
  VT_SetDecimal0Pron_ENG();
  FUN_10028430(1,param_3,0);
  uVar1 = FUN_10012210(param_1,param_2);
  if ((short)uVar1 < 0) {
    return CONCAT22((short)((uint)uVar1 >> 0x10),0xffff);
  }
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  uVar1 = FUN_10012280(param_1,param_2);
  if ((short)uVar1 < 0) {
    return CONCAT22((short)((uint)uVar1 >> 0x10),0xffff);
  }
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  FUN_10028430(1,param_3,10);
  uVar2 = FUN_10039fe0((int)local_204);
  if ((short)uVar2 == 0) {
    *param_2 = 5;
    return CONCAT22((short)(uVar2 >> 0x10),0xffff);
  }
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  FUN_10028430(1,param_3,0x14);
  FUN_100122f0();
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  FUN_10028430(1,param_3,0x1e);
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  FUN_10028430(1,param_3,0x23);
  FUN_1002e800(param_1);
  FUN_1002e860(extraout_ECX,extraout_EDX,param_1);
  FUN_1002e830(param_1);
  piVar3 = FUN_100154d0(&PTR_DAT_1007eb88);
  *(int **)(param_1 + 0x20428) = piVar3;
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  FUN_10028430(1,param_3,0x28);
  DAT_100a749c = DAT_100a74b8;
  return CONCAT22(extraout_var,1);
}



/* ===== FUN_10012a60 ===== */
/* Entry: 10012a60 */

undefined4 __cdecl FUN_10012a60(int *param_1,undefined2 *param_2,int param_3,HWND param_4)

{
  undefined4 uVar1;

  DAT_100a74b8 = 0;
  *(undefined4 *)(&DAT_100a74a0 + param_3 * 4) = 0;
  DAT_100a74bc = 0;
  VT_SetDecimal0Pron_ENG();
  FUN_10028430(1,param_4,0x32);
  VT_SetDecimal0Pron_ENG();
  uVar1 = FUN_10012580((int)param_1,param_2);
  if ((short)uVar1 < 0) {
    return CONCAT22((short)((uint)uVar1 >> 0x10),0xffff);
  }
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  FUN_10028430(1,param_4,0x3c);
  uVar1 = FUN_100126e0((int)param_1,param_2);
  if ((short)uVar1 < 0) {
    return CONCAT22((short)((uint)uVar1 >> 0x10),0xffff);
  }
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  FUN_10028430(1,param_4,0x46);
  uVar1 = FUN_10012790((char *)param_1,param_2);
  if ((short)uVar1 < 0) {
    return CONCAT22((short)((uint)uVar1 >> 0x10),0xffff);
  }
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  if (DAT_100a045d != '\0') {
    uVar1 = FUN_100125d0((int)param_1);
    if ((short)uVar1 < 0) {
      return CONCAT22((short)((uint)uVar1 >> 0x10),0xffff);
    }
  }
  FUN_10028430(1,param_4,0x50);
  uVar1 = FUN_100128c0(param_1,param_2);
  if ((short)uVar1 < 0) {
    return CONCAT22((short)((uint)uVar1 >> 0x10),0xffff);
  }
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  FUN_10028430(1,param_4,100);
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  uVar1 = DAT_100a74b8;
  *(undefined4 *)(&DAT_100a74a0 + param_3 * 4) = DAT_100a74b8;
  return CONCAT22((short)((uint)uVar1 >> 0x10),1);
}



/* ===== FUN_10012580 ===== */
/* Entry: 10012580 */

undefined4 __cdecl FUN_10012580(int param_1,undefined2 *param_2)

{
  int iVar1;

  iVar1 = FUN_10001050(param_1 + 0x494,param_1);
  if (iVar1 < 0) {
    *param_2 = 8;
    return CONCAT22((short)((uint)param_2 >> 0x10),0xffff);
  }
  return CONCAT22((short)((uint)iVar1 >> 0x10),1);
}



/* ===== FUN_100126e0 ===== */
/* Entry: 100126e0 */

undefined4 __cdecl FUN_100126e0(int param_1,undefined2 *param_2)

{
  undefined4 uVar1;
  int iVar2;
  undefined2 extraout_var;
  undefined1 local_204 [512];

  FUN_10063ed6(local_204,(byte *)s__sdist_tbl_cepdist_tbl_100795d8);
  uVar1 = FUN_1001ad20((short *)(param_1 + 8),local_204);
  if ((short)uVar1 < 0) {
    *param_2 = 6;
    return CONCAT22((short)((uint)uVar1 >> 0x10),0xffff);
  }
  VT_SetDecimal0Pron_ENG();
  iVar2 = FUN_1001ae70(param_1);
  if (iVar2 < 0) {
    *param_2 = 6;
    return CONCAT22((short)((uint)iVar2 >> 0x10),0xffff);
  }
  VT_SetDecimal0Pron_ENG();
  return CONCAT22(extraout_var,1);
}



/* ===== FUN_10012790 ===== */
/* Entry: 10012790 */

undefined4 __cdecl FUN_10012790(char *param_1,undefined2 *param_2)

{
  bool bVar1;
  int iVar2;
  undefined3 extraout_var;
  undefined2 extraout_var_00;
  undefined1 local_204 [512];

  VT_SetDecimal0Pron_ENG();
  FUN_10063ed6(local_204,(byte *)s__sdblist_idx_100796e8);
  iVar2 = FUN_10019f40(param_1,(int)local_204);
  if (iVar2 == 0) {
    *param_2 = 7;
    return 0xffff;
  }
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  FUN_10063ed6(local_204,(byte *)s__sclass_idx_10079690);
  bVar1 = FUN_1001aca0((int)param_1,local_204);
  if (CONCAT31(extraout_var,bVar1) == 0) {
    *param_2 = 7;
    return CONCAT22((short)((uint3)extraout_var >> 8),0xffff);
  }
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  FUN_10063ed6(local_204,(byte *)s__sclasshp_idx_1007963c);
  iVar2 = FUN_1001ac60((int)param_1,(int)local_204);
  if (iVar2 == 0) {
    *param_2 = 7;
    return 0xffff;
  }
  VT_SetDecimal0Pron_ENG();
  VT_SetDecimal0Pron_ENG();
  return CONCAT22(extraout_var_00,1);
}



/* ===== FUN_100125d0 ===== */
/* Entry: 100125d0 */

undefined4 __cdecl FUN_100125d0(int param_1)

{
  int in_EAX;
  undefined4 uVar1;
  int iVar2;
  undefined1 local_208 [512];
  int local_8;

  iVar2 = 0;
  local_8 = 0;
  if (0 < *(short *)(param_1 + 0xe94)) {
    do {
      FUN_10063ed6(local_208,(byte *)s__s_s_his_1007958c);
      uVar1 = FUN_1001b240(*(int *)(param_1 + 0x84) + iVar2,(int)local_208);
      if ((short)uVar1 == 0) {
        return CONCAT22((short)((uint)uVar1 >> 0x10),0xffff);
      }
      in_EAX = local_8 + 1;
      iVar2 = iVar2 + 100;
      local_8 = in_EAX;
    } while (in_EAX < *(short *)(param_1 + 0xe94));
  }
  return CONCAT22((short)((uint)in_EAX >> 0x10),1);
}



/* ===== FUN_100128c0 ===== */
/* Entry: 100128c0 */

undefined4 __cdecl FUN_100128c0(int *param_1,undefined2 *param_2)

{
  int iVar1;
  undefined *puVar2;
  int iVar3;
  undefined1 local_204 [512];

  iVar1 = FUN_1001d9c0((int)(short)param_1[8] << 2);
  *param_1 = iVar1;
  puVar2 = (undefined *)0x0;
  if (iVar1 != 0) {
    iVar1 = 0;
    if (0 < (short)param_1[8]) {
      do {
        FUN_10063ed6(local_204,(byte *)s__s_s_dat_10079714);
        puVar2 = FUN_10024d20(local_204,&DAT_100771a8);
        *(undefined **)(*param_1 + iVar1 * 4) = puVar2;
        if (*(int *)(*param_1 + iVar1 * 4) == 0) goto LAB_100128e5;
        iVar1 = iVar1 + 1;
      } while (iVar1 < (short)param_1[8]);
    }
    iVar1 = FUN_1001d9c0((int)(short)param_1[8] << 2);
    param_1[1] = iVar1;
    if (iVar1 == 0) {
      *param_2 = 10;
      return CONCAT22((short)((uint)param_2 >> 0x10),0xffff);
    }
    iVar3 = 0;
    if (0 < (short)param_1[8]) {
      do {
        FUN_10063ed6(local_204,(byte *)s__s_s_upm_10079708);
        puVar2 = FUN_10024d20(local_204,&DAT_100771a8);
        *(undefined **)(param_1[1] + iVar3 * 4) = puVar2;
        iVar1 = param_1[1];
        if (*(int *)(iVar1 + iVar3 * 4) == 0) {
          *param_2 = 10;
          return CONCAT22((short)((uint)iVar1 >> 0x10),0xffff);
        }
        iVar3 = iVar3 + 1;
      } while (iVar3 < (short)param_1[8]);
    }
    return CONCAT22((short)((uint)iVar1 >> 0x10),1);
  }
LAB_100128e5:
  *param_2 = 9;
  return CONCAT22((short)((uint)puVar2 >> 0x10),0xffff);
}



/* ===== FUN_10001050 ===== */
/* Entry: 10001050 */

int __cdecl FUN_10001050(undefined4 param_1,int param_2)

{
  short sVar1;
  undefined4 uVar2;
  char local_204 [512];

  FUN_10063ed6(local_204,(byte *)s__spitch_nbt_tree2_10077194);
  uVar2 = FUN_10001900((int *)(param_2 + 0x1b0),local_204);
  if ((short)uVar2 == 0) {
    return -1;
  }
  FUN_10063ed6(local_204,(byte *)s__spitch_nbf_tree2_10077180);
  sVar1 = FUN_10001af0((int *)(param_2 + 0x220),local_204);
  if (sVar1 == 0) {
    return -1;
  }
  FUN_10063ed6(local_204,(byte *)s__spitch_bt_tree2_1007716c);
  uVar2 = FUN_10001900((int *)(param_2 + 0x1cc),local_204);
  if ((short)uVar2 == 0) {
    return -1;
  }
  FUN_10063ed6(local_204,(byte *)s__spitch_bf_tree2_10077158);
  sVar1 = FUN_10001af0((int *)(param_2 + 0x23c),local_204);
  if (sVar1 == 0) {
    return -1;
  }
  FUN_10063ed6(local_204,(byte *)s__spitch_qbt_tree2_10077144);
  uVar2 = FUN_10001900((int *)(param_2 + 0x204),local_204);
  if ((short)uVar2 == 0) {
    return -1;
  }
  FUN_10063ed6(local_204,(byte *)s__spitch_qbf_tree2_10077130);
  sVar1 = FUN_10001af0((int *)(param_2 + 0x274),local_204);
  if (sVar1 == 0) {
    return -1;
  }
  FUN_10063ed6(local_204,(byte *)s__spitch_sbt_tree2_1007711c);
  uVar2 = FUN_10001900((int *)(param_2 + 0x1e8),local_204);
  if ((short)uVar2 == 0) {
    return -1;
  }
  FUN_10063ed6(local_204,(byte *)s__spitch_sbf_tree2_10077108);
  sVar1 = FUN_10001af0((int *)(param_2 + 600),local_204);
  if (sVar1 == 0) {
    return -1;
  }
  FUN_10063ed6(local_204,(byte *)s__sduration_vshort_tree2_100770f0);
  uVar2 = FUN_10001900((int *)(param_2 + 0xb4),local_204);
  if ((short)uVar2 == 0) {
    return -1;
  }
  FUN_10063ed6(local_204,(byte *)s__sduration_vlong_tree2_100770d8);
  uVar2 = FUN_10001900((int *)(param_2 + 0xd0),local_204);
  if ((short)uVar2 == 0) {
    return -1;
  }
  FUN_10063ed6(local_204,(byte *)s__sduration_vdi_tree2_100770c0);
  uVar2 = FUN_10001900((int *)(param_2 + 0xec),local_204);
  if ((short)uVar2 == 0) {
    return -1;
  }
  FUN_10063ed6(local_204,(byte *)s__sduration_vsch_tree2_100770a8);
  uVar2 = FUN_10001900((int *)(param_2 + 0x108),local_204);
  if ((short)uVar2 == 0) {
    return -1;
  }
  FUN_10063ed6(local_204,(byte *)s__sduration_cstop_tree2_10077090);
  uVar2 = FUN_10001900((int *)(param_2 + 0x124),local_204);
  if ((short)uVar2 == 0) {
    return -1;
  }
  FUN_10063ed6(local_204,(byte *)s__sduration_cfri_tree2_10077078);
  uVar2 = FUN_10001900((int *)(param_2 + 0x140),local_204);
  if ((short)uVar2 == 0) {
    return -1;
  }
  FUN_10063ed6(local_204,(byte *)s__sduration_caff_tree2_10077060);
  uVar2 = FUN_10001900((int *)(param_2 + 0x15c),local_204);
  if ((short)uVar2 == 0) {
    return -1;
  }
  FUN_10063ed6(local_204,(byte *)s__sduration_cnas_tree2_10077048);
  uVar2 = FUN_10001900((int *)(param_2 + 0x178),local_204);
  if ((short)uVar2 == 0) {
    return -1;
  }
  FUN_10063ed6(local_204,(byte *)s__sduration_capp_tree2_10077030);
  uVar2 = FUN_10001900((int *)(param_2 + 0x194),local_204);
  return (-(uint)((short)uVar2 != 0) & 2) - 1;
}



/* ===== FUN_10001900 ===== */
/* Entry: 10001900 */

undefined4 __cdecl FUN_10001900(int *param_1,char *param_2)

{
  char cVar1;
  size_t sVar2;
  undefined *puVar3;
  char *pcVar4;
  undefined2 extraout_var;
  uint uVar5;

  uVar5 = 0xffffffff;
  pcVar4 = param_2;
  do {
    if (uVar5 == 0) break;
    uVar5 = uVar5 - 1;
    cVar1 = *pcVar4;
    pcVar4 = pcVar4 + 1;
  } while (cVar1 != '\0');
  param_2[~uVar5 - 2] = '3';
  puVar3 = FUN_10024d20(param_2,&DAT_100771a8);
  if (puVar3 == (undefined *)0x0) {
    return 0;
  }
  sVar2 = *(size_t *)(puVar3 + 0x10);
  pcVar4 = (char *)FUN_1001d9c0(sVar2);
  param_1[3] = (int)pcVar4;
  if (pcVar4 == (char *)0x0) {
    return 0;
  }
  FUN_10025440((int)puVar3,0,0,pcVar4,1,sVar2);
  FUN_100253d0(puVar3);
  param_1[4] = 0;
  FUN_100016e0(param_1);
  FUN_1001da30((undefined *)param_1[3]);
  return CONCAT22(extraout_var,1);
}



/* ===== FUN_1001ad20 ===== */
/* Entry: 1001ad20 */

undefined4 __cdecl FUN_1001ad20(short *param_1,undefined4 param_2)

{
  int iVar1;
  size_t sVar2;
  undefined *puVar3;
  char *pcVar4;
  undefined2 extraout_var;
  undefined4 uVar5;
  undefined4 *puVar6;
  undefined2 extraout_var_00;
  int iVar7;
  int iVar8;

  puVar3 = FUN_10024d20(param_2,&DAT_100771a8);
  if (puVar3 == (undefined *)0x0) {
    return 0xffff;
  }
  sVar2 = *(size_t *)(puVar3 + 0x10);
  pcVar4 = (char *)FUN_1001d9c0(sVar2);
  if (pcVar4 == (char *)0x0) {
    return 0;
  }
  FUN_10025440((int)puVar3,0,0,pcVar4,1,sVar2);
  FUN_100253d0(puVar3);
  FUN_100014f0((undefined4 *)param_1,(undefined4 *)pcVar4,2,1);
  iVar7 = (int)*param_1;
  iVar8 = (iVar7 / 2) * (iVar7 + 1);
  if ((int)sVar2 < iVar8 * 4 + 2) {
    FUN_1001da30(pcVar4);
    return CONCAT22(extraout_var,0xffff);
  }
  uVar5 = FUN_1001d9c0(iVar7 << 2);
  *(undefined4 *)(param_1 + 2) = uVar5;
  puVar6 = (undefined4 *)FUN_1001d9c0(iVar8 * 4);
  *(undefined4 **)(param_1 + 4) = puVar6;
  FUN_100014f0(puVar6,(undefined4 *)(pcVar4 + 2),4,iVar8);
  iVar8 = 0;
  iVar7 = 0;
  if (0 < *param_1) {
    do {
      iVar1 = iVar7 * 4;
      iVar7 = iVar7 + 1 + iVar8;
      *(int *)(*(int *)(param_1 + 2) + iVar8 * 4) = *(int *)(param_1 + 4) + iVar1;
      iVar8 = iVar8 + 1;
    } while (iVar8 < *param_1);
  }
  FUN_1001da30(pcVar4);
  return CONCAT22(extraout_var_00,1);
}



/* ===== FUN_1001ae70 ===== */
/* Entry: 1001ae70 */

undefined4 __cdecl FUN_1001ae70(int param_1)

{
  short *psVar1;
  int iVar2;
  int iVar3;
  undefined4 uVar4;
  int iVar5;
  int iVar6;
  int iVar7;
  float10 fVar8;
  undefined4 local_8;

  iVar3 = param_1;
  *(undefined2 *)(param_1 + 0x14) = 0x100;
  uVar4 = FUN_1001d9c0(0x400);
  *(undefined4 *)(param_1 + 0x18) = uVar4;
  uVar4 = FUN_1001d9c0(0x20200);
  *(undefined4 *)(param_1 + 0x1c) = uVar4;
  iVar7 = 0;
  iVar6 = 0;
  iVar5 = 0;
  if (0 < *(short *)(param_1 + 0x14)) {
    do {
      iVar2 = iVar6 * 4;
      iVar6 = iVar6 + 1 + iVar5;
      *(int *)(*(int *)(param_1 + 0x18) + iVar5 * 4) = *(int *)(param_1 + 0x1c) + iVar2;
      iVar5 = iVar5 + 1;
    } while (iVar5 < *(short *)(param_1 + 0x14));
  }
  psVar1 = (short *)(param_1 + 0x14);
  param_1 = 0;
  if (0 < *psVar1) {
    do {
      local_8 = 0;
      if (-1 < iVar7) {
        do {
          fVar8 = FUN_1001ae20((float)param_1,(float)local_8);
          local_8 = local_8 + 1;
          *(float *)(*(int *)(iVar7 + *(int *)(iVar3 + 0x18)) + -4 + local_8 * 4) = (float)fVar8;
        } while (local_8 <= param_1);
      }
      iVar7 = iVar7 + 4;
      param_1 = param_1 + 1;
    } while (param_1 < *(short *)(iVar3 + 0x14));
  }
  return 1;
}



/* ===== FUN_10019f40 ===== */
/* Entry: 10019f40 */

undefined4 __cdecl FUN_10019f40(char *param_1,int param_2)

{
  int *piVar1;
  char cVar2;
  byte bVar3;
  char *pcVar4;
  short sVar5;
  undefined *puVar6;
  char *pcVar7;
  byte *pbVar8;
  int iVar9;
  int iVar10;
  undefined4 uVar11;
  undefined4 uVar12;
  void *this;
  uint uVar13;
  uint uVar14;
  short *psVar15;
  byte *pbVar16;
  short *psVar17;
  char *pcVar18;
  bool bVar19;
  undefined1 local_23c [512];
  undefined *local_3c;
  undefined *local_38;
  undefined4 local_34;
  undefined4 local_30;
  undefined4 local_2c;
  int local_28;
  int local_24;
  char *local_20;
  undefined *local_1c;
  char *local_18;
  int local_14;
  short *local_10;
  byte *local_c;
  char *local_8;

  pcVar4 = param_1;
  puVar6 = FUN_10024d20(param_2,&DAT_100771a8);
  if (puVar6 == (undefined *)0x0) {
    return 0;
  }
  iVar9 = *(int *)(puVar6 + 0x10);
  pcVar7 = (char *)FUN_1001d9c0(iVar9 + 1);
  local_18 = pcVar7;
  FUN_10025440((int)puVar6,0,0,pcVar7,1,iVar9);
  FUN_100253d0(puVar6);
  pbVar16 = (byte *)(pcVar7 + iVar9);
  *pbVar16 = 0;
  local_20 = pcVar7;
  local_c = pbVar16;
  pbVar8 = (byte *)FUN_1002e150(&local_20,(char *)pbVar16);
  if (pbVar8 == (byte *)0x0) {
    return 0;
  }
  sVar5 = FUN_10064645(this,pbVar8);
  *(short *)(param_1 + 0xe94) = sVar5;
  if (0x1e < sVar5) {
    return 0;
  }
  param_2 = 0;
  if (0 < sVar5) {
    local_10 = (short *)0x0;
    local_8 = param_1 + 0x4a96;
    param_1 = param_1 + 0xe96;
    do {
      pcVar7 = FUN_1002e150(&local_20,(char *)pbVar16);
      if (pcVar7 == (char *)0x0) {
        return 0;
      }
      uVar13 = 0xffffffff;
      do {
        pcVar18 = pcVar7;
        if (uVar13 == 0) break;
        uVar13 = uVar13 - 1;
        pcVar18 = pcVar7 + 1;
        cVar2 = *pcVar7;
        pcVar7 = pcVar18;
      } while (cVar2 != '\0');
      uVar13 = ~uVar13;
      pcVar7 = pcVar18 + -uVar13;
      pcVar18 = param_1;
      for (uVar14 = uVar13 >> 2; uVar14 != 0; uVar14 = uVar14 - 1) {
        *(undefined4 *)pcVar18 = *(undefined4 *)pcVar7;
        pcVar7 = pcVar7 + 4;
        pcVar18 = pcVar18 + 4;
      }
      for (uVar13 = uVar13 & 3; uVar13 != 0; uVar13 = uVar13 - 1) {
        *pcVar18 = *pcVar7;
        pcVar7 = pcVar7 + 1;
        pcVar18 = pcVar18 + 1;
      }
      pcVar7 = FUN_1002e150(&local_20,(char *)local_c);
      if (pcVar7 == (char *)0x0) {
        return 0;
      }
      uVar13 = 0xffffffff;
      do {
        pcVar18 = pcVar7;
        if (uVar13 == 0) break;
        uVar13 = uVar13 - 1;
        pcVar18 = pcVar7 + 1;
        cVar2 = *pcVar7;
        pcVar7 = pcVar18;
      } while (cVar2 != '\0');
      uVar13 = ~uVar13;
      pcVar7 = pcVar18 + -uVar13;
      pcVar18 = local_8;
      for (uVar14 = uVar13 >> 2; uVar14 != 0; uVar14 = uVar14 - 1) {
        *(undefined4 *)pcVar18 = *(undefined4 *)pcVar7;
        pcVar7 = pcVar7 + 4;
        pcVar18 = pcVar18 + 4;
      }
      for (uVar13 = uVar13 & 3; uVar13 != 0; uVar13 = uVar13 - 1) {
        *pcVar18 = *pcVar7;
        pcVar7 = pcVar7 + 1;
        pcVar18 = pcVar18 + 1;
      }
      uVar13 = 0xffffffff;
      pcVar7 = param_1;
      do {
        if (uVar13 == 0) break;
        uVar13 = uVar13 - 1;
        cVar2 = *pcVar7;
        pcVar7 = pcVar7 + 1;
      } while (cVar2 != '\0');
      pbVar8 = &DAT_1007c5b0;
      pbVar16 = (byte *)(pcVar4 + (int)local_10 + ~uVar13 + 0xe91);
      do {
        bVar3 = *pbVar16;
        bVar19 = bVar3 < *pbVar8;
        if (bVar3 != *pbVar8) {
LAB_1001a0b3:
          iVar9 = (1 - (uint)bVar19) - (uint)(bVar19 != 0);
          goto LAB_1001a0b8;
        }
        if (bVar3 == 0) break;
        bVar3 = pbVar16[1];
        bVar19 = bVar3 < pbVar8[1];
        if (bVar3 != pbVar8[1]) goto LAB_1001a0b3;
        pbVar16 = pbVar16 + 2;
        pbVar8 = pbVar8 + 2;
      } while (bVar3 != 0);
      iVar9 = 0;
LAB_1001a0b8:
      if (iVar9 == 0) {
        pcVar4[(int)local_10 + ~uVar13 + 0xe91] = 0;
      }
      local_8 = local_8 + 0x14;
      param_2 = param_2 + 1;
      local_10 = (short *)((int)local_10 + 0x200);
      param_1 = param_1 + 0x200;
      pcVar7 = local_18;
      pbVar16 = local_c;
    } while (param_2 < *(short *)(pcVar4 + 0xe94));
  }
  FUN_1001da30(pcVar7);
  local_30 = 1000000;
  local_38 = (undefined *)FUN_1001d9c0(1000000);
  local_10 = (short *)FUN_1001d9c0(*(short *)(pcVar4 + 0xe94) * 100);
  *(short **)(pcVar4 + 0x84) = local_10;
  pcVar4[0x6c] = '\0';
  pcVar4[0x6d] = '\0';
  pcVar4[0x6e] = '\0';
  pcVar4[0x6f] = '\0';
  param_2 = 0;
  if (0 < *(short *)(pcVar4 + 0xe94)) {
    do {
      FUN_10063ed6(local_23c,(byte *)s__s_s_idx_1007c5a4);
      puVar6 = FUN_10024d20(local_23c,&DAT_100771a8);
      if (puVar6 == (undefined *)0x0) {
        return 0;
      }
      iVar9 = *(int *)(puVar6 + 0x10);
      psVar17 = local_10 + param_2 * 0x32;
      iVar10 = FUN_10019e80((uint)puVar6,(int)psVar17);
      if (iVar10 == 0) {
        return 0;
      }
      local_2c = 0;
      local_34 = 0;
      local_28 = (int)psVar17[0x30];
      local_3c = puVar6;
      local_24 = iVar9;
      FUN_10019700(psVar17,(int *)&local_3c,0);
      FUN_10019600((undefined4 *)(psVar17 + 0x26),(int *)&local_3c,4,(undefined4 *)0x1);
      FUN_100253d0(puVar6);
      *(int *)(pcVar4 + 0x6c) = *(int *)(pcVar4 + 0x6c) + *(int *)(psVar17 + 0x26);
      if (param_2 == 0) {
        local_10[0x2e] = 0;
        local_10[0x2f] = 0;
      }
      else {
        *(int *)(psVar17 + 0x2e) = *(int *)(psVar17 + -0xc) + *(int *)(psVar17 + -4);
      }
      param_2 = param_2 + 1;
    } while (param_2 < *(short *)(pcVar4 + 0xe94));
  }
  psVar17 = local_10;
  FUN_10019a70((undefined4 *)(pcVar4 + 0x30),*(size_t *)(pcVar4 + 0x6c));
  local_14 = 0;
  local_18 = (char *)0x0;
  param_2 = 0;
  if (0 < *(short *)(pcVar4 + 0xe94)) {
    local_c = (byte *)(pcVar4 + 0xe96);
    psVar15 = psVar17 + 0x29;
    do {
      FUN_10063ed6(local_23c,(byte *)s__s_s_idx_1007c5a4);
      local_3c = FUN_10024d20(local_23c,&DAT_100771a8);
      if (local_3c == (undefined *)0x0) {
        return 0;
      }
      iVar9 = *(int *)(local_3c + 0x10);
      local_2c = 0;
      local_34 = 0;
      sVar5 = psVar15[7];
      local_28 = (int)sVar5;
      local_24 = iVar9;
      local_1c = local_3c;
      iVar10 = FUN_10019700(psVar15 + -0x29,(int *)&local_3c,1);
      local_8 = (char *)(sVar5 + iVar10);
      piVar1 = (int *)(psVar15 + -3);
      iVar10 = FUN_10019600(piVar1,(int *)&local_3c,4,(undefined4 *)0x1);
      local_8 = (char *)((int)local_8 + iVar10);
      if ((char)psVar15[8] == '\0') {
        psVar15[-1] = (short)local_8;
        *psVar15 = (short)((iVar9 - (short)local_8) / *piVar1);
      }
      else {
        iVar9 = FUN_10019600((undefined4 *)psVar15,(int *)&local_3c,2,(undefined4 *)0x1);
        psVar15[-1] = (short)iVar9 + (short)local_8;
      }
      FUN_10019b10((int *)(psVar15 + -0x21),(int *)(pcVar4 + 0x30),local_14);
      iVar9 = FUN_100199f0((int)(psVar15 + -0x29),(int *)&local_3c);
      if (iVar9 == 0) {
        return 0;
      }
      iVar9 = 0;
      psVar15[-0x23] = (short)local_18;
      if (0 < *piVar1) {
        do {
          *(char *)(*(int *)(psVar15 + -5) + iVar9) =
               *(char *)(*(int *)(psVar15 + -5) + iVar9) + (char)local_18;
          iVar9 = iVar9 + 1;
        } while (iVar9 < *piVar1);
      }
      local_14 = local_14 + *piVar1;
      local_18 = (char *)((int)local_18 + CONCAT22((short)((uint)local_14 >> 0x10),psVar15[-0x29]));
      local_c = local_c + 0x200;
      *(undefined **)(psVar15 + 3) = local_1c;
      param_2 = param_2 + 1;
      psVar15 = psVar15 + 0x32;
      psVar17 = local_10;
    } while (param_2 < *(short *)(pcVar4 + 0xe94));
  }
  iVar9 = 0;
  pcVar4[0x20] = '\0';
  pcVar4[0x21] = '\0';
  iVar10 = 0;
  psVar15 = psVar17;
  if (0 < *(short *)(pcVar4 + 0xe94)) {
    do {
      *(short *)(pcVar4 + 0x20) = *(short *)(pcVar4 + 0x20) + *psVar15;
      iVar10 = iVar10 + 1;
      psVar15 = psVar15 + 0x32;
    } while (iVar10 < *(short *)(pcVar4 + 0xe94));
  }
  uVar11 = FUN_1001d9c0((int)*(short *)(pcVar4 + 0x20) << 2);
  *(undefined4 *)(pcVar4 + 0x24) = uVar11;
  uVar11 = FUN_1001d9c0((int)*(short *)(pcVar4 + 0x20));
  *(undefined4 *)(pcVar4 + 0x28) = uVar11;
  local_14 = 0;
  param_2 = 0;
  if (0 < *(short *)(pcVar4 + 0xe94)) {
    local_c = (byte *)(pcVar4 + 0x4a96);
    psVar15 = psVar17;
    do {
      uVar11 = FUN_10019cb0(local_c);
      local_8 = (char *)0x0;
      if (0 < *psVar15) {
        do {
          uVar13 = 0xffffffff;
          pcVar7 = *(char **)(*(int *)(psVar15 + 2) + (int)local_8 * 4);
          do {
            if (uVar13 == 0) break;
            uVar13 = uVar13 - 1;
            cVar2 = *pcVar7;
            pcVar7 = pcVar7 + 1;
          } while (cVar2 != '\0');
          local_1c = (undefined *)~uVar13;
          uVar12 = FUN_1001d9c0((size_t)local_1c);
          iVar9 = local_14;
          *(undefined4 *)(*(int *)(pcVar4 + 0x24) + local_14 * 4) = uVar12;
          FUN_10063f30(*(undefined4 **)(*(int *)(pcVar4 + 0x24) + local_14 * 4),
                       *(undefined4 **)(*(int *)(psVar15 + 2) + (int)local_8 * 4),(uint)local_1c);
          *(char *)(*(int *)(pcVar4 + 0x28) + iVar9) = (char)uVar11;
          iVar9 = iVar9 + 1;
          local_8 = local_8 + 1;
          local_14 = iVar9;
        } while ((int)local_8 < (int)*psVar15);
      }
      FUN_10019a20(psVar15);
      *(int *)(psVar15 + 2) = *(int *)(pcVar4 + 0x24) + (iVar9 - *psVar15) * 4;
      param_2 = param_2 + 1;
      *(int *)(psVar15 + 4) = (*(int *)(pcVar4 + 0x28) - (int)*psVar15) + iVar9;
      local_c = local_c + 0x14;
      psVar17 = local_10;
      psVar15 = psVar15 + 0x32;
    } while (param_2 < *(short *)(pcVar4 + 0xe94));
  }
  FUN_1001da30(local_38);
  iVar9 = 1;
  if (1 < *(short *)(pcVar4 + 0xe94)) {
    psVar15 = psVar17 + 99;
    do {
      if ((char)*psVar15 != (char)psVar17[0x31]) {
        return 0;
      }
      iVar9 = iVar9 + 1;
      psVar15 = psVar15 + 0x32;
    } while (iVar9 < *(short *)(pcVar4 + 0xe94));
  }
  sVar5 = psVar17[0x31];
  pcVar4[0x82] = (char)sVar5;
  if ((char)sVar5 == '\0') {
    pcVar4[0x74] = '\b';
    pcVar4[0x75] = '\0';
    pcVar4[0x76] = '\x04';
    pcVar4[0x77] = '\0';
  }
  sVar5 = psVar17[0x29];
  *(short *)(pcVar4 + 0x72) = sVar5;
  if (sVar5 < 0xc9) {
    iVar9 = *(int *)(pcVar4 + 0x6c);
    iVar10 = 0;
    if (0 < iVar9) {
      do {
        if (*(short *)(pcVar4 + 0x20) <= (short)(ushort)*(byte *)(*(int *)(pcVar4 + 0x68) + iVar10))
        {
          return 0;
        }
        iVar10 = iVar10 + 1;
      } while (iVar10 < iVar9);
    }
    iVar10 = 0;
    if (0 < iVar9) {
      iVar9 = 0;
      do {
        if (*(char *)((uint)*(byte *)(iVar10 + *(int *)(pcVar4 + 0x68)) + *(int *)(pcVar4 + 0x28))
            == '\f') {
          *(byte *)(*(int *)(pcVar4 + 100) + 6 + iVar9) =
               *(byte *)(*(int *)(pcVar4 + 100) + 6 + iVar9) | 0x20;
        }
        iVar10 = iVar10 + 1;
        iVar9 = iVar9 + 7;
      } while (iVar10 < *(int *)(pcVar4 + 0x6c));
    }
    return 1;
  }
  return 0;
}



/* ===== FUN_10019e80 ===== */
/* Entry: 10019e80 */

undefined4 __cdecl FUN_10019e80(uint param_1,int param_2)

{
  byte bVar1;
  uint uVar2;
  int iVar3;
  byte *pbVar4;
  byte *pbVar5;
  bool bVar6;
  byte local_104 [256];

  uVar2 = param_1;
  FUN_10025440(param_1,0,0,(char *)&param_1,1,1);
  FUN_10025440(uVar2,1,0,(char *)local_104,1,param_1 & 0xff);
  iVar3 = _strncmp((char *)local_104,&DAT_1007c59c,4);
  if (iVar3 == 0) {
    *(undefined1 *)(param_2 + 0x62) = 1;
    *(ushort *)(param_2 + 0x60) = (byte)param_1 + 1;
    pbVar5 = &DAT_1007c590;
    pbVar4 = local_104;
    do {
      bVar1 = *pbVar4;
      bVar6 = bVar1 < *pbVar5;
      if (bVar1 != *pbVar5) {
LAB_10019f16:
        iVar3 = (1 - (uint)bVar6) - (uint)(bVar6 != 0);
        goto LAB_10019f1b;
      }
      if (bVar1 == 0) break;
      bVar1 = pbVar4[1];
      bVar6 = bVar1 < pbVar5[1];
      if (bVar1 != pbVar5[1]) goto LAB_10019f16;
      pbVar4 = pbVar4 + 2;
      pbVar5 = pbVar5 + 2;
    } while (bVar1 != 0);
    iVar3 = 0;
LAB_10019f1b:
    if (iVar3 != 0) {
      return 0;
    }
  }
  else {
    *(undefined1 *)(param_2 + 0x62) = 0;
    *(undefined2 *)(param_2 + 0x60) = 0;
  }
  return 1;
}



/* ===== FUN_10019700 ===== */
/* Entry: 10019700 */

int __cdecl FUN_10019700(short *param_1,int *param_2,short param_3)

{
  short *psVar1;
  int iVar2;
  undefined4 uVar3;
  int iVar4;
  int iVar5;
  int iVar6;
  int iVar7;

  psVar1 = param_1;
  iVar2 = FUN_10019600((undefined4 *)param_1,param_2,2,(undefined4 *)0x1);
  if (param_3 == 0) {
    uVar3 = FUN_1001d9c0((int)*psVar1 << 2);
    *(undefined4 *)(psVar1 + 2) = uVar3;
    uVar3 = FUN_1001d9c0((int)*psVar1);
    *(undefined4 *)(psVar1 + 4) = uVar3;
  }
  iVar7 = 0;
  if (0 < *psVar1) {
    do {
      iVar4 = FUN_10019600((undefined4 *)((int)&param_1 + 2),param_2,2,(undefined4 *)0x1);
      if (param_3 == 0) {
        uVar3 = FUN_1001d9c0((int)param_1._2_2_ + 1);
        *(undefined4 *)(*(int *)(psVar1 + 2) + iVar7 * 4) = uVar3;
      }
      iVar5 = FUN_10019600(*(undefined4 **)(*(int *)(psVar1 + 2) + iVar7 * 4),param_2,1,
                           (undefined4 *)(int)param_1._2_2_);
      *(undefined1 *)(*(int *)(*(int *)(psVar1 + 2) + iVar7 * 4) + (int)param_1._2_2_) = 0;
      iVar6 = FUN_10019600((undefined4 *)(*(int *)(psVar1 + 4) + iVar7),param_2,1,(undefined4 *)0x1)
      ;
      iVar2 = iVar2 + iVar4 + iVar5 + iVar6;
      iVar7 = iVar7 + 1;
    } while (iVar7 < *psVar1);
  }
  return iVar2;
}



/* ===== FUN_10019600 ===== */
/* Entry: 10019600 */

int __cdecl FUN_10019600(undefined4 *param_1,int *param_2,uint param_3,undefined4 *param_4)

{
  undefined4 *puVar1;
  int iVar2;
  int iVar3;
  int iVar4;
  int iVar5;

  puVar1 = param_1;
  iVar5 = param_3 * (int)param_4;
  iVar2 = param_2[5];
  iVar4 = (param_2[4] - iVar2) + param_2[2];
  if (iVar4 < 1) {
    iVar2 = FUN_100195b0(param_2);
    if (iVar2 == 0) {
      return 0;
    }
    iVar2 = param_2[5];
    iVar4 = (param_2[4] - iVar2) + param_2[2];
  }
  if (iVar5 - iVar4 != 0 && iVar4 <= iVar5) {
    param_1 = (undefined4 *)0x0;
    param_4 = puVar1;
    while( true ) {
      iVar3 = iVar5 - (int)param_1;
      if (iVar4 < iVar5 - (int)param_1) {
        iVar3 = iVar4;
      }
      iVar4 = (iVar3 / (int)param_3) * param_3;
      FUN_100014f0(param_4,(undefined4 *)((iVar2 - param_2[4]) + param_2[1]),param_3,
                   iVar3 / (int)param_3);
      param_4 = (undefined4 *)((int)param_4 + iVar4);
      param_2[5] = param_2[5] + iVar4;
      param_1 = (undefined4 *)((int)param_1 + iVar4);
      if (iVar5 <= (int)param_1) break;
      iVar2 = FUN_100195b0(param_2);
      if (iVar2 == 0) {
        return 0;
      }
      iVar2 = param_2[5];
      iVar4 = (param_2[4] - iVar2) + param_2[2];
    }
    return iVar5;
  }
  FUN_100014f0(param_1,(undefined4 *)((param_2[1] - param_2[4]) + iVar2),param_3,(int)param_4);
  param_2[5] = param_2[5] + iVar5;
  return iVar5;
}



/* ===== FUN_100199f0 ===== */
/* Entry: 100199f0 */

void __cdecl FUN_100199f0(int param_1,int *param_2)

{
  if (*(char *)(param_1 + 0x62) == '\0') {
    FUN_100197d0(param_1,param_2);
    return;
  }
  FUN_10019940(param_1,param_2);
  return;
}



/* ===== FUN_10019b10 ===== */
/* Entry: 10019b10 */

undefined4 __cdecl FUN_10019b10(int *param_1,int *param_2,int param_3)

{
  int iVar1;
  int iVar2;
  int *piVar3;
  int *piVar4;

  iVar2 = param_3;
  iVar1 = param_3 * 4;
  *param_1 = *param_2 + iVar1;
  if (DAT_100a045d != '\0') {
    param_1[1] = param_2[1] + iVar1;
    param_1[2] = param_2[2] + iVar1;
  }
  iVar1 = param_3 * 2;
  param_3 = 3;
  piVar3 = param_1 + 3;
  piVar4 = param_2 + 9;
  do {
    *piVar3 = *(int *)(((int)param_2 - (int)param_1) + (int)piVar3) + iVar1;
    piVar3[3] = piVar4[-3] + iVar2;
    piVar3[6] = *piVar4 + iVar2;
    param_3 = param_3 + -1;
    piVar3 = piVar3 + 1;
    piVar4 = piVar4 + 1;
  } while (param_3 != 0);
  param_1[0xc] = param_2[0xc] + iVar2;
  param_1[0xd] = param_2[0xd] + iVar2 * 7;
  param_1[0xe] = param_2[0xe] + iVar2;
  return 1;
}



/* ===== FUN_10024d20 ===== */
/* Entry: 10024d20 */

undefined * __cdecl FUN_10024d20(undefined4 param_1,byte *param_2)

{
  byte bVar1;
  undefined *puVar2;
  undefined4 *puVar3;
  LPCSTR pCVar4;
  int iVar5;
  uint uVar6;
  uint uVar7;
  byte *pbVar8;
  byte *pbVar9;
  uint local_c;
  byte *local_8;

  puVar2 = (undefined *)FUN_1001d9c0(0x420);
  if (puVar2 != (undefined *)0x0) {
    uVar6 = 0xffffffff;
    local_8 = puVar2 + 0x14;
    pbVar8 = param_2;
    do {
      pbVar9 = pbVar8;
      if (uVar6 == 0) break;
      uVar6 = uVar6 - 1;
      pbVar9 = pbVar8 + 1;
      bVar1 = *pbVar8;
      pbVar8 = pbVar9;
    } while (bVar1 != 0);
    uVar6 = ~uVar6;
    pbVar8 = pbVar9 + -uVar6;
    pbVar9 = local_8;
    for (uVar7 = uVar6 >> 2; uVar7 != 0; uVar7 = uVar7 - 1) {
      *(undefined4 *)pbVar9 = *(undefined4 *)pbVar8;
      pbVar8 = pbVar8 + 4;
      pbVar9 = pbVar9 + 4;
    }
    for (uVar6 = uVar6 & 3; uVar6 != 0; uVar6 = uVar6 - 1) {
      *pbVar9 = *pbVar8;
      pbVar8 = pbVar8 + 1;
      pbVar9 = pbVar9 + 1;
    }
    puVar3 = FUN_10024e20(puVar2,&local_c,param_2);
    if (puVar3 == (undefined4 *)0x0) {
      pCVar4 = (LPCSTR)FUN_10025360((int)puVar2,param_1);
      if ((pCVar4 != (LPCSTR)0x0) && (*pCVar4 != '\0')) {
        iVar5 = FUN_10025370((int)puVar2,pCVar4,local_c,0x80);
        if (-1 < iVar5) {
          FUN_10025400((int)puVar2,0,0);
          FUN_10025400((int)puVar2,0,2);
          *(undefined4 *)(puVar2 + 0x10) = *(undefined4 *)(puVar2 + 0xc);
          FUN_10025400((int)puVar2,0,0);
          if (*param_2 == 0x72) {
            iVar5 = FUN_10024cd0((int)puVar2);
            if (iVar5 == 0) {
              DAT_100a74b8 = DAT_100a74b8 + *(int *)(puVar2 + 0x10);
            }
          }
          return puVar2;
        }
      }
    }
    FUN_1001da30(puVar2);
  }
  return (undefined *)0x0;
}



/* ===== FUN_10024e20 ===== */
/* Entry: 10024e20 */

undefined4 * __cdecl FUN_10024e20(undefined4 param_1,undefined4 *param_2,byte *param_3)

{
  byte bVar1;
  byte *pbVar2;
  int iVar3;
  byte *pbVar4;
  bool bVar5;

  bVar1 = *param_3;
  if (bVar1 == 0x61) {
    pbVar4 = &DAT_10077718;
    pbVar2 = param_3;
    do {
      bVar1 = *pbVar2;
      bVar5 = bVar1 < *pbVar4;
      if (bVar1 != *pbVar4) {
LAB_100251d5:
        iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
        goto LAB_100251da;
      }
      if (bVar1 == 0) break;
      bVar1 = pbVar2[1];
      bVar5 = bVar1 < pbVar4[1];
      if (bVar1 != pbVar4[1]) goto LAB_100251d5;
      pbVar2 = pbVar2 + 2;
      pbVar4 = pbVar4 + 2;
    } while (bVar1 != 0);
    iVar3 = 0;
LAB_100251da:
    if (iVar3 != 0) {
      pbVar4 = &DAT_10077694;
      pbVar2 = param_3;
      do {
        bVar1 = *pbVar2;
        bVar5 = bVar1 < *pbVar4;
        if (bVar1 != *pbVar4) {
LAB_10025211:
          iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
          goto LAB_10025216;
        }
        if (bVar1 == 0) break;
        bVar1 = pbVar2[1];
        bVar5 = bVar1 < pbVar4[1];
        if (bVar1 != pbVar4[1]) goto LAB_10025211;
        pbVar2 = pbVar2 + 2;
        pbVar4 = pbVar4 + 2;
      } while (bVar1 != 0);
      iVar3 = 0;
LAB_10025216:
      if (iVar3 != 0) {
        pbVar4 = &DAT_1007d6c0;
        pbVar2 = param_3;
        do {
          bVar1 = *pbVar2;
          bVar5 = bVar1 < *pbVar4;
          if (bVar1 != *pbVar4) {
LAB_1002524d:
            iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
            goto LAB_10025252;
          }
          if (bVar1 == 0) break;
          bVar1 = pbVar2[1];
          bVar5 = bVar1 < pbVar4[1];
          if (bVar1 != pbVar4[1]) goto LAB_1002524d;
          pbVar2 = pbVar2 + 2;
          pbVar4 = pbVar4 + 2;
        } while (bVar1 != 0);
        iVar3 = 0;
LAB_10025252:
        if (iVar3 != 0) {
          pbVar4 = &DAT_1007d6bc;
          pbVar2 = param_3;
          do {
            bVar1 = *pbVar2;
            bVar5 = bVar1 < *pbVar4;
            if (bVar1 != *pbVar4) {
LAB_10025289:
              iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
              goto LAB_1002528e;
            }
            if (bVar1 == 0) break;
            bVar1 = pbVar2[1];
            bVar5 = bVar1 < pbVar4[1];
            if (bVar1 != pbVar4[1]) goto LAB_10025289;
            pbVar2 = pbVar2 + 2;
            pbVar4 = pbVar4 + 2;
          } while (bVar1 != 0);
          iVar3 = 0;
LAB_1002528e:
          if (iVar3 != 0) {
            pbVar4 = &DAT_1007d6b8;
            pbVar2 = param_3;
            do {
              bVar1 = *pbVar2;
              bVar5 = bVar1 < *pbVar4;
              if (bVar1 != *pbVar4) {
LAB_100252c5:
                iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
                goto LAB_100252ca;
              }
              if (bVar1 == 0) break;
              bVar1 = pbVar2[1];
              bVar5 = bVar1 < pbVar4[1];
              if (bVar1 != pbVar4[1]) goto LAB_100252c5;
              pbVar2 = pbVar2 + 2;
              pbVar4 = pbVar4 + 2;
            } while (bVar1 != 0);
            iVar3 = 0;
LAB_100252ca:
            if (iVar3 == 0) {
              *param_2 = 0x8109;
              return (undefined4 *)0x0;
            }
            pbVar2 = &DAT_1007d6b4;
            do {
              bVar1 = *param_3;
              bVar5 = bVar1 < *pbVar2;
              if (bVar1 != *pbVar2) {
LAB_1002530d:
                iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
                goto LAB_10025312;
              }
              if (bVar1 == 0) break;
              bVar1 = param_3[1];
              bVar5 = bVar1 < pbVar2[1];
              if (bVar1 != pbVar2[1]) goto LAB_1002530d;
              param_3 = param_3 + 2;
              pbVar2 = pbVar2 + 2;
            } while (bVar1 != 0);
            iVar3 = 0;
LAB_10025312:
            if (iVar3 != 0) {
              return param_2;
            }
            *param_2 = 0x810a;
            return (undefined4 *)0x0;
          }
        }
        *param_2 = 0x410a;
        return (undefined4 *)0x1;
      }
    }
    *param_2 = 0x4109;
    return (undefined4 *)0x1;
  }
  if (bVar1 != 0x72) {
    if (bVar1 != 0x77) {
      *param_2 = 0x4000;
      return (undefined4 *)0x1;
    }
    pbVar4 = &DAT_1007d348;
    pbVar2 = param_3;
    do {
      bVar1 = *pbVar2;
      bVar5 = bVar1 < *pbVar4;
      if (bVar1 != *pbVar4) {
LAB_10024e81:
        iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
        goto LAB_10024e86;
      }
      if (bVar1 == 0) break;
      bVar1 = pbVar2[1];
      bVar5 = bVar1 < pbVar4[1];
      if (bVar1 != pbVar4[1]) goto LAB_10024e81;
      pbVar2 = pbVar2 + 2;
      pbVar4 = pbVar4 + 2;
    } while (bVar1 != 0);
    iVar3 = 0;
LAB_10024e86:
    if (iVar3 != 0) {
      pbVar4 = &DAT_1007d534;
      pbVar2 = param_3;
      do {
        bVar1 = *pbVar2;
        bVar5 = bVar1 < *pbVar4;
        if (bVar1 != *pbVar4) {
LAB_10024ebd:
          iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
          goto LAB_10024ec2;
        }
        if (bVar1 == 0) break;
        bVar1 = pbVar2[1];
        bVar5 = bVar1 < pbVar4[1];
        if (bVar1 != pbVar4[1]) goto LAB_10024ebd;
        pbVar2 = pbVar2 + 2;
        pbVar4 = pbVar4 + 2;
      } while (bVar1 != 0);
      iVar3 = 0;
LAB_10024ec2:
      if (iVar3 != 0) {
        pbVar4 = &DAT_1007d6d8;
        pbVar2 = param_3;
        do {
          bVar1 = *pbVar2;
          bVar5 = bVar1 < *pbVar4;
          if (bVar1 != *pbVar4) {
LAB_10024ef9:
            iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
            goto LAB_10024efe;
          }
          if (bVar1 == 0) break;
          bVar1 = pbVar2[1];
          bVar5 = bVar1 < pbVar4[1];
          if (bVar1 != pbVar4[1]) goto LAB_10024ef9;
          pbVar2 = pbVar2 + 2;
          pbVar4 = pbVar4 + 2;
        } while (bVar1 != 0);
        iVar3 = 0;
LAB_10024efe:
        if (iVar3 != 0) {
          pbVar4 = &DAT_1007d6d4;
          pbVar2 = param_3;
          do {
            bVar1 = *pbVar2;
            bVar5 = bVar1 < *pbVar4;
            if (bVar1 != *pbVar4) {
LAB_10024f35:
              iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
              goto LAB_10024f3a;
            }
            if (bVar1 == 0) break;
            bVar1 = pbVar2[1];
            bVar5 = bVar1 < pbVar4[1];
            if (bVar1 != pbVar4[1]) goto LAB_10024f35;
            pbVar2 = pbVar2 + 2;
            pbVar4 = pbVar4 + 2;
          } while (bVar1 != 0);
          iVar3 = 0;
LAB_10024f3a:
          if (iVar3 != 0) {
            pbVar4 = &DAT_10079290;
            pbVar2 = param_3;
            do {
              bVar1 = *pbVar2;
              bVar5 = bVar1 < *pbVar4;
              if (bVar1 != *pbVar4) {
LAB_10024f71:
                iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
                goto LAB_10024f76;
              }
              if (bVar1 == 0) break;
              bVar1 = pbVar2[1];
              bVar5 = bVar1 < pbVar4[1];
              if (bVar1 != pbVar4[1]) goto LAB_10024f71;
              pbVar2 = pbVar2 + 2;
              pbVar4 = pbVar4 + 2;
            } while (bVar1 != 0);
            iVar3 = 0;
LAB_10024f76:
            if (iVar3 == 0) {
              *param_2 = 0x8301;
              return (undefined4 *)0x0;
            }
            pbVar2 = &DAT_1007d60c;
            do {
              bVar1 = *param_3;
              bVar5 = bVar1 < *pbVar2;
              if (bVar1 != *pbVar2) {
LAB_10024fb9:
                iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
                goto LAB_10024fbe;
              }
              if (bVar1 == 0) break;
              bVar1 = param_3[1];
              bVar5 = bVar1 < pbVar2[1];
              if (bVar1 != pbVar2[1]) goto LAB_10024fb9;
              param_3 = param_3 + 2;
              pbVar2 = pbVar2 + 2;
            } while (bVar1 != 0);
            iVar3 = 0;
LAB_10024fbe:
            if (iVar3 != 0) {
              return param_2;
            }
            *param_2 = 0x8302;
            return (undefined4 *)0x0;
          }
        }
        *param_2 = 0x4302;
        return (undefined4 *)0x1;
      }
    }
    *param_2 = 0x4301;
    return (undefined4 *)0x1;
  }
  pbVar4 = &DAT_1007d358;
  pbVar2 = param_3;
  do {
    bVar1 = *pbVar2;
    bVar5 = bVar1 < *pbVar4;
    if (bVar1 != *pbVar4) {
LAB_1002502b:
      iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
      goto LAB_10025030;
    }
    if (bVar1 == 0) break;
    bVar1 = pbVar2[1];
    bVar5 = bVar1 < pbVar4[1];
    if (bVar1 != pbVar4[1]) goto LAB_1002502b;
    pbVar2 = pbVar2 + 2;
    pbVar4 = pbVar4 + 2;
  } while (bVar1 != 0);
  iVar3 = 0;
LAB_10025030:
  if (iVar3 != 0) {
    pbVar4 = &DAT_1007d6d0;
    pbVar2 = param_3;
    do {
      bVar1 = *pbVar2;
      bVar5 = bVar1 < *pbVar4;
      if (bVar1 != *pbVar4) {
LAB_10025067:
        iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
        goto LAB_1002506c;
      }
      if (bVar1 == 0) break;
      bVar1 = pbVar2[1];
      bVar5 = bVar1 < pbVar4[1];
      if (bVar1 != pbVar4[1]) goto LAB_10025067;
      pbVar2 = pbVar2 + 2;
      pbVar4 = pbVar4 + 2;
    } while (bVar1 != 0);
    iVar3 = 0;
LAB_1002506c:
    if (iVar3 != 0) {
      pbVar4 = &DAT_1007d6cc;
      pbVar2 = param_3;
      do {
        bVar1 = *pbVar2;
        bVar5 = bVar1 < *pbVar4;
        if (bVar1 != *pbVar4) {
LAB_100250a3:
          iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
          goto LAB_100250a8;
        }
        if (bVar1 == 0) break;
        bVar1 = pbVar2[1];
        bVar5 = bVar1 < pbVar4[1];
        if (bVar1 != pbVar4[1]) goto LAB_100250a3;
        pbVar2 = pbVar2 + 2;
        pbVar4 = pbVar4 + 2;
      } while (bVar1 != 0);
      iVar3 = 0;
LAB_100250a8:
      if (iVar3 != 0) {
        pbVar4 = &DAT_1007d6c8;
        pbVar2 = param_3;
        do {
          bVar1 = *pbVar2;
          bVar5 = bVar1 < *pbVar4;
          if (bVar1 != *pbVar4) {
LAB_100250df:
            iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
            goto LAB_100250e4;
          }
          if (bVar1 == 0) break;
          bVar1 = pbVar2[1];
          bVar5 = bVar1 < pbVar4[1];
          if (bVar1 != pbVar4[1]) goto LAB_100250df;
          pbVar2 = pbVar2 + 2;
          pbVar4 = pbVar4 + 2;
        } while (bVar1 != 0);
        iVar3 = 0;
LAB_100250e4:
        if (iVar3 != 0) {
          pbVar4 = &DAT_100771a8;
          pbVar2 = param_3;
          do {
            bVar1 = *pbVar2;
            bVar5 = bVar1 < *pbVar4;
            if (bVar1 != *pbVar4) {
LAB_1002511b:
              iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
              goto LAB_10025120;
            }
            if (bVar1 == 0) break;
            bVar1 = pbVar2[1];
            bVar5 = bVar1 < pbVar4[1];
            if (bVar1 != pbVar4[1]) goto LAB_1002511b;
            pbVar2 = pbVar2 + 2;
            pbVar4 = pbVar4 + 2;
          } while (bVar1 != 0);
          iVar3 = 0;
LAB_10025120:
          if (iVar3 == 0) {
            *param_2 = 0x8000;
            return (undefined4 *)0x0;
          }
          pbVar2 = &DAT_1007d6c4;
          do {
            bVar1 = *param_3;
            bVar5 = bVar1 < *pbVar2;
            if (bVar1 != *pbVar2) {
LAB_10025163:
              iVar3 = (1 - (uint)bVar5) - (uint)(bVar5 != 0);
              goto LAB_10025168;
            }
            if (bVar1 == 0) break;
            bVar1 = param_3[1];
            bVar5 = bVar1 < pbVar2[1];
            if (bVar1 != pbVar2[1]) goto LAB_10025163;
            param_3 = param_3 + 2;
            pbVar2 = pbVar2 + 2;
          } while (bVar1 != 0);
          iVar3 = 0;
LAB_10025168:
          if (iVar3 != 0) {
            return param_2;
          }
          *param_2 = 0x8002;
          return (undefined4 *)0x0;
        }
      }
      *param_2 = 0x4002;
      return (undefined4 *)0x1;
    }
  }
  *param_2 = 0x4000;
  return (undefined4 *)0x1;
}



/* ===== FUN_10025440 ===== */
/* Entry: 10025440 */

uint __cdecl
FUN_10025440(int param_1,int param_2,DWORD param_3,char *param_4,uint param_5,int param_6)

{
  uint uVar1;

  EnterCriticalSection((LPCRITICAL_SECTION)&DAT_100a8400);
  if (param_3 == 0) {
    if (param_2 == *(int *)(param_1 + 0xc)) goto LAB_10025477;
    param_3 = 0;
  }
  FUN_10025400(param_1,param_2,param_3);
LAB_10025477:
  uVar1 = FUN_100254a0(param_4,param_5,param_6,param_1);
  LeaveCriticalSection((LPCRITICAL_SECTION)&DAT_100a8400);
  return uVar1;
}



/* ===== FUN_10063ed6 ===== */
/* Entry: 10063ed6 */

int __cdecl FUN_10063ed6(undefined1 *param_1,byte *param_2)

{
  int iVar1;
  undefined1 *local_24;
  int local_20;
  undefined1 *local_1c;
  undefined4 local_18;

  local_1c = param_1;
  local_24 = param_1;
  local_18 = 0x42;
  local_20 = 0x7fffffff;
  iVar1 = FUN_100653b4((int *)&local_24,param_2,(undefined4 *)&stack0x0000000c);
  local_20 = local_20 + -1;
  if (local_20 < 0) {
    FUN_1006529c(0,(int *)&local_24);
  }
  else {
    *local_24 = 0;
  }
  return iVar1;
}
