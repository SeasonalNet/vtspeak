set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-repeat-2013-slot0-pool-to-a000-v2/adapted-candidate-scores-gdb.log
set logging overwrite on
set logging enabled on

break *0x100235f6
condition $bpnum (*(int *)($ebp + 8) == 0)
commands
  silent
  set $state = *(unsigned int *)($ebp + 0xc)
  set $count = *(short *)($state + 0xae988)
  printf "SLOT0_ALIAS_PRE_SCORE count=%d", $count
  set $index = 0
  while $index < $count && $index < 100
    set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
    printf " [%d]=%u(rank=%g,span=%d,weighted=%d)", $index, *(unsigned int *)($node + 8), *(float *)($node + 4), *(short *)($node + 0xe), *(short *)($node + 0x10)
    set $index = $index + 1
  end
  printf "\n"
  continue
end

break *0x1002389a
condition $bpnum (*(int *)($ebp + 8) == 0)
commands
  silent
  set $node = *(unsigned int *)$esi
  printf "SLOT0_ALIAS_LOCAL id=%u rank=%g span=%d weighted=%d cost=%g\n", *(unsigned int *)($node + 8), *(float *)($node + 4), *(short *)($node + 0xe), *(short *)($node + 0x10), $st0
  continue
end

hbreak *0x408187
commands
  silent
  printf "SLOT0_ALIAS_SCORE_TRACE_READY\n"
  continue
end

continue
