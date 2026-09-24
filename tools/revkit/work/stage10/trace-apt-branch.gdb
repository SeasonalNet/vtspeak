set pagination off
set confirm off
break *0x10043d52
commands
  silent
  printf "ADDRESS_BRANCH index=%d kind=%d key=", *(int *)($ebp-4), *(int *)($ebp-0x64)
  x/s ($ebp-0x60)
  continue
end
break *0x10044003
commands
  silent
  printf "ADDRESS_JOIN index=%d class=%#x replacement=", *(int *)($ebp-4), *(unsigned char *)($ebp+$ebx-0xa0)
  x/s ($ebp-0x1ac)
  continue
end
continue
