set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

disassemble /r 0x100214d0, 0x10021540

set $progress_message_count = 0
set $progress_poll_count = 0
set $progress_limit = 12
set $progress_total_bytes = 0
set $progress_buffer = 0
set $progress_length = 0
set $progress_text = 0
set $progress_stack = (unsigned int)malloc(64)
set $progress_stack = $progress_stack & 0xfffffff0
set $progress_return = 0

break *0x10021513
commands
  silent
  set $progress_context = $eax
  set $progress_context_wparam = *(unsigned int *)$eax
  printf "PROGRESS_SOURCE context=%#x first_dword=%#x\n", $progress_context, $progress_context_wparam
  continue
end
disable 1

break *0x10021519
commands
  silent
  set $progress_message_count = $progress_message_count + 1
  printf "PROGRESS_POST index=%d hwnd=%#x message=%#x wparam=%#x lparam=%#x context=%#x context_first=%#x source_matches=%d\n", $progress_message_count, *(unsigned int *)($ebp + 8), *(unsigned int *)($ebp + 12), *(unsigned int *)$esp, *(unsigned int *)($esp + 4), $progress_context, $progress_context_wparam, (*(unsigned int *)$esp == $progress_context_wparam)
  continue
end
disable 2

break *0x4016d9
commands
  silent
  set $progress_total_bytes = $progress_total_bytes + *$progress_length
  printf "PROGRESS_RETURN poll=%d raw=%d output_bytes=%d notifications=%d\n", $progress_poll_count, (short)$eax, *$progress_length, $progress_message_count
  if (short)$eax == 0
    if $progress_poll_count >= $progress_limit
      printf "PROGRESS_ABORT poll_limit=%d\n", $progress_limit
      kill
    else
      set $progress_poll_count = $progress_poll_count + 1
      set *$progress_length = 60000
      set {unsigned int}$progress_stack = $progress_return
      set {unsigned int}($progress_stack + 4) = 0
      set {unsigned int}($progress_stack + 8) = 0x8005
      set {unsigned int}($progress_stack + 12) = $progress_text
      set {unsigned int}($progress_stack + 16) = $progress_buffer
      set {unsigned int}($progress_stack + 20) = $progress_length
      set {int}($progress_stack + 24) = 1
      set {int}($progress_stack + 28) = 0
      set {int}($progress_stack + 32) = 1
      set {int}($progress_stack + 36) = -1
      set {int}($progress_stack + 40) = -1
      set {int}($progress_stack + 44) = -1
      set {int}($progress_stack + 48) = -1
      set {int}($progress_stack + 52) = -1
      set {int}($progress_stack + 56) = -1
      set $esp = $progress_stack
      set $eip = 0x1001dda0
      continue
    end
  else
    printf "PROGRESS_COMPLETE raw=%d polls=%d notifications=%d total_bytes=%d\n", (short)$eax, $progress_poll_count, $progress_message_count, $progress_total_bytes
    kill
  end
end
disable 3

break *0x1001da50
commands
  silent
disable 4
  set $progress_return = *(unsigned int *)$esp
  set $source_text = *(unsigned int *)($esp + 8)
  set $source_length = 0
  while *(unsigned char *)($source_text + $source_length) != 0
    set $source_length = $source_length + 1
  end
  set $progress_text = (unsigned int)malloc(64)
  set $copy_position = 0
  set $repeat_index = 0
  while $repeat_index < 4
    set $copy_index = 0
    while $copy_index < $source_length
      set {unsigned char}($progress_text + $copy_position) = *(unsigned char *)($source_text + $copy_index)
      set $copy_position = $copy_position + 1
      set $copy_index = $copy_index + 1
    end
    if $repeat_index < 3
      set {unsigned char}($progress_text + $copy_position) = 32
      set $copy_position = $copy_position + 1
    end
    set $repeat_index = $repeat_index + 1
  end
  set {unsigned char}($progress_text + $copy_position) = 0
  set $allocation = (char *)malloc(60032)
  set $progress_buffer = (unsigned int)($allocation + 16)
  set $progress_length = (int *)malloc(4)
  set $guard_index = 0
  while $guard_index < 16
    set {unsigned char}($allocation + $guard_index) = 0xa5
    set {unsigned char}($progress_buffer + 60000 + $guard_index) = 0xa5
    set $guard_index = $guard_index + 1
  end
  set *$progress_length = 60000
  set $progress_poll_count = 0
  set $progress_message_count = 0
  set {unsigned int}$progress_stack = $progress_return
  set {unsigned int}($progress_stack + 4) = 0
  set {unsigned int}($progress_stack + 8) = 0x8005
  set {unsigned int}($progress_stack + 12) = $progress_text
  set {unsigned int}($progress_stack + 16) = $progress_buffer
  set {unsigned int}($progress_stack + 20) = $progress_length
  set {int}($progress_stack + 24) = 0
  set {int}($progress_stack + 28) = 0
  set {int}($progress_stack + 32) = 1
  set {int}($progress_stack + 36) = -1
  set {int}($progress_stack + 40) = -1
  set {int}($progress_stack + 44) = -1
  set {int}($progress_stack + 48) = -1
  set {int}($progress_stack + 52) = -1
  set {int}($progress_stack + 56) = -1
  printf "PROGRESS_TARGET_BEGIN hwnd=0 message=0x8005 source_bytes=%d text_bytes=%d text=%s stack=%#x\n", $source_length, $copy_position, $progress_text, $progress_stack
  enable 1
  enable 2
  enable 3
  enable 5
  set $esp = $progress_stack
  set $eip = 0x1001dda0
  continue
end

break *0x10021395
commands
  silent
  disable 5
  set $progress_context = $eax
  set $progress_context_first = *(unsigned int *)$progress_context
  set $progress_context_prev = $progress_context_first
  printf "PROGRESS_CONTEXT_INIT context=%#x first_dword=%#x\n", $progress_context, $progress_context_first
  watch -l *(unsigned int *)$progress_context
  commands 7
    silent
    printf "PROGRESS_CONTEXT_WRITE eip=%#x old=%#x new=%#x\n", $eip, $progress_context_prev, *(unsigned int *)$progress_context
    set $progress_context_prev = *(unsigned int *)$progress_context
    continue
  end
  continue
end
disable 5

break *0x10021439
commands
  silent
  set $progress_helper_string = *(unsigned int *)($esp + 8)
  set $progress_helper_string_length = 0
  while *(unsigned char *)($progress_helper_string + $progress_helper_string_length) != 0
    set $progress_helper_string_length = $progress_helper_string_length + 1
  end
  printf "PROGRESS_HELPER_STRING pointer=%#x bytes=%d text=%s\n", $progress_helper_string, $progress_helper_string_length, $progress_helper_string
  continue
end

continue
