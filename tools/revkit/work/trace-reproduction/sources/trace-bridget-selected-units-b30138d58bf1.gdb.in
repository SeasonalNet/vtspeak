set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $query_calls = 0
set $selected_calls = 0

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
    printf "BRIDGET_2013_QUERY call=%u signature=%02x,%02x,%02x,%02x,%02x,%02x,%02x key=%02x,%02x,%02x,%02x,%02x count=%u\n", $query_calls, *(unsigned char *)$signature, *(unsigned char *)($signature + 1), *(unsigned char *)($signature + 2), *(unsigned char *)($signature + 3), *(unsigned char *)($signature + 4), *(unsigned char *)($signature + 5), *(unsigned char *)($signature + 6), $key0, $key1, $key2, $key3, $key4, $count
    continue
  end
  continue
end

break *0x100235f6
commands
  silent
  set $state = *(unsigned int *)($ebp + 0xc)
  set $position = *(int *)($ebp + 8)
  set $count = *(short *)($state + 0xae988 + $position * 0xfc)
  printf "BRIDGET_CANDIDATE_LIST position=%d count=%d", $position, $count
  set $index = 0
  while $index < $count
    set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
    printf " [%d]=%u(span=%d,weighted=%d,key=%g)", $index, *(unsigned int *)($node + 8), *(short *)($node + 0xe), *(short *)($node + 0x10), *(float *)($node + 4)
    set $index = $index + 1
  end
  printf "\n"
  continue
end

hbreak *0x10018c80
commands
  silent
  set $context = *(int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $return = *(unsigned int *)$esp
  set $record = $state + 0xae894 + $context * 0xfc
  tbreak *$return
  commands
    silent
    set $count = *(short *)($record + 0xf4)
    printf "BRIDGET_TRANSITION context=%d count=%d units:", $context, $count
    set $index = 0
    while $index < $count && $index < 64
      printf " %u", *(unsigned int *)($record + 0x7c + $index * 4)
      set $index = $index + 1
    end
    printf "\n"
    continue
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  printf "BRIDGET_SELECTED_UNIT_HARDWARE n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  continue
end

break *0x10023c81
commands
  silent
  printf "BRIDGET_RELAXED_CLASS_INPUT args=%08x,%08x,%08x,%08x,%08x,%08x,%08x,%08x\n", *(unsigned int *)($ebp + 8), *(unsigned int *)($ebp + 0xc), *(unsigned int *)($ebp + 0x10), *(unsigned int *)($ebp + 0x14), *(unsigned int *)($ebp + 0x18), *(unsigned int *)($ebp + 0x1c), *(unsigned int *)($ebp + 0x20), *(unsigned int *)($ebp + 0x24)
  continue
end

hbreak *0x408187
commands
  silent
  printf "BRIDGET_HARDWARE_TRACE_READY\n"
  continue
end

continue
