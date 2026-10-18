set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/legacy-neighbor-pool-boundary-hardware-gdb.log
set logging overwrite on
set logging enabled on
set $legacy_dp_calls = 0

hbreak *0x1001cff0
commands
  silent
  set $legacy_dp_calls = $legacy_dp_calls + 1
  set $index = *(unsigned int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $return = *(unsigned int *)$esp
  set $context = $state + 0xbb154 + $index * 0xfc
  set $entry_count = *(unsigned short *)$context
  printf "LEGACY_NEIGHBOR_ENTRY call=%u position=%u shortlist_classes=%u ids:", $legacy_dp_calls, $index, $entry_count
  set $i = 0
  while $i < $entry_count && $i < 30
    printf " %u", *(unsigned int *)($context + 4 + $i * 4)
    set $i = $i + 1
  end
  printf "\n"
  thbreak *$return
  commands
    silent
    set $return_count = *(unsigned short *)($context + 0xf4)
    printf "LEGACY_NEIGHBOR_RETURN call=%u position=%u candidates=%u links:", $legacy_dp_calls, $index, $return_count
    set $i = 0
    while $i < $return_count && $i < 30
      set $unit = *(unsigned int *)($context + 0x7c + $i * 4)
      set $pred_index = *(unsigned short *)($state + 0xaa7b4 + $index * 0x3c + $i * 2)
      if $index > 0
        set $previous_context = $state + 0xbb154 + ($index - 1) * 0xfc
        set $pred_unit = *(unsigned int *)($previous_context + 0x7c + $pred_index * 4)
        printf " {%u <- %u@%u}", $unit, $pred_unit, $pred_index
      else
        printf " {%u}", $unit
      end
      set $i = $i + 1
    end
    printf "\n"
    continue
  end
  continue
end

break *0x408187
commands
  silent
  printf "LEGACY_NEIGHBOR_POOL_BOUNDARY_TRACE_READY\n"
  continue
end

continue
