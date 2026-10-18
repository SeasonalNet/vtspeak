set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $context = -1
set $state = 0

break *0x10023350
commands
  silent
  set $context = *(int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  continue
end

# Candidate list after continuity ordering/full-span filtering, before scoring.
break *0x10023814
commands
  silent
  if $context == 4
    set $count = *(unsigned short *)($state + 0xae988 + $context * 0xfc)
    set $index = 0
    set $target_rank = -1
    while $index < $count && $index < 10000
      set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
      if *(unsigned int *)($node + 8) == 150652
        set $target_rank = $index
      end
      set $index = $index + 1
    end
    printf "BRIDGET_RANK150652_AFTER_COVERAGE context=4 count=%u index=%d\n", $count, $target_rank
  end
  continue
end

# Capture the result of the target row's local feature scoring if it reaches it.
hbreak *0x1002389a
commands
  silent
  if (*(int *)($ebp + 8) == 4)
    set $node = *(unsigned int *)$esi
    if *(unsigned int *)($node + 8) == 150652
      printf "BRIDGET_RANK150652_LOCAL_SCORE score=%g key=%g right=%d span=%d weighted=%d flags=%#x\n", $st0, *(float *)($node + 4), *(short *)($node + 0xc), *(short *)($node + 0xe), *(short *)($node + 0x10), *(unsigned char *)($node + 0x12)
    end
  end
  continue
end

# Capture target order around the final suffix sort and its resulting top 30.
break *0x10023948
commands
  silent
  if (*(int *)($ebp + 8) == 4)
    set $count = $ebx
    set $prefix = *(int *)($ebp - 0x20)
    set $index = 0
    while $index < $count && $index < 10000
      set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
      if *(unsigned int *)($node + 8) == 150652
        printf "BRIDGET_RANK150652_SORT_INPUT total=%u prefix=%u index=%u score=%g key=%g span=%d weighted=%d\n", $count, $prefix, $index, *(float *)$node, *(float *)($node + 4), *(short *)($node + 0xe), *(short *)($node + 0x10)
      end
      set $index = $index + 1
    end
    tbreak *0x1002394d
    commands
      silent
      set $index = 0
      set $target_rank = -1
      while $index < 30
        set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
        if *(unsigned int *)($node + 8) == 150652
          set $target_rank = $index
        end
        set $index = $index + 1
      end
      printf "BRIDGET_RANK150652_SORTED_TOP30 target_index=%d\n", $target_rank
      continue
    end
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "BRIDGET_RANK150652_TRACE_READY\n"
  continue
end

continue
