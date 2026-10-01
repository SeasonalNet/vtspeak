set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

catch signal SIGSEGV
commands
  silent
  if $eip >= 0x10000000 && $eip < 0x10110000
    printf "UNIT_HISTORY_FAULT eip=%p\n", $eip
    x/i $eip
    bt 6
  end
  continue
end

break *0x10027af0
commands
  silent
  set $load_return = *(unsigned int *)$esp
  printf "UNIT_HISTORY_PRESENT_DBROOT arg=%#x\n", *(unsigned int *)($esp + 12)
  call ((void (*)(unsigned char))0x100283f0)(1)
  printf "UNIT_HISTORY_PRESENT_MODE value=%u\n", *(unsigned char *)0x100a045d
  disable 1
  tbreak *$load_return
  commands
    silent
    printf "UNIT_HISTORY_PRESENT_LOAD ax=%d\n", (short)$eax
    continue
  end
  continue
end

break *0x1001b240
commands
  silent
  set $history_item = *(unsigned int *)($esp + 4)
  set $history_path = (char *)*(unsigned int *)($esp + 8)
  set $history_return = *(unsigned int *)$esp
  printf "UNIT_HISTORY_PRESENT_ENTRY count=%d path=%s\n", *(int *)($history_item + 0x4c), $history_path
  tbreak *$history_return
  commands
    silent
    printf "UNIT_HISTORY_PRESENT_RETURN ax=%d\n", (short)$eax
    continue
  end
  continue
end

break *0x1001b25c
commands
  silent
  set $history_item = *(unsigned int *)($ebp + 8)
  printf "UNIT_HISTORY_PRESENT_OPEN source=%p\n", $edi
  continue
end

break *0x1001b28b
commands
  silent
  set $history_item = *(unsigned int *)($ebp + 8)
  set $history_count = *(int *)($history_item + 0x4c)
  set $history_a = (unsigned int *)*(unsigned int *)($history_item + 0x14)
  set $history_b = (unsigned int *)*(unsigned int *)($history_item + 0x18)
  printf "UNIT_HISTORY_PRESENT_DATA count=%d a_first=%u a_last=%u b_first=%u b_last=%u\n", $history_count, *$history_a, *($history_a + $history_count - 1), *$history_b, *($history_b + $history_count - 1)
  continue
end

break *0x1001b343
commands
  silent
  printf "UNIT_HISTORY_PRESENT_PARSE_FAILURE\n"
  continue
end

break *0x1001264f
commands
  silent
  printf "UNIT_HISTORY_PRESENT_BANK_LOOP_FAILURE result_ax=%d\n", (short)$eax
  continue
end

break *0x10012659
commands
  silent
  printf "UNIT_HISTORY_PRESENT_BANK_LOOP_RETURN ax=%d\n", (short)$eax
  continue
end

break *0x1001da50
commands
  silent
  set $text_return = *(unsigned int *)$esp
  disable 9
  tbreak *$text_return
  commands
    silent
    printf "UNIT_HISTORY_PRESENT_TEXT ax=%d\n", (short)$eax
    continue
  end
  continue
end

continue
