set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  disable 1
  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, "Z:/work/stage16/flag10-preprocess-isolated-20260926", 10, 1, -1, -1, -1, -1, 0, 0)
  printf "PREPROCESS_FLAG_ISOLATED flag=10 raw_eax=%#x low_ax=%d\n", $raw, ((short)$raw)
  continue
end

continue
