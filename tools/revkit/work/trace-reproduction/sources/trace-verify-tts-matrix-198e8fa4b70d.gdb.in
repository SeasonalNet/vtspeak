set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $verify_text = *(unsigned char **)($esp + 8)
  set $i = 0
  while $i < 6
    if $i == 0
      set $dict_index = -2
    end
    if $i == 1
      set $dict_index = -1
    end
    if $i == 2
      set $dict_index = 0
    end
    if $i == 3
      set $dict_index = 1
    end
    if $i == 4
      set $dict_index = 1023
    end
    if $i == 5
      set $dict_index = 1024
    end
    set $j = 0
    while $j < 7
      if $j == 0
        set $text_type = -1
      end
      if $j == 1
        set $text_type = 0
      end
      if $j == 2
        set $text_type = 1
      end
      if $j == 3
        set $text_type = 4
      end
      if $j == 4
        set $text_type = 6
      end
      if $j == 5
        set $text_type = 7
      end
      if $j == 6
        set $text_type = 255
      end
      set $verify_result = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($verify_text, 1, $dict_index, $text_type)
      printf "VERIFY_TTS slot=1 dict_index=%d text_type=%d eax=%d\n", $dict_index, $text_type, $verify_result
      set $j = $j + 1
    end
    set $i = $i + 1
  end
  continue
end

continue
