set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1

  set $result = ((int (*)(int))0x1002a5b0)(-2147483648)
  printf "USERDICT_LIMIT selector=INT_MIN result=%d\n", $result
  set $result = ((int (*)(int))0x1002a5b0)(-1)
  printf "USERDICT_LIMIT selector=-1 result=%d\n", $result
  set $result = ((int (*)(int))0x1002a5b0)(0)
  printf "USERDICT_LIMIT selector=0 result=%d\n", $result
  set $result = ((int (*)(int))0x1002a5b0)(1)
  printf "USERDICT_LIMIT selector=1 result=%d\n", $result
  set $result = ((int (*)(int))0x1002a5b0)(2)
  printf "USERDICT_LIMIT selector=2 result=%d\n", $result
  set $result = ((int (*)(int))0x1002a5b0)(3)
  printf "USERDICT_LIMIT selector=3 result=%d\n", $result
  set $result = ((int (*)(int))0x1002a5b0)(4)
  printf "USERDICT_LIMIT selector=4 result=%d\n", $result
  set $result = ((int (*)(int))0x1002a5b0)(5)
  printf "USERDICT_LIMIT selector=5 result=%d\n", $result
  set $result = ((int (*)(int))0x1002a5b0)(6)
  printf "USERDICT_LIMIT selector=6 result=%d\n", $result
  set $result = ((int (*)(int))0x1002a5b0)(2147483647)
  printf "USERDICT_LIMIT selector=INT_MAX result=%d\n", $result

  kill
  quit
end

continue
