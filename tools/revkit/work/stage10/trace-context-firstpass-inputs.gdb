set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10018770
commands
  silent
  set $call_site = *(unsigned int *)$esp
  if $call_site == 0x10024119 || $call_site == 0x100240f4
    set $input_ptr = *(unsigned int *)($esp + 4)
    set $input_index = *(unsigned int *)($esp + 8)
    set $input_state = *(unsigned int *)($esp + 12)
    set $input_mode = *(short *)($esp + 20)
    set $input_row = $input_state + 0xec628 + $input_index * 6
    printf "FIRSTPASS_INPUT context=%u mode=%d caller=%#x features=", $input_index, $input_mode, $call_site
    x/7ub $input_ptr
    printf "FIRSTPASS_ROW context=%u bytes=", $input_index
    x/6ub $input_row
  end
  continue
end

continue
