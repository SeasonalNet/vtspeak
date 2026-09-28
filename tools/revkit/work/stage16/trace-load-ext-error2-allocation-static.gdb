set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10027af0
commands 1
  silent
  printf "ERROR2_PROBE loader_entered\n"
  disable 1
  set $load_return = *(unsigned int *)$esp
  tbreak *$load_return
  commands 2
    silent
    printf "ERROR2_PROBE exported_load_ax=%d eax=%#x\n", (short)$ax, $eax
    continue
  end
  break *0x10025e88
  commands 3
    silent
    printf "ERROR2_PROBE shared_allocation_returned original_eax=%#x forced_eax=0\n", $eax
    set $eax = 0
    disable 3
    continue
  end
  continue
end

continue
