set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  set $speaker = *(int *)($esp + 16)
  set $pitch = *(int *)($esp + 20)
  set $speed = *(int *)($esp + 24)
  set $volume = *(int *)($esp + 28)
  set $pause = *(int *)($esp + 32)
  set $dict = *(int *)($esp + 36)
  set $texttype = *(int *)($esp + 40)
  disable 1
  set $buffer = (char *)malloc(1048576)
  set $length = (int *)malloc(4)
  set *$length = 1048576
  set $result = ((int (*)(int, char *, char *, int *, int, int, int, int, int, int, int, int, int))0x1001dc40)(4, (char *)$text, $buffer, $length, 0, 0, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "BUFFER_ERROR invalid-format result=%d\n", $result
  set $result = ((int (*)(int, char *, char *, int *, int, int, int, int, int, int, int, int, int))0x1001dc40)(0, (char *)0, $buffer, $length, 0, 0, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "BUFFER_ERROR null-text result=%d\n", $result
  set $first_byte = *(unsigned char *)$text
  set *(unsigned char *)$text = 0
  set $result = ((int (*)(int, char *, char *, int *, int, int, int, int, int, int, int, int, int))0x1001dc40)(0, (char *)$text, $buffer, $length, 0, 0, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "BUFFER_ERROR empty-text result=%d\n", $result
  set *(unsigned char *)$text = $first_byte
  set $result = ((int (*)(int, char *, char *, int *, int, int, int, int, int, int, int, int, int))0x1001dc40)(0, (char *)$text, (char *)0, $length, 0, 0, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "BUFFER_ERROR null-buffer result=%d\n", $result
  continue
end

continue
