set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
break *0x1001da50
commands
  silent
  disable 1
  set $sequence = ((char *(*)(unsigned int))0x1001d9c0)(300)
  set $token = 0
  set $position = 0
  while $token < 131
    set {char}($sequence + $position) = 66
    set $position = $position + 1
    set $token = $token + 1
    if $token < 131
      set {char}($sequence + $position) = 32
      set $position = $position + 1
    end
  end
  set {char}($sequence + $position) = 0
  set $bad = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($bad) = 91
  set {char}($bad + 1) = 0
  set {char}($bad + 2) = 0
  set $invalid = ((char *(*)(unsigned int))0x1001d9c0)(12)
  set {char}($invalid) = 72
  set {char}($invalid + 1) = 72
  set {char}($invalid + 2) = 32
  set {char}($invalid + 3) = 69
  set {char}($invalid + 4) = 88
  set {char}($invalid + 5) = 65
  set {char}($invalid + 6) = 77
  set {char}($invalid + 7) = 80
  set {char}($invalid + 8) = 76
  set {char}($invalid + 9) = 69
  set {char}($invalid + 10) = 0
  set $large_first = ((short (*)(char *))0x1002a590)($sequence)
  printf "TARGETPHON order=large-first length=%d result=%d\n", $position, $large_first
  set $bracket_result = ((short (*)(char *))0x1002a590)($bad)
  printf "TARGETPHON order=bracket result=%d\n", $bracket_result
  set $large_after_bracket = ((short (*)(char *))0x1002a590)($sequence)
  printf "TARGETPHON order=large-after-bracket result=%d\n", $large_after_bracket
  set $invalid_result = ((short (*)(char *))0x1002a590)($invalid)
  printf "TARGETPHON order=invalid-token result=%d\n", $invalid_result
  set $large_after_invalid = ((short (*)(char *))0x1002a590)($sequence)
  printf "TARGETPHON order=large-after-invalid result=%d\n", $large_after_invalid
  call ((void (*)(void *))0x1001da30)($sequence)
  call ((void (*)(void *))0x1001da30)($bad)
  call ((void (*)(void *))0x1001da30)($invalid)
  continue
end
continue
