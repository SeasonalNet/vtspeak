set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-repeat-2013-completed-path-rows/adapted-path-rows-gdb.log
set logging overwrite on
set logging enabled on

# Keep the 2006 first choice in the controlled 2013 run.
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
      if *(unsigned int *)($node + 8) == 273369
        set {float}($node + 4) = -2.0
      end
      set $index = $index + 1
    end
  end
  continue
end

# FUN_10018c80 stores the best cumulative cost and predecessor index for each
# retained transition candidate. Dump each completed row before backtracking.
break *0x10018c80
commands
  silent
  set $context = *(int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $return = *(unsigned int *)$esp
  if $context >= 1 && $context <= 11
    set $record = $state + 0xae894 + $context * 0xfc
    set $previous_record = $record - 0xfc
    tbreak *$return
    commands
      silent
      set $count = *(short *)($record + 0xf4)
      set $previous_count = *(short *)($previous_record + 0xf4)
      printf "COMPLETED_PATH_ROW context=%d candidates=%d predecessors=%d\n", $context, $count, $previous_count
      set $i = 0
      while $i < $count && $i < 30
        set $unit = *(unsigned int *)($record + 0x7c + $i * 4)
        set $cost = *(float *)($state + 0x8212c + ($context * 30 + $i) * 4)
        set $pred_index = *(short *)($state + 0x9f664 + ($context * 30 + $i) * 2)
        if $pred_index >= 0 && $pred_index < $previous_count
          set $pred = *(unsigned int *)($previous_record + 0x7c + $pred_index * 4)
          printf "COMPLETED_PATH context=%d unit=%u cumulative=%g predecessor_index=%d predecessor=%u\n", $context, $unit, $cost, $pred_index, $pred
        else
          printf "COMPLETED_PATH context=%d unit=%u cumulative=%g predecessor_index=%d predecessor=OUT_OF_RANGE\n", $context, $unit, $cost, $pred_index
        end
        set $i = $i + 1
      end
      continue
    end
  end
  continue
end

set $selected_calls = 0
break *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  if $selected_calls <= 16
    printf "COMPLETED_PATH_SELECTED context=%u unit=%u\n", $selected_calls - 1, *(unsigned int *)($esp + 8)
  end
  continue
end

break *0x408187
commands
  silent
  printf "COMPLETED_PATH_TRACE_READY\n"
  continue
end

continue
