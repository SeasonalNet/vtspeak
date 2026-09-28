set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $i = 0
  while $i < 8
    if $i == 0
      set $flag = -2147483648
    end
    if $i == 1
      set $flag = -1
    end
    if $i == 2
      set $flag = 0
    end
    if $i == 3
      set $flag = 1
    end
    if $i == 4
      set $flag = 2
    end
    if $i == 5
      set $flag = 3
    end
    if $i == 6
      set $flag = 255
    end
    if $i == 7
      set $flag = 2147483647
    end
    set $parser = ((void *(*)(void))0x10016bc0)()
    set $line = (char *)malloc(64)
    set {char[64]}$line = "a,\"b,c\",d"
    set $parsed = ((int (*)(void *, char *, int))0x10016bd0)($parser, $line, $flag)
    set $count = ((int (*)(void *))0x10016c10)($parser)
    printf "CSV_FLAG flag=%d result=%d count=%d source=<%s>", $flag, $parsed, $count, $line
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
    set $i = $i + 1
  end
  continue
end

continue
