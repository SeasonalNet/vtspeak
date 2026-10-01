set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/apple-selector-byte0-context3-weights-only-2026-09-29/adapted-context3-weights-only-gdb.log
set logging overwrite on
set logging enabled on
set $context3_edges = 0
set $selected_count = 0

# Apply the measured 2006 context-3 coefficient schedule to every 2013 edge
# at context 3, with the adapted index left unmarked and all other weights
# unchanged.
hbreak *0x100191bc
commands
  silent
  set $context = *(int *)($ebp + 8)
  if $context == 3
    set $raw = *(float *)($ebp - 0x38)
    set $derived_a = *(float *)($ebp - 0x34)
    set $derived_b = *(float *)($ebp - 0x30)
    set *(float *)($ebp - 0x4c) = $raw * 10.0 + $derived_a * 2.0 + $derived_b * 10.0 + 2.0
    set $context3_edges = $context3_edges + 1
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  set $selected_count = $selected_count + 1
  if $selected_count <= 15
    printf "APPLE_2013_CONTEXT3_WEIGHTS_ONLY_SELECTED n=%u id=%u\n", $selected_count, *(unsigned int *)($esp + 8)
  end
  continue
end

continue
