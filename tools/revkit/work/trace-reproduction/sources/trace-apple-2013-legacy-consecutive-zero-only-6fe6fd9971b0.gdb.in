set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/apple-selector-byte0-legacy-consecutive-zero-only-2026-09-29/adapted-consecutive-zero-only-gdb.log
set logging overwrite on
set logging enabled on

# Apply the legacy zero-cost shortcut to every consecutive global-ID pair in
# context 4. Leave the 2013 feature weights and every other score term intact.
hbreak *0x100191bc
commands
  silent
  set $context = *(int *)($ebp + 8)
  set $current = *(unsigned int *)($ebp - 8)
  set $previous_ptr = *(unsigned int *)($ebp - 0x18)
  set $previous = *(unsigned int *)$previous_ptr
  if ($context == 4) && ($current == $previous + 1)
    set *(float *)($ebp - 0x4c) = 0.0
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  printf "APPLE_2013_LEGACY_CONSECUTIVE_ZERO_ONLY_SELECTED id=%u\n", *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_2013_LEGACY_CONSECUTIVE_ZERO_ONLY_TRACE_READY\n"
  continue
end

continue
