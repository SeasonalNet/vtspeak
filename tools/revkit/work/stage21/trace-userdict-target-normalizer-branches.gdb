set pagination off
set confirm off
set debuginfo enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $target_input = (unsigned int)malloc(256)
  set $target_case = 0
  while $target_case < 7
    set $target_offset = 0
    while $target_offset < 256
      set {unsigned char}($target_input + $target_offset) = 0
      set $target_offset = $target_offset + 1
    end
    if $target_case == 0
      set {unsigned char}($target_input + 0) = 60
      set {unsigned char}($target_input + 1) = 65
      set {unsigned char}($target_input + 2) = 66
      set {unsigned char}($target_input + 3) = 62
    else
      if $target_case == 1
        set {unsigned char}($target_input + 0) = 60
        set {unsigned char}($target_input + 1) = 65
        set {unsigned char}($target_input + 2) = 66
      else
        if $target_case == 2
          set {unsigned char}($target_input + 0) = 60
          set {unsigned char}($target_input + 1) = 65
          set {unsigned char}($target_input + 2) = 32
          set {unsigned char}($target_input + 3) = 66
          set {unsigned char}($target_input + 4) = 62
        else
          if $target_case == 3
            set $target_length = 30
            while $target_length <= 31
              set $target_offset = 0
              while $target_offset < 256
                set {unsigned char}($target_input + $target_offset) = 0
                set $target_offset = $target_offset + 1
              end
              set $target_offset = 0
              while $target_offset < $target_length
                set {unsigned char}($target_input + $target_offset) = 65
                set $target_offset = $target_offset + 1
              end
              set {unsigned char}($target_input + $target_length) = 32
              set {unsigned char}($target_input + $target_length + 1) = 66
              set $target_result = ((int (*)(char *))0x1002a570)($target_input)
              printf "TARGET_NORM_BRANCH name=segment_len_%d result=%d bytes=", $target_length, $target_result
              x/36bx $target_input
              set $target_length = $target_length + 1
            end
          else
            if $target_case == 4
              set {unsigned char}($target_input + 0) = 65
              set {unsigned char}($target_input + 1) = 0xa1
              set {unsigned char}($target_input + 2) = 0xa1
            else
              if $target_case == 5
                set {unsigned char}($target_input + 0) = 65
                set {unsigned char}($target_input + 1) = 0xae
                set {unsigned char}($target_input + 2) = 0xa1
              else
                set {unsigned char}($target_input + 0) = 65
                set {unsigned char}($target_input + 1) = 0xfd
                set {unsigned char}($target_input + 2) = 0xfe
              end
            end
          end
        end
      end
    end
    if $target_case < 3
      set $target_result = ((int (*)(char *))0x1002a570)($target_input)
      printf "TARGET_NORM_BRANCH case=%d result=%d bytes=", $target_case, $target_result
      x/12bx $target_input
    else
      if $target_case > 3
        set $target_result = ((int (*)(char *))0x1002a570)($target_input)
        printf "TARGET_NORM_BRANCH case=%d result=%d bytes=", $target_case, $target_result
        x/12bx $target_input
      end
    end
    set $target_case = $target_case + 1
  end
  printf "TARGET_NORM_BRANCHES_COMPLETE cases=7 plus_segment_boundaries=2\n"
  kill
end

continue
