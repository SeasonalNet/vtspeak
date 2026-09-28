set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

define ex_drain_pause
  set $pause_value = $arg0
  set $total = 0
  set *$length = 60000
  set $result = ((int (*)(int, char *, char *, int *, int, int, unsigned int, unsigned int *, unsigned int *, int, int, int, int, int, int))0x1001ddf0)(1, (char *)$text, $buffer, $length, 0, 0, 1, 0, 0, -1, -1, -1, $pause_value, 0, 0)
  set $total = $total + *$length
  printf "SCALAR_PAUSE start value=%d result=%d bytes=%d total=%d\n", $pause_value, $result, *$length, $total
  set $poll = 1
  while ($result == 0) && ($poll < 32)
    set *$length = 60000
    set $result = ((int (*)(int, char *, char *, int *, int, int, unsigned int, unsigned int *, unsigned int *, int, int, int, int, int, int))0x1001ddf0)(1, (char *)$text, $buffer, $length, 1, 0, 1, 0, 0, -1, -1, -1, $pause_value, 0, 0)
    set $total = $total + *$length
    printf "SCALAR_PAUSE poll value=%d poll=%d result=%d bytes=%d total=%d\n", $pause_value, $poll, $result, *$length, $total
    set $poll = $poll + 1
  end
  printf "SCALAR_PAUSE_DONE selector=1 value=%d result=%d polls=%d total_bytes=%d\n", $pause_value, $result, $poll, $total
end

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  disable 1
  set $buffer = (char *)malloc(60000)
  set $length = (int *)malloc(4)
  ex_drain_pause 0
  ex_drain_pause 120
  ex_drain_pause 250
  ex_drain_pause 65535
  continue
end

continue
