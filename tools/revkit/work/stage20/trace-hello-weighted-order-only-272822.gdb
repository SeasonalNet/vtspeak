set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $selected_calls = 0

# First remove the natural coverage contribution so FUN_1001b5f0 orders
# 273369 ahead of 272822 at position 0.
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
        printf "HELLO_ORDER_ONLY_PRE_SORT id=272822 span=%d weighted=%d rank=%g\n", *(short *)($node + 0xe), *(short *)($node + 0x10), *(float *)($node + 4)
        set {short}($node + 0x10) = 0
        set {float}($node + 4) = -6
      end
      set $index = $index + 1
    end
  end
  continue
end

# Preserve the altered node order, but restore the original metric and key
# before FUN_10023350 scores candidates. This isolates order from the later
# +0x10 distance-scaling use.
break *0x100235f6
condition $bpnum (*(int *)($ebp + 8) == 0)
commands
  silent
  set $state = *(unsigned int *)($ebp + 0xc)
  set $count = *(short *)($state + 0xae988)
  set $index = 0
  printf "HELLO_ORDER_ONLY_PRE_SCORE count=%d", $count
  while $index < $count && $index < 8
    set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
    printf " [%d]=%u(key=%g,span=%d,weighted=%d)", $index, *(unsigned int *)($node + 8), *(float *)($node + 4), *(short *)($node + 0xe), *(short *)($node + 0x10)
    set $index = $index + 1
  end
  set $index = 0
  while $index < $count
    set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
    if *(unsigned int *)($node + 8) == 272822
      set {short}($node + 0x10) = 4
      set {float}($node + 4) = -406
      printf "HELLO_ORDER_ONLY_RESTORED id=272822 index=%d span=%d weighted=%d rank=%g\n", $index, *(short *)($node + 0xe), *(short *)($node + 0x10), *(float *)($node + 4)
    end
    set $index = $index + 1
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  printf "HELLO_ORDER_ONLY_SELECTED n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "HELLO_ORDER_ONLY_TRACE_READY\n"
  continue
end

continue
