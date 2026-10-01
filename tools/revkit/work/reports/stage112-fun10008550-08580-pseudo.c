===== 0x10008550 =====
Function: FUN_10008550 @ 10008550

bool __cdecl FUN_10008550(char *param_1)

{
  uint uVar1;
  
  uVar1 = FUN_100560a0(&PTR_s_billion_10078768,param_1,DAT_1007883c,0x49);
  return uVar1 < 0x80000000;
}



===== 0x10008580 =====
Function: FUN_10008580 @ 10008580

undefined4 __cdecl FUN_10008580(int param_1,int param_2)

{
  byte *pbVar1;
  uint in_EAX;
  int iVar2;
  undefined4 uVar3;
  int iVar4;
  byte *pbVar5;
  
  if (param_2 == 0) {
    return CONCAT22((short)(in_EAX >> 0x10),1);
  }
  if (1 < param_2) {
    iVar4 = param_2 * 0x70 + param_1;
    pbVar1 = (byte *)(iVar4 + -0xd5);
    iVar2 = FUN_1001c2c0(pbVar1,&DAT_10077750);
    if (iVar2 == 0) {
      return 1;
    }
    pbVar5 = (byte *)(iVar4 + -0x65);
    uVar3 = FUN_10008440(pbVar1,pbVar5);
    if ((short)uVar3 != 0) {
      return CONCAT22((short)((uint)uVar3 >> 0x10),1);
    }
    uVar3 = FUN_10008440((byte *)(param_1 + 0xb),(byte *)(param_1 + 0x7b));
    if ((short)uVar3 != 0) {
      return CONCAT22((short)((uint)uVar3 >> 0x10),1);
    }
    in_EAX = FUN_10008440(pbVar1,&DAT_1009f948);
    if ((short)in_EAX != 0) {
      iVar4 = FUN_1001c2c0(pbVar5,&DAT_10077568);
      if (iVar4 != 0) {
        in_EAX = FUN_1001c2c0(pbVar5,&DAT_10077608);
        if (in_EAX != 0) goto LAB_1000864e;
      }
      return 1;
    }
  }
LAB_1000864e:
  if (0 < param_2) {
    pbVar1 = (byte *)(param_2 * 0x70 + -0x65 + param_1);
    iVar4 = FUN_1001c2c0(pbVar1,&DAT_1007786c);
    if (iVar4 == 0) {
      return 1;
    }
    iVar4 = FUN_1001c2c0(pbVar1,(byte *)&DAT_10077380);
    if (iVar4 == 0) {
      return 1;
    }
    in_EAX = FUN_10008440(pbVar1,&DAT_1009f948);
    if ((short)in_EAX != 0) {
      return CONCAT22((short)(in_EAX >> 0x10),1);
    }
  }
  return in_EAX & 0xffff0000;
}



