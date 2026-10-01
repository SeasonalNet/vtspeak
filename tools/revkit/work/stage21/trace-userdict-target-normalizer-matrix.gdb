set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

set $target_case = -1
set $target_input = 0
set $target_stack = 0
set $target_return = 0

break *0x4016d9
commands
  silent
  printf "TARGET_NORM_RESULT case=%d return=%d input=", $target_case, (int)$eax
  x/s $target_input
  if $target_case == 14
    printf "TARGET_NORM_MATRIX_COMPLETE cases=16\n"
    kill
  else
    set $target_case = $target_case + 1
    set $target_index = 0
    while $target_index < 64
      set {unsigned char}($target_input + $target_index) = 0
      set $target_index = $target_index + 1
    end
    if $target_case == 0
      set {unsigned char}($target_input + 0) = 104
      set {unsigned char}($target_input + 1) = 101
      set {unsigned char}($target_input + 2) = 108
      set {unsigned char}($target_input + 3) = 108
      set {unsigned char}($target_input + 4) = 111
    else
      if $target_case == 1
        set {unsigned char}($target_input + 0) = 72
        set {unsigned char}($target_input + 1) = 69
        set {unsigned char}($target_input + 2) = 76
        set {unsigned char}($target_input + 3) = 76
        set {unsigned char}($target_input + 4) = 79
      else
        if $target_case == 2
          set {unsigned char}($target_input + 0) = 72
          set {unsigned char}($target_input + 1) = 101
          set {unsigned char}($target_input + 2) = 108
          set {unsigned char}($target_input + 3) = 108
          set {unsigned char}($target_input + 4) = 111
          set {unsigned char}($target_input + 5) = 32
          set {unsigned char}($target_input + 6) = 87
          set {unsigned char}($target_input + 7) = 111
          set {unsigned char}($target_input + 8) = 114
          set {unsigned char}($target_input + 9) = 108
          set {unsigned char}($target_input + 10) = 100
        else
          if $target_case == 3
            set {unsigned char}($target_input + 0) = 32
            set {unsigned char}($target_input + 1) = 32
            set {unsigned char}($target_input + 2) = 104
            set {unsigned char}($target_input + 3) = 101
            set {unsigned char}($target_input + 4) = 108
            set {unsigned char}($target_input + 5) = 108
            set {unsigned char}($target_input + 6) = 111
            set {unsigned char}($target_input + 7) = 32
            set {unsigned char}($target_input + 8) = 32
            set {unsigned char}($target_input + 9) = 32
          else
            if $target_case == 4
              set {unsigned char}($target_input + 0) = 104
              set {unsigned char}($target_input + 1) = 101
              set {unsigned char}($target_input + 2) = 108
              set {unsigned char}($target_input + 3) = 108
              set {unsigned char}($target_input + 4) = 111
              set {unsigned char}($target_input + 5) = 44
              set {unsigned char}($target_input + 6) = 32
              set {unsigned char}($target_input + 7) = 119
              set {unsigned char}($target_input + 8) = 111
              set {unsigned char}($target_input + 9) = 114
              set {unsigned char}($target_input + 10) = 108
              set {unsigned char}($target_input + 11) = 100
              set {unsigned char}($target_input + 12) = 33
            else
              if $target_case == 5
                set {unsigned char}($target_input + 0) = 104
                set {unsigned char}($target_input + 1) = 101
                set {unsigned char}($target_input + 2) = 108
                set {unsigned char}($target_input + 3) = 108
                set {unsigned char}($target_input + 4) = 111
                set {unsigned char}($target_input + 5) = 45
                set {unsigned char}($target_input + 6) = 119
                set {unsigned char}($target_input + 7) = 111
                set {unsigned char}($target_input + 8) = 114
                set {unsigned char}($target_input + 9) = 108
                set {unsigned char}($target_input + 10) = 100
              else
                if $target_case == 6
                  set {unsigned char}($target_input + 0) = 91
                  set {unsigned char}($target_input + 1) = 83
                  set {unsigned char}($target_input + 2) = 75
                  set {unsigned char}($target_input + 3) = 73
                  set {unsigned char}($target_input + 4) = 80
                  set {unsigned char}($target_input + 5) = 93
                else
                  if $target_case == 7
                    set {unsigned char}($target_input + 0) = 91
                    set {unsigned char}($target_input + 1) = 67
                    set {unsigned char}($target_input + 2) = 73
                    set {unsigned char}($target_input + 3) = 93
                  else
                    if $target_case == 8
                      set {unsigned char}($target_input + 0) = 104
                      set {unsigned char}($target_input + 1) = 101
                      set {unsigned char}($target_input + 2) = 108
                      set {unsigned char}($target_input + 3) = 108
                      set {unsigned char}($target_input + 4) = 111
                      set {unsigned char}($target_input + 5) = 91
                      set {unsigned char}($target_input + 6) = 67
                      set {unsigned char}($target_input + 7) = 73
                      set {unsigned char}($target_input + 8) = 93
                    else
                      if $target_case == 9
                        set {unsigned char}($target_input + 0) = 91
                        set {unsigned char}($target_input + 1) = 79
                        set {unsigned char}($target_input + 2) = 84
                        set {unsigned char}($target_input + 3) = 72
                        set {unsigned char}($target_input + 4) = 69
                        set {unsigned char}($target_input + 5) = 82
                        set {unsigned char}($target_input + 6) = 93
                      else
                        if $target_case == 10
                          set {unsigned char}($target_input + 0) = 60
                          set {unsigned char}($target_input + 1) = 118
                          set {unsigned char}($target_input + 2) = 116
                          set {unsigned char}($target_input + 3) = 109
                          set {unsigned char}($target_input + 4) = 108
                          set {unsigned char}($target_input + 5) = 95
                          set {unsigned char}($target_input + 6) = 115
                          set {unsigned char}($target_input + 7) = 117
                          set {unsigned char}($target_input + 8) = 98
                          set {unsigned char}($target_input + 9) = 62
                        else
                          if $target_case == 11
                            set {unsigned char}($target_input + 0) = 60
                          else
                            if $target_case == 12
                              set {unsigned char}($target_input + 0) = 35
                            else
                              if $target_case == 13
                                set $target_index = 0
                                while $target_index < 64
                                  set {unsigned char}($target_input + $target_index) = 65
                                  set $target_index = $target_index + 1
                                end
                              else
                                if $target_case == 14
                                  set $target_index = 0
                                  while $target_index < 65
                                    set {unsigned char}($target_input + $target_index) = 65
                                    set $target_index = $target_index + 1
                                  end
                                end
                              end
                            end
                          end
                        end
                      end
                    end
                  end
                end
              end
            end
          end
        end
      end
    end
    set {unsigned char}($target_input + 127) = 0
    set {unsigned int}$target_stack = $target_return
    set {unsigned int}($target_stack + 4) = $target_input
    set $esp = $target_stack
    set $eip = 0x1002a570
    continue
  end
end

break *0x1001da50
commands
  silent
  disable 2
  set $target_return = *(unsigned int *)$esp
  set $target_input = (unsigned int)malloc(256)
  set $target_stack = $esp
  set $target_case = -1
  set $target_index = 0
  while $target_index < 128
    set {unsigned char}($target_input + $target_index) = 0
    set $target_index = $target_index + 1
  end
  set {unsigned int}$target_stack = $target_return
  set {unsigned int}($target_stack + 4) = $target_input
  set $esp = $target_stack
  set $eip = 0x1002a570
  continue
end

continue
