set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
break *0x1001df44
commands 1
  silent
  set $text = *(unsigned int *)($esp + 8)
  disable 1
  set $name = (char *)0x1009f948
  x/s $name
  set $result = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de50)((char *)$text, $name, 1, -1, -1, -1, -1, -1, -1)
  printf "LIPSYNC_SENTINEL raw_eax=%#x signed=%d\n", $result, $result
  continue
end
commands 2
  silent
  printf "LIPSYNC_SENTINEL_DEFAULT_PATH="
  x/s $ebp-0x400
  disable 2
  continue
end

continue
