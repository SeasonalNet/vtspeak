set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10027af0
commands
  silent
  printf "LOAD speaker=%d db_path=", *(int *)($esp + 8)
  set $db_path = *(char **)($esp + 12)
  if $db_path
    x/s $db_path
  else
    printf "(null)\n"
  end
  continue
end

continue
quit
