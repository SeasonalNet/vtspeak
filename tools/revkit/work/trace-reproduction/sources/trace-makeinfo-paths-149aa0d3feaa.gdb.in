set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  disable 1
  set $result = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de90)((char *)$text, "Z:/tmp/vtspeak-makeinfo-paths-run/base", 1, -1, -1, -1, -1, -1, -1)
  printf "MAKEINFO_PATHPROBE speaker=1 output_prefix=Z:/tmp/vtspeak-makeinfo-paths-run/base raw_eax=%#x signed=%d\n", $result, $result
  continue
end

continue
