set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $trace_target = 0

break *0x10025dc0 if $trace_target == 1
commands
  silent
  printf "PREPROCESS_FILE_HELPER_ENTRY arg1 path=%#x arg2 mode=%#x\n", *(unsigned int *)($esp + 4), *(unsigned int *)($esp + 8)
  x/s *(char **)($esp + 4)
  x/s *(char **)($esp + 8)
  set $trace_target = 0
  disable 1
  continue
end

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  disable 2
  set $trace_target = 1
  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, "Z:/work/stage16/preprocess-path-inspect-20260926", 5, 1, -1, -1, -1, -1, 0, 0)
  printf "PREPROCESS_PATH_INSPECT raw_eax=%#x low_ax=%d\n", $raw, ((short)$raw)
  continue
end

continue
