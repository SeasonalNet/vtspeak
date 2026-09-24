set pagination off
set confirm off
break *0x10043810
commands
  silent
  set $return_address = *(unsigned int *)($ebp+4)
  printf "ADDRESS_PARSE_ENTRY text="
  x/s *(char **)($ebp+0x10)
  tbreak *$return_address
  commands
    silent
    printf "ADDRESS_PARSE_RETURN eax=%d\n", $eax
    continue
  end
  continue
end
continue
