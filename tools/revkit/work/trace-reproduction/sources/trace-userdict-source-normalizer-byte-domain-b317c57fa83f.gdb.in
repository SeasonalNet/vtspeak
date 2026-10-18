set pagination off
set confirm off
set debuginfo enabled off
handle SIGSEGV nostop noprint pass

set $source_byte = 0
set $source_input = 0
set $source_stack = 0
set $source_return = 0

break *0x1002a567
commands
  silent
  printf "SOURCE_NORM_BYTE byte=%#x helper_ret=%d scratch=", $source_byte, (int)$eax
  x/4bx $ebp-0x34
  continue
end

break *0x4016d9
commands
  silent
  if $source_byte == 255
    printf "SOURCE_NORM_BYTE_DOMAIN_COMPLETE first=0 last=255 count=256\n"
    kill
  else
    set $source_byte = $source_byte + 1
    set $case_index = 0
    while $case_index < 64
      set {unsigned char}($source_input + $case_index) = 0
      set $case_index = $case_index + 1
    end
    set {unsigned char}$source_input = $source_byte
    set {unsigned int}$source_stack = $source_return
    set {unsigned int}($source_stack + 4) = $source_input
    set $esp = $source_stack
    set $eip = 0x1002a550
    continue
  end
end

break *0x1001da50
commands
  silent
  disable 3
  set $source_return = *(unsigned int *)$esp
  set $source_input = (unsigned int)malloc(128)
  set $source_stack = $esp
  set $source_byte = 0
  set $case_index = 0
  while $case_index < 64
    set {unsigned char}($source_input + $case_index) = 0
    set $case_index = $case_index + 1
  end
  set {unsigned int}$source_stack = $source_return
  set {unsigned int}($source_stack + 4) = $source_input
  set $esp = $source_stack
  set $eip = 0x1002a550
  continue
end

continue
