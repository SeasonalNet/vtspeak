set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/apple-selector-byte0-context2-legacy-subtotals-2026-09-29/adapted-context2-legacy-subtotals-gdb.log
set logging overwrite on
set logging enabled on

# Replace the two measured context-2 pair-feature subtotals with native 2006
# values. All later context terms and other 2013 scoring stay intact.
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
  if $replacement >= 0.0
    set $feature_address = $ebp - 0x4c
    set $before = *(float *)$feature_address
    set *(float *)$feature_address = $replacement
    printf "APPLE_2013_CONTEXT2_LEGACY_SUBTOTAL current=%u previous=%u before=%g after=%g\n", $current, $previous, $before, $replacement
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  printf "APPLE_2013_CONTEXT2_LEGACY_SUBTOTAL_SELECTED id=%u\n", *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_2013_CONTEXT2_LEGACY_SUBTOTAL_TRACE_READY\n"
  continue
end

continue
