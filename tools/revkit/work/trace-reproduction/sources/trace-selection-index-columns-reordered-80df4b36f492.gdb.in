set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/versioned-index-columns-gdb.log
set logging overwrite on
set logging enabled on
set $index_loads = 0

break *0x10019940
commands
  silent
  set $index_loads = $index_loads + 1
  set $index = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  set $units = *(unsigned int *)($index + 0x4c)
  printf "INDEX_LOAD_BEGIN call=%u units=%u attr_a_ptr=%#x key_ptr=%#x attr_b_ptr=%#x\n", $index_loads, $units, *(unsigned int *)($index + 0x48), *(unsigned int *)($index + 0x44), *(unsigned int *)($index + 0x40)
  tbreak *$return
  commands
    silent
    printf "INDEX_LOAD_END call=%u attr_a=", $index_loads
    x/8bx *(unsigned int *)($index + 0x48)
    printf "INDEX_LOAD_END call=%u key=", $index_loads
    x/16bx *(unsigned int *)($index + 0x44)
    printf "INDEX_LOAD_END call=%u attr_b=", $index_loads
    x/8bx *(unsigned int *)($index + 0x40)
    continue
  end
  continue
end

continue
