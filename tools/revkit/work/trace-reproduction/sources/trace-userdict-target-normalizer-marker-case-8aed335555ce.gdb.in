set pagination off
set confirm off
set debuginfo enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $target_input = (unsigned int)malloc(256)
  set $target_mask = 0
  while $target_mask < 16
    set $target_offset = 0
    while $target_offset < 256
      set {unsigned char}($target_input + $target_offset) = 0
      set $target_offset = $target_offset + 1
    end
    set {unsigned char}($target_input + 0) = 91
    set $target_index = 0
    while $target_index < 4
      if $target_index == 0
        set $target_char = 83
      else
        if $target_index == 1
          set $target_char = 75
        else
          if $target_index == 2
            set $target_char = 73
          else
            set $target_char = 80
          end
        end
      end
      if ($target_mask & (1 << $target_index)) != 0
        set $target_char = $target_char + 32
      end
      set {unsigned char}($target_input + $target_index + 1) = $target_char
      set $target_index = $target_index + 1
    end
    set {unsigned char}($target_input + 5) = 93
    set $target_result = ((int (*)(char *))0x1002a570)($target_input)
    printf "TARGET_NORM_SKIP_CASE mask=%d result=%d bytes=", $target_mask, $target_result
    x/7bx $target_input
    set $target_mask = $target_mask + 1
  end
  set $target_mask = 0
  while $target_mask < 4
    set $target_offset = 0
    while $target_offset < 256
      set {unsigned char}($target_input + $target_offset) = 0
      set $target_offset = $target_offset + 1
    end
    set {unsigned char}($target_input + 0) = 65
    set {unsigned char}($target_input + 1) = 91
    set $target_char = 67
    if ($target_mask & 1) != 0
      set $target_char = 99
    end
    set {unsigned char}($target_input + 2) = $target_char
    set $target_char = 73
    if ($target_mask & 2) != 0
      set $target_char = 105
    end
    set {unsigned char}($target_input + 3) = $target_char
    set {unsigned char}($target_input + 4) = 93
    set $target_result = ((int (*)(char *))0x1002a570)($target_input)
    printf "TARGET_NORM_CI_CASE mask=%d result=%d bytes=", $target_mask, $target_result
    x/7bx $target_input
    set $target_mask = $target_mask + 1
  end
  printf "TARGET_NORM_MARKER_CASE_COMPLETE skip=16 ci=4 total=20\n"
  kill
end

continue
