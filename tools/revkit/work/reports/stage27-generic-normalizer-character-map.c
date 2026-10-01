===== 0x10002d70 =====
Function: FUN_10002d70 @ 10002d70

uint __cdecl FUN_10002d70(char param_1)

{
  undefined4 in_EAX;
  
  if (('@' < param_1) && (param_1 < '[')) {
    return (int)param_1 - 0x40;
  }
  if (('`' < param_1) && (param_1 < '{')) {
    return CONCAT22(param_1 >> 7,*(undefined2 *)(&DAT_1007e688 + param_1 * 2)) - 0x40;
  }
  return CONCAT22((short)((uint)in_EAX >> 0x10),(param_1 != '\'') - 1) & 0xffff001c;
}



===== 0x100026f0 =====
Function: FUN_100026f0 @ 100026f0

undefined4 __cdecl FUN_100026f0(char *param_1,int param_2,int param_3,ushort param_4)

{
  char cVar1;
  uint in_EAX;
  int iVar2;
  uint uVar3;
  uint uVar4;
  char *pcVar5;
  char *pcVar6;
  char local_44 [31];
  char acStack_25 [33];
  
  if (((param_1 == (char *)0x0) || (*param_1 == '\0')) || (8 < param_2)) goto LAB_10002c4a;
  in_EAX = FUN_10002680(param_1,param_4);
  if (in_EAX < 0x80000000) goto LAB_10002c28;
  if ((param_2 < 3) || (param_4 != 1)) {
    if (param_4 == 2) {
      if (param_2 == 4) {
        _strncpy(acStack_25 + 1,param_1,2);
        acStack_25[3] = 0;
        pcVar5 = param_1 + 2;
      }
      else if (param_2 == 5) {
        _strncpy(acStack_25 + 1,param_1,2);
        uVar3 = 0xffffffff;
        pcVar5 = param_1 + 2;
        do {
          pcVar6 = pcVar5;
          if (uVar3 == 0) break;
          uVar3 = uVar3 - 1;
          pcVar6 = pcVar5 + 1;
          cVar1 = *pcVar5;
          pcVar5 = pcVar6;
        } while (cVar1 != '\0');
        uVar3 = ~uVar3;
        acStack_25[3] = 0;
        pcVar5 = pcVar6 + -uVar3;
        pcVar6 = local_44;
        for (uVar4 = uVar3 >> 2; uVar4 != 0; uVar4 = uVar4 - 1) {
          *(undefined4 *)pcVar6 = *(undefined4 *)pcVar5;
          pcVar5 = pcVar5 + 4;
          pcVar6 = pcVar6 + 4;
        }
        for (uVar3 = uVar3 & 3; uVar3 != 0; uVar3 = uVar3 - 1) {
          *pcVar6 = *pcVar5;
          pcVar5 = pcVar5 + 1;
          pcVar6 = pcVar6 + 1;
        }
        iVar2 = FUN_10002680(acStack_25 + 1,1);
        if ((-1 < iVar2) && (in_EAX = FUN_10002680(local_44,4), in_EAX < 0x80000000))
        goto LAB_10002c28;
        _strncpy(acStack_25 + 1,param_1,3);
        acStack_25[4] = 0;
        pcVar5 = param_1 + 3;
      }
      else if (param_2 == 6) {
        _strncpy(acStack_25 + 1,param_1,3);
        uVar3 = 0xffffffff;
        pcVar5 = param_1 + 3;
        do {
          pcVar6 = pcVar5;
          if (uVar3 == 0) break;
          uVar3 = uVar3 - 1;
          pcVar6 = pcVar5 + 1;
          cVar1 = *pcVar5;
          pcVar5 = pcVar6;
        } while (cVar1 != '\0');
        uVar3 = ~uVar3;
        acStack_25[4] = 0;
        pcVar5 = pcVar6 + -uVar3;
        pcVar6 = local_44;
        for (uVar4 = uVar3 >> 2; uVar4 != 0; uVar4 = uVar4 - 1) {
          *(undefined4 *)pcVar6 = *(undefined4 *)pcVar5;
          pcVar5 = pcVar5 + 4;
          pcVar6 = pcVar6 + 4;
        }
        for (uVar3 = uVar3 & 3; uVar3 != 0; uVar3 = uVar3 - 1) {
          *pcVar6 = *pcVar5;
          pcVar5 = pcVar5 + 1;
          pcVar6 = pcVar6 + 1;
        }
        iVar2 = FUN_10002680(acStack_25 + 1,1);
        if ((-1 < iVar2) && (in_EAX = FUN_10002680(local_44,4), in_EAX < 0x80000000))
        goto LAB_10002c28;
        _strncpy(acStack_25 + 1,param_1,2);
        uVar3 = 0xffffffff;
        pcVar5 = param_1 + 2;
        do {
          pcVar6 = pcVar5;
          if (uVar3 == 0) break;
          uVar3 = uVar3 - 1;
          pcVar6 = pcVar5 + 1;
          cVar1 = *pcVar5;
          pcVar5 = pcVar6;
        } while (cVar1 != '\0');
        uVar3 = ~uVar3;
        acStack_25[3] = 0;
        pcVar5 = pcVar6 + -uVar3;
        pcVar6 = local_44;
        for (uVar4 = uVar3 >> 2; uVar4 != 0; uVar4 = uVar4 - 1) {
          *(undefined4 *)pcVar6 = *(undefined4 *)pcVar5;
          pcVar5 = pcVar5 + 4;
          pcVar6 = pcVar6 + 4;
        }
        for (uVar3 = uVar3 & 3; uVar3 != 0; uVar3 = uVar3 - 1) {
          *pcVar6 = *pcVar5;
          pcVar5 = pcVar5 + 1;
          pcVar6 = pcVar6 + 1;
        }
        iVar2 = FUN_10002680(acStack_25 + 1,1);
        if ((-1 < iVar2) && (in_EAX = FUN_10002680(local_44,4), in_EAX < 0x80000000))
        goto LAB_10002c28;
        _strncpy(acStack_25 + 1,param_1,4);
        acStack_25[5] = 0;
        pcVar5 = param_1 + 4;
      }
      else if (param_2 == 7) {
        _strncpy(acStack_25 + 1,param_1,2);
        uVar3 = 0xffffffff;
        pcVar5 = param_1 + 2;
        do {
          pcVar6 = pcVar5;
          if (uVar3 == 0) break;
          uVar3 = uVar3 - 1;
          pcVar6 = pcVar5 + 1;
          cVar1 = *pcVar5;
          pcVar5 = pcVar6;
        } while (cVar1 != '\0');
        uVar3 = ~uVar3;
        acStack_25[3] = 0;
        pcVar5 = pcVar6 + -uVar3;
        pcVar6 = local_44;
        for (uVar4 = uVar3 >> 2; uVar4 != 0; uVar4 = uVar4 - 1) {
          *(undefined4 *)pcVar6 = *(undefined4 *)pcVar5;
          pcVar5 = pcVar5 + 4;
          pcVar6 = pcVar6 + 4;
        }
        for (uVar3 = uVar3 & 3; uVar3 != 0; uVar3 = uVar3 - 1) {
          *pcVar6 = *pcVar5;
          pcVar5 = pcVar5 + 1;
          pcVar6 = pcVar6 + 1;
        }
        iVar2 = FUN_10002680(acStack_25 + 1,1);
        if ((-1 < iVar2) && (in_EAX = FUN_10002680(local_44,4), in_EAX < 0x80000000))
        goto LAB_10002c28;
        _strncpy(acStack_25 + 1,param_1,3);
        uVar3 = 0xffffffff;
        pcVar5 = param_1 + 3;
        do {
          pcVar6 = pcVar5;
          if (uVar3 == 0) break;
          uVar3 = uVar3 - 1;
          pcVar6 = pcVar5 + 1;
          cVar1 = *pcVar5;
          pcVar5 = pcVar6;
        } while (cVar1 != '\0');
        uVar3 = ~uVar3;
        acStack_25[4] = 0;
        pcVar5 = pcVar6 + -uVar3;
        pcVar6 = local_44;
        for (uVar4 = uVar3 >> 2; uVar4 != 0; uVar4 = uVar4 - 1) {
          *(undefined4 *)pcVar6 = *(undefined4 *)pcVar5;
          pcVar5 = pcVar5 + 4;
          pcVar6 = pcVar6 + 4;
        }
        for (uVar3 = uVar3 & 3; uVar3 != 0; uVar3 = uVar3 - 1) {
          *pcVar6 = *pcVar5;
          pcVar5 = pcVar5 + 1;
          pcVar6 = pcVar6 + 1;
        }
        iVar2 = FUN_10002680(acStack_25 + 1,1);
        if ((-1 < iVar2) && (in_EAX = FUN_10002680(local_44,4), in_EAX < 0x80000000))
        goto LAB_10002c28;
        _strncpy(acStack_25 + 1,param_1,4);
        uVar3 = 0xffffffff;
        pcVar5 = param_1 + 4;
        do {
          pcVar6 = pcVar5;
          if (uVar3 == 0) break;
          uVar3 = uVar3 - 1;
          pcVar6 = pcVar5 + 1;
          cVar1 = *pcVar5;
          pcVar5 = pcVar6;
        } while (cVar1 != '\0');
        uVar3 = ~uVar3;
        acStack_25[5] = 0;
        pcVar5 = pcVar6 + -uVar3;
        pcVar6 = local_44;
        for (uVar4 = uVar3 >> 2; uVar4 != 0; uVar4 = uVar4 - 1) {
          *(undefined4 *)pcVar6 = *(undefined4 *)pcVar5;
          pcVar5 = pcVar5 + 4;
          pcVar6 = pcVar6 + 4;
        }
        for (uVar3 = uVar3 & 3; uVar3 != 0; uVar3 = uVar3 - 1) {
          *pcVar6 = *pcVar5;
          pcVar5 = pcVar5 + 1;
          pcVar6 = pcVar6 + 1;
        }
        iVar2 = FUN_10002680(acStack_25 + 1,1);
        if ((-1 < iVar2) && (in_EAX = FUN_10002680(local_44,4), in_EAX < 0x80000000))
        goto LAB_10002c28;
        _strncpy(acStack_25 + 1,param_1,5);
        acStack_25[6] = 0;
        pcVar5 = param_1 + 5;
      }
      else {
        if (param_2 != 8) goto LAB_10002c39;
        _strncpy(acStack_25 + 1,param_1,2);
        uVar3 = 0xffffffff;
        pcVar5 = param_1 + 2;
        do {
          pcVar6 = pcVar5;
          if (uVar3 == 0) break;
          uVar3 = uVar3 - 1;
          pcVar6 = pcVar5 + 1;
          cVar1 = *pcVar5;
          pcVar5 = pcVar6;
        } while (cVar1 != '\0');
        uVar3 = ~uVar3;
        acStack_25[3] = 0;
        pcVar5 = pcVar6 + -uVar3;
        pcVar6 = local_44;
        for (uVar4 = uVar3 >> 2; uVar4 != 0; uVar4 = uVar4 - 1) {
          *(undefined4 *)pcVar6 = *(undefined4 *)pcVar5;
          pcVar5 = pcVar5 + 4;
          pcVar6 = pcVar6 + 4;
        }
        for (uVar3 = uVar3 & 3; uVar3 != 0; uVar3 = uVar3 - 1) {
          *pcVar6 = *pcVar5;
          pcVar5 = pcVar5 + 1;
          pcVar6 = pcVar6 + 1;
        }
        iVar2 = FUN_10002680(acStack_25 + 1,1);
        if ((-1 < iVar2) && (iVar2 = FUN_10002680(local_44,4), -1 < iVar2)) {
          return CONCAT22((short)((uint)iVar2 >> 0x10),1);
        }
        _strncpy(acStack_25 + 1,param_1,3);
        uVar3 = 0xffffffff;
        pcVar5 = param_1 + 3;
        do {
          pcVar6 = pcVar5;
          if (uVar3 == 0) break;
          uVar3 = uVar3 - 1;
          pcVar6 = pcVar5 + 1;
          cVar1 = *pcVar5;
          pcVar5 = pcVar6;
        } while (cVar1 != '\0');
        uVar3 = ~uVar3;
        acStack_25[4] = 0;
        pcVar5 = pcVar6 + -uVar3;
        pcVar6 = local_44;
        for (uVar4 = uVar3 >> 2; uVar4 != 0; uVar4 = uVar4 - 1) {
          *(undefined4 *)pcVar6 = *(undefined4 *)pcVar5;
          pcVar5 = pcVar5 + 4;
          pcVar6 = pcVar6 + 4;
        }
        for (uVar3 = uVar3 & 3; uVar3 != 0; uVar3 = uVar3 - 1) {
          *pcVar6 = *pcVar5;
          pcVar5 = pcVar5 + 1;
          pcVar6 = pcVar6 + 1;
        }
        iVar2 = FUN_10002680(acStack_25 + 1,1);
        if ((-1 < iVar2) && (iVar2 = FUN_10002680(local_44,4), -1 < iVar2)) {
          return CONCAT22((short)((uint)iVar2 >> 0x10),1);
        }
        _strncpy(acStack_25 + 1,param_1,4);
        uVar3 = 0xffffffff;
        pcVar5 = param_1 + 4;
        do {
          pcVar6 = pcVar5;
          if (uVar3 == 0) break;
          uVar3 = uVar3 - 1;
          pcVar6 = pcVar5 + 1;
          cVar1 = *pcVar5;
          pcVar5 = pcVar6;
        } while (cVar1 != '\0');
        uVar3 = ~uVar3;
        acStack_25[5] = 0;
        pcVar5 = pcVar6 + -uVar3;
        pcVar6 = local_44;
        for (uVar4 = uVar3 >> 2; uVar4 != 0; uVar4 = uVar4 - 1) {
          *(undefined4 *)pcVar6 = *(undefined4 *)pcVar5;
          pcVar5 = pcVar5 + 4;
          pcVar6 = pcVar6 + 4;
        }
        for (uVar3 = uVar3 & 3; uVar3 != 0; uVar3 = uVar3 - 1) {
          *pcVar6 = *pcVar5;
          pcVar5 = pcVar5 + 1;
          pcVar6 = pcVar6 + 1;
        }
        iVar2 = FUN_10002680(acStack_25 + 1,1);
        if ((-1 < iVar2) && (iVar2 = FUN_10002680(local_44,4), -1 < iVar2)) {
          return CONCAT22((short)((uint)iVar2 >> 0x10),1);
        }
        _strncpy(acStack_25 + 1,param_1,5);
        uVar3 = 0xffffffff;
        pcVar5 = param_1 + 5;
        do {
          pcVar6 = pcVar5;
          if (uVar3 == 0) break;
          uVar3 = uVar3 - 1;
          pcVar6 = pcVar5 + 1;
          cVar1 = *pcVar5;
          pcVar5 = pcVar6;
        } while (cVar1 != '\0');
        uVar3 = ~uVar3;
        acStack_25[6] = 0;
        pcVar5 = pcVar6 + -uVar3;
        pcVar6 = local_44;
        for (uVar4 = uVar3 >> 2; uVar4 != 0; uVar4 = uVar4 - 1) {
          *(undefined4 *)pcVar6 = *(undefined4 *)pcVar5;
          pcVar5 = pcVar5 + 4;
          pcVar6 = pcVar6 + 4;
        }
        for (uVar3 = uVar3 & 3; uVar3 != 0; uVar3 = uVar3 - 1) {
          *pcVar6 = *pcVar5;
          pcVar5 = pcVar5 + 1;
          pcVar6 = pcVar6 + 1;
        }
        iVar2 = FUN_10002680(acStack_25 + 1,1);
        if ((-1 < iVar2) && (iVar2 = FUN_10002680(local_44,4), -1 < iVar2)) {
          return CONCAT22((short)((uint)iVar2 >> 0x10),1);
        }
        _strncpy(acStack_25 + 1,param_1,6);
        acStack_25[7] = 0;
        pcVar5 = param_1 + 6;
      }
      uVar3 = 0xffffffff;
      do {
        pcVar6 = pcVar5;
        if (uVar3 == 0) break;
        uVar3 = uVar3 - 1;
        pcVar6 = pcVar5 + 1;
        cVar1 = *pcVar5;
        pcVar5 = pcVar6;
      } while (cVar1 != '\0');
      uVar3 = ~uVar3;
      pcVar5 = pcVar6 + -uVar3;
      pcVar6 = local_44;
      for (uVar4 = uVar3 >> 2; uVar4 != 0; uVar4 = uVar4 - 1) {
        *(undefined4 *)pcVar6 = *(undefined4 *)pcVar5;
        pcVar5 = pcVar5 + 4;
        pcVar6 = pcVar6 + 4;
      }
      for (uVar3 = uVar3 & 3; uVar3 != 0; uVar3 = uVar3 - 1) {
        *pcVar6 = *pcVar5;
        pcVar5 = pcVar5 + 1;
        pcVar6 = pcVar6 + 1;
      }
      in_EAX = FUN_10002680(acStack_25 + 1,1);
      if ((-1 < (int)in_EAX) && (in_EAX = FUN_10002680(local_44,4), -1 < (int)in_EAX)) {
LAB_10002c28:
        return CONCAT22((short)(in_EAX >> 0x10),1);
      }
    }
  }
  else {
    in_EAX = (uint)param_1[param_2 + -1];
    if (*(short *)(&DAT_1007e388 + in_EAX * 2) == 0x73) {
      _strncpy(acStack_25 + 1,param_1,param_2 - 1);
      acStack_25[param_2] = '\0';
      in_EAX = FUN_10002680(acStack_25 + 1,1);
      if (-1 < (int)in_EAX) {
        return CONCAT22((short)(in_EAX >> 0x10),1);
      }
    }
  }
LAB_10002c39:
  if (6 < param_3) {
    if (param_2 < 5) {
      if (param_4 != 8) {
LAB_10002c5a:
        return CONCAT22((short)(in_EAX >> 0x10),1);
      }
    }
    else if (param_4 == 2) goto LAB_10002c5a;
  }
LAB_10002c4a:
  return in_EAX & 0xffff0000;
}



