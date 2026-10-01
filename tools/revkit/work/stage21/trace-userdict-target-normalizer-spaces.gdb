set pagination off
set confirm off
set debuginfo enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $target_input = (unsigned int)malloc(256)
  set $target_spaces = 1
  while $target_spaces <= 35
    set $target_offset = 0
    while $target_offset < 256
      set {unsigned char}($target_input + $target_offset) = 0
      set $target_offset = $target_offset + 1
    end
    set {unsigned char}($target_input) = 65
    set $target_offset = 0
    while $target_offset < $target_spaces
      set {unsigned char}($target_input + $target_offset + 1) = 32
      set $target_offset = $target_offset + 1
    end
    set {unsigned char}($target_input + $target_spaces + 1) = 65
    set $target_result = ((int (*)(char *))0x1002a570)($target_input)
    printf "TARGET_NORM_SPACES spaces=%d result=%d input=", $target_spaces, $target_result
    x/24bx $target_input
    set $target_spaces = $target_spaces + 1
  end
  printf "TARGET_NORM_SPACES_COMPLETE first=1 last=35 count=35\n"
  kill
end

continue
