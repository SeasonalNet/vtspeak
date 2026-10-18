set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10024060
commands
  silent
  set $slot = *(unsigned int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $row = $state + ($slot * 3 + 0x76314) * 2
  set $bank = *(unsigned char *)$row
  set $leaf = *(unsigned char *)($row + 2)
  set $voice = *(unsigned int *)($state + 0x4c)
  set $signature = $voice + $bank * 0x3c0 + 0x6e2 + $leaf * 7
  set $return = *(unsigned int *)$esp
  printf "APPLE_SELECTOR_SIGNATURE slot=%u bank=%u leaf=%u:", $slot, $bank, $leaf
  x/7ub $signature
  tbreak *$return
  commands
    silent
    printf "APPLE_SELECTOR_SPLIT_RETURN slot=%u accepted=%u\n", $slot, $al
    continue
  end
  continue
end

break *0x10023060
commands
  silent
  set $caller = *(unsigned int *)$esp
  if $caller >= 0x10024060 && $caller < 0x100242a0
    set $ids = *(unsigned int *)($esp + 4)
    set $count = *(unsigned short *)($esp + 8)
    set $context = *(unsigned int *)($esp + 12)
    set $weights = *(unsigned int *)($context + 0x8c)
    set $return = $caller
    printf "APPLE_SELECTOR_CLASS_ENTRY count=%u ids_weights:", $count
    set $i = 0
    while $i < $count && $i < 64
      set $id = *(unsigned int *)($ids + $i * 4)
      printf " %u:%u", $id, *(unsigned short *)($weights + $id * 2)
      set $i = $i + 1
    end
    printf "\n"
    tbreak *$return
    commands
      silent
      printf "APPLE_SELECTOR_CLASS_SUM=%u\n", $eax
      continue
    end
  end
  continue
end

break *0x100242a0
commands
  silent
  set $slot = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  printf "APPLE_SELECTOR_FALLBACK slot=%u\n", $slot
  tbreak *$return
  commands
    silent
    printf "APPLE_SELECTOR_FALLBACK_ROWS=%u\n", $eax
    continue
  end
  continue
end

break *0x10024680
commands
  silent
  set $return = *(unsigned int *)$esp
  tbreak *$return
  commands
    silent
    printf "APPLE_SELECTOR_POSITION_COUNT=%u\n", $eax
    continue
  end
  continue
end

break *0x408187
commands
  silent
  printf "APPLE_SELECTOR_TRACE_READY\n"
  continue
end

continue
