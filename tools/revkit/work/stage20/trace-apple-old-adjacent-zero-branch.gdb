set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/apple-selector-byte0-old-adjacent-zero-branch-2026-09-29/legacy-adjacent-zero-branch-gdb.log
set logging overwrite on
set logging enabled on

# 0x1001e952 is the old FUN_1001e470 branch target that writes a zero pair
# feature term. Capture the exact row IDs, selector byte, and prior-row
# continuation byte for the three late Apple candidate pairs.
hbreak *0x1001e952
commands
  silent
  set $current = *(unsigned int *)($esp + 0x3c)
  set $previous = $ebp
  set $state = *(unsigned int *)($esp + 0x28)
  set $selector = *(unsigned char *)($state + 5)
  set $continuation_table = *(unsigned int *)($ebx + 0x68)
  set $continuation = *(unsigned char *)($continuation_table + $previous)
  if (($current == 59559) && ($previous == 59558)) || (($current == 82037) && ($previous == 82036)) || (($current == 2555) && ($previous == 2554))
    printf "APPLE_2006_ZERO_BRANCH current=%u previous=%u selector_byte=%u previous_continuation=%u\n", $current, $previous, $selector, $continuation
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_2006_ADJACENT_ZERO_BRANCH_TRACE_READY\n"
  continue
end

continue
