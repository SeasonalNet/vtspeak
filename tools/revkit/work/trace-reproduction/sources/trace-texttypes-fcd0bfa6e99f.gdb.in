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
  disable 1
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, (char *)$text, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, -1)
  printf "TEXTTYPE value=-1 return=%d\n", $result
  shell cp /work/stage5/output.wav /work/stage16/texttype--1.wav
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, (char *)$text, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, 0)
  printf "TEXTTYPE value=0 return=%d\n", $result
  shell cp /work/stage5/output.wav /work/stage16/texttype-0.wav
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, (char *)$text, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, 4)
  printf "TEXTTYPE value=4 return=%d\n", $result
  shell cp /work/stage5/output.wav /work/stage16/texttype-4.wav
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, (char *)$text, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, 6)
  printf "TEXTTYPE value=6 return=%d\n", $result
  shell cp /work/stage5/output.wav /work/stage16/texttype-6.wav
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, (char *)$text, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, 8)
  printf "TEXTTYPE value=8 return=%d\n", $result
  shell cp /work/stage5/output.wav /work/stage16/texttype-8.wav
  continue
end

continue
