set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  set $path = *(unsigned int *)($esp + 12)
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
  printf "BUFFER_STREAM_TEXT chars=78 repeats=6 sample=Hello-world\n"

  set $file_result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(0, $long_text, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "BUFFER_STREAM_FILE return=%d\n", $file_result
  shell cp /work/stage5/output.wav /work/stage16/buffer-stream-file.bin

  set $buffer = (char *)malloc(1048576)
  set $length = (int *)malloc(4)
  set *$length = 1048576
  set $result = ((int (*)(int, char *, char *, int *, int, int, int, int, int, int, int, int, int))0x1001dc40)(0, $long_text, $buffer, $length, 0, 0, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "BUFFER_STREAM_CALL0 return=%d output_len=%d\n", $result, *$length
  if (*$length > 0) && (*$length <= 1048576)
    dump binary memory /work/stage16/buffer-stream-0.bin $buffer $buffer+*$length
  end

  if $result == 0
    set *$length = 1048576
    set $result = ((int (*)(int, char *, char *, int *, int, int, int, int, int, int, int, int, int))0x1001dc40)(0, $long_text, $buffer, $length, 1, 0, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
    printf "BUFFER_STREAM_POLL1 return=%d output_len=%d\n", $result, *$length
    if (*$length > 0) && (*$length <= 1048576)
      dump binary memory /work/stage16/buffer-stream-1.bin $buffer $buffer+*$length
    end
  end

  if $result == 0
    set *$length = 1048576
    set $result = ((int (*)(int, char *, char *, int *, int, int, int, int, int, int, int, int, int))0x1001dc40)(0, $long_text, $buffer, $length, 1, 0, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
    printf "BUFFER_STREAM_POLL2 return=%d output_len=%d\n", $result, *$length
    if (*$length > 0) && (*$length <= 1048576)
      dump binary memory /work/stage16/buffer-stream-2.bin $buffer $buffer+*$length
    end
  end

  if $result == 0
    set *$length = 1048576
    set $result = ((int (*)(int, char *, char *, int *, int, int, int, int, int, int, int, int, int))0x1001dc40)(0, $long_text, $buffer, $length, 1, 0, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
    printf "BUFFER_STREAM_POLL3 return=%d output_len=%d\n", $result, *$length
    if (*$length > 0) && (*$length <= 1048576)
      dump binary memory /work/stage16/buffer-stream-3.bin $buffer $buffer+*$length
    end
  end

  if $result == 0
    set *$length = 1048576
    set $result = ((int (*)(int, char *, char *, int *, int, int, int, int, int, int, int, int, int))0x1001dc40)(0, $long_text, $buffer, $length, 1, 0, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
    printf "BUFFER_STREAM_POLL4 return=%d output_len=%d\n", $result, *$length
    if (*$length > 0) && (*$length <= 1048576)
      dump binary memory /work/stage16/buffer-stream-4.bin $buffer $buffer+*$length
    end
  end

  if $result == 0
    set *$length = 1048576
    set $result = ((int (*)(int, char *, char *, int *, int, int, int, int, int, int, int, int, int))0x1001dc40)(0, $long_text, $buffer, $length, 2, 0, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
    printf "BUFFER_STREAM_CANCEL return=%d output_len=%d\n", $result, *$length
  end

  if $result == 1
    set $file_after_result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(0, $long_text, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
    printf "BUFFER_STREAM_FILE_AFTER return=%d\n", $file_after_result
    shell cp /work/stage5/output.wav /work/stage16/buffer-stream-file-after.bin
  end
  continue
end

continue
