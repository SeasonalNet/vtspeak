set pagination off
set confirm off
break *0x10007520
commands
  silent
  set $ctx = *(unsigned int *)($esp + 4)
  set $upstream = *(unsigned int *)($esp + 8)
  set $rows = $ctx + 0x429a2
  set $return_address = *(unsigned int *)$esp
  printf "CONTEXT_ENTRY count=%u rows=%#x upstream=%#x\n", *(unsigned short *)$rows, $rows, $upstream
  set $j = 0
  while $j < *(unsigned short *)$rows && $j < 128
    set $entry_row = $rows + 4 + $j * 0x70
    set $source_index = *(unsigned short *)($entry_row + 2)
    set $token = $upstream + $source_index * 0x94
    printf "UPSTREAM row=%u source=%u type=%#x class=%#x status=%#x surface=", $j, $source_index, *(unsigned char *)($token + 0x23), *(unsigned char *)($token + 0x24), *(unsigned char *)($token + 0x30)
    x/s ($token + 0x34)
    printf "UPSTREAM_CONTEXT row=%u value=", $j
    x/s ($token + 0x52)
    printf "UPSTREAM_CONTEXT_BYTES row=%u bytes:\n", $j
    x/12bx ($token + 0x52)
    set $j = $j + 1
  end
  tbreak *$return_address
  commands
    silent
    printf "CONTEXT_RETURN count=%u\n", *(unsigned short *)$rows
    set $i = 0
    while $i < *(unsigned short *)$rows && $i < 128
      set $row = $rows + 4 + $i * 0x70
      printf "ROW index=%u flags=%#x source=%u class=%#x marker=%#x surface=", $i, *(unsigned short *)$row, *(unsigned short *)($row + 2), *(unsigned char *)($row + 6), *(unsigned char *)($row + 0x66)
      x/s ($row + 7)
      printf "ROW_CODE index=%u bytes:", $i
      x/24bx ($row + 0x25)
      printf "ROW_META index=%u words:", $i
      x/4uh ($row + 0x68)
      set $i = $i + 1
    end
    continue
  end
  continue
end
continue
