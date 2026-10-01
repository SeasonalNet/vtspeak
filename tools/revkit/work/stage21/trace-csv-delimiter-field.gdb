set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands 1
  silent
  disable 1

  set $parser = ((void *(*)(void))0x10016bc0)()
  set $line = (char *)malloc(64)
  set {char[64]}$line = "a;b;c"
  printf "CSV_DELIMITER default_pointer=%#x default_set=<%s>\n", *(unsigned int *)($parser + 0x14), *(char **)((char *)$parser + 0x14)
  set $parsed = ((int (*)(void *, char *, int))0x10016bd0)($parser, $line, 0)
  set $count = ((int (*)(void *))0x10016c10)($parser)
  printf "CSV_DELIMITER case=default parsed=%d count=%d field0=<%s>\n", $parsed, $count, ((char *(*)(void *, int))0x10016c30)($parser, 0)
  call ((void (*)(void *))0x10016bf0)($parser)

  set $parser = ((void *(*)(void))0x10016bc0)()
  set $delimiter = (char *)malloc(2)
  set {char[2]}$delimiter = ";"
  set {unsigned int}((char *)$parser + 0x14) = (unsigned int)$delimiter
  set $line = (char *)malloc(64)
  set {char[64]}$line = "a;b;c"
  printf "CSV_DELIMITER override_set=<%s>\n", *(char **)((char *)$parser + 0x14)
  set $parsed = ((int (*)(void *, char *, int))0x10016bd0)($parser, $line, 0)
  set $count = ((int (*)(void *))0x10016c10)($parser)
  printf "CSV_DELIMITER case=semicolon parsed=%d count=%d", $parsed, $count
  set $i = 0
  while $i < $count
    printf " field%d=<%s>", $i, ((char *(*)(void *, int))0x10016c30)($parser, $i)
    set $i = $i + 1
  end
  printf "\n"
  call ((void (*)(void *))0x10016bf0)($parser)

  set $parser = ((void *(*)(void))0x10016bc0)()
  set $delimiter = (char *)malloc(2)
  set {char[2]}$delimiter = ";"
  set {unsigned int}((char *)$parser + 0x14) = (unsigned int)$delimiter
  set $line = (char *)malloc(64)
  set {char[64]}$line = "a;\"b;c\";d"
  set $parsed = ((int (*)(void *, char *, int))0x10016bd0)($parser, $line, 0)
  set $count = ((int (*)(void *))0x10016c10)($parser)
  printf "CSV_DELIMITER case=semicolon_quoted parsed=%d count=%d", $parsed, $count
  set $i = 0
  while $i < $count
    printf " field%d=<%s>", $i, ((char *(*)(void *, int))0x10016c30)($parser, $i)
    set $i = $i + 1
  end
  printf "\n"
  call ((void (*)(void *))0x10016bf0)($parser)
  continue
end

continue
