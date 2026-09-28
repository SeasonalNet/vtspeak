set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
break *0x1001da50
commands
  silent
  disable 1
  set $exact = ((char *(*)(unsigned int))0x1001d9c0)(8)
  set {char}($exact + 0) = 72
  set {char}($exact + 1) = 72
  set {char}($exact + 2) = 32
  set {char}($exact + 3) = 91
  set {char}($exact + 4) = 67
  set {char}($exact + 5) = 73
  set {char}($exact + 6) = 93
  set {char}($exact + 7) = 0
  set $lowercase = ((char *(*)(unsigned int))0x1001d9c0)(8)
  set {char}($lowercase + 0) = 72
  set {char}($lowercase + 1) = 72
  set {char}($lowercase + 2) = 32
  set {char}($lowercase + 3) = 91
  set {char}($lowercase + 4) = 99
  set {char}($lowercase + 5) = 105
  set {char}($lowercase + 6) = 93
  set {char}($lowercase + 7) = 0
  set $mixedcase = ((char *(*)(unsigned int))0x1001d9c0)(8)
  set {char}($mixedcase + 0) = 72
  set {char}($mixedcase + 1) = 72
  set {char}($mixedcase + 2) = 32
  set {char}($mixedcase + 3) = 91
  set {char}($mixedcase + 4) = 67
  set {char}($mixedcase + 5) = 105
  set {char}($mixedcase + 6) = 93
  set {char}($mixedcase + 7) = 0
  set $two_spaces = ((char *(*)(unsigned int))0x1001d9c0)(9)
  set {char}($two_spaces + 0) = 72
  set {char}($two_spaces + 1) = 72
  set {char}($two_spaces + 2) = 32
  set {char}($two_spaces + 3) = 32
  set {char}($two_spaces + 4) = 91
  set {char}($two_spaces + 5) = 67
  set {char}($two_spaces + 6) = 73
  set {char}($two_spaces + 7) = 93
  set {char}($two_spaces + 8) = 0
  set $tab_separator = ((char *(*)(unsigned int))0x1001d9c0)(8)
  set {char}($tab_separator + 0) = 72
  set {char}($tab_separator + 1) = 72
  set {char}($tab_separator + 2) = 9
  set {char}($tab_separator + 3) = 91
  set {char}($tab_separator + 4) = 67
  set {char}($tab_separator + 5) = 73
  set {char}($tab_separator + 6) = 93
  set {char}($tab_separator + 7) = 0
  set $trailing_space = ((char *(*)(unsigned int))0x1001d9c0)(9)
  set {char}($trailing_space + 0) = 72
  set {char}($trailing_space + 1) = 72
  set {char}($trailing_space + 2) = 32
  set {char}($trailing_space + 3) = 91
  set {char}($trailing_space + 4) = 67
  set {char}($trailing_space + 5) = 73
  set {char}($trailing_space + 6) = 93
  set {char}($trailing_space + 7) = 32
  set {char}($trailing_space + 8) = 0
  set $trailing_tab = ((char *(*)(unsigned int))0x1001d9c0)(9)
  set {char}($trailing_tab + 0) = 72
  set {char}($trailing_tab + 1) = 72
  set {char}($trailing_tab + 2) = 32
  set {char}($trailing_tab + 3) = 91
  set {char}($trailing_tab + 4) = 67
  set {char}($trailing_tab + 5) = 73
  set {char}($trailing_tab + 6) = 93
  set {char}($trailing_tab + 7) = 9
  set {char}($trailing_tab + 8) = 0
  set $following_phone = ((char *(*)(unsigned int))0x1001d9c0)(11)
  set {char}($following_phone + 0) = 72
  set {char}($following_phone + 1) = 72
  set {char}($following_phone + 2) = 32
  set {char}($following_phone + 3) = 91
  set {char}($following_phone + 4) = 67
  set {char}($following_phone + 5) = 73
  set {char}($following_phone + 6) = 93
  set {char}($following_phone + 7) = 32
  set {char}($following_phone + 8) = 72
  set {char}($following_phone + 9) = 72
  set {char}($following_phone + 10) = 0
  set $repeated_marker = ((char *(*)(unsigned int))0x1001d9c0)(12)
  set {char}($repeated_marker + 0) = 72
  set {char}($repeated_marker + 1) = 72
  set {char}($repeated_marker + 2) = 32
  set {char}($repeated_marker + 3) = 91
  set {char}($repeated_marker + 4) = 67
  set {char}($repeated_marker + 5) = 73
  set {char}($repeated_marker + 6) = 93
  set {char}($repeated_marker + 7) = 91
  set {char}($repeated_marker + 8) = 67
  set {char}($repeated_marker + 9) = 73
  set {char}($repeated_marker + 10) = 93
  set {char}($repeated_marker + 11) = 0
  set $adjacent_marker = ((char *(*)(unsigned int))0x1001d9c0)(7)
  set {char}($adjacent_marker + 0) = 72
  set {char}($adjacent_marker + 1) = 72
  set {char}($adjacent_marker + 2) = 91
  set {char}($adjacent_marker + 3) = 67
  set {char}($adjacent_marker + 4) = 73
  set {char}($adjacent_marker + 5) = 93
  set {char}($adjacent_marker + 6) = 0
  set $r_exact = ((short (*)(char *))0x1002a590)($exact)
  printf "TARGETPHON boundary=exact result=%d\n", $r_exact
  set $r_lowercase = ((short (*)(char *))0x1002a590)($lowercase)
  printf "TARGETPHON boundary=lowercase result=%d\n", $r_lowercase
  set $r_mixedcase = ((short (*)(char *))0x1002a590)($mixedcase)
  printf "TARGETPHON boundary=mixedcase result=%d\n", $r_mixedcase
  set $r_two_spaces = ((short (*)(char *))0x1002a590)($two_spaces)
  printf "TARGETPHON boundary=two-spaces result=%d\n", $r_two_spaces
  set $r_tab_separator = ((short (*)(char *))0x1002a590)($tab_separator)
  printf "TARGETPHON boundary=tab-separator result=%d\n", $r_tab_separator
  set $r_trailing_space = ((short (*)(char *))0x1002a590)($trailing_space)
  printf "TARGETPHON boundary=trailing-space result=%d\n", $r_trailing_space
  set $r_trailing_tab = ((short (*)(char *))0x1002a590)($trailing_tab)
  printf "TARGETPHON boundary=trailing-tab result=%d\n", $r_trailing_tab
  set $r_following_phone = ((short (*)(char *))0x1002a590)($following_phone)
  printf "TARGETPHON boundary=following-phone result=%d\n", $r_following_phone
  set $r_repeated_marker = ((short (*)(char *))0x1002a590)($repeated_marker)
  printf "TARGETPHON boundary=repeated-marker result=%d\n", $r_repeated_marker
  set $r_adjacent_marker = ((short (*)(char *))0x1002a590)($adjacent_marker)
  printf "TARGETPHON boundary=adjacent-marker result=%d\n", $r_adjacent_marker
  call ((void (*)(void *))0x1001da30)($exact)
  call ((void (*)(void *))0x1001da30)($lowercase)
  call ((void (*)(void *))0x1001da30)($mixedcase)
  call ((void (*)(void *))0x1001da30)($two_spaces)
  call ((void (*)(void *))0x1001da30)($tab_separator)
  call ((void (*)(void *))0x1001da30)($trailing_space)
  call ((void (*)(void *))0x1001da30)($trailing_tab)
  call ((void (*)(void *))0x1001da30)($following_phone)
  call ((void (*)(void *))0x1001da30)($repeated_marker)
  call ((void (*)(void *))0x1001da30)($adjacent_marker)
  continue
end
continue
