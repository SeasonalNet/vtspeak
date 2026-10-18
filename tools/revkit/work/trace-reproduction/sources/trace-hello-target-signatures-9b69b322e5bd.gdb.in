set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10024060
commands
  silent
  set $slot = *(unsigned int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $return = *(unsigned int *)$esp
  set $context_row = $state + ($slot * 3 + 0x76314) * 2
  set $bank = *(unsigned char *)$context_row
  set $voice = *(unsigned int *)($state + 0x4c)
  set $sig = $voice + $bank * 0x3c0 + 0x6e2 + *(unsigned char *)($context_row + 2) * 7
  printf "HELLO_TARGET_SIGNATURE slot=%u bank=%u phone_row=%u:", $slot, $bank, *(unsigned char *)($context_row + 2)
  x/7ub $sig
  tbreak *$return
  commands
    silent
    printf "HELLO_TARGET_SPLIT_RETURN slot=%u accepted=%u\n", $slot, $al
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
    printf "HELLO_TARGET_METRIC_ENTRY count=%u ids_weights:", $count
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
      printf "HELLO_TARGET_METRIC_RETURN sum=%u\n", $eax
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
  printf "HELLO_TARGET_FALLBACK_ENTER slot=%u\n", $slot
  tbreak *$return
  commands
    silent
    printf "HELLO_TARGET_FALLBACK_RETURN positions=%u\n", $eax
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
    printf "HELLO_TARGET_POSITION_COUNT=%u\n", $eax
    continue
  end
  continue
end

break *0x408187
commands
  silent
  printf "HELLO_TARGET_SIGNATURE_TRACE_READY\n"
  continue
end

continue
