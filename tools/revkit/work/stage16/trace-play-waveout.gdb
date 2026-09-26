set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10026c1d
commands
  silent
  printf "PLAY_WAVEOUT_OPEN_CALL handle_ptr=%#x device_id=%d format_ptr=%#x callback=%#x instance=%#x flags=%#x\n", *(unsigned int *)$esp, *(int *)($esp + 4), *(unsigned int *)($esp + 8), *(unsigned int *)($esp + 12), *(unsigned int *)($esp + 16), *(unsigned int *)($esp + 20)
  printf "PLAY_WAVEOUT_FORMAT tag=%u channels=%u rate=%u avg_bytes=%u align=%u bits=%u\n", *(unsigned short *)0x100a9d20, *(unsigned short *)0x100a9d22, *(unsigned int *)0x100a9d24, *(unsigned int *)0x100a9d28, *(unsigned short *)0x100a9d2c, *(unsigned short *)0x100a9d2e
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
  printf "PLAY_WAVEOUT_API_RETURN ret=%d handle=%#x play_state=%d\n", (short)$eax, *(unsigned int *)0x100a7490, *(int *)0x100a7498
  call ((void (*)(void))0x10027010)()
  printf "PLAY_PAUSE_RETURNED\n"
  call ((void (*)(void))0x10027030)()
  printf "PLAY_RESTART_RETURNED\n"
  call ((void (*)(void))0x10027560)()
  printf "PLAY_STOP_RETURNED\n"
  continue
end
disable 3

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
  printf "PLAY_WAVEOUT_REDIRECT retaddr=%#x text=%#x speaker=1\n", $ret, $text
  continue
end

continue
