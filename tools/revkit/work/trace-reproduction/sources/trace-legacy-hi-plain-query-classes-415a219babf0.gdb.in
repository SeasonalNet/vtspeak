set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/hi-plain-2006-query-key-trace-v2.log
set logging overwrite on
set logging enabled on

hbreak *0x1001dfb0
commands
  silent
  set $slot = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  printf "HI_PLAIN_2006_SPLIT_ENTER slot=%u\n", $slot
  thbreak *$return
  commands
    silent
    printf "HI_PLAIN_2006_SPLIT_RETURN slot=%u accepted=%u\n", $slot, $al
    continue
  end
  continue
end

hbreak *0x1001da80
commands
  silent
  set $caller = *(unsigned int *)$esp
  if $caller >= 0x1001dfb0 && $caller < 0x1001e110
    set $key = *(unsigned int *)($esp + 4)
    set $slot = *(unsigned int *)($esp + 8)
    set $state = *(unsigned int *)($esp + 12)
    set $mode = *(unsigned short *)($esp + 20)
    set $return = $caller
    printf "HI_PLAIN_2006_CLASS_LOOKUP slot=%u mode=%u key:", $slot, $mode
    x/5ub $key
    thbreak *$return
    commands
      silent
      printf "HI_PLAIN_2006_CLASS_RETURN slot=%u mode=%u appended=%u total=%u\n", $slot, $mode, $eax, *(unsigned short *)($state + 0xf8e64)
      continue
    end
  end
  continue
end

hbreak *0x1001ed7f
commands
  silent
  printf "HI_PLAIN_2006_SELECTED_UNIT index=%d id=%u\n", *(short *)$edi, $edx
  continue
end

hbreak *0x408187
commands
  silent
  printf "HI_PLAIN_2006_QUERY_TRACE_READY\n"
  continue
end

continue
