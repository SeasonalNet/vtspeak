set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/hello-window-inputs-2006/legacy-hello-window-inputs-gdb.log
set logging overwrite on
set logging enabled on
set $join_calls = 0
set $saved0 = 0
set $saved1 = 0
set $saved2 = 0
set $saved3 = 0

hbreak *0x10023b40
commands
  silent
  set $state = *(unsigned int *)($esp + 8)
  set $slot = *(short *)($esp + 12)
  set $count = *(int *)($state + 11000 + $slot * 4)
  set $source = *(unsigned int *)($state + $slot * 4)
  set $span = *(short *)($state + 0x6fbc + $slot * 2)
  set $join_calls = $join_calls + 1
  printf "LEGACY_HELLO_WINDOW_JOIN call=%u slot=%d count=%d span=%d source=%p\n", $join_calls, $slot, $count, $span, $source
  if $slot == 0 && $saved0 == 0
    dump binary memory /probe/stage20/hello-window-inputs-2006/row0-source.raw $source $source + $count * 2
    set $saved0 = 1
  end
  if $slot == 1 && $saved1 == 0
    dump binary memory /probe/stage20/hello-window-inputs-2006/row1-source.raw $source $source + $count * 2
    set $saved1 = 1
  end
  if $slot == 2 && $saved2 == 0
    dump binary memory /probe/stage20/hello-window-inputs-2006/row2-source.raw $source $source + $count * 2
    set $saved2 = 1
  end
  if $slot == 3 && $saved3 == 0
    dump binary memory /probe/stage20/hello-window-inputs-2006/row3-source.raw $source $source + $count * 2
    set $saved3 = 1
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "LEGACY_HELLO_WINDOW_TRACE_READY\n"
  continue
end

continue
