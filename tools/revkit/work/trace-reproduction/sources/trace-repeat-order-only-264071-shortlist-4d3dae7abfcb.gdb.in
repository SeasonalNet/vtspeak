set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-repeat-2013-order-only-264071-shortlist/adapted-order-only-264071-gdb.log
set logging overwrite on
set logging enabled on

# Keep the controlled 2006 first choice by changing only its initial rank key.
break *0x10023331
commands
  silent
  set $position = *(int *)($ebp + 8)
  set $state = *(unsigned int *)($ebp + 0xc)
  if $position == 0
    set $count = *(short *)($state + 0xae988)
    set $index = 0
    while $index < $count
      set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
      if *(unsigned int *)($node + 8) == 273369
        set {float}($node + 4) = -2.0
      end
      set $index = $index + 1
    end
  end
  continue
end

# Immediately after FUN_100230a0 has built and initially sorted the 79-row
# candidate list, move only 264071's pointer to index 0. Leave node fields
# unchanged before FUN_10023350 computes its preserved prefix and local costs.
break *0x100235f6
commands
  silent
  set $position = *(int *)($ebp + 8)
  set $state = *(unsigned int *)($ebp + 0xc)
  if $position == 6
    set $count = *(short *)($state + 0xae988 + $position * 0xfc)
    set $array = $state + 0x477ac
    set $index = 0
    while $index < $count
      set $node = *(unsigned int *)($array + $index * 4)
      if *(unsigned int *)($node + 8) == 264071
        set $target = $node
        set $target_index = $index
      end
      set $index = $index + 1
    end
    if $target_index > 0
      set $first = *(unsigned int *)$array
      set *(unsigned int *)$array = $target
      set *(unsigned int *)($array + $target_index * 4) = $first
      printf "ORDER_ONLY_SWAP position=6 unit=264071 old_index=%d new_index=0 unit_key=%g unit_span=%d unit_weighted=%d swapped_with=%u swapped_key=%g\n", $target_index, *(float *)($target + 4), *(short *)($target + 0xe), *(short *)($target + 0x10), *(unsigned int *)($first + 8), *(float *)($first + 4)
    end
  end
  continue
end

# Record final local scores for the old/new rows after unit scoring.
break *0x1002389a
commands
  silent
  set $position = *(int *)($ebp + 8)
  if $position == 6
    set $node = *(unsigned int *)$esi
    set $unit = *(unsigned int *)($node + 8)
    if $unit == 264071 || $unit == 49874
      printf "SHORTLIST_UNIT_SCORE position=6 unit=%u score=%g key=%g\n", $unit, $st0, *(float *)($node + 4)
    end
  end
  continue
end

# Capture the actual score sort, then the 30 rows sent to the DP.
break *0x10023948
commands
  silent
  set $position = *(int *)($ebp + 8)
  if $position == 6
    set $state = *(unsigned int *)($ebp + 0xc)
    set $total = $ebx
    set $prefix = *(int *)($ebp - 0x20)
    set $sort_count = *(int *)$esp
    printf "SHORTLIST_SORT_INPUT position=6 total=%d presorted_prefix=%d suffix_count=%d\n", $total, $prefix, $sort_count
    set $i = 0
    while $i < $total && $i < 100
      set $node = *(unsigned int *)($state + 0x477ac + $i * 4)
      set $unit = *(unsigned int *)($node + 8)
      if $unit == 264071 || $unit == 49874
        printf "SHORTLIST_SORT_TARGET position=6 index=%d unit=%u score=%g key=%g flags=%#x\n", $i, $unit, *(float *)$node, *(float *)($node + 4), *(unsigned char *)($node + 0x12)
      end
      set $i = $i + 1
    end
    tbreak *0x1002394d
    commands
      silent
      set $i = 0
      printf "SHORTLIST_SORTED_TOP30 position=6"
      while $i < 30
        set $node = *(unsigned int *)($state + 0x477ac + $i * 4)
        printf " [%d]=%u(score=%g,key=%g)", $i, *(unsigned int *)($node + 8), *(float *)$node, *(float *)($node + 4)
        set $i = $i + 1
      end
      printf "\n"
      continue
    end
  end
  continue
end

break *0x408187
commands
  silent
  printf "ORDER_ONLY_SHORTLIST_TRACE_READY\n"
  continue
end

continue
