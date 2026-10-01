set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hi-2013-local-scores/adapted-local-scores-gdb.log
set logging overwrite on
set logging enabled on
set $selected_calls = 0

break *0x1002389a
commands
  silent
  set $position = *(int *)($ebp + 8)
  set $node = *(unsigned int *)$esi
  set $unit = *(unsigned int *)($node + 8)
  if $position == 0 && ($unit == 272820 || $unit == 65331)
    printf "HI_2013_LOCAL_SCORE position=%d unit=%u key=%g span=%d weighted=%d score=%g\n", $position, $unit, *(float *)($node + 4), *(short *)($node + 0xe), *(short *)($node + 0x10), $st0
  end
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
  printf "HI_2013_LOCAL_SCORE_TRACE_READY\n"
  continue
end

continue
