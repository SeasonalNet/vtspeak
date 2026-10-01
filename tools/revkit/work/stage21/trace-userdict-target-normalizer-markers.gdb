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
  while $target_case < 8
    set $target_offset = 0
    while $target_offset < 256
      set {unsigned char}($target_input + $target_offset) = 0
      set $target_offset = $target_offset + 1
    end
    if $target_case == 0
      set {unsigned char}($target_input + 0) = 91
      set {unsigned char}($target_input + 1) = 83
      set {unsigned char}($target_input + 2) = 75
      set {unsigned char}($target_input + 3) = 73
      set {unsigned char}($target_input + 4) = 80
      set {unsigned char}($target_input + 5) = 93
      set {unsigned char}($target_input + 6) = 120
    else
      if $target_case == 1
        set {unsigned char}($target_input + 0) = 91
        set {unsigned char}($target_input + 1) = 83
        set {unsigned char}($target_input + 2) = 75
        set {unsigned char}($target_input + 3) = 73
        set {unsigned char}($target_input + 4) = 80
        set {unsigned char}($target_input + 5) = 93
        set {unsigned char}($target_input + 6) = 32
      else
        if $target_case == 2
          set {unsigned char}($target_input + 0) = 91
          set {unsigned char}($target_input + 1) = 115
          set {unsigned char}($target_input + 2) = 107
          set {unsigned char}($target_input + 3) = 105
          set {unsigned char}($target_input + 4) = 112
          set {unsigned char}($target_input + 5) = 93
        else
          if $target_case == 3
            set {unsigned char}($target_input + 0) = 65
            set {unsigned char}($target_input + 1) = 91
            set {unsigned char}($target_input + 2) = 67
            set {unsigned char}($target_input + 3) = 73
            set {unsigned char}($target_input + 4) = 93
          else
            if $target_case == 4
              set {unsigned char}($target_input + 0) = 65
              set {unsigned char}($target_input + 1) = 91
              set {unsigned char}($target_input + 2) = 99
              set {unsigned char}($target_input + 3) = 105
              set {unsigned char}($target_input + 4) = 93
            else
              if $target_case == 5
                set {unsigned char}($target_input + 0) = 65
                set {unsigned char}($target_input + 1) = 32
                set {unsigned char}($target_input + 2) = 91
                set {unsigned char}($target_input + 3) = 67
                set {unsigned char}($target_input + 4) = 73
                set {unsigned char}($target_input + 5) = 93
              else
                if $target_case == 6
                  set {unsigned char}($target_input + 0) = 65
                  set {unsigned char}($target_input + 1) = 91
                  set {unsigned char}($target_input + 2) = 79
                  set {unsigned char}($target_input + 3) = 84
                  set {unsigned char}($target_input + 4) = 72
                  set {unsigned char}($target_input + 5) = 69
                  set {unsigned char}($target_input + 6) = 82
                  set {unsigned char}($target_input + 7) = 93
                else
                  set {unsigned char}($target_input + 0) = 65
                  set {unsigned char}($target_input + 1) = 91
                  set {unsigned char}($target_input + 2) = 67
                  set {unsigned char}($target_input + 3) = 73
                  set {unsigned char}($target_input + 4) = 93
                  set {unsigned char}($target_input + 5) = 66
                end
              end
            end
          end
        end
      end
    end
    set $target_result = ((int (*)(char *))0x1002a570)($target_input)
    printf "TARGET_NORM_MARKER case=%d result=%d bytes=", $target_case, $target_result
    x/12bx $target_input
    set $target_case = $target_case + 1
  end
  printf "TARGET_NORM_MARKERS_COMPLETE cases=8\n"
  kill
end

continue
