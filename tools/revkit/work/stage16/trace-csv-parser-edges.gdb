set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1

  set $parser = ((void *(*)(void))0x10016bc0)()
  set $line = (char *)malloc(64)
  set {char[64]}$line = "\"a\"\"b\",c"
  set $parsed = ((int (*)(void *, char *, int))0x10016bd0)($parser, $line, 0)
  set $count = ((int (*)(void *))0x10016c10)($parser)
  printf "CSV_EDGE case=doubled_quotes flag=0 low_ax=%d count=%d", $parsed, $count
  if $count > 0
    printf " field0=<%s>", ((char *(*)(void *, int))0x10016c30)($parser, 0)
  end
  if $count > 1
    printf " field1=<%s>", ((char *(*)(void *, int))0x10016c30)($parser, 1)
  end
  printf "\n"
  call ((void (*)(void *))0x10016bf0)($parser)

  set $parser = ((void *(*)(void))0x10016bc0)()
  set $line = (char *)malloc(64)
  set {char[64]}$line = "a\"b,c"
  set $parsed = ((int (*)(void *, char *, int))0x10016bd0)($parser, $line, 0)
  set $count = ((int (*)(void *))0x10016c10)($parser)
  printf "CSV_EDGE case=quote_in_unquoted flag=0 low_ax=%d count=%d", $parsed, $count
  if $count > 0
    printf " field0=<%s>", ((char *(*)(void *, int))0x10016c30)($parser, 0)
  end
  if $count > 1
    printf " field1=<%s>", ((char *(*)(void *, int))0x10016c30)($parser, 1)
  end
  printf "\n"
  call ((void (*)(void *))0x10016bf0)($parser)

  set $parser = ((void *(*)(void))0x10016bc0)()
  set $line = (char *)malloc(64)
  set {char[64]}$line = "a,"
  set $parsed = ((int (*)(void *, char *, int))0x10016bd0)($parser, $line, 0)
  set $count = ((int (*)(void *))0x10016c10)($parser)
  printf "CSV_EDGE case=trailing_delimiter flag=0 low_ax=%d count=%d", $parsed, $count
  if $count > 0
    printf " field0=<%s>", ((char *(*)(void *, int))0x10016c30)($parser, 0)
  end
  if $count > 1
    printf " field1=<%s>", ((char *(*)(void *, int))0x10016c30)($parser, 1)
  end
  printf "\n"
  call ((void (*)(void *))0x10016bf0)($parser)

  set $parser = ((void *(*)(void))0x10016bc0)()
  set $line = (char *)malloc(64)
  set {char[64]}$line = "a,\r\nb"
  set $parsed = ((int (*)(void *, char *, int))0x10016bd0)($parser, $line, 0)
  set $count = ((int (*)(void *))0x10016c10)($parser)
  printf "CSV_EDGE case=crlf flag=0 low_ax=%d count=%d", $parsed, $count
  if $count > 0
    printf " field0=<%s>", ((char *(*)(void *, int))0x10016c30)($parser, 0)
  end
  if $count > 1
    printf " field1=<%s>", ((char *(*)(void *, int))0x10016c30)($parser, 1)
  end
  printf "\n"
  call ((void (*)(void *))0x10016bf0)($parser)

  set $parser = ((void *(*)(void))0x10016bc0)()
  set $line = (char *)malloc(64)
  set {char[64]}$line = "a,b\nc,d"
  set $parsed = ((int (*)(void *, char *, int))0x10016bd0)($parser, $line, 0)
  set $count = ((int (*)(void *))0x10016c10)($parser)
  printf "CSV_EDGE case=two_lines flag=0 low_ax=%d count=%d", $parsed, $count
  if $count > 0
    printf " field0=<%s>", ((char *(*)(void *, int))0x10016c30)($parser, 0)
  end
  if $count > 1
    printf " field1=<%s>", ((char *(*)(void *, int))0x10016c30)($parser, 1)
  end
  if $count > 2
    printf " field2=<%s>", ((char *(*)(void *, int))0x10016c30)($parser, 2)
  end
  if $count > 3
    printf " field3=<%s>", ((char *(*)(void *, int))0x10016c30)($parser, 3)
  end
  printf "\n"
  call ((void (*)(void *))0x10016bf0)($parser)

  set $parser = ((void *(*)(void))0x10016bc0)()
  set $line = (char *)malloc(64)
  set {char[64]}$line = "a,\"b,c\",d"
  set $parsed = ((int (*)(void *, char *, int))0x10016bd0)($parser, $line, 0)
  set $count = ((int (*)(void *))0x10016c10)($parser)
  printf "CSV_EDGE case=quoted_comma flag=0 low_ax=%d count=%d", $parsed, $count
  if $count > 0
    printf " field0=<%s>", ((char *(*)(void *, int))0x10016c30)($parser, 0)
  end
  if $count > 1
    printf " field1=<%s>", ((char *(*)(void *, int))0x10016c30)($parser, 1)
  end
  if $count > 2
    printf " field2=<%s>", ((char *(*)(void *, int))0x10016c30)($parser, 2)
  end
  printf "\n"
  call ((void (*)(void *))0x10016bf0)($parser)

  set $parser = ((void *(*)(void))0x10016bc0)()
  set $line = (char *)malloc(64)
  set {char[64]}$line = "a,\"b,c\",d"
  set $parsed = ((int (*)(void *, char *, int))0x10016bd0)($parser, $line, 1)
  set $count = ((int (*)(void *))0x10016c10)($parser)
  printf "CSV_EDGE case=quoted_comma flag=1 low_ax=%d count=%d", $parsed, $count
  if $count > 0
    printf " field0=<%s>", ((char *(*)(void *, int))0x10016c30)($parser, 0)
  end
  if $count > 1
    printf " field1=<%s>", ((char *(*)(void *, int))0x10016c30)($parser, 1)
  end
  if $count > 2
    printf " field2=<%s>", ((char *(*)(void *, int))0x10016c30)($parser, 2)
  end
  printf "\n"
  call ((void (*)(void *))0x10016bf0)($parser)

  set $parser = ((void *(*)(void))0x10016bc0)()
  set $line = (char *)malloc(64)
  set {char[64]}$line = "a,\"b,c\",d"
  set $parsed = ((int (*)(void *, char *, int))0x10016bd0)($parser, $line, 2)
  set $count = ((int (*)(void *))0x10016c10)($parser)
  printf "CSV_EDGE case=quoted_comma flag=2 low_ax=%d count=%d", $parsed, $count
  if $count > 0
    printf " field0=<%s>", ((char *(*)(void *, int))0x10016c30)($parser, 0)
  end
  if $count > 1
    printf " field1=<%s>", ((char *(*)(void *, int))0x10016c30)($parser, 1)
  end
  if $count > 2
    printf " field2=<%s>", ((char *(*)(void *, int))0x10016c30)($parser, 2)
  end
  printf "\n"
  call ((void (*)(void *))0x10016bf0)($parser)

  continue
end

continue
