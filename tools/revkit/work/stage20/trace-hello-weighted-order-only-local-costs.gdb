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

# FUN_100182e0 returns its local score in ST(0); the caller stores it into
# candidate-node +0 at 0x100238a3. Restrict this trace to Hello position 0.
break *0x1002389a
condition $bpnum (*(int *)($ebp + 8) == 0)
commands
  silent
  set $node = *(unsigned int *)$esi
  printf "HELLO_ORDER_ONLY_LOCAL id=%u key=%g span=%d weighted=%d cost=%g\n", *(unsigned int *)($node + 8), *(float *)($node + 4), *(short *)($node + 0xe), *(short *)($node + 0x10), $st0
  continue
end

# FUN_10018c80 writes each current unit's best cumulative path cost and the
# index of its best predecessor. Capture the complete bounded rows so the
# chosen backtrack can be compared with the local costs above.
break *0x10018c80
commands
  silent
  set $transition_context = *(int *)($esp + 4)
  set $transition_state = *(unsigned int *)($esp + 8)
  set $transition_return = *(unsigned int *)$esp
  set $transition_record = $transition_state + 0xae894 + $transition_context * 0xfc
  printf "HELLO_ORDER_ONLY_TRANSITION_ENTER context=%d current_count=%d previous_count=%d\n", $transition_context, *(short *)($transition_record + 0xf4), *(short *)($transition_record - 0xfc + 0xf4)
  tbreak *$transition_return
  commands
    silent
    set $transition_count = *(short *)($transition_record + 0xf4)
    printf "HELLO_ORDER_ONLY_TRANSITION_RETURN context=%d candidates=%d\n", $transition_context, $transition_count
    set $transition_i = 0
    while $transition_i < $transition_count && $transition_i < 30
      set $transition_id = *(unsigned int *)($transition_record + 0x7c + $transition_i * 4)
      set $transition_cost = *(float *)($transition_state + 0x8212c + ($transition_context * 30 + $transition_i) * 4)
      set $transition_pred = *(short *)($transition_state + 0x9f664 + ($transition_context * 30 + $transition_i) * 2)
      set $transition_prev_record = $transition_record - 0xfc
      set $transition_prev_id = *(unsigned int *)($transition_prev_record + 0x7c + $transition_pred * 4)
      printf "HELLO_ORDER_ONLY_PATH context=%d unit=%u cumulative=%g predecessor=%u predecessor_index=%d\n", $transition_context, $transition_id, $transition_cost, $transition_prev_id, $transition_pred
      set $transition_i = $transition_i + 1
    end
    continue
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
