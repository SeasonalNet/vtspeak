set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
break *0x1001da50
commands
  silent
  disable 1
  set $phon = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set $first = 65
  while $first <= 90
    set $second = 65
    while $second <= 90
      set {unsigned char}($phon) = $first
      set {unsigned char}($phon + 1) = $second
      set {char}($phon + 2) = 0
      set $result = ((short (*)(char *))0x1002a590)($phon)
      printf "TARGETPHON pair=%c%c result=%d\n", $first, $second, $result
      set $second = $second + 1
    end
    set $first = $first + 1
  end
  call ((void (*)(void *))0x1001da30)($phon)
  continue
end
continue
