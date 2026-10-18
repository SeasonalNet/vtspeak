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
    printf "UPSTREAM row=%u source=%u status=%#x surface=", $j, $source_index, *(unsigned char *)($token + 0x30)
    x/s ($token + 0x34)
    printf "PHONE_CONTEXT_BEFORE row=%u metadata=%u,%u,%u,%u surface=", $j, *(unsigned short *)($entry_row + 0x68), *(unsigned short *)($entry_row + 0x6a), *(unsigned short *)($entry_row + 0x6c), *(unsigned short *)($entry_row + 0x6e)
    x/s ($entry_row + 7)
    printf "PHONE_CONTEXT_BEFORE_CODE row=%u ", $j
    x/s ($entry_row + 0x25)
    set {unsigned short}($entry_row + 0x68) = 0
    set {unsigned short}($entry_row + 0x6a) = 0
    set {unsigned short}($entry_row + 0x6c) = 0
    set {unsigned short}($entry_row + 0x6e) = 0
    set $j = $j + 1
  end
  tbreak *$return_address
  commands
    silent
    printf "CONTEXT_RETURN count=%u\n", *(unsigned short *)$rows
    set $j = 0
    while $j < *(unsigned short *)$rows && $j < 128
      set $entry_row = $rows + 4 + $j * 0x70
      printf "PHONE_CONTEXT_AFTER row=%u flags=%#x metadata=%u,%u,%u,%u surface=", $j, *(unsigned short *)$entry_row, *(unsigned short *)($entry_row + 0x68), *(unsigned short *)($entry_row + 0x6a), *(unsigned short *)($entry_row + 0x6c), *(unsigned short *)($entry_row + 0x6e)
      x/s ($entry_row + 7)
      printf "PHONE_CONTEXT_AFTER_CODE row=%u ", $j
      x/s ($entry_row + 0x25)
      set $j = $j + 1
    end
    continue
  end
  continue
end
continue
