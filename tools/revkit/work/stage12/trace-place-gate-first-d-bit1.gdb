set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1003a7b0
commands
  silent
  set $out = *(unsigned int *)($esp + 4)
  set $key = *(unsigned int *)($esp + 8)
  set $selector = *(unsigned int *)($esp + 12)
  set $ret = *(unsigned int *)$esp
  set $length = 0
  while *(unsigned char *)($key + $length) != 0 && $length < 64
    set $length = $length + 1
  end
  if $selector == 0x42 && $length == 13 && *(unsigned int *)$key == 0x6f626241
    tbreak *$ret
    commands
      silent
      printf "PLACE_GATE_BIT_OVERRIDE key="
      x/s $key
      set {unsigned char}($out + 1) = 0x31
      printf "PLACE_GATE_BITS output="
      x/s $out
      continue
    end
  end
  continue
end

break *0x10034180
commands
  silent
  set $context = *(unsigned int *)($esp + 4)
  set $arg = *(int *)($esp + 8)
  set $components = *(unsigned int *)($context + 0x39e4)
  set $count = *(unsigned int *)$components
  set $ret = *(unsigned int *)$esp
  if $count == 2
    set {unsigned char}($components + 0x1a) = 0x64
    printf "PLACE_GATE_SET component=0 marker=d arg=%d\n", $arg
  end
  printf "PLACE_CALL count=%u arg=%d\n", $count, $arg
  tbreak *$ret
  commands
    silent
    printf "PLACE_CALL_RETURN eax=%#x\n", $eax
    continue
  end
  continue
end

continue
