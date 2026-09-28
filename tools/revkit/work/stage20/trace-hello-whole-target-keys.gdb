set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10018770
commands
  silent
  set $caller = *(unsigned int *)$esp
  if $caller >= 0x10024060 && $caller < 0x100242a0
    set $key = *(unsigned int *)($esp + 4)
    set $slot = *(unsigned int *)($esp + 8)
    set $mode = *(unsigned short *)($esp + 20)
    printf "HELLO_WHOLE_TARGET slot=%u mode=%u caller=%#x:", $slot, $mode, $caller
    x/7ub $key
  end
  continue
end

break *0x408187
commands
  silent
  printf "HELLO_WHOLE_TARGET_TRACE_READY\n"
  continue
end

continue
