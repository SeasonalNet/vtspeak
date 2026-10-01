set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $result = ((char *(*)(int))0x1001c690)(-2)
  printf "PATHKEY_OOB selector=-2 pointer=%p value=%s\n", $result, $result
  kill
  quit
end

continue
