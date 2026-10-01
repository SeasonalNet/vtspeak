set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $target_byte = 32
  set $target_input = (unsigned int)malloc(16)
  while $target_byte <= 126
    set $target_offset = 0
    while $target_offset < 16
      set {unsigned char}($target_input + $target_offset) = 0
      set $target_offset = $target_offset + 1
    end
    set {unsigned char}$target_input = $target_byte
    set $target_result = ((int (*)(char *))0x1002a570)($target_input)
    printf "TARGET_NORM_ASCII byte=%#x result=%d input=", $target_byte, $target_result
    x/s $target_input
    set $target_byte = $target_byte + 1
  end
  printf "TARGET_NORM_ASCII_COMPLETE first=0x20 last=0x7e count=95\n"
  kill
end

continue
