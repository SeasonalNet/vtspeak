set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
break *0x1001da50
commands
  silent
  disable 1
  set $phon = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set $byte = 9
  while $byte <= 13
    set {unsigned char}($phon) = $byte
    set {char}($phon + 1) = 0
    set $result = ((short (*)(char *))0x1002a590)($phon)
    printf "TARGETPHON single_byte=0x%02x result=%d\n", $byte, $result
    set $byte = $byte + 1
  end
  set $byte = 32
  while $byte <= 126
    set {unsigned char}($phon) = $byte
    set {char}($phon + 1) = 0
    set $result = ((short (*)(char *))0x1002a590)($phon)
    printf "TARGETPHON single_byte=0x%02x result=%d\n", $byte, $result
    set $byte = $byte + 1
  end
  call ((void (*)(void *))0x1001da30)($phon)
  continue
end
continue
