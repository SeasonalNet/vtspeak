set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break LoadLibraryA
commands 1
  silent
  printf "ERROR2_PROBE dynamic_module_load path="
  x/s *(char **)($esp + 4)
  finish
  printf "ERROR2_PROBE dynamic_module_handle=%#x\n", $eax
  disable 1
  break *0x10027af0
  commands 2
    silent
    printf "ERROR2_PROBE loader_entered\n"
    disable 2
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
end

continue
