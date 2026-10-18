set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/old-cutoff-experiment-v4-gdb.log
set logging overwrite on
set logging enabled on
set $selected_calls = 0

break *0x10024060
commands
  silent
  set $slot = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  set $last_sum = -1
  printf "OLD_CUTOFF_SPLIT_ENTER slot=%u\n", $slot
  tbreak *$return
  commands
    silent
    set $native = $al
    if $native == 0 && $last_sum > 2
      printf "OLD_CUTOFF_OVERRIDE slot=%u sum=%u native=%u overridden=1\n", $slot, $last_sum, $native
      set $eax = 1
    else
      printf "OLD_CUTOFF_SPLIT_RETURN slot=%u sum=%d accepted=%u\n", $slot, $last_sum, $native
    end
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
    printf "OLD_CUTOFF_METRIC_ENTRY count=%u ids_weights:", $count
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
      set $last_sum = $eax
      printf "OLD_CUTOFF_METRIC_RETURN sum=%u\n", $eax
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
  printf "OLD_CUTOFF_FALLBACK_ENTER slot=%u\n", $slot
  tbreak *$return
  commands
    silent
    printf "OLD_CUTOFF_FALLBACK_RETURN slot=%u positions=%u\n", $slot, $eax
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
    printf "OLD_CUTOFF_POSITION_BUILDER_RETURN positions=%u\n", $eax
    continue
  end
  continue
end

break *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  if $selected_calls <= 24
    printf "OLD_CUTOFF_SELECTED_UNIT n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  end
  continue
end

break *0x408187
commands
  silent
  printf "OLD_CUTOFF_TRACE_READY\n"
  continue
end

continue
