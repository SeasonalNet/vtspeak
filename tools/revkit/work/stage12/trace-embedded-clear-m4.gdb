set pagination off
set confirm off
break *0x10007520
commands
  silent
  set $ctx = *(unsigned int *)($esp + 4)
  set $rows = $ctx + 0x429a2
  set $return_address = *(unsigned int *)$esp
  set $entry_row = $rows + 4
  printf "METADATA_BEFORE=%u,%u,%u,%u\n", *(unsigned short *)($entry_row + 0x68), *(unsigned short *)($entry_row + 0x6a), *(unsigned short *)($entry_row + 0x6c), *(unsigned short *)($entry_row + 0x6e)
  set {unsigned short}($entry_row + 0x6e) = 0
  tbreak *$return_address
  commands
    silent
    printf "METADATA_AFTER flags=%#x tuple=%u,%u,%u,%u surface=", *(unsigned short *)$entry_row, *(unsigned short *)($entry_row + 0x68), *(unsigned short *)($entry_row + 0x6a), *(unsigned short *)($entry_row + 0x6c), *(unsigned short *)($entry_row + 0x6e)
    x/s ($entry_row + 7)
    printf "PHONE_CODE "
    x/s ($entry_row + 0x25)
    continue
  end
  continue
end
continue
