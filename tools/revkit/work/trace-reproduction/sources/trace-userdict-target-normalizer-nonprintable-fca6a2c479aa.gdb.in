set pagination off
set confirm off
set debuginfo enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $target_input = (unsigned int)malloc(16)
  set $target_byte = 1
  while $target_byte <= 255
    if $target_byte < 32 || $target_byte >= 127
      set $target_offset = 0
      while $target_offset < 16
        set {unsigned char}($target_input + $target_offset) = 0
        set $target_offset = $target_offset + 1
      end
      set {unsigned char}($target_input) = $target_byte
      set $target_result = ((int (*)(char *))0x1002a570)($target_input)
      printf "TARGET_NORM_BYTE byte=%#x result=%d buffer=", $target_byte, $target_result
      x/4bx $target_input
    end
    set $target_byte = $target_byte + 1
  end
  printf "TARGET_NORM_NONPRINTABLE_COMPLETE controls=31 del=1 high=129 total=161\n"
  kill
end

continue
