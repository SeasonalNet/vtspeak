set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
break *0x100086c0
commands
  silent
  set $input = *(unsigned int *)($esp + 4)
  set $rows = *(unsigned int *)($esp + 8)
  set $index = *(unsigned int *)($esp + 12)
  set $optional = *(unsigned int *)($esp + 16)
  set $return_address = *(unsigned int *)$esp
  set $entry_row = $rows + $index * 0x70 + 4
  printf "HELPER_ENTRY index=%u optional=%#x metadata4=%u text=", $index, $optional, *(unsigned short *)($entry_row + 0x6e)
  x/s $input
  tbreak *$return_address
  commands
    silent
    printf "HELPER_RETURN index=%u eax=%#x low=%u row_flags=%#x metadata4=%u\n", $index, $eax, $eax & 0xffff, *(unsigned short *)$entry_row, *(unsigned short *)($entry_row + 0x6e)
    continue
  end
  continue
end
continue
