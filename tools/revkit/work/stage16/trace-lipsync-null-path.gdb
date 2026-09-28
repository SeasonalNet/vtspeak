set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  disable 1
  set $result = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de50)((char *)$text, 0, 1, -1, -1, -1, -1, -1, -1)
  printf "LIPSYNC_NULL_PATH raw_eax=%#x signed=%d\n", $result, $result
  continue
end

continue
