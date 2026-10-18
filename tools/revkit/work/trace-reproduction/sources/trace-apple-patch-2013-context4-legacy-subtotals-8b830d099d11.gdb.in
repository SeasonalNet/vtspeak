set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

# Replace only the three measured context-4 pair-feature subtotals with the
# native 2006 zero return. All candidate generation, other edges, local costs,
# and transition selection remain in the 2013 DLL.
set logging file /work/corpus-parity/stage20/apple-selector-byte0-context4-legacy-subtotals-2026-09-29/adapted-context4-legacy-subtotals-gdb.log
set logging overwrite on
set logging enabled on

hbreak *0x100191bc
commands
  silent
  set $context = *(int *)($ebp + 8)
  set $current = *(unsigned int *)($ebp - 8)
  set $previous_ptr = *(unsigned int *)($ebp - 0x18)
  set $previous = *(unsigned int *)$previous_ptr
  if ($context == 4) && ((($current == 59559) && ($previous == 59558)) || (($current == 82037) && ($previous == 82036)) || (($current == 2555) && ($previous == 2554)))
    set $feature_address = $ebp - 0x4c
    set $before = *(float *)$feature_address
    set *(float *)$feature_address = 0.0
    printf "APPLE_2013_CONTEXT4_LEGACY_SUBTOTAL current=%u previous=%u before=%g after=0\n", $current, $previous, $before
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  printf "APPLE_2013_CONTEXT4_LEGACY_SUBTOTAL_SELECTED id=%u\n", *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_2013_CONTEXT4_LEGACY_SUBTOTAL_TRACE_READY\n"
  continue
end

continue
