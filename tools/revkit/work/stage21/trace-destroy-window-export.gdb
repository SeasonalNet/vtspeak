set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $destroy_hwnd_before = *(unsigned int *)0x100a8418
  printf "DESTROY_WINDOW_BEFORE hwnd=%#x\n", $destroy_hwnd_before
  call ((void (*)(void))0x10027fd0)()
  printf "DESTROY_WINDOW_FIRST raw_eax=%#x global_after=%#x\n", $eax, *(unsigned int *)0x100a8418
  call ((void (*)(void))0x10027fd0)()
  printf "DESTROY_WINDOW_SECOND raw_eax=%#x global_after=%#x\n", $eax, *(unsigned int *)0x100a8418
  continue
end

continue
