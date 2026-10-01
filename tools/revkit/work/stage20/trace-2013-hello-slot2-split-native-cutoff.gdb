set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-slot2-split-native/adapted-query-selector-gdb.log
set logging overwrite on
set logging enabled on

set $query_calls = 0
set $split_calls = 0
set $fallback_calls = 0
set $selected_calls = 0
set $last_sum = -1

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
    printf "SLOT2_SPLIT_NATIVE_QUERY call=%u key=%02x,%02x,%02x,%02x,%02x classes=%u\n", $query_calls, $key0, $key1, $key2, $key3, $key4, $eax & 0xffff
    continue
  end
  continue
end

break *0x10024060
commands
  silent
  set $split_calls = $split_calls + 1
  set $last_sum = -1
  set $slot = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  printf "SLOT2_SPLIT_NATIVE_SPLIT_ENTER call=%u slot=%u\n", $split_calls, $slot
  tbreak *$return
  commands
    silent
    printf "SLOT2_SPLIT_NATIVE_SPLIT_RETURN slot=%u sum=%d accepted=%u\n", $slot, $last_sum, $al
    continue
  end
  continue
end

break *0x10023060
commands
  silent
  set $caller = *(unsigned int *)$esp
  if $caller >= 0x10024060 && $caller < 0x100242a0
    set $ids = *(unsigned int *)($esp + 4)
    set $count = *(unsigned short *)($esp + 8)
    set $context = *(unsigned int *)($esp + 12)
    set $weights = *(unsigned int *)($context + 0x8c)
    set $return = $caller
    printf "SLOT2_SPLIT_NATIVE_METRIC_ENTRY count=%u ids_weights:", $count
    set $i = 0
    while $i < $count && $i < 96
      set $id = *(unsigned int *)($ids + $i * 4)
      printf " %u:%u", $id, *(unsigned short *)($weights + $id * 2)
      set $i = $i + 1
    end
    printf "\n"
    tbreak *$return
    commands
      silent
      set $last_sum = $eax
      printf "SLOT2_SPLIT_NATIVE_METRIC_RETURN sum=%u\n", $eax
      continue
    end
  end
  continue
end

break *0x100242a0
commands
  silent
  set $fallback_calls = $fallback_calls + 1
  set $slot = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  printf "SLOT2_SPLIT_NATIVE_FALLBACK_ENTER call=%u slot=%u\n", $fallback_calls, $slot
  tbreak *$return
  commands
    silent
    printf "SLOT2_SPLIT_NATIVE_FALLBACK_RETURN slot=%u positions=%u\n", $slot, $eax
    continue
  end
  continue
end

break *0x10024680
commands
  silent
  set $return = *(unsigned int *)$esp
  tbreak *$return
  commands
    silent
    printf "SLOT2_SPLIT_NATIVE_POSITION_BUILDER_RETURN positions=%u\n", $eax
    continue
  end
  continue
end

break *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  if $selected_calls <= 32
    printf "SLOT2_SPLIT_NATIVE_SELECTED_UNIT n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "SLOT2_SPLIT_NATIVE_TRACE_READY\n"
  continue
end

continue
