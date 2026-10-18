set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $transition_calls = 0
set $backtrack_calls = 0

hbreak *0x10018c80
commands
  silent
  set $transition_calls = $transition_calls + 1
  set $context = *(int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $record = $state + 0xae894 + $context * 0xfc
  set $previous = $record - 0xfc
  set $return = *(unsigned int *)$esp
  printf "APPLE_TRANSITION_ENTRY n=%u context=%d current_count=%u previous_count=%u\n", $transition_calls, $context, *(unsigned short *)($record + 0xf4), *(unsigned short *)($previous + 0xf4)
  thbreak *$return
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

hbreak *0x10024900
commands
  silent
  set $backtrack_calls = $backtrack_calls + 1
  set $first = *(int *)($esp + 4)
  set $last = *(int *)($esp + 8)
  set $state = *(unsigned int *)($esp + 12)
  set $return = *(unsigned int *)$esp
  printf "APPLE_BACKTRACK_ENTRY n=%u first=%d last=%d\n", $backtrack_calls, $first, $last
  thbreak *$return
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

continue
