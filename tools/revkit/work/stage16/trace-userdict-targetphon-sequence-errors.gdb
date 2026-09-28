set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
break *0x1001da50
commands
  silent
  disable 1
  set $empty = ((char *(*)(unsigned int))0x1001d9c0)(1)
  set {char}($empty + 0) = 0
  set $spaces = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($spaces + 0) = 32
  set {char}($spaces + 1) = 32
  set {char}($spaces + 2) = 32
  set {char}($spaces + 3) = 0
  set $single = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($single + 0) = 72
  set {char}($single + 1) = 72
  set {char}($single + 2) = 0
  set $sequence = ((char *(*)(unsigned int))0x1001d9c0)(12)
  set {char}($sequence + 0) = 72
  set {char}($sequence + 1) = 72
  set {char}($sequence + 2) = 32
  set {char}($sequence + 3) = 66
  set {char}($sequence + 4) = 32
  set {char}($sequence + 5) = 65
  set {char}($sequence + 6) = 65
  set {char}($sequence + 7) = 49
  set {char}($sequence + 8) = 32
  set {char}($sequence + 9) = 67
  set {char}($sequence + 10) = 72
  set {char}($sequence + 11) = 0
  set $repeated_spaces = ((char *(*)(unsigned int))0x1001d9c0)(11)
  set {char}($repeated_spaces + 0) = 32
  set {char}($repeated_spaces + 1) = 32
  set {char}($repeated_spaces + 2) = 72
  set {char}($repeated_spaces + 3) = 72
  set {char}($repeated_spaces + 4) = 32
  set {char}($repeated_spaces + 5) = 32
  set {char}($repeated_spaces + 6) = 32
  set {char}($repeated_spaces + 7) = 66
  set {char}($repeated_spaces + 8) = 32
  set {char}($repeated_spaces + 9) = 32
  set {char}($repeated_spaces + 10) = 0
  set $tabs = ((char *(*)(unsigned int))0x1001d9c0)(7)
  set {char}($tabs + 0) = 9
  set {char}($tabs + 1) = 72
  set {char}($tabs + 2) = 72
  set {char}($tabs + 3) = 9
  set {char}($tabs + 4) = 66
  set {char}($tabs + 5) = 9
  set {char}($tabs + 6) = 0
  set $bad_token = ((char *(*)(unsigned int))0x1001d9c0)(11)
  set {char}($bad_token + 0) = 72
  set {char}($bad_token + 1) = 72
  set {char}($bad_token + 2) = 32
  set {char}($bad_token + 3) = 69
  set {char}($bad_token + 4) = 88
  set {char}($bad_token + 5) = 65
  set {char}($bad_token + 6) = 77
  set {char}($bad_token + 7) = 80
  set {char}($bad_token + 8) = 76
  set {char}($bad_token + 9) = 69
  set {char}($bad_token + 10) = 0
  set $hash = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($hash + 0) = 35
  set {char}($hash + 1) = 0
  set $bracket = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($bracket + 0) = 91
  set {char}($bracket + 1) = 0
  set $result_null = ((short (*)(char *))0x1002a590)(0)
  printf "TARGETPHON case=null result=%d\n", $result_null
  set $result_empty = ((short (*)(char *))0x1002a590)($empty)
  printf "TARGETPHON case=empty result=%d\n", $result_empty
  set $result_spaces = ((short (*)(char *))0x1002a590)($spaces)
  printf "TARGETPHON case=spaces result=%d\n", $result_spaces
  set $result_single = ((short (*)(char *))0x1002a590)($single)
  printf "TARGETPHON case=single result=%d\n", $result_single
  set $result_sequence = ((short (*)(char *))0x1002a590)($sequence)
  printf "TARGETPHON case=sequence result=%d\n", $result_sequence
  set $result_repeated = ((short (*)(char *))0x1002a590)($repeated_spaces)
  printf "TARGETPHON case=repeated-spaces result=%d\n", $result_repeated
  set $result_tabs = ((short (*)(char *))0x1002a590)($tabs)
  printf "TARGETPHON case=tabs result=%d\n", $result_tabs
  set $result_bad = ((short (*)(char *))0x1002a590)($bad_token)
  printf "TARGETPHON case=bad-token result=%d\n", $result_bad
  set $result_hash = ((short (*)(char *))0x1002a590)($hash)
  printf "TARGETPHON case=hash-marker result=%d\n", $result_hash
  set $result_bracket = ((short (*)(char *))0x1002a590)($bracket)
  printf "TARGETPHON case=open-bracket result=%d\n", $result_bracket

  set $long = ((char *(*)(unsigned int))0x1001d9c0)(263)
  set $i = 0
  while $i < 129
    set {char}($long + $i * 2) = 66
    set {char}($long + $i * 2 + 1) = 32
    set $i = $i + 1
  end
  set {char}($long + 258) = 66
  set {char}($long + 259) = 0
  set $result_259 = ((short (*)(char *))0x1002a590)($long)
  printf "TARGETPHON case=length-259 result=%d\n", $result_259
  set {char}($long + 259) = 32
  set {char}($long + 260) = 0
  set $result_260 = ((short (*)(char *))0x1002a590)($long)
  printf "TARGETPHON case=length-260 result=%d\n", $result_260
  set {char}($long + 260) = 66
  set {char}($long + 261) = 0
  set $result_261 = ((short (*)(char *))0x1002a590)($long)
  printf "TARGETPHON case=length-261 result=%d\n", $result_261
  call ((void (*)(void *))0x1001da30)($long)
  call ((void (*)(void *))0x1001da30)($empty)
  call ((void (*)(void *))0x1001da30)($spaces)
  call ((void (*)(void *))0x1001da30)($single)
  call ((void (*)(void *))0x1001da30)($sequence)
  call ((void (*)(void *))0x1001da30)($repeated_spaces)
  call ((void (*)(void *))0x1001da30)($tabs)
  call ((void (*)(void *))0x1001da30)($bad_token)
  call ((void (*)(void *))0x1001da30)($hash)
  call ((void (*)(void *))0x1001da30)($bracket)
  continue
end
continue
