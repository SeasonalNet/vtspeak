set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

# Zero only the natural +0x10 value for unit 272822 at position 0, immediately
# before FUN_1001b5f0 sorts the candidate list.
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
      if *(unsigned int *)($node + 8) == 272822
        set $span = *(short *)($node + 0xe)
        set $weighted = *(short *)($node + 0x10)
        printf "HELLO_WEIGHTED_ZERO position=%d id=272822 old_index=%d before_span=%d before_weighted=%d before_rank=%g\n", $position, $index, $span, $weighted, *(float *)($node + 4)
        set {short}($node + 0x10) = 0
        set {float}($node + 4) = -$span
        printf "HELLO_WEIGHTED_ZERO_APPLIED position=%d id=272822 after_span=%d after_weighted=%d after_rank=%g\n", $position, *(short *)($node + 0xe), *(short *)($node + 0x10), *(float *)($node + 4)
      end
      set $index = $index + 1
    end
  end
  continue
end

# Observe the reordered list as soon as FUN_100230a0 returns, before local
# candidate scoring changes the ordering.
break *0x100235f6
condition $bpnum (*(int *)($ebp + 8) == 0)
commands
  silent
  set $state = *(unsigned int *)($ebp + 0xc)
  set $count = *(short *)($state + 0xae988)
  set $index = 0
  set $target_index = -1
  printf "HELLO_WEIGHTED_ZERO_ORDER position=0 count=%d", $count
  while $index < $count && $index < 12
    set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
    set $unit = *(unsigned int *)($node + 8)
    printf " [%d]=%u(key=%g,span=%d,weighted=%d)", $index, $unit, *(float *)($node + 4), *(short *)($node + 0xe), *(short *)($node + 0x10)
    if $unit == 272822
      set $target_index = $index
    end
    set $index = $index + 1
  end
  set $index = 12
  while $index < $count
    set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
    if *(unsigned int *)($node + 8) == 272822
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
  printf "HELLO_WEIGHTED_ZERO_TRACE_READY\n"
  continue
end

continue
