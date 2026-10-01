set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/apple-selector-byte0-context3-4-legacy-subtotals-2026-09-29/adapted-context3-4-legacy-subtotals-gdb.log
set logging overwrite on
set logging enabled on

# Replace the three measured context-3 pair-feature subtotals and the three
# context-4 subtotals with native 2006 values. Context 2 remains 2013-scored.
hbreak *0x100191bc
commands
  silent
  set $context = *(int *)($ebp + 8)
  set $current = *(unsigned int *)($ebp - 8)
  set $previous_ptr = *(unsigned int *)($ebp - 0x18)
  set $previous = *(unsigned int *)$previous_ptr
  set $replacement = -1.0
  if ($context == 3) && ($current == 59558) && ($previous == 100180)
    set $replacement = 11.9439
  end
  if ($context == 3) && ($current == 2554) && ($previous == 64719)
    set $replacement = 14.2417
  end
  if ($context == 3) && ($current == 82036) && ($previous == 64719)
    set $replacement = 16.3335
  end
  if ($context == 4) && ((($current == 59559) && ($previous == 59558)) || (($current == 82037) && ($previous == 82036)) || (($current == 2555) && ($previous == 2554)))
    set $replacement = 0.0
  end
  if $replacement >= 0.0
    set $feature_address = $ebp - 0x4c
    set $before = *(float *)$feature_address
    set *(float *)$feature_address = $replacement
    printf "APPLE_2013_CONTEXT3_4_LEGACY_SUBTOTAL context=%d current=%u previous=%u before=%g after=%g\n", $context, $current, $previous, $before, $replacement
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  printf "APPLE_2013_CONTEXT3_4_LEGACY_SUBTOTAL_SELECTED id=%u\n", *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_2013_CONTEXT3_4_LEGACY_SUBTOTAL_TRACE_READY\n"
  continue
end

continue
