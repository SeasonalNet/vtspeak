set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-repeat-2013-promote-264071/adapted-promoted-264071-gdb.log
set logging overwrite on
set logging enabled on

# Keep the 2006 first choice, then promote its selected unit at position 6
# through the same pre-sort key field used by the earlier order intervention.
break *0x10023331
commands
  silent
  set $position = *(int *)($ebp + 8)
  set $state = *(unsigned int *)($ebp + 0xc)
  if $position == 0 || $position == 6
    set $count = *(short *)($state + 0xae988 + $position * 0xfc)
    set $index = 0
    while $index < $count
      set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
      if $position == 0 && *(unsigned int *)($node + 8) == 273369
        printf "ORDER_OVERRIDE position=0 id=273369 before=%g", *(float *)($node + 4)
        set {float}($node + 4) = -2.0
        printf " after=%g\n", *(float *)($node + 4)
      end
      if $position == 6 && *(unsigned int *)($node + 8) == 264071
        printf "ORDER_OVERRIDE position=6 id=264071 before=%g", *(float *)($node + 4)
        set {float}($node + 4) = -2.0
        printf " after=%g\n", *(float *)($node + 4)
      end
      set $index = $index + 1
    end
  end
  continue
end

# Verify whether the order intervention carries 264071 into the final
# transition shortlist after ranking and the fixed 30-row cap.
break *0x100235f6
commands
  silent
  set $position = *(int *)($ebp + 8)
  set $state = *(unsigned int *)($ebp + 0xc)
  if $position == 6
    set $count = *(short *)($state + 0xae988 + $position * 0xfc)
    printf "PROMOTED_RANKED_POOL position=6 count=%d", $count
    set $index = 0
    while $index < $count && $index < 100
      set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
      printf " [%d]=%u(local=%g,key=%g,span=%d,weighted=%d)", $index, *(unsigned int *)($node + 8), *(float *)$node, *(float *)($node + 4), *(short *)($node + 0xe), *(short *)($node + 0x10)
      set $index = $index + 1
    end
    printf "\n"
  end
  continue
end

# Dump completed path costs and predecessors for all rows.
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
    printf "PROMOTED_SELECTED context=%u unit=%u\n", $selected_calls - 1, *(unsigned int *)($esp + 8)
  end
  continue
end

break *0x408187
commands
  silent
  printf "PROMOTE_264071_TRACE_READY\n"
  continue
end

continue
