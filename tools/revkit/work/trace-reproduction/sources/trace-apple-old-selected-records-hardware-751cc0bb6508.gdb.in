set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $record_calls = 0

hbreak *0x10013cf0
commands
  silent
  set $record_calls = $record_calls + 1
  if $record_calls <= 64
    printf "APPLE_OLD_SELECTED_RECORD_HARDWARE n=%u unit_id=%u\n", $record_calls, *(unsigned int *)($esp + 8)
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_OLD_HARDWARE_TRACE_READY\n"
  continue
end

continue
