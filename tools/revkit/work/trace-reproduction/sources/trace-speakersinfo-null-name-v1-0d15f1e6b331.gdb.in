set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $path = (char *)malloc(512)
  set $result = ((int (*)(int, char *, char *))0x1002aa60)(0, (char *)0, $path)
  printf "SPEAKERSINFO_NULL_POINTER case=null-name result=%d path=%s\n", $result, $path
  kill
  quit
end

continue
