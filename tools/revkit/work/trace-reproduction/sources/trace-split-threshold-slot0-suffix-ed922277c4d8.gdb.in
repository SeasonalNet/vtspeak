set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/split-threshold-slot0-suffix-gdb.log
set logging overwrite on
set logging enabled on

break *0x10024060
commands
  silent
  set $slot = *(unsigned int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $return = *(unsigned int *)$esp
  if $slot == 0
    set $context_row = $state + ($slot * 3 + 0x76314) * 2
    set $bank = *(unsigned char *)$context_row
    set $voice = *(unsigned int *)($state + 0x4c)
    set $sig = $voice + $bank * 0x3c0 + 0x6e2 + *(unsigned char *)($context_row + 2) * 7
    set $original_byte5 = *(unsigned char *)($sig + 5)
    printf "PATCH_SLOT0_SIGNATURE_BEFORE:"
    x/7ub $sig
    set {unsigned char}($sig + 5) = 30
    printf "PATCH_SLOT0_SIGNATURE_DURING:"
    x/7ub $sig
  end
  printf "SLOT0_TRACE_SPLIT_ENTER slot=%u\n", $slot
  tbreak *$return
  commands
    silent
    if $slot == 0
      set {unsigned char}($sig + 5) = $original_byte5
      printf "PATCH_SLOT0_SIGNATURE_RESTORED:"
      x/7ub $sig
    end
    printf "SLOT0_TRACE_SPLIT_RETURN slot=%u accepted=%u\n", $slot, $al
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
    printf "SLOT0_TRACE_METRIC_ENTRY caller=%#x count=%u ids_weights:", $caller, $count
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
      printf "SLOT0_TRACE_METRIC_RETURN sum=%u\n", $eax
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
  printf "SLOT0_TRACE_FALLBACK_ENTER slot=%u\n", $slot
  tbreak *$return
  commands
    silent
    printf "SLOT0_TRACE_FALLBACK_RETURN slot=%u positions=%u\n", $slot, $eax
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
    printf "SLOT0_TRACE_POSITION_BUILDER_RETURN positions=%u\n", $eax
    continue
  end
  continue
end

break *0x408187
commands
  silent
  printf "SLOT0_SUFFIX_TRACE_READY\n"
  continue
end

continue
