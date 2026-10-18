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

break *0x1002746f
commands
  silent
  printf "PLAY_CALLBACK_STATE_CLEAR eip=%#x handle=%#x state=%d message_index=%d\n", $eip, *(unsigned int *)0x100a7490, *(int *)0x100a7498, $delivered
  query_play_state 4
  continue
end

break *0x4016d9
commands
  silent
  printf "PLAY_API_RETURN ret=%d handle=%#x play_state=%d callback_window=%#x\n", (short)$eax, *(unsigned int *)0x100a7490, *(int *)0x100a7498, *(unsigned int *)0x100a8418
  query_play_state 1
  set $delivered = 0
  set $msg = (void *)($esp - 0x100)
  set $window = *(unsigned int *)0x100a8418
  set $i = 0
  set $delivered = 0
  watch -l *(int *)0x100a7498
  commands
    silent
    printf "PLAY_STATE_WRITE eip=%#x value=%d message_index=%d\n", $eip, *(int *)0x100a7498, $delivered
    query_play_state 5
    continue
  end
  while $i < 1000 && *(int *)0x100a7498 != 0
    set $got = ((int (*)(void *, unsigned int, unsigned int, unsigned int, unsigned int))*(void **)0x1006d120)($msg, $window, 0, 0, 1)
    if $got
      set $delivered = $delivered + 1
      set $message = *(unsigned int *)((char *)$msg + 4)
      set $wparam = *(unsigned int *)((char *)$msg + 8)
      set $lparam = *(unsigned int *)((char *)$msg + 12)
      printf "PLAY_MESSAGE index=%d message=%#x wparam=%#x lparam=%#x state_before=%d\n", $delivered, $message, $wparam, $lparam, *(int *)0x100a7498
      disable 1 4
      call ((int (*)(void *))*(void **)0x1006d124)($msg)
      call ((int (*)(void *))*(void **)0x1006d128)($msg)
      enable 1 4
      query_play_state 2
    else
      call ((void (*)(unsigned int))*(void **)0x1006d028)(10)
    end
    set $i = $i + 1
  end
  printf "PLAY_PUMP_END iterations=%d delivered=%d state=%d\n", $i, $delivered, *(int *)0x100a7498
  query_play_state 3
  set $eip = 0x10027560
  continue
end

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  disable 3
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
  printf "PLAY_REDIRECT text=%#x speaker=1\n", $text
  continue
end

continue
