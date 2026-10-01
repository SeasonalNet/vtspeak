set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

# Capture each literal coefficient triple consumed by the 2013 edge scorer
# while synthesizing plain Hello. The coefficient table layout was mapped in
# the Stage 20 Apple selector investigation.
hbreak *0x10018fbb
commands
  silent
  set $context = *(int *)($ebp + 8)
  set $current = *(unsigned int *)($ebp - 8)
  set $previous_ptr = *(unsigned int *)($ebp - 0x18)
  set $previous = *(unsigned int *)$previous_ptr
  set $slot = *(signed char *)($ebp - 1)
  set $table_offset = $slot * 12
  printf "HELLO_2013_WEIGHT_TABLE context=%u current=%u previous=%u slot=%d mode=%u weights=(%g,%g,%g) distances=(%g,%g,%g)\n", $context, $current, $previous, $slot, *(unsigned char *)($edx + 4), *(float *)(0x1007c288 + $table_offset), *(float *)(0x1007c290 + $table_offset), *(float *)(0x1007c28c + $table_offset), *(float *)($ebp - 0x38), *(float *)($ebp - 0x34), *(float *)($ebp - 0x30)
  continue
end

hbreak *0x1001b200
commands
  silent
  printf "HELLO_SELECTED_UNIT_HARDWARE id=%u\n", *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "HELLO_WEIGHT_TRACE_READY\n"
  continue
end

continue
