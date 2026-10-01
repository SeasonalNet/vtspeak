set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $selected_calls = 0

hbreak *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  printf "HELLO_SELECTED_UNIT_HARDWARE n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "HELLO_HARDWARE_TRACE_READY\n"
  continue
end

continue
