set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $parser_row_calls = 0

# FUN_1003e240 receives parser state after FUN_1003d3d0/FUN_1003e070 and before
# it mutates the state. Record source offsets, selected row bytes, and the two
# source strings consumed by FUN_1000d190.
break *0x1003e240
commands
  silent
  set $parser_row_calls = $parser_row_calls + 1
  set $state = *(unsigned int *)($esp + 4)
  set $count = *(short *)$state
  printf "MODEL_PARSER_ROWS call=%u state=%08x count=%d\n", $parser_row_calls, $state, $count
  set $row_index = 0
  while $row_index < $count && $row_index < 100
    set $row = $state + 0x14 + $row_index * 0x94
    printf "MODEL_PARSER_ROW call=%u index=%d start=%d end=%d field20=%08x field23=%02x field29=%02x field30=%02x field52=%02x field66=%02x\n", $parser_row_calls, $row_index, *(int *)$row, *(int *)($row + 4), *(unsigned int *)($row + 0x20), *(unsigned char *)($row + 0x23), *(unsigned char *)($row + 0x29), *(unsigned char *)($row + 0x30), *(unsigned char *)($row + 0x52), *(unsigned char *)($row + 0x66)
    printf "MODEL_PARSER_ROW_PRIMARY_TEXT call=%u index=%d\n", $parser_row_calls, $row_index
    x/s $row + 0x34
    printf "MODEL_PARSER_ROW_AUX_TEXT call=%u index=%d\n", $parser_row_calls, $row_index
    x/s $row + 0x52
    set $row_index = $row_index + 1
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "MODEL_PARSER_TRACE_READY\n"
  continue
end

continue
