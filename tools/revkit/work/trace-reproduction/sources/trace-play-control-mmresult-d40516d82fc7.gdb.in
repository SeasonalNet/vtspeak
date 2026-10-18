set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

set $pause_count = 0
set $restart_count = 0

break *0x10027010
commands
  silent
  set $pause_count = $pause_count + 1
  printf "PLAY_CONTROL pause_enter=%d handle=%#x\n", $pause_count, *(unsigned int *)0x100a7490
  continue
end

break *0x10027020
commands
  silent
  printf "PLAY_CONTROL pause_mmresult_call=%d result=%u\n", $pause_count, $eax
  continue
end

break *0x10027030
commands
  silent
  set $restart_count = $restart_count + 1
  printf "PLAY_CONTROL restart_enter=%d handle=%#x\n", $restart_count, *(unsigned int *)0x100a7490
  continue
end

break *0x10027040
commands
  silent
  printf "PLAY_CONTROL restart_mmresult_call=%d result=%u\n", $restart_count, $eax
  continue
end

break *0x4016d9
commands
  silent
  printf "PLAY_CONTROL playback_return ret=%d handle=%#x state=%d\n", (short)$eax, *(unsigned int *)0x100a7490, *(int *)0x100a7498
  call ((void (*)(void))0x10027010)()
  printf "PLAY_CONTROL pause_call=1 complete\n"
  call ((void (*)(void))0x10027010)()
  printf "PLAY_CONTROL pause_call=2 complete\n"
  call ((void (*)(void))0x10027030)()
  printf "PLAY_CONTROL restart_call=1 complete\n"
  call ((void (*)(void))0x10027030)()
  printf "PLAY_CONTROL restart_call=2 complete\n"
  call ((void (*)(void))0x10027560)()
  printf "PLAY_CONTROL stop_call complete\n"
  continue
end
disable 5

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  set $ret = *(unsigned int *)$esp
  disable 6
  enable 5
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
  printf "PLAY_CONTROL redirect retaddr=%#x text=%#x speaker=1\n", $ret, $text
  continue
end

continue
