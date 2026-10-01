set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1

  set $null_count = ((int (*)(void *))0x10016c10)((void *)0)
  printf "CSV_GET_EDGE case=null_parser_count raw=%d\n", $null_count
  set $null_field = ((char *(*)(void *, int))0x10016c30)((void *)0, 0)
  printf "CSV_GET_EDGE case=null_parser_field index=0 value=%p\n", $null_field

  set $p = ((void *(*)(void))0x10016bc0)()
  set $empty_count = ((int (*)(void *))0x10016c10)($p)
  set $empty_field = ((char *(*)(void *, int))0x10016c30)($p, 0)
  printf "CSV_GET_EDGE case=unparsed count=%d index0=%p\n", $empty_count, $empty_field

  set $line = (char *)malloc(32)
  set {char[32]}$line = "first,second"
  set $parse = ((int (*)(void *, char *, int))0x10016bd0)($p, $line, 0)
  set $count = ((int (*)(void *))0x10016c10)($p)
  printf "CSV_GET_EDGE case=parsed parse=%d count=%d\n", $parse, $count

  set $i = 0
  while $i < 6
    if $i == 0
      set $index = -1
    end
    if $i == 1
      set $index = 0
    end
    if $i == 2
      set $index = 1
    end
    if $i == 3
      set $index = 2
    end
    if $i == 4
      set $index = 3
    end
    if $i == 5
      set $index = 2147483647
    end
    set $field = ((char *(*)(void *, int))0x10016c30)($p, $index)
    printf "CSV_GET_EDGE case=parsed_field index=%d value=%p\n", $index, $field
    set $i = $i + 1
  end

  call ((void (*)(void *))0x10016bf0)($p)
  printf "CSV_GET_EDGE case=exit completed=1\n"
  kill
  quit
end

continue
