set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/old-class-features-gdb.log
set logging overwrite on
set logging enabled on

break *0x408187
commands
  silent
  printf "OLD_THRESHOLD_TRACE_READY\n"

break *0x1001dfb0
commands
  silent
  set $split_slot = *(unsigned int *)($esp + 4)
  set $split_state = *(unsigned int *)($esp + 8)
  set $split_context = *(unsigned int *)($esp + 12)
  set $split_return = *(unsigned int *)$esp
  printf "OLD_SPLIT_ENTER slot=%u state=%#x context=%#x\n", $split_slot, $split_state, $split_context
  tbreak *$split_return
  commands
    silent
    printf "OLD_SPLIT_RETURN slot=%u accepted=%u\n", $split_slot, $al
    continue
  end
  continue
end

break *0x1001d5c0
commands
  silent
  set $sum_caller = *(unsigned int *)$esp
  if $sum_caller >= 0x1001dfb0 && $sum_caller < 0x1001e110
    set $sum_ids = *(unsigned int *)($esp + 4)
    set $sum_count = *(unsigned short *)($esp + 8)
    set $sum_context = *(unsigned int *)($esp + 12)
    set $sum_weights = *(unsigned int *)($sum_context + 0x8c)
    set $sum_return = $sum_caller
    printf "OLD_SPLIT_METRIC_ENTRY caller=%#x count=%u ids_weights:", $sum_caller, $sum_count
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
      printf "OLD_SPLIT_METRIC_RETURN sum=%u\n", $eax
      continue
    end
  end
  continue
end

break *0x1001da80
commands
  silent
  set $lookup_caller = *(unsigned int *)$esp
  if $lookup_caller >= 0x1001dfb0 && $lookup_caller < 0x1001e110
    set $lookup_key = *(unsigned int *)($esp + 4)
    set $lookup_slot = *(unsigned int *)($esp + 8)
    set $lookup_state = *(unsigned int *)($esp + 12)
    set $lookup_mode = *(unsigned short *)($esp + 20)
    set $lookup_return = $lookup_caller
    printf "OLD_CLASS_LOOKUP_ENTER slot=%u mode=%u key:", $lookup_slot, $lookup_mode
    x/5ub $lookup_key
    tbreak *$lookup_return
    commands
      silent
      set $lookup_count = *(unsigned short *)($lookup_state + 0xf8e64)
      printf "OLD_CLASS_LOOKUP_RETURN slot=%u mode=%u appended=%u total=%u ids:", $lookup_slot, $lookup_mode, $eax, $lookup_count
      set $lookup_i = 0
      while $lookup_i < $lookup_count && $lookup_i < 32
        set $lookup_id = *(unsigned int *)($lookup_state + 0xf89b4 + $lookup_i * 4)
        printf " %u", $lookup_id
        set $lookup_features = *(unsigned int *)($lookup_state + 0x98)
        printf "\nOLD_CLASS_FEATURE id=%u:", $lookup_id
        x/7ub ($lookup_features + $lookup_id * 7)
        set $lookup_i = $lookup_i + 1
      end
      printf "\n"
      continue
    end
  end
  continue
end

break *0x1001e110
commands
  silent
  set $fallback_slot = *(unsigned int *)($esp + 4)
  set $fallback_return = *(unsigned int *)$esp
  printf "OLD_SPLIT_FALLBACK_ENTER slot=%u\n", $fallback_slot
  tbreak *$fallback_return
  commands
    silent
    printf "OLD_SPLIT_FALLBACK_RETURN slot=%u positions=%u\n", $fallback_slot, $eax
    continue
  end
  continue
end

break *0x1001e470
commands
  silent
  set $selector_context = *(unsigned int *)($esp + 8)
  set $selector_return = *(unsigned int *)$esp
  tbreak *$selector_return
  commands
    silent
    printf "OLD_SELECTOR_RETURN positions=%u\n", $eax
    continue
  end
  continue
end

  continue
end

continue
