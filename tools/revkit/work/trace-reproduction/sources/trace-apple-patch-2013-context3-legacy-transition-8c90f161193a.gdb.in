set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

# Test sufficiency of the context-3 pair-feature change alone. Values are
# the measured 2006 subtotals for the aligned candidate edges; every other
# 2013 transition and local score remains unchanged.
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
  if $replacement >= 0.0
    set $feature_address = $ebp - 0x4c
    set $before = *(float *)$feature_address
    set *(float *)$feature_address = $replacement
    printf "APPLE_2013_CONTEXT3_LEGACY_FEATURE context=%d current=%u previous=%u before=%g after=%g\n", $context, $current, $previous, $before, $replacement
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  printf "APPLE_2013_CONTEXT3_SELECTED id=%u\n", *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_2013_CONTEXT3_TRACE_READY\n"
  continue
end

continue
