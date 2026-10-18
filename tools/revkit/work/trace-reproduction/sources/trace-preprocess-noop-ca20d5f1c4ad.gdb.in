set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, unsigned int, int, int))0x1001dd60)(0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
  printf "PREPROCESS_ZERO_FLAG all_other_arguments_null raw_eax=%#x signed=%d\n", $raw, $raw
  continue
end

continue
