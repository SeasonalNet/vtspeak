set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-repeat-2013-producer/adapted-producer-keyclasses-gdb.log
set logging overwrite on
set logging enabled on

set $producer_calls = 0
hbreak *0x10018770
commands
  silent
  set $producer_calls = $producer_calls + 1
  set $signature = *(unsigned int *)($esp + 4)
  set $slot = *(unsigned int *)($esp + 8)
  set $state = *(unsigned int *)($esp + 12)
  set $mode = *(unsigned short *)($esp + 20)
  set $row = $state + ($slot * 3 + 0x76314) * 2
  printf "REPEAT_2013_PRODUCER n=%u slot=%u mode=%u row:", $producer_calls, $slot, $mode
  x/6ub $row
  printf "REPEAT_2013_SOURCE_SIGNATURE:"
  x/7ub $signature
  continue
end

set $query_calls = 0
hbreak *0x10023dc0
commands
  silent
  set $query_calls = $query_calls + 1
  set $signature = *(unsigned int *)($ebp + 8)
  set $return = *(unsigned int *)($ebp + 4)
  set $key0 = *(unsigned char *)(0x1007b7ec + *(unsigned char *)($signature + 1))
  set $key1 = *(unsigned char *)(0x1007b788 + *(unsigned char *)($signature + 2))
  set $key2 = *(unsigned char *)(0x1007b850 + *(unsigned char *)($signature + 3))
  set $key3 = *(unsigned char *)($signature + 5)
  set $key4 = *(unsigned char *)($signature + 6) & 0x20
  tbreak *$return
  commands
    silent
    set $count = $eax & 0xffff
    printf "REPEAT_2013_QUERY n=%u signature=%02x,%02x,%02x,%02x,%02x,%02x,%02x key=%02x,%02x,%02x,%02x,%02x count=%u\n", $query_calls, *(unsigned char *)$signature, *(unsigned char *)($signature + 1), *(unsigned char *)($signature + 2), *(unsigned char *)($signature + 3), *(unsigned char *)($signature + 4), *(unsigned char *)($signature + 5), *(unsigned char *)($signature + 6), $key0, $key1, $key2, $key3, $key4, $count
    continue
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "REPEAT_2013_PRODUCER_TRACE_READY\n"
  continue
end

continue
