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
  printf "PLACE_CALL count=%u arg=%d", $count, $arg
  set $i = 0
  while $i < $count && $i < 5
    set $entry = $components + $i * 0x140
    printf " component=%u class=%#x surface=", $i, *(unsigned char *)($entry + 0x1a)
    x/s ($entry + 0x46)
    set $i = $i + 1
  end
  tbreak *$ret
  commands
    silent
    printf "PLACE_CALL_RETURN eax=%#x\n", $eax
    continue
  end
  continue
end

continue
