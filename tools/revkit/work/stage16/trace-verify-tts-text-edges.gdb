set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $empty_text = ((char *(*)(unsigned int))0x1001d9c0)(1)
  set {char}($empty_text) = 0
  set $empty_result = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($empty_text, 1, -1, 0)
  printf "VERIFY_TTS empty_nonnull slot=1 dict_index=-1 text_type=0 eax=%d\n", $empty_result
  set $short_markup = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($short_markup) = 60
  set {char}($short_markup + 1) = 0
  set $markup_result = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($short_markup, 1, -1, 0)
  printf "VERIFY_TTS text_less_than_only slot=1 dict_index=-1 text_type=0 eax=%d\n", $markup_result
  continue
end

continue
