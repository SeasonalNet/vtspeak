set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
break *0x1001da50
commands
  silent
  disable 1
  set $n130 = ((char *(*)(unsigned int))0x1001d9c0)(300)
  set $n130_space = ((char *(*)(unsigned int))0x1001d9c0)(300)
  set $n131 = ((char *(*)(unsigned int))0x1001d9c0)(300)
  set $token = 0
  set $position = 0
  while $token < 130
    set {char}($n130 + $position) = 66
    set $position = $position + 1
    set $token = $token + 1
    if $token < 130
      set {char}($n130 + $position) = 32
      set $position = $position + 1
    end
  end
  set {char}($n130 + $position) = 0
  printf "TARGETPHON fresh-case=n130 bytes=%d result=%d\n", $position, ((short (*)(char *))0x1002a590)($n130)

  set $position = 0
  while $token > 0
    set {char}($n130_space + $position) = 66
    set $position = $position + 1
    set $token = $token - 1
    if $token > 0
      set {char}($n130_space + $position) = 32
      set $position = $position + 1
    end
  end
  set {char}($n130_space + $position) = 32
  set $position = $position + 1
  set {char}($n130_space + $position) = 0
  printf "TARGETPHON fresh-case=n130-trailing-space bytes=%d result=%d\n", $position, ((short (*)(char *))0x1002a590)($n130_space)

  set $token = 0
  set $position = 0
  while $token < 131
    set {char}($n131 + $position) = 66
    set $position = $position + 1
    set $token = $token + 1
    if $token < 131
      set {char}($n131 + $position) = 32
      set $position = $position + 1
    end
  end
  set {char}($n131 + $position) = 0
  printf "TARGETPHON fresh-case=n131 bytes=%d result=%d\n", $position, ((short (*)(char *))0x1002a590)($n131)
  call ((void (*)(void *))0x1001da30)($n130)
  call ((void (*)(void *))0x1001da30)($n130_space)
  call ((void (*)(void *))0x1001da30)($n131)
  continue
end
continue
