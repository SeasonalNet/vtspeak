set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $verify_text = *(unsigned char **)($esp + 8)
  set $verify_valid = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($verify_text, 1, 0, 0)
  printf "VERIFY_TTS valid_text slot=1 args_3_4=0,0 eax=%d\n", $verify_valid
  set $verify_invalid_slot = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($verify_text, -1, 0, 0)
  printf "VERIFY_TTS valid_text slot=-1 args_3_4=0,0 eax=%d\n", $verify_invalid_slot
  set $verify_unloaded_slot = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($verify_text, 0, 0, 0)
  printf "VERIFY_TTS valid_text slot=0 args_3_4=0,0 eax=%d\n", $verify_unloaded_slot
  set $verify_null = ((int (*)(unsigned char *, int, int, int))0x1001ded0)((unsigned char *)0, 1, 0, 0)
  printf "VERIFY_TTS null_text slot=1 args_3_4=0,0 eax=%d\n", $verify_null
  continue
end

continue
