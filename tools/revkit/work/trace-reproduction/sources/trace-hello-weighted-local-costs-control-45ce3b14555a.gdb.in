set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $selected_calls = 0

# Matched same-stack control: no node or sort mutation.
break *0x100235f6
commands
  silent
  set $state = *(unsigned int *)($ebp + 0xc)
  set $count = *(short *)($state + 0xae988)
  set $index = 0
  printf "HELLO_ORDER_COST_PRE_SCORE count=%d", $count
  while $index < $count && $index < 8
    set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
    printf " [%d]=%u(key=%g,span=%d,weighted=%d)", $index, *(unsigned int *)($node + 8), *(float *)($node + 4), *(short *)($node + 0xe), *(short *)($node + 0x10)
    set $index = $index + 1
  end
  continue
end

# At 0x1002389a, FUN_100182e0 has returned its local score in ST(0),
# immediately before the caller stores it into candidate-node +0.
break *0x1002389a
condition $bpnum (*(int *)($ebp + 8) == 0)
commands
  silent
  set $node = *(unsigned int *)$esi
  printf "HELLO_ORDER_COST_LOCAL id=%u key=%g span=%d weighted=%d cost=%g\n", *(unsigned int *)($node + 8), *(float *)($node + 4), *(short *)($node + 0xe), *(short *)($node + 0x10), $st0
  continue
end

# Trace each following position's best predecessor and cumulative path score.
break *0x10018c80
commands
  silent
  set $transition_context = *(int *)($esp + 4)
  set $transition_state = *(unsigned int *)($esp + 8)
  set $transition_return = *(unsigned int *)$esp
  set $transition_record = $transition_state + 0xae894 + $transition_context * 0xfc
  printf "HELLO_CONTROL_TRANSITION_ENTER context=%d current_count=%d previous_count=%d\n", $transition_context, *(short *)($transition_record + 0xf4), *(short *)($transition_record - 0xfc + 0xf4)
  tbreak *$transition_return
  commands
    silent
    set $transition_count = *(short *)($transition_record + 0xf4)
    printf "HELLO_CONTROL_TRANSITION_RETURN context=%d candidates=%d\n", $transition_context, $transition_count
    set $transition_i = 0
    while $transition_i < $transition_count && $transition_i < 30
      set $transition_id = *(unsigned int *)($transition_record + 0x7c + $transition_i * 4)
      set $transition_cost = *(float *)($transition_state + 0x8212c + ($transition_context * 30 + $transition_i) * 4)
      set $transition_pred = *(short *)($transition_state + 0x9f664 + ($transition_context * 30 + $transition_i) * 2)
      set $transition_prev_record = $transition_record - 0xfc
      set $transition_prev_id = *(unsigned int *)($transition_prev_record + 0x7c + $transition_pred * 4)
      printf "HELLO_CONTROL_PATH context=%d unit=%u cumulative=%g predecessor=%u predecessor_index=%d\n", $transition_context, $transition_id, $transition_cost, $transition_prev_id, $transition_pred
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
  printf "HELLO_ORDER_COST_SELECTED n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "HELLO_ORDER_COST_TRACE_READY\n"
  continue
end

continue
