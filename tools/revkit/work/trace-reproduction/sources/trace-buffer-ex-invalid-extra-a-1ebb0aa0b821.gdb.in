set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands 1
  silent
  set $text = *(unsigned int *)($esp + 8)
  disable 1
  set $buffer = (char *)malloc(60000)
  set $length = (int *)malloc(4)
  set $extra_a = (unsigned int *)1
  set $extra_b = (unsigned int *)0
  set *$length = 60000
  printf "INVALID_EXTRA_CALL slot=a extra_a=%#x extra_b=%#x\n", $extra_a, $extra_b
  set $result = ((int (*)(int, char *, char *, int *, int, int, unsigned int, unsigned int *, unsigned int *, int, int, int, int, int, int))0x1001ddf0)(0, (char *)$text, $buffer, $length, 0, 0, 1, $extra_a, $extra_b, -1, -1, -1, -1, 0, 0)
  printf "INVALID_EXTRA_RETURN slot=a result=%d length=%d\n", $result, *$length
  continue
end

continue
