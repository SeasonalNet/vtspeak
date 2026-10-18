set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-repeat-2013-slots0-3-5-8-ranked-pools/adapted-ranked-pools-gdb.log
set logging overwrite on
set logging enabled on

# Promote the 2006-selected first unit without changing its acoustic fields.
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
        printf "ORDER_OVERRIDE id=273369 before=%g", *(float *)($node + 4)
        set {float}($node + 4) = -2.0
        printf " after=%g\n", *(float *)($node + 4)
      end
      set $index = $index + 1
    end
  end
  continue
end

# Capture each position's post-score, pre-transition retained candidate list.
break *0x100235f6
commands
  silent
  set $position = *(int *)($ebp + 8)
  set $state = *(unsigned int *)($ebp + 0xc)
  set $count = *(short *)($state + 0xae988 + $position * 0xfc)
  printf "RANKED_POOL position=%d count=%d", $position, $count
  set $index = 0
  while $index < $count && $index < 100
    set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
    printf " [%d]=%u(local=%g,key=%g,span=%d,weighted=%d)", $index, *(unsigned int *)($node + 8), *(float *)$node, *(float *)($node + 4), *(short *)($node + 0xe), *(short *)($node + 0x10)
    set $index = $index + 1
  end
  printf "\n"
  continue
end

# Record every position's transition candidate set for member-set matching.
break *0x10018c80
commands
  silent
  set $context = *(int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $return = *(unsigned int *)$esp
  set $record = $state + 0xae894 + $context * 0xfc
  tbreak *$return
  commands
    silent
    set $count = *(short *)($record + 0xf4)
    printf "TRANSITION_POOL context=%d count=%d units:", $context, $count
    set $index = 0
    while $index < $count && $index < 100
      printf " %u", *(unsigned int *)($record + 0x7c + $index * 4)
      set $index = $index + 1
    end
    printf "\n"
    continue
  end
  continue
end

# Capture the selected backtracked unit at every context.
break *0x10024900
commands
  silent
  set $first = *(int *)($esp + 4)
  set $last = *(int *)($esp + 8)
  set $state = *(unsigned int *)($esp + 12)
  set $return = *(unsigned int *)$esp
  tbreak *$return
  commands
    silent
    set $context = $first
    while $context <= $last
      set $record = $state + 0xae894 + $context * 0xfc
      set $count = *(short *)($record + 0xf4)
      set $index = *(short *)($state + 0xae0c4 + $context * 2)
      set $unit = *(unsigned int *)($record + 0x7c + $index * 4)
      printf "SELECTED context=%d count=%d index=%d unit=%u\n", $context, $count, $index, $unit
      set $context = $context + 1
    end
    continue
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "RANKED_POOL_TRACE_READY\n"
  continue
end

continue
