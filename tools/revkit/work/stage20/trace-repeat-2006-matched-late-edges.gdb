set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/hello-repeat-2006-matched-late-edges/legacy-matched-late-edges-gdb.log
set logging overwrite on
set logging enabled on

# Record native 2006 transition scores for later rows that also occur in the
# controlled 2013 repeated-Hello comparison.
hbreak *0x1001ec07
commands
  silent
  set $context = *(int *)($esp + 0x2c)
  set $current = *(unsigned int *)($esp + 0x3c)
  set $previous = $ebp
  if ($current == 232670) || ($current == 266023) || ($current == 264071)
    printf "LEGACY_MATCHED_EDGE context=%d current=%u previous=%u transition=%g previous_path=%g local=%g total=%g penalty=%g\n", $context, $current, $previous, *(float *)($esp + 0x44), *(float *)(*(unsigned int *)($esp + 0x6c) - 4), *(float *)($esp + 0x94), $st0, *(float *)($esp + 0x14)
  end
  continue
end

break *0x408187
commands
  silent
  printf "LEGACY_MATCHED_EDGE_TRACE_READY\n"
  continue
end

continue
