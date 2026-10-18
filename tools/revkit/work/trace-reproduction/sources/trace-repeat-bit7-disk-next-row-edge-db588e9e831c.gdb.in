set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-repeat-2013-bit7-next-row-edge/adapted-next-row-edge-gdb.log
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

# Score the next selected-row candidate under the disk-level continuity correction.
hbreak *0x100191bc
commands
  silent
  set $context = *(int *)($ebp + 8)
  set $current = *(unsigned int *)($ebp - 8)
  set $previous_ptr = *(unsigned int *)($ebp - 0x18)
  set $previous = *(unsigned int *)$previous_ptr
  if ($context == 7) && ($current == 264072)
    set $feature = *(float *)($ebp - 0x4c)
    set $divisor = *(int *)($ebp - 0x10)
    set $penalty = *(int *)($ebp - 0xc)
    set $previous_cost_ptr = *(unsigned int *)($ebp - 0x54)
    set $previous_cost = *(float *)$previous_cost_ptr
    set $local_cost = *(float *)($ebp - 0x74)
    printf "BIT7_NEXT_EDGE context=%d current=%u previous=%u feature=%g divisor=%d penalty=%d previous_path=%g local=%g total=%g\n", $context, $current, $previous, $feature, $divisor, $penalty, $previous_cost, $local_cost, $feature / $divisor + $penalty + $previous_cost + $local_cost
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "BIT7_NEXT_EDGE_TRACE_READY\n"
  continue
end

continue
