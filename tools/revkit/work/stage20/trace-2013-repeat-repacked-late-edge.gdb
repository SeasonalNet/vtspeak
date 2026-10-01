set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-repeat-2013-repacked-late-edge-control/adapted-late-edge-gdb.log
set logging overwrite on
set logging enabled on

hbreak *0x100191bc
commands
  silent
  set $context = *(int *)($ebp + 8)
  set $current = *(unsigned int *)($ebp - 8)
  if ($context == 9) && ($current == 273370)
    set $previous_ptr = *(unsigned int *)($ebp - 0x18)
    set $previous = *(unsigned int *)$previous_ptr
    if $previous == 181914 || $previous == 76688
      set $category = *(signed char *)($ebp - 1)
      set $feature = *(float *)($ebp - 0x4c)
      set $divisor = *(int *)($ebp - 0x10)
      set $penalty = *(int *)($ebp - 0xc)
      set $previous_cost_ptr = *(unsigned int *)($ebp - 0x54)
      set $previous_cost = *(float *)$previous_cost_ptr
      set $local_cost = *(float *)($ebp - 0x74)
      printf "LATE_EDGE_CONTROL context=%d current=%u previous=%u category=%d feature=%g divisor=%d penalty=%d previous_path=%g local=%g total=%g\n", $context, $current, $previous, $category, $feature, $divisor, $penalty, $previous_cost, $local_cost, $feature / $divisor + $penalty + $previous_cost + $local_cost
    end
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "LATE_EDGE_CONTROL_TRACE_READY\n"
  continue
end

continue
