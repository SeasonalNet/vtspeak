set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-no-marker-window-inputs/hello-no-marker-window-inputs-gdb.log
set logging overwrite on
set logging enabled on
set $selected_calls = 0
set $join_calls = 0

hbreak *0x10027fe0
commands
  silent
  printf "HELLO_NO_MARKER_VOLUME old=%d new=100\n", *(int *)($esp + 12)
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
  printf "HELLO_NO_MARKER_JOIN call=%u slot=%d count=%d leading=%d trailing=%d source=%p\n", $join_calls, $slot, $count, $leading, $trailing, $source
  if $join_calls <= 8
    if $slot == 0
      dump binary memory /work/corpus-parity/stage20/hello-no-marker-window-inputs/row0-source.raw $source $source + $count * 2
    end
    if $slot == 1
      dump binary memory /work/corpus-parity/stage20/hello-no-marker-window-inputs/row1-source.raw $source $source + $count * 2
    end
    if $slot == 2
      dump binary memory /work/corpus-parity/stage20/hello-no-marker-window-inputs/row2-source.raw $source $source + $count * 2
    end
    if $slot == 3
      dump binary memory /work/corpus-parity/stage20/hello-no-marker-window-inputs/row3-source.raw $source $source + $count * 2
    end
    if $slot == 4
      dump binary memory /work/corpus-parity/stage20/hello-no-marker-window-inputs/row4-source.raw $source $source + $count * 2
    end
    if $slot == 5
      dump binary memory /work/corpus-parity/stage20/hello-no-marker-window-inputs/row5-source.raw $source $source + $count * 2
    end
    if $slot == 6
      dump binary memory /work/corpus-parity/stage20/hello-no-marker-window-inputs/row6-source.raw $source $source + $count * 2
    end
    if $slot == 7
      dump binary memory /work/corpus-parity/stage20/hello-no-marker-window-inputs/row7-source.raw $source $source + $count * 2
    end
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  printf "HELLO_NO_MARKER_SELECTED n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "HELLO_NO_MARKER_TRACE_READY\n"
  continue
end

continue
