set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hi-window-inputs-gdb.log
set logging overwrite on
set logging enabled on

set $hi_metric = -1
set $hi_join_calls = 0

# Accept the measured slot-1 total of 9, retaining the selected two-row chain.
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
    set $i = 0
    while $i < $count && $i < 96
      set $id = *(unsigned int *)($ids + $i * 4)
      printf "HI_WINDOW_METRIC id=%u weight=%u\n", $id, *(unsigned short *)($weights + $id * 2)
      set $i = $i + 1
    end
    tbreak *$return
    commands
      silent
      set $hi_metric = $eax
      printf "HI_WINDOW_METRIC_SUM value=%u\n", $eax
      continue
    end
  end
  continue
end

break *0x10024060
commands
  silent
  set $hi_metric = -1
  set $slot = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  tbreak *$return
  commands
    silent
    if $slot == 1 && $hi_metric == 9 && $al == 0
      printf "HI_WINDOW_CUTOFF_OVERRIDE slot=1 metric=9\n"
      set $eax = 1
    end
    continue
  end
  continue
end

# Match the 2006 caller's observed effective per-row volume scale.
hbreak *0x10027fe0
commands
  silent
  printf "HI_WINDOW_VOLUME old=%d new=100\n", *(int *)($esp + 12)
  set *(int *)($esp + 12) = 100
  continue
end

# Capture the actual input windows and pending carry before the direct join.
# The breakpoint only reads memory and writes evidence files; it does not edit
# engine buffers or registers.
hbreak *0x1002aac0
commands
  silent
  set $state = *(unsigned int *)($esp + 4)
  set $slot = *(short *)($esp + 8)
  set $hi_join_calls = $hi_join_calls + 1
  set $count = *(int *)($state + 0xede04 + $slot * 0x34)
  set $leading = *(short *)($state + 0xede1a + $slot * 0x34)
  set $trailing = *(short *)($state + 0xeddf8 + $slot * 0x34 + 0x24)
  set $carry = *(short *)($state + 0x47754)
  set $source = $state + 0x10eef4
  printf "HI_WINDOW_JOIN call=%u slot=%d count=%d leading=%d trailing=%d pending=%d source=%p\n", $hi_join_calls, $slot, $count, $leading, $trailing, $carry, $source
  if $slot == 0
    dump binary memory /work/corpus-parity/stage20/hi-window-inputs/row0-source.raw $source $source + $count * 2
  else
    if $slot == 1
      dump binary memory /work/corpus-parity/stage20/hi-window-inputs/row1-source.raw $source $source + $count * 2
    end
  end
  if $carry > 0
    dump binary memory /work/corpus-parity/stage20/hi-window-inputs/carry-before.raw $state + 0x10cd58 $state + 0x10cd58 + $carry * 2
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "HI_WINDOW_INPUT_TRACE_READY\n"
  continue
end

continue
