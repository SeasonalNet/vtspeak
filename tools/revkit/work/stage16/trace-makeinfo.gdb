set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  disable 1
  set $result = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de90)((char *)$text, "Z:/work/stage16/lead6-makeinfo-probe", 1, -1, -1, -1, -1, -1, -1)
  printf "MAKEINFO speaker=1 output_prefix=Z:/work/stage16/lead6-makeinfo-probe options=-1 raw_eax=%#x signed=%d\n", $result, $result
  continue
end

continue
