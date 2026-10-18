set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-repeat-2013-slot0-pool-to-a000-v4/adapted-score-precision-gdb.log
set logging overwrite on
set logging enabled on

break *0x1002389a
condition $bpnum (*(int *)($ebp + 8) == 0)
commands
  silent
  set $node = *(unsigned int *)$esi
  set $id = *(unsigned int *)($node + 8)
  if $id == 273369 || $id == 149799
    printf "SLOT0_ALIAS_SCORE id=%u local=%.12g node_key=%.12g node_score=%.12g span=%d weighted=%d\n", $id, $st0, *(float *)($node + 4), *(float *)$node, *(short *)($node + 0xe), *(short *)($node + 0x10)
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "SLOT0_ALIAS_PRECISION_TRACE_READY\n"
  continue
end

continue
