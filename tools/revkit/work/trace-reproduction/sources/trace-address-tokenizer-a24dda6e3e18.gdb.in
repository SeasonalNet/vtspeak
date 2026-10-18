set pagination off
set confirm off
break *0x100438bf
commands
  silent
  printf "ADDRESS_TOKEN index=%d parser_bytes=%d kind=%d subkind=%d local24=%d local8=%d local28=%d previous=", *(int *)($ebp-4), *(int *)($ebp-0x1c), *(int *)($ebp-0x64), *(int *)($ebp-0x70), *(int *)($ebp-0x20), *(int *)($ebp-8), *(int *)($ebp-0x28)
  x/8bx ($ebp-0x60)
  printf "ADDRESS_INPUT=\n"
  x/s *(unsigned int *)($ebp+0x10)
  continue
end
continue
