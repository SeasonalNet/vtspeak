set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

set $heap_access_count = 0
set $heap_last_value = *(unsigned int *)0x100ff11c
awatch *(unsigned int *)0x100ff11c
commands
  silent
  set $heap_access_count = $heap_access_count + 1
  printf "HEAP_EXPORT_ACCESS count=%d pc=%#x before=%#x after=%#x\n", $heap_access_count, $eip, $heap_last_value, *(unsigned int *)0x100ff11c
  set $heap_last_value = *(unsigned int *)0x100ff11c
  bt 4
  continue
end

break *0x10027e70
commands
  silent
  printf "HEAP_EXPORT_BEFORE_LOAD value=%#x\n", *(unsigned int *)0x100ff11c
  disable 2
  continue
end

break *0x4016d9
commands
  silent
  printf "HEAP_EXPORT_AFTER_SYNTHESIS eax=%#x value=%#x accesses=%d\n", $eax, *(unsigned int *)0x100ff11c, $heap_access_count
  disable 3
  call ((void (*)(int))0x10027f80)(1)
  printf "HEAP_EXPORT_AFTER_UNLOAD value=%#x accesses=%d\n", *(unsigned int *)0x100ff11c, $heap_access_count
  continue
end

continue
