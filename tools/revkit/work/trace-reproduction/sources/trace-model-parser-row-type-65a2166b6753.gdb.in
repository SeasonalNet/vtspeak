set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $parser_row_calls = 0

# FUN_1003e240 receives the parser state after the source parser has written
# each row. This trace only reads the row type consumed later by FUN_1000ea20.
break *0x1003e240
commands
  silent
  set $parser_row_calls = $parser_row_calls + 1
  set $state = *(unsigned int *)($esp + 4)
  set $count = *(short *)$state
  printf "MODEL_PARSER_ROWS call=%u count=%d\n", $parser_row_calls, $count
  set $row_index = 0
  while $row_index < $count && $row_index < 100
    set $row = $state + 0x14 + $row_index * 0x94
    printf "MODEL_PARSER_ROW call=%u index=%d start=%d end=%d discriminator=%02x context_type=%08x\n", $parser_row_calls, $row_index, *(int *)$row, *(int *)($row + 4), *(unsigned char *)($row + 0x24), *(unsigned int *)($row + 0x2c)
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
