set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $selected_component_calls = 0
set $selected_calls = 0

break *0x10018c80
condition $bpnum (*(int *)($esp + 4) == 3)
commands
  silent
  set $entry_position = *(int *)($esp + 4)
  set $entry_state = *(unsigned int *)($esp + 8)
  set $entry_count = *(short *)($entry_state + 0xae988 + $entry_position * 0xfc)
  set $entry_units = $entry_state + 0xae894 + $entry_position * 0xfc + 0x7c
  set $entry_index = 0
  printf "APPLE_CONTEXT3_CANDIDATES count=%d ids=", $entry_count
  while $entry_index < $entry_count
    printf "%u,", *(unsigned int *)($entry_units + $entry_index * 4)
    set $entry_index = $entry_index + 1
  end
  printf "\n"
  continue
end

break *0x1002328a
condition $bpnum ((*(int *)($ebp + 8) == 3) && (*(unsigned int *)(*(unsigned int *)$esi + 8) == 64729))
commands
  silent
  set $force_e_record = *(unsigned int *)$esi
  set $force_e_old = $edx
  set $edx = 4
  printf "APPLE_CONTEXT3_FORCE_E4 id=64729 old_computed_e=%d forced_e=%d\\n", $force_e_old, $edx
  continue
end

break *0x100235f6
condition $bpnum (*(int *)($ebp + 8) == 3)
commands
  silent
  set $meta_position = *(int *)($ebp + 8)
  set $meta_state = *(unsigned int *)($ebp + 0xc)
  set $meta_count = *(short *)($meta_state + 0xae988 + $meta_position * 0xfc)
  set $meta_index = 0
  set $meta_found = 0
  printf "APPLE_CONTEXT3_CONTINUITY_METADATA count=%d", $meta_count
  while $meta_index < $meta_count
    set $meta_ptr = *(unsigned int *)($meta_state + 0x477ac + $meta_index * 4)
    set $meta_id = *(unsigned int *)($meta_ptr + 8)
    if $meta_index < 12
      printf "APPLE_CONTEXT3_ORDER_ROW index=%d id=%u rank_term=%d score=%g fields_c_e_10=%d,%d,%d\\n", $meta_index, $meta_id, -(*(short *)($meta_ptr + 0xe) + 100 * *(short *)($meta_ptr + 0x10)), *(float *)$meta_ptr, *(short *)($meta_ptr + 0xc), *(short *)($meta_ptr + 0xe), *(short *)($meta_ptr + 0x10)
    end
    if ($meta_id == 64729) || ($meta_id == 28631)
      printf " candidate=%u list_index=%d pre_score=%g fields_c_e_10=%d,%d,%d", $meta_id, $meta_index, *(float *)$meta_ptr, *(short *)($meta_ptr + 0xc), *(short *)($meta_ptr + 0xe), *(short *)($meta_ptr + 0x10)
      set $meta_found = 1
    end
    set $meta_index = $meta_index + 1
  end
  printf "\n"
  continue
end

break *0x100235f1
condition $bpnum (*(int *)($ebp + 8) == 3)
commands
  silent
  set $pre_sort_state = *(unsigned int *)($ebp + 0xc)
  set $pre_sort_count = *(short *)($pre_sort_state + 0xae988 + 3 * 0xfc)
  set $pre_sort_index = 0
  printf "APPLE_CONTEXT3_PRE_230A0_ORDER count=%d", $pre_sort_count
  while $pre_sort_index < $pre_sort_count
    set $pre_sort_record = *(unsigned int *)($pre_sort_state + 0x477ac + $pre_sort_index * 4)
    set $pre_sort_id = *(unsigned int *)($pre_sort_record + 8)
    if ($pre_sort_id == 64729) || ($pre_sort_id == 28631)
      printf " candidate=%u list_index=%d score=%g fields_c_e_10=%d,%d,%d", $pre_sort_id, $pre_sort_index, *(float *)$pre_sort_record, *(short *)($pre_sort_record + 0xc), *(short *)($pre_sort_record + 0xe), *(short *)($pre_sort_record + 0x10)
    end
    set $pre_sort_index = $pre_sort_index + 1
  end
  printf "\n"
  continue
end

break *0x100238a7
condition $bpnum (*(int *)($ebp + 8) == 3)
commands
  silent
  set $rank_position = *(int *)($ebp + 8)
  set $rank_state = *(unsigned int *)($ebp + 0xc)
  set $rank_count = *(short *)($rank_state + 0xae988 + $rank_position * 0xfc)
  set $rank_index = 0
  printf "APPLE_CONTEXT3_PRE_CUTOFF_RANKS count=%d", $rank_count
  while $rank_index < $rank_count
    set $rank_ptr = *(unsigned int *)($rank_state + 0x477ac + $rank_index * 4)
    set $rank_id = *(unsigned int *)($rank_ptr + 8)
    if ($rank_id == 64729) || ($rank_id == 28631)
      printf " candidate=%u score=%g fields_c_e_10=%d,%d,%d", $rank_id, *(float *)$rank_ptr, *(short *)($rank_ptr + 0xc), *(short *)($rank_ptr + 0xe), *(short *)($rank_ptr + 0x10)
    end
    set $rank_index = $rank_index + 1
  end
  printf "\n"
  continue
