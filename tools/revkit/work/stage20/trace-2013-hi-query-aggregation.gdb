set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /tmp/hi-query-aggregation.gdb.log
set logging overwrite on
set logging enabled on
set $hi_query_calls = 0
set $hi_selector_calls = 0
set $hi_split_calls = 0
set $hi_fallback_calls = 0
set $hi_selected_calls = 0
set $hi_last_metric = -1

# Follow each concrete slot through the direct/relaxed wrapper.
break *0x10023f90
commands
  silent
  set $hi_selector_calls = $hi_selector_calls + 1
  set $hi_selector_slot = *(unsigned int *)($ebp + 0xc)
  set $hi_selector_state = *(unsigned int *)($ebp + 0x10)
  set $hi_selector_signature = *(unsigned int *)($ebp + 8)
  set $hi_selector_cursor = *(short *)($hi_selector_state + 0xec620)
  set $hi_selector_flag = *(unsigned char *)($hi_selector_state + 0xec62c + $hi_selector_slot * 6)
  printf "HI_SELECTOR_ENTER call=%u slot=%u flag=%u class_cursor=%d target=", $hi_selector_calls, $hi_selector_slot, $hi_selector_flag, $hi_selector_cursor
  x/7ub $hi_selector_signature
  set $hi_selector_return = *(unsigned int *)$esp
  tbreak *$hi_selector_return
  commands
    silent
    printf "HI_SELECTOR_RETURN call=%u slot=%u result=%#x low16=%u class_cursor=%d\n", $hi_selector_calls, $hi_selector_slot, $eax, $eax & 0xffff, *(short *)($hi_selector_state + 0xec620)
    continue
  end
  continue
end

# Record the transformed exact key and class IDs returned by each lookup.
break *0x10023dc0
commands
  silent
  set $hi_query_calls = $hi_query_calls + 1
  set $hi_query_signature = *(unsigned int *)($ebp + 8)
  set $hi_query_classes = *(unsigned int *)($ebp + 0xc)
  set $hi_query_return = *(unsigned int *)($ebp + 4)
  set $hi_q0 = *(unsigned char *)(0x1007b7ec + *(unsigned char *)($hi_query_signature + 1))
  set $hi_q1 = *(unsigned char *)(0x1007b788 + *(unsigned char *)($hi_query_signature + 2))
  set $hi_q2 = *(unsigned char *)(0x1007b850 + *(unsigned char *)($hi_query_signature + 3))
  set $hi_q3 = *(unsigned char *)($hi_query_signature + 5)
  set $hi_q4 = *(unsigned char *)($hi_query_signature + 6) & 0x20
  tbreak *$hi_query_return
  commands
    silent
    set $hi_class_count = $eax & 0xffff
    printf "HI_EXACT_QUERY n=%u key=%02x%02x%02x%02x%02x count=%u class_ids=", $hi_query_calls, $hi_q0, $hi_q1, $hi_q2, $hi_q3, $hi_q4, $hi_class_count
    set $hi_class_index = 0
    while $hi_class_index < $hi_class_count && $hi_class_index < 32
      printf "%u ", *(unsigned int *)($hi_query_classes + $hi_class_index * 4)
      set $hi_class_index = $hi_class_index + 1
    end
    printf "\n"
    continue
  end
  continue
end

# Capture each metric list and the sum consumed by FUN_10024060.
break *0x10023060
commands
  silent
  set $hi_metric_caller = *(unsigned int *)$esp
  if $hi_metric_caller >= 0x10024060 && $hi_metric_caller < 0x100242a0
    set $hi_metric_ids = *(unsigned int *)($esp + 4)
    set $hi_metric_count = *(unsigned short *)($esp + 8)
    set $hi_metric_context = *(unsigned int *)($esp + 12)
    set $hi_metric_weights = *(unsigned int *)($hi_metric_context + 0x8c)
    set $hi_metric_return = $hi_metric_caller
    printf "HI_METRIC_INPUT count=%u id_weight=", $hi_metric_count
    set $hi_metric_index = 0
    while $hi_metric_index < $hi_metric_count && $hi_metric_index < 96
      set $hi_metric_id = *(unsigned int *)($hi_metric_ids + $hi_metric_index * 4)
      printf "%u:%u ", $hi_metric_id, *(unsigned short *)($hi_metric_weights + $hi_metric_id * 2)
      set $hi_metric_index = $hi_metric_index + 1
    end
    printf "\n"
    tbreak *$hi_metric_return
    commands
      silent
      set $hi_last_metric = $eax
      printf "HI_METRIC_SUM value=%u\n", $eax
      continue
    end
  end
  continue
end

break *0x10024060
commands
  silent
  set $hi_split_calls = $hi_split_calls + 1
  set $hi_last_metric = -1
  set $hi_split_slot = *(unsigned int *)($esp + 4)
  set $hi_split_return = *(unsigned int *)$esp
  printf "HI_SPLIT_ENTER call=%u slot=%u\n", $hi_split_calls, $hi_split_slot
  tbreak *$hi_split_return
  commands
    silent
    printf "HI_SPLIT_RETURN call=%u slot=%u metric=%d accepted=%u\n", $hi_split_calls, $hi_split_slot, $hi_last_metric, $al
    continue
  end
  continue
end

break *0x100242a0
commands
  silent
  set $hi_fallback_calls = $hi_fallback_calls + 1
  set $hi_fallback_slot = *(unsigned int *)($esp + 4)
  set $hi_fallback_return = *(unsigned int *)$esp
  printf "HI_FALLBACK_ENTER call=%u slot=%u\n", $hi_fallback_calls, $hi_fallback_slot
  tbreak *$hi_fallback_return
  commands
    silent
    printf "HI_FALLBACK_RETURN call=%u slot=%u positions=%u\n", $hi_fallback_calls, $hi_fallback_slot, $eax
    continue
  end
  continue
end

break *0x10024680
commands
  silent
  set $hi_builder_return = *(unsigned int *)$esp
  tbreak *$hi_builder_return
  commands
    silent
    printf "HI_POSITION_BUILDER_RETURN positions=%u\n", $eax
    continue
  end
  continue
end

break *0x100235f6
commands
  silent
  set $hi_candidate_state = *(unsigned int *)($ebp + 0xc)
  set $hi_candidate_position = *(int *)($ebp + 8)
  if $hi_candidate_position < 4
    set $hi_candidate_count = *(short *)($hi_candidate_state + 0xae988 + $hi_candidate_position * 0xfc)
    printf "HI_FINAL_CANDIDATES position=%d count=%d ids=", $hi_candidate_position, $hi_candidate_count
    set $hi_candidate_index = 0
    while $hi_candidate_index < $hi_candidate_count && $hi_candidate_index < 80
      set $hi_candidate_node = *(unsigned int *)($hi_candidate_state + 0x477ac + $hi_candidate_index * 4)
      printf "%u ", *(unsigned int *)($hi_candidate_node + 8)
      set $hi_candidate_index = $hi_candidate_index + 1
    end
    printf "\n"
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  set $hi_selected_calls = $hi_selected_calls + 1
  printf "HI_AGG_SELECTED n=%u id=%u\n", $hi_selected_calls, *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "HI_QUERY_AGGREGATION_TRACE_READY\n"
  continue
end

continue
