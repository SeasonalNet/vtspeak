set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $name = (char *)malloc(512)
  set $path = (char *)malloc(512)
  set $result = ((int (*)(int, char *, char *))0x1002aa60)(6, $name, $path)
  printf "SPEAKERSINFO_OOB selector=6 result=%d name=%s path=%s\n", $result, $name, $path
  kill
  quit
end

continue
