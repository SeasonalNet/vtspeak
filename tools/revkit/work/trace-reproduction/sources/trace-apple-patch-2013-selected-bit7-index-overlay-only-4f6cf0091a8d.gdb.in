set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/apple-selector-byte0-selected-bit7-only-2026-09-29/adapted-selected-bit7-only-gdb.log
set logging overwrite on
set logging enabled on
set $selected_count = 0

# Leave all 2013 transition weights unchanged; observe only selected units
# while the disposable index contains the three targeted continuation bits.
hbreak *0x10018ed3
commands
  silent
  if (($edx == 59559) && ($edi == 59558)) || (($edx == 82037) && ($edi == 82036)) || (($edx == 2555) && ($edi == 2554))
    printf "APPLE_2013_MARKER_ONLY_BRANCH_INPUT current=%u previous=%u signature_byte6=%u\n", $edx, $edi, *(unsigned char *)($esi + 6)
  end
  continue
end

hbreak *0x10018ed9
commands
  silent
  if (($edx == 59559) && ($edi == 59558)) || (($edx == 82037) && ($edi == 82036)) || (($edx == 2555) && ($edi == 2554))
    printf "APPLE_2013_MARKER_ONLY_BRANCH_TAKEN current=%u previous=%u\n", $edx, $edi
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  set $selected_count = $selected_count + 1
  if $selected_count <= 15
    printf "APPLE_2013_MARKER_ONLY_SELECTED n=%u id=%u\n", $selected_count, *(unsigned int *)($esp + 8)
  end
  continue
end

continue
