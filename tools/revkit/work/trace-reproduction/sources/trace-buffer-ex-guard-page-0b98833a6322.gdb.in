set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV stop print pass
set $guard_started = 0

break *0x1001da50
commands 1
  silent
  set $text = *(unsigned int *)($esp + 8)
  disable 1
  set $allocation = ((char *(*)(void *, unsigned int, unsigned int, unsigned int))0x7b68d8a0)((char *)0, 0x3000, 0x3000, 0x04)
  set $old_protect = (unsigned int *)malloc(4)
  set $protected = ((int (*)(void *, unsigned int, unsigned int, unsigned int *))0x7b68df20)($allocation + 0x1000, 0x1000, 0x01, $old_protect)
  set $buffer = $allocation + 0xfff
  set {unsigned char}$buffer = 0xa5
  printf "GUARD_SETUP allocation=%#x protect_ok=%d accessible_bytes=1 guard=%#x text=%s\n", $allocation, $protected, $buffer + 1, $text
  set $guard_started = 1
  set $length = (int *)malloc(4)
  set $extra_a = (unsigned int *)malloc(4)
  set $extra_b = (unsigned int *)malloc(4)
  set *$length = 1
  set *$extra_a = 0x55555555
  set *$extra_b = 0x66666666
  set $result = ((int (*)(int, char *, char *, int *, int, int, unsigned int, unsigned int *, unsigned int *, int, int, int, int, int, int))0x1001ddf0)(@SELECTOR@, (char *)$text, $buffer, $length, 0, 0, 1, $extra_a, $extra_b, -1, -1, -1, -1, 0, 0)
  printf "GUARD_RETURN selector=@SELECTOR@ result=%d length=%d byte0=%#x\n", $result, *$length, *(unsigned char *)$buffer
  continue
end

catch signal SIGSEGV
commands 2
  silent
  if $guard_started
    printf "GUARD_SIGSEGV eip=%#x buffer=%#x\n", $eip, $buffer
    info registers eip esp
    x/8i $eip-8
  end
  continue
end

catch signal SIGTRAP
commands 3
  silent
  if $guard_started
    printf "GUARD_SIGTRAP eip=%#x buffer=%#x\n", $eip, $buffer
    info registers eip esp
    x/8i $eip-8
  end
  continue
end

continue