end

break *0x10023814
condition $bpnum (*(int *)($ebp + 8) == 3)
commands
  silent
  set $score_prefix_count = $ebx
  set $score_prefix_state = *(unsigned int *)($ebp + 0xc)
  set $score_prefix_index = 0
  set $score_prefix_64729 = -1
  set $score_prefix_28631 = -1
  while $score_prefix_index < $score_prefix_count
    set $score_prefix_record = *(unsigned int *)($score_prefix_state + 0x477ac + $score_prefix_index * 4)
    set $score_prefix_id = *(unsigned int *)($score_prefix_record + 8)
    if $score_prefix_id == 64729
      set $score_prefix_64729 = $score_prefix_index
    end
    if $score_prefix_id == 28631
      set $score_prefix_28631 = $score_prefix_index
    end
    set $score_prefix_index = $score_prefix_index + 1
  end
  printf "APPLE_CONTEXT3_SCORE_PREFIX count=%d candidate64729_index=%d candidate28631_index=%d first_id=%u\n", $score_prefix_count, $score_prefix_64729, $score_prefix_28631, *(unsigned int *)(*(unsigned int *)($score_prefix_state + 0x477ac) + 8)
  continue
end

break *0x10023895
condition $bpnum ((*(int *)($ebp + 8) == 3) && (($ecx == 64729) || ($ecx == 28631)))
commands
  silent
  set $score_call_record = *(unsigned int *)$esi
  printf "APPLE_CONTEXT3_SCORE_CALL id=%u record=%p fields_c_e_10=%d,%d,%d norm=%g weight=%g\n", $ecx, $score_call_record, *(short *)($score_call_record + 0xc), *(short *)($score_call_record + 0xe), *(short *)($score_call_record + 0x10), *(float *)($esp + 0x10), *(float *)($esp + 0x14)
  continue
end

set $score_split_calls = 0
hbreak *0x100182e0
commands
  silent
  set $score_split_unit = *(unsigned int *)($esp + 4)
  if ($score_split_unit == 64729) || ($score_split_unit == 28631)
    set $score_split_calls = $score_split_calls + 1
    set $score_split_return = *(unsigned int *)$esp
    printf "APPLE_UNIT_SCORE_ENTRY n=%u id=%u state=%p target=%p context=%p norm=%g weight=%g\n", $score_split_calls, $score_split_unit, *(void **)($esp + 8), *(void **)($esp + 12), *(void **)($esp + 16), *(float *)($esp + 20), *(float *)($esp + 24)
    tbreak *0x1001875c
    commands
      silent
      printf "APPLE_UNIT_SCORE_SPLIT id=%u local_penalty=%d distance_intermediate=%g final_factor=%g distance_term=%g returned_score=%g\n", $score_split_unit, *(int *)($ebp - 8), *(float *)($ebp + 0x10), *(float *)($ebp + 0x18), *(float *)($ebp + 0x10) * *(float *)($ebp + 0x18), *(float *)($ebp + 0x10) * *(float *)($ebp + 0x18) + *(int *)($ebp - 8)
      continue
    end
  end
  continue
end

break *0x100191bc
condition $bpnum ((*(int *)($ebp + 8) == 2) && ((*(unsigned int *)($ebp - 8) == 100180) || (*(unsigned int *)($ebp - 8) == 128863) || (*(unsigned int *)($ebp - 8) == 64728)) && ((*(unsigned int *)(*(unsigned int *)($ebp - 0x18)) == 64727) || (*(unsigned int *)(*(unsigned int *)($ebp - 0x18)) == 188825) || (*(unsigned int *)(*(unsigned int *)($ebp - 0x18)) == 142230) || (*(unsigned int *)(*(unsigned int *)($ebp - 0x18)) == 177774))) || (*(unsigned int *)($ebp - 8) == 64729)
commands
  silent
  set $selected_component_calls = $selected_component_calls + 1
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
  printf "APPLE_SELECTED_EDGE_COMPONENTS n=%u context=%d current=%u previous=%u previous_cost=%g local_cost=%g raw=%g derived_a=%g derived_b=%g weights_raw_a_b=%g,%g,%g subtotal=%g reconstructed_total=%g\n", $selected_component_calls, $pair_context, $pair_current, $pair_previous, $pair_previous_cost, $pair_local_cost, $pair_raw_distance, $pair_feature_a, $pair_feature_b, $pair_raw_weight, $pair_a_weight, $pair_b_weight, $pair_feature, $pair_feature / $pair_divisor + $pair_penalty + $pair_previous_cost + $pair_local_cost
  continue
end

hbreak *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  printf "APPLE_SELECTED_UNIT_HARDWARE n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_SELECTED_EDGE_TRACE_READY\n"
  continue
end

continue
