set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10027af0
commands
  silent
  disable 1
  printf "UNIT_HISTORY_LOAD_REQUEST prior=%u loaded=%u requested=1\n", *(unsigned char *)0x100a045d, *(unsigned char *)0x100a747c
  call ((void (*)(unsigned char))0x100283f0)(1)
  printf "UNIT_HISTORY_LOAD_SET value=%u loaded=%u\n", *(unsigned char *)0x100a045d, *(unsigned char *)0x100a747c
  continue
end

break *0x1001b240
commands
  silent
  set $history_arg0 = *(unsigned int *)($esp + 4)
  set $history_path = (char *)*(unsigned int *)($esp + 8)
  set $history_return = *(unsigned int *)$esp
  printf "UNIT_HISTORY_ENTRY arg0=%p path=%s\n", $history_arg0, $history_path
  tbreak *$history_return
  commands
    silent
    printf "UNIT_HISTORY_RETURN eax=0x%08x ax=%d\n", $eax, (short)$eax
    continue
  end
  continue
end

break *0x1001b25c
commands
  silent
  set $history_item = *(unsigned int *)($ebp + 8)
  printf "UNIT_HISTORY_SOURCE handle=%p count=%d\n", $edi, *(int *)($history_item + 0x4c)
  continue
end

break *0x1001b28b
commands
  silent
  set $history_item = *(unsigned int *)($ebp + 8)
  set $history_count = *(int *)($history_item + 0x4c)
  set $history_a = (unsigned int *)*(unsigned int *)($history_item + 0x14)
  set $history_b = (unsigned int *)*(unsigned int *)($history_item + 0x18)
  printf "UNIT_HISTORY_FALLBACK count=%d a_first=%u a_last=%u b_first=%u b_last=%u\n", $history_count, *$history_a, *($history_a + $history_count - 1), *$history_b, *($history_b + $history_count - 1)
  continue
end

break *0x1001da50
commands
  silent
  set $history_text_return = *(unsigned int *)$esp
  printf "UNIT_HISTORY_TEXT_ENTRY format=%d text=%#x path=%#x\n", *(int *)($esp + 4), *(unsigned int *)($esp + 8), *(unsigned int *)($esp + 12)
  disable 3
  tbreak *$history_text_return
  commands
    silent
    printf "UNIT_HISTORY_TEXT_RETURN ax=%d\n", (short)$eax
    continue
  end
  continue
end

continue
