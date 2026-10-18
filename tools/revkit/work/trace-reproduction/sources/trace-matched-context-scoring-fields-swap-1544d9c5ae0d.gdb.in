set pagination off
set confirm off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/matched-context-scoring-fields-swap-gdb.log
set logging overwrite on
set logging enabled on
set $rank_calls = 0
set $final_calls = 0
set $prosody_calls = 0

break *0x10023af0
commands
  silent
  set $rank_calls = $rank_calls + 1
  if $rank_calls <= 3
    set $limit = *(int *)($esp + 4)
    set $budget = *(int *)($esp + 8)
    set $output = *(unsigned int *)($esp + 12)
    set $signature = *(unsigned int *)($esp + 16)
    set $state = *(unsigned int *)($esp + 20)
    set $model = *(unsigned int *)($esp + 24)
    set $count = *(unsigned short *)($state + 0xec620)
    set $ids = $state + 0xec170
    set $keys = *(unsigned int *)($model + 0x90)
    set $durations = *(unsigned int *)($model + 0x8c)
    set $ret = *(unsigned int *)$esp
    printf "KEY_RANK_ENTRY call=%d limit=%d budget=%d count=%u signature=", $rank_calls, $limit, $budget, $count
    x/7bx $signature
    set $i = 0
    while $i < $count && $i < 8
      set $id = *(unsigned int *)($ids + $i * 4)
      printf "KEY_CANDIDATE call=%d index=%d class_id=%u duration=%u key=", $rank_calls, $i, $id, *(unsigned short *)($durations + $id * 2)
      x/5bx ($keys + $id * 5)
      set $i = $i + 1
    end
    tbreak *$ret
    commands
      silent
      set $selected = (unsigned short)$eax
      printf "KEY_RANK_RETURN call=%d selected=%u class_ids=", $rank_calls, $selected
      set $i = 0
      while $i < $selected && $i < 8
        printf "%u ", *(unsigned int *)($output + $i * 4)
        set $i = $i + 1
      end
      printf "\n"
      continue
    end
  end
  continue
end

break *0x10023350
commands
  silent
  set $final_calls = $final_calls + 1
  if $final_calls <= 3
    set $context = *(int *)($esp + 4)
    set $state = *(unsigned int *)($esp + 8)
    set $record = $state + 0xae894 + $context * 0xfc
    set $input_count = *(unsigned short *)$record
    set $ret = *(unsigned int *)$esp
    printf "FINAL_RANK_ENTRY call=%d context=%d input_count=%u input_ids=", $final_calls, $context, $input_count
    set $i = 0
    while $i < $input_count && $i < 8
      printf "%u ", *(unsigned int *)($record + 4 + $i * 4)
      set $i = $i + 1
    end
    printf "\n"
    tbreak *$ret
    commands
      silent
      set $selected = *(unsigned short *)($record + 0xf4)
      printf "FINAL_RANK_RETURN call=%d selected=%u unit_ids=", $final_calls, $selected
      set $i = 0
      while $i < $selected && $i < 10
        printf "%u ", *(unsigned int *)($record + 0x7c + $i * 4)
        set $i = $i + 1
      end
      printf "\nFINAL_RANK_COSTS call=%d\n", $final_calls
      set $i = 0
      while $i < $selected && $i < 10
        set $node = *(unsigned int *)($state + 0x477ac + $i * 4)
        printf "FINAL_COST position=%d unit=%u score=%g duration_penalty=%d context_penalty=%d\n", $i, *(unsigned int *)($node + 8), *(float *)$node, *(short *)($node + 0xe), *(short *)($node + 0x10)
        set $i = $i + 1
      end
      continue
    end
  end
  continue
end

break *0x100182e0
commands
  silent
  set $prosody_calls = $prosody_calls + 1
  if $prosody_calls <= 8
    set $unit = *(int *)($esp + 4)
    set $target = *(unsigned int *)($esp + 12)
    set $feature_view = *(unsigned int *)($esp + 16)
    set $scale = *(float *)($esp + 20)
    set $ret = *(unsigned int *)$esp
    printf "CEP_PROSODY_ENTRY call=%d unit=%d scale=%g target=", $prosody_calls, $unit, $scale
    x/7bx $target
    printf "CEP_PROSODY_FEATURE_VIEW="
    x/10bx $feature_view
    tbreak *$ret
    commands
      silent
      printf "CEP_PROSODY_RETURN call=%d cost=%g\n", $prosody_calls, $st0
      continue
    end
  end
  continue
end


set $selected_calls = 0

