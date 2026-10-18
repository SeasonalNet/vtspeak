set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

set $unprepare_count = 0
break *0x100275a8
commands
  silent
  printf "STOP_WAVEOUT_RESET eax=%u handle=%#x play_state=%d\n", $eax, *(unsigned int *)0x100a7490, *(int *)0x100a7498
  continue
end

break *0x100275ca
commands
  silent
  set $unprepare_count = $unprepare_count + 1
  printf "STOP_WAVEOUT_UNPREPARE phase=1 count=%d eax=%u handle=%#x header=%#x\n", $unprepare_count, $eax, *(unsigned int *)0x100a7490, *(unsigned int *)($esi - 0xc)
  continue
end

break *0x100275f7
commands
  silent
  set $unprepare_count = $unprepare_count + 1
  printf "STOP_WAVEOUT_UNPREPARE phase=2 count=%d eax=%u handle=%#x header=%#x\n", $unprepare_count, $eax, *(unsigned int *)0x100a7490, *(unsigned int *)($esi - 0xc)
  continue
end

break *0x10027618
commands
  silent
  printf "STOP_WAVEOUT_CLOSE eax=%u handle_before_clear=%#x\n", $eax, *(unsigned int *)0x100a7490
  continue
end

break *0x10027622
commands
  silent
  printf "STOP_AFTER_CLOSE handle=%#x play_state=%d unprepare_calls=%d\n", *(unsigned int *)0x100a7490, *(int *)0x100a7498, $unprepare_count
  continue
end

break *0x100276a7
commands
  silent
  printf "STOP_API_EPILOGUE eax=%#x handle=%#x play_state=%d unprepare_calls=%d\n", $eax, *(unsigned int *)0x100a7490, *(int *)0x100a7498, $unprepare_count
  set $value = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {unsigned int}$value = 0x5a5a5a5a
  set $info_status = ((int (*)(int, char *, void *, int))0x1002a690)(101, (char *)0, (void *)$value, 4)
  printf "STOP_REQUEST101 status=%d output=%#x global=%d\n", $info_status, *(unsigned int *)$value, *(int *)0x100a7498
  continue
end

break *0x4016d9
commands
  silent
  printf "PLAY_API_RETURN ret=%d handle=%#x play_state=%d\n", (short)$eax, *(unsigned int *)0x100a7490, *(int *)0x100a7498
  call ((void (*)(void))0x10027010)()
  printf "PLAY_PAUSE_RETURNED\n"
  call ((void (*)(void))0x10027030)()
  printf "PLAY_RESTART_RETURNED\n"
  set $esp = $esp - 4
  set {unsigned int}$esp = 0x4016f9
  set $eip = 0x10027560
  continue
end
disable 7

break *0x4016f9
commands
  silent
  printf "PLAY_AFTER_STOP_RETURN play_state=%d\n", *(int *)0x100a7498
  continue
end
break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  disable 9
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
  printf "PLAY_STOP_REDIRECT text=%#x speaker=1\n", $text
  continue
end

continue
