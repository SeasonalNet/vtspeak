set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

set $destroy_window_calls = 0
break *0x10027fdc
commands
  silent
  set $destroy_window_calls = $destroy_window_calls + 1
  printf "DESTROY_WINDOW_CALL count=%d raw_destroywindow_result=%u hwnd=%#x\n", $destroy_window_calls, $eax, *(unsigned int *)0x100a8418
  if $destroy_window_calls == 1
    set $eip = 0x10027fd0
  else
    continue
  end
end
disable 1

break *0x10027fd6
commands
  silent
  printf "DESTROY_WINDOW_BEFORE_USER32 hwnd=%#x\n", *(unsigned int *)0x100a8418
  continue
end

break *0x1001da50
commands
  silent
  disable 3
  printf "DESTROY_WINDOW_TARGET_ENTRY hwnd=%#x\n", *(unsigned int *)0x100a8418
  set $eip = 0x10027fd0
  continue
end

continue
