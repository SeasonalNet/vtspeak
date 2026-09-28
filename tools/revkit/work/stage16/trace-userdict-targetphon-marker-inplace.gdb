set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
break *0x1001da50
commands
  silent
  disable 1
  set $padded = ((char *(*)(unsigned int))0x1001d9c0)(7)
  set {char}($padded + 0) = 9
  set {char}($padded + 1) = 32
  set {char}($padded + 2) = 72
  set {char}($padded + 3) = 72
  set {char}($padded + 4) = 32
  set {char}($padded + 5) = 13
  set {char}($padded + 6) = 0
  set $marker_space = ((char *(*)(unsigned int))0x1001d9c0)(9)
  set {char}($marker_space + 0) = 72
  set {char}($marker_space + 1) = 72
  set {char}($marker_space + 2) = 32
  set {char}($marker_space + 3) = 91
  set {char}($marker_space + 4) = 67
  set {char}($marker_space + 5) = 73
  set {char}($marker_space + 6) = 93
  set {char}($marker_space + 7) = 32
  set {char}($marker_space + 8) = 0
  set $marker_adjacent = ((char *(*)(unsigned int))0x1001d9c0)(7)
  set {char}($marker_adjacent + 0) = 72
  set {char}($marker_adjacent + 1) = 72
  set {char}($marker_adjacent + 2) = 91
  set {char}($marker_adjacent + 3) = 67
  set {char}($marker_adjacent + 4) = 73
  set {char}($marker_adjacent + 5) = 93
  set {char}($marker_adjacent + 6) = 0
  set $marker_extra = ((char *(*)(unsigned int))0x1001d9c0)(11)
  set {char}($marker_extra + 0) = 72
  set {char}($marker_extra + 1) = 72
  set {char}($marker_extra + 2) = 32
  set {char}($marker_extra + 3) = 91
  set {char}($marker_extra + 4) = 67
  set {char}($marker_extra + 5) = 73
  set {char}($marker_extra + 6) = 93
  set {char}($marker_extra + 7) = 32
  set {char}($marker_extra + 8) = 72
  set {char}($marker_extra + 9) = 72
  set {char}($marker_extra + 10) = 0
  set $invalid_marker = ((char *(*)(unsigned int))0x1001d9c0)(10)
  set {char}($invalid_marker + 0) = 72
  set {char}($invalid_marker + 1) = 72
  set {char}($invalid_marker + 2) = 32
  set {char}($invalid_marker + 3) = 91
  set {char}($invalid_marker + 4) = 83
  set {char}($invalid_marker + 5) = 75
  set {char}($invalid_marker + 6) = 73
  set {char}($invalid_marker + 7) = 80
  set {char}($invalid_marker + 8) = 93
  set {char}($invalid_marker + 9) = 0
  set $r = ((short (*)(char *))0x1002a590)($padded)
  printf "TARGETPHON inplace=padded result=%d text=<%s>\n", $r, $padded
  set $r = ((short (*)(char *))0x1002a590)($marker_space)
  printf "TARGETPHON inplace=marker-space result=%d text=<%s>\n", $r, $marker_space
  set $r = ((short (*)(char *))0x1002a590)($marker_adjacent)
  printf "TARGETPHON inplace=marker-adjacent result=%d text=<%s>\n", $r, $marker_adjacent
  set $r = ((short (*)(char *))0x1002a590)($marker_extra)
  printf "TARGETPHON inplace=marker-extra result=%d text=<%s>\n", $r, $marker_extra
  set $r = ((short (*)(char *))0x1002a590)($invalid_marker)
  printf "TARGETPHON inplace=invalid-marker result=%d text=<%s>\n", $r, $invalid_marker
  call ((void (*)(void *))0x1001da30)($padded)
  call ((void (*)(void *))0x1001da30)($marker_space)
  call ((void (*)(void *))0x1001da30)($marker_adjacent)
  call ((void (*)(void *))0x1001da30)($marker_extra)
  call ((void (*)(void *))0x1001da30)($invalid_marker)
  continue
end
continue
