set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $result = ((int (*)(int, char *, char *))0x1002aa60)(0, (char *)0, (char *)0)
  printf "SPEAKERSINFO_NULL_POINTER case=both-null result=%d\n", $result
  kill
  quit
end

continue
