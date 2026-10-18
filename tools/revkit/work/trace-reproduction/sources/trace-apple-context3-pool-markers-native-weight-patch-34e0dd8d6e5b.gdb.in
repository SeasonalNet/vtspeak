set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/apple-selector-byte0-context3-native-weight-patch-2026-09-29/adapted-context3-native-weight-patch-gdb.log
set logging overwrite on
set logging enabled on
set $selected_count = 0

# Observe the unmodified 2013 continuation branch and selected IDs. The DLL
# supplies the corrected context-3 coefficient from its patched data table.
hbreak *0x10018ed3
commands
  silent
  if (($edx == 59559) && ($edi == 59558)) || (($edx == 82037) && ($edi == 82036)) || (($edx == 2555) && ($edi == 2554))
    printf "APPLE_PATCHED_DLL_MARKER current=%u previous=%u signature_byte6=%u\n", $edx, $edi, *(unsigned char *)($esi + 6)
  end
  continue
end

hbreak *0x10018ed9
commands
  silent
  if (($edx == 59559) && ($edi == 59558)) || (($edx == 82037) && ($edi == 82036)) || (($edx == 2555) && ($edi == 2554))
    printf "APPLE_PATCHED_DLL_CONTINUITY_BRANCH current=%u previous=%u\n", $edx, $edi
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  set $selected_count = $selected_count + 1
  if $selected_count <= 15
    printf "APPLE_PATCHED_DLL_SELECTED n=%u id=%u\n", $selected_count, *(unsigned int *)($esp + 8)
  end
  continue
end

continue
