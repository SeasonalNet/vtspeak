set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/apple-selector-byte0-context3-bit7-weights-debug-2026-09-29/adapted-context3-bit7-weights-debug-gdb.log
set logging overwrite on
set logging enabled on
set $marker_active = 0
set $marker_edges = 0
set $context3_edges = 0
set $scoring_calls = 0
set $selected_count = 0

# Feed bit 0x80 to the existing 2013 adjacent-row branch for the three known
# legacy continuation edges, then restore the source byte immediately.
hbreak *0x10018ed3
commands
  silent
  if (($edx == 59559) && ($edi == 59558)) || (($edx == 82037) && ($edi == 82036)) || (($edx == 2555) && ($edi == 2554))
    set $marker_pointer = $esi + 6
    set $marker_original = *(unsigned char *)$marker_pointer
    set *(unsigned char *)$marker_pointer = $marker_original | 0x80
    set $marker_active = 1
    set $marker_edges = $marker_edges + 1
    printf "APPLE_2013_MARKER current=%u previous=%u before=%u after=%u\n", $edx, $edi, $marker_original, *(unsigned char *)$marker_pointer
  end
  continue
end

hbreak *0x10018ed9
commands
  silent
  if $marker_active
    printf "APPLE_2013_MARKER_BRANCH_TAKEN current=%u previous=%u\n", $edx, $edi
    set *(unsigned char *)$marker_pointer = $marker_original
    set $marker_active = 0
  end
  continue
end

hbreak *0x100191bc
commands
  silent
  set $scoring_calls = $scoring_calls + 1
  set $context = *(int *)($ebp + 8)
  set $current = *(unsigned int *)($ebp - 8)
  set $previous_ptr = *(unsigned int *)($ebp - 0x18)
  set $previous = *(unsigned int *)$previous_ptr
  if $scoring_calls <= 30
    printf "APPLE_2013_SCORER_CALL n=%u context=%d current=%u previous=%u feature=%g\n", $scoring_calls, $context, $current, $previous, *(float *)($ebp - 0x4c)
  end
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
  if $selected_count == 1
    printf "APPLE_2013_DEBUG_SUMMARY calls=%u context3_edges=%u marker_edges=%u\n", $scoring_calls, $context3_edges, $marker_edges
  end
  if $selected_count <= 15
    printf "APPLE_2013_DEBUG_SELECTED n=%u id=%u\n", $selected_count, *(unsigned int *)($esp + 8)
  end
  continue
end

continue
