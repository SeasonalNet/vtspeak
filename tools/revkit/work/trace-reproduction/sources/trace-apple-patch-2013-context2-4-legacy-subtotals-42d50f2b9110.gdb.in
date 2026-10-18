set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

# Replace the five measured context-2 and context-4 pair-feature subtotals
# with native 2006 values. Context 3 and all other 2013 scoring stays intact.
set logging file /work/corpus-parity/stage20/apple-selector-byte0-context2-4-legacy-subtotals-2026-09-29/adapted-context2-4-legacy-subtotals-gdb.log
set logging overwrite on
set logging enabled on

hbreak *0x100191bc
commands
  silent
  set $context = *(int *)($ebp + 8)
  set $current = *(unsigned int *)($ebp - 8)
  set $previous_ptr = *(unsigned int *)($ebp - 0x18)
  set $previous = *(unsigned int *)$previous_ptr
  set $replacement = -1.0
  if ($context == 2) && ($current == 100180) && ($previous == 177774)
    set $replacement = 6.84252
  end
  if ($context == 2) && ($current == 64719) && ($previous == 177774)
    set $replacement = 5.73442
  end
  if ($context == 4) && ((($current == 59559) && ($previous == 59558)) || (($current == 82037) && ($previous == 82036)) || (($current == 2555) && ($previous == 2554)))
    set $replacement = 0.0
  end
  if $replacement >= 0.0
    set $feature_address = $ebp - 0x4c
    set $before = *(float *)$feature_address
    set *(float *)$feature_address = $replacement
    printf "APPLE_2013_CONTEXT2_4_LEGACY_SUBTOTAL context=%d current=%u previous=%u before=%g after=%g\n", $context, $current, $previous, $before, $replacement
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  printf "APPLE_2013_CONTEXT2_4_LEGACY_SUBTOTAL_SELECTED id=%u\n", *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_2013_CONTEXT2_4_LEGACY_SUBTOTAL_TRACE_READY\n"
  continue
end

continue
