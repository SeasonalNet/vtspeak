set pagination off
set confirm off
break *0x10043b61
commands
  silent
  printf "ADDRESS_GUARD token=%d pending=%d kind=%d key=", *(int *)($ebp-4), *(int *)($ebp-0x24), *(int *)($ebp-0x64)
  x/s ($ebp-0x60)
  continue
end
continue
