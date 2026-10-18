set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10027e70
commands
  silent
  printf "HEAP_START before_load=%#x\n", *(unsigned int *)0x100ff11c
  disable 1
  continue
end

break *0x1001da50
commands
  silent
  printf "HEAP_START after_load=%#x\n", *(unsigned int *)0x100ff11c
  disable 2
  set $unload = ((void (*)(int))0x10027f80)(1)
  printf "HEAP_START after_unload=%#x unload_raw_eax=%#x\n", *(unsigned int *)0x100ff11c, $eax
  continue
end

continue
