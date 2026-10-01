set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-repeat-2013-slot8-split-threshold/adapted-split-threshold-gdb.log
set logging overwrite on
set logging enabled on

break *0x10024060
commands
  silent
  set $slot = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  printf "REPEAT_SPLIT_ENTER slot=%u\n", $slot
  tbreak *$return
  commands
    silent
    printf "REPEAT_SPLIT_RETURN slot=%u accepted=%u\n", $slot, $al
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
    printf "REPEAT_SPLIT_METRIC_ENTRY caller=%#x count=%u ids_weights:", $caller, $count
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
      printf "REPEAT_SPLIT_METRIC_RETURN sum=%u\n", $eax
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
  printf "REPEAT_SPLIT_FALLBACK_ENTER slot=%u\n", $slot
  tbreak *$return
  commands
    silent
    printf "REPEAT_SPLIT_FALLBACK_RETURN slot=%u positions=%u\n", $slot, $eax
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
    printf "REPEAT_POSITION_BUILDER_RETURN positions=%u\n", $eax
    continue
  end
  continue
end

break *0x408187
commands
  silent
  printf "REPEAT_SPLIT_TRACE_READY\n"
  continue
end

continue
