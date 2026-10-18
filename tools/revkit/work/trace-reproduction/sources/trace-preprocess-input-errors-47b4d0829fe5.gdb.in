set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $null_result = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)0, "Z:/work/stage16/preprocess-null-text-20260926", 5, 1, -1, -1, -1, -1, 0, 0)
  printf "PREPROCESS_NULL_TEXT raw_eax=%#x low_ax=%d\n", $null_result, ((short)$null_result)
  set $empty_result = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)"", "Z:/work/stage16/preprocess-empty-text-20260926", 5, 1, -1, -1, -1, -1, 0, 0)
  printf "PREPROCESS_EMPTY_TEXT raw_eax=%#x low_ax=%d\n", $empty_result, ((short)$empty_result)
  continue
end

continue
