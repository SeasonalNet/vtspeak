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
  set $ret = *(unsigned int *)$esp
  if $selector == 0x41
    tbreak *$ret
    commands
      silent
      if $eax != 0 && *(unsigned char *)$out == 0x58 && *(unsigned char *)($out + 1) == 0
        printf "TPP_AX_OVERRIDE original=X replacement=miss key="
        x/s $key
        set {unsigned char}($out) = 0
        set $eax = 0
        printf "TPP_AX_CHANGED output="
        x/s $out
      end
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
