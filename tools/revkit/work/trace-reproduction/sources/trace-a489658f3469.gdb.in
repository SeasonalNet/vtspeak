set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10034180
commands
  silent
  set $state = *(unsigned int *)($esp + 4)
  set $second = *(unsigned int *)($esp + 8)
  set $records = *(unsigned int *)($state + 0x39e4)
  set $count = *(unsigned int *)$records
  set $return_address = *(unsigned int *)$esp
  printf "PROPER_NAME_ENTRY count=%u second_arg=%u\n", $count, $second
  set $i = 0
  while $i < $count && $i < 8
    printf "PROPER_NAME_PART index=%u text=", $i
    x/s ($records + $i * 0x140 + 0x46)
    set $i = $i + 1
  end
  tbreak *$return_address
  commands
    silent
    printf "PROPER_NAME_RETURN eax=%#x\n", $eax
    continue
  end
  continue
end

break *0x1000e0c0
commands
  silent
  set $tokens = *(unsigned int *)($esp + 4)
  set $numeric_code = *(unsigned char *)($esp + 8)
  set $first_index = *(unsigned int *)($esp + 12)
  set $last_index = *(unsigned int *)($esp + 16)
  set $return_address = *(unsigned int *)$esp
  printf "TPP_NUMERIC_APPLY code=%u first=%u last=%u\n", $numeric_code, $first_index, $last_index
  tbreak *$return_address
  commands
    silent
    set $i = $first_index
    if $last_index < $first_index
      set $last_index = $first_index
    end
    while $i <= $last_index && $i < 64
      set $token = $tokens + $i * 0x94
      printf "TPP_NUMERIC_TOKEN index=%u code_byte_25=%u class_byte_26=%u surface=", $i, *(unsigned char *)($token + 0x25), *(unsigned char *)($token + 0x26)
      x/s ($token + 0x34)
      set $i = $i + 1
    end
    continue
  end
  continue
end

break *0x1003a7b0
commands
  silent
  set $out = *(unsigned int *)($esp + 4)
  set $key = *(unsigned int *)($esp + 8)
  set $selector = *(unsigned int *)($esp + 12)
  set $ret = *(unsigned int *)$esp
  printf "TPP_LOOKUP selector=%#x key=", $selector
  x/s $key
  tbreak *$ret
  commands
    silent
    printf "TPP_RESULT eax=%#x output=", $eax
    x/s $out
    continue
  end
  continue
end

break *0x10007520
commands
  silent
  set $ctx = *(unsigned int *)($esp + 4)
  set $upstream = *(unsigned int *)($esp + 8)
  set $rows = $ctx + 0x429a2
  set $ret = *(unsigned int *)$esp
  printf "CONTEXT_ENTRY count=%u\n", *(unsigned short *)$rows
  set $j = 0
  while $j < *(unsigned short *)$rows && $j < 128
    set $entry = $rows + 4 + $j * 0x70
    set $source = *(unsigned short *)($entry + 2)
    set $token = $upstream + $source * 0x94
    printf "UPSTREAM row=%u source=%u type=%#x class=%#x code_byte_25=%u class_byte_26=%u status=%#x surface=", $j, $source, *(unsigned char *)($token + 0x23), *(unsigned char *)($token + 0x24), *(unsigned char *)($token + 0x25), *(unsigned char *)($token + 0x26), *(unsigned char *)($token + 0x30)
    x/s ($token + 0x34)
    set $j = $j + 1
  end
  tbreak *$ret
  commands
    silent
    printf "CONTEXT_RETURN count=%u\n", *(unsigned short *)$rows
    set $i = 0
    while $i < *(unsigned short *)$rows && $i < 128
      set $row = $rows + 4 + $i * 0x70
      printf "ROW index=%u flags=%#x source=%u class=%#x surface=", $i, *(unsigned short *)$row, *(unsigned short *)($row + 2), *(unsigned char *)($row + 6)
      x/s ($row + 7)
      printf "ROW_CODE index=%u bytes:", $i
      x/16bx ($row + 0x25)
      set $i = $i + 1
    end
    continue
  end
  continue
end

continue
