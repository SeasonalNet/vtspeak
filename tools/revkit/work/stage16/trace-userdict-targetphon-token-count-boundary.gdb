set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
break *0x1001da50
commands
  silent
  disable 1
  set $sequence = ((char *(*)(unsigned int))0x1001d9c0)(600)
  set $count = 259
  while $count <= 261
    set $position = 0
    set $token = 0
    while $token < $count
      set {char}($sequence + $position) = 66
      set $position = $position + 1
      set $token = $token + 1
      if $token < $count
        set {char}($sequence + $position) = 32
        set $position = $position + 1
      end
    end
    set {char}($sequence + $position) = 0
    set $result = ((short (*)(char *))0x1002a590)($sequence)
    printf "TARGETPHON token_count=%d bytes=%d result=%d\n", $count, $position, $result
    set $count = $count + 1
  end
  call ((void (*)(void *))0x1001da30)($sequence)
  continue
end
continue
