set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $transition_calls = 0
set $backtrack_calls = 0

break *0x10018c80
commands
  silent
  set $transition_calls = $transition_calls + 1
  set $context = *(int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $record = $state + 0xae894 + $context * 0xfc
  set $previous = $record - 0xfc
  set $return = *(unsigned int *)$esp
  printf "APPLE_TRANSITION_ENTRY n=%u context=%d current_count=%u previous_count=%u\n", $transition_calls, $context, *(unsigned short *)($record + 0xf4), *(unsigned short *)($previous + 0xf4)
  tbreak *$return
  commands
    silent
    set $count = *(unsigned short *)($record + 0xf4)
    set $i = 0
    while $i < $count && $i < 30
      set $slot = $context * 30 + $i
      set $unit = *(unsigned int *)($record + 0x7c + $i * 4)
      set $cost = *(float *)($state + 0x8212c + $slot * 4)
      set $pred = *(short *)($state + 0x9f664 + $slot * 2)
      set $pred_unit = -1
      if $pred >= 0 && $pred < *(unsigned short *)($previous + 0xf4)
        set $pred_unit = *(unsigned int *)($previous + 0x7c + $pred * 4)
      end
      printf "APPLE_TRANSITION_CANDIDATE context=%d index=%d unit=%u cumulative_cost=%g predecessor_index=%d predecessor_unit=%d\n", $context, $i, $unit, $cost, $pred, $pred_unit
      set $i = $i + 1
    end
    continue
  end
  continue
end

break *0x10024900
commands
  silent
  set $backtrack_calls = $backtrack_calls + 1
  set $first = *(int *)($esp + 4)
  set $last = *(int *)($esp + 8)
  set $state = *(unsigned int *)($esp + 12)
  set $return = *(unsigned int *)$esp
  printf "APPLE_BACKTRACK_ENTRY n=%u first=%d last=%d\n", $backtrack_calls, $first, $last
  tbreak *$return
  commands
    silent
    set $context = $first
    while $context <= $last
      set $record = $state + 0xae894 + $context * 0xfc
      set $index = *(short *)($state + 0xae0c4 + $context * 2)
      set $count = *(unsigned short *)($record + 0xf4)
      set $unit = -1
      set $cost = 0.0
      if $index >= 0 && $index < $count
        set $unit = *(unsigned int *)($record + 0x7c + $index * 4)
        set $slot = $context * 30 + $index
        set $cost = *(float *)($state + 0x8212c + $slot * 4)
      end
      printf "APPLE_BACKTRACK_RESULT context=%d index=%d count=%u unit=%d cumulative_cost=%g\n", $context, $index, $count, $unit, $cost
      set $context = $context + 1
    end
    continue
  end
  continue
end

break *0x100191bc
condition $bpnum (*(int *)($ebp + 8) == 2) && ((*(unsigned int *)($ebp - 8) == 100180) || (*(unsigned int *)($ebp - 8) == 128863)) && (*(unsigned int *)(*(unsigned int *)($ebp - 0x18)) == 142230)
commands
  silent
  set $pair_context = *(int *)($ebp + 8)
  set $pair_current = *(unsigned int *)($ebp - 8)
  set $pair_previous_ptr = *(unsigned int *)($ebp - 0x18)
  set $pair_previous = *(unsigned int *)$pair_previous_ptr
  set $pair_feature = *(float *)($ebp - 0x4c)
  set $pair_divisor = *(int *)($ebp - 0x10)
  set $pair_penalty = *(int *)($ebp - 0xc)
  set $pair_previous_cost_ptr = *(unsigned int *)($ebp - 0x54)
  set $pair_previous_cost = *(float *)$pair_previous_cost_ptr
  set $pair_local_cost = *(float *)($ebp - 0x74)
  set $pair_raw_distance = *(float *)($ebp - 0x38)
  set $pair_feature_a = *(float *)($ebp - 0x34)
  set $pair_feature_b = *(float *)($ebp - 0x30)
  set $pair_weight_group = *(signed char *)($ebp - 1)
  set $pair_raw_weight = *(float *)(0x1007c288 + $pair_weight_group * 12)
  set $pair_a_weight = *(float *)(0x1007c290 + $pair_weight_group * 12)
  set $pair_b_weight = *(float *)(0x1007c28c + $pair_weight_group * 12)
  set $pair_current_metric_group = *(short *)($ebp - 0x14)
  set $pair_previous_metric_group = *(short *)($ebp - 0x28)
  set $pair_current_metric_array = *(unsigned int *)($ebx + $pair_current_metric_group * 4 + 0x3c)
  set $pair_previous_metric_array = *(unsigned int *)($ebx + $pair_previous_metric_group * 4 + 0x3c)
  set $pair_current_metric_word = *(unsigned short *)($pair_current_metric_array + $pair_current * 2)
  set $pair_previous_metric_word = *(unsigned short *)($pair_previous_metric_array + $pair_previous * 2)
  printf "APPLE_TRANSITION_COMPONENTS context=%d current=%u previous=%u weight_group=%d metric_groups_current_previous=%d,%d metric_words_current_previous=%u,%u metric_indices_current_previous=%u,%u weighted_subtotal=%g weights_raw_a_b=%g,%g,%g divisor=%d categorical=%d previous_cost=%g local_cost=%g raw_distance=%g derived_distance_a=%g derived_distance_b=%g reconstructed=%g\n", $pair_context, $pair_current, $pair_previous, $pair_weight_group, $pair_current_metric_group, $pair_previous_metric_group, $pair_current_metric_word, $pair_previous_metric_word, $pair_current_metric_word & 0x3fff, $pair_previous_metric_word & 0x3fff, $pair_feature, $pair_raw_weight, $pair_a_weight, $pair_b_weight, $pair_divisor, $pair_penalty, $pair_previous_cost, $pair_local_cost, $pair_raw_distance, $pair_feature_a, $pair_feature_b, $pair_feature / $pair_divisor + $pair_penalty + $pair_previous_cost + $pair_local_cost
  continue
end

break *0x408187
commands
  silent
  printf "APPLE_PATH_TRACE_READY\n"
  continue
end

continue
