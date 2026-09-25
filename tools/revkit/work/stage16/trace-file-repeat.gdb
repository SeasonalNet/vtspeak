set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $fmt = *(int *)($esp + 4)
  set $text = *(unsigned int *)($esp + 8)
  set $path = *(unsigned int *)($esp + 12)
  set $speaker = *(int *)($esp + 16)
  set $pitch = *(int *)($esp + 20)
  set $speed = *(int *)($esp + 24)
  set $volume = *(int *)($esp + 28)
  set $pause = *(int *)($esp + 32)
  set $dict = *(int *)($esp + 36)
  set $texttype = *(int *)($esp + 40)
  printf "REPEAT_FIRST_ENTRY fmt=%d speaker=%d pitch=%d speed=%d volume=%d pause=%d dict=%d texttype=%d\n", $fmt, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype
  disable 1
  set $first_result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, (char *)$text, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "REPEAT_FIRST_RETURN result=%d\n", $first_result
  shell cp /work/stage5/output.wav /work/stage16/repeat-first.wav
  set $second_result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, (char *)$text, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "REPEAT_SECOND_RETURN result=%d\n", $second_result
  shell cp /work/stage5/output.wav /work/stage16/repeat-second.wav
  continue
end

continue
