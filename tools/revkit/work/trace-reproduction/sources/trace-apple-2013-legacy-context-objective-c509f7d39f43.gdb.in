set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/apple-selector-byte0-legacy-context-objective-2026-09-29/adapted-context-objective-gdb.log
set logging overwrite on
set logging enabled on

# Replace the 2013 pair subtotal for every edge in Apple contexts 2 and 3
# with the measured 2006 raw/A/B schedules. At context 4, apply the old
# consecutive-unit zero-cost branch to every adjacent global-ID pair.
# This experiment deliberately does not special-case candidate IDs.
hbreak *0x100191bc
commands
  silent
  set $context = *(int *)($ebp + 8)
  set $current = *(unsigned int *)($ebp - 8)
  set $previous_ptr = *(unsigned int *)($ebp - 0x18)
  set $previous = *(unsigned int *)$previous_ptr
  set $feature_address = $ebp - 0x4c
  set $before = *(float *)$feature_address
  set $raw = *(float *)($ebp - 0x38)
  set $derived_a = *(float *)($ebp - 0x34)
  set $derived_b = *(float *)($ebp - 0x30)
  set $changed = 0
  if $context == 2
    set $after = $raw * 10.0 + $derived_a * 1.0 + $derived_b * 5.0 + 2.0
    set *(float *)$feature_address = $after
    set $changed = 1
  end
  if $context == 3
    set $after = $raw * 10.0 + $derived_a * 2.0 + $derived_b * 10.0 + 2.0
    set *(float *)$feature_address = $after
    set $changed = 1
  end
  if ($context == 4) && ($current == $previous + 1)
    set $after = 0.0
    set *(float *)$feature_address = $after
    set $changed = 1
  end
  if $changed
    printf "APPLE_2013_LEGACY_CONTEXT_OBJECTIVE context=%d current=%u previous=%u raw=%g derived_a=%g derived_b=%g before=%g after=%g\n", $context, $current, $previous, $raw, $derived_a, $derived_b, $before, $after
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  printf "APPLE_2013_LEGACY_CONTEXT_OBJECTIVE_SELECTED id=%u\n", *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_2013_LEGACY_CONTEXT_OBJECTIVE_TRACE_READY\n"
  continue
end

continue
