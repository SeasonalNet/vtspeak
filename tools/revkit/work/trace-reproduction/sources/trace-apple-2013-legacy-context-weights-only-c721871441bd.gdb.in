set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/apple-selector-byte0-legacy-context-weights-only-2026-09-29/adapted-context-weights-only-gdb.log
set logging overwrite on
set logging enabled on

# Apply the measured old raw/A/B schedules to every 2013 pair edge in
# contexts 2 and 3. Leave context 4 and every categorical/local term intact.
hbreak *0x100191bc
commands
  silent
  set $context = *(int *)($ebp + 8)
  set $raw = *(float *)($ebp - 0x38)
  set $derived_a = *(float *)($ebp - 0x34)
  set $derived_b = *(float *)($ebp - 0x30)
  if $context == 2
    set *(float *)($ebp - 0x4c) = $raw * 10.0 + $derived_a + $derived_b * 5.0 + 2.0
  end
  if $context == 3
    set *(float *)($ebp - 0x4c) = $raw * 10.0 + $derived_a * 2.0 + $derived_b * 10.0 + 2.0
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  printf "APPLE_2013_LEGACY_CONTEXT_WEIGHTS_ONLY_SELECTED id=%u\n", *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_2013_LEGACY_CONTEXT_WEIGHTS_ONLY_TRACE_READY\n"
  continue
end

continue
