set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x100242a0
commands
  silent
  set $build_index = *(unsigned int *)($esp + 4)
  set $build_state = *(unsigned int *)($esp + 8)
  set $build_row = $build_state + 0xec628 + $build_index * 6
  printf "CONTEXT_ROW_BUILD index=%u row=%#x bytes=", $build_index, $build_row
  x/6ub $build_row
  bt 5
  continue
end

continue
