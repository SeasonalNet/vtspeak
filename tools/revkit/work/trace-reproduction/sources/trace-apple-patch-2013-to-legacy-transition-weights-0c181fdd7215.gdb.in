set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

# Intervention: replace only the 2013 pair-feature subtotal for aligned
# candidate edges with the measured 2006 subtotal. The write is to a live
# stack float in this process and leaves DLL/resource files untouched.
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
    printf "APPLE_2013_LEGACY_FEATURE_OVERRIDE context=%d current=%u previous=%u before=%g after=%g\n", $context, $current, $previous, $before, $replacement
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  printf "APPLE_2013_LEGACY_WEIGHT_SELECTED id=%u\n", *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_2013_LEGACY_WEIGHT_TRACE_READY\n"
  continue
end

continue
