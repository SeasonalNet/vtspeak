set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/hi-window-inputs-2006/legacy-hi-window-inputs-gdb.log
set logging overwrite on
set logging enabled on
set $hi_window_calls = 0
set $hi_row0_saved = 0
set $hi_row1_saved = 0

# FUN_10023b40 is the direct 2006 join reached by the selected Hi chain.
# Keep this capture read-only: dump the decoded row buffer before the routine
# changes its carry/output state.
hbreak *0x10023b40
commands
  silent
  set $state = *(unsigned int *)($esp + 8)
  set $slot = *(short *)($esp + 12)
  set $count = *(int *)($state + 11000 + $slot * 4)
  set $source = *(unsigned int *)($state + $slot * 4)
  set $carry = *(short *)($state + 0x6fbc + $slot * 2)
  set $hi_window_calls = $hi_window_calls + 1
  printf "LEGACY_HI_WINDOW_JOIN call=%u slot=%d count=%d span=%d source=%p\n", $hi_window_calls, $slot, $count, $carry, $source
  if $slot == 0 && $hi_row0_saved == 0
    dump binary memory /probe/stage20/hi-window-inputs-2006/row0-source.raw $source $source + $count * 2
    set $hi_row0_saved = 1
  end
  if $slot == 1 && $hi_row1_saved == 0
    dump binary memory /probe/stage20/hi-window-inputs-2006/row1-source.raw $source $source + $count * 2
    set $hi_row1_saved = 1
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "LEGACY_HI_WINDOW_TRACE_READY\n"
  continue
end

continue
