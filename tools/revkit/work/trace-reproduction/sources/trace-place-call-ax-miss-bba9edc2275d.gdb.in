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
  if $selector == 0x41
    tbreak *$ret
    commands
      silent
      if $eax != 0 && *(unsigned char *)$out == 0x58 && *(unsigned char *)($out + 1) == 0
        printf "AX_MISS key="
        x/s $key
        set {unsigned char}($out) = 0
        set $eax = 0
      end
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
