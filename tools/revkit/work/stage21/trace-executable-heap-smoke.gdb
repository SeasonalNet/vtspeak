set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands 1
  silent
  disable 1
  set $code = (unsigned char *)malloc(1)
  set $code[0] = 0xc3
  set $returned = ((int (*)(void))$code)()
  printf "HEAP_CODE_RETURNED value=%d address=%p\n", $returned, $code
  continue
end

continue
