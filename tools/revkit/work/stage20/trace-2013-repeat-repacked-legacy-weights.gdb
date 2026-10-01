set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-repeat-2013-repacked-legacy-weights/adapted-legacy-weights-gdb.log
set logging overwrite on
set logging enabled on
set $weight_override_done = 0
set $selected_calls = 0

# Coefficient-only probe: preserve candidate generation and local scores, then
# set category-1's effective coefficients to 10,10,2 in the 2013 metric-slot
# order. This does not undo 2013's swapped second/third metric slots, so it is
# not the complete 2006 transition formula.
break *0x10018c80
commands
  silent
  set $context = *(int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  if $weight_override_done == 0
    set $row = 0x1007c288 + 12
    printf "REPACKED_WEIGHT_OVERRIDE before=%g,%g,%g context=%d\n", *(float *)$row, *(float *)($row + 4), *(float *)($row + 8), $context
    set {float} $row = 10.0
    set {float} ($row + 4) = 2.0
    set {float} ($row + 8) = 10.0
    set $weight_override_done = 1
    printf "REPACKED_WEIGHT_OVERRIDE after=%g,%g,%g\n", *(float *)$row, *(float *)($row + 4), *(float *)($row + 8)
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
    set $feature = *(float *)($ebp - 0x4c)
    set $divisor = *(int *)($ebp - 0x10)
    set $penalty = *(int *)($ebp - 0xc)
    set $previous_cost_ptr = *(unsigned int *)($ebp - 0x54)
    set $previous_cost = *(float *)$previous_cost_ptr
    set $local_cost = *(float *)($ebp - 0x74)
    printf "REPACKED_OLD_WEIGHT_PAIR context=%d current=%u previous=%u feature=%g divisor=%d penalty=%d previous_path=%g local=%g total=%g\n", $pair_context, $current, $previous, $feature, $divisor, $penalty, $previous_cost, $local_cost, $feature / $divisor + $penalty + $previous_cost + $local_cost
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  printf "REPACKED_OLD_WEIGHT_SELECTED n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "REPACKED_OLD_WEIGHT_TRACE_READY\n"
  continue
end

continue
