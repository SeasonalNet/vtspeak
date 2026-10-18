set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

define query_play_state
  set $state_ptr = (int *)malloc(4)
  set *$state_ptr = 0x5a5a5a5a
  set $state_result = ((int (*)(int, char *, void *, int))0x1002a690)(101, (char *)0, (void *)$state_ptr, 4)
  printf "PLAY_STATE_QUERY phase=%d result=%d output=%#x global=%d\n", $arg0, $state_result, *$state_ptr, *(int *)0x100a7498
end

break *0x10026c1d
commands
  silent
  printf "PLAY_WAVEOUT_OPEN_CALL\n"
  continue
end
disable 1

break *0x10026c23
commands
  silent
  printf "PLAY_WAVEOUT_OPEN_RETURN mmresult=%u handle=%#x\n", $eax, *(unsigned int *)0x100a7490
  continue
end
disable 2

break *0x4016d9
commands
  silent
  printf "PLAY_API_RETURN ret=%d handle=%#x global=%d\n", (short)$eax, *(unsigned int *)0x100a7490, *(int *)0x100a7498
  query_play_state 1
  call ((void (*)(void))0x10027010)()
  printf "PLAY_PAUSE_RETURNED global=%d\n", *(int *)0x100a7498
  query_play_state 2
  call ((void (*)(void))0x10027030)()
  printf "PLAY_RESTART_RETURNED global=%d\n", *(int *)0x100a7498
  query_play_state 3
  enable 5
  set $esp = $esp - 4
  set {unsigned int}$esp = 0x4016f9
  set $eip = 0x10027560
  continue
end
disable 3

break *0x4016f9
commands
  silent
  printf "PLAY_AFTER_STOP_RETURN global=%d\n", *(int *)0x100a7498
  query_play_state 4
  continue
end
disable 5

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  set $ret = *(unsigned int *)$esp
  disable 4
  enable 1
  enable 2
  enable 3
  set {int}($esp + 4) = 0
  set {int}($esp + 8) = 0
  set {int}($esp + 12) = $text
  set {int}($esp + 16) = 1
  set {int}($esp + 20) = -1
  set {int}($esp + 24) = -1
  set {int}($esp + 28) = -1
  set {int}($esp + 32) = -1
  set {int}($esp + 36) = -1
  set {int}($esp + 40) = 0
  set $eip = 0x10027050
  printf "PLAY_REDIRECT retaddr=%#x text=%#x speaker=1\n", $ret, $text
  continue
end

continue
