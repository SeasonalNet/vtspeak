set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

# This is the instruction after FUN_10024060 returns in FUN_10024680.
break *0x10024808
commands
  silent
  set $firstpass_index = *(unsigned int *)$esp
  set $firstpass_state = *(unsigned int *)($esp + 4)
  set $firstpass_row = $firstpass_state + 0xec628 + $firstpass_index * 6
  printf "CONTEXT_FIRSTPASS index=%u return=%d row4=%u\n", $firstpass_index, (short)$eax, *(unsigned char *)($firstpass_row + 4)
  continue
end

continue
