set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1

  set $p1 = ((void *(*)(void))0x10016bc0)()
  set $p2 = ((void *(*)(void))0x10016bc0)()
  printf "CSV_LIFECYCLE case=init p1=%p p2=%p distinct=%d size=24 p1_fields=%u,%u,%u,%u,%u delim_ptr=%p delim=<%s> p2_delim_ptr=%p p2_delim=<%s>\n", $p1, $p2, $p1 != $p2, *(unsigned int *)$p1, *(unsigned int *)($p1 + 4), *(unsigned int *)($p1 + 8), *(unsigned int *)($p1 + 12), *(unsigned int *)($p1 + 16), *(char **)($p1 + 20), *(char **)($p1 + 20), *(char **)($p2 + 20), *(char **)($p2 + 20)

  call ((void (*)(void *))0x10016bf0)(0)
  printf "CSV_LIFECYCLE case=exit_null completed=1\n"
  call ((void (*)(void *))0x10016bf0)($p1)
  printf "CSV_LIFECYCLE case=exit_unparsed completed=1\n"
  call ((void (*)(void *))0x10016bf0)($p2)

  set $p3 = ((void *(*)(void))0x10016bc0)()
  set $line = (char *)malloc(32)
  set {char[32]}$line = "first,second"
  set $parse = ((int (*)(void *, char *, int))0x10016bd0)($p3, $line, 0)
  set $count = ((int (*)(void *))0x10016c10)($p3)
  printf "CSV_LIFECYCLE case=parse_copy low_ax=%d count=%d copy_ptr=%p fields_ptr=%p field0=<%s> field1=<%s>\n", $parse, $count, *(char **)($p3 + 4), *(char ***)($p3 + 8), ((char *(*)(void *, int))0x10016c30)($p3, 0), ((char *(*)(void *, int))0x10016c30)($p3, 1)
  call ((void (*)(void *))0x10016bf0)($p3)
  printf "CSV_LIFECYCLE case=exit_parsed completed=1\n"

  kill
  quit
end

continue
