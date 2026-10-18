set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $weight_override_done = 0
set $backtrack_calls = 0

# The 2013 table applies its second and third entries to metric slots 3 and 2.
# {10,10,2} therefore restores both the old coefficients and the old metric
# ordering for this category-1 transition comparison.
break *0x10018c80
commands
  silent
  if $weight_override_done == 0
    set $row = 0x1007c288 + 12
    printf "FULL_LEGACY_MAP_OVERRIDE before=%g,%g,%g\n", *(float *)$row, *(float *)($row + 4), *(float *)($row + 8)
    set {float} $row = 10.0
    set {float} ($row + 4) = 10.0
    set {float} ($row + 8) = 2.0
    set $weight_override_done = 1
    printf "FULL_LEGACY_MAP_OVERRIDE after=%g,%g,%g\n", *(float *)$row, *(float *)($row + 4), *(float *)($row + 8)
  end
  continue
end

hbreak *0x100191bc
commands
  silent
  set $pair_context = *(int *)($ebp + 8)
  set $current = *(unsigned int *)($ebp - 8)
  if ($pair_context == 2) && ($current == 272823)
    set $previous_ptr = *(unsigned int *)($ebp - 0x18)
    set $previous = *(unsigned int *)$previous_ptr
    if $previous == 273369 || $previous == 272822
      set $feature = *(float *)($ebp - 0x4c)
      set $divisor = *(int *)($ebp - 0x10)
      set $penalty = *(int *)($ebp - 0xc)
      set $previous_cost_ptr = *(unsigned int *)($ebp - 0x54)
      set $previous_cost = *(float *)$previous_cost_ptr
      set $local_cost = *(float *)($ebp - 0x74)
      printf "FULL_LEGACY_MAP_PAIR context=%d current=%u previous=%u feature=%g divisor=%d penalty=%d previous_path=%g local=%g total=%g\n", $pair_context, $current, $previous, $feature, $divisor, $penalty, $previous_cost, $local_cost, $feature / $divisor + $penalty + $previous_cost + $local_cost
    end
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
  tbreak *$return
  commands
    silent
    printf "FULL_LEGACY_MAP_BACKTRACK call=%u first=%d last=%d\n", $backtrack_calls, $first, $last
    set $context = $first
    while $context <= $last
      set $record = $state + 0xae894 + $context * 0xfc
      set $count = *(short *)($state + 0xae988 + $context * 0xfc)
      set $index = *(short *)($state + 0xae0c4 + $context * 2)
      set $unit = *(unsigned int *)($record + 0x7c + $index * 4)
      printf "FULL_LEGACY_MAP_ROW context=%d count=%d selected_index=%d selected_unit=%u\n", $context, $count, $index, $unit
      set $context = $context + 1
    end
    continue
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "FULL_LEGACY_MAP_BACKTRACK_TRACE_READY\n"
  continue
end

continue
