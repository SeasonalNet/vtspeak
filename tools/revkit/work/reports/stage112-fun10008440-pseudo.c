===== 0x10008440 =====
Function: FUN_10008440 @ 10008440

undefined4 __cdecl FUN_10008440(byte *param_1,byte *param_2)

{
  uint in_EAX;
  int iVar1;
  
  if ((param_2 == (byte *)0x0) || (*param_2 == 0)) {
    in_EAX = FUN_100560a0(&PTR_DAT_10078044,(char *)param_1,DAT_100780c4,0x49);
    if (-1 < (int)in_EAX) {
      return CONCAT22((short)(in_EAX >> 0x10),1);
    }
  }
  if (param_2 == (byte *)0x0) goto LAB_10008543;
  iVar1 = FUN_1001c2c0(param_1,&DAT_1007776c);
  if (iVar1 == 0) {
LAB_100084b5:
    iVar1 = FUN_1001c2c0(param_2,&DAT_10077568);
    if (iVar1 == 0) {
      return 1;
    }
  }
  else {
    iVar1 = FUN_1001c2c0(param_1,&DAT_10077768);
    if (iVar1 == 0) goto LAB_100084b5;
    iVar1 = FUN_1001c2c0(param_1,&DAT_10077764);
    if (iVar1 == 0) goto LAB_100084b5;
  }
  iVar1 = FUN_1001c2c0(param_1,(byte *)s_ought_1007775c);
  if (iVar1 != 0) {
    iVar1 = FUN_1001c2c0(param_1,&DAT_10077394);
    if (iVar1 != 0) {
      iVar1 = FUN_1001c2c0(param_1,&DAT_10077390);
      if (iVar1 != 0) {
        iVar1 = FUN_1001c2c0(param_1,&DAT_1007738c);
        if (iVar1 != 0) {
          in_EAX = FUN_1001c2c0(param_1,&DAT_10077754);
          if (in_EAX != 0) goto LAB_10008543;
        }
      }
    }
  }
  in_EAX = FUN_1001c2c0(param_2,(byte *)&DAT_10077380);
  if (in_EAX == 0) {
    return 1;
  }
LAB_10008543:
  return in_EAX & 0xffff0000;
}



