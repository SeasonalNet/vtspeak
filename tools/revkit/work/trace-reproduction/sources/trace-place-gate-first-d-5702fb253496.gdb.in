set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

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
