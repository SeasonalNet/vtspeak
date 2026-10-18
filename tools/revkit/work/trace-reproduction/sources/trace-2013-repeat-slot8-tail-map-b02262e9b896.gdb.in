set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-repeat-2013-slot8-tail-map/adapted-query-keyclasses-gdb.log
set logging overwrite on
set logging enabled on

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
    printf "REPEAT_SLOT8_QUERY n=%u signature=%02x,%02x,%02x,%02x,%02x,%02x,%02x key=%02x,%02x,%02x,%02x,%02x count=%u\n", $query_calls, *(unsigned char *)$signature, *(unsigned char *)($signature + 1), *(unsigned char *)($signature + 2), *(unsigned char *)($signature + 3), *(unsigned char *)($signature + 4), *(unsigned char *)($signature + 5), *(unsigned char *)($signature + 6), $key0, $key1, $key2, $key3, $key4, $count
    continue
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "REPEAT_SLOT8_TRACE_READY\n"
  continue
end

continue
