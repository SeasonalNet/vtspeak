set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hi-2013-pair-control/adapted-pair-control-gdb.log
set logging overwrite on
set logging enabled on
set $selected_calls = 0

# Record the ordered candidates entering each context's unit scorer. This
# control intentionally avoids breakpoints inside the quadratic pair pass.
break *0x100235f6
commands
  silent
  set $state = *(unsigned int *)($ebp + 0xc)
  set $position = *(int *)($ebp + 8)
  set $count = *(short *)($state + 0xae988 + $position * 0xfc)
  printf "HI_CONTROL_CANDIDATE_LIST position=%d count=%d ids=", $position, $count
  set $index = 0
  while $index < $count && $index < 40
    set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
    printf "%u ", *(unsigned int *)($node + 8)
    set $index = $index + 1
  end
  printf "\n"
  continue
end

hbreak *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  printf "HI_CONTROL_SELECTED n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "HI_PAIR_CONTROL_TRACE_READY\n"
  continue
end

continue
