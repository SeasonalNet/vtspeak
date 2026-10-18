set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hi-plain-2013-metric4-volume100.log
set logging overwrite on
set logging enabled on

set $hi_metric = -1
set $hi_split = 0
set $hi_fallback = 0
set $hi_selected = 0

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
    printf "HI4_METRIC_INPUT count=%u id_weight=", $count
    set $i = 0
    while $i < $count && $i < 96
      set $id = *(unsigned int *)($ids + $i * 4)
      printf "%u:%u ", $id, *(unsigned short *)($weights + $id * 2)
      set $i = $i + 1
    end
    printf "\n"
    tbreak *$return
    commands
      silent
      set $hi_metric = $eax
      printf "HI4_METRIC_SUM value=%u\n", $eax
      continue
    end
  end
  continue
end

break *0x10024060
commands
  silent
  set $hi_split = $hi_split + 1
  set $hi_metric = -1
  set $slot = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  printf "HI4_SPLIT_ENTER call=%u slot=%u\n", $hi_split, $slot
  tbreak *$return
  commands
    silent
    set $native = $al
    if $slot == 1 && $hi_metric == 4 && $native == 0
      printf "HI4_CUTOFF_OVERRIDE slot=%u metric=%u native=%u forced=1\n", $slot, $hi_metric, $native
      set $eax = 1
    else
      printf "HI4_SPLIT_RETURN slot=%u metric=%d accepted=%u\n", $slot, $hi_metric, $native
    end
    continue
  end
  continue
end

break *0x100242a0
commands
  silent
  set $hi_fallback = $hi_fallback + 1
  set $slot = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  printf "HI4_FALLBACK_ENTER call=%u slot=%u\n", $hi_fallback, $slot
  tbreak *$return
  commands
    silent
    printf "HI4_FALLBACK_RETURN slot=%u positions=%u\n", $slot, $eax
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
    printf "HI4_POSITION_BUILDER_RETURN positions=%u\n", $eax
    continue
  end
  continue
end

break *0x100235f6
commands
  silent
  set $state = *(unsigned int *)($ebp + 0xc)
  set $position = *(int *)($ebp + 8)
  if $position < 4
    set $count = *(short *)($state + 0xae988 + $position * 0xfc)
    printf "HI4_FINAL_CANDIDATES position=%d count=%d ids=", $position, $count
    set $i = 0
    while $i < $count && $i < 80
      set $node = *(unsigned int *)($state + 0x477ac + $i * 4)
      printf "%u ", *(unsigned int *)($node + 8)
      set $i = $i + 1
    end
    printf "\n"
  end
  continue
end

hbreak *0x10027fe0
commands
  silent
  printf "HI4_VOLUME old=%d new=100\n", *(int *)($esp + 12)
  set *(int *)($esp + 12) = 100
  continue
end

hbreak *0x1001b200
commands
  silent
  set $hi_selected = $hi_selected + 1
  printf "HI4_SELECTED n=%u id=%u\n", $hi_selected, *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "HI4_CUTOFF_OVERRIDE_TRACE_READY\n"
  continue
end

continue
