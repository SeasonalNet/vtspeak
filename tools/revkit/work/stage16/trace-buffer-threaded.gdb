set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10025fc0
commands
  silent
  printf "BUFFER_CONTEXT_ENTER out_status=%#x dict=%d speaker=%d\n", *(unsigned int *)($esp + 4), *(int *)($esp + 8), *(int *)($esp + 12)
  continue
end
disable 1

break *0x10020225
commands
  silent
  printf "BUFFER_CONTEXT_RETURN result=%#x status_hi=%d\n", $eax, *(short *)($ebp + 0x16)
  continue
end
disable 2

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
  disable 3
  enable 1
  enable 2
  set $buffer = (char *)malloc(1048576)
  set $length = (int *)malloc(4)
  set *$length = 1048576
  set $start = ((int (*)(int, char *, char *, int *, int, int, int, int, int, int, int, int, int))0x1001dc40)(0, (char *)$text, $buffer, $length, 2, 0, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "BUFFER_THREAD_START return=%d output_len=%d\n", $start, *$length
  set $poll = ((int (*)(int, char *, char *, int *, int, int, int, int, int, int, int, int, int))0x1001dc40)(0, (char *)$text, $buffer, $length, 1, 0, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "BUFFER_THREAD_POLL1 return=%d output_len=%d\n", $poll, *$length
  if ($poll == 1) && (*$length > 0) && (*$length <= 1048576)
    dump binary memory /work/stage16/buffer-threaded.bin $buffer $buffer+*$length
  end
  disable 1
  disable 2
  continue
end

continue
