set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/hello-repeat-2006-pairmatrix/legacy-pairmatrix-gdb.log
set logging overwrite on
set logging enabled on

hbreak *0x1001ec07
commands
  silent
  set $ctx = *(int *)($esp + 0x2c)
  set $current = *(unsigned int *)($esp + 0x3c)
  set $previous = $ebp
  if ($ctx == 1) && ($current == 272823)
    printf "LEGACY_PAIR_COST position=%d current=%u previous=%u transition=%g previous_path=%g local=%g total=%g penalty=%g\n", $ctx, $current, $previous, *(float *)($esp + 0x44), *(float *)(*(unsigned int *)($esp + 0x6c) - 4), *(float *)($esp + 0x94), $st0, *(float *)($esp + 0x14)
  end
  continue
end

break *0x408187
commands
  silent
  printf "LEGACY_PAIR_COST_TRACE_READY\n"
  continue
end

continue
