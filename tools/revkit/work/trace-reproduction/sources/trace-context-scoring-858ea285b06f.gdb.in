set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

set $rank_calls = 0
set $pair_calls = 0
set $prune_calls = 0
set $backtrack_calls = 0
set $active_context = -1
set $scorer_calls = 0
set $metadata_calls = 0

# Observe the candidate metadata pass immediately before and after it writes
# node fields consumed by the local scorer.
break *0x100230a0
commands
  silent
  set $metadata_calls = $metadata_calls + 1
  set $context = *(int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $third_arg = *(unsigned int *)($esp + 12)
  set $record = $state + 0xae894 + $context * 0xfc
  set $count = *(unsigned short *)($record + 0xf4)
  set $ret = *(unsigned int *)$esp
  printf "NODE_METADATA_ENTRY call=%u context=%d count=%u arg3=%#x\n", $metadata_calls, $context, $count, $third_arg
  set $i = 0
  while $i < $count && $i < 30
    set $node = *(unsigned int *)($state + 0x477ac + $i * 4)
    printf "NODE_METADATA_BEFORE context=%d index=%d unit=%u duration=%d penalty=%d flag=%d\n", $context, $i, *(unsigned int *)($node + 8), *(short *)($node + 0xe), *(short *)($node + 0x10), *(short *)($node + 0xc)
    set $i = $i + 1
  end
  tbreak *$ret
  commands
    silent
    printf "NODE_METADATA_RETURN call=%u context=%d\n", $metadata_calls, $context
    set $i = 0
    while $i < $count && $i < 30
      set $node = *(unsigned int *)($state + 0x477ac + $i * 4)
      printf "NODE_METADATA_AFTER context=%d index=%d unit=%u duration=%d penalty=%d flag=%d\n", $context, $i, *(unsigned int *)($node + 8), *(short *)($node + 0xe), *(short *)($node + 0x10), *(short *)($node + 0xc)
      set $i = $i + 1
    end
    continue
  end
  continue
end

# Inspect the two accumulated stack counters and the final stored metadata
# for the candidate that changes in the right P/B control.
break *0x100232ee
commands
  silent
  set $context = *(int *)($ebp + 8)
  set $node = $esi
  if $context == 0 && *(unsigned int *)($node + 8) == 3837
    printf "NODE_METADATA_COUNTERS context=%d unit=3837 stack14=%d stack18=%d stack04=%d stack10=%d duration=%d field0c=%d field10=%d\n", $context, *(short *)($ebp - 0x14), *(short *)($ebp - 0x18), *(short *)($ebp - 4), *(short *)($ebp + 0x10), *(short *)($node + 0xe), *(short *)($node + 0xc), *(short *)($node + 0x10)
  end
  continue
end

# Verify the table index and selected value for following context lookups
# that contribute to candidate 3837 in the first span.
break *0x1002326b
commands
  silent
  set $context = *(int *)($ebp + 8)
  set $node_list = *(unsigned int *)($ebp - 8)
  set $node = *(unsigned int *)$node_list
  if $context == 0 && *(unsigned int *)($node + 8) == 3837
    set $row4_ptr = *(unsigned int *)($ebp - 0x20)
    set $lookup_addr = $ebp + $ecx * 2 - 0x40
    printf "FORWARD_TABLE_LOOKUP unit=3837 context_index=%d row4_index=%u row4_ptr=%#x row4_value=%u table_addr=%#x table_value=%u\n", $eax, $ecx, $row4_ptr, *(unsigned char *)$row4_ptr, $lookup_addr, *(unsigned short *)$lookup_addr
  end
  continue
end

break *0x10023270
commands
  silent
  set $context = *(int *)($ebp + 8)
  set $node_list = *(unsigned int *)($ebp - 8)
  set $node = *(unsigned int *)$node_list
  if $context == 0 && *(unsigned int *)($node + 8) == 3837
    set $state = *(unsigned int *)($ebp + 12)
    set $next_context = $eax
    set $candidate_row = $state + 0xec628 + $next_context * 6
    set $next_record = $state + 0xae894 + $next_context * 0xfc
    set $row4_ptr = *(unsigned int *)($ebp - 0x20)
    printf "FORWARD_LOOKUP_HIT unit=3837 next_context=%d ecx_after_table=%u contribution=%u accumulator_before=%d candidate_count=%u row4_ptr=%#x row4_value=%u context_row=", $next_context, $ecx, ($edx & 0xffff), *(short *)($ebp - 0x18), *(short *)$next_record, $row4_ptr, *(unsigned char *)$row4_ptr
    printf "CONTEXT_ROW="
    x/6bx $candidate_row
    set $i = 0
    printf "NEXT_CONTEXT_IDS="
    while $i < *(short *)$next_record && $i < 30
      printf "%u ", *(unsigned int *)($next_record + 4 + $i * 4)
      set $i = $i + 1
    end
    printf "\n"
  end
  continue
end

# Per-context candidate generation and local unit costs.
break *0x10023350
commands
  silent
  set $rank_calls = $rank_calls + 1
  set $context = *(int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $record = $state + 0xae894 + $context * 0xfc
  set $active_context = $context
  set $context_row = $state + ($context * 3 + 0x76314) * 2
  set $input_count = *(unsigned short *)$record
  set $ret = *(unsigned int *)$esp
  printf "CONTEXT_ROW context=%d bytes:", $context
  x/6bx $context_row
  printf "CANDIDATES_ENTRY context=%d input_count=%u ids=", $context, $input_count
  set $i = 0
  while $i < $input_count && $i < 30
    printf "%u ", *(unsigned int *)($record + 4 + $i * 4)
    set $i = $i + 1
  end
  printf "\n"
  tbreak *$ret
  commands
    silent
    set $count = *(unsigned short *)($record + 0xf4)
    printf "CANDIDATES_RETURN context=%d count=%u\n", $context, $count
    set $i = 0
    while $i < $count && $i < 30
      set $node = *(unsigned int *)($state + 0x477ac + $i * 4)
      printf "LOCAL_COST context=%d index=%d unit=%u cost=%g duration_penalty=%d context_penalty=%d flag=%d\n", $context, $i, *(unsigned int *)($node + 8), *(float *)$node, *(short *)($node + 0xe), *(short *)($node + 0x10), *(short *)($node + 0xc)
      set $i = $i + 1
    end
    continue
  end
  continue
end

# Verify the local scorer scale for candidates whose path changes in the
# fixed-T P/B contrast. The first stack word is printed raw, not retyped.
break *0x100182e0
commands
  silent
  set $raw_unit = *(unsigned int *)($esp + 4)
  if $raw_unit == 3837 || $raw_unit == 152011 || $raw_unit == 152012 || $raw_unit == 3838
    set $scorer_calls = $scorer_calls + 1
    set $target = *(unsigned int *)($esp + 12)
    set $feature_view = *(unsigned int *)($esp + 16)
    set $scale = *(float *)($esp + 20)
    set $ret = *(unsigned int *)$esp
    printf "FOCUSED_UNIT_SCORE_ENTRY context=%d raw_unit=%u scale=%g target=", $active_context, $raw_unit, $scale
    x/7bx $target
    printf "FOCUSED_UNIT_SCORE_VIEW="
    x/10bx $feature_view
    tbreak *$ret
    commands
      silent
      printf "FOCUSED_UNIT_SCORE_RETURN context=%d raw_unit=%u cost=%g\n", $active_context, $raw_unit, $st0
      continue
    end
  end
  continue
end

# Dynamic-programming score for each current candidate against prior context.
break *0x10018c80
commands
  silent
  set $pair_calls = $pair_calls + 1
  set $context = *(int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $record = $state + 0xae894 + $context * 0xfc
  set $previous = $state + 0xae894 + ($context - 1) * 0xfc
  set $count = *(unsigned short *)($record + 0xf4)
  set $previous_count = *(unsigned short *)($previous + 0xf4)
  set $ret = *(unsigned int *)$esp
  printf "TRANSITION_ENTRY context=%d current_count=%u previous_count=%u\n", $context, $count, $previous_count
  set $i = 0
  while $i < $count && $i < 30
    set $node = *(unsigned int *)($state + 0x477ac + $i * 4)
    printf "TRANSITION_INPUT context=%d index=%d unit=%u local_cost=%g\n", $context, $i, *(unsigned int *)($node + 8), *(float *)$node
    set $i = $i + 1
  end
  tbreak *$ret
  commands
    silent
    set $count = *(unsigned short *)($record + 0xf4)
    printf "TRANSITION_RETURN context=%d count=%u\n", $context, $count
    set $i = 0
    while $i < $count && $i < 30
      set $slot = $context * 30 + $i
      set $node = *(unsigned int *)($state + 0x477ac + $i * 4)
      printf "TRANSITION_COST context=%d index=%d unit=%u predecessor_index=%d cumulative_cost=%g\n", $context, $i, *(unsigned int *)($node + 8), *(short *)($state + 0x9f664 + $slot * 2), *(float *)($state + 0x8212c + $slot * 4)
      set $i = $i + 1
    end
    continue
  end
  continue
end

# Pruning result and retained candidate IDs.
break *0x10024510
commands
  silent
  set $prune_calls = $prune_calls + 1
  set $context = *(int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $record = $state + 0xae894 + $context * 0xfc
  set $count = *(unsigned short *)($record + 0xf4)
  set $ret = *(unsigned int *)$esp
  printf "PRUNE_ENTRY context=%d count=%u\n", $context, $count
  tbreak *$ret
  commands
    silent
    set $count = *(unsigned short *)($record + 0xf4)
    printf "PRUNE_RETURN context=%d count=%u ids=", $context, $count
    set $i = 0
    while $i < $count && $i < 30
      printf "%u ", *(unsigned int *)($record + 0x7c + $i * 4)
      set $i = $i + 1
    end
    printf "\n"
    continue
  end
  continue
end

# Chosen candidate index and cumulative cost for each backtracked context.
break *0x10024900
commands
  silent
  set $backtrack_calls = $backtrack_calls + 1
  set $first = *(int *)($esp + 4)
  set $last = *(int *)($esp + 8)
  set $state = *(unsigned int *)($esp + 12)
  set $ret = *(unsigned int *)$esp
  printf "BACKTRACK_ENTRY first=%d last=%d\n", $first, $last
  tbreak *$ret
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
      printf "BACKTRACK_RESULT context=%d index=%d count=%u unit=%d cumulative_cost=%g\n", $context, $index, $count, $unit, $cost
      set $context = $context + 1
    end
    continue
  end
  continue
end

continue
