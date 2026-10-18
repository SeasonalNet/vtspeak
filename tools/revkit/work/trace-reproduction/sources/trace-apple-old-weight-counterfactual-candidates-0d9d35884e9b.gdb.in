set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/apple-selector-byte0-old-weight-counterfactual-candidates-2026-09-29/legacy-weight-counterfactual-candidates-gdb.log
set logging overwrite on
set logging enabled on

# Check whether the alternate 2013 chain admitted by legacy weights exists in
# the native 2006 transition graph, and capture its native edge terms.
hbreak *0x1001ec07
commands
  silent
  set $context = *(int *)($esp + 0x2c)
  set $current = *(unsigned int *)($esp + 0x3c)
  set $previous = $ebp
  if (($context == 2) && ($current == 188826)) || (($context == 3) && (($current == 129559) || ($current == 129560))) || (($context == 4) && (($current == 129559) || ($current == 129560)))
    printf "APPLE_OLD_COUNTERFACTUAL_EDGE context=%d current=%u previous=%u transition=%g previous_path=%g local=%g total=%g penalty=%g\n", $context, $current, $previous, *(float *)($esp + 0x44), *(float *)(*(unsigned int *)($esp + 0x6c) - 4), *(float *)($esp + 0x94), $st0, *(float *)($esp + 0x14)
  end
  continue
end

hbreak *0x1001ed7f
commands
  silent
  printf "APPLE_OLD_COUNTERFACTUAL_SELECTED id=%u\n", $edx
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_OLD_COUNTERFACTUAL_TRACE_READY\n"
  continue
end

continue
