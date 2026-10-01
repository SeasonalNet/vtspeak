set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $name = (char *)malloc(512)
  set $result = ((int (*)(int, char *, char *))0x1002aa60)(0, $name, (char *)0)
  printf "SPEAKERSINFO_NULL_POINTER case=null-path result=%d name=%s\n", $result, $name
  kill
  quit
end

continue
