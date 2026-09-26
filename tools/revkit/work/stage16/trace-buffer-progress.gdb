set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  disable 1
  set $allocation = (char *)malloc(60032)
  set $buffer = $allocation + 16
  set $length = (int *)malloc(4)
  set $i = 0
  while $i < 16
    set {unsigned char}($allocation + $i) = 0xa5
    set {unsigned char}($buffer + 60000 + $i) = 0xa5
    set $i = $i + 1
  end
  set *$length = 60000
  printf "BUFFER_PROGRESS_ALLOC output=%#x length_ptr=%#x capacity=60000 hwnd=0 message=0 thread=0 speaker=1 options=-1\n", $buffer, $length
  set $raw = ((unsigned int (*)(void *, unsigned int, char *, unsigned int *, int *, int, int, int, int, int, int, int, int, int))0x1001dda0)(0, 0, (char *)$text, (unsigned int *)$buffer, $length, 0, 0, 1, -1, -1, -1, -1, -1, -1)
  printf "BUFFER_PROGRESS_START raw_eax=%#x signed=%d output_len=%d\n", $raw, $raw, *$length
  if (*$length > 0) && (*$length <= 60000)
    dump binary memory /work/stage16/buffer-progress-start.bin $buffer $buffer+*$length
  end
  set $poll = 1
  while ($raw == 0) && ($poll <= 4)
    set *$length = 60000
    set $raw = ((unsigned int (*)(void *, unsigned int, char *, unsigned int *, int *, int, int, int, int, int, int, int, int, int))0x1001dda0)(0, 0, (char *)$text, (unsigned int *)$buffer, $length, 1, 0, 1, -1, -1, -1, -1, -1, -1)
    printf "BUFFER_PROGRESS_POLL poll=%d raw_eax=%#x signed=%d output_len=%d\n", $poll, $raw, $raw, *$length
    if (*$length > 0) && (*$length <= 60000)
      dump binary memory /work/stage16/buffer-progress-poll.bin $buffer $buffer+*$length
    end
    set $poll = $poll + 1
  end
  printf "BUFFER_PROGRESS_GUARDS before=%02x,%02x,%02x,%02x after=%02x,%02x,%02x,%02x\n", *(unsigned char *)($buffer - 4), *(unsigned char *)($buffer - 3), *(unsigned char *)($buffer - 2), *(unsigned char *)($buffer - 1), *(unsigned char *)($buffer + 60000), *(unsigned char *)($buffer + 60001), *(unsigned char *)($buffer + 60002), *(unsigned char *)($buffer + 60003)
  continue
end

continue
