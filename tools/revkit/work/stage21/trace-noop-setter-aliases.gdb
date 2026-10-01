set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x4016d9
commands
  silent
  printf "NOOP_ALIAS_CONTROL text_to_file_eax=%#x speaker1=%#x play_state=%d history_mode=%d\n", $eax, *(unsigned int *)0x100a0468, *(int *)0x100a7498, *(unsigned char *)0x100a045d

  set $eax = 0x11223344
  set $esp_before = $esp
  call ((void (*)(void))0x10028420)()
  printf "NOOP_ALIAS name=VT_SetDecimal0Pron_ENG eax=%#x esp_before=%#x esp_after=%#x speaker1=%#x play_state=%d history_mode=%d\n", $eax, $esp_before, $esp, *(unsigned int *)0x100a0468, *(int *)0x100a7498, *(unsigned char *)0x100a045d

  set $eax = 0x55667788
  set $esp_before = $esp
  call ((void (*)(void))0x10028420)()
  printf "NOOP_ALIAS name=VT_SetPhone0Pron_ENG eax=%#x esp_before=%#x esp_after=%#x speaker1=%#x play_state=%d history_mode=%d\n", $eax, $esp_before, $esp, *(unsigned int *)0x100a0468, *(int *)0x100a7498, *(unsigned char *)0x100a045d

  set $eax = 0x13579bdf
  set $esp_before = $esp
  call ((void (*)(void))0x10028420)()
  printf "NOOP_ALIAS name=VT_SetVirtualTagMode_ENG eax=%#x esp_before=%#x esp_after=%#x speaker1=%#x play_state=%d history_mode=%d\n", $eax, $esp_before, $esp, *(unsigned int *)0x100a0468, *(int *)0x100a7498, *(unsigned char *)0x100a045d

  disable 1
  continue
end

continue
