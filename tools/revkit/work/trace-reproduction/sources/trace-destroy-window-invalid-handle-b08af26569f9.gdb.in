set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

set $destroy_invalid_stage = 0
break *0x10027fdc
commands
  silent
  set $destroy_invalid_stage = $destroy_invalid_stage + 1
  printf "DESTROY_WINDOW_INVALID_RESULT stage=%d eax=%u handle=%#x\n", $destroy_invalid_stage, $eax, *(unsigned int *)0x100a8418
  if $destroy_invalid_stage == 1
    set *(unsigned int *)0x100a8418 = 0xdeadbeef
    set $eip = 0x10027fd0
    continue
  else
    set *(unsigned int *)0x100a8418 = $destroy_original_handle
    printf "DESTROY_WINDOW_RESTORED handle=%#x\n", *(unsigned int *)0x100a8418
    continue
  end
end

break *0x1001da50
commands
  silent
  disable 2
  set $destroy_original_handle = *(unsigned int *)0x100a8418
  printf "DESTROY_WINDOW_ORIGINAL handle=%#x\n", $destroy_original_handle
  set *(unsigned int *)0x100a8418 = 0
  set $eip = 0x10027fd0
  continue
end

continue
