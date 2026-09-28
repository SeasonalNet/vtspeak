set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
break *0x1001da50
commands
  silent
  disable 1
  set $ci = ((char *(*)(unsigned int))0x1001d9c0)(5)
  set {char}($ci + 0) = 91
  set {char}($ci + 1) = 67
  set {char}($ci + 2) = 73
  set {char}($ci + 3) = 93
  set {char}($ci + 4) = 0
  set $ci_after = ((char *(*)(unsigned int))0x1001d9c0)(8)
  set {char}($ci_after + 0) = 72
  set {char}($ci_after + 1) = 72
  set {char}($ci_after + 2) = 32
  set {char}($ci_after + 3) = 91
  set {char}($ci_after + 4) = 67
  set {char}($ci_after + 5) = 73
  set {char}($ci_after + 6) = 93
  set {char}($ci_after + 7) = 0
  set $ci_embedded = ((char *(*)(unsigned int))0x1001d9c0)(8)
  set {char}($ci_embedded + 0) = 72
  set {char}($ci_embedded + 1) = 72
  set {char}($ci_embedded + 2) = 91
  set {char}($ci_embedded + 3) = 67
  set {char}($ci_embedded + 4) = 73
  set {char}($ci_embedded + 5) = 93
  set {char}($ci_embedded + 6) = 66
  set {char}($ci_embedded + 7) = 0
  set $skip = ((char *(*)(unsigned int))0x1001d9c0)(7)
  set {char}($skip + 0) = 91
  set {char}($skip + 1) = 83
  set {char}($skip + 2) = 75
  set {char}($skip + 3) = 73
  set {char}($skip + 4) = 80
  set {char}($skip + 5) = 93
  set {char}($skip + 6) = 0
  set $skip_after = ((char *(*)(unsigned int))0x1001d9c0)(10)
  set {char}($skip_after + 0) = 72
  set {char}($skip_after + 1) = 72
  set {char}($skip_after + 2) = 32
  set {char}($skip_after + 3) = 91
  set {char}($skip_after + 4) = 83
  set {char}($skip_after + 5) = 75
  set {char}($skip_after + 6) = 73
  set {char}($skip_after + 7) = 80
  set {char}($skip_after + 8) = 93
  set {char}($skip_after + 9) = 0
  set $bracket = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($bracket + 0) = 91
  set {char}($bracket + 1) = 0
  set $bracket_word = ((char *(*)(unsigned int))0x1001d9c0)(8)
  set {char}($bracket_word + 0) = 91
  set {char}($bracket_word + 1) = 79
  set {char}($bracket_word + 2) = 84
  set {char}($bracket_word + 3) = 72
  set {char}($bracket_word + 4) = 69
  set {char}($bracket_word + 5) = 82
  set {char}($bracket_word + 6) = 93
  set {char}($bracket_word + 7) = 0
  set $hash_sequence = ((char *(*)(unsigned int))0x1001d9c0)(5)
  set {char}($hash_sequence + 0) = 35
  set {char}($hash_sequence + 1) = 32
  set {char}($hash_sequence + 2) = 72
  set {char}($hash_sequence + 3) = 72
  set {char}($hash_sequence + 4) = 0
  set $result_ci = ((short (*)(char *))0x1002a590)($ci)
  printf "TARGETPHON marker=ci-only result=%d\n", $result_ci
  set $result_ci_after = ((short (*)(char *))0x1002a590)($ci_after)
  printf "TARGETPHON marker=ci-after-phone result=%d\n", $result_ci_after
  set $result_ci_embedded = ((short (*)(char *))0x1002a590)($ci_embedded)
  printf "TARGETPHON marker=ci-embedded result=%d\n", $result_ci_embedded
  set $result_skip = ((short (*)(char *))0x1002a590)($skip)
  printf "TARGETPHON marker=skip-only result=%d\n", $result_skip
  set $result_skip_after = ((short (*)(char *))0x1002a590)($skip_after)
  printf "TARGETPHON marker=skip-after-phone result=%d\n", $result_skip_after
  set $result_bracket = ((short (*)(char *))0x1002a590)($bracket)
  printf "TARGETPHON marker=open-bracket result=%d\n", $result_bracket
  set $result_other = ((short (*)(char *))0x1002a590)($bracket_word)
  printf "TARGETPHON marker=other-bracket result=%d\n", $result_other
  set $result_hash = ((short (*)(char *))0x1002a590)($hash_sequence)
  printf "TARGETPHON marker=hash-sequence result=%d\n", $result_hash
  call ((void (*)(void *))0x1001da30)($ci)
  call ((void (*)(void *))0x1001da30)($ci_after)
  call ((void (*)(void *))0x1001da30)($ci_embedded)
  call ((void (*)(void *))0x1001da30)($skip)
  call ((void (*)(void *))0x1001da30)($skip_after)
  call ((void (*)(void *))0x1001da30)($bracket)
  call ((void (*)(void *))0x1001da30)($bracket_word)
  call ((void (*)(void *))0x1001da30)($hash_sequence)
  continue
end
continue
