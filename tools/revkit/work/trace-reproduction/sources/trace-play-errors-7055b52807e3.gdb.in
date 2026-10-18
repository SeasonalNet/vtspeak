set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  disable 1
  set $null_result = ((short (*)(void *, unsigned int, char *, int, int, int, int, int, int, int))0x10027050)(0, 0, (char *)0, -1, -1, -1, -1, -1, -1, -1)
  printf "PLAY_NULL_TEXT_RETURN ret=%d\n", $null_result
  set $empty = (char *)malloc(1)
  set *$empty = 0
  set $empty_result = ((short (*)(void *, unsigned int, char *, int, int, int, int, int, int, int))0x10027050)(0, 0, $empty, -1, -1, -1, -1, -1, -1, -1)
  printf "PLAY_EMPTY_TEXT_RETURN ret=%d\n", $empty_result
  set $play_result = ((short (*)(void *, unsigned int, char *, int, int, int, int, int, int, int))0x10027050)(0, 0, (char *)$text, -1, -1, -1, -1, -1, -1, -1)
  printf "PLAY_VALID_TEXT_RETURN ret=%d\n", $play_result
  call ((void (*)(void))0x10027560)()
  printf "PLAY_POST_VALID_STOP_RETURN\n"
  call ((void (*)(void))0x10027010)()
  printf "PLAY_PAUSE_RETURN\n"
  call ((void (*)(void))0x10027030)()
  printf "PLAY_RESTART_RETURN\n"
  call ((void (*)(void))0x10027560)()
  printf "PLAY_STOP_RETURN\n"
  continue
end

continue
