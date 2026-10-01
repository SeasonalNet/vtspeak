set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $v0 = (unsigned int *)0x1007d6f8
  set $v1 = (unsigned int *)0x1007d6fc
  set $v2 = (unsigned int *)0x1007d700
  set $v3 = (unsigned int *)0x1007d704
  set $orig0 = *$v0
  set $orig1 = *$v1
  set $orig2 = *$v2
  set $orig3 = *$v3
  printf "VERSION_EXPORT_BASE first=%u second=%u third=%u fourth=%u\n", *$v0, *$v1, *$v2, *$v3
  set *$v0 = 203
  set *$v1 = 204
  set *$v2 = 205
  set *$v3 = 206
  printf "VERSION_EXPORT_MUTATED first=%u second=%u third=%u fourth=%u\n", *$v0, *$v1, *$v2, *$v3
  printf "DEFAULT_VERSION_WITH_MUTATED_DATA value=%s\n", ((char *(*)(void))0x10028260)()
  set *$v0 = $orig0
  set *$v1 = $orig1
  set *$v2 = $orig2
  set *$v3 = $orig3
  printf "VERSION_EXPORT_RESTORED first=%u second=%u third=%u fourth=%u\n", *$v0, *$v1, *$v2, *$v3
  continue
end

continue
