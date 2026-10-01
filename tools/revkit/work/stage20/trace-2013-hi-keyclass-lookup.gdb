set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hi-2013-keyclass-lookup/adapted-keyclass-lookup-gdb.log
set logging overwrite on
set logging enabled on
set $lookup_calls = 0

break *0x10023dc0
commands
  silent
  set $lookup_calls = $lookup_calls + 1
  set $signature = *(unsigned int *)($ebp + 8)
  set $result_array = *(unsigned int *)($ebp + 12)
  set $engine_state = *(unsigned int *)($ebp + 16)
  set $model = *(unsigned int *)($ebp + 20)
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
    printf "HI_2013_KEYCLASS_LOOKUP call=%u signature=%02x,%02x,%02x,%02x,%02x,%02x,%02x key=%02x,%02x,%02x,%02x,%02x count=%u ids:", $lookup_calls, *(unsigned char *)$signature, *(unsigned char *)($signature + 1), *(unsigned char *)($signature + 2), *(unsigned char *)($signature + 3), *(unsigned char *)($signature + 4), *(unsigned char *)($signature + 5), *(unsigned char *)($signature + 6), $key0, $key1, $key2, $key3, $key4, $count
    set $i = 0
    while $i < $count && $i < 80
      printf " %u", *(unsigned int *)($result_array + $i * 4)
      set $i = $i + 1
    end
    printf "\n"
    continue
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "HI_2013_KEYCLASS_LOOKUP_TRACE_READY\n"
  continue
end

continue
