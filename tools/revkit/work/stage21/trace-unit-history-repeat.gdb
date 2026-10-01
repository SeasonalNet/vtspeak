set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10027af0
commands
  silent
  set $history_load_return = *(unsigned int *)$esp
  call ((void (*)(unsigned char))0x100283f0)(1)
  printf "UNIT_HISTORY_REPEAT_MODE value=%u\n", *(unsigned char *)0x100a045d
  disable 1
  tbreak *$history_load_return
  commands
    silent
    printf "UNIT_HISTORY_REPEAT_LOAD ax=%d\n", (short)$eax
    continue
  end
  continue
end

break *0x1001b240
commands
  silent
  set $history_item = *(unsigned int *)($esp + 4)
  set $history_path = (char *)*(unsigned int *)($esp + 8)
  set $history_return = *(unsigned int *)$esp
  printf "UNIT_HISTORY_REPEAT_ENTRY count=%d path=%s\n", *(int *)($history_item + 0x4c), $history_path
  tbreak *$history_return
  commands
    silent
    printf "UNIT_HISTORY_REPEAT_RETURN ax=%d\n", (short)$eax
    continue
  end
  continue
end

break *0x1001b25c
commands
  silent
  printf "UNIT_HISTORY_REPEAT_OPEN source=%p\n", $edi
  continue
end

break *0x1001b28b
commands
  silent
  set $history_item = *(unsigned int *)($ebp + 8)
  set $history_count = *(int *)($history_item + 0x4c)
  set $history_a = (unsigned int *)*(unsigned int *)($history_item + 0x14)
  set $history_b = (unsigned int *)*(unsigned int *)($history_item + 0x18)
  printf "UNIT_HISTORY_REPEAT_DATA count=%d a_first=%u a_last=%u b_first=%u b_last=%u\n", $history_count, *$history_a, *($history_a + $history_count - 1), *$history_b, *($history_b + $history_count - 1)
  continue
end

break *0x1001da50
commands
  silent
  set $history_text_return = *(unsigned int *)$esp
  set $history_fmt = *(int *)($esp + 4)
  set $history_text = (char *)*(unsigned int *)($esp + 8)
  set $history_speaker = *(int *)($esp + 16)
  set $history_pitch = *(int *)($esp + 20)
  set $history_speed = *(int *)($esp + 24)
  set $history_volume = *(int *)($esp + 28)
  set $history_pause = *(int *)($esp + 32)
  set $history_dict = *(int *)($esp + 36)
  set $history_texttype = *(int *)($esp + 40)
  disable 5
  tbreak *$history_text_return
  commands
    silent
    set $history_first_eax = $eax
    printf "UNIT_HISTORY_REPEAT_CALL index=0 ax=%d eax=%#x\n", (short)$eax, $eax
    set $history_repeat_1 = ((int (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($history_fmt, $history_text, (char *)"repeat1.wav", $history_speaker, $history_pitch, $history_speed, $history_volume, $history_pause, $history_dict, $history_texttype)
    printf "UNIT_HISTORY_REPEAT_CALL index=1 ax=%d\n", (short)$history_repeat_1
    set $history_repeat_2 = ((int (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($history_fmt, $history_text, (char *)"repeat2.wav", $history_speaker, $history_pitch, $history_speed, $history_volume, $history_pause, $history_dict, $history_texttype)
    printf "UNIT_HISTORY_REPEAT_CALL index=2 ax=%d\n", (short)$history_repeat_2
    set $history_repeat_3 = ((int (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($history_fmt, $history_text, (char *)"repeat3.wav", $history_speaker, $history_pitch, $history_speed, $history_volume, $history_pause, $history_dict, $history_texttype)
    printf "UNIT_HISTORY_REPEAT_CALL index=3 ax=%d\n", (short)$history_repeat_3
    set $eax = $history_first_eax
    continue
  end
  continue
end

continue