break *0x10024060
commands
  silent
  set $slot = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  set $last_sum = -1
  printf "OLD_CUTOFF_SPLIT_ENTER slot=%u\n", $slot
  tbreak *$return
  commands
    silent
    set $native = $al
    if $native == 0 && $last_sum > 2
      printf "OLD_CUTOFF_OVERRIDE slot=%u sum=%u native=%u overridden=1\n", $slot, $last_sum, $native
      set $eax = 1
    else
      printf "OLD_CUTOFF_SPLIT_RETURN slot=%u sum=%d accepted=%u\n", $slot, $last_sum, $native
    end
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
    printf "OLD_CUTOFF_METRIC_ENTRY count=%u ids_weights:", $count
    set $i = 0
    while $i < $count && $i < 64
      set $id = *(unsigned int *)($ids + $i * 4)
      printf " %u:%u", $id, *(unsigned short *)($weights + $id * 2)
      set $i = $i + 1
    end
    printf "\n"
    tbreak *$return
    commands
      silent
      set $last_sum = $eax
      printf "OLD_CUTOFF_METRIC_RETURN sum=%u\n", $eax
      continue
    end
  end
  continue
end

break *0x100242a0
commands
  silent
  set $slot = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  printf "OLD_CUTOFF_FALLBACK_ENTER slot=%u\n", $slot
  tbreak *$return
  commands
    silent
    printf "OLD_CUTOFF_FALLBACK_RETURN slot=%u positions=%u\n", $slot, $eax
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
    printf "OLD_CUTOFF_POSITION_BUILDER_RETURN positions=%u\n", $eax
    continue
  end
  continue
end

break *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  if $selected_calls <= 24
    printf "OLD_CUTOFF_SELECTED_UNIT n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  end
  continue
end

break *0x408187
commands
  silent
  printf "OLD_CUTOFF_TRACE_READY\n"
  continue
end

break *0x10018c80
commands
  silent
  set $transition_context = *(int *)($esp + 4)
  set $transition_state = *(unsigned int *)($esp + 8)
  set $transition_return = *(unsigned int *)$esp
  set $transition_record = $transition_state + 0xae894 + $transition_context * 0xfc
  printf "TRANSITION_ENTER context=%d current_count=%d previous_count=%d\n", $transition_context, *(short *)($transition_record + 0xf4), *(short *)($transition_record - 0xfc + 0xf4)
  tbreak *$transition_return
  commands
    silent
    set $transition_count = *(short *)($transition_record + 0xf4)
    printf "TRANSITION_RETURN context=%d candidates=%d cumulative_costs_and_predecessors:\n", $transition_context, $transition_count
    set $transition_i = 0
    while $transition_i < $transition_count && $transition_i < 30
      set $transition_id = *(unsigned int *)($transition_record + 0x7c + $transition_i * 4)
      set $transition_cost = *(float *)($transition_state + 0x8212c + ($transition_context * 30 + $transition_i) * 4)
      set $transition_pred = *(short *)($transition_state + 0x9f664 + ($transition_context * 30 + $transition_i) * 2)
      set $transition_prev_record = $transition_record - 0xfc
      set $transition_prev_id = *(unsigned int *)($transition_prev_record + 0x7c + $transition_pred * 4)
      printf "TRANSITION_CANDIDATE context=%d unit=%u cost=%g predecessor=%u predecessor_index=%d\n", $transition_context, $transition_id, $transition_cost, $transition_prev_id, $transition_pred
      set $transition_i = $transition_i + 1
    end
    continue
  end
  continue
end

break *0x100191bc
commands
  silent
  set $pair_context = *(int *)($ebp + 8)
  set $pair_current = *(unsigned int *)($ebp - 8)
  set $pair_previous_ptr = *(unsigned int *)($ebp - 0x18)
  set $pair_previous = *(unsigned int *)$pair_previous_ptr
  if (($pair_context == 1 && ($pair_current == 272823 || $pair_current == 273370) && ($pair_previous == 272822 || $pair_previous == 273369)) || ($pair_context == 2 && ($pair_current == 272824 || $pair_current == 273371) && ($pair_previous == 272823 || $pair_previous == 273370)) || ($pair_context == 3 && ($pair_current == 272825 || $pair_current == 255282) && ($pair_previous == 272824 || $pair_previous == 273371)))
    set $pair_feature = *(float *)($ebp - 0x4c)
    set $pair_divisor = *(int *)($ebp - 0x10)
    set $pair_penalty = *(int *)($ebp - 0xc)
    set $pair_previous_cost_ptr = *(unsigned int *)($ebp - 0x54)
    set $pair_previous_cost = *(float *)$pair_previous_cost_ptr
    set $pair_local_cost = *(float *)($ebp - 0x74)
    printf "TRANSITION_COMPONENTS context=%d current=%u previous=%u feature=%g divisor=%d categorical=%d previous_cost=%g local_cost=%g sum=%g\n", $pair_context, $pair_current, $pair_previous, $pair_feature, $pair_divisor, $pair_penalty, $pair_previous_cost, $pair_local_cost, $pair_feature / $pair_divisor + $pair_penalty + $pair_previous_cost + $pair_local_cost
  end
  continue
end

continue
