set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/apple-selector-byte0-context3-weights-context4-continuity-2026-09-29/adapted-context3-weights-context4-continuity-gdb.log
set logging overwrite on
set logging enabled on
set $context3_edges = 0
set $context4_edges = 0
set $selected_count = 0

# Apply the measured native 2006 context-3 coefficient schedule to every 2013
# edge at that position. Separately model the verified old adjacent-row
# continuation shortcut for the three compared context-4 pairs.
hbreak *0x100191bc
commands
  silent
  set $context = *(int *)($ebp + 8)
  set $current = *(unsigned int *)($ebp - 8)
  set $previous_ptr = *(unsigned int *)($ebp - 0x18)
  set $previous = *(unsigned int *)$previous_ptr
  if $context == 3
    set $raw = *(float *)($ebp - 0x38)
    set $derived_a = *(float *)($ebp - 0x34)
    set $derived_b = *(float *)($ebp - 0x30)
    set *(float *)($ebp - 0x4c) = $raw * 10.0 + $derived_a * 2.0 + $derived_b * 10.0 + 2.0
    set $context3_edges = $context3_edges + 1
    if (($current == 59558) && ($previous == 100180)) || (($current == 2554) && ($previous == 64719)) || (($current == 82036) && ($previous == 64719))
      printf "APPLE_2013_CONTEXT3_OLD_SCHEDULE_EDGE current=%u previous=%u raw=%g derived_a=%g derived_b=%g replacement=%g\n", $current, $previous, $raw, $derived_a, $derived_b, *(float *)($ebp - 0x4c)
    end
  end
  if ($context == 4) && ((($current == 59559) && ($previous == 59558)) || (($current == 82037) && ($previous == 82036)) || (($current == 2555) && ($previous == 2554)))
    set $before = *(float *)($ebp - 0x4c)
    set *(float *)($ebp - 0x4c) = 0.0
    set $context4_edges = $context4_edges + 1
    printf "APPLE_2013_CONTEXT4_OLD_CONTINUITY_EDGE current=%u previous=%u before=%g after=0\n", $current, $previous, $before
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  set $selected_count = $selected_count + 1
  if $selected_count <= 15
    printf "APPLE_2013_CONTEXT3_OLD_SCHEDULE_CONTEXT4_CONTINUITY_SELECTED n=%u id=%u\n", $selected_count, *(unsigned int *)($esp + 8)
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_2013_CONTEXT3_OLD_SCHEDULE_CONTEXT4_CONTINUITY_COUNTS context3_edges=%u context4_edges=%u\n", $context3_edges, $context4_edges
  printf "APPLE_2013_CONTEXT3_OLD_SCHEDULE_CONTEXT4_CONTINUITY_TRACE_READY\n"
  continue
end

continue
