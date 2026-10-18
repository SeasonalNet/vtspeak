set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

# FUN_100230a0 has finished building every node at this call site and is about
# to order them by the float at node+4. Change only unit 266885's weighted
# coverage field and the corresponding rank key for a controlled counterfactual.
break *0x10023331
commands
  silent
  set $position = *(int *)($ebp + 8)
  set $state = *(unsigned int *)($ebp + 0xc)
  set $count = *(short *)($state + 0xae988 + $position * 0xfc)
  set $index = 0
  while $index < $count
    set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
    set $unit = *(unsigned int *)($node + 8)
    if $unit == 266885
      set $span = *(short *)($node + 0xe)
      set $weighted = *(short *)($node + 0x10)
      if $weighted > 0
        printf "APPLE_WEIGHTED_ZERO position=%d id=%u old_index=%d before_span=%d before_weighted=%d before_rank=%g\n", $position, $unit, $index, $span, $weighted, *(float *)($node + 4)
        set {short}($node + 0x10) = 0
        set {float}($node + 4) = -$span
        printf "APPLE_WEIGHTED_ZERO_APPLIED position=%d id=%u after_span=%d after_weighted=%d after_rank=%g\n", $position, $unit, *(short *)($node + 0xe), *(short *)($node + 0x10), *(float *)($node + 4)
      end
    end
    set $index = $index + 1
  end
  continue
end

# This is immediately after FUN_100230a0 returns, so it observes the ordered
# list before FUN_10023350 applies local candidate scoring.
break *0x100235f6
condition $bpnum ((*(int *)($ebp + 8) == 5) || (*(int *)($ebp + 8) == 6))
commands
  silent
  set $position = *(int *)($ebp + 8)
  set $state = *(unsigned int *)($ebp + 0xc)
  set $count = *(short *)($state + 0xae988 + $position * 0xfc)
  set $index = 0
  set $target_index = -1
  printf "APPLE_WEIGHTED_ZERO_ORDER position=%d count=%d", $position, $count
  while $index < $count && $index < 12
    set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
    set $unit = *(unsigned int *)($node + 8)
    printf " [%d]=%u(key=%g,span=%d,weighted=%d)", $index, $unit, *(float *)($node + 4), *(short *)($node + 0xe), *(short *)($node + 0x10)
    if $unit == 266885
      set $target_index = $index
    end
    set $index = $index + 1
  end
  set $index = 12
  while $index < $count
    set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
    if *(unsigned int *)($node + 8) == 266885
      set $target_index = $index
    end
    set $index = $index + 1
  end
  printf " target_index=%d\n", $target_index
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_WEIGHTED_ZERO_TRACE_READY\n"
  continue
end

continue
