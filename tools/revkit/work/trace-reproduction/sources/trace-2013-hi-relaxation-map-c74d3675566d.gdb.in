set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /tmp/hi-relaxation-map.gdb.log
set logging overwrite on
set logging enabled on
set $lookup_calls = 0
set $selected_calls = 0

# Report exact key-class results for the producer's four second-side keys.
break *0x10023dc0
commands
  silent
  set $lookup_calls = $lookup_calls + 1
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
    printf "HI_RELAX_LOOKUP n=%u key=%02x%02x%02x%02x%02x count=%u\n", $lookup_calls, $key0, $key1, $key2, $key3, $key4, $eax & 0xffff
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
  printf "HI_RELAX_CANDIDATES position=%d count=%d ids=", $position, $count
  set $index = 0
  while $index < $count && $index < 60
    set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
    printf "%u ", *(unsigned int *)($node + 8)
    set $index = $index + 1
  end
  printf "\n"
  continue
end

hbreak *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  printf "HI_RELAX_SELECTED n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "HI_RELAX_TRACE_READY\n"
  continue
end

continue
