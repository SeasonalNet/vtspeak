set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/apple-selector-byte0-old-path-costs-2026-09-29/legacy-path-costs-gdb.log
set logging overwrite on
set logging enabled on

# FUN_1001e470's transition update: retain only edges touching the competing
# IDs at Apple contexts 2-4, where the 2006 and 2013 selected chains diverge.
hbreak *0x1001ec07
commands
  silent
  set $context = *(int *)($esp + 0x2c)
  set $current = *(unsigned int *)($esp + 0x3c)
  set $previous = $ebp
  if (($context == 2) && (($current == 100180) || ($current == 64719))) || (($context == 3) && (($current == 59558) || ($current == 2554) || ($current == 82036)) && (($previous == 100180) || ($previous == 64719) || ($previous == 189455))) || (($context == 4) && (($current == 59559) || ($current == 82037) || ($current == 2555)) && (($previous == 59558) || ($previous == 2554) || ($previous == 82036)))
    printf "APPLE_OLD_PATH_EDGE context=%d current=%u previous=%u transition=%g previous_path=%g local=%g total=%g penalty=%g\n", $context, $current, $previous, *(float *)($esp + 0x44), *(float *)(*(unsigned int *)($esp + 0x6c) - 4), *(float *)($esp + 0x94), $st0, *(float *)($esp + 0x14)
  end
  continue
end

# Capture the predecessor index selected for each target row after its
# previous-position scan has completed.
hbreak *0x1001ec44
commands
  silent
  set $position = *(int *)($esp + 0x2c)
  set $candidate = *(unsigned int *)($esp + 0x3c)
  if (($position == 2) && (($candidate == 100180) || ($candidate == 64719))) || (($position == 3) && (($candidate == 59558) || ($candidate == 2554) || ($candidate == 82036))) || (($position == 4) && (($candidate == 59559) || ($candidate == 82037) || ($candidate == 2555)))
    set $state = *(unsigned int *)($esp + 0xa0)
    set $pred_index = $eax
    set $pred_array = $state + 0xbb154 + ($position - 1) * 0xfc + 0x7c
    set $pred_id = *(unsigned int *)($pred_array + $pred_index * 4)
    printf "APPLE_OLD_PATH_BACKPOINTER position=%d candidate=%u predecessor_index=%u predecessor=%u\n", $position, $candidate, $pred_index, $pred_id
  end
  continue
end

# Tie the DP result to the IDs handed to the legacy selected-record reader.
hbreak *0x1001ed7f
commands
  silent
  printf "APPLE_OLD_FINAL_CHAIN selected_index=%d unit_id=%u\n", *(short *)$edi, $edx
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_OLD_PATH_TRACE_READY\n"
  continue
end

continue
