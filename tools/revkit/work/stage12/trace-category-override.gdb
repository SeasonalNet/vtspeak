set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1003a7b0
commands
  silent
  set $out = *(unsigned int *)($esp + 4)
  set $key = *(unsigned int *)($esp + 8)
  set $selector = *(unsigned int *)($esp + 12)
  set $length = 0
  while *(unsigned char *)($key + $length) != 0 && $length < 64
    set $length = $length + 1
  end
  set $prefix = *(unsigned int *)$key
  set $offset = -1
  set $replacement = 0
  if $selector == 0x41 && $length == 6 && $prefix == 0x6f636341
    set $offset = 0
    set $replacement = 0x31
  end
  if $selector == 0x41 && $length == 6 && $prefix == 0x746e6543
    set $offset = 0
    set $replacement = 0x30
  end
  if $selector == 0x41 && $length == 7 && $prefix == 0x6c65704f
    set $offset = 0
    set $replacement = 0x30
  end
  if $selector == 0x41 && $length == 9 && $prefix == 0x6c6c6154
    set $offset = 0
    set $replacement = 0x30
  end
  if $selector == 0x42 && $length == 13 && $prefix == 0x6f626241
    set $offset = 2
    set $replacement = 0x30
  end
  if $selector == 0x43 && $length == 18 && $prefix == 0x6c656d41
    set $offset = 2
    set $replacement = 0x30
  end
  if $selector == 0x44 && $length == 20 && $prefix == 0x65706143
    set $offset = 1
    set $replacement = 0x30
  end
  if $offset >= 0
    printf "TPP_CODE_OVERRIDE selector=%#x key=", $selector
    x/s $key
    set $return = *(unsigned int *)$esp
    tbreak *$return
    commands
      silent
      set {unsigned char}($out + $offset) = $replacement
      printf "TPP_CODE_CHANGED output="
      x/s $out
      continue
    end
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
  set $j = 0
  while $j < *(unsigned short *)$rows && $j < 128
    set $entry = $rows + 4 + $j * 0x70
    set $source = *(unsigned short *)($entry + 2)
    set $token = $upstream + $source * 0x94
    printf "EFFECT_UPSTREAM row=%u source=%u code25=%u class26=%u surface=", $j, $source, *(unsigned char *)($token + 0x25), *(unsigned char *)($token + 0x26)
    x/s ($token + 0x34)
    set $j = $j + 1
  end
  tbreak *$ret
  commands
    silent
    set $i = 0
    while $i < *(unsigned short *)$rows && $i < 128
      set $row = $rows + 4 + $i * 0x70
      printf "EFFECT_ROW index=%u source=%u flags=%#x surface=", $i, *(unsigned short *)($row + 2), *(unsigned short *)$row
      x/s ($row + 7)
      printf "EFFECT_ROW_CODE index=%u bytes:", $i
      x/16bx ($row + 0x25)
      set $i = $i + 1
    end
    continue
  end
  continue
end

break *0x1001b200
commands
  silent
  printf "EFFECT_SELECTED_UNIT id=%u\n", *(unsigned int *)($esp + 8)
  continue
end

continue
