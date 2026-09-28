set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

hbreak *0x1001e470
commands
  silent
  set $return = *(unsigned int *)$esp
  thbreak *$return
  commands
    silent
    printf "APPLE_OLD_POSITION_COUNT=%u\n", $eax
    continue
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_OLD_POSITION_TRACE_READY\n"
  continue
end

continue
