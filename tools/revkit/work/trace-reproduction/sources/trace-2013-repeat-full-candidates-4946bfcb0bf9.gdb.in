set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x100235f6
commands
  silent
  set $state = *(unsigned int *)($ebp + 0xc)
  set $position = *(int *)($ebp + 8)
  set $count = *(short *)($state + 0xae988 + $position * 0xfc)
  printf "BASELINE_2013_CANDIDATE_LIST position=%d count=%d", $position, $count
  set $index = 0
  while $index < $count
    set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
    printf " [%d]=%u(span=%d,weighted=%d,key=%g)", $index, *(unsigned int *)($node + 8), *(short *)($node + 0xe), *(short *)($node + 0x10), *(float *)($node + 4)
    set $index = $index + 1
  end
  printf "\n"
  continue
end

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
    printf "BASELINE_2013_TRANSITION context=%d current_count=%d current_units:", $context, $count
    set $index = 0
    while $index < $count
      printf " %u", *(unsigned int *)($record + 0x7c + $index * 4)
      set $index = $index + 1
    end
    printf "\n"
    continue
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "BASELINE_2013_FULL_CANDIDATE_TRACE_READY\n"
  continue
end

continue
