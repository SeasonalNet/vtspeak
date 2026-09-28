set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $verify_text = *(unsigned char **)($esp + 8)
  set $text_type = 0
  set $successes = 0
  set $failures = 0
  while $text_type < 256
    set $verify_result = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($verify_text, 1, -1, $text_type)
    if $verify_result == 1
      set $successes = $successes + 1
    else
      set $failures = $failures + 1
      printf "VERIFY_TTS type_failure type=%d eax=%d\n", $text_type, $verify_result
    end
    set $text_type = $text_type + 1
  end
  printf "VERIFY_TTS text_types=0..255 dict_index=-1 success=%d failure=%d\n", $successes, $failures
  continue
end

continue
