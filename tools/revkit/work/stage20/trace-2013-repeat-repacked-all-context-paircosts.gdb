set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-repeat-2013-repacked-all-context-paircosts/adapted-pairmatrix-gdb.log
set logging overwrite on
set logging enabled on

hbreak *0x100191bc
commands
  silent
  set $context = *(int *)($ebp + 8)
  set $current = *(unsigned int *)($ebp - 8)
  if $current == 272823
    set $previous_ptr = *(unsigned int *)($ebp - 0x18)
    set $previous = *(unsigned int *)$previous_ptr
    set $feature = *(float *)($ebp - 0x4c)
    set $divisor = *(int *)($ebp - 0x10)
    set $penalty = *(int *)($ebp - 0xc)
    set $previous_cost_ptr = *(unsigned int *)($ebp - 0x54)
    set $previous_cost = *(float *)$previous_cost_ptr
    set $local_cost = *(float *)($ebp - 0x74)
    printf "REPACKED_NEW_PAIR_COST context=%d current=%u previous=%u feature=%g divisor=%d penalty=%d previous_path=%g local=%g total=%g\n", $context, $current, $previous, $feature, $divisor, $penalty, $previous_cost, $local_cost, $feature / $divisor + $penalty + $previous_cost + $local_cost
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "REPACKED_NEW_ALL_CONTEXT_TRACE_READY\n"
  continue
end

continue
