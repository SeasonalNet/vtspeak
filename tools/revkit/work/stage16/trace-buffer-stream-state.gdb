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

  set $long_text = (char *)malloc(79)
  set $i = 0
  while $i < 78
    if ($i % 13) == 12
      set {char}($long_text + $i) = 32
    else
      set {char}($long_text + $i) = *(char *)($text + ($i % 13))
    end
    set $i = $i + 1
  end
  set {char}($long_text + 78) = 0

  set $buffer = (char *)malloc(1048576)
  set $length = (int *)malloc(4)
  set *$length = 1048576
  set $start = ((int (*)(int, char *, char *, int *, int, int, int, int, int, int, int, int, int))0x1001dc40)(0, $long_text, $buffer, $length, 0, 0, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "BUFFER_STATE_START return=%d output_len=%d\n", $start, *$length

  set *$length = 1048576
  set $busy = ((int (*)(int, char *, char *, int *, int, int, int, int, int, int, int, int, int))0x1001dc40)(0, $long_text, $buffer, $length, 0, 0, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "BUFFER_STATE_SECOND_START return=%d output_len=%d\n", $busy, *$length

  set *$length = 1048576
  set $cancel = ((int (*)(int, char *, char *, int *, int, int, int, int, int, int, int, int, int))0x1001dc40)(0, $long_text, $buffer, $length, 2, 0, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "BUFFER_STATE_CANCEL return=%d output_len=%d\n", $cancel, *$length

  set *$length = 1048576
  set $after_cancel = ((int (*)(int, char *, char *, int *, int, int, int, int, int, int, int, int, int))0x1001dc40)(0, $long_text, $buffer, $length, 1, 0, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "BUFFER_STATE_POLL_AFTER_CANCEL return=%d output_len=%d\n", $after_cancel, *$length
  continue
end

continue
