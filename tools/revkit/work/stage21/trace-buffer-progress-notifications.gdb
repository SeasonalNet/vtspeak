set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $long_text = (char *)malloc(52)
  set {char[52]}$long_text = "Hello world. Hello world. Hello world. Hello world."
  set $hwnd = 0
  set $message_id = 0x8005
  set $allocation = (char *)malloc(60032)
  set $buffer = $allocation + 16
  set $length = (int *)malloc(4)
  set $message = (char *)malloc(28)
  set $index = 0
  while $index < 16
    set {unsigned char}($allocation + $index) = 0xa5
    set {unsigned char}($buffer + 60000 + $index) = 0xa5
    set $index = $index + 1
  end
  set *$length = 60000
  printf "PROGRESS_BEGIN hwnd=%#x message=%#x text_bytes=%d text=%s\n", $hwnd, $message_id, strlen($long_text), $long_text
  set $total_bytes = 0
  set $raw = ((unsigned int (*)(void *, unsigned int, char *, unsigned int *, int *, int, int, int, int, int, int, int, int, int))0x1001dda0)((void *)$hwnd, $message_id, $long_text, (unsigned int *)$buffer, $length, 0, 0, 1, -1, -1, -1, -1, -1, -1)
  set $total_bytes = $total_bytes + *$length
  printf "PROGRESS_START raw=%d bytes=%d total=%d\n", $raw, *$length, $total_bytes
  if (*$length > 0) && (*$length <= 60000)
    dump binary memory /work/stage21/progress-notify-start.bin $buffer $buffer+*$length
  end
  set $poll = 1
  while ($raw == 0) && ($poll <= 12)
    set *$length = 60000
    set $raw = ((unsigned int (*)(void *, unsigned int, char *, unsigned int *, int *, int, int, int, int, int, int, int, int, int))0x1001dda0)((void *)$hwnd, $message_id, $long_text, (unsigned int *)$buffer, $length, 1, 0, 1, -1, -1, -1, -1, -1, -1)
    set $total_bytes = $total_bytes + *$length
    printf "PROGRESS_POLL poll=%d raw=%d bytes=%d total=%d\n", $poll, $raw, *$length, $total_bytes
    if (*$length > 0) && (*$length <= 60000)
      if $poll == 1
        dump binary memory /work/stage21/progress-notify-poll-1.bin $buffer $buffer+*$length
      end
      if $poll == 2
        dump binary memory /work/stage21/progress-notify-poll-2.bin $buffer $buffer+*$length
      end
      if $poll == 3
        dump binary memory /work/stage21/progress-notify-poll-3.bin $buffer $buffer+*$length
      end
      if $poll == 4
        dump binary memory /work/stage21/progress-notify-poll-4.bin $buffer $buffer+*$length
      end
    end
    set $poll = $poll + 1
  end
  printf "PROGRESS_FINAL raw=%d polls=%d total_bytes=%d\n", $raw, $poll - 1, $total_bytes
  set $messages = 0
  set $peek = ((int (*)(void *, void *, unsigned int, unsigned int, unsigned int))0x1006d120)($message, (void *)$hwnd, $message_id, $message_id, 1)
  while ($peek != 0) && ($messages < 32)
    set $messages = $messages + 1
    printf "PROGRESS_MESSAGE index=%d hwnd=%#x id=%#x wparam=%#x lparam=%#x\n", $messages, *(unsigned int *)$message, *(unsigned int *)($message + 4), *(unsigned int *)($message + 8), *(unsigned int *)($message + 12)
    set $peek = ((int (*)(void *, void *, unsigned int, unsigned int, unsigned int))0x1006d120)($message, (void *)$hwnd, $message_id, $message_id, 1)
  end
  printf "PROGRESS_MESSAGE_COUNT=%d\n", $messages
  printf "PROGRESS_GUARDS before=%02x,%02x,%02x,%02x after=%02x,%02x,%02x,%02x\n", *(unsigned char *)($buffer - 4), *(unsigned char *)($buffer - 3), *(unsigned char *)($buffer - 2), *(unsigned char *)($buffer - 1), *(unsigned char *)($buffer + 60000), *(unsigned char *)($buffer + 60001), *(unsigned char *)($buffer + 60002), *(unsigned char *)($buffer + 60003)
  continue
end

continue
