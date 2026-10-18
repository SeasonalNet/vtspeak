set pagination off
set confirm off
break *0x10056150
commands
  silent
  set $key = *(unsigned int *)($esp + 8)
  set $ret = *(unsigned int *)$esp
  printf "LOOKUP_56150 caller=%#x key=", $ret
  x/s $key
  tbreak *$ret
  commands
    silent
    printf "LOOKUP_56150_RETURN value=%d\n", $eax
    continue
  end
  continue
end
break *0x100560a0
commands
  silent
  set $key = *(unsigned int *)($esp + 8)
  set $ret = *(unsigned int *)$esp
  printf "LOOKUP_560a0 caller=%#x key=", $ret
  x/s $key
  tbreak *$ret
  commands
    silent
    printf "LOOKUP_560a0_RETURN value=%d\n", $eax
    continue
  end
  continue
end
continue
