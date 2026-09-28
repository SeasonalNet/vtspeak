set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  disable 1
  set $buffer = (char *)malloc(60000)
  set $length = (int *)malloc(4)
  set *$length = 60000
  set $armed = 1
  break *0x10021f12
  commands
    silent
    if $armed
      set $normal_cursor = *(int *)($ebp - 4)
      set $copy_limit = *(int *)($ebp - 8)
      set *(int *)($ebp - 4) = $copy_limit
      set $armed = 0
      printf "FORCED_MINUS8_INJECT selector=0 normal_cursor=%d copy_limit=%d\n", $normal_cursor, $copy_limit
    end
    continue
  end
  set $result = ((int (*)(int, char *, char *, int *, int, int, unsigned int, unsigned int *, unsigned int *, int, int, int, int, int, int))0x1001ddf0)(0, (char *)$text, $buffer, $length, 0, 0, 1, 0, 0, -1, -1, -1, -1, 0, 0)
  printf "FORCED_MINUS8_RETURN selector=0 result=%d length=%d\n", $result, *$length
  set *$length = 60000
  set $cancel = ((int (*)(int, char *, char *, int *, int, int, unsigned int, unsigned int *, unsigned int *, int, int, int, int, int, int))0x1001ddf0)(0, $text, $buffer, $length, 2, 0, 1, 0, 0, -1, -1, -1, -1, 0, -1)
  printf "FORCED_MINUS8_CANCEL selector=0 result=%d length=%d\n", $cancel, *$length
  set *$length = 60000
  set $after = ((int (*)(int, char *, char *, int *, int, int, unsigned int, unsigned int *, unsigned int *, int, int, int, int, int, int))0x1001ddf0)(0, $text, $buffer, $length, 1, 0, 1, 0, 0, -1, -1, -1, -1, 0, 0)
  printf "FORCED_MINUS8_AFTER_CANCEL selector=0 result=%d length=%d\n", $after, *$length
  continue
end

continue
