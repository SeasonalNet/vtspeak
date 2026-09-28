set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/split-threshold-new-gdb.log
set logging overwrite on
set logging enabled on

break *0x10024060
commands
  silent
  set $split_slot = *(unsigned int *)($esp + 4)
  set $split_state = *(unsigned int *)($esp + 8)
  set $split_return = *(unsigned int *)$esp
  set $split_row = $split_state + ($split_slot * 3 + 0x76314) * 2
  set $split_bank = *(unsigned char *)$split_row
  set $split_voice = *(unsigned int *)($split_state + 0x4c)
  set $split_class = *(unsigned char *)($split_voice + $split_bank * 0x3c0 + 0x92b)
  printf "NEW_SPLIT_ENTER slot=%u bank=%u class_byte=%u\n", $split_slot, $split_bank, $split_class
  tbreak *$split_return
  commands
    silent
    printf "NEW_SPLIT_RETURN slot=%u accepted=%u\n", $split_slot, $al
    continue
  end
  continue
end

break *0x10023060
commands
  silent
  set $sum_caller = *(unsigned int *)$esp
  if $sum_caller >= 0x10024060 && $sum_caller < 0x100242a0
    set $sum_ids = *(unsigned int *)($esp + 4)
    set $sum_count = *(unsigned short *)($esp + 8)
    set $sum_context = *(unsigned int *)($esp + 12)
    set $sum_weights = *(unsigned int *)($sum_context + 0x8c)
    set $sum_return = $sum_caller
    printf "NEW_SPLIT_METRIC_ENTRY caller=%#x count=%u ids_weights:", $sum_caller, $sum_count
    set $sum_i = 0
    while $sum_i < $sum_count && $sum_i < 64
      set $sum_id = *(unsigned int *)($sum_ids + $sum_i * 4)
      printf " %u:%u", $sum_id, *(unsigned short *)($sum_weights + $sum_id * 2)
      set $sum_i = $sum_i + 1
    end
    printf "\n"
    tbreak *$sum_return
    commands
      silent
      printf "NEW_SPLIT_METRIC_RETURN sum=%u\n", $eax
      continue
    end
  end
  continue
end

break *0x10018770
commands
  silent
  set $lookup_caller = *(unsigned int *)$esp
  if $lookup_caller >= 0x10024060 && $lookup_caller < 0x10024700
    set $lookup_key = *(unsigned int *)($esp + 4)
    set $lookup_slot = *(unsigned int *)($esp + 8)
    set $lookup_state = *(unsigned int *)($esp + 12)
    set $lookup_mode = *(unsigned short *)($esp + 20)
    set $lookup_return = $lookup_caller
    printf "NEW_CLASS_LOOKUP_ENTER slot=%u mode=%u key:", $lookup_slot, $lookup_mode
    x/7ub $lookup_key
    tbreak *$lookup_return
    commands
      silent
      set $lookup_count = *(unsigned short *)($lookup_state + 0xec620)
      printf "NEW_CLASS_LOOKUP_RETURN slot=%u mode=%u return=%u total=%u ids:", $lookup_slot, $lookup_mode, $eax, $lookup_count
      set $lookup_i = 0
      while $lookup_i < $lookup_count && $lookup_i < 32
        printf " %u", *(unsigned int *)($lookup_state + 0xec170 + $lookup_i * 4)
        set $lookup_i = $lookup_i + 1
      end
      printf "\n"
      continue
    end
  end
  continue
end

break *0x100242a0
commands
  silent
  set $fallback_slot = *(unsigned int *)($esp + 4)
  set $fallback_return = *(unsigned int *)$esp
  printf "NEW_SPLIT_FALLBACK_ENTER slot=%u\n", $fallback_slot
  tbreak *$fallback_return
  commands
    silent
    printf "NEW_SPLIT_FALLBACK_RETURN slot=%u positions=%u\n", $fallback_slot, $eax
    continue
  end
  continue
end

break *0x10024680
commands
  silent
  set $selector_return = *(unsigned int *)$esp
  tbreak *$selector_return
  commands
    silent
    printf "NEW_POSITION_BUILDER_RETURN positions=%u\n", $eax
    continue
  end
  continue
end

break *0x408187
commands
  silent
  printf "NEW_THRESHOLD_TRACE_READY\n"
  continue
end

continue
