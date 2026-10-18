set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/apple-selector-byte0-context3-weights-bit7-continuity-2026-09-29/adapted-context3-weights-bit7-continuity-gdb.log
set logging overwrite on
set logging enabled on
set $marker_active = 0
set $marker_edges = 0
set $context3_edges = 0
set $selected_count = 0

# At the existing 2013 adjacent-row branch, temporarily supply bit 0x80 on
# the previous signature row for the three legacy-marked Apple pairs. Restore
# the byte immediately after the branch decision. The edge scorer then runs
# its own zero-term path. Apply the old measured coefficients to every context-3
# edge, leaving context 2 and all other scoring unchanged.
hbreak *0x10018ed3
commands
  silent
  set $current = $edx
  set $previous = $edi
  if (($current == 59559) && ($previous == 59558)) || (($current == 82037) && ($previous == 82036)) || (($current == 2555) && ($previous == 2554))
    set $marker_pointer = $esi + 6
    set $marker_original = *(unsigned char *)$marker_pointer
    set $marker_updated = $marker_original | 0x80
    set *(unsigned char *)$marker_pointer = $marker_updated
    set $marker_active = 1
    set $marker_edges = $marker_edges + 1
    printf "APPLE_2013_CONTINUITY_BRANCH_INPUT current=%u previous=%u previous_signature_byte6_before=%u supplied=%u\n", $current, $previous, $marker_original, $marker_updated
  end
  continue
end

hbreak *0x10018ed9
commands
  silent
  if $marker_active
    printf "APPLE_2013_CONTINUITY_BRANCH_TAKEN current=%u previous=%u previous_signature_byte6=%u restoring=%u\n", $edx, $edi, *(unsigned char *)$marker_pointer, $marker_original
    set *(unsigned char *)$marker_pointer = $marker_original
    set $marker_active = 0
  end
  continue
end

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
    set $current = *(unsigned int *)($ebp - 8)
    set $previous_ptr = *(unsigned int *)($ebp - 0x18)
    set $previous = *(unsigned int *)$previous_ptr
    if (($current == 59558) && ($previous == 100180)) || (($current == 2554) && ($previous == 64719)) || (($current == 82036) && ($previous == 64719))
      printf "APPLE_2013_CONTEXT3_OLD_SCHEDULE_EDGE current=%u previous=%u replacement=%g\n", $current, $previous, *(float *)($ebp - 0x4c)
    end
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  set $selected_count = $selected_count + 1
  if $selected_count == 1
    printf "APPLE_2013_CONTEXT3_EDGE_CALLS=%u CONTINUITY_BRANCHES=%u\n", $context3_edges, $marker_edges
  end
  if $selected_count <= 15
    printf "APPLE_2013_CONTEXT3_WEIGHTS_BIT7_CONTINUITY_SELECTED n=%u id=%u\n", $selected_count, *(unsigned int *)($esp + 8)
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_2013_CONTEXT3_WEIGHTS_BIT7_CONTINUITY_TRACE_READY\n"
  continue
end

continue
