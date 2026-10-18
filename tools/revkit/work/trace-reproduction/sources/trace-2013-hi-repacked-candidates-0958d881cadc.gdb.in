set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $selected_calls = 0

break *0x100235f6
commands
  silent
  set $state = *(unsigned int *)($ebp + 0xc)
  set $position = *(int *)($ebp + 8)
  set $count = *(short *)($state + 0xae988 + $position * 0xfc)
  printf "HI_2013_CANDIDATE_LIST position=%d count=%d", $position, $count
  set $index = 0
  while $index < $count
    set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
    printf " [%d]=%u(span=%d,weighted=%d)", $index, *(unsigned int *)($node + 8), *(short *)($node + 0xe), *(short *)($node + 0x10)
    set $index = $index + 1
  end
  printf "\n"
  continue
end

hbreak *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  printf "HI_2013_SELECTED_UNIT n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "HI_2013_REPACKED_TRACE_READY\n"
  continue
end

continue
