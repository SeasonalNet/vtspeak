set pagination off
set confirm off
break *0x1003a7b0
commands
  silent
  set $out = *(unsigned int *)($esp + 4)
  set $key = *(unsigned int *)($esp + 8)
  set $type = *(unsigned int *)($esp + 12)
  set $ret = *(unsigned int *)$esp
  printf "TPP_LOOKUP selector=%#x key=", $type
  x/s $key
  tbreak *$ret
  commands
    silent
    printf "TPP_RESULT eax=%#x output=", $eax
    x/s $out
    continue
  end
  continue
end
continue
