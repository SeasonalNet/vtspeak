set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

set $source_case = -1
set $source_input = 0
set $source_stack = 0
set $source_return = 0

break *0x1002a567
commands
  silent
  printf "SOURCE_NORM_RESULT case=%d helper_ret=%d input=", $source_case, (int)$eax
  x/s $source_input
  printf "SOURCE_NORM_SCRATCH case=%d bytes=", $source_case
  x/52bx $ebp-0x34
  continue
end

break *0x4016d9
commands
  silent
  if $source_case == 15
    printf "SOURCE_NORM_MATRIX_COMPLETE cases=17\n"
    kill
  else
    set $source_case = $source_case + 1
    set $case_index = 0
    while $case_index < 64
      set {unsigned char}($source_input + $case_index) = 0
      set $case_index = $case_index + 1
    end
    if $source_case == 0
      set {unsigned char}($source_input + 0) = 72
      set {unsigned char}($source_input + 1) = 101
      set {unsigned char}($source_input + 2) = 108
      set {unsigned char}($source_input + 3) = 108
      set {unsigned char}($source_input + 4) = 111
      set {unsigned char}($source_input + 5) = 32
      set {unsigned char}($source_input + 6) = 119
      set {unsigned char}($source_input + 7) = 111
      set {unsigned char}($source_input + 8) = 114
      set {unsigned char}($source_input + 9) = 108
      set {unsigned char}($source_input + 10) = 100
    else
      if $source_case == 1
        set {unsigned char}($source_input + 0) = 72
        set {unsigned char}($source_input + 1) = 101
        set {unsigned char}($source_input + 2) = 108
        set {unsigned char}($source_input + 3) = 108
        set {unsigned char}($source_input + 4) = 111
        set {unsigned char}($source_input + 5) = 44
        set {unsigned char}($source_input + 6) = 32
        set {unsigned char}($source_input + 7) = 119
        set {unsigned char}($source_input + 8) = 111
        set {unsigned char}($source_input + 9) = 114
        set {unsigned char}($source_input + 10) = 108
        set {unsigned char}($source_input + 11) = 100
        set {unsigned char}($source_input + 12) = 33
      else
        if $source_case == 2
          set {unsigned char}($source_input + 0) = 32
          set {unsigned char}($source_input + 1) = 32
          set {unsigned char}($source_input + 2) = 72
          set {unsigned char}($source_input + 3) = 101
          set {unsigned char}($source_input + 4) = 108
          set {unsigned char}($source_input + 5) = 108
          set {unsigned char}($source_input + 6) = 111
          set {unsigned char}($source_input + 7) = 32
          set {unsigned char}($source_input + 8) = 32
          set {unsigned char}($source_input + 9) = 32
          set {unsigned char}($source_input + 10) = 119
          set {unsigned char}($source_input + 11) = 111
          set {unsigned char}($source_input + 12) = 114
          set {unsigned char}($source_input + 13) = 108
          set {unsigned char}($source_input + 14) = 100
          set {unsigned char}($source_input + 15) = 33
          set {unsigned char}($source_input + 16) = 32
          set {unsigned char}($source_input + 17) = 32
        else
          if $source_case == 3
            set {unsigned char}($source_input + 0) = 77
            set {unsigned char}($source_input + 1) = 105
            set {unsigned char}($source_input + 2) = 88
            set {unsigned char}($source_input + 3) = 101
            set {unsigned char}($source_input + 4) = 68
            set {unsigned char}($source_input + 5) = 95
            set {unsigned char}($source_input + 6) = 99
            set {unsigned char}($source_input + 7) = 97
            set {unsigned char}($source_input + 8) = 115
            set {unsigned char}($source_input + 9) = 101
            set {unsigned char}($source_input + 10) = 45
            set {unsigned char}($source_input + 11) = 49
            set {unsigned char}($source_input + 12) = 50
            set {unsigned char}($source_input + 13) = 51
          else
            if $source_case == 4
              set {unsigned char}($source_input + 0) = 99
              set {unsigned char}($source_input + 1) = 97
              set {unsigned char}($source_input + 2) = 110
              set {unsigned char}($source_input + 3) = 39
              set {unsigned char}($source_input + 4) = 116
              set {unsigned char}($source_input + 5) = 32
              set {unsigned char}($source_input + 6) = 115
              set {unsigned char}($source_input + 7) = 116
              set {unsigned char}($source_input + 8) = 111
              set {unsigned char}($source_input + 9) = 112
            else
              if $source_case == 5
                set {unsigned char}($source_input + 0) = 118
                set {unsigned char}($source_input + 1) = 116
                set {unsigned char}($source_input + 2) = 115
                set {unsigned char}($source_input + 3) = 112
                set {unsigned char}($source_input + 4) = 101
                set {unsigned char}($source_input + 5) = 97
                set {unsigned char}($source_input + 6) = 107
                set {unsigned char}($source_input + 7) = 112
                set {unsigned char}($source_input + 8) = 114
                set {unsigned char}($source_input + 9) = 111
                set {unsigned char}($source_input + 10) = 98
                set {unsigned char}($source_input + 11) = 101
              else
                if $source_case == 6
                  set {unsigned char}($source_input + 0) = 99
                  set {unsigned char}($source_input + 1) = 97
                  set {unsigned char}($source_input + 2) = 102
                  set {unsigned char}($source_input + 3) = 195
                  set {unsigned char}($source_input + 4) = 169
                else
                  if $source_case == 7
                    set {unsigned char}($source_input + 0) = 9
                    set {unsigned char}($source_input + 1) = 13
                    set {unsigned char}($source_input + 2) = 10
                    set {unsigned char}($source_input + 3) = 72
                    set {unsigned char}($source_input + 4) = 101
                    set {unsigned char}($source_input + 5) = 108
                    set {unsigned char}($source_input + 6) = 108
                    set {unsigned char}($source_input + 7) = 111
                    set {unsigned char}($source_input + 8) = 13
                    set {unsigned char}($source_input + 9) = 10
                    set {unsigned char}($source_input + 10) = 9
                  else
                    if $source_case == 8
                      set $case_index = 0
                      while $case_index < 49
                        set {unsigned char}($source_input + $case_index) = 65
                        set $case_index = $case_index + 1
                      end
                    else
                      if $source_case == 9
                        set $case_index = 0
                        while $case_index < 50
                          set {unsigned char}($source_input + $case_index) = 65
                          set $case_index = $case_index + 1
                        end
                      else
                        if $source_case == 10
                          set $case_index = 0
                          while $case_index < 51
                            set {unsigned char}($source_input + $case_index) = 65
                            set $case_index = $case_index + 1
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
    if $source_case == 11
      set {unsigned char}($source_input + 0) = 161
      set {unsigned char}($source_input + 1) = 161
    end
    if $source_case == 12
      set {unsigned char}($source_input + 0) = 174
      set {unsigned char}($source_input + 1) = 161
    end
    if $source_case == 13
      set {unsigned char}($source_input + 0) = 253
      set {unsigned char}($source_input + 1) = 254
    end
    if $source_case == 14
      set {unsigned char}($source_input + 0) = 161
      set {unsigned char}($source_input + 1) = 160
    end
    if $source_case == 15
      set {unsigned char}($source_input + 0) = 32
      set {unsigned char}($source_input + 1) = 9
      set {unsigned char}($source_input + 2) = 13
      set {unsigned char}($source_input + 3) = 10
    end
    set {unsigned char}($source_input + 63) = 0
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
  set $source_case = -1
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
