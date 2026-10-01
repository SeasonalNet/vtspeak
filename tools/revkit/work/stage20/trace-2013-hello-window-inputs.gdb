set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-window-inputs/hello-window-inputs-gdb.log
set logging overwrite on
set logging enabled on
set $last_sum = -1
set $selected_calls = 0
set $join_calls = 0

break *0x10024060
commands
  silent
  set $slot = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  set $last_sum = -1
  tbreak *$return
  commands
    silent
    set $native = $al
    if $native == 0 && $last_sum > 2
      printf "HELLO_WINDOW_CUTOFF_OVERRIDE slot=%u metric_sum=%u\n", $slot, $last_sum
      set $eax = 1
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
    set $last_sum = 0
    set $i = 0
    while $i < $count && $i < 64
      set $id = *(unsigned int *)($ids + $i * 4)
      set $last_sum = $last_sum + *(unsigned short *)($weights + $id * 2)
      set $i = $i + 1
    end
    tbreak *$return
    commands
      silent
      set $last_sum = $eax
      continue
    end
  end
  continue
end

hbreak *0x10027fe0
commands
  silent
  printf "HELLO_WINDOW_VOLUME old=%d new=100\n", *(int *)($esp + 12)
  set *(int *)($esp + 12) = 100
  continue
end

hbreak *0x1002aac0
commands
  silent
  set $state = *(unsigned int *)($esp + 4)
  set $slot = *(short *)($esp + 8)
  set $count = *(int *)($state + 0xede04 + $slot * 0x34)
  set $leading = *(short *)($state + 0xede1a + $slot * 0x34)
  set $trailing = *(short *)($state + 0xeddf8 + $slot * 0x34 + 0x24)
  set $join_calls = $join_calls + 1
  set $source = $state + 0x10eef4
  printf "HELLO_WINDOW_JOIN call=%u slot=%d count=%d leading=%d trailing=%d source=%p\n", $join_calls, $slot, $count, $leading, $trailing, $source
  if $join_calls <= 4
    if $slot == 0
      dump binary memory /work/corpus-parity/stage20/hello-window-inputs/row0-source.raw $source $source + $count * 2
    end
    if $slot == 1
      dump binary memory /work/corpus-parity/stage20/hello-window-inputs/row1-source.raw $source $source + $count * 2
    end
    if $slot == 2
      dump binary memory /work/corpus-parity/stage20/hello-window-inputs/row2-source.raw $source $source + $count * 2
    end
    if $slot == 3
      dump binary memory /work/corpus-parity/stage20/hello-window-inputs/row3-source.raw $source $source + $count * 2
    end
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  printf "HELLO_WINDOW_SELECTED n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "HELLO_WINDOW_TRACE_READY\n"
  continue
end

continue
