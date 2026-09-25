set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10018770
commands
  silent
  set $middle_caller = *(unsigned int *)$esp
  if $middle_caller == 0x10024119 || $middle_caller == 0x100240f4
    set $middle_input = *(unsigned int *)($esp + 4)
    set $middle_index = *(unsigned int *)($esp + 8)
    set $middle_state = *(unsigned int *)($esp + 12)
    set $middle_mode = *(short *)($esp + 20)
    set $middle_row = $middle_state + 0xec628 + $middle_index * 6
    printf "MIDDLE_INPUT index=%u mode=%d caller=%#x features=", $middle_index, $middle_mode, $middle_caller
    x/7ub $middle_input
    printf "MIDDLE_ROW index=%u bytes=", $middle_index
    x/6ub $middle_row
  end
  continue
end

break *0x10024808
commands
  silent
  set $middle_return_index = *(unsigned int *)$esp
  set $middle_return_state = *(unsigned int *)($esp + 4)
  set $middle_return_row = $middle_return_state + 0xec628 + $middle_return_index * 6
  printf "MIDDLE_FIRSTPASS index=%u return=%d row4=%u\n", $middle_return_index, (short)$eax, *(unsigned char *)($middle_return_row + 4)
  continue
end

continue
