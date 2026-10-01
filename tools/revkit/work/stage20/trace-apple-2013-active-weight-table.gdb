set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/apple-selector-byte0-active-weight-table-2026-09-29/adapted-active-weight-table-gdb.log
set logging overwrite on
set logging enabled on

# Record the selector index and the literal coefficient triple consumed by
# the 2013 scorer for the Apple edges compared with the 2006 path.
hbreak *0x10018fbb
commands
  silent
  set $context = *(int *)($ebp + 8)
  set $current = *(unsigned int *)($ebp - 8)
  set $previous_ptr = *(unsigned int *)($ebp - 0x18)
  set $previous = *(unsigned int *)$previous_ptr
  set $target = 0
  if ($context == 2) && ($current == 100180) && ($previous == 177774)
    set $target = 1
  end
  if ($context == 3) && ($current == 59558) && ($previous == 100180)
    set $target = 1
  end
  if ($context == 3) && ($current == 2554) && ($previous == 64719)
    set $target = 1
  end
  if ($context == 3) && ($current == 82036) && ($previous == 64719)
    set $target = 1
  end
  if ($context == 4) && ($current == 59559) && ($previous == 59558)
    set $target = 1
  end
  if ($context == 4) && ($current == 82037) && ($previous == 82036)
    set $target = 1
  end
  if ($context == 4) && ($current == 2555) && ($previous == 2554)
    set $target = 1
  end
  if $target
    set $slot = *(signed char *)($ebp - 1)
    set $table_offset = $slot * 12
    printf "APPLE_2013_WEIGHT_TABLE context=%u current=%u previous=%u slot=%d mode=%u weights=(%g,%g,%g) distances=(%g,%g,%g)\n", $context, $current, $previous, $slot, *(unsigned char *)($edx + 4), *(float *)(0x1007c288 + $table_offset), *(float *)(0x1007c290 + $table_offset), *(float *)(0x1007c28c + $table_offset), *(float *)($ebp - 0x38), *(float *)($ebp - 0x34), *(float *)($ebp - 0x30)
  end
  continue
end

continue
